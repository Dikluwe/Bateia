/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_possibilities_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestInsightPossibilitiesSeparateLeadsReadingsAndOperationalCoverage(t *testing.T) {
	report := InsightSiteReport{OpportunityCoverage: &InsightOpportunityCoverage{AttemptedPerspectives: []string{"DIRECT"}, CompletedPerspectives: []string{}, AttemptedSources: []string{ConnectorGitHubCatalog}}, ReferenceLeads: []InsightReferenceLead{{ID: "lead", SourceCandidateID: "parent", Status: "REFERENCE_NOT_READ", Target: ProductIdentity{ConnectorID: ConnectorGitHubCatalog, ExternalID: "other/idea"}}}, Tasks: []InsightSiteTask{{ReferenceLeadRefs: []string{"lead"}, ReadingCardRefs: []string{"card"}, CatalogReferenceIDs: []string{"parent"}}}}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var copy InsightSiteReport
	if err := json.Unmarshal(data, &copy); err != nil || !reflect.DeepEqual(report, copy) {
		t.Fatal("scope lost", err)
	}
	empty, _ := json.Marshal(InsightSiteReport{Tasks: []InsightSiteTask{{}}})
	for _, name := range []string{"reference_leads", "reference_lead_refs", "opportunity_coverage", "scheduling_policy", "execution_order"} {
		if strings.Contains(string(empty), name) {
			t.Fatal("empty extension changed legacy format", name)
		}
	}
}
