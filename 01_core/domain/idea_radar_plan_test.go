/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_radar_plan_test.md
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

func TestIdeaRadarPlanPublicWireContract(t *testing.T) {
	probe := domain.IdeaRadarProbe{
		ID:                "candidate#concept::p0001",
		Kind:              domain.IdeaRadarProbeKind("BIND_EVIDENCE"),
		Priority:          60,
		TargetID:          "candidate#concept",
		ObservationIDs:    []string{"obs-a"},
		ProjectionKind:    domain.IdeaProjectionKind(""),
		OriginRequirement: domain.IdeaRadarOriginRequirement("RECORDED"),
		RequiredRoles: []domain.EvidenceCoalitionRole{
			domain.EvidenceCoalitionRoleImplementation,
			domain.EvidenceCoalitionRoleIndependentTest,
		},
		Reason: "bind_physical_evidence",
	}
	report := domain.IdeaRadarPlanReport{
		SchemaVersion: "idea-radar-plan-report-v1",
		Status:        "COMPLETE",
		Reason:        "planned",
		Plans: []domain.IdeaRadarTargetPlan{{
			TargetID:            "candidate#concept",
			RadarOutcome:        domain.IdeaRadarOutcomeConvergent,
			Frontier:            []domain.IdeaRadarProbe{probe},
			EstimatedProbeCount: 1,
			Reason:              "evidence_binding_only",
		}},
		Metrics: domain.IdeaRadarPlanMetrics{
			Targets:            1,
			Plans:              1,
			Probes:             1,
			RadarOperations:    7,
			PlanningOperations: 5,
			Operations:         12,
			LimitExceeded:      false,
			VerdictChanges:     0,
		},
	}

	wantTags(t, reflect.TypeOf(probe), []string{
		"id", "kind", "priority", "target_id", "observation_ids", "projection_kind",
		"origin_requirement", "required_roles", "reason",
	})
	wantTags(t, reflect.TypeOf(report.Plans[0]), []string{
		"target_id", "radar_outcome", "frontier", "estimated_probe_count", "reason",
	})
	wantTags(t, reflect.TypeOf(report.Metrics), []string{
		"targets", "plans", "probes", "radar_operations", "planning_operations",
		"operations", "limit_exceeded", "verdict_changes",
	})
	wantTags(t, reflect.TypeOf(report), []string{"schema_version", "status", "reason", "plans", "metrics"})

	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var wire map[string]any
	if err := json.Unmarshal(encoded, &wire); err != nil {
		t.Fatalf("unmarshal wire: %v", err)
	}
	if wire["schema_version"] != "idea-radar-plan-report-v1" || wire["status"] != "COMPLETE" {
		t.Fatalf("unexpected report identity: %s", encoded)
	}
	if report.Metrics.VerdictChanges != 0 {
		t.Fatalf("planner changed a verdict: %+v", report.Metrics)
	}
	for _, forbidden := range []string{"ready", "ready_for_0065", "executed", "decision_impact", "evidence_source"} {
		if _, found := wire[forbidden]; found {
			t.Fatalf("forbidden top-level assertion %q in %s", forbidden, encoded)
		}
	}

	var roundTrip domain.IdeaRadarPlanReport
	if err := json.Unmarshal(encoded, &roundTrip); err != nil {
		t.Fatalf("round trip: %v", err)
	}
	if !reflect.DeepEqual(roundTrip, report) {
		t.Fatalf("round trip changed report\n got: %#v\nwant: %#v", roundTrip, report)
	}
}

func TestIdeaRadarPlanEnumWireValues(t *testing.T) {
	kinds := []domain.IdeaRadarProbeKind{
		"RESOLVE_CONFLICT", "RESOLVE_AMBIGUITY", "RECORD_ORIGIN",
		"ADD_PERSPECTIVE", "ADD_INDEPENDENT_ORIGIN", "BIND_EVIDENCE",
	}
	wantKinds := []string{
		"RESOLVE_CONFLICT", "RESOLVE_AMBIGUITY", "RECORD_ORIGIN",
		"ADD_PERSPECTIVE", "ADD_INDEPENDENT_ORIGIN", "BIND_EVIDENCE",
	}
	for i := range kinds {
		if string(kinds[i]) != wantKinds[i] {
			t.Fatalf("probe kind %d = %q, want %q", i, kinds[i], wantKinds[i])
		}
	}
	requirements := []domain.IdeaRadarOriginRequirement{"NONE", "DISTINCT_RECORDED", "RECORDED"}
	for i, want := range []string{"NONE", "DISTINCT_RECORDED", "RECORDED"} {
		if string(requirements[i]) != want {
			t.Fatalf("origin requirement %d = %q, want %q", i, requirements[i], want)
		}
	}
}

func wantTags(t *testing.T, typ reflect.Type, want []string) {
	t.Helper()
	if typ.NumField() != len(want) {
		t.Fatalf("%s has %d fields, want %d", typ.Name(), typ.NumField(), len(want))
	}
	for i, tag := range want {
		field := typ.Field(i)
		if got := field.Tag.Get("json"); got != tag {
			t.Fatalf("%s.%s json tag = %q, want %q", typ.Name(), field.Name, got, tag)
		}
		if field.Tag.Get("json") == "-" {
			t.Fatalf("%s.%s is hidden from JSON", typ.Name(), field.Name)
		}
	}
}
