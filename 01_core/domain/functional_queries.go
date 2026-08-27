/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/functional_queries.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

type FunctionalFrameID string
type FunctionalFrame struct {
	ID            FunctionalFrameID  `json:"id"`
	ConceptID     DiscoveryConceptID `json:"concept_id"`
	SourceIDs     []string           `json:"source_ids"`
	Inputs        []string           `json:"inputs"`
	Actions       []string           `json:"actions"`
	Objects       []string           `json:"objects"`
	Outputs       []string           `json:"outputs"`
	Timing        []string           `json:"timing,omitempty"`
	Constraints   []string           `json:"constraints,omitempty"`
	NegativeJobs  []string           `json:"negative_jobs,omitempty"`
	RuleID        string             `json:"rule_id"`
	UnknownReason string             `json:"unknown_reason,omitempty"`
}
type FunctionalLexicon struct {
	SchemaVersion string                                 `json:"schema_version"`
	Version       string                                 `json:"version"`
	Frames        map[DiscoveryConceptID]FunctionalFrame `json:"frames"`
	Forbidden     []string                               `json:"forbidden"`
	Hash          string                                 `json:"hash"`
}
type FunctionalQueryTemplate struct {
	ID            string   `json:"id"`
	Pattern       string   `json:"pattern"`
	RequiredSlots []string `json:"required_slots"`
	MaxVariants   int      `json:"max_variants"`
	Surface       string   `json:"surface"`
}
type FunctionalQuery struct {
	ID          string             `json:"id"`
	FrameID     FunctionalFrameID  `json:"frame_id"`
	ConceptID   DiscoveryConceptID `json:"concept_id"`
	TemplateID  string             `json:"template_id"`
	Text        string             `json:"text"`
	Inputs      []string           `json:"inputs"`
	Actions     []string           `json:"actions"`
	Objects     []string           `json:"objects"`
	Outputs     []string           `json:"outputs"`
	Timing      []string           `json:"timing,omitempty"`
	SourceIDs   []string           `json:"source_ids"`
	LexiconHash string             `json:"lexicon_hash"`
}
type FunctionalExpansion struct {
	SchemaVersion string            `json:"schema_version"`
	LexiconHash   string            `json:"lexicon_hash"`
	Frames        []FunctionalFrame `json:"frames"`
	Queries       []FunctionalQuery `json:"queries"`
	Unknown       []FunctionalFrame `json:"unknown"`
}
type FunctionalQueryFinding struct {
	QueryID       string   `json:"query_id"`
	Verdict       string   `json:"verdict"`
	MissingSlots  []string `json:"missing_slots,omitempty"`
	ForbiddenHits []string `json:"forbidden_hits,omitempty"`
	Reasons       []string `json:"reasons"`
}
