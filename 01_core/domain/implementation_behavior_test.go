/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/implementation_behavior_test.md
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

func TestBehavior57JSONSchemaAndRoundTrip(t *testing.T) {
	b := domain.ImplementationBehavior{Source: domain.ImplementationReference{Path: "a.go", Symbol: "F", Line: 2, SHA256: "abc"}, EndLine: 4, Status: "observed", Reason: "syntax", Facts: []domain.BehaviorFact{{Kind: "return", Line: 3, EndLine: 3, Target: "", Snippet: "return"}}, Limitations: []string{"unverified"}}
	r := domain.ImplementationSearchReport{Behavior: []domain.ImplementationBehavior{b}}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var got domain.ImplementationSearchReport
	if err = json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Behavior, r.Behavior) {
		t.Fatalf("behavior lost in JSON: %s", data)
	}
	var raw map[string]json.RawMessage
	_ = json.Unmarshal(data, &raw)
	var bs []map[string]json.RawMessage
	_ = json.Unmarshal(raw["behavior"], &bs)
	if len(bs) != 1 {
		t.Fatalf("missing behavior: %s", data)
	}
	for _, key := range []string{"source", "end_line", "status", "reason", "facts", "limitations"} {
		if _, ok := bs[0][key]; !ok {
			t.Errorf("missing %s", key)
		}
	}
	var fs []map[string]json.RawMessage
	_ = json.Unmarshal(bs[0]["facts"], &fs)
	if len(fs) != 1 {
		t.Fatal("missing facts")
	}
	for _, key := range []string{"kind", "line", "end_line", "target", "snippet"} {
		if _, ok := fs[0][key]; !ok {
			t.Errorf("missing fact field %s", key)
		}
	}
}
