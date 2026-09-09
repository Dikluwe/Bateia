/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/intent_preservation_test.md
 * @layer L3
 * @updated 2026-09-08
 */

package sourcegraph_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/03_infra/sourcegraph"
)

// This adapter uses only public ports and public domain fields.
func intentCall(ctx context.Context, endpoint, token, port, query, language, scope string) (int, []domain.SearchAttempt, error) {
	client := sourcegraph.NewClientWithEndpoint(token, endpoint)
	if port == "direct" {
		candidates, attempts, err := client.SearchCandidatesWithTrace(ctx, domain.QueryFilter{Query: query, Language: language, MaxResults: 5})
		return len(candidates), attempts, err
	}
	result, err := client.Discover(ctx, domain.DiscoveryRequest{Query: query, Language: language, RepositoryScope: scope, Discriminators: []string{"Como criar"}, MaxResults: 5, MaxAttempts: 8})
	attempts := make([]domain.SearchAttempt, len(result.Attempts))
	for i, a := range result.Attempts {
		attempts[i] = domain.SearchAttempt{Query: a.Query, Candidates: a.ResultCount, Error: a.Error}
	}
	return len(result.Candidates), attempts, err
}

func TestIntent0054HTTPPreservesWholeQuestion(t *testing.T) {
	for _, port := range []string{"direct", "discover"} {
		for _, c := range []struct{ name, query, language, prefix string }{
			{"late_negation_symbols", "Como criar cache distribuído sem Redis usando CompareAndSwap<T>() != nil", "go", "lang:Go "},
			{"cpp", "Como obter renderização sem ordenar fragmentos com std::atomic<T>", "cpp", "lang:C++ "},
			{"no_language", "Como criar  índice sem perder cauda: x->y && !z", "", ""},
			{"whitespace_integral", "  Como obter resultado sem truncar o final?  ", "", ""},
		} {
			for _, successful := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%s/success=%v", port, c.name, successful), func(t *testing.T) {
					var calls atomic.Int32
					queries := make(chan string, 16)
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						calls.Add(1)
						var body struct {
							Variables map[string]string `json:"variables"`
						}
						if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
							t.Error(err)
						}
						queries <- body.Variables["query"]
						results := []map[string]interface{}{}
						if successful {
							results = append(results, map[string]interface{}{"__typename": "Repository", "name": "github.com/org/evidence", "url": "/github.com/org/evidence"})
						}
						_, _ = w.Write(mockSGResponse(results))
					}))
					defer server.Close()
					scope := ""
					want := c.prefix + c.query
					if port == "discover" {
						scope = "org/repo"
						want = "repo:^org/repo$ " + want
					}
					count, attempts, err := intentCall(t.Context(), server.URL, "", port, c.query, c.language, scope)
					if err != nil {
						t.Fatal(err)
					}
					expectedCount := 0
					if successful {
						expectedCount = 1
					}
					if calls.Load() != 1 || len(attempts) != 1 || count != expectedCount {
						t.Fatalf("calls=%d attempts=%+v candidates=%d", calls.Load(), attempts, count)
					}
					if sent := <-queries; sent != want {
						t.Fatalf("sent %q; want %q", sent, want)
					}
					if attempts[0].Query != want || attempts[0].Candidates != expectedCount || attempts[0].Error != "" {
						t.Fatalf("trace=%+v", attempts)
					}
					if port == "direct" && attempts[0].OriginalQuery != c.query {
						t.Fatalf("original changed: %+v", attempts)
					}
				})
			}
		}
	}
}

func TestIntent0054RejectsBeforeHTTP(t *testing.T) {
	for _, port := range []string{"direct", "discover"} {
		for _, input := range []string{"", " \t\n ", "cancelled"} {
			t.Run(port+fmt.Sprintf("/%q", input), func(t *testing.T) {
				var calls atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); _, _ = w.Write(mockSGResponse(nil)) }))
				defer server.Close()
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				if input == "cancelled" {
					cancel()
				}
				count, _, err := intentCall(ctx, server.URL, "", port, input, "go", "")
				if err == nil || calls.Load() != 0 || count != 0 {
					t.Fatalf("err=%v calls=%d count=%d", err, calls.Load(), count)
				}
			})
		}
	}
}

func TestIntent0054FailureIsNotEmptySuccess(t *testing.T) {
	for _, port := range []string{"direct", "discover"} {
		for _, mode := range []string{"http", "graphql", "invalid_json", "deadline", "cancel"} {
			t.Run(port+"/"+mode, func(t *testing.T) {
				const token = "intent0054-private-token"
				var calls atomic.Int32
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				if mode == "deadline" {
					var stop context.CancelFunc
					ctx, stop = context.WithTimeout(ctx, 100*time.Millisecond)
					defer stop()
				}
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					switch mode {
					case "http":
						http.Error(w, token, http.StatusInternalServerError)
					case "graphql":
						_, _ = fmt.Fprintf(w, `{"errors":[{"message":%q}]}`, token)
					case "invalid_json":
						_, _ = fmt.Fprintf(w, "invalid %s", token)
					case "deadline":
						select {
						case <-r.Context().Done():
						case <-time.After(time.Second):
						}
					case "cancel":
						cancel()
						select {
						case <-r.Context().Done():
						case <-time.After(time.Second):
						}
					}
				}))
				defer server.Close()
				count, attempts, err := intentCall(ctx, server.URL, token, port, "Como obter dados sem perder o alvo final", "go", "")
				if err == nil || count != 0 || calls.Load() != 1 || len(attempts) != 1 {
					t.Fatalf("err=%v count=%d calls=%d trace=%+v", err, count, calls.Load(), attempts)
				}
				if attempts[0].Error == "" || attempts[0].Candidates != 0 || attempts[0].Query != "lang:Go Como obter dados sem perder o alvo final" {
					t.Fatalf("missing failure trace: %+v", attempts)
				}
				if strings.Contains(err.Error(), token) || strings.Contains(fmt.Sprint(attempts), token) {
					t.Fatal("private token leaked")
				}
			})
		}
	}
}
