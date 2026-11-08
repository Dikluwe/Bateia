/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_frontier_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"reflect"
	"testing"
)

func TestInsightNeighborhoodCyclesAndDiscardedIntermediatesDoNotEscapeFocus(t *testing.T) {
	first, reader, bindings := neighborhoodFixture(t)
	reader.contents["branch/alpha"] = "# Alpha\n[Self](https://github.com/branch/alpha)\n[Root](https://github.com/org/repo)\n[Beta](https://github.com/branch/beta)\n"
	reader.contents["branch/beta"] = "# Beta\n[Alpha](https://github.com/branch/alpha)\n"
	one := 1
	ctx := domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{})
	second, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Focus: "https://github.com/branch/alpha", FollowLinks: &one}, bindings)
	if err != nil {
		t.Fatal(err)
	}
	memory, err := UpdateInsightExplorationMemory(*second.ExplorationMemory, "https://github.com/branch/beta", "discarded", "Skip this branch")
	if err != nil {
		t.Fatal(err)
	}
	closed, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: second, Memory: &memory, Focus: "https://github.com/branch/alpha"}, bindings)
	if err != nil || closed.Status != "FRONTIER_EXHAUSTED" || len(reader.reads) != 1 {
		t.Fatal("cycle or discarded branch was read", err, reader.reads)
	}
	if _, err := ValidateInsightNeighborhood(domain.InsightNeighborhoodRequest{Input: closed, Focus: "https://github.com/branch/beta"}); err == nil {
		t.Fatal("discarded focus accepted without reopening")
	}
	third, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: second, Focus: "https://github.com/branch/alpha"}, bindings)
	if err != nil || !reflect.DeepEqual(reader.reads, []string{"branch/alpha", "branch/beta"}) {
		t.Fatal("reference from alpha to beta was not followed", err)
	}
	end, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: third, Focus: "https://github.com/branch/alpha"}, bindings)
	if err != nil || end.Status != "FRONTIER_EXHAUSTED" || len(reader.reads) != 2 {
		t.Fatal("cycle reread an examined node", err)
	}
}

func TestInsightNeighborhoodKeepsForgeNamespaceAndSourceRestrictions(t *testing.T) {
	reader, bindings := followFixture()
	reader.contents["org/repo"] = "# Roots\n[GitLab](https://gitlab.com/team/sub/project)\n"
	reader.contents["team/sub/project"] = "# Project\nA namespace-specific mechanism.\n"
	bindings = append(bindings, domain.InsightSiteBinding{Connector: &insightSiteFixtureConnector{id: domain.ConnectorGitLabCatalog, empty: true}, RepositoryReader: reader})
	request, zero := insightSiteRequest(), 0
	request.FollowLinks = &zero
	ctx := domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{})
	first, err := DiscoverInsightSites(ctx, request, bindings)
	if err != nil {
		t.Fatal(err)
	}
	one := 1
	reader.reads = nil
	second, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Focus: "https://gitlab.com/team/sub/project", FollowLinks: &one}, bindings)
	if err != nil || len(second.Followups) != 1 || second.Followups[0].Cards[0].Source.Repository != "team/sub/project" || second.Followups[0].Identity.ConnectorID != domain.ConnectorGitLabCatalog {
		t.Fatal("GitLab namespace lost or treated as GitHub", err)
	}
	// A reader in the factory does not authorize a site absent from the inherited selection.
	request.Sources = []string{"github"}
	first, err = DiscoverInsightSites(ctx, request, bindings)
	if err != nil {
		t.Fatal(err)
	}
	reader.reads = nil
	second, err = ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, FollowLinks: &one}, bindings)
	if err != nil || len(reader.reads) != 0 || second.Followups[0].Status != "SOURCE_NOT_SELECTED" {
		t.Fatal("source boundary crossed", err)
	}
}
