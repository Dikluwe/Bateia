/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_discovery_selection_test.md
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

func TestInsightDiscoverySelectionCoversRoutesAndOwnersIgnoringStars(t *testing.T) {
	refs := mergeInsightCatalogCandidates(nil, []domain.CatalogCandidate{insightCatalogCandidate("org", "a")}, "q1", nil)
	refs = mergeInsightCatalogCandidates(refs, []domain.CatalogCandidate{insightCatalogCandidate("org", "a")}, "q2", nil)
	popular := insightCatalogCandidate("org", "b")
	popular.Repository.Stars = 999999
	refs = mergeInsightCatalogCandidates(refs, []domain.CatalogCandidate{popular, insightCatalogCandidate("other", "c")}, "q3", []string{"queue"})
	before, _ := json.Marshal(refs)
	got := selectInsightDiscoveryReferences(refs, 2)
	selected := map[string]bool{}
	for _, ref := range got {
		if ref.Selected {
			selected[repositoryKey(ref.Repository)] = true
		}
	}
	if !selected["org/a"] || !selected["other/c"] || selected["org/b"] {
		t.Fatalf("selection follows popularity or loses coverage/owner: %+v", got)
	}
	after, _ := json.Marshal(refs)
	if string(before) != string(after) {
		t.Fatal("selection mutated inputs")
	}
}

func TestInsightDiscoveryMergeTracksAllRevisionsAndCanonicalMetadata(t *testing.T) {
	a := insightCatalogCandidate("Org", "Repo")
	b := a
	b.Repository.CommitID = ""
	b.Repository.Stars = 0
	c := a
	c.Repository.CommitID = strings.Repeat("b", 40)
	forward := mergeInsightCatalogCandidates(nil, []domain.CatalogCandidate{a, b, c}, "query", nil)
	backward := mergeInsightCatalogCandidates(nil, []domain.CatalogCandidate{c, b, a}, "query", nil)
	if !reflect.DeepEqual(forward, backward) || len(forward) != 1 || !forward[0].RevisionConflict || len(forward[0].ObservedCommits) != 2 {
		t.Fatal("opaque observation erased revision conflict or order changes metadata")
	}
	merged := mergeInsightCatalogCandidates(forward, []domain.CatalogCandidate{a}, "another", nil)
	if len(merged[0].QueryIDs) != 2 || len(forward[0].QueryIDs) != 1 {
		t.Fatal("query provenance lost or prior input mutated")
	}
	if validInsightCatalogCandidates([]domain.CatalogCandidate{a, c}, 1) {
		t.Fatal("oversized response admitted")
	}
	b.Repository.Owner = "../org"
	if validInsightCatalogCandidates([]domain.CatalogCandidate{b}, 1) {
		t.Fatal("malformed identity admitted")
	}
}
