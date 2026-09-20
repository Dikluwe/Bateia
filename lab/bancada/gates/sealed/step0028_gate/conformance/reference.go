/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/reference.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sort"
)

type ReferenceSubject struct{}

func emptyObservation() Observation {
	return Observation{
		WarningCodes: []string{}, IntegrationUnknowns: []string{}, Tools: []Tool{},
		LeakedFields: []string{}, ExternalEffects: []string{},
	}
}

func (ReferenceSubject) Observe(_ context.Context, testCase Case) (Observation, error) {
	input := testCase.Input
	result := emptyObservation()
	if input.ScenarioVersion == "O3" || input.ScenarioVersion == "O4" || testCase.Operation == OperationArtifact || testCase.Operation == OperationPagination || testCase.Operation == OperationSchema || testCase.Operation == OperationFull {
		projectInputFields(&result, input)
	}
	switch testCase.Operation {
	case OperationExtract:
		result.Novelty = "UNKNOWN"
		if input.Owner != "" || input.Repo != "" || input.Commit != "" || input.Path != "" {
			result.CanonicalSHA256 = canonicalDigest([]string{input.Owner, input.Repo, input.Commit, input.Path, input.Language, input.Ecosystem})
		}
		if input.CatalogState == "INVALID" {
			result.ErrorCode = "INTEGRITY_ERROR"
			return result, nil
		}
		if input.RuleState == "UNKNOWN" || input.RuleFromContent {
			result.Decision = "CONTROLLED_GAP"
			result.MechanismStatus = "UNKNOWN"
			result.GapCode = "RULE_NOT_IN_MECHANISM_CATALOG"
			return result, nil
		}
		if input.SourceKindState == "INCOMPATIBLE" {
			result.Decision = "CONTROLLED_GAP"
			result.MechanismStatus = "UNKNOWN"
			result.GapCode = "INCOMPATIBLE_WITNESS_KIND"
			return result, nil
		}
		if input.ProductOnly {
			result.Decision = "OBSERVATION_ONLY"
			result.MechanismStatus = "NONE"
			return result, nil
		}
		if !input.SameImmutableSource || !input.ExactBundleBinding || !input.ContentHashValid {
			result.Decision = "CONTROLLED_GAP"
			result.MechanismStatus = "UNKNOWN"
			return result, nil
		}
		result.BindingVerified = true
		result.ProvenanceComplete = true
		if input.CentralConflict {
			result.Decision = "MECHANISM_PRESERVED"
			result.MechanismStatus = "INCONCLUSIVE"
			return result, nil
		}
		if input.DocumentationOnly {
			result.Decision = "MECHANISM_PRESERVED"
			result.MechanismStatus = "PARTIAL"
			return result, nil
		}
		if input.Implementation && input.Independent && input.ParentStatus == "PROVEN" {
			result.Decision = "MECHANISM_PRESERVED"
			result.MechanismStatus = "PROVEN"
			if input.ClassInBaseline {
				result.Novelty = "CONVENTIONAL"
			} else {
				result.Novelty = "DIVERGENT"
			}
			return result, nil
		}
		result.Decision = "MECHANISM_PRESERVED"
		result.MechanismStatus = "PARTIAL"
		return result, nil

	case OperationStart:
		if !input.SourceExists {
			result.ErrorCode = "RUN_NOT_FOUND"
			return result, nil
		}
		if !input.SourceTerminal || !input.SourcePublished {
			result.ErrorCode = "SOURCE_BASIS_UNAVAILABLE"
			return result, nil
		}
		if !input.SourceHashMatches || !input.CanonicalDocument {
			result.ErrorCode = "INTEGRITY_ERROR"
			return result, nil
		}
		result.ResolveCalls = 1
		if input.PriorRequestHash != "" {
			if input.PriorRequestHash != input.RequestHash {
				result.ErrorCode = "IDEMPOTENCY_CONFLICT"
				return result, nil
			}
			result.RunID = input.PriorRunID
			result.Lifecycle = "QUEUED"
			return result, nil
		}
		result.ReserveCalls, result.IDCalls, result.CreateCalls = 1, 1, 1
		if input.ForceCollision {
			result.ErrorCode = "ID_COLLISION"
			result.CollisionDetected = true
			result.CreateCalls = 0
			return result, nil
		}
		result.RunID, result.Lifecycle = "run_reference", "QUEUED"
		if input.CanonicalItems != nil {
			canonicalItems := append([]string(nil), input.CanonicalItems...)
			sort.Strings(canonicalItems)
			encoded, _ := json.Marshal(canonicalItems)
			digest := sha256.Sum256(encoded)
			result.CanonicalSHA256 = "sha256:" + hex.EncodeToString(digest[:])
		}
		if input.ConcurrentStarts > 0 {
			result.RunIDs = make([]string, input.ConcurrentStarts)
			for index := range result.RunIDs {
				result.RunIDs[index] = result.RunID
			}
		}
		return result, nil

	case OperationLifecycle:
		result.Lifecycle = input.Lifecycle
		result.AdoptionDisposition = "NOT_AN_ADOPTION_DECISION"
		if input.TargetState == "UNKNOWN" {
			result.HypothesisStatus = "PROPOSED"
		} else if input.MechanismStatus == "PROVEN" && input.CounterSearch == "NONE_FOUND_IN_SCOPE" {
			result.HypothesisStatus = "SUPPORTED"
		} else {
			result.HypothesisStatus = "INCONCLUSIVE"
		}
		if input.Cancelled {
			result.Lifecycle = "CANCELLED"
			result.CompletedSteps = 0
		}
		if input.RecoveryRequested {
			result.RecoveryState = "INTERRUPTED"
		}
		return result, nil

	case OperationReplay:
		if !input.CatalogHashMatches || !input.CaptureHashMatches {
			result.ErrorCode = "INTEGRITY_ERROR"
			return result, nil
		}
		items := append([]string(nil), input.ReplayItems...)
		sort.Strings(items)
		encoded, _ := json.Marshal(items)
		digest := sha256.Sum256(encoded)
		result.ReplaySHA256 = "sha256:" + hex.EncodeToString(digest[:])
		result.NetworkCalls = 0
		return result, nil

	case OperationRuntime:
		result.RuntimeUnique = singleIdentity(input.StoreIDs) && singleIdentity(input.QueueIDs) && singleIdentity(input.RootIDs) && singleIdentity(input.SecretIDs) && singleIdentity(input.WorkerIDs)
		if input.ScenarioVersion == "O3" || input.ScenarioVersion == "O4" {
			result.StoreIdentity = firstIdentity(input.StoreIDs)
			result.QueueIdentity = firstIdentity(input.QueueIDs)
			result.RootIdentity = firstIdentity(input.RootIDs)
			result.SecretIdentity = firstIdentity(input.SecretIDs)
			result.WorkerIdentity = firstIdentity(input.WorkerIDs)
		}
		return result, nil

	case OperationTools:
		result.BasalHashVerified = input.BasalToolsSHA256 != "" && input.BasalToolsSHA256 == input.ActualBasalToolsSHA256
		wantNames := []string{"bateia_search", "bateia_benchmark", "bateia_recommend", "bateia_start_recommendation", "bateia_get_run", "bateia_get_artifact", "bateia_cancel_run", "bateia_list_connectors", "bateia_start_market_discovery"}
		if len(input.BasalTools) != len(wantNames) {
			result.BasalHashVerified = false
		}
		for index := range input.BasalTools {
			if index >= len(wantNames) || input.BasalTools[index].Name != wantNames[index] {
				result.BasalHashVerified = false
			}
		}
		if !result.BasalHashVerified {
			result.ErrorCode = "INTEGRITY_ERROR"
			return result, nil
		}
		result.Tools = append([]Tool(nil), input.BasalTools...)
		if input.AppendTool10 {
			result.Tools = append(result.Tools, Tool{Name: "bateia_start_improvement_synthesis", Definition: "sha256:frozen-tool-10-schema"})
		}
		return result, nil

	case OperationRefinement:
		result.RefinementsPreserved = reflect.DeepEqual(input.RefinementHashes, input.RequiredRefinementHashes)
		return result, nil

	case OperationHypothesis:
		result.AdoptionDisposition = "NOT_AN_ADOPTION_DECISION"
		result.HypothesisStatus = "PROPOSED"
		if input.ScenarioVersion == "O3" || input.ScenarioVersion == "O4" {
			result.MechanismStatus = input.MechanismStatus
			result.Novelty = input.ClaimedNovelty
			if input.CanonicalItems != nil {
				result.CanonicalSHA256 = canonicalDigest(sortedCopy(input.CanonicalItems))
			}
		}
		if input.ProductOnly {
			result.Decision = "OBSERVATION_ONLY"
			result.MechanismStatus = "NONE"
		}
		if input.TargetState == "UNKNOWN" {
			result.HypothesisStatus = "PROPOSED"
		}
		if input.CounterSearch == "CENTRAL_CONFLICT" {
			result.MechanismStatus = "INCONCLUSIVE"
			result.HypothesisStatus = "INCONCLUSIVE"
		}
		if input.CounterSearch == "BUDGET_EXHAUSTED" {
			result.HypothesisStatus = "INCONCLUSIVE"
		}
		result.TransferMode = input.RequestedTransfer
		if input.License == "UNRECOGNIZED" {
			result.WarningCodes = []string{"LICENSE_NOT_RECOGNIZED"}
			if result.TransferMode == "COMPONENT_REUSE" {
				result.TransferMode = "PATTERN_ADAPTATION"
			}
		}
		if input.MechanismCount > 1 && !input.CompatibilityWitnessed {
			result.IntegrationUnknowns = []string{"failure model compatibility", "interface compatibility", "license compatibility", "performance compatibility", "state model compatibility"}
		}
		if input.ArtifactsPublished {
			result.ArtifactNames = append([]string(nil), input.ArtifactNames...)
			result.SourceRefs = append([]string(nil), input.SourceRefs...)
			result.CounterScope = append([]string(nil), input.CounterScope...)
			entries := append([]string(nil), input.ArtifactEntries...)
			if input.ArtifactsCanonical {
				sort.Strings(entries)
			}
			encoded, _ := json.Marshal(entries)
			digest := sha256.Sum256(encoded)
			result.ArtifactSHA256 = "sha256:" + hex.EncodeToString(digest[:])
		}
		return result, nil

	case OperationSafety:
		if len(input.RawFields) > 0 || len(input.ExternalEffects) > 0 {
			result.Decision = "SANITIZED_NO_EFFECT"
		}
		return result, nil

	case OperationArtifact:
		result.Lifecycle = input.Lifecycle
		if input.ArtifactState == "PUBLISHED" && (len(input.ArtifactNames) == 0 || len(input.ArtifactHashes) == 0) {
			result.ErrorCode = "INTEGRITY_ERROR"
			return result, nil
		}
		if input.ArtifactState == "PUBLISHED" && !reflect.DeepEqual(input.ArtifactTransitions, []string{"STAGED", "PUBLISHED"}) {
			result.ErrorCode = "INVALID_ARTIFACT_TRANSITION"
			return result, nil
		}
		if input.ArtifactState == "FAILED" && !reflect.DeepEqual(input.ArtifactTransitions, []string{"STAGED", "FAILED"}) {
			result.ErrorCode = "INVALID_ARTIFACT_TRANSITION"
			return result, nil
		}
		result.ArtifactSHA256 = canonicalMapDigest(input.ArtifactHashes)
		return result, nil

	case OperationPagination:
		if limit := input.AppliedLimits["page_size"]; limit > 0 && len(input.PageItems) > limit {
			result.ErrorCode = "LIMIT_EXCEEDED"
		}
		return result, nil

	case OperationSchema:
		if len(input.SchemaDefinitions) == 0 || len(input.ToolDefinitions) == 0 {
			result.ErrorCode = "INTEGRITY_ERROR"
		}
		return result, nil

	case OperationFull:
		result.RuntimeUnique = singleIdentity(input.StoreIDs) && singleIdentity(input.QueueIDs) && singleIdentity(input.RootIDs) && singleIdentity(input.SecretIDs) && singleIdentity(input.WorkerIDs)
		result.StoreIdentity = firstIdentity(input.StoreIDs)
		result.QueueIdentity = firstIdentity(input.QueueIDs)
		result.RootIdentity = firstIdentity(input.RootIDs)
		result.SecretIdentity = firstIdentity(input.SecretIDs)
		result.WorkerIdentity = firstIdentity(input.WorkerIDs)
		result.RefinementsPreserved = reflect.DeepEqual(input.RefinementReceipts, input.RequiredRefinementHashes)
		result.ArtifactSHA256 = canonicalMapDigest(input.ArtifactHashes)
		result.BasalHashVerified = input.BasalToolsSHA256 != "" && input.BasalToolsSHA256 == input.ActualBasalToolsSHA256
		result.Tools = append([]Tool(nil), input.BasalTools...)
		return result, nil
	}
	result.ErrorCode = "INVALID_OPERATION"
	return result, nil
}

