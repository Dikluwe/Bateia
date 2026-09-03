/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o15_refinement_test.md
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

func TestO15MinimumWitnessRelationDistinguishesTechnicalAndDocumentaryPartial(t *testing.T) {
	rule := extractionCatalogMechanism{MinimumWitnesses: extractionMinimumWitnesses{DistinctContentHashes: 2}, AcceptedEvidenceRules: []extractionCatalogEvidence{
		{RuleID: "implementation-rule", Role: "IMPLEMENTATION", WitnessKind: "SOURCE"},
		{RuleID: "independent-rule", Role: "INDEPENDENT", WitnessKind: "TEST"},
		{RuleID: "documentary-rule", Role: "DOCUMENTARY", WitnessKind: "DOCUMENTATION"},
	}}
	technical := &extractionExpectedGroup{Rule: rule, Status: "PARTIAL", Positive: []extractionWitness{{RuleID: "implementation-rule", Kind: "SOURCE", AssessmentStatus: "PARTIAL", ContentHash: strings.Repeat("a", 64)}}}
	if !o15GroupHasSemanticWitnessLimit(technical) {
		t.Fatal("technical minimum-witness limit was not recognized")
	}
	documentary := &extractionExpectedGroup{Rule: rule, Status: "PARTIAL", Positive: []extractionWitness{{RuleID: "documentary-rule", Kind: "DOCUMENTATION", AssessmentStatus: "PARTIAL", ContentHash: strings.Repeat("b", 64)}}}
	if o15GroupHasSemanticWitnessLimit(documentary) {
		t.Fatal("accepted documentary PARTIAL became opaque")
	}
}

func TestO15PublicOpaqueArtifactRelationsAreUnknown(t *testing.T) {
	plan := []any{map[string]any{"missing_evidence": []any{"frozen observation"}}}
	tests := []struct {
		name, artifact string
		entry          map[string]any
	}{
		{"target baseline", "architecture_comparison", map[string]any{"target_property": "bounded-target", "mechanism_id": "mech-a", "source_run_id": "run-a", "target_state": "UNKNOWN", "baseline_absent": true, "validation_plan": plan}},
		{"composition", "improvement_hypotheses", map[string]any{"id": "h-composition", "mechanism_refs": []any{"mech-a", "mech-b"}, "integration_unknowns": []any{"unmeasured relation"}, "status": "INCONCLUSIVE", "validation_plan": plan}},
		{"license", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-license", "code": "LICENSE_NOT_RECOGNIZED", "source_ref": "source-license"}},
		{"counterevidence budget", "improvement_hypotheses", map[string]any{"id": "h-budget", "counterevidence_search": map[string]any{"status": "BUDGET_EXHAUSTED"}, "status": "INCONCLUSIVE", "validation_plan": plan}},
		{"identity", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-identity", "code": "IMMUTABLE_IDENTITY_AMBIGUOUS", "source_ref": "source-a"}},
		{"provider capture", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-provider", "code": "PROVIDER_AVAILABILITY_NOT_EVIDENCED_BY_FROZEN_CAPTURE", "source_ref": "source-b"}},
		{"utility baseline", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-utility", "code": "UTILITY_COST_BASELINE_MISSING", "source_ref": "source-c"}},
		{"security baseline", "hypothesis_gaps", map[string]any{"hypothesis_id": "h-security", "code": "SECURITY_EFFECT_BASELINE_MISSING", "source_ref": "source-d"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			control := o15ArtifactExchange(t, tc.artifact, tc.entry)
			if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Unknown {
				t.Fatalf("public protocol relation was not Unknown: %+v", finding)
			}
		})
	}
}

func TestO15SafeLicenseAndIsolatedUnknownTextRemainPreserved(t *testing.T) {
	for _, entry := range []map[string]any{
		{"id": "h-adapt", "license": "UNRECOGNIZED", "transfer_mode": "PATTERN_ADAPTATION", "status": "PROPOSED", "validation_plan": []any{"validate"}},
		{"id": "h-reference", "license": "UNRECOGNIZED", "transfer_mode": "REFERENCE_ONLY", "status": "PROPOSED", "validation_plan": []any{"validate"}},
		{"id": "h-text", "title": "provider availability unknown and baseline missing", "status": "PROPOSED"},
	} {
		control := o15ArtifactExchange(t, "improvement_hypotheses", entry)
		if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
			t.Fatalf("non-controlling state became opaque: %+v", finding)
		}
	}
}

func TestO15ViolationPrecedesPublicOpacity(t *testing.T) {
	control := o15ArtifactExchange(t, "hypothesis_gaps", map[string]any{"hypothesis_id": "h", "code": "PROVIDER_AVAILABILITY_NOT_EVIDENCED", "source_ref": "s"})
	var output ContractOutput
	var returned o6ArtifactResult
	if decodeStrict(control.ContractOutput, &output) != nil || decodeStrict(output.ReturnValue, &returned) != nil {
		t.Fatal("could not decode opacity control")
	}
	returned.SHA256 = strings.Repeat("0", 64)
	output.ReturnValue = mustJSON(t, returned)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("opacity masked integrity violation: %+v", finding)
	}
}

func TestO15PublicJSONArtifactsAreStrictAndDuplicateFree(t *testing.T) {
	paths := [][]string{
		{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o15.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o14-to-o15.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o15.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o15.json"},
	}
	for _, path := range paths {
		if _, err := o6ParseUniqueJSON(readPublicO5(t, path...)); err != nil {
			t.Fatalf("O15 artifact %s is not strict duplicate-free JSON: %v", strings.Join(path, "/"), err)
		}
	}
}

func TestO15ProtectedManifestRemainsAnImmutableHistoricalAnchor(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o15.json")
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
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v15" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O15" || len(manifest.Files) != 30 {
		t.Fatal("unexpected O15 protected manifest identity")
	}
	if digestBytes(document) != "ba502012af763a47420ca1d7b80b493d13620aea5cace203eac6ec4602b8f25b" {
		t.Fatal("O15 protected manifest historical bytes changed")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || filepath.IsAbs(file.Path) || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe O15 protected manifest entry %q", file.Path)
		}
		seen[file.Path] = true
	}
}

func o15ArtifactExchange(t *testing.T, artifact string, entry map[string]any) Exchange {
	t.Helper()
	entries, err := json.Marshal([]any{entry})
	if err != nil {
		t.Fatal(err)
	}
	contents := []byte(`{"schema_version":"1","artifact":"` + artifact + `","entries":` + string(entries) + `}`)
	request := map[string]any{"run_id": runID('f'), "artifact": artifact, "cursor": nil}
	returned := map[string]any{"artifact": artifact, "sha256": digestBytes(contents), "size": len(contents), "canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil}
	return exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
}
