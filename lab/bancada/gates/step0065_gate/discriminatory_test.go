/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0065_gate/discriminatory_test.md
 * @layer L4
 * @updated 2026-09-11
 */
package step0065gate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type suite0065 struct {
	SchemaVersion          string            `json:"schema_version"`
	CreatedBeforeCandidate bool              `json:"created_before_candidate"`
	Scope                  string            `json:"scope"`
	FrozenInputs           map[string]string `json:"frozen_inputs"`
	Controls               []control0065     `json:"controls"`
	Mutations              []mutation0065    `json:"mutations"`
	Determinism            determinism0065   `json:"determinism"`
}

type control0065 struct {
	ID             string          `json:"id"`
	Class          string          `json:"class"`
	BaselineStatus string          `json:"baseline_status"`
	ExpectedGate   string          `json:"expected_gate"`
	Expected       json.RawMessage `json:"expected"`
}

type mutation0065 struct {
	ID                string         `json:"id"`
	Rule              string         `json:"rule"`
	Stimulus          string         `json:"stimulus"`
	Reference         map[string]any `json:"reference"`
	Mutant            map[string]any `json:"mutant"`
	ExpectedDiffPaths []string       `json:"expected_diff_paths"`
}

type determinism0065 struct {
	RequiredRepetitions   int      `json:"required_repetitions"`
	PermutationJSONHashes []string `json:"permutation_json_hashes"`
	NilJSONHash           string   `json:"nil_json_hash"`
	EmptyJSONHash         string   `json:"empty_json_hash"`
}

type expectedObservation0065 struct {
	ReportStatus          string   `json:"report_status"`
	ReportReason          string   `json:"report_reason"`
	Outcome               string   `json:"outcome"`
	Independence          string   `json:"independence"`
	DecisionImpact        string   `json:"decision_impact"`
	SupportingWitnessIDs  []string `json:"supporting_witness_ids"`
	ConflictingWitnessIDs []string `json:"conflicting_witness_ids"`
	MissingRoles          []string `json:"missing_roles"`
	DistinctOrigins       []string `json:"distinct_origins"`
	Operations            int      `json:"operations"`
	VerdictChanges        int      `json:"verdict_changes"`
}

type mutationResult0065 struct {
	Verdict string
	Witness string
}

func TestStep0065PreCandidateDiscriminatoryGate(t *testing.T) {
	suite := loadSuite0065(t)
	checkFrozenInputs0065(t, suite.FrozenInputs)
	checkControls0065(t, suite.Controls)
	checkDeterminism0065(t, suite.Determinism)

	if len(suite.Mutations) != 12 {
		t.Fatalf("valid mutation denominator: want 12, got %d", len(suite.Mutations))
	}
	byID := make(map[string]mutation0065, len(suite.Mutations))
	for _, mutation := range suite.Mutations {
		if _, exists := byID[mutation.ID]; exists {
			t.Fatalf("duplicate mutation %s", mutation.ID)
		}
		byID[mutation.ID] = mutation
	}

	baseline := make(map[string]mutationResult0065, 12)
	rejected := 0
	for repetition := 0; repetition < suite.Determinism.RequiredRepetitions; repetition++ {
		for _, id := range mutationOrder0065(repetition) {
			mutation, ok := byID[id]
			if !ok {
				t.Fatalf("missing required mutation %s", id)
			}
			got := evaluateMutation0065(t, mutation)
			if prior, exists := baseline[id]; exists && prior != got {
				t.Fatalf("unstable mutation %s: %#v != %#v", id, prior, got)
			}
			baseline[id] = got
			if repetition == 0 && got.Verdict == "Violated" {
				rejected++
			}
		}
	}
	if rejected != len(suite.Mutations) {
		t.Fatalf("mutation_score=%d/%d; surviving mutation", rejected, len(suite.Mutations))
	}
	t.Logf("mutation_score=%.1f (%d/%d), repetitions=%d", float64(rejected)/float64(len(suite.Mutations)), rejected, len(suite.Mutations), suite.Determinism.RequiredRepetitions)
}

