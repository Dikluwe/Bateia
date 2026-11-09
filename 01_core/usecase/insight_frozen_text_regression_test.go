/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_frozen_text_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"testing"
)

func TestFrozenSourcesShareCollectedTextProfile(t *testing.T) {
	for _, kind := range []string{"github", "site", "followup"} {
		for _, content := range []string{"documentation", "\tline\r\nnext", "", "  \n\t", "line\x00hidden", "line\x7fhidden", "line\xffhidden"} {
			t.Run(kind+"/"+content, func(t *testing.T) {
				input := frozenNeighborhoodOriginsFixture(t, kind, []int{1})
				var b *domain.EvidenceBundle
				var receipt *domain.InsightReadmeReceipt
				switch kind {
				case "github":
					b = &input.Bundles[0]
					receipt = input.Readmes[0].Receipt
				case "site":
					b = &input.SiteContexts[0].Bundle
					receipt = input.Readmes[0].Receipt
				case "followup":
					b = input.Followups[0].Bundle
					receipt = input.Followups[0].Receipt
				}
				b.Items[0].Content = content
				b.Items[0].ContentHash = "sha256:" + contentHash(content)
				b.BytesRead = int64(len(content))
				receipt.Bytes = len(content)
				input.ExplorationMemory.Entries[0].LastRead.ContentHash = b.Items[0].ContentHash
				sealed := sealInsightMemory(*input.ExplorationMemory)
				input.ExplorationMemory = &sealed
				input.ID = ""
				input.ID = "site_discovery_" + insightHash(input)
				valid := content == "documentation" || content == "\tline\r\nnext"
				_, err := ValidateInsightNeighborhood(domain.InsightNeighborhoodRequest{Input: input})
				if valid {
					if err != nil {
						t.Fatal("valid frozen text rejected", err)
					}
				} else {
					if !errors.Is(err, domain.ErrIntegrity) {
						t.Fatalf("invalid frozen text admitted: %q err=%v", content, err)
					}
					assertFrozenOriginRejectedBeforeRead(t, input)
				}
				// Public collection must agree with admission on the textual profile.
				reader := &insightTestReader{mutate: func(extracted *domain.ExtractedRepository) { extracted.Files[0].Content = content }}
				zero := 0
				fresh, err := DiscoverInsightSites(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "queue"}, Sources: []string{"github"}, FollowLinks: &zero}, []domain.InsightSiteBinding{{Connector: &insightSiteFixtureConnector{id: domain.ConnectorGitHubCatalog}, Reader: reader}})
				if err != nil {
					t.Fatal(err)
				}
				if (len(fresh.Bundles) == 1) != valid {
					t.Fatal("collector/admission profile disagree")
				}
			})
		}
	}
}
