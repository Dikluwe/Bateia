/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/discovery_ontology.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

type DiscoveryConceptID string

type ConceptQuery struct {
	Text           string   `json:"text"`
	Discriminators []string `json:"discriminators"`
	NegativeTerms  []string `json:"negative_terms,omitempty"`
	Language       string   `json:"language,omitempty"`
}
type DiscoveryConcept struct {
	ID        DiscoveryConceptID   `json:"id"`
	Triggers  []string             `json:"triggers"`
	Queries   []ConceptQuery       `json:"queries"`
	Related   []DiscoveryConceptID `json:"related,omitempty"`
	Forbidden []string             `json:"forbidden,omitempty"`
	Rationale string               `json:"rationale"`
}
type DiscoveryOntology struct {
	SchemaVersion string             `json:"schema_version"`
	Version       string             `json:"version"`
	Concepts      []DiscoveryConcept `json:"concepts"`
	Hash          string             `json:"hash"`
}
type ConceptMatch struct {
	ConceptID DiscoveryConceptID `json:"concept_id"`
	SourceID  string             `json:"source_id"`
	Trigger   string             `json:"trigger"`
	Strength  float64            `json:"strength"`
}
type OntologyExpansionResult struct {
	SchemaVersion   string            `json:"schema_version"`
	OntologyVersion string            `json:"ontology_version"`
	OntologyHash    string            `json:"ontology_hash"`
	Matches         []ConceptMatch    `json:"matches"`
	Expansions      []SearchExpansion `json:"expansions"`
	Uncovered       []string          `json:"uncovered"`
}
type QueryQualityVerdict string

const (
	QueryAccepted QueryQualityVerdict = "ACCEPTED"
	QueryWeak     QueryQualityVerdict = "WEAK"
	QueryRejected QueryQualityVerdict = "REJECTED"
)

type QueryQualityFinding struct {
	ExpansionID    string              `json:"expansion_id"`
	Query          string              `json:"query"`
	Verdict        QueryQualityVerdict `json:"verdict"`
	Discriminators []string            `json:"discriminators"`
	Missing        []string            `json:"missing,omitempty"`
	Reasons        []string            `json:"reasons"`
}
type QueryQualityReport struct {
	SchemaVersion   string                `json:"schema_version"`
	OntologyVersion string                `json:"ontology_version"`
	OntologyHash    string                `json:"ontology_hash"`
	Findings        []QueryQualityFinding `json:"findings"`
	Accepted        []SearchExpansion     `json:"accepted"`
	Rejected        int                   `json:"rejected"`
}
type OntologyEvaluationCase struct {
	ID                     string               `json:"id"`
	Group                  string               `json:"group"`
	Idea                   IdeaSpec             `json:"idea"`
	ExpectedConcepts       []DiscoveryConceptID `json:"expected_concepts"`
	ForbiddenConcepts      []DiscoveryConceptID `json:"forbidden_concepts"`
	RequiredDiscriminators []string             `json:"required_discriminators"`
}
type OntologyEvaluationCorpus struct {
	Cases []OntologyEvaluationCase `json:"cases"`
}
type OntologyMetrics struct {
	ConceptPrecision       *float64 `json:"concept_precision"`
	ConceptRecall          *float64 `json:"concept_recall"`
	QueryAcceptanceRate    float64  `json:"query_acceptance_rate"`
	DiscriminatorRetention float64  `json:"discriminator_retention"`
	GenericFallbackRate    float64  `json:"generic_fallback_rate"`
	ForbiddenLeakRate      float64  `json:"forbidden_leak_rate"`
	ControlContamination   float64  `json:"control_contamination"`
	UncoveredRate          float64  `json:"uncovered_rate"`
}
type OntologyCaseResult struct {
	ID        string               `json:"id"`
	Found     []DiscoveryConceptID `json:"found"`
	Missing   []DiscoveryConceptID `json:"missing"`
	Forbidden []DiscoveryConceptID `json:"forbidden"`
	Quality   QueryQualityReport   `json:"quality"`
}
type OntologyEvaluationReport struct {
	SchemaVersion   string               `json:"schema_version"`
	OntologyVersion string               `json:"ontology_version"`
	OntologyHash    string               `json:"ontology_hash"`
	Cases           []OntologyCaseResult `json:"cases"`
	Metrics         OntologyMetrics      `json:"metrics"`
	Passed          bool                 `json:"passed"`
}
