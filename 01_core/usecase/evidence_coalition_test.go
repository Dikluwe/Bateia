/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_coalition_test.md
 * @layer L1
 * @updated 2026-09-11
 */

package usecase_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

const (
	roleImplementation domain.EvidenceCoalitionRole     = "IMPLEMENTATION"
	roleTest           domain.EvidenceCoalitionRole     = "INDEPENDENT_TEST"
	roleExecution      domain.EvidenceCoalitionRole     = "EXECUTION"
	roleDocumentation  domain.EvidenceCoalitionRole     = "DOCUMENTATION"
	support            domain.EvidenceCoalitionPolarity = "SUPPORT"
	conflict           domain.EvidenceCoalitionPolarity = "CONFLICT"
	recorded           domain.EvidenceOriginStatus      = "RECORDED"
	unknownOrigin      domain.EvidenceOriginStatus      = "UNKNOWN_ORIGIN"
)

func source(key string) domain.EvidenceSource {
	return domain.EvidenceSource{Kind: domain.EvidenceFile, Path: key + ".go", Symbol: key, RuleID: "source-rule", ContentHash: "hash-" + key}
}

func witness(id string, role domain.EvidenceCoalitionRole, polarity domain.EvidenceCoalitionPolarity, originStatus domain.EvidenceOriginStatus, originID, sourceKey string) domain.EvidenceCoalitionWitness {
	return domain.EvidenceCoalitionWitness{ID: id, Role: role, Polarity: polarity, OriginStatus: originStatus, OriginID: originID, Source: source(sourceKey)}
}

func input(ruleID, claimID string, status domain.EvidenceStatus, roles []domain.EvidenceCoalitionRole, minimum int, witnesses []domain.EvidenceCoalitionWitness) domain.EvidenceCoalitionInput {
	return domain.EvidenceCoalitionInput{
		SchemaVersion: "evidence-coalition-input-v1",
		Rules:         []domain.EvidenceCoalitionRule{{ID: ruleID, RequiredRoles: roles, MinimumDistinctOrigins: minimum}},
		Claims:        []domain.EvidenceCoalitionClaim{{ID: claimID, RuleID: ruleID, BaselineStatus: status, Witnesses: witnesses}},
	}
}

