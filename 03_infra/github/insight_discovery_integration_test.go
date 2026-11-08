/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/insight_discovery_integration_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package github

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type insightHTTPRuntime struct{}

func (insightHTTPRuntime) Now() time.Time { return time.Now() }
func (insightHTTPRuntime) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}

func TestInsightDiscoveryGitHubHTTPQueriesPinnedReadmeAndPartialQuota(t *testing.T) {
	const commit = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	queries := []string{}
	pinnedReads := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encode := func(value any) {
			if err := json.NewEncoder(w).Encode(value); err != nil {
				t.Error(err)
			}
		}
		if r.URL.Path == "/search/repositories" {
			query := r.URL.Query().Get("q")
			queries = append(queries, query)
			if strings.Contains(query, "failure") {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			if strings.Contains(query, "tradeoffs") {
				encode(map[string]any{"items": []any{}})
				return
			}
			owner, name := "org", "queue"
			if strings.Contains(query, "alternative") {
				owner, name = "other", "channel"
			}
			encode(map[string]any{"items": []any{map[string]any{"name": name, "owner": map[string]any{"login": owner}, "html_url": "https://github.com/" + owner + "/" + name, "description": "backpressure alternative", "default_branch": "main"}}})
			return
		}
		if strings.HasSuffix(r.URL.Path, "/commits/main") {
			encode(map[string]any{"sha": commit})
			return
		}
		if strings.Contains(r.URL.Path, "/git/trees/") {
			if !strings.Contains(r.URL.Path, commit) {
				t.Error("tree not pinned")
			}
			encode(map[string]any{"tree": []any{map[string]any{"path": "README.md", "mode": "100644", "type": "blob", "size": 80}}})
			return
		}
		if strings.HasSuffix(r.URL.Path, "/contents/README.md") {
			if r.URL.Query().Get("ref") != commit {
				t.Error("README not pinned")
			}
			pinnedReads++
			content := "# Reference\nA bounded backpressure mechanism from " + r.URL.Path
			encode(map[string]any{"type": "file", "path": "README.md", "size": len(content), "encoding": "base64", "content": base64.StdEncoding.EncodeToString([]byte(content))})
			return
		}
		if r.URL.Path == "/repos/org/queue" || r.URL.Path == "/repos/other/channel" {
			encode(map[string]any{"default_branch": "main", "description": "backpressure alternative"})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()
	client := NewClientWithBaseURL("", server.URL)
	report, err := usecase.DiscoverInsightReferences(domain.WithExecutionRuntime(t.Context(), insightHTTPRuntime{}), domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "backpressure queue cancellation", MaxPerQuery: 1, MaxReferences: 2}, client, client)
	if err != nil {
		t.Fatal(err)
	}
	if len(queries) != 5 || pinnedReads != 2 || len(report.Bundles) != 2 || len(report.Exploration.Cards) != 2 || report.Status != "PARTIAL" || report.Attempts[4].Status != "RATE_LIMITED" {
		t.Fatalf("HTTP workflow incomplete: %+v", report)
	}
	for _, bundle := range report.Bundles {
		if bundle.Repository.Commit != commit || bundle.Complete {
			t.Fatal("revision or scope promoted")
		}
	}
}
