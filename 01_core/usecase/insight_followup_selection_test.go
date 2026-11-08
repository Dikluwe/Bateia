/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_followup_selection_test.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"testing"
)

func TestFollowupSelectionKeepsEligibilityBeforeDistanceWithoutEffects(t *testing.T) {
	reader, _ := followFixture()
	entries := []domain.InsightMemoryEntry{
		{ID: "focus", Identity: domain.ProductIdentity{ConnectorID: domain.ConnectorGitHubCatalog}, ReadAttempts: 0, FirstGeneration: 1},
		{ID: "far", Identity: domain.ProductIdentity{ConnectorID: domain.ConnectorCodebergCatalog}, ReadAttempts: 0, FirstGeneration: 1},
		{ID: "near", Identity: domain.ProductIdentity{ConnectorID: domain.ConnectorCodebergCatalog}, ReadAttempts: 5, FirstGeneration: 2},
	}
	report := domain.InsightSiteReport{Plan: domain.InsightSitePlan{Sources: []domain.DiscoveryConnectorCapabilities{{ID: domain.ConnectorGitHubCatalog}, {ID: domain.ConnectorCodebergCatalog}}}, ExplorationMemory: &domain.InsightExplorationMemory{Entries: entries}}
	bindings := map[string]domain.InsightSiteBinding{domain.ConnectorGitHubCatalog: {}, domain.ConnectorCodebergCatalog: {RepositoryReader: reader}}
	queue := append([]domain.InsightMemoryEntry{}, entries...)
	ids := func() []string { return []string{queue[0].ID, queue[1].ID, queue[2].ID} }
	sortInsightFollowupQueue(&report, bindings, queue, nil, map[string]int{"focus": 0, "near": 1, "far": 2})
	if !reflect.DeepEqual(ids(), []string{"near", "far", "focus"}) {
		t.Fatal("unavailable focus outranked an eligible neighbor", ids())
	}
	sortInsightFollowupQueue(&report, bindings, queue, nil, nil)
	if !reflect.DeepEqual(ids(), []string{"far", "near", "focus"}) || !reflect.DeepEqual(report.ExplorationMemory.Entries, entries) || len(reader.resolves) != 0 || len(reader.reads) != 0 {
		t.Fatal("general ordering mutated memory or performed reading", ids())
	}
}

func TestFollowupSelectionTieBreaksAcrossPermutations(t *testing.T) {
	cases := []struct {
		name        string
		generations []int
		references  map[string]int
		want        []string
	}{
		{"age before references and ID", []int{1, 2, 3}, map[string]int{"a": 0, "b": 1, "z": 2}, []string{"z", "b", "a"}},
		{"reference presence before ID", []int{1, 1, 1}, map[string]int{"z": 1, "b": 0}, []string{"b", "z", "a"}},
		{"reference position before ID", []int{1, 1, 1}, map[string]int{"z": 0, "b": 2, "a": 1}, []string{"z", "a", "b"}},
		{"ID breaks final tie", []int{1, 1, 1}, nil, []string{"a", "b", "z"}},
	}
	permutations := [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, focused := range []bool{false, true} {
				entries := make([]domain.InsightMemoryEntry, 3)
				order := map[string]int{}
				var distances map[string]int
				if focused {
					distances = map[string]int{"z": 1, "b": 1, "a": 1}
				}
				for i, id := range []string{"z", "b", "a"} {
					identity := domain.ProductIdentity{ConnectorID: domain.ConnectorGitHubCatalog, ExternalID: "org/" + id, CanonicalURL: "https://github.com/org/" + id}
					entries[i] = domain.InsightMemoryEntry{ID: id, Identity: identity, ReadAttempts: 2, FirstGeneration: tc.generations[i]}
					if position, ok := tc.references[id]; ok {
						order[insightMemoryKey(identity)] = position
					}
				}
				reader, _ := followFixture()
				bindings := map[string]domain.InsightSiteBinding{domain.ConnectorGitHubCatalog: {FollowReader: reader}}
				report := domain.InsightSiteReport{Plan: domain.InsightSitePlan{Sources: []domain.DiscoveryConnectorCapabilities{{ID: domain.ConnectorGitHubCatalog}}}, ExplorationMemory: &domain.InsightExplorationMemory{Entries: entries}}
				before, _ := json.Marshal(report)
				beforeOrder, _ := json.Marshal(order)
				beforeDistances, _ := json.Marshal(distances)
				for _, permutation := range permutations {
					queue := []domain.InsightMemoryEntry{entries[permutation[0]], entries[permutation[1]], entries[permutation[2]]}
					sortInsightFollowupQueue(&report, bindings, queue, order, distances)
					got := []string{queue[0].ID, queue[1].ID, queue[2].ID}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("focused=%t permutation=%v: got %v want %v", focused, permutation, got, tc.want)
					}
				}
				after, _ := json.Marshal(report)
				afterOrder, _ := json.Marshal(order)
				afterDistances, _ := json.Marshal(distances)
				if string(before) != string(after) || string(beforeOrder) != string(afterOrder) || string(beforeDistances) != string(afterDistances) || len(reader.reads) != 0 || len(reader.resolves) != 0 {
					t.Fatal("sorting changed history/metadata or called a reader")
				}
			}
		})
	}
}
