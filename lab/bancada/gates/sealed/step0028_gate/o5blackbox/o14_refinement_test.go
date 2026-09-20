/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o14_refinement_test.md
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

func TestO14MechanismPreimageUsesCatalogRuleRatherThanEvidenceRule(t *testing.T) {
	control := documentaryExtractionExchange(t)
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("catalog-rule documentary control rejected: %+v", finding)
	}
	var output ContractOutput
	var result extractionResult
	if decodeStrict(control.ContractOutput, &output) != nil || decodeStrict(output.ReturnValue, &result) != nil {
		t.Fatal("could not decode documentary control")
	}
	if len(result.Mechanisms) != 1 || !extractionEqualStrings(result.Mechanisms[0].ExtractionRuleIDs, []string{"append-only-architecture-documentation"}) {
		t.Fatal("documentary source RuleID was not preserved independently")
	}
	result.Mechanisms[0].ID = extractionMechanismID(result.Mechanisms[0], "append-only-architecture-documentation")
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("evidence RuleID incorrectly survived as catalog mechanism RuleID: %+v", finding)
	}
}

func TestO14MinimumWitnessSemanticLimitIsUnknown(t *testing.T) {
	control := documentaryExtractionExchange(t)
	var output ContractOutput
	var result extractionResult
	if decodeStrict(control.ContractOutput, &output) != nil || decodeStrict(output.ReturnValue, &result) != nil {
		t.Fatal("could not decode documentary control")
	}
	mechanism := result.Mechanisms[0]
	result.Gaps = []extractionGap{{
		Code: "MINIMUM_WITNESS_SEMANTICALLY_INSUFFICIENT", Source: mechanism.Source,
		RuleID: mechanism.Evidence[0].RuleID, WitnessHash: mechanism.Evidence[0].ContentHash, StatusCeiling: "PARTIAL",
	}}
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Unknown {
		t.Fatalf("minimum witness epistemic limit was not Unknown: %+v", finding)
	}
}

func TestO14PermittedContractStatesAreRelationallyUnknown(t *testing.T) {
	plan := []any{map[string]any{"missing_evidence": "measure the unresolved relation"}}
	tests := map[string]map[string]any{
		"target baseline": {
			"target_state": "UNKNOWN", "baseline_absent": true, "target_property": "latency", "validation_plan": plan,
		},
		"composition": {
			"mechanism_refs": []any{"mech_a", "mech_b"}, "integration_unknowns": []any{"interface compatibility"}, "status": "INCONCLUSIVE", "validation_plan": plan,
		},
		"license": {
			"license": "UNRECOGNIZED", "transfer_mode": "UNKNOWN", "status": "INCONCLUSIVE", "validation_plan": plan,
		},
		"counterevidence budget": {
			"counter_evidence_search": map[string]any{"status": "BUDGET_EXHAUSTED"}, "status": "INCONCLUSIVE", "validation_plan": plan,
		},
		"identity": {
			"code": "IMMUTABLE_IDENTITY_MISSING", "status_ceiling": "INCONCLUSIVE", "source_ref": "unresolved-source",
		},
		"provider capture": {
			"provider_availability_evidenced": false, "frozen_capture_present": false, "source_ref": "provider-source",
		},
		"utility baseline": {
			"effect_kind": "UTILITY_COST", "baseline_absent": true, "target_property": "utility cost", "validation_plan": plan,
		},
		"security baseline": {
			"effect_kind": "SECURITY_EFFECT", "baseline_absent": true, "target_property": "security effect", "validation_plan": plan,
		},
	}
	for name, entry := range tests {
		t.Run(name, func(t *testing.T) {
			control := o14ArtifactExchange(t, entry)
			if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Unknown {
				t.Fatalf("relational opacity was not Unknown: %+v", finding)
			}
		})
	}
}

func TestO14ExplicitSafeLicenseTransfersRemainPreserved(t *testing.T) {
	for _, transfer := range []string{"PATTERN_ADAPTATION", "REFERENCE_ONLY"} {
		entry := map[string]any{
			"license": "UNRECOGNIZED", "transfer_mode": transfer, "status": "PROPOSED",
			"target_property": "bounded recovery", "validation_plan": []any{map[string]any{"missing_evidence": "human license review"}},
		}
		control := o14ArtifactExchange(t, entry)
		if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
			t.Fatalf("safe explicit transfer %s became opaque: %+v", transfer, finding)
		}
	}
}

func TestO14ViolationPrecedesRelationalOpacity(t *testing.T) {
	control := o14ArtifactExchange(t, map[string]any{
		"target_state": "UNKNOWN", "baseline_absent": true, "target_property": "latency", "validation_plan": []any{"measure baseline"},
	})
	var output ContractOutput
	if decodeStrict(control.ContractOutput, &output) != nil {
		t.Fatal("could not decode artifact control")
	}
	var returned o6ArtifactResult
	if decodeStrict(output.ReturnValue, &returned) != nil {
		t.Fatal("could not decode artifact result")
	}
	returned.SHA256 = strings.Repeat("0", 64)
	output.ReturnValue = mustJSON(t, returned)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("relational opacity masked an integrity violation: %+v", finding)
	}
}

func TestO14ProtectedManifestRemainsAnImmutableHistoricalAnchor(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o14.json")
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
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v14" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O14" || len(manifest.Files) != 28 {
		t.Fatal("unexpected O14 protected manifest identity")
	}
	if digestBytes(document) != "ac8d10da663d301f4f10beb3c35113db7b27f61a50b2e037086642b7c63b1343" {
		t.Fatal("O14 protected manifest historical bytes changed")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || strings.HasPrefix(file.Path, "/") || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe O14 protected manifest entry %q", file.Path)
		}
		seen[file.Path] = true
	}
}

func TestO14PublicJSONArtifactsAreStrictAndDuplicateFree(t *testing.T) {
	paths := [][]string{
		{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o14.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o13-to-o14.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o14.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o14.json"},
	}
	for _, path := range paths {
		if _, err := o6ParseUniqueJSON(readPublicO5(t, path...)); err != nil {
			t.Fatalf("O14 artifact %s is not strict duplicate-free JSON: %v", strings.Join(path, "/"), err)
		}
	}
}

func o14ArtifactExchange(t *testing.T, entry map[string]any) Exchange {
	t.Helper()
	entries, err := json.Marshal([]any{entry})
	if err != nil {
		t.Fatal(err)
	}
	contents := []byte(`{"schema_version":"1","artifact":"architecture_comparison","entries":` + string(entries) + `}`)
	request := map[string]any{"run_id": runID('e'), "artifact": "architecture_comparison", "cursor": nil}
	returned := map[string]any{
		"artifact": "architecture_comparison", "sha256": digestBytes(contents), "size": len(contents),
		"canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil,
	}
	return exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
}
