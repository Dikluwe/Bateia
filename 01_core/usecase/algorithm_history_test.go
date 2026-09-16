/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/algorithm_history_test.md
 * @layer L1
 * @updated 2026-09-12
 */

package usecase_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

const historyBaseline = `{
  "schema_version":"algorithm-history-input-v1",
  "algorithms":[
    {"id":"alpha-sort","name":"Alpha Sort","chronology":"KNOWN","year_start":1960,"year_end":1960,"chronology_reason":"dated design record","provenance_ids":["archive-alpha"]},
    {"id":"beta-sort","name":"Beta Sort","chronology":"APPROXIMATE","year_start":1970,"year_end":1972,"chronology_reason":"bounded archive records","provenance_ids":["archive-beta","paper-beta"]},
    {"id":"gamma-sort","name":"Gamma Sort","chronology":"UNKNOWN","year_start":null,"year_end":null,"chronology_reason":"surviving source is undated","provenance_ids":["archive-gamma"]}
  ],
  "provenance":[
    {"id":"archive-alpha","kind":"PRIMARY_SOURCE","origin_id":"archive:alpha","citation":"Synthetic Alpha design record","locator":"archive://alpha/design"},
    {"id":"archive-beta","kind":"PRIMARY_SOURCE","origin_id":"archive:beta","citation":"Synthetic Beta notebook","locator":"archive://beta/notebook"},
    {"id":"paper-beta","kind":"SECONDARY_SOURCE","origin_id":"journal:beta","citation":"Synthetic independent Beta study","locator":"doi:10.0000/example.beta"},
    {"id":"archive-gamma","kind":"EMPIRICAL_ARTIFACT","origin_id":"artifact:gamma","citation":"Synthetic Gamma source artifact","locator":"artifact://gamma/source"}
  ],
  "claims":[
    {"id":"beta-from-alpha","subject_id":"beta-sort","relation":"DERIVED_FROM","object_id":"alpha-sort","statement":"Beta Sort derived its partition schedule from Alpha Sort"},
    {"id":"gamma-improves-beta","subject_id":"gamma-sort","relation":"IMPROVES","object_id":"beta-sort","statement":"Gamma Sort improves the bounded-memory behavior of Beta Sort"},
    {"id":"gamma-resembles-alpha","subject_id":"gamma-sort","relation":"RESEMBLES","object_id":"alpha-sort","statement":"Gamma Sort resembles Alpha Sort in its partition boundary"}
  ],
  "bases":[
    {"id":"beta-support-primary","claim_id":"beta-from-alpha","kind":"EVIDENCE","stance":"SUPPORT","provenance_id":"archive-beta","summary":"Notebook attributes the schedule to Alpha","derived_from_ids":[]},
    {"id":"beta-support-secondary","claim_id":"beta-from-alpha","kind":"EVIDENCE","stance":"SUPPORT","provenance_id":"paper-beta","summary":"Independent study reports the derivation","derived_from_ids":[]},
    {"id":"beta-inference","claim_id":"beta-from-alpha","kind":"INFERENCE","stance":"SUPPORT","provenance_id":"paper-beta","summary":"Ordering of changes is consistent with derivation","derived_from_ids":["beta-support-primary"]},
    {"id":"gamma-improve-support","claim_id":"gamma-improves-beta","kind":"EVIDENCE","stance":"SUPPORT","provenance_id":"archive-gamma","summary":"Artifact uses less bounded workspace in the fixture","derived_from_ids":[]},
    {"id":"gamma-improve-conflict","claim_id":"gamma-improves-beta","kind":"EVIDENCE","stance":"CONFLICT","provenance_id":"paper-beta","summary":"Study reports no improvement under its workload","derived_from_ids":[]},
    {"id":"gamma-resemblance-ambiguous","claim_id":"gamma-resembles-alpha","kind":"EVIDENCE","stance":"AMBIGUOUS","provenance_id":"archive-alpha","summary":"Boundary notation admits more than one reading","derived_from_ids":[]},
    {"id":"gamma-resemblance-inference","claim_id":"gamma-resembles-alpha","kind":"INFERENCE","stance":"SUPPORT","provenance_id":"archive-gamma","summary":"Structural similarity is inferred from the artifact","derived_from_ids":["gamma-resemblance-ambiguous"]}
  ]
}`

