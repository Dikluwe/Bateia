/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/discovery_coverage_test.md
 * @layer L3
 * @updated 2026-09-08
 */
package sourcegraph_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/sourcegraph"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type coverage0056Result struct {
	candidates []domain.RepositoryCandidate
	attempts   []domain.SearchAttempt
	matches    []domain.ProviderCandidateMatch
}

func coverage0056Call(t *testing.T, port string, limit int, items []map[string]interface{}, matchCount interface{}, mode string, budget int) coverage0056Result {
	t.Helper()
	var calls atomic.Int32
	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Variables map[string]string `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
		n := calls.Add(1)
		want := "lang:Go serializar dados"
		if n == 2 {
			count := 20
			if limit >= 20 {
				count = 200
			} else if limit > 2 {
				count = limit * 10
			}
			want = fmt.Sprintf("lang:Go patterntype:regexp (serializ|marshal) AND (data) count:%d", count)
		}
		if body.Variables["query"] != want {
			t.Errorf("query=%q want %q", body.Variables["query"], want)
		}
		if n == 1 {
			_, _ = w.Write(mockSGResponse([]map[string]interface{}{}))
			return
		}
		switch mode {
		case "http":
			http.Error(w, "unavailable", 503)
			return
		case "graphql":
			_, _ = w.Write([]byte(`{"errors":[{"message":"unavailable"}]}`))
			return
		case "invalid":
			_, _ = w.Write([]byte(`broken`))
			return
		case "deadline":
			select {
			case <-r.Context().Done():
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
		results := map[string]interface{}{"results": items}
		if matchCount != "absent" {
			results["matchCount"] = matchCount
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"search": map[string]interface{}{"results": results}}})
	}))
	t.Cleanup(func() { close(done); server.CloseClientConnections(); server.Close() })
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	client := sourcegraph.NewClientWithEndpoint("", server.URL)
	var out coverage0056Result
	var err error
	if port == "direct" {
		out.candidates, out.attempts, err = client.SearchCandidatesWithTrace(ctx, domain.QueryFilter{Query: "serializar dados", Language: "go", MaxResults: limit})
		if len(out.attempts) > 1 {
			out.matches = out.attempts[1].AdmissionEvidence
		}
	} else {
		result, e := client.Discover(ctx, domain.DiscoveryRequest{Query: "serializar dados", Language: "go", MaxResults: limit, MaxAttempts: budget, AttemptTimeout: 100 * time.Millisecond})
		err = e
		out.candidates = result.Candidates
		out.matches = result.Matches
		for _, a := range result.Attempts {
			out.attempts = append(out.attempts, domain.SearchAttempt{Query: a.Query, Candidates: a.ResultCount, RawCandidates: a.RawCandidates, RejectedCandidates: a.RejectedCandidates, Coverage: a.Coverage, CandidateDiagnostics: a.CandidateDiagnostics, Error: a.Error})
		}
	}
	wantCalls := 2
	if port == "discover" && budget == 1 {
		wantCalls = 1
	}
	if int(calls.Load()) != wantCalls || len(out.attempts) != wantCalls {
		t.Fatalf("calls=%d attempts=%+v", calls.Load(), out.attempts)
	}
	if out.attempts[0].Coverage != nil || len(out.attempts[0].CandidateDiagnostics) != 0 {
		t.Fatal("literal acquired derived coverage")
	}
	if mode != "" {
		if err == nil || out.attempts[len(out.attempts)-1].Error == "" || len(out.candidates) != 0 {
			t.Fatalf("failure became empty success: %v %+v", err, out)
		}
		return out
	}
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func coverage0056Check(t *testing.T, r coverage0056Result, received, observed, examined, eligible, selected, rejected int) {
	t.Helper()
	a := r.attempts[1]
	c := a.Coverage
	if c == nil {
		t.Fatal("missing coverage")
	}
	if c.ReceivedItems != received || c.ObservedRepositories != observed || c.ExaminedRepositories != examined || c.DeferredRepositories != observed-examined || c.EligibleCandidates != eligible || c.SelectedCandidates != selected || c.Completeness != "unknown" || a.RawCandidates != examined || a.RejectedCandidates != rejected || a.Candidates != selected || len(r.candidates) != selected || len(a.CandidateDiagnostics) != examined {
		t.Fatalf("incorrect coverage=%+v attempt=%+v", c, a)
	}
	if len(r.matches) != selected {
		t.Fatalf("matches include nonselected: %+v", r.matches)
	}
	for _, candidate := range r.candidates {
		if len(candidate.SearchAttempts) != 2 || !reflect.DeepEqual(candidate.SearchAttempts[1].Coverage, c) || !reflect.DeepEqual(candidate.SearchAttempts[1].CandidateDiagnostics, a.CandidateDiagnostics) {
			t.Fatalf("candidate lost trace: %+v", candidate)
		}
	}
	for _, m := range r.matches {
		if m.AttemptIndex != 1 || m.EffectiveQuery != a.Query || len(m.Evidence) == 0 {
			t.Fatalf("lost provenance %+v", m)
		}
	}
}
func TestDiscovery0056SelectionAfterAdmission(t *testing.T) {
	items := []map[string]interface{}{discovery0055File("irrelevant", "a.go", "unrelated"), discovery0055File("winner", "b.go", "marshal data"), discovery0055File("eligible", "c.go", "serialize data")}
	var previous *domain.SearchAttempt
	for _, port := range []string{"direct", "discover"} {
		t.Run(port, func(t *testing.T) {
			r := coverage0056Call(t, port, 1, items, 97, "", 2)
			coverage0056Check(t, r, 3, 3, 3, 2, 1, 1)
			a := r.attempts[1]
			if r.candidates[0].Repo != "winner" || a.Coverage.RequestedMatches != 20 || a.Coverage.ProviderMatches == nil || *a.Coverage.ProviderMatches != 97 {
				t.Fatalf("selection/count %+v", r)
			}
			for i, want := range []struct{ status, reason string }{{"insufficient_evidence", "missing_both"}, {"selected", "matched_same_file"}, {"candidate_limit", "candidate_limit"}} {
				d := a.CandidateDiagnostics[i]
				if d.Status != want.status || d.Reason != want.reason {
					t.Errorf("diagnostic=%+v want=%+v", d, want)
				}
			}
			if previous != nil && (!reflect.DeepEqual(previous.Coverage, a.Coverage) || !reflect.DeepEqual(previous.CandidateDiagnostics, a.CandidateDiagnostics)) {
				t.Fatal("ports disagree")
			}
			previous = &a
		})
	}
}
func TestDiscovery0056ObservedDeferredAndLateEvidence(t *testing.T) {
	items := []map[string]interface{}{discovery0055File("winner", "same.go", "marshal")}
	for i := 1; i < 22; i++ {
		items = append(items, discovery0055File(fmt.Sprintf("r%d", i), "a.go", "unrelated"))
	}
	items = append(items, discovery0055File("winner", "same.go", "data"), map[string]interface{}{"__typename": "Unknown"}, map[string]interface{}{"__typename": "Repository", "name": ""})
	for _, port := range []string{"direct", "discover"} {
		t.Run(port, func(t *testing.T) {
			r := coverage0056Call(t, port, 1, items, "absent", "", 2)
			coverage0056Check(t, r, 25, 22, 20, 1, 1, 19)
			if r.candidates[0].Repo != "winner" || r.attempts[1].Coverage.ProviderMatches != nil {
				t.Fatal("late evidence or absence lost")
			}
			if len(r.attempts[1].CandidateDiagnostics[0].Evidence) != 2 {
				t.Fatal("later same-repo evidence omitted")
			}
		})
	}
}
func TestDiscovery0056ReasonsAndMetadata(t *testing.T) {
	items := []map[string]interface{}{{"__typename": "Repository", "name": "github.com/control/marshal-data", "description": "marshal data"}, discovery0055File("empty", "a.go"), discovery0055File("path", "marshal-data.go", "unrelated"), discovery0055File("action", "a.go", "marshal"), discovery0055File("object", "a.go", "data"), discovery0055File("split", "a.go", "marshal"), discovery0055File("split", "b.go", "data"), discovery0055File("same", "a.go", "marshal", "data")}
	for _, port := range []string{"direct", "discover"} {
		t.Run(port, func(t *testing.T) {
			r := coverage0056Call(t, port, 10, items, 0, "", 2)
			coverage0056Check(t, r, 8, 7, 7, 1, 1, 6)
			for i, w := range []struct {
				reason string
				a, o   bool
			}{{"no_file_evidence", false, false}, {"no_file_evidence", false, false}, {"missing_both", false, false}, {"missing_object", true, false}, {"missing_action", false, true}, {"split_across_files", true, true}, {"matched_same_file", true, true}} {
				d := r.attempts[1].CandidateDiagnostics[i]
				if d.Reason != w.reason || d.ActionFound != w.a || d.ObjectFound != w.o {
					t.Errorf("diag %d=%+v expected=%+v", i, d, w)
				}
				for _, e := range d.Evidence {
					if e.Path == "" || e.Snippet == "" || e.MatchType == "" {
						t.Errorf("not retained file evidence %+v", e)
					}
				}
			}
		})
	}
}
func TestDiscovery0056SnippetRetention(t *testing.T) {
	previews := []string{}
	for i := 0; i < 32; i++ {
		previews = append(previews, fmt.Sprintf("irrelevant %d", i))
	}
	previews = append(previews, "marshal data", "marshal data")
	items := []map[string]interface{}{discovery0055File("omitted", "a.go", previews...), discovery0055File("tail", "a.go", strings.Repeat("界", 240)+" marshal data"), discovery0055File("retained", "a.go", "  marshal \t data  "+strings.Repeat("界", 260))}
	for _, port := range []string{"direct", "discover"} {
		t.Run(port, func(t *testing.T) {
			r := coverage0056Call(t, port, 1, items, nil, "", 2)
			coverage0056Check(t, r, 3, 3, 3, 1, 1, 2)
			a := r.attempts[1]
			if a.Coverage.ProviderMatches != nil {
				t.Fatal("null is not zero")
			}
			if r.candidates[0].Repo != "retained" {
				t.Fatal("discarded bytes admitted")
			}
			for i, w := range []struct{ retained, omitted, truncated int }{{32, 2, 0}, {1, 0, 1}, {1, 0, 1}} {
				d := a.CandidateDiagnostics[i]
				if len(d.Evidence) != w.retained || d.OmittedSnippets != w.omitted || d.TruncatedSnippets != w.truncated {
					t.Errorf("retention %+v want %+v", d, w)
				}
				for _, e := range d.Evidence {
					if len([]rune(e.Snippet)) > 240 || strings.Join(strings.Fields(e.Snippet), " ") != e.Snippet {
						t.Errorf("unbounded/unnormalized %+v", e)
					}
				}
			}
		})
	}
}
func TestDiscovery0056CountNullableEmptyAndOverflow(t *testing.T) {
	for _, port := range []string{"direct", "discover"} {
		for _, limit := range []int{-3, 0, 1, 2, 3, 10, 19, 20, int(^uint(0) >> 1)} {
			t.Run(fmt.Sprintf("%s/%d", port, limit), func(t *testing.T) {
				r := coverage0056Call(t, port, limit, []map[string]interface{}{}, 0, "", 2)
				coverage0056Check(t, r, 0, 0, 0, 0, 0, 0)
				c := r.attempts[1].Coverage
				want := 20
				if limit >= 20 {
					want = 200
				} else if limit > 2 {
					want = limit * 10
				}
				if c.RequestedMatches != want || c.ProviderMatches == nil || *c.ProviderMatches != 0 {
					t.Fatalf("coverage=%+v count=%d", c, want)
				}
			})
		}
	}
}
func TestDiscovery0056DerivedFailureAndBudget(t *testing.T) {
	for _, port := range []string{"direct", "discover"} {
		for _, mode := range []string{"http", "graphql", "invalid", "deadline"} {
			t.Run(port+"/"+mode, func(t *testing.T) { coverage0056Call(t, port, 1, nil, nil, mode, 2) })
		}
	}
	coverage0056Call(t, "discover", 1, []map[string]interface{}{}, 0, "", 1)
}
