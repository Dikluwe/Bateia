/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_selection.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"strings"
)

// Only a pair already present on the exact same host is recognized as an alias.
func insightNeighborhoodAliasTarget(memory *domain.InsightExplorationMemory, identity domain.ProductIdentity) domain.ProductIdentity {
	if !strings.HasSuffix(identity.ExternalID, ".git") {
		return identity
	}
	base, ok := insightLeadTarget(strings.TrimSuffix(identity.CanonicalURL, ".git"))
	if !ok || strings.HasSuffix(base.ExternalID, ".git") || base.ConnectorID != identity.ConnectorID || insightMemoryFind(memory, identity) == nil {
		return identity
	}
	entry := insightMemoryFind(memory, base)
	if entry != nil && entry.Identity == base {
		return base
	}
	return identity
}

func insightNeighborhoodAliasDisposition(memory *domain.InsightExplorationMemory, entry domain.InsightMemoryEntry) (string, string) {
	target := insightNeighborhoodAliasTarget(memory, entry.Identity)
	if target != entry.Identity {
		base := insightMemoryFind(memory, target)
		switch base.State {
		case "EXAMINED":
			return "ALIAS_ALREADY_EXAMINED", base.ID
		case "DISCARDED":
			return "ALIAS_DISCARDED", base.ID
		default:
			return "DUPLICATE_ALIAS", base.ID
		}
	}
	if !strings.HasSuffix(entry.Identity.ExternalID, ".git") {
		identity, ok := insightLeadTarget(entry.Identity.CanonicalURL + ".git")
		alias := insightMemoryFind(memory, identity)
		if ok && alias != nil && insightNeighborhoodAliasTarget(memory, alias.Identity) == entry.Identity && entry.State == "PENDING" {
			if alias.State == "DISCARDED" {
				return "ALIAS_STATE_CONFLICT", alias.ID
			}
			if alias.State == "EXAMINED" && entry.Reason == "" {
				return "ALIAS_ALREADY_EXAMINED", alias.ID
			}
		}
	}
	return "FRONTIER", ""
}

func selectInsightNeighborhoodNeighbors(report *domain.InsightSiteReport, bindings map[string]domain.InsightSiteBinding, all []domain.InsightMemoryEntry, order, distances map[string]int) []domain.InsightMemoryEntry {
	sortInsightFollowupQueue(report, bindings, all, order, distances)
	queue := []domain.InsightMemoryEntry{}
	selected, suppressed := []domain.InsightNeighborhoodNeighbor{}, []domain.InsightNeighborhoodNeighbor{}
	for _, entry := range all {
		status, related := insightNeighborhoodAliasDisposition(report.ExplorationMemory, entry)
		neighbor := domain.InsightNeighborhoodNeighbor{EntryID: entry.ID, Identity: entry.Identity, ReadAttempts: entry.ReadAttempts, SelectionStatus: status, RelatedEntryID: related}
		if distances != nil {
			distance := distances[entry.ID]
			neighbor.Distance = &distance
		}
		if status != "FRONTIER" {
			suppressed = append(suppressed, neighbor)
			if status == "ALIAS_STATE_CONFLICT" {
				report.Gaps = append(report.Gaps, "ALIAS_STATE_CONFLICT:"+entry.ID)
			}
			continue
		}
		queue = append(queue, entry)
		selected = append(selected, neighbor)
		report.Neighborhood.FrontierEntryIDs = append(report.Neighborhood.FrontierEntryIDs, entry.ID)
	}
	report.Neighborhood.SelectionPolicy = "eligible-distance-attempts-age-reference-v1"
	report.Neighborhood.Neighbors = append(selected, suppressed...)
	return queue
}
