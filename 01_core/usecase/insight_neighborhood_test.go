/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_test.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func neighborhoodFixture(t *testing.T) (domain.InsightSiteReport, *memoryFollowReader, []domain.InsightSiteBinding) {
	t.Helper()
	r, bindings := followFixture()
	r.contents["org/repo"] = "# Roots\n[A](https://github.com/branch/alpha)\n[B](https://github.com/branch/beta)\n"
	r.contents["branch/alpha"] = "# Alpha\n[Next](https://github.com/deeper/one)\n"
	r.contents["branch/beta"] = "# Beta\n[Other](https://github.com/unrelated/one)\n"
	r.contents["deeper/one"] = "# Leaf\nNo outgoing repository references.\n"
	zero := 0
	request := insightSiteRequest()
	request.FollowLinks = &zero
	first, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings)
	if err != nil {
		t.Fatal(err)
	}
	r.reads, r.resolves = nil, nil
	// Any catalog invocation during expansion is an immediate test failure.
	bindings[0].Connector = &insightSiteFixtureConnector{id: domain.ConnectorGitHubCatalog, auth: true, hook: func() { t.Fatal("expansion queried a catalog") }}
	return first, r, bindings
}

func TestInsightNeighborhoodFocusContinuesThroughThreeStepsWithoutSearch(t *testing.T) {
	first, reader, bindings := neighborhoodFixture(t)
	before, _ := json.Marshal(first)
	one := 1
	ctx := domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{})
	second, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Focus: "https://github.com/branch/alpha", FollowLinks: &one}, bindings)
	if err != nil || !reflect.DeepEqual(reader.reads, []string{"branch/alpha"}) || len(second.Followups) != 1 || second.Followups[0].Bundle == nil {
		t.Fatal("did not read only focused pending entry", err, reader.reads)
	}
	if second.Calls != 0 || len(second.Attempts) != 0 || len(second.Candidates) != 0 || len(second.Bundles) != 0 || second.Neighborhood.InputReportID != first.ID || second.Plan.SchedulingPolicy != "memory-neighborhood-v1" {
		t.Fatal("history became a current catalog/read result")
	}
	for _, entry := range second.ExplorationMemory.Entries {
		if entry.Identity.ExternalID == "deeper/one" {
			t.Fatal("child links were extracted in same call")
		}
	}
	after, _ := json.Marshal(first)
	if string(before) != string(after) {
		t.Fatal("input mutated")
	}
	third, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: second, Focus: "https://github.com/branch/alpha"}, bindings)
	if err != nil || !reflect.DeepEqual(reader.reads, []string{"branch/alpha", "deeper/one"}) || len(third.Followups) != 1 || len(third.Tasks[0].ReadingCardRefs) != 1 {
		t.Fatal("grandchild not reached from frozen child", err, reader.reads)
	}
	if third.Plan.Request.Discovery.Question != first.Plan.Request.Discovery.Question || third.ExplorationMemory.Generation != first.ExplorationMemory.Generation+2 || third.NoveltyState != "NOT_ASSESSED" || third.EffectState != "UNVERIFIED" {
		t.Fatal("intent/generation/evidence scope changed")
	}
	fourth, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: third, Focus: "https://github.com/branch/alpha"}, bindings)
	if err != nil || fourth.Status != "FRONTIER_EXHAUSTED" || len(fourth.Followups) != 0 || len(reader.reads) != 2 {
		t.Fatal("exhausted focus spilled into beta", err, fourth.Status, reader.reads)
	}
	if _, err := ValidateInsightExplorationMemory(fourth.ExplorationMemory); err != nil {
		t.Fatal("generated memory cannot be loaded", err)
	}
}

func TestInsightNeighborhoodCancellationAndCardQuota(t *testing.T) {
	first, reader, bindings := neighborhoodFixture(t)
	two := 2
	result, err := ExpandInsightNeighborhood(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), domain.InsightNeighborhoodRequest{Input: first, FollowLinks: &two, MaxPassages: 1}, bindings)
	if err != nil || len(reader.reads) != 2 || len(result.Tasks[0].ReadingCardRefs) != 1 {
		t.Fatal("reading/card quota lost", err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	reader.hook = cancel
	result, err = ExpandInsightNeighborhood(domain.WithExecutionRuntime(ctx, &insightTestRuntime{}), domain.InsightNeighborhoodRequest{Input: first, FollowLinks: &two}, bindings)
	if !errors.Is(err, context.Canceled) || result.Status != "CANCELLED" || result.ID == "" || result.ExplorationMemory == nil || result.Followups[0].Status != "CANCELLED" {
		t.Fatal("cancelled expansion discarded state", err, result.Status)
	}
}

func TestInsightNeighborhoodCanReopenPrimaryWithoutInventingLeadAndEnforcesCapacity(t *testing.T) {
	first, reader, bindings := neighborhoodFixture(t)
	// A root may need another explicit examination even without an incoming citation.
	memory, err := UpdateInsightExplorationMemory(*first.ExplorationMemory, "https://github.com/org/repo", "pending", "Inspect a new revision explicitly")
	if err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{"https://github.com/branch/alpha", "https://github.com/branch/beta"} {
		memory, err = UpdateInsightExplorationMemory(memory, target, "discarded", "Inspect only reopened root")
		if err != nil {
			t.Fatal(err)
		}
	}
	one := 1
	result, err := ExpandInsightNeighborhood(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), domain.InsightNeighborhoodRequest{Input: first, Memory: &memory, Focus: "https://github.com/org/repo", FollowLinks: &one}, bindings)
	if err != nil || len(reader.reads) != 1 || reader.reads[0] != "org/repo" || result.Followups[0].LeadID != "" {
		t.Fatal("explicitly reopened primary could not be read", err)
	}
	memory.Generation = 999999
	memory = sealInsightMemory(memory)
	if _, err := ValidateInsightNeighborhood(domain.InsightNeighborhoodRequest{Input: first, Memory: &memory}); err == nil {
		t.Fatal("generation overflow accepted")
	}
	// Maximum-sized documents still respect the per-read and shared quotas.
	reader.contents["branch/alpha"] = strings.Repeat("x", 32<<10)
	reader.contents["branch/beta"] = strings.Repeat("y", 32<<10)
	two := 2
	result, err = ExpandInsightNeighborhood(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), domain.InsightNeighborhoodRequest{Input: first, FollowLinks: &two}, bindings)
	if err != nil || len(result.Followups) != 2 || result.Followups[0].Bundle == nil || result.Followups[1].Bundle == nil {
		t.Fatal("bounded document quota failed", err)
	}
}
