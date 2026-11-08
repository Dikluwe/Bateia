/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/insight_sites_integration_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package forge

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type insightForgeTransport func(*http.Request) (*http.Response, error)

func (f insightForgeTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type insightForgeRuntime struct{}

func (insightForgeRuntime) Now() time.Time { return time.Now() }
func (insightForgeRuntime) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}
func TestInsightSitesRealGiteeAtomGitConnectorsKeepAuthAndHostsSeparate(t *testing.T) {
	counts := map[string]int{}
	giteeToken, atomToken := "fixture-gitee-token", "fixture-atom-token"
	client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
		counts[r.URL.Host]++
		host := "gitee.com"
		if r.URL.Host == "gitee.com" {
			if r.URL.Query().Get("access_token") != giteeToken || r.Header.Get("Authorization") != "" {
				t.Error("Gitee credentials crossed boundary")
			}
		} else if r.URL.Host == "api.atomgit.com" {
			host = "atomgit.com"
			if r.Header.Get("Authorization") != "Bearer "+atomToken || r.URL.Query().Has("access_token") {
				t.Error("AtomGit credentials crossed boundary")
			}
		} else {
			t.Fatal("unexpected host", r.URL.Host)
		}
		if r.URL.Path != "/api/v5/search/repositories" {
			t.Error("unexpected endpoint", r.URL.Path)
		}
		body := `[{"full_name":"org/repo","html_url":"https://` + host + `/org/repo","description":"queue alternative","private":false,"public":true}]`
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})}
	request := domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "queue cancellation"}, Sources: []string{"gitee", "atomgit"}, MaxCalls: 10}
	bindings := []domain.InsightSiteBinding{{Connector: NewGiteeConnectorWithToken(giteeToken, client)}, {Connector: NewAtomGitConnector(atomToken, client)}}
	report, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), insightForgeRuntime{}), request, bindings)
	if err != nil || len(report.Candidates) != 2 || report.Calls != 10 || counts["gitee.com"] != 5 || counts["api.atomgit.com"] != 5 || len(report.Bundles) != 0 || len(report.Tasks[0].CatalogReferenceIDs) != 2 {
		t.Fatalf("multisite HTTP workflow failed: %v %+v", err, report)
	}
	if report.Candidates[0].ID == report.Candidates[1].ID {
		t.Fatal("cross-host owner/repo merged")
	}
	counts = map[string]int{}
	bindings[1].Connector = NewAtomGitConnector("", client)
	report, err = usecase.DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), insightForgeRuntime{}), request, bindings)
	if err != nil || counts["api.atomgit.com"] != 0 || counts["gitee.com"] != 5 || report.Status != "PARTIAL" || report.Attempts[0].Status != "AUTH_REQUIRED" || report.Attempts[2].Status != "AUTH_BLOCKED" {
		t.Fatal("missing auth was retried or hid Gitee results")
	}
}
