/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_coalition_test.md
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

func TestEvidenceCoalitionContractValues(t *testing.T) {
	t.Parallel()

	values := map[string]string{
		"implementation":       string(domain.EvidenceCoalitionRole("IMPLEMENTATION")),
		"independent test":     string(domain.EvidenceCoalitionRole("INDEPENDENT_TEST")),
		"execution":            string(domain.EvidenceCoalitionRole("EXECUTION")),
		"documentation":        string(domain.EvidenceCoalitionRole("DOCUMENTATION")),
		"support":              string(domain.EvidenceCoalitionPolarity("SUPPORT")),
		"conflict":             string(domain.EvidenceCoalitionPolarity("CONFLICT")),
		"recorded":             string(domain.EvidenceOriginStatus("RECORDED")),
		"same origin":          string(domain.EvidenceOriginStatus("SAME_ORIGIN")),
		"distinct origin":      string(domain.EvidenceOriginStatus("DISTINCT_ORIGIN")),
		"unknown origin":       string(domain.EvidenceOriginStatus("UNKNOWN_ORIGIN")),
		"satisfied":            string(domain.EvidenceCoalitionOutcome("SATISFIED")),
		"missing roles":        string(domain.EvidenceCoalitionOutcome("MISSING_ROLES")),
		"unknown independence": string(domain.EvidenceCoalitionOutcome("UNKNOWN_INDEPENDENCE")),
		"insufficient":         string(domain.EvidenceCoalitionOutcome("INSUFFICIENT_INDEPENDENCE")),
		"invalid":              string(domain.EvidenceCoalitionOutcome("INVALID_INPUT")),
	}
	want := map[string]string{
		"implementation": "IMPLEMENTATION", "independent test": "INDEPENDENT_TEST",
		"execution": "EXECUTION", "documentation": "DOCUMENTATION", "support": "SUPPORT",
		"conflict": "CONFLICT", "recorded": "RECORDED", "same origin": "SAME_ORIGIN",
		"distinct origin": "DISTINCT_ORIGIN", "unknown origin": "UNKNOWN_ORIGIN",
		"satisfied": "SATISFIED", "missing roles": "MISSING_ROLES",
		"unknown independence": "UNKNOWN_INDEPENDENCE", "insufficient": "INSUFFICIENT_INDEPENDENCE",
		"invalid": "INVALID_INPUT",
	}
	for name, expected := range want {
		if values[name] != expected {
			t.Errorf("%s = %q, want %q", name, values[name], expected)
		}
	}
}

func TestEvidenceCoalitionJSONTagsAreSnakeCaseAndNeverOmitObservableSlices(t *testing.T) {
	t.Parallel()

	cases := []struct {
		value any
		tags  map[string]string
	}{
		{domain.EvidenceCoalitionInput{}, map[string]string{"SchemaVersion": "schema_version", "Rules": "rules", "Claims": "claims"}},
		{domain.EvidenceCoalitionRule{}, map[string]string{"ID": "id", "RequiredRoles": "required_roles", "MinimumDistinctOrigins": "minimum_distinct_origins"}},
		{domain.EvidenceCoalitionClaim{}, map[string]string{"ID": "id", "RuleID": "rule_id", "BaselineStatus": "baseline_status", "Witnesses": "witnesses"}},
		{domain.EvidenceCoalitionWitness{}, map[string]string{"ID": "id", "Role": "role", "Polarity": "polarity", "OriginStatus": "origin_status", "OriginID": "origin_id", "Source": "source"}},
		{domain.EvidenceCoalitionAssessment{}, map[string]string{"ClaimID": "claim_id", "RuleID": "rule_id", "BaselineStatus": "baseline_status", "DecisionImpact": "decision_impact", "Outcome": "outcome", "Independence": "independence", "SupportingWitnessIDs": "supporting_witness_ids", "ConflictingWitnessIDs": "conflicting_witness_ids", "MissingRoles": "missing_roles", "DistinctOrigins": "distinct_origins", "Reason": "reason"}},
		{domain.EvidenceCoalitionMetrics{}, map[string]string{"Claims": "claims", "Witnesses": "witnesses", "Satisfied": "satisfied", "Conflicts": "conflicts", "MissingRoles": "missing_roles", "UnknownIndependence": "unknown_independence", "InsufficientIndependence": "insufficient_independence", "Invalid": "invalid", "Operations": "operations", "LimitExceeded": "limit_exceeded", "VerdictChanges": "verdict_changes"}},
		{domain.EvidenceCoalitionReport{}, map[string]string{"SchemaVersion": "schema_version", "Status": "status", "Reason": "reason", "Assessments": "assessments", "Metrics": "metrics"}},
	}

	for _, tc := range cases {
		typ := reflect.TypeOf(tc.value)
		for fieldName, wantTag := range tc.tags {
			field, ok := typ.FieldByName(fieldName)
			if !ok {
				t.Fatalf("%s lacks public field %s", typ.Name(), fieldName)
			}
			if got := field.Tag.Get("json"); got != wantTag {
				t.Errorf("%s.%s json tag = %q, want %q (omitempty is forbidden)", typ.Name(), fieldName, got, wantTag)
			}
		}
	}
}

func TestEvidenceCoalitionJSONRoundTripPreservesBaselineAndEmptyArrays(t *testing.T) {
	t.Parallel()

	report := domain.EvidenceCoalitionReport{
		SchemaVersion: "evidence-coalition-report-v1",
		Status:        "COMPLETE",
		Reason:        "analyzed",
		Assessments: []domain.EvidenceCoalitionAssessment{{
			ClaimID: "claim", RuleID: "rule", BaselineStatus: domain.EvidenceConflict,
			DecisionImpact: "NONE", Outcome: domain.EvidenceCoalitionOutcome("MISSING_ROLES"),
			Independence:         domain.EvidenceOriginStatus("UNKNOWN_ORIGIN"),
			SupportingWitnessIDs: []string{}, ConflictingWitnessIDs: []string{},
			MissingRoles: []domain.EvidenceCoalitionRole{}, DistinctOrigins: []string{}, Reason: "missing_required_roles",
		}},
	}
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		t.Fatal(err)
	}
	assessments, ok := object["assessments"].([]any)
	if !ok || len(assessments) != 1 {
		t.Fatalf("assessments must be a non-null JSON array: %s", raw)
	}
	assessment := assessments[0].(map[string]any)
	for _, key := range []string{"supporting_witness_ids", "conflicting_witness_ids", "missing_roles", "distinct_origins"} {
		if array, ok := assessment[key].([]any); !ok || len(array) != 0 {
			t.Errorf("%s must be encoded as []: %s", key, raw)
		}
	}
	var decoded domain.EvidenceCoalitionReport
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Assessments[0].BaselineStatus != domain.EvidenceConflict {
		t.Fatalf("baseline changed on JSON round trip: got %q", decoded.Assessments[0].BaselineStatus)
	}
}
