/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_schedule.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import "bateia/01_core/domain"

const insightSiteSchedulingPolicy = "PERSPECTIVE_SOURCE_ROUND_ROBIN_V1"

func insightSiteSchedule(plan domain.InsightSitePlan) []int {
	sources, queries := len(plan.Sources), len(plan.Discovery.Queries)
	order := []int{}
	if sources == 0 || queries == 0 {
		return order
	}
	stride := queries
	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	for gcd(stride, sources) != 1 {
		stride++
	}
	for round := 0; round < sources; round++ {
		for query := 0; query < queries; query++ {
			order = append(order, query*sources+(round*stride+query)%sources)
		}
	}
	return order
}

func recordInsightOpportunityCoverage(report *domain.InsightSiteReport) {
	attempted, completed, sources := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, attempt := range report.Attempts {
		if attempt.ExecutionOrder > 0 {
			attempted[attempt.QueryID] = true
			sources[attempt.ConnectorID] = true
		}
		if attempt.Status == "SUCCEEDED" || attempt.Status == "EMPTY" {
			completed[attempt.QueryID] = true
		}
	}
	coverage := &domain.InsightOpportunityCoverage{AttemptedPerspectives: []string{}, CompletedPerspectives: []string{}, AttemptedSources: []string{}}
	for _, query := range report.Plan.Discovery.Queries {
		if attempted[query.ID] {
			coverage.AttemptedPerspectives = append(coverage.AttemptedPerspectives, query.Perspective)
		} else {
			report.Gaps = append(report.Gaps, "NO_ATTEMPT_FOR_PERSPECTIVE:"+query.Perspective)
		}
		if completed[query.ID] {
			coverage.CompletedPerspectives = append(coverage.CompletedPerspectives, query.Perspective)
		}
	}
	for _, source := range report.Plan.Sources {
		if sources[source.ID] {
			coverage.AttemptedSources = append(coverage.AttemptedSources, source.ID)
		}
	}
	report.OpportunityCoverage = coverage
}