type historyReportView struct {
	SchemaVersion string                  `json:"schema_version"`
	Status        string                  `json:"status"`
	Reason        string                  `json:"reason"`
	Assessments   []historyAssessmentView `json:"assessments"`
	Metrics       historyMetricsView      `json:"metrics"`
}

type historyAssessmentView struct {
	ClaimID                   string   `json:"claim_id"`
	SubjectID                 string   `json:"subject_id"`
	Relation                  string   `json:"relation"`
	ObjectID                  string   `json:"object_id"`
	Statement                 string   `json:"statement"`
	Outcome                   string   `json:"outcome"`
	SupportingEvidenceIDs     []string `json:"supporting_evidence_ids"`
	ConflictingEvidenceIDs    []string `json:"conflicting_evidence_ids"`
	AmbiguousEvidenceIDs      []string `json:"ambiguous_evidence_ids"`
	SupportingInferenceIDs    []string `json:"supporting_inference_ids"`
	ConflictingInferenceIDs   []string `json:"conflicting_inference_ids"`
	AmbiguousInferenceIDs     []string `json:"ambiguous_inference_ids"`
	DistinctSupportingOrigins []string `json:"distinct_supporting_origins"`
	Reason                    string   `json:"reason"`
}

type historyMetricsView struct {
	InputAlgorithms     int  `json:"input_algorithms"`
	InputProvenance     int  `json:"input_provenance"`
	InputClaims         int  `json:"input_claims"`
	InputBases          int  `json:"input_bases"`
	Algorithms          int  `json:"algorithms"`
	Provenance          int  `json:"provenance"`
	Claims              int  `json:"claims"`
	Bases               int  `json:"bases"`
	DuplicateAlgorithms int  `json:"duplicate_algorithms"`
	DuplicateProvenance int  `json:"duplicate_provenance"`
	DuplicateClaims     int  `json:"duplicate_claims"`
	DuplicateBases      int  `json:"duplicate_bases"`
	Evidence            int  `json:"evidence"`
	Inferences          int  `json:"inferences"`
	Corroborated        int  `json:"corroborated"`
	Supported           int  `json:"supported"`
	Contested           int  `json:"contested"`
	Contradicted        int  `json:"contradicted"`
	Unknown             int  `json:"unknown"`
	Operations          int  `json:"operations"`
	LimitExceeded       bool `json:"limit_exceeded"`
	VerdictChanges      int  `json:"verdict_changes"`
}

func decodeHistoryInput(t *testing.T, raw []byte) domain.AlgorithmHistoryInput {
	t.Helper()
	var input domain.AlgorithmHistoryInput
	if err := json.Unmarshal(raw, &input); err != nil {
		t.Fatalf("decode public input fixture: %v", err)
	}
	return input
}

func analyzeHistory(t *testing.T, input domain.AlgorithmHistoryInput) (historyReportView, any) {
	t.Helper()
	report := usecase.AnalyzeAlgorithmHistory(input)
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal public report: %v", err)
	}
	var view historyReportView
	if err := json.Unmarshal(raw, &view); err != nil {
		t.Fatalf("decode public report: %v", err)
	}
	var generic any
	if err := json.Unmarshal(raw, &generic); err != nil {
		t.Fatalf("decode comparable report: %v", err)
	}
	return view, generic
}

func baselineObject(t *testing.T) map[string]any {
	t.Helper()
	var object map[string]any
	if err := json.Unmarshal([]byte(historyBaseline), &object); err != nil {
		t.Fatal(err)
	}
	return object
}

