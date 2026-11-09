/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/integrity_0170_metrics_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestIntegrity0170TrueProvenRequiresImmutableSources(t *testing.T) {
	for _, kind := range []string{"absent-sources", "missing-evidence", "commit-mismatch", "mutable", "valid"} {
		for _, split := range []string{"CONTRACT", "VALIDATION", "TEST"} {
			t.Run(kind+"/"+split, func(t *testing.T) {
				o := evidenceCalibrationObservation("one", split, domain.EvidenceProven, domain.EvidenceProven)
				switch kind {
				case "absent-sources":
					o.Observed.Sources = nil
				case "missing-evidence":
					o.Evidence = nil
				case "commit-mismatch":
					o.Evidence[0].Ref.Commit = "different"
				case "mutable":
					o.Observed.Sources[0].Ref.Commit = "HEAD"
				}
				before, _ := json.Marshal(o)
				got := EvaluateEvidenceCalibration([]domain.EvidenceCalibrationObservation{o})
				after, _ := json.Marshal(o)
				if string(before) != string(after) {
					t.Error("observation mutated")
				}
				for name, m := range map[string]domain.EvidenceCalibrationMetrics{"aggregate": got.Aggregate, "split": got.BySplit[split]} {
					want := 0
					invalid := 1
					if kind == "valid" {
						want = 1
						invalid = 0
					}
					if m.ExpectedProven != 1 || m.ObservedProven != 1 || m.TrueProven != want || m.InvalidSources != invalid || m.ProvenPrecision == nil || m.ProvenRecall == nil {
						t.Fatalf("%s wrong denominator/true positive: %+v", name, m)
					}
					if *m.ProvenPrecision != float64(want) || *m.ProvenRecall != float64(want) {
						t.Errorf("%s invalid sources inflate precision/recall: %+v", name, m)
					}
				}
			})
		}
	}
	t.Run("mixed", func(t *testing.T) {
		good := evidenceCalibrationObservation("good", "TEST", domain.EvidenceProven, domain.EvidenceProven)
		bad := evidenceCalibrationObservation("bad", "TEST", domain.EvidenceProven, domain.EvidenceProven)
		bad.Evidence[0].Ref.Commit = "wrong"
		got := EvaluateEvidenceCalibration([]domain.EvidenceCalibrationObservation{bad, good})
		if got.Aggregate.TrueProven != 1 || got.Aggregate.InvalidSources != 1 || *got.Aggregate.ProvenPrecision != .5 || *got.Aggregate.ProvenRecall != .5 {
			t.Fatal("mixed invalid source inflated rate", got)
		}
	})
}

func TestIntegrity0170AlternativeMetadataPermutationTotalOrder(t *testing.T) {
	for _, field := range []string{"stale", "archived", "default-branch", "updated-at", "search-attempts"} {
		t.Run(field, func(t *testing.T) {
			in := newAlternativeTestInput(t)
			a := alternativeTestCandidate("org", "duplicate", "abc123", 10)
			b := a
			switch field {
			case "stale":
				b.Stale = true
			case "archived":
				b.Archived = true
			case "default-branch":
				a.DefaultBranch = "a"
				b.DefaultBranch = "b"
			case "updated-at":
				a.UpdatedAt = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				b.UpdatedAt = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			case "search-attempts":
				a.SearchAttempts = []domain.SearchAttempt{{Provider: "fixture", Query: "first"}}
				b.SearchAttempts = []domain.SearchAttempt{{Provider: "fixture", Query: "second"}}
			}
			in.assessments = []domain.CandidateAssessment{alternativeTestAssessment("org", "duplicate", "abc123")}
			var control domain.SoftwareAlternativeReport
			for index, candidates := range [][]domain.DiscoveredCandidate{{{Repository: a}, {Repository: b}}, {{Repository: b}, {Repository: a}}, {{Repository: b}, {Repository: a}, {Repository: a}}, {{Repository: a}, {Repository: b}, {Repository: a}}} {
				in.discovery.Candidates = candidates
				before, _ := json.Marshal(in)
				got := rankAlternativeTestInput(t, in)
				after, _ := json.Marshal(in)
				if string(before) != string(after) {
					t.Error("rank mutated snapshots")
				}
				if len(got.Excluded) != 1 || got.Excluded[0].Fit != domain.AlternativeInconclusive || got.TopChoice != nil {
					t.Error("duplicate became eligible", got)
				}
				if field == "stale" && len(got.Excluded) == 1 {
					for _, reason := range got.Excluded[0].Reasons {
						if strings.Contains(reason, "identidade local não verificável") {
							t.Error("JSON tie-break selected stale=true over stale=false")
						}
					}
				}
				if index == 0 {
					control = got
				} else if !reflect.DeepEqual(control, got) {
					t.Errorf("snapshot permutation changed representative/reasons for %s: first=%+v got=%+v", field, control, got)
				}
			}
		})
	}
	t.Run("unique-valid-control", func(t *testing.T) {
		in := newAlternativeTestInput(t)
		in.add(alternativeTestCandidate("org", "unique", "abc123", 10), alternativeTestAssessment("org", "unique", "abc123"))
		got := rankAlternativeTestInput(t, in)
		if len(got.Alternatives) != 1 || got.Alternatives[0].Fit != domain.AlternativeRecommended || got.Alternatives[0].Coverage != 1 || got.TopChoice == nil || got.TopChoice.Repo != "unique" {
			t.Fatal("valid unique control changed", got)
		}
	})
}
