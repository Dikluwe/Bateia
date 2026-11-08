/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_followup_selection.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"sort"
)

func sortInsightFollowupQueue(report *domain.InsightSiteReport, bindings map[string]domain.InsightSiteBinding, pending []domain.InsightMemoryEntry, order, distances map[string]int) {
	selected := map[string]bool{}
	for _, source := range report.Plan.Sources {
		selected[source.ID] = true
	}
	eligible := func(entry domain.InsightMemoryEntry) bool {
		b := bindings[entry.Identity.ConnectorID]
		return selected[entry.Identity.ConnectorID] && (b.RepositoryReader != nil || b.FollowReader != nil)
	}
	sort.Slice(pending, func(i, j int) bool {
		a, b := pending[i], pending[j]
		if eligible(a) != eligible(b) {
			return eligible(a)
		}
		if distances != nil && distances[a.ID] != distances[b.ID] {
			return distances[a.ID] < distances[b.ID]
		}
		if a.ReadAttempts != b.ReadAttempts {
			return a.ReadAttempts < b.ReadAttempts
		}
		if a.FirstGeneration != b.FirstGeneration {
			return a.FirstGeneration < b.FirstGeneration
		}
		ai, aok := order[insightMemoryKey(a.Identity)]
		bi, bok := order[insightMemoryKey(b.Identity)]
		if aok != bok {
			return aok
		}
		if aok && ai != bi {
			return ai < bi
		}
		return a.ID < b.ID
	})
}
