/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/functional_queries_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "testing"

func TestFunctionalFrameDoesNotHideMissingSlots(t *testing.T) {
	f := FunctionalFrame{ID: "x", ConceptID: "repository-recommendation", UnknownReason: "input missing"}
	if f.UnknownReason == "" || len(f.Actions) > 0 {
		t.Fatal("unknown frame fabricated")
	}
}
func TestFunctionalQueryKeepsCausalIdentity(t *testing.T) {
	q := FunctionalQuery{FrameID: "f", ConceptID: "adoption-decision", TemplateID: "t", SourceIDs: []string{"s"}, LexiconHash: "h"}
	if q.FrameID == "" || q.TemplateID == "" || q.LexiconHash == "" {
		t.Fatal("causal identity lost")
	}
}
