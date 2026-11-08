/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_discovery_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInsightDiscoveryWireKeepsOperationalAndEvidenceScopes(t *testing.T) {
	report := InsightDiscoveryReport{Status: "PARTIAL", NoveltyState: "NOT_ASSESSED", EffectState: "UNVERIFIED", Attempts: []InsightDiscoveryAttempt{{QueryID: "q", Status: "EMPTY"}}, Candidates: []InsightDiscoveredReference{{ID: "ref", QueryIDs: []string{"q"}, RevisionConflict: true, ObservedCommits: []string{"a", "b"}}}, Readmes: []InsightReadmeReceipt{{ReferenceID: "ref", Status: "REVISION_CONFLICT", Scope: "ROOT_README_ONLY"}}, Exploration: InsightExplorationReport{Status: "SEARCH_NEEDED"}}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var got InsightDiscoveryReport
	if err := json.Unmarshal(data, &got); err != nil || !reflect.DeepEqual(got, report) {
		t.Fatalf("scopes lost: %v", err)
	}
}
