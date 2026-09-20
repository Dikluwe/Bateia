/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o9_refinement_test.md
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

func TestO9HistoricalProtectedManifestRemainsFrozenAfterO11Invalidation(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o9.json")
	if digestBytes(document) != "b7960417c3e9de2c12b42fa56b5f9b73a12bc87585541800b3bd92d1596f4b0c" {
		t.Fatal("historical O9 protected manifest changed")
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
	if err := json.Unmarshal(document, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v9" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O9" || len(manifest.Files) != 22 {
		t.Fatal("unexpected O9 protected manifest identity")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || filepath.IsAbs(file.Path) || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe O9 protected manifest entry %q", file.Path)
		}
		seen[file.Path] = true
	}
}

func TestO9StartRejectsInventedOrCrossCategoryReferences(t *testing.T) {
	request := validStartRequest()
	request["proposals"] = []any{validProposal(map[string]any{
		"product_observation_refs": []string{"mech_" + strings.Repeat("a", 64)},
	})}
	assertStartSuccessViolated(t, request)

	request = validStartRequest()
	request["proposals"] = []any{validProposal(map[string]any{
		"product_observation_refs": []string{"market-observation-1"},
	})}
	assertStartSuccessViolated(t, request)
}

func TestO9StartRejectsInternallyInconsistentSourceReferences(t *testing.T) {
	request := validStartRequest()
	sources := request["source_artifacts"].([]any)
	sources = append(sources, map[string]any{
		"run_id": runID('1'), "run_kind": "MARKET_DISCOVERY", "artifact": "market_snapshot", "sha256": strings.Repeat("2", 64),
	})
	request["source_artifacts"] = sources
	assertStartSuccessViolated(t, request)

	request = validStartRequest()
	sources = request["source_artifacts"].([]any)
	sources = append(sources, map[string]any{
		"run_id": runID('1'), "run_kind": "RECOMMENDATION", "artifact": "mechanism_gaps", "sha256": strings.Repeat("1", 64),
	})
	request["source_artifacts"] = sources
	assertStartSuccessViolated(t, request)
}

func TestO9CancelPreservesTerminalVersionAndSingleCASTransition(t *testing.T) {
	request := map[string]any{"run_id": runID('9'), "expected_version": 7}
	terminal := map[string]any{
		"run_id": runID('9'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "SUCCEEDED", "version": 7,
		"recommendation_status": nil, "market_discovery_status": nil, "improvement_status": "PARTIAL",
		"progress": map[string]any{"planned": 5, "completed": 5, "stage": "terminal"},
	}
	control := exchange(t, RunCancel, request, terminal, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("idempotent terminal cancel rejected: %+v", finding)
	}
	terminal["version"] = 9
	control = exchange(t, RunCancel, request, terminal, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("terminal cancel mutation survived: %+v", finding)
	}
}

func TestO9ArtifactRecognizesRealOpacityWithoutMaskingCanonicalFailure(t *testing.T) {
	contents := []byte(`{"schema_version":"1","artifact":"architecture_comparison","entries":[{"target_property":"latency","mechanism_id":"mech_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","source_run_id":"run_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","novelty":"UNKNOWN","transfer_mode":"UNKNOWN"}]}`)
	request := map[string]any{"run_id": runID('a'), "artifact": "architecture_comparison", "cursor": nil}
	returned := map[string]any{"artifact": "architecture_comparison", "sha256": digestBytes(contents), "size": len(contents), "canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil}
	control := exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("novelty/transfer shape alone was misclassified as opaque: %+v", finding)
	}

	noncanonical := []byte(`{"artifact":"architecture_comparison","schema_version":"1","entries":[{"target_property":"latency","mechanism_id":"mech_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","source_run_id":"run_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","novelty":"UNKNOWN"}]}`)
	returned["sha256"] = digestBytes(noncanonical)
	returned["size"] = len(noncanonical)
	returned["canonical_json"] = base64.StdEncoding.EncodeToString(noncanonical)
	control = exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("opacity masked noncanonical artifact bytes: %+v", finding)
	}
}

func assertStartSuccessViolated(t *testing.T, request map[string]any) {
	t.Helper()
	receipt := map[string]any{"schema_version": "1", "run_id": runID('c'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "QUEUED", "version": 1}
	control := exchange(t, SynthesisStart, request, receipt, nil, calls(len(request["source_artifacts"].([]any)), 1, 1, 1), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("invalid start reference survived: %+v", finding)
	}
}
