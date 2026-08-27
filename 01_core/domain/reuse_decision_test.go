/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/reuse_decision_test.md
 * @layer L1
 * @updated 2026-08-26
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReuseDecisionReportKeepsPolicyAndGatesObservable(t *testing.T) {
	report := ReuseDecisionReport{SchemaVersion: "1", Policy: ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}}, PolicyHash: "hash", Decision: DecisionUnknown, Matrix: []CapabilityComparison{}, Gaps: []string{}, Reasons: []DecisionReason{{RuleID: "DEC-INCONCLUSIVE-001", Effect: "block", Message: "missing evidence"}}, Alternatives: []DecisionAlternative{}, Effort: EffortAssessment{Level: EffortUnknown}}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"policy_hash", "matrix", "reasons", "effort", "INCONCLUSIVE"} {
		if !strings.Contains(string(data), field) {
			t.Fatalf("missing %s in %s", field, data)
		}
	}
}
