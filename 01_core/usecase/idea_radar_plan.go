/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_radar_plan.md
 * @layer L1
 * @updated 2026-09-11
 */

package usecase

import (
	"fmt"
	"sort"

	"bateia/01_core/domain"
)

const (
	ideaRadarPlanMaxPlans      = 512
	ideaRadarPlanMaxProbes     = 1024
	ideaRadarPlanMaxOperations = 8192
)

var ideaRadarPlanKindOrder = []domain.IdeaProjectionKind{
	domain.IdeaProjectionProblem,
	domain.IdeaProjectionBehavior,
	domain.IdeaProjectionArchitecture,
	domain.IdeaProjectionConstraint,
	domain.IdeaProjectionOperation,
}

type ideaRadarPlanProjection struct {
	kind         domain.IdeaProjectionKind
	originStatus domain.EvidenceOriginStatus
}

// PlanIdeaRadarProbes describes only the minimum structural frontier for the
// blocker classified by AnalyzeIdeaRadar. It does not execute any probe.
func PlanIdeaRadarProbes(input domain.IdeaRadarInput) domain.IdeaRadarPlanReport {
	radar := AnalyzeIdeaRadar(input)
	report := newIdeaRadarPlanReport(radar)
	if radar.Status != domain.IdeaRadarReportComplete {
		report.Reason = radar.Reason
		return report
	}
	if len(radar.Assessments) > ideaRadarPlanMaxPlans {
		return unknownIdeaRadarPlanReport(report)
	}

	projections := make(map[string]ideaRadarPlanProjection, len(input.Projections))
	for _, projection := range input.Projections {
		id, _ := canonicalIdeaRadarID(projection.ID)
		projections[id] = ideaRadarPlanProjection{kind: projection.Kind, originStatus: projection.OriginStatus}
		if !countIdeaRadarPlanningOperation(&report, 1) {
			return unknownIdeaRadarPlanReport(report)
		}
	}

	observations := make(map[string]ideaRadarObservation, len(input.Observations))
	for _, observation := range input.Observations {
		normalized, _ := normalizeIdeaRadarObservation(observation)
		if prior, exists := observations[normalized.id]; !exists || normalized.projectionID < prior.projectionID {
			observations[normalized.id] = normalized
		}
		if !countIdeaRadarPlanningOperation(&report, 1) {
			return unknownIdeaRadarPlanReport(report)
		}
	}

	plans := make([]domain.IdeaRadarTargetPlan, 0, len(radar.Assessments))
	probeCount := 0
	for _, assessment := range radar.Assessments {
		if !countIdeaRadarPlanningOperation(&report, 1) {
			return unknownIdeaRadarPlanReport(report)
		}
		frontier := planIdeaRadarTarget(assessment, observations, projections)
		probeCount += len(frontier)
		if probeCount > ideaRadarPlanMaxProbes || !countIdeaRadarPlanningOperation(&report, len(frontier)) {
			return unknownIdeaRadarPlanReport(report)
		}
		sortIdeaRadarProbes(frontier)
		for index := range frontier {
			frontier[index].ID = fmt.Sprintf("%s::p%04d", assessment.TargetID, index+1)
		}
		reason := "active_blocker_frontier"
		if assessment.Outcome == domain.IdeaRadarConvergent {
			reason = "evidence_binding_only"
		}
		plans = append(plans, domain.IdeaRadarTargetPlan{
			TargetID:            assessment.TargetID,
			RadarOutcome:        assessment.Outcome,
			Frontier:            frontier,
			EstimatedProbeCount: len(frontier),
			Reason:              reason,
		})
	}

	sort.Slice(plans, func(i, j int) bool { return plans[i].TargetID < plans[j].TargetID })
	report.Status = domain.IdeaRadarPlanReportComplete
	report.Reason = "planned"
	report.Plans = plans
	report.Metrics.Plans = len(plans)
	report.Metrics.Probes = probeCount
	return report
}

func newIdeaRadarPlanReport(radar domain.IdeaRadarReport) domain.IdeaRadarPlanReport {
	return domain.IdeaRadarPlanReport{
		SchemaVersion: domain.IdeaRadarPlanReportSchemaVersion,
		Status:        domain.IdeaRadarPlanReportUnknown,
		Plans:         make([]domain.IdeaRadarTargetPlan, 0),
		Metrics: domain.IdeaRadarPlanMetrics{
			Targets:         radar.Metrics.Targets,
			RadarOperations: radar.Metrics.Operations,
			Operations:      radar.Metrics.Operations,
			LimitExceeded:   radar.Metrics.LimitExceeded,
		},
	}
}

