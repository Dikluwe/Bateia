/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_exploration_memory.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

type InsightReadObservation struct {
	Source      CommitRef `json:"source"`
	Path        string    `json:"path"`
	ContentHash string    `json:"content_hash"`
}
type InsightMemoryEntry struct {
	ID              string                  `json:"id"`
	Identity        ProductIdentity         `json:"identity"`
	State           string                  `json:"state"`
	Reason          string                  `json:"reason,omitempty"`
	FirstGeneration int                     `json:"first_generation"`
	LastQuestion    string                  `json:"last_question"`
	ReadAttempts    int                     `json:"read_attempts"`
	LastStatus      string                  `json:"last_status,omitempty"`
	Lead            *InsightReferenceLead   `json:"lead,omitempty"`
	LastRead        *InsightReadObservation `json:"last_read,omitempty"`
}
type InsightExplorationMemory struct {
	SchemaVersion string               `json:"schema_version"`
	ID            string               `json:"id"`
	Generation    int                  `json:"generation"`
	Entries       []InsightMemoryEntry `json:"entries"`
}
type InsightFollowup struct {
	EntryID  string                `json:"entry_id"`
	LeadID   string                `json:"lead_id"`
	Identity ProductIdentity       `json:"identity"`
	Status   string                `json:"status"`
	Receipt  *InsightReadmeReceipt `json:"receipt,omitempty"`
	Bundle   *EvidenceBundle       `json:"bundle,omitempty"`
	Cards    []InsightReadingCard  `json:"cards"`
}
