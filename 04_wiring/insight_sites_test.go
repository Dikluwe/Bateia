/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/insight_sites_test.md
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

func TestInsightSiteCompositionUsesFederatedConnectorsAndOnlyGitHubReader(t *testing.T) {
	t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "")
	t.Setenv("BATEIA_GITEE_TOKEN", "")
	t.Setenv("BATEIA_ATOMGIT_TOKEN", "")
	bindings, err := insightDiscoverySites(github.NewClient(""))
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, binding := range bindings {
		id := binding.Connector.Capabilities().ID
		seen[id] = true
		if (binding.Reader != nil) != (id == domain.ConnectorGitHubCatalog) {
			t.Fatal("reader crossed site boundary")
		}
		if (id == domain.ConnectorGitHubCatalog || id == domain.ConnectorSourcegraph) && binding.Connector.Capabilities().Version != "insight-repository-bridge-v1" {
			t.Fatal("capability-only source not replaced")
		}
	}
	if len(bindings) != 11 || !seen[domain.ConnectorGiteeCatalog] || !seen[domain.ConnectorAtomGitCatalog] || !seen[domain.ConnectorCodebergCatalog] || !seen[domain.ConnectorJOSS] {
		t.Fatal("federated sources missing")
	}
	if _, err := usecase.PlanInsightSites(domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "queue"}, Sources: []string{"all"}}, bindings); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "example")
	bindings, err = insightDiscoverySites(github.NewClient(""))
	if err != nil || len(bindings) != 12 {
		t.Fatal("configured Bitbucket omitted", err)
	}
}
