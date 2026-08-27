/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/functional_grammar_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func TestRealizeAndAuditFunctionalGrammar(t *testing.T) {
	g := domain.FunctionalGrammar{
		SchemaVersion: "2",
		Version:       "v2",
		Morphology: []domain.MorphologyRule{
			{ID: "m1", SubjectKind: "singular_3rd", Source: "evaluate", Realized: "evaluates"},
		},
	}
	exp := domain.FunctionalExpansion{
		Queries: []domain.FunctionalQuery{
			{ID: "q1", Text: "tool that evaluate code"},
		},
	}

	realized, err := RealizeFunctionalGrammar(exp, g)
	if err != nil {
		t.Fatalf("realize error: %v", err)
	}
	if realized.Queries[0].Text != "tool that evaluates code" {
		t.Errorf("expected tool that evaluates code, got %s", realized.Queries[0].Text)
	}

	findings, err := AuditFunctionalGrammar(realized, g)
	if err != nil {
		t.Fatalf("audit error: %v", err)
	}
	if len(findings) != 1 || findings[0].Verdict != "PASS" {
		t.Errorf("expected PASS verdict")
	}
}
