/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_sites_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInsightSitesWireKeepsCatalogAndReadingsDistinct(t *testing.T) {
	r := InsightSiteReport{Status: "PARTIAL", NoveltyState: "NOT_ASSESSED", EffectState: "UNVERIFIED", Attempts: []InsightSiteAttempt{{ConnectorID: ConnectorAtomGitCatalog, Status: "AUTH_REQUIRED"}}, Candidates: []InsightSiteCandidate{{ID: "catalog", Origin: "CATALOG_METADATA_UNVERIFIED", Product: ProductCandidate{Identity: ProductIdentity{ConnectorID: ConnectorGiteeCatalog, ExternalID: "org/repo"}}}}, Tasks: []InsightSiteTask{{ReadingCardRefs: []string{"card"}, CatalogReferenceIDs: []string{"catalog"}}}}
	b, e := json.Marshal(r)
	if e != nil {
		t.Fatal(e)
	}
	var got InsightSiteReport
	if e = json.Unmarshal(b, &got); e != nil || !reflect.DeepEqual(r, got) {
		t.Fatal("site/evidence scope lost", e)
	}
}
