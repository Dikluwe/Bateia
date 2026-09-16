/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_radar_test.md
 * @layer L1
 * @updated 2026-09-11
 */

package usecase_test

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

const (
	inputSchema  = "idea-radar-input-v1"
	reportSchema = "idea-radar-report-v1"
)

type radarCase struct {
	id       string
	family   string
	input    domain.IdeaRadarInput
	expected string
	forbid   string
	outcomes map[string]string
	metrics  domain.IdeaRadarMetrics
}

func TestAnalyzeIdeaRadarFrozenCorpus(t *testing.T) {
	cases := frozenRadarCorpus()
	if len(cases) != 6 {
		t.Fatalf("corpus cases = %d, want 6", len(cases))
	}
	families := map[string]bool{}
	for _, tc := range cases {
		families[tc.family] = true
		t.Run(tc.id, func(t *testing.T) {
			before := mustJSON(t, tc.input)
			report := usecase.AnalyzeIdeaRadar(tc.input)
			assertRadarCompleteReport(t, report)
			if after := mustJSON(t, tc.input); string(after) != string(before) {
				t.Fatalf("input mutated\nbefore: %s\n after: %s", before, after)
			}
			if !reflect.DeepEqual(report.Metrics, tc.metrics) {
				t.Fatalf("metrics = %#v, want %#v", report.Metrics, tc.metrics)
			}
			if len(report.Assessments) != len(tc.outcomes) {
				t.Fatalf("assessments = %d, want %d: %#v", len(report.Assessments), len(tc.outcomes), report.Assessments)
			}
			previous := ""
			for _, assessment := range report.Assessments {
				wantOutcome, ok := tc.outcomes[assessment.TargetID]
				if !ok {
					t.Errorf("unexpected target %q", assessment.TargetID)
				}
				if string(assessment.Outcome) != wantOutcome {
					t.Errorf("%s outcome = %q, want %q", assessment.TargetID, assessment.Outcome, wantOutcome)
				}
				if assessment.TargetID < previous {
					t.Errorf("assessments not canonical: %q before %q", previous, assessment.TargetID)
				}
				previous = assessment.TargetID
				assertRadarCanonicalAssessment(t, assessment)
			}
		})
	}
	if len(families) != 3 {
		t.Fatalf("families = %v, want exactly 3", families)
	}
}

func TestAnalyzeIdeaRadarCanonicalDeduplication(t *testing.T) {
	report := usecase.AnalyzeIdeaRadar(frozenRadarCorpus()[0].input)
	if report.Metrics.Observations != 4 || report.Metrics.CanonicalObservations != 3 || report.Metrics.DuplicateObservations != 1 {
		t.Fatalf("dedup metrics = %#v", report.Metrics)
	}
	a := report.Assessments[0]
	want := []string{"obs-a", "obs-b-a", "obs-p"}
	if !reflect.DeepEqual(a.SupportingObservationIDs, want) {
		t.Fatalf("support IDs = %v, want retained lowest semantic ID %v", a.SupportingObservationIDs, want)
	}
	if a.MissingProjectionKinds != 0 || a.MissingOrigins != 0 {
		t.Fatalf("missing counts = kinds:%d origins:%d", a.MissingProjectionKinds, a.MissingOrigins)
	}
}

func TestAnalyzeIdeaRadarPrecedenceAndIndependentThresholds(t *testing.T) {
	checks := map[string]struct {
		outcome                      string
		reason                       string
		missingKinds, missingOrigins int
	}{
		"noise-cache#cache":             {"INSUFFICIENT_PERSPECTIVES", "projection_diversity_below_minimum", 1, 1},
		"noise-auth#acl":                {"UNKNOWN_ORIGIN", "origin_unrecorded", 0, 1},
		"cand-rbac#role-based-access":   {"INSUFFICIENT_ORIGINS", "distinct_origins_below_minimum", 0, 2},
		"noise-role#role-map":           {"CONFLICT", "conflicting_projection", 0, 2},
		"noise-sync#timestamp-copy":     {"INSUFFICIENT_ORIGINS", "distinct_origins_below_minimum", 0, 2},
		"cand-crdt#conflict-free-merge": {"AMBIGUOUS", "ambiguous_projection", 0, 0},
		"noise-merge#last-write-wins":   {"AMBIGUOUS", "ambiguous_projection", 0, 0},
	}
	seen := map[string]bool{}
	for _, tc := range frozenRadarCorpus() {
		for _, a := range usecase.AnalyzeIdeaRadar(tc.input).Assessments {
			want, ok := checks[a.TargetID]
			if !ok {
				continue
			}
			seen[a.TargetID] = true
			if string(a.Outcome) != want.outcome || a.Reason != want.reason || a.MissingProjectionKinds != want.missingKinds || a.MissingOrigins != want.missingOrigins {
				t.Errorf("%s = (%s,%s,%d,%d), want (%s,%s,%d,%d)", a.TargetID, a.Outcome, a.Reason, a.MissingProjectionKinds, a.MissingOrigins, want.outcome, want.reason, want.missingKinds, want.missingOrigins)
			}
		}
	}
	if len(seen) != len(checks) {
		t.Fatalf("threshold witnesses seen = %v, want all %v", seen, checks)
	}
}

