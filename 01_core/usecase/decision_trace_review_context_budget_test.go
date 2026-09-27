/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_review_context_budget_test.md
 * @layer L1
 * @updated 2026-09-27
 */
package usecase

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestDecisionTraceReviewContextExpandedProjectionBudget(t *testing.T) {
	// This deliberately tests structural output admission, not document coherence.
	projection := traceReviewContextRun(t, traceReviewContextFixture(t))
	projection.Ledger.Tasks = nil
	projection.TaskContexts = nil
	roots := make([]string, 100)
	for i := range roots {
		roots[i] = fmt.Sprintf("%064x", i+1)
	}
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("%064x", i+1001)
		task := domain.DecisionTraceReviewTask{NodeID: id, Kind: domain.TraceClaim, ChangedAncestorIDs: append([]string{}, roots...), Reason: "RECORDED_DEPENDENCY_CHANGED"}
		projection.Ledger.Tasks = append(projection.Ledger.Tasks, domain.DecisionTraceReviewProgress{Task: task, TaskHash: strings.Repeat("a", 64), Status: "IN_PROGRESS", Note: strings.Repeat("n", 8192), Evidence: []domain.DecisionTraceReviewEvidence{}})
		projection.TaskContexts = append(projection.TaskContexts, domain.DecisionTraceReviewTaskContext{NodeID: id, ContextNodeIDs: append(append([]string{}, roots...), id)})
	}
	if _, err := traceBoundedJSON(projection.Ledger); err != nil {
		t.Fatalf("ledger should fit alone: %v", err)
	}
	if _, err := traceBoundedJSON(projection.Entries); err != nil {
		t.Fatalf("index should fit alone: %v", err)
	}
	out, err := boundedReviewContextProjection(projection)
	if err == nil || !strings.Contains(err.Error(), "review context output exceeds budget") || !reflect.DeepEqual(out, domain.DecisionTraceReviewContext{}) {
		t.Fatalf("oversized join leaked a partial result: %v", err)
	}
	for i := range projection.Ledger.Tasks {
		projection.Ledger.Tasks[i].Note = strings.Repeat("n", 4096)
	}
	out, err = boundedReviewContextProjection(projection)
	if err != nil || !reflect.DeepEqual(out, projection) {
		t.Fatalf("bounded positive projection rejected: %v", err)
	}
}
