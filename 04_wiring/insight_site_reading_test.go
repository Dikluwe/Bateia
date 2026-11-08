/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/insight_site_reading_test.md
 * @layer L4
 * @updated 2026-10-28
 */
package main

import (
	"bateia/01_core/domain"
	"bateia/03_infra/github"
	"testing"
)

func TestInsightSiteReadmeCompositionBindsDedicatedForgeReaders(t *testing.T) {
	t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "")
	t.Setenv("BATEIA_GITEE_TOKEN", "")
	t.Setenv("BATEIA_ATOMGIT_TOKEN", "")
	bindings, err := insightDiscoverySites(github.NewClient(""))
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, binding := range bindings {
		id := binding.Connector.Capabilities().ID
		expected := id == domain.ConnectorGiteeCatalog || id == domain.ConnectorAtomGitCatalog || id == domain.ConnectorCodebergCatalog || id == domain.ConnectorGitLabCatalog
		if (binding.RepositoryReader != nil) != expected {
			t.Fatal("wrong site reader", id)
		}
		if expected {
			count++
			if binding.Reader != nil {
				t.Fatal("legacy reader also bound")
			}
		}
	}
	if count != 4 {
		t.Fatal("forge readers missing")
	}
}
