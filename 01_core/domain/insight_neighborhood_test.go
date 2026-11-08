/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_neighborhood_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestNeighborhoodExtensionPreservesLegacyJSONAndStrictRoundtrip(t *testing.T) {
	report := InsightSiteReport{SchemaVersion: "insight-site-discovery-v1", ID: "previous-id", NoveltyState: "NOT_ASSESSED", EffectState: "UNVERIFIED"}
	legacy, err := json.Marshal(report)
	if err != nil || bytes.Contains(legacy, []byte("neighborhood")) {
		t.Fatal("optional field changed legacy shape", err)
	}
	report.Neighborhood = &InsightNeighborhood{InputReportID: "input-id", Focus: &ProductIdentity{ConnectorID: ConnectorGitLabCatalog, ExternalID: "team/sub/repo", CanonicalURL: "https://gitlab.com/team/sub/repo"}, FrontierEntryIDs: []string{"pending-1", "pending-2"}}
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var roundtrip InsightSiteReport
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&roundtrip); err != nil || !reflect.DeepEqual(roundtrip.Neighborhood, report.Neighborhood) {
		t.Fatal("expanded report lost public metadata", err)
	}
	roundtrip.Neighborhood = nil
	previous, err := json.Marshal(roundtrip)
	if err != nil || !bytes.Equal(legacy, previous) {
		t.Fatal("removing optional extension changed prior serialization", err)
	}
}

func TestNeighborhoodNeighborMetadataPreservesAbsentFieldsAndZeroDistance(t *testing.T) {
	old := InsightNeighborhood{InputReportID: "previous", FrontierEntryIDs: []string{"entry"}}
	legacy, _ := json.Marshal(old)
	if bytes.Contains(legacy, []byte("neighbors")) || bytes.Contains(legacy, []byte("selection_policy")) {
		t.Fatal("new fields changed previous JSON")
	}
	distance := 0
	old.SelectionPolicy = "eligible-distance-attempts-age-reference-v1"
	old.Neighbors = []InsightNeighborhoodNeighbor{{EntryID: "entry", Identity: ProductIdentity{ConnectorID: ConnectorGitHubCatalog, ExternalID: "org/repo", CanonicalURL: "https://github.com/org/repo"}, Distance: &distance, ReadAttempts: 3, SelectionStatus: "FRONTIER"}, {EntryID: "alias", SelectionStatus: "DUPLICATE_ALIAS", RelatedEntryID: "entry"}}
	raw, _ := json.Marshal(old)
	var actual InsightNeighborhood
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&actual); err != nil || !reflect.DeepEqual(old, actual) || !bytes.Contains(raw, []byte(`"distance":0`)) {
		t.Fatal("neighbor metadata lost zero distance or alias relation", err)
	}
	actual.SelectionPolicy = ""
	actual.Neighbors = nil
	previous, _ := json.Marshal(actual)
	if !bytes.Equal(previous, legacy) {
		t.Fatal("absent extension does not roundtrip historical shape")
	}
}
