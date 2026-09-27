/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_compare.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"fmt"
	"reflect"
	"sort"

	"bateia/01_core/domain"
)

// CompareDecisionTraces compares frozen declarations and plans review against
// the historical graph. The candidate decision is supplied, never promoted.
func CompareDecisionTraces(input domain.DecisionTraceCompareInput) (domain.DecisionTraceComparison, error) {
	fail := func(err error) (domain.DecisionTraceComparison, error) {
		return domain.DecisionTraceComparison{}, err
	}
	if _, err := traceBoundedJSON(input); err != nil {
		return fail(err)
	}
	if input.SchemaVersion != "1" {
		return fail(fmt.Errorf("invalid comparison schema"))
	}
	if input.Before.GraphHash == "" || input.After.GraphHash == "" {
		return fail(fmt.Errorf("comparison requires two frozen graph hashes"))
	}
	before, err := FreezeDecisionTrace(input.Before)
	if err != nil {
		return fail(fmt.Errorf("invalid historical graph: %w", err))
	}
	after, err := FreezeDecisionTrace(input.After)
	if err != nil {
		return fail(fmt.Errorf("invalid candidate graph: %w", err))
	}
	result := domain.DecisionTraceComparison{
		SchemaVersion:   domain.DecisionTraceComparisonSchemaVersion,
		BeforeGraphHash: before.GraphHash, AfterGraphHash: after.GraphHash,
		Scope:              domain.DecisionTraceScope,
		HistoricalDecision: before.Decision, CandidateDecision: after.Decision,
		DecisionImpact:  "NONE",
		MetadataChanged: !reflect.DeepEqual(before.Limitations, after.Limitations),
		NodeChanges:     []domain.DecisionTraceNodeDelta{},
		AddedEdges:      []domain.DecisionTraceEdge{}, RemovedEdges: []domain.DecisionTraceEdge{},
		Limitations: append(append([]string{}, traceLimitations...),
			"Added and removed nodes or edges describe snapshot declarations, not source expiration, revocation or verified changes in the world.",
			"The candidate decision is a supplied declaration, not a recomputed or approved replacement of the historical decision.",
			"Added declarations and changed graph metadata conservatively require review of the historical decision; traversal uses only historical recorded dependencies."),
	}
	oldNodes := make(map[string]domain.DecisionTraceNode, len(before.Nodes))
	newNodes := make(map[string]domain.DecisionTraceNode, len(after.Nodes))
	decisionID := ""
	for _, node := range before.Nodes {
		oldNodes[node.ID] = node
		if node.Kind == domain.TraceDecision {
			decisionID = node.ID
		}
	}
	for _, node := range after.Nodes {
		newNodes[node.ID] = node
	}
	roots := make(map[string]bool)
	for _, node := range before.Nodes {
		next, exists := newNodes[node.ID]
		if exists && reflect.DeepEqual(node, next) {
			continue
		}
		prior := node
		delta := domain.DecisionTraceNodeDelta{NodeID: node.ID, Status: "REMOVED", Before: &prior}
		if exists {
			delta.Status = "DECLARATION_CHANGED"
			delta.After = &next
		}
		result.NodeChanges = append(result.NodeChanges, delta)
		roots[node.ID] = true
	}
	for _, node := range after.Nodes {
		if _, exists := oldNodes[node.ID]; exists {
			continue
		}
		next := node
		result.NodeChanges = append(result.NodeChanges, domain.DecisionTraceNodeDelta{NodeID: node.ID, Status: "ADDED", After: &next})
		// An added ID cannot be traversed in the historical graph. Do not
		// fabricate an old dependency or silently report no recorded impact.
		roots[decisionID] = true
	}
	sort.Slice(result.NodeChanges, func(i, j int) bool { return result.NodeChanges[i].NodeID < result.NodeChanges[j].NodeID })
	oldEdges := make(map[domain.DecisionTraceEdge]bool, len(before.Edges))
	newEdges := make(map[domain.DecisionTraceEdge]bool, len(after.Edges))
	for _, edge := range before.Edges {
		oldEdges[edge] = true
	}
	for _, edge := range after.Edges {
		newEdges[edge] = true
	}
	markConsumer := func(edge domain.DecisionTraceEdge) {
		if _, exists := oldNodes[edge.To]; exists {
			roots[edge.To] = true
		} else {
			roots[decisionID] = true
		}
	}
	for _, edge := range before.Edges {
		if !newEdges[edge] {
			result.RemovedEdges = append(result.RemovedEdges, edge)
			markConsumer(edge)
		}
	}
	for _, edge := range after.Edges {
		if !oldEdges[edge] {
			result.AddedEdges = append(result.AddedEdges, edge)
			markConsumer(edge)
		}
	}
	if result.MetadataChanged {
		roots[decisionID] = true
	}
	if len(roots) > traceMaxChanges {
		return fail(fmt.Errorf("comparison exceeds historical review root budget (%d)", traceMaxChanges))
	}
	changes := make([]domain.DecisionTraceChange, 0, len(roots))
	for id := range roots {
		// Snapshot differences alone cannot establish CONTENT_CHANGED,
		// EXPIRED, REVOKED or a rule/policy change in the real world.
		changes = append(changes, domain.DecisionTraceChange{NodeID: id, Reason: "INPUT_CHANGED"})
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].NodeID < changes[j].NodeID })
	result.Impact, err = PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: before, Changes: changes})
	if err != nil {
		return fail(fmt.Errorf("comparison review plan: %w", err))
	}
	if _, err := traceBoundedJSON(result); err != nil {
		return fail(fmt.Errorf("decision trace comparison output exceeds budget: %w", err))
	}
	return result, nil
}
