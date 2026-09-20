/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o4_predicates.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type O4Execution struct {
	SemanticVariantID     string      `json:"semantic_variant_id"`
	EntryCommitmentSHA256 string      `json:"entry_commitment_sha256"`
	MutationStage         string      `json:"mutation_stage"`
	InputSHA256           string      `json:"input_sha256"`
	Observation           Observation `json:"observation"`
}

type O4PredicateFinding struct {
	ID                  string   `json:"id"`
	Outcome             string   `json:"outcome"`
	Relation            string   `json:"relation"`
	ChangedFields       []string `json:"changed_fields"`
	EntryHashVerified   bool     `json:"entry_hash_verified"`
	InputMutationProven bool     `json:"input_mutation_proven"`
}

type O4ExecutableSubject interface {
	Execute(context.Context, Case, O4Commitment) (O4Execution, error)
}

type O4ReferenceSubject struct{}

func (O4ReferenceSubject) Execute(ctx context.Context, testCase Case, commitment O4Commitment) (O4Execution, error) {
	observation, err := ReferenceSubject{}.Observe(ctx, testCase)
	return O4Execution{
		SemanticVariantID: commitment.SemanticVariantID, EntryCommitmentSHA256: commitment.CanonicalAttackEntrySHA256,
		MutationStage: "NONE", InputSHA256: o4InputHash(testCase.Input), Observation: observation,
	}, err
}

type O4OpaqueSubject struct{}

func (O4OpaqueSubject) Execute(_ context.Context, testCase Case, commitment O4Commitment) (O4Execution, error) {
	if !testCase.Input.Opaque || testCase.Input.OpacityReason == "" {
		return O4Execution{}, fmt.Errorf("opaque input is not explicit")
	}
	observation := frozenProjection(testCase.Input)
	observation.Classification = "Unknown"
	return O4Execution{
		SemanticVariantID: commitment.SemanticVariantID, EntryCommitmentSHA256: commitment.CanonicalAttackEntrySHA256,
		MutationStage: "OPAQUE_INPUT_CLASSIFICATION", InputSHA256: o4InputHash(testCase.Input), Observation: observation,
	}, nil
}

func EvaluateO4Negative(commitment O4Commitment, frozen Observation, baseline, variant O4Execution) O4PredicateFinding {
	finding := O4PredicateFinding{ID: commitment.ID, Relation: commitment.ViolatedPredicate.RequiredRelation}
	baselineChanges := changedRequiredFields(frozen, baseline.Observation, commitment.RequiredObservableFields)
	variantChanges := changedRequiredFields(baseline.Observation, variant.Observation, commitment.RequiredObservableFields)
	finding.ChangedFields = variantChanges
	finding.EntryHashVerified = variant.EntryCommitmentSHA256 == commitment.CanonicalAttackEntrySHA256 && variant.SemanticVariantID == commitment.SemanticVariantID
	finding.InputMutationProven = variant.MutationStage == "INPUT_BEFORE_OBSERVATION" && variant.InputSHA256 != baseline.InputSHA256
	baselinePredicate := commitment.BaselinePredicate.ExpectedOutcome == "Preserved" && commitment.BaselinePredicate.RequiredRelation == "ALL_REQUIRED_OBSERVABLES_EQUAL_FROZEN_CONTRACT" && len(baselineChanges) == 0
	violatedPredicate := commitment.ViolatedPredicate.ExpectedOutcome == "Violated" && commitment.ViolatedPredicate.SubjectKind == "INDEPENDENT_SEMANTIC_VARIANT" && commitment.ViolatedPredicate.SemanticVariantID == commitment.SemanticVariantID && commitment.ViolatedPredicate.RequiredRelation == "INPUT_DRIVEN_OBLIGATION_SPECIFIC_DEVIATION" && commitment.ViolatedPredicate.ChangedFieldCardinality == len(variantChanges) && len(variantChanges) == 1
	if baselinePredicate && violatedPredicate && finding.EntryHashVerified && finding.InputMutationProven {
		finding.Outcome = "Violated"
	} else {
		finding.Outcome = "Survived"
	}
	return finding
}

func EvaluateO4Positive(commitment O4Commitment, frozen Observation, execution O4Execution) O4PredicateFinding {
	changes := changedRequiredFields(frozen, execution.Observation, commitment.RequiredObservableFields)
	finding := O4PredicateFinding{ID: commitment.ID, Relation: commitment.BaselinePredicate.RequiredRelation, ChangedFields: changes}
	finding.EntryHashVerified = execution.EntryCommitmentSHA256 == commitment.CanonicalAttackEntrySHA256 && execution.SemanticVariantID == commitment.SemanticVariantID
	if commitment.BaselinePredicate.ExpectedOutcome == "Preserved" && commitment.BaselinePredicate.SubjectKind == "INDEPENDENT_REFERENCE_CONTROL" && commitment.BaselinePredicate.RequiredRelation == "ALL_REQUIRED_OBSERVABLES_EQUAL_FROZEN_CONTRACT" && len(changes) == 0 && finding.EntryHashVerified {
		finding.Outcome = "Preserved"
	} else {
		finding.Outcome = "Violated"
	}
	return finding
}

func EvaluateO4Opaque(commitment O4Commitment, execution O4Execution) O4PredicateFinding {
	finding := O4PredicateFinding{ID: commitment.ID, Relation: commitment.BaselinePredicate.RequiredRelation}
	finding.EntryHashVerified = execution.EntryCommitmentSHA256 == commitment.CanonicalAttackEntrySHA256 && execution.SemanticVariantID == commitment.SemanticVariantID
	executable := commitment.BaselinePredicate.Executable != nil && *commitment.BaselinePredicate.Executable && commitment.ViolatedPredicate.Executable != nil && *commitment.ViolatedPredicate.Executable
	if commitment.BaselinePredicate.ExpectedOutcome == "Unknown" && commitment.ViolatedPredicate.ExpectedOutcome == "Unknown" && commitment.BaselinePredicate.RequiredRelation == "EXPLICIT_OPACITY_CLASSIFICATION" && commitment.ViolatedPredicate.RequiredRelation == "NOT_COERCED_TO_PRESERVED_OR_VIOLATED" && executable && execution.Observation.Classification == "Unknown" && finding.EntryHashVerified && execution.MutationStage == "OPAQUE_INPUT_CLASSIFICATION" {
		finding.Outcome = "Unknown"
	} else {
		finding.Outcome = "Violated"
	}
	return finding
}

func changedRequiredFields(want, got Observation, fields []string) []string {
	wantFields := observationJSONFields(want)
	gotFields := observationJSONFields(got)
	changed := make([]string, 0, len(fields))
	for _, field := range fields {
		if !bytes.Equal(wantFields[field], gotFields[field]) {
			changed = append(changed, field)
		}
	}
	return changed
}

func observationJSONFields(observation Observation) map[string]json.RawMessage {
	document, _ := json.Marshal(observation)
	fields := map[string]json.RawMessage{}
	_ = json.Unmarshal(document, &fields)
	return fields
}

func o4InputHash(input Input) string {
	document, _ := json.Marshal(input)
	digest := sha256.Sum256(document)
	return "sha256:" + hex.EncodeToString(digest[:])
}
