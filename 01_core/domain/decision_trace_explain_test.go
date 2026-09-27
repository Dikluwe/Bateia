/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_explain_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"testing"
)

func TestDecisionTraceExplainContract(t *testing.T) {
	exp := domain.DecisionTraceExplanation{
		SchemaVersion:     domain.DecisionTraceExplanationSchemaVersion,
		GraphHash:         "gh-123",
		ReportFingerprint: "rf-456",
		Scope:             "LOCAL",
		Entries: []domain.DecisionTraceExplanationEntry{
			{
				NodeID:      "claim-1",
				Kind:        domain.TraceClaim,
				Label:       "Licença MIT",
				Fingerprint: "fp-mit",
				Statement:   "Licença MIT confirmada por análise confinada",
				Locations: []domain.DecisionTraceLocation{
					{Document: "bundles.json", Pointer: "/0/items/0"},
				},
			},
		},
		Limitations: []string{"consultivo"},
	}
	data, err := json.Marshal(exp)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.DecisionTraceExplanation
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.SchemaVersion != domain.DecisionTraceExplanationSchemaVersion || len(restored.Entries) != 1 {
		t.Fatalf("restored explanation mismatch: %+v", restored)
	}
	if restored.Entries[0].Locations[0].Pointer != "/0/items/0" {
		t.Fatalf("location pointer mismatch: %+v", restored)
	}
}
