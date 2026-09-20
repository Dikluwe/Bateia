/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0070_gate/discriminatory_test.md
 * @layer L4
 * @updated 2026-09-12
 */

package step0070_gate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

type fixtureSet struct {
	SchemaVersion          string            `json:"schema_version"`
	CreatedBeforeCandidate bool              `json:"created_before_candidate"`
	BindingFree            bool              `json:"binding_free"`
	FrozenInputs           map[string]string `json:"frozen_inputs"`
	Requirements           requirements      `json:"requirements"`
	Cost                   fixtureCost       `json:"cost"`
	Cases                  []fixtureCase     `json:"cases"`
}

type requirements struct {
	MutationScore float64 `json:"mutation_score"`
	Repetitions   int     `json:"repetitions"`
	Orders        int     `json:"orders"`
}

type fixtureCost struct {
	Cases                     int `json:"cases"`
	Controls                  int `json:"controls"`
	Opaque                    int `json:"opaque"`
	ValidMutations            int `json:"valid_mutations"`
	ClassificationEvaluations int `json:"classification_evaluations"`
	ExternalProcesses         int `json:"external_processes"`
}

type fixtureCase struct {
	ID          string      `json:"id"`
	Kind        string      `json:"kind"`
	Expected    string      `json:"expected"`
	Observation observation `json:"observation"`
}

type observation struct {
	Status                  string   `json:"status"`
	Reason                  string   `json:"reason"`
	AssessmentCount         int      `json:"assessment_count"`
	LimitExceeded           bool     `json:"limit_exceeded"`
	DeclaredClaimIDs        []string `json:"declared_claim_ids"`
	EmittedClaimIDs         []string `json:"emitted_claim_ids"`
	DeclaredRelation        string   `json:"declared_relation"`
	ReportedRelation        string   `json:"reported_relation"`
	Outcome                 string   `json:"outcome"`
	SupportingEvidenceIDs   []string `json:"supporting_evidence_ids"`
	ConflictingEvidenceIDs  []string `json:"conflicting_evidence_ids"`
	AmbiguousEvidenceIDs    []string `json:"ambiguous_evidence_ids"`
	SupportingInferenceIDs  []string `json:"supporting_inference_ids"`
	ConflictingInferenceIDs []string `json:"conflicting_inference_ids"`
	AmbiguousInferenceIDs   []string `json:"ambiguous_inference_ids"`
	PromotedInferenceIDs    []string `json:"promoted_inference_ids"`
	SupportingOriginIDs     []string `json:"supporting_origin_ids"`
	OrphanReferences        []string `json:"orphan_references"`
	ClaimsWithoutEvidence   []string `json:"claims_without_evidence"`
	OrderFingerprints       []string `json:"order_fingerprints"`
	Effects                 []string `json:"effects"`
	InputSHA256Before       string   `json:"input_sha256_before"`
	InputSHA256After        string   `json:"input_sha256_after"`
}

type classification struct {
	Verdict string
	Witness string
}

