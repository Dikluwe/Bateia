/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_compare_test.md
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

func TestDecisionTraceComparisonJSONContract(t *testing.T) {
	node := domain.DecisionTraceNode{ID: "source", Kind: domain.TraceSource, Label: "declared source"}
	in := domain.DecisionTraceComparison{SchemaVersion: domain.DecisionTraceComparisonSchemaVersion, Scope: domain.DecisionTraceScope, HistoricalDecision: domain.DecisionUnknown, CandidateDecision: domain.DecisionBuild, DecisionImpact: "NONE", NodeChanges: []domain.DecisionTraceNodeDelta{{NodeID: node.ID, Status: "ADDED", After: &node}, {NodeID: node.ID, Status: "REMOVED", Before: &node}, {NodeID: node.ID, Status: "DECLARATION_CHANGED", Before: &node, After: &node}}}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out domain.DecisionTraceComparison
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) || out.NodeChanges[0].Before != nil || out.NodeChanges[1].After != nil {
		t.Fatal("comparison JSON lost snapshot distinctions")
	}
}
