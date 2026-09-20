/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0066_gate/discriminatory_test.md
 * @layer L4
 * @updated 2026-09-11
 */

package step0066_gate

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
	Controls            []radarCase      `json:"controls"`
	Mutations           []radarCase      `json:"mutations"`
	UtilityControls     []utilityControl `json:"utility_controls"`
}

type radarCase struct {
	ID           string      `json:"id"`
	Class        string      `json:"class"`
	Witness      string      `json:"witness"`
	InputState   string      `json:"input_state"`
	Opacity      string      `json:"opacity"`
	PhysicalRefs []string    `json:"physical_refs"`
	Supports     [][5]string `json:"supports"`
	VariantJSON  []string    `json:"variant_json"`
	Claim        claim       `json:"claim"`
}

type claim struct {
	Status                string   `json:"status"`
	Reason                string   `json:"reason"`
	Outcome               string   `json:"outcome"`
	CanonicalObservations int      `json:"canonical_observations"`
	MissingKinds          int      `json:"missing_kinds"`
	MissingOrigins        int      `json:"missing_origins"`
	DraftReady            bool     `json:"draft_ready"`
	DraftReason           string   `json:"draft_reason"`
	RequiredRoles         []string `json:"required_roles"`
	VerdictChanges        int      `json:"verdict_changes"`
	LimitExceeded         bool     `json:"limit_exceeded"`
}

type counts struct {
	TP int `json:"tp"`
	FP int `json:"fp"`
	FN int `json:"fn"`
}

type utilityControl struct {
	ID               string `json:"id"`
	Class            string `json:"class"`
	BindingComplete  bool   `json:"binding_complete"`
	Baseline         counts `json:"baseline"`
	Radar            counts `json:"radar"`
	ImprovedCases    int    `json:"improved_cases"`
	ImprovedFamilies int    `json:"improved_families"`
	ClaimedVerdict   string `json:"claimed_verdict"`
}

func loadFixture(t *testing.T) fixture {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve gate source")
	}
	path := filepath.Join(filepath.Dir(source), "..", "..", "..", "..", "00_nucleo", "benchmark", "0066", "adversarial-fixtures.json")
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

func evaluateRadar(c radarCase) (string, string) {
	if c.InputState == "opaque" {
		if c.Opacity == "" {
			return violated, "opaque control lacks public opacity reason"
		}
		return unknown, c.Opacity
	}

	if c.InputState == "invalid" {
		if c.Claim.Status != "UNKNOWN" || c.Claim.Outcome != "" || c.Claim.LimitExceeded {
			return violated, "invalid input was not rejected atomically"
		}
		return preserved, "invalid input remains UNKNOWN"
	}
	if c.InputState == "budget_exceeded" {
		if c.Claim.Status != "UNKNOWN" || c.Claim.Reason != "budget_exceeded" || c.Claim.Outcome != "" || !c.Claim.LimitExceeded {
			return violated, "budget excess produced partial or non-UNKNOWN output"
		}
		return preserved, "budget excess is atomic UNKNOWN"
	}
	if c.InputState != "valid" {
		return unknown, "unsupported input state"
	}

	canonical := make(map[string][5]string)
	for _, observation := range c.Supports {
		key := strings.Join(observation[:4], "\x00")
		prior, exists := canonical[key]
		if !exists || observation[4] < prior[4] {
			canonical[key] = observation
		}
	}
	kinds := map[string]bool{}
	recordedKinds := map[string]bool{}
	origins := map[string]bool{}
	hasConflict, hasAmbiguous := false, false
	for _, observation := range canonical {
		kind, originStatus, origin, stance := observation[0], observation[1], observation[2], observation[3]
		switch stance {
		case "CONFLICT":
			hasConflict = true
		case "AMBIGUOUS":
			hasAmbiguous = true
		case "SUPPORT":
			kinds[kind] = true
			if originStatus == "RECORDED" {
				recordedKinds[kind] = true
				origins[origin] = true
			}
		}
	}

	expectedOutcome := "CONVERGENT"
	switch {
	case hasConflict:
		expectedOutcome = "CONFLICT"
	case hasAmbiguous:
		expectedOutcome = "AMBIGUOUS"
	case len(kinds) < 3:
		expectedOutcome = "INSUFFICIENT_PERSPECTIVES"
	case len(recordedKinds) < 3:
		expectedOutcome = "UNKNOWN_ORIGIN"
	case len(origins) < 3:
		expectedOutcome = "INSUFFICIENT_ORIGINS"
	}
	expectedMissingKinds := max(0, 3-len(kinds))
	expectedMissingOrigins := max(0, 3-len(origins))
	expectedRoles := []string{"IMPLEMENTATION", "INDEPENDENT_TEST"}

	var defects []string
	check := func(ok bool, name string) {
		if !ok {
			defects = append(defects, name)
		}
	}
	check(c.Claim.Status == "COMPLETE" && c.Claim.Reason == "analyzed", "report status/reason")
	check(c.Claim.Outcome == expectedOutcome, "outcome")
	check(c.Claim.CanonicalObservations == len(canonical), "canonical observation count")
	check(c.Claim.MissingKinds == expectedMissingKinds, "missing kinds")
	check(c.Claim.MissingOrigins == expectedMissingOrigins, "missing origins")
	check(!c.Claim.DraftReady && c.Claim.DraftReason == "requires_evidence_binding", "0065 draft boundary")
	check(equalStrings(c.Claim.RequiredRoles, expectedRoles), "0065 required roles")
	check(c.Claim.VerdictChanges == 0, "zero verdict changes")
	check(!c.Claim.LimitExceeded, "valid input limit flag")
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
	return preserved, "all public invariants preserved"
}

