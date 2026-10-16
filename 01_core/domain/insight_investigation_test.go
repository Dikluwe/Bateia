/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_investigation_test.md
 * @layer L1
 * @updated 2026-10-06
 */

package domain

import (
	"encoding/json"
	"testing"
)

func TestInsightReportWireContractPreservesDistinctStatesAndOrigin(t *testing.T) {
	raw := []byte(`{"schema_version":"insight-investigation-v1","scope":"FROZEN_EVALUATED_EVIDENCE","status":"PARTIAL","mechanisms":[{"status":"PROVEN"}],"hypotheses":[{"status":"PROPOSED","target_state":"UNKNOWN","adoption_disposition":"NOT_AN_ADOPTION_DECISION"}],"provenance":[{"assessment_origin":"DECLARED_NOT_REASSESSED","source":{"repository":"example/library","commit":"abc"}}]}`)
	var report InsightInvestigationReport
	if err := json.Unmarshal(raw, &report); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var again InsightInvestigationReport
	if err := json.Unmarshal(encoded, &again); err != nil {
		t.Fatal(err)
	}
	if again.Status != ImprovementPartial || again.Mechanisms[0].Status != MechanismProven || again.Hypotheses[0].Status != HypothesisProposed || again.Hypotheses[0].TargetState != TargetUnknown || again.Provenance[0].AssessmentOrigin != "DECLARED_NOT_REASSESSED" || again.Provenance[0].Source.Repository != "example/library" {
		t.Fatal("wire contract conflated states or lost provenance")
	}
	var fields map[string]json.RawMessage
	json.Unmarshal(encoded, &fields)
	for _, key := range []string{"run_id", "decision", "policy"} {
		if _, ok := fields[key]; ok {
			t.Fatalf("unexpected reuse coupling: %s", key)
		}
	}
}
