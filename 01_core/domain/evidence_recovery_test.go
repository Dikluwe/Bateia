/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_recovery_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRecoveryContractsKeepCauseFindingAndPolicySeparate(t *testing.T) {
	p := EvidenceRecoveryPlan{SchemaVersion: "1", Round: 1, DiagnosisHash: "d", PolicyHash: "p", Policy: RecoveryPolicy{SchemaVersion: "1", MaxRounds: 2}, Targets: []RecoveryTarget{{FindingID: "f", Candidate: CommitRef{Owner: "o", Repo: "r", Commit: "c"}, Kind: EvidenceTest, RuleID: "rule", Expected: []EvidenceStatus{EvidenceProven}}}, Skipped: []InconclusiveFinding{}}
	b, e := json.Marshal(p)
	if e != nil {
		t.Fatal(e)
	}
	for _, x := range []string{"finding_id", "diagnosis_hash", "policy_hash", "max_rounds"} {
		if !strings.Contains(string(b), x) {
			t.Fatalf("missing %s", x)
		}
	}
}