func TestStep0070DiscriminatoryContract(t *testing.T) {
	fixtures := loadFixtures(t)
	validateFixtureEnvelope(t, fixtures)
	verifyFrozenInputs(t, fixtures.FrozenInputs)

	baseline := make(map[string]classification, len(fixtures.Cases))
	controls, opaque, mutations, rejected := 0, 0, 0, 0
	for _, item := range fixtures.Cases {
		result := classify(item.Observation)
		baseline[item.ID] = result
		if result.Verdict != item.Expected {
			t.Fatalf("case %s: got %s (%s), want %s", item.ID, result.Verdict, result.Witness, item.Expected)
		}
		switch item.Kind {
		case "control":
			controls++
		case "opaque":
			opaque++
		case "mutation":
			mutations++
			if result.Verdict == "Violated" {
				rejected++
			}
		default:
			t.Fatalf("case %s has unknown kind %q", item.ID, item.Kind)
		}
	}

	if controls != fixtures.Cost.Controls || opaque != fixtures.Cost.Opaque || mutations != fixtures.Cost.ValidMutations {
		t.Fatalf("fixture cost mismatch: controls=%d opaque=%d mutations=%d", controls, opaque, mutations)
	}
	score := float64(rejected) / float64(mutations)
	if score != fixtures.Requirements.MutationScore || score != 1.0 {
		t.Fatalf("mutation score %.3f, want 1.0", score)
	}

	evaluations := 0
	seenOrders := make(map[string]struct{}, fixtures.Requirements.Orders)
	for repetition := 0; repetition < fixtures.Requirements.Repetitions; repetition++ {
		for order := 0; order < fixtures.Requirements.Orders; order++ {
			indices := deterministicOrder(len(fixtures.Cases), order)
			keyBytes, _ := json.Marshal(indices)
			seenOrders[string(keyBytes)] = struct{}{}
			for _, index := range indices {
				item := fixtures.Cases[index]
				result := classify(item.Observation)
				evaluations++
				if result != baseline[item.ID] {
					t.Fatalf("case %s changed at repetition %d order %d: %#v != %#v", item.ID, repetition, order, result, baseline[item.ID])
				}
			}
		}
	}
	if len(seenOrders) != fixtures.Requirements.Orders {
		t.Fatalf("exercised %d distinct orders, want %d", len(seenOrders), fixtures.Requirements.Orders)
	}
	if evaluations != fixtures.Cost.ClassificationEvaluations {
		t.Fatalf("classification cost %d, want %d", evaluations, fixtures.Cost.ClassificationEvaluations)
	}
}

