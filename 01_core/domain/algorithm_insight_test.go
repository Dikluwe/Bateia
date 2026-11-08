/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/algorithm_insight_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAlgorithmInsightWireKeepsDeclaredAndUnverifiedStates(t *testing.T) {
	r := AlgorithmInsightReport{SchemaVersion: "insight-algorithms-v1", InterpretationState: "DECLARED_NOT_REASSESSED", TargetRelationState: "DECLARED_NOT_REASSESSED", EffectState: "UNVERIFIED", Ideas: []LocatedAlgorithmIdea{{ID: "id", AlgorithmIdea: AlgorithmIdea{TargetRelation: "NOT_ESTABLISHED", Decision: AlgorithmDecision{Current: "current", Proposed: "candidate", Benefit: "hypothesis", Cost: "work"}}}}}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var got AlgorithmInsightReport
	if err := json.Unmarshal(data, &got); err != nil || !reflect.DeepEqual(r, got) {
		t.Fatal("states lost", err)
	}
}
