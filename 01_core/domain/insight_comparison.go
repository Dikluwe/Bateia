/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_comparison.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

type InsightBaselineFile struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	ContentHash string `json:"content_hash"`
}

// InsightTargetBaseline describes supplied files, not a Git commit or executed behavior.
type InsightTargetBaseline struct {
	Repository     string                     `json:"repository"`
	SnapshotSHA256 string                     `json:"snapshot_sha256"`
	Files          []InsightBaselineFile      `json:"files"`
	Observations   []SourceInsightObservation `json:"observations"`
}

type InsightComparisonRequest struct {
	SchemaVersion  string                      `json:"schema_version"`
	Investigation  InsightInvestigationRequest `json:"investigation"`
	TargetBaseline *InsightTargetBaseline      `json:"target_baseline"`
}

type TargetInsightObservation struct {
	MechanismClass string `json:"mechanism_class"`
	LocatedInsightObservation
}

type LocatedInsightBaseline struct {
	Repository     string                     `json:"repository"`
	SnapshotSHA256 string                     `json:"snapshot_sha256"`
	Scope          string                     `json:"scope"`
	Origin         string                     `json:"origin"`
	Observations   []TargetInsightObservation `json:"observations"`
}

// InsightComparisonRow records declaration availability, not semantic equivalence.
type InsightComparisonRow struct {
	ID                    string                 `json:"id"`
	Source                ImmutableCandidateRef  `json:"source"`
	MechanismClass        string                 `json:"mechanism_class"`
	Kind                  InsightObservationKind `json:"kind"`
	DeclarationStatus     string                 `json:"declaration_status"`
	SourceObservationRefs []string               `json:"source_observation_refs"`
	TargetObservationRefs []string               `json:"target_observation_refs"`
	Gaps                  []string               `json:"gaps"`
	NextStep              ValidationStep         `json:"next_step"`
}

type InsightComparisonReport struct {
	SchemaVersion        string                     `json:"schema_version"`
	ID                   string                     `json:"id"`
	Status               ImprovementSynthesisStatus `json:"status"`
	Question             string                     `json:"question"`
	Target               InsightTarget              `json:"target"`
	SourceReportID       string                     `json:"source_report_id"`
	SourceInputSHA256    string                     `json:"source_input_sha256"`
	CatalogSHA256        string                     `json:"catalog_sha256"`
	SourceAnalyses       []SourceInsightAnalysis    `json:"source_analyses"`
	SourceProvenance     []InsightSourceProvenance  `json:"source_provenance"`
	SourceExtractionGaps []MechanismGap             `json:"source_extraction_gaps"`
	TargetBaseline       *LocatedInsightBaseline    `json:"target_baseline"`
	Comparisons          []InsightComparisonRow     `json:"comparisons"`
	EffectState          string                     `json:"effect_state"`
	AdoptionDisposition  string                     `json:"adoption_disposition"`
	Limitations          []string                   `json:"limitations"`
}
