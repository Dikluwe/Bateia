/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/public_insight_reader_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package forge

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func publicReadmeFixtureBody(request *http.Request, external string) (string, string) {
	gitlab := request.URL.Host == "gitlab.com"
	path := request.URL.Path
	phase := "metadata"
	if strings.HasSuffix(path, "/repos/search") {
		return `{"data":[{"full_name":"` + external + `","html_url":"https://codeberg.org/` + external + `","private":false}]}`, "catalog"
	}
	if path == "/api/v4/projects" {
		return `[{"path_with_namespace":"` + external + `","web_url":"https://gitlab.com/` + external + `","visibility":"public"}]`, "catalog"
	}
	var body any
	switch {
	case strings.Contains(path, "/branches/"):
		phase = "commit"
		body = map[string]any{"name": "main", "commit": map[string]string{"id": strings.Repeat("a", 40)}}
	case strings.Contains(path, "/repository/commits/"):
		phase = "commit"
		body = map[string]string{"id": strings.Repeat("a", 40)}
	case strings.Contains(path, "/git/trees/") || strings.HasSuffix(path, "/repository/tree"):
		phase = "tree"
		entries := []map[string]string{{"path": "README.md", "type": "blob", "mode": "100644"}}
		if gitlab {
			body = entries
		} else {
			body = map[string]any{"tree": entries, "truncated": false, "total_count": 1}
		}
	case strings.Contains(path, "/contents/") || strings.Contains(path, "/repository/files/"):
		phase = "content"
		content := []byte(readmeFixtureContent)
		sha := fmt.Sprintf("%x", sha1.Sum(append([]byte(fmt.Sprintf("blob %d\x00", len(content))), content...)))
		body = map[string]any{"type": "file", "path": "README.md", "file_path": "README.md", "encoding": "base64", "size": len(content), "sha": sha, "blob_id": sha, "commit_id": strings.Repeat("a", 40), "content_sha256": fmt.Sprintf("%x", sha256.Sum256(content)), "content": base64.StdEncoding.EncodeToString(content)}
	default:
		if gitlab {
			body = map[string]string{"path_with_namespace": external, "web_url": "https://gitlab.com/" + external, "visibility": "public", "default_branch": "main"}
		} else {
			body = map[string]any{"full_name": external, "html_url": "https://codeberg.org/" + external, "private": false, "default_branch": "main"}
		}
	}
	b, _ := json.Marshal(body)
	return string(b), phase
}
func publicReadmeFixtureReader(gitlab bool, client *http.Client) *PublicInsightReader {
	if gitlab {
		return NewGitLabInsightReader(client)
	}
	return NewCodebergInsightReader(client)
}
func publicReadmeFixtureConnector(gitlab bool, client *http.Client) domain.MarketDiscoveryConnector {
	if gitlab {
		return NewGitLabConnector(client)
	}
	return NewCodebergConnector(client)
}
func publicReadmeFixtureIdentity(gitlab bool, external string) domain.ProductIdentity {
	id, host := domain.ConnectorCodebergCatalog, "codeberg.org"
	if gitlab {
		id, host = domain.ConnectorGitLabCatalog, "gitlab.com"
	}
	return domain.ProductIdentity{ConnectorID: id, ExternalID: external, CanonicalURL: "https://" + host + "/" + external}
}
func TestPublicForgeReadmesDiscoverReadAndPrepareNamespacedCards(t *testing.T) {
	for _, nested := range []bool{false, true} {
		t.Run(fmt.Sprintf("nested-%v", nested), func(t *testing.T) {
			counts := map[string]int{}
			client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
				counts[r.URL.Host]++
				external := "org/repo"
				if r.URL.Host != "codeberg.org" && r.URL.Host != "gitlab.com" {
					t.Fatal("wrong host", r.URL.Host)
				}
				if r.URL.Host == "gitlab.com" && nested {
					external = "org/sub/repo"
				}
				body, phase := publicReadmeFixtureBody(r, external)
				prefix := "/api/v1/repos/" + external
				allowed := map[string]bool{"/api/v1/repos/search": true}
				if r.URL.Host == "gitlab.com" {
					prefix = "/api/v4/projects/" + url.PathEscape(external)
					allowed = map[string]bool{"/api/v4/projects": true, prefix + "/repository/commits/main": true, prefix + "/repository/tree": true, prefix + "/repository/files/README.md": true}
				} else {
					allowed[prefix+"/branches/main"] = true
					allowed[prefix+"/git/trees/"+strings.Repeat("a", 40)] = true
					allowed[prefix+"/contents/README.md"] = true
				}
				allowed[prefix] = true
				if !allowed[r.URL.EscapedPath()] {
					t.Fatal("unexpected endpoint", r.URL.EscapedPath())
				}
				if phase == "tree" && r.URL.Query().Get("per_page") != "100" {
					t.Fatal("tree limit omitted")
				}
				if r.Header.Get("Authorization") != "" || r.Header.Get("Private-Token") != "" || r.URL.Query().Has("access_token") {
					t.Fatal("anonymous reader sent credentials")
				}
				if r.URL.Host == "gitlab.com" && phase != "catalog" && !strings.Contains(r.URL.EscapedPath(), url.PathEscape(external)) {
					t.Fatal("project namespace not encoded", r.URL.EscapedPath())
				}
				if phase == "tree" || phase == "content" {
					if r.URL.Host == "gitlab.com" || phase == "content" {
						if r.URL.Query().Get("ref") != strings.Repeat("a", 40) {
							t.Fatal("mutable branch used")
						}
					} else if !strings.HasSuffix(r.URL.Path, strings.Repeat("a", 40)) {
						t.Fatal("tree not pinned")
					}
					if r.URL.Query().Has("recursive") {
						t.Fatal("recursive tree requested")
					}
				}
				return readmeFixtureResponse(r, body), nil
			})}
			bindings := []domain.InsightSiteBinding{{Connector: NewCodebergConnector(client), RepositoryReader: NewCodebergInsightReader(client)}, {Connector: NewGitLabConnector(client), RepositoryReader: NewGitLabInsightReader(client)}}
			report, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), insightForgeRuntime{}), domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "queue cancellation"}, Sources: []string{"codeberg", "gitlab"}}, bindings)
			if err != nil || report.Status != "CONTEXT_PREPARED" || len(report.SiteContexts) != 2 || counts["codeberg.org"] != 9 || counts["gitlab.com"] != 9 {
				t.Fatalf("collection failed: %v counts=%v report=%+v", err, counts, report)
			}
			ids := map[string]bool{}
			for _, scoped := range report.SiteContexts {
				if len(scoped.Cards) != 1 || ids[scoped.Cards[0].ID] || scoped.Bundle.Items[0].Content != readmeFixtureContent {
					t.Fatal("invalid cards or origin collision")
				}
				ids[scoped.Cards[0].ID] = true
				if nested && scoped.Identity.ConnectorID == domain.ConnectorGitLabCatalog && (scoped.Bundle.Repository.Owner != "org/sub" || scoped.Cards[0].Source.Repository != "org/sub/repo") {
					t.Fatal("subgroup lost")
				}
			}
			for _, task := range report.Tasks {
				if len(task.ReadingCardRefs) != 2 {
					t.Fatal("cards missing")
				}
				for _, id := range task.ReadingCardRefs {
					if !ids[id] {
						t.Fatal("dangling ref")
					}
				}
			}
		})
	}
}
func TestPublicForgeReadmesRejectMalformedAndIncompleteResponses(t *testing.T) {
	cases := []struct{ mode, phase, status string }{
		{"private", "metadata", "INTEGRITY_ERROR"}, {"missing-public", "metadata", "INTEGRITY_ERROR"}, {"identity", "metadata", "INTEGRITY_ERROR"}, {"branch", "metadata", "INTEGRITY_ERROR"}, {"body-limit", "metadata", "INTEGRITY_ERROR"}, {"opaque", "commit", "INTEGRITY_ERROR"},
		{"pagination", "tree", "TREE_LIMIT"}, {"hundred", "tree", "TREE_LIMIT"}, {"symlink", "tree", "README_NOT_FOUND"}, {"null", "tree", "INTEGRITY_ERROR"}, {"rate", "tree", "RATE_LIMIT"},
		{"auth", "content", "AUTH_REQUIRED"}, {"404", "content", "NOT_FOUND"}, {"base64", "content", "INTEGRITY_ERROR"}, {"size", "content", "INTEGRITY_ERROR"}, {"commit-id", "content", "INTEGRITY_ERROR"}, {"content-hash", "content", "INTEGRITY_ERROR"}, {"blob", "content", "INTEGRITY_ERROR"}, {"path", "content", "INTEGRITY_ERROR"}, {"large", "content", "CONTENT_LIMIT"}, {"body-limit", "content", "INTEGRITY_ERROR"}, {"redirect", "metadata", "READ_FAILED"},
	}
	for _, gitlab := range []bool{false, true} {
		for _, tc := range cases {
			if !gitlab && (tc.mode == "commit-id" || tc.mode == "content-hash") {
				continue
			}
			t.Run(fmt.Sprintf("gitlab-%v-%s-%s", gitlab, tc.mode, tc.phase), func(t *testing.T) {
				calls := 0
				client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
					calls++
					body, phase := publicReadmeFixtureBody(r, "org/repo")
					resp := readmeFixtureResponse(r, body)
					if phase == tc.phase {
						var value map[string]any
						if phase != "tree" {
							_ = json.Unmarshal([]byte(body), &value)
						}
						switch tc.mode {
						case "private":
							value["private"] = true
							value["visibility"] = "private"
						case "missing-public":
							delete(value, "private")
							delete(value, "visibility")
						case "identity":
							value["full_name"] = "other/repo"
							value["path_with_namespace"] = "other/repo"
						case "branch":
							value["default_branch"] = ".."
						case "opaque":
							if gitlab {
								value["id"] = "main"
							} else {
								value["commit"] = map[string]string{"id": "main"}
							}
						case "pagination":
							resp.Header.Set("X-Next-Page", "2")
						case "hundred":
							entries := make([]map[string]string, 100)
							for i := range entries {
								entries[i] = map[string]string{"path": "README.md", "type": "blob", "mode": "100644"}
							}
							var v any = entries
							if !gitlab {
								v = map[string]any{"tree": entries}
							}
							b, _ := json.Marshal(v)
							body = string(b)
						case "symlink":
							body = strings.Replace(body, "100644", "120000", 1)
						case "null":
							body = "null"
						case "rate":
							resp.StatusCode = 429
						case "auth":
							resp.StatusCode = 401
						case "404":
							resp.StatusCode = 404
						case "base64":
							value["content"] = "!!!"
						case "size":
							value["size"] = 1
						case "commit-id":
							value["commit_id"] = strings.Repeat("b", 40)
						case "content-hash":
							value["content_sha256"] = strings.Repeat("b", 64)
						case "blob":
							value["sha"] = strings.Repeat("b", 40)
							value["blob_id"] = strings.Repeat("b", 40)
						case "path":
							value["path"] = "other.md"
							value["file_path"] = "other.md"
						case "large":
							value["size"] = 32769
						case "body-limit":
							limit := 256 << 10
							if phase == "content" {
								limit = 64 << 10
							}
							body = strings.Repeat("x", limit+1)
							value = nil
						case "redirect":
							resp.StatusCode = 302
							resp.Header.Set("Location", "https://unexpected.invalid/README.md")
						}
						if value != nil {
							b, _ := json.Marshal(value)
							body = string(b)
						}
					}
					return &http.Response{StatusCode: resp.StatusCode, Header: resp.Header, Body: readmeFixtureResponse(r, body).Body, Request: r}, nil
				})}
				// Catalog remains valid to isolate each failure in the reader itself.
				catalogClient := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
					body, _ := publicReadmeFixtureBody(r, "org/repo")
					return readmeFixtureResponse(r, body), nil
				})}
				binding := domain.InsightSiteBinding{Connector: publicReadmeFixtureConnector(gitlab, catalogClient), RepositoryReader: publicReadmeFixtureReader(gitlab, client)}
				report, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), insightForgeRuntime{}), domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "queue"}, Sources: []string{"all"}}, []domain.InsightSiteBinding{binding})
				if err != nil || report.Status != "PARTIAL" || len(report.SiteContexts) != 0 || len(report.Readmes) != 1 || report.Readmes[0].Status != tc.status || calls > 4 {
					t.Fatalf("bad response accepted: %v calls=%d receipts=%+v", err, calls, report.Readmes)
				}
			})
		}
	}
}
func TestPublicForgeReadmesRejectInvalidOriginsBeforeNetwork(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) { calls++; return nil, errors.New("unexpected") })}
	for _, gitlab := range []bool{false, true} {
		reader := publicReadmeFixtureReader(gitlab, client)
		for _, bad := range []string{"https://elsewhere.invalid/org/repo", "http://gitlab.com/org/repo", "https://gitlab.com/org/repo?secret=x", "https://gitlab.com/org/%72epo"} {
			id := publicReadmeFixtureIdentity(gitlab, "org/repo")
			id.CanonicalURL = bad
			if _, err := reader.ResolveRepository(t.Context(), id); err == nil {
				t.Fatal("bad URL admitted")
			}
		}
	}
	if calls != 0 {
		t.Fatal("invalid identity reached network")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := NewGitLabInsightReader(client).ResolveRepository(ctx, publicReadmeFixtureIdentity(true, "org/repo")); !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation lost", err)
	}
	if calls != 0 {
		t.Fatal("cancelled request reached network")
	}
}

