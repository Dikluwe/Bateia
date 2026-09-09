/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/discovery_query.md
 * @layer L1
 * @updated 2026-09-08
 */
package domain

// DiscoveryQueryPlan is a lexical hypothesis, never proof of relevance.
type DiscoveryQueryPlan struct {
	Original    string   `json:"original"`
	Version     string   `json:"version"`
	Status      string   `json:"status"`
	Reason      string   `json:"reason"`
	ActionTerms []string `json:"action_terms"`
	ObjectTerms []string `json:"object_terms"`
	RuleIDs     []string `json:"rule_ids"`
	Unverified  []string `json:"unverified"`
}
