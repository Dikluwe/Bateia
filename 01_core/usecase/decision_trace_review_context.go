/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_review_context.md
 * @layer L1
 * @updated 2026-09-27
 */
package usecase

import (
	"fmt"
	"sort"

	"bateia/01_core/domain"
)

// BuildDecisionTraceReviewContext binds progress to original recorded documents.
// A task context contains itself and changed roots, not a complete causal path.
func BuildDecisionTraceReviewContext(input domain.DecisionTraceReviewContextInput) (domain.DecisionTraceReviewContext, error) {
	fail := func(err error) (domain.DecisionTraceReviewContext, error) {
		return domain.DecisionTraceReviewContext{}, err
	}
	if input.SchemaVersion != "1" {
		return fail(fmt.Errorf("invalid review context schema"))
	}
	if _, err := traceBoundedJSON(input); err != nil {
		return fail(err)
	}
	ledger, err := TrackDecisionTraceReview(input.ReviewInput)
	if err != nil {
		return fail(fmt.Errorf("invalid review context progress: %w", err))
	}
	contexts := make([]domain.DecisionTraceReviewTaskContext, 0, len(ledger.Tasks))
	selected := make(map[string]bool)
	for _, progress := range ledger.Tasks {
		ids := map[string]bool{progress.Task.NodeID: true}
		for _, id := range progress.Task.ChangedAncestorIDs {
			ids[id] = true
		}
		context := domain.DecisionTraceReviewTaskContext{NodeID: progress.Task.NodeID, ContextNodeIDs: []string{}}
		for id := range ids {
			context.ContextNodeIDs = append(context.ContextNodeIDs, id)
			selected[id] = true
		}
		sort.Strings(context.ContextNodeIDs)
		contexts = append(contexts, context)
	}
	nodeIDs := make([]string, 0, len(selected))
	for id := range selected {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)
	explanation, err := ExplainDecisionTrace(domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: input.ReviewInput.ImpactInput.Graph, BuildInput: input.BuildInput, NodeIDs: nodeIDs})
	if err != nil {
		return fail(fmt.Errorf("invalid review context documents: %w", err))
	}
	result := domain.DecisionTraceReviewContext{
		SchemaVersion: domain.DecisionTraceReviewContextSchemaVersion,
		GraphHash:     ledger.GraphHash, PlanHash: ledger.PlanHash, ReportFingerprint: explanation.ReportFingerprint,
		Scope: domain.DecisionTraceScope, ContextScope: domain.DecisionTraceReviewContextScope,
		Ledger: ledger, TaskContexts: contexts, Entries: []domain.DecisionTraceExplanationEntry{},
		Limitations: append(append([]string{}, explanation.Limitations...),
			"Task contexts contain only the task and its recorded changed ancestors, not every dependency or a complete causal path.",
			"Locations identify declarations in supplied historical documents; they are not paths to fetch, authenticated facts or instructions to execute.",
			"Progress remains declared and consultative; document context does not approve completion or replace the historical decision."),
	}
	for _, entry := range explanation.Entries {
		if selected[entry.NodeID] {
			result.Entries = append(result.Entries, entry)
		}
	}
	if len(result.Entries) != len(selected) {
		return fail(fmt.Errorf("incomplete review context declaration coverage"))
	}
	return boundedReviewContextProjection(result)
}

// boundedReviewContextProjection admits aggregate output size, not factual validity.
func boundedReviewContextProjection(result domain.DecisionTraceReviewContext) (domain.DecisionTraceReviewContext, error) {
	if _, err := traceBoundedJSON(result); err != nil {
		return domain.DecisionTraceReviewContext{}, fmt.Errorf("decision trace review context output exceeds budget: %w", err)
	}
	return result, nil
}
