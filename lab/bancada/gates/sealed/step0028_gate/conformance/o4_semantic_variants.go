/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o4_semantic_variants.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"fmt"
	"strings"
)

type O4SemanticVariantSubject struct{}

func (O4SemanticVariantSubject) Execute(_ context.Context, testCase Case, commitment O4Commitment) (O4Execution, error) {
	fault, err := o4SemanticFault(commitment.SemanticVariantID)
	if err != nil {
		return O4Execution{}, err
	}
	mutated := testCase
	mutated.Input.ScenarioVersion = "O4"
	mutated.Input.SemanticVariantID = commitment.SemanticVariantID
	mutated.Input.EntryCommitmentSHA256 = commitment.CanonicalAttackEntrySHA256
	mutated.Input.SemanticFault = fault
	observation, err := observeO4SemanticInput(mutated)
	return O4Execution{
		SemanticVariantID: commitment.SemanticVariantID, EntryCommitmentSHA256: commitment.CanonicalAttackEntrySHA256,
		MutationStage: "INPUT_BEFORE_OBSERVATION", InputSHA256: o4InputHash(mutated.Input), Observation: observation,
	}, err
}

func o4SemanticFault(semanticVariantID string) (string, error) {
	parts := strings.Split(semanticVariantID, ".")
	if len(parts) != 6 || parts[0] != "SV" || parts[1] != "O4" || parts[2] != "NEGATIVE" {
		return "", fmt.Errorf("invalid negative semantic variant %q", semanticVariantID)
	}
	return parts[4], nil
}