func marshalObject(t *testing.T, object map[string]any) []byte {
	t.Helper()
	raw, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func assessmentByID(t *testing.T, report historyReportView, id string) historyAssessmentView {
	t.Helper()
	for _, assessment := range report.Assessments {
		if assessment.ClaimID == id {
			return assessment
		}
	}
	t.Fatalf("assessment %q not found", id)
	return historyAssessmentView{}
}

func TestAnalyzeAlgorithmHistoryBaseline53Operations(t *testing.T) {
	input := decodeHistoryInput(t, []byte(historyBaseline))
	before, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	report, _ := analyzeHistory(t, input)
	after, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) {
		t.Fatal("AnalyzeAlgorithmHistory mutated the input or a reachable slice")
	}
	if report.SchemaVersion != "algorithm-history-report-v1" || report.Status != "COMPLETE" || report.Reason != "analyzed" {
		t.Fatalf("unexpected envelope: %+v", report)
	}
	if len(report.Assessments) != 3 {
		t.Fatalf("got %d assessments, want 3", len(report.Assessments))
	}
	beta := assessmentByID(t, report, "beta-from-alpha")
	if beta.Relation != "DERIVED_FROM" || beta.Outcome != "CORROBORATED" || beta.Reason != "independent_origins_corroborate" {
		t.Fatalf("beta assessment: %+v", beta)
	}
	if !reflect.DeepEqual(beta.SupportingEvidenceIDs, []string{"beta-support-primary", "beta-support-secondary"}) || !reflect.DeepEqual(beta.SupportingInferenceIDs, []string{"beta-inference"}) || !reflect.DeepEqual(beta.DistinctSupportingOrigins, []string{"archive:beta", "journal:beta"}) {
		t.Fatalf("beta evidence/inference separation: %+v", beta)
	}
	improve := assessmentByID(t, report, "gamma-improves-beta")
	if improve.Outcome != "CONTESTED" || improve.Reason != "evidence_contested" || !reflect.DeepEqual(improve.SupportingEvidenceIDs, []string{"gamma-improve-support"}) || !reflect.DeepEqual(improve.ConflictingEvidenceIDs, []string{"gamma-improve-conflict"}) {
		t.Fatalf("contested assessment: %+v", improve)
	}
	resembles := assessmentByID(t, report, "gamma-resembles-alpha")
	if resembles.Relation != "RESEMBLES" || resembles.Outcome != "UNKNOWN" || resembles.Reason != "evidence_ambiguous" || !reflect.DeepEqual(resembles.SupportingInferenceIDs, []string{"gamma-resemblance-inference"}) {
		t.Fatalf("epistemic UNKNOWN assessment: %+v", resembles)
	}
	wantMetrics := historyMetricsView{InputAlgorithms: 3, InputProvenance: 4, InputClaims: 3, InputBases: 7, Algorithms: 3, Provenance: 4, Claims: 3, Bases: 7, Evidence: 5, Inferences: 2, Corroborated: 1, Contested: 1, Unknown: 1, Operations: 53}
	if report.Metrics != wantMetrics {
		t.Fatalf("metrics mismatch:\n got %+v\nwant %+v", report.Metrics, wantMetrics)
	}
	assertNonNullAssessmentLists(t, usecase.AnalyzeAlgorithmHistory(input))
}

func assertNonNullAssessmentLists(t *testing.T, report any) {
	t.Helper()
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		t.Fatal(err)
	}
	assessments, ok := object["assessments"].([]any)
	if !ok {
		t.Fatalf("assessments must be a non-null array: %s", raw)
	}
	listFields := []string{"supporting_evidence_ids", "conflicting_evidence_ids", "ambiguous_evidence_ids", "supporting_inference_ids", "conflicting_inference_ids", "ambiguous_inference_ids", "distinct_supporting_origins"}
	for _, rawAssessment := range assessments {
		assessment := rawAssessment.(map[string]any)
		for _, field := range listFields {
			if _, ok := assessment[field].([]any); !ok {
				t.Errorf("assessment %v field %s must be a non-null array", assessment["claim_id"], field)
			}
		}
	}
}

