/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_execution.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"sort"
)

type MarketDiscoveryService struct{ registry *DiscoveryConnectorRegistry }
type MarketDiscoveryExecution struct {
	Result    domain.MarketDiscoveryResult
	Artifacts domain.MarketDiscoveryArtifacts
}
type marketExecutionFailure struct{ code, message string }

func (e marketExecutionFailure) Error() string { return e.message }
func NewMarketDiscoveryService(registry *DiscoveryConnectorRegistry) *MarketDiscoveryService {
	return &MarketDiscoveryService{registry: registry}
}
func (s *MarketDiscoveryService) DiscoverMarket(ctx context.Context, r domain.MarketDiscoveryStartRequest) (domain.MarketDiscoveryResult, error) {
	return s.DiscoverMarketObserved(ctx, r, nil)
}
func (s *MarketDiscoveryService) DiscoverMarketObserved(ctx context.Context, r domain.MarketDiscoveryStartRequest, p domain.MarketDiscoveryProgressReporter) (domain.MarketDiscoveryResult, error) {
	out, err := s.DiscoverMarketExecutionObserved(ctx, r, p)
	return out.Result, err
}
func selectedMarketCapabilities(registry *DiscoveryConnectorRegistry, request domain.MarketDiscoveryStartRequest) []domain.DiscoveryConnectorCapabilities {
	out := []domain.DiscoveryConnectorCapabilities{}
	for _, capability := range registry.Capabilities() {
		if marketContains(request.ConnectorIDs, capability.ID) {
			out = append(out, capability)
		}
	}
	return out
}
func reportMarketProgress(reporter domain.MarketDiscoveryProgressReporter, stage domain.RecommendationRunStage, completed int) error {
	if reporter != nil {
		reporter.ReportMarketDiscoveryProgress(domain.MarketDiscoveryProgressEvent{Stage: stage, Progress: domain.RecommendationRunProgress{Planned: 4, Completed: completed}})
		if state, ok := reporter.(interface{ MarketProgressError() error }); ok {
			return state.MarketProgressError()
		}
	}
	return nil
}
func (s *MarketDiscoveryService) DiscoverMarketExecutionObserved(ctx context.Context, request domain.MarketDiscoveryStartRequest, reporter domain.MarketDiscoveryProgressReporter) (MarketDiscoveryExecution, error) {
	if s == nil || s.registry == nil {
		return MarketDiscoveryExecution{}, domain.ErrInvalidRunInput
	}
	normalized, validationErr := ValidateMarketDiscoveryStartRequest(request, s.registry)
	if validationErr != nil {
		return MarketDiscoveryExecution{}, marketExecutionFailure{"INVALID_PLAN", "market plan invalid"}
	}
	request = normalized

	selectedCaps := selectedMarketCapabilities(s.registry, request)
	plan, err := PlanMarketDiscovery(request, s.registry)
	if err != nil {
		return MarketDiscoveryExecution{}, marketExecutionFailure{"INVALID_PLAN", "market plan invalid"}
	}
	plan = AllocateMarketDiscoveryBudget(plan, request.Budget)
	if err := reportMarketProgress(reporter, domain.RunStageMarketDiscover, 1); err != nil {
		return MarketDiscoveryExecution{}, err
	}
	discoveries := map[string]domain.ConnectorDiscovery{}
	queries := map[string]domain.MarketQueryProposal{}
	for _, query := range request.QueryProposals {
		queries[query.ID] = query
	}
	accepted := 0
	for _, attempt := range plan {
		if ctx.Err() != nil {
			return MarketDiscoveryExecution{}, ctx.Err()
		}
		if attempt.Status != domain.AttemptPlanned && (attempt.Status != domain.AttemptMarketNotFiltered || attempt.PlannedCost <= 0) {
			continue
		}
		connector, ok := s.registry.Resolve(attempt.ConnectorID)
		if !ok {
			discoveries[attempt.AttemptID] = domain.ConnectorDiscovery{AttemptID: attempt.AttemptID, Status: domain.AttemptUnsupportedCapability, Candidates: []domain.ProductCandidate{}, ErrorClass: "UNSUPPORTED_CAPABILITY"}
			continue
		}
		remaining := request.Budget.MaxTotalCandidates - accepted
		if remaining <= 0 {
			discoveries[attempt.AttemptID] = domain.ConnectorDiscovery{AttemptID: attempt.AttemptID, Status: domain.AttemptBudgetExhausted, Candidates: []domain.ProductCandidate{}, ErrorClass: "BUDGET_EXHAUSTED"}
			continue
		}
		limit := request.Budget.MaxResultsPerQuery
		if limit > remaining {
			limit = remaining
		}
		query := queries[attempt.QueryID]
		discovery, discoverErr := connector.Discover(ctx, domain.MarketDiscoveryRequest{QueryID: attempt.QueryID, QueryLocale: attempt.QueryLocale, Query: query.Query, Market: attempt.Market, Limit: limit})
		if ctx.Err() != nil {
			return MarketDiscoveryExecution{}, ctx.Err()
		}
		if discoverErr != nil {
			discovery = domain.ConnectorDiscovery{Status: domain.AttemptUnknownProvider, Candidates: []domain.ProductCandidate{}, ErrorClass: "UNKNOWN_PROVIDER"}
		}
		discovery.AttemptID = attempt.AttemptID
		if discovery.Candidates == nil {
			discovery.Candidates = []domain.ProductCandidate{}
		}
		if discovery.Status == "" {
			if len(discovery.Candidates) == 0 {
				discovery.Status = domain.AttemptEmptyObserved
			} else {
				discovery.Status = domain.AttemptSucceeded
			}
		}
		if len(discovery.Candidates) > limit {
			discovery.Candidates = discovery.Candidates[:limit]
		}
		accepted += len(discovery.Candidates)
		discoveries[attempt.AttemptID] = discovery
	}
	if err := reportMarketProgress(reporter, domain.RunStageMarketResolve, 2); err != nil {
		return MarketDiscoveryExecution{}, err
	}
	discoveryList := make([]domain.ConnectorDiscovery, 0, len(discoveries))
	for _, value := range discoveries {
		discoveryList = append(discoveryList, value)
	}
	sort.Slice(discoveryList, func(i, j int) bool { return discoveryList[i].AttemptID < discoveryList[j].AttemptID })
	artifacts, err := ProjectMarketDiscoveryArtifacts(selectedCaps, plan, discoveryList)
	if err != nil {
		return MarketDiscoveryExecution{}, marketExecutionFailure{"INTEGRITY_ERROR", "market projection failed"}
	}
	if err := reportMarketProgress(reporter, domain.RunStageMarketPublish, 3); err != nil {
		return MarketDiscoveryExecution{}, err
	}
	result := domain.MarketDiscoveryResult{SchemaVersion: "1", Status: marketStatusFor(artifacts.MarketSnapshot, artifacts.MarketGaps), Capabilities: artifacts.ConnectorCapabilities, Plan: artifacts.MarketPlan, Discoveries: discoveryList, Candidates: artifacts.MarketCandidates, Gaps: artifacts.MarketGaps}
	return MarketDiscoveryExecution{Result: result, Artifacts: artifacts}, nil
}
func marketStatusFor(snapshots []domain.MarketSnapshotEntry, gaps []domain.MarketGap) domain.MarketDiscoveryStatus {
	selected, completed := 0, 0
	for _, snapshot := range snapshots {
		if snapshot.PlannedCost > 0 {
			selected++
			if snapshot.Status == domain.AttemptSucceeded || snapshot.Status == domain.AttemptEmptyObserved {
				completed++
			}
		}
	}
	if selected == 0 || completed*100 < selected*80 {
		return domain.MarketDiscoveryInconclusive
	}
	if len(gaps) > 0 {
		return domain.MarketDiscoveryPartial
	}
	return domain.MarketDiscoveryComplete
}
