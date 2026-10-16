/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_investigation.md
 * @layer L1
 * @updated 2026-10-06
 */

package domain

// InsightTarget describes a desired effect, not an observed target capability.
type InsightTarget struct {
	Name          string   `json:"name"`
	Context       string   `json:"context"`
	Property      string   `json:"property"`
	DesiredEffect string   `json:"desired_effect"`
	Constraints   []string `json:"constraints"`
}

type InsightEvidenceInput struct {
	Bundle      EvidenceBundle        `json:"bundle"`
	Assessments []CandidateAssessment `json:"assessments"`
}

type InsightInvestigationRequest struct {
	SchemaVersion    string                 `json:"schema_version"`
	Question         string                 `json:"question"`
	Target           InsightTarget          `json:"target"`
	MechanismClasses []string               `json:"mechanism_classes"`
	Sources          []InsightEvidenceInput `json:"sources"`
	MaxHypotheses    int                    `json:"max_hypotheses,omitempty"`
}

type InsightSourceProvenance struct {
	Source            ImmutableCandidateRef `json:"source"`
	BundleSHA256      string                `json:"bundle_sha256"`
	AssessmentsSHA256 string                `json:"assessments_sha256"`
	Complete          bool                  `json:"complete"`
	AssessmentOrigin  string                `json:"assessment_origin"`
}

type InsightInvestigationReport struct {
	SchemaVersion  string                     `json:"schema_version"`
	ID             string                     `json:"id"`
	Scope          string                     `json:"scope"`
	Status         ImprovementSynthesisStatus `json:"status"`
	Question       string                     `json:"question"`
	Target         InsightTarget              `json:"target"`
	InputSHA256    string                     `json:"input_sha256"`
	CatalogSHA256  string                     `json:"catalog_sha256"`
	Mechanisms     []ArchitecturalMechanism   `json:"mechanisms"`
	Hypotheses     []ImprovementHypothesis    `json:"hypotheses"`
	ExtractionGaps []MechanismGap             `json:"extraction_gaps"`
	Provenance     []InsightSourceProvenance  `json:"provenance"`
	Limitations    []string                   `json:"limitations"`
}
