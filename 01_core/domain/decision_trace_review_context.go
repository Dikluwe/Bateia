/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_review_context.md
 * @layer L1
 * @updated 2026-09-27
 */
package domain

const DecisionTraceReviewContextSchemaVersion = "decision-trace-review-context-v1"
const DecisionTraceReviewContextScope = "TASK_AND_CHANGED_ANCESTORS"

type DecisionTraceReviewContextInput struct {
	SchemaVersion string                   `json:"schema_version"`
	ReviewInput   DecisionTraceReviewInput `json:"review_input"`
	BuildInput    DecisionTraceBuildInput  `json:"build_input"`
}

type DecisionTraceReviewTaskContext struct {
	NodeID         string   `json:"node_id"`
	ContextNodeIDs []string `json:"context_node_ids"`
}

// Contexts link declared tasks to historical documents, never to verified facts.
type DecisionTraceReviewContext struct {
	SchemaVersion     string                           `json:"schema_version"`
	GraphHash         string                           `json:"graph_hash"`
	PlanHash          string                           `json:"plan_hash"`
	ReportFingerprint string                           `json:"report_fingerprint"`
	Scope             string                           `json:"scope"`
	ContextScope      string                           `json:"context_scope"`
	Ledger            DecisionTraceReviewLedger        `json:"ledger"`
	TaskContexts      []DecisionTraceReviewTaskContext `json:"task_contexts"`
	Entries           []DecisionTraceExplanationEntry  `json:"entries"`
	Limitations       []string                         `json:"limitations"`
}
