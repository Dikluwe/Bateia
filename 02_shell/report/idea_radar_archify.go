/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/idea_radar_archify.md
 * @layer L2
 * @updated 2026-09-12
 */

package report

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

const (
	ideaRadarArchifyMaxOrigins      = 5
	ideaRadarArchifyMaxPerspectives = 5
	ideaRadarArchifyMaxTargets      = 5
	ideaRadarArchifyMaxNodes        = 15
	ideaRadarArchifyMaxFlows        = 1029
	ideaRadarArchifyMaxCards        = 5
	ideaRadarArchifyMaxProbes       = 1024
	ideaRadarArchifyMaxBytes        = 4 * 1024 * 1024
)

type ideaRadarArchifyDocument struct {
	SchemaVersion int                     `json:"schema_version"`
	DiagramType   string                  `json:"diagram_type"`
	Meta          ideaRadarArchifyMeta    `json:"meta"`
	Stages        []ideaRadarArchifyStage `json:"stages"`
	Nodes         []ideaRadarArchifyNode  `json:"nodes"`
	Flows         []ideaRadarArchifyFlow  `json:"flows"`
	Cards         []ideaRadarArchifyCard  `json:"cards"`
}

type ideaRadarArchifyMeta struct {
	Title   string `json:"title"`
	ViewBox [2]int `json:"viewBox"`
}

type ideaRadarArchifyStage struct {
	Label string `json:"label"`
}

type ideaRadarArchifyNode struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Stage int    `json:"stage"`
	Row   int    `json:"row"`

	identity string
}

type ideaRadarArchifyFlow struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Label   string `json:"label"`
	Variant string `json:"variant"`
	Route   string `json:"route,omitempty"`
	LabelAt []int  `json:"labelAt,omitempty"`

	identity string
}

type ideaRadarArchifyCard struct {
	Dot   string   `json:"dot"`
	Title string   `json:"title"`
	Items []string `json:"items"`

	identity string
}

type canonicalRadarProjection struct {
	id           string
	kind         domain.IdeaProjectionKind
	originStatus domain.EvidenceOriginStatus
	originID     string
}

type canonicalRadarObservation struct {
	id           string
	projectionID string
	candidateID  string
	conceptID    string
	stance       domain.IdeaRadarStance
}

// ExportIdeaRadarArchify recalculates the radar and probe frontier and projects
// only their declared topology into an Archify Dataflow v1 document.
func ExportIdeaRadarArchify(input domain.IdeaRadarInput) ([]byte, error) {
	radar := usecase.AnalyzeIdeaRadar(input)
	if radar.Status != domain.IdeaRadarReportComplete {
		return nil, fmt.Errorf("idea radar UNKNOWN: %s", radar.Reason)
	}
	plan := usecase.PlanIdeaRadarProbes(input)
	if plan.Status != domain.IdeaRadarPlanReportComplete {
		return nil, fmt.Errorf("idea radar plan UNKNOWN: %s", plan.Reason)
	}
	if err := validateIdeaRadarReports(radar, plan); err != nil {
		return nil, err
	}

	projections := canonicalIdeaRadarProjections(input.Projections)
	observations := canonicalIdeaRadarObservations(input.Observations)
	if err := validateIdeaRadarRenderBudget(projections, radar.Assessments); err != nil {
		return nil, err
	}
	document, err := buildIdeaRadarArchifyDocument(projections, observations, radar, plan)
	if err != nil {
		return nil, err
	}
	encoded, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("serialize Archify Dataflow: %w", err)
	}
	encoded = append(encoded, '\n')
	if len(encoded) > ideaRadarArchifyMaxBytes {
		return nil, fmt.Errorf("render_budget_exceeded: serialized bytes %d > %d", len(encoded), ideaRadarArchifyMaxBytes)
	}
	return encoded, nil
}

func validateIdeaRadarRenderBudget(projections []canonicalRadarProjection, assessments []domain.IdeaRadarAssessment) error {
	if len(projections) > ideaRadarArchifyMaxPerspectives {
		return fmt.Errorf("render_budget_exceeded: perspectives %d > %d", len(projections), ideaRadarArchifyMaxPerspectives)
	}
	if len(assessments) > ideaRadarArchifyMaxTargets {
		return fmt.Errorf("render_budget_exceeded: targets %d > %d", len(assessments), ideaRadarArchifyMaxTargets)
	}
	origins := make(map[string]struct{}, len(projections))
	for _, projection := range projections {
		identity := "recorded\x00" + projection.originID
		if projection.originStatus == domain.EvidenceOriginUnknown {
			identity = "unknown\x00" + projection.id
		}
		origins[identity] = struct{}{}
	}
	if len(origins) > ideaRadarArchifyMaxOrigins {
		return fmt.Errorf("render_budget_exceeded: unique origins %d > %d", len(origins), ideaRadarArchifyMaxOrigins)
	}
	return nil
}