func evaluateUtility(c utilityControl) (string, string) {
	if !c.BindingComplete {
		return unknown, "candidate/result binding unavailable"
	}
	bp, br, bf, bok := metrics(c.Baseline)
	rp, rr, rf, rok := metrics(c.Radar)
	expected := "INCONCLUSIVE"
	if bok && rok {
		expected = "NOT_OBSERVED"
		precisionDelta := rp - bp
		recallDelta := rr - br
		f1Delta := rf - bf
		if precisionDelta >= -1e-12 && (recallDelta >= .1-1e-12 || f1Delta >= .1-1e-12) && c.ImprovedCases >= 2 && c.ImprovedFamilies >= 2 {
			expected = "OBSERVED"
		}
	}
	if c.ClaimedVerdict != expected {
		return violated, fmt.Sprintf("claimed %s, external contract requires %s", c.ClaimedVerdict, expected)
	}
	return preserved, "external verdict " + expected
}

func metrics(c counts) (precision, recall, f1 float64, ok bool) {
	if c.TP+c.FP == 0 || c.TP+c.FN == 0 || 2*c.TP+c.FP+c.FN == 0 {
		return 0, 0, 0, false
	}
	precision = float64(c.TP) / float64(c.TP+c.FP)
	recall = float64(c.TP) / float64(c.TP+c.FN)
	f1 = float64(2*c.TP) / float64(2*c.TP+c.FP+c.FN)
	return precision, recall, f1, !(math.IsNaN(precision) || math.IsNaN(recall) || math.IsNaN(f1))
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

func TestStep0066FixtureShapeAndMutationScore(t *testing.T) {
	f := loadFixture(t)
	if f.SchemaVersion != "tekt.step0066.binding-free-fixtures.v1" || f.Step != "0066" || len(f.ManifestSHA256) != 64 {
		t.Fatalf("invalid fixture identity: %+v", f)
	}
	if f.RequiredRepetitions != 20 {
		t.Fatalf("required repetitions = %d, want 20", f.RequiredRepetitions)
	}
	if len(f.Mutations) != 14 {
		t.Fatalf("mutation denominator = %d, want 14", len(f.Mutations))
	}
	for i, mutation := range f.Mutations {
		wantID := fmt.Sprintf("M%02d", i+1)
		if mutation.ID != wantID || mutation.Class != "N" || mutation.Witness == "" {
			t.Fatalf("invalid mutation %d: %+v", i, mutation)
		}
		classification, reason := evaluateRadar(mutation)
		if classification != violated {
			t.Fatalf("%s survived as %s (%s)", mutation.ID, classification, reason)
		}
	}
	if score := float64(len(f.Mutations)) / float64(len(f.Mutations)); score != 1.0 {
		t.Fatalf("mutation_score=%v, want 1.0", score)
	}
}

func TestStep0066PositiveNegativeOpaqueAndExternalVerdicts(t *testing.T) {
	f := loadFixture(t)
	seen := map[string]int{}
	for _, c := range f.Controls {
		classification, reason := evaluateRadar(c)
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
}

func TestStep0066TwentyOrdersAreIdentical(t *testing.T) {
	f := loadFixture(t)
	all := append(append([]radarCase{}, f.Controls...), f.Mutations...)
	baseline := classificationVector(all)
	for repetition := 0; repetition < f.RequiredRepetitions; repetition++ {
		ordered := rotate(all, repetition%len(all))
		if repetition%2 == 1 {
			reverse(ordered)
		}
		got := classificationVector(ordered)
		if !equalMap(got, baseline) {
			t.Fatalf("order/repetition %02d changed classification vector", repetition+1)
		}
	}
}

func classificationVector(cases []radarCase) map[string]string {
	result := make(map[string]string, len(cases))
	for _, c := range cases {
		classification, _ := evaluateRadar(c)
		result[c.ID] = classification
	}
	return result
}

func rotate(in []radarCase, n int) []radarCase {
	out := append([]radarCase{}, in[n:]...)
	return append(out, in[:n]...)
}

func reverse(in []radarCase) {
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
