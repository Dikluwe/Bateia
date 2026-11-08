/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_sites_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type insightSiteFixtureConnector struct {
	id          string
	calls       []domain.MarketDiscoveryRequest
	auth, empty bool
	hook        func()
	unique      bool
}

func (p *insightSiteFixtureConnector) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: p.id, Version: "fixture-v1", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, Markets: []string{}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"catalog_metadata"}, PublicOnly: true, Replayable: true}
}
func (p *insightSiteFixtureConnector) Discover(_ context.Context, r domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	p.calls = append(p.calls, r)
	if p.hook != nil {
		p.hook()
	}
	if p.auth {
		return domain.ConnectorDiscovery{Status: domain.AttemptUnknownProvider, ErrorClass: "AUTH_REQUIRED"}, nil
	}
	values := []domain.ProductCandidate{}
	if !p.empty {
		count := 1
		if p.unique {
			count = r.Limit
		}
		for i := 0; i < count; i++ {
			v := insightSiteProduct(p.id)
			if p.unique {
				v.Identity.ExternalID = fmt.Sprintf("org/repo%d-%d", len(p.calls), i)
			}
			values = append(values, v)
		}
	}
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, Candidates: values, ResultCount: len(values)}, nil
}
func insightSiteProduct(id string) domain.ProductCandidate {
	host := "gitee.com"
	if id == domain.ConnectorGitHubCatalog {
		host = "github.com"
	}
	if id == domain.ConnectorAtomGitCatalog {
		host = "atomgit.com"
	}
	return domain.ProductCandidate{Identity: domain.ProductIdentity{ConnectorID: id, ExternalID: "org/repo", CanonicalURL: "https://" + host + "/org/repo"}, Name: "org/repo", Summary: "queue reference", Publisher: "org", ProductKind: "repository", RepositoryLinks: []domain.ProductRepositoryLink{{OwnerRepo: "org/repo", Commit: strings.Repeat("a", 40), Status: domain.RepositoryLinkUnverified}}}
}
func insightSiteRequest() domain.InsightSiteRequest {
	return domain.InsightSiteRequest{Discovery: insightDiscoveryFixture(), Sources: []string{"all"}}
}

func TestInsightSitesGlobalBudgetFairnessAuthAndContentIsolation(t *testing.T) {
	a, g, h := &insightSiteFixtureConnector{id: domain.ConnectorAtomGitCatalog, auth: true}, &insightSiteFixtureConnector{id: domain.ConnectorGiteeCatalog}, &insightSiteFixtureConnector{id: domain.ConnectorGitHubCatalog}
	reader := &insightTestReader{}
	r := insightSiteRequest()
	r.MaxCalls = 4
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), r, []domain.InsightSiteBinding{{Connector: h, Reader: reader}, {Connector: g}, {Connector: a}})
	if err != nil || report.Calls != 4 || len(a.calls) != 1 || len(g.calls) != 2 || len(h.calls) != 1 || len(report.Attempts) != 15 || len(report.Candidates) != 2 || report.Status != "PARTIAL" {
		t.Fatalf("fairness/budget/auth lost: %v %+v", err, report)
	}
	if report.Attempts[3].Status != "AUTH_BLOCKED" || reader.treeCalls != 1 || len(report.Bundles) != 1 || len(report.Exploration.Cards) != 1 {
		t.Fatal("auth retried or non-GitHub metadata read via GitHub")
	}
	seenCatalogOnly := false
	for _, receipt := range report.Readmes {
		if receipt.ConnectorID == domain.ConnectorGiteeCatalog {
			seenCatalogOnly = receipt.Status == "CATALOG_ONLY"
		}
	}
	if !seenCatalogOnly || report.NoveltyState != "NOT_ASSESSED" || report.EffectState != "UNVERIFIED" {
		t.Fatal("catalog promoted")
	}
	for _, task := range report.Tasks {
		if len(task.CatalogReferenceIDs) != 2 || len(task.ReadingCardRefs) != 1 {
			t.Fatal("catalog/reading refs missing or mixed")
		}
	}
}

