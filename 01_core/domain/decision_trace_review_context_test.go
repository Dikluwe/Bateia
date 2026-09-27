/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_review_context_test.md
 * @layer L1
 * @updated 2026-09-27
 */
package domain_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func TestDecisionTraceReviewContextJSONContract(t *testing.T) {
	in := domain.DecisionTraceReviewContextInput{SchemaVersion: "1", ReviewInput: domain.DecisionTraceReviewInput{SchemaVersion: "1"}, BuildInput: domain.DecisionTraceBuildInput{SchemaVersion: "1"}}
	encoded, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.DecisionTraceReviewContextInput
	if err := json.Unmarshal(encoded, &restored); err != nil || !reflect.DeepEqual(in, restored) {
		t.Fatalf("input roundtrip: %v", err)
	}
	out := domain.DecisionTraceReviewContext{SchemaVersion: domain.DecisionTraceReviewContextSchemaVersion, Scope: domain.DecisionTraceScope, ContextScope: domain.DecisionTraceReviewContextScope, Ledger: domain.DecisionTraceReviewLedger{DecisionBefore: domain.DecisionUnknown, DecisionAfter: domain.DecisionUnknown, DecisionImpact: "NONE"}, TaskContexts: []domain.DecisionTraceReviewTaskContext{{NodeID: "claim", ContextNodeIDs: []string{"claim", "source"}}}, Entries: []domain.DecisionTraceExplanationEntry{{NodeID: "claim", Locations: []domain.DecisionTraceLocation{{Document: "report", Pointer: "/reasons/0"}}}, {NodeID: "source", Locations: []domain.DecisionTraceLocation{{Document: "report", Pointer: "/reasons/0/evidence/0"}}}}}
	encoded, err = json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	var decoded domain.DecisionTraceReviewContext
	if err := json.Unmarshal(encoded, &decoded); err != nil || !reflect.DeepEqual(out, decoded) {
		t.Fatalf("context roundtrip: %v", err)
	}
}
