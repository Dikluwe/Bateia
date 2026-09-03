/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o18_refinement_test.md
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

func TestO18SevenPublishedRelationsAreNecessaryAndSufficient(t *testing.T) {
	rows := []struct {
		name, artifact   string
		opaque, neighbor map[string]any
	}{
		{"composition", "improvement_hypotheses",
			map[string]any{"id": "h", "mechanism_refs": []any{"a", "b"}, "integration_unknowns": []any{"u"}, "status": "INCONCLUSIVE", "validation_plan": []any{"v"}},
			map[string]any{"id": "h", "mechanism_refs": []any{"a", "b"}, "integration_unknowns": []any{"u"}, "status": "INCONCLUSIVE", "validation_plan": []any{}}},
		{"license", "hypothesis_gaps",
			map[string]any{"hypothesis_id": "h", "code": "LICENSE_NOT_RECOGNIZED", "source_ref": "s"},
			map[string]any{"hypothesis_id": "h", "code": "LICENSE_RECOGNIZED", "source_ref": "s"}},
		{"counterevidence", "improvement_hypotheses",
			map[string]any{"id": "h", "counterevidence_search": map[string]any{"status": "BUDGET_EXHAUSTED"}, "status": "INCONCLUSIVE", "validation_plan": []any{"v"}},
			map[string]any{"id": "h", "counterevidence_search": map[string]any{"status": "NONE_FOUND_IN_SCOPE"}, "status": "INCONCLUSIVE", "validation_plan": []any{"v"}}},
		{"identity", "mechanism_gaps",
			map[string]any{"code": "IDENTITY_BINDING_AMBIGUOUS"}, map[string]any{"code": "IDENTITY_BINDING_PRESENT"}},
		{"provider", "hypothesis_gaps",
			map[string]any{"hypothesis_id": "h", "code": "PROVIDER_CAPTURE_NOT_EVIDENCED", "source_ref": "s"},
			map[string]any{"hypothesis_id": "h", "code": "PROVIDER_CAPTURE_EVIDENCED", "source_ref": "s"}},
		{"utility", "hypothesis_gaps",
			map[string]any{"hypothesis_id": "h", "code": "UTILITY_COST_WITHOUT_BASELINE", "source_ref": "s"},
			map[string]any{"hypothesis_id": "h", "code": "UTILITY_COST_BASELINE_PRESENT", "source_ref": "s"}},
		{"security", "hypothesis_gaps",
			map[string]any{"hypothesis_id": "h", "code": "SECURITY_EFFECT_BASELINE_MISSING", "source_ref": "s"},
			map[string]any{"hypothesis_id": "h", "code": "SECURITY_EFFECT_BASELINE_PRESENT", "source_ref": "s"}},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			if finding := Evaluate(o15ArtifactExchange(t, row.artifact, row.opaque).ContractInput, o15ArtifactExchange(t, row.artifact, row.opaque).ContractOutput); finding.Classification != Unknown {
				t.Fatalf("published relation not Unknown: %+v", finding)
			}
			control := o15ArtifactExchange(t, row.artifact, row.neighbor)
			if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
				t.Fatalf("one-requirement neighbor not Preserved: %+v", finding)
			}
		})
	}
}

func TestO18ViolationPrecedesPublishedRelation(t *testing.T) {
	control := o15ArtifactExchange(t, "hypothesis_gaps", map[string]any{"hypothesis_id": "h", "code": "LICENSE_NOT_RECOGNIZED", "source_ref": "s"})
	var output ContractOutput
	var result o6ArtifactResult
	if decodeStrict(control.ContractOutput, &output) != nil || decodeStrict(output.ReturnValue, &result) != nil {
		t.Fatal("decode")
	}
	result.SHA256 = strings.Repeat("0", 64)
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("opacity masked violation: %+v", finding)
	}
}

func TestO18PublicArtifactsAndHistoricalManifest(t *testing.T) {
	paths := [][]string{{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o18.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o17-to-o18.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o18.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o18.json"}}
	for _, path := range paths {
		if _, err := o6ParseUniqueJSON(readPublicO5(t, path...)); err != nil {
			t.Fatalf("strict JSON %s: %v", strings.Join(path, "/"), err)
		}
	}
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o18.json")
	var manifest struct {
		SchemaVersion  string `json:"schema_version"`
		Visibility     string `json:"visibility"`
		OracleRevision string `json:"oracle_revision"`
		Files          []struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"files"`
	}
	if json.Unmarshal(document, &manifest) != nil || manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v18" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O18" || len(manifest.Files) != 36 {
		t.Fatal("manifest identity")
	}
	if digestBytes(document) != "2ab4d381ca15b6ac4961ee0fc4211a5473c2e4ca9bf4e9b427636cfaf50bdc5a" {
		t.Fatal("O18 historical manifest changed")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || filepath.IsAbs(file.Path) || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe manifest entry %q", file.Path)
		}
		seen[file.Path] = true
	}
}
