/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_reading_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"strings"
	"testing"
)

type siteReadmeFixture struct {
	insightTestReader
	resolves int
	err      error
	alter    func(*domain.RepositoryCandidate)
	hook     func()
}

func (r *siteReadmeFixture) ResolveRepository(_ context.Context, id domain.ProductIdentity) (domain.RepositoryCandidate, error) {
	r.resolves++
	if r.hook != nil {
		r.hook()
	}
	parts := strings.Split(id.ExternalID, "/")
	repo := domain.RepositoryCandidate{Owner: parts[0], Repo: parts[1], URL: id.CanonicalURL, CommitID: strings.Repeat("a", 40)}
	if r.alter != nil {
		r.alter(&repo)
	}
	return repo, r.err
}
func siteReadmeBindings() []domain.InsightSiteBinding {
	return []domain.InsightSiteBinding{
		{Connector: &insightSiteFixtureConnector{id: domain.ConnectorGitHubCatalog}, Reader: &insightTestReader{}},
		{Connector: &insightSiteFixtureConnector{id: domain.ConnectorGiteeCatalog}, RepositoryReader: &siteReadmeFixture{}},
		{Connector: &insightSiteFixtureConnector{id: domain.ConnectorAtomGitCatalog}, RepositoryReader: &siteReadmeFixture{}},
	}
}
func TestSiteReadmeContextsSeparateSameRepositoryCommitAndBytes(t *testing.T) {
	request := insightSiteRequest()
	bindings := siteReadmeBindings()
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings)
	if err != nil || report.Status != "CONTEXT_PREPARED" || len(report.SiteContexts) != 2 || len(report.Bundles) != 1 || len(report.Exploration.Cards) != 1 {
		t.Fatalf("reading lost: %v %+v", err, report)
	}
	ids := map[string]bool{report.Exploration.Cards[0].ID: true}
	for _, scoped := range report.SiteContexts {
		if scoped.Bundle.Repository != report.Bundles[0].Repository || scoped.Bundle.Items[0].Content != report.Bundles[0].Items[0].Content || len(scoped.Cards) != 1 {
			t.Fatal("fixture not identical or context absent")
		}
		card := scoped.Cards[0]
		if ids[card.ID] || card.Reference.ContentHash != scoped.Bundle.Items[0].ContentHash || scoped.CandidateID == "" || scoped.Identity.ConnectorID == domain.ConnectorGitHubCatalog {
			t.Fatal("cross-site collision or unbound bytes")
		}
		ids[card.ID] = true
	}
	for _, task := range report.Tasks {
		if len(task.ReadingCardRefs) != 3 || len(task.CatalogReferenceIDs) != 3 {
			t.Fatal("new reading refs not available to AI")
		}
		for _, id := range task.ReadingCardRefs {
			if !ids[id] {
				t.Fatal("dangling card", id)
			}
		}
	}
	if report.NoveltyState != "NOT_ASSESSED" || report.EffectState != "UNVERIFIED" {
		t.Fatal("operation promoted hypothesis")
	}
	request.Discovery.MaxPassages = 1
	report, err = DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings)
	if err != nil || len(report.Tasks[0].ReadingCardRefs) != 1 || len(report.SiteContexts[0].Cards) != 0 || len(report.SiteContexts) != 2 {
		t.Fatal("card budget exceeded or bytes discarded")
	}
}
func TestSiteReadmeRejectsCrossSiteReadersAndUnresolvedIdentity(t *testing.T) {
	for _, mode := range []string{"wrong-site", "both"} {
		bindings := siteReadmeBindings()
		if mode == "wrong-site" {
			bindings[0].RepositoryReader = &siteReadmeFixture{}
		} else {
			bindings[1].Reader = &insightTestReader{}
		}
		if _, err := PlanInsightSites(insightSiteRequest(), bindings); err == nil {
			t.Fatal("invalid reader admitted")
		}
	}
	for _, mode := range []string{"url", "owner", "commit", "auth", "raw-error", "invalid-url"} {
		reader := &siteReadmeFixture{}
		candidate := domain.InsightSiteCandidate{ID: "candidate", Product: insightSiteProduct(domain.ConnectorGiteeCatalog)}
		switch mode {
		case "url":
			reader.alter = func(r *domain.RepositoryCandidate) { r.URL = "https://github.com/org/repo" }
		case "owner":
			reader.alter = func(r *domain.RepositoryCandidate) { r.Owner = "other" }
		case "commit":
			reader.alter = func(r *domain.RepositoryCandidate) { r.CommitID = "main" }
		case "auth":
			reader.err = domain.InsightReadmeFailure{Code: "AUTH_REQUIRED"}
		case "raw-error":
			reader.err = errors.New("token-secret-provider-message")
		case "invalid-url":
			candidate.Product.Identity.CanonicalURL = "https://gitee.com/org/repo?token=secret"
		}
		receipt, bundle := collectInsightSiteReadme(t.Context(), candidate, reader)
		if bundle != nil || receipt.Status == "COLLECTED" || reader.treeCalls != 0 || strings.Contains(receipt.Status, "secret") {
			t.Fatal("untrusted resolution admitted", mode)
		}
		if mode == "invalid-url" && reader.resolves != 0 {
			t.Fatal("bad identity reached reader")
		}
		if mode == "auth" && receipt.Status != "AUTH_REQUIRED" {
			t.Fatal("auth hidden")
		}
	}
}

type siteReadmeManyConnector struct{ insightSiteFixtureConnector }

func (p *siteReadmeManyConnector) Discover(_ context.Context, r domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	p.calls = append(p.calls, r)
	v := insightSiteProduct(p.id)
	v.Identity.ExternalID = "org/repo" + strings.Repeat("x", len(p.calls))
	v.Identity.CanonicalURL = "https://gitee.com/" + v.Identity.ExternalID
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, Candidates: []domain.ProductCandidate{v}}, nil
}
func TestSiteReadmeGlobalByteBudgetAndPartialCancellation(t *testing.T) {
	reader := &siteReadmeFixture{insightTestReader: insightTestReader{content: strings.Repeat("a", 30<<10)}}
	connector := &siteReadmeManyConnector{insightSiteFixtureConnector: insightSiteFixtureConnector{id: domain.ConnectorGiteeCatalog}}
	request := insightSiteRequest()
	request.Discovery.MaxReferences = 8
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, []domain.InsightSiteBinding{{Connector: connector, RepositoryReader: reader}})
	if err != nil || len(report.SiteContexts) != 4 || report.Status != "PARTIAL" || report.Readmes[4].Status != "CONTENT_BUDGET" {
		t.Fatalf("byte cap failed %v %+v", err, report.Readmes)
	}
	ctx, cancel := context.WithCancel(t.Context())
	bindings := siteReadmeBindings()
	bindings[2].RepositoryReader = &siteReadmeFixture{hook: cancel}
	report, err = DiscoverInsightSites(domain.WithExecutionRuntime(ctx, &insightTestRuntime{}), insightSiteRequest(), bindings)
	if !errors.Is(err, context.Canceled) || report.Status != "CANCELLED" || len(report.SiteContexts) != 0 || report.Readmes[0].Status != "CANCELLED" || report.Readmes[1].Status != "NOT_STARTED" {
		t.Fatal("cancellation not preserved", err, report.Readmes)
	}
}
