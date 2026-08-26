/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/jsonpath_extractor_test.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestQueryJSONPath(t *testing.T) {
	jsonStr := `{
		"name": "bateia",
		"version": "1.0.0",
		"dependencies": {
			"express": "^4.17.1",
			"react": "^18.0.0"
		}
	}`

	val, ok := QueryJSONPath(jsonStr, "$.dependencies.express")
	if !ok || val != "^4.17.1" {
		t.Errorf("Esperava '^4.17.1', obteve %s (ok=%v)", val, ok)
	}

	valName, okName := QueryJSONPath(jsonStr, "name")
	if !okName || valName != "bateia" {
		t.Errorf("Esperava 'bateia', obteve %s (ok=%v)", valName, okName)
	}

	_, okMissing := QueryJSONPath(jsonStr, "$.dependencies.nonexistent")
	if okMissing {
		t.Errorf("Esperava false para campo inexistente")
	}
}
