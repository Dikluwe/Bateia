/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_compare.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain

const DecisionTraceComparisonSchemaVersion = "decision-trace-comparison-v1"

// Before and After are supplied snapshots, not observed source events.
type DecisionTraceCompareInput struct {
	SchemaVersion string             `json:"schema_version"`
	Before        DecisionTraceGraph `json:"before"`
	After         DecisionTraceGraph `json:"after"`
}

type DecisionTraceNodeDelta struct {
	NodeID string             `json:"node_id"`
	Status string             `json:"status"`
	Before *DecisionTraceNode `json:"before,omitempty"`
	After  *DecisionTraceNode `json:"after,omitempty"`
}

// CandidateDecision is a declaration from After, never a promoted verdict.
// Impact remains anchored to Before and preserves its historical decision.
type DecisionTraceComparison struct {
	SchemaVersion      string                   `json:"schema_version"`
	BeforeGraphHash    string                   `json:"before_graph_hash"`
	AfterGraphHash     string                   `json:"after_graph_hash"`
	Scope              string                   `json:"scope"`
	HistoricalDecision ReuseDecision            `json:"historical_decision"`
	CandidateDecision  ReuseDecision            `json:"candidate_decision"`
	DecisionImpact     string                   `json:"decision_impact"`
	MetadataChanged    bool                     `json:"metadata_changed"`
	NodeChanges        []DecisionTraceNodeDelta `json:"node_changes"`
	AddedEdges         []DecisionTraceEdge      `json:"added_edges"`
	RemovedEdges       []DecisionTraceEdge      `json:"removed_edges"`
	Impact             DecisionTraceImpactPlan  `json:"impact"`
	Limitations        []string                 `json:"limitations"`
}
