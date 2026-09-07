/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query_evidence_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestQueryEvidenceFollowsCallsAndRejectsChangedBytes(t *testing.T) {
	content := "void entry() { /* packet decoding */\n stage();\n}\nvoid stage() {\n finish();\n}\nvoid finish() {}\n"
	files := []domain.ExtractedFile{{Path: "src/main.cpp", Content: content}}
	symbols := []domain.SemanticSymbol{
		{ID: "a", Name: "entry", Kind: "FunctionDecl", File: "/p/src/main.cpp", Line: 1, EndLine: 3, Definition: true},
		{ID: "b", Name: "stage", Kind: "FunctionDecl", File: "/p/src/main.cpp", Line: 4, EndLine: 6, Definition: true},
		{ID: "c", Name: "finish", Kind: "FunctionDecl", File: "/p/src/main.cpp", Line: 7, EndLine: 7, Definition: true},
	}
	u := domain.SemanticUnit{Status: "complete", Symbols: symbols, Inputs: []domain.SemanticInput{{Path: "/p/src/main.cpp", SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(content)))}}, Calls: []domain.SemanticCall{
		{CallerID: "a", CalleeID: "b", Caller: "entry", Callee: "stage", File: "/p/src/main.cpp", Line: 2},
		{CallerID: "b", CalleeID: "c", Caller: "stage", Callee: "finish", File: "/p/src/main.cpp", Line: 5},
	}}
	semantic := &domain.SemanticReport{Status: "partial", Units: []domain.SemanticUnit{u, u}}
	ev := BuildQueryEvidence("packet decoding", semantic, files, "/p")
	if ev.Status != "implementation_evidence" || len(ev.Definitions) != 3 || len(ev.Calls) != 2 {
		t.Fatalf("lost graph/dedup: %+v", ev)
	}
	if !strings.Contains(ev.Definitions[0].Snippet, "packet decoding") || ev.Definitions[0].SnippetFirstLine != 1 {
		t.Fatal("snippet lost anchor")
	}
	if !reflect.DeepEqual(ev, BuildQueryEvidence("packet decoding", semantic, files, "/p")) {
		t.Fatal("nondeterministic")
	}
	if got := BuildQueryEvidence("transparency opengl", semantic, files, "/p"); got.Status != "insufficient" {
		t.Fatal("unrelated candidate promoted")
	}
	files[0].Content += "// changed"
	if got := BuildQueryEvidence("packet decoding", semantic, files, "/p"); len(got.Definitions) != 0 {
		t.Fatal("stale hash accepted")
	}
}
func TestSpecificQueryDoesNotFillWithUnrelatedFiles(t *testing.T) {
	paths := SelectEvidenceFiles("packet decoding", []domain.DirectoryNode{{Path: "README.md", Type: "file"}, {Path: "xml/a.h", Type: "file"}, {Path: "src/packet.cpp", Type: "file"}}, 12)
	if strings.Join(paths, ",") != "README.md,src/packet.cpp" {
		t.Fatalf("irrelevant fill: %v", paths)
	}
}
