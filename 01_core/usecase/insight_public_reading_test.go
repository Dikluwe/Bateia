/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_public_reading_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"strings"
	"testing"
)

type publicSiteFixtureConnector struct {
	insightSiteFixtureConnector
	identity domain.ProductIdentity
}

func (c *publicSiteFixtureConnector) Discover(_ context.Context, _ domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	product := insightSiteProduct(c.id)
	product.Identity = c.identity
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, Candidates: []domain.ProductCandidate{product}}, nil
}
func TestPublicInsightReadingNamespacesAndContextRefs(t *testing.T) {
	code := &publicSiteFixtureConnector{insightSiteFixtureConnector: insightSiteFixtureConnector{id: domain.ConnectorCodebergCatalog}, identity: domain.ProductIdentity{ConnectorID: domain.ConnectorCodebergCatalog, ExternalID: "org/repo", CanonicalURL: "https://codeberg.org/org/repo"}}
	git := &publicSiteFixtureConnector{insightSiteFixtureConnector: insightSiteFixtureConnector{id: domain.ConnectorGitLabCatalog}, identity: domain.ProductIdentity{ConnectorID: domain.ConnectorGitLabCatalog, ExternalID: "org/sub/repo", CanonicalURL: "https://gitlab.com/org/sub/repo"}}
	codeReader := &siteReadmeFixture{}
	gitReader := &siteReadmeFixture{alter: func(r *domain.RepositoryCandidate) { r.Owner = "org/sub"; r.Repo = "repo" }}
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), insightSiteRequest(), []domain.InsightSiteBinding{{Connector: code, RepositoryReader: codeReader}, {Connector: git, RepositoryReader: gitReader}})
	if err != nil || report.Status != "CONTEXT_PREPARED" || len(report.SiteContexts) != 2 || len(report.Tasks[0].ReadingCardRefs) != 2 {
		t.Fatalf("public contexts lost %v %+v", err, report)
	}
	ids := map[string]bool{}
	for _, scoped := range report.SiteContexts {
		if len(scoped.Cards) != 1 || ids[scoped.Cards[0].ID] {
			t.Fatal("card collision")
		}
		ids[scoped.Cards[0].ID] = true
		if scoped.Identity.ConnectorID == domain.ConnectorGitLabCatalog && (scoped.Bundle.Repository.Owner != "org/sub" || scoped.Cards[0].Source.Repository != "org/sub/repo") {
			t.Fatal("namespace lost or invented")
		}
	}
	for _, id := range report.Tasks[0].ReadingCardRefs {
		if !ids[id] {
			t.Fatal("dangling card ref")
		}
	}
	for _, bad := range []string{"https://gitlab.com/org//repo", "https://gitlab.com/org/../repo", "https://gitlab.com/org/sub/repo?token=x", "https://gitlab.com/org/sub/%72epo", "https://gitlab.com:443/org/sub/repo"} {
		identity := git.identity
		identity.CanonicalURL = bad
		if validInsightSiteRepository(identity) {
			t.Fatal("bad identity admitted", bad)
		}
	}
}
func TestPublicInsightReadingPreservesEarlierCardIDsAndRejectsTampering(t *testing.T) {
	bindings := siteReadmeBindings()
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), insightSiteRequest(), bindings)
	if err != nil {
		t.Fatal(err)
	}
	for _, scoped := range report.SiteContexts {
		explored, err := ExploreInsights(t.Context(), domain.InsightExplorationRequest{SchemaVersion: "1", Question: report.Plan.Request.Discovery.Question, FamiliarApproaches: report.Plan.Request.Discovery.FamiliarApproaches, MaxPassages: 1, Bundles: []domain.EvidenceBundle{scoped.Bundle}})
		if err != nil || len(explored.Cards) != 1 {
			t.Fatal(err)
		}
		expected := "site_card_" + insightHash(struct {
			Identity domain.ProductIdentity
			CardID   string
		}{scoped.Identity, explored.Cards[0].ID})
		if scoped.Cards[0].ID != expected {
			t.Fatal("prior card identity changed")
		}
	}
	scoped := report.SiteContexts[0]
	scoped.Identity = domain.ProductIdentity{ConnectorID: domain.ConnectorGitLabCatalog, ExternalID: "org/sub/repo", CanonicalURL: "https://gitlab.com/org/sub/repo"}
	scoped.Bundle.Repository.Owner = "org/sub"
	scoped.Bundle.Items[0].Ref = scoped.Bundle.Repository
	scoped.Bundle.Items[0].Content = "tampered without updating hash"
	report.SiteContexts = []domain.InsightSiteContext{scoped}
	if err := prepareInsightSiteContexts(&report); err == nil {
		t.Fatal("tampered nested context admitted")
	}
	identity := domain.ProductIdentity{ConnectorID: domain.ConnectorCodebergCatalog, ExternalID: "org/sub/repo", CanonicalURL: "https://codeberg.org/org/sub/repo"}
	if validInsightSiteRepository(identity) {
		t.Fatal("subgroup admitted for Codeberg")
	}
	identity = domain.ProductIdentity{ConnectorID: domain.ConnectorGitLabCatalog, ExternalID: "org/" + strings.Repeat("a", 101) + "/repo", CanonicalURL: "https://gitlab.com/org/" + strings.Repeat("a", 101) + "/repo"}
	if validInsightSiteRepository(identity) {
		t.Fatal("oversized namespace admitted")
	}
}
