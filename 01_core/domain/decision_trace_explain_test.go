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
	"strings"
	"testing"
)

func TestDecisionTraceExplainContract(t *testing.T) {
	exp := domain.DecisionTraceExplanation{
		SchemaVersion:     domain.DecisionTraceExplanationSchemaVersion,
		GraphHash:         strings.Repeat("a", 64),
		ReportFingerprint: strings.Repeat("b", 64),
		Scope:             domain.DecisionTraceScope,
		Entries: []domain.DecisionTraceExplanationEntry{
			{
				NodeID:      strings.Repeat("c", 64),
				Kind:        domain.TraceClaim,
				Label:       "Razão declarada",
				Fingerprint: strings.Repeat("d", 64),
				Statement:   "Registro consultivo de uma razão do relatório",
				Locations: []domain.DecisionTraceLocation{
					{Document: "report", Pointer: "/reasons/0"},
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
	if restored.Entries[0].Locations[0].Pointer != "/reasons/0" {
		t.Fatalf("location pointer mismatch: %+v", restored)
	}
}
