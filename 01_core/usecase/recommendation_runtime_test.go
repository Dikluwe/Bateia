/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_runtime_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"testing"
	"time"
)

type recommendationRuntimeStub struct{}

func (*recommendationRuntimeStub) Now() time.Time { return time.Time{} }
func (*recommendationRuntimeStub) WithTimeout(ctx context.Context, _ time.Duration) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}
func TestRecommendationRuntimeIsExplicit(t *testing.T) {
	var absent *RecommendationService
	if absent.WithExecutionRuntime(nil) != nil {
		t.Fatal("nil service changed")
	}
	service := NewRecommendationService(nil, nil, nil)
	var runtime domain.ExecutionRuntime = &recommendationRuntimeStub{}
	if service.WithExecutionRuntime(runtime) != service || service.dependencies.Runtime != runtime {
		t.Fatal("runtime not injected")
	}
	service.WithExecutionRuntime(nil)
	if service.dependencies.Runtime != nil {
		t.Fatal("runtime not cleared")
	}
}
