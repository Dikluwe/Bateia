/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_radar_plan_test.md
 * @layer L1
 * @updated 2026-09-11
 */

package usecase_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

type planCase0067 struct {
	id          string
	family      string
	input       domain.IdeaRadarInput
	outcome     domain.IdeaRadarOutcome
	probes      []string
	naiveProbes int
	postOutcome domain.IdeaRadarOutcome
}

func TestPlanIdeaRadarProbesFrozenCorpus(t *testing.T) {
	cases := frozenPlanCases0067()
	plannerTotal, naiveTotal := 0, 0
	eligible, plannerResolved, naiveResolved, plannerConvergent, naiveConvergent := 0, 0, 0, 0, 0
	improvedCases := make([]string, 0)
	improvedFamilies := make(map[string]struct{})

	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			report := usecase.PlanIdeaRadarProbes(tc.input)
			if tc.id == "c10" {
				assertUnknownPlan0067(t, report, "invalid_schema")
				return
			}
			if report.SchemaVersion != "idea-radar-plan-report-v1" || report.Status != "COMPLETE" || report.Reason != "planned" {
				t.Fatalf("unexpected report identity: %+v", report)
			}
			if len(report.Plans) != 1 {
				t.Fatalf("plans = %d, want 1: %+v", len(report.Plans), report)
			}
			plan := report.Plans[0]
			if plan.RadarOutcome != tc.outcome {
				t.Fatalf("outcome = %s, want %s", plan.RadarOutcome, tc.outcome)
			}
			if got := probeDescriptors0067(plan.Frontier); !reflect.DeepEqual(got, tc.probes) {
				t.Fatalf("probes = %#v, want %#v\nplan=%+v", got, tc.probes, plan)
			}
			assertCanonicalProbeIDs0067(t, plan)
			assertProbeShapes0067(t, plan)
			wantPlanReason := "active_blocker_frontier"
			if tc.outcome == domain.IdeaRadarOutcomeConvergent {
				wantPlanReason = "evidence_binding_only"
			}
			if plan.Reason != wantPlanReason {
				t.Fatalf("plan reason = %q, want %q", plan.Reason, wantPlanReason)
			}
			assertPlanAccounting0067(t, tc.input, report)
			assertNonNilPlanSlices0067(t, report)
			if report.Metrics.VerdictChanges != 0 {
				t.Fatalf("verdict changes = %d", report.Metrics.VerdictChanges)
			}

			if tc.outcome == domain.IdeaRadarOutcomeConvergent {
				if plan.Reason != "evidence_binding_only" || len(plan.Frontier) != 1 {
					t.Fatalf("convergent target is not binding-only: %+v", plan)
				}
				return
			}

			post := usecase.AnalyzeIdeaRadar(applyPlan0067(t, tc.input, plan))
			if post.Status != domain.IdeaRadarReportComplete || len(post.Assessments) != 1 {
				t.Fatalf("planner simulation did not yield one assessment: %+v", post)
			}
			if post.Assessments[0].Outcome != tc.postOutcome {
				t.Fatalf("post outcome = %s, want %s", post.Assessments[0].Outcome, tc.postOutcome)
			}
			naiveInput, gotNaiveProbes := applyNaivePlan0067(t, tc.input)
			if gotNaiveProbes != tc.naiveProbes {
				t.Fatalf("naive probes = %d, want %d", gotNaiveProbes, tc.naiveProbes)
			}
			naivePost := usecase.AnalyzeIdeaRadar(naiveInput)
			if naivePost.Status != domain.IdeaRadarReportComplete || len(naivePost.Assessments) != 1 {
				t.Fatalf("naive simulation did not yield one assessment: %+v", naivePost)
			}

			eligible++
			plannerTotal += len(plan.Frontier)
			naiveTotal += gotNaiveProbes
			if post.Assessments[0].Outcome != tc.outcome {
				plannerResolved++
			}
			if naivePost.Assessments[0].Outcome != tc.outcome {
				naiveResolved++
			}
			if post.Assessments[0].Outcome == domain.IdeaRadarOutcomeConvergent {
				plannerConvergent++
			}
			if naivePost.Assessments[0].Outcome == domain.IdeaRadarOutcomeConvergent {
				naiveConvergent++
			}
			if len(plan.Frontier) < gotNaiveProbes {
				improvedCases = append(improvedCases, tc.id)
				improvedFamilies[tc.family] = struct{}{}
			}
		})
	}

	sort.Strings(improvedCases)
	if eligible != 8 || plannerResolved != 8 || naiveResolved != 8 {
		t.Fatalf("blocker resolution = eligible:%d planner:%d naive:%d, want 8/8 and 8/8", eligible, plannerResolved, naiveResolved)
	}
	if plannerConvergent != 7 || naiveConvergent != 7 {
		t.Fatalf("convergence = planner:%d naive:%d, want 7/7", plannerConvergent, naiveConvergent)
	}
	if plannerTotal != 11 || naiveTotal != 15 {
		t.Fatalf("eligible effort = planner:%d naive:%d, want 11/15", plannerTotal, naiveTotal)
	}
	if got := float64(naiveTotal-plannerTotal) / float64(naiveTotal); got != 4.0/15.0 {
		t.Fatalf("effort reduction = %v, want %v", got, 4.0/15.0)
	}
	if !reflect.DeepEqual(improvedCases, []string{"c01", "c02", "c04"}) || len(improvedFamilies) != 2 {
		t.Fatalf("improvements = %v across %d families", improvedCases, len(improvedFamilies))
	}
}

