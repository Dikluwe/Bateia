/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_execution_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"testing"
)

type configuredMarketSpy struct {
	inner       *MarketDiscoveryService
	calls       int
	alterResult bool
}

func (s *configuredMarketSpy) DiscoverMarket(ctx context.Context, r domain.MarketDiscoveryStartRequest) (domain.MarketDiscoveryResult, error) {
	return s.DiscoverMarketObserved(ctx, r, nil)
}
func (s *configuredMarketSpy) DiscoverMarketObserved(ctx context.Context, r domain.MarketDiscoveryStartRequest, p domain.MarketDiscoveryProgressReporter) (domain.MarketDiscoveryResult, error) {
	s.calls++
	result, err := s.inner.DiscoverMarketObserved(ctx, r, p)
	if s.alterResult {
		result.Candidates = []domain.ProductCandidate{}
		result.Gaps = []domain.MarketDiscoveryGap{}
		result.Status = domain.MarketDiscoveryInconclusive
	}
	return result, err
}
func TestDurableMarketUsesConfiguredInvestigationExecutor(t *testing.T) {
	for _, alter := range []bool{false, true} {
		connector := &globalCatalogProbe{companionConnector: companionConnector{domain.DiscoveryConnectorCapabilities{ID: domain.ConnectorGiteeCatalog, Version: "v1", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"repository_identity"}, Replayable: true, PublicOnly: true}}}
		registry, err := NewMarketDiscoveryRegistry([]domain.MarketDiscoveryConnector{connector})
		if err != nil {
			t.Fatal(err)
		}
		executor := &configuredMarketSpy{inner: NewMarketDiscoveryService(registry), alterResult: alter}
		store := newStep0026MemoryStore()
		service := NewRecommendationRunService(store, &step0026Queue{capacity: 2}, &step0026IDGenerator{}, nil).ConfigureMarketDiscovery(executor, registry)
		request := domain.MarketDiscoveryStartRequest{Idea: domain.IdeaSpec{Problem: "investigate tools", Capabilities: []string{"search"}}, QueryProposals: []domain.MarketQueryProposal{{Locale: "en-US", Query: "tools", SourceText: "tools", Derivation: domain.QueryExternalAgent, Rationale: "functional query"}}, Markets: []string{"JP"}, ConnectorIDs: []string{domain.ConnectorGiteeCatalog}, Budget: domain.MarketDiscoveryBudget{MaxQueries: 1, MaxResultsPerQuery: 1, MaxTotalCandidates: 1}}
		receipt, err := service.StartMarketDiscovery(context.Background(), "configured-market", request)
		if err != nil {
			t.Fatal(err)
		}
		if err := service.Execute(context.Background(), receipt.RunID); err != nil {
			t.Fatal(err)
		}
		if executor.calls != 1 || connector.calls != 1 {
			t.Fatalf("executor ignored/duplicated: %d %d", executor.calls, connector.calls)
		}
		run, err := store.Get(context.Background(), receipt.RunID)
		if err != nil || run.Lifecycle != domain.RunSucceeded {
			t.Fatalf("lifecycle: %+v %v", run, err)
		}
		for _, name := range []domain.RunArtifactName{domain.ArtifactConnectorCapabilities, domain.ArtifactMarketPlan, domain.ArtifactMarketSnapshot, domain.ArtifactMarketCandidates, domain.ArtifactMarketGaps} {
			if len(store.artifacts[receipt.RunID][name]) == 0 {
				t.Fatalf("missing %s", name)
			}
		}
		if alter {
			var document struct {
				Entries []domain.ProductCandidate `json:"entries"`
			}
			if err := json.Unmarshal(store.artifacts[receipt.RunID][domain.ArtifactMarketCandidates], &document); err != nil {
				t.Fatal(err)
			}
			if len(document.Entries) != 0 {
				t.Fatal("runtime re-derived candidates rejected by configured investigator")
			}
		}
	}

}

type queryRecordingMarketConnector struct {
	companionConnector
	last domain.MarketDiscoveryRequest
}

func (c *queryRecordingMarketConnector) Discover(_ context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	c.last = request
	return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved}, nil
}
func TestMarketExecutorNormalizesDirectRequestsBeforeDispatch(t *testing.T) {
	connector := &queryRecordingMarketConnector{companionConnector: companionConnector{domain.DiscoveryConnectorCapabilities{ID: domain.ConnectorGiteeCatalog, Version: "v1", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"repository_identity"}, Replayable: true, PublicOnly: true}}}
	registry, err := NewMarketDiscoveryRegistry([]domain.MarketDiscoveryConnector{connector})
	if err != nil {
		t.Fatal(err)
	}
	request := domain.MarketDiscoveryStartRequest{Idea: domain.IdeaSpec{Problem: "investigate tools", Capabilities: []string{"search"}}, QueryProposals: []domain.MarketQueryProposal{{Locale: "en-US", Query: "tools", SourceText: "tools", Derivation: domain.QueryExternalAgent, Rationale: "functional query"}}, Markets: []string{"JP"}, ConnectorIDs: []string{domain.ConnectorGiteeCatalog}, Budget: domain.MarketDiscoveryBudget{MaxQueries: 1, MaxResultsPerQuery: 1, MaxTotalCandidates: 1}}
	_, err = NewMarketDiscoveryService(registry).DiscoverMarket(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if connector.last.QueryID == "" || connector.last.Query != "tools" {
		t.Fatalf("raw query not normalized: %+v", connector.last)
	}
}
