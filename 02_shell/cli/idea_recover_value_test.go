/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recover_value_test.md
 * @layer L2
 * @updated 2026-09-24
 */

package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func TestIdeaRecoverValueProducesOfflineCounterfactualPlan(t *testing.T) {
	ref := domain.CommitRef{Owner: "acme", Repo: "tool", Commit: "abc"}
	request := domain.DecisionValueRequest{
		Idea: domain.IdeaSpec{Problem: "index", Capabilities: []string{"index"}},
		Assessment: domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{{
			Repository: ref, Capabilities: []domain.CapabilityEvidence{{Capability: "index", Status: domain.EvidenceInconclusive}},
			License: domain.LicenseAssessment{Status: domain.EvidenceInconclusive}, Maintenance: domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive},
			Operation: domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive},
		}}},
		DecisionPolicy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .8, IntegrationMinValue: .5},
		Diagnosis:      domain.InconclusiveDiagnosis{SchemaVersion: "1", Findings: []domain.InconclusiveFinding{{ID: "license", Cause: domain.CauseUnknownLicense, Candidate: ref.Key(), RuleID: "GATE-LICENSE-001", Recoverable: true}}},
		RecoveryPolicy: usecase.DefaultRecoveryPolicy(), Round: 1, PreviousTargets: []domain.RecoveryTarget{},
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(path, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := RunIdeaRecoverCommand([]string{"value", "--input", path}, &stdout, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"method": "COUNTERFACTUAL_DECISION_REACHABILITY_V1"`) || strings.Contains(stdout.String(), `"probability"`) {
		t.Fatalf("resultado inesperado: %s", stdout.String())
	}
}

func TestIdeaRecoverValueRejectsUnknownFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.json")
	if err := os.WriteFile(path, []byte(`{"round":1,"invented_probability":0.9}`), 0600); err != nil {
		t.Fatal(err)
	}
	err := RunIdeaRecoverCommand([]string{"value", "--input", path}, &bytes.Buffer{}, nil)
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("campo desconhecido deveria ser rejeitado: %v", err)
	}
}
