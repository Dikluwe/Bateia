/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/source_insight.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

type InsightObservationKind string

const (
	InsightStrategy     InsightObservationKind = "STRATEGY"
	InsightPrecondition InsightObservationKind = "PRECONDITION"
	InsightFailureMode  InsightObservationKind = "FAILURE_MODE"
	InsightTradeoff     InsightObservationKind = "TRADEOFF"
	InsightLimitation   InsightObservationKind = "LIMITATION"
)

type InsightObservationReference struct {
	Path        string `json:"path"`
	ContentHash string `json:"content_hash"`
	StartLine   int    `json:"start_line"`
	EndLine     int    `json:"end_line"`
}

// SourceInsightObservation is a supplied interpretation, not a verified effect.
type SourceInsightObservation struct {
	MechanismClass string                        `json:"mechanism_class"`
	Kind           InsightObservationKind        `json:"kind"`
	Statement      string                        `json:"statement"`
	References     []InsightObservationReference `json:"references"`
}

type LocatedInsightReference struct {
	InsightObservationReference
	Excerpt string `json:"excerpt"`
}

type LocatedInsightObservation struct {
	ID         string                    `json:"id"`
	Kind       InsightObservationKind    `json:"kind"`
	Statement  string                    `json:"statement"`
	Origin     string                    `json:"origin"`
	References []LocatedInsightReference `json:"references"`
}

type InsightCatalogDefinition struct {
	Origin         string `json:"origin"`
	MechanismClass string `json:"mechanism_class"`
	Name           string `json:"name"`
	Purpose        string `json:"purpose"`
	Strategy       string `json:"strategy"`
}

type InsightExtractionScope struct {
	Languages  []string `json:"languages"`
	Ecosystems []string `json:"ecosystems"`
}

type SourceInsightAnalysis struct {
	Source            ImmutableCandidateRef       `json:"source"`
	CatalogDefinition InsightCatalogDefinition    `json:"catalog_definition"`
	ExtractionScope   InsightExtractionScope      `json:"extraction_scope"`
	MechanismRef      string                      `json:"mechanism_ref,omitempty"`
	MechanismStatus   MechanismEvidenceStatus     `json:"mechanism_status"`
	AssessmentOrigin  string                      `json:"assessment_origin"`
	Evidence          []MechanismWitness          `json:"evidence"`
	CounterEvidence   []MechanismWitness          `json:"counter_evidence"`
	Observations      []LocatedInsightObservation `json:"observations"`
	Gaps              []string                    `json:"gaps"`
}

type SourceInsightInvestigationReport struct {
	SchemaVersion  string                     `json:"schema_version"`
	ID             string                     `json:"id"`
	Scope          string                     `json:"scope"`
	Status         ImprovementSynthesisStatus `json:"status"`
	Question       string                     `json:"question"`
	Target         InsightTarget              `json:"target"`
	InputSHA256    string                     `json:"input_sha256"`
	CatalogSHA256  string                     `json:"catalog_sha256"`
	SourceAnalyses []SourceInsightAnalysis    `json:"source_analyses"`
	Hypotheses     []ImprovementHypothesis    `json:"hypotheses"`
	ExtractionGaps []MechanismGap             `json:"extraction_gaps"`
	Provenance     []InsightSourceProvenance  `json:"provenance"`
	Limitations    []string                   `json:"limitations"`
}
