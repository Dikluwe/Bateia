/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/investigation_plan.md
 * @layer L1
 * @updated 2026-09-07
 */
package domain

// InvestigationPlan distinguishes observed discovery from work still required.
type InvestigationPlan struct {
	SchemaVersion  string                     `json:"schema_version"`
	Requirements   []InvestigationRequirement `json:"requirements"`
	Searches       []InvestigationSearch      `json:"searches"`
	Actions        []InvestigationAction      `json:"actions"`
	StopReason     string                     `json:"stop_reason"`
	Scope          string                     `json:"scope"`
	TelemetryLabel string                     `json:"telemetry_label"`
}
type InvestigationRequirement struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Text       string `json:"text"`
	Question   string `json:"question"`
	Acceptance string `json:"acceptance"`
}
type InvestigationSearch struct {
	RequirementIDs []string `json:"requirement_ids"`
	ExpansionID    string   `json:"expansion_id"`
	Purpose        string   `json:"purpose"`
	OriginalQuery  string   `json:"original_query"`
	EffectiveQuery string   `json:"effective_query"`
	Relaxation     string   `json:"relaxation"`
	Provider       string   `json:"provider"`
	Error          string   `json:"error,omitempty"`
}
type InvestigationAction struct {
	Question       string           `json:"question"`
	Gap            string           `json:"gap"`
	RequirementID  string           `json:"requirement_id"`
	Candidate      string           `json:"candidate,omitempty"`
	Classification string           `json:"classification"`
	Status         EvidenceStatus   `json:"status"`
	Evidence       []EvidenceSource `json:"evidence"`
	NextAction     string           `json:"next_action"`
	Acceptance     string           `json:"acceptance"`
}
