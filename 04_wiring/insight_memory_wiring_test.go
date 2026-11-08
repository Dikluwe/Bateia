/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/insight_memory_wiring_test.md
 * @layer L4
 * @updated 2026-10-28
 */
package main

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/03_infra/github"
	"testing"
)

func TestInsightMemoryWiringHasDedicatedGitHubFollowReader(t *testing.T) {
	t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "")
	bindings, err := insightDiscoverySites(github.NewClient(""))
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, binding := range bindings {
		expected := binding.Connector.Capabilities().ID == domain.ConnectorGitHubCatalog
		if (binding.FollowReader != nil) != expected {
			t.Fatal("follow reader crossed sources")
		}
		if expected {
			count++
		}
	}
	if count != 1 {
		t.Fatal("missing GitHub follow reader")
	}
	if _, err := usecase.PlanInsightSites(domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "new possibilities"}, Sources: []string{"all"}}, bindings); err != nil {
		t.Fatal(err)
	}
}