func TestPlanIdeaRadarExternalUtilityVerdicts(t *testing.T) {
	if got := utilityVerdict0067(8, 8, 8, 11, 15, 3, 2, false, false); got != "OBSERVED" {
		t.Fatalf("frozen utility verdict = %s, want OBSERVED", got)
	}
	if got := utilityVerdict0067(8, 7, 8, 15, 15, 0, 0, false, false); got != "NOT_OBSERVED" {
		t.Fatalf("valid non-improvement verdict = %s, want NOT_OBSERVED", got)
	}
	for name, got := range map[string]string{
		"undefined": utilityVerdict0067(0, 0, 0, 0, 0, 0, 0, false, false),
		"unknown":   utilityVerdict0067(8, 8, 8, 11, 15, 3, 2, true, false),
		"budget":    utilityVerdict0067(8, 8, 8, 11, 15, 3, 2, false, true),
	} {
		if got != "INCONCLUSIVE" {
			t.Fatalf("%s verdict = %s, want INCONCLUSIVE", name, got)
		}
	}
}

func TestPlanIdeaRadarProbesActiveBlockerAndEveryOutcome(t *testing.T) {
	input := radarInput0067("precedence", []domain.IdeaRadarProjection{
		projection0067("p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginUnknown, ""),
		projection0067("p2", domain.IdeaProjectionKindProblem, domain.EvidenceOriginUnknown, ""),
	}, []domain.IdeaRadarObservation{
		observation0067("z-conflict", "p2", "precedence", domain.IdeaRadarStanceConflict),
		observation0067("a-conflict", "p", "precedence", domain.IdeaRadarStanceConflict),
		observation0067("hidden-ambiguous", "p", "precedence", domain.IdeaRadarStanceAmbiguous),
	})
	report := usecase.PlanIdeaRadarProbes(input)
	plan := requireSinglePlan0067(t, report)
	if plan.RadarOutcome != domain.IdeaRadarOutcomeConflict || !reflect.DeepEqual(probeDescriptors0067(plan.Frontier), []string{
		"RESOLVE_CONFLICT:a-conflict:NONE", "RESOLVE_CONFLICT:z-conflict:NONE",
	}) {
		t.Fatalf("conflict did not exclusively hide lower blockers: %+v", plan)
	}
	for _, probe := range plan.Frontier {
		if probe.Priority != 10 || probe.Reason != "resolve_conflicting_observation" || probe.ProjectionKind != "" || len(probe.RequiredRoles) != 0 {
			t.Fatalf("invalid conflict probe: %+v", probe)
		}
	}

	input.Observations = input.Observations[2:]
	plan = requireSinglePlan0067(t, usecase.PlanIdeaRadarProbes(input))
	if plan.RadarOutcome != domain.IdeaRadarOutcomeAmbiguous || !reflect.DeepEqual(probeDescriptors0067(plan.Frontier), []string{"RESOLVE_AMBIGUITY:hidden-ambiguous:NONE"}) {
		t.Fatalf("ambiguity did not exclusively hide lower blockers: %+v", plan)
	}
	if plan.Frontier[0].Priority != 20 || plan.Frontier[0].Reason != "resolve_ambiguous_observation" {
		t.Fatalf("invalid ambiguity probe: %+v", plan.Frontier[0])
	}
}

func TestPlanIdeaRadarProbesMinimalUnknownSubsetAndBinding(t *testing.T) {
	for _, id := range []string{"c04", "c08", "c09"} {
		tc := caseByID0067(t, id)
		plan := requireSinglePlan0067(t, usecase.PlanIdeaRadarProbes(tc.input))
		switch id {
		case "c04":
			if len(plan.Frontier) != 1 || !reflect.DeepEqual(plan.Frontier[0].ObservationIDs, []string{"obs-arch-1"}) || plan.Frontier[0].Priority != 30 {
				t.Fatalf("unknown-origin frontier is not the minimal canonical subset: %+v", plan)
			}
		case "c08":
			probe := plan.Frontier[0]
			if probe.Kind != "BIND_EVIDENCE" || probe.Priority != 60 || probe.OriginRequirement != "RECORDED" || probe.ProjectionKind != "" {
				t.Fatalf("binding probe is not conservative: %+v", probe)
			}
			if !reflect.DeepEqual(probe.ObservationIDs, []string{"obs-a", "obs-b", "obs-p"}) || !reflect.DeepEqual(probe.RequiredRoles, []domain.EvidenceCoalitionRole{
				domain.EvidenceCoalitionRoleImplementation, domain.EvidenceCoalitionRoleIndependentTest,
			}) {
				t.Fatalf("binding probe did not reuse the 0066 draft: %+v", probe)
			}
		case "c09":
			if plan.RadarOutcome != domain.IdeaRadarOutcomeConflict || len(plan.Frontier) != 1 || plan.Reason != "active_blocker_frontier" {
				t.Fatalf("C09 falsely promises convergence: %+v", plan)
			}
		}
	}
}

func TestPlanIdeaRadarProbesUnknownOriginThenResidualOrigin(t *testing.T) {
	input := fixture0067("mixed-origin", []projectionSpec0067{
		{"p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, "shared"},
		{"b", domain.IdeaProjectionKindBehavior, domain.EvidenceOriginRecorded, "shared"},
		{"a", domain.IdeaProjectionKindArchitecture, domain.EvidenceOriginUnknown, ""},
	}, []observationSpec0067{
		{"obs-p", "p", domain.IdeaRadarStanceSupport},
		{"obs-b", "b", domain.IdeaRadarStanceSupport},
		{"obs-a", "a", domain.IdeaRadarStanceSupport},
	})
	plan := requireSinglePlan0067(t, usecase.PlanIdeaRadarProbes(input))
	if got := probeDescriptors0067(plan.Frontier); !reflect.DeepEqual(got, []string{
		"RECORD_ORIGIN:obs-a:DISTINCT_RECORDED",
		"ADD_INDEPENDENT_ORIGIN:PROBLEM:DISTINCT_RECORDED",
	}) {
		t.Fatalf("mixed origin frontier = %#v, want record then one residual origin", got)
	}
	post := usecase.AnalyzeIdeaRadar(applyPlan0067(t, input, plan))
	if post.Status != domain.IdeaRadarReportComplete || len(post.Assessments) != 1 || post.Assessments[0].Outcome != domain.IdeaRadarOutcomeConvergent {
		t.Fatalf("mixed origin probes did not converge under simulation: %+v", post)
	}
}

func TestPlanIdeaRadarProbesDoesNotInferTextAndDeduplicatesSemantics(t *testing.T) {
	tc := caseByID0067(t, "c02")
	want := usecase.PlanIdeaRadarProbes(tc.input)
	textual := cloneInput0067(t, tc.input)
	textual.Projections[0].Terms = []string{"architecture", "constraint", "operation"}
	textual.Projections[0].Actions = []string{"invent independent origin"}
	textual.Projections[0].Objects = []string{"behavior"}
	textual.Projections[0].Required = []string{"convergent"}
	if got := usecase.PlanIdeaRadarProbes(textual); !reflect.DeepEqual(got.Plans, want.Plans) {
		t.Fatalf("free text changed structural plan\n got=%+v\nwant=%+v", got.Plans, want.Plans)
	}
	repeated := cloneInput0067(t, tc.input)
	duplicate := repeated.Observations[0]
	duplicate.ID = "zz-semantic-duplicate"
	repeated.Observations = append(repeated.Observations, duplicate)
	if got := usecase.PlanIdeaRadarProbes(repeated); !reflect.DeepEqual(got.Plans, want.Plans) {
		t.Fatalf("semantic repetition changed plan\n got=%+v\nwant=%+v", got.Plans, want.Plans)
	}
}

func TestPlanIdeaRadarProbesUnknownNilAndNoPartial(t *testing.T) {
	inputs := []struct {
		name   string
		input  domain.IdeaRadarInput
		reason string
	}{
		{"invalid-schema", domain.IdeaRadarInput{SchemaVersion: "idea-radar-input-v999"}, "invalid_schema"},
		{"nil", domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion}, "no_projections"},
		{"empty", domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion, Projections: []domain.IdeaRadarProjection{}, Observations: []domain.IdeaRadarObservation{}}, "no_projections"},
		{"no-observations", domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion, Projections: []domain.IdeaRadarProjection{projection0067("p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, "o")}}, "no_observations"},
	}
	var nilJSON, emptyJSON []byte
	for _, tc := range inputs {
		t.Run(tc.name, func(t *testing.T) {
			report := usecase.PlanIdeaRadarProbes(tc.input)
			assertUnknownPlan0067(t, report, tc.reason)
			encoded, err := json.Marshal(report)
			if err != nil {
				t.Fatal(err)
			}
			if tc.name == "nil" {
				nilJSON = encoded
			}
			if tc.name == "empty" {
				emptyJSON = encoded
			}
		})
	}
	if !reflect.DeepEqual(nilJSON, emptyJSON) {
		t.Fatalf("nil and empty inputs differ:\n%s\n%s", nilJSON, emptyJSON)
	}
}

func TestPlanIdeaRadarProbesInclusiveAndExclusiveLimits(t *testing.T) {
	projections := []domain.IdeaRadarProjection{
		projection0067("p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, "origin"),
		projection0067("b", domain.IdeaProjectionKindBehavior, domain.EvidenceOriginRecorded, "origin-b"),
	}
	observations := make([]domain.IdeaRadarObservation, 0, 1024)
	for i := 0; i < 512; i++ {
		candidate := fmt.Sprintf("target-%04d", i)
		observations = append(observations,
			observation0067(fmt.Sprintf("conflict-p-%04d", i), "p", candidate, domain.IdeaRadarStanceConflict),
			observation0067(fmt.Sprintf("conflict-b-%04d", i), "b", candidate, domain.IdeaRadarStanceConflict),
		)
	}
	inclusive := domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion, Projections: projections, Observations: observations}
	report := usecase.PlanIdeaRadarProbes(inclusive)
	if report.Status != "COMPLETE" || len(report.Plans) != 512 || report.Metrics.Probes != 1024 || report.Metrics.LimitExceeded {
		t.Fatalf("inclusive plan/probe limits rejected: status=%s plans=%d metrics=%+v", report.Status, len(report.Plans), report.Metrics)
	}
	assertPlanAccounting0067(t, inclusive, report)

	overObservations := append(cloneObservations0067(observations), observation0067("extra", "p", "extra", domain.IdeaRadarStanceConflict))
	assertBudgetUnknown0067(t, usecase.PlanIdeaRadarProbes(domain.IdeaRadarInput{
		SchemaVersion: domain.IdeaRadarInputSchemaVersion, Projections: projections, Observations: overObservations,
	}))
	overProjections := make([]domain.IdeaRadarProjection, 65)
	for i := range overProjections {
		overProjections[i] = projection0067(fmt.Sprintf("p-%02d", i), domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, fmt.Sprintf("o-%02d", i))
	}
	assertBudgetUnknown0067(t, usecase.PlanIdeaRadarProbes(domain.IdeaRadarInput{
		SchemaVersion: domain.IdeaRadarInputSchemaVersion,
		Projections:   overProjections,
		Observations:  []domain.IdeaRadarObservation{observation0067("obs", "p-00", "candidate", domain.IdeaRadarStanceSupport)},
	}))

	maxID := strings.Repeat("a", 128)
	accepted := radarInput0067("limit", []domain.IdeaRadarProjection{
		projection0067(maxID, domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, strings.Repeat("o", 128)),
	}, []domain.IdeaRadarObservation{{ID: maxID, ProjectionID: maxID, CandidateID: maxID, ConceptID: maxID, Stance: domain.IdeaRadarStanceSupport}})
	if got := usecase.PlanIdeaRadarProbes(accepted); got.Status != "COMPLETE" || got.Metrics.LimitExceeded {
		t.Fatalf("inclusive scalar limit rejected: %+v", got)
	}
	rejected := accepted
	rejected.Observations = cloneObservations0067(accepted.Observations)
	rejected.Observations[0].CandidateID += "a"
	assertBudgetUnknown0067(t, usecase.PlanIdeaRadarProbes(rejected))

	textBoundary := caseByID0067(t, "c02").input
	textBoundary.Projections[0].Terms = make([]string, 64)
	for i := range textBoundary.Projections[0].Terms {
		textBoundary.Projections[0].Terms[i] = "x"
	}
	if got := usecase.PlanIdeaRadarProbes(textBoundary); got.Status != "COMPLETE" {
		t.Fatalf("inclusive text item limit rejected: %+v", got)
	}
	textExceeded := cloneInput0067(t, textBoundary)
	textExceeded.Projections[0].Terms = append(textExceeded.Projections[0].Terms, "x")
	assertBudgetUnknown0067(t, usecase.PlanIdeaRadarProbes(textExceeded))
	textBoundary.Projections[0].Terms = []string{strings.Repeat("x", 4096)}
	if got := usecase.PlanIdeaRadarProbes(textBoundary); got.Status != "COMPLETE" {
		t.Fatalf("inclusive text byte limit rejected: %+v", got)
	}
	textBoundary.Projections[0].Terms[0] += "x"
	assertBudgetUnknown0067(t, usecase.PlanIdeaRadarProbes(textBoundary))
}

func TestPlanIdeaRadarProbesOperationsAndImmutability(t *testing.T) {
	tc := caseByID0067(t, "c02")
	before := mustJSON0067(t, tc.input)
	report := usecase.PlanIdeaRadarProbes(tc.input)
	after := mustJSON0067(t, tc.input)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("planner mutated its input\nbefore=%s\nafter=%s", before, after)
	}
	assertPlanAccounting0067(t, tc.input, report)
	if report.Metrics.RadarOperations != 6 || report.Metrics.PlanningOperations != 6 || report.Metrics.Operations != 12 {
		t.Fatalf("C02 operations = %+v, want radar=6 planning=6 total=12", report.Metrics)
	}
}

func TestPlanIdeaRadarProbesTwentyPermutationsAndNilEmptyText(t *testing.T) {
	base := combinedCorpusInput0067(t)
	baseReport := usecase.PlanIdeaRadarProbes(base)
	assertSortedPlans0067(t, baseReport)
	want := mustJSON0067(t, baseReport)
	for i := 0; i < 20; i++ {
		permuted := cloneInput0067(t, base)
		rotateProjections0067(permuted.Projections, (i*7+3)%len(permuted.Projections))
		rotateObservations0067(permuted.Observations, (i*11+5)%len(permuted.Observations))
		if i%2 == 1 {
			reverseProjections0067(permuted.Projections)
		}
		if i%3 == 1 {
			reverseObservations0067(permuted.Observations)
		}
		for p := range permuted.Projections {
			if (i+p)%2 == 0 {
				permuted.Projections[p].Terms = []string{}
				permuted.Projections[p].Actions = []string{}
				permuted.Projections[p].Objects = []string{}
				permuted.Projections[p].Required = []string{}
			} else {
				permuted.Projections[p].Terms = nil
				permuted.Projections[p].Actions = nil
				permuted.Projections[p].Objects = nil
				permuted.Projections[p].Required = nil
			}
		}
		got := mustJSON0067(t, usecase.PlanIdeaRadarProbes(permuted))
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("permutation %d changed output\n got=%s\nwant=%s", i, got, want)
		}
	}
}

func assertSortedPlans0067(t *testing.T, report domain.IdeaRadarPlanReport) {
	t.Helper()
	if report.Status != "COMPLETE" || len(report.Plans) != 9 {
		t.Fatalf("combined corpus report = %+v", report)
	}
	for i := 1; i < len(report.Plans); i++ {
		if report.Plans[i-1].TargetID >= report.Plans[i].TargetID {
			t.Fatalf("plans are not strictly sorted at %d: %q then %q", i, report.Plans[i-1].TargetID, report.Plans[i].TargetID)
		}
	}
}

func frozenPlanCases0067() []planCase0067 {
	return []planCase0067{
		{"c01", "coverage-gap", fixture0067("c01", []projectionSpec0067{{"p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, "or-a"}}, []observationSpec0067{{"obs-p", "p", domain.IdeaRadarStanceSupport}}), domain.IdeaRadarOutcomeInsufficientPerspectives, []string{"ADD_PERSPECTIVE:BEHAVIOR:DISTINCT_RECORDED", "ADD_PERSPECTIVE:ARCHITECTURE:DISTINCT_RECORDED"}, 4, domain.IdeaRadarOutcomeConvergent},
		{"c02", "coverage-gap", fixture0067("c02", []projectionSpec0067{{"p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, "or-a"}, {"b", domain.IdeaProjectionKindBehavior, domain.EvidenceOriginRecorded, "or-b"}}, []observationSpec0067{{"obs-p", "p", domain.IdeaRadarStanceSupport}, {"obs-b", "b", domain.IdeaRadarStanceSupport}}), domain.IdeaRadarOutcomeInsufficientPerspectives, []string{"ADD_PERSPECTIVE:ARCHITECTURE:DISTINCT_RECORDED"}, 2, domain.IdeaRadarOutcomeConvergent},
		{"c03", "provenance-gap", fixture0067("c03", standardProjections0067(false, 1), standardSupports0067(1)), domain.IdeaRadarOutcomeUnknownOrigin, []string{"RECORD_ORIGIN:obs-arch-1:DISTINCT_RECORDED"}, 1, domain.IdeaRadarOutcomeConvergent},
		{"c04", "provenance-gap", fixture0067("c04", standardProjections0067(false, 2), standardSupports0067(2)), domain.IdeaRadarOutcomeUnknownOrigin, []string{"RECORD_ORIGIN:obs-arch-1:DISTINCT_RECORDED"}, 2, domain.IdeaRadarOutcomeConvergent},
		{"c05", "provenance-gap", fixture0067("c05", []projectionSpec0067{{"p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, "shared"}, {"b", domain.IdeaProjectionKindBehavior, domain.EvidenceOriginRecorded, "shared"}, {"a", domain.IdeaProjectionKindArchitecture, domain.EvidenceOriginRecorded, "shared"}}, []observationSpec0067{{"obs-p", "p", domain.IdeaRadarStanceSupport}, {"obs-b", "b", domain.IdeaRadarStanceSupport}, {"obs-a", "a", domain.IdeaRadarStanceSupport}}), domain.IdeaRadarOutcomeInsufficientOrigins, []string{"ADD_INDEPENDENT_ORIGIN:PROBLEM:DISTINCT_RECORDED", "ADD_INDEPENDENT_ORIGIN:BEHAVIOR:DISTINCT_RECORDED"}, 2, domain.IdeaRadarOutcomeConvergent},
		{"c06", "signal-resolution", fixture0067("c06", standardProjections0067(true, 1), append(standardRecordedSupports0067(), observationSpec0067{"obs-conflict-2", "b", domain.IdeaRadarStanceConflict}, observationSpec0067{"obs-conflict-1", "p", domain.IdeaRadarStanceConflict})), domain.IdeaRadarOutcomeConflict, []string{"RESOLVE_CONFLICT:obs-conflict-1:NONE", "RESOLVE_CONFLICT:obs-conflict-2:NONE"}, 2, domain.IdeaRadarOutcomeConvergent},
		{"c07", "signal-resolution", fixture0067("c07", standardProjections0067(true, 1), append(standardRecordedSupports0067(), observationSpec0067{"obs-ambiguous-1", "p", domain.IdeaRadarStanceAmbiguous})), domain.IdeaRadarOutcomeAmbiguous, []string{"RESOLVE_AMBIGUITY:obs-ambiguous-1:NONE"}, 1, domain.IdeaRadarOutcomeConvergent},
		{"c08", "evidence-binding", fixture0067("c08", standardProjections0067(true, 1), standardRecordedSupports0067()), domain.IdeaRadarOutcomeConvergent, []string{"BIND_EVIDENCE:RECORDED:IMPLEMENTATION,INDEPENDENT_TEST"}, 1, ""},
		{"c09", "safety-control", fixture0067("c09", []projectionSpec0067{{"p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, "or-a"}}, []observationSpec0067{{"obs-p", "p", domain.IdeaRadarStanceSupport}, {"obs-conflict-1", "p", domain.IdeaRadarStanceConflict}}), domain.IdeaRadarOutcomeConflict, []string{"RESOLVE_CONFLICT:obs-conflict-1:NONE"}, 1, domain.IdeaRadarOutcomeInsufficientPerspectives},
		{"c10", "opaque-control", domain.IdeaRadarInput{SchemaVersion: "idea-radar-input-v999"}, "", nil, 0, ""},
	}
}

type projectionSpec0067 struct {
	id       string
	kind     domain.IdeaProjectionKind
	status   domain.EvidenceOriginStatus
	originID string
}

type observationSpec0067 struct {
	id         string
	projection string
	stance     domain.IdeaRadarStance
}

func fixture0067(candidate string, projections []projectionSpec0067, observations []observationSpec0067) domain.IdeaRadarInput {
	result := domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion, Projections: make([]domain.IdeaRadarProjection, 0, len(projections)), Observations: make([]domain.IdeaRadarObservation, 0, len(observations))}
	for _, projection := range projections {
		result.Projections = append(result.Projections, projection0067(projection.id, projection.kind, projection.status, projection.originID))
	}
	for _, observation := range observations {
		result.Observations = append(result.Observations, observation0067(observation.id, observation.projection, candidate, observation.stance))
	}
	return result
}

func standardProjections0067(recordedArchitecture bool, architectures int) []projectionSpec0067 {
	result := []projectionSpec0067{{"p", domain.IdeaProjectionKindProblem, domain.EvidenceOriginRecorded, "or-a"}, {"b", domain.IdeaProjectionKindBehavior, domain.EvidenceOriginRecorded, "or-b"}}
	for i := 1; i <= architectures; i++ {
		status, origin := domain.EvidenceOriginUnknown, ""
		if recordedArchitecture {
			status, origin = domain.EvidenceOriginRecorded, fmt.Sprintf("or-arch-%d", i)
		}
		result = append(result, projectionSpec0067{fmt.Sprintf("a%d", i), domain.IdeaProjectionKindArchitecture, status, origin})
	}
	return result
}

func standardSupports0067(architectures int) []observationSpec0067 {
	result := []observationSpec0067{{"obs-p", "p", domain.IdeaRadarStanceSupport}, {"obs-b", "b", domain.IdeaRadarStanceSupport}}
	for i := 1; i <= architectures; i++ {
		result = append(result, observationSpec0067{fmt.Sprintf("obs-arch-%d", i), fmt.Sprintf("a%d", i), domain.IdeaRadarStanceSupport})
	}
	return result
}

func standardRecordedSupports0067() []observationSpec0067 {
	return []observationSpec0067{{"obs-p", "p", domain.IdeaRadarStanceSupport}, {"obs-b", "b", domain.IdeaRadarStanceSupport}, {"obs-a", "a1", domain.IdeaRadarStanceSupport}}
}

func radarInput0067(candidate string, projections []domain.IdeaRadarProjection, observations []domain.IdeaRadarObservation) domain.IdeaRadarInput {
	return domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion, Projections: projections, Observations: observations}
}

func projection0067(id string, kind domain.IdeaProjectionKind, status domain.EvidenceOriginStatus, origin string) domain.IdeaRadarProjection {
	return domain.IdeaRadarProjection{ID: id, Kind: kind, OriginStatus: status, OriginID: origin}
}

func observation0067(id, projection, candidate string, stance domain.IdeaRadarStance) domain.IdeaRadarObservation {
	return domain.IdeaRadarObservation{ID: id, ProjectionID: projection, CandidateID: candidate, ConceptID: "concept", Stance: stance}
}

func probeDescriptors0067(probes []domain.IdeaRadarProbe) []string {
	result := make([]string, 0, len(probes))
	for _, probe := range probes {
		switch probe.Kind {
		case "RESOLVE_CONFLICT", "RESOLVE_AMBIGUITY", "RECORD_ORIGIN":
			observation := ""
			if len(probe.ObservationIDs) > 0 {
				observation = probe.ObservationIDs[0]
			}
			result = append(result, fmt.Sprintf("%s:%s:%s", probe.Kind, observation, probe.OriginRequirement))
		case "ADD_PERSPECTIVE", "ADD_INDEPENDENT_ORIGIN":
			result = append(result, fmt.Sprintf("%s:%s:%s", probe.Kind, probe.ProjectionKind, probe.OriginRequirement))
		case "BIND_EVIDENCE":
			roles := make([]string, len(probe.RequiredRoles))
			for i, role := range probe.RequiredRoles {
				roles[i] = string(role)
			}
			result = append(result, fmt.Sprintf("%s:%s:%s", probe.Kind, probe.OriginRequirement, strings.Join(roles, ",")))
		default:
			result = append(result, "UNKNOWN:"+string(probe.Kind))
		}
	}
	return result
}

func assertCanonicalProbeIDs0067(t *testing.T, plan domain.IdeaRadarTargetPlan) {
	t.Helper()
	if plan.EstimatedProbeCount != len(plan.Frontier) {
		t.Fatalf("estimated probes = %d, frontier = %d", plan.EstimatedProbeCount, len(plan.Frontier))
	}
	for i, probe := range plan.Frontier {
		want := fmt.Sprintf("%s::p%04d", plan.TargetID, i+1)
		if probe.ID != want || probe.TargetID != plan.TargetID {
			t.Fatalf("probe %d identity = %q/%q, want %q", i, probe.ID, probe.TargetID, want)
		}
		if probe.ObservationIDs == nil || probe.RequiredRoles == nil {
			t.Fatalf("probe %d contains null arrays: %+v", i, probe)
		}
	}
}

func assertProbeShapes0067(t *testing.T, plan domain.IdeaRadarTargetPlan) {
	t.Helper()
	priorities := map[domain.IdeaRadarProbeKind]int{
		"RESOLVE_CONFLICT":       10,
		"RESOLVE_AMBIGUITY":      20,
		"RECORD_ORIGIN":          30,
		"ADD_PERSPECTIVE":        40,
		"ADD_INDEPENDENT_ORIGIN": 50,
		"BIND_EVIDENCE":          60,
	}
	reasons := map[domain.IdeaRadarProbeKind]string{
		"RESOLVE_CONFLICT":       "resolve_conflicting_observation",
		"RESOLVE_AMBIGUITY":      "resolve_ambiguous_observation",
		"RECORD_ORIGIN":          "record_distinct_origin",
		"ADD_PERSPECTIVE":        "add_missing_projection_with_distinct_origin",
		"ADD_INDEPENDENT_ORIGIN": "add_independent_origin",
		"BIND_EVIDENCE":          "bind_physical_evidence",
	}
	for _, probe := range plan.Frontier {
		if probe.Priority != priorities[probe.Kind] || probe.Reason != reasons[probe.Kind] {
			t.Fatalf("wrong priority/reason for %s: %+v", probe.Kind, probe)
		}
		switch probe.Kind {
		case "RESOLVE_CONFLICT", "RESOLVE_AMBIGUITY":
			if probe.OriginRequirement != "NONE" || probe.ProjectionKind != "" || len(probe.ObservationIDs) != 1 || len(probe.RequiredRoles) != 0 {
				t.Fatalf("invalid signal-resolution probe: %+v", probe)
			}
		case "RECORD_ORIGIN":
			if probe.OriginRequirement != "DISTINCT_RECORDED" || len(probe.ObservationIDs) != 1 || len(probe.RequiredRoles) != 0 {
				t.Fatalf("invalid record-origin probe: %+v", probe)
			}
		case "ADD_PERSPECTIVE", "ADD_INDEPENDENT_ORIGIN":
			if probe.OriginRequirement != "DISTINCT_RECORDED" || probe.ProjectionKind == "" || len(probe.ObservationIDs) != 0 || len(probe.RequiredRoles) != 0 {
				t.Fatalf("invalid additive probe: %+v", probe)
			}
		case "BIND_EVIDENCE":
			if probe.OriginRequirement != "RECORDED" || probe.ProjectionKind != "" || len(probe.ObservationIDs) == 0 || len(probe.RequiredRoles) != 2 {
				t.Fatalf("invalid binding probe: %+v", probe)
			}
		default:
			t.Fatalf("unknown probe kind: %+v", probe)
		}
	}
}

func assertPlanAccounting0067(t *testing.T, input domain.IdeaRadarInput, report domain.IdeaRadarPlanReport) {
	t.Helper()
	radar := usecase.AnalyzeIdeaRadar(input)
	wantPlanning := len(input.Projections) + len(input.Observations) + len(radar.Assessments) + report.Metrics.Probes
	if report.Metrics.Targets != len(radar.Assessments) || report.Metrics.Plans != len(report.Plans) || report.Metrics.RadarOperations != radar.Metrics.Operations || report.Metrics.PlanningOperations != wantPlanning || report.Metrics.Operations != radar.Metrics.Operations+wantPlanning {
		t.Fatalf("accounting mismatch: radar=%+v plan=%+v want planning=%d", radar.Metrics, report.Metrics, wantPlanning)
	}
}

func assertNonNilPlanSlices0067(t *testing.T, report domain.IdeaRadarPlanReport) {
	t.Helper()
	if report.Plans == nil {
		t.Fatal("Plans must not be nil")
	}
	for i, plan := range report.Plans {
		if plan.Frontier == nil {
			t.Fatalf("plan %d frontier is nil", i)
		}
		assertCanonicalProbeIDs0067(t, plan)
	}
}

func assertUnknownPlan0067(t *testing.T, report domain.IdeaRadarPlanReport, reason string) {
	t.Helper()
	if report.SchemaVersion != "idea-radar-plan-report-v1" || report.Status != "UNKNOWN" || report.Reason != reason || report.Plans == nil || len(report.Plans) != 0 || report.Metrics.Plans != 0 || report.Metrics.Probes != 0 || report.Metrics.VerdictChanges != 0 {
		t.Fatalf("unexpected UNKNOWN report: %+v", report)
	}
}

func assertBudgetUnknown0067(t *testing.T, report domain.IdeaRadarPlanReport) {
	t.Helper()
	assertUnknownPlan0067(t, report, "budget_exceeded")
	if !report.Metrics.LimitExceeded {
		t.Fatalf("budget report did not mark limit: %+v", report)
	}
}

func requireSinglePlan0067(t *testing.T, report domain.IdeaRadarPlanReport) domain.IdeaRadarTargetPlan {
	t.Helper()
	if report.Status != "COMPLETE" || len(report.Plans) != 1 {
		t.Fatalf("want one complete plan: %+v", report)
	}
	return report.Plans[0]
}

func caseByID0067(t *testing.T, id string) planCase0067 {
	t.Helper()
	for _, tc := range frozenPlanCases0067() {
		if tc.id == id {
			return tc
		}
	}
	t.Fatalf("missing frozen case %s", id)
	return planCase0067{}
}

func applyPlan0067(t *testing.T, input domain.IdeaRadarInput, plan domain.IdeaRadarTargetPlan) domain.IdeaRadarInput {
	t.Helper()
	result := cloneInput0067(t, input)
	for i, probe := range plan.Frontier {
		switch probe.Kind {
		case "RESOLVE_CONFLICT":
			result.Observations = removeObservations0067(result.Observations, probe.ObservationIDs)
		case "RESOLVE_AMBIGUITY":
			for j := range result.Observations {
				if contains0067(probe.ObservationIDs, result.Observations[j].ID) {
					result.Observations[j].Stance = domain.IdeaRadarStanceSupport
				}
			}
		case "RECORD_ORIGIN":
			for _, observation := range result.Observations {
				if contains0067(probe.ObservationIDs, observation.ID) {
					for j := range result.Projections {
						if result.Projections[j].ID == observation.ProjectionID {
							result.Projections[j].OriginStatus = domain.EvidenceOriginRecorded
							result.Projections[j].OriginID = fmt.Sprintf("sim-origin-%04d", i+1)
						}
					}
				}
			}
		case "ADD_PERSPECTIVE", "ADD_INDEPENDENT_ORIGIN":
			candidate, concept := splitTarget0067(t, plan.TargetID)
			projectionID := fmt.Sprintf("sim-projection-%04d", i+1)
			result.Projections = append(result.Projections, projection0067(projectionID, probe.ProjectionKind, domain.EvidenceOriginRecorded, fmt.Sprintf("sim-origin-%04d", i+1)))
			result.Observations = append(result.Observations, domain.IdeaRadarObservation{ID: fmt.Sprintf("sim-observation-%04d", i+1), ProjectionID: projectionID, CandidateID: candidate, ConceptID: concept, Stance: domain.IdeaRadarStanceSupport})
		case "BIND_EVIDENCE":
		default:
			t.Fatalf("unsupported probe kind %q", probe.Kind)
		}
	}
	return result
}

func applyNaivePlan0067(t *testing.T, input domain.IdeaRadarInput) (domain.IdeaRadarInput, int) {
	t.Helper()
	result := cloneInput0067(t, input)
	radar := usecase.AnalyzeIdeaRadar(result)
	if radar.Status != domain.IdeaRadarReportComplete || len(radar.Assessments) != 1 {
		t.Fatalf("naive input is not one complete assessment: %+v", radar)
	}
	assessment := radar.Assessments[0]
	count := 0
	switch assessment.Outcome {
	case domain.IdeaRadarOutcomeConflict:
		count = len(assessment.ConflictingObservationIDs)
		result.Observations = removeObservations0067(result.Observations, assessment.ConflictingObservationIDs)
	case domain.IdeaRadarOutcomeAmbiguous:
		count = len(assessment.AmbiguousObservationIDs)
		for i := range result.Observations {
			if contains0067(assessment.AmbiguousObservationIDs, result.Observations[i].ID) {
				result.Observations[i].Stance = domain.IdeaRadarStanceSupport
			}
		}
	case domain.IdeaRadarOutcomeInsufficientPerspectives:
		missing := missingKinds0067(assessment.ProjectionKinds, assessment.MissingProjectionKinds)
		for _, kind := range missing {
			result = appendNaiveSupport0067(t, result, assessment.TargetID, kind, count)
			count++
		}
		for i := 0; i < assessment.MissingOrigins; i++ {
			result = appendNaiveSupport0067(t, result, assessment.TargetID, assessment.ProjectionKinds[0], count)
			count++
		}
	case domain.IdeaRadarOutcomeUnknownOrigin:
		for _, id := range assessment.UnknownOriginObservationIDs {
			for _, observation := range result.Observations {
				if observation.ID == id {
					for i := range result.Projections {
						if result.Projections[i].ID == observation.ProjectionID {
							result.Projections[i].OriginStatus = domain.EvidenceOriginRecorded
							result.Projections[i].OriginID = fmt.Sprintf("naive-origin-%04d", count+1)
							count++
						}
					}
				}
			}
		}
		intermediate := usecase.AnalyzeIdeaRadar(result)
		if intermediate.Status == domain.IdeaRadarReportComplete && len(intermediate.Assessments) == 1 && intermediate.Assessments[0].Outcome == domain.IdeaRadarOutcomeInsufficientOrigins {
			for i := 0; i < intermediate.Assessments[0].MissingOrigins; i++ {
				result = appendNaiveSupport0067(t, result, assessment.TargetID, assessment.ProjectionKinds[0], count)
				count++
			}
		}
	case domain.IdeaRadarOutcomeInsufficientOrigins:
		for i := 0; i < assessment.MissingOrigins; i++ {
			result = appendNaiveSupport0067(t, result, assessment.TargetID, assessment.ProjectionKinds[i%len(assessment.ProjectionKinds)], count)
			count++
		}
	}
	return result, count
}

func appendNaiveSupport0067(t *testing.T, input domain.IdeaRadarInput, target string, kind domain.IdeaProjectionKind, sequence int) domain.IdeaRadarInput {
	t.Helper()
	candidate, concept := splitTarget0067(t, target)
	id := fmt.Sprintf("naive-%04d", sequence+1)
	input.Projections = append(input.Projections, projection0067(id, kind, domain.EvidenceOriginRecorded, "origin-"+id))
	input.Observations = append(input.Observations, domain.IdeaRadarObservation{ID: "obs-" + id, ProjectionID: id, CandidateID: candidate, ConceptID: concept, Stance: domain.IdeaRadarStanceSupport})
	return input
}

func missingKinds0067(present []domain.IdeaProjectionKind, count int) []domain.IdeaProjectionKind {
	order := []domain.IdeaProjectionKind{domain.IdeaProjectionKindProblem, domain.IdeaProjectionKindBehavior, domain.IdeaProjectionKindArchitecture, domain.IdeaProjectionKindConstraint, domain.IdeaProjectionKindOperation}
	found := make(map[domain.IdeaProjectionKind]bool, len(present))
	for _, kind := range present {
		found[kind] = true
	}
	result := make([]domain.IdeaProjectionKind, 0, count)
	for _, kind := range order {
		if !found[kind] && len(result) < count {
			result = append(result, kind)
		}
	}
	return result
}

func combinedCorpusInput0067(t *testing.T) domain.IdeaRadarInput {
	t.Helper()
	combined := domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion, Projections: []domain.IdeaRadarProjection{}, Observations: []domain.IdeaRadarObservation{}}
	for _, tc := range frozenPlanCases0067() {
		if tc.id == "c10" {
			continue
		}
		for _, projection := range tc.input.Projections {
			old := projection.ID
			projection.ID = tc.id + "-" + old
			combined.Projections = append(combined.Projections, projection)
		}
		for _, observation := range tc.input.Observations {
			observation.ID = tc.id + "-" + observation.ID
			observation.ProjectionID = tc.id + "-" + observation.ProjectionID
			combined.Observations = append(combined.Observations, observation)
		}
	}
	return combined
}

func cloneInput0067(t *testing.T, input domain.IdeaRadarInput) domain.IdeaRadarInput {
	t.Helper()
	encoded := mustJSON0067(t, input)
	var result domain.IdeaRadarInput
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatalf("clone input: %v", err)
	}
	return result
}

func cloneObservations0067(input []domain.IdeaRadarObservation) []domain.IdeaRadarObservation {
	result := make([]domain.IdeaRadarObservation, len(input))
	copy(result, input)
	return result
}

func mustJSON0067(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return encoded
}

func removeObservations0067(input []domain.IdeaRadarObservation, ids []string) []domain.IdeaRadarObservation {
	result := make([]domain.IdeaRadarObservation, 0, len(input))
	for _, observation := range input {
		if !contains0067(ids, observation.ID) {
			result = append(result, observation)
		}
	}
	return result
}

func contains0067(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

func splitTarget0067(t *testing.T, target string) (string, string) {
	t.Helper()
	parts := strings.Split(target, "#")
	if len(parts) != 2 {
		t.Fatalf("invalid target ID %q", target)
	}
	return parts[0], parts[1]
}

func rotateProjections0067(values []domain.IdeaRadarProjection, n int) {
	copy(values, append(append([]domain.IdeaRadarProjection{}, values[n:]...), values[:n]...))
}

func rotateObservations0067(values []domain.IdeaRadarObservation, n int) {
	copy(values, append(append([]domain.IdeaRadarObservation{}, values[n:]...), values[:n]...))
}

func reverseProjections0067(values []domain.IdeaRadarProjection) {
	for left, right := 0, len(values)-1; left < right; left, right = left+1, right-1 {
		values[left], values[right] = values[right], values[left]
	}
}

func reverseObservations0067(values []domain.IdeaRadarObservation) {
	for left, right := 0, len(values)-1; left < right; left, right = left+1, right-1 {
		values[left], values[right] = values[right], values[left]
	}
}

func utilityVerdict0067(targets, plannerResolved, naiveResolved, plannerProbes, naiveProbes, improvedCases, improvedFamilies int, unknown, forbidden bool) string {
	if targets == 0 || naiveProbes == 0 || targets < 2 || improvedFamilies < 0 || unknown || forbidden {
		return "INCONCLUSIVE"
	}
	plannerRate := float64(plannerResolved) / float64(targets)
	naiveRate := float64(naiveResolved) / float64(targets)
	reduction := float64(naiveProbes-plannerProbes) / float64(naiveProbes)
	resolutionGain := plannerRate - naiveRate
	effortDelta := plannerProbes - naiveProbes
	if plannerRate >= naiveRate && (reduction >= 0.20 || resolutionGain >= 0.10 && effortDelta <= 0) && improvedCases >= 2 && improvedFamilies >= 2 {
		return "OBSERVED"
	}
	return "NOT_OBSERVED"
}
