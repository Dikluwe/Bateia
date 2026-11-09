/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/integrity_0171_metrics_test.md
 * @layer L1
 * @updated 2026-11-09
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"testing"
)

func TestIntegrity0171InvalidReviewsCannotInflateMetrics(t *testing.T) {
	for _, ai := range []string{"PARTIAL", "UNCERTAIN", "CONFLICT"} {
		for _, review := range []domain.HumanReview{
			{Action: "INVALID", Verdict: "PARTIAL", Reason: "explanation"},
			{Action: "OVERRIDDEN", Verdict: "RELEVANT"},
			{Action: "RESOLVED", Verdict: "IRRELEVANT", Reason: " \t\n "},
			{Action: "OVERRIDDEN", Verdict: "UNKNOWN", Reason: "explanation"},
			{Action: "RESOLVED", Verdict: "NOT_A_VERDICT", Reason: "explanation"},
		} {
			t.Run(ai+"/"+review.Action+"/"+review.Verdict+"/"+review.Reason, func(t *testing.T) {
				pkg := integrity0170Package()
				j := integrity0170Judgment("c")
				j.Verdict = ai
				j.Confidence = .5
				review.CandidateID = "c"
				before, _ := json.Marshal([]any{pkg, j, review})
				got := ValidateAssistedJudgments(pkg, []domain.AssistedJudgment{j}, []domain.HumanReview{review})
				m := got.Metrics
				if m.Reviewed != 0 || m.DefinitiveReviewed != 0 || m.DefinitiveAgreements != 0 || m.DefinitiveOverrides != 0 || m.ResolvedUncertain != 0 || m.RawAgreementRate != 0 || m.OverrideRate != 0 {
					t.Errorf("invalid review counted: %+v", m)
				}
				if got.EffectiveVerdict("c") != "UNKNOWN" || got.Verdict != "INCONCLUSIVE_HUMAN_REVIEW" || m.EvidenceGroundingRate != 1 || m.Escalated != 1 {
					t.Errorf("invalid effective/grounding: %+v", got)
				}
				if len(got.Record("c").AIHistory) != 1 || len(got.Record("c").HumanHistory) != 1 {
					t.Error("history lost")
				}
				after, _ := json.Marshal([]any{pkg, j, review})
				if string(before) != string(after) {
					t.Error("inputs mutated")
				}
			})
		}
	}
}

func TestIntegrity0171ValidReviewMetricsAndMixedDenominators(t *testing.T) {
	pkg := domain.BlindCandidatePackage{SchemaVersion: "1"}
	var js []domain.AssistedJudgment
	var rs []domain.HumanReview
	for _, c := range []struct{ id, ai, action, verdict, reason string }{
		{"confirm", "PARTIAL", "CONFIRMED", "IRRELEVANT", ""},
		{"override", "PARTIAL", "OVERRIDDEN", "RELEVANT", " correction "},
		{"uncertain", "UNCERTAIN", "RESOLVED", "IRRELEVANT", "resolved"},
		{"conflict", "CONFLICT", "RESOLVED", "PARTIAL", "resolved"},
		{"bad", "PARTIAL", "INVALID", "PARTIAL", "reason"},
	} {
		pkg.Candidates = append(pkg.Candidates, domain.BlindCandidate{CandidateID: c.id, Fragments: []domain.BlindEvidenceFragment{{FragmentID: "f", Content: "literal evidence"}}})
		j := integrity0170Judgment(c.id)
		j.Verdict = c.ai
		j.Confidence = .5
		js = append(js, j)
		rs = append(rs, domain.HumanReview{CandidateID: c.id, Action: c.action, Verdict: c.verdict, Reason: c.reason})
	}
	got := ValidateAssistedJudgments(pkg, js, rs)
	m := got.Metrics
	if m.Reviewed != 4 || m.DefinitiveReviewed != 2 || m.DefinitiveAgreements != 1 || m.DefinitiveOverrides != 1 || m.ResolvedUncertain != 2 || m.RawAgreementRate != .5 || m.OverrideRate != .5 || m.EvidenceGroundingRate != 1 || m.Escalated != 5 {
		t.Fatalf("wrong mixed denominators: %+v", m)
	}
	for id, want := range map[string]string{"confirm": "PARTIAL", "override": "RELEVANT", "uncertain": "IRRELEVANT", "conflict": "PARTIAL", "bad": "UNKNOWN"} {
		if got.EffectiveVerdict(id) != want {
			t.Errorf("%s effective=%s", id, got.EffectiveVerdict(id))
		}
	}
	// Idempotence applies to metrics/effective while duplicate declarations remain in history.
	duplicate := ValidateAssistedJudgments(pkg, append(js, js...), append(rs, rs...))
	if !reflect.DeepEqual(m, duplicate.Metrics) {
		t.Fatal("identical declarations changed metrics", duplicate.Metrics)
	}
	for _, r := range got.Records {
		if duplicate.EffectiveVerdict(r.CandidateID) != r.Effective {
			t.Fatal("duplicate changed effective")
		}
	}
	// Invalid review cannot hide evidence failure or promote grounding.
	js[4].Evidence[0].FragmentID = "missing"
	failed := ValidateAssistedJudgments(pkg, js, rs)
	if failed.Verdict != "FAIL_EVIDENCE_GROUNDING" || failed.Metrics.EvidenceGroundingRate != .8 {
		t.Fatal("grounding precedence changed", failed)
	}
}
