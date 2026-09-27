/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_review_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func TestDecisionTraceReviewJSONContract(t *testing.T) {
	in := domain.DecisionTraceReviewInput{SchemaVersion: "1", ImpactInput: domain.DecisionTraceImpactInput{SchemaVersion: "1"}, Records: []domain.DecisionTraceReviewRecord{{NodeID: "claim", PlanHash: "declared-plan-hash", TaskHash: "declared-task-hash", Status: "COMPLETED", Note: "declared result", Evidence: []domain.DecisionTraceReviewEvidence{{Ref: "logical-document", ContentHash: "declared-content-hash"}}}}}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out domain.DecisionTraceReviewInput
	if err := json.Unmarshal(data, &out); err != nil || !reflect.DeepEqual(in, out) {
		t.Fatalf("input roundtrip: %v", err)
	}
	ledger := domain.DecisionTraceReviewLedger{SchemaVersion: domain.DecisionTraceReviewSchemaVersion, PlanHash: in.Records[0].PlanHash, Scope: domain.DecisionTraceScope, Status: "RECORDED_REVIEW_COMPLETE", DecisionBefore: domain.DecisionUnknown, DecisionAfter: domain.DecisionUnknown, DecisionImpact: "NONE", Counts: domain.DecisionTraceReviewCounts{Completed: 1}, Tasks: []domain.DecisionTraceReviewProgress{{Task: domain.DecisionTraceReviewTask{NodeID: "claim", Kind: domain.TraceClaim}, TaskHash: in.Records[0].TaskHash, Status: "COMPLETED", Note: in.Records[0].Note, Evidence: in.Records[0].Evidence}}}
	data, err = json.Marshal(ledger)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.DecisionTraceReviewLedger
	if err := json.Unmarshal(data, &restored); err != nil || !reflect.DeepEqual(ledger, restored) {
		t.Fatalf("ledger roundtrip: %v", err)
	}
}