func validateIdeaRadarReports(radar domain.IdeaRadarReport, plan domain.IdeaRadarPlanReport) error {
	if len(radar.Assessments) != len(plan.Plans) {
		return fmt.Errorf("inconsistent idea radar reports: %d assessments and %d plans", len(radar.Assessments), len(plan.Plans))
	}
	assessments := make(map[string]domain.IdeaRadarOutcome, len(radar.Assessments))
	for _, assessment := range radar.Assessments {
		if _, exists := assessments[assessment.TargetID]; exists {
			return fmt.Errorf("inconsistent idea radar reports: duplicate assessment target %q", assessment.TargetID)
		}
		assessments[assessment.TargetID] = assessment.Outcome
	}
	plans := make(map[string]struct{}, len(plan.Plans))
	probeCount := 0
	for _, targetPlan := range plan.Plans {
		if _, exists := plans[targetPlan.TargetID]; exists {
			return fmt.Errorf("inconsistent idea radar reports: duplicate plan target %q", targetPlan.TargetID)
		}
		plans[targetPlan.TargetID] = struct{}{}
		outcome, exists := assessments[targetPlan.TargetID]
		if !exists || outcome != targetPlan.RadarOutcome {
			return fmt.Errorf("inconsistent idea radar reports for target %q", targetPlan.TargetID)
		}
		probeCount += len(targetPlan.Frontier)
	}
	if probeCount > ideaRadarArchifyMaxProbes {
		return fmt.Errorf("render_budget_exceeded: probes %d > %d", probeCount, ideaRadarArchifyMaxProbes)
	}
	return nil
}

