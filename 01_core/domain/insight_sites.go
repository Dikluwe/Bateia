/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_sites.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

type InsightSiteBinding struct {
	FollowReader     InsightSiteRepositoryReader
	Connector        MarketDiscoveryConnector
	Reader           InsightReferenceReader
	RepositoryReader InsightSiteRepositoryReader
}
type InsightSiteFactory func() ([]InsightSiteBinding, error)
type InsightSiteRequest struct {
	Memory      *InsightExplorationMemory `json:"memory,omitempty"`
	FollowLinks *int                      `json:"follow_links,omitempty"`
	Discovery   InsightDiscoveryRequest   `json:"discovery"`
	Sources     []string                  `json:"sources"`
	MaxCalls    int                       `json:"max_calls"`
	Locale      string                    `json:"locale"`
}
type InsightSitePlan struct {
	SchemaVersion    string                           `json:"schema_version"`
	ID               string                           `json:"id"`
	Request          InsightSiteRequest               `json:"request"`
	Discovery        InsightDiscoveryPlan             `json:"discovery"`
	SchedulingPolicy string                           `json:"scheduling_policy,omitempty"`
	Sources          []DiscoveryConnectorCapabilities `json:"sources"`
}
type InsightSiteAttempt struct {
	ExecutionOrder int    `json:"execution_order,omitempty"`
	ConnectorID    string `json:"connector_id"`
	QueryID        string `json:"query_id"`
	Market         string `json:"market"`
	Status         string `json:"status"`
	ErrorCode      string `json:"error_code,omitempty"`
	ReportedCount  int    `json:"reported_count"`
	AdmittedCount  int    `json:"admitted_count"`
	ResponseHash   string `json:"response_hash,omitempty"`
}
type InsightSiteCandidate struct {
	ID                       string           `json:"id"`
	Product                  ProductCandidate `json:"product"`
	QueryIDs                 []string         `json:"query_ids"`
	Origin                   string           `json:"origin"`
	FamiliarSubstringMatches []string         `json:"familiar_substring_matches"`
	Selected                 bool             `json:"selected"`
	SelectionOrder           int              `json:"selection_order"`
	SelectionReasons         []string         `json:"selection_reasons"`
}
type InsightSiteReadme struct {
	CandidateID string                `json:"candidate_id"`
	ConnectorID string                `json:"connector_id"`
	Status      string                `json:"status"`
	Receipt     *InsightReadmeReceipt `json:"receipt,omitempty"`
}
type InsightSiteTask struct {
	MemoryEntryRefs     []string `json:"memory_entry_refs,omitempty"`
	ReferenceLeadRefs   []string `json:"reference_lead_refs,omitempty"`
	Perspective         string   `json:"perspective"`
	Question            string   `json:"question"`
	ReadingCardRefs     []string `json:"reading_card_refs"`
	CatalogReferenceIDs []string `json:"catalog_reference_ids"`
}
type InsightSiteReport struct {
	SchemaVersion       string                      `json:"schema_version"`
	ID                  string                      `json:"id"`
	Status              string                      `json:"status"`
	Plan                InsightSitePlan             `json:"plan"`
	Calls               int                         `json:"calls"`
	Attempts            []InsightSiteAttempt        `json:"attempts"`
	Candidates          []InsightSiteCandidate      `json:"candidates"`
	Readmes             []InsightSiteReadme         `json:"readmes"`
	Bundles             []EvidenceBundle            `json:"bundles"`
	Exploration         InsightExplorationReport    `json:"exploration"`
	SiteContexts        []InsightSiteContext        `json:"site_contexts,omitempty"`
	OpportunityCoverage *InsightOpportunityCoverage `json:"opportunity_coverage,omitempty"`
	ExplorationMemory   *InsightExplorationMemory   `json:"exploration_memory,omitempty"`
	Followups           []InsightFollowup           `json:"followups,omitempty"`
	ReferenceLeads      []InsightReferenceLead      `json:"reference_leads,omitempty"`
	Tasks               []InsightSiteTask           `json:"tasks"`
	AIInstructions      []string                    `json:"ai_instructions"`
	Gaps                []string                    `json:"gaps"`
	NoveltyState        string                      `json:"novelty_state"`
	EffectState         string                      `json:"effect_state"`
	Neighborhood        *InsightNeighborhood        `json:"neighborhood,omitempty"`
}