func TestAnalyzeEvidenceCoalitionsContractCorpus(t *testing.T) {
	t.Parallel()

	duplicateA := witness("z-copy", roleImplementation, support, recorded, "A", "same-physical-source")
	duplicateB := duplicateA
	duplicateB.ID = "a-canonical"
	contradictory := duplicateB
	contradictory.ID = "conflicting-copy"
	contradictory.Polarity = conflict

	tests := []struct {
		name         string
		input        domain.EvidenceCoalitionInput
		outcome      string
		independence string
		reason       string
		supportIDs   []string
		conflictIDs  []string
		missing      []domain.EvidenceCoalitionRole
		origins      []string
		operations   int
		wantInvalid  int
	}{
		{"P01 distinct origins", input("r1", "c1", domain.EvidenceProven, []domain.EvidenceCoalitionRole{roleImplementation, roleTest}, 2, []domain.EvidenceCoalitionWitness{witness("impl", roleImplementation, support, recorded, "A", "impl"), witness("test", roleTest, support, recorded, "B", "test")}), "SATISFIED", "DISTINCT_ORIGIN", "coalition_satisfied", []string{"impl", "test"}, []string{}, []domain.EvidenceCoalitionRole{}, []string{"A", "B"}, 4, 0},
		{"P02 threshold one", input("r2", "c2", domain.EvidencePartial, []domain.EvidenceCoalitionRole{roleImplementation}, 1, []domain.EvidenceCoalitionWitness{witness("impl", roleImplementation, support, recorded, "A", "impl")}), "SATISFIED", "RECORDED", "coalition_satisfied", []string{"impl"}, []string{}, []domain.EvidenceCoalitionRole{}, []string{"A"}, 2, 0},
		{"N01 same origin", input("r3", "c3", domain.EvidenceAbsent, []domain.EvidenceCoalitionRole{roleImplementation, roleTest}, 2, []domain.EvidenceCoalitionWitness{witness("impl", roleImplementation, support, recorded, "A", "impl"), witness("test", roleTest, support, recorded, "A", "test")}), "INSUFFICIENT_INDEPENDENCE", "SAME_ORIGIN", "distinct_origins_below_minimum", []string{"impl", "test"}, []string{}, []domain.EvidenceCoalitionRole{}, []string{"A"}, 4, 0},
		{"O01 unknown origin", input("r4", "c4", domain.EvidenceInconclusive, []domain.EvidenceCoalitionRole{roleImplementation, roleTest}, 2, []domain.EvidenceCoalitionWitness{witness("impl", roleImplementation, support, recorded, "A", "impl"), witness("test", roleTest, support, unknownOrigin, "", "different-path-hash-and-id")}), "UNKNOWN_INDEPENDENCE", "UNKNOWN_ORIGIN", "origin_unrecorded", []string{"impl", "test"}, []string{}, []domain.EvidenceCoalitionRole{}, []string{"A"}, 4, 0},
		{"N02 missing implementation", input("r5", "c5", domain.EvidenceConflict, []domain.EvidenceCoalitionRole{roleImplementation, roleTest}, 2, []domain.EvidenceCoalitionWitness{witness("test", roleTest, support, recorded, "B", "test")}), "MISSING_ROLES", "UNKNOWN_ORIGIN", "missing_required_roles", []string{"test"}, []string{}, []domain.EvidenceCoalitionRole{roleImplementation}, []string{"B"}, 3, 0},
		{"N03 documentation is not implementation", input("r6", "c6", domain.EvidenceProven, []domain.EvidenceCoalitionRole{roleImplementation, roleTest}, 2, []domain.EvidenceCoalitionWitness{witness("docs", roleDocumentation, support, recorded, "A", "docs"), witness("test", roleTest, support, recorded, "B", "test")}), "MISSING_ROLES", "UNKNOWN_ORIGIN", "missing_required_roles", []string{"docs", "test"}, []string{}, []domain.EvidenceCoalitionRole{roleImplementation}, []string{"B"}, 4, 0},
		{"N04 conflict prevails", input("r7", "c7", domain.EvidencePartial, []domain.EvidenceCoalitionRole{roleImplementation, roleTest}, 2, []domain.EvidenceCoalitionWitness{witness("impl", roleImplementation, support, recorded, "A", "impl"), witness("test", roleTest, support, recorded, "B", "test"), witness("counter", roleExecution, conflict, unknownOrigin, "", "counter")}), "CONFLICT", "DISTINCT_ORIGIN", "conflicting_evidence", []string{"impl", "test"}, []string{"counter"}, []domain.EvidenceCoalitionRole{}, []string{"A", "B"}, 5, 0},
		{"P03 identical duplicate", input("r8", "c8", domain.EvidenceAbsent, []domain.EvidenceCoalitionRole{roleImplementation}, 1, []domain.EvidenceCoalitionWitness{duplicateA, duplicateB}), "SATISFIED", "RECORDED", "coalition_satisfied", []string{"a-canonical"}, []string{}, []domain.EvidenceCoalitionRole{}, []string{"A"}, 2, 0},
		{"N05 contradictory duplicate", input("r9", "c9", domain.EvidenceInconclusive, []domain.EvidenceCoalitionRole{roleImplementation}, 1, []domain.EvidenceCoalitionWitness{duplicateA, contradictory}), "INVALID_INPUT", "UNKNOWN_ORIGIN", "invalid_claim", []string{}, []string{}, []domain.EvidenceCoalitionRole{}, []string{}, 2, 1},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			report := usecase.AnalyzeEvidenceCoalitions(tc.input)
			if report.SchemaVersion != "evidence-coalition-report-v1" || report.Status != "COMPLETE" {
				t.Fatalf("report header = %q/%q, want evidence-coalition-report-v1/COMPLETE", report.SchemaVersion, report.Status)
			}
			wantReportReason := "analyzed"
			if tc.wantInvalid != 0 {
				wantReportReason = "analyzed_with_invalid_claims"
			}
			if report.Reason != wantReportReason {
				t.Errorf("report reason = %q, want %q", report.Reason, wantReportReason)
			}
			if len(report.Assessments) != 1 {
				t.Fatalf("got %d assessments, want 1", len(report.Assessments))
			}
			a := report.Assessments[0]
			if string(a.Outcome) != tc.outcome || string(a.Independence) != tc.independence || a.Reason != tc.reason {
				t.Errorf("classification = %q/%q/%q, want %q/%q/%q", a.Outcome, a.Independence, a.Reason, tc.outcome, tc.independence, tc.reason)
			}
			if a.BaselineStatus != tc.input.Claims[0].BaselineStatus || a.DecisionImpact != "NONE" {
				t.Errorf("shadow analysis changed decision: baseline=%q impact=%q", a.BaselineStatus, a.DecisionImpact)
			}
			if !reflect.DeepEqual(a.SupportingWitnessIDs, tc.supportIDs) || !reflect.DeepEqual(a.ConflictingWitnessIDs, tc.conflictIDs) || !reflect.DeepEqual(a.MissingRoles, tc.missing) || !reflect.DeepEqual(a.DistinctOrigins, tc.origins) {
				t.Errorf("assessment details = supports %#v conflicts %#v missing %#v origins %#v", a.SupportingWitnessIDs, a.ConflictingWitnessIDs, a.MissingRoles, a.DistinctOrigins)
			}
			if report.Metrics.Operations != tc.operations || report.Metrics.Invalid != tc.wantInvalid || report.Metrics.VerdictChanges != 0 || report.Metrics.LimitExceeded {
				t.Errorf("metrics = %+v, want operations=%d invalid=%d verdict_changes=0 limit_exceeded=false", report.Metrics, tc.operations, tc.wantInvalid)
			}
			assertNonNilSlices(t, report)
		})
	}
}

