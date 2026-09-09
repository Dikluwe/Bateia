/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/discovery_coverage_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"testing"
)

func TestDiscovery0056CoverageJSONContract(t *testing.T) {
	zero := 0
	for _, matches := range []*int{nil, &zero} {
		c := domain.DiscoveryCoverage{RequestedMatches: 20, ReceivedItems: 29, ObservedRepositories: 22, ExaminedRepositories: 20, DeferredRepositories: 2, EligibleCandidates: 3, SelectedCandidates: 1, Completeness: "unknown", ProviderMatches: matches}
		b, err := json.Marshal(c)
		if err != nil {
			t.Fatal(err)
		}
		var got map[string]interface{}
		if err = json.Unmarshal(b, &got); err != nil {
			t.Fatal(err)
		}
		for k, v := range map[string]interface{}{"requested_matches": float64(20), "received_items": float64(29), "observed_repositories": float64(22), "examined_repositories": float64(20), "deferred_repositories": float64(2), "eligible_candidates": float64(3), "selected_candidates": float64(1), "completeness": "unknown"} {
			if got[k] != v {
				t.Errorf("%s=%v want=%v JSON=%s", k, got[k], v, b)
			}
		}
		if matches == nil {
			if got["provider_matches"] != nil {
				t.Fatalf("absence became zero: %s", b)
			}
		} else if got["provider_matches"] != float64(0) {
			t.Fatalf("zero lost: %s", b)
		}
		var round domain.DiscoveryCoverage
		if err = json.Unmarshal(b, &round); err != nil || !reflect.DeepEqual(c, round) {
			t.Fatalf("roundtrip %+v err=%v", round, err)
		}
	}
}
func TestDiscovery0056AttemptJSONOptionalAndComplete(t *testing.T) {
	diag := domain.DiscoveryCandidateDiagnostic{Repository: "control/example", Status: "insufficient_evidence", Reason: "split_across_files", ActionFound: true, ObjectFound: true, TruncatedSnippets: 2, OmittedSnippets: 4, Evidence: []domain.MatchEvidence{{MatchType: "FileMatch", Path: "a.go", Snippet: "marshal"}, {MatchType: "FileMatch", Path: "b.go", Snippet: "data"}}}
	coverage := &domain.DiscoveryCoverage{Completeness: "unknown", RequestedMatches: 20}
	for _, a := range []interface{}{domain.SearchAttempt{}, domain.ProviderAttempt{}} {
		b, _ := json.Marshal(a)
		var got map[string]interface{}
		_ = json.Unmarshal(b, &got)
		for _, k := range []string{"coverage", "candidate_diagnostics"} {
			if _, ok := got[k]; ok {
				t.Fatalf("optional %s present %s", k, b)
			}
		}
	}
	for _, a := range []interface{}{domain.SearchAttempt{Coverage: coverage, CandidateDiagnostics: []domain.DiscoveryCandidateDiagnostic{diag}}, domain.ProviderAttempt{Coverage: coverage, CandidateDiagnostics: []domain.DiscoveryCandidateDiagnostic{diag}}} {
		b, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}
		var envelope struct {
			Coverage    *domain.DiscoveryCoverage             `json:"coverage"`
			Diagnostics []domain.DiscoveryCandidateDiagnostic `json:"candidate_diagnostics"`
		}
		if err = json.Unmarshal(b, &envelope); err != nil || !reflect.DeepEqual(envelope.Coverage, coverage) || !reflect.DeepEqual(envelope.Diagnostics, []domain.DiscoveryCandidateDiagnostic{diag}) {
			t.Fatalf("lost diagnostics %s err=%v", b, err)
		}
		var raw map[string]json.RawMessage
		_ = json.Unmarshal(b, &raw)
		var list []map[string]interface{}
		_ = json.Unmarshal(raw["candidate_diagnostics"], &list)
		for k, w := range map[string]interface{}{"repository": "control/example", "status": "insufficient_evidence", "reason": "split_across_files", "action_found": true, "object_found": true, "truncated_snippets": float64(2), "omitted_snippets": float64(4)} {
			if list[0][k] != w {
				t.Errorf("%s=%v want=%v", k, list[0][k], w)
			}
		}
	}
}