func minimalHistory(relations []string) map[string]any {
	algorithms := []any{map[string]any{"id": "root", "name": "Root", "chronology": "KNOWN", "year_start": 1, "year_end": 1, "chronology_reason": "dated", "provenance_ids": []any{"p"}}}
	claims := make([]any, 0, len(relations))
	bases := make([]any, 0, len(relations))
	for i, relation := range relations {
		id := fmt.Sprintf("a-%02d", i)
		claimID := fmt.Sprintf("c-%02d", i)
		algorithms = append(algorithms, map[string]any{"id": id, "name": id, "chronology": "KNOWN", "year_start": 2 + i, "year_end": 2 + i, "chronology_reason": "dated", "provenance_ids": []any{"p"}})
		claims = append(claims, map[string]any{"id": claimID, "subject_id": id, "relation": relation, "object_id": "root", "statement": "explicit " + relation})
		bases = append(bases, map[string]any{"id": "b-" + claimID, "claim_id": claimID, "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "p", "summary": "explicit source", "derived_from_ids": []any{}})
	}
	return map[string]any{
		"schema_version": "algorithm-history-input-v1",
		"algorithms":     algorithms,
		"provenance":     []any{map[string]any{"id": "p", "kind": "PRIMARY_SOURCE", "origin_id": "origin:one", "citation": "citation", "locator": "archive://one"}},
		"claims":         claims,
		"bases":          bases,
	}
}

func TestAnalyzeAlgorithmHistoryPreservesAllEightExplicitRelations(t *testing.T) {
	relations := []string{"DERIVED_FROM", "INFLUENCED_BY", "IMPROVES", "GENERALIZES", "SPECIALIZES", "COMBINES", "REDISCOVERS", "RESEMBLES"}
	report, _ := analyzeHistory(t, decodeHistoryInput(t, marshalObject(t, minimalHistory(relations))))
	if report.Status != "COMPLETE" || len(report.Assessments) != len(relations) {
		t.Fatalf("unexpected report: %+v", report)
	}
	got := make([]string, 0, len(report.Assessments))
	for _, assessment := range report.Assessments {
		got = append(got, assessment.Relation)
		if assessment.Outcome != "SUPPORTED" || assessment.Reason != "single_origin_support" {
			t.Errorf("relation %s outcome: %+v", assessment.Relation, assessment)
		}
	}
	sort.Strings(got)
	want := append([]string(nil), relations...)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("relations changed or inferred: got %v want %v", got, want)
	}
}

func TestAnalyzeAlgorithmHistoryOutcomePrecedence(t *testing.T) {
	object := minimalHistory([]string{"DERIVED_FROM", "IMPROVES", "GENERALIZES", "SPECIALIZES", "RESEMBLES"})
	object["provenance"] = []any{
		map[string]any{"id": "p", "kind": "PRIMARY_SOURCE", "origin_id": "origin:one", "citation": "one", "locator": "archive://one"},
		map[string]any{"id": "q", "kind": "SECONDARY_SOURCE", "origin_id": "origin:two", "citation": "two", "locator": "archive://two"},
	}
	object["bases"] = []any{
		map[string]any{"id": "c0-s", "claim_id": "c-00", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "p", "summary": "support", "derived_from_ids": []any{}},
		map[string]any{"id": "c0-c", "claim_id": "c-00", "kind": "EVIDENCE", "stance": "CONFLICT", "provenance_id": "q", "summary": "conflict", "derived_from_ids": []any{}},
		map[string]any{"id": "c1-a", "claim_id": "c-01", "kind": "EVIDENCE", "stance": "AMBIGUOUS", "provenance_id": "p", "summary": "ambiguous", "derived_from_ids": []any{}},
		map[string]any{"id": "c1-c", "claim_id": "c-01", "kind": "EVIDENCE", "stance": "CONFLICT", "provenance_id": "q", "summary": "conflict too", "derived_from_ids": []any{}},
		map[string]any{"id": "c2-c", "claim_id": "c-02", "kind": "EVIDENCE", "stance": "CONFLICT", "provenance_id": "p", "summary": "conflict", "derived_from_ids": []any{}},
		map[string]any{"id": "c3-s1", "claim_id": "c-03", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "p", "summary": "support one", "derived_from_ids": []any{}},
		map[string]any{"id": "c3-s2", "claim_id": "c-03", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "q", "summary": "support two", "derived_from_ids": []any{}},
		map[string]any{"id": "c4-s", "claim_id": "c-04", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "p", "summary": "support", "derived_from_ids": []any{}},
	}
	report, _ := analyzeHistory(t, decodeHistoryInput(t, marshalObject(t, object)))
	want := map[string][2]string{
		"c-00": {"CONTESTED", "evidence_contested"},
		"c-01": {"UNKNOWN", "evidence_ambiguous"},
		"c-02": {"CONTRADICTED", "evidence_contradicts_claim"},
		"c-03": {"CORROBORATED", "independent_origins_corroborate"},
		"c-04": {"SUPPORTED", "single_origin_support"},
	}
	for id, expected := range want {
		got := assessmentByID(t, report, id)
		if got.Outcome != expected[0] || got.Reason != expected[1] {
			t.Errorf("%s: got %s/%s want %s/%s", id, got.Outcome, got.Reason, expected[0], expected[1])
		}
	}
}

func TestAnalyzeAlgorithmHistoryInferenceNeverCounts(t *testing.T) {
	object := minimalHistory([]string{"DERIVED_FROM"})
	object["provenance"] = append(object["provenance"].([]any), map[string]any{"id": "q", "kind": "SECONDARY_SOURCE", "origin_id": "origin:two", "citation": "two", "locator": "archive://two"})
	object["bases"] = []any{
		map[string]any{"id": "e", "claim_id": "c-00", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "p", "summary": "observed", "derived_from_ids": []any{}},
		map[string]any{"id": "i-support", "claim_id": "c-00", "kind": "INFERENCE", "stance": "SUPPORT", "provenance_id": "q", "summary": "interpreted support", "derived_from_ids": []any{"e"}},
		map[string]any{"id": "i-conflict", "claim_id": "c-00", "kind": "INFERENCE", "stance": "CONFLICT", "provenance_id": "q", "summary": "interpreted conflict", "derived_from_ids": []any{"e"}},
	}
	report, _ := analyzeHistory(t, decodeHistoryInput(t, marshalObject(t, object)))
	a := assessmentByID(t, report, "c-00")
	if a.Outcome != "SUPPORTED" || a.Reason != "single_origin_support" || !reflect.DeepEqual(a.SupportingEvidenceIDs, []string{"e"}) || !reflect.DeepEqual(a.SupportingInferenceIDs, []string{"i-support"}) || !reflect.DeepEqual(a.ConflictingInferenceIDs, []string{"i-conflict"}) {
		t.Fatalf("inference changed evidentiary outcome: %+v", a)
	}
}

func TestAnalyzeAlgorithmHistorySameOriginDoesNotCorroborate(t *testing.T) {
	object := minimalHistory([]string{"DERIVED_FROM"})
	object["provenance"] = []any{
		map[string]any{"id": "p", "kind": "PRIMARY_SOURCE", "origin_id": "same:opaque", "citation": "edition one", "locator": "book://one"},
		map[string]any{"id": "q", "kind": "SECONDARY_SOURCE", "origin_id": "same:opaque", "citation": "edition two", "locator": "book://two"},
	}
	object["bases"] = []any{
		map[string]any{"id": "e1", "claim_id": "c-00", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "p", "summary": "first", "derived_from_ids": []any{}},
		map[string]any{"id": "e2", "claim_id": "c-00", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "q", "summary": "second", "derived_from_ids": []any{}},
	}
	report, _ := analyzeHistory(t, decodeHistoryInput(t, marshalObject(t, object)))
	a := assessmentByID(t, report, "c-00")
	if a.Outcome != "SUPPORTED" || !reflect.DeepEqual(a.DistinctSupportingOrigins, []string{"same:opaque"}) {
		t.Fatalf("same physical/editorial origin was counted twice: %+v", a)
	}
}

func TestAnalyzeAlgorithmHistoryCanonicalDedupAndRemapping(t *testing.T) {
	object := minimalHistory([]string{"DERIVED_FROM"})
	algorithms := object["algorithms"].([]any)
	object["algorithms"] = append(algorithms, algorithms[1])
	object["provenance"] = []any{
		map[string]any{"id": "z-prov", "kind": "PRIMARY_SOURCE", "origin_id": "origin:one", "citation": "same", "locator": "archive://same"},
		map[string]any{"id": "a-prov", "kind": "PRIMARY_SOURCE", "origin_id": "origin:one", "citation": "same", "locator": "archive://same"},
	}
	for _, algorithm := range object["algorithms"].([]any) {
		algorithm.(map[string]any)["provenance_ids"] = []any{"z-prov", "a-prov", "z-prov"}
	}
	claim := object["claims"].([]any)[0].(map[string]any)
	claim["id"] = "z-claim"
	duplicateClaim := map[string]any{"id": "a-claim", "subject_id": claim["subject_id"], "relation": claim["relation"], "object_id": claim["object_id"], "statement": claim["statement"]}
	object["claims"] = []any{claim, duplicateClaim}
	object["bases"] = []any{
		map[string]any{"id": "z-base", "claim_id": "z-claim", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "z-prov", "summary": "same evidence", "derived_from_ids": []any{}},
		map[string]any{"id": "a-base", "claim_id": "a-claim", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "a-prov", "summary": "same evidence", "derived_from_ids": nil},
	}
	report, _ := analyzeHistory(t, decodeHistoryInput(t, marshalObject(t, object)))
	if report.Status != "COMPLETE" || len(report.Assessments) != 1 {
		t.Fatalf("dedup report: %+v", report)
	}
	a := report.Assessments[0]
	if a.ClaimID != "a-claim" || !reflect.DeepEqual(a.SupportingEvidenceIDs, []string{"a-base"}) || !reflect.DeepEqual(a.DistinctSupportingOrigins, []string{"origin:one"}) {
		t.Fatalf("references were not remapped to minimum representatives: %+v", a)
	}
	if report.Metrics.DuplicateAlgorithms != 1 || report.Metrics.DuplicateProvenance != 1 || report.Metrics.DuplicateClaims != 1 || report.Metrics.DuplicateBases != 1 {
		t.Fatalf("duplicate metrics: %+v", report.Metrics)
	}
}

func TestAnalyzeAlgorithmHistoryInvalidAmbiguousAndBudgetsFailClosed(t *testing.T) {
	cases := []struct {
		name, reason string
		mutate       func(map[string]any)
		limit        bool
	}{
		{"invalid schema wins", "invalid_schema", func(o map[string]any) { o["schema_version"] = "wrong"; o["algorithms"] = make([]any, 257) }, false},
		{"algorithm collection budget", "budget_exceeded", func(o map[string]any) { o["algorithms"] = make([]any, 257) }, true},
		{"reference list budget", "budget_exceeded", func(o map[string]any) {
			o["algorithms"].([]any)[0].(map[string]any)["provenance_ids"] = make([]any, 65)
		}, true},
		{"id byte budget", "budget_exceeded", func(o map[string]any) { o["algorithms"].([]any)[0].(map[string]any)["id"] = strings.Repeat("a", 129) }, true},
		{"text byte budget", "budget_exceeded", func(o map[string]any) {
			o["claims"].([]any)[0].(map[string]any)["statement"] = strings.Repeat("x", 4097)
		}, true},
		{"no algorithms", "no_algorithms", func(o map[string]any) { o["algorithms"] = []any{} }, false},
		{"no provenance", "no_provenance", func(o map[string]any) { o["provenance"] = []any{} }, false},
		{"no claims", "no_claims", func(o map[string]any) { o["claims"] = []any{} }, false},
		{"no bases", "no_bases", func(o map[string]any) { o["bases"] = []any{} }, false},
		{"orphan claim endpoint", "invalid_input", func(o map[string]any) { o["claims"].([]any)[0].(map[string]any)["object_id"] = "missing" }, false},
		{"claim without canonical evidence", "invalid_input", func(o map[string]any) {
			kept := make([]any, 0)
			for _, rawBase := range o["bases"].([]any) {
				if rawBase.(map[string]any)["claim_id"] != "beta-from-alpha" {
					kept = append(kept, rawBase)
				}
			}
			o["bases"] = kept
		}, false},
		{"inference from inference", "invalid_input", func(o map[string]any) {
			o["bases"].([]any)[2].(map[string]any)["derived_from_ids"] = []any{"gamma-resemblance-inference"}
		}, false},
		{"ambiguous canonical id", "ambiguous_input", func(o map[string]any) {
			a := o["algorithms"].([]any)[0].(map[string]any)
			clone := make(map[string]any)
			for k, v := range a {
				clone[k] = v
			}
			clone["id"] = " ALPHA-SORT "
			clone["name"] = "different"
			o["algorithms"] = append(o["algorithms"].([]any), clone)
		}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object := baselineObject(t)
			tc.mutate(object)
			report, _ := analyzeHistory(t, decodeHistoryInput(t, marshalObject(t, object)))
			if report.Status != "UNKNOWN" || report.Reason != tc.reason || len(report.Assessments) != 0 || report.Assessments == nil || report.Metrics.LimitExceeded != tc.limit {
				t.Fatalf("got status=%s reason=%s assessments=%v limit=%v", report.Status, report.Reason, report.Assessments, report.Metrics.LimitExceeded)
			}
		})
	}
}

