/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/discovery_coverage.md
 * @layer L1
 * @updated 2026-09-08
 */
package domain

// DiscoveryCoverage describes only the received response, never index coverage.
type DiscoveryCoverage struct {
	RequestedMatches     int    `json:"requested_matches"`
	ReceivedItems        int    `json:"received_items"`
	ObservedRepositories int    `json:"observed_repositories"`
	ExaminedRepositories int    `json:"examined_repositories"`
	DeferredRepositories int    `json:"deferred_repositories"`
	EligibleCandidates   int    `json:"eligible_candidates"`
	SelectedCandidates   int    `json:"selected_candidates"`
	Completeness         string `json:"completeness"`
	ProviderMatches      *int   `json:"provider_matches"`
}

// DiscoveryCandidateDiagnostic records retained lexical observations, not capability absence.
type DiscoveryCandidateDiagnostic struct {
	Repository        string          `json:"repository"`
	Status            string          `json:"status"`
	Reason            string          `json:"reason"`
	Evidence          []MatchEvidence `json:"evidence"`
	ActionFound       bool            `json:"action_found"`
	ObjectFound       bool            `json:"object_found"`
	TruncatedSnippets int             `json:"truncated_snippets"`
	OmittedSnippets   int             `json:"omitted_snippets"`
}
