/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_result_test.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"errors"
	"testing"
)

func TestInsightNeighborhoodFeedbackRetryAndDryRun(t *testing.T) {
	first, reader, bindings := neighborhoodFixture(t)
	memory, err := UpdateInsightExplorationMemory(*first.ExplorationMemory, "https://github.com/branch/beta", "discarded", "Another investigation covers this branch")
	if err != nil {
		t.Fatal(err)
	}
	ctx := domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{})
	one := 1
	reader.fail["branch/alpha"] = true
	failed, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Memory: &memory, FollowLinks: &one}, bindings)
	if !errors.Is(err, domain.ErrInsightDiscoveryUnavailable) || failed.Status != "FAILED" || len(failed.Followups) != 1 || failed.Followups[0].Status != "AUTH_REQUIRED" {
		t.Fatal("failed reading did not preserve report", err, failed.Status)
	}
	alpha := failed.Followups[0].Identity
	entry := insightMemoryFind(failed.ExplorationMemory, alpha)
	if entry.State != "PENDING" || entry.ReadAttempts != 1 {
		t.Fatal("failure became examined/discarded")
	}
	reader.fail["branch/alpha"] = false
	resumed, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: failed}, bindings)
	if err != nil || len(reader.reads) != 1 || reader.reads[0] != "branch/alpha" {
		t.Fatal("pending failure could not resume", err, reader.reads)
	}
	zero := 0
	dry, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: resumed, FollowLinks: &zero}, bindings)
	if err != nil || dry.Status != "FRONTIER_PENDING" || len(reader.reads) != 1 || len(dry.Neighborhood.FrontierEntryIDs) != 1 || dry.Followups[0].Status != "FOLLOWUP_BUDGET" {
		t.Fatal("dry run had effects or lost frontier", err, dry.Status)
	}
}
