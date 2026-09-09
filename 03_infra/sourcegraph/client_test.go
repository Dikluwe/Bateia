/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/client_test.md
 * @layer L3
 * @updated 2026-09-08
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
	// A valid GraphQL empty result is an array; malformed envelopes use explicit JSON.
	if results == nil {
		results = []map[string]interface{}{}
	}
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
	candidates, err := client.SearchCandidates(context.Background(), domain.QueryFilter{
		Query:      "test",
		Language:   "go",
		MaxResults: 5,
	})
	if err == nil {
		t.Fatal("esperava erro após todas as tentativas ao Sourcegraph falharem")
	}
	if len(candidates) != 0 {
		t.Errorf("Esperado 0 candidatos em erro de servidor, recebeu: %d", len(candidates))
	}
}

// Step 0054 supersedes destructive fallback expectations; baseline hashes are in intent_preservation_test.md.
func TestSearchTraceRetainsSingleOriginalQuery(t *testing.T) {
	var queries []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Variables map[string]string `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
		}
		queries = append(queries, req.Variables["query"])
		var results []map[string]interface{}
		if len(queries) == 3 {
			results = []map[string]interface{}{{"__typename": "Repository", "name": "github.com/example/USalign"}}
		}
		w.Write(mockSGResponse(results))
	}))
	defer server.Close()
	original := "order independent transparency opengl"
	candidates, attempts, err := sourcegraph.NewClientWithEndpoint("", server.URL).SearchCandidatesWithTrace(context.Background(), domain.QueryFilter{Query: original, Language: "cpp", MaxResults: 5})
	if err != nil || len(candidates) != 0 || len(attempts) != 1 || len(queries) != 1 {
		t.Fatalf("%v %v %v", candidates, attempts, err)
	}
	for i, a := range attempts {
		if a.OriginalQuery != original || a.Query != queries[i] {
			t.Fatalf("trace mismatch %+v", a)
		}
	}
	if attempts[0].Query != "lang:C++ "+original || attempts[0].Relaxation != "" {
		t.Fatalf("unexpected relaxation %+v", attempts)
	}
}
