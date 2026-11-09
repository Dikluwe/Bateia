/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/multicriteria_extreme_scale_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func TestMulticriteriaExtremeFiniteScalesRemainSerializable(t *testing.T) {
	for _, bound := range []float64{1, 1e308, math.MaxFloat64} {
		for _, direction := range []domain.CriterionDirection{domain.CriterionMaximize, domain.CriterionMinimize} {
			for _, tc := range []struct {
				low, high float64
				missing   bool
			}{{-1, -1, false}, {-.5, -.5, false}, {0, 0, false}, {.5, .5, false}, {1, 1, false}, {-.5, .5, false}, {0, 0, true}} {
				t.Run(fmt.Sprintf("%g/%s/%g/%g/%v", bound, direction, tc.low, tc.high, tc.missing), func(t *testing.T) {
					r := multicriteriaFixture()
					r.Assessment.Assessments = r.Assessment.Assessments[:1]
					a := r.Assessment.Assessments[0]
					r.Criteria = []domain.MulticriteriaCriterion{{ID: "scale", Direction: direction, ScaleMin: -bound, ScaleMax: bound}}
					r.Scenarios = []domain.PreferenceScenario{{ID: "explicit", Weights: []domain.CriterionWeight{{CriterionID: "scale", Weight: 1}}}}
					r.Measurements = nil
					low, high := (tc.low+1)/2, (tc.high+1)/2
					if direction == domain.CriterionMinimize {
						low, high = 1-high, 1-low
					}
					if tc.missing {
						low, high = 0, 1
					} else {
						status := domain.EvidencePartial
						if tc.low == tc.high {
							status = domain.EvidenceProven
						}
						r.Measurements = []domain.CriterionMeasurement{{Candidate: a.Repository.Key(), CriterionID: "scale", Status: status, Interval: domain.CriterionInterval{Low: tc.low * bound, High: tc.high * bound}, Sources: []domain.EvidenceSource{multicriteriaSource(a.Repository, domain.EvidenceFile, "measurement.go", "declared measurement")}}}
					}
					before, e := json.Marshal(r)
					if e != nil {
						t.Fatal(e)
					}
					got, e := AnalyzeMulticriteriaDecision(r)
					if e != nil || len(got.Scenarios) != 1 || len(got.Scenarios[0].Scores) != 1 {
						t.Fatal(e, got)
					}
					score := got.Scenarios[0].Scores[0]
					utility := score.Contributions[0].Utility
					if math.IsNaN(utility.Low) || math.IsNaN(utility.High) || math.Abs(utility.Low-low) > 1e-14 || math.Abs(utility.High-high) > 1e-14 || score.Score != utility || score.Contributions[0].Weighted != utility {
						t.Errorf("invalid normalization: got=%+v want=[%g,%g]", score, low, high)
					}
					if _, e = json.Marshal(got); e != nil {
						t.Errorf("successful report not JSON serializable: %v", e)
					}
					after, _ := json.Marshal(r)
					if string(before) != string(after) {
						t.Fatal("input mutated")
					}
				})
			}
		}
	}
}
func TestMulticriteriaOrdinaryScaleKeepsFrozenReport(t *testing.T) {
	r := multicriteriaFixture()
	r.Scenarios = []domain.PreferenceScenario{{ID: "balanced", Weights: []domain.CriterionWeight{{CriterionID: "coverage", Weight: .5}, {CriterionID: "adaptation_cost", Weight: .5}}}}
	got, e := AnalyzeMulticriteriaDecision(r)
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(got)
	if fmt.Sprintf("%x", sha256.Sum256(raw)) != "8e1172c41dfac65b5ef847c784405ac85049a0e2055eaa1eecd75cf34de31bc1" {
		t.Fatal("ordinary multicriteria report changed")
	}
}
