/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/multicriteria_decision_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

import "testing"

func TestMulticriteriaContractKeepsPreferenceSeparateFromMeasurement(t *testing.T) {
	measurement := CriterionMeasurement{Candidate: "acme/tool@abc", CriterionID: "coverage", Status: EvidenceProven, Interval: CriterionInterval{Low: .8, High: .8}}
	scenario := PreferenceScenario{ID: "delivery", Weights: []CriterionWeight{{CriterionID: "coverage", Weight: 1}}}
	if measurement.Status != EvidenceProven || scenario.Weights[0].Weight != 1 || PairwiseIncomparable == PairwiseUncertain {
		t.Fatal("contrato misturou medição, preferência ou incerteza")
	}
}
