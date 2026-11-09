/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_reference_collection.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func collectInsightReadme(ctx context.Context, reference domain.InsightDiscoveredReference, reader domain.InsightReferenceReader) (domain.InsightReadmeReceipt, *domain.EvidenceBundle) {
	repo := reference.Repository
	ref := domain.CommitRef{Owner: repo.Owner, Repo: repo.Repo, Commit: repo.CommitID}
	receipt := domain.InsightReadmeReceipt{ReferenceID: reference.ID, Source: ref, Status: "NOT_STARTED", Scope: "ROOT_README_ONLY"}
	if reference.RevisionConflict {
		receipt.Status = "REVISION_CONFLICT"
		return receipt, nil
	}
	if !immutableCommitPattern.MatchString(repo.CommitID) {
		if ctx.Err() != nil {
			receipt.Status = "CANCELLED"
			return receipt, nil
		}
		resolver, ok := reader.(domain.InsightSiteRepositoryReader)
		if !ok {
			receipt.Status = "OPAQUE_REVISION"
			return receipt, nil
		}
		id := domain.ProductIdentity{ConnectorID: domain.ConnectorGitHubCatalog, ExternalID: repo.Owner + "/" + repo.Repo, CanonicalURL: repo.URL}
		if !validInsightSiteRepository(id) {
			receipt.Status = "INTEGRITY_ERROR"
			return receipt, nil
		}
		resolved, err := resolver.ResolveRepository(ctx, id)
		if err != nil {
			receipt.Status = insightReadmeFailureStatus(err)
			return receipt, nil
		}
		if resolved.Owner != repo.Owner || resolved.Repo != repo.Repo || resolved.URL != repo.URL || !immutableCommitPattern.MatchString(resolved.CommitID) {
			receipt.Status = "INTEGRITY_ERROR"
			return receipt, nil
		}
		repo.CommitID = resolved.CommitID
		ref.Commit = resolved.CommitID
		receipt.Source = ref
	}
	if ctx.Err() != nil {
		receipt.Status = "CANCELLED"
		return receipt, nil
	}
	tree, err := reader.FetchTree(ctx, repo)
	if err != nil {
		receipt.Status = insightReadmeFailureStatus(err)
		return receipt, nil
	}
	paths := []string{}
	rank := map[string]int{"readme.md": 0, "readme.rst": 1, "readme.txt": 2, "readme": 3}
	for _, node := range tree {
		if _, ok := rank[strings.ToLower(node.Path)]; ok && node.Type == "file" && validMechanismPath(node.Path) {
			paths = append(paths, node.Path)
		}
	}
	sort.Slice(paths, func(i, j int) bool {
		a, b := rank[strings.ToLower(paths[i])], rank[strings.ToLower(paths[j])]
		if a != b {
			return a < b
		}
		return paths[i] < paths[j]
	})
	if len(paths) == 0 {
		receipt.Status = "README_NOT_FOUND"
		return receipt, nil
	}
	receipt.Path = paths[0]
	if ctx.Err() != nil {
		receipt.Status = "CANCELLED"
		return receipt, nil
	}
	extracted, err := reader.ExtractRepoFiles(ctx, repo, []string{paths[0]})
	if err != nil {
		receipt.Status = insightReadmeFailureStatus(err)
		return receipt, nil
	}
	if ctx.Err() != nil {
		receipt.Status = "CANCELLED"
		return receipt, nil
	}
	if extracted == nil || extracted.Ref != ref || len(extracted.Files) != 1 || extracted.Files[0].Path != paths[0] {
		receipt.Status = "INTEGRITY_ERROR"
		return receipt, nil
	}
	content := extracted.Files[0].Content
	if len(content) > 32<<10 {
		receipt.Status = "CONTENT_LIMIT"
		return receipt, nil
	}
	if !validInsightReadmeText(content) {
		receipt.Status = "INTEGRITY_ERROR"
		return receipt, nil
	}

	bundle := &domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: false, TreeComplete: false, BytesRead: int64(len(content)), Items: []domain.CollectedEvidence{{Kind: domain.EvidenceDocumentation, Ref: ref, Path: paths[0], Content: content, ContentHash: "sha256:" + contentHash(content)}}, Targets: []domain.CollectionTargetResult{{Target: domain.ExtractionTarget{Path: paths[0], Kind: domain.EvidenceDocumentation, Reason: "ROOT_README_ONLY"}, Status: domain.TargetCollected, Bytes: int64(len(content))}}}
	bundle.Plan = domain.EvidenceExtractionPlan{SchemaVersion: "1", Repository: ref, Candidate: repo, Targets: []domain.ExtractionTarget{{Path: paths[0], Kind: domain.EvidenceDocumentation, Reason: "ROOT_README_ONLY"}}, Budget: domain.ExtractionBudget{Documents: 1, TotalBytes: 32 << 10}}
	// Coverage concerns this document only, not a software capability or repository audit.
	bundle.RequirementCoverage = []domain.RequirementCollectionCoverage{{RequirementKind: "document", Requirement: "root README", RuleID: "insight-readme-v1", Status: domain.RequirementTargetedExamined, SearchComplete: true, CandidatePaths: []string{paths[0]}, ExaminedPaths: []string{paths[0]}, SkippedPaths: []string{}, ReasonCode: "ROOT_README_ONLY"}}
	receipt.Status = "COLLECTED"
	receipt.Bytes = len(content)
	return receipt, bundle
}

// The same textual profile governs fresh and frozen ROOT_README_ONLY documents.
func validInsightReadmeText(content string) bool {
	if strings.TrimSpace(content) == "" || !utf8.ValidString(content) {
		return false
	}
	for _, r := range content {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return false
		}
	}
	return true
}
