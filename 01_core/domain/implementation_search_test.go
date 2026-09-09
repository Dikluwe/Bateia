/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/implementation_search_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"testing"
)

func TestImplementationSearchJSONUnknownRoundTrip(t *testing.T) {
	report := &domain.ImplementationSearchReport{Status: "not_evaluated", Method: "go_ast_lexical", Reason: "unsupported input remains unknown", Budget: domain.ImplementationSearchBudget{MaxFiles: 2, MaxBytes: 400, MaxFileBytes: 200, MaxResults: 1}, Limited: true, Issues: []string{"parse unavailable"}, Examined: []domain.ImplementationInput{{Path: "broken.go", SHA256: "abcdef", Bytes: 12}}, FilesRead: 1, BytesRead: 12}
	original := domain.ExtractedRepository{ImplementationSearch: report}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.ExtractedRepository
	if err = json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(report, restored.ImplementationSearch) {
		t.Fatalf("lost report: %s", data)
	}
	bundle := domain.EvidenceBundle{ImplementationSearch: report}
	data, err = json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	var restoredBundle domain.EvidenceBundle
	if err = json.Unmarshal(data, &restoredBundle); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(report, restoredBundle.ImplementationSearch) {
		t.Fatalf("lost bundle report: %s", data)
	}
}
