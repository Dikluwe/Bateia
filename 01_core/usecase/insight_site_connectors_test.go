/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_connectors_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"testing"
)

type insightSiteCodeFixture struct {
	request domain.DiscoveryRequest
	bad     bool
}

func (p *insightSiteCodeFixture) Discover(_ context.Context, r domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	p.request = r
	provider := "sourcegraph"
	if p.bad {
		provider = "foreign"
	}
	return domain.ProviderDiscovery{Provider: provider, Attempts: []domain.ProviderAttempt{{Status: "SUCCEEDED"}}, Candidates: []domain.RepositoryCandidate{insightCatalogCandidate("org", "repo").Repository}}, nil
}
func TestInsightSiteBridgesPreserveMetadataAndBudgets(t *testing.T) {
	r := domain.MarketDiscoveryRequest{QueryID: "q", Query: "queue", QueryLocale: "en-US", Limit: 1}
	catalog := &insightTestCatalog{}
	bridge := NewInsightGitHubConnector(catalog)
	got, err := bridge.Discover(context.Background(), r)
	if err != nil || len(got.Candidates) != 1 || got.Candidates[0].Identity.ConnectorID != domain.ConnectorGitHubCatalog || got.Candidates[0].RepositoryLinks[0].Status != domain.RepositoryLinkUnverified || catalog.calls[0].MaxResults != 1 {
		t.Fatal("GitHub bridge loses scope/budget")
	}
	code := &insightSiteCodeFixture{}
	got, err = NewInsightSourcegraphConnector(code).Discover(context.Background(), r)
	if err != nil || len(got.Candidates) != 1 || code.request.MaxAttempts != 1 || code.request.MaxResults != 1 || got.Candidates[0].Identity.ConnectorID != domain.ConnectorSourcegraph {
		t.Fatal("Sourcegraph bridge loses provenance/budget")
	}
	code.bad = true
	got, _ = NewInsightSourcegraphConnector(code).Discover(context.Background(), r)
	if got.Status != domain.AttemptIntegrityError {
		t.Fatal("foreign code provider admitted")
	}
	catalog.respond = func(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
		out := insightCatalogResult(r)
		out.Attempt.Status = "RATE_LIMITED"
		return out, errors.New("private text")
	}
	got, _ = bridge.Discover(context.Background(), r)
	if got.ErrorClass != "RATE_LIMIT" || len(got.Candidates) != 0 {
		t.Fatal("quota lost or unsafe payload admitted")
	}
}
