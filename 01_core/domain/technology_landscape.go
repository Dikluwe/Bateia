/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/technology_landscape.md
 * @layer L1
 * @updated 2026-09-24
 */
package domain

const TechnologyLandscapeSchemaVersion = "1"

type LandscapeCandidateKind string

const (
	LandscapeRepository LandscapeCandidateKind = "REPOSITORY"
	LandscapeProduct    LandscapeCandidateKind = "PRODUCT"
	LandscapeTechnology LandscapeCandidateKind = "TECHNOLOGY"
)

type LandscapeRole string

const (
	LandscapeSubstitute LandscapeRole = "SUBSTITUTE"
	LandscapeComplement LandscapeRole = "COMPLEMENT"
	LandscapeAnalogue   LandscapeRole = "ANALOGUE"
	LandscapeUnresolved LandscapeRole = "UNRESOLVED"
)

type LandscapeAssertion struct {
	Status  EvidenceStatus `json:"status"`
	Sources []string       `json:"sources,omitempty"`
}

type LandscapeCapability struct {
	ID      string         `json:"id"`
	Status  EvidenceStatus `json:"status"`
	Sources []string       `json:"sources,omitempty"`
}

type LandscapeSignalKind string

const (
	LandscapeReleases      LandscapeSignalKind = "RELEASES"
	LandscapeResponseHours LandscapeSignalKind = "MEDIAN_RESPONSE_HOURS"
	LandscapeBusFactor     LandscapeSignalKind = "BUS_FACTOR"
	LandscapeArchived      LandscapeSignalKind = "ARCHIVED"
)

// Signals are observations, not a composite maturity score.
type LandscapeSignal struct {
	Kind       LandscapeSignalKind `json:"kind"`
	Value      float64             `json:"value"`
	Unit       string              `json:"unit"`
	WindowDays int                 `json:"window_days"`
	ObservedAt string              `json:"observed_at"`
	Sources    []string            `json:"sources"`
}

type LandscapeCandidate struct {
	ID           string                 `json:"id"`
	Kind         LandscapeCandidateKind `json:"kind"`
	Name         string                 `json:"name"`
	Capabilities []LandscapeCapability  `json:"capabilities"`
	Standalone   LandscapeAssertion     `json:"standalone"`
	Integration  LandscapeAssertion     `json:"integration"`
	Analogy      LandscapeAssertion     `json:"analogy"`
	Signals      []LandscapeSignal      `json:"signals,omitempty"`
}

type LandscapeSnapshot struct {
	AsOf       string               `json:"as_of"`
	Need       []string             `json:"need"`
	Candidates []LandscapeCandidate `json:"candidates"`
}

type TechnologyLandscapeRequest struct {
	SchemaVersion string             `json:"schema_version"`
	Current       LandscapeSnapshot  `json:"current"`
	Previous      *LandscapeSnapshot `json:"previous,omitempty"`
}

type LandscapeAssessment struct {
	CandidateID     string                 `json:"candidate_id"`
	Name            string                 `json:"name"`
	Kind            LandscapeCandidateKind `json:"kind"`
	Role            LandscapeRole          `json:"role"`
	Capabilities    []LandscapeCapability  `json:"capabilities"`
	Standalone      LandscapeAssertion     `json:"standalone"`
	Integration     LandscapeAssertion     `json:"integration"`
	Analogy         LandscapeAssertion     `json:"analogy"`
	Proven          []string               `json:"proven"`
	Absent          []string               `json:"absent"`
	Unknown         []string               `json:"unknown"`
	MissingEvidence []string               `json:"missing_evidence"`
	Signals         []LandscapeSignal      `json:"signals"`
}

type LandscapeSignalMovement struct {
	Kind     LandscapeSignalKind `json:"kind"`
	Movement string              `json:"movement"`
	Before   *float64            `json:"before,omitempty"`
	After    *float64            `json:"after,omitempty"`
}

type LandscapeChange struct {
	CandidateID string                    `json:"candidate_id"`
	Kind        string                    `json:"kind"`
	BeforeRole  LandscapeRole             `json:"before_role,omitempty"`
	AfterRole   LandscapeRole             `json:"after_role,omitempty"`
	Signals     []LandscapeSignalMovement `json:"signals,omitempty"`
}

type TechnologyLandscapeReport struct {
	SchemaVersion string                `json:"schema_version"`
	AsOf          string                `json:"as_of"`
	Need          []string              `json:"need"`
	Assessments   []LandscapeAssessment `json:"assessments"`
	Changes       []LandscapeChange     `json:"changes"`
	Counts        map[LandscapeRole]int `json:"counts"`
	Caveats       []string              `json:"caveats"`
}
