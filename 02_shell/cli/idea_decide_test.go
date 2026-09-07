/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_decide_test.md
 * @layer L2
 * @updated 2026-09-07
 */

package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func TestRunIdeaDecideIsOfflineAndWritesReport(t *testing.T) {
	dir := t.TempDir()
	idea := domain.IdeaSpec{Problem: "x", Capabilities: []string{"core"}}
	src := domain.EvidenceSource{Kind: domain.EvidenceFile, Ref: domain.CommitRef{Owner: "o", Repo: "r", Commit: "abc"}, Path: "main.go", RuleID: "cap", Strength: 1}
	a := domain.CandidateAssessment{Repository: src.Ref, ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceProven}, Capabilities: []domain.CapabilityEvidence{{Capability: "core", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{src}}}, License: domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT", Sources: []domain.EvidenceSource{src}}, Operation: domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{src}}}
	assessment := domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{a}}
	policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, RequireLocal: true, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
	write := func(name string, v any) string {
		p := filepath.Join(dir, name)
		data, _ := json.Marshal(v)
		if err := os.WriteFile(p, data, 0600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	t.Setenv("HOME", dir)
	t.Setenv("HTTP_PROXY", "http://127.0.0.1:1")
	var out bytes.Buffer
	err := RunIdeaCommand([]string{"decide", "--idea", write("idea.json", idea), "--assessment", write("assessment.json", assessment), "--policy", write("policy.json", policy)}, &out)
	if err != nil {
		t.Fatal(err)
	}
	var got domain.ReuseDecisionReport
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Decision != domain.DecisionUnknown || got.PolicyHash == "" {
		t.Fatalf("got %#v", got)
	}
}
func TestRunIdeaDecideRequiresAllInputs(t *testing.T) {
	if err := RunIdeaCommand([]string{"decide"}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected required argument error")
	}
}
