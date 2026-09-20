/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0068_gate/discriminatory_test.md
 * @layer L4
 * @updated 2026-09-11
 */

package step0068_gate

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

var canonicalQueries = []string{
	"robust geometric orientation and in-circle predicates for nearly degenerate cases",
	"orient2d incircle",
	"Shewchuk exact predicates",
}

var canonicalRules = []string{
	"robust-predicates.translation.en",
	"robust-predicates.symbol-pair",
	"robust-predicates.reference-pair",
}

var canonicalDiscriminators = [][]string{
	{"orientation", "in-circle"},
	{"orient2d", "incircle"},
	{"Shewchuk", "predicates"},
}

type fixture struct {
	SchemaVersion       string           `json:"schema_version"`
	Step                string           `json:"step"`
	ManifestSHA256      string           `json:"manifest_sha256"`
	RequiredRepetitions int              `json:"required_repetitions"`
	ReferenceClaim      compilerClaim    `json:"reference_claim"`
	Controls            []control        `json:"controls"`
	Mutations           []mutation       `json:"mutations"`
	UtilityControls     []utilityControl `json:"utility_controls"`
}

type control struct {
	ID         string        `json:"id"`
	Class      string        `json:"class"`
	InputState string        `json:"input_state"`
	Opacity    string        `json:"opacity,omitempty"`
	Claim      compilerClaim `json:"claim"`
}

type mutation struct {
	ID             string                     `json:"id"`
	Class          string                     `json:"class"`
	Witness        string                     `json:"witness"`
	ExpectedDefect string                     `json:"expected_defect"`
	InputState     string                     `json:"input_state"`
	Patch          map[string]json.RawMessage `json:"patch"`
}

type compilerClaim struct {
	SourcePhrase           string     `json:"source_phrase"`
	Translation            string     `json:"translation"`
	Literal                string     `json:"literal"`
	LiteralPreserved       bool       `json:"literal_preserved"`
	CatalogMatched         bool       `json:"catalog_matched"`
	ProposalQueries        []string   `json:"proposal_queries"`
	ProposalStrategies     []string   `json:"proposal_strategies"`
	ProposalRuleIDs        []string   `json:"proposal_rule_ids"`
	ProposalSourceIDs      []string   `json:"proposal_source_ids"`
	ProposalDiscriminators [][]string `json:"proposal_discriminators"`
	ProductiveLimit        int        `json:"productive_limit"`
	MinimumAccepted        int        `json:"minimum_discriminators_accepted"`
	RequiresSameFile       bool       `json:"requires_same_file"`
	RejectsSingle          bool       `json:"rejects_single_discriminator"`
	RejectsSplit           bool       `json:"rejects_split_files"`
	MetadataCounts         bool       `json:"metadata_counts"`
	PreservesScopeLanguage bool       `json:"preserves_scope_language"`
	UsesLegacyPlanner      bool       `json:"uses_legacy_planner"`
	GenericFallback        bool       `json:"generic_fallback"`
	AllTechnicalExecuted   bool       `json:"all_technical_executed"`
	ProviderAttempts       int        `json:"provider_attempts"`
	UnionOrigins           []string   `json:"union_origins"`
	LegacyByteIdentical    bool       `json:"legacy_byte_identical"`
	ResponseStatus         string     `json:"response_status"`
	ErrorReturned          bool       `json:"error_returned"`
	EmptyOnError           bool       `json:"empty_on_error"`
	EvidenceAddressable    bool       `json:"evidence_addressable"`
	CapabilityProven       bool       `json:"capability_proven"`
	InputMutated           bool       `json:"input_mutated"`
}

type metrics struct {
	RecallAt5                 *float64 `json:"recall_at_5"`
	ContaminationAt5          *float64 `json:"contamination_at_5"`
	EvidenceAddressable       int      `json:"evidence_addressable_capabilities"`
	ProviderAttempts          int      `json:"provider_attempts"`
	TechnicalOriginsForTarget int      `json:"technical_origins_for_target"`
}

type utilityControl struct {
	ID               string  `json:"id"`
	Class            string  `json:"class"`
	BindingComplete  bool    `json:"binding_complete"`
	ReplayAvailable  bool    `json:"replay_available"`
	Baseline         metrics `json:"baseline"`
	Candidate        metrics `json:"candidate"`
	ImprovedCases    int     `json:"improved_cases"`
	ImprovedFamilies int     `json:"improved_families"`
	ForbiddenQuery   bool    `json:"forbidden_query"`
	Error            bool    `json:"error"`
	BudgetExceeded   bool    `json:"budget_exceeded"`
	Regression       bool    `json:"regression"`
	CapabilityProven bool    `json:"capability_proven"`
	ClaimedVerdict   string  `json:"claimed_verdict"`
}

func loadFixture(t *testing.T) fixture {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve gate source")
	}
	path := filepath.Join(filepath.Dir(source), "..", "..", "..", "..", "lab", "bancada", "history", "benchmarks", "0068", "adversarial-fixtures.json")
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
	return got
}

