/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/improvement_synthesis_test.md
 * @layer L1
 * @updated 2026-09-02
 */

package domain

import "testing"

func TestImprovementSynthesisEnumsRemainSeparate(t *testing.T) {
	if HypothesisSupported != "SUPPORTED" || ImprovementComplete != "COMPLETE" || ImprovementLicenseCompatible != "COMPATIBLE" {
		t.Fatal("public synthesis enums changed")
	}
	if string(HypothesisSupported) == string(ImprovementComplete) {
		t.Fatal("hypothesis and lifecycle result must remain distinct")
	}
}
