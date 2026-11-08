/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0034_android_and_forges_test.md
 * @layer L4
 * @updated 2026-10-28
 */
package main

import (
	"bateia/01_core/domain"
	"testing"
)

func TestStep0034FederatedConnectorsComposition(t *testing.T) {
	t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "")
	connectors, err := federatedMarketDiscoveryConnectors()
	if err != nil {
		t.Fatalf("erro ao compor conectores federados: %v", err)
	}

	foundIDs := map[string]bool{}
	for _, c := range connectors {
		if foundIDs[c.Capabilities().ID] {
			t.Fatal("duplicate connector")
		}
		foundIDs[c.Capabilities().ID] = true
	}

	expectedIDs := []string{
		"apple-itunes-search-software",
		"google-play-store-apps",
		"fdroid-foss-apps",
		"github-repository-catalog",
		"sourcegraph-repository-catalog",
		"gitee-repository-catalog",
		domain.ConnectorAtomGitCatalog,
		"codeberg-repository-catalog",
		"gitlab-repository-catalog",
		domain.ConnectorZBMath,
		domain.ConnectorJOSS,
	}

	for _, wantID := range expectedIDs {
		if !foundIDs[wantID] {
			t.Fatalf("conector esperado %s não encontrado no conjunto federado", wantID)
		}
	}
}

func TestBitbucketConfiguredInDefaultAndFederatedDiscovery(t *testing.T) {
	t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "example,another")
	for _, build := range []func() ([]domain.MarketDiscoveryConnector, error){defaultMarketDiscoveryConnectors, federatedMarketDiscoveryConnectors} {
		connectors, err := build()
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, connector := range connectors {
			if connector.Capabilities().ID == domain.ConnectorBitbucketCatalog {
				found = true
			}
		}
		if !found {
			t.Fatal("Bitbucket connector missing")
		}
	}
}