func TestAnalyzeIdeaRadarTwentyPermutationsAndNilEmptyText(t *testing.T) {
	base := frozenRadarCorpus()[4].input
	want := mustJSON(t, usecase.AnalyzeIdeaRadar(base))
	for iteration := 0; iteration < 20; iteration++ {
		candidate := cloneInput(base)
		rotateProjections(candidate.Projections, iteration%len(candidate.Projections))
		rotateObservations(candidate.Observations, (iteration*5)%len(candidate.Observations))
		if iteration%2 == 1 {
			reverseProjections(candidate.Projections)
		}
		if iteration%3 == 1 {
			reverseObservations(candidate.Observations)
		}
		for i := range candidate.Projections {
			if (iteration+i)%2 == 0 {
				candidate.Projections[i].Actions = nil
				candidate.Projections[i].Objects = []string{}
				candidate.Projections[i].Required = nil
			} else {
				candidate.Projections[i].Actions = []string{}
				candidate.Projections[i].Objects = nil
				candidate.Projections[i].Required = []string{}
			}
		}
		got := mustJSON(t, usecase.AnalyzeIdeaRadar(candidate))
		if string(got) != string(want) {
			t.Fatalf("permutation %d differs\n got: %s\nwant: %s", iteration, got, want)
		}
	}
}

func TestAnalyzeIdeaRadarRejectsInvalidAndAmbiguousInputs(t *testing.T) {
	valid := simpleInput()
	tests := []struct {
		name   string
		input  domain.IdeaRadarInput
		reason string
	}{
		{"unknown schema", withSchema(valid, "idea-radar-input-v2"), "invalid_schema"},
		{"nil projections", domain.IdeaRadarInput{SchemaVersion: inputSchema, Observations: []domain.IdeaRadarObservation{}}, "no_projections"},
		{"empty projections", domain.IdeaRadarInput{SchemaVersion: inputSchema, Projections: []domain.IdeaRadarProjection{}, Observations: []domain.IdeaRadarObservation{{}}}, "no_projections"},
		{"nil observations", domain.IdeaRadarInput{SchemaVersion: inputSchema, Projections: valid.Projections}, "no_observations"},
		{"empty observations", domain.IdeaRadarInput{SchemaVersion: inputSchema, Projections: valid.Projections, Observations: []domain.IdeaRadarObservation{}}, "no_observations"},
		{"unicode ID", mutateInput(valid, func(in *domain.IdeaRadarInput) { in.Observations[0].ConceptID = "café" }), "invalid_input"},
		{"separator in ID", mutateInput(valid, func(in *domain.IdeaRadarInput) { in.Observations[0].CandidateID = "bad#id" }), "invalid_input"},
		{"origin surrounding whitespace", mutateInput(valid, func(in *domain.IdeaRadarInput) { in.Projections[0].OriginID = " Origin-A" }), "invalid_input"},
		{"derived origin state", mutateInput(valid, func(in *domain.IdeaRadarInput) { in.Projections[0].OriginStatus = domain.EvidenceOriginSame }), "invalid_input"},
		{"unknown origin with ID", mutateInput(valid, func(in *domain.IdeaRadarInput) { in.Projections[0].OriginStatus = domain.EvidenceOriginUnknown }), "invalid_input"},
		{"missing projection reference", mutateInput(valid, func(in *domain.IdeaRadarInput) { in.Observations[0].ProjectionID = "missing" }), "invalid_input"},
		{"contradictory duplicate ID", contradictoryObservationID(valid), "ambiguous_input"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			report := usecase.AnalyzeIdeaRadar(tc.input)
			assertRadarUnknownReport(t, report, tc.reason, false)
		})
	}
}

func TestAnalyzeIdeaRadarUsesExactASCIIIdentityWithoutFuzzyOrUnicodeFolding(t *testing.T) {
	in := domain.IdeaRadarInput{
		SchemaVersion: inputSchema,
		Projections: []domain.IdeaRadarProjection{
			projection(" P ", "PROBLEM", "Origin-A"),
			projection("B", "BEHAVIOR", "Origin-B"),
			projection("a", "ARCHITECTURE", "Origin-C"),
		},
		Observations: []domain.IdeaRadarObservation{
			observation(" O-1 ", "p", " Candidate ", " Cache ", "SUPPORT"),
			observation("o-2", "b", "candidate", "caches", "SUPPORT"),
			observation("o-3", "A", "CANDIDATE", "CACHE", "SUPPORT"),
		},
	}
	report := usecase.AnalyzeIdeaRadar(in)
	assertRadarCompleteReport(t, report)
	if len(report.Assessments) != 2 || report.Assessments[0].TargetID != "candidate#cache" || report.Assessments[1].TargetID != "candidate#caches" {
		t.Fatalf("exact targets = %#v", report.Assessments)
	}
	if report.Assessments[0].Outcome != domain.IdeaRadarOutcome("INSUFFICIENT_PERSPECTIVES") || report.Assessments[0].MissingProjectionKinds != 1 {
		t.Fatalf("cache assessment = %#v", report.Assessments[0])
	}
}

