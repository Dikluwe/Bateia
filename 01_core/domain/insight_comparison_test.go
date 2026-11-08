/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_comparison_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInsightComparisonWireKeepsSnapshotAndDeclarationScopes(t *testing.T) {
	report := InsightComparisonReport{SchemaVersion: "insight-baseline-comparison-v1", EffectState: "UNVERIFIED", TargetBaseline: &LocatedInsightBaseline{
		Repository: "org/target", SnapshotSHA256: "sha256:snapshot", Scope: "SUPPLIED_FILES_ONLY", Origin: "DECLARED_NOT_REASSESSED",
		Observations: []TargetInsightObservation{{MechanismClass: "class", LocatedInsightObservation: LocatedInsightObservation{ID: "target_obs_1", Kind: InsightStrategy, Statement: "Declared", Origin: "DECLARED_NOT_REASSESSED"}}},
	}, Comparisons: []InsightComparisonRow{{SourceObservationRefs: []string{"obs_1"}, TargetObservationRefs: []string{"target_obs_1"}, DeclarationStatus: "BOTH_DECLARED"}}}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var again InsightComparisonReport
	if err := json.Unmarshal(data, &again); err != nil || !reflect.DeepEqual(report, again) {
		t.Fatalf("snapshot, observations or reference scopes lost: %v", err)
	}
	var request InsightComparisonRequest
	if err := json.Unmarshal([]byte(`{"schema_version":"1","target_baseline":null}`), &request); err != nil || request.TargetBaseline != nil {
		t.Fatal("unknown baseline not preserved")
	}
}
