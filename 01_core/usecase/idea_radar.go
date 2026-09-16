/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_radar.md
 * @layer L1
 * @updated 2026-09-11
 */

package usecase

import (
	"sort"

	"bateia/01_core/domain"
)

const (
	ideaRadarMaxProjections  = 64
	ideaRadarMaxObservations = 1024
	ideaRadarMaxTargets      = 512
	ideaRadarMaxIDBytes      = 128
	ideaRadarMaxTextItems    = 64
	ideaRadarMaxTextBytes    = 4096
	ideaRadarMaxOperations   = 4096
	ideaRadarMinimumKinds    = 3
	ideaRadarMinimumOrigins  = 3
	ideaRadarDraftReason     = "requires_evidence_binding"
)

type ideaRadarProjection struct {
	kind         domain.IdeaProjectionKind
	originStatus domain.EvidenceOriginStatus
	originID     string
}

type ideaRadarObservation struct {
	id           string
	projectionID string
	candidateID  string
	conceptID    string
	stance       domain.IdeaRadarStance
}

type ideaRadarTarget struct {
	candidateID string
	conceptID   string
	items       []ideaRadarObservation
}

// AnalyzeIdeaRadar classifies explicitly declared conceptual observations without I/O
// or changes to any existing decision path.
func AnalyzeIdeaRadar(input domain.IdeaRadarInput) domain.IdeaRadarReport {
	report := newIdeaRadarReport(len(input.Projections), len(input.Observations))
	if input.SchemaVersion != domain.IdeaRadarInputSchemaVersion {
		return unknownIdeaRadarReport(report, "invalid_schema", false)
	}
	if len(input.Projections) > ideaRadarMaxProjections || len(input.Observations) > ideaRadarMaxObservations {
		return unknownIdeaRadarReport(report, "budget_exceeded", true)
	}
	if projectionTextLimitExceeded(input.Projections) {
		return unknownIdeaRadarReport(report, "budget_exceeded", true)
	}
	if ideaRadarIDLimitExceeded(input) {
		return unknownIdeaRadarReport(report, "budget_exceeded", true)
	}
	if len(input.Projections) == 0 {
		return unknownIdeaRadarReport(report, "no_projections", false)
	}
	if len(input.Observations) == 0 {
		return unknownIdeaRadarReport(report, "no_observations", false)
	}

	projections := make(map[string]ideaRadarProjection, len(input.Projections))
	for _, projection := range input.Projections {
		id, ok := canonicalIdeaRadarID(projection.ID)
		if !ok || !validIdeaProjectionKind(projection.Kind) || !validIdeaRadarOrigin(projection.OriginStatus, projection.OriginID) {
			return unknownIdeaRadarReport(report, "invalid_input", false)
		}
		if _, exists := projections[id]; exists {
			return unknownIdeaRadarReport(report, "ambiguous_input", false)
		}
		projections[id] = ideaRadarProjection{
			kind:         projection.Kind,
			originStatus: projection.OriginStatus,
			originID:     projection.OriginID,
		}
		report.Metrics.Operations++
	}

	byObservationID := make(map[string]ideaRadarObservation, len(input.Observations))
	canonical := make(map[string]ideaRadarObservation, len(input.Observations))
	for _, observation := range input.Observations {
		normalized, ok := normalizeIdeaRadarObservation(observation)
		if !ok || !validIdeaRadarStance(normalized.stance) {
			return unknownIdeaRadarReport(report, "invalid_input", false)
		}
		if _, exists := projections[normalized.projectionID]; !exists {
			return unknownIdeaRadarReport(report, "invalid_input", false)
		}
		report.Metrics.Operations++

		if prior, exists := byObservationID[normalized.id]; exists {
			if !sameIdeaRadarObservation(prior, normalized) {
				return unknownIdeaRadarReport(report, "ambiguous_input", false)
			}
		} else {
			byObservationID[normalized.id] = normalized
		}

		key := ideaRadarSemanticKey(normalized)
		if prior, exists := canonical[key]; !exists || normalized.id < prior.id {
			canonical[key] = normalized
		}
	}

	report.Metrics.CanonicalObservations = len(canonical)
	report.Metrics.DuplicateObservations = len(input.Observations) - len(canonical)
	report.Metrics.Operations += len(canonical)
	if report.Metrics.Operations > ideaRadarMaxOperations {
		return unknownIdeaRadarReport(report, "budget_exceeded", true)
	}

	targets := make(map[string]*ideaRadarTarget, len(canonical))
	for _, observation := range canonical {
		targetID := observation.candidateID + "#" + observation.conceptID
		target := targets[targetID]
		if target == nil {
			if len(targets) == ideaRadarMaxTargets {
				return unknownIdeaRadarReport(report, "budget_exceeded", true)
			}
			target = &ideaRadarTarget{candidateID: observation.candidateID, conceptID: observation.conceptID}
			targets[targetID] = target
		}
		target.items = append(target.items, observation)
	}

	targetIDs := make([]string, 0, len(targets))
	for targetID := range targets {
		targetIDs = append(targetIDs, targetID)
	}
	sort.Strings(targetIDs)
	report.Metrics.Targets = len(targetIDs)
	report.Assessments = make([]domain.IdeaRadarAssessment, 0, len(targetIDs))
	for _, targetID := range targetIDs {
		assessment := assessIdeaRadarTarget(targetID, targets[targetID], projections)
		report.Assessments = append(report.Assessments, assessment)
		countIdeaRadarOutcome(&report.Metrics, assessment.Outcome)
	}
	report.Status = domain.IdeaRadarReportComplete
	report.Reason = "analyzed"
	return report
}

