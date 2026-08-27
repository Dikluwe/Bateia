/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/competitive_autophagy_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "testing"

func TestCompetitiveAutophagyJSONEnumsRemainStable(t *testing.T) {
	if InternalKeep != "KEEP" || InternalConsolidate != "CONSOLIDATE" || InternalGap != "GAP" || InternalUnknown != "UNKNOWN" {
		t.Fatal("autophagy dispositions changed")
	}
}

func TestKatamariPolicyDefaultsDoNotPermitUnknown(t *testing.T) {
	if (KatamariPolicy{}).AllowUnknown {
		t.Fatal("zero policy must not absorb unknown evidence")
	}
}