func unknownIdeaRadarPlanReport(report domain.IdeaRadarPlanReport) domain.IdeaRadarPlanReport {
	report.Status = domain.IdeaRadarPlanReportUnknown
	report.Reason = "budget_exceeded"
	report.Plans = make([]domain.IdeaRadarTargetPlan, 0)
	report.Metrics.Plans = 0
	report.Metrics.Probes = 0
	report.Metrics.LimitExceeded = true
	return report
}

func countIdeaRadarPlanningOperation(report *domain.IdeaRadarPlanReport, count int) bool {
	report.Metrics.PlanningOperations += count
	report.Metrics.Operations = report.Metrics.RadarOperations + report.Metrics.PlanningOperations
	return report.Metrics.Operations <= ideaRadarPlanMaxOperations
}

func planIdeaRadarTarget(
	assessment domain.IdeaRadarAssessment,
	observations map[string]ideaRadarObservation,
	projections map[string]ideaRadarPlanProjection,
) []domain.IdeaRadarProbe {
	frontier := make([]domain.IdeaRadarProbe, 0)
	switch assessment.Outcome {
	case domain.IdeaRadarConflict:
		for _, observationID := range assessment.ConflictingObservationIDs {
			frontier = append(frontier, newIdeaRadarProbe(assessment.TargetID, domain.IdeaRadarProbeResolveConflict,
				domain.IdeaRadarProbePriorityResolveConflict, []string{observationID}, "", domain.IdeaRadarOriginRequirementNone,
				"resolve_conflicting_observation"))
		}
	case domain.IdeaRadarAmbiguous:
		for _, observationID := range assessment.AmbiguousObservationIDs {
			frontier = append(frontier, newIdeaRadarProbe(assessment.TargetID, domain.IdeaRadarProbeResolveAmbiguity,
				domain.IdeaRadarProbePriorityResolveAmbiguity, []string{observationID}, "", domain.IdeaRadarOriginRequirementNone,
				"resolve_ambiguous_observation"))
		}
	case domain.IdeaRadarInsufficientPerspectives:
		present := ideaRadarPlanKindSet(assessment.ProjectionKinds)
		missing := 0
		for _, kind := range ideaRadarPlanKindOrder {
			if _, exists := present[kind]; exists || missing == assessment.MissingProjectionKinds {
				continue
			}
			frontier = append(frontier, newIdeaRadarProbe(assessment.TargetID, domain.IdeaRadarProbeAddPerspective,
				domain.IdeaRadarProbePriorityAddPerspective, nil, kind, domain.IdeaRadarOriginRequirementDistinctRecorded,
				"add_missing_projection_with_distinct_origin"))
			missing++
		}
		frontier = appendIndependentOriginProbes(frontier, assessment.TargetID,
			assessment.MissingOrigins-assessment.MissingProjectionKinds, assessment.ProjectionKinds)
	case domain.IdeaRadarUnknownOrigin:
		recordedKinds := make(map[domain.IdeaProjectionKind]struct{})
		unknownByKind := make(map[domain.IdeaProjectionKind][]string)
		for _, observationID := range assessment.SupportingObservationIDs {
			observation := observations[observationID]
			projection := projections[observation.projectionID]
			if projection.originStatus == domain.EvidenceOriginRecorded {
				recordedKinds[projection.kind] = struct{}{}
			} else {
				unknownByKind[projection.kind] = append(unknownByKind[projection.kind], observationID)
			}
		}
		neededKinds := ideaRadarMinimumKinds - len(recordedKinds)
		for _, kind := range ideaRadarPlanKindOrder {
			ids := unknownByKind[kind]
			if neededKinds == 0 || len(ids) == 0 {
				continue
			}
			sort.Strings(ids)
			frontier = append(frontier, newIdeaRadarProbe(assessment.TargetID, domain.IdeaRadarProbeRecordOrigin,
				domain.IdeaRadarProbePriorityRecordOrigin, []string{ids[0]}, kind,
				domain.IdeaRadarOriginRequirementDistinctRecorded, "record_distinct_origin"))
			neededKinds--
		}
		recordProbes := 0
		for _, probe := range frontier {
			if probe.Kind == domain.IdeaRadarProbeRecordOrigin {
				recordProbes++
			}
		}
		frontier = appendIndependentOriginProbes(frontier, assessment.TargetID,
			assessment.MissingOrigins-recordProbes, assessment.ProjectionKinds)
	case domain.IdeaRadarInsufficientOrigins:
		frontier = appendIndependentOriginProbes(frontier, assessment.TargetID,
			assessment.MissingOrigins, assessment.ProjectionKinds)
	case domain.IdeaRadarConvergent:
		probe := newIdeaRadarProbe(assessment.TargetID, domain.IdeaRadarProbeBindEvidence,
			domain.IdeaRadarProbePriorityBindEvidence, assessment.SupportingObservationIDs, "",
			domain.IdeaRadarOriginRequirementRecorded, "bind_physical_evidence")
		probe.RequiredRoles = make([]domain.EvidenceCoalitionRole, len(assessment.CoalitionDraft.RequiredRoles))
		copy(probe.RequiredRoles, assessment.CoalitionDraft.RequiredRoles)
		frontier = append(frontier, probe)
	}
	return frontier
}

