/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_collection.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"strings"
)

func CollectInsightEvidence(ctx context.Context, plans []domain.EvidenceExtractionPlan, provider domain.ContentProvider, request domain.InsightExplorationRequest) (domain.InsightCollectionReport, error) {
	if ctx == nil || len(plans) < 1 || len(plans) > 8 || len(request.Bundles) != 0 {
		return domain.InsightCollectionReport{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.InsightCollectionReport{}, err
	}
	request.ReviewAlternatives = true
	request, err := normalizeInsightExploration(request)
	if err != nil {
		return domain.InsightCollectionReport{}, err
	}
	type originKey struct{ repository, commit string }
	seen := map[originKey]bool{}
	for _, plan := range plans {
		if err := validateEvidenceCollectionPlan(plan, true); err != nil {
			return domain.InsightCollectionReport{}, err
		}
		commit := plan.Repository.Commit
		if mutableIdeaRunCommit(commit) {
			commit = "HEAD"
		}
		key := originKey{strings.ToLower(plan.Repository.Owner + "/" + plan.Repository.Repo), commit}
		if seen[key] {
			return domain.InsightCollectionReport{}, domain.ErrInvalidRunInput
		}
		seen[key] = true
	}
	bundles := make([]domain.EvidenceBundle, 0, len(plans))
	for _, plan := range plans {
		if err := ctx.Err(); err != nil {
			return domain.InsightCollectionReport{}, err
		}
		bundle, collectErr := CollectEvidence(ctx, plan, provider)
		if err := ctx.Err(); err != nil {
			return domain.InsightCollectionReport{}, err
		}
		if errors.Is(collectErr, domain.ErrIntegrity) {
			return domain.InsightCollectionReport{}, domain.ErrIntegrity
		}
		if collectErr != nil {
			bundle = domain.EvidenceBundle{SchemaVersion: "1", Repository: plan.Repository, Plan: plan, Complete: false, TreeComplete: false, Failures: []domain.ExecutionFailure{{Scope: plan.Repository.Key(), Stage: "resolution", Message: collectErr.Error()}}}
		}
		bundles = append(bundles, bundle)
	}
	unexplored := []string{}
	for _, bundle := range bundles {
		if !immutableCommitPattern.MatchString(bundle.Repository.Commit) {
			unexplored = append(unexplored, bundle.Repository.Key())
			continue
		}
		request.Bundles = append(request.Bundles, bundle)
	}
	exploration, err := ExploreInsights(ctx, request)
	if err != nil {
		return domain.InsightCollectionReport{}, err
	}
	if len(unexplored) > 0 {
		exploration.Gaps = append(exploration.Gaps, "UNRESOLVED_SOURCE_BUNDLES")
	}
	return domain.InsightCollectionReport{UnexploredRepositories: unexplored, SchemaVersion: "insight-collection-v1", Bundles: bundles, Exploration: exploration}, nil
}
