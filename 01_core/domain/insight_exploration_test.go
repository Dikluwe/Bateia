/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_exploration_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInsightExplorationWireSeparatesReadingsFromConclusions(t *testing.T) {
	report := InsightExplorationReport{SchemaVersion: "insight-exploration-v1", Status: "CONTEXT_PREPARED", NoveltyState: "NOT_ASSESSED", EffectState: "UNVERIFIED",
		Cards: []InsightReadingCard{{ID: "card_1", Origin: "UNINTERPRETED_EXCERPT", Reference: LocatedInsightReference{InsightObservationReference: InsightObservationReference{Path: "queue.go", ContentHash: "sha256:fixture", StartLine: 2, EndLine: 2}, Excerpt: "reject()"}, FamiliarSubstringMatches: []string{"reject"}}},
		Tasks: []InsightExplorationTask{{Perspective: "INVERSION", CardRefs: []string{"card_1"}}}, ResearchDirections: []InsightResearchDirection{{Status: "PLANNED_NOT_EXECUTED"}}}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var got InsightExplorationReport
	if err := json.Unmarshal(data, &got); err != nil || !reflect.DeepEqual(report, got) {
		t.Fatalf("reading/task scopes lost: %v", err)
	}
}
