/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_rules.md
 * @layer L1
 * @updated 2026-09-16
 */

package usecase

import (
	"bateia/01_core/domain"
)

func NewRecommendationServiceWithEvidenceRules(
	semantic SemanticExpander,
	discovery domain.ObservableDiscoveryProvider,
	content domain.ContentProvider,
	catalog domain.CapabilityRuleCatalog,
) (*RecommendationService, error) {
	if err := validateRuleCatalog(catalog, true); err != nil {
		return nil, err
	}
	catalog.Rules = canonicalCapabilityRules(catalog.Rules)
	catalog.Hash = ruleCatalogHash(catalog)
	return &RecommendationService{dependencies: IdeaRunDependencies{
		Semantic: semantic, Discovery: discovery, Content: content, EvidenceRules: catalog,
	}}, nil
}
