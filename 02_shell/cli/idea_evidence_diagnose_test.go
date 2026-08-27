/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_diagnose_test.md
 * @layer L2
 * @updated 2026-08-26
 */

package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestEvidenceDiagnoseOffline(t *testing.T) {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	report := domain.AssessmentReport{SchemaVersion: "1", RulesVersion: "test", Assessments: []domain.CandidateAssessment{{Repository: ref, Capabilities: []domain.CapabilityEvidence{{Capability: "static analysis", Status: domain.EvidenceInconclusive}}}}}
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true, Targets: []domain.CollectionTargetResult{{Target: domain.ExtractionTarget{Path: "analyzer.go", Capabilities: []string{"static analysis"}}, Status: domain.TargetCollected}}}
	judgments := []domain.EvidenceJudgment{{Repository: "org/repo", Capability: "static analysis", Expected: domain.EvidenceProven}}
	reportPath := writeEvidenceFixture(t, "assessment.json", report)
	bundlePath := writeEvidenceFixture(t, "bundles.json", []domain.EvidenceBundle{bundle})
	judgmentPath := writeEvidenceFixture(t, "judgments.json", judgments)
	var out bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "diagnose", "--assessment", reportPath, "--evidence", bundlePath, "--judgments", judgmentPath}, &out); err != nil {
		t.Fatal(err)
	}
	var got domain.EvidenceCoverageReport
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Metrics.ExpectedProvenCount != 1 || len(got.Diagnoses) != 1 || got.Diagnoses[0].PrimaryCause != domain.MissPattern {
		t.Fatalf("relatório inesperado: %+v", got)
	}
}

func TestEvidenceDiagnoseRequiresFrozenInputs(t *testing.T) {
	err := RunIdeaCommand([]string{"evidence", "diagnose"}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "--assessment") {
		t.Fatalf("erro = %v", err)
	}
}
