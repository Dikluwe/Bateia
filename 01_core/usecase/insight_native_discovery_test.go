/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_native_discovery_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestNativeDiscoveryContextualRoutesPreserveContextOutsideInventory(t *testing.T) {
	r := domain.InsightDiscoveryRequest{SchemaVersion: "1", QueryPolicy: "context-v2", Question: "encrypted incremental backup restore integrity"}
	p, err := PlanInsightDiscovery(r)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"encrypted incremental backup", "encrypted incremental", "encrypted integrity", "incremental backup", "backup integrity"}
	for i, q := range p.Queries {
		if q.Query != want[i] || q.Derivation != "CONTEXTUAL_TERM_CONJUNCTION_V2" || len(q.AnchorTerms) < 2 {
			t.Fatalf("context lost: %+v", q)
		}
		for _, v := range explorationWords(q.Query) {
			if !strings.Contains(r.Question, v) {
				t.Fatal("invented vocabulary")
			}
		}
	}
	if p.Request.MaxPerQuery != 5 {
		t.Fatal("new catalog breadth missing")
	}
	r.QueryPolicy = "not-supported"
	if _, err := PlanInsightDiscovery(r); err == nil {
		t.Fatal("unknown policy accepted")
	}
}

func TestNativeDiscoveryMetadataPrecedesQueryCoverageOwnerAndPopularity(t *testing.T) {
	noise := insightCatalogCandidate("new-owner", "chat")
	noise.Repository.Description = "a chatting application"
	noise.Repository.Stars = 1000000
	noise.Topics = nil
	exact := insightCatalogCandidate("same-owner", "archive")
	exact.Repository.Description = "encrypted incremental backup restore integrity"
	exact.Topics = nil
	weak := insightCatalogCandidate("third-owner", "backup")
	weak.Repository.Description = "backup tasks"
	weak.Topics = nil
	refs := mergeInsightCatalogCandidates(nil, []domain.CatalogCandidate{noise}, "q1", nil)
	refs = mergeInsightCatalogCandidates(refs, []domain.CatalogCandidate{noise}, "q2", nil)
	refs = mergeInsightCatalogCandidates(refs, []domain.CatalogCandidate{exact, weak}, "q3", nil)
	out := selectInsightDiscoveryReferences(refs, 1, "encrypted incremental backup restore integrity")
	for _, c := range out {
		if c.Selected != (c.Repository.Repo == "archive") {
			t.Fatal("diversity or stars defeated correspondence")
		}
		if c.Repository.Repo == "chat" && !containsText(c.SelectionReasons, "NO_METADATA_QUESTION_MATCH") {
			t.Fatal("zero match not explained")
		}
	}
}

type nativeResolvingReader struct {
	insightTestReader
	resolutions int
	resolved    domain.RepositoryCandidate
	resolveErr  error
}

func (r *nativeResolvingReader) ResolveRepository(_ context.Context, id domain.ProductIdentity) (domain.RepositoryCandidate, error) {
	r.resolutions++
	if r.resolveErr != nil {
		return domain.RepositoryCandidate{}, r.resolveErr
	}
	if r.resolved.Owner != "" {
		return r.resolved, nil
	}
	parts := strings.Split(id.ExternalID, "/")
	return domain.RepositoryCandidate{Owner: parts[0], Repo: parts[1], URL: id.CanonicalURL, CommitID: strings.Repeat("b", 40)}, nil
}

func TestNativeDiscoveryResolvesOnlyChosenReferenceBeforeReading(t *testing.T) {
	cat := &insightTestCatalog{}
	cat.respond = func(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
		if !r.MetadataOnly || !r.DeferRevisionResolution {
			t.Fatal("native policy not sent to port")
		}
		a := insightCatalogCandidate("org", "archive")
		a.Repository.CommitID = ""
		a.Repository.URL = "https://github.com/org/archive"
		a.Repository.Description = "encrypted incremental backup"
		b := insightCatalogCandidate("other", "unrelated")
		b.Repository.CommitID = ""
		b.Repository.URL = "https://github.com/other/unrelated"
		b.Repository.Description = "chat toolkit"
		b.Topics = nil
		return insightCatalogResult(r, a, b), nil
	}
	reader := &nativeResolvingReader{}
	request := domain.InsightDiscoveryRequest{SchemaVersion: "1", QueryPolicy: "context-v2", Question: "encrypted incremental backup", MaxReferences: 1}
	report, err := DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), request, cat, reader)
	if err != nil || reader.resolutions != 1 || reader.treeCalls != 1 || len(report.Bundles) != 1 || report.Bundles[0].Repository.Commit != strings.Repeat("b", 40) {
		t.Fatalf("lazy pin/read failed: %+v %v", report, err)
	}
	for _, bad := range []*nativeResolvingReader{
		{resolved: domain.RepositoryCandidate{Owner: "wrong", Repo: "archive", URL: "https://github.com/wrong/archive", CommitID: strings.Repeat("b", 40)}},
		{resolveErr: errors.New("private provider body")},
	} {
		report, err := DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), request, cat, bad)
		if err != nil || bad.treeCalls != 0 || len(report.Bundles) != 0 || report.Status != "PARTIAL" {
			t.Fatal("unresolved identity read or failure lost")
		}
	}
}
