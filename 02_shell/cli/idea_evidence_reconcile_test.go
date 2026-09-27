/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_reconcile_test.md
 * @layer L2
 * @updated 2026-09-24
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

func TestEvidenceReconcileAndOptionalDecisionAttachment(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, value any) string {
		path := filepath.Join(dir, name)
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	batch := domain.TechnicalEvidenceBatch{SchemaVersion: "1", Findings: []domain.TechnicalEvidence{{
		ID: "lic-1", Technique: domain.TechniqueLicense, Candidate: "o/r", Claim: "licença compatível", Result: domain.TechnicalInconclusive,
		Source: domain.TechnicalEvidenceSource{Ref: "LICENSE", Revision: "abc"}, Method: "inspeção", ObservedAt: "2026-09-20T00:00:00Z",
	}}}
	path := write("technical.json", batch)
	var out bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "reconcile", "--input", path, "--as-of", "2026-09-24T00:00:00Z"}, &out); err != nil {
		t.Fatal(err)
	}
	var reviewed domain.TechnicalEvidenceReview
	if err := json.Unmarshal(out.Bytes(), &reviewed); err != nil || len(reviewed.Claims) != 1 {
		t.Fatalf("invalid reconciliation: %v %s", err, out.String())
	}
	idea := write("idea.json", domain.IdeaSpec{Problem: "x", Capabilities: []string{"core"}})
	assessment := write("assessment.json", domain.AssessmentReport{SchemaVersion: "1"})
	policy := write("policy.json", domain.ReuseDecisionPolicy{SchemaVersion: "1", MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3})
	args := []string{"decide", "--idea", idea, "--assessment", assessment, "--policy", policy}
	out.Reset()
	if err := RunIdeaCommand(args, &out); err != nil {
		t.Fatal(err)
	}
	var baseline domain.ReuseDecisionReport
	if err := json.Unmarshal(out.Bytes(), &baseline); err != nil {
		t.Fatal(err)
	}
	if baseline.TechnicalEvidence != nil {
		t.Fatal("simple decision unexpectedly has specialized evidence")
	}
	out.Reset()
	if err := RunIdeaCommand(append(args, "--technical-evidence", path, "--evidence-as-of", "2026-09-24T00:00:00Z"), &out); err != nil {
		t.Fatal(err)
	}
	var advanced domain.ReuseDecisionReport
	if err := json.Unmarshal(out.Bytes(), &advanced); err != nil {
		t.Fatal(err)
	}
	if advanced.Decision != baseline.Decision || advanced.TechnicalEvidence == nil || advanced.TechnicalEvidence.DecisionImpact != "ADVISORY_ONLY" {
		t.Fatalf("technical evidence changed deterministic decision or was omitted: %+v", advanced)
	}
}
