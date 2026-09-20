/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o3_controls_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"testing"
)

func TestStep0028O3ReferencePreservesTwelveFullPositiveControls(t *testing.T) {
	cases := O3PositiveCases()
	if len(cases) != 12 {
		t.Fatalf("positive controls = %d, want 12", len(cases))
	}
	findings := Evaluate(context.Background(), ReferenceSubject{}, cases)
	seen := map[string]bool{}
	for index, finding := range findings {
		if finding.Outcome != "Preserved" {
			t.Fatalf("reference failed %s: %s", finding.CaseID, finding.Detail)
		}
		seen[finding.CaseID] = true
		requireO3ControlSemantics(t, cases[index].ID, cases[index].Expected)
	}
	for number := 1; number <= 12; number++ {
		id := positiveID(number)
		if !seen[id] {
			t.Fatalf("positive control %s not demonstrated", id)
		}
	}
}

func requireO3ControlSemantics(t *testing.T, id string, observed Observation) {
	t.Helper()
	switch id {
	case "P01":
		if observed.MechanismStatus != "PROVEN" || !observed.BindingVerified || !observed.ProvenanceComplete || len(observed.SourceRefs) != 2 || observed.Owner == "" || observed.Repo == "" || observed.Commit == "" || observed.Path == "" {
			t.Fatalf("P01 incomplete: %+v", observed)
		}
	case "P02":
		if observed.MechanismStatus != "PARTIAL" || observed.Novelty != "UNKNOWN" {
			t.Fatalf("P02 incomplete: %+v", observed)
		}
	case "P03":
		if observed.Decision != "OBSERVATION_ONLY" || observed.MechanismStatus != "NONE" || len(observed.ArtifactNames) != 3 || len(observed.SourceRefs) != 2 || len(observed.ValidationPlan) != 3 {
			t.Fatalf("P03 incomplete: %+v", observed)
		}
	case "P04":
		if observed.TargetState != "UNKNOWN" || observed.HypothesisStatus != "PROPOSED" || len(observed.TargetGapRefs) == 0 || len(observed.ValidationPlan) == 0 {
			t.Fatalf("P04 incomplete: %+v", observed)
		}
	case "P05":
		if observed.License != "UNRECOGNIZED" || observed.TransferMode != "PATTERN_ADAPTATION" || len(observed.WarningCodes) == 0 || len(observed.SourceRefs) == 0 {
			t.Fatalf("P05 incomplete: %+v", observed)
		}
	case "P06":
		if len(observed.MechanismRefs) != 2 || len(observed.SourceRefs) != 2 || len(observed.IntegrationUnknowns) != 5 {
			t.Fatalf("P06 incomplete: %+v", observed)
		}
	case "P07":
		if observed.MechanismStatus != "INCONCLUSIVE" || observed.HypothesisStatus != "INCONCLUSIVE" || len(observed.CounterScope) == 0 || len(observed.CounterEvidenceRefs) == 0 {
			t.Fatalf("P07 incomplete: %+v", observed)
		}
	case "P08":
		if len(observed.Ordering) != 2 || observed.License == "" || observed.Fit == "" || observed.CanonicalSHA256 == "" || observed.ArtifactSHA256 == "" {
			t.Fatalf("P08 incomplete: %+v", observed)
		}
	case "P09":
		if len(observed.RunIDs) != 4 || observed.ResolveCalls != 1 || observed.ReserveCalls != 1 || observed.IDCalls != 1 || observed.CreateCalls != 1 || observed.CanonicalSHA256 == "" {
			t.Fatalf("P09 incomplete: %+v", observed)
		}
	case "P10":
		if observed.Lifecycle != "CANCELLED" || observed.RecoveryState != "INTERRUPTED" || len(observed.ArtifactNames) == 0 || len(observed.ExternalEffects) != 0 {
			t.Fatalf("P10 incomplete: %+v", observed)
		}
	case "P11":
		if observed.NetworkCalls != 0 || observed.ReplaySHA256 == "" || len(observed.ArtifactNames) != 7 || len(observed.ArtifactHashes) != 7 || len(observed.IndexHashes) != 2 {
			t.Fatalf("P11 incomplete: %+v", observed)
		}
	case "P12":
		if !observed.RuntimeUnique || !observed.RefinementsPreserved || !observed.BasalHashVerified || len(observed.Tools) != 9 || len(observed.RefinementReceipts) != 3 || len(observed.LayerFindings) != 3 || len(observed.SchemaDefinitions) == 0 {
			t.Fatalf("P12 incomplete: %+v", observed)
		}
	}
}

func positiveID(number int) string {
	if number < 10 {
		return "P0" + string(rune('0'+number))
	}
	return "P" + string(rune('0'+number/10)) + string(rune('0'+number%10))
}
