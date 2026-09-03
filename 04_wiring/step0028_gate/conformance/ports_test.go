/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/ports_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"testing"
)

func TestStep0028O2PortSubjectRequiresAndDispatchesConcretePorts(t *testing.T) {
	called := false
	subject := PortSubject{Ports: Ports{Start: func(_ context.Context, _ Case) (Observation, error) {
		called = true
		return Observation{Lifecycle: "QUEUED"}, nil
	}}}
	got, err := subject.Observe(context.Background(), Case{Operation: OperationStart})
	if err != nil || !called || got.Lifecycle != "QUEUED" {
		t.Fatalf("start port not dispatched: got=%+v called=%v err=%v", got, called, err)
	}
	if _, err := subject.Observe(context.Background(), Case{Operation: OperationReplay}); err == nil {
		t.Fatal("unbound replay port was silently accepted")
	}
}