func TestPublicForgeReadmesExactLimitSHA256AndPartialFailure(t *testing.T) {
	content := strings.Repeat("a", 32<<10)
	for _, gitlab := range []bool{false, true} {
		identity := publicReadmeFixtureIdentity(gitlab, "org/repo")
		client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
			if r.URL.Query().Get("ref") != strings.Repeat("a", 64) {
				t.Fatal("revision lost")
			}
			object := append([]byte(fmt.Sprintf("blob %d\x00", len(content))), []byte(content)...)
			sha := fmt.Sprintf("%x", sha256.Sum256(object))
			body, _ := json.Marshal(map[string]any{"type": "file", "path": "README.md", "file_path": "README.md", "encoding": "base64", "size": len(content), "sha": sha, "blob_id": sha, "commit_id": strings.Repeat("a", 64), "content_sha256": fmt.Sprintf("%x", sha256.Sum256([]byte(content))), "content": base64.StdEncoding.EncodeToString([]byte(content))})
			return readmeFixtureResponse(r, string(body)), nil
		})}
		repo := domain.RepositoryCandidate{Owner: "org", Repo: "repo", URL: identity.CanonicalURL, CommitID: strings.Repeat("a", 64)}
		extracted, err := publicReadmeFixtureReader(gitlab, client).ExtractRepoFiles(t.Context(), repo, []string{"README.md"})
		if err != nil || len(extracted.Files) != 1 || extracted.Files[0].Content != content {
			t.Fatal("exact-limit/SHA256 rejected", err)
		}
	}
	client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
		body, phase := publicReadmeFixtureBody(r, "org/repo")
		response := readmeFixtureResponse(r, body)
		if r.URL.Host == "gitlab.com" && phase == "metadata" {
			response.StatusCode = 401
		}
		return response, nil
	})}
	report, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), insightForgeRuntime{}), domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "queue"}, Sources: []string{"all"}}, []domain.InsightSiteBinding{{Connector: NewCodebergConnector(client), RepositoryReader: NewCodebergInsightReader(client)}, {Connector: NewGitLabConnector(client), RepositoryReader: NewGitLabInsightReader(client)}})
	if err != nil || report.Status != "PARTIAL" || len(report.SiteContexts) != 1 || report.SiteContexts[0].Identity.ConnectorID != domain.ConnectorCodebergCatalog || len(report.Tasks[0].ReadingCardRefs) != 1 {
		t.Fatal("source failure lost other source", err)
	}
}
