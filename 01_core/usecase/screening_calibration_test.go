/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/probabilistic_screening.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func calibrationGroup(id string, relevantFirst bool) domain.ScreeningDatasetGroup {
	snapshot := screeningFixture()
	snapshot.Candidates = snapshot.Candidates[:2]
	first, second := domain.RelevanceRelevant, domain.RelevanceIrrelevant
	if !relevantFirst {
		first, second = second, first
	}
	return domain.ScreeningDatasetGroup{IdeaID: id, Snapshot: snapshot, Judgments: domain.DiscoveryJudgments{SchemaVersion: "1", Judgments: []domain.CandidateJudgment{
		{Repository: "org/a", Label: first, Reason: "fixture"}, {Repository: "org/b", Label: second, Reason: "fixture"}, {Repository: "org/c", Label: domain.RelevanceUncertain, Reason: "ignored"},
	}}}
}

func TestCalibrateScreeningKeepsIdeasOutsideTrainingFold(t *testing.T) {
	groups := []domain.ScreeningDatasetGroup{calibrationGroup("one", true), calibrationGroup("two", true), calibrationGroup("three", false)}
	report, err := CalibrateScreening(groups, DefaultScreeningModel())
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Folds) != 3 || report.ObservationCount != 6 {
		t.Fatalf("folds/observações incorretos: %+v", report)
	}
	for _, fold := range report.Folds {
		if fold.TrainingCount != 4 || fold.TestCount != 2 {
			t.Fatalf("vazamento ou contagem incorreta: %+v", fold)
		}
	}
}

func TestCalibrationMetricsHandleEmptySets(t *testing.T) {
	if brier, loss := calibrationMetrics(nil, domain.CalibrationParameters{}); brier != 0 || loss != 0 {
		t.Fatalf("métrica vazia inválida: %v %v", brier, loss)
	}
	if got := precisionAtK(nil, domain.CalibrationParameters{}, 10); got != 0 {
		t.Fatalf("precision@k vazia inválida: %v", got)
	}
}

func TestSelectAuditCandidatesIncludesKnownAndNegativeControls(t *testing.T) {
	screening := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{
		{Repository: "o/relevant", RankingScore: .8}, {Repository: "o/uncertain", RankingScore: .5}, {Repository: "o/low", RankingScore: .1}, {Repository: "o/boundary", RankingScore: .4},
	}}
	judgments := domain.DiscoveryJudgments{Judgments: []domain.CandidateJudgment{
		{Repository: "o/relevant", Label: domain.RelevanceRelevant}, {Repository: "o/uncertain", Label: domain.RelevanceUncertain}, {Repository: "o/low", Label: domain.RelevanceIrrelevant}, {Repository: "o/boundary", Label: domain.RelevanceIrrelevant},
	}}
	selection, err := SelectAuditCandidates(screening, judgments, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(selection.Candidates) != 4 {
		t.Fatalf("seleção incompleta: %+v", selection)
	}
}