func TestAnalyzeIdeaRadarInclusiveAndExclusiveLimits(t *testing.T) {
	assertRadarCompleteReport(t, usecase.AnalyzeIdeaRadar(projectionCountInput(64)))
	assertRadarUnknownReport(t, usecase.AnalyzeIdeaRadar(projectionCountInput(65)), "budget_exceeded", true)
	assertRadarCompleteReport(t, usecase.AnalyzeIdeaRadar(observationCountInput(1024)))
	assertRadarUnknownReport(t, usecase.AnalyzeIdeaRadar(observationCountInput(1025)), "budget_exceeded", true)
	assertRadarCompleteReport(t, usecase.AnalyzeIdeaRadar(targetCountInput(512)))
	assertRadarUnknownReport(t, usecase.AnalyzeIdeaRadar(targetCountInput(513)), "budget_exceeded", true)

	validID := mutateInput(simpleInput(), func(in *domain.IdeaRadarInput) {
		in.Projections[0].ID = strings.Repeat("a", 128)
		in.Observations[0].ProjectionID = strings.Repeat("a", 128)
	})
	assertRadarCompleteReport(t, usecase.AnalyzeIdeaRadar(validID))
	invalidID := mutateInput(validID, func(in *domain.IdeaRadarInput) { in.Projections[0].ID += "a"; in.Observations[0].ProjectionID += "a" })
	assertRadarUnknownReport(t, usecase.AnalyzeIdeaRadar(invalidID), "budget_exceeded", true)
	validOrigin := mutateInput(simpleInput(), func(in *domain.IdeaRadarInput) { in.Projections[0].OriginID = strings.Repeat("O", 128) })
	assertRadarCompleteReport(t, usecase.AnalyzeIdeaRadar(validOrigin))
	invalidOrigin := mutateInput(validOrigin, func(in *domain.IdeaRadarInput) { in.Projections[0].OriginID += "O" })
	assertRadarUnknownReport(t, usecase.AnalyzeIdeaRadar(invalidOrigin), "budget_exceeded", true)

	for _, dimension := range []string{"terms", "actions", "objects", "required"} {
		t.Run(dimension+" items", func(t *testing.T) {
			assertRadarCompleteReport(t, usecase.AnalyzeIdeaRadar(textDimensionInput(dimension, 64, 64)))
			assertRadarUnknownReport(t, usecase.AnalyzeIdeaRadar(textDimensionInput(dimension, 65, 1)), "budget_exceeded", true)
		})
		t.Run(dimension+" bytes", func(t *testing.T) {
			assertRadarCompleteReport(t, usecase.AnalyzeIdeaRadar(textDimensionInput(dimension, 64, 64)))
			assertRadarUnknownReport(t, usecase.AnalyzeIdeaRadar(textDimensionInput(dimension, 64, 65)), "budget_exceeded", true)
		})
	}
}

func TestAnalyzeIdeaRadarLogicalOperationsAreExactAndBounded(t *testing.T) {
	tests := []struct {
		name string
		in   domain.IdeaRadarInput
		want int
	}{
		{"semantic duplicate associates once", frozenRadarCorpus()[0].input, 10},
		{"single target", simpleInput(), 3},
		{"maximum observations", observationCountInput(1024), 1026},
		{"maximum targets", targetCountInput(512), 1025},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			report := usecase.AnalyzeIdeaRadar(tc.in)
			assertRadarCompleteReport(t, report)
			if report.Metrics.Operations != tc.want {
				t.Fatalf("operations = %d, want %d", report.Metrics.Operations, tc.want)
			}
			if report.Metrics.Operations > 4096 {
				t.Fatalf("operations exceeded logical budget: %d", report.Metrics.Operations)
			}
		})
	}
}

func TestExternalUtilityOracle(t *testing.T) {
	cases := frozenRadarCorpus()
	truth, forbidden, isolated := frozenUtilityData()
	bestKind, best := selectBestIsolated(cases, truth, forbidden, isolated)
	assertMetric(t, "baseline", best, score{TP: 3, FP: 1, FN: 3, Precision: ptr(0.75), Recall: ptr(0.5), F1: ptr(0.6), FalseConvergence: ptr(0.25), Contamination: ptr(0.25)})
	if bestKind != "ARCHITECTURE" {
		t.Fatalf("best isolated = %q, want ARCHITECTURE", bestKind)
	}

	radarPredictions := map[string][]string{}
	for _, tc := range cases {
		for _, a := range usecase.AnalyzeIdeaRadar(tc.input).Assessments {
			if string(a.Outcome) == "CONVERGENT" {
				radarPredictions[tc.id] = append(radarPredictions[tc.id], a.TargetID)
			}
		}
	}
	radar := microScore(cases, truth, forbidden, radarPredictions)
	assertMetric(t, "radar", radar, score{TP: 4, FP: 0, FN: 2, Precision: ptr(1), Recall: ptr(2.0 / 3), F1: ptr(0.8), FalseConvergence: ptr(0), Contamination: ptr(0)})
	if delta(radar.Precision, best.Precision) != 0.25 || !near(delta(radar.Recall, best.Recall), 1.0/6) || !near(delta(radar.F1, best.F1), 0.2) {
		t.Fatalf("deltas = precision %.12f recall %.12f f1 %.12f", delta(radar.Precision, best.Precision), delta(radar.Recall, best.Recall), delta(radar.F1, best.F1))
	}
	verdict, improvedCases, improvedFamilies := utilityVerdict(cases, truth, isolated[bestKind], radarPredictions, false)
	if verdict != "OBSERVED" || !reflect.DeepEqual(improvedCases, []string{"c02", "c05"}) || !reflect.DeepEqual(improvedFamilies, []string{"event-recovery", "offline-sync"}) {
		t.Fatalf("utility = %s cases=%v families=%v", verdict, improvedCases, improvedFamilies)
	}

	notObserved, _, _ := utilityVerdict(cases, truth, isolated[bestKind], isolated[bestKind], false)
	if notObserved != "NOT_OBSERVED" {
		t.Fatalf("unchanged valid vector = %s, want NOT_OBSERVED", notObserved)
	}
	inconclusive, _, _ := utilityVerdict(cases[:1], truth, isolated[bestKind], radarPredictions, false)
	if inconclusive != "INCONCLUSIVE" {
		t.Fatalf("one-case vector = %s, want INCONCLUSIVE", inconclusive)
	}
	budgetInconclusive, _, _ := utilityVerdict(cases, truth, isolated[bestKind], radarPredictions, true)
	if budgetInconclusive != "INCONCLUSIVE" {
		t.Fatalf("budget-exceeded vector = %s, want INCONCLUSIVE", budgetInconclusive)
	}
	undefined, _, _ := utilityVerdict(cases, truth, isolated[bestKind], map[string][]string{}, false)
	if undefined != "INCONCLUSIVE" {
		t.Fatalf("zero-prediction vector = %s, want INCONCLUSIVE", undefined)
	}
}