func loadFixtures(t *testing.T) fixtureSet {
	t.Helper()
	path := filepath.Join("..", "..", "..", "..", "00_nucleo", "benchmark", "0070", "adversarial-fixtures.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fixtures fixtureSet
	decoderErr := json.Unmarshal(data, &fixtures)
	if decoderErr != nil {
		t.Fatal(decoderErr)
	}
	return fixtures
}

func validateFixtureEnvelope(t *testing.T, fixtures fixtureSet) {
	t.Helper()
	if fixtures.SchemaVersion != "algorithm-history-adversarial-fixtures-v1" || !fixtures.CreatedBeforeCandidate || !fixtures.BindingFree {
		t.Fatal("fixture envelope is not frozen and binding-free")
	}
	if fixtures.Requirements.MutationScore != 1.0 || fixtures.Requirements.Repetitions != 2 || fixtures.Requirements.Orders != 20 {
		t.Fatalf("invalid discriminatory requirements: %#v", fixtures.Requirements)
	}
	if fixtures.Cost.Cases != len(fixtures.Cases) || fixtures.Cost.ExternalProcesses != 0 {
		t.Fatalf("invalid recorded cost: %#v", fixtures.Cost)
	}
}

func verifyFrozenInputs(t *testing.T, frozen map[string]string) {
	t.Helper()
	paths := make([]string, 0, len(frozen))
	for path := range frozen {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join("..", "..", "..", "..", filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("read frozen input %s: %v", path, err)
		}
		digest := sha256.Sum256(data)
		actual := hex.EncodeToString(digest[:])
		if actual != frozen[path] {
			t.Fatalf("frozen input %s changed: %s != %s", path, actual, frozen[path])
		}
	}
}

func classify(value observation) classification {
	if len(value.Effects) > 0 {
		return classification{"Violated", "forbidden_io:" + value.Effects[0]}
	}
	if value.InputSHA256Before != value.InputSHA256After {
		return classification{"Violated", "input_mutated"}
	}
	if !allEqual(value.OrderFingerprints) {
		return classification{"Violated", "order_dependent"}
	}
	if value.Status == "UNKNOWN" {
		if value.AssessmentCount != 0 {
			return classification{"Violated", "partial_unknown_output"}
		}
		if isGlobalUnknownReason(value.Reason) {
			return classification{"Unknown", value.Reason}
		}
		return classification{"Violated", "unknown_without_public_reason"}
	}
	if value.Status != "COMPLETE" || value.Reason != "analyzed" {
		return classification{"Violated", "invalid_global_status"}
	}
	if len(value.OrphanReferences) > 0 {
		return classification{"Violated", "orphan_reference:" + value.OrphanReferences[0]}
	}
	if len(value.ClaimsWithoutEvidence) > 0 {
		return classification{"Violated", "claim_without_evidence:" + value.ClaimsWithoutEvidence[0]}
	}
	if invented := firstMissing(value.EmittedClaimIDs, value.DeclaredClaimIDs); invented != "" {
		return classification{"Violated", "invented_relation:" + invented}
	}
	if value.ReportedRelation != value.DeclaredRelation {
		return classification{"Violated", fmt.Sprintf("relation_promoted:%s_to_%s", value.DeclaredRelation, value.ReportedRelation)}
	}
	if len(value.PromotedInferenceIDs) > 0 || overlaps(value.SupportingEvidenceIDs, value.SupportingInferenceIDs) ||
		overlaps(value.ConflictingEvidenceIDs, value.ConflictingInferenceIDs) || overlaps(value.AmbiguousEvidenceIDs, value.AmbiguousInferenceIDs) {
		return classification{"Violated", "inference_promoted_to_evidence"}
	}
	if value.AssessmentCount != len(value.EmittedClaimIDs) {
		return classification{"Violated", "assessment_cardinality"}
	}
	expected := evidenceOutcome(value)
	if value.Outcome != expected {
		return classification{"Violated", fmt.Sprintf("evidence_outcome:%s_not_%s", value.Outcome, expected)}
	}
	return classification{"Preserved", "public_contract_preserved"}
}

func evidenceOutcome(value observation) string {
	supports := len(value.SupportingEvidenceIDs)
	conflicts := len(value.ConflictingEvidenceIDs)
	ambiguous := len(value.AmbiguousEvidenceIDs)
	switch {
	case supports > 0 && conflicts > 0:
		return "CONTESTED"
	case ambiguous > 0:
		return "UNKNOWN"
	case conflicts > 0:
		return "CONTRADICTED"
	case supports > 0 && uniqueCount(value.SupportingOriginIDs) >= 2:
		return "CORROBORATED"
	case supports > 0 && uniqueCount(value.SupportingOriginIDs) == 1:
		return "SUPPORTED"
	default:
		return "UNKNOWN"
	}
}

func uniqueCount(values []string) int {
	unique := make(map[string]struct{}, len(values))
	for _, value := range values {
		unique[value] = struct{}{}
	}
	return len(unique)
}

func firstMissing(values, allowed []string) string {
	set := make(map[string]struct{}, len(allowed))
	for _, value := range allowed {
		set[value] = struct{}{}
	}
	for _, value := range values {
		if _, ok := set[value]; !ok {
			return value
		}
	}
	return ""
}

func overlaps(left, right []string) bool {
	set := make(map[string]struct{}, len(left))
	for _, value := range left {
		set[value] = struct{}{}
	}
	for _, value := range right {
		if _, ok := set[value]; ok {
			return true
		}
	}
	return false
}

func allEqual(values []string) bool {
	for index := 1; index < len(values); index++ {
		if values[index] != values[0] {
			return false
		}
	}
	return true
}

func isGlobalUnknownReason(reason string) bool {
	switch reason {
	case "invalid_schema", "budget_exceeded", "no_algorithms", "no_provenance", "no_claims", "no_bases", "invalid_input", "ambiguous_input":
		return true
	default:
		return false
	}
}

func deterministicOrder(size, order int) []int {
	result := make([]int, size)
	for index := range result {
		result[index] = index
	}
	if order >= size {
		for left, right := 0, size-1; left < right; left, right = left+1, right-1 {
			result[left], result[right] = result[right], result[left]
		}
		order -= size
	}
	rotated := make([]int, size)
	for index := range result {
		rotated[index] = result[(index+order)%size]
	}
	return rotated
}
