/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_compare_boundary_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestDecisionTraceCompareRootBudgetDoesNotTruncate(t *testing.T) {
	g := traceFixture(t)
	g.GraphHash = ""
	for i := 0; i < 101; i++ {
		id := fmt.Sprintf("bounded-input-%03d", i)
		g.Nodes = append(g.Nodes, domain.DecisionTraceNode{ID: id, Kind: domain.TraceInput, Fingerprint: strings.Repeat("a", 64)})
		g.Edges = append(g.Edges, domain.DecisionTraceEdge{From: id, To: "decision", Relation: domain.TraceDependsOn})
	}
	before, err := FreezeDecisionTrace(g)
	if err != nil {
		t.Fatal(err)
	}
	for _, count := range []int{100, 101} {
		after := traceCloneGraph(t, before)
		after.GraphHash = ""
		changed := 0
		for i := range after.Nodes {
			if strings.HasPrefix(after.Nodes[i].ID, "bounded-input-") && changed < count {
				after.Nodes[i].Label = "different declared label"
				changed++
			}
		}
		after, err = FreezeDecisionTrace(after)
		if err != nil {
			t.Fatal(err)
		}
		result, err := CompareDecisionTraces(domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: before, After: after})
		if count == 100 {
			if err != nil || len(result.NodeChanges) != 100 || len(result.Impact.ChangedNodeIDs) != 100 {
				t.Fatalf("100 roots rejected/truncated: %v %+v", err, result)
			}
		} else if err == nil || !reflect.DeepEqual(result, domain.DecisionTraceComparison{}) {
			t.Fatalf("oversized roots leaked a partial result: %v %+v", err, result)
		}
	}
}

func TestDecisionTraceCompareAggregateBudget(t *testing.T) {
	g := traceFixture(t)
	g.GraphHash = ""
	for i := 0; i < 160; i++ {
		g.Nodes = append(g.Nodes, domain.DecisionTraceNode{ID: fmt.Sprintf("payload-%03d", i), Kind: domain.TraceClaim, Fingerprint: strings.Repeat("a", 64), Label: strings.Repeat("x", 7000)})
	}
	g, err := FreezeDecisionTrace(g)
	if err != nil {
		t.Fatalf("individual graph should fit: %v", err)
	}
	result, err := CompareDecisionTraces(domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: g, After: g})
	if err == nil || !reflect.DeepEqual(result, domain.DecisionTraceComparison{}) {
		t.Fatalf("aggregate request must be bounded: %v", err)
	}
}
