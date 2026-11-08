/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/insight_reader_test.md
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
	"io"
	"net/http"
	"strings"
	"testing"
)

const readmeFixtureContent = "# Queue\nA bounded mechanism with cooperative cancellation.\n"

func readmeFixtureIdentity(id string) domain.ProductIdentity {
	host := "gitee.com"
	if id == domain.ConnectorAtomGitCatalog {
		host = "atomgit.com"
	}
	return domain.ProductIdentity{ConnectorID: id, ExternalID: "org/repo", CanonicalURL: "https://" + host + "/org/repo"}
}
func readmeFixtureBody(path, host string) string {
	switch {
	case strings.HasSuffix(path, "/search/repositories"):
		return `[{"full_name":"org/repo","html_url":"https://` + host + `/org/repo","private":false,"public":true}]`
	case strings.HasSuffix(path, "/repos/org/repo"):
		return `{"full_name":"org/repo","html_url":"https://` + host + `/org/repo","private":false,"public":true,"default_branch":"main"}`
	case strings.HasSuffix(path, "/commits/main"):
		return `{"sha":"` + strings.Repeat("a", 40) + `"}`
	case strings.Contains(path, "/git/trees/"):
		return `{"tree":[{"path":"README.md","type":"blob","mode":"100644"}],"truncated":false}`
	case strings.Contains(path, "/contents/"):
		object := fmt.Sprintf("blob %d\x00%s", len(readmeFixtureContent), readmeFixtureContent)
		data, _ := json.Marshal(map[string]any{"type": "file", "encoding": "base64", "path": "README.md", "size": len(readmeFixtureContent), "sha": fmt.Sprintf("%x", sha1.Sum([]byte(object))), "content": base64.StdEncoding.EncodeToString([]byte(readmeFixtureContent))})
		return string(data)
	}
	return `{}`
}
func readmeFixtureResponse(r *http.Request, body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}
}
func TestForgeInsightReadmesEndToEndKeepHostsAuthAndPinnedCommits(t *testing.T) {
	counts := map[string]int{}
	client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
		counts[r.URL.Host]++
		host := "gitee.com"
		switch r.URL.Host {
		case "gitee.com":
			if r.URL.Query().Get("access_token") != "gitee-fixture-token" || r.Header.Get("Authorization") != "" {
				t.Fatal("wrong Gitee auth")
			}
		case "api.atomgit.com":
			host = "atomgit.com"
			if r.Header.Get("Authorization") != "Bearer atom-fixture-token" || r.URL.Query().Has("access_token") {
				t.Fatal("wrong AtomGit auth")
			}
		default:
			t.Fatal("credentials crossed host", r.URL.Host)
		}
		if strings.Contains(r.URL.Path, "/git/trees/") && !strings.HasSuffix(r.URL.Path, strings.Repeat("a", 40)) {
			t.Fatal("tree not pinned")
		}
		if strings.Contains(r.URL.Path, "/contents/") && r.URL.Query().Get("ref") != strings.Repeat("a", 40) {
			t.Fatal("branch move affected README")
		}
		return readmeFixtureResponse(r, readmeFixtureBody(r.URL.Path, host)), nil
	})}
	bindings := []domain.InsightSiteBinding{
		{Connector: NewGiteeConnectorWithToken("gitee-fixture-token", client), RepositoryReader: NewGiteeInsightReader("gitee-fixture-token", client)},
		{Connector: NewAtomGitConnector("atom-fixture-token", client), RepositoryReader: NewAtomGitInsightReader("atom-fixture-token", client)},
	}
	request := domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "queue cancellation"}, Sources: []string{"all"}}
	report, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), insightForgeRuntime{}), request, bindings)
	if err != nil || report.Status != "CONTEXT_PREPARED" || len(report.SiteContexts) != 2 || len(report.Bundles) != 0 || len(report.Exploration.Cards) != 0 || counts["gitee.com"] != 9 || counts["api.atomgit.com"] != 9 {
		t.Fatalf("workflow: %v counts=%v report=%+v", err, counts, report)
	}
	ids := map[string]bool{}
	for _, scoped := range report.SiteContexts {
		if len(scoped.Cards) != 1 || ids[scoped.Cards[0].ID] || scoped.Bundle.Items[0].Content != readmeFixtureContent || scoped.Bundle.Repository.Commit != strings.Repeat("a", 40) {
			t.Fatal("cross-site origin lost")
		}
		ids[scoped.Cards[0].ID] = true
	}
	for _, task := range report.Tasks {
		if len(task.ReadingCardRefs) != 2 {
			t.Fatal("missing reading cards")
		}
		for _, id := range task.ReadingCardRefs {
			if !ids[id] {
				t.Fatal("dangling card ref")
			}
		}
	}
	encoded, _ := json.Marshal(report)
	if strings.Contains(string(encoded), "fixture-token") {
		t.Fatal("secret leaked")
	}
}

