/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_site_reading.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import "context"

// InsightSiteRepositoryReader resolves an origin before reading only that forge.
type InsightSiteRepositoryReader interface {
	InsightReferenceReader
	ResolveRepository(context.Context, ProductIdentity) (RepositoryCandidate, error)
}

// InsightReadmeFailure never carries a remote response or credential.
type InsightReadmeFailure struct{ Code string }

func (e InsightReadmeFailure) Error() string { return e.Code }

// Legacy bundles have no site field; this wrapper owns their origin explicitly.
type InsightSiteContext struct {
	CandidateID string               `json:"candidate_id"`
	Identity    ProductIdentity      `json:"identity"`
	Bundle      EvidenceBundle       `json:"bundle"`
	Cards       []InsightReadingCard `json:"cards"`
}
