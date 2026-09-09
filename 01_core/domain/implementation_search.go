/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/implementation_search.md
 * @layer L1
 * @updated 2026-09-08
 */
package domain

import "context"

type ImplementationSearchBudget struct {
	MaxFiles     int   `json:"max_files"`
	MaxBytes     int64 `json:"max_bytes"`
	MaxFileBytes int64 `json:"max_file_bytes"`
	MaxResults   int   `json:"max_results"`
}
type ImplementationReference struct {
	Path   string `json:"path"`
	Symbol string `json:"symbol"`
	Line   int    `json:"line"`
	SHA256 string `json:"sha256"`
}
type ImplementationTermWeight struct {
	Term              string  `json:"term,omitempty"`
	Field             string  `json:"field,omitempty"`
	DocumentFrequency int     `json:"document_frequency,omitempty"`
	CorpusDocuments   int     `json:"corpus_documents,omitempty"`
	Weight            float64 `json:"weight,omitempty"`
}
type ImplementationHit struct {
	DocumentationLine    int                        `json:"documentation_line,omitempty"`
	DocumentationEndLine int                        `json:"documentation_end_line,omitempty"`
	TermWeights          []ImplementationTermWeight `json:"term_weights,omitempty"`
	Retrieval            string                     `json:"retrieval,omitempty"`
	RelatedTo            []ImplementationReference  `json:"related_to,omitempty"`
	Path                 string                     `json:"path"`
	Symbol               string                     `json:"symbol"`
	Kind                 string                     `json:"kind"`
	SHA256               string                     `json:"sha256"`
	Reason               string                     `json:"reason"`
	Line                 int                        `json:"line"`
	EndLine              int                        `json:"end_line"`
	Score                int                        `json:"score"`
	MatchedTerms         []string                   `json:"matched_terms"`
}
type ImplementationInput struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Bytes  int64  `json:"bytes"`
}

// ImplementationQuery distinguishes lexical terms from explicitly requested symbols.
type ImplementationQuery struct {
	Original    string                          `json:"original"`
	Mode        string                          `json:"mode"`
	Terms       []string                        `json:"terms"`
	Symbols     []string                        `json:"symbols"`
	Constraints []ImplementationQueryConstraint `json:"constraints,omitempty"`
	Expansions  []ImplementationQueryExpansion  `json:"expansions,omitempty"`
}

// ImplementationQueryExpansion records lexical hypotheses, not behavioral equivalence.
type ImplementationQueryExpansion struct {
	Source string   `json:"source"`
	Terms  []string `json:"terms"`
	RuleID string   `json:"rule_id"`
	Kind   string   `json:"kind"`
	Status string   `json:"status"`
}

// ImplementationQueryConstraint records wording, never proof of absence or behavior.
type ImplementationQueryConstraint struct {
	Marker string `json:"marker"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

// RequirementRetrieval describes lexical evidence, never satisfaction of a requirement.
type RequirementRetrieval struct {
	DocumentationSources []ImplementationReference       `json:"documentation_sources,omitempty"`
	Constraints          []ImplementationQueryConstraint `json:"constraints,omitempty"`
	Requirement          string                          `json:"requirement"`
	Status               string                          `json:"status"`
	Reason               string                          `json:"reason"`
	MatchedTerms         []string                        `json:"matched_terms"`
	MissingTerms         []string                        `json:"missing_terms"`
	Sources              []ImplementationReference       `json:"sources"`
}
type ImplementationSearchReport struct {
	ActionAssessments   []RequirementActionAssessment `json:"action_assessments,omitempty"`
	Behavior            []ImplementationBehavior      `json:"behavior,omitempty"`
	QueryInterpretation []ImplementationQuery         `json:"query_interpretation,omitempty"`
	Requirements        []RequirementRetrieval        `json:"requirements,omitempty"`
	Status              string                        `json:"status"`
	Method              string                        `json:"method"`
	Reason              string                        `json:"reason"`
	Hits                []ImplementationHit           `json:"hits"`
	Examined            []ImplementationInput         `json:"examined"`
	FilesRead           int                           `json:"files_read"`
	BytesRead           int64                         `json:"bytes_read"`
	Budget              ImplementationSearchBudget    `json:"budget"`
	Limited             bool                          `json:"limited"`
	Issues              []string                      `json:"issues,omitempty"`
}
type ImplementationSearcher interface {
	SearchImplementations(context.Context, RepositoryCandidate, string, ImplementationSearchBudget) (ImplementationSearchReport, error)
}
