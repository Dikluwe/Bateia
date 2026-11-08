/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/decision_trace_artifact.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

const decisionTraceArtifactContractVersion = "decision-trace-projection-v2"

func buildDecisionTraceStage(root string, graph domain.DecisionTraceGraph, index domain.DecisionTraceExplanation, versions map[string]string) (domain.IdeaRunStage, error) {
	return usecase.BuildRecommendationTraceStage(checkpointArtifactHashes{root: root}, graph, index, versions)
}
func markDecisionTraceFailed(manifest *domain.IdeaRunManifest, err error) {
	usecase.MarkRecommendationTraceFailed(manifest, err)
}