func TestAnalyzeAlgorithmHistoryOperationBudget(t *testing.T) {
	object := minimalHistory([]string{"DERIVED_FROM"})
	bases := make([]any, 0, 314)
	derived := make([]any, 64)
	for i := 0; i < 64; i++ {
		id := fmt.Sprintf("e-%02d", i)
		derived[i] = id
		bases = append(bases, map[string]any{"id": id, "claim_id": "c-00", "kind": "EVIDENCE", "stance": "SUPPORT", "provenance_id": "p", "summary": "evidence " + id, "derived_from_ids": []any{}})
	}
	for i := 0; i < 250; i++ {
		id := fmt.Sprintf("i-%03d", i)
		bases = append(bases, map[string]any{"id": id, "claim_id": "c-00", "kind": "INFERENCE", "stance": "SUPPORT", "provenance_id": "p", "summary": "inference " + id, "derived_from_ids": derived})
	}
	object["bases"] = bases
	report, _ := analyzeHistory(t, decodeHistoryInput(t, marshalObject(t, object)))
	if report.Status != "UNKNOWN" || report.Reason != "budget_exceeded" || !report.Metrics.LimitExceeded || len(report.Assessments) != 0 {
		t.Fatalf("operation budget did not fail closed: %+v", report)
	}
}