func loadSuite0065(t *testing.T) suite0065 {
	t.Helper()
	path := filepath.Join("..", "..", "..", "..", "lab", "bancada", "history", "benchmarks", "0065", "adversarial-cases.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var suite suite0065
	if err := json.Unmarshal(data, &suite); err != nil {
		t.Fatal(err)
	}
	if suite.SchemaVersion != "tekt-step0065-adversarial-v1" || !suite.CreatedBeforeCandidate {
		t.Fatalf("unexpected suite identity: %q pre_candidate=%v", suite.SchemaVersion, suite.CreatedBeforeCandidate)
	}
	if suite.Scope != "binding-free discriminatory model; no candidate import or execution" {
		t.Fatalf("unexpected suite scope %q", suite.Scope)
	}
	return suite
}

func checkFrozenInputs0065(t *testing.T, inputs map[string]string) {
	t.Helper()
	if len(inputs) != 5 {
		t.Fatalf("want five frozen inputs, got %d", len(inputs))
	}
	for relative, expected := range inputs {
		resolved := strings.TrimPrefix(relative, "00_nucleo/benchmark/")
		if resolved != relative {
			resolved = filepath.Join("lab", "bancada", "history", "benchmarks", filepath.FromSlash(resolved))
		}
		data, err := os.ReadFile(filepath.Join("..", "..", "..", "..", resolved))
		if err != nil {
			t.Fatalf("read frozen input %s: %v", relative, err)
		}
		digest := sha256.Sum256(data)
		got := hex.EncodeToString(digest[:])
		if got != expected {
			t.Fatalf("frozen input changed: %s want %s got %s", relative, expected, got)
		}
	}
}

func checkControls0065(t *testing.T, controls []control0065) {
	t.Helper()
	if len(controls) != 12 {
		t.Fatalf("want twelve contractual controls, got %d", len(controls))
	}
	wantIDs := []string{"N01", "N02", "N03", "N04", "N05", "O01", "O02", "O03", "P01", "P02", "P03", "P04"}
	gotIDs := make([]string, 0, len(controls))
	statuses := map[string]bool{}
	classes := map[string]int{}
	for _, control := range controls {
		gotIDs = append(gotIDs, control.ID)
		statuses[control.BaselineStatus] = true
		classes[control.Class]++
		var observation expectedObservation0065
		if err := json.Unmarshal(control.Expected, &observation); err != nil {
			t.Fatalf("control %s: %v", control.ID, err)
		}
		if observation.DecisionImpact != "NONE" || observation.VerdictChanges != 0 {
			t.Fatalf("control %s changes a legacy verdict", control.ID)
		}
		if observation.SupportingWitnessIDs == nil || observation.ConflictingWitnessIDs == nil || observation.MissingRoles == nil || observation.DistinctOrigins == nil {
			t.Fatalf("control %s has a null public slice", control.ID)
		}
		if !sorted0065(observation.SupportingWitnessIDs) || !sorted0065(observation.ConflictingWitnessIDs) || !sorted0065(observation.MissingRoles) || !sorted0065(observation.DistinctOrigins) {
			t.Fatalf("control %s has non-canonical public ordering", control.ID)
		}
		switch control.Class {
		case "positive":
			if control.ExpectedGate != "Preserved" || observation.Outcome != "SATISFIED" || observation.ReportStatus != "COMPLETE" {
				t.Fatalf("positive control %s is not preserved", control.ID)
			}
		case "negative":
			if control.ExpectedGate != "Preserved" || observation.Outcome == "SATISFIED" {
				t.Fatalf("negative control %s is not conservatively preserved", control.ID)
			}
		case "opaque":
			if control.ExpectedGate != "Unknown" || (observation.ReportStatus != "UNKNOWN" && observation.Outcome != "UNKNOWN_INDEPENDENCE") {
				t.Fatalf("opaque control %s did not remain Unknown", control.ID)
			}
		default:
			t.Fatalf("control %s has unknown class %q", control.ID, control.Class)
		}
	}
	sort.Strings(gotIDs)
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Fatalf("contractual controls: want %v got %v", wantIDs, gotIDs)
	}
	wantStatuses := []string{"ABSENT", "CONFLICT", "INCONCLUSIVE", "PARTIAL", "PROVEN"}
	for _, status := range wantStatuses {
		if !statuses[status] {
			t.Fatalf("missing baseline EvidenceStatus %s", status)
		}
	}
	if classes["positive"] == 0 || classes["negative"] == 0 || classes["opaque"] == 0 {
		t.Fatalf("missing case class: %#v", classes)
	}
}

func checkDeterminism0065(t *testing.T, determinism determinism0065) {
	t.Helper()
	if determinism.RequiredRepetitions != 20 || len(determinism.PermutationJSONHashes) != 20 {
		t.Fatalf("want 20 deterministic repetitions and hashes")
	}
	for i, digest := range determinism.PermutationJSONHashes {
		if digest != determinism.PermutationJSONHashes[0] {
			t.Fatalf("permutation %d changed canonical JSON hash", i)
		}
	}
	if determinism.NilJSONHash != determinism.EmptyJSONHash || determinism.NilJSONHash != determinism.PermutationJSONHashes[0] {
		t.Fatal("nil, empty, or permutation serialization differs")
	}
}

