/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_candidates_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestInsightSiteCandidatesPreserveNamespaceAndSealOrigins(t *testing.T) {
	values := []domain.InsightSiteCandidate{}
	for _, id := range []string{domain.ConnectorGitHubCatalog, domain.ConnectorGiteeCatalog, domain.ConnectorAtomGitCatalog} {
		product := insightSiteProduct(id)
		product.Origins = []domain.MarketCandidateOrigin{{ConnectorID: "forged", QueryID: "forged"}}
		cap := (&insightSiteFixtureConnector{id: id}).Capabilities()
		values = mergeInsightSiteCandidates(values, []domain.ProductCandidate{product}, cap, "q1", "en-US", "sha256:"+strings.Repeat("a", 64), nil)
	}
	if len(values) != 3 || values[0].ID == values[1].ID {
		t.Fatal("same owner/repo merged across hosts")
	}
	for _, candidate := range values {
		if candidate.Product.Origins[0].ConnectorID != candidate.Product.Identity.ConnectorID || candidate.Product.Origins[0].QueryID != "q1" || candidate.Product.RepositoryLinks[0].Status != domain.RepositoryLinkUnverified {
			t.Fatal("origins/links promoted")
		}
	}
	before, _ := json.Marshal(values)
	selected := selectInsightSiteCandidates(values, 3)
	after, _ := json.Marshal(values)
	if string(before) != string(after) {
		t.Fatal("selection mutated caller")
	}
	for _, c := range selected {
		if !c.Selected || c.SelectionOrder == 0 {
			t.Fatal("site diversity not selected")
		}
	}
}
func TestInsightSiteCandidatesCanonicalMetadataAndBounds(t *testing.T) {
	id := domain.ConnectorGiteeCatalog
	cap := (&insightSiteFixtureConnector{id: id}).Capabilities()
	a := insightSiteProduct(id)
	b := a
	b.Summary = "alternative"
	b.ScholarlyEvidence = &domain.ScholarlyEvidence{Catalogue: "declared", Assessment: "INDEXED"}
	forward := mergeInsightSiteCandidates(nil, []domain.ProductCandidate{a, b}, cap, "q1", "en-US", "", nil)
	backward := mergeInsightSiteCandidates(nil, []domain.ProductCandidate{b, a}, cap, "q1", "en-US", "", nil)
	if !reflect.DeepEqual(forward, backward) {
		t.Fatal("payload order changes metadata")
	}
	merged := mergeInsightSiteCandidates(forward, []domain.ProductCandidate{a}, cap, "q2", "en-US", "", nil)
	if len(merged) != 1 || len(merged[0].QueryIDs) != 2 || len(forward[0].QueryIDs) != 1 {
		t.Fatal("duplicate origins lost or mutated")
	}
	a.Identity.ConnectorID = "foreign"
	if validInsightSiteCandidates([]domain.ProductCandidate{a}, id, 1) {
		t.Fatal("foreign origin admitted")
	}
	a = insightSiteProduct(id)
	a.Summary = strings.Repeat("a", 5000)
	if validInsightSiteCandidates([]domain.ProductCandidate{a}, id, 1) {
		t.Fatal("oversized metadata admitted")
	}
}