func TestForgeInsightReaderRefusesIdentityAndMissingTokenWithoutHTTP(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) { calls++; return nil, errors.New("must not call") })}
	reader := NewGiteeInsightReader("token", client)
	for _, rawURL := range []string{"https://github.com/org/repo", "https://gitee.com/org/repo?secret=x", "https://user@gitee.com/org/repo", "http://gitee.com/org/repo", "https://gitee.com:443/org/repo", "https://gitee.com/org/%72epo", "https://gitee.com/org/repo#fragment"} {
		id := readmeFixtureIdentity(domain.ConnectorGiteeCatalog)
		id.CanonicalURL = rawURL
		if _, err := reader.ResolveRepository(t.Context(), id); err == nil {
			t.Fatal("bad identity admitted", rawURL)
		}
	}
	if _, err := NewAtomGitInsightReader("", client).ResolveRepository(t.Context(), readmeFixtureIdentity(domain.ConnectorAtomGitCatalog)); err == nil {
		t.Fatal("missing auth admitted")
	}
	if calls != 0 {
		t.Fatal("invalid request reached transport")
	}
}

func TestForgeInsightReaderRejectsUntrustedResponsesAtEachPhase(t *testing.T) {
	cases := []struct{ mode, phase, code string }{
		{"private", "metadata", "INTEGRITY_ERROR"}, {"identity", "metadata", "INTEGRITY_ERROR"}, {"branch", "metadata", "INTEGRITY_ERROR"}, {"bad-json", "metadata", "INTEGRITY_ERROR"}, {"large", "metadata", "INTEGRITY_ERROR"},
		{"opaque", "commit", "INTEGRITY_ERROR"}, {"auth", "commit", "AUTH_REQUIRED"}, {"404", "tree", "NOT_FOUND"}, {"rate", "tree", "RATE_LIMIT"}, {"truncated", "tree", "TREE_LIMIT"}, {"next-page", "tree", "TREE_LIMIT"}, {"hundred", "tree", "TREE_LIMIT"}, {"symlink", "tree", "README_NOT_FOUND"},
		{"body-limit", "content", "INTEGRITY_ERROR"}, {"base64", "content", "INTEGRITY_ERROR"}, {"size", "content", "INTEGRITY_ERROR"}, {"blob", "content", "INTEGRITY_ERROR"}, {"large", "content", "CONTENT_LIMIT"}, {"binary", "content", "INTEGRITY_ERROR"}, {"redirect", "metadata", "READ_FAILED"},
	}
	for _, tc := range cases {
		t.Run(tc.mode+"-"+tc.phase, func(t *testing.T) {
			calls := 0
			client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.URL.Host != "gitee.com" {
					t.Fatal("redirect leaked credentials")
				}
				path := r.URL.Path
				phase := "metadata"
				if strings.Contains(path, "/commits/") {
					phase = "commit"
				}
				if strings.Contains(path, "/git/trees/") {
					phase = "tree"
				}
				if strings.Contains(path, "/contents/") {
					phase = "content"
				}
				body := readmeFixtureBody(path, "gitee.com")
				response := readmeFixtureResponse(r, body)
				if phase == tc.phase {
					switch tc.mode {
					case "private":
						body = strings.Replace(body, `"private":false`, `"private":true`, 1)
					case "identity":
						body = strings.ReplaceAll(body, "org/repo", "other/repo")
					case "branch":
						body = strings.Replace(body, `"main"`, `".."`, 1)
					case "bad-json":
						body = "null"
					case "large":
						if phase == "metadata" {
							body = strings.Repeat("x", (256<<10)+1)
						} else {
							body = strings.Replace(body, fmt.Sprintf(`"size":%d`, len(readmeFixtureContent)), `"size":32769`, 1)
						}
					case "opaque":
						body = `{"sha":"main"}`
					case "auth":
						response.StatusCode = 401
					case "404":
						response.StatusCode = 404
					case "rate":
						response.StatusCode = 429
					case "truncated":
						body = strings.Replace(body, `"truncated":false`, `"truncated":true`, 1)
					case "next-page":
						response.Header.Set("X-Next-Page", "2")
					case "hundred":
						entries := make([]map[string]string, 100)
						for i := range entries {
							entries[i] = map[string]string{"path": "README.md", "type": "blob", "mode": "100644"}
						}
						b, _ := json.Marshal(map[string]any{"tree": entries})
						body = string(b)
					case "body-limit":
						body = strings.Repeat("x", (64<<10)+1)
					case "symlink":
						body = strings.Replace(body, "100644", "120000", 1)
					case "base64":
						body = strings.Replace(body, base64.StdEncoding.EncodeToString([]byte(readmeFixtureContent)), "!!!", 1)
					case "size":
						body = strings.Replace(body, fmt.Sprintf(`"size":%d`, len(readmeFixtureContent)), `"size":1`, 1)
					case "blob":
						var v map[string]any
						_ = json.Unmarshal([]byte(body), &v)
						v["sha"] = strings.Repeat("b", 40)
						b, _ := json.Marshal(v)
						body = string(b)
					case "binary":
						data := []byte("bad\x00text")
						object := append([]byte(fmt.Sprintf("blob %d\x00", len(data))), data...)
						v := map[string]any{"type": "file", "encoding": "base64", "path": "README.md", "size": len(data), "sha": fmt.Sprintf("%x", sha1.Sum(object)), "content": base64.StdEncoding.EncodeToString(data)}
						b, _ := json.Marshal(v)
						body = string(b)
					case "redirect":
						response.StatusCode = 302
						response.Header.Set("Location", "https://elsewhere.invalid/secret")
					}
				}
				response.Body = io.NopCloser(strings.NewReader(body))
				return response, nil
			})}
			reader := NewGiteeInsightReader("token", client)
			// Run the real L1 collector through the actual reader and connector.
			connector := NewGiteeConnectorWithToken("token", &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
				return readmeFixtureResponse(r, readmeFixtureBody(r.URL.Path, "gitee.com")), nil
			})})
			request := domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "queue"}, Sources: []string{"gitee"}}
			report, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), insightForgeRuntime{}), request, []domain.InsightSiteBinding{{Connector: connector, RepositoryReader: reader}})
			if err != nil || len(report.Readmes) != 1 || report.Readmes[0].Status != tc.code || len(report.SiteContexts) != 0 || report.Status != "PARTIAL" || calls > 4 {
				t.Fatalf("untrusted response accepted %v calls=%d report=%+v", err, calls, report.Readmes)
			}
		})
	}
}
func TestForgeInsightReaderCancellationAndGitCodeAlias(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	calls := 0
	client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		cancel()
		return readmeFixtureResponse(r, readmeFixtureBody(r.URL.Path, "gitee.com")), nil
	})}
	_, err := NewGiteeInsightReader("", client).ResolveRepository(ctx, readmeFixtureIdentity(domain.ConnectorGiteeCatalog))
	if !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatal("cancellation lost")
	}
	client.Transport = insightForgeTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "api.gitcode.com" || r.Header.Get("Authorization") != "Bearer atom-token" {
			t.Fatal("wrong alias/auth")
		}
		return readmeFixtureResponse(r, readmeFixtureBody(r.URL.Path, "gitcode.com")), nil
	})
	id := readmeFixtureIdentity(domain.ConnectorAtomGitCatalog)
	id.CanonicalURL = "https://gitcode.com/org/repo"
	repo, err := NewAtomGitInsightReader("atom-token", client).ResolveRepository(t.Context(), id)
	if err != nil || repo.URL != id.CanonicalURL {
		t.Fatal("known canonical alias lost", err)
	}
}

