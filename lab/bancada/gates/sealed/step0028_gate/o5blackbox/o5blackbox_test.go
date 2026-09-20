/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o5blackbox_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	protocolSHA256 = "c86b04d028bd6c0f33ddcff78ca323f70ab2b72d96bc86a6046b2ec6598755bb"
	rulesSHA256    = "52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142"
)

func TestO5BindingFreeReferenceControlsCoverEveryPort(t *testing.T) {
	protocol := readPublicO5(t, "00_nucleo", "benchmark", "0028", "executable-mutant-protocol-o8.json")
	if digestBytes(protocol) != protocolSHA256 {
		t.Fatal("public O5 protocol hash changed")
	}
	for _, exchange := range referenceExchanges(t) {
		finding := Evaluate(exchange.ContractInput, exchange.ContractOutput)
		if finding.Classification != Preserved {
			t.Fatalf("%s reference control: %+v", finding.Port, finding)
		}
		if finding.Predicates == 0 {
			t.Fatalf("%s evaluated no contract predicate", finding.Port)
		}
	}
}

func TestO5DerivesViolationsAndConcreteUnknownsWithoutBindings(t *testing.T) {
	controls := referenceExchanges(t)
	var output ContractOutput
	if err := decodeStrict(controls[0].ContractOutput, &output); err != nil {
		t.Fatal(err)
	}
	output.EffectObservation.NetworkCalls = 1
	if finding := Evaluate(controls[0].ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("external effect did not violate contract: %+v", finding)
	}

	unknown := unknownExtractionExchange(t)
	if finding := Evaluate(unknown.ContractInput, unknown.ContractOutput); finding.Classification != Unknown {
		t.Fatalf("catalog opacity was not preserved as Unknown: %+v", finding)
	}
	unknown = unavailableSourceExchange(t)
	if finding := Evaluate(unknown.ContractInput, unknown.ContractOutput); finding.Classification != Unknown {
		t.Fatalf("historical source opacity was not preserved as Unknown: %+v", finding)
	}
}

func TestO5SequenceRequiresCanonicalReplayAndStableRuntimeIdentity(t *testing.T) {
	controls := referenceExchanges(t)
	tools := controls[6]
	var changed ContractOutput
	if err := decodeStrict(tools.ContractOutput, &changed); err != nil {
		t.Fatal(err)
	}
	var returned map[string]any
	if err := json.Unmarshal(changed.ReturnValue, &returned); err != nil {
		t.Fatal(err)
	}
	returnedTools := returned["tools"].([]any)
	returnedTools[0].(map[string]any)["description"] = "different but schema-valid"
	changed.ReturnValue = mustJSON(t, returned)
	findings := EvaluateSequence([]Exchange{tools, {ContractInput: tools.ContractInput, ContractOutput: mustJSON(t, changed)}})
	if findings[0].Classification != Violated || findings[1].Classification != Violated {
		t.Fatalf("non-deterministic replay survived: %+v", findings)
	}
}

func TestO5ProductionSurfaceContainsNoRecordBindingOrExpectation(t *testing.T) {
	for _, name := range []string{"api.go", "json.go", "evaluator.go", "extraction.go", "ports.go"} {
		source := string(readPublicO5(t, "04_wiring", "step0028_gate", "o5blackbox", name))
		for _, forbidden := range []string{"attack_entry_sha256", "semantic_variant_id", "entry_commitment_sha256", "semantic_fault", "opacity_reason", "record.id", "Case.Expected", "ExpectedOutcome", "M001", "P01", "O01"} {
			if strings.Contains(source, forbidden) {
				t.Fatalf("%s contains forbidden binding/expectation token %q", name, forbidden)
			}
		}
	}
}

