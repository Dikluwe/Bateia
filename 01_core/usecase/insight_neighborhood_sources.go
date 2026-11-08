/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_sources.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
)

func insightNeighborhoodOrigins(input domain.InsightSiteReport) ([]insightLeadOrigin, error) {
	origins := []insightLeadOrigin{}
	seen, total := map[string]bool{}, 0
	add := func(identity domain.ProductIdentity, id string, bundle domain.EvidenceBundle, receipt *domain.InsightReadmeReceipt) error {
		entry := insightMemoryFind(input.ExplorationMemory, identity)
		if entry == nil || entry.LastRead == nil || !validInsightMemoryIdentity(identity) || seen[entry.ID] || len(bundle.Items) != 1 || validateEvidenceBundle(bundle) != nil || receipt == nil {
			return domain.ErrIntegrity
		}
		item := bundle.Items[0]
		observation := domain.InsightReadObservation{Source: bundle.Repository, Path: item.Path, ContentHash: item.ContentHash}
		if !validInsightMemoryObservation(identity, observation) || *entry.LastRead != observation || item.Kind != domain.EvidenceDocumentation || len(item.Content) > 32<<10 || int(bundle.BytesRead) != len(item.Content) || receipt.Status != "COLLECTED" || receipt.Source != bundle.Repository || receipt.Path != item.Path || receipt.Bytes != len(item.Content) || receipt.ReferenceID != id || receipt.Scope != "ROOT_README_ONLY" {
			return domain.ErrIntegrity
		}
		total += len(item.Content)
		if len(origins) >= 8 || total > 128<<10 {
			return domain.ErrInvalidRunInput
		}
		seen[entry.ID] = true
		origins = append(origins, insightLeadOrigin{candidate: domain.InsightSiteCandidate{ID: id, Product: domain.ProductCandidate{Identity: identity}}, bundle: bundle})
		return nil
	}
	receiptFor := func(id, connector string) *domain.InsightReadmeReceipt {
		for _, r := range input.Readmes {
			if r.CandidateID == id && r.ConnectorID == connector && r.Status == "COLLECTED" {
				return r.Receipt
			}
		}
		return nil
	}
	for _, bundle := range input.Bundles {
		var candidate *domain.InsightSiteCandidate
		for i := range input.Candidates {
			c := &input.Candidates[i]
			if c.Selected && c.Product.Identity.ConnectorID == domain.ConnectorGitHubCatalog && c.Product.Identity.ExternalID == bundle.Repository.Owner+"/"+bundle.Repository.Repo {
				candidate = c
				break
			}
		}
		if candidate == nil || add(candidate.Product.Identity, candidate.ID, bundle, receiptFor(candidate.ID, domain.ConnectorGitHubCatalog)) != nil {
			return nil, domain.ErrIntegrity
		}
	}
	for _, scoped := range input.SiteContexts {
		if err := add(scoped.Identity, scoped.CandidateID, scoped.Bundle, receiptFor(scoped.CandidateID, scoped.Identity.ConnectorID)); err != nil {
			return nil, err
		}
	}
	for _, follow := range input.Followups {
		if follow.Bundle == nil {
			if follow.Status == "COLLECTED" {
				return nil, domain.ErrIntegrity
			}
			continue
		}
		if follow.Status != "COLLECTED" || follow.EntryID != insightMemoryEntryID(follow.Identity) {
			return nil, domain.ErrIntegrity
		}
		if err := add(follow.Identity, follow.EntryID, *follow.Bundle, follow.Receipt); err != nil {
			return nil, err
		}
	}
	return origins, nil
}

// validateInsightNeighborhoodReceipts closes the reverse edge of the historical
// origin validator: every primary collected receipt must own exactly one bundle.
func validateInsightNeighborhoodReceipts(input domain.InsightSiteReport) error {
	type receiptKey struct{ connector, candidate string }
	seen := map[receiptKey]bool{}
	for _, r := range input.Readmes {
		innerCollected := r.Receipt != nil && r.Receipt.Status == "COLLECTED"
		if r.Status != "COLLECTED" && !innerCollected {
			continue
		}
		key := receiptKey{r.ConnectorID, r.CandidateID}
		if r.Status != "COLLECTED" || !innerCollected || seen[key] {
			return domain.ErrIntegrity
		}
		seen[key] = true
		matches := 0
		if r.ConnectorID == domain.ConnectorGitHubCatalog {
			for _, bundle := range input.Bundles {
				for _, candidate := range input.Candidates {
					if candidate.Selected && candidate.ID == r.CandidateID && candidate.Product.Identity.ConnectorID == r.ConnectorID && candidate.Product.Identity.ExternalID == bundle.Repository.Owner+"/"+bundle.Repository.Repo {
						matches++
						break
					}
				}
			}
		}
		for _, scoped := range input.SiteContexts {
			if scoped.CandidateID == r.CandidateID && scoped.Identity.ConnectorID == r.ConnectorID {
				matches++
			}
		}
		if matches != 1 {
			return domain.ErrIntegrity
		}
	}
	return nil
}
