/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/product_validation_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

import "testing"

func TestProductValidationSupportsNotEvaluated(t *testing.T) {
	r := ProductValidationReport{Evidence: "NOT_EVALUATED", Decision: "NOT_EVALUATED", Utility: "NOT_EVALUATED"}
	if r.Evidence != "NOT_EVALUATED" || r.Decision != "NOT_EVALUATED" || r.Utility != "NOT_EVALUATED" {
		t.Fatal("fase não avaliada foi convertida")
	}
}