func buildIdeaRadarArchifyDocument(
	projections []canonicalRadarProjection,
	observations []canonicalRadarObservation,
	radar domain.IdeaRadarReport,
	plan domain.IdeaRadarPlanReport,
) (ideaRadarArchifyDocument, error) {
	document := ideaRadarArchifyDocument{
		SchemaVersion: 1,
		DiagramType:   "dataflow",
		Meta:          ideaRadarArchifyMeta{Title: "Bateia Idea Radar", ViewBox: [2]int{1400, 850}},
		Stages: []ideaRadarArchifyStage{
			{Label: "Origins"},
			{Label: "Perspectives"},
			{Label: "Targets"},
		},
		Nodes: make([]ideaRadarArchifyNode, 0, len(projections)*2+len(radar.Assessments)),
		Flows: make([]ideaRadarArchifyFlow, 0, len(projections)+len(observations)),
		Cards: make([]ideaRadarArchifyCard, 0, len(radar.Assessments)),
	}

	originNodes := make(map[string]ideaRadarArchifyNode, len(projections))
	projectionRows := make(map[string]int, len(projections))
	for projectionRow, projection := range projections {
		projectionRows[projection.id] = projectionRow
		perspectiveID := ideaRadarArchifyID("perspective_", projection.id)
		originID, originLabel, originIdentity, projectionFlowLabel := "", "", "", ""
		if projection.originStatus == domain.EvidenceOriginRecorded {
			originID = ideaRadarArchifyID("origin_", projection.originID)
			originLabel = "O " + shortIdeaRadarArchifyLabel(projection.originID)
			originIdentity = originID
			projectionFlowLabel = "projection " + projection.id + " · origin " + projection.originID
		} else {
			originID = ideaRadarArchifyID("origin_unknown_", projection.id)
			originLabel = "O unknown"
			originIdentity = originID
			projectionFlowLabel = "projection " + projection.id + " · origin unknown"
		}
		originNodes[originID] = ideaRadarArchifyNode{
			ID: originID, Type: "external", Label: originLabel,
			Stage: 0, identity: originIdentity,
		}
		document.Nodes = append(document.Nodes, ideaRadarArchifyNode{
			ID: perspectiveID, Type: "external", Label: "P " + shortIdeaRadarProjectionKind(projection.kind),
			Stage: 1, identity: projection.id,
		})
		document.Flows = append(document.Flows, ideaRadarArchifyFlow{
			ID: ideaRadarArchifyID("flow_projection_", projection.id), From: originID, To: perspectiveID,
			Label: projectionFlowLabel, Variant: "default",
			LabelAt: ideaRadarArchifyLabelPoint(208, projectionRow), identity: "0\x00" + projection.id,
		})
	}
	for _, node := range originNodes {
		document.Nodes = append(document.Nodes, node)
	}

	for _, assessment := range radar.Assessments {
		document.Nodes = append(document.Nodes, ideaRadarArchifyNode{
			ID: ideaRadarArchifyID("target_", assessment.TargetID), Type: "external",
			Label: "T " + shortIdeaRadarArchifyLabel(assessment.TargetID), Stage: 2,
			identity: assessment.TargetID,
		})
	}
	for _, observation := range observations {
		targetIdentity := observation.candidateID + "#" + observation.conceptID
		document.Flows = append(document.Flows, ideaRadarArchifyFlow{
			ID:       ideaRadarArchifyID("flow_observation_", observation.id),
			From:     ideaRadarArchifyID("perspective_", observation.projectionID),
			To:       ideaRadarArchifyID("target_", targetIdentity),
			Label:    "observation " + observation.id + " · " + string(observation.stance),
			Variant:  ideaRadarArchifyStanceVariant(observation.stance),
			LabelAt:  ideaRadarArchifyObservationLabelPoint(projectionRows[observation.projectionID]),
			identity: "1\x00" + observation.id,
		})
	}

	planByTarget := make(map[string]domain.IdeaRadarTargetPlan, len(plan.Plans))
	for _, targetPlan := range plan.Plans {
		planByTarget[targetPlan.TargetID] = targetPlan
	}
	for _, assessment := range radar.Assessments {
		items := []string{
			"Outcome: " + string(assessment.Outcome),
			"Reason: " + assessment.Reason,
			formatIdeaRadarProvenance(assessment),
		}
		for _, probe := range planByTarget[assessment.TargetID].Frontier {
			items = append(items, formatIdeaRadarProbe(probe))
		}
		document.Cards = append(document.Cards, ideaRadarArchifyCard{
			Dot: ideaRadarArchifyOutcomeDot(assessment.Outcome), Title: assessment.TargetID,
			Items: items, identity: assessment.TargetID,
		})
	}

	sort.Slice(document.Nodes, func(i, j int) bool {
		if document.Nodes[i].Stage != document.Nodes[j].Stage {
			return document.Nodes[i].Stage < document.Nodes[j].Stage
		}
		return document.Nodes[i].identity < document.Nodes[j].identity
	})
	rows := [3]int{}
	for index := range document.Nodes {
		stage := document.Nodes[index].Stage
		document.Nodes[index].Row = rows[stage]
		rows[stage]++
	}
	sort.Slice(document.Flows, func(i, j int) bool { return document.Flows[i].identity < document.Flows[j].identity })
	sort.Slice(document.Cards, func(i, j int) bool { return document.Cards[i].identity < document.Cards[j].identity })

	if len(document.Nodes) > ideaRadarArchifyMaxNodes {
		return ideaRadarArchifyDocument{}, fmt.Errorf("render_budget_exceeded: nodes %d > %d", len(document.Nodes), ideaRadarArchifyMaxNodes)
	}
	if len(document.Flows) > ideaRadarArchifyMaxFlows {
		return ideaRadarArchifyDocument{}, fmt.Errorf("render_budget_exceeded: flows %d > %d", len(document.Flows), ideaRadarArchifyMaxFlows)
	}
	if len(document.Cards) > ideaRadarArchifyMaxCards {
		return ideaRadarArchifyDocument{}, fmt.Errorf("render_budget_exceeded: cards %d > %d", len(document.Cards), ideaRadarArchifyMaxCards)
	}
	return document, nil
}

func canonicalIdeaRadarProjections(input []domain.IdeaRadarProjection) []canonicalRadarProjection {
	result := make([]canonicalRadarProjection, 0, len(input))
	for _, projection := range input {
		result = append(result, canonicalRadarProjection{
			id: canonicalIdeaRadarArchifyID(projection.ID), kind: projection.Kind,
			originStatus: projection.OriginStatus, originID: projection.OriginID,
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].id < result[j].id })
	return result
}

func canonicalIdeaRadarObservations(input []domain.IdeaRadarObservation) []canonicalRadarObservation {
	bySemanticKey := make(map[string]canonicalRadarObservation, len(input))
	for _, observation := range input {
		canonical := canonicalRadarObservation{
			id:           canonicalIdeaRadarArchifyID(observation.ID),
			projectionID: canonicalIdeaRadarArchifyID(observation.ProjectionID),
			candidateID:  canonicalIdeaRadarArchifyID(observation.CandidateID),
			conceptID:    canonicalIdeaRadarArchifyID(observation.ConceptID),
			stance:       observation.Stance,
		}
		key := canonical.projectionID + "\x00" + canonical.candidateID + "\x00" + canonical.conceptID + "\x00" + string(canonical.stance)
		prior, exists := bySemanticKey[key]
		if !exists || canonical.id < prior.id {
			bySemanticKey[key] = canonical
		}
	}
	result := make([]canonicalRadarObservation, 0, len(bySemanticKey))
	for _, observation := range bySemanticKey {
		result = append(result, observation)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].id < result[j].id })
	return result
}

