/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/capability_compilation.md
 * @layer L1
 * @updated 2026-09-24
 */
package domain

// CapabilityCompilation is a plan of verification, never evidence about a candidate.
type CapabilityCompilation struct {
	SchemaVersion string                     `json:"schema_version"`
	IntakeHash    string                     `json:"intake_hash"`
	IntakeStatus  IntakeStatus               `json:"intake_status"`
	Goal          string                     `json:"goal"`
	Users         []string                   `json:"users"`
	Requirements  []CompiledRequirement      `json:"requirements"`
	Questions     []CompiledEvidenceQuestion `json:"questions"`
}

type CompiledRequirement struct {
	ID         string         `json:"id"`
	Kind       string         `json:"kind"`
	Text       string         `json:"text"`
	SourcePath string         `json:"source_path"`
	Authority  InputAuthority `json:"authority"`
	Priority   string         `json:"priority"`
}

type CompiledEvidenceQuestion struct {
	ID            string `json:"id"`
	RequirementID string `json:"requirement_id"`
	Channel       string `json:"channel"`
	Question      string `json:"question"`
	Support       string `json:"support"`
	Refute        string `json:"refute"`
	Limit         string `json:"limit"`
	Status        string `json:"status"`
}
