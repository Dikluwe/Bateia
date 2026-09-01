/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/regional_discovery.md
 * @layer L1
 * @updated 2026-09-01
 */

package domain

import (
	"fmt"
	"strings"
)

type RegionalQueryCatalog struct {
	SchemaVersion string                   `json:"schema_version"`
	Version       string                   `json:"version"`
	Concepts      []RegionalCatalogConcept `json:"concepts"`
}

type RegionalCatalogConcept struct {
	Concept            string                 `json:"concept"`
	RuleID             string                 `json:"rule_id"`
	SourceQuery        string                 `json:"source_query"`
	DerivationStrategy string                 `json:"derivation_strategy"`
	Queries            []RegionalCatalogQuery `json:"queries"`
}

type RegionalCatalogQuery struct {
	Locale string `json:"locale"`
	Query  string `json:"query"`
}

type RegionalDiscoveryPlan struct {
	Requested []string               `json:"requested"`
	Realized  []string               `json:"realized"`
	Gaps      []RegionalDiscoveryGap `json:"gaps"`
	Proposals []SearchExpansion      `json:"proposals"`
}

type RegionalDiscoveryGap struct {
	Locale  string `json:"locale"`
	Concept string `json:"concept"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
}

type RegionalQueryProvenance struct {
	Locale             string `json:"locale"`
	CatalogVersion     string `json:"catalog_version"`
	Concept            string `json:"concept"`
	RuleID             string `json:"rule_id"`
	SourceQuery        string `json:"source_query"`
	DerivationStrategy string `json:"derivation_strategy"`
	ProposalID         string `json:"proposal_id"`
}

func ValidateRegionalCandidateOrigin(origin CandidateOrigin) error {
	locale := strings.TrimSpace(origin.SearchLocale)
	proposalID := strings.TrimSpace(origin.RegionalProposalID)
	if locale == "" && proposalID == "" {
		return nil
	}
	if locale == "" || proposalID == "" {
		return fmt.Errorf("origem regional exige search_locale e regional_proposal_id")
	}
	if strings.TrimSpace(origin.PlannedQuery) == "" {
		return fmt.Errorf("origem regional exige consulta-fonte planejada")
	}
	return nil
}