func canonicalIdeaRadarArchifyID(value string) string {
	value = strings.Trim(value, " \t\n\r\v\f")
	bytes := []byte(value)
	for index, char := range bytes {
		if char >= 'A' && char <= 'Z' {
			bytes[index] = char + ('a' - 'A')
		}
	}
	return string(bytes)
}

func ideaRadarArchifyID(prefix, identity string) string {
	return prefix + hex.EncodeToString([]byte(identity))
}

func ideaRadarArchifyLabelPoint(x, row int) []int {
	return []int{x, 200 + row*114}
}

func ideaRadarArchifyObservationLabelPoint(row int) []int {
	return []int{423, 220 + row*114}
}

func shortIdeaRadarArchifyLabel(value string) string {
	short := []byte(value)
	for index, char := range short {
		if !isIdeaRadarArchifyASCIIAlphanumeric(char) {
			short[index] = '_'
		}
	}
	if len(short) <= 8 {
		return string(short)
	}
	return string(short[:7]) + "~"
}

func isIdeaRadarArchifyASCIIAlphanumeric(char byte) bool {
	return char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' || char >= '0' && char <= '9'
}

func shortIdeaRadarProjectionKind(kind domain.IdeaProjectionKind) string {
	switch kind {
	case domain.IdeaProjectionProblem:
		return "PROB"
	case domain.IdeaProjectionBehavior:
		return "BEHV"
	case domain.IdeaProjectionArchitecture:
		return "ARCH"
	case domain.IdeaProjectionConstraint:
		return "CONS"
	case domain.IdeaProjectionOperation:
		return "OPER"
	default:
		return "UNKN"
	}
}

func ideaRadarArchifyStanceVariant(stance domain.IdeaRadarStance) string {
	switch stance {
	case domain.IdeaRadarStanceConflict:
		return "security"
	case domain.IdeaRadarStanceAmbiguous:
		return "dashed"
	default:
		return "default"
	}
}

func ideaRadarArchifyOutcomeDot(outcome domain.IdeaRadarOutcome) string {
	switch outcome {
	case domain.IdeaRadarConvergent:
		return "emerald"
	case domain.IdeaRadarConflict:
		return "rose"
	case domain.IdeaRadarAmbiguous:
		return "amber"
	case domain.IdeaRadarUnknownOrigin:
		return "orange"
	case domain.IdeaRadarInsufficientPerspectives, domain.IdeaRadarInsufficientOrigins:
		return "violet"
	default:
		return "slate"
	}
}

func formatIdeaRadarProvenance(assessment domain.IdeaRadarAssessment) string {
	return strings.Join([]string{
		"Provenance: support=" + nonEmptyIdeaRadarArchifyList(assessment.SupportingObservationIDs),
		"conflict=" + nonEmptyIdeaRadarArchifyList(assessment.ConflictingObservationIDs),
		"ambiguous=" + nonEmptyIdeaRadarArchifyList(assessment.AmbiguousObservationIDs),
		"unknown_origin=" + nonEmptyIdeaRadarArchifyList(assessment.UnknownOriginObservationIDs),
		"origins=" + nonEmptyIdeaRadarArchifyList(assessment.DistinctOrigins),
	}, " | ")
}

func nonEmptyIdeaRadarArchifyList(items []string) string {
	if len(items) == 0 {
		return "-"
	}
	return strings.Join(items, ",")
}

func formatIdeaRadarProbe(probe domain.IdeaRadarProbe) string {
	parts := []string{"Probe " + probe.ID + ": " + string(probe.Kind), "reason=" + probe.Reason}
	if len(probe.ObservationIDs) > 0 {
		parts = append(parts, "observations="+strings.Join(probe.ObservationIDs, ","))
	}
	if probe.ProjectionKind != "" {
		parts = append(parts, "projection_kind="+string(probe.ProjectionKind))
	}
	if probe.OriginRequirement != domain.IdeaRadarOriginRequirementNone {
		parts = append(parts, "origin_requirement="+string(probe.OriginRequirement))
	}
	if len(probe.RequiredRoles) > 0 {
		roles := make([]string, len(probe.RequiredRoles))
		for index, role := range probe.RequiredRoles {
			roles[index] = string(role)
		}
		parts = append(parts, "required_roles="+strings.Join(roles, ","))
	}
	return strings.Join(parts, " | ")
}
