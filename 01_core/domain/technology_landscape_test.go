/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/technology_landscape_test.md
 * @layer L1
 * @updated 2026-09-24
 */
package domain

import (
	"encoding/json"
	"testing"
)

func TestTechnologyLandscapeContractRoundTrip(t *testing.T) {
	in := TechnologyLandscapeRequest{SchemaVersion: TechnologyLandscapeSchemaVersion, Current: LandscapeSnapshot{AsOf: "2026-09-24", Need: []string{"search"}, Candidates: []LandscapeCandidate{{ID: "org/tool", Kind: LandscapeRepository, Name: "Tool", Capabilities: []LandscapeCapability{{ID: "search", Status: EvidenceProven, Sources: []string{"assessment:1"}}}, Standalone: LandscapeAssertion{Status: EvidenceProven, Sources: []string{"test:1"}}}}}}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out TechnologyLandscapeRequest
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if out.Current.Candidates[0].Capabilities[0].Sources[0] != "assessment:1" || out.Current.AsOf != in.Current.AsOf {
		t.Fatalf("round trip perdeu proveniência: %+v", out)
	}
}
