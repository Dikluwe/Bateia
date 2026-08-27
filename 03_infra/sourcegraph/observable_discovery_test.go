/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/observable_discovery_test.md
 * @layer L3
 * @updated 2026-08-24
 */

package sourcegraph

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func TestObservableQueriesPreserveFallbackOrder(t *testing.T) {
	queries := observableQueries("code impact analysis", "go", nil, "")
	want := []string{"lang:Go code impact analysis", "lang:Go code impact", "code impact"}
	if len(queries) != len(want) {
		t.Fatalf("tentativas incorretas: %+v", queries)
	}
	for index := range want {
		if queries[index] != want[index] {
			t.Fatalf("ordem incorreta: %+v", queries)
		}
	}
}

func TestObservableQueriesPreserveDeclaredDiscriminator(t *testing.T) {
	queries := observableQueries("GitHub repository recommendation", "", []string{"repository recommendation"}, "")
	want := []string{"GitHub repository recommendation", "repository recommendation"}
	if len(queries) != len(want) || queries[1] != want[1] {
		t.Fatalf("discriminator perdido no fallback: %+v", queries)
	}
}

func TestObservableQueriesKeepRepositoryScopeInEveryFallback(t *testing.T) {
	queries := observableQueries("software reuse discovery", "", []string{"software reuse"}, "org/repo")
	for _, query := range queries {
		if !strings.HasPrefix(query, "repo:^org/repo$ ") {
			t.Fatalf("scope escaped: %q", query)
		}
	}
}

func TestDiscoverClassifiesAttemptTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		_, _ = writer.Write([]byte(`{"data":{"search":{"results":{"matchCount":0,"results":[]}}}}`))
	}))
	defer server.Close()
	client := NewClientWithEndpoint("", server.URL)
	result, err := client.Discover(t.Context(), domain.DiscoveryRequest{Query: "slow", MaxResults: 1, AttemptTimeout: 20 * time.Millisecond})
	if err == nil {
		t.Fatal("timeout must fail the provider expansion")
	}
	if len(result.Attempts) != 2 || result.Attempts[0].Status != "TIMED_OUT" || result.Attempts[1].Status != "TIMED_OUT" || result.Attempts[0].ErrorClass != "TIMEOUT" {
		t.Fatalf("timeout misclassified: %#v", result.Attempts)
	}
}

func TestSafeErrorRedactsToken(t *testing.T) {
	client := NewClient("secret-token")
	message := client.safeError(errors.New("secret-token failed"))
	if strings.Contains(message, "secret-token") {
		t.Fatal("token exposto")
	}
}

func TestClientImplementsObservableDiscoveryProvider(t *testing.T) {
	var _ domain.ObservableDiscoveryProvider = NewClient("")
}

func TestMatchSnippetIsNormalizedAndTruncated(t *testing.T) {
	input := strings.Repeat("word   ", 100)
	got := truncateMatchSnippet(input)
	if len([]rune(got)) != 240 || strings.Contains(got, "  ") {
		t.Fatalf("snippet inválido: len=%d value=%q", len([]rune(got)), got)
	}
}

func TestDiscoverRecordsEmptyAttemptBeforeSuccessfulFallback(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		calls++
		results := []map[string]any{}
		if calls == 2 {
			results = append(results, map[string]any{"__typename": "Repository", "name": "github.com/org/repo", "url": "/github.com/org/repo", "description": "match"})
		}
		response := map[string]any{"data": map[string]any{"search": map[string]any{"results": map[string]any{"matchCount": len(results), "results": results}}}}
		_ = json.NewEncoder(writer).Encode(response)
	}))
	defer server.Close()
	client := NewClientWithEndpoint("", server.URL)
	result, err := client.Discover(t.Context(), domain.DiscoveryRequest{Query: "code impact analysis", Language: "go", MaxResults: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Attempts) != 2 || result.Attempts[0].ResultCount != 0 || result.Attempts[0].Error != "" || result.Attempts[1].ResultCount != 1 {
		t.Fatalf("tentativas incorretas: %+v", result.Attempts)
	}
}
