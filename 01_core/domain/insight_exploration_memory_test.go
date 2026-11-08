/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_exploration_memory_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInsightMemoryWireKeepsHistoryApartFromCurrentContent(t *testing.T) {
	report := InsightSiteReport{ExplorationMemory: &InsightExplorationMemory{SchemaVersion: "insight-exploration-memory-v1", Entries: []InsightMemoryEntry{{ID: "entry", State: "DISCARDED", Reason: "Already investigated"}}}, Followups: []InsightFollowup{{EntryID: "entry2", Status: "COLLECTED", Cards: []InsightReadingCard{{ID: "read-card"}}}}}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var copy InsightSiteReport
	if json.Unmarshal(data, &copy) != nil || !reflect.DeepEqual(report, copy) {
		t.Fatal("history and reading conflated")
	}
}
