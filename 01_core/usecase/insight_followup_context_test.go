/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_followup_context_test.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"testing"
)

func TestInsightFollowupsPreservePartialOnCancellationAndShareCardBudget(t *testing.T) {
	reader, bindings := followFixture()
	request := insightSiteRequest()
	request.Discovery.MaxPassages = 1
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings)
	if err != nil || len(report.Tasks[0].ReadingCardRefs) != 1 {
		t.Fatal("card budget exceeded")
	}
	ctx, cancel := context.WithCancel(t.Context())
	reader.hook = cancel
	report, err = DiscoverInsightSites(domain.WithExecutionRuntime(ctx, &insightTestRuntime{}), request, bindings)
	if err == nil || report.Status != "CANCELLED" || len(report.Bundles) != 1 || report.ExplorationMemory == nil {
		t.Fatal("partial report lost on cancel")
	}
}
