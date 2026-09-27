/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/composition_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

import (
	"encoding/json"
	"testing"
)

func TestCompositionContractPreservesUnknownCostAndExecutionAttestation(t *testing.T) {
	input := CompositionInput{Verifications: []CompositionVerification{{Members: []string{"a", "b"}, Status: EvidenceProven, Executed: true, Passed: true}}}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CompositionInput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Verifications) != 1 || !decoded.Verifications[0].Executed || !decoded.Verifications[0].Passed {
		t.Fatalf("attestation lost: %+v", decoded)
	}
	candidate := CompositionCandidate{Status: CompositionProposed}
	data, err = json.Marshal(candidate)
	if err != nil {
		t.Fatal(err)
	}
	var roundtrip CompositionCandidate
	if err := json.Unmarshal(data, &roundtrip); err != nil || roundtrip.IntegrationCost != nil {
		t.Fatalf("unknown integration cost became numeric: %+v, %v", roundtrip, err)
	}
}
