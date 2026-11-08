/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
)

func ExpandInsightNeighborhood(ctx context.Context, request domain.InsightNeighborhoodRequest, bindings []domain.InsightSiteBinding) (domain.InsightSiteReport, error) {
	if ctx == nil {
		return domain.InsightSiteReport{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.InsightSiteReport{}, err
	}
	admission, err := admitInsightNeighborhood(request)
	if err != nil {
		return domain.InsightSiteReport{}, err
	}
	request = admission.Request
	plan, err := PlanInsightSites(neighborhoodSiteRequest(request), bindings)
	if err != nil {
		return domain.InsightSiteReport{}, err
	}
	plan.SchedulingPolicy, plan.ID = "memory-neighborhood-v1", ""
	plan.ID = "site_plan_" + insightHash(plan)
	ctx, stop, err := domain.ExecutionTimeout(ctx, plan.Request.Discovery.TotalTimeout)
	if err != nil {
		return domain.InsightSiteReport{}, err
	}
	defer stop()
	report := newInsightNeighborhoodReport(plan, request.Input.ID)
	beginInsightExplorationMemory(&report)
	if err := prepareInsightNeighborhoodContext(&report); err != nil {
		return domain.InsightSiteReport{}, err
	}
	order := prepareInsightNeighborhoodLeads(&report, admission.Origins)
	byID := map[string]domain.InsightSiteBinding{}
	for _, binding := range bindings {
		byID[binding.Connector.Capabilities().ID] = binding
	}
	pending, distances := insightNeighborhoodFrontier(&report, request.Focus, byID, order)
	failed, err := runInsightFollowupQueue(ctx, &report, byID, pending, order, distances)
	if err != nil {
		return domain.InsightSiteReport{}, err
	}
	err = completeInsightNeighborhood(&report, len(pending), failed, ctx.Err())
	return report, err
}
