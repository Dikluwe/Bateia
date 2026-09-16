/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0067_gate/discriminatory_test.md
 * @layer L4
 * @updated 2026-09-11
 */

package step0067_gate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

const (
	preserved = "Preserved"
	violated  = "Violated"
	unknown   = "Unknown"
)

type fixture struct {
	SchemaVersion       string           `json:"schema_version"`
	Step                string           `json:"step"`
	ManifestSHA256      string           `json:"manifest_sha256"`
	RequiredRepetitions int              `json:"required_repetitions"`
	Controls            []planCase       `json:"controls"`
	Mutations           []planCase       `json:"mutations"`
	UtilityControls     []utilityControl `json:"utility_controls"`
}

type planCase struct {
	ID               string   `json:"id"`
	Class            string   `json:"class"`
	Witness          string   `json:"witness"`
	InputState       string   `json:"input_state"`
	Opacity          string   `json:"opacity"`
	Outcome          string   `json:"outcome"`
	TargetID         string   `json:"target_id"`
	RequiredFrontier []string `json:"required_frontier"`
	VariantJSON      []string `json:"variant_json"`
	Claim            claim    `json:"claim"`
}

type claim struct {
	Status              string   `json:"status"`
	Reason              string   `json:"reason"`
	PlanReason          string   `json:"plan_reason"`
	Outcome             string   `json:"outcome"`
	Frontier            []string `json:"frontier"`
	Plans               int      `json:"plans"`
	Probes              int      `json:"probes"`
	RadarOperations     int      `json:"radar_operations"`
	PlanningOperations  int      `json:"planning_operations"`
	Operations          int      `json:"operations"`
	LimitExceeded       bool     `json:"limit_exceeded"`
	ArraysNonNull       bool     `json:"arrays_non_null"`
	PromisesConvergence bool     `json:"promises_convergence"`
	FabricatesEvidence  bool     `json:"fabricates_evidence"`
	ReadyFor0065        bool     `json:"ready_for_0065"`
	Calls0065           bool     `json:"calls_0065"`
	MutatesInput        bool     `json:"mutates_input"`
	VerdictChanges      int      `json:"verdict_changes"`
	ConsumerChanged     bool     `json:"consumer_changed"`
	ThresholdChanged    bool     `json:"threshold_changed"`
}

type utilityControl struct {
	ID                string `json:"id"`
	Class             string `json:"class"`
	BindingComplete   bool   `json:"binding_complete"`
	Eligible          int    `json:"eligible"`
	PlannerResolved   int    `json:"planner_resolved"`
	NaiveResolved     int    `json:"naive_resolved"`
	PlannerConvergent int    `json:"planner_convergent"`
	NaiveConvergent   int    `json:"naive_convergent"`
	PlannerProbes     int    `json:"planner_probes"`
	NaiveProbes       int    `json:"naive_probes"`
	ImprovedCases     int    `json:"improved_cases"`
	ImprovedFamilies  int    `json:"improved_families"`
	ForbiddenProbe    bool   `json:"forbidden_probe"`
	FalsePromise      bool   `json:"false_promise"`
	BudgetExceeded    bool   `json:"budget_exceeded"`
	VerdictChanges    int    `json:"verdict_changes"`
	ClaimedVerdict    string `json:"claimed_verdict"`
}

func loadFixture(t *testing.T) fixture {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve gate source")
	}
	path := filepath.Join(filepath.Dir(source), "..", "..", "00_nucleo", "benchmark", "0067", "adversarial-fixtures.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read frozen fixture: %v", err)
	}
	var got fixture
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("decode frozen fixture: %v", err)
	}
	if dec.More() {
		t.Fatal("trailing fixture value")
	}
	return got
}

