/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o13_refinement_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestO13CanonicalMechanismIDUsesLiteralR4Preimage(t *testing.T) {
	hashA := "sha256:" + strings.Repeat("a", 64)
	hashB := "sha256:" + strings.Repeat("b", 64)
	mechanism := extractionMechanism{
		Source:         extractionImmutable{Repository: "acme/<&>", Commit: strings.Repeat("c", 40)},
		MechanismClass: "CLASS<&>",
		Purpose:        "Cafe\u0301 <&>",
		Strategy:       "\u1100\u1161",
		Evidence: []extractionWitness{
			{ContentHash: hashB},
			{ContentHash: hashA},
		},
		CounterEvidence: []extractionWitness{{ContentHash: hashA}},
	}
	ruleID := "RULE<&>"
	expectedPreimage := []byte(`{"schema_version":"1","source":{"repository":"acme/<&>","commit":"cccccccccccccccccccccccccccccccccccccccc"},"rule_id":"RULE<&>","mechanism_class":"CLASS<&>","purpose":"Café <&>","strategy":"가","witness_hashes":["` + hashA + `","` + hashB + `"]}`)
	expectedID := "mech_" + digestBytes(expectedPreimage)
	if actual := extractionMechanismID(mechanism, ruleID); actual != expectedID {
		t.Fatalf("causal ID did not use the literal R4 bytes: got %s want %s", actual, expectedID)
	}

	composed := mechanism
	composed.Purpose = "Café <&>"
	composed.Strategy = "가"
	if extractionMechanismID(composed, ruleID) != expectedID {
		t.Fatal("canonically equivalent NFC text changed mechanism identity")
	}

	changedCounter := mechanism
	changedCounter.CounterEvidence = []extractionWitness{{ContentHash: "sha256:" + strings.Repeat("d", 64)}}
	if extractionMechanismID(changedCounter, ruleID) == expectedID {
		t.Fatal("counter-evidence hash was omitted from the causal identity")
	}
}

func TestO13MissingBaselineIsRealExtractionOpacityForProvenNovelty(t *testing.T) {
	control := referenceExtractionExchange(t)
	var input ContractInput
	var request extractionRequest
	var output ContractOutput
	var result extractionResult
	if decodeStrict(control.ContractInput, &input) != nil || decodeStrict(input.Request, &request) != nil ||
		decodeStrict(control.ContractOutput, &output) != nil || decodeStrict(output.ReturnValue, &result) != nil {
		t.Fatal("could not decode the public extraction control")
	}
	request.ConventionalBaseline = nil
	result.Mechanisms[0].Novelty = "UNKNOWN"
	input.Request = mustJSON(t, request)
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(mustJSON(t, input), mustJSON(t, output)); finding.Classification != Unknown {
		t.Fatalf("missing baseline did not remain epistemically opaque: %+v", finding)
	}

	result.Mechanisms[0].Novelty = "DIVERGENT"
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(mustJSON(t, input), mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("opacity masked an invented novelty classification: %+v", finding)
	}
}

func TestO13ConformantLicenseAndIntegrationUnknownsAreNotGateOpacity(t *testing.T) {
	entries := `[{"id":"hyp_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","license":"UNRECOGNIZED","integration_unknowns":["failure model compatibility","interface compatibility","license compatibility","performance compatibility","state model compatibility"]}]`
	assertO13ArtifactClassification(t, "improvement_hypotheses", entries, Preserved)
}

func TestO13HistoricalProtectedManifestRemainsFrozenAfterO14Invalidation(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o13.json")
	if digestBytes(document) != "06d8c7a88910b4ef2e99e9db83b5a63384546d7fbdf6d83b5c21506bfc76afb1" {
		t.Fatal("historical O13 protected manifest changed")
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
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v13" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O13" || len(manifest.Files) != 26 {
		t.Fatal("unexpected O13 protected manifest identity")
	}
}

func TestO13PublicJSONArtifactsAreStrictAndDuplicateFree(t *testing.T) {
	paths := [][]string{
		{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o13.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o11-to-o13.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o13.json"},
		{"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o13.json"},
	}
	for _, path := range paths {
		document := readPublicO5(t, path...)
		if _, err := o6ParseUniqueJSON(document); err != nil {
			t.Fatalf("O13 artifact %s is not strict duplicate-free JSON: %v", strings.Join(path, "/"), err)
		}
	}
}

func assertO13ArtifactClassification(t *testing.T, artifact, entries string, expected Classification) {
	t.Helper()
	compactEntries := new(bytes.Buffer)
	if err := json.Compact(compactEntries, []byte(entries)); err != nil {
		t.Fatal(err)
	}
	contents := []byte(`{"schema_version":"1","artifact":"` + artifact + `","entries":` + compactEntries.String() + `}`)
	request := map[string]any{"run_id": runID('d'), "artifact": artifact, "cursor": nil}
	returned := map[string]any{
		"artifact": artifact, "sha256": digestBytes(contents), "size": len(contents),
		"canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil,
	}
	control := exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), emptyPublication())
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != expected {
		t.Fatalf("artifact %s classification: %+v", artifact, finding)
	}
}
