/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_frontier.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import "bateia/01_core/domain"

func insightNeighborhoodDistances(report *domain.InsightSiteReport, focus string) map[string]int {
	if focus == "" {
		return nil
	}
	children := map[string][]domain.ProductIdentity{}
	add := func(lead domain.InsightReferenceLead) {
		source := insightNeighborhoodAliasTarget(report.ExplorationMemory, lead.SourceIdentity)
		target := insightNeighborhoodAliasTarget(report.ExplorationMemory, lead.Target)
		key := insightMemoryKey(source)
		children[key] = append(children[key], target)
	}
	var root domain.ProductIdentity
	for _, entry := range report.ExplorationMemory.Entries {
		if entry.Lead != nil {
			add(*entry.Lead)
		}
		if entry.Identity.CanonicalURL == focus {
			root = entry.Identity
		}
	}
	report.Neighborhood.Focus = &root
	for _, lead := range report.ReferenceLeads {
		add(lead)
	}
	reachable := map[string]int{insightMemoryKey(root): 0}
	queue := []domain.ProductIdentity{root}
	for len(queue) > 0 {
		parent := queue[0]
		queue = queue[1:]
		entry := insightMemoryFind(report.ExplorationMemory, parent)
		if entry == nil || entry.State == "DISCARDED" {
			continue
		}
		status, _ := insightNeighborhoodAliasDisposition(report.ExplorationMemory, *entry)
		if status == "ALIAS_STATE_CONFLICT" {
			continue
		}
		for _, child := range children[insightMemoryKey(parent)] {
			key := insightMemoryKey(child)
			if _, seen := reachable[key]; !seen && insightMemoryFind(report.ExplorationMemory, child) != nil {
				reachable[key] = reachable[insightMemoryKey(parent)] + 1
				queue = append(queue, child)
			}
		}
	}
	return reachable
}

func insightNeighborhoodFrontier(report *domain.InsightSiteReport, focus string, bindings map[string]domain.InsightSiteBinding, order map[string]int) ([]domain.InsightMemoryEntry, map[string]int) {
	reachable := insightNeighborhoodDistances(report, focus)
	all := []domain.InsightMemoryEntry{}
	var distances map[string]int
	if focus != "" {
		distances = map[string]int{}
	}
	for _, entry := range report.ExplorationMemory.Entries {
		if entry.State != "PENDING" {
			continue
		}
		target := insightNeighborhoodAliasTarget(report.ExplorationMemory, entry.Identity)
		distance, known := reachable[insightMemoryKey(target)]
		if focus != "" && !known {
			continue
		}
		all = append(all, entry)
		if distances != nil {
			distances[entry.ID] = distance
		}
	}
	return selectInsightNeighborhoodNeighbors(report, bindings, all, order, distances), distances
}

func prepareInsightNeighborhoodLeads(report *domain.InsightSiteReport, origins []insightLeadOrigin) map[string]int {
	activeOrigins := []insightLeadOrigin{}
	for _, origin := range origins {
		if entry := insightMemoryFind(report.ExplorationMemory, origin.candidate.Product.Identity); entry != nil && entry.State != "DISCARDED" {
			activeOrigins = append(activeOrigins, origin)
		}
	}
	known := map[string]bool{}
	for _, entry := range report.ExplorationMemory.Entries {
		if entry.State == "EXAMINED" || entry.State == "DISCARDED" {
			known[insightMemoryKey(entry.Identity)] = true
		}
	}
	for _, origin := range origins {
		known[insightMemoryKey(origin.candidate.Product.Identity)] = true
	}
	deriveInsightReferenceLeads(report, activeOrigins, known)
	return rememberInsightReferenceLeads(report)
}
