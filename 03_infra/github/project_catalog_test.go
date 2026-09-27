/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/project_catalog_test.md
 * @layer L3
 * @updated 2026-08-26
 */
package github

import (
	"bateia/01_core/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitHubCatalogSearchReturnsResolvedRepository(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/search/repositories" {
			_ = json.NewEncoder(w).Encode(map[string]any{"items": []any{map[string]any{"name": "reuse", "html_url": "https://github.com/org/reuse", "description": "software reuse discovery", "stargazers_count": 42, "default_branch": "main", "owner": map[string]any{"login": "org"}, "topics": []string{"software-reuse"}}}})
			return
		}
		if r.URL.Path == "/repos/org/reuse" {
			_ = json.NewEncoder(w).Encode(map[string]any{"description": "software reuse discovery", "default_branch": "main"})
			return
		}
		if r.URL.Path == "/repos/org/reuse/commits/main" {
			_ = json.NewEncoder(w).Encode(map[string]any{"sha": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()
	c := NewClientWithBaseURL("", server.URL)
	got, err := c.SearchCatalog(t.Context(), domain.CatalogSearchRequest{ExpansionID: "x", Query: "software reuse discovery", Discriminators: []string{"software reuse"}, MaxResults: 3})
	if err != nil {
		t.Fatal(err)
	}
	if got.Role != domain.RoleProjectCatalog || len(got.Candidates) != 1 || got.Candidates[0].Repository.CommitID != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("invalid catalog: %+v", got)
	}
	if len(got.Candidates[0].MatchedFields) == 0 {
		t.Fatal("declared discriminator was not matched")
	}
}

func TestCatalogFieldMatchRejectsIsolatedGenericToken(t *testing.T) {
	repo := domain.RepositoryCandidate{Repo: "component", Description: "unrelated UI library"}
	if got := matchedCatalogFields([]string{"component selection"}, repo, nil); len(got) != 0 {
		t.Fatalf("generic token leaked: %v", got)
	}
}

func TestGitHubCatalogClassifiesRateLimit(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTooManyRequests) }))
	defer s.Close()
	c := NewClientWithBaseURL("", s.URL)
	got, err := c.SearchCatalog(t.Context(), domain.CatalogSearchRequest{Query: "software reuse", MaxResults: 1})
	if err == nil || got.Attempt.Status != "RATE_LIMITED" {
		t.Fatalf("rate limit hidden: %+v %v", got, err)
	}
}
