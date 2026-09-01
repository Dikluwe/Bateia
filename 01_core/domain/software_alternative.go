/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/software_alternative.md
 * @layer L1
 * @updated 2026-08-31
 */

package domain

type SoftwareAlternativeFit string

const (
	AlternativeRecommended    SoftwareAlternativeFit = "RECOMMENDED"
	AlternativePartial        SoftwareAlternativeFit = "PARTIAL"
	AlternativeReference      SoftwareAlternativeFit = "REFERENCE"
	AlternativeBlocked        SoftwareAlternativeFit = "BLOCKED"
	AlternativeInconclusive   SoftwareAlternativeFit = "INCONCLUSIVE"
	AlternativeNotAlternative SoftwareAlternativeFit = "NOT_ALTERNATIVE"
)

type SoftwareAlternativeTarget struct {
	Name       string              `json:"name"`
	Repository *RepositoryIdentity `json:"repository,omitempty"`
}

type SoftwareAlternativeFilters struct {
	AllowedLicenses        []string `json:"allowed_licenses"`
	RequireLocal           bool     `json:"require_local"`
	AllowArchived          bool     `json:"allow_archived"`
	RecommendedMinCoverage float64  `json:"recommended_min_coverage"`
	PartialMinCoverage     float64  `json:"partial_min_coverage"`
}

type SoftwareAlternativeCapability struct {
	Capability string           `json:"capability"`
	Status     EvidenceStatus   `json:"status"`
	Evidence   []EvidenceSource `json:"evidence"`
}

type SoftwareAlternative struct {
	Repository       CommitRef                       `json:"repository"`
	URL              string                          `json:"url"`
	Description      string                          `json:"description"`
	Stars            int                             `json:"stars"`
	Fit              SoftwareAlternativeFit          `json:"fit"`
	Coverage         float64                         `json:"coverage"`
	ProblemAlignment EvidenceStatus                  `json:"problem_alignment"`
	Capabilities     []SoftwareAlternativeCapability `json:"capabilities"`
	Gaps             []string                        `json:"gaps"`
	License          LicenseAssessment               `json:"license"`
	Maintenance      MaintenanceAssessment           `json:"maintenance"`
	Operation        OperationAssessment             `json:"operation"`
	Reasons          []string                        `json:"reasons"`
	Blockers         []string                        `json:"blockers"`
	Evidence         []EvidenceSource                `json:"evidence"`
	DiscoveryLocales []string                        `json:"discovery_locales,omitempty"`
}

type SoftwareAlternativeReport struct {
	SchemaVersion string                     `json:"schema_version"`
	Target        SoftwareAlternativeTarget  `json:"target"`
	Filters       SoftwareAlternativeFilters `json:"filters"`
	TopChoice     *CommitRef                 `json:"top_choice,omitempty"`
	Alternatives  []SoftwareAlternative      `json:"alternatives"`
	Excluded      []SoftwareAlternative      `json:"excluded"`
}