func TestAnalyzeEvidenceCoalitionsRejectsOpaqueGlobalInputs(t *testing.T) {
	t.Parallel()

	validClaim := domain.EvidenceCoalitionClaim{ID: "claim", RuleID: "rule", BaselineStatus: domain.EvidenceProven, Witnesses: []domain.EvidenceCoalitionWitness{}}
	tests := []struct {
		name, schema, reason string
		rules                []domain.EvidenceCoalitionRule
		claims               []domain.EvidenceCoalitionClaim
	}{
		{"unknown schema", "future-schema", "invalid_schema", []domain.EvidenceCoalitionRule{}, []domain.EvidenceCoalitionClaim{validClaim}},
		{"empty claims", "evidence-coalition-input-v1", "no_claims", []domain.EvidenceCoalitionRule{}, []domain.EvidenceCoalitionClaim{}},
		{"nil claims", "evidence-coalition-input-v1", "no_claims", nil, nil},
		{"duplicate rule IDs", "evidence-coalition-input-v1", "ambiguous_input", []domain.EvidenceCoalitionRule{{ID: "rule", RequiredRoles: []domain.EvidenceCoalitionRole{roleImplementation}, MinimumDistinctOrigins: 1}, {ID: "rule", RequiredRoles: []domain.EvidenceCoalitionRole{roleTest}, MinimumDistinctOrigins: 1}}, []domain.EvidenceCoalitionClaim{validClaim}},
		{"duplicate claim IDs", "evidence-coalition-input-v1", "ambiguous_input", []domain.EvidenceCoalitionRule{{ID: "rule", RequiredRoles: []domain.EvidenceCoalitionRole{roleImplementation}, MinimumDistinctOrigins: 1}}, []domain.EvidenceCoalitionClaim{validClaim, validClaim}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			report := usecase.AnalyzeEvidenceCoalitions(domain.EvidenceCoalitionInput{SchemaVersion: tc.schema, Rules: tc.rules, Claims: tc.claims})
			assertUnknownReport(t, report, tc.reason)
		})
	}
}

