/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/integrity_0170_judgment_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"sort"
	"testing"
)

func integrity0170Judgment(id string) domain.AssistedJudgment {
	return domain.AssistedJudgment{CandidateID: id, Verdict: "PARTIAL", Confidence: .99, RubricVersion: AssistedRubricVersion, Reasons: []string{"literal evidence"}, Evidence: []domain.EvidencePointer{{CandidateID: id, FragmentID: "f"}}}
}
func integrity0170Package() domain.BlindCandidatePackage {
	return domain.BlindCandidatePackage{SchemaVersion: "1", Candidates: []domain.BlindCandidate{{CandidateID: "c", Fragments: []domain.BlindEvidenceFragment{{FragmentID: "f", Content: "literal evidence"}}}}}
}
func integrity0170Has(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
func integrity0170SortedJSON[T any](xs []T) []string {
	out := []string{}
	for _, x := range xs {
		b, _ := json.Marshal(x)
		out = append(out, string(b))
	}
	sort.Strings(out)
	return out
}
func integrity0170NoFabricatedReview(t *testing.T, r domain.AssistedOracleReport) {
	t.Helper()
	m := r.Metrics
	if m.Escalated != 1 || m.Reviewed != 0 || m.DefinitiveReviewed != 0 || m.DefinitiveAgreements != 0 || m.DefinitiveOverrides != 0 || m.ResolvedUncertain != 0 {
		t.Errorf("ambiguous candidate manufactured resolution metrics: %+v", m)
	}
}

func TestIntegrity0170AmbiguousJudgmentsKeepAllDeclarations(t *testing.T) {
	for _, difference := range []string{"verdict", "confidence", "reason", "evidence-order"} {
		t.Run(difference, func(t *testing.T) {
			pkg := integrity0170Package()
			a := integrity0170Judgment("c")
			b := integrity0170Clone(a)
			switch difference {
			case "verdict":
				b.Verdict = "IRRELEVANT"
			case "confidence":
				b.Confidence = .98
			case "reason":
				b.Reasons = []string{"other explanation"}
			case "evidence-order":
				pkg.Candidates[0].Fragments = append(pkg.Candidates[0].Fragments, domain.BlindEvidenceFragment{FragmentID: "g", Content: "more evidence"})
				a.Evidence = append(a.Evidence, domain.EvidencePointer{CandidateID: "c", FragmentID: "g"})
				b = integrity0170Clone(a)
				b.Evidence[0], b.Evidence[1] = b.Evidence[1], b.Evidence[0]
			}
			reviews := []domain.HumanReview{{CandidateID: "c", Action: "CONFIRMED"}}
			var control domain.AssistedOracleReport
			for index, js := range [][]domain.AssistedJudgment{{a, b, a}, {b, a, a}, {a, a, b}} {
				before, _ := json.Marshal(js)
				got := ValidateAssistedJudgments(pkg, js, reviews)
				after, _ := json.Marshal(js)
				if string(before) != string(after) {
					t.Error("judgments mutated")
				}
				record := got.Record("c")
				if got.Verdict != "INCONCLUSIVE_HUMAN_REVIEW" || record.Effective != "UNKNOWN" || !record.Escalated || !integrity0170Has(record.EscalationCauses, "AMBIGUOUS_AI_JUDGMENTS") {
					t.Errorf("ambiguity chose a declaration: %+v", got)
				}
				if len(record.AIHistory) != 3 || !reflect.DeepEqual(integrity0170SortedJSON(record.AIHistory), integrity0170SortedJSON(js)) {
					t.Error("lost AI declaration", record.AIHistory)
				}
				actual := []string{}
				for _, j := range record.AIHistory {
					raw, _ := json.Marshal(j)
					actual = append(actual, string(raw))
				}
				if !sort.StringsAreSorted(actual) {
					t.Error("history does not use full declaration order")
				}
				if len(record.HumanHistory) != 1 {
					t.Error("review lost")
				}
				integrity0170NoFabricatedReview(t, got)
				if got.Metrics.EvidenceGroundingRate != 1 {
					t.Error("valid distinct declarations lost grounding")
				}
				if index == 0 {
					control = got
				} else if !reflect.DeepEqual(got, control) {
					t.Error("arrival order changed report")
				}
				if !integrity0170Has(DetermineHumanEscalations(pkg, js), "c") {
					t.Error("ambiguous judgments not escalated")
				}
			}
		})
	}
}

func TestIntegrity0170AmbiguousReviewsNeverResolveOrOverride(t *testing.T) {
	for _, aiVerdict := range []string{"PARTIAL", "UNCERTAIN"} {
		for _, difference := range []string{"action", "reason", "verdict"} {
			t.Run(aiVerdict+"/"+difference, func(t *testing.T) {
				pkg := integrity0170Package()
				j := integrity0170Judgment("c")
				j.Verdict = aiVerdict
				a := domain.HumanReview{CandidateID: "c", Action: "RESOLVED", Verdict: "PARTIAL", Reason: "checked first"}
				b := a
				switch difference {
				case "action":
					b.Action = "OVERRIDDEN"
				case "reason":
					b.Reason = "checked second"
				case "verdict":
					b.Verdict = "IRRELEVANT"
				}
				var control domain.AssistedOracleReport
				for index, reviews := range [][]domain.HumanReview{{b, a, a}, {a, b, a}, {a, a, b}} {
					before, _ := json.Marshal(reviews)
					got := ValidateAssistedJudgments(pkg, []domain.AssistedJudgment{j}, reviews)
					after, _ := json.Marshal(reviews)
					if string(before) != string(after) {
						t.Error("reviews mutated")
					}
					record := got.Record("c")
					if got.Verdict != "INCONCLUSIVE_HUMAN_REVIEW" || record.Effective != "UNKNOWN" || !record.Escalated || !integrity0170Has(record.EscalationCauses, "AMBIGUOUS_HUMAN_REVIEWS") {
						t.Errorf("ambiguous reviews selected last: %+v", got)
					}
					if len(record.HumanHistory) != 3 || !reflect.DeepEqual(integrity0170SortedJSON(record.HumanHistory), integrity0170SortedJSON(reviews)) {
						t.Error("review history lost")
					}
					actual := []string{}
					for _, r := range record.HumanHistory {
						raw, _ := json.Marshal(r)
						actual = append(actual, string(raw))
					}
					if !sort.StringsAreSorted(actual) {
						t.Error("review history not ordered by full declaration")
					}
					integrity0170NoFabricatedReview(t, got)
					if index == 0 {
						control = got
					} else if !reflect.DeepEqual(got, control) {
						t.Error("review arrival order changed report")
					}
				}
			})
		}
	}
}

func TestIntegrity0170IdenticalDeclarationsAreMetricIdempotent(t *testing.T) {
	pkg := integrity0170Package()
	j := integrity0170Judgment("c")
	r := domain.HumanReview{CandidateID: "c", Action: "CONFIRMED"}
	one := ValidateAssistedJudgments(pkg, []domain.AssistedJudgment{j}, []domain.HumanReview{r})
	got := ValidateAssistedJudgments(pkg, []domain.AssistedJudgment{j, j, j}, []domain.HumanReview{r, r})
	if one.Verdict != "ORACLE_ACCEPTED" || got.Verdict != one.Verdict || !reflect.DeepEqual(got.Metrics, one.Metrics) || got.EffectiveVerdict("c") != one.EffectiveVerdict("c") {
		t.Fatalf("identical repeats changed effective metrics: %+v %+v", one, got)
	}
	record := got.Record("c")
	if len(record.AIHistory) != 3 || len(record.HumanHistory) != 2 {
		t.Fatal("identical declaration history collapsed", record)
	}
	if integrity0170Has(record.EscalationCauses, "AMBIGUOUS_AI_JUDGMENTS") || integrity0170Has(record.EscalationCauses, "AMBIGUOUS_HUMAN_REVIEWS") {
		t.Fatal("identical declarations deemed ambiguous")
	}
	if !reflect.DeepEqual(DetermineHumanEscalations(pkg, []domain.AssistedJudgment{j}), DetermineHumanEscalations(pkg, []domain.AssistedJudgment{j, j, j})) {
		t.Fatal("identical repeats changed escalation sample")
	}
}

func TestIntegrity0170AnyInvalidJudgmentDefeatsGrounding(t *testing.T) {
	pkg := integrity0170Package()
	valid := integrity0170Judgment("c")
	invalid := integrity0170Clone(valid)
	invalid.Evidence[0].FragmentID = "missing"
	r := domain.HumanReview{CandidateID: "c", Action: "CONFIRMED"}
	for _, js := range [][]domain.AssistedJudgment{{invalid, valid}, {valid, invalid}} {
		got := ValidateAssistedJudgments(pkg, js, []domain.HumanReview{r, r})
		if got.Verdict != "FAIL_EVIDENCE_GROUNDING" || got.Metrics.EvidenceGroundingRate != 0 {
			t.Fatalf("last valid declaration masked invalid grounding: %+v", got)
		}
		if got.EffectiveVerdict("c") != "UNKNOWN" {
			t.Error("ambiguous declaration became effective")
		}
		integrity0170NoFabricatedReview(t, got)
	}
}

func TestIntegrity0170EscalationCannotHideDuplicateOutsideSample(t *testing.T) {
	pkg := domain.BlindCandidatePackage{SchemaVersion: "1"}
	js := []domain.AssistedJudgment{}
	for _, id := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t"} {
		pkg.Candidates = append(pkg.Candidates, domain.BlindCandidate{CandidateID: id, Fragments: []domain.BlindEvidenceFragment{{FragmentID: "f", Content: "literal evidence"}}})
		js = append(js, integrity0170Judgment(id))
	}
	sample := DetermineHumanEscalations(pkg, js)
	index := -1
	for i, j := range js {
		if !integrity0170Has(sample, j.CandidateID) {
			index = i
			break
		}
	}
	if index < 0 {
		t.Fatal("no residual candidate")
	}
	other := integrity0170Clone(js[index])
	other.Reasons = []string{"different statement"}
	for _, input := range [][]domain.AssistedJudgment{append(append([]domain.AssistedJudgment{}, js...), other), append([]domain.AssistedJudgment{other}, js...)} {
		if !integrity0170Has(DetermineHumanEscalations(pkg, input), other.CandidateID) {
			t.Fatal("distinct high confidence duplicate escaped mandatory escalation")
		}
	}
}
