/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_observed_test.md
 * @layer L1
 * @updated 2026-09-01
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestRecommendationServiceSatisfiesObservableExecutor(t *testing.T) {
	var executor domain.ObservableRecommendationExecutor = (*RecommendationService)(nil)
	if _, ok := executor.(*RecommendationService); !ok {
		t.Fatal("recommendation service no longer satisfies observable executor")
	}
}
