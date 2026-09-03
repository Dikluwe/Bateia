/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/suite_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"testing"
)

type deviatingSubject struct{ ReferenceSubject }

func (subject deviatingSubject) Observe(ctx context.Context, testCase Case) (Observation, error) {
	result, err := subject.ReferenceSubject.Observe(ctx, testCase)
	result.BindingVerified = false
	return result, err
}

func TestStep0028O2EvaluatorRejectsAnObservableDeviation(t *testing.T) {
	testCase := Case{
		ID: "DEVIATION-CONTROL", Obligation: P01ExtractionBinding, Operation: OperationExtract,
		Input:    Input{Implementation: true, Independent: true, SameImmutableSource: true, ExactBundleBinding: true, ContentHashValid: true, ParentStatus: "PROVEN"},
		Expected: Observation{Decision: "MECHANISM_PRESERVED", MechanismStatus: "PROVEN", Novelty: "DIVERGENT", WarningCodes: []string{}, IntegrationUnknowns: []string{}, Tools: []Tool{}, LeakedFields: []string{}, ExternalEffects: []string{}, BindingVerified: true, ProvenanceComplete: true},
	}
	findings := Evaluate(context.Background(), deviatingSubject{}, []Case{testCase})
	if len(findings) != 1 || findings[0].Outcome != "Violated" || findings[0].Detail == "" {
		t.Fatalf("observable deviation was not rejected: %+v", findings)
	}
}

func TestStep0028O2DecoderRejectsTrailingDocument(t *testing.T) {
	document := []byte(`{"schema_version":"tekt.step0028.conformance-cases.v2","contract_sha256":"f8d8649a9bb3dd1707ae01f3cedc93648b1bd3c501d62a506f745b66a03c928b","cases":[]} {}`)
	if _, err := DecodeCases(document); err == nil {
		t.Fatal("decoder accepted a trailing JSON value")
	}
}
