/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/discovery_evaluator.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

func validRelevanceLabel(label domain.RelevanceLabel) bool {
	return label == domain.RelevanceRelevant || label == domain.RelevanceIrrelevant || label == domain.RelevanceUncertain
}

// EvaluateDiscovery aplica julgamentos explícitos a um snapshot congelado.
func EvaluateDiscovery(snapshot domain.DiscoverySnapshot, judgments domain.DiscoveryJudgments) (domain.DiscoveryEvaluation, error) {
	if snapshot.SchemaVersion != "1" && snapshot.SchemaVersion != "2" {
		return domain.DiscoveryEvaluation{}, fmt.Errorf("schema de snapshot não suportado: %q", snapshot.SchemaVersion)
	}
	if judgments.SchemaVersion != "1" {
		return domain.DiscoveryEvaluation{}, fmt.Errorf("schema de julgamentos não suportado: %q", judgments.SchemaVersion)
	}
	labels := make(map[string]domain.CandidateJudgment, len(judgments.Judgments))
	for _, judgment := range judgments.Judgments {
		key := strings.ToLower(strings.TrimSpace(judgment.Repository))
		if key == "" || strings.TrimSpace(judgment.Reason) == "" || !validRelevanceLabel(judgment.Label) {
			return domain.DiscoveryEvaluation{}, fmt.Errorf("julgamento inválido para %q", judgment.Repository)
		}
		if _, exists := labels[key]; exists {
			return domain.DiscoveryEvaluation{}, fmt.Errorf("julgamento duplicado para %s", judgment.Repository)
		}
		labels[key] = judgment
	}

	evaluation := domain.DiscoveryEvaluation{SchemaVersion: "1", Strategy: snapshot.Strategy}
	relevantRepos := make(map[string]bool)
	overlap := 0
	for _, candidate := range snapshot.Candidates {
		key := repositoryKey(candidate.Repository)
		if len(candidate.ExpansionIDs) > 1 {
			overlap++
		}
		judgment, exists := labels[key]
		if !exists {
			evaluation.Metrics.UnjudgedCount++
			continue
		}
		switch judgment.Label {
		case domain.RelevanceRelevant:
			evaluation.Metrics.RelevantCount++
			relevantRepos[key] = true
		case domain.RelevanceIrrelevant:
			evaluation.Metrics.IrrelevantCount++
		case domain.RelevanceUncertain:
			evaluation.Metrics.UncertainCount++
		}
	}
	judged := evaluation.Metrics.RelevantCount + evaluation.Metrics.IrrelevantCount
	if judged > 0 {
		value := ratio(evaluation.Metrics.RelevantCount, judged)
		evaluation.Metrics.JudgedPrecision = &value
	}
	evaluation.Metrics.UncertainRate = ratio(evaluation.Metrics.UncertainCount, len(snapshot.Candidates))
	evaluation.Metrics.UnjudgedRate = ratio(evaluation.Metrics.UnjudgedCount, len(snapshot.Candidates))
	evaluation.Metrics.UniqueYield = ratio(len(snapshot.Candidates), len(snapshot.PerExpansion))
	evaluation.Metrics.OverlapRate = ratio(overlap, len(snapshot.Candidates))
	evaluation.Metrics.FailureRate = ratio(len(snapshot.Failures), len(snapshot.PerExpansion))

	hits := 0
	for _, discovery := range snapshot.PerExpansion {
		relevant := 0
		for _, candidate := range discovery.Provider.Candidates {
			if relevantRepos[repositoryKey(candidate)] {
				relevant++
			}
		}
		if relevant > 0 {
			hits++
		}
		evaluation.PerExpansion = append(evaluation.PerExpansion, domain.ExpansionDiscoveryEvaluation{
			ExpansionID: discovery.Expansion.ID, CandidateCount: len(discovery.Provider.Candidates), RelevantCount: relevant, Hit: relevant > 0,
		})
	}
	evaluation.Metrics.ExpansionHitRate = ratio(hits, len(snapshot.PerExpansion))
	return evaluation, nil
}
