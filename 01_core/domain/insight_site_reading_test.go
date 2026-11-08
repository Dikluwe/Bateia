/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_site_reading_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestInsightSiteContextOriginIsExplicitAndExtensionIsOptional(t *testing.T) {
	report := InsightSiteReport{}
	legacy, err := json.Marshal(report)
	if err != nil || strings.Contains(string(legacy), "site_contexts") {
		t.Fatal("empty extension changed legacy format")
	}
	report.SiteContexts = []InsightSiteContext{{CandidateID: "candidate", Identity: ProductIdentity{ConnectorID: ConnectorGiteeCatalog, ExternalID: "org/repo", CanonicalURL: "https://gitee.com/org/repo"}}}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var copy InsightSiteReport
	if err := json.Unmarshal(encoded, &copy); err != nil || len(copy.SiteContexts) != 1 || copy.SiteContexts[0].Identity != report.SiteContexts[0].Identity {
		t.Fatal("origin lost in report")
	}
	failure := InsightReadmeFailure{Code: "AUTH_REQUIRED"}
	if failure.Error() != "AUTH_REQUIRED" {
		t.Fatal("error contract changed")
	}
}
