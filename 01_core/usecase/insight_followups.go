/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_followups.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
)

func runInsightFollowups(ctx context.Context, report *domain.InsightSiteReport, bindings map[string]domain.InsightSiteBinding) (bool, error) {
	order := rememberInsightReferenceLeads(report)
	pending := []domain.InsightMemoryEntry{}
	for _, entry := range report.ExplorationMemory.Entries {
		if entry.State == "PENDING" && entry.Lead != nil {
			pending = append(pending, entry)
		}
	}
	return runInsightFollowupQueue(ctx, report, bindings, pending, order)
}

func runInsightFollowupQueue(ctx context.Context, report *domain.InsightSiteReport, bindings map[string]domain.InsightSiteBinding, pending []domain.InsightMemoryEntry, order map[string]int, priorities ...map[string]int) (bool, error) {
	var distances map[string]int
	if len(priorities) > 0 {
		distances = priorities[0]
	}
	sortInsightFollowupQueue(report, bindings, pending, order, distances)
	selected := map[string]bool{}
	for _, source := range report.Plan.Sources {
		selected[source.ID] = true
	}
	bytes, documents := 0, 0
	for _, b := range report.Bundles {
		bytes += int(b.BytesRead)
		documents++
	}
	for _, s := range report.SiteContexts {
		bytes += int(s.Bundle.BytesRead)
		documents++
	}
	failed, attempts := false, 0
	for _, entry := range pending {
		if len(report.Followups) >= 16 {
			report.Gaps = append(report.Gaps, "FOLLOWUP_PLAN_LIMITED")
			break
		}
		follow := domain.InsightFollowup{EntryID: entry.ID, Identity: entry.Identity, Status: "FOLLOWUP_BUDGET", Cards: []domain.InsightReadingCard{}}
		if entry.Lead != nil {
			follow.LeadID = entry.Lead.ID
		}
		reader := bindings[entry.Identity.ConnectorID].RepositoryReader
		if entry.Identity.ConnectorID == domain.ConnectorGitHubCatalog {
			reader = bindings[entry.Identity.ConnectorID].FollowReader
		}
		switch {
		case ctx.Err() != nil:
			follow.Status = "NOT_STARTED"
		case !selected[entry.Identity.ConnectorID]:
			follow.Status = "SOURCE_NOT_SELECTED"
		case reader == nil:
			follow.Status = "READER_UNAVAILABLE"
		case attempts >= *report.Plan.Request.FollowLinks:
		case documents >= 8:
			follow.Status = "DOCUMENT_BUDGET"
		case bytes >= 128<<10:
			follow.Status = "CONTENT_BUDGET"
		default:
			operation, stop, err := domain.ExecutionTimeout(ctx, report.Plan.Request.Discovery.AttemptTimeout)
			if err != nil {
				return failed, err
			}
			attempts++
			candidate := domain.InsightSiteCandidate{ID: entry.ID, Product: domain.ProductCandidate{Identity: entry.Identity}}
			receipt, bundle := collectInsightSiteReadme(operation, candidate, reader)
			operationErr := operation.Err()
			stop()
			if operationErr != nil {
				bundle = nil
				receipt.Status = "CANCELLED"
				if errors.Is(operationErr, context.DeadlineExceeded) {
					receipt.Status = "TIMED_OUT"
				}
			}
			if bundle != nil && bytes+receipt.Bytes > 128<<10 {
				receipt.Status = "CONTENT_BUDGET"
				receipt.Bytes = 0
				bundle = nil
			}
			follow.Status, follow.Receipt = receipt.Status, &receipt
			follow.Bundle = bundle
			recordInsightMemoryObservation(report, entry.Identity, receipt.Status, bundle)
			if bundle != nil {
				bytes += receipt.Bytes
				documents++
			} else {
				failed = true
			}
		}
		report.Followups = append(report.Followups, follow)
	}
	if err := prepareInsightFollowupContext(report); err != nil {
		return failed, err
	}
	memory := sealInsightMemory(*report.ExplorationMemory)
	report.ExplorationMemory = &memory
	return failed, nil
}