func TestAnalyzeAlgorithmHistoryTwentyPermutationsNilEmptyAndImmutable(t *testing.T) {
	baseline := baselineObject(t)
	baseInput := decodeHistoryInput(t, marshalObject(t, baseline))
	_, want := analyzeHistory(t, baseInput)
	for iteration := 0; iteration < 20; iteration++ {
		object := baselineObject(t)
		for _, field := range []string{"algorithms", "provenance", "claims", "bases"} {
			items := object[field].([]any)
			shift := iteration % len(items)
			items = append(items[shift:], items[:shift]...)
			if iteration%2 == 1 {
				for left, right := 0, len(items)-1; left < right; left, right = left+1, right-1 {
					items[left], items[right] = items[right], items[left]
				}
			}
			object[field] = items
		}
		for _, rawBase := range object["bases"].([]any) {
			base := rawBase.(map[string]any)
			if ids, ok := base["derived_from_ids"].([]any); ok && len(ids) == 0 && iteration%2 == 1 {
				delete(base, "derived_from_ids")
			}
		}
		raw := marshalObject(t, object)
		input := decodeHistoryInput(t, raw)
		before, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		_, got := analyzeHistory(t, input)
		after, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("permutation %d changed report", iteration)
		}
		if !reflect.DeepEqual(before, after) {
			t.Fatalf("permutation %d mutated input", iteration)
		}
	}
}
