/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_schedule_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"fmt"
	"strings"
	"testing"
)

func TestInsightScheduleIsCompletePermutationAndBalancesPerspectives(t *testing.T) {
	for count := 1; count <= 12; count++ {
		plan := domain.InsightSitePlan{Sources: make([]domain.DiscoveryConnectorCapabilities, count), Discovery: domain.InsightDiscoveryPlan{Queries: make([]domain.InsightDiscoveryQuery, 5)}}
		order := insightSiteSchedule(plan)
		if len(order) != count*5 {
			t.Fatal("missing cells", count)
		}
		seen, first := map[int]bool{}, map[int]bool{}
		for position, index := range order {
			if index < 0 || index >= count*5 || seen[index] {
				t.Fatal("duplicate/outside cell", count, index)
			}
			seen[index] = true
			if position < 5 {
				first[index/count] = true
			}
		}
		if len(first) != 5 {
			t.Fatal("perspective starvation", count)
		}
	}
}
func TestInsightScheduleTwentyCallsCoverElevenSourcesAndFivePerspectives(t *testing.T) {
	bindings := []domain.InsightSiteBinding{}
	for i := 0; i < 11; i++ {
		bindings = append(bindings, domain.InsightSiteBinding{Connector: &insightSiteFixtureConnector{id: fmt.Sprintf("fixture-%02d", i), empty: true}})
	}
	request := insightSiteRequest()
	request.MaxCalls = 20
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings)
	if err != nil || report.Calls != 20 || len(report.Attempts) != 55 || report.Plan.SchedulingPolicy != insightSiteSchedulingPolicy {
		t.Fatalf("schedule failed %v %+v", err, report)
	}
	coverage := report.OpportunityCoverage
	if coverage == nil || len(coverage.AttemptedPerspectives) != 5 || len(coverage.CompletedPerspectives) != 5 || len(coverage.AttemptedSources) != 11 {
		t.Fatal("opportunities lost", coverage)
	}
	counts := map[string]int{}
	orders := map[int]bool{}
	for _, attempt := range report.Attempts {
		if attempt.ExecutionOrder > 0 {
			counts[attempt.QueryID]++
			if orders[attempt.ExecutionOrder] {
				t.Fatal("duplicate execution order")
			}
			orders[attempt.ExecutionOrder] = true
		}
	}
	for _, query := range report.Plan.Discovery.Queries {
		if counts[query.ID] != 4 {
			t.Fatal("unequal opportunity", counts)
		}
	}
	if len(orders) != 20 {
		t.Fatal("logical calls not recorded")
	}
	// The previous query-major policy spent the same 20 cells on only two queries.
	oldPerspectives := map[int]bool{}
	for i := 0; i < 20; i++ {
		oldPerspectives[i/11] = true
	}
	if len(oldPerspectives) != 2 {
		t.Fatal("legacy control changed")
	}
	request.MaxCalls = 2
	report, err = DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings)
	if err != nil || len(report.OpportunityCoverage.AttemptedPerspectives) != 2 {
		t.Fatal("small quota misreported")
	}

	missing := 0
	for _, gap := range report.Gaps {
		if strings.HasPrefix(gap, "NO_ATTEMPT_FOR_PERSPECTIVE:") {
			missing++
		}
	}
	if missing != 3 {
		t.Fatal("missing perspectives not explicit", report.Gaps)
	}
}
