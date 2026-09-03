/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_discovery_test.md
 * @layer L1
 * @updated 2026-09-01
 */

package usecase

import (
	"context"
	"testing"

	"bateia/01_core/domain"
)

type companionConnector struct {
	capability domain.DiscoveryConnectorCapabilities
}

func (c companionConnector) Capabilities() domain.DiscoveryConnectorCapabilities { return c.capability }
func (companionConnector) Discover(context.Context, domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	return domain.ConnectorDiscovery{}, nil
}

func TestMarketDiscoveryGlobalConnectorPreservesMarketGap(t *testing.T) {
	connector := companionConnector{domain.DiscoveryConnectorCapabilities{ID: "global", Version: "v1", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, Markets: []string{}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{}, Replayable: true, PublicOnly: true}}
	registry, err := NewDiscoveryConnectorRegistry([]domain.MarketDiscoveryConnector{connector})
	if err != nil {
		t.Fatal(err)
	}
	request := domain.MarketDiscoveryStartRequest{Idea: domain.IdeaSpec{Problem: "p", Capabilities: []string{"c"}}, QueryProposals: []domain.MarketQueryProposal{{Locale: "en-US", Query: "asset maintenance", SourceText: "maintenance", Derivation: domain.QueryExternalAgent, Rationale: "functional wording"}}, Markets: []string{"JP"}, ConnectorIDs: []string{"global"}, Budget: domain.MarketDiscoveryBudget{MaxQueries: 1, MaxResultsPerQuery: 1, MaxTotalCandidates: 1}}
	plan, err := PlanMarketDiscovery(request, registry)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != 1 || plan[0].Status != domain.AttemptMarketNotFiltered || plan[0].QueryLocale != "en-US" || plan[0].Market != "JP" {
		t.Fatalf("unexpected plan: %+v", plan)
	}
}
