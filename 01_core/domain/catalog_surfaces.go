/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/catalog_surfaces.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "context"

type CatalogSurfaceID string

const (
	SurfaceRepository CatalogSurfaceID = "REPOSITORY"
	SurfaceReadme     CatalogSurfaceID = "README"
	SurfaceTopics     CatalogSurfaceID = "TOPICS"
	SurfacePackage    CatalogSurfaceID = "PACKAGE_REGISTRY"
	SurfaceExternal   CatalogSurfaceID = "EXTERNAL_PROJECT_CATALOG"
)

type CatalogSurfaceCapabilities struct {
	ID                   CatalogSurfaceID `json:"id"`
	Provider             string           `json:"provider"`
	Unit                 string           `json:"unit"`
	SupportsExactPhrase  bool             `json:"supports_exact_phrase"`
	SupportsNegativeTerm bool             `json:"supports_negative_term"`
	ResolvesRepository   bool             `json:"resolves_repository"`
	ResolvesRevision     bool             `json:"resolves_revision"`
}
type SurfaceSearchRequest struct {
	ExpansionID    string           `json:"expansion_id"`
	Surface        CatalogSurfaceID `json:"surface"`
	Query          string           `json:"query"`
	Discriminators []string         `json:"discriminators"`
	NegativeTerms  []string         `json:"negative_terms,omitempty"`
	MaxResults     int              `json:"max_results"`
}
type SurfaceQueryPlan struct {
	ExpansionID             string           `json:"expansion_id"`
	Surface                 CatalogSurfaceID `json:"surface"`
	Original                string           `json:"original"`
	Effective               string           `json:"effective"`
	PreservedDiscriminators []string         `json:"preserved_discriminators"`
	Removed                 []string         `json:"removed,omitempty"`
	Valid                   bool             `json:"valid"`
	Reason                  string           `json:"reason"`
}
type SurfaceMatch struct {
	Surface        CatalogSurfaceID `json:"surface"`
	Provider       string           `json:"provider"`
	Repository     string           `json:"repository,omitempty"`
	Package        string           `json:"package,omitempty"`
	Field          string           `json:"field"`
	Snippet        string           `json:"snippet,omitempty"`
	Discriminators []string         `json:"discriminators"`
	Query          string           `json:"query"`
	ExpansionID    string           `json:"expansion_id"`
	CommitID       string           `json:"commit_id,omitempty"`
	UnknownReason  string           `json:"unknown_reason,omitempty"`
}
type SurfaceDiscovery struct {
	Capabilities CatalogSurfaceCapabilities `json:"capabilities"`
	Plan         SurfaceQueryPlan           `json:"plan"`
	Matches      []SurfaceMatch             `json:"matches"`
	Attempt      CatalogAttempt             `json:"attempt"`
}
type CatalogSurface interface {
	SurfaceCapabilities() CatalogSurfaceCapabilities
	SearchSurface(context.Context, SurfaceSearchRequest) (SurfaceDiscovery, error)
}
type MultiViewCandidate struct {
	Repository RepositoryCandidate  `json:"repository"`
	Matches    []SurfaceMatch       `json:"matches"`
	Surfaces   []CatalogSurfaceID   `json:"surfaces"`
	Concepts   []DiscoveryConceptID `json:"concepts"`
}
type MultiViewSnapshot struct {
	SchemaVersion string               `json:"schema_version"`
	QueryPlans    []SurfaceQueryPlan   `json:"query_plans"`
	Discoveries   []SurfaceDiscovery   `json:"discoveries"`
	Candidates    []MultiViewCandidate `json:"candidates"`
	Unknown       []SurfaceDiscovery   `json:"unknown"`
}

func SameSurface(a, b CatalogSurfaceID) bool { return a != "" && a == b }
