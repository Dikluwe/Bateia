/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/project_catalog.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import (
	"context"
	"time"
)

type DiscoveryProviderRole string

const (
	RoleProjectCatalog DiscoveryProviderRole = "PROJECT_CATALOG"
	RoleCodeProof      DiscoveryProviderRole = "CODE_PROOF"
)

type ProviderCapabilities struct {
	Provider       string                  `json:"provider"`
	Roles          []DiscoveryProviderRole `json:"roles"`
	SearchUnits    []string                `json:"search_units"`
	IdentityLevels []string                `json:"identity_levels"`
}

type ProviderCapabilityReporter interface{ Capabilities() ProviderCapabilities }

type CatalogSearchRequest struct {
	ExpansionID    string        `json:"expansion_id"`
	Query          string        `json:"query"`
	Discriminators []string      `json:"discriminators"`
	MaxResults     int           `json:"max_results"`
	AttemptTimeout time.Duration `json:"attempt_timeout,omitempty"`
}

type CatalogCandidate struct {
	Repository    RepositoryCandidate `json:"repository"`
	Topics        []string            `json:"topics,omitempty"`
	ReadmeSnippet string              `json:"readme_snippet,omitempty"`
	License       string              `json:"license,omitempty"`
	MatchedFields []string            `json:"matched_fields"`
	Origins       []CandidateOrigin   `json:"origins"`
}

type CatalogAttempt struct {
	ExpansionID string        `json:"expansion_id"`
	Query       string        `json:"query"`
	Status      string        `json:"status"`
	Duration    time.Duration `json:"duration,omitempty"`
	ResultCount int           `json:"result_count"`
	ErrorClass  string        `json:"error_class,omitempty"`
	SafeMessage string        `json:"safe_message,omitempty"`
}

type CatalogDiscovery struct {
	Provider   string                `json:"provider"`
	Role       DiscoveryProviderRole `json:"role"`
	SearchUnit string                `json:"search_unit"`
	Candidates []CatalogCandidate    `json:"candidates"`
	Attempt    CatalogAttempt        `json:"attempt"`
}

type ProjectCatalog interface {
	ProviderCapabilityReporter
	SearchCatalog(context.Context, CatalogSearchRequest) (CatalogDiscovery, error)
}

type CatalogSnapshot struct {
	SchemaVersion string               `json:"schema_version"`
	Idea          IdeaSpec             `json:"idea"`
	Plan          ExpansionPlan        `json:"plan"`
	Provider      ProviderCapabilities `json:"provider"`
	PerExpansion  []CatalogDiscovery   `json:"per_expansion"`
	Candidates    []CatalogCandidate   `json:"candidates"`
	Receipt       *DiscoveryRunReceipt `json:"receipt,omitempty"`
}

type CatalogVerdict string

const (
	CatalogAdvance CatalogVerdict = "ADVANCE"
	CatalogReject  CatalogVerdict = "REJECT"
	CatalogReview  CatalogVerdict = "REVIEW"
)

type CatalogSignal struct {
	ID           string  `json:"id"`
	Contribution float64 `json:"contribution"`
	Reason       string  `json:"reason"`
}

type CatalogDecision struct {
	Repository RepositoryCandidate `json:"repository"`
	Verdict    CatalogVerdict      `json:"verdict"`
	Score      float64             `json:"score"`
	Signals    []CatalogSignal     `json:"signals"`
	Origins    []CandidateOrigin   `json:"origins"`
}

type CatalogScreeningReport struct {
	SchemaVersion string            `json:"schema_version"`
	Decisions     []CatalogDecision `json:"decisions"`
}

func SupportsRole(capabilities ProviderCapabilities, role DiscoveryProviderRole) bool {
	for _, candidate := range capabilities.Roles {
		if candidate == role {
			return true
		}
	}
	return false
}
