/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_radar_plan.md
 * @layer L1
 * @updated 2026-09-11
 */

package domain

const IdeaRadarPlanReportSchemaVersion = "idea-radar-plan-report-v1"

const (
	IdeaRadarPlanReportComplete = "COMPLETE"
	IdeaRadarPlanReportUnknown  = "UNKNOWN"
)

type IdeaRadarProbeKind string

const (
	IdeaRadarProbeResolveConflict          IdeaRadarProbeKind = "RESOLVE_CONFLICT"
	IdeaRadarProbeResolveAmbiguity         IdeaRadarProbeKind = "RESOLVE_AMBIGUITY"
	IdeaRadarProbeRecordOrigin             IdeaRadarProbeKind = "RECORD_ORIGIN"
	IdeaRadarProbeAddPerspective           IdeaRadarProbeKind = "ADD_PERSPECTIVE"
	IdeaRadarProbeAddIndependentOrigin     IdeaRadarProbeKind = "ADD_INDEPENDENT_ORIGIN"
	IdeaRadarProbeBindEvidence             IdeaRadarProbeKind = "BIND_EVIDENCE"
	IdeaRadarProbeKindResolveConflict                         = IdeaRadarProbeResolveConflict
	IdeaRadarProbeKindResolveAmbiguity                        = IdeaRadarProbeResolveAmbiguity
	IdeaRadarProbeKindRecordOrigin                            = IdeaRadarProbeRecordOrigin
	IdeaRadarProbeKindAddPerspective                          = IdeaRadarProbeAddPerspective
	IdeaRadarProbeKindAddIndependentOrigin                    = IdeaRadarProbeAddIndependentOrigin
	IdeaRadarProbeKindBindEvidence                            = IdeaRadarProbeBindEvidence
)

const (
	IdeaRadarProbePriorityResolveConflict      = 10
	IdeaRadarProbePriorityResolveAmbiguity     = 20
	IdeaRadarProbePriorityRecordOrigin         = 30
	IdeaRadarProbePriorityAddPerspective       = 40
	IdeaRadarProbePriorityAddIndependentOrigin = 50
	IdeaRadarProbePriorityBindEvidence         = 60
)

type IdeaRadarOriginRequirement string

const (
	IdeaRadarOriginRequirementNone             IdeaRadarOriginRequirement = "NONE"
	IdeaRadarOriginRequirementDistinctRecorded IdeaRadarOriginRequirement = "DISTINCT_RECORDED"
	IdeaRadarOriginRequirementRecorded         IdeaRadarOriginRequirement = "RECORDED"
)

type IdeaRadarProbe struct {
	ID                string                     `json:"id"`
	Kind              IdeaRadarProbeKind         `json:"kind"`
	Priority          int                        `json:"priority"`
	TargetID          string                     `json:"target_id"`
	ObservationIDs    []string                   `json:"observation_ids"`
	ProjectionKind    IdeaProjectionKind         `json:"projection_kind"`
	OriginRequirement IdeaRadarOriginRequirement `json:"origin_requirement"`
	RequiredRoles     []EvidenceCoalitionRole    `json:"required_roles"`
	Reason            string                     `json:"reason"`
}

type IdeaRadarTargetPlan struct {
	TargetID            string           `json:"target_id"`
	RadarOutcome        IdeaRadarOutcome `json:"radar_outcome"`
	Frontier            []IdeaRadarProbe `json:"frontier"`
	EstimatedProbeCount int              `json:"estimated_probe_count"`
	Reason              string           `json:"reason"`
}

type IdeaRadarPlanMetrics struct {
	Targets            int  `json:"targets"`
	Plans              int  `json:"plans"`
	Probes             int  `json:"probes"`
	RadarOperations    int  `json:"radar_operations"`
	PlanningOperations int  `json:"planning_operations"`
	Operations         int  `json:"operations"`
	LimitExceeded      bool `json:"limit_exceeded"`
	VerdictChanges     int  `json:"verdict_changes"`
}

type IdeaRadarPlanReport struct {
	SchemaVersion string                `json:"schema_version"`
	Status        string                `json:"status"`
	Reason        string                `json:"reason"`
	Plans         []IdeaRadarTargetPlan `json:"plans"`
	Metrics       IdeaRadarPlanMetrics  `json:"metrics"`
}