func frozenRadarCorpus() []radarCase {
	c01 := domain.IdeaRadarInput{SchemaVersion: inputSchema,
		Projections: []domain.IdeaRadarProjection{projection("p1", "PROBLEM", "or-a"), projection("b1", "BEHAVIOR", "or-b"), projection("a1", "ARCHITECTURE", "or-c")},
		Observations: []domain.IdeaRadarObservation{
			observation("obs-p", "p1", "cand-log", "event-log-recovery", "SUPPORT"),
			observation("obs-b-z", "b1", "cand-log", "event-log-recovery", "SUPPORT"),
			observation("obs-a", "a1", "cand-log", "event-log-recovery", "SUPPORT"),
			observation("obs-b-a", "b1", "cand-log", "event-log-recovery", "SUPPORT"),
		}}
	c02 := domain.IdeaRadarInput{SchemaVersion: inputSchema,
		Projections: []domain.IdeaRadarProjection{projection("b-main", "BEHAVIOR", "or-d"), projection("c-main", "CONSTRAINT", "or-e"), projection("o-main", "OPERATION", "or-f"), projection("p-noise", "PROBLEM", "or-x"), projection("b-noise", "BEHAVIOR", "or-y")},
		Observations: []domain.IdeaRadarObservation{
			observation("o-b", "b-main", "cand-journal", "write-ahead-log", "SUPPORT"), observation("o-c", "c-main", "cand-journal", "write-ahead-log", "SUPPORT"), observation("o-o", "o-main", "cand-journal", "write-ahead-log", "SUPPORT"),
			observation("o-np", "p-noise", "noise-cache", "cache", "SUPPORT"), observation("o-nb", "b-noise", "noise-cache", "cache", "SUPPORT"),
		}}
	c03 := domain.IdeaRadarInput{SchemaVersion: inputSchema,
		Projections: []domain.IdeaRadarProjection{projection(" P-MAIN ", "PROBLEM", "or-g"), projection("A-MAIN", "ARCHITECTURE", "or-h"), projection("c-main", "CONSTRAINT", "or-i"), projection("p-noise", "PROBLEM", "or-j"), unknownProjection("b-noise", "BEHAVIOR"), projection("a-noise", "ARCHITECTURE", "or-k")},
		Observations: []domain.IdeaRadarObservation{
			observation("O-P", "p-main", " CAND-POLICY ", "Policy-Enforcement", "SUPPORT"), observation("o-a", " a-main ", "cand-policy", "policy-enforcement", "SUPPORT"), observation("o-c", "C-MAIN", "cand-policy", "policy-enforcement", "SUPPORT"),
			observation("o-np", "p-noise", "noise-auth", "acl", "SUPPORT"), observation("o-nb", "b-noise", "noise-auth", "acl", "SUPPORT"), observation("o-na", "a-noise", "noise-auth", "acl", "SUPPORT"),
		}}
	shared := []domain.IdeaRadarProjection{projection("p", "PROBLEM", "shared"), projection("b", "BEHAVIOR", "shared"), projection("a", "ARCHITECTURE", "shared")}
	c04 := domain.IdeaRadarInput{SchemaVersion: inputSchema, Projections: shared,
		Observations: []domain.IdeaRadarObservation{
			observation("m-p", "p", "cand-rbac", "role-based-access", "SUPPORT"), observation("m-b", "b", "cand-rbac", "role-based-access", "SUPPORT"), observation("m-a", "a", "cand-rbac", "role-based-access", "SUPPORT"),
			observation("n-p", "p", "noise-role", "role-map", "SUPPORT"), observation("n-b", "b", "noise-role", "role-map", "SUPPORT"), observation("n-a", "a", "noise-role", "role-map", "SUPPORT"), observation("n-conflict", "a", "noise-role", "role-map", "CONFLICT"),
		}}
	c05 := domain.IdeaRadarInput{SchemaVersion: inputSchema,
		Projections: []domain.IdeaRadarProjection{projectionWithTerms("b-main", "BEHAVIOR", "or-l", []string{"operation log"}), projectionWithTerms("c-main", "CONSTRAINT", "or-m", []string{"operation log"}), projectionWithTerms("o-main", "OPERATION", "or-n", []string{"operation log"}), projectionWithTerms("p-file-1", "PROBLEM", "physical-same", []string{"timestamp"}), projectionWithTerms("a-file-2", "ARCHITECTURE", "physical-same", []string{"timestamp"}), projectionWithTerms("c-hash-3", "CONSTRAINT", "physical-same", []string{"timestamp"})},
		Observations: []domain.IdeaRadarObservation{
			observation("m-b", "b-main", "cand-sync", "operation-log-sync", "SUPPORT"), observation("m-c", "c-main", "cand-sync", "operation-log-sync", "SUPPORT"), observation("m-o", "o-main", "cand-sync", "operation-log-sync", "SUPPORT"),
			observation("n-p", "p-file-1", "noise-sync", "timestamp-copy", "SUPPORT"), observation("n-a", "a-file-2", "noise-sync", "timestamp-copy", "SUPPORT"), observation("n-c", "c-hash-3", "noise-sync", "timestamp-copy", "SUPPORT"),
		}}
	c06 := domain.IdeaRadarInput{SchemaVersion: inputSchema,
		Projections: []domain.IdeaRadarProjection{projection("p", "PROBLEM", "or-p"), projection("b", "BEHAVIOR", "or-q"), projection("a", "ARCHITECTURE", "or-r")},
		Observations: []domain.IdeaRadarObservation{
			observation("m-p", "p", "cand-crdt", "conflict-free-merge", "SUPPORT"), observation("m-b", "b", "cand-crdt", "conflict-free-merge", "SUPPORT"), observation("m-a", "a", "cand-crdt", "conflict-free-merge", "SUPPORT"), observation("m-amb", "a", "cand-crdt", "conflict-free-merge", "AMBIGUOUS"),
			observation("n-p", "p", "noise-merge", "last-write-wins", "SUPPORT"), observation("n-b", "b", "noise-merge", "last-write-wins", "SUPPORT"), observation("n-a", "a", "noise-merge", "last-write-wins", "SUPPORT"), observation("n-amb", "p", "noise-merge", "last-write-wins", "AMBIGUOUS"),
		}}
	return []radarCase{
		{"c01", "event-recovery", c01, "cand-log#event-log-recovery", "noise-log#snapshot-only", map[string]string{"cand-log#event-log-recovery": "CONVERGENT"}, metrics(3, 4, 3, 1, 1, 1, 0, 0, 0, 0, 0, 10)},
		{"c02", "event-recovery", c02, "cand-journal#write-ahead-log", "noise-cache#cache", map[string]string{"cand-journal#write-ahead-log": "CONVERGENT", "noise-cache#cache": "INSUFFICIENT_PERSPECTIVES"}, metrics(5, 5, 5, 0, 2, 1, 0, 0, 1, 0, 0, 15)},
		{"c03", "authorization", c03, "cand-policy#policy-enforcement", "noise-auth#acl", map[string]string{"cand-policy#policy-enforcement": "CONVERGENT", "noise-auth#acl": "UNKNOWN_ORIGIN"}, metrics(6, 6, 6, 0, 2, 1, 0, 0, 0, 1, 0, 18)},
		{"c04", "authorization", c04, "cand-rbac#role-based-access", "noise-role#role-map", map[string]string{"cand-rbac#role-based-access": "INSUFFICIENT_ORIGINS", "noise-role#role-map": "CONFLICT"}, metrics(3, 7, 7, 0, 2, 0, 1, 0, 0, 0, 1, 17)},
		{"c05", "offline-sync", c05, "cand-sync#operation-log-sync", "noise-sync#timestamp-copy", map[string]string{"cand-sync#operation-log-sync": "CONVERGENT", "noise-sync#timestamp-copy": "INSUFFICIENT_ORIGINS"}, metrics(6, 6, 6, 0, 2, 1, 0, 0, 0, 0, 1, 18)},
		{"c06", "offline-sync", c06, "cand-crdt#conflict-free-merge", "noise-merge#last-write-wins", map[string]string{"cand-crdt#conflict-free-merge": "AMBIGUOUS", "noise-merge#last-write-wins": "AMBIGUOUS"}, metrics(3, 8, 8, 0, 2, 0, 0, 2, 0, 0, 0, 19)},
	}
}

