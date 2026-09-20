/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/probabilistic_screening_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"math"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func screeningFixture() domain.DiscoverySnapshot {
	return domain.DiscoverySnapshot{SchemaVersion: "2", Strategy: "lexical", Idea: domain.IdeaSpec{
		Problem: "analisar impacto estrutural do código", Capabilities: []string{"grafo de dependências"}, NonGoals: []string{"runtime monitoring"},
	}, Candidates: []domain.DiscoveredCandidate{
		{Repository: candidate("org", "a"), Queries: []string{"code impact analysis"}, Origins: []domain.CandidateOrigin{
			{ExpansionID: "one", PlannedQuery: "code impact analysis", EffectiveQuery: "code impact analysis", AttemptIndex: 0, Evidence: []domain.MatchEvidence{{MatchType: "file", Path: "analyzer.go"}}},
			{ExpansionID: "two", PlannedQuery: "dependency graph", EffectiveQuery: "dependency graph", AttemptIndex: 0, Evidence: []domain.MatchEvidence{{MatchType: "symbol"}}},
		}},
		{Repository: domain.RepositoryCandidate{Owner: "org", Repo: "b", Description: "game toolkit"}, Queries: []string{"impact"}, Origins: []domain.CandidateOrigin{{ExpansionID: "one", PlannedQuery: "code impact analysis", EffectiveQuery: "impact", AttemptIndex: 2, Evidence: []domain.MatchEvidence{{MatchType: "repository"}}}}},
		{Repository: candidate("org", "c"), Origins: []domain.CandidateOrigin{{ExpansionID: "one", PlannedQuery: "code impact analysis", EffectiveQuery: "impact", AttemptIndex: 2}}},
		{Repository: domain.RepositoryCandidate{Owner: "org", Repo: "d", Description: "runtime monitoring and code impact"}, Queries: []string{"code impact"}, Origins: []domain.CandidateOrigin{{ExpansionID: "one", PlannedQuery: "code impact", EffectiveQuery: "code impact", AttemptIndex: 0, Evidence: []domain.MatchEvidence{{MatchType: "file"}}}}},
	}}
}

func TestScreenDiscoveryExplainsControlledCandidates(t *testing.T) {
	report, err := ScreenDiscovery(screeningFixture(), DefaultScreeningModel())
	if err != nil {
		t.Fatal(err)
	}
	byRepo := make(map[string]domain.ScreeningDecision)
	for _, decision := range report.Decisions {
		byRepo[decision.Repository] = decision
	}
	a, b, c, d := byRepo["org/a"], byRepo["org/b"], byRepo["org/c"], byRepo["org/d"]
	if a.RankingScore <= b.RankingScore {
		t.Fatalf("A deveria superar B: A=%v B=%v", a.RankingScore, b.RankingScore)
	}
	if c.Decision != "AUDIT" {
		t.Fatalf("ausência deveria gerar AUDIT: %+v", c)
	}
	if d.Decision != "AUDIT" || d.Signals.NonGoalConflict == 0 {
		t.Fatalf("NonGoal sem calibração deveria gerar AUDIT: %+v", d)
	}
	if a.Probability != nil || a.CalibrationStatus != "uncalibrated" {
		t.Fatalf("score não calibrado virou probabilidade: %+v", a)
	}
	raw := 0.0
	for _, item := range a.Contributions {
		raw += item.Effect
	}
	if math.Abs(clamp(raw)-a.RankingScore) > 1e-12 {
		t.Fatalf("score não reconstruível: raw=%v score=%v", raw, a.RankingScore)
	}
}