func evaluatePlan(c planCase) (string, string) {
	if c.InputState == "opaque" {
		if c.Opacity == "" {
			return violated, "opaque control lacks public opacity reason"
		}
		return unknown, c.Opacity
	}

	var defects []string
	check := func(ok bool, name string) {
		if !ok {
			defects = append(defects, name)
		}
	}
	atomicUnknown := c.InputState == "radar_unknown" || c.InputState == "invalid_text_inference" || c.InputState == "budget_exceeded"
	if atomicUnknown {
		check(c.Claim.Status == "UNKNOWN", "atomic UNKNOWN status")
		if c.InputState == "budget_exceeded" {
			check(c.Claim.Reason == "budget_exceeded" && c.Claim.LimitExceeded, "budget reason/flag")
		} else {
			check(!c.Claim.LimitExceeded, "non-budget limit flag")
		}
		check(c.Claim.Plans == 0 && c.Claim.Probes == 0 && len(c.Claim.Frontier) == 0, "UNKNOWN partial plan")
		check(c.Claim.PlanReason == "" && c.Claim.Outcome == "", "UNKNOWN leaked assessment")
	} else if c.InputState == "valid" {
		check(c.Claim.Status == "COMPLETE" && c.Claim.Reason == "planned", "complete status/reason")
		check(c.Claim.Outcome == c.Outcome, "0066 outcome conservation")
		wantReason := "active_blocker_frontier"
		if c.Outcome == "CONVERGENT" {
			wantReason = "evidence_binding_only"
		}
		check(c.Claim.PlanReason == wantReason, "active plan reason")
		check(equalStrings(c.Claim.Frontier, c.RequiredFrontier), "minimal canonical frontier")
		check(c.Claim.Plans == 1 && c.Claim.Probes == len(c.RequiredFrontier), "plan/probe counts")
		check(c.Claim.Operations == c.Claim.RadarOperations+c.Claim.PlanningOperations, "operation sum")
		check(c.Claim.Operations <= 8192 && c.Claim.Probes <= 1024 && c.Claim.Plans <= 512, "planner budget")
		check(!c.Claim.LimitExceeded, "valid input limit flag")
	} else {
		return unknown, "unsupported input state"
	}

	check(c.Claim.ArraysNonNull, "nonnull output arrays")
	check(!c.Claim.PromisesConvergence, "no convergence promise")
	check(!c.Claim.FabricatesEvidence && !c.Claim.ReadyFor0065 && !c.Claim.Calls0065, "0065/evidence boundary")
	check(!c.Claim.MutatesInput, "input immutability")
	check(c.Claim.VerdictChanges == 0 && !c.Claim.ConsumerChanged && !c.Claim.ThresholdChanged, "zero decision impact")
	if len(c.VariantJSON) > 0 {
		first := c.VariantJSON[0]
		for _, variant := range c.VariantJSON[1:] {
			check(variant == first, "byte-identical order/nil variants")
		}
	}
	if len(defects) > 0 {
		sort.Strings(defects)
		return violated, strings.Join(defects, ", ")
	}
	return preserved, "all public planning invariants preserved"
}

