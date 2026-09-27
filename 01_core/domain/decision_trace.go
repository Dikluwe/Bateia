/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_trace.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain

const DecisionTraceSchemaVersion = "decision-trace-v1"
const DecisionTraceScope = "RECORDED_DEPENDENCIES_ONLY"

type DecisionTraceNodeKind string

const (
	TraceInput    DecisionTraceNodeKind = "INPUT"
	TraceSource   DecisionTraceNodeKind = "SOURCE"
	TraceClaim    DecisionTraceNodeKind = "CLAIM"
	TraceDecision DecisionTraceNodeKind = "DECISION"
)

type DecisionTraceRelation string

const (
	TraceDependsOn DecisionTraceRelation = "DEPENDS_ON"
	TraceSupports  DecisionTraceRelation = "SUPPORTS"
	TraceRefutes   DecisionTraceRelation = "REFUTES"
	TraceInforms   DecisionTraceRelation = "INFORMS"
)

// Edges point from a dependency to its consumer, including refutations.
// These records do not authenticate sources or assert completeness/proof.
type DecisionTraceNode struct {
	ID          string                `json:"id"`
	Kind        DecisionTraceNodeKind `json:"kind"`
	Label       string                `json:"label"`
	Fingerprint string                `json:"fingerprint"`
	Source      *EvidenceSource       `json:"source,omitempty"`
}

type DecisionTraceEdge struct {
	From     string                `json:"from"`
	To       string                `json:"to"`
	Relation DecisionTraceRelation `json:"relation"`
}

type DecisionTraceGraph struct {
	SchemaVersion string              `json:"schema_version"`
	GraphHash     string              `json:"graph_hash"`
	Decision      ReuseDecision       `json:"decision"`
	Scope         string              `json:"scope"`
	Nodes         []DecisionTraceNode `json:"nodes"`
	Edges         []DecisionTraceEdge `json:"edges"`
	Limitations   []string            `json:"limitations"`
}

type DecisionTraceBuildInput struct {
	SchemaVersion string              `json:"schema_version"`
	Idea          *IdeaSpec           `json:"idea,omitempty"`
	Report        ReuseDecisionReport `json:"report"`
	RuleVersions  map[string]string   `json:"rule_versions,omitempty"`
}

type DecisionTraceChange struct {
	NodeID         string `json:"node_id"`
	Reason         string `json:"reason"`
	NewFingerprint string `json:"new_fingerprint,omitempty"`
}

type DecisionTraceImpactInput struct {
	SchemaVersion string                `json:"schema_version"`
	Graph         DecisionTraceGraph    `json:"graph"`
	Changes       []DecisionTraceChange `json:"changes"`
}

type DecisionTraceReviewTask struct {
	NodeID             string                `json:"node_id"`
	Kind               DecisionTraceNodeKind `json:"kind"`
	ChangedAncestorIDs []string              `json:"changed_ancestor_ids"`
	Reason             string                `json:"reason"`
}

type DecisionTraceImpactPlan struct {
	SchemaVersion     string                    `json:"schema_version"`
	GraphHash         string                    `json:"graph_hash"`
	Status            string                    `json:"status"`
	Scope             string                    `json:"scope"`
	DecisionBefore    ReuseDecision             `json:"decision_before"`
	DecisionAfter     ReuseDecision             `json:"decision_after"`
	DecisionImpact    string                    `json:"decision_impact"`
	ChangedNodeIDs    []string                  `json:"changed_node_ids"`
	ImpactedNodeIDs   []string                  `json:"impacted_node_ids"`
	UnaffectedNodeIDs []string                  `json:"unaffected_node_ids"`
	Tasks             []DecisionTraceReviewTask `json:"tasks"`
	Limitations       []string                  `json:"limitations"`
}
