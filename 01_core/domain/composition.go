/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/composition.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

// CompositionInput contains externally checked architecture claims. Omitting it
// leaves every functional combination a hypothesis, never a verified solution.
type CompositionInput struct {
	Connections   []CompositionConnection   `json:"connections,omitempty"`
	Exclusions    []CompositionExclusion    `json:"exclusions,omitempty"`
	Requirements  []CompositionRequirement  `json:"requirements,omitempty"`
	Verifications []CompositionVerification `json:"verifications,omitempty"`
}

type CompositionConnection struct {
	From    string           `json:"from"`
	To      string           `json:"to"`
	Status  EvidenceStatus   `json:"status"`
	Sources []EvidenceSource `json:"sources,omitempty"`
	Cost    *IntegrationCost `json:"cost,omitempty"`
}

type CompositionExclusion struct {
	Left    string           `json:"left"`
	Right   string           `json:"right"`
	Sources []EvidenceSource `json:"sources,omitempty"`
	Reason  string           `json:"reason"`
}

type CompositionRequirement struct {
	Component string `json:"component"`
	Requires  string `json:"requires"`
}

type CompositionVerification struct {
	Members  []string         `json:"members"`
	Status   EvidenceStatus   `json:"status"`
	Executed bool             `json:"executed"`
	Passed   bool             `json:"passed"`
	Sources  []EvidenceSource `json:"sources,omitempty"`
}

type IntegrationCost struct {
	Unit string  `json:"unit"`
	Low  float64 `json:"low"`
	High float64 `json:"high"`
}

type CompositionStatus string

const (
	CompositionProposed     CompositionStatus = "PROPOSED"
	CompositionVerified     CompositionStatus = "VERIFIED"
	CompositionIncompatible CompositionStatus = "INCOMPATIBLE"
)

type CompositionMember struct {
	ID         string            `json:"id"`
	Repository CommitRef         `json:"repository"`
	Component  ReusableComponent `json:"component"`
}

type CompositionCandidate struct {
	ID              string              `json:"id"`
	Members         []CompositionMember `json:"members"`
	Coverage        []string            `json:"coverage"`
	Overlap         []string            `json:"overlap,omitempty"`
	Gaps            []string            `json:"gaps"`
	Status          CompositionStatus   `json:"status"`
	Unknowns        []string            `json:"unknowns,omitempty"`
	IntegrationCost *IntegrationCost    `json:"integration_cost,omitempty"`
	Evidence        []EvidenceSource    `json:"evidence,omitempty"`
}

type CompositionReport struct {
	SchemaVersion string                 `json:"schema_version"`
	Method        string                 `json:"method"`
	Candidates    []CompositionCandidate `json:"candidates"`
	Truncated     bool                   `json:"truncated"`
	StopReason    string                 `json:"stop_reason,omitempty"`
}
