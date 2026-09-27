/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace_explain.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain

const DecisionTraceExplanationSchemaVersion = "decision-trace-explanation-v1"

// Documents are supplied declarations, not paths to fetch or authenticated facts.
// Pointer is a JSON Pointer relative to that document, including "" for its root.
type DecisionTraceLocation struct {
	Document string `json:"document"`
	Pointer  string `json:"pointer"`
}

type DecisionTraceExplainInput struct {
	SchemaVersion string                  `json:"schema_version"`
	Graph         DecisionTraceGraph      `json:"graph"`
	BuildInput    DecisionTraceBuildInput `json:"build_input"`
	NodeIDs       []string                `json:"node_ids,omitempty"`
}

type DecisionTraceExplanationEntry struct {
	NodeID             string                  `json:"node_id"`
	Kind               DecisionTraceNodeKind   `json:"kind"`
	Label              string                  `json:"label"`
	Fingerprint        string                  `json:"fingerprint"`
	Statement          string                  `json:"statement"`
	StatementTruncated bool                    `json:"statement_truncated"`
	Locations          []DecisionTraceLocation `json:"locations"`
}

type DecisionTraceExplanation struct {
	SchemaVersion     string                          `json:"schema_version"`
	GraphHash         string                          `json:"graph_hash"`
	ReportFingerprint string                          `json:"report_fingerprint"`
	Scope             string                          `json:"scope"`
	Entries           []DecisionTraceExplanationEntry `json:"entries"`
	Limitations       []string                        `json:"limitations"`
}
