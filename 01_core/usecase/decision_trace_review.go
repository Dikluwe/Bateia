/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_review.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"fmt"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

// TrackDecisionTraceReview projects declared receipts onto a recomputed plan.
// No receipt authenticates work or promotes the historical decision.
func TrackDecisionTraceReview(input domain.DecisionTraceReviewInput) (domain.DecisionTraceReviewLedger, error) {
	fail := func(err error) (domain.DecisionTraceReviewLedger, error) {
		return domain.DecisionTraceReviewLedger{}, err
	}
	if input.SchemaVersion != "1" || len(input.Records) > traceMaxNodes {
		return fail(fmt.Errorf("invalid review schema or record budget"))
	}
	if _, err := traceBoundedJSON(input); err != nil {
		return fail(err)
	}
	plan, err := PlanDecisionTraceImpact(input.ImpactInput)
	if err != nil {
		return fail(fmt.Errorf("invalid review impact input: %w", err))
	}
	changes := append([]domain.DecisionTraceChange{}, input.ImpactInput.Changes...)
	sort.Slice(changes, func(i, j int) bool { return changes[i].NodeID < changes[j].NodeID })
	planHash, err := traceFingerprint(struct {
		GraphHash string
		Changes   []domain.DecisionTraceChange
		Plan      domain.DecisionTraceImpactPlan
	}{plan.GraphHash, changes, plan})
	if err != nil {
		return fail(err)
	}
	result := domain.DecisionTraceReviewLedger{
		SchemaVersion: domain.DecisionTraceReviewSchemaVersion, PlanHash: planHash,
		GraphHash: plan.GraphHash, Scope: plan.Scope,
		DecisionBefore: plan.DecisionBefore, DecisionAfter: plan.DecisionAfter, DecisionImpact: "NONE",
		Tasks: []domain.DecisionTraceReviewProgress{},
		Limitations: append(append([]string{}, traceLimitations...),
			"Progress, notes, source references and their hashes are supplied declarations; no work, bytes, author or approval is authenticated.",
			"Recorded review completion does not resolve uncertainty, approve evidence or replace the historical decision.",
			"This is a progress snapshot, not an append-only event history; transitions, timestamps and concurrent revisions are not enforced."),
	}
	taskHashes := make(map[string]string, len(plan.Tasks))
	for _, task := range plan.Tasks {
		hash, err := traceFingerprint(struct {
			PlanHash string
			Task     domain.DecisionTraceReviewTask
		}{planHash, task})
		if err != nil {
			return fail(err)
		}
		taskHashes[task.NodeID] = hash
	}
	records := make(map[string]domain.DecisionTraceReviewRecord, len(input.Records))
	for _, record := range input.Records {
		hash, exists := taskHashes[record.NodeID]
		if _, duplicate := records[record.NodeID]; duplicate || !exists {
			return fail(fmt.Errorf("unknown or duplicate review task"))
		}
		if record.PlanHash != planHash || record.TaskHash != hash {
			return fail(fmt.Errorf("review record does not match current plan/task hash"))
		}
		switch record.Status {
		case "PENDING", "IN_PROGRESS":
		case "BLOCKED":
			if strings.TrimSpace(record.Note) == "" {
				return fail(fmt.Errorf("blocked review requires a note"))
			}
		case "COMPLETED":
			if strings.TrimSpace(record.Note) == "" || len(record.Evidence) == 0 {
				return fail(fmt.Errorf("completed review requires a note and evidence reference"))
			}
		default:
			return fail(fmt.Errorf("invalid declared review status"))
		}
		if len(record.Evidence) > 20 {
			return fail(fmt.Errorf("review evidence reference budget exceeded"))
		}
		for _, evidence := range record.Evidence {
			if strings.TrimSpace(evidence.Ref) == "" || !traceValidDigest(evidence.ContentHash) {
				return fail(fmt.Errorf("invalid declared review evidence reference/hash"))
			}
		}
		records[record.NodeID] = record
	}
	for _, task := range plan.Tasks {
		progress := domain.DecisionTraceReviewProgress{Task: task, TaskHash: taskHashes[task.NodeID], Status: "PENDING", Evidence: []domain.DecisionTraceReviewEvidence{}}
		if record, exists := records[task.NodeID]; exists {
			progress.Status, progress.Note = record.Status, record.Note
			progress.Evidence = append(progress.Evidence, record.Evidence...)
		}
		switch progress.Status {
		case "PENDING":
			result.Counts.Pending++
		case "IN_PROGRESS":
			result.Counts.InProgress++
		case "BLOCKED":
			result.Counts.Blocked++
		case "COMPLETED":
			result.Counts.Completed++
		}
		result.Tasks = append(result.Tasks, progress)
	}
	switch {
	case len(result.Tasks) == 0:
		result.Status = "NO_REVIEW_TASKS"
	case result.Counts.Blocked > 0:
		result.Status = "REVIEW_BLOCKED"
	case result.Counts.Completed == len(result.Tasks):
		result.Status = "RECORDED_REVIEW_COMPLETE"
	case result.Counts.InProgress > 0 || result.Counts.Completed > 0:
		result.Status = "REVIEW_IN_PROGRESS"
	default:
		result.Status = "REVIEW_PENDING"
	}
	if _, err := traceBoundedJSON(result); err != nil {
		return fail(fmt.Errorf("decision trace review output exceeds budget: %w", err))
	}
	return result, nil
}
