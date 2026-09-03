/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o11_refinement_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestO11HistoricalProtectedManifestRemainsFrozenAfterO13Invalidation(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o11.json")
	if digestBytes(document) != "0c7766e895338fb1cdba077e04b29531dd920819247e9d6fbb1179e64be20a7a" {
		t.Fatal("historical O11 protected manifest changed")
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
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v11" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O11" || len(manifest.Files) != 24 {
		t.Fatal("unexpected O11 protected manifest identity")
	}
}

func TestO11StartRejectsInventedEvidenceObservation(t *testing.T) {
	request := validStartRequest()
	receipt := map[string]any{
		"schema_version": "1", "run_id": runID('c'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "QUEUED", "version": 1,
	}
	published := emptyPublication()
	published.LayerFindings = []string{"proposal supplied evidence"}
	control := exchange(t, SynthesisStart, request, receipt, nil, calls(1, 1, 1, 1), published)
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("synthesis start evidence invention survived: %+v", finding)
	}
}

func TestO11StartRejectsSourceRunIdentityCollision(t *testing.T) {
	request := validStartRequest()
	receipt := map[string]any{
		"schema_version": "1", "run_id": runID('1'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "QUEUED", "version": 1,
	}
	control := exchange(t, SynthesisStart, request, receipt, nil, calls(1, 1, 1, 1), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("source/new-run identity collision survived: %+v", finding)
	}
}

func TestO11MechanismOpacitySurvivesAlongsideDefiniteMechanism(t *testing.T) {
	definite := referenceExtractionExchange(t)
	opaque := unknownExtractionExchange(t)

	var definiteInput, opaqueInput ContractInput
	var definiteRequest, opaqueRequest extractionRequest
	var definiteOutput, opaqueOutput ContractOutput
	var definiteResult, opaqueResult extractionResult
	if decodeStrict(definite.ContractInput, &definiteInput) != nil || decodeStrict(opaque.ContractInput, &opaqueInput) != nil ||
		decodeStrict(definiteInput.Request, &definiteRequest) != nil || decodeStrict(opaqueInput.Request, &opaqueRequest) != nil ||
		decodeStrict(definite.ContractOutput, &definiteOutput) != nil || decodeStrict(opaque.ContractOutput, &opaqueOutput) != nil ||
		decodeStrict(definiteOutput.ReturnValue, &definiteResult) != nil || decodeStrict(opaqueOutput.ReturnValue, &opaqueResult) != nil {
		t.Fatal("could not compose independent extraction controls")
	}
	definiteRequest.Bundle.Items = append(definiteRequest.Bundle.Items, opaqueRequest.Bundle.Items...)
	definiteRequest.Assessments = append(definiteRequest.Assessments, opaqueRequest.Assessments...)
	definiteResult.Gaps = append(definiteResult.Gaps, opaqueResult.Gaps...)
	definiteInput.Request = mustJSON(t, definiteRequest)
	definiteOutput.ReturnValue = mustJSON(t, definiteResult)
	if finding := Evaluate(mustJSON(t, definiteInput), mustJSON(t, definiteOutput)); finding.Classification != Unknown {
		t.Fatalf("mixed resolved/unresolved extraction was not Unknown: %+v", finding)
	}
}

func TestO11ArtifactOpacityRequiresConcreteEpistemicUncertainty(t *testing.T) {
	contents := []byte(`{"schema_version":"1","artifact":"architecture_comparison","entries":[{"target_property":"availability","mechanism_id":"mech_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","source_run_id":"run_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","novelty":"UNKNOWN","transfer_mode":"UNKNOWN"},{"target_property":"latency","mechanism_id":"mech_cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc","source_run_id":"run_dddddddddddddddddddddddddddddddd","target_state":"UNKNOWN","target_gap_refs":["baseline-not-observed"],"integration_unknowns":["state model compatibility unmeasured"]}]}`)
	request := map[string]any{"run_id": runID('a'), "artifact": "architecture_comparison", "cursor": nil}
	returned := map[string]any{
		"artifact": "architecture_comparison", "sha256": digestBytes(contents), "size": len(contents),
		"canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil,
	}
	control := exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Unknown {
		t.Fatalf("mixed artifact with a concrete unresolved baseline was not Unknown: %+v", finding)
	}
}
