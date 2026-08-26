/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/discovery_evaluator_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func evaluationSnapshot() domain.DiscoverySnapshot {
	return domain.DiscoverySnapshot{SchemaVersion: "1", Strategy: "lexical",
		PerExpansion: []domain.ExpansionDiscovery{
			{Expansion: domain.SearchExpansion{ID: "a"}, Provider: domain.ProviderDiscovery{Candidates: []domain.RepositoryCandidate{candidate("o", "one"), candidate("o", "two")}}},
			{Expansion: domain.SearchExpansion{ID: "b"}, Provider: domain.ProviderDiscovery{Candidates: []domain.RepositoryCandidate{candidate("o", "one")}}},
		},
		Candidates: []domain.DiscoveredCandidate{
			{Repository: candidate("o", "one"), ExpansionIDs: []string{"a", "b"}},
			{Repository: candidate("o", "two"), ExpansionIDs: []string{"a"}},
			{Repository: candidate("o", "three"), ExpansionIDs: []string{"a"}},
		},
		Failures: []domain.ExecutionFailure{{Scope: "b"}},
	}
}

func TestEvaluateDiscoverySeparatesUnjudgedAndComputesMetrics(t *testing.T) {
	judgments := domain.DiscoveryJudgments{SchemaVersion: "1", Judgments: []domain.CandidateJudgment{
		{Repository: "o/one", Label: domain.RelevanceRelevant, Reason: "resolve a capacidade"},
		{Repository: "o/two", Label: domain.RelevanceIrrelevant, Reason: "domínio diferente"},
	}}
	evaluation, err := EvaluateDiscovery(evaluationSnapshot(), judgments)
	if err != nil {
		t.Fatal(err)
	}
	if evaluation.Metrics.JudgedPrecision == nil || *evaluation.Metrics.JudgedPrecision != 0.5 {
		t.Fatalf("precisão incorreta: %+v", evaluation.Metrics)
	}
	if evaluation.Metrics.UnjudgedCount != 1 || evaluation.Metrics.OverlapRate != 1.0/3.0 || evaluation.Metrics.ExpansionHitRate != 1 {
		t.Fatalf("métricas incorretas: %+v", evaluation.Metrics)
	}
}

func TestEvaluateDiscoveryUsesNullPrecisionWithoutJudgedCandidates(t *testing.T) {
	evaluation, err := EvaluateDiscovery(evaluationSnapshot(), domain.DiscoveryJudgments{SchemaVersion: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if evaluation.Metrics.JudgedPrecision != nil || evaluation.Metrics.UnjudgedCount != 3 {
		t.Fatalf("não julgados confundidos: %+v", evaluation.Metrics)
	}
}

func TestEvaluateDiscoveryRejectsDuplicateJudgment(t *testing.T) {
	judgments := domain.DiscoveryJudgments{SchemaVersion: "1", Judgments: []domain.CandidateJudgment{
		{Repository: "o/one", Label: domain.RelevanceRelevant, Reason: "one"},
		{Repository: "O/ONE", Label: domain.RelevanceIrrelevant, Reason: "two"},
	}}
	if _, err := EvaluateDiscovery(evaluationSnapshot(), judgments); err == nil {
		t.Fatal("julgamento duplicado aceito")
	}
}