func TestAnalyzeEvidenceCoalitionsCanonicalAcrossTwentyPermutationsAndNilEmpty(t *testing.T) {
	t.Parallel()

	rules := []domain.EvidenceCoalitionRule{
		{ID: "rule-b", RequiredRoles: []domain.EvidenceCoalitionRole{roleTest, roleImplementation}, MinimumDistinctOrigins: 2},
		{ID: "rule-a", RequiredRoles: []domain.EvidenceCoalitionRole{roleExecution}, MinimumDistinctOrigins: 1},
	}
	claims := []domain.EvidenceCoalitionClaim{
		{ID: "claim-b", RuleID: "rule-b", BaselineStatus: domain.EvidencePartial, Witnesses: []domain.EvidenceCoalitionWitness{witness("z", roleTest, support, recorded, "B", "z"), witness("a", roleImplementation, support, recorded, "A", "a")}},
		{ID: "claim-a", RuleID: "rule-a", BaselineStatus: domain.EvidenceConflict, Witnesses: []domain.EvidenceCoalitionWitness{witness("x", roleExecution, support, recorded, "X", "x")}},
	}
	canonicalJSON := marshalReport(t, usecase.AnalyzeEvidenceCoalitions(domain.EvidenceCoalitionInput{SchemaVersion: "evidence-coalition-input-v1", Rules: rules, Claims: claims}))
	for permutation := 0; permutation < 20; permutation++ {
		permutedRules := append([]domain.EvidenceCoalitionRule(nil), rules...)
		permutedClaims := append([]domain.EvidenceCoalitionClaim(nil), claims...)
		if permutation%2 == 0 {
			permutedRules[0], permutedRules[1] = permutedRules[1], permutedRules[0]
		}
		if permutation%3 == 0 {
			permutedClaims[0], permutedClaims[1] = permutedClaims[1], permutedClaims[0]
		}
		for i := range permutedRules {
			permutedRules[i].RequiredRoles = append([]domain.EvidenceCoalitionRole(nil), permutedRules[i].RequiredRoles...)
			if permutation%5 == 0 && len(permutedRules[i].RequiredRoles) == 2 {
				permutedRules[i].RequiredRoles[0], permutedRules[i].RequiredRoles[1] = permutedRules[i].RequiredRoles[1], permutedRules[i].RequiredRoles[0]
			}
		}
		for i := range permutedClaims {
			permutedClaims[i].Witnesses = append([]domain.EvidenceCoalitionWitness(nil), permutedClaims[i].Witnesses...)
			if permutation%7 == 0 && len(permutedClaims[i].Witnesses) == 2 {
				permutedClaims[i].Witnesses[0], permutedClaims[i].Witnesses[1] = permutedClaims[i].Witnesses[1], permutedClaims[i].Witnesses[0]
			}
		}
		got := marshalReport(t, usecase.AnalyzeEvidenceCoalitions(domain.EvidenceCoalitionInput{SchemaVersion: "evidence-coalition-input-v1", Rules: permutedRules, Claims: permutedClaims}))
		if string(got) != string(canonicalJSON) {
			t.Fatalf("permutation %d changed JSON\nbase: %s\ngot:  %s", permutation, canonicalJSON, got)
		}
	}

	nilInput := input("rule", "claim", domain.EvidenceAbsent, []domain.EvidenceCoalitionRole{roleImplementation}, 1, nil)
	emptyInput := input("rule", "claim", domain.EvidenceAbsent, []domain.EvidenceCoalitionRole{roleImplementation}, 1, []domain.EvidenceCoalitionWitness{})
	if got, want := string(marshalReport(t, usecase.AnalyzeEvidenceCoalitions(nilInput))), string(marshalReport(t, usecase.AnalyzeEvidenceCoalitions(emptyInput))); got != want {
		t.Fatalf("nil and empty witness slices differ\nnil:   %s\nempty: %s", got, want)
	}
}

func TestAnalyzeEvidenceCoalitionsDoesNotMutateInput(t *testing.T) {
	t.Parallel()
	in := input("rule", "claim", domain.EvidenceProven, []domain.EvidenceCoalitionRole{roleTest, roleImplementation}, 2, []domain.EvidenceCoalitionWitness{witness("z", roleTest, support, recorded, "B", "z"), witness("a", roleImplementation, support, recorded, "A", "a")})
	before, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	_ = usecase.AnalyzeEvidenceCoalitions(in)
	after, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatalf("input mutated\nbefore: %s\nafter:  %s", before, after)
	}
}