func newIdeaRadarReport(projections, observations int) domain.IdeaRadarReport {
	return domain.IdeaRadarReport{
		SchemaVersion: domain.IdeaRadarReportSchemaVersion,
		Status:        domain.IdeaRadarReportUnknown,
		Assessments:   make([]domain.IdeaRadarAssessment, 0),
		Metrics: domain.IdeaRadarMetrics{
			Projections:  projections,
			Observations: observations,
		},
	}
}

func unknownIdeaRadarReport(report domain.IdeaRadarReport, reason string, limitExceeded bool) domain.IdeaRadarReport {
	report.Status = domain.IdeaRadarReportUnknown
	report.Reason = reason
	report.Assessments = make([]domain.IdeaRadarAssessment, 0)
	report.Metrics.LimitExceeded = limitExceeded
	return report
}

func projectionTextLimitExceeded(projections []domain.IdeaRadarProjection) bool {
	for _, projection := range projections {
		for _, dimension := range [][]string{projection.Terms, projection.Actions, projection.Objects, projection.Required} {
			if len(dimension) > ideaRadarMaxTextItems {
				return true
			}
			bytes := 0
			for _, item := range dimension {
				bytes += len(item)
				if bytes > ideaRadarMaxTextBytes {
					return true
				}
			}
		}
	}
	return false
}

func ideaRadarIDLimitExceeded(input domain.IdeaRadarInput) bool {
	for _, projection := range input.Projections {
		if trimmedASCIIBytes(projection.ID) > ideaRadarMaxIDBytes || len(projection.OriginID) > ideaRadarMaxIDBytes {
			return true
		}
	}
	for _, observation := range input.Observations {
		if trimmedASCIIBytes(observation.ID) > ideaRadarMaxIDBytes ||
			trimmedASCIIBytes(observation.ProjectionID) > ideaRadarMaxIDBytes ||
			trimmedASCIIBytes(observation.CandidateID) > ideaRadarMaxIDBytes ||
			trimmedASCIIBytes(observation.ConceptID) > ideaRadarMaxIDBytes {
			return true
		}
	}
	return false
}

func trimmedASCIIBytes(value string) int {
	start, end := 0, len(value)
	for start < end && isASCIIWhitespace(value[start]) {
		start++
	}
	for end > start && isASCIIWhitespace(value[end-1]) {
		end--
	}
	return end - start
}

