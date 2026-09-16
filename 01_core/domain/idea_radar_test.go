/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_radar_test.md
 * @layer L1
 * @updated 2026-09-11
 */

package domain_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func TestIdeaRadarPublicWireContract(t *testing.T) {
	if domain.IdeaRadarInputSchemaVersion != "idea-radar-input-v1" {
		t.Fatalf("input schema = %q", domain.IdeaRadarInputSchemaVersion)
	}
	if domain.IdeaRadarReportSchemaVersion != "idea-radar-report-v1" {
		t.Fatalf("report schema = %q", domain.IdeaRadarReportSchemaVersion)
	}

	assertStringValues(t, "projection kinds", []string{
		string(domain.IdeaProjectionKindProblem),
		string(domain.IdeaProjectionKindBehavior),
		string(domain.IdeaProjectionKindArchitecture),
		string(domain.IdeaProjectionKindConstraint),
		string(domain.IdeaProjectionKindOperation),
	}, []string{"PROBLEM", "BEHAVIOR", "ARCHITECTURE", "CONSTRAINT", "OPERATION"})
	assertStringValues(t, "stances", []string{
		string(domain.IdeaRadarStanceSupport),
		string(domain.IdeaRadarStanceConflict),
		string(domain.IdeaRadarStanceAmbiguous),
	}, []string{"SUPPORT", "CONFLICT", "AMBIGUOUS"})
	assertStringValues(t, "outcomes", []string{
		string(domain.IdeaRadarOutcomeConvergent),
		string(domain.IdeaRadarOutcomeConflict),
		string(domain.IdeaRadarOutcomeAmbiguous),
		string(domain.IdeaRadarOutcomeInsufficientPerspectives),
		string(domain.IdeaRadarOutcomeUnknownOrigin),
		string(domain.IdeaRadarOutcomeInsufficientOrigins),
	}, []string{"CONVERGENT", "CONFLICT", "AMBIGUOUS", "INSUFFICIENT_PERSPECTIVES", "UNKNOWN_ORIGIN", "INSUFFICIENT_ORIGINS"})
	assertStringValues(t, "report states", []string{
		domain.IdeaRadarReportComplete,
		domain.IdeaRadarReportUnknown,
	}, []string{"COMPLETE", "UNKNOWN"})
}

func TestIdeaRadarJSONShapeAndRoundTrip(t *testing.T) {
	report := domain.IdeaRadarReport{
		SchemaVersion: "idea-radar-report-v1",
		Status:        "COMPLETE",
		Reason:        "analyzed",
		Assessments: []domain.IdeaRadarAssessment{{
			TargetID:                    "candidate#concept",
			CandidateID:                 "candidate",
			ConceptID:                   "concept",
			Outcome:                     domain.IdeaRadarOutcome("CONVERGENT"),
			SupportingObservationIDs:    []string{"obs-1"},
			ConflictingObservationIDs:   []string{},
			AmbiguousObservationIDs:     []string{},
			UnknownOriginObservationIDs: []string{},
			ProjectionKinds:             []domain.IdeaProjectionKind{domain.IdeaProjectionKind("ARCHITECTURE")},
			DistinctOrigins:             []string{"origin-A"},
			MissingProjectionKinds:      0,
			MissingOrigins:              0,
			CoalitionDraft: domain.IdeaRadarCoalitionDraft{
				ClaimID:        "candidate#concept",
				ObservationIDs: []string{"obs-1"},
				RequiredRoles: []domain.EvidenceCoalitionRole{
					domain.EvidenceCoalitionRoleImplementation,
					domain.EvidenceCoalitionRoleIndependentTest,
				},
				ReadyFor0065: false,
				Reason:       "requires_evidence_binding",
			},
			Reason: "idea_convergent",
		}},
		Metrics: domain.IdeaRadarMetrics{VerdictChanges: 0},
	}

	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var shape map[string]any
	if err := json.Unmarshal(raw, &shape); err != nil {
		t.Fatal(err)
	}
	assertKeys(t, shape, []string{"schema_version", "status", "reason", "assessments", "metrics"})
	assessment := shape["assessments"].([]any)[0].(map[string]any)
	assertKeys(t, assessment, []string{
		"target_id", "candidate_id", "concept_id", "outcome",
		"supporting_observation_ids", "conflicting_observation_ids",
		"ambiguous_observation_ids", "unknown_origin_observation_ids",
		"projection_kinds", "distinct_origins", "missing_projection_kinds",
		"missing_origins", "coalition_draft", "reason",
	})
	draft := assessment["coalition_draft"].(map[string]any)
	assertKeys(t, draft, []string{"claim_id", "observation_ids", "required_roles", "ready_for_0065", "reason"})
	metrics := shape["metrics"].(map[string]any)
	assertKeys(t, metrics, []string{
		"projections", "observations", "canonical_observations", "duplicate_observations",
		"targets", "convergent", "conflicts", "ambiguous", "insufficient_perspectives",
		"unknown_origins", "insufficient_origins", "operations", "limit_exceeded", "verdict_changes",
	})

	var roundTrip domain.IdeaRadarReport
	if err := json.Unmarshal(raw, &roundTrip); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(roundTrip, report) {
		t.Fatalf("round trip mismatch\n got: %#v\nwant: %#v", roundTrip, report)
	}
}

func TestIdeaRadarInputUsesExistingOriginTypeAndNonOmittingTags(t *testing.T) {
	input := domain.IdeaRadarInput{
		SchemaVersion: "idea-radar-input-v1",
		Projections: []domain.IdeaRadarProjection{{
			ID:           "p",
			Kind:         domain.IdeaProjectionKind("PROBLEM"),
			OriginStatus: domain.EvidenceOriginRecorded,
			OriginID:     "Origin-A",
			Terms:        []string{},
			Actions:      []string{},
			Objects:      []string{},
			Required:     []string{},
		}},
		Observations: []domain.IdeaRadarObservation{{
			ID:           "o",
			ProjectionID: "p",
			CandidateID:  "candidate",
			ConceptID:    "concept",
			Stance:       domain.IdeaRadarStance("SUPPORT"),
		}},
	}
	raw, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var shape map[string]any
	if err := json.Unmarshal(raw, &shape); err != nil {
		t.Fatal(err)
	}
	assertKeys(t, shape, []string{"schema_version", "projections", "observations"})
	projection := shape["projections"].([]any)[0].(map[string]any)
	assertKeys(t, projection, []string{"id", "kind", "origin_status", "origin_id", "terms", "actions", "objects", "required"})
	observation := shape["observations"].([]any)[0].(map[string]any)
	assertKeys(t, observation, []string{"id", "projection_id", "candidate_id", "concept_id", "stance"})
}

func assertStringValues(t *testing.T, label string, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
}

func assertKeys(t *testing.T, object map[string]any, want []string) {
	t.Helper()
	if len(object) != len(want) {
		t.Fatalf("keys = %v, want exactly %v", object, want)
	}
	for _, key := range want {
		if _, ok := object[key]; !ok {
			t.Errorf("missing JSON key %q in %v", key, object)
		}
	}
}
