/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/rule_dsl.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestEvaluateDSL(t *testing.T) {
	rules := []DSLRule{
		{Name: "NO_RAW_SQL", Severity: "HIGH", Pattern: "SELECT * FROM"},
		{Name: "HARDCODED_SECRET", Severity: "CRITICAL", Pattern: "API_SECRET="},
	}

	violations := EvaluateDSL(rules, "const query = 'SELECT * FROM users';")
	if len(violations) != 1 {
		t.Fatalf("Esperava 1 violação, obteve %d", len(violations))
	}

	if violations[0] != "NO_RAW_SQL" {
		t.Errorf("Esperava violação NO_RAW_SQL, obteve %s", violations[0])
	}
}
