/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/investigation_plan_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"testing"
)

func TestStep0043InvestigationUnknownJSONRoundtrip(t *testing.T) {
	input := domain.InvestigationPlan{SchemaVersion: "1", Scope: "one recorded query", StopReason: "NO_CANDIDATES_IN_SCOPE", Requirements: []domain.InvestigationRequirement{{ID: "req", Kind: "mandatory", Text: "bounded queue", Question: "What bounds admission?", Acceptance: "Observe rejection at capacity"}}, Actions: []domain.InvestigationAction{{RequirementID: "req", Classification: "UNVERIFIED", Status: domain.EvidenceInconclusive, NextAction: "Inspect capacity enforcement and execute a saturation test", Acceptance: "Admission never exceeds configured capacity"}}}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var got domain.InvestigationPlan
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input, got) {
		t.Fatalf("unknown contract changed through JSON: %+v", got)
	}
}