func patchedClaim(t *testing.T, base compilerClaim, patch map[string]json.RawMessage) compilerClaim {
	t.Helper()
	raw, err := json.Marshal(base)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		t.Fatal(err)
	}
	for key, value := range patch {
		if _, known := object[key]; !known {
			t.Fatalf("mutation patch uses unknown claim field %q", key)
		}
		object[key] = value
	}
	raw, err = json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	var claim compilerClaim
	if err := json.Unmarshal(raw, &claim); err != nil {
		t.Fatal(err)
	}
	return claim
}

func evaluateCompiler(inputState string, claim compilerClaim, opacity string) (string, []string) {
	if inputState == "opaque" {
		if opacity == "" {
			return violated, []string{"opacity_reason"}
		}
		return unknown, []string{opacity}
	}
	defects := make([]string, 0)
	check := func(ok bool, code string) {
		if !ok {
			defects = append(defects, code)
		}
	}
	check(!claim.CapabilityProven, "no_proven_claim")
	check(!claim.InputMutated, "input_immutability")
	switch inputState {
	case "canonical":
		check(claim.SourcePhrase == "predicados geometricos de orientacao e in-circle robustos a casos quase degenerados", "source_phrase")
		check(claim.Translation == canonicalQueries[0], "translation")
		check(claim.LiteralPreserved && claim.Literal == claim.SourcePhrase, "literal_preserved")
		check(claim.CatalogMatched, "finite_catalog")
		check(equalStrings(claim.ProposalQueries, canonicalQueries), "proposal_queries")
		check(equalStrings(claim.ProposalStrategies, []string{"technical-lexicon-v1", "technical-lexicon-v1", "technical-lexicon-v1"}), "strategy")
		check(equalStrings(claim.ProposalRuleIDs, canonicalRules) && equalStrings(claim.ProposalSourceIDs, []string{"exp-capability", "exp-capability", "exp-capability"}), "provenance")
		check(equalNestedStrings(claim.ProposalDiscriminators, canonicalDiscriminators), "discriminators")
		check(claim.ProductiveLimit == 3 && len(claim.ProposalQueries) <= 3, "proposal_limit")
		for _, query := range claim.ProposalQueries {
			check(len([]byte(query)) <= 256, "query_budget")
			check(query != "predicados geometricos" && query != "geometric predicates", "forbidden_query")
		}
		for _, pair := range claim.ProposalDiscriminators {
			check(len(pair) == 2 && pair[0] != "" && pair[1] != "" && !strings.EqualFold(pair[0], pair[1]), "discriminators")
			for _, item := range pair {
				check(len([]byte(item)) <= 64, "discriminator_budget")
			}
		}
		check(claim.MinimumAccepted == 2, "discriminator_gate")
		check(claim.RequiresSameFile && claim.RejectsSingle && claim.RejectsSplit, "same_file")
		check(!claim.MetadataCounts, "metadata")
		check(claim.PreservesScopeLanguage && !claim.UsesLegacyPlanner && !claim.GenericFallback, "technical_request")
		check(claim.AllTechnicalExecuted, "all_executed")
		check(claim.ProviderAttempts == 4, "provider_attempts")
		check(equalStrings(claim.UnionOrigins, []string{"orient2d incircle", "Shewchuk exact predicates"}), "union_origins")
		check(claim.ResponseStatus == "SUCCEEDED" && !claim.ErrorReturned && !claim.EmptyOnError, "canonical_status")
		check(claim.EvidenceAddressable, "evidence_addressable")
	case "partial":
		check(!claim.CatalogMatched && len(claim.ProposalQueries) == 0, "finite_catalog")
	case "legacy":
		check(claim.LegacyByteIdentical, "legacy_identity")
		check(!claim.CatalogMatched && len(claim.ProposalQueries) == 0, "legacy_catalog_boundary")
	case "error":
		check(claim.ResponseStatus == "FAILED" && claim.ErrorReturned && !claim.EmptyOnError, "error_propagation")
	default:
		return unknown, []string{"unsupported_input_state"}
	}
	if len(defects) > 0 {
		sort.Strings(defects)
		return violated, defects
	}
	return preserved, []string{"all_public_invariants_preserved"}
}

func evaluateUtility(c utilityControl) (string, string) {
	if !c.BindingComplete || !c.ReplayAvailable || c.Baseline.RecallAt5 == nil || c.Baseline.ContaminationAt5 == nil || c.Candidate.RecallAt5 == nil || c.Candidate.ContaminationAt5 == nil || c.Error || c.BudgetExceeded {
		return unknown, "required metric, replay, or result binding unavailable"
	}
	expected := "NOT_OBSERVED"
	recallDelta := *c.Candidate.RecallAt5 - *c.Baseline.RecallAt5
	contaminationDelta := *c.Candidate.ContaminationAt5 - *c.Baseline.ContaminationAt5
	evidenceDelta := c.Candidate.EvidenceAddressable - c.Baseline.EvidenceAddressable
	guards := !c.ForbiddenQuery && !c.Regression && !c.CapabilityProven && c.Candidate.ProviderAttempts <= 4
	if recallDelta >= 1.0-1e-12 && contaminationDelta <= 1e-12 && evidenceDelta >= 1 && c.Candidate.TechnicalOriginsForTarget >= 2 && c.ImprovedCases >= 2 && c.ImprovedFamilies >= 2 && guards {
		expected = "OBSERVED"
	}
	if c.ClaimedVerdict != expected {
		return violated, fmt.Sprintf("claimed %s, external contract requires %s", c.ClaimedVerdict, expected)
	}
	return preserved, "external verdict " + expected
}