func projection(id, kind, origin string) domain.IdeaRadarProjection {
	return projectionWithTerms(id, kind, origin, nil)
}

func projectionWithTerms(id, kind, origin string, terms []string) domain.IdeaRadarProjection {
	return domain.IdeaRadarProjection{ID: id, Kind: domain.IdeaProjectionKind(kind), OriginStatus: domain.EvidenceOriginRecorded, OriginID: origin, Terms: terms}
}

func unknownProjection(id, kind string) domain.IdeaRadarProjection {
	return domain.IdeaRadarProjection{ID: id, Kind: domain.IdeaProjectionKind(kind), OriginStatus: domain.EvidenceOriginUnknown}
}

func observation(id, projectionID, candidateID, conceptID, stance string) domain.IdeaRadarObservation {
	return domain.IdeaRadarObservation{ID: id, ProjectionID: projectionID, CandidateID: candidateID, ConceptID: conceptID, Stance: domain.IdeaRadarStance(stance)}
}

func metrics(p, o, canonical, duplicate, targets, convergent, conflicts, ambiguous, insufficientPerspectives, unknownOrigins, insufficientOrigins, operations int) domain.IdeaRadarMetrics {
	return domain.IdeaRadarMetrics{Projections: p, Observations: o, CanonicalObservations: canonical, DuplicateObservations: duplicate, Targets: targets, Convergent: convergent, Conflicts: conflicts, Ambiguous: ambiguous, InsufficientPerspectives: insufficientPerspectives, UnknownOrigins: unknownOrigins, InsufficientOrigins: insufficientOrigins, Operations: operations, VerdictChanges: 0}
}

