/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_radar.md
 * @layer L1
 * @updated 2026-09-11
 */

package domain

const (
	IdeaRadarInputSchemaVersion  = "idea-radar-input-v1"
	IdeaRadarReportSchemaVersion = "idea-radar-report-v1"
)

type IdeaProjectionKind string

const (
	IdeaProjectionProblem      IdeaProjectionKind = "PROBLEM"
	IdeaProjectionBehavior     IdeaProjectionKind = "BEHAVIOR"
	IdeaProjectionArchitecture IdeaProjectionKind = "ARCHITECTURE"
	IdeaProjectionConstraint   IdeaProjectionKind = "CONSTRAINT"
	IdeaProjectionOperation    IdeaProjectionKind = "OPERATION"

	IdeaProjectionKindProblem      = IdeaProjectionProblem
	IdeaProjectionKindBehavior     = IdeaProjectionBehavior
	IdeaProjectionKindArchitecture = IdeaProjectionArchitecture
	IdeaProjectionKindConstraint   = IdeaProjectionConstraint
	IdeaProjectionKindOperation    = IdeaProjectionOperation
)

type IdeaRadarStance string

const (
	IdeaRadarStanceSupport   IdeaRadarStance = "SUPPORT"
	IdeaRadarStanceConflict  IdeaRadarStance = "CONFLICT"
	IdeaRadarStanceAmbiguous IdeaRadarStance = "AMBIGUOUS"
)

type IdeaRadarOutcome string

const (
	IdeaRadarConvergent               IdeaRadarOutcome = "CONVERGENT"
	IdeaRadarConflict                 IdeaRadarOutcome = "CONFLICT"
	IdeaRadarAmbiguous                IdeaRadarOutcome = "AMBIGUOUS"
	IdeaRadarInsufficientPerspectives IdeaRadarOutcome = "INSUFFICIENT_PERSPECTIVES"
	IdeaRadarUnknownOrigin            IdeaRadarOutcome = "UNKNOWN_ORIGIN"
	IdeaRadarInsufficientOrigins      IdeaRadarOutcome = "INSUFFICIENT_ORIGINS"

	IdeaRadarOutcomeConvergent               = IdeaRadarConvergent
	IdeaRadarOutcomeConflict                 = IdeaRadarConflict
	IdeaRadarOutcomeAmbiguous                = IdeaRadarAmbiguous
	IdeaRadarOutcomeInsufficientPerspectives = IdeaRadarInsufficientPerspectives
	IdeaRadarOutcomeUnknownOrigin            = IdeaRadarUnknownOrigin
	IdeaRadarOutcomeInsufficientOrigins      = IdeaRadarInsufficientOrigins
)

const (
	IdeaRadarReportComplete = "COMPLETE"
	IdeaRadarReportUnknown  = "UNKNOWN"
)

type IdeaRadarInput struct {
	SchemaVersion string                 `json:"schema_version"`
	Projections   []IdeaRadarProjection  `json:"projections"`
	Observations  []IdeaRadarObservation `json:"observations"`
}

type IdeaRadarProjection struct {
	ID           string               `json:"id"`
	Kind         IdeaProjectionKind   `json:"kind"`
	OriginStatus EvidenceOriginStatus `json:"origin_status"`
	OriginID     string               `json:"origin_id"`
	Terms        []string             `json:"terms"`
	Actions      []string             `json:"actions"`
	Objects      []string             `json:"objects"`
	Required     []string             `json:"required"`
}

type IdeaRadarObservation struct {
	ID           string          `json:"id"`
	ProjectionID string          `json:"projection_id"`
	CandidateID  string          `json:"candidate_id"`
	ConceptID    string          `json:"concept_id"`
	Stance       IdeaRadarStance `json:"stance"`
}

type IdeaRadarAssessment struct {
	TargetID                    string                  `json:"target_id"`
	CandidateID                 string                  `json:"candidate_id"`
	ConceptID                   string                  `json:"concept_id"`
	Outcome                     IdeaRadarOutcome        `json:"outcome"`
	SupportingObservationIDs    []string                `json:"supporting_observation_ids"`
	ConflictingObservationIDs   []string                `json:"conflicting_observation_ids"`
	AmbiguousObservationIDs     []string                `json:"ambiguous_observation_ids"`
	UnknownOriginObservationIDs []string                `json:"unknown_origin_observation_ids"`
	ProjectionKinds             []IdeaProjectionKind    `json:"projection_kinds"`
	DistinctOrigins             []string                `json:"distinct_origins"`
	MissingProjectionKinds      int                     `json:"missing_projection_kinds"`
	MissingOrigins              int                     `json:"missing_origins"`
	CoalitionDraft              IdeaRadarCoalitionDraft `json:"coalition_draft"`
	Reason                      string                  `json:"reason"`
}

type IdeaRadarCoalitionDraft struct {
	ClaimID        string                  `json:"claim_id"`
	ObservationIDs []string                `json:"observation_ids"`
	RequiredRoles  []EvidenceCoalitionRole `json:"required_roles"`
	ReadyFor0065   bool                    `json:"ready_for_0065"`
	Reason         string                  `json:"reason"`
}

type IdeaRadarMetrics struct {
	Projections              int  `json:"projections"`
	Observations             int  `json:"observations"`
	CanonicalObservations    int  `json:"canonical_observations"`
	DuplicateObservations    int  `json:"duplicate_observations"`
	Targets                  int  `json:"targets"`
	Convergent               int  `json:"convergent"`
	Conflicts                int  `json:"conflicts"`
	Ambiguous                int  `json:"ambiguous"`
	InsufficientPerspectives int  `json:"insufficient_perspectives"`
	UnknownOrigins           int  `json:"unknown_origins"`
	InsufficientOrigins      int  `json:"insufficient_origins"`
	Operations               int  `json:"operations"`
	LimitExceeded            bool `json:"limit_exceeded"`
	VerdictChanges           int  `json:"verdict_changes"`
}

type IdeaRadarReport struct {
	SchemaVersion string                `json:"schema_version"`
	Status        string                `json:"status"`
	Reason        string                `json:"reason"`
	Assessments   []IdeaRadarAssessment `json:"assessments"`
	Metrics       IdeaRadarMetrics      `json:"metrics"`
}