func evaluateMutation0065(t *testing.T, mutation mutation0065) mutationResult0065 {
	t.Helper()
	if mutation.Rule == "" || mutation.Stimulus == "" || len(mutation.Reference) == 0 || len(mutation.Mutant) == 0 {
		t.Fatalf("incomplete mutation %s", mutation.ID)
	}
	diffs := diffPaths0065(mutation.Reference, mutation.Mutant)
	want := append([]string(nil), mutation.ExpectedDiffPaths...)
	sort.Strings(want)
	if len(want) == 0 || !reflect.DeepEqual(diffs, want) {
		t.Fatalf("mutation %s: observable diffs want %v got %v", mutation.ID, want, diffs)
	}
	checkSpecialMutation0065(t, mutation)
	payload, err := json.Marshal(struct {
		ID, Rule, Stimulus string
		Diffs              []string
		Reference, Mutant  map[string]any
	}{mutation.ID, mutation.Rule, mutation.Stimulus, diffs, mutation.Reference, mutation.Mutant})
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(payload)
	return mutationResult0065{Verdict: "Violated", Witness: "sha256:" + hex.EncodeToString(digest[:])}
}

func checkSpecialMutation0065(t *testing.T, mutation mutation0065) {
	t.Helper()
	switch mutation.ID {
	case "M06":
		if !allEqual0065(stringSlice0065(t, mutation.Reference["canonical_json_hashes"])) || allEqual0065(stringSlice0065(t, mutation.Mutant["canonical_json_hashes"])) {
			t.Fatal("M06 does not represent order-dependent serialization")
		}
	case "M11":
		if mutation.Reference["nil_json_hash"] != mutation.Reference["empty_json_hash"] || mutation.Mutant["nil_json_hash"] == mutation.Mutant["empty_json_hash"] {
			t.Fatal("M11 does not represent nil/empty divergence")
		}
	case "M12":
		refOps := int(mutation.Reference["operations"].(float64))
		wantOps := int(mutation.Reference["normalized_witnesses"].(float64) + mutation.Reference["required_role_pairs"].(float64))
		mutOps := int(mutation.Mutant["operations"].(float64))
		if refOps != wantOps || mutOps == wantOps || mutation.Mutant["enumerates_subsets"] != true {
			t.Fatal("M12 does not violate the public operation formula")
		}
	}
}

func diffPaths0065(reference, mutant map[string]any) []string {
	keys := make(map[string]struct{}, len(reference)+len(mutant))
	for key := range reference {
		keys[key] = struct{}{}
	}
	for key := range mutant {
		keys[key] = struct{}{}
	}
	diffs := make([]string, 0, len(keys))
	for key := range keys {
		if !reflect.DeepEqual(reference[key], mutant[key]) {
			diffs = append(diffs, key)
		}
	}
	sort.Strings(diffs)
	return diffs
}

func mutationOrder0065(repetition int) []string {
	ids := make([]string, 12)
	for i := 1; i <= 12; i++ {
		ids[i-1] = fmt.Sprintf("M%02d", i)
	}
	if repetition%3 == 1 {
		for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
			ids[i], ids[j] = ids[j], ids[i]
		}
	} else if repetition%3 == 2 {
		sort.Slice(ids, func(i, j int) bool {
			a := sha256.Sum256([]byte(fmt.Sprintf("step0065:%d:%s", repetition, ids[i])))
			b := sha256.Sum256([]byte(fmt.Sprintf("step0065:%d:%s", repetition, ids[j])))
			return strings.Compare(string(a[:]), string(b[:])) < 0
		})
	}
	return ids
}

func sorted0065(values []string) bool {
	return sort.StringsAreSorted(values)
}

func stringSlice0065(t *testing.T, value any) []string {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("want JSON array, got %T", value)
	}
	result := make([]string, len(items))
	for i, item := range items {
		result[i], ok = item.(string)
		if !ok {
			t.Fatalf("want string at %d, got %T", i, item)
		}
	}
	return result
}

func allEqual0065(values []string) bool {
	if len(values) == 0 {
		return false
	}
	for _, value := range values[1:] {
		if value != values[0] {
			return false
		}
	}
	return true
}
