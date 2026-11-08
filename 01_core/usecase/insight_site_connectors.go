/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_connectors.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
)

type insightRepositoryConnector struct {
	catalog domain.ProjectCatalog
	code    domain.ObservableDiscoveryProvider
	id      string
}

func NewInsightGitHubConnector(catalog domain.ProjectCatalog) domain.MarketDiscoveryConnector {
	return insightRepositoryConnector{catalog: catalog, id: domain.ConnectorGitHubCatalog}
}
func NewInsightSourcegraphConnector(provider domain.ObservableDiscoveryProvider) domain.MarketDiscoveryConnector {
	return insightRepositoryConnector{code: provider, id: domain.ConnectorSourcegraph}
}
func (c insightRepositoryConnector) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: c.id, Version: "insight-repository-bridge-v1", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US", "pt-BR"}, Markets: []string{}, IdentityLevels: []string{"owner_repo", "commit"}, EvidenceSurfaces: []string{"catalog_metadata"}, Replayable: true, PublicOnly: true}
}
func (c insightRepositoryConnector) Discover(ctx context.Context, r domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	out := domain.ConnectorDiscovery{Candidates: []domain.ProductCandidate{}, Status: domain.AttemptUnknownProvider, ErrorClass: "PROVIDER"}
	values := []domain.CatalogCandidate{}
	if c.catalog != nil {
		result, err := c.catalog.SearchCatalog(ctx, domain.CatalogSearchRequest{MetadataOnly: r.QueryPolicy == "context-v2", DeferRevisionResolution: r.QueryPolicy == "context-v2", ExpansionID: r.QueryID, Query: r.Query, MaxResults: r.Limit})
		if err != nil {
			if result.Attempt.Status == "RATE_LIMITED" {
				out.ErrorClass = "RATE_LIMIT"
			}
			return out, nil
		}
		if result.Role != domain.RoleProjectCatalog || result.Provider != c.catalog.Capabilities().Provider || result.Attempt.Query != r.Query || result.Attempt.ExpansionID != r.QueryID || result.Attempt.Status != "SUCCEEDED" {
			out.Status = domain.AttemptIntegrityError
			out.ErrorClass = "INTEGRITY"
			return out, nil
		}
		values = result.Candidates
	} else if c.code != nil {
		result, err := c.code.Discover(ctx, domain.DiscoveryRequest{ExpansionID: r.QueryID, Query: r.Query, MaxResults: r.Limit, MaxAttempts: 1})
		if err != nil {
			return out, nil
		}
		if result.Provider != "sourcegraph" || len(result.Attempts) != 1 || (result.Attempts[0].Status != "SUCCEEDED" && result.Attempts[0].Status != "EMPTY") {
			out.Status = domain.AttemptIntegrityError
			out.ErrorClass = "INTEGRITY"
			return out, nil
		}
		for _, repo := range result.Candidates {
			values = append(values, domain.CatalogCandidate{Repository: repo})
		}
	} else {
		return out, nil
	}
	if !validInsightCatalogCandidates(values, r.Limit) {
		out.Status = domain.AttemptIntegrityError
		out.ErrorClass = "INTEGRITY"
		return out, nil
	}
	for _, value := range values {
		repo := value.Repository
		identity := repo.Owner + "/" + repo.Repo
		out.Candidates = append(out.Candidates, domain.ProductCandidate{Identity: domain.ProductIdentity{ConnectorID: c.id, ExternalID: identity, CanonicalURL: repo.URL}, Name: identity, Summary: repo.Description, Publisher: repo.Owner, ProductKind: "repository", ListingLocales: []string{r.QueryLocale}, ObservedMarkets: []string{}, PublisherCountries: []string{}, RepositoryLinks: []domain.ProductRepositoryLink{{OwnerRepo: identity, Commit: repo.CommitID, Status: domain.RepositoryLinkUnverified, Source: c.id}}, Origins: []domain.MarketCandidateOrigin{}, Witnesses: []domain.IdentityWitness{}, Resolution: domain.IdentityUnknown})
	}
	out.Status = domain.AttemptSucceeded
	if len(out.Candidates) == 0 {
		out.Status = domain.AttemptEmptyObserved
	}
	out.ErrorClass = ""
	out.ResultCount = len(out.Candidates)
	out.AcceptedCount = len(out.Candidates)
	out.ResponseHash = "sha256:" + insightHash(out.Candidates)
	return out, nil
}
