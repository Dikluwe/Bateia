/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o3_ports_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"testing"
)

func TestStep0028O3AllConformancePortsAreExecutable(t *testing.T) {
	port := func(ctx context.Context, testCase Case) (Observation, error) {
		return ReferenceSubject{}.Observe(ctx, testCase)
	}
	subject := PortSubject{Ports: Ports{
		Extract: port, Start: port, Lifecycle: port, Replay: port, Runtime: port,
		Tools: port, Refinement: port, Hypothesis: port, Safety: port,
		Artifact: port, Pagination: port, Schema: port, Full: port,
	}}
	for _, testCase := range O3PositiveCases() {
		findings := Evaluate(context.Background(), subject, []Case{testCase})
		if len(findings) != 1 || findings[0].Outcome != "Preserved" {
			t.Fatalf("port subject failed %s: %+v", testCase.ID, findings)
		}
	}
	operations := []Operation{OperationExtract, OperationStart, OperationLifecycle, OperationReplay, OperationRuntime, OperationTools, OperationRefinement, OperationHypothesis, OperationSafety, OperationArtifact, OperationPagination, OperationSchema, OperationFull}
	for _, operation := range operations {
		if _, err := subject.Observe(context.Background(), Case{Operation: operation}); err != nil {
			t.Fatalf("operation %s has no executable port: %v", operation, err)
		}
	}
}
