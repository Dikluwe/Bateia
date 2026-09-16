/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/algorithm_history_test.md
 * @layer L1
 * @updated 2026-09-12
 */

package domain_test

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"bateia/01_core/domain"
)

func roundTripHistoryValue(t *testing.T, raw string, destination any) map[string]any {
	t.Helper()
	if err := json.Unmarshal([]byte(raw), destination); err != nil {
		t.Fatalf("construct public domain value: %v", err)
	}
	encoded, err := json.Marshal(destination)
	if err != nil {
		t.Fatalf("marshal public domain value: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("decode marshaled public domain value: %v", err)
	}
	return got
}

func decodeHistoryObject(t *testing.T, raw string) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func requireExactHistoryShape(t *testing.T, got, want map[string]any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		gotJSON, _ := json.MarshalIndent(got, "", "  ")
		wantJSON, _ := json.MarshalIndent(want, "", "  ")
		t.Fatalf("public JSON tags/shape differ:\n got %s\nwant %s", gotJSON, wantJSON)
	}
}

func TestAlgorithmHistoryInputPublicConstantsAndJSONShape(t *testing.T) {
	const raw = `{
  "schema_version":"algorithm-history-input-v1",
  "algorithms":[
    {"id":"known","name":"Known","chronology":"KNOWN","year_start":1,"year_end":1,"chronology_reason":"dated","provenance_ids":["primary"]},
    {"id":"approximate","name":"Approximate","chronology":"APPROXIMATE","year_start":2,"year_end":3,"chronology_reason":"bounded","provenance_ids":["secondary"]},
    {"id":"unknown","name":"Unknown","chronology":"UNKNOWN","year_start":null,"year_end":null,"chronology_reason":"undated","provenance_ids":["artifact"]}
  ],
  "provenance":[
    {"id":"primary","kind":"PRIMARY_SOURCE","origin_id":"origin:primary","citation":"Primary citation","locator":"archive://primary"},
    {"id":"secondary","kind":"SECONDARY_SOURCE","origin_id":"origin:secondary","citation":"Secondary citation","locator":"book://secondary"},
    {"id":"artifact","kind":"EMPIRICAL_ARTIFACT","origin_id":"origin:artifact","citation":"Artifact citation","locator":"artifact://sample"}
  ],
  "claims":[
    {"id":"c-derived","subject_id":"approximate","relation":"DERIVED_FROM","object_id":"known","statement":"derived"},
    {"id":"c-influenced","subject_id":"approximate","relation":"INFLUENCED_BY","object_id":"known","statement":"influenced"},
    {"id":"c-improves","subject_id":"unknown","relation":"IMPROVES","object_id":"known","statement":"improves"},
    {"id":"c-generalizes","subject_id":"unknown","relation":"GENERALIZES","object_id":"known","statement":"generalizes"},
    {"id":"c-specializes","subject_id":"unknown","relation":"SPECIALIZES","object_id":"known","statement":"specializes"},
    {"id":"c-combines","subject_id":"unknown","relation":"COMBINES","object_id":"known","statement":"combines"},
    {"id":"c-rediscovers","subject_id":"unknown","relation":"REDISCOVERS","object_id":"known","statement":"rediscovers"},
    {"id":"c-resembles","subject_id":"unknown","relation":"RESEMBLES","object_id":"known","statement":"resembles"}
  ],
  "bases":[
    {"id":"e-support","claim_id":"c-derived","kind":"EVIDENCE","stance":"SUPPORT","provenance_id":"primary","summary":"support evidence","derived_from_ids":[]},
    {"id":"e-conflict","claim_id":"c-influenced","kind":"EVIDENCE","stance":"CONFLICT","provenance_id":"secondary","summary":"conflicting evidence","derived_from_ids":[]},
    {"id":"e-ambiguous","claim_id":"c-improves","kind":"EVIDENCE","stance":"AMBIGUOUS","provenance_id":"artifact","summary":"ambiguous evidence","derived_from_ids":[]},
    {"id":"i-support","claim_id":"c-derived","kind":"INFERENCE","stance":"SUPPORT","provenance_id":"primary","summary":"support inference","derived_from_ids":["e-support"]},
    {"id":"i-conflict","claim_id":"c-influenced","kind":"INFERENCE","stance":"CONFLICT","provenance_id":"secondary","summary":"conflict inference","derived_from_ids":["e-conflict"]},
    {"id":"i-ambiguous","claim_id":"c-improves","kind":"INFERENCE","stance":"AMBIGUOUS","provenance_id":"artifact","summary":"ambiguous inference","derived_from_ids":["e-ambiguous"]}
  ]
}`
	var input domain.AlgorithmHistoryInput
	got := roundTripHistoryValue(t, raw, &input)
	requireExactHistoryShape(t, got, decodeHistoryObject(t, raw))

	if got["schema_version"] != "algorithm-history-input-v1" {
		t.Fatalf("input schema constant changed: %v", got["schema_version"])
	}
	assertHistoryTokens(t, got["algorithms"], "chronology", []string{"APPROXIMATE", "KNOWN", "UNKNOWN"})
	assertHistoryTokens(t, got["provenance"], "kind", []string{"EMPIRICAL_ARTIFACT", "PRIMARY_SOURCE", "SECONDARY_SOURCE"})
	assertHistoryTokens(t, got["claims"], "relation", []string{"COMBINES", "DERIVED_FROM", "GENERALIZES", "IMPROVES", "INFLUENCED_BY", "REDISCOVERS", "RESEMBLES", "SPECIALIZES"})
	assertHistoryTokens(t, got["bases"], "kind", []string{"EVIDENCE", "INFERENCE"})
	assertHistoryTokens(t, got["bases"], "stance", []string{"AMBIGUOUS", "CONFLICT", "SUPPORT"})
}