func projectInputFields(result *Observation, input Input) {
	result.Owner, result.Repo, result.Commit, result.Path = input.Owner, input.Repo, input.Commit, input.Path
	result.Language, result.Ecosystem = input.Language, input.Ecosystem
	result.SourceHashes = append([]string(nil), input.SourceHashes...)
	result.SourceKinds = append([]string(nil), input.SourceKinds...)
	result.RuleIDs = append([]string(nil), input.RuleIDs...)
	result.WitnessRoles = append([]string(nil), input.WitnessRoles...)
	result.AppliedLimits = cloneMap(input.AppliedLimits)
	result.MechanismIDs = append([]string(nil), input.MechanismIDs...)
	result.HypothesisIDs = append([]string(nil), input.HypothesisIDs...)
	result.ArtifactIDs = append([]string(nil), input.ArtifactIDs...)
	result.ValidationPlan = append([]string(nil), input.ValidationPlan...)
	result.TargetState = input.TargetState
	result.TargetGapRefs = append([]string(nil), input.TargetGapRefs...)
	result.MechanismRefs = append([]string(nil), input.MechanismRefs...)
	result.CounterEvidenceRefs = append([]string(nil), input.CounterEvidenceRefs...)
	result.Assumptions = append([]string(nil), input.Assumptions...)
	result.Preconditions = append([]string(nil), input.Preconditions...)
	result.Risks = append([]string(nil), input.Risks...)
	result.Tradeoffs = append([]string(nil), input.Tradeoffs...)
	result.Ordering = append([]string(nil), input.Ordering...)
	result.License, result.Fit = input.License, input.Fit
	result.ImprovementStatus = input.ImprovementStatus
	result.ArtifactState = input.ArtifactState
	result.ArtifactTransitions = append([]string(nil), input.ArtifactTransitions...)
	result.ArtifactNames = append([]string(nil), input.ArtifactNames...)
	result.ArtifactHashes = cloneMap(input.ArtifactHashes)
	result.IndexHashes = cloneMap(input.IndexHashes)
	result.SourceRefs = append([]string(nil), input.SourceRefs...)
	result.CounterScope = append([]string(nil), input.CounterScope...)
	result.SnapshotSHA256, result.Cursor, result.NextCursor = input.SnapshotSHA256, input.Cursor, input.NextCursor
	result.PageItems = append([]string(nil), input.PageItems...)
	result.ToolDefinitions = cloneMap(input.ToolDefinitions)
	result.SchemaDefinitions = cloneMap(input.SchemaDefinitions)
	result.RefinementReceipts = cloneMap(input.RefinementReceipts)
	result.LayerFindings = append([]string(nil), input.LayerFindings...)
}

func cloneMap[T any](source map[string]T) map[string]T {
	if source == nil {
		return nil
	}
	result := make(map[string]T, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

func canonicalDigest(items []string) string {
	encoded, _ := json.Marshal(items)
	digest := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func canonicalMapDigest(values map[string]string) string {
	if values == nil {
		return ""
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	items := make([]string, 0, len(keys))
	for _, key := range keys {
		items = append(items, key+"="+values[key])
	}
	return canonicalDigest(items)
}

func firstIdentity(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func singleIdentity(values []string) bool {
	if len(values) == 0 {
		return false
	}
	first := values[0]
	if first == "" {
		return false
	}
	for _, value := range values[1:] {
		if value != first {
			return false
		}
	}
	return true
}
