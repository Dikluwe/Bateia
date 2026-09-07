/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/query_evidence_test.md
 * @layer L2
 * @updated 2026-09-07
 */
package report

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestQueryProjectionBeyondInitialExcerpt0036(t *testing.T) {
	content := strings.Repeat("// license\n", 80) + "void decode() { /* packet decoding */\n finish();\n}\nvoid finish() {}\n"
	file := domain.ExtractedFile{Path: "src/data.cpp", Content: content}
	unit := domain.SemanticUnit{Status: "complete", Path: "/p/src/data.cpp", Inputs: []domain.SemanticInput{{Path: "/p/src/data.cpp", SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(content)))}}, Symbols: []domain.SemanticSymbol{
		{ID: "decode", Name: "decode", Kind: "FunctionDecl", File: "/p/src/data.cpp", Line: 81, EndLine: 83, Definition: true},
		{ID: "finish", Name: "finish", Kind: "FunctionDecl", File: "/p/src/data.cpp", Line: 84, EndLine: 84, Definition: true},
	}, Calls: []domain.SemanticCall{{CallerID: "decode", CalleeID: "finish", Caller: "decode", Callee: "finish", File: "/p/src/data.cpp", Line: 82}}}
	sem := &domain.SemanticReport{Status: "partial", Units: []domain.SemanticUnit{unit, unit}, TotalUnits: 200, SelectedUnits: 2}
	ev := parser.BuildQueryEvidence("packet decoding", sem, []domain.ExtractedFile{file}, "/p")
	md := GenerateMarkdownReport(&domain.SearchReport{Query: domain.QueryFilter{Query: "packet decoding"}, Results: []domain.ExtractedRepository{{Files: []domain.ExtractedFile{file}, Semantic: sem, QueryEvidence: ev, HealthScore: -1}}})
	for _, want := range []string{"src/data.cpp:81–83", "decode → finish", "packet decoding", "Não verificado", "partial", "2/200"} {
		if !strings.Contains(md, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(md, "// license") || strings.Contains(md, "Soluções Exatas Encontradas") {
		t.Fatal("boilerplate or overclaim")
	}
	if strings.Count(md, "definição `decode`") != 1 {
		t.Fatal("duplicate definition")
	}
}