func TestScreenDiscoveryRequiresValidCalibrationForDrop(t *testing.T) {
	decisions := func(model domain.ScreeningModel) map[string]domain.ScreeningDecision {
		t.Helper()
		report, err := ScreenDiscovery(screeningFixture(), model)
		if err != nil {
			t.Fatal(err)
		}
		result := make(map[string]domain.ScreeningDecision, len(report.Decisions))
		for _, decision := range report.Decisions {
			result[decision.Repository] = decision
		}
		return result
	}

	uncalibrated := DefaultScreeningModel()
	uncalibrated.AuditSampleRate = 0
	contained := decisions(uncalibrated)["org/d"]
	if contained.Decision != "AUDIT" || contained.Probability != nil ||
		!containsReason(contained.UncertaintyReasons, "modelo sem calibração válida não autoriza DROP") {
		t.Fatalf("modelo não calibrado descartou ou ocultou a contenção: %+v", contained)
	}

	calibrated := uncalibrated
	calibrated.CalibrationStatus = "calibrated"
	calibrated.Calibration = &domain.CalibrationParameters{A: 1, B: 0}
	authorized := decisions(calibrated)["org/d"]
	if authorized.Decision != "DROP" || authorized.Probability == nil {
		t.Fatalf("calibração válida não preservou o descarte: %+v", authorized)
	}

	malformed := uncalibrated
	malformed.CalibrationStatus = "calibrated"
	for repository, decision := range decisions(malformed) {
		if decision.Decision != "AUDIT" || decision.Probability != nil ||
			!containsReason(decision.UncertaintyReasons, "estado de calibração opaco ou malformado exige AUDIT") {
			t.Fatalf("estado malformado não foi contido em %s: %+v", repository, decision)
		}
	}
}

func TestScreenDiscoveryRejectsNonFiniteCalibrationAsAuthority(t *testing.T) {
	model := DefaultScreeningModel()
	model.CalibrationStatus = "calibrated"
	model.Calibration = &domain.CalibrationParameters{A: math.Inf(1), B: 0}
	report, err := ScreenDiscovery(screeningFixture(), model)
	if err != nil {
		t.Fatal(err)
	}
	for _, decision := range report.Decisions {
		if decision.Decision != "AUDIT" || decision.Probability != nil {
			t.Fatalf("calibrador não finito ganhou autoridade: %+v", decision)
		}
	}
}

func containsReason(reasons []string, want string) bool {
	for _, reason := range reasons {
		if strings.EqualFold(reason, want) {
			return true
		}
	}
	return false
}

func TestMigrationV1DoesNotInventEvidence(t *testing.T) {
	snapshot := evaluationSnapshot()
	snapshot.SchemaVersion = "1"
	migrated, err := MigrateDiscoverySnapshotV2(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if migrated.SchemaVersion != "2" || len(migrated.Candidates[0].Origins) == 0 || len(migrated.Candidates[0].Origins[0].Evidence) != 0 {
		t.Fatalf("migração fabricou ou perdeu origem: %+v", migrated.Candidates[0])
	}
}

func TestRemovingOriginOnlyChangesConvergenceDependentContribution(t *testing.T) {
	snapshot := screeningFixture()
	model := DefaultScreeningModel()
	before := ExtractScreeningSignals(snapshot, snapshot.Candidates[0], model)
	snapshot.Candidates[0].Origins = snapshot.Candidates[0].Origins[:1]
	after := ExtractScreeningSignals(snapshot, snapshot.Candidates[0], model)
	if before.OriginConvergence != 2 || after.OriginConvergence != 1 || before.QueryRetention != after.QueryRetention || before.DescriptionAlignment != after.DescriptionAlignment {
		t.Fatalf("remoção contaminou sinais: before=%+v after=%+v", before, after)
	}
}

func TestValidateAuditRevisionsCountsFalseNegative(t *testing.T) {
	report, err := ValidateAuditRevisions([]domain.JudgmentRevision{{Repository: "org/repo", PreviousLabel: domain.RelevanceIrrelevant, NewLabel: domain.RelevanceRelevant, Evidence: "README declara capacidade", Reason: "revisão profunda"}})
	if err != nil {
		t.Fatal(err)
	}
	if report.ReviewedCount != 1 || report.FalseNegativeCount != 1 {
		t.Fatalf("auditoria incorreta: %+v", report)
	}
}

func TestDiagnoseScreeningAdvancesKnownRelevantWithoutChangingScore(t *testing.T) {
	report := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "org/repo", RankingScore: .1, Decision: "DROP"}}}
	judgments := domain.DiscoveryJudgments{Judgments: []domain.CandidateJudgment{{Repository: "org/repo", Label: domain.RelevanceRelevant}}}
	got := DiagnoseScreening(report, judgments)
	if got.Decisions[0].Decision != "ADVANCE" || got.Decisions[0].RankingScore != .1 || report.Decisions[0].Decision != "DROP" {
		t.Fatalf("diagnóstico inválido: %+v", got)
	}
}