func TestStep0068FixtureShapeAndMutationScore(t *testing.T) {
	f := loadFixture(t)
	if f.SchemaVersion != "tekt.step0068.binding-free-fixtures.v1" || f.Step != "0068" || len(f.ManifestSHA256) != 64 {
		t.Fatalf("invalid fixture identity: %+v", f)
	}
	if f.RequiredRepetitions != 20 || len(f.Mutations) != 16 {
		t.Fatalf("repetitions=%d mutations=%d, want 20 and 16", f.RequiredRepetitions, len(f.Mutations))
	}
	killed := 0
	for index, mutation := range f.Mutations {
		wantID := fmt.Sprintf("M%02d", index+1)
		if mutation.ID != wantID || mutation.Class != "N" || mutation.Witness == "" || mutation.ExpectedDefect == "" {
			t.Fatalf("invalid mutation %d: %+v", index, mutation)
		}
		claim := patchedClaim(t, f.ReferenceClaim, mutation.Patch)
		classification, defects := evaluateCompiler(mutation.InputState, claim, "")
		if classification != violated || !contains(defects, mutation.ExpectedDefect) {
			t.Fatalf("%s survived as %s defects=%v, expected witness %s", mutation.ID, classification, defects, mutation.ExpectedDefect)
		}
		killed++
	}
	score := float64(killed) / float64(len(f.Mutations))
	if math.Abs(score-1.0) > 1e-12 {
		t.Fatalf("mutation_score=%v, want 1.0", score)
	}
}

func TestStep0068PositiveNegativeOpaqueAndExternalMetrics(t *testing.T) {
	f := loadFixture(t)
	seen := map[string]int{}
	for _, c := range f.Controls {
		classification, defects := evaluateCompiler(c.InputState, c.Claim, c.Opacity)
		seen[c.Class]++
		want := map[string]string{"P": preserved, "O": unknown}[c.Class]
		if classification != want {
			t.Errorf("%s = %s (%v), want %s", c.ID, classification, defects, want)
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
	for _, verdict := range []string{"OBSERVED", "NOT_OBSERVED"} {
		if !verdicts[verdict] {
			t.Errorf("external verdict %s was not positively discriminated", verdict)
		}
	}
	observed := f.UtilityControls[0]
	if *observed.Baseline.RecallAt5 != 0 || *observed.Candidate.RecallAt5 != 1 || *observed.Baseline.ContaminationAt5 != 1 || *observed.Candidate.ContaminationAt5 != .5 || observed.Baseline.EvidenceAddressable != 0 || observed.Candidate.EvidenceAddressable != 1 || observed.Baseline.ProviderAttempts != 1 || observed.Candidate.ProviderAttempts != 4 || observed.Candidate.TechnicalOriginsForTarget != 2 {
		t.Fatalf("frozen A/B metrics drifted: %+v", observed)
	}
}

func TestStep0068TwentyOrdersAreIdentical(t *testing.T) {
	f := loadFixture(t)
	baseline := classificationVector(t, f, append([]mutation(nil), f.Mutations...))
	for repetition := 0; repetition < f.RequiredRepetitions; repetition++ {
		ordered := rotate(f.Mutations, repetition%len(f.Mutations))
		if repetition%2 == 1 {
			reverse(ordered)
		}
		if got := classificationVector(t, f, ordered); !equalMap(got, baseline) {
			t.Fatalf("order/repetition %02d changed classification vector", repetition+1)
		}
	}
}

func classificationVector(t *testing.T, f fixture, mutations []mutation) map[string]string {
	t.Helper()
	result := make(map[string]string, len(mutations)+len(f.Controls)+len(f.UtilityControls))
	for _, c := range f.Controls {
		result[c.ID], _ = evaluateCompiler(c.InputState, c.Claim, c.Opacity)
	}
	for _, m := range mutations {
		result[m.ID], _ = evaluateCompiler(m.InputState, patchedClaim(t, f.ReferenceClaim, m.Patch), "")
	}
	for _, c := range f.UtilityControls {
		result[c.ID], _ = evaluateUtility(c)
	}
	return result
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

func equalNestedStrings(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !equalStrings(a[i], b[i]) {
			return false
		}
	}
	return true
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func rotate(in []mutation, n int) []mutation {
	out := append([]mutation{}, in[n:]...)
	return append(out, in[:n]...)
}

func reverse(in []mutation) {
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