func TestAnalyzeEvidenceCoalitionsLimitsAndExactOperations(t *testing.T) {
	t.Parallel()

	t.Run("inclusive rules and claims", func(t *testing.T) {
		rules := make([]domain.EvidenceCoalitionRule, 128)
		claims := make([]domain.EvidenceCoalitionClaim, 128)
		for i := range rules {
			id := fmt.Sprintf("rule-%03d", i)
			rules[i] = domain.EvidenceCoalitionRule{ID: id, RequiredRoles: []domain.EvidenceCoalitionRole{roleImplementation}, MinimumDistinctOrigins: 1}
			claims[i] = domain.EvidenceCoalitionClaim{ID: fmt.Sprintf("claim-%03d", i), RuleID: id, BaselineStatus: domain.EvidenceProven, Witnesses: []domain.EvidenceCoalitionWitness{witness(fmt.Sprintf("w-%03d", i), roleImplementation, support, recorded, fmt.Sprintf("origin-%03d", i), fmt.Sprintf("source-%03d", i))}}
		}
		report := usecase.AnalyzeEvidenceCoalitions(domain.EvidenceCoalitionInput{SchemaVersion: "evidence-coalition-input-v1", Rules: rules, Claims: claims})
		if report.Status != "COMPLETE" || report.Metrics.Operations != 256 || report.Metrics.Satisfied != 128 || report.Metrics.LimitExceeded {
			t.Fatalf("inclusive 128/128 boundary rejected: %+v", report.Metrics)
		}
	})

	t.Run("inclusive four roles and 64 witnesses", func(t *testing.T) {
		roles := []domain.EvidenceCoalitionRole{roleDocumentation, roleExecution, roleTest, roleImplementation}
		witnesses := make([]domain.EvidenceCoalitionWitness, 64)
		for i := range witnesses {
			witnesses[i] = witness(fmt.Sprintf("w-%02d", i), roles[i%4], support, recorded, fmt.Sprintf("origin-%02d", i), fmt.Sprintf("source-%02d", i))
		}
		report := usecase.AnalyzeEvidenceCoalitions(input("rule", "claim", domain.EvidencePartial, roles, 4, witnesses))
		if report.Status != "COMPLETE" || report.Metrics.Operations != 68 || report.Metrics.Witnesses != 64 || report.Metrics.Satisfied != 1 {
			t.Fatalf("inclusive role/witness boundary rejected: %+v", report.Metrics)
		}
	})

	t.Run("inclusive 4096 total witnesses", func(t *testing.T) {
		rules := make([]domain.EvidenceCoalitionRule, 64)
		claims := make([]domain.EvidenceCoalitionClaim, 64)
		for i := range claims {
			ruleID := fmt.Sprintf("r-%02d", i)
			rules[i] = domain.EvidenceCoalitionRule{ID: ruleID, RequiredRoles: []domain.EvidenceCoalitionRole{roleImplementation}, MinimumDistinctOrigins: 1}
			items := make([]domain.EvidenceCoalitionWitness, 64)
			for j := range items {
				items[j] = witness(fmt.Sprintf("w-%02d-%02d", i, j), roleImplementation, support, recorded, fmt.Sprintf("o-%02d-%02d", i, j), fmt.Sprintf("s-%02d-%02d", i, j))
			}
			claims[i] = domain.EvidenceCoalitionClaim{ID: fmt.Sprintf("c-%02d", i), RuleID: ruleID, BaselineStatus: domain.EvidenceInconclusive, Witnesses: items}
		}
		report := usecase.AnalyzeEvidenceCoalitions(domain.EvidenceCoalitionInput{SchemaVersion: "evidence-coalition-input-v1", Rules: rules, Claims: claims})
		if report.Status != "COMPLETE" || report.Metrics.Witnesses != 4096 || report.Metrics.Operations != 4160 {
			t.Fatalf("inclusive total witness boundary rejected: status=%s metrics=%+v", report.Status, report.Metrics)
		}
		if report.Metrics.Operations > 8192 {
			t.Fatalf("operation budget exceeded: %d", report.Metrics.Operations)
		}
	})

	for _, tc := range []struct {
		name string
		in   domain.EvidenceCoalitionInput
	}{
		{"129 rules", overRulesInput()}, {"129 claims", overClaimsInput()}, {"65 witnesses per claim", overWitnessesPerClaimInput()}, {"4097 total witnesses", overTotalWitnessesInput()}, {"five required roles", input("rule", "claim", domain.EvidenceProven, []domain.EvidenceCoalitionRole{roleImplementation, roleTest, roleExecution, roleDocumentation, roleImplementation}, 1, []domain.EvidenceCoalitionWitness{})},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertUnknownReport(t, usecase.AnalyzeEvidenceCoalitions(tc.in), "budget_exceeded")
		})
	}
}

