/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/flow_test.md
 * @layer L2
 * @updated 2026-09-07
 */
package report

import (
	"bateia/01_core/domain"
	"strings"
	"testing"
)

func TestFlowProjectionStates0038(t *testing.T) {
	next := 2
	repo := domain.ExtractedRepository{QueryEvidence: &domain.QueryEvidence{Context: "x.cpp", Definitions: []domain.ImplementationEvidence{{Symbol: domain.SemanticSymbol{ID: "f", File: "x.cpp", Line: 1}}}}, Semantic: &domain.SemanticReport{Units: []domain.SemanticUnit{{Path: "x.cpp", Status: "complete", Flow: &domain.StaticFlow{Schema: "bateia-flow-v1", Status: "complete", Functions: []domain.FlowFunction{{FlowLocation: domain.FlowLocation{File: "x.cpp", Line: 1}, Name: "render", Status: "complete", Blocks: []domain.FlowBlock{{ID: 1, Successors: []*int{&next, nil}, Terminator: &domain.FlowLocation{Kind: "IfStmt", File: "x.cpp", Line: 2}, Statements: []domain.FlowLocation{{Kind: "CallExpr", File: "x.cpp", Line: 3, Resolution: "unresolved"}}}}}}}}}}}
	var sb strings.Builder
	writeFlowEvidence(&sb, repo, []string{"f"})
	for _, s := range []string{"CFG observado", "B2", "indisponível", "IfStmt", "alvo não resolvido", "não comprovam viabilidade"} {
		if !strings.Contains(sb.String(), s) {
			t.Errorf("missing %s", s)
		}
	}
	repo.Semantic.Units[0].Status = "not_evaluated"
	sb.Reset()
	writeFlowEvidence(&sb, repo, []string{"f"})
	if strings.Contains(sb.String(), "CFG observado") || !strings.Contains(sb.String(), "não avaliado") {
		t.Fatal("invalid unit promoted")
	}
}
