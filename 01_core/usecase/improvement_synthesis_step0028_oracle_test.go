/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/improvement_synthesis_step0028_oracle_test.md
 * @layer L1
 * @updated 2026-09-02
 */

package usecase_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type step0028SynthesisFixture struct {
	SchemaVersion string                     `json:"schema_version"`
	Cases         map[string]json.RawMessage `json:"cases"`
}

type step0028RuntimeFixture struct {
	SchemaVersion string                     `json:"schema_version"`
	Cases         map[string]json.RawMessage `json:"cases"`
}

func TestStep0028OracleProtectedSynthesisAndRuntimeFixturesAreComplete(t *testing.T) {
	read := func(name string, target any) {
		t.Helper()
		content, err := os.ReadFile(filepath.Join("..", "..", "00_nucleo", "benchmark", "0028", "oracle", "protected", "fixtures", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(content, target); err != nil {
			t.Fatal(err)
		}
	}
	var synthesis step0028SynthesisFixture
	var runtime step0028RuntimeFixture
	read("synthesis.json", &synthesis)
	read("runtime.json", &runtime)
	if synthesis.SchemaVersion != "tekt.step0028.synthesis-fixtures.v1" || runtime.SchemaVersion != "tekt.step0028.runtime-fixtures.v1" {
		t.Fatalf("fixture schemas changed: %q %q", synthesis.SchemaVersion, runtime.SchemaVersion)
	}
	wantSynthesis := []string{"composition_unknown", "popularity_novelty_invariant", "product_only", "unknown_target", "unrecognized_license"}
	wantRuntime := []string{"offline_replay", "source_hash_prequeue", "tool_registry"}
	keys := func(values map[string]json.RawMessage) []string {
		out := make([]string, 0, len(values))
		for key := range values {
			out = append(out, key)
		}
		sort.Strings(out)
		return out
	}
	if !reflect.DeepEqual(keys(synthesis.Cases), wantSynthesis) || !reflect.DeepEqual(keys(runtime.Cases), wantRuntime) {
		t.Fatalf("protected coverage changed: synthesis=%v runtime=%v", keys(synthesis.Cases), keys(runtime.Cases))
	}
	composition := string(synthesis.Cases["composition_unknown"])
	for _, required := range []string{"interface compatibility", "state model compatibility", "failure model compatibility", "performance compatibility", "license compatibility"} {
		if !strings.Contains(composition, required) {
			t.Fatalf("composition fixture lost %q", required)
		}
	}
}

// Step0028OracleCheckHypothesisArtifact checks a candidate artifact captured by
// the verifier without depending on a concrete synthesis constructor.
func Step0028OracleCheckHypothesisArtifact(t *testing.T, document []byte) {
	t.Helper()
	var envelope struct {
		SchemaVersion string           `json:"schema_version"`
		Artifact      string           `json:"artifact"`
		Entries       []map[string]any `json:"entries"`
	}
	if err := json.Unmarshal(document, &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.SchemaVersion != "1" || envelope.Artifact != "improvement_hypotheses" || envelope.Entries == nil {
		t.Fatalf("invalid hypothesis artifact envelope: %s", document)
	}
	for _, hypothesis := range envelope.Entries {
		if hypothesis["adoption_disposition"] != "NOT_AN_ADOPTION_DECISION" {
			t.Fatalf("adoption disclaimer missing: %#v", hypothesis)
		}
		for _, key := range []string{"assumptions", "preconditions", "integration_unknowns", "risks", "tradeoffs", "evidence_refs", "counter_evidence_refs", "validation_plan"} {
			if value, ok := hypothesis[key]; !ok || value == nil {
				t.Fatalf("hypothesis lacks explicit %s: %#v", key, hypothesis)
			}
		}
		if hypothesis["target_state"] == "UNKNOWN" && hypothesis["status"] == "SUPPORTED" {
			t.Fatalf("UNKNOWN target became SUPPORTED: %#v", hypothesis)
		}
		if hypothesis["license"] == "UNRECOGNIZED" && hypothesis["transfer_mode"] == "COMPONENT_REUSE" {
			t.Fatalf("unknown license enabled reuse: %#v", hypothesis)
		}
		if refs, ok := hypothesis["mechanism_refs"].([]any); ok && len(refs) > 1 {
			encoded, _ := json.Marshal(hypothesis["integration_unknowns"])
			for _, required := range []string{"interface compatibility", "state model compatibility", "failure model compatibility", "performance compatibility", "license compatibility"} {
				if !strings.Contains(string(encoded), required) {
					t.Fatalf("composition lost %q: %#v", required, hypothesis)
				}
			}
		}
	}
	for _, forbidden := range []string{"ORACLE_TEST_CANARY", `"snippet"`, `"content"`, `"rationale"`, `"idempotency_key"`} {
		if strings.Contains(string(document), forbidden) {
			t.Fatalf("unsafe field/content %q crossed artifact", forbidden)
		}
	}
}
