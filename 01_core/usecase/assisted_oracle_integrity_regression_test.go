/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/assisted_oracle_integrity_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"testing"
)

func assistedIntegritySnapshot(arm, text string) domain.WebSearchSnapshot {
	s := sealedSnapshot(arm, "https://github.com/acme/tool")
	s.Responses[0].Results[0].Title, s.Responses[0].Results[0].Snippet = text, text+" mechanism"
	s.Responses[0] = SealWebResponse(s.Responses[0].Request, s.Responses[0])
	s.SnapshotHash = WebSnapshotHash(s)
	return s
}
func TestAssistedUnionRejectsInvalidSnapshotBeforePublication(t *testing.T) {
	for name, mutate := range map[string]func(*domain.WebSearchSnapshot){
		"missing seal": func(s *domain.WebSearchSnapshot) { s.SnapshotHash = "" },
		"tampered":     func(s *domain.WebSearchSnapshot) { s.Responses[0].Results[0].Title = "tampered" },
		"inner response": func(s *domain.WebSearchSnapshot) {
			s.Responses[0].ResponseHash = "wrong"
			s.SnapshotHash = WebSnapshotHash(*s)
		},
		"missing response identity": func(s *domain.WebSearchSnapshot) {
			s.Responses[0].Provider = ""
			s.Responses[0] = SealWebResponse(s.Responses[0].Request, s.Responses[0])
			s.SnapshotHash = WebSnapshotHash(*s)
		},
		"schema": func(s *domain.WebSearchSnapshot) { s.SchemaVersion = "other"; s.SnapshotHash = WebSnapshotHash(*s) },
	} {
		for _, reverse := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "/first", true: "/last"}[reverse], func(t *testing.T) {
				good, bad := assistedIntegritySnapshot("good", "literal"), assistedIntegritySnapshot("bad", "literal")
				mutate(&bad)
				snapshots := []domain.WebSearchSnapshot{bad, good}
				if reverse {
					snapshots[0], snapshots[1] = good, bad
				}
				before, _ := json.Marshal(snapshots)
				got, err := PrepareBlindCandidateUnion(snapshots)
				if !errors.Is(err, domain.ErrIntegrity) || !reflect.DeepEqual(got, domain.BlindCandidatePackage{}) {
					t.Errorf("invalid snapshot admitted: %v %+v", err, got)
				}
				after, _ := json.Marshal(snapshots)
				if string(before) != string(after) {
					t.Fatal("snapshot mutated")
				}
			})
		}
	}
}
func TestAssistedUnionPreservesAllVariantsInCanonicalOrder(t *testing.T) {
	a, b := assistedIntegritySnapshot("A", "Zulu"), assistedIntegritySnapshot("B", "Alpha")
	var control domain.BlindCandidatePackage
	for i, input := range [][]domain.WebSearchSnapshot{{a, b, a}, {a, a, b}, {b, a, a}} {
		before, _ := json.Marshal(input)
		got, err := PrepareBlindCandidateUnion(input)
		if err != nil {
			t.Fatal(err)
		}
		if len(got.Candidates) != 1 || len(got.Candidates[0].Fragments) != 4 {
			t.Errorf("lost variant: %+v", got)
			continue
		}
		if i == 0 {
			control = got
		} else if !reflect.DeepEqual(got, control) {
			t.Error("snapshot order changed package")
		}
		expected := map[string]bool{"Alpha": true, "Zulu": true, "Alpha mechanism": true, "Zulu mechanism": true}
		ids := map[string]bool{}
		for _, f := range got.Candidates[0].Fragments {
			if !expected[f.Content] || ids[f.FragmentID] {
				t.Fatal("invalid variant or duplicate fragment ID", f)
			}
			delete(expected, f.Content)
			ids[f.FragmentID] = true
			id := got.Candidates[0].CandidateID
			r := ValidateAssistedJudgments(got, []domain.AssistedJudgment{{CandidateID: id, Verdict: "PARTIAL", Confidence: 1, RubricVersion: AssistedRubricVersion, Evidence: []domain.EvidencePointer{{CandidateID: id, FragmentID: f.FragmentID}}}}, []domain.HumanReview{{CandidateID: id, Action: "CONFIRMED"}})
			if r.Verdict != "ORACLE_ACCEPTED" {
				t.Fatal("variant cannot be cited", r)
			}
		}
		after, _ := json.Marshal(input)
		if string(before) != string(after) {
			t.Fatal("input mutated")
		}
	}
	// Reordering responses within a valid newly constructed snapshot is equivalent.
	combined := a
	combined.Responses = []domain.RawWebResponse{a.Responses[0], b.Responses[0]}
	combined.SnapshotHash = WebSnapshotHash(combined)
	got, err := PrepareBlindCandidateUnion([]domain.WebSearchSnapshot{combined})
	if err != nil || !reflect.DeepEqual(got, control) {
		t.Error("response union differs", err)
	}
	combined.Responses[0], combined.Responses[1] = combined.Responses[1], combined.Responses[0]
	combined.SnapshotHash = WebSnapshotHash(combined)
	got, err = PrepareBlindCandidateUnion([]domain.WebSearchSnapshot{combined})
	if err != nil || !reflect.DeepEqual(got, control) {
		t.Error("response order differs", err)
	}
}
func TestAssistedCanonicalControlKeepsPackageHash(t *testing.T) {
	s := assistedIntegritySnapshot("A", "literal")
	one, err := PrepareBlindCandidateUnion([]domain.WebSearchSnapshot{s})
	if err != nil {
		t.Fatal(err)
	}
	if one.PackageHash != "5a3f155b2af8a6020bc0d4bdc7155ca6a5cfb262b3db9af6e82b51a66276a506" {
		t.Fatal("canonical package hash changed", one.PackageHash)
	}
	same, err := PrepareBlindCandidateUnion([]domain.WebSearchSnapshot{s, s})
	if err != nil || !reflect.DeepEqual(same, one) {
		t.Fatal("identical occurrence changed package", err)
	}
	empty, err := PrepareBlindCandidateUnion(nil)
	if err != nil || len(empty.Candidates) != 0 {
		t.Fatal("empty input changed", err)
	}
}
func TestAssistedEvidenceFailurePrecedenceAndGrounding(t *testing.T) {
	for _, validEvidence := range []bool{false, true} {
		for _, action := range []string{"CONFIRMED", "INVALID", "OVERRIDDEN"} {
			pkg := domain.BlindCandidatePackage{SchemaVersion: "1", Candidates: []domain.BlindCandidate{{CandidateID: "a", Fragments: []domain.BlindEvidenceFragment{{FragmentID: "f", Content: "literal"}}}}}
			fragment := "f"
			if !validEvidence {
				fragment = "missing"
			}
			judgments := []domain.AssistedJudgment{{CandidateID: "a", Verdict: "PARTIAL", Confidence: 1, RubricVersion: AssistedRubricVersion, Evidence: []domain.EvidencePointer{{CandidateID: "a", FragmentID: fragment}}}}
			reviews := []domain.HumanReview{{CandidateID: "a", Action: action, Verdict: "RELEVANT"}} // missing reason is intentional for OVERRIDDEN
			got := ValidateAssistedJudgments(pkg, judgments, reviews)
			want := "ORACLE_ACCEPTED"
			if action != "CONFIRMED" {
				want = "INCONCLUSIVE_HUMAN_REVIEW"
			}
			if !validEvidence {
				want = "FAIL_EVIDENCE_GROUNDING"
			}
			rate := 0.0
			if validEvidence {
				rate = 1
			}
			if got.Verdict != want || got.Metrics.EvidenceGroundingRate != rate {
				t.Errorf("evidence=%v action=%s got=%s rate=%g want=%s/%g", validEvidence, action, got.Verdict, got.Metrics.EvidenceGroundingRate, want, rate)
			}
			if action != "CONFIRMED" && got.Records[0].Effective != "UNKNOWN" {
				t.Fatal("invalid review became effective")
			}
			if len(got.Records[0].AIHistory) != 1 || len(got.Records[0].HumanHistory) != 1 {
				t.Fatal("history erased")
			}
		}
	}
	pkg := domain.BlindCandidatePackage{SchemaVersion: "1"}
	for _, id := range []string{"a", "b", "c"} {
		pkg.Candidates = append(pkg.Candidates, domain.BlindCandidate{CandidateID: id, Fragments: []domain.BlindEvidenceFragment{{FragmentID: "f", Content: "literal"}}})
	}
	js := []domain.AssistedJudgment{{CandidateID: "a", Verdict: "PARTIAL", Confidence: 1, RubricVersion: AssistedRubricVersion, Evidence: []domain.EvidencePointer{{CandidateID: "a", FragmentID: "f"}}}, {CandidateID: "b", Verdict: "PARTIAL", Confidence: 1, RubricVersion: AssistedRubricVersion, Evidence: []domain.EvidencePointer{{CandidateID: "b", FragmentID: "missing"}}}}
	reviews := []domain.HumanReview{{CandidateID: "a", Action: "CONFIRMED"}, {CandidateID: "b", Action: "INVALID"}}
	for i := 0; i < 2; i++ {
		got := ValidateAssistedJudgments(pkg, js, reviews)
		if got.Verdict != "FAIL_EVIDENCE_GROUNDING" || got.Metrics.EvidenceGroundingRate != 1.0/3 {
			t.Errorf("mixed grounding/precedence incorrect: %+v", got)
		}
		pkg.Candidates[0], pkg.Candidates[2] = pkg.Candidates[2], pkg.Candidates[0]
	}
	empty := ValidateAssistedJudgments(domain.BlindCandidatePackage{SchemaVersion: "1"}, nil, nil)
	if empty.Metrics.EvidenceGroundingRate != 0 || math.IsNaN(empty.Metrics.EvidenceGroundingRate) {
		t.Fatal("empty metric invalid")
	}
}
