/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/probabilistic_screening.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"math"
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
	if d.Decision != "DROP" || d.Signals.NonGoalConflict == 0 {
		t.Fatalf("NonGoal não bloqueou avanço: %+v", d)
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
