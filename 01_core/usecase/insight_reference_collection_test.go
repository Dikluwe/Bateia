/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_reference_collection_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"strings"
	"testing"
)

func insightReadmeReference() domain.InsightDiscoveredReference {
	return domain.InsightDiscoveredReference{ID: "ref", Repository: insightCatalogCandidate("org", "queue").Repository}
}

func TestInsightReadmeCollectsOnlyRootDocumentAtPinnedRevision(t *testing.T) {
	reader := &insightTestReader{tree: []domain.DirectoryNode{{Path: "docs/README.md", Type: "file"}, {Path: "README.rst", Type: "file"}, {Path: "README.md", Type: "symlink"}, {Path: "readme.md", Type: "file"}}}
	receipt, bundle := collectInsightReadme(context.Background(), insightReadmeReference(), reader)
	if receipt.Status != "COLLECTED" || receipt.Path != "readme.md" || bundle == nil || bundle.Repository.Commit != strings.Repeat("a", 40) || bundle.Complete || bundle.TreeComplete || len(bundle.Items) != 1 {
		t.Fatal("readme scope/revision lost")
	}
	item := bundle.Items[0]
	if item.ContentHash != "sha256:"+contentHash(item.Content) || item.Ref != bundle.Repository || reader.extractCalls != 1 {
		t.Fatal("unbound content or extra extraction")
	}
	if len(bundle.RequirementCoverage) != 1 || bundle.RequirementCoverage[0].RequirementKind != "document" || bundle.RequirementCoverage[0].Status != domain.RequirementTargetedExamined {
		t.Fatal("document coverage lost or confused with capability")
	}
}

func TestInsightReadmeRefusesOpaqueConflictingAndUntrustedContent(t *testing.T) {
	for _, mode := range []string{"opaque", "conflict", "missing", "ref", "path", "cardinality", "utf8", "control", "large", "limit"} {
		t.Run(mode, func(t *testing.T) {
			ref := insightReadmeReference()
			reader := &insightTestReader{}
			switch mode {
			case "opaque":
				ref.Repository.CommitID = "HEAD"
			case "conflict":
				ref.RevisionConflict = true
			case "missing":
				reader.tree = []domain.DirectoryNode{{Path: "docs/README.md", Type: "file"}}
			case "ref":
				reader.mutate = func(out *domain.ExtractedRepository) { out.Ref.Commit = strings.Repeat("b", 40) }
			case "path":
				reader.mutate = func(out *domain.ExtractedRepository) { out.Files[0].Path = "different.md" }
			case "cardinality":
				reader.mutate = func(out *domain.ExtractedRepository) { out.Files = append(out.Files, out.Files[0]) }
			case "utf8":
				reader.content = string([]byte{255})
			case "control":
				reader.content = "text\x00binary"
			case "large":
				reader.content = strings.Repeat("a", (32<<10)+1)
			case "limit":
				reader.content = strings.Repeat("é", 16<<10)
			}
			receipt, bundle := collectInsightReadme(context.Background(), ref, reader)
			if mode == "limit" {
				if receipt.Status != "COLLECTED" || receipt.Bytes != 32<<10 || bundle == nil {
					t.Fatal("exact byte limit rejected")
				}
				return
			}
			if receipt.Status == "COLLECTED" || bundle != nil {
				t.Fatal("invalid collection admitted")
			}
			if (mode == "opaque" || mode == "conflict") && reader.treeCalls != 0 {
				t.Fatal("unresolved revision was read")
			}
		})
	}
}
