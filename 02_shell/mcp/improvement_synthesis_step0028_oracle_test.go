/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/improvement_synthesis_step0028_oracle_test.md
 * @layer L2
 * @updated 2026-09-02
 */

package mcp

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Step0028OracleCheckToolsList accepts a captured JSON-RPC tools/list response.
// The verifier supplies the capture, so the oracle requires no non-contractual
// constructor or private server name.
func Step0028OracleCheckToolsList(t *testing.T, basalNineJSON, candidateTenJSON []byte) {
	t.Helper()
	decode := func(content []byte) []map[string]any {
		var response struct {
			Result struct {
				Tools []map[string]any `json:"tools"`
			} `json:"result"`
		}
		if err := json.Unmarshal(content, &response); err != nil {
			t.Fatal(err)
		}
		return response.Result.Tools
	}
	basal, candidate := decode(basalNineJSON), decode(candidateTenJSON)
	if len(basal) != 9 || len(candidate) != 10 {
		t.Fatalf("tool counts basal=%d candidate=%d", len(basal), len(candidate))
	}
	for index := 0; index < 9; index++ {
		if index == 5 {
			continue
		}
		if !reflect.DeepEqual(basal[index], candidate[index]) {
			t.Fatalf("certified tool %d changed", index+1)
		}
	}
	clone := func(value map[string]any) map[string]any {
		encoded, _ := json.Marshal(value)
		var result map[string]any
		if err := json.Unmarshal(encoded, &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	basalArtifact, candidateArtifact := clone(basal[5]), clone(candidate[5])
	basalEnum := basalArtifact["inputSchema"].(map[string]any)["properties"].(map[string]any)["artifact"].(map[string]any)["enum"].([]any)
	candidateSchema := candidateArtifact["inputSchema"].(map[string]any)["properties"].(map[string]any)["artifact"].(map[string]any)
	wantSuffix := []any{"mechanism_catalog", "mechanism_gaps", "architecture_comparison", "improvement_hypotheses", "hypothesis_gaps", "validation_plans", "synthesis_provenance"}
	candidateEnum := candidateSchema["enum"].([]any)
	if len(candidateEnum) != len(basalEnum)+len(wantSuffix) || !reflect.DeepEqual(candidateEnum[len(basalEnum):], wantSuffix) {
		t.Fatalf("artifact enum suffix = %v", candidateEnum)
	}
	candidateSchema["enum"] = append([]any(nil), basalEnum...)
	if !reflect.DeepEqual(basalArtifact, candidateArtifact) {
		t.Fatal("tool 6 changed outside the seven additive artifact enum entries")
	}
	wantNames := []string{"bateia_search", "bateia_benchmark", "bateia_recommend", "bateia_start_recommendation", "bateia_get_run", "bateia_get_artifact", "bateia_cancel_run", "bateia_list_connectors", "bateia_start_market_discovery", "bateia_start_improvement_synthesis"}
	for index, want := range wantNames {
		if candidate[index]["name"] != want {
			t.Fatalf("tool[%d]=%v, want %s", index, candidate[index]["name"], want)
		}
	}
	tool := candidate[9]
	if tool["description"] != "Inicia síntese durável de hipóteses de melhoria a partir de mecanismos e observações imutáveis" {
		t.Fatalf("tool 10 description changed: %v", tool["description"])
	}
	schema, ok := tool["inputSchema"].(map[string]any)
	if !ok || schema["type"] != "object" || schema["additionalProperties"] != false {
		t.Fatalf("tool 10 schema root: %#v", tool["inputSchema"])
	}
	if !reflect.DeepEqual(schema["required"], []any{"idempotency_key", "idea", "source_artifacts", "proposals", "policy"}) {
		t.Fatalf("tool 10 required: %v", schema["required"])
	}
}
