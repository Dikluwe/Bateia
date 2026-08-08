/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/client.md
 * @layer L3
 * @updated 2026-08-08
 */

package sourcegraph_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bateia/01_core/domain"
	"bateia/03_infra/sourcegraph"
)

// mockSGResponse gera uma resposta GraphQL do Sourcegraph simulada.
func mockSGResponse(results []map[string]interface{}) []byte {
	resp := map[string]interface{}{
		"data": map[string]interface{}{
			"search": map[string]interface{}{
				"results": map[string]interface{}{
					"matchCount": len(results),
					"results":    results,
				},
			},
		},
	}
	data, _ := json.Marshal(resp)
	return data
}

func TestClient_SearchCandidates_RepositoryMatch(t *testing.T) {
	results := []map[string]interface{}{
		{
			"__typename":  "Repository",
			"name":        "github.com/gofiber/fiber",
			"url":         "/github.com/gofiber/fiber",
			"description": "Express inspired web framework",
		},
		{
			"__typename":  "Repository",
			"name":        "github.com/gin-gonic/gin",
			"url":         "/github.com/gin-gonic/gin",
			"description": "HTTP web framework in Go",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockSGResponse(results))
	}))
	defer server.Close()

	client := sourcegraph.NewClientWithEndpoint("test-token", server.URL)
	candidates, err := client.SearchCandidates(context.Background(), domain.QueryFilter{
		Query:      "web framework",
		Language:   "go",
		MaxResults: 5,
	})
	if err != nil {
		t.Fatalf("SearchCandidates falhou: %v", err)
	}

	if len(candidates) != 2 {
		t.Fatalf("Esperados 2 candidatos, recebeu: %d", len(candidates))
	}
	if candidates[0].Owner != "gofiber" || candidates[0].Repo != "fiber" {
		t.Errorf("Primeiro candidato inesperado: %s/%s", candidates[0].Owner, candidates[0].Repo)
	}
	if candidates[1].Owner != "gin-gonic" || candidates[1].Repo != "gin" {
		t.Errorf("Segundo candidato inesperado: %s/%s", candidates[1].Owner, candidates[1].Repo)
	}
}

func TestClient_SearchCandidates_FileMatch(t *testing.T) {
	results := []map[string]interface{}{
		{
			"__typename": "FileMatch",
			"repository": map[string]interface{}{
				"name":        "github.com/stretchr/testify",
				"url":         "/github.com/stretchr/testify",
				"description": "Testing toolkit",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockSGResponse(results))
	}))
	defer server.Close()

	client := sourcegraph.NewClientWithEndpoint("", server.URL)
	candidates, err := client.SearchCandidates(context.Background(), domain.QueryFilter{
		Query:      "assert testing",
		Language:   "go",
		MaxResults: 10,
	})
	if err != nil {
		t.Fatalf("SearchCandidates falhou: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("Esperado 1 candidato, recebeu: %d", len(candidates))
	}
	if candidates[0].Repo != "testify" {
		t.Errorf("Repo inesperado: %s", candidates[0].Repo)
	}
}

func TestClient_SearchCandidates_Deduplication(t *testing.T) {
	results := []map[string]interface{}{
		{"__typename": "Repository", "name": "github.com/test/dup", "url": "/test/dup", "description": ""},
		{"__typename": "FileMatch", "repository": map[string]interface{}{"name": "github.com/test/dup", "url": "/test/dup", "description": ""}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockSGResponse(results))
	}))
	defer server.Close()

	client := sourcegraph.NewClientWithEndpoint("", server.URL)
	candidates, _ := client.SearchCandidates(context.Background(), domain.QueryFilter{
		Query:      "dup test",
		Language:   "go",
		MaxResults: 10,
	})
	if len(candidates) != 1 {
		t.Errorf("Esperado 1 candidato (dedup), recebeu: %d", len(candidates))
	}
}

func TestClient_SearchCandidates_EmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockSGResponse([]map[string]interface{}{}))
	}))
	defer server.Close()

	client := sourcegraph.NewClientWithEndpoint("", server.URL)
	candidates, err := client.SearchCandidates(context.Background(), domain.QueryFilter{
		Query:      "nonexistent obscure thing",
		Language:   "cobol",
		MaxResults: 5,
	})
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if len(candidates) != 0 {
		t.Errorf("Esperado 0 candidatos, recebeu: %d", len(candidates))
	}
}

func TestClient_SearchCandidates_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := sourcegraph.NewClientWithEndpoint("", server.URL)
	candidates, _ := client.SearchCandidates(context.Background(), domain.QueryFilter{
		Query:      "test",
		Language:   "go",
		MaxResults: 5,
	})
	// Sem panic — deve retornar resultado vazio pelo fallback
	if len(candidates) != 0 {
		t.Errorf("Esperado 0 candidatos em erro de servidor, recebeu: %d", len(candidates))
	}
}
