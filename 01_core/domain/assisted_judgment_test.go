/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/assisted_judgment_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

import "testing"

func TestAssistedOracleEffectiveVerdict(t *testing.T) {
	r := AssistedOracleReport{Records: []AssistedOracleRecord{{CandidateID: "c", Effective: "UNKNOWN"}}}
	if r.EffectiveVerdict("c") != "UNKNOWN" {
		t.Fatal("veredito efetivo não preservado")
	}
}
