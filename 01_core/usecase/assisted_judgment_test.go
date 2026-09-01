/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/assisted_judgment_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func TestAssistedJudgmentEmptyCorpus(t *testing.T) {
	if got := ValidateAssistedJudgments(domain.BlindCandidatePackage{SchemaVersion: "1"}, nil, nil).Verdict; got != "INCONCLUSIVE_NO_CANDIDATES" {
		t.Fatalf("veredito=%s", got)
	}
}

func TestAssistedJudgmentMetricsSeparateResolutionFromDefinitiveAgreement(t *testing.T) {
	pkg := domain.BlindCandidatePackage{SchemaVersion: "1", Candidates: []domain.BlindCandidate{
		{CandidateID: "a", Fragments: []domain.BlindEvidenceFragment{{FragmentID: "f", Content: "evidence"}}},
		{CandidateID: "b", Fragments: []domain.BlindEvidenceFragment{{FragmentID: "f", Content: "evidence"}}},
	}}
	judgments := []domain.AssistedJudgment{
		{CandidateID: "a", Verdict: "PARTIAL", Confidence: .9, Evidence: []domain.EvidencePointer{{CandidateID: "a", FragmentID: "f"}}, RubricVersion: AssistedRubricVersion},
		{CandidateID: "b", Verdict: "UNCERTAIN", Confidence: .5, Evidence: []domain.EvidencePointer{{CandidateID: "b", FragmentID: "f"}}, RubricVersion: AssistedRubricVersion},
	}
	reviews := []domain.HumanReview{
		{CandidateID: "a", Action: "CONFIRMED"},
		{CandidateID: "b", Action: "RESOLVED", Verdict: "IRRELEVANT", Reason: "insuficiente"},
	}
	report := ValidateAssistedJudgments(pkg, judgments, reviews)
	if report.Verdict != "ORACLE_ACCEPTED" || report.Metrics.RawAgreementRate != 1 || report.Metrics.OverrideRate != 0 || report.Metrics.ResolvedUncertain != 1 {
		t.Fatalf("métricas inesperadas: %#v", report)
	}
}
