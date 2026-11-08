/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/forges_test.md
 * @layer L3
 * @updated 2026-10-28
 */

package forge

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

type forgeTransport func(*http.Request) (*http.Response, error)

func (f forgeTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func forgeResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body))}
}

func TestV5SearchProtocolAndEvidence(t *testing.T) {
	for _, provider := range []string{"gitee", "atomgit"} {
		for _, token := range []string{"", "test-token"} {
			if provider == "atomgit" && token == "" {
				continue
			}
			t.Run(provider+"/"+token, func(t *testing.T) {
				host := provider + ".com"
				apiHost := host
				if provider == "atomgit" {
					apiHost = "api." + host
				}
				body := fmt.Sprintf(`[{"full_name":"org/sub/repo","web_url":"https://%s/org/sub/repo","description":"catalogue","default_branch":"main","private":false,"public":true},{"full_name":"org/other","html_url":"https://%s/org/other","private":0,"public":1}]`, host, host)
				calls := 0
				client := &http.Client{Transport: forgeTransport(func(r *http.Request) (*http.Response, error) {
					calls++
					if r.Method != "GET" || r.URL.Host != apiHost || r.URL.Path != "/api/v5/search/repositories" || r.URL.Query().Get("q") != "代码 & tools" || r.URL.Query().Get("per_page") != "1" || r.URL.Query().Get("page") != "1" {
						t.Fatalf("unexpected request: %s", r.URL)
					}
					if provider == "gitee" && (r.URL.Query().Get("access_token") != token || r.Header.Get("Authorization") != "") {
						t.Fatal("incorrect Gitee auth")
					}
					if provider == "atomgit" && (r.Header.Get("Authorization") != "Bearer "+token || r.URL.Query().Get("access_token") != "") {
						t.Fatal("incorrect AtomGit auth")
					}
					return forgeResponse(200, body), nil
				})}
				var c domain.MarketDiscoveryConnector = NewGiteeConnectorWithToken(token, client)
				if provider == "atomgit" {
					c = NewAtomGitConnector(token, client)
				}
				cap := c.Capabilities()
				if strings.Join(cap.IdentityLevels, ",") != "owner_repo" || strings.Join(cap.EvidenceSurfaces, ",") != "repository_identity" {
					t.Fatalf("overstated evidence: %+v", cap)
				}
				result, err := c.Discover(context.Background(), domain.MarketDiscoveryRequest{Query: "代码 & tools", QueryLocale: "zh-CN", Limit: 1})
				if err != nil || calls != 1 || result.Status != domain.AttemptSucceeded || len(result.Candidates) != 1 {
					t.Fatalf("unexpected discovery: %+v %v", result, err)
				}
				if result.ResponseHash != fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(body))) {
					t.Fatal("missing response provenance")
				}
				if _, err := time.Parse(time.RFC3339Nano, result.ObservedAt); err != nil {
					t.Fatal(err)
				}
				candidate := result.Candidates[0]
				if candidate.Identity.ExternalID != "org/sub/repo" || candidate.Identity.ConnectorID != cap.ID || len(candidate.ObservedMarkets) != 0 {
					t.Fatalf("bad identity: %+v", candidate)
				}
				if candidate.RepositoryLinks[0].Commit != "" || candidate.RepositoryLinks[0].Status != domain.RepositoryLinkUnverified {
					t.Fatal("catalogue promoted to verified commit")
				}
				raw, _ := json.Marshal(result)
				if token != "" && strings.Contains(string(raw), token) {
					t.Fatal("credential leaked")
				}
			})
		}
	}
}