func evaluateUtility(c utilityControl) (string, string) {
	if !c.BindingComplete {
		return unknown, "candidate/result binding unavailable"
	}
	expected := "INCONCLUSIVE"
	if c.Eligible > 0 && c.NaiveProbes > 0 {
		expected = "NOT_OBSERVED"
		plannerRate := float64(c.PlannerResolved) / float64(c.Eligible)
		naiveRate := float64(c.NaiveResolved) / float64(c.Eligible)
		reduction := float64(c.NaiveProbes-c.PlannerProbes) / float64(c.NaiveProbes)
		resolutionGain := plannerRate - naiveRate
		noForbidden := !c.ForbiddenProbe && !c.FalsePromise && !c.BudgetExceeded && c.VerdictChanges == 0
		usefulDelta := reduction >= .2-1e-12 || resolutionGain >= .1-1e-12 && c.PlannerProbes <= c.NaiveProbes
		if plannerRate+1e-12 >= naiveRate && usefulDelta && c.ImprovedCases >= 2 && c.ImprovedFamilies >= 2 && noForbidden {
			expected = "OBSERVED"
		}
	}
	if c.ClaimedVerdict != expected {
		return violated, fmt.Sprintf("claimed %s, external contract requires %s", c.ClaimedVerdict, expected)
	}
	return preserved, "external verdict " + expected
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestStep0067FixtureShapeAndMutationScore(t *testing.T) {
	f := loadFixture(t)
	if f.SchemaVersion != "tekt.step0067.binding-free-fixtures.v1" || f.Step != "0067" || len(f.ManifestSHA256) != 64 {
		t.Fatalf("invalid fixture identity: %+v", f)
	}
	if f.RequiredRepetitions != 20 {
		t.Fatalf("required repetitions = %d, want 20", f.RequiredRepetitions)
	}
	if len(f.Mutations) != 16 {
		t.Fatalf("mutation denominator = %d, want 16", len(f.Mutations))
	}
	killed := 0
	for i, mutation := range f.Mutations {
		wantID := fmt.Sprintf("M%02d", i+1)
		if mutation.ID != wantID || mutation.Class != "N" || mutation.Witness == "" {
			t.Fatalf("invalid mutation %d: %+v", i, mutation)
		}
		classification, reason := evaluatePlan(mutation)
		if classification != violated {
			t.Fatalf("%s survived as %s (%s)", mutation.ID, classification, reason)
		}
		killed++
	}
	score := float64(killed) / float64(len(f.Mutations))
	if math.Abs(score-1.0) > 1e-12 {
		t.Fatalf("mutation_score=%v, want 1.0", score)
	}
}

func TestStep0067PositiveNegativeOpaqueAndExternalMetrics(t *testing.T) {
	f := loadFixture(t)
	seen := map[string]int{}
	for _, c := range f.Controls {
		classification, reason := evaluatePlan(c)
		seen[c.Class]++
		want := map[string]string{"P": preserved, "O": unknown}[c.Class]
		if classification != want {
			t.Errorf("%s = %s (%s), want %s", c.ID, classification, reason, want)
		}
	}
	verdicts := map[string]bool{}
	for _, c := range f.UtilityControls {
		classification, reason := evaluateUtility(c)
		seen[c.Class]++
		want := map[string]string{"P": preserved, "N": violated, "O": unknown}[c.Class]
		if classification != want {
			t.Errorf("%s = %s (%s), want %s", c.ID, classification, reason, want)
		}
		if classification == preserved {
			verdicts[c.ClaimedVerdict] = true
		}
	}
	for _, class := range []string{"P", "N", "O"} {
		if seen[class] == 0 {
			t.Errorf("missing %s control", class)
		}
	}
	for _, verdict := range []string{"OBSERVED", "NOT_OBSERVED", "INCONCLUSIVE"} {
		if !verdicts[verdict] {
			t.Errorf("external verdict %s was not positively discriminated", verdict)
		}
	}
	u := f.UtilityControls[0]
	if u.Eligible != 8 || u.PlannerResolved != 8 || u.NaiveResolved != 8 || u.PlannerConvergent != 7 || u.NaiveConvergent != 7 || u.PlannerProbes != 11 || u.NaiveProbes != 15 {
		t.Fatalf("frozen external metrics drifted: %+v", u)
	}
	if reduction := float64(u.NaiveProbes-u.PlannerProbes) / float64(u.NaiveProbes); math.Abs(reduction-4.0/15.0) > 1e-12 {
		t.Fatalf("effort reduction=%v, want 4/15", reduction)
	}
}

func TestStep0067TwentyOrdersAreIdentical(t *testing.T) {
	f := loadFixture(t)
	all := append(append([]planCase{}, f.Controls...), f.Mutations...)
	baseline := classificationVector(all)
	for repetition := 0; repetition < f.RequiredRepetitions; repetition++ {
		ordered := rotate(all, repetition%len(all))
		if repetition%2 == 1 {
			reverse(ordered)
		}
		if got := classificationVector(ordered); !equalMap(got, baseline) {
			t.Fatalf("order/repetition %02d changed classification vector", repetition+1)
		}
	}
}

func classificationVector(cases []planCase) map[string]string {
	result := make(map[string]string, len(cases))
	for _, c := range cases {
		classification, _ := evaluatePlan(c)
		result[c.ID] = classification
	}
	return result
}

func rotate(in []planCase, n int) []planCase {
	out := append([]planCase{}, in[n:]...)
	return append(out, in[:n]...)
}

func reverse(in []planCase) {
	for left, right := 0, len(in)-1; left < right; left, right = left+1, right-1 {
		in[left], in[right] = in[right], in[left]
	}
}

func equalMap(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if b[key] != value {
			return false
		}
	}
	return true
}
