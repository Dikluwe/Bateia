/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_verify_test.md
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

func TestEvidenceVerifyCanFillOnlyVerifiedLicenseGate(t *testing.T) {
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
	ref := domain.CommitRef{Owner: "o", Repo: "r", Commit: "abc123"}
	codeSource := func(rule string) domain.EvidenceSource {
		snippet := "func Core() {} // " + rule
		return domain.EvidenceSource{Kind: domain.EvidenceSymbol, Ref: ref, Path: "main.go", Symbol: "Core", Line: 1, Snippet: snippet, RuleID: rule, ContentHash: fmt.Sprintf("%x", sha256.Sum256([]byte(snippet))), Strength: 1}
	}
	meta := `{"archived":false}`
	maintenance := domain.EvidenceSource{Kind: domain.EvidenceRepository, Ref: ref, Path: "@repository-metadata", RuleID: "maintenance", ContentHash: fmt.Sprintf("%x", sha256.Sum256([]byte(meta))), Strength: 1}
	a := domain.CandidateAssessment{
		Repository:       ref,
		ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource("alignment")}},
		Capabilities:     []domain.CapabilityEvidence{{Capability: "core", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource("core")}}},
		License:          domain.LicenseAssessment{Status: domain.EvidenceInconclusive},
		Maintenance:      domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{maintenance}},
		Operation:        domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource("local")}},
	}
	idea := write("idea.json", domain.IdeaSpec{Problem: "x", Capabilities: []string{"core"}})
	assessment := write("assessment.json", domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{a}})
	policy := write("policy.json", domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, RequireLocal: true, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3})
	licenseContent := securityMITLicense
	digest := fmt.Sprintf("%x", sha256.Sum256([]byte(licenseContent)))
	item := domain.CollectedEvidence{Kind: domain.EvidenceLicense, Ref: ref, Path: "LICENSE", Content: licenseContent, ContentHash: digest}
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Plan: domain.EvidenceExtractionPlan{Repository: ref}, Items: []domain.CollectedEvidence{item}, Complete: true}
	batch := domain.TechnicalEvidenceBatch{SchemaVersion: "1", Findings: []domain.TechnicalEvidence{{
		ID: "license", Technique: domain.TechniqueLicense, Candidate: ref.Key(), Claim: "licença raiz MIT", Predicate: domain.TechnicalLicenseID, Value: "MIT", Result: domain.TechnicalSupports,
		Source: domain.TechnicalEvidenceSource{Ref: "LICENSE", Revision: ref.Commit, Digest: "sha256:" + digest}, Method: "verificar bytes da licença", ObservedAt: "2026-09-20T00:00:00Z",
	}}}
	batchPath := write("technical.json", batch)
	bundlePath := write("bundle.json", bundle)
	var out bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "verify", "--input", batchPath, "--bundles", bundlePath, "--as-of", "2026-09-24T00:00:00Z"}, &out); err != nil {
		t.Fatal(err)
	}
	var verified domain.TechnicalEvidenceReview
	if err := json.Unmarshal(out.Bytes(), &verified); err != nil || len(verified.Verifications) != 1 || verified.Verifications[0].Status != domain.TechnicalVerified {
		t.Fatalf("technical evidence not verified: %v %s", err, out.String())
	}
	base := []string{"decide", "--idea", idea, "--assessment", assessment, "--policy", policy}
	out.Reset()
	if err := RunIdeaCommand(base, &out); err != nil {
		t.Fatal(err)
	}
	var before domain.ReuseDecisionReport
	if err := json.Unmarshal(out.Bytes(), &before); err != nil || before.Decision != domain.DecisionUnknown {
		t.Fatalf("baseline should abstain: %v %s", err, out.String())
	}
	out.Reset()
	withEvidence := append(append([]string(nil), base...), "--technical-evidence", batchPath, "--technical-bundles", bundlePath, "--evidence-as-of", "2026-09-24T00:00:00Z")
	if err := RunIdeaCommand(withEvidence, &out); err != nil {
		t.Fatal(err)
	}
	var after domain.ReuseDecisionReport
	if err := json.Unmarshal(out.Bytes(), &after); err != nil || after.Decision != domain.DecisionFork || after.TechnicalEvidence == nil || after.TechnicalEvidence.DecisionImpact != "VERIFIED_LICENSE_GATE" {
		t.Fatalf("verified license did not affect gate: %v %s", err, out.String())
	}
	bundle.Items[0].Content += " tampered"
	tamperedPath := write("tampered.json", bundle)
	withEvidence[len(withEvidence)-3] = tamperedPath
	out.Reset()
	if err := RunIdeaCommand(withEvidence, &out); err != nil {
		t.Fatal(err)
	}
	var tampered domain.ReuseDecisionReport
	if err := json.Unmarshal(out.Bytes(), &tampered); err != nil || tampered.Decision != domain.DecisionUnknown || tampered.TechnicalEvidence.Verifications[0].Status == domain.TechnicalVerified {
		t.Fatalf("tampered bundle influenced decision: %v %s", err, out.String())
	}
}
