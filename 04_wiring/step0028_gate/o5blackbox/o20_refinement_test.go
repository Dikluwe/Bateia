/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o20_refinement_test.md
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

type o20E2EExchange struct {
	ContractInput    json.RawMessage        `json:"contract_input"`
	ArtifactDocument o19CalibrationDocument `json:"artifact_document"`
}

func TestO20PublicEndToEndReplayIsSevenUnknown(t *testing.T) {
	var corpus struct {
		SchemaVersion           string           `json:"schema_version"`
		ContainsCaseIDs         bool             `json:"contains_case_ids"`
		ContainsAttackHashes    bool             `json:"contains_attack_hashes"`
		ContainsMarkers         bool             `json:"contains_markers"`
		ContainsVerdictMetadata bool             `json:"contains_verdict_metadata"`
		Exchanges               []o20E2EExchange `json:"exchanges"`
	}
	raw := readPublicO5(t, "00_nucleo", "benchmark", "0028", "public-end-to-end-calibration-o19.json")
	if _, err := o6ParseUniqueJSON(raw); err != nil || json.Unmarshal(raw, &corpus) != nil || corpus.SchemaVersion != "tekt.step0028.public-e2e-calibration-source.v1" || corpus.ContainsCaseIDs || corpus.ContainsAttackHashes || corpus.ContainsMarkers || corpus.ContainsVerdictMetadata || len(corpus.Exchanges) != 7 {
		t.Fatal("invalid public end-to-end corpus")
	}
	unknown, preserved := 0, 0
	for i, item := range corpus.Exchanges {
		output := o20DocumentOutput(t, item.ArtifactDocument)
		switch Evaluate(item.ContractInput, output).Classification {
		case Unknown:
			unknown++
		case Preserved:
			preserved++
		default:
			t.Fatalf("public exchange %d was Violated", i)
		}
		neighbor := item.ArtifactDocument
		neighbor.Entries[0]["validation_plan"] = []any{map[string]any{"observable": "bounded observation"}}
		if finding := Evaluate(item.ContractInput, o20DocumentOutput(t, neighbor)); finding.Classification != Preserved {
			t.Fatalf("conformant neighbor %d not Preserved: %+v", i, finding)
		}
	}
	if unknown != 7 || preserved != 0 {
		t.Fatalf("public replay totals Unknown=%d Preserved=%d", unknown, preserved)
	}
}

func TestO20E2EOpacityDoesNotMaskViolation(t *testing.T) {
	var corpus struct {
		Exchanges []o20E2EExchange `json:"exchanges"`
	}
	if json.Unmarshal(readPublicO5(t, "00_nucleo", "benchmark", "0028", "public-end-to-end-calibration-o19.json"), &corpus) != nil || len(corpus.Exchanges) == 0 {
		t.Fatal("corpus")
	}
	var output ContractOutput
	if decodeStrict(o20DocumentOutput(t, corpus.Exchanges[0].ArtifactDocument), &output) != nil {
		t.Fatal("output")
	}
	var result o6ArtifactResult
	if decodeStrict(output.ReturnValue, &result) != nil {
		t.Fatal("result")
	}
	result.SHA256 = strings.Repeat("0", 64)
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(corpus.Exchanges[0].ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("opacity masked violation: %+v", finding)
	}
}

func TestO20PublicArtifactsAndHistoricalManifest(t *testing.T) {
	paths := [][]string{{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o20.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o19-to-o20.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o20.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o20.json"}}
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
	if json.Unmarshal(readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o20.json"), &manifest) != nil || manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v20" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O20" || len(manifest.Files) != 40 {
		t.Fatal("manifest identity")
	}
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o20.json")
	if digestBytes(document) != "ab21537deb3614c91323cb8531e224bc42739eb74a3f6392895f0a5508b631f1" {
		t.Fatal("O20 historical manifest changed")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || filepath.IsAbs(file.Path) || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe %q", file.Path)
		}
		seen[file.Path] = true
	}
}

func o20DocumentOutput(t *testing.T, document o19CalibrationDocument) []byte {
	t.Helper()
	contents := mustJSON(t, document)
	returned := map[string]any{"artifact": document.Artifact, "sha256": digestBytes(contents), "size": len(contents), "canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil}
	return mustJSON(t, ContractOutput{ReturnValue: mustJSON(t, returned), EffectObservation: calls(0, 0, 0, 0), PublicationObservation: emptyPublication()})
}
