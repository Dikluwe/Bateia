/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_review.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain

const DecisionTraceReviewSchemaVersion = "decision-trace-review-v1"

// References and hashes are supplied declarations, never opened or authenticated.
type DecisionTraceReviewEvidence struct {
	Ref         string `json:"ref"`
	ContentHash string `json:"content_hash"`
}

type DecisionTraceReviewRecord struct {
	NodeID   string                        `json:"node_id"`
	PlanHash string                        `json:"plan_hash"`
	TaskHash string                        `json:"task_hash"`
	Status   string                        `json:"status"`
	Note     string                        `json:"note"`
	Evidence []DecisionTraceReviewEvidence `json:"evidence"`
}

type DecisionTraceReviewInput struct {
	SchemaVersion string                      `json:"schema_version"`
	ImpactInput   DecisionTraceImpactInput    `json:"impact_input"`
	Records       []DecisionTraceReviewRecord `json:"records,omitempty"`
}

type DecisionTraceReviewProgress struct {
	Task     DecisionTraceReviewTask       `json:"task"`
	TaskHash string                        `json:"task_hash"`
	Status   string                        `json:"status"`
	Note     string                        `json:"note"`
	Evidence []DecisionTraceReviewEvidence `json:"evidence"`
}

type DecisionTraceReviewCounts struct {
	Pending    int `json:"pending"`
	InProgress int `json:"in_progress"`
	Blocked    int `json:"blocked"`
	Completed  int `json:"completed"`
}

// All progress is declared. Completion never approves or recalculates a decision.
type DecisionTraceReviewLedger struct {
	SchemaVersion  string                        `json:"schema_version"`
	PlanHash       string                        `json:"plan_hash"`
	GraphHash      string                        `json:"graph_hash"`
	Scope          string                        `json:"scope"`
	DecisionBefore ReuseDecision                 `json:"decision_before"`
	DecisionAfter  ReuseDecision                 `json:"decision_after"`
	DecisionImpact string                        `json:"decision_impact"`
	Status         string                        `json:"status"`
	Counts         DecisionTraceReviewCounts     `json:"counts"`
	Tasks          []DecisionTraceReviewProgress `json:"tasks"`
	Limitations    []string                      `json:"limitations"`
}
