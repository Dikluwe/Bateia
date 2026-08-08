/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/client.md
 * @layer L3
 * @updated 2026-08-08
 */

package github_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/03_infra/github"
)

func TestClient_FetchRepoMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/repos/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resp := map[string]interface{}{
			"stargazers_count": 1234,
			"description":     "A great project",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := github.NewClientWithBaseURL("test-token", server.URL)
	candidate := domain.RepositoryCandidate{
		Owner: "testowner",
		Repo:  "testrepo",
		Stars: 0,
	}

	result, err := client.FetchRepoMetadata(context.Background(), candidate)
	if err != nil {
		t.Fatalf("FetchRepoMetadata falhou: %v", err)
	}
	if result.Stars != 1234 {
		t.Errorf("Stars esperado 1234, recebeu: %d", result.Stars)
	}
	if result.Description != "A great project" {
		t.Errorf("Description inesperada: %s", result.Description)
	}
}

func TestClient_FetchRepoMetadata_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := github.NewClientWithBaseURL("", server.URL)
	candidate := domain.RepositoryCandidate{Owner: "x", Repo: "y"}

	result, err := client.FetchRepoMetadata(context.Background(), candidate)
	if err != nil {
		t.Fatalf("Não deveria retornar erro, mas retornou: %v", err)
	}
	// Em caso de erro HTTP, retorna candidato original sem modificação
	if result.Stars != 0 {
		t.Errorf("Stars deveria permanecer 0, recebeu: %d", result.Stars)
	}
}

func TestClient_FetchTree(t *testing.T) {
	items := []map[string]interface{}{
		{"name": "go.mod", "path": "go.mod", "type": "file", "size": 200},
		{"name": "main.go", "path": "main.go", "type": "file", "size": 500},
		{"name": "cmd", "path": "cmd", "type": "dir", "size": 0},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	client := github.NewClientWithBaseURL("", server.URL)
	candidate := domain.RepositoryCandidate{Owner: "test", Repo: "repo"}

	nodes, err := client.FetchTree(context.Background(), candidate)
	if err != nil {
		t.Fatalf("FetchTree falhou: %v", err)
	}
	if len(nodes) != 3 {
		t.Fatalf("Esperados 3 nós, recebeu: %d", len(nodes))
	}
	if nodes[0].Path != "go.mod" || nodes[0].Type != "file" {
		t.Errorf("Primeiro nó inesperado: %+v", nodes[0])
	}
}

func TestClient_FetchManifest_GoMod(t *testing.T) {
	goModContent := "module github.com/test/repo\n\ngo 1.21\n\nrequire github.com/stretchr/testify v1.8.4"
	encoded := base64.StdEncoding.EncodeToString([]byte(goModContent))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/contents/go.mod") {
			resp := map[string]interface{}{
				"name":     "go.mod",
				"path":     "go.mod",
				"type":     "file",
				"size":     len(goModContent),
				"content":  encoded,
				"encoding": "base64",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := github.NewClientWithBaseURL("", server.URL)
	candidate := domain.RepositoryCandidate{Owner: "test", Repo: "repo", CommitID: "HEAD"}

	manifest, err := client.FetchManifest(context.Background(), candidate)
	if err != nil {
		t.Fatalf("FetchManifest falhou: %v", err)
	}
	if manifest == nil {
		t.Fatal("Manifesto retornou nil")
	}
	if manifest.Language != "go" {
		t.Errorf("Language esperada 'go', recebeu: %s", manifest.Language)
	}
}

func TestClient_FetchManifest_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := github.NewClientWithBaseURL("", server.URL)
	candidate := domain.RepositoryCandidate{Owner: "empty", Repo: "repo"}

	_, err := client.FetchManifest(context.Background(), candidate)
	if err == nil {
		t.Fatal("Esperava erro MANIFEST_NOT_FOUND, mas não recebeu")
	}
}

func TestClient_AuthorizationHeader(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		resp := map[string]interface{}{
			"stargazers_count": 0,
			"description":     "",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := github.NewClientWithBaseURL("my-secret-token", server.URL)
	candidate := domain.RepositoryCandidate{Owner: "a", Repo: "b"}
	client.FetchRepoMetadata(context.Background(), candidate)

	if gotAuth != "token my-secret-token" {
		t.Errorf("Authorization header inesperado: %q", gotAuth)
	}
}