func assertRadarCompleteReport(t *testing.T, report domain.IdeaRadarReport) {
	t.Helper()
	if report.SchemaVersion != reportSchema || report.Status != "COMPLETE" || report.Reason != "analyzed" || report.Metrics.LimitExceeded || report.Metrics.VerdictChanges != 0 {
		t.Fatalf("unexpected report envelope: %#v", report)
	}
	if report.Assessments == nil {
		t.Fatal("assessments must be a non-null array")
	}
}

func assertRadarUnknownReport(t *testing.T, report domain.IdeaRadarReport, reason string, limit bool) {
	t.Helper()
	if report.SchemaVersion != reportSchema || report.Status != "UNKNOWN" || report.Reason != reason || report.Metrics.LimitExceeded != limit || report.Metrics.VerdictChanges != 0 {
		t.Fatalf("unexpected UNKNOWN envelope: %#v", report)
	}
	if report.Assessments == nil || len(report.Assessments) != 0 {
		t.Fatalf("UNKNOWN assessments must be empty non-null array: %#v", report.Assessments)
	}
}

func assertRadarCanonicalAssessment(t *testing.T, a domain.IdeaRadarAssessment) {
	t.Helper()
	if a.TargetID != a.CandidateID+"#"+a.ConceptID || a.CoalitionDraft.ClaimID != a.TargetID {
		t.Errorf("target/draft identity mismatch: %#v", a)
	}
	if a.CoalitionDraft.ReadyFor0065 || a.CoalitionDraft.Reason != "requires_evidence_binding" {
		t.Errorf("unsafe 0065 draft: %#v", a.CoalitionDraft)
	}
	wantRoles := []domain.EvidenceCoalitionRole{domain.EvidenceCoalitionRoleImplementation, domain.EvidenceCoalitionRoleIndependentTest}
	if !reflect.DeepEqual(a.CoalitionDraft.RequiredRoles, wantRoles) || !reflect.DeepEqual(a.CoalitionDraft.ObservationIDs, a.SupportingObservationIDs) {
		t.Errorf("draft bridge = %#v, support = %v", a.CoalitionDraft, a.SupportingObservationIDs)
	}
	if a.SupportingObservationIDs == nil || a.ConflictingObservationIDs == nil || a.AmbiguousObservationIDs == nil || a.UnknownOriginObservationIDs == nil || a.ProjectionKinds == nil || a.DistinctOrigins == nil || a.CoalitionDraft.ObservationIDs == nil || a.CoalitionDraft.RequiredRoles == nil {
		t.Errorf("assessment contains null array: %#v", a)
	}
	assertSortedStrings(t, a.SupportingObservationIDs)
	assertSortedStrings(t, a.ConflictingObservationIDs)
	assertSortedStrings(t, a.AmbiguousObservationIDs)
	assertSortedStrings(t, a.UnknownOriginObservationIDs)
	assertSortedStrings(t, a.DistinctOrigins)
	kinds := make([]string, len(a.ProjectionKinds))
	for i := range a.ProjectionKinds {
		kinds[i] = string(a.ProjectionKinds[i])
	}
	assertSortedStrings(t, kinds)
}

func assertSortedStrings(t *testing.T, values []string) {
	t.Helper()
	if !sort.StringsAreSorted(values) {
		t.Errorf("not canonically sorted: %v", values)
	}
}

func simpleInput() domain.IdeaRadarInput {
	return domain.IdeaRadarInput{SchemaVersion: inputSchema, Projections: []domain.IdeaRadarProjection{projection("p", "PROBLEM", "Origin-A")}, Observations: []domain.IdeaRadarObservation{observation("o", "p", "candidate", "concept", "SUPPORT")}}
}

func cloneInput(in domain.IdeaRadarInput) domain.IdeaRadarInput {
	out := in
	out.Projections = append([]domain.IdeaRadarProjection(nil), in.Projections...)
	for i := range out.Projections {
		out.Projections[i].Terms = append([]string(nil), in.Projections[i].Terms...)
		out.Projections[i].Actions = append([]string(nil), in.Projections[i].Actions...)
		out.Projections[i].Objects = append([]string(nil), in.Projections[i].Objects...)
		out.Projections[i].Required = append([]string(nil), in.Projections[i].Required...)
	}
	out.Observations = append([]domain.IdeaRadarObservation(nil), in.Observations...)
	return out
}

