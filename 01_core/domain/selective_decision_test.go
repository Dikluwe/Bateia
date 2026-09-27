/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/selective_decision_test.md
 * @layer L1
 * @updated 2026-09-24
 */
package domain

import (
	"encoding/json"
	"testing"
)

func TestSelectiveDecisionContractRoundTrip(t *testing.T) {
	score := 0.8
	request := SelectiveDecisionRequest{SchemaVersion: "1", Split: SplitTest, Cases: []SelectiveDecisionCase{{ID: "a", Observed: DecisionBuild, Expected: DecisionBuild, Reliability: &score}}}
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var decoded SelectiveDecisionRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Split != SplitTest || len(decoded.Cases) != 1 || decoded.Cases[0].Reliability == nil || *decoded.Cases[0].Reliability != score {
		t.Fatalf("round trip perdeu dados: %+v", decoded)
	}
}
