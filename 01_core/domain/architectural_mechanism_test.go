/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/architectural_mechanism_test.md
 * @layer L1
 * @updated 2026-09-02
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestArchitecturalMechanismKeepsExplicitEmptyCollections(t *testing.T) {
	value := ArchitecturalMechanism{Components: []MechanismComponentRef{}, DataFlow: []string{}, StateModel: []string{}, FailureModel: []string{}, Boundaries: []string{}, Preconditions: []string{}, Tradeoffs: []string{}, Evidence: []MechanismWitness{}, CounterEvidence: []MechanismWitness{}, ExtractionRuleIDs: []string{}}
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), ":null") {
		t.Fatalf("mechanism collections must be explicit arrays: %s", encoded)
	}
	if MechanismProven != "PROVEN" || NoveltyUnknown != "UNKNOWN" || WitnessDocumentation != "DOCUMENTATION" {
		t.Fatal("public mechanism enums changed")
	}
}