func canonicalIdeaRadarID(value string) (string, bool) {
	start, end := 0, len(value)
	for start < end && isASCIIWhitespace(value[start]) {
		start++
	}
	for end > start && isASCIIWhitespace(value[end-1]) {
		end--
	}
	if end-start == 0 || end-start > ideaRadarMaxIDBytes {
		return "", false
	}
	canonical := make([]byte, end-start)
	for i, char := range []byte(value[start:end]) {
		if char >= 'A' && char <= 'Z' {
			char += 'a' - 'A'
		}
		if !isIdeaRadarIDChar(char) {
			return "", false
		}
		canonical[i] = char
	}
	return string(canonical), true
}

func isASCIIWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r' || char == '\v' || char == '\f'
}

func isIdeaRadarIDChar(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= '0' && char <= '9' ||
		char == '.' || char == '_' || char == ':' || char == '/' || char == '-'
}

func validIdeaProjectionKind(kind domain.IdeaProjectionKind) bool {
	switch kind {
	case domain.IdeaProjectionProblem, domain.IdeaProjectionBehavior, domain.IdeaProjectionArchitecture,
		domain.IdeaProjectionConstraint, domain.IdeaProjectionOperation:
		return true
	default:
		return false
	}
}

func validIdeaRadarOrigin(status domain.EvidenceOriginStatus, originID string) bool {
	switch status {
	case domain.EvidenceOriginRecorded:
		if len(originID) == 0 || len(originID) > ideaRadarMaxIDBytes ||
			isASCIIWhitespace(originID[0]) || isASCIIWhitespace(originID[len(originID)-1]) {
			return false
		}
		for i := 0; i < len(originID); i++ {
			if originID[i] > 0x7f {
				return false
			}
		}
		return true
	case domain.EvidenceOriginUnknown:
		return originID == ""
	default:
		return false
	}
}

func normalizeIdeaRadarObservation(observation domain.IdeaRadarObservation) (ideaRadarObservation, bool) {
	id, idOK := canonicalIdeaRadarID(observation.ID)
	projectionID, projectionOK := canonicalIdeaRadarID(observation.ProjectionID)
	candidateID, candidateOK := canonicalIdeaRadarID(observation.CandidateID)
	conceptID, conceptOK := canonicalIdeaRadarID(observation.ConceptID)
	return ideaRadarObservation{
		id:           id,
		projectionID: projectionID,
		candidateID:  candidateID,
		conceptID:    conceptID,
		stance:       observation.Stance,
	}, idOK && projectionOK && candidateOK && conceptOK
}

func validIdeaRadarStance(stance domain.IdeaRadarStance) bool {
	return stance == domain.IdeaRadarStanceSupport || stance == domain.IdeaRadarStanceConflict || stance == domain.IdeaRadarStanceAmbiguous
}

func sameIdeaRadarObservation(left, right ideaRadarObservation) bool {
	return left.projectionID == right.projectionID && left.candidateID == right.candidateID &&
		left.conceptID == right.conceptID && left.stance == right.stance
}

func ideaRadarSemanticKey(observation ideaRadarObservation) string {
	return observation.projectionID + "#" + observation.candidateID + "#" + observation.conceptID + "#" + string(observation.stance)
}

