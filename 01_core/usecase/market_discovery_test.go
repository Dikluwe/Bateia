/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_discovery_test.md
 * @layer L1
 * @updated 2026-10-28
 */

package usecase

import (
	"context"
	"encoding/json"
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

// This port observes actual durable dispatch without performing I/O in L1.
type globalCatalogProbe struct {
	companionConnector
	calls int
}

func (c *globalCatalogProbe) Discover(_ context.Context, r domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	c.calls++
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, ResponseHash: "sha256:fixture", Candidates: []domain.ProductCandidate{{Identity: domain.ProductIdentity{ConnectorID: c.capability.ID, ExternalID: "org/repo"}, Name: "org/repo", ProductKind: "repository", ObservedMarkets: []string{}}}}, nil
}

func TestGlobalCatalogExecutesWithinBudgetWithoutClaimingMarketCoverage(t *testing.T) {
	for _, budget := range []int{1, 2} {
		t.Run(string(rune('0'+budget)), func(t *testing.T) {
			probes := []*globalCatalogProbe{}
			connectors := []domain.MarketDiscoveryConnector{}
			for _, id := range []string{domain.ConnectorGiteeCatalog, domain.ConnectorAtomGitCatalog} {
				probe := &globalCatalogProbe{companionConnector: companionConnector{domain.DiscoveryConnectorCapabilities{ID: id, Version: "v1", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"repository_identity"}, Replayable: true, PublicOnly: true}}}
				probes = append(probes, probe)
				connectors = append(connectors, probe)
			}
			registry, err := NewMarketDiscoveryRegistry(connectors)
			if err != nil {
				t.Fatal(err)
			}
			store := newStep0026MemoryStore()
			service := NewRecommendationRunService(store, &step0026Queue{capacity: 2}, &step0026IDGenerator{}, nil)
			service.ConfigureMarketDiscovery(nil, registry)
			request := domain.MarketDiscoveryStartRequest{Idea: domain.IdeaSpec{Problem: "investigate tools", Capabilities: []string{"search"}}, QueryProposals: []domain.MarketQueryProposal{{Locale: "en-US", Query: "tools", SourceText: "tools", Derivation: domain.QueryExternalAgent, Rationale: "functional query"}}, Markets: []string{"JP"}, ConnectorIDs: []string{domain.ConnectorGiteeCatalog, domain.ConnectorAtomGitCatalog}, Budget: domain.MarketDiscoveryBudget{MaxQueries: budget, MaxResultsPerQuery: 1, MaxTotalCandidates: 2}}
			receipt, err := service.StartMarketDiscovery(context.Background(), "global-catalog-0126", request)
			if err != nil {
				t.Fatal(err)
			}
			if err := service.Execute(context.Background(), receipt.RunID); err != nil {
				t.Fatal(err)
			}
			if probes[0].calls+probes[1].calls != budget {
				t.Fatalf("dispatch ignored global sources or budget: %d %d", probes[0].calls, probes[1].calls)
			}
			var candidates struct {
				Entries []domain.ProductCandidate `json:"entries"`
			}
			if err := json.Unmarshal(store.artifacts[receipt.RunID][domain.ArtifactMarketCandidates], &candidates); err != nil {
				t.Fatal(err)
			}
			if len(candidates.Entries) != budget {
				t.Fatalf("candidates not published: %+v", candidates)
			}
			for _, c := range candidates.Entries {
				if len(c.ObservedMarkets) != 0 || len(c.Origins) != 1 || c.Origins[0].ObservedMarket != "" || c.Origins[0].ResponseHash != "sha256:fixture" {
					t.Fatalf("false geography or missing provenance: %+v", c)
				}
			}
			var gaps struct {
				Entries []domain.MarketGap `json:"entries"`
			}
			if err := json.Unmarshal(store.artifacts[receipt.RunID][domain.ArtifactMarketGaps], &gaps); err != nil {
				t.Fatal(err)
			}
			if len(gaps.Entries) != 2 {
				t.Fatalf("missing geographic limitations: %+v", gaps)
			}
			for _, g := range gaps.Entries {
				if g.Status != domain.AttemptMarketNotFiltered || g.Market != "JP" {
					t.Fatalf("wrong gap: %+v", g)
				}
			}
		})
	}
}

func TestMarketProjectionPreservesFilteredCoverage(t *testing.T) {
	for _, filtered := range []bool{false, true} {
		status := domain.AttemptMarketNotFiltered
		if filtered {
			status = domain.AttemptPlanned
		}
		plan := []domain.MarketDiscoveryAttempt{{AttemptID: "a", ConnectorID: "catalog", Market: "JP", Status: status, MarketFilterApplied: filtered, PlannedCost: 1}}
		discoveries := []domain.ConnectorDiscovery{{AttemptID: "a", Status: domain.AttemptSucceeded, Candidates: []domain.ProductCandidate{{Identity: domain.ProductIdentity{ConnectorID: "catalog", ExternalID: "repo"}}}}}
		result, err := ProjectMarketDiscoveryArtifacts(nil, plan, discoveries)
		if err != nil {
			t.Fatal(err)
		}
		c := result.MarketCandidates[0]
		if filtered {
			if len(c.ObservedMarkets) != 1 || c.ObservedMarkets[0] != "JP" || c.Origins[0].ObservedMarket != "JP" || len(result.MarketGaps) != 0 {
				t.Fatalf("filtered coverage lost: %+v", result)
			}
		} else if len(c.ObservedMarkets) != 0 || c.Origins[0].ObservedMarket != "" || len(result.MarketGaps) != 1 {
			t.Fatalf("unfiltered coverage invented: %+v", result)
		}
	}
}
