/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/algorithm_history.md
 * @layer L1
 * @updated 2026-09-12
 */

package domain

const (
	AlgorithmHistoryInputSchemaVersion  = "algorithm-history-input-v1"
	AlgorithmHistoryReportSchemaVersion = "algorithm-history-report-v1"
	AlgorithmHistoryReportComplete      = "COMPLETE"
	AlgorithmHistoryReportUnknown       = "UNKNOWN"
)

type AlgorithmChronology string

const (
	AlgorithmChronologyKnown       AlgorithmChronology = "KNOWN"
	AlgorithmChronologyApproximate AlgorithmChronology = "APPROXIMATE"
	AlgorithmChronologyUnknown     AlgorithmChronology = "UNKNOWN"
)

type AlgorithmHistoryProvenanceKind string

const (
	AlgorithmHistoryPrimarySource     AlgorithmHistoryProvenanceKind = "PRIMARY_SOURCE"
	AlgorithmHistorySecondarySource   AlgorithmHistoryProvenanceKind = "SECONDARY_SOURCE"
	AlgorithmHistoryEmpiricalArtifact AlgorithmHistoryProvenanceKind = "EMPIRICAL_ARTIFACT"
)

type AlgorithmHistoryRelation string

const (
	AlgorithmHistoryDerivedFrom  AlgorithmHistoryRelation = "DERIVED_FROM"
	AlgorithmHistoryInfluencedBy AlgorithmHistoryRelation = "INFLUENCED_BY"
	AlgorithmHistoryImproves     AlgorithmHistoryRelation = "IMPROVES"
	AlgorithmHistoryGeneralizes  AlgorithmHistoryRelation = "GENERALIZES"
	AlgorithmHistorySpecializes  AlgorithmHistoryRelation = "SPECIALIZES"
	AlgorithmHistoryCombines     AlgorithmHistoryRelation = "COMBINES"
	AlgorithmHistoryRediscovers  AlgorithmHistoryRelation = "REDISCOVERS"
	AlgorithmHistoryResembles    AlgorithmHistoryRelation = "RESEMBLES"
)

type AlgorithmHistoryBasisKind string

const (
	AlgorithmHistoryEvidence  AlgorithmHistoryBasisKind = "EVIDENCE"
	AlgorithmHistoryInference AlgorithmHistoryBasisKind = "INFERENCE"
)

type AlgorithmHistoryStance string

const (
	AlgorithmHistorySupport   AlgorithmHistoryStance = "SUPPORT"
	AlgorithmHistoryConflict  AlgorithmHistoryStance = "CONFLICT"
	AlgorithmHistoryAmbiguous AlgorithmHistoryStance = "AMBIGUOUS"
)

type AlgorithmHistoryOutcome string

const (
	AlgorithmHistoryCorroborated AlgorithmHistoryOutcome = "CORROBORATED"
	AlgorithmHistorySupported    AlgorithmHistoryOutcome = "SUPPORTED"
	AlgorithmHistoryContested    AlgorithmHistoryOutcome = "CONTESTED"
	AlgorithmHistoryContradicted AlgorithmHistoryOutcome = "CONTRADICTED"
	AlgorithmHistoryUnknown      AlgorithmHistoryOutcome = "UNKNOWN"
)

type AlgorithmHistoryInput struct {
	SchemaVersion string                       `json:"schema_version"`
	Algorithms    []AlgorithmHistoryAlgorithm  `json:"algorithms"`
	Provenance    []AlgorithmHistoryProvenance `json:"provenance"`
	Claims        []AlgorithmHistoryClaim      `json:"claims"`
	Bases         []AlgorithmHistoryBasis      `json:"bases"`
}

type AlgorithmHistoryAlgorithm struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Chronology       AlgorithmChronology `json:"chronology"`
	YearStart        *int                `json:"year_start"`
	YearEnd          *int                `json:"year_end"`
	ChronologyReason string              `json:"chronology_reason"`
	ProvenanceIDs    []string            `json:"provenance_ids"`
}

type AlgorithmHistoryProvenance struct {
	ID       string                         `json:"id"`
	Kind     AlgorithmHistoryProvenanceKind `json:"kind"`
	OriginID string                         `json:"origin_id"`
	Citation string                         `json:"citation"`
	Locator  string                         `json:"locator"`
}

type AlgorithmHistoryClaim struct {
	ID        string                   `json:"id"`
	SubjectID string                   `json:"subject_id"`
	Relation  AlgorithmHistoryRelation `json:"relation"`
	ObjectID  string                   `json:"object_id"`
	Statement string                   `json:"statement"`
}

type AlgorithmHistoryBasis struct {
	ID             string                    `json:"id"`
	ClaimID        string                    `json:"claim_id"`
	Kind           AlgorithmHistoryBasisKind `json:"kind"`
	Stance         AlgorithmHistoryStance    `json:"stance"`
	ProvenanceID   string                    `json:"provenance_id"`
	Summary        string                    `json:"summary"`
	DerivedFromIDs []string                  `json:"derived_from_ids"`
}

type AlgorithmHistoryAssessment struct {
	ClaimID                   string                   `json:"claim_id"`
	SubjectID                 string                   `json:"subject_id"`
	Relation                  AlgorithmHistoryRelation `json:"relation"`
	ObjectID                  string                   `json:"object_id"`
	Statement                 string                   `json:"statement"`
	Outcome                   AlgorithmHistoryOutcome  `json:"outcome"`
	SupportingEvidenceIDs     []string                 `json:"supporting_evidence_ids"`
	ConflictingEvidenceIDs    []string                 `json:"conflicting_evidence_ids"`
	AmbiguousEvidenceIDs      []string                 `json:"ambiguous_evidence_ids"`
	SupportingInferenceIDs    []string                 `json:"supporting_inference_ids"`
	ConflictingInferenceIDs   []string                 `json:"conflicting_inference_ids"`
	AmbiguousInferenceIDs     []string                 `json:"ambiguous_inference_ids"`
	DistinctSupportingOrigins []string                 `json:"distinct_supporting_origins"`
	Reason                    string                   `json:"reason"`
}

type AlgorithmHistoryMetrics struct {
	InputAlgorithms     int  `json:"input_algorithms"`
	InputProvenance     int  `json:"input_provenance"`
	InputClaims         int  `json:"input_claims"`
	InputBases          int  `json:"input_bases"`
	Algorithms          int  `json:"algorithms"`
	Provenance          int  `json:"provenance"`
	Claims              int  `json:"claims"`
	Bases               int  `json:"bases"`
	DuplicateAlgorithms int  `json:"duplicate_algorithms"`
	DuplicateProvenance int  `json:"duplicate_provenance"`
	DuplicateClaims     int  `json:"duplicate_claims"`
	DuplicateBases      int  `json:"duplicate_bases"`
	Evidence            int  `json:"evidence"`
	Inferences          int  `json:"inferences"`
	Corroborated        int  `json:"corroborated"`
	Supported           int  `json:"supported"`
	Contested           int  `json:"contested"`
	Contradicted        int  `json:"contradicted"`
	Unknown             int  `json:"unknown"`
	Operations          int  `json:"operations"`
	LimitExceeded       bool `json:"limit_exceeded"`
	VerdictChanges      int  `json:"verdict_changes"`
}

type AlgorithmHistoryReport struct {
	SchemaVersion string                       `json:"schema_version"`
	Status        string                       `json:"status"`
	Reason        string                       `json:"reason"`
	Assessments   []AlgorithmHistoryAssessment `json:"assessments"`
	Metrics       AlgorithmHistoryMetrics      `json:"metrics"`
}
