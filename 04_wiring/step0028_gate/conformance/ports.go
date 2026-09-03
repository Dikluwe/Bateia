/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/ports.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"fmt"
)

type ObservePort func(context.Context, Case) (Observation, error)

type Ports struct {
	Extract    ObservePort
	Start      ObservePort
	Lifecycle  ObservePort
	Replay     ObservePort
	Runtime    ObservePort
	Tools      ObservePort
	Refinement ObservePort
	Hypothesis ObservePort
	Safety     ObservePort
	Artifact   ObservePort
	Pagination ObservePort
	Schema     ObservePort
	Full       ObservePort
}

type PortSubject struct{ Ports Ports }

func (subject PortSubject) Observe(ctx context.Context, testCase Case) (Observation, error) {
	var port ObservePort
	switch testCase.Operation {
	case OperationExtract:
		port = subject.Ports.Extract
	case OperationStart:
		port = subject.Ports.Start
	case OperationLifecycle:
		port = subject.Ports.Lifecycle
	case OperationReplay:
		port = subject.Ports.Replay
	case OperationRuntime:
		port = subject.Ports.Runtime
	case OperationTools:
		port = subject.Ports.Tools
	case OperationRefinement:
		port = subject.Ports.Refinement
	case OperationHypothesis:
		port = subject.Ports.Hypothesis
	case OperationSafety:
		port = subject.Ports.Safety
	case OperationArtifact:
		port = subject.Ports.Artifact
	case OperationPagination:
		port = subject.Ports.Pagination
	case OperationSchema:
		port = subject.Ports.Schema
	case OperationFull:
		port = subject.Ports.Full
	}
	if port == nil {
		return Observation{}, fmt.Errorf("unbound conformance port: %s", testCase.Operation)
	}
	return port(ctx, testCase)
}
