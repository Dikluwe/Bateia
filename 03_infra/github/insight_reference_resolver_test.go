/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/insight_reference_resolver_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package github

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type resolverTransport func(*http.Request) (*http.Response, error)

func (f resolverTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func resolverResponse(r *http.Request, status int, body string) *http.Response {
	return &http.Response{Request: r, StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}
func resolverIdentity() domain.ProductIdentity {
	return domain.ProductIdentity{ConnectorID: domain.ConnectorGitHubCatalog, ExternalID: "child/idea", CanonicalURL: "https://github.com/child/idea"}
}

type resolverCatalog struct{}

func (resolverCatalog) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: domain.ConnectorGitHubCatalog, Version: "fixture", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"catalog"}, PublicOnly: true, Replayable: true}
}
func (resolverCatalog) Discover(context.Context, domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, Candidates: []domain.ProductCandidate{{Identity: domain.ProductIdentity{ConnectorID: domain.ConnectorGitHubCatalog, ExternalID: "org/repo", CanonicalURL: "https://github.com/org/repo"}, Name: "org/repo", RepositoryLinks: []domain.ProductRepositoryLink{{OwnerRepo: "org/repo", Commit: strings.Repeat("a", 40)}}}}}, nil
}
func TestInsightResolverReadsReferencedPublicRepositoryAtPinnedCommit(t *testing.T) {
	client := NewClient("fixture-token")
	childCalls := 0
	client.httpClient.Transport = resolverTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "api.github.com" || r.Header.Get("Authorization") != "token fixture-token" {
			t.Fatal("wrong host/auth")
		}
		if strings.Contains(r.URL.Path, "/child/idea") {
			childCalls++
		}
		body := ""
		switch {
		case r.URL.Path == "/repos/child/idea":
			body = `{"full_name":"child/idea","html_url":"https://github.com/child/idea","private":false,"default_branch":"main"}`
		case strings.HasSuffix(r.URL.Path, "/commits/main"):
			body = `{"sha":"` + strings.Repeat("b", 40) + `"}`
		case strings.Contains(r.URL.Path, "/git/trees/"):
			commit := strings.Repeat("a", 40)
			if strings.Contains(r.URL.Path, "/child/idea") {
				commit = strings.Repeat("b", 40)
			}
			if !strings.HasSuffix(r.URL.Path, commit) {
				t.Fatal("tree not pinned")
			}
			body = `{"tree":[{"path":"README.md","mode":"100644","type":"blob","size":80}]}`
		case strings.HasSuffix(r.URL.Path, "/contents/README.md"):
			commit, content := strings.Repeat("a", 40), "# Index\n[Child](https://github.com/child/idea)\nAn alternative.\n"
			if strings.Contains(r.URL.Path, "/child/idea") {
				commit = strings.Repeat("b", 40)
				content = "# Child\nAn independent mechanism.\n"
			}
			if r.URL.Query().Get("ref") != commit {
				t.Fatal("file not pinned")
			}
			body = `{"type":"file","path":"README.md","size":` + strconv.Itoa(len(content)) + `,"sha":"` + fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("blob %d\x00%s", len(content), content)))) + `","encoding":"base64","content":"` + base64.StdEncoding.EncodeToString([]byte(content)) + `"}`
		default:
			t.Fatal("unexpected route", r.URL.Path)
		}
		return resolverResponse(r, 200, body), nil
	})
	report, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), insightHTTPRuntime{}), domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "alternative discovery"}, Sources: []string{"github"}}, []domain.InsightSiteBinding{{Connector: resolverCatalog{}, Reader: client, FollowReader: NewInsightLinkedReader(client)}})
	if err != nil || childCalls != 4 || len(report.Followups) != 1 || report.Followups[0].Status != "COLLECTED" || report.Followups[0].Bundle.Repository.Commit != strings.Repeat("b", 40) || len(report.Followups[0].Cards) != 1 {
		t.Fatal("real reader workflow failed", err, childCalls, report.Followups)
	}
	for _, entry := range report.ExplorationMemory.Entries {
		if entry.Identity.ExternalID == "child/idea" && (entry.State != "EXAMINED" || entry.LastRead == nil) {
			t.Fatal("memory missing read witness")
		}
	}
}
func TestInsightResolverRejectsBadIdentityAndUnsafeResponses(t *testing.T) {
	for _, mode := range []string{"private", "missing-private", "wrong-identity", "opaque-commit", "auth", "404", "redirect", "bad-json"} {
		t.Run(mode, func(t *testing.T) {
			client := NewClient("fixture-token")
			calls := 0
			client.httpClient.Transport = resolverTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				body := `{"full_name":"child/idea","html_url":"https://github.com/child/idea","private":false,"default_branch":"main"}`
				status := 200
				if strings.Contains(r.URL.Path, "/commits/") {
					body = `{"sha":"main"}`
				} else {
					switch mode {
					case "private":
						body = strings.Replace(body, `"private":false`, `"private":true`, 1)
					case "missing-private":
						body = strings.Replace(body, `"private":false,`, "", 1)
					case "wrong-identity":
						body = strings.ReplaceAll(body, "child/idea", "other/idea")
					case "auth":
						status = 401
					case "404":
						status = 404
					case "redirect":
						status = 302
					case "bad-json":
						body = "null"
					}
				}
				response := resolverResponse(r, status, body)
				if mode == "redirect" {
					response.Header.Set("Location", "https://unexpected.invalid/private")
				}
				return response, nil
			})
			_, err := client.ResolveRepository(t.Context(), resolverIdentity())
			if err == nil || strings.Contains(err.Error(), "fixture-token") || calls > 2 {
				t.Fatal("bad metadata accepted or secret leaked", err, calls)
			}
		})
	}
	client := NewClient("")
	client.httpClient.Transport = resolverTransport(func(*http.Request) (*http.Response, error) { t.Fatal("bad identity reached HTTP"); return nil, nil })
	identity := resolverIdentity()
	identity.CanonicalURL = "https://gitee.com/child/idea"
	if _, err := client.ResolveRepository(t.Context(), identity); err == nil {
		t.Fatal("foreign host admitted")
	}
}

