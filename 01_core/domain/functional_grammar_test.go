/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/functional_grammar_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "testing"

func TestFunctionalGrammarTypes(t *testing.T) {
	g := FunctionalGrammar{SchemaVersion: "2", Version: "v2"}
	if g.SchemaVersion != "2" {
		t.Errorf("expected 2, got %s", g.SchemaVersion)
	}
}
