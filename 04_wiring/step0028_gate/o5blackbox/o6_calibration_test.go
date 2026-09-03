/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o6_calibration_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestO6StartAcceptsBothContractualAdmissionPaths(t *testing.T) {
	request := validStartRequest()
	request["source_artifacts"] = append(request["source_artifacts"].([]any), map[string]any{
		"run_id": runID('b'), "run_kind": "MARKET_DISCOVERY", "artifact": "market_snapshot", "sha256": strings.Repeat("b", 64),
	})
	receipt := map[string]any{"schema_version": "1", "run_id": runID('c'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "QUEUED", "version": 1}
	for name, effects := range map[string]EffectObservation{
		"idempotent reuse": calls(2, 0, 0, 0),
		"fresh admission":  calls(2, 1, 1, 1),
	} {
		t.Run(name, func(t *testing.T) {
			control := exchange(t, SynthesisStart, request, receipt, nil, effects, emptyPublication())
			if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
				t.Fatalf("contractual start path rejected: %+v", finding)
			}
		})
	}
}

func TestO6StartRejectsSemanticInvalidityBeforeAdmission(t *testing.T) {
	tests := map[string]func(map[string]any){
		"empty problem":      func(request map[string]any) { request["idea"].(map[string]any)["problem"] = "" },
		"missing capability": func(request map[string]any) { request["idea"].(map[string]any)["capabilities"] = []string{} },
		"invalid proposal identity": func(request map[string]any) {
			request["proposals"] = []any{validProposal(map[string]any{"proposal_id": "-invalid"})}
		},
		"invalid mechanism reference": func(request map[string]any) {
			request["proposals"] = []any{validProposal(map[string]any{"mechanism_refs": []string{"not-a-mechanism"}})}
		},
		"empty proposal title": func(request map[string]any) {
			request["proposals"] = []any{validProposal(map[string]any{"title": ""})}
		},
		"control text": func(request map[string]any) { request["idea"].(map[string]any)["problem"] = "unsafe\ntext" },
		"duplicate source": func(request map[string]any) {
			source := request["source_artifacts"].([]any)[0]
			request["source_artifacts"] = []any{source, source}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			request := validStartRequest()
			mutate(request)
			success := map[string]any{"schema_version": "1", "run_id": runID('d'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "QUEUED", "version": 1}
			control := exchange(t, SynthesisStart, request, success, nil, calls(1, 1, 1, 1), emptyPublication())
			if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
				t.Fatalf("semantic invalidity survived: %+v", finding)
			}
		})
	}
}

func TestO6RunLifecycleDoesNotEquateCompletionClaims(t *testing.T) {
	request := map[string]any{"run_id": runID('e')}
	degradedComplete := map[string]any{
		"run_id": runID('e'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "CANCELLED", "version": 2,
		"recommendation_status": nil, "market_discovery_status": nil, "improvement_status": "COMPLETE",
		"progress": map[string]any{"planned": 5, "completed": 2, "stage": "terminal"},
	}
	control := exchange(t, RunRead, request, degradedComplete, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("degraded lifecycle overclaimed COMPLETE: %+v", finding)
	}
}

func TestO6CancelAcceptsDegradedReturnWithoutInventedStatusAndForbidsPublication(t *testing.T) {
	request := map[string]any{"run_id": runID('f'), "expected_version": 3}
	returned := map[string]any{
		"run_id": runID('f'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "CANCELLED", "version": 3,
		"recommendation_status": nil, "market_discovery_status": nil, "improvement_status": nil,
		"progress": map[string]any{"planned": 5, "completed": 2, "stage": "terminal"},
	}
	control := exchange(t, RunCancel, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("degraded cancellation representation rejected: %+v", finding)
	}
	control = exchange(t, RunCancel, request, returned, nil, calls(0, 0, 0, 0), publication([]string{"improvement_hypotheses"}))
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("cancel publication survived: %+v", finding)
	}
}

func TestO6ArtifactFixedOrderIntegrityAndOpacity(t *testing.T) {
	control := referenceArtifactExchange(t)
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("fixed-order artifact rejected: %+v", finding)
	}
	var output ContractOutput
	if err := decodeStrict(control.ContractOutput, &output); err != nil {
		t.Fatal(err)
	}
	var result artifactResult
	if err := decodeStrict(output.ReturnValue, &result); err != nil {
		t.Fatal(err)
	}
	result.Size++
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("artifact integrity deviation survived: %+v", finding)
	}

	contents := []byte(`{"schema_version":"1","artifact":"hypothesis_gaps","entries":[{"target_state":"UNKNOWN","target_gap_refs":["baseline-not-observed"]}]}`)
	request := map[string]any{"run_id": runID('1'), "artifact": "hypothesis_gaps", "cursor": nil}
	returned := map[string]any{"artifact": "hypothesis_gaps", "sha256": digestBytes(contents), "size": len(contents), "canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil}
	control = exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Unknown {
		t.Fatalf("explicit artifact opacity was not Unknown: %+v", finding)
	}
}