func mutateInput(in domain.IdeaRadarInput, mutate func(*domain.IdeaRadarInput)) domain.IdeaRadarInput {
	out := cloneInput(in)
	mutate(&out)
	return out
}
func withSchema(in domain.IdeaRadarInput, schema string) domain.IdeaRadarInput {
	return mutateInput(in, func(out *domain.IdeaRadarInput) { out.SchemaVersion = schema })
}
func contradictoryObservationID(in domain.IdeaRadarInput) domain.IdeaRadarInput {
	return mutateInput(in, func(out *domain.IdeaRadarInput) {
		duplicate := out.Observations[0]
		duplicate.Stance = domain.IdeaRadarStance("CONFLICT")
		out.Observations = append(out.Observations, duplicate)
	})
}
func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func rotateProjections(values []domain.IdeaRadarProjection, n int) {
	if len(values) == 0 {
		return
	}
	copy(values, append(append([]domain.IdeaRadarProjection{}, values[n:]...), values[:n]...))
}
func rotateObservations(values []domain.IdeaRadarObservation, n int) {
	if len(values) == 0 {
		return
	}
	copy(values, append(append([]domain.IdeaRadarObservation{}, values[n:]...), values[:n]...))
}
func reverseProjections(values []domain.IdeaRadarProjection) {
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}
func reverseObservations(values []domain.IdeaRadarObservation) {
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}

func projectionCountInput(count int) domain.IdeaRadarInput {
	in := simpleInput()
	in.Projections = make([]domain.IdeaRadarProjection, count)
	for i := range in.Projections {
		in.Projections[i] = projection(fmt.Sprintf("p-%03d", i), []string{"PROBLEM", "BEHAVIOR", "ARCHITECTURE", "CONSTRAINT", "OPERATION"}[i%5], fmt.Sprintf("Origin-%03d", i))
	}
	in.Observations[0].ProjectionID = "p-000"
	return in
}
func observationCountInput(count int) domain.IdeaRadarInput {
	in := simpleInput()
	in.Observations = make([]domain.IdeaRadarObservation, count)
	for i := range in.Observations {
		in.Observations[i] = observation(fmt.Sprintf("o-%04d", i), "p", "candidate", "concept", "SUPPORT")
	}
	return in
}
func targetCountInput(count int) domain.IdeaRadarInput {
	in := simpleInput()
	in.Observations = make([]domain.IdeaRadarObservation, count)
	for i := range in.Observations {
		in.Observations[i] = observation(fmt.Sprintf("o-%04d", i), "p", fmt.Sprintf("candidate-%04d", i), "concept", "SUPPORT")
	}
	return in
}
func textDimensionInput(dimension string, items, bytesPerItem int) domain.IdeaRadarInput {
	in := simpleInput()
	values := make([]string, items)
	for i := range values {
		values[i] = strings.Repeat("x", bytesPerItem)
	}
	switch dimension {
	case "terms":
		in.Projections[0].Terms = values
	case "actions":
		in.Projections[0].Actions = values
	case "objects":
		in.Projections[0].Objects = values
	case "required":
		in.Projections[0].Required = values
	}
	return in
}

type score struct {
	TP, FP, FN                                             int
	Precision, Recall, F1, FalseConvergence, Contamination *float64
}

func frozenUtilityData() (map[string][]string, map[string][]string, map[string]map[string][]string) {
	truth := map[string][]string{"c01": {"cand-log#event-log-recovery"}, "c02": {"cand-journal#write-ahead-log"}, "c03": {"cand-policy#policy-enforcement"}, "c04": {"cand-rbac#role-based-access"}, "c05": {"cand-sync#operation-log-sync"}, "c06": {"cand-crdt#conflict-free-merge"}}
	forbidden := map[string][]string{"c01": {"noise-log#snapshot-only"}, "c02": {"noise-cache#cache"}, "c03": {"noise-auth#acl"}, "c04": {"noise-role#role-map"}, "c05": {"noise-sync#timestamp-copy"}, "c06": {"noise-merge#last-write-wins"}}
	isolated := map[string]map[string][]string{
		"PROBLEM":      {"c01": {"cand-log#event-log-recovery"}, "c02": {"noise-cache#cache"}, "c03": {"cand-policy#policy-enforcement"}, "c04": {}, "c05": {"noise-sync#timestamp-copy"}, "c06": {}},
		"BEHAVIOR":     {"c01": {"cand-log#event-log-recovery"}, "c02": {"cand-journal#write-ahead-log"}, "c03": {}, "c04": {"noise-role#role-map"}, "c05": {"cand-sync#operation-log-sync"}, "c06": {"noise-merge#last-write-wins"}},
		"ARCHITECTURE": {"c01": {"cand-log#event-log-recovery"}, "c02": {}, "c03": {"cand-policy#policy-enforcement"}, "c04": {"noise-role#role-map"}, "c05": {}, "c06": {"cand-crdt#conflict-free-merge"}},
		"CONSTRAINT":   {"c01": {"noise-log#snapshot-only"}, "c02": {}, "c03": {"cand-policy#policy-enforcement"}, "c04": {"cand-rbac#role-based-access"}, "c05": {}, "c06": {}},
		"OPERATION":    {"c01": {"cand-log#event-log-recovery"}, "c02": {"cand-journal#write-ahead-log"}, "c03": {}, "c04": {}, "c05": {"cand-sync#operation-log-sync"}, "c06": {"noise-merge#last-write-wins"}},
	}
	return truth, forbidden, isolated
}

