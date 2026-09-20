/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o16_refinement_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestO16DecidablePartialEvidenceRemainsPreserved(t *testing.T) {
	rule := extractionCatalogMechanism{AcceptedEvidenceRules: []extractionCatalogEvidence{
		{RuleID: "implementation-rule", Role: "IMPLEMENTATION", WitnessKind: "IMPLEMENTATION"},
		{RuleID: "independent-rule", Role: "INDEPENDENT", WitnessKind: "TEST"},
	}}
	for _, witness := range []extractionWitness{
		{RuleID: "implementation-rule", Kind: "IMPLEMENTATION", AssessmentStatus: "PROVEN", ContentHash: strings.Repeat("a", 64)},
		{RuleID: "independent-rule", Kind: "TEST", AssessmentStatus: "PROVEN", ContentHash: strings.Repeat("b", 64)},
	} {
		group := &extractionExpectedGroup{Rule: rule, Status: "PARTIAL", Positive: []extractionWitness{witness}}
		if o15GroupHasSemanticWitnessLimit(group) {
			t.Fatalf("decidable PROVEN witness role became opaque: %s", witness.RuleID)
		}
	}
}

func TestO16ConformantHypothesisLimitsRemainPreserved(t *testing.T) {
	entries := []map[string]any{
		{"id": "h-target", "target_state": "UNKNOWN", "validation_plan": []any{map[string]any{"missing_evidence": "measure target"}}, "status": "PROPOSED"},
		{"id": "h-composition", "mechanism_refs": []any{"mech-a", "mech-b"}, "integration_unknowns": []any{"interface compatibility", "state model compatibility"}, "status": "PROPOSED", "validation_plan": []any{"validate"}},
		{"id": "h-license", "license": "UNRECOGNIZED", "transfer_mode": "PATTERN_ADAPTATION", "status": "PROPOSED"},
	}
	for _, entry := range entries {
		control := o15ArtifactExchange(t, "improvement_hypotheses", entry)
		if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
			t.Fatalf("conformant hypothesis limit became opaque: %+v", finding)
		}
	}
}

func TestO16OpaqueLicenseGapIsUnknown(t *testing.T) {
	control := o15ArtifactExchange(t, "hypothesis_gaps", map[string]any{
		"hypothesis_id": "h-license", "code": "LICENSE_NOT_RECOGNIZED", "source_ref": "license-assessment",
	})
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Unknown {
		t.Fatalf("typed opaque license gap was not Unknown: %+v", finding)
	}
}

func TestO16ViolationPrecedesOpaqueLicenseGap(t *testing.T) {
	control := o15ArtifactExchange(t, "hypothesis_gaps", map[string]any{
		"hypothesis_id": "h-license", "code": "LICENSE_NOT_RECOGNIZED", "source_ref": "license-assessment",
	})
	var output ContractOutput
	var result o6ArtifactResult
	if decodeStrict(control.ContractOutput, &output) != nil || decodeStrict(output.ReturnValue, &result) != nil {
		t.Fatal("could not decode license-gap control")
	}
	result.SHA256 = strings.Repeat("0", 64)
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("license opacity masked integrity violation: %+v", finding)
	}
}

func TestO16PublicJSONArtifactsAreStrictAndDuplicateFree(t *testing.T) {
	paths := [][]string{
		{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o16.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o15-to-o16.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o16.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o16.json"},
	}
	for _, path := range paths {
		if _, err := o6ParseUniqueJSON(readPublicO5(t, path...)); err != nil {
			t.Fatalf("O16 artifact %s is not strict duplicate-free JSON: %v", strings.Join(path, "/"), err)
		}
	}
}

func TestO16ProtectedManifestRemainsAnImmutableHistoricalAnchor(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o16.json")
	var manifest struct {
		SchemaVersion  string `json:"schema_version"`
		Visibility     string `json:"visibility"`
		OracleRevision string `json:"oracle_revision"`
		Files          []struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"files"`
	}
	if err := json.Unmarshal(document, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v16" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O16" || len(manifest.Files) != 32 {
		t.Fatal("unexpected O16 protected manifest identity")
	}
	if digestBytes(document) != "f038cea1a4be217a0452c85bcc0abe808982388a4ef3e2995263c419d5772da7" {
		t.Fatal("O16 protected manifest historical bytes changed")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || filepath.IsAbs(file.Path) || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe O16 manifest entry %q", file.Path)
		}
		seen[file.Path] = true
	}
}
