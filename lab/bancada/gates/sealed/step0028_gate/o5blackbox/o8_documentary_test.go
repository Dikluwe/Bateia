/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o8_documentary_test.md
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

func TestO8HistoricalProtectedManifestRemainsFrozenAfterO9Invalidation(t *testing.T) {
	document := readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o8.json")
	if digestBytes(document) != "30a5af21d089938c35f5fed4b001cd6e714cdfae850057409ec0d47d5e8b916c" {
		t.Fatal("historical O8 protected manifest changed")
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
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v8" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O8" || len(manifest.Files) != 20 {
		t.Fatal("unexpected O8 protected manifest identity")
	}
}

func TestO8DocumentaryEvidenceMaterializesOnlyPartial(t *testing.T) {
	control := documentaryExtractionExchange(t)
	if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("literal documentary evidence was not preserved: %+v", finding)
	}

	var output ContractOutput
	if err := decodeStrict(control.ContractOutput, &output); err != nil {
		t.Fatal(err)
	}
	var result extractionResult
	if err := decodeStrict(output.ReturnValue, &result); err != nil {
		t.Fatal(err)
	}
	result.Mechanisms[0].Status = "PROVEN"
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("documentary-only evidence was promoted to PROVEN: %+v", finding)
	}
}

func TestO8DocumentaryRoleRequiresLiteralCatalogShapeAndDocumentationKind(t *testing.T) {
	rules := readPublicO5(t, "00_nucleo", "benchmark", "0028", "mechanism-rules.json")
	var catalog extractionCatalog
	if err := decodeStrict(rules, &catalog); err != nil {
		t.Fatal(err)
	}
	if !extractionCatalogIsExact(catalog) {
		t.Fatal("reference v3 catalog was not accepted")
	}
	catalog.Rules[0].AcceptedEvidenceRules[5].StatusMax = "PROVEN"
	if extractionCatalogIsExact(catalog) {
		t.Fatal("documentary status_max other than PARTIAL was accepted")
	}

	control := documentaryExtractionExchange(t)
	var input ContractInput
	if err := decodeStrict(control.ContractInput, &input); err != nil {
		t.Fatal(err)
	}
	var request extractionRequest
	if err := decodeStrict(input.Request, &request); err != nil {
		t.Fatal(err)
	}
	request.Bundle.Items[0].Kind = "file"
	request.Assessments[0].Sources[0].Kind = "file"
	input.Request = mustJSON(t, request)
	if finding := Evaluate(mustJSON(t, input), control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("DOCUMENTARY accepted a non-documentation source kind: %+v", finding)
	}
}

func TestO8CatalogHashBindingRejectsThePredecessorCatalog(t *testing.T) {
	control := documentaryExtractionExchange(t)
	var input ContractInput
	if err := decodeStrict(control.ContractInput, &input); err != nil {
		t.Fatal(err)
	}
	var request extractionRequest
	if err := decodeStrict(input.Request, &request); err != nil {
		t.Fatal(err)
	}
	request.RuleCatalogSHA256 = strings.Repeat("0", 64)
	input.Request = mustJSON(t, request)
	if finding := Evaluate(mustJSON(t, input), control.ContractOutput); finding.Classification != Violated {
		t.Fatalf("non-v3 catalog binding survived: %+v", finding)
	}
}

func documentaryExtractionExchange(t *testing.T) Exchange {
	t.Helper()
	rules := readPublicO5(t, "00_nucleo", "benchmark", "0028", "mechanism-rules.json")
	content := []byte("architecture boundary documentation")
	contentHash := "sha256:" + digestBytes(content)
	commit := strings.Repeat("b", 40)
	ref := "acme/documented-log@" + commit
	ruleID := "append-only-architecture-documentation"
	request := map[string]any{
		"bundle": map[string]any{"run_id": runID('8'), "items": []any{
			map[string]any{"kind": "documentation", "ref": ref, "path": "architecture.md", "content_hash": contentHash, "content_utf8": string(content)},
		}},
		"assessments": []any{map[string]any{"repository": "acme/documented-log", "status": "PROVEN", "sources": []any{
			map[string]any{"kind": "documentation", "ref": ref, "path": "architecture.md", "content_hash": contentHash, "rule_id": ruleID, "symbol": "", "line": 0},
		}}},
		"conventional_baseline": nil,
		"rule_catalog_json":     base64.StdEncoding.EncodeToString(rules),
		"rule_catalog_sha256":   rulesSHA256,
	}
	witness := map[string]any{
		"repository": "acme/documented-log", "commit": commit, "path": "architecture.md", "symbol": "", "line": 0,
		"source_kind": "documentation", "kind": "DOCUMENTATION", "assessment_status": "PROVEN", "rule_id": ruleID, "content_hash": contentHash,
	}
	mechanism := map[string]any{
		"id": "mech_" + strings.Repeat("0", 64), "mechanism_class": "APPEND_ONLY_LOG_RECOVERY", "name": "append-only log with explicit recovery",
		"purpose": "preserve an ordered history that can reconstruct confirmed state after interruption", "trigger": "documented recovery boundary",
		"strategy":   "append immutable records and rebuild state by deterministic replay with explicit integrity handling",
		"components": []any{}, "data_flow": []string{}, "state_model": []string{}, "failure_model": []string{},
		"boundaries": []string{"repository-backed source bundle"}, "preconditions": []string{}, "tradeoffs": []string{},
		"source": map[string]any{"repository": "acme/documented-log", "commit": commit}, "evidence": []any{witness}, "counter_evidence": []any{},
		"status": "PARTIAL", "novelty": "UNKNOWN", "extraction_rule_ids": []string{ruleID},
	}
	var typed extractionMechanism
	if err := decodeStrict(mustJSON(t, mechanism), &typed); err != nil {
		t.Fatal(err)
	}
	mechanism["id"] = extractionMechanismID(typed, "MR-APPEND-ONLY-LOG-V1")
	return exchange(t, MechanismExtract, request, map[string]any{"mechanisms": []any{mechanism}, "gaps": []any{}}, nil, calls(0, 0, 0, 0), emptyPublication())
}