func assertNonNilSlices(t *testing.T, report domain.EvidenceCoalitionReport) {
	t.Helper()
	if report.Assessments == nil {
		t.Error("report assessments is nil")
	}
	for _, a := range report.Assessments {
		if a.SupportingWitnessIDs == nil || a.ConflictingWitnessIDs == nil || a.MissingRoles == nil || a.DistinctOrigins == nil {
			t.Errorf("assessment contains nil public slice: %+v", a)
		}
	}
}

func assertUnknownReport(t *testing.T, report domain.EvidenceCoalitionReport, reason string) {
	t.Helper()
	if report.SchemaVersion != "evidence-coalition-report-v1" || report.Status != "UNKNOWN" || report.Reason != reason || len(report.Assessments) != 0 || report.Assessments == nil || report.Metrics.Satisfied != 0 || report.Metrics.VerdictChanges != 0 {
		t.Fatalf("UNKNOWN report mismatch: %+v", report)
	}
	if reason == "budget_exceeded" && !report.Metrics.LimitExceeded {
		t.Fatalf("budget report must set limit_exceeded: %+v", report.Metrics)
	}
}

func marshalReport(t *testing.T, report domain.EvidenceCoalitionReport) []byte {
	t.Helper()
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func overRulesInput() domain.EvidenceCoalitionInput {
	rules := make([]domain.EvidenceCoalitionRule, 129)
	for i := range rules {
		rules[i] = domain.EvidenceCoalitionRule{ID: fmt.Sprintf("r-%03d", i), RequiredRoles: []domain.EvidenceCoalitionRole{roleImplementation}, MinimumDistinctOrigins: 1}
	}
	return domain.EvidenceCoalitionInput{SchemaVersion: "evidence-coalition-input-v1", Rules: rules, Claims: []domain.EvidenceCoalitionClaim{{ID: "claim", RuleID: "r-000", BaselineStatus: domain.EvidenceProven, Witnesses: []domain.EvidenceCoalitionWitness{}}}}
}

func overClaimsInput() domain.EvidenceCoalitionInput {
	claims := make([]domain.EvidenceCoalitionClaim, 129)
	for i := range claims {
		claims[i] = domain.EvidenceCoalitionClaim{ID: fmt.Sprintf("c-%03d", i), RuleID: "rule", BaselineStatus: domain.EvidenceProven, Witnesses: []domain.EvidenceCoalitionWitness{}}
	}
	return domain.EvidenceCoalitionInput{SchemaVersion: "evidence-coalition-input-v1", Rules: []domain.EvidenceCoalitionRule{{ID: "rule", RequiredRoles: []domain.EvidenceCoalitionRole{roleImplementation}, MinimumDistinctOrigins: 1}}, Claims: claims}
}

func overWitnessesPerClaimInput() domain.EvidenceCoalitionInput {
	witnesses := make([]domain.EvidenceCoalitionWitness, 65)
	for i := range witnesses {
		witnesses[i] = witness(fmt.Sprintf("w-%02d", i), roleImplementation, support, recorded, fmt.Sprintf("o-%02d", i), fmt.Sprintf("s-%02d", i))
	}
	return input("rule", "claim", domain.EvidenceProven, []domain.EvidenceCoalitionRole{roleImplementation}, 1, witnesses)
}

func overTotalWitnessesInput() domain.EvidenceCoalitionInput {
	rules := make([]domain.EvidenceCoalitionRule, 65)
	claims := make([]domain.EvidenceCoalitionClaim, 65)
	for i := range claims {
		ruleID := fmt.Sprintf("r-%02d", i)
		rules[i] = domain.EvidenceCoalitionRule{ID: ruleID, RequiredRoles: []domain.EvidenceCoalitionRole{roleImplementation}, MinimumDistinctOrigins: 1}
		count := 64
		if i == 64 {
			count = 1
		}
		items := make([]domain.EvidenceCoalitionWitness, count)
		for j := range items {
			items[j] = witness(fmt.Sprintf("w-%02d-%02d", i, j), roleImplementation, support, recorded, fmt.Sprintf("o-%02d-%02d", i, j), fmt.Sprintf("s-%02d-%02d", i, j))
		}
		claims[i] = domain.EvidenceCoalitionClaim{ID: fmt.Sprintf("c-%02d", i), RuleID: ruleID, BaselineStatus: domain.EvidenceProven, Witnesses: items}
	}
	return domain.EvidenceCoalitionInput{SchemaVersion: "evidence-coalition-input-v1", Rules: rules, Claims: claims}
}
