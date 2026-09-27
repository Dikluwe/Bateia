/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"testing"
)

func TestDecisionTraceGraphContract(t *testing.T) {
	graph := domain.DecisionTraceGraph{
		SchemaVersion: domain.DecisionTraceSchemaVersion,
		GraphHash:     "hash-123",
		Decision:      domain.DecisionFork,
		Scope:         domain.DecisionTraceScope,
		Nodes: []domain.DecisionTraceNode{
			{ID: "input-1", Kind: domain.TraceInput, Label: "IdeaSpec", Fingerprint: "fp-1"},
			{ID: "decision-1", Kind: domain.TraceDecision, Label: "FORK", Fingerprint: "fp-2"},
		},
		Edges: []domain.DecisionTraceEdge{
			{From: "input-1", To: "decision-1", Relation: domain.TraceDependsOn},
		},
		Limitations: []string{"consultivo"},
	}
	data, err := json.Marshal(graph)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.DecisionTraceGraph
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.SchemaVersion != domain.DecisionTraceSchemaVersion || len(restored.Nodes) != 2 || len(restored.Edges) != 1 {
		t.Fatalf("unexpected restored graph: %+v", restored)
	}
	if restored.Nodes[0].Kind != domain.TraceInput || restored.Edges[0].Relation != domain.TraceDependsOn {
		t.Fatalf("node or edge corrupted: %+v", restored)
	}
}