func TestV5PublicVisibilityAndIdentity(t *testing.T) {
	cases := []struct {
		name, body string
		ok         bool
	}{
		{"boolean", `{"full_name":"org/repo","web_url":"https://atomgit.com/org/repo","private":false,"public":true}`, true},
		{"integer legacy nested", `{"full_name":"org/sub/repo","web_url":"https://gitcode.com/org/sub/repo","private":0,"public":1,"internal":0}`, true},
		{"private", `{"full_name":"org/repo","web_url":"https://atomgit.com/org/repo","private":1,"public":1}`, false},
		{"internal", `{"full_name":"org/repo","web_url":"https://atomgit.com/org/repo","private":false,"public":true,"internal":true}`, false},
		{"unknown visibility", `{"full_name":"org/repo","web_url":"https://atomgit.com/org/repo"}`, false},
		{"missing public", `{"full_name":"org/repo","web_url":"https://atomgit.com/org/repo","private":false}`, false},
		{"private string", `{"full_name":"org/repo","web_url":"https://atomgit.com/org/repo","private":"false","public":true}`, false},
		{"other host", `{"full_name":"org/repo","web_url":"https://evil.example/org/repo","private":false,"public":true}`, false},
		{"userinfo", `{"full_name":"org/repo","web_url":"https://secret@atomgit.com/org/repo","private":false,"public":true}`, false},
		{"different path", `{"full_name":"org/repo","web_url":"https://atomgit.com/other/repo","private":false,"public":true}`, false},
		{"query", `{"full_name":"org/repo","web_url":"https://atomgit.com/org/repo?token=secret","private":false,"public":true}`, false},
		{"no namespace", `{"name":"repo","web_url":"https://atomgit.com/org/repo","private":false,"public":true}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseForgeResponse(domain.ConnectorAtomGitCatalog, []byte("["+tc.body+"]"), domain.MarketDiscoveryRequest{Limit: 10})
			if err != nil {
				t.Fatal(err)
			}
			if tc.ok {
				if result.Status != domain.AttemptSucceeded || len(result.Candidates) != 1 {
					t.Fatalf("public identity rejected: %+v", result)
				}
			} else if result.Status != domain.AttemptIntegrityError || len(result.Candidates) != 0 {
				t.Fatalf("unsafe identity admitted: %+v", result)
			}
		})
	}
	// Gitee does not require AtomGit's public flag, but does require private=false.
	for _, flag := range []string{"false", "0", "true", "1", "null"} {
		result, _ := parseForgeResponse(domain.ConnectorGiteeCatalog, []byte(fmt.Sprintf(`[{"full_name":"org/repo","html_url":"https://gitee.com/org/repo","private":%s}]`, flag)), domain.MarketDiscoveryRequest{})
		if (len(result.Candidates) == 1) != (flag == "false" || flag == "0") {
			t.Fatalf("Gitee visibility %s: %+v", flag, result)
		}
	}
}

func TestForgeFailuresAreNotEmptySearches(t *testing.T) {
	for _, provider := range []string{"gitee", "atomgit"} {
		for _, tc := range []struct {
			name      string
			status    int
			body      string
			transport bool
			want      domain.MarketAttemptStatus
			class     string
		}{
			{"empty", 200, "[]", false, domain.AttemptEmptyObserved, ""},
			{"unauthorized", 401, "secret", false, domain.AttemptUnknownProvider, "HTTP_401"},
			{"forbidden", 403, "secret", false, domain.AttemptUnknownProvider, "HTTP_403"},
			{"not found", 404, "secret", false, domain.AttemptUnknownProvider, "HTTP_404"},
			{"throttle", 429, "secret", false, domain.AttemptUnknownProvider, "HTTP_429"},
			{"server", 503, "secret", false, domain.AttemptUnknownProvider, "HTTP_503"},
			{"transport", 0, "", true, domain.AttemptUnknownProvider, "TRANSPORT_ERROR"},
			{"malformed", 200, "[", false, domain.AttemptIntegrityError, "INVALID_RESPONSE"},
			{"error object", 200, `{"message":"secret"}`, false, domain.AttemptIntegrityError, "INVALID_RESPONSE"},
			{"null", 200, "null", false, domain.AttemptIntegrityError, "INVALID_RESPONSE"},
			{"oversized", 200, strings.Repeat("x", maximumBodyBytes+1), false, domain.AttemptIntegrityError, "INVALID_RESPONSE"},
		} {
			t.Run(provider+"/"+tc.name, func(t *testing.T) {
				client := &http.Client{Transport: forgeTransport(func(*http.Request) (*http.Response, error) {
					if tc.transport {
						return nil, errors.New("secret URL")
					}
					return forgeResponse(tc.status, tc.body), nil
				})}
				c := NewGiteeConnectorWithToken("secret", client)
				if provider == "atomgit" {
					c = NewAtomGitConnector("secret", client)
				}
				result, err := c.Discover(context.Background(), domain.MarketDiscoveryRequest{Query: "tools"})
				if err != nil || result.Status != tc.want || result.ErrorClass != tc.class || len(result.Candidates) != 0 {
					t.Fatalf("unexpected failure: %+v %v", result, err)
				}
				raw, _ := json.Marshal(result)
				if strings.Contains(string(raw), "secret") {
					t.Fatal("remote details leaked")
				}
			})
		}
	}
}

func TestForgeMissingTokenCancellationAndRedirectIsolation(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: forgeTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		response := forgeResponse(302, "")
		response.Header.Set("Location", "https://other.example/stolen")
		return response, nil
	})}
	result, err := NewAtomGitConnector("", client).Discover(context.Background(), domain.MarketDiscoveryRequest{Query: "tools"})
	if err != nil || result.ErrorClass != "AUTH_REQUIRED" || calls != 0 {
		t.Fatalf("missing auth: %+v %v", result, err)
	}
	for _, c := range []domain.MarketDiscoveryConnector{NewGiteeConnectorWithToken("secret", client), NewAtomGitConnector("secret", client)} {
		before := calls
		result, err = c.Discover(context.Background(), domain.MarketDiscoveryRequest{Query: "tools"})
		if err != nil || result.ErrorClass != "TRANSPORT_ERROR" || calls != before+1 {
			t.Fatalf("redirect followed: %+v %v", result, err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err = c.Discover(ctx, domain.MarketDiscoveryRequest{Query: "tools"})
		if err != nil || result.Status != domain.AttemptCancelled || calls != before+1 {
			t.Fatalf("cancel ignored: %+v %v", result, err)
		}
		_, err = c.Discover(context.Background(), domain.MarketDiscoveryRequest{Query: " "})
		if !errors.Is(err, domain.ErrInvalidRunInput) {
			t.Fatal("invalid query accepted")
		}
	}
	if client.CheckRedirect != nil || client.Timeout != 0 {
		t.Fatal("caller HTTP client mutated")
	}
}

func TestCatalogsDoNotMergeByBranchOrName(t *testing.T) {
	candidates := []domain.ProductCandidate{}
	for _, provider := range []string{"gitee", "atomgit"} {
		result, _ := parseForgeResponse(provider+"-repository-catalog", []byte(fmt.Sprintf(`[{"full_name":"org/repo","html_url":"https://%s.com/org/repo","private":false,"public":true,"default_branch":"main"}]`, provider)), domain.MarketDiscoveryRequest{})
		candidates = append(candidates, result.Candidates...)
	}
	if len(candidates) != 2 {
		t.Fatalf("missing identities: %+v", candidates)
	}
	resolution, _ := usecase.ResolveProductIdentity(candidates[0], candidates[1])
	if resolution == domain.IdentitySame {
		t.Fatal("same name/branch merged across platforms")
	}
}