func newIdeaRadarProbe(
	targetID string,
	kind domain.IdeaRadarProbeKind,
	priority int,
	observationIDs []string,
	projectionKind domain.IdeaProjectionKind,
	originRequirement domain.IdeaRadarOriginRequirement,
	reason string,
) domain.IdeaRadarProbe {
	ids := make([]string, len(observationIDs))
	copy(ids, observationIDs)
	return domain.IdeaRadarProbe{
		Kind:              kind,
		Priority:          priority,
		TargetID:          targetID,
		ObservationIDs:    ids,
		ProjectionKind:    projectionKind,
		OriginRequirement: originRequirement,
		RequiredRoles:     make([]domain.EvidenceCoalitionRole, 0),
		Reason:            reason,
	}
}

func appendIndependentOriginProbes(
	frontier []domain.IdeaRadarProbe,
	targetID string,
	count int,
	presentKinds []domain.IdeaProjectionKind,
) []domain.IdeaRadarProbe {
	if count <= 0 {
		return frontier
	}
	present := ideaRadarPlanKindSet(presentKinds)
	kinds := make([]domain.IdeaProjectionKind, 0, len(present))
	for _, kind := range ideaRadarPlanKindOrder {
		if _, exists := present[kind]; exists {
			kinds = append(kinds, kind)
		}
	}
	for index := 0; index < count; index++ {
		frontier = append(frontier, newIdeaRadarProbe(targetID, domain.IdeaRadarProbeAddIndependentOrigin,
			domain.IdeaRadarProbePriorityAddIndependentOrigin, nil, kinds[index%len(kinds)],
			domain.IdeaRadarOriginRequirementDistinctRecorded, "add_independent_origin"))
	}
	return frontier
}

func ideaRadarPlanKindSet(kinds []domain.IdeaProjectionKind) map[domain.IdeaProjectionKind]struct{} {
	result := make(map[domain.IdeaProjectionKind]struct{}, len(kinds))
	for _, kind := range kinds {
		result[kind] = struct{}{}
	}
	return result
}

func sortIdeaRadarProbes(probes []domain.IdeaRadarProbe) {
	sort.Slice(probes, func(i, j int) bool {
		left, right := probes[i], probes[j]
		if left.Priority != right.Priority {
			return left.Priority < right.Priority
		}
		if left.Kind != right.Kind {
			return left.Kind < right.Kind
		}
		leftKind, rightKind := ideaRadarPlanKindRank(left.ProjectionKind), ideaRadarPlanKindRank(right.ProjectionKind)
		if leftKind != rightKind {
			return leftKind < rightKind
		}
		leftObservation, rightObservation := "", ""
		if len(left.ObservationIDs) > 0 {
			leftObservation = left.ObservationIDs[0]
		}
		if len(right.ObservationIDs) > 0 {
			rightObservation = right.ObservationIDs[0]
		}
		if leftObservation != rightObservation {
			return leftObservation < rightObservation
		}
		return left.Reason < right.Reason
	})
}

func ideaRadarPlanKindRank(kind domain.IdeaProjectionKind) int {
	for index, candidate := range ideaRadarPlanKindOrder {
		if kind == candidate {
			return index
		}
	}
	return -1
}
