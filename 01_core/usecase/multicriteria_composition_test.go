/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/multicriteria_composition_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestMulticriteriaFrontierIncludesVerifiedComposition(t *testing.T) {
	reuse, aID, bID := compositionFixture()
	composedID := "compose:" + aID + "+" + bID
	request := domain.MulticriteriaDecisionRequest{
		Idea:        reuse.Idea,
		Assessment:  domain.AssessmentReport{SchemaVersion: "1", Assessments: reuse.Assessments},
		Policy:      reuse.Policy,
		Composition: verifiedComposition(aID, bID),
		Criteria: []domain.MulticriteriaCriterion{
			{ID: "coverage", Direction: domain.CriterionMaximize, ScaleMin: 0, ScaleMax: 1},
			{ID: "effort", Direction: domain.CriterionMinimize, ScaleMin: 0, ScaleMax: 10},
		},
		Measurements: []domain.CriterionMeasurement{
			{Candidate: "o/a@abc", CriterionID: "coverage", Status: domain.EvidenceProven, Interval: domain.CriterionInterval{Low: .33, High: .33}, Sources: []domain.EvidenceSource{decisionSource("a", "cap-core")}},
			{Candidate: "o/a@abc", CriterionID: "effort", Status: domain.EvidenceProven, Interval: domain.CriterionInterval{Low: 2, High: 2}, Sources: []domain.EvidenceSource{decisionSource("a", "cap-core")}},
			{Candidate: "o/b@abc", CriterionID: "coverage", Status: domain.EvidenceProven, Interval: domain.CriterionInterval{Low: .67, High: .67}, Sources: []domain.EvidenceSource{decisionSource("b", "cap-export")}},
			{Candidate: "o/b@abc", CriterionID: "effort", Status: domain.EvidenceProven, Interval: domain.CriterionInterval{Low: 3, High: 3}, Sources: []domain.EvidenceSource{decisionSource("b", "cap-export")}},
			{Candidate: composedID, CriterionID: "coverage", Status: domain.EvidenceProven, Interval: domain.CriterionInterval{Low: 1, High: 1}, Sources: []domain.EvidenceSource{compositionTestSource("a")}},
			{Candidate: composedID, CriterionID: "effort", Status: domain.EvidenceProven, Interval: domain.CriterionInterval{Low: 4, High: 4}, Sources: []domain.EvidenceSource{compositionTestSource("a")}},
		},
	}
	report, err := AnalyzeMulticriteriaDecision(request)
	if err != nil {
		t.Fatal(err)
	}
	if report.Composition == nil || len(report.Composition.Candidates) != 1 || !contains(report.Frontier, composedID) {
		t.Fatalf("composition omitted from frontier: %+v", report)
	}
	request.Composition.Verifications[0].Passed = false
	report, err = AnalyzeMulticriteriaDecision(request)
	if err != nil {
		t.Fatal(err)
	}
	if contains(report.Frontier, composedID) {
		t.Fatal("failed composition entered frontier")
	}
	found := false
	for _, eligibility := range report.Eligibility {
		if eligibility.Candidate == composedID {
			found = eligibility.Status == domain.CandidateIneligible
		}
	}
	if !found {
		t.Fatal("missing composition conflict in eligibility")
	}
}
