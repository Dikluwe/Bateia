/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/suite.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
)

func DecodeCases(document []byte) ([]Case, error) {
	var envelope struct {
		SchemaVersion  string `json:"schema_version"`
		ContractSHA256 string `json:"contract_sha256"`
		Cases          []Case `json:"cases"`
	}
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return nil, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return nil, fmt.Errorf("trailing conformance document")
	}
	if envelope.SchemaVersion != "tekt.step0028.conformance-cases.v2" || envelope.ContractSHA256 != "f8d8649a9bb3dd1707ae01f3cedc93648b1bd3c501d62a506f745b66a03c928b" || len(envelope.Cases) < 12 {
		return nil, fmt.Errorf("invalid conformance envelope")
	}
	seenCases := map[string]bool{}
	seenObligations := map[Obligation]bool{}
	allowedObligations := map[Obligation]bool{
		P01ExtractionBinding: true, P02EpistemicSeparation: true, P03ProductOnly: true,
		P04UnknownTarget: true, P05License: true, P06Composition: true,
		P07CounterEvidence: true, P08PopularityNovelty: true, P09Idempotency: true,
		P10CancelRecovery: true, P11OfflineReplay: true, P12RefinementRuntime: true,
	}
	allowedOperations := map[Operation]bool{
		OperationExtract: true, OperationStart: true, OperationLifecycle: true,
		OperationReplay: true, OperationRuntime: true, OperationTools: true,
		OperationRefinement: true, OperationHypothesis: true, OperationSafety: true,
		OperationArtifact: true, OperationPagination: true, OperationSchema: true,
		OperationFull: true,
	}
	for _, testCase := range envelope.Cases {
		if testCase.ID == "" || seenCases[testCase.ID] {
			return nil, fmt.Errorf("duplicate or empty case id")
		}
		if !allowedObligations[testCase.Obligation] || !allowedOperations[testCase.Operation] {
			return nil, fmt.Errorf("unknown obligation or operation in %s", testCase.ID)
		}
		seenCases[testCase.ID] = true
		seenObligations[testCase.Obligation] = true
	}
	if len(seenObligations) != 12 {
		return nil, fmt.Errorf("not all obligations are represented")
	}
	return envelope.Cases, nil
}

func Evaluate(ctx context.Context, subject Subject, cases []Case) []Finding {
	findings := make([]Finding, 0, len(cases))
	for _, testCase := range cases {
		got, err := subject.Observe(ctx, testCase)
		finding := Finding{CaseID: testCase.ID, Obligation: testCase.Obligation}
		if err != nil {
			finding.Outcome = "Violated"
			finding.Detail = err.Error()
		} else if !reflect.DeepEqual(got, testCase.Expected) {
			finding.Outcome = "Violated"
			finding.Detail = firstDifference(testCase.Expected, got)
		} else {
			finding.Outcome = "Preserved"
		}
		findings = append(findings, finding)
	}
	sort.Slice(findings, func(i, j int) bool { return findings[i].CaseID < findings[j].CaseID })
	return findings
}

func firstDifference(want, got Observation) string {
	wantBytes, _ := json.Marshal(want)
	gotBytes, _ := json.Marshal(got)
	return fmt.Sprintf("want=%s got=%s", wantBytes, gotBytes)
}
