/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o17_refinement_test.md
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

func TestO17TenPublicOpacityStatesAndNeighborControls(t *testing.T) {
	tests := []struct {
		name, artifact string
		entry          map[string]any
		opaque         bool
	}{
		{"unknown catalog class", "mechanism_gaps", map[string]any{"code": "MECHANISM_CLASS_UNKNOWN"}, true},
		{"partial technical assessment", "mechanism_gaps", map[string]any{"code": "MINIMUM_WITNESS_INSUFFICIENT"}, true},
		{"unknown target with inconclusive baseline", "improvement_hypotheses", map[string]any{"id": "h-target", "target_state": "UNKNOWN", "status": "INCONCLUSIVE"}, true},
		{"unmeasured composition", "improvement_hypotheses", map[string]any{"id": "h-compose", "mechanism_refs": []any{"m-a", "m-b"}, "integration_unknowns": []any{"compatibility"}, "status": "INCONCLUSIVE", "validation_plan": []any{"measure"}}, true},
		{"unresolved license gap", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-license", "code": "LICENSE_NOT_RECOGNIZED", "source_ref": "s-license"}, true},
		{"exhausted counterevidence", "improvement_hypotheses", map[string]any{"id": "h-budget", "counterevidence_search": map[string]any{"status": "BUDGET_EXHAUSTED"}, "status": "INCONCLUSIVE", "validation_plan": []any{"search"}}, true},
		{"ambiguous immutable identity", "mechanism_gaps", map[string]any{"code": "IMMUTABLE_IDENTITY_AMBIGUOUS"}, true},
		{"provider without frozen capture", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-provider", "code": "PROVIDER_CAPTURE_NOT_EVIDENCED", "source_ref": "s-provider"}, true},
		{"utility without baseline", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-utility", "code": "UTILITY_COST_WITHOUT_BASELINE", "source_ref": "s-utility"}, true},
		{"security without baseline", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-security", "code": "SECURITY_EFFECT_WITHOUT_BASELINE", "source_ref": "s-security"}, true},
		{"unknown target proposed", "improvement_hypotheses", map[string]any{"id": "h-target-control", "target_state": "UNKNOWN", "status": "PROPOSED"}, false},
		{"composition proposed", "improvement_hypotheses", map[string]any{"id": "h-compose-control", "mechanism_refs": []any{"m-a", "m-b"}, "integration_unknowns": []any{"compatibility"}, "status": "PROPOSED", "validation_plan": []any{"validate"}}, false},
		{"pattern adaptation", "improvement_hypotheses", map[string]any{"id": "h-license-control", "license": "UNRECOGNIZED", "transfer_mode": "PATTERN_ADAPTATION", "status": "PROPOSED"}, false},
		{"noncontrolling prose", "improvement_hypotheses", map[string]any{"id": "h-text-control", "title": "baseline unavailable provider unknown", "status": "PROPOSED"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			control := o15ArtifactExchange(t, tc.artifact, tc.entry)
			want := Preserved
			if tc.opaque {
				want = Unknown
			}
			if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != want {
				t.Fatalf("classification = %s, want %s: %+v", finding.Classification, want, finding)
			}
		})
	}
}

func TestO17ViolationPrecedesEveryTypedOpacityState(t *testing.T) {
	control := o15ArtifactExchange(t, "hypothesis_gaps", map[string]any{"hypothesis_id": "h", "code": "UTILITY_COST_WITHOUT_BASELINE", "source_ref": "s"})
	var output ContractOutput
	var result o6ArtifactResult
	if decodeStrict(control.ContractOutput, &output) != nil || decodeStrict(output.ReturnValue, &result) != nil {
		t.Fatal("could not decode typed opacity control")
	}
	result.SHA256 = strings.Repeat("0", 64)
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("typed opacity masked violation: %+v", finding)
	}
}

func TestO17PublicJSONArtifactsAreStrictAndDuplicateFree(t *testing.T) {
	paths := [][]string{
		{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o17.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o16-to-o17.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o17.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o17.json"},
	}
	for _, path := range paths {
		if _, err := o6ParseUniqueJSON(readPublicO5(t, path...)); err != nil {
			t.Fatalf("O17 artifact %s is not strict duplicate-free JSON: %v", strings.Join(path, "/"), err)
		}
	}
}

func TestO17ProtectedManifestRemainsAnImmutableHistoricalAnchor(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o17.json")
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
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v17" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O17" || len(manifest.Files) != 34 {
		t.Fatal("unexpected O17 protected manifest identity")
	}
	if digestBytes(document) != "736019da8a8b087bedfbae7d92b8c13fd83ac2aa69a6e50489abdcdca4686090" {
		t.Fatal("O17 protected manifest historical bytes changed")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || filepath.IsAbs(file.Path) || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe O17 manifest entry %q", file.Path)
		}
		seen[file.Path] = true
	}
}
