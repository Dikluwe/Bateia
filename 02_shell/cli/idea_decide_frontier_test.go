/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_decide_frontier_test.md
 * @layer L2
 * @updated 2026-09-24
 */

package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func TestRunIdeaDecideFrontierWritesStrictOfflineReport(t *testing.T) {
	dir := t.TempDir()
	ref := domain.CommitRef{Owner: "acme", Repo: "tool", Commit: "abc123"}
	sourceFor := func(kind domain.EvidenceKind, path, snippet string) domain.EvidenceSource {
		source := domain.EvidenceSource{Kind: kind, Ref: ref, Path: path, RuleID: "CLI-MCDA", ContentHash: fmt.Sprintf("%x", sha256.Sum256([]byte(snippet))), Strength: 1}
		if kind != domain.EvidenceRepository {
			source.Line = 1
			source.Snippet = snippet
		}
		return source
	}
	codeSource := sourceFor(domain.EvidenceFile, "evidence.go", "func Core() {}")
	licenseSource := sourceFor(domain.EvidenceLicense, "LICENSE", "MIT license")
	maintenanceSource := sourceFor(domain.EvidenceRepository, "@repository-metadata", "active repository")
	assessment := domain.CandidateAssessment{
		Repository:       ref,
		ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource}},
		Capabilities:     []domain.CapabilityEvidence{{Capability: "core", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource}}},
		License:          domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT", Sources: []domain.EvidenceSource{licenseSource}},
		Maintenance:      domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{maintenanceSource}},
		Operation:        domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource}},
	}
	request := domain.MulticriteriaDecisionRequest{
		Idea:         domain.IdeaSpec{Problem: "x", Capabilities: []string{"core"}},
		Assessment:   domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{assessment}},
		Policy:       domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, RequireLocal: true, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3},
		Criteria:     []domain.MulticriteriaCriterion{{ID: "coverage", Label: "Coverage", Direction: domain.CriterionMaximize, ScaleMin: 0, ScaleMax: 1}},
		Measurements: []domain.CriterionMeasurement{{Candidate: ref.Key(), CriterionID: "coverage", Status: domain.EvidenceProven, Interval: domain.CriterionInterval{Low: 1, High: 1}, Sources: []domain.EvidenceSource{codeSource}}},
		Scenarios:    []domain.PreferenceScenario{},
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(dir, "request.json")
	if err := os.WriteFile(input, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"decide", "frontier", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var report domain.MulticriteriaDecisionReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.StopReason != "SINGLE_PARETO_CANDIDATE" || len(report.Recommended) != 1 || report.Recommended[0] != ref.Key() {
		t.Fatalf("relatório inesperado: %+v", report)
	}
}

func TestRunIdeaDecideFrontierRequiresInput(t *testing.T) {
	if err := RunIdeaCommand([]string{"decide", "frontier"}, &bytes.Buffer{}); err == nil {
		t.Fatal("esperava erro de input obrigatório")
	}
}
