/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/discovery_query_test.md
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

func TestDiscovery0055PlanJSONTrace(t *testing.T) {
	p := domain.DiscoveryQueryPlan{Original: "serializar dados", Version: "0055-v2", Status: "hypothesis", Reason: "catalog", ActionTerms: []string{"serializ", "marshal"}, ObjectTerms: []string{"data"}, RuleIDs: []string{"synthetic-rule"}, Unverified: []string{"serializar dados"}}
	for _, v := range []any{domain.SearchAttempt{QueryPlan: &p, RawCandidates: 7, RejectedCandidates: 6, Candidates: 1}, domain.ProviderAttempt{QueryPlan: &p, RawCandidates: 7, RejectedCandidates: 6, ResultCount: 1}} {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		var m map[string]json.RawMessage
		if err = json.Unmarshal(b, &m); err != nil {
			t.Fatal(err)
		}
		var got domain.DiscoveryQueryPlan
		if err = json.Unmarshal(m["query_plan"], &got); err != nil || !reflect.DeepEqual(got, p) {
			t.Fatalf("plan lost %s err=%v", b, err)
		}
		if string(m["raw_candidates"]) != "7" || string(m["rejected_candidates"]) != "6" {
			t.Fatalf("counters lost %s", b)
		}
	}
	for _, v := range []any{domain.SearchAttempt{}, domain.ProviderAttempt{}} {
		b, _ := json.Marshal(v)
		var m map[string]any
		_ = json.Unmarshal(b, &m)
		if _, ok := m["query_plan"]; ok {
			t.Fatalf("literal carries plan %s", b)
		}
	}
}
