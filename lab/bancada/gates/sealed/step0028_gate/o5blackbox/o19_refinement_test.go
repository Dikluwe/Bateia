/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o19_refinement_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

type o19CalibrationDocument struct {
	SchemaVersion string           `json:"schema_version"`
	Artifact      string           `json:"artifact"`
	Entries       []map[string]any `json:"entries"`
}

func TestO19ConsumesEverySanitizedPublicProjection(t *testing.T) {
	var calibration struct {
		SchemaVersion          string `json:"schema_version"`
		BindingFree            bool   `json:"binding_free"`
		ContainsFixturePayload bool   `json:"contains_fixture_payload"`
		ContainsAttackPayload  bool   `json:"contains_attack_payload"`
		Projections            []struct {
			ArtifactDocument o19CalibrationDocument `json:"artifact_document"`
		} `json:"projections"`
	}
	raw := readPublicO5(t, "00_nucleo", "benchmark", "0028", "public-opacity-calibration-o18.json")
	_, uniqueErr := o6ParseUniqueJSON(raw)
	if uniqueErr != nil || json.Unmarshal(raw, &calibration) != nil || calibration.SchemaVersion != "tekt.step0028.public-opacity-calibration.v1" || !calibration.BindingFree || calibration.ContainsFixturePayload || calibration.ContainsAttackPayload || len(calibration.Projections) != 7 {
		t.Fatal("invalid public calibration")
	}
	for i, projection := range calibration.Projections {
		document := projection.ArtifactDocument
		if document.SchemaVersion != "1" || document.Artifact != "hypothesis_gaps" || len(document.Entries) != 1 {
			t.Fatalf("projection %d shape", i)
		}
		control := o19DocumentExchange(t, document)
		if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Unknown {
			t.Fatalf("projection %d not Unknown: %+v", i, finding)
		}
		delete(document.Entries[0], "source_ref")
		neighbor := o19DocumentExchange(t, document)
		if finding := Evaluate(neighbor.ContractInput, neighbor.ContractOutput); finding.Classification != Preserved {
			t.Fatalf("projection %d neighbor not Preserved: %+v", i, finding)
		}
	}
}

func TestO19ExactGapCodeDoesNotMaskViolation(t *testing.T) {
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

func TestO19PublicArtifactsAndHistoricalManifest(t *testing.T) {
	paths := [][]string{{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o19.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o18-to-o19.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o19.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o19.json"}}
	for _, path := range paths {
		if _, err := o6ParseUniqueJSON(readPublicO5(t, path...)); err != nil {
			t.Fatalf("strict JSON %s: %v", strings.Join(path, "/"), err)
		}
	}
	var manifest struct {
		SchemaVersion  string `json:"schema_version"`
		Visibility     string `json:"visibility"`
		OracleRevision string `json:"oracle_revision"`
		Files          []struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"files"`
	}
	if json.Unmarshal(readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o19.json"), &manifest) != nil || manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v19" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O19" || len(manifest.Files) != 38 {
		t.Fatal("manifest identity")
	}
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o19.json")
	if digestBytes(document) != "2b3ba97bb485bc5fd449354d6cdeb4feab9eb3ce314d533b760237a569e81157" {
		t.Fatal("O19 historical manifest changed")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || filepath.IsAbs(file.Path) || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe %q", file.Path)
		}
		seen[file.Path] = true
	}
}

func o19DocumentExchange(t *testing.T, document o19CalibrationDocument) Exchange {
	t.Helper()
	contents := mustJSON(t, document)
	request := map[string]any{"run_id": runID('9'), "artifact": document.Artifact, "cursor": nil}
	returned := map[string]any{"artifact": document.Artifact, "sha256": digestBytes(contents), "size": len(contents), "canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil}
	return exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
}