func assertHistoryTokens(t *testing.T, raw any, field string, want []string) {
	t.Helper()
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("expected array while checking %s, got %T", field, raw)
	}
	seen := map[string]bool{}
	for _, item := range items {
		object, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("expected object while checking %s, got %T", field, item)
		}
		value, ok := object[field].(string)
		if !ok {
			t.Fatalf("missing string field %s in %v", field, object)
		}
		seen[value] = true
	}
	got := make([]string, 0, len(seen))
	for value := range seen {
		got = append(got, value)
	}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("public tokens for %s: got %v want %v", field, got, want)
	}
}

func TestAlgorithmHistoryReportPublicConstantsAndJSONShape(t *testing.T) {
	const complete = `{
  "schema_version":"algorithm-history-report-v1",
  "status":"COMPLETE",
  "reason":"analyzed",
  "assessments":[
    {"claim_id":"c-contested","subject_id":"a","relation":"DERIVED_FROM","object_id":"b","statement":"contested","outcome":"CONTESTED","supporting_evidence_ids":["e1"],"conflicting_evidence_ids":["e2"],"ambiguous_evidence_ids":[],"supporting_inference_ids":[],"conflicting_inference_ids":[],"ambiguous_inference_ids":[],"distinct_supporting_origins":["origin:one"],"reason":"evidence_contested"},
    {"claim_id":"c-unknown-a","subject_id":"a","relation":"RESEMBLES","object_id":"b","statement":"ambiguous","outcome":"UNKNOWN","supporting_evidence_ids":[],"conflicting_evidence_ids":[],"ambiguous_evidence_ids":["e3"],"supporting_inference_ids":["i1"],"conflicting_inference_ids":[],"ambiguous_inference_ids":[],"distinct_supporting_origins":[],"reason":"evidence_ambiguous"},
    {"claim_id":"c-contradicted","subject_id":"a","relation":"IMPROVES","object_id":"b","statement":"contradicted","outcome":"CONTRADICTED","supporting_evidence_ids":[],"conflicting_evidence_ids":["e4"],"ambiguous_evidence_ids":[],"supporting_inference_ids":[],"conflicting_inference_ids":[],"ambiguous_inference_ids":[],"distinct_supporting_origins":[],"reason":"evidence_contradicts_claim"},
    {"claim_id":"c-corroborated","subject_id":"a","relation":"GENERALIZES","object_id":"b","statement":"corroborated","outcome":"CORROBORATED","supporting_evidence_ids":["e5","e6"],"conflicting_evidence_ids":[],"ambiguous_evidence_ids":[],"supporting_inference_ids":[],"conflicting_inference_ids":[],"ambiguous_inference_ids":[],"distinct_supporting_origins":["origin:one","origin:two"],"reason":"independent_origins_corroborate"},
    {"claim_id":"c-supported","subject_id":"a","relation":"SPECIALIZES","object_id":"b","statement":"supported","outcome":"SUPPORTED","supporting_evidence_ids":["e7"],"conflicting_evidence_ids":[],"ambiguous_evidence_ids":[],"supporting_inference_ids":[],"conflicting_inference_ids":[],"ambiguous_inference_ids":[],"distinct_supporting_origins":["origin:one"],"reason":"single_origin_support"},
    {"claim_id":"c-unknown-i","subject_id":"a","relation":"COMBINES","object_id":"b","statement":"insufficient","outcome":"UNKNOWN","supporting_evidence_ids":[],"conflicting_evidence_ids":[],"ambiguous_evidence_ids":[],"supporting_inference_ids":[],"conflicting_inference_ids":[],"ambiguous_inference_ids":[],"distinct_supporting_origins":[],"reason":"insufficient_knowledge"}
  ],
  "metrics":{"input_algorithms":1,"input_provenance":2,"input_claims":3,"input_bases":4,"algorithms":5,"provenance":6,"claims":7,"bases":8,"duplicate_algorithms":9,"duplicate_provenance":10,"duplicate_claims":11,"duplicate_bases":12,"evidence":13,"inferences":14,"corroborated":15,"supported":16,"contested":17,"contradicted":18,"unknown":19,"operations":20,"limit_exceeded":true,"verdict_changes":0}
}`
	var report domain.AlgorithmHistoryReport
	got := roundTripHistoryValue(t, complete, &report)
	requireExactHistoryShape(t, got, decodeHistoryObject(t, complete))
	if got["schema_version"] != "algorithm-history-report-v1" || got["status"] != "COMPLETE" {
		t.Fatalf("report schema/status constants changed: schema=%v status=%v", got["schema_version"], got["status"])
	}
	assertHistoryTokens(t, got["assessments"], "outcome", []string{"CONTESTED", "CONTRADICTED", "CORROBORATED", "SUPPORTED", "UNKNOWN"})
	assertHistoryTokens(t, got["assessments"], "reason", []string{"evidence_ambiguous", "evidence_contested", "evidence_contradicts_claim", "independent_origins_corroborate", "insufficient_knowledge", "single_origin_support"})

	const unknown = `{"schema_version":"algorithm-history-report-v1","status":"UNKNOWN","reason":"invalid_input","assessments":[],"metrics":{"input_algorithms":0,"input_provenance":0,"input_claims":0,"input_bases":0,"algorithms":0,"provenance":0,"claims":0,"bases":0,"duplicate_algorithms":0,"duplicate_provenance":0,"duplicate_claims":0,"duplicate_bases":0,"evidence":0,"inferences":0,"corroborated":0,"supported":0,"contested":0,"contradicted":0,"unknown":0,"operations":0,"limit_exceeded":false,"verdict_changes":0}}`
	var unknownReport domain.AlgorithmHistoryReport
	unknownGot := roundTripHistoryValue(t, unknown, &unknownReport)
	requireExactHistoryShape(t, unknownGot, decodeHistoryObject(t, unknown))
	if unknownGot["status"] != "UNKNOWN" {
		t.Fatalf("UNKNOWN status constant changed: %v", unknownGot["status"])
	}
}
