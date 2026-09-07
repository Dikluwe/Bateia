/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/mechanism_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestMechanismDoesNotJoinCompilerContexts(t *testing.T) {
	content := "void entry(){ /* packet decoding */\n stage();\n}\nvoid stage(){\n finish();\n}\nvoid finish(){}\n"
	syms := []domain.SemanticSymbol{{ID: "a", Name: "entry", Kind: "FunctionDecl", File: "/p/x.cpp", Line: 1, EndLine: 3, Definition: true}, {ID: "b", Name: "stage", Kind: "FunctionDecl", File: "/p/x.cpp", Line: 4, EndLine: 6, Definition: true}, {ID: "c", Name: "finish", Kind: "FunctionDecl", File: "/p/x.cpp", Line: 7, EndLine: 7, Definition: true}}
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	u := domain.SemanticUnit{Path: "/p/one.cpp", Status: "complete", Arguments: []string{"-DFIRST"}, Symbols: syms, Inputs: []domain.SemanticInput{{Path: "/p/x.cpp", SHA256: hash}}, Calls: []domain.SemanticCall{{CallerID: "a", CalleeID: "b", File: "/p/x.cpp", Line: 2}}}
	v := u
	v.Path = "/p/two.cpp"
	v.Arguments = []string{"-DSECOND"}
	v.Calls = []domain.SemanticCall{{CallerID: "b", CalleeID: "c", File: "/p/x.cpp", Line: 5}}
	e := BuildQueryEvidence("packet decoding", &domain.SemanticReport{Units: []domain.SemanticUnit{u, v}}, []domain.ExtractedFile{{Path: "x.cpp", Content: content}}, "/p")
	if e.Context == "" || e.ContextSHA256 == "" {
		t.Fatal("missing context")
	}
	for _, c := range e.Calls {
		if c.CallerID == "b" {
			t.Fatal("incompatible unit completed chain")
		}
	}
	for _, d := range e.Definitions {
		for _, x := range d.Excerpts {
			if x.Content != strings.Join(strings.Split(content, "\n")[x.FirstLine-1:x.LastLine], "\n") {
				t.Fatal("fabricated snippet")
			}
		}
	}
}
func TestDistantMisleadingCommentDoesNotAnchor(t *testing.T) {
	content := "void unrelated(){ // packet decoding\n" + strings.Repeat("// filler\n", 30) + "other();\n}\n"
	u := domain.SemanticUnit{Status: "complete", Symbols: []domain.SemanticSymbol{{ID: "x", Name: "unrelated", Kind: "FunctionDecl", File: "/p/x.cpp", Line: 1, EndLine: 33, Definition: true}}, Inputs: []domain.SemanticInput{{Path: "/p/x.cpp", SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(content)))}}, Calls: []domain.SemanticCall{{CallerID: "x", CalleeID: "other", File: "/p/x.cpp", Line: 32}}}
	e := BuildQueryEvidence("packet decoding", &domain.SemanticReport{Units: []domain.SemanticUnit{u}}, []domain.ExtractedFile{{Path: "x.cpp", Content: content}}, "/p")
	if len(e.Definitions) != 0 {
		t.Fatal("distant comment promoted unrelated function")
	}
}