func observeO4SemanticInput(testCase Case) (Observation, error) {
	input := testCase.Input
	result := frozenProjection(input)
	switch testCase.Operation {
	case OperationExtract:
		decision := "MECHANISM_PRESERVED"
		status, novelty := "PROVEN", "CONVENTIONAL"
		if input.DocumentationOnly {
			status, novelty = "PARTIAL", "UNKNOWN"
		}
		result.Decision = o4FaultString(input, "DECISION", decision, "CONTROLLED_GAP")
		result.MechanismStatus, result.Novelty = status, novelty
		result.BindingVerified, result.ProvenanceComplete = true, true
		result.CanonicalSHA256 = canonicalDigest([]string{input.Owner, input.Repo, input.Commit, input.Path, input.Language, input.Ecosystem})

	case OperationHypothesis:
		status := input.MechanismStatus
		if status == "" {
			status = "PROVEN"
		}
		hypothesis := "PROPOSED"
		result.MechanismStatus = o4FaultString(input, "MECHANISM_STATUS", status, "INCONCLUSIVE")
		result.HypothesisStatus = o4FaultString(input, "HYPOTHESIS_STATUS", hypothesis, "SUPPORTED")
		result.Lifecycle = o4FaultString(input, "LIFECYCLE", input.Lifecycle, "SUCCEEDED")
		result.Novelty, result.TransferMode = input.ClaimedNovelty, input.RequestedTransfer
		result.AdoptionDisposition = "NOT_AN_ADOPTION_DECISION"
		if input.CanonicalItems != nil {
			result.CanonicalSHA256 = canonicalDigest(sortedCopy(input.CanonicalItems))
		}
		if input.ArtifactsPublished {
			result.ArtifactSHA256 = canonicalDigest(sortedCopy(input.ArtifactEntries))
		}

	case OperationStart:
		result.ErrorCode = o4FaultString(input, "ERROR_CODE", "", "INTEGRITY_ERROR")
		result.ResolveCalls, result.ReserveCalls, result.IDCalls, result.CreateCalls = 1, 1, 1, 1
		result.RunID = "run_reference"
		result.RunIDs = make([]string, input.ConcurrentStarts)
		for index := range result.RunIDs {
			result.RunIDs[index] = result.RunID
		}
		result.CanonicalSHA256 = canonicalDigest(sortedCopy(input.CanonicalItems))
		result.Lifecycle = o4FaultString(input, "LIFECYCLE", "QUEUED", "RUNNING")

	case OperationPagination:
		result.ErrorCode = o4FaultString(input, "ERROR_CODE", "", "LIMIT_EXCEEDED")

	case OperationRuntime:
		result.RuntimeUnique = o4FaultBool(input, "RUNTIME_UNIQUE", true, false)
		result.StoreIdentity, result.QueueIdentity = firstIdentity(input.StoreIDs), firstIdentity(input.QueueIDs)
		result.RootIdentity, result.SecretIdentity, result.WorkerIdentity = firstIdentity(input.RootIDs), firstIdentity(input.SecretIDs), firstIdentity(input.WorkerIDs)

	case OperationLifecycle:
		result.Lifecycle = o4FaultString(input, "LIFECYCLE", "CANCELLED", "SUCCEEDED")
		result.HypothesisStatus, result.AdoptionDisposition = "PROPOSED", "NOT_AN_ADOPTION_DECISION"
		result.RecoveryState = "INTERRUPTED"

	case OperationReplay:
		baseline := canonicalDigest(sortedCopy(input.ReplayItems))
		result.ReplaySHA256 = o4FaultString(input, "REPLAY_SHA256", baseline, "sha256:integrity-bypassed")

	case OperationSafety:
		result.Lifecycle = o4FaultString(input, "LIFECYCLE", input.Lifecycle, "SUCCEEDED")
		result.LeakedFields = o4FaultSlice(input, "LEAKED_FIELDS", []string{}, []string{"raw-content"})
		result.LayerFindings = o4FaultSlice(input, "LAYER_FINDINGS", input.LayerFindings, []string{"L1_EXTERNAL_EFFECT"})

	case OperationTools:
		result.BasalHashVerified = true
		result.Tools = o4FaultTools(input, input.BasalTools)

	case OperationSchema:
		result.Tools = o4FaultTools(input, []Tool{})

	case OperationArtifact:
		result.Lifecycle = input.Lifecycle
		result.ArtifactNames = o4FaultSlice(input, "ARTIFACT_NAMES", input.ArtifactNames, []string{})
		result.ArtifactSHA256 = canonicalMapDigest(input.ArtifactHashes)
		result.Tools = o4FaultTools(input, []Tool{})

	case OperationRefinement:
		result.RefinementsPreserved = o4FaultBool(input, "REFINEMENTS_PRESERVED", true, false)

	case OperationFull:
		result.RuntimeUnique, result.RefinementsPreserved, result.BasalHashVerified = true, true, true
		result.StoreIdentity, result.QueueIdentity = firstIdentity(input.StoreIDs), firstIdentity(input.QueueIDs)
		result.RootIdentity, result.SecretIdentity, result.WorkerIdentity = firstIdentity(input.RootIDs), firstIdentity(input.SecretIDs), firstIdentity(input.WorkerIDs)
		result.Tools = append([]Tool(nil), input.BasalTools...)
		result.ArtifactSHA256 = canonicalMapDigest(input.ArtifactHashes)
		result.LayerFindings = o4FaultSlice(input, "LAYER_FINDINGS", input.LayerFindings, appendO4Finding(input.LayerFindings, "L1_EXTERNAL_EFFECT"))

	default:
		return Observation{}, fmt.Errorf("unsupported O4 semantic port %s", testCase.Operation)
	}
	return result, nil
}

func o4FaultString(input Input, field, baseline, violated string) string {
	if input.SemanticFault == field {
		return violated
	}
	return baseline
}

func o4FaultBool(input Input, field string, baseline, violated bool) bool {
	if input.SemanticFault == field {
		return violated
	}
	return baseline
}

func o4FaultSlice(input Input, field string, baseline, violated []string) []string {
	if input.SemanticFault == field {
		return append([]string{}, violated...)
	}
	return append([]string{}, baseline...)
}

func o4FaultTools(input Input, baseline []Tool) []Tool {
	if input.SemanticFault == "TOOLS" {
		if len(baseline) == 0 {
			return []Tool{{Name: "uncommitted_tool", Definition: "invalid"}}
		}
		return append([]Tool(nil), baseline[:len(baseline)-1]...)
	}
	return append([]Tool(nil), baseline...)
}

func appendO4Finding(values []string, finding string) []string {
	result := append([]string(nil), values...)
	return append(result, finding)
}