func TestO6ReplaySeparatesFrozenSourcesFromProducedArtifacts(t *testing.T) {
	request := map[string]any{
		"catalog_sha256": rulesSHA256, "capture_sha256": strings.Repeat("2", 64),
		"artifact_hashes":  map[string]string{"source": strings.Repeat("3", 64)},
		"canonical_inputs": []string{base64.StdEncoding.EncodeToString([]byte(`{"schema_version":"1"}`))},
	}
	artifactHashes := map[string]string{}
	indexHashes := map[string]string{}
	for _, name := range o6ProducedArtifacts {
		artifactHashes[name] = strings.Repeat("4", 64)
		indexHashes[name] = strings.Repeat("5", 64)
	}
	returned := map[string]any{"artifact_hashes": artifactHashes, "index_hashes": indexHashes, "network_calls": 0}
	control := exchange(t, Replay, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("produced replay hashes were equated to frozen sources: %+v", finding)
	}

	opaque := exchange(t, Replay, request, nil, &ControlledError{Code: "INTEGRITY_ERROR", Message: "INTEGRITY_ERROR"}, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(opaque.ContractInput, opaque.ContractOutput); finding.Classification != Unknown {
		t.Fatalf("unobservable replay integrity cause was not Unknown: %+v", finding)
	}
}

func TestO6RefinementIncompleteEvidenceIsUnknownButCannotOverclaim(t *testing.T) {
	request := map[string]any{
		"receipt_hashes":       map[string]string{"0025": receipt0025Digest, "0026": "", "0027": receipt0027Digest},
		"runtime_instance_ids": map[string]string{"store": "store-1", "queue": "", "root": "root-1", "cursor_secret": "secret-1", "worker": "worker-1"},
		"layer_findings":       []string{"L1_PURE", "L4_WIRING_ONLY"},
	}
	result := map[string]any{
		"receipts_verified": map[string]bool{"0025": true, "0026": false, "0027": true},
		"runtime_unique":    false, "layer_findings": []string{"L1_PURE", "L4_WIRING_ONLY"},
	}
	control := exchange(t, RefinementAudit, request, result, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Unknown {
		t.Fatalf("incomplete refinement evidence was not Unknown: %+v", finding)
	}
	result["runtime_unique"] = true
	control = exchange(t, RefinementAudit, request, result, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("unsupported runtime uniqueness overclaim survived: %+v", finding)
	}
}

func TestO6ViolationPrecedesSimultaneousOpacity(t *testing.T) {
	control := unknownExtractionExchange(t)
	var output ContractOutput
	if err := decodeStrict(control.ContractOutput, &output); err != nil {
		t.Fatal(err)
	}
	output.EffectObservation.NetworkCalls = 1
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("opacity masked a decidable violation: %+v", finding)
	}
}

func validProposal(overrides map[string]any) map[string]any {
	proposal := map[string]any{
		"proposal_id": "proposal-1", "title": "bounded recovery", "target_property": "recovery latency",
		"expected_effect": "lower bounded latency", "mechanism_refs": []string{}, "product_observation_refs": []string{},
		"assumptions": []string{}, "derivation": "EXTERNAL_AGENT_PROPOSAL", "rationale": "validate the causal claim",
	}
	for key, value := range overrides {
		proposal[key] = value
	}
	return proposal
}

func TestO6ProductionRemainsBindingFree(t *testing.T) {
	for _, name := range []string{"api.go", "json.go", "evaluator.go", "extraction.go", "ports.go", "o6_extraction_artifact.go", "o6_runtime.go"} {
		source := string(readPublicO5(t, "04_wiring", "step0028_gate", "o5blackbox", name))
		for _, forbidden := range []string{"case_id", "attack_entry_sha256", "semantic_variant_id", "entry_commitment_sha256", "expected_result", "expected_outcome", "marker"} {
			if strings.Contains(strings.ToLower(source), forbidden) {
				t.Fatalf("%s contains forbidden evaluator dependency %q", name, forbidden)
			}
		}
	}
}

func TestO6HistoricalProtectedManifestRemainsFrozenAfterO8Invalidation(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o6.json")
	if digestBytes(document) != "f48e81288a30ca220169ab4b733af1166be0f87d436cc7695302f416f997ff5e" {
		t.Fatal("historical O6 protected manifest changed")
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
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v6" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O6" || len(manifest.Files) != 18 {
		t.Fatal("unexpected O6 protected manifest identity")
	}
}

func TestO6StrictJSONStillRejectsUnknownEnvelopeFields(t *testing.T) {
	control := referenceExchanges(t)[6]
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(control.ContractInput, &envelope); err != nil {
		t.Fatal(err)
	}
	envelope["binding"] = json.RawMessage(`"forbidden"`)
	if finding := Evaluate(mustJSON(t, envelope), control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("unknown evaluator input field survived: %+v", finding)
	}
}
