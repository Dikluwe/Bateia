/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_runtime.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import "bateia/01_core/domain"

func (s *RecommendationService) WithExecutionRuntime(runtime domain.ExecutionRuntime) *RecommendationService {
	if s != nil {
		s.dependencies.Runtime = runtime
	}
	return s
}