func microScore(cases []radarCase, truth, forbidden map[string][]string, predictions map[string][]string) score {
	result := score{}
	for _, tc := range cases {
		want, predicted := set(truth[tc.id]), set(predictions[tc.id])
		for target := range predicted {
			if want[target] {
				result.TP++
			} else {
				result.FP++
			}
		}
		for target := range want {
			if !predicted[target] {
				result.FN++
			}
		}
	}
	result.Precision = ratio(result.TP, result.TP+result.FP)
	result.Recall = ratio(result.TP, result.TP+result.FN)
	result.F1 = ratio(2*result.TP, 2*result.TP+result.FP+result.FN)
	result.FalseConvergence = ratio(result.FP, result.TP+result.FP)
	forbiddenPredicted, predictedTotal := 0, 0
	for _, tc := range cases {
		forbiddenSet := set(forbidden[tc.id])
		for target := range set(predictions[tc.id]) {
			predictedTotal++
			if forbiddenSet[target] {
				forbiddenPredicted++
			}
		}
	}
	result.Contamination = ratio(forbiddenPredicted, predictedTotal)
	return result
}

func selectBestIsolated(cases []radarCase, truth, forbidden map[string][]string, isolated map[string]map[string][]string) (string, score) {
	kinds := make([]string, 0, len(isolated))
	for kind := range isolated {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	bestKind := ""
	var best score
	for _, kind := range kinds {
		candidate := microScore(cases, truth, forbidden, isolated[kind])
		if bestKind == "" || scoreGreater(candidate, best) {
			bestKind, best = kind, candidate
		}
	}
	return bestKind, best
}
func scoreGreater(a, b score) bool {
	for _, pair := range [][2]*float64{{a.F1, b.F1}, {a.Precision, b.Precision}, {a.Recall, b.Recall}} {
		av, bv := numeric(pair[0]), numeric(pair[1])
		if !near(av, bv) {
			return av > bv
		}
	}
	return false
}

func utilityVerdict(cases []radarCase, truth map[string][]string, baseline, radar map[string][]string, budgetExceeded bool) (string, []string, []string) {
	_, forbidden, _ := frozenUtilityData()
	baselineScore := microScore(cases, truth, forbidden, baseline)
	radarScore := microScore(cases, truth, forbidden, radar)
	if budgetExceeded || len(cases) < 2 || baselineScore.Precision == nil || baselineScore.Recall == nil || baselineScore.F1 == nil || radarScore.Precision == nil || radarScore.Recall == nil || radarScore.F1 == nil {
		return "INCONCLUSIVE", nil, nil
	}
	improvedCases, familySet := []string{}, map[string]bool{}
	for _, tc := range cases {
		base := microScore([]radarCase{tc}, truth, forbidden, baseline)
		candidate := microScore([]radarCase{tc}, truth, forbidden, radar)
		if candidate.F1 != nil && base.F1 != nil && *candidate.F1 > *base.F1 {
			improvedCases = append(improvedCases, tc.id)
			familySet[tc.family] = true
		}
	}
	families := make([]string, 0, len(familySet))
	for family := range familySet {
		families = append(families, family)
	}
	sort.Strings(families)
	observed := *radarScore.Precision >= *baselineScore.Precision && (delta(radarScore.Recall, baselineScore.Recall) >= 0.1 || delta(radarScore.F1, baselineScore.F1) >= 0.1) && len(improvedCases) >= 2 && len(families) >= 2
	if observed {
		return "OBSERVED", improvedCases, families
	}
	return "NOT_OBSERVED", improvedCases, families
}

func set(values []string) map[string]bool {
	result := map[string]bool{}
	for _, value := range values {
		result[value] = true
	}
	return result
}
func ratio(numerator, denominator int) *float64 {
	if denominator == 0 {
		return nil
	}
	value := float64(numerator) / float64(denominator)
	return &value
}
func ptr(value float64) *float64 { return &value }
func numeric(value *float64) float64 {
	if value == nil {
		return math.Inf(-1)
	}
	return *value
}
func delta(a, b *float64) float64 {
	if a == nil || b == nil {
		return math.NaN()
	}
	return *a - *b
}
func near(a, b float64) bool { return math.Abs(a-b) < 1e-12 }
func assertMetric(t *testing.T, label string, got, want score) {
	t.Helper()
	if got.TP != want.TP || got.FP != want.FP || got.FN != want.FN || !sameFloat(got.Precision, want.Precision) || !sameFloat(got.Recall, want.Recall) || !sameFloat(got.F1, want.F1) || !sameFloat(got.FalseConvergence, want.FalseConvergence) || !sameFloat(got.Contamination, want.Contamination) {
		t.Fatalf("%s = %#v, want %#v", label, got, want)
	}
}
func sameFloat(a, b *float64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return near(*a, *b)
}