func TestInsightSitesCatalogOnlyEmptyAndCancellation(t *testing.T) {
	g := &insightSiteFixtureConnector{id: domain.ConnectorGiteeCatalog}
	ctx := domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{})
	report, err := DiscoverInsightSites(ctx, insightSiteRequest(), []domain.InsightSiteBinding{{Connector: g}})
	if err != nil || report.Status != "CATALOG_REFERENCES" || len(report.Bundles) != 0 || len(report.Tasks[0].CatalogReferenceIDs) != 1 || len(report.Tasks[0].ReadingCardRefs) != 0 {
		t.Fatal("catalog-only result hidden or fabricated reading")
	}
	g.empty = true
	report, err = DiscoverInsightSites(ctx, insightSiteRequest(), []domain.InsightSiteBinding{{Connector: g}})
	if err != nil || report.Status != "NO_CANDIDATES" {
		t.Fatal("empty is not a valid scoped result")
	}
	parent, cancel := context.WithCancel(context.Background())
	g = &insightSiteFixtureConnector{id: domain.ConnectorGiteeCatalog}
	h := &insightSiteFixtureConnector{id: domain.ConnectorGitHubCatalog, hook: cancel}
	report, err = DiscoverInsightSites(domain.WithExecutionRuntime(parent, &insightTestRuntime{}), insightSiteRequest(), []domain.InsightSiteBinding{{Connector: h}, {Connector: g}})
	if !errors.Is(err, context.Canceled) || report.Status != "CANCELLED" || report.Calls != 2 || len(report.Candidates) != 1 || report.Attempts[2].Status != "NOT_STARTED" {
		t.Fatal("cancelled work continued or discarded prior data")
	}
}

func TestInsightSitesCandidateBudgetIsGlobal(t *testing.T) {
	bindings := []domain.InsightSiteBinding{}
	for _, id := range []string{domain.ConnectorGiteeCatalog, domain.ConnectorAtomGitCatalog, domain.ConnectorCodebergCatalog} {
		bindings = append(bindings, domain.InsightSiteBinding{Connector: &insightSiteFixtureConnector{id: id, unique: true}})
	}
	r := insightSiteRequest()
	r.MaxCalls = 60
	r.Discovery.MaxPerQuery = 5
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), r, bindings)
	if err != nil || len(report.Candidates) != 64 || report.Calls != 13 || report.Attempts[14].Status != "CANDIDATE_BUDGET" {
		t.Fatal("candidate cap exceeded or budget omission not recorded")
	}
	encoded, _ := json.Marshal(report)
	if strings.Contains(string(encoded), "SUPPORTED") {
		t.Fatal("operational selection promoted hypotheses")
	}
}

type insightMarketSiteFixture struct{ *insightSiteFixtureConnector }

func (p insightMarketSiteFixture) Capabilities() domain.DiscoveryConnectorCapabilities {
	cap := p.insightSiteFixtureConnector.Capabilities()
	cap.Kind = domain.ConnectorApplicationCatalog
	cap.SupportsMarketFilter = true
	cap.Markets = []string{"US"}
	cap.QueryLocales = []string{"en-US", "pt-BR"}
	return cap
}
func TestInsightSitesMarketFilterUsesSupportedLocaleCountry(t *testing.T) {
	connector := &insightSiteFixtureConnector{id: domain.ConnectorAppleItunes}
	bindings := []domain.InsightSiteBinding{{Connector: insightMarketSiteFixture{connector}}}
	ctx := domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{})
	report, err := DiscoverInsightSites(ctx, insightSiteRequest(), bindings)
	if err != nil || report.Calls != 5 || connector.calls[0].Market != "US" || report.Attempts[0].Market != "US" || report.Candidates[0].Product.Origins[0].ObservedMarket != "US" {
		t.Fatal("market omitted or origin lost")
	}
	r := insightSiteRequest()
	r.Locale = "pt-BR"
	report, err = DiscoverInsightSites(ctx, r, bindings)
	if !errors.Is(err, domain.ErrInsightDiscoveryUnavailable) || report.Calls != 0 || report.Attempts[0].Status != "UNSUPPORTED_MARKET" || len(connector.calls) != 5 {
		t.Fatal("unsupported market called provider or became emptiness")
	}
}