func TestInsightLinkedReaderDoesNotFallbackOrAdmitBadBlob(t *testing.T) {
	for _, mode := range []string{"404", "blob", "size", "symlink"} {
		t.Run(mode, func(t *testing.T) {
			client := NewClient("secret-fixture")
			calls := 0
			client.httpClient.Transport = resolverTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.URL.Host != "api.github.com" {
					t.Fatal("raw/redirect fallback")
				}
				if mode == "404" {
					return resolverResponse(r, 404, "not found"), nil
				}
				if mode == "symlink" {
					return resolverResponse(r, 200, `{"tree":[{"path":"README.md","type":"blob","mode":"120000"}]}`), nil
				}
				size := 1
				if mode == "size" {
					size = 32769
				}
				body := fmt.Sprintf(`{"type":"file","path":"README.md","encoding":"base64","size":%d,"sha":"%s","content":"YQ=="}`, size, strings.Repeat("a", 40))
				return resolverResponse(r, 200, body), nil
			})
			reader := NewInsightLinkedReader(client)
			repo := domain.RepositoryCandidate{Owner: "child", Repo: "idea", URL: "https://github.com/child/idea", CommitID: strings.Repeat("b", 40)}
			if mode == "symlink" {
				nodes, err := reader.FetchTree(t.Context(), repo)
				if err != nil || len(nodes) != 0 {
					t.Fatal("symlink treated as regular")
				}
			} else {
				if _, err := reader.ExtractRepoFiles(t.Context(), repo, []string{"README.md"}); err == nil {
					t.Fatal("untrusted bytes admitted")
				}
			}
			if calls != 1 {
				t.Fatal("hidden retry or fallback", calls)
			}
		})
	}
}
