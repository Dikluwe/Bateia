/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_observed.md
 * @layer L1
 * @updated 2026-09-01
 */

package usecase

import (
	"context"

	"bateia/01_core/domain"
)

// RecommendObserved executa o mesmo pipeline de Recommend e publica somente
// marcos enumerados pelo Bateia.
func (s *RecommendationService) RecommendObserved(ctx context.Context, request domain.RecommendationRequest, reporter domain.RecommendationProgressReporter) (domain.RecommendationResult, error) {
	execution, err := s.RecommendObservedWithRun(ctx, request, reporter)
	return execution.Recommendation, err
}

// RecommendObservedWithRun keeps the frozen evidence in memory for optional
// mechanism extraction during the durable publish stage.
func (s *RecommendationService) RecommendObservedWithRun(ctx context.Context, request domain.RecommendationRequest, reporter domain.RecommendationProgressReporter) (RecommendationExecution, error) {
	stages := []domain.RecommendationRunStage{
		domain.RunStageExpand, domain.RunStageDiscover, domain.RunStageScreen, domain.RunStageCollect,
		domain.RunStageAssess, domain.RunStageDecide, domain.RunStageAlternatives, domain.RunStagePublish,
	}
	if reporter != nil {
		reporter.ReportRecommendationProgress(domain.RecommendationProgressEvent{Stage: stages[0], Progress: domain.RecommendationRunProgress{Planned: len(stages)}})
	}
	execution, err := s.RecommendWithRun(ctx, request)
	if ctx.Err() != nil {
		return execution, ctx.Err()
	}
	if err != nil {
		return execution, err
	}
	if reporter != nil {
		for index, stage := range stages {
			reporter.ReportRecommendationProgress(domain.RecommendationProgressEvent{Stage: stage, Progress: domain.RecommendationRunProgress{Planned: len(stages), Completed: index + 1}})
		}
	}
	return execution, nil
}
