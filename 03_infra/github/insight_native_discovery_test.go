/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/insight_native_discovery_test.md
 * @layer L3
 * @updated 2026-11-08
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
)

func TestNativeDiscoveryHTTPMetadataSearchAndSelectedPin(t *testing.T) {
	const pin = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	searches, resolutions, reads := 0, 0, 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		send := func(v any) {
			if err := json.NewEncoder(w).Encode(v); err != nil {
				t.Error(err)
			}
		}
		if r.URL.Path == "/search/repositories" {
			searches++
			q := r.URL.Query().Get("q")
			if !strings.HasSuffix(q, " in:name,description,topics") || strings.Contains(q, "readme") {
				t.Error("incidental README search still active", q)
			}
			send(map[string]any{"items": []any{
				map[string]any{"name": "archive", "html_url": "https://github.com/org/archive", "description": "encrypted incremental backup", "owner": map[string]any{"login": "org"}},
				map[string]any{"name": "chat", "html_url": "https://github.com/other/chat", "description": "unrelated chat client", "owner": map[string]any{"login": "other"}},
			}})
			return
		}
		if r.URL.Path == "/repos/org/archive" {
			if searches != 5 {
				t.Error("metadata fetched during catalog search")
			}
			resolutions++
			send(map[string]any{"full_name": "org/archive", "html_url": "https://github.com/org/archive", "private": false, "default_branch": "main"})
			return
		}
		if r.URL.Path == "/repos/org/archive/commits/main" {
			send(map[string]any{"sha": pin})
			return
		}
		if strings.Contains(r.URL.Path, "/git/trees/") {
			if !strings.HasSuffix(r.URL.Path, pin) {
				t.Error("unpinned tree")
			}
			send(map[string]any{"tree": []any{map[string]any{"path": "README.md", "mode": "100644", "type": "blob", "size": 40}}})
			return
		}
		if strings.HasSuffix(r.URL.Path, "/contents/README.md") {
			if r.URL.Query().Get("ref") != pin {
				t.Error("unpinned content")
			}
			reads++
			content := "Encrypted incremental backups with restore verification."
			send(map[string]any{"type": "file", "path": "README.md", "size": len(content), "encoding": "base64", "content": base64.StdEncoding.EncodeToString([]byte(content))})
			return
		}
		t.Error("unexpected or unselected repository request", r.URL.Path)
		http.NotFound(w, r)
	}))
	defer server.Close()
	c := NewClientWithBaseURL("", server.URL)
	report, err := usecase.DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), insightHTTPRuntime{}), domain.InsightDiscoveryRequest{SchemaVersion: "1", QueryPolicy: "context-v2", Question: "encrypted incremental backup", MaxPerQuery: 2, MaxReferences: 1}, c, c)
	if err != nil || resolutions != 1 || reads != 1 || len(report.Bundles) != 1 || report.Bundles[0].Repository.Commit != pin {
		t.Fatalf("native chain failed: %v %+v", err, report)
	}
}