func TestForgeInsightReaderExactDocumentLimitAndSHA256(t *testing.T) {
	content := strings.Repeat("a", 32<<10)
	client := &http.Client{Transport: insightForgeTransport(func(r *http.Request) (*http.Response, error) {
		object := append([]byte(fmt.Sprintf("blob %d\x00", len(content))), []byte(content)...)
		sum := sha256.Sum256(object)
		encoded, _ := json.Marshal(map[string]any{"type": "file", "encoding": "base64", "path": "README.md", "size": len(content), "sha": fmt.Sprintf("%x", sum), "content": base64.StdEncoding.EncodeToString([]byte(content))})
		if r.URL.Query().Get("ref") != strings.Repeat("a", 64) {
			t.Fatal("64-digit revision lost")
		}
		return readmeFixtureResponse(r, string(encoded)), nil
	})}
	reader := NewAtomGitInsightReader("token", client)
	repo := domain.RepositoryCandidate{Owner: "org", Repo: "repo", URL: "https://atomgit.com/org/repo", CommitID: strings.Repeat("a", 64)}
	files, err := reader.ExtractRepoFiles(t.Context(), repo, []string{"README.md"})
	if err != nil || len(files.Files) != 1 || files.Files[0].Content != content {
		t.Fatal("valid exact-limit SHA256 file rejected", err)
	}
}