func assessIdeaRadarTarget(targetID string, target *ideaRadarTarget, projections map[string]ideaRadarProjection) domain.IdeaRadarAssessment {
	supports := make([]string, 0)
	conflicts := make([]string, 0)
	ambiguous := make([]string, 0)
	unknownOrigins := make([]string, 0)
	kinds := make(map[domain.IdeaProjectionKind]struct{})
	recordedKinds := make(map[domain.IdeaProjectionKind]struct{})
	origins := make(map[string]struct{})

	for _, observation := range target.items {
		projection := projections[observation.projectionID]
		switch observation.stance {
		case domain.IdeaRadarStanceSupport:
			supports = append(supports, observation.id)
			kinds[projection.kind] = struct{}{}
			if projection.originStatus == domain.EvidenceOriginRecorded {
				recordedKinds[projection.kind] = struct{}{}
				origins[projection.originID] = struct{}{}
			} else {
				unknownOrigins = append(unknownOrigins, observation.id)
			}
		case domain.IdeaRadarStanceConflict:
			conflicts = append(conflicts, observation.id)
		case domain.IdeaRadarStanceAmbiguous:
			ambiguous = append(ambiguous, observation.id)
		}
	}

	sort.Strings(supports)
	sort.Strings(conflicts)
	sort.Strings(ambiguous)
	sort.Strings(unknownOrigins)
	projectionKinds := sortedIdeaProjectionKinds(kinds)
	distinctOrigins := make([]string, 0, len(origins))
	for origin := range origins {
		distinctOrigins = append(distinctOrigins, origin)
	}
	sort.Strings(distinctOrigins)

	missingKinds := ideaRadarMinimumKinds - len(kinds)
	if missingKinds < 0 {
		missingKinds = 0
	}
	missingOrigins := ideaRadarMinimumOrigins - len(origins)
	if missingOrigins < 0 {
		missingOrigins = 0
	}
	outcome, reason := classifyIdeaRadarTarget(len(conflicts), len(ambiguous), len(kinds), len(recordedKinds), len(origins))
	draftObservationIDs := make([]string, len(supports))
	copy(draftObservationIDs, supports)

	return domain.IdeaRadarAssessment{
		TargetID:                    targetID,
		CandidateID:                 target.candidateID,
		ConceptID:                   target.conceptID,
		Outcome:                     outcome,
		SupportingObservationIDs:    supports,
		ConflictingObservationIDs:   conflicts,
		AmbiguousObservationIDs:     ambiguous,
		UnknownOriginObservationIDs: unknownOrigins,
		ProjectionKinds:             projectionKinds,
		DistinctOrigins:             distinctOrigins,
		MissingProjectionKinds:      missingKinds,
		MissingOrigins:              missingOrigins,
		CoalitionDraft: domain.IdeaRadarCoalitionDraft{
			ClaimID:        targetID,
			ObservationIDs: draftObservationIDs,
			RequiredRoles: []domain.EvidenceCoalitionRole{
				domain.EvidenceCoalitionRoleImplementation,
				domain.EvidenceCoalitionRoleIndependentTest,
			},
			ReadyFor0065: false,
			Reason:       ideaRadarDraftReason,
		},
		Reason: reason,
	}
}

func sortedIdeaProjectionKinds(kinds map[domain.IdeaProjectionKind]struct{}) []domain.IdeaProjectionKind {
	result := make([]domain.IdeaProjectionKind, 0, len(kinds))
	for kind := range kinds {
		result = append(result, kind)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func classifyIdeaRadarTarget(conflicts, ambiguous, kinds, recordedKinds, origins int) (domain.IdeaRadarOutcome, string) {
	switch {
	case conflicts > 0:
		return domain.IdeaRadarConflict, "conflicting_projection"
	case ambiguous > 0:
		return domain.IdeaRadarAmbiguous, "ambiguous_projection"
	case kinds < ideaRadarMinimumKinds:
		return domain.IdeaRadarInsufficientPerspectives, "projection_diversity_below_minimum"
	case recordedKinds < ideaRadarMinimumKinds:
		return domain.IdeaRadarUnknownOrigin, "origin_unrecorded"
	case origins < ideaRadarMinimumOrigins:
		return domain.IdeaRadarInsufficientOrigins, "distinct_origins_below_minimum"
	default:
		return domain.IdeaRadarConvergent, "idea_convergent"
	}
}

func countIdeaRadarOutcome(metrics *domain.IdeaRadarMetrics, outcome domain.IdeaRadarOutcome) {
	switch outcome {
	case domain.IdeaRadarConvergent:
		metrics.Convergent++
	case domain.IdeaRadarConflict:
		metrics.Conflicts++
	case domain.IdeaRadarAmbiguous:
		metrics.Ambiguous++
	case domain.IdeaRadarInsufficientPerspectives:
		metrics.InsufficientPerspectives++
	case domain.IdeaRadarUnknownOrigin:
		metrics.UnknownOrigins++
	case domain.IdeaRadarInsufficientOrigins:
		metrics.InsufficientOrigins++
	}
}