func TestO5HistoricalProtectedManifestRemainsFrozenAfterInvalidation(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o5.json")
	if digestBytes(document) != "07337d19aa36a4e85bd4a5933e167896769e0fbea56824cfc42d10a3b6e2956b" {
		t.Fatal("historical O5 protected manifest changed")
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
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v5" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O5" || len(manifest.Files) != 12 {
		t.Fatal("unexpected O5 protected manifest identity")
	}
}

func referenceExchanges(t *testing.T) []Exchange {
	t.Helper()
	return []Exchange{
		referenceExtractionExchange(t),
		exchange(t, SynthesisStart, validStartRequest(), map[string]any{"schema_version": "1", "run_id": runID('2'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "QUEUED", "version": 1}, nil, calls(1, 1, 1, 1), emptyPublication()),
		exchange(t, SynthesisExecute, map[string]any{"run_id": runID('3')}, nil, nil, calls(0, 0, 0, 0), synthesisPublication()),
		exchange(t, RunRead, map[string]any{"run_id": runID('4')}, map[string]any{"run_id": runID('4'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "QUEUED", "version": 1, "recommendation_status": nil, "market_discovery_status": nil, "improvement_status": nil, "progress": map[string]any{"planned": 5, "completed": 0, "stage": "queued"}}, nil, calls(0, 0, 0, 0), emptyPublication()),
		referenceArtifactExchange(t),
		exchange(t, RunCancel, map[string]any{"run_id": runID('6'), "expected_version": 2}, map[string]any{"run_id": runID('6'), "run_kind": "IMPROVEMENT_SYNTHESIS", "lifecycle": "CANCELLED", "version": 2, "recommendation_status": nil, "market_discovery_status": nil, "improvement_status": nil, "progress": map[string]any{"planned": 5, "completed": 2, "stage": "terminal"}}, nil, calls(0, 0, 0, 0), emptyPublication()),
		referenceToolsExchange(t),
		referenceReplayExchange(t),
		exchange(t, RefinementAudit, map[string]any{"receipt_hashes": map[string]string{"0025": receipt0025Digest, "0026": receipt0026Digest, "0027": receipt0027Digest}, "runtime_instance_ids": map[string]string{"store": "store-1", "queue": "queue-1", "root": "root-1", "cursor_secret": "secret-1", "worker": "worker-1"}, "layer_findings": []string{"L1_PURE", "L2_NO_L3_OR_L4", "L3_NO_L2_OR_L4", "L4_WIRING_ONLY"}}, map[string]any{"receipts_verified": map[string]bool{"0025": true, "0026": true, "0027": true}, "runtime_unique": true, "layer_findings": []string{"L1_PURE", "L2_NO_L3_OR_L4", "L3_NO_L2_OR_L4", "L4_WIRING_ONLY"}}, nil, runtimeCalls(), emptyPublication()),
	}
}

func referenceExtractionExchange(t *testing.T) Exchange {
	t.Helper()
	rules := readPublicO5(t, "00_nucleo", "benchmark", "0028", "mechanism-rules.json")
	implementation := []byte("immutable append implementation")
	independent := []byte("independent recovery test")
	commit := strings.Repeat("a", 40)
	ref := "acme/log@" + commit
	implementationHash := "sha256:" + digestBytes(implementation)
	independentHash := "sha256:" + digestBytes(independent)
	witness := func(kind, path, hash, rule, role string) map[string]any {
		return map[string]any{"repository": "acme/log", "commit": commit, "path": path, "symbol": "", "line": 0, "source_kind": kind, "kind": role, "assessment_status": "PROVEN", "rule_id": rule, "content_hash": hash}
	}
	request := map[string]any{
		"bundle": map[string]any{"run_id": runID('1'), "items": []any{
			map[string]any{"kind": "file", "ref": ref, "path": "log.go", "content_hash": implementationHash, "content_utf8": string(implementation)},
			map[string]any{"kind": "test", "ref": ref, "path": "log_test.go", "content_hash": independentHash, "content_utf8": string(independent)},
		}},
		"assessments": []any{map[string]any{"repository": "acme/log", "status": "PROVEN", "sources": []any{
			map[string]any{"kind": "file", "ref": ref, "path": "log.go", "content_hash": implementationHash, "rule_id": "append-only-write", "symbol": "", "line": 0},
			map[string]any{"kind": "test", "ref": ref, "path": "log_test.go", "content_hash": independentHash, "rule_id": "append-only-recovery-test", "symbol": "", "line": 0},
		}}},
		"conventional_baseline": map[string]any{"baseline_id": "baseline", "idea_hash": strings.Repeat("d", 64), "mechanism_classes": []string{"APPEND_ONLY_LOG_RECOVERY"}},
		"rule_catalog_json":     base64.StdEncoding.EncodeToString(rules), "rule_catalog_sha256": rulesSHA256,
	}
	mechanism := map[string]any{
		"id": "mech_" + strings.Repeat("0", 64), "mechanism_class": "APPEND_ONLY_LOG_RECOVERY", "name": "append-only log with explicit recovery",
		"purpose": "preserve an ordered history that can reconstruct confirmed state after interruption", "trigger": "recovery", "strategy": "append immutable records and rebuild state by deterministic replay with explicit integrity handling",
		"components": []any{}, "data_flow": []string{}, "state_model": []string{}, "failure_model": []string{}, "boundaries": []string{"immutable source"}, "preconditions": []string{}, "tradeoffs": []string{},
		"source":           map[string]any{"repository": "acme/log", "commit": commit},
		"evidence":         []any{witness("file", "log.go", implementationHash, "append-only-write", "IMPLEMENTATION"), witness("test", "log_test.go", independentHash, "append-only-recovery-test", "TEST")},
		"counter_evidence": []any{}, "status": "PROVEN", "novelty": "CONVENTIONAL", "extraction_rule_ids": []string{"append-only-recovery-test", "append-only-write"},
	}
	var typed extractionMechanism
	if err := decodeStrict(mustJSON(t, mechanism), &typed); err != nil {
		t.Fatal(err)
	}
	mechanism["id"] = extractionMechanismID(typed, "MR-APPEND-ONLY-LOG-V1")
	return exchange(t, MechanismExtract, request, map[string]any{"mechanisms": []any{mechanism}, "gaps": []any{}}, nil, calls(0, 0, 0, 0), emptyPublication())
}

func unknownExtractionExchange(t *testing.T) Exchange {
	t.Helper()
	rules := readPublicO5(t, "00_nucleo", "benchmark", "0028", "mechanism-rules.json")
	content := []byte("structured but uncatalogued evidence")
	hash := "sha256:" + digestBytes(content)
	commit := strings.Repeat("e", 40)
	request := map[string]any{
		"bundle":                map[string]any{"run_id": runID('a'), "items": []any{map[string]any{"kind": "file", "ref": "acme/unknown@" + commit, "path": "unknown.go", "content_hash": hash, "content_utf8": string(content)}}},
		"assessments":           []any{map[string]any{"repository": "acme/unknown", "status": "PROVEN", "sources": []any{map[string]any{"kind": "file", "ref": "acme/unknown@" + commit, "path": "unknown.go", "content_hash": hash, "rule_id": "not-in-catalog", "symbol": "", "line": 0}}}},
		"conventional_baseline": nil, "rule_catalog_json": base64.StdEncoding.EncodeToString(rules), "rule_catalog_sha256": rulesSHA256,
	}
	gap := map[string]any{"code": "RULE_NOT_IN_MECHANISM_CATALOG", "source": map[string]any{"repository": "acme/unknown", "commit": commit}, "rule_id": "not-in-catalog", "witness_hash": hash, "status_ceiling": "INCONCLUSIVE"}
	return exchange(t, MechanismExtract, request, map[string]any{"mechanisms": []any{}, "gaps": []any{gap}}, nil, calls(0, 0, 0, 0), emptyPublication())
}

func unavailableSourceExchange(t *testing.T) Exchange {
	return exchange(t, SynthesisStart, validStartRequest(), nil, &ControlledError{Code: "SOURCE_BASIS_UNAVAILABLE", Message: "SOURCE_BASIS_UNAVAILABLE"}, calls(1, 0, 0, 0), emptyPublication())
}

func referenceReplayExchange(t *testing.T) Exchange {
	artifacts := map[string]string{}
	indices := map[string]string{}
	for _, name := range o6ProducedArtifacts {
		artifacts[name] = strings.Repeat("8", 64)
		indices[name] = strings.Repeat("9", 64)
	}
	request := map[string]any{"catalog_sha256": rulesSHA256, "capture_sha256": strings.Repeat("7", 64), "artifact_hashes": map[string]string{"source": strings.Repeat("8", 64)}, "canonical_inputs": []string{base64.StdEncoding.EncodeToString([]byte(`{"schema_version":"1"}`))}}
	result := map[string]any{"artifact_hashes": artifacts, "index_hashes": indices, "network_calls": 0}
	return exchange(t, Replay, request, result, nil, calls(0, 0, 0, 0), emptyPublication())
}

func referenceArtifactExchange(t *testing.T) Exchange {
	contents := []byte(`{"schema_version":"1","artifact":"improvement_hypotheses","entries":[]}`)
	request := map[string]any{"run_id": runID('5'), "artifact": "improvement_hypotheses", "cursor": nil}
	returned := map[string]any{"artifact": "improvement_hypotheses", "sha256": digestBytes(contents), "size": len(contents), "canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil}
	return exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
}

func referenceToolsExchange(t *testing.T) Exchange {
	names := []string{"bateia_search", "bateia_benchmark", "bateia_recommend", "bateia_start_recommendation", "bateia_get_run", "bateia_get_artifact", "bateia_cancel_run", "bateia_list_connectors", "bateia_start_market_discovery", "bateia_start_improvement_synthesis"}
	tools := make([]any, len(names))
	for index, name := range names {
		description := "certified basal tool"
		if index == 9 {
			description = "Inicia síntese durável de hipóteses de melhoria a partir de mecanismos e observações imutáveis"
		}
		schema := any(map[string]any{"type": "object"})
		if index == 5 {
			values := []string{"request", "recommendation", "decision", "alternatives", "failures", "provenance", "connector_capabilities", "market_plan", "market_snapshot", "market_candidates", "market_gaps", "mechanism_catalog", "mechanism_gaps", "architecture_comparison", "improvement_hypotheses", "hypothesis_gaps", "validation_plans", "synthesis_provenance"}
			schema = map[string]any{"properties": map[string]any{"artifact": map[string]any{"enum": values}}}
		}
		if index == 9 {
			if err := json.Unmarshal([]byte(tool10Schema), &schema); err != nil {
				t.Fatal(err)
			}
		}
		tools[index] = map[string]any{"name": name, "description": description, "input_schema": schema}
	}
	return exchange(t, ToolsList, map[string]any{}, map[string]any{"tools": tools}, nil, calls(0, 0, 0, 0), emptyPublication())
}

func validStartRequest() map[string]any {
	return map[string]any{
		"idempotency_key": "stable-key", "idea": map[string]any{"problem": "bounded recovery", "users": []string{}, "capabilities": []string{"recover"}, "constraints": []string{}, "non_goals": []string{}, "preferred_stack": []string{}},
		"source_artifacts": []any{map[string]any{"run_id": runID('1'), "run_kind": "RECOMMENDATION", "artifact": "mechanism_catalog", "sha256": strings.Repeat("1", 64)}},
		"proposals":        []any{}, "policy": map[string]any{"max_hypotheses": 4, "max_mechanisms_per_hypothesis": 2, "require_counterevidence_search": true},
	}
}

func exchange(t *testing.T, port Port, request, returned any, contractError *ControlledError, effects EffectObservation, published PublicationObservation) Exchange {
	t.Helper()
	input := ContractInput{Port: port, Request: mustJSON(t, request)}
	returnRaw := json.RawMessage("null")
	if returned != nil {
		returnRaw = mustJSON(t, returned)
	}
	output := ContractOutput{ReturnValue: returnRaw, Error: contractError, EffectObservation: effects, PublicationObservation: published}
	return Exchange{ContractInput: mustJSON(t, input), ContractOutput: mustJSON(t, output)}
}

func calls(resolve, reserve, id, create int) EffectObservation {
	return EffectObservation{ResolveCalls: resolve, ReserveCalls: reserve, IDCalls: id, CreateCalls: create, StoreInstanceIDs: []string{}, QueueInstanceIDs: []string{}, RootInstanceIDs: []string{}, CursorSecretInstanceIDs: []string{}, WorkerInstanceIDs: []string{}}
}

func runtimeCalls() EffectObservation {
	result := calls(0, 0, 0, 0)
	result.StoreInstanceIDs = []string{"store-1", "store-1"}
	result.QueueInstanceIDs = []string{"queue-1", "queue-1"}
	result.RootInstanceIDs = []string{"root-1", "root-1"}
	result.CursorSecretInstanceIDs = []string{"secret-1", "secret-1"}
	result.WorkerInstanceIDs = []string{"worker-1", "worker-1"}
	return result
}

func emptyPublication() PublicationObservation {
	return PublicationObservation{ArtifactNames: []string{}, ArtifactHashes: map[string]string{}, IndexHashes: map[string]string{}, LeakedFieldNames: []string{}, LayerFindings: []string{}}
}

func publication(names []string) PublicationObservation {
	result := emptyPublication()
	result.ArtifactNames = append([]string(nil), names...)
	for _, name := range names {
		result.ArtifactHashes[name] = strings.Repeat("f", 64)
	}
	result.IndexHashes["artifact-index"] = strings.Repeat("e", 64)
	return result
}

func synthesisPublication() PublicationObservation {
	result := publication(synthesisArtifactOrder)
	result.IndexHashes = map[string]string{}
	for _, name := range synthesisArtifactOrder {
		result.IndexHashes[name] = strings.Repeat("e", 64)
	}
	return result
}

func runID(fill byte) string { return "run_" + strings.Repeat(string(fill), 32) }

func mustJSON(t *testing.T, value any) json.RawMessage {
	t.Helper()
	document, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return document
}

func readPublicO5(t *testing.T, parts ...string) []byte {
	t.Helper()
	_, current, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate O5 package")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(current), "..", "..", ".."))
	document, err := os.ReadFile(filepath.Join(append([]string{root}, parts...)...))
	if err != nil {
		t.Fatal(err)
	}
	return document
}
