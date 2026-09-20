//go:build step0028_o3_invalidated

/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o3_variants.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type VariantSpec struct {
	ID              string `json:"id"`
	ControlID       string `json:"control_id"`
	ObservableField string `json:"observable_field"`
	VariantKind     string `json:"variant_kind"`
}

type VariantSubject struct {
	Base    Subject
	Variant VariantSpec
}

func DecodeO3Variants(document []byte) ([]VariantSpec, error) {
	var envelope struct {
		SchemaVersion           string        `json:"schema_version"`
		Step                    string        `json:"step"`
		OracleRevision          string        `json:"oracle_revision"`
		ContractRevision        string        `json:"contract_revision"`
		SourceCoverageGapSHA256 string        `json:"source_coverage_gap_sha256"`
		SelectionPolicy         string        `json:"selection_policy"`
		Entries                 []VariantSpec `json:"entries"`
	}
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return nil, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return nil, fmt.Errorf("trailing O3 variant document")
	}
	if envelope.SchemaVersion != "tekt.step0028.oracle-variants.v3" || envelope.Step != "0028" || envelope.OracleRevision != "O3" || envelope.ContractRevision != "R3" || envelope.SourceCoverageGapSHA256 != "92da449660fdd1f0fba1906cfe9077168820d1e838c671633c33fd0c83c69f3b" || envelope.SelectionPolicy != "ROUND_ROBIN_BY_PUBLIC_DIAGNOSTIC_GROUP" || len(envelope.Entries) != 160 {
		return nil, fmt.Errorf("invalid O3 variant envelope")
	}
	seen := map[string]bool{}
	for _, variant := range envelope.Entries {
		if variant.ID == "" || variant.ControlID == "" || variant.ObservableField == "" || variant.VariantKind != "ALTER_OBSERVABLE" || seen[variant.ID] {
			return nil, fmt.Errorf("invalid or duplicate O3 variant %q", variant.ID)
		}
		seen[variant.ID] = true
	}
	return envelope.Entries, nil
}

func (subject VariantSubject) Observe(ctx context.Context, testCase Case) (Observation, error) {
	if subject.Base == nil {
		return Observation{}, fmt.Errorf("variant %s has no base subject", subject.Variant.ID)
	}
	result, err := subject.Base.Observe(ctx, testCase)
	if err != nil {
		return result, err
	}
	if err := alterObservable(&result, subject.Variant.ObservableField, subject.Variant.ID); err != nil {
		return Observation{}, err
	}
	return result, nil
}

func alterObservable(result *Observation, field, variantID string) error {
	altered := "ALTERED_BY_" + variantID
	switch field {
	case "decision":
		result.Decision = altered
	case "mechanism_status":
		result.MechanismStatus = altered
	case "hypothesis_status":
		result.HypothesisStatus = altered
	case "improvement_status":
		result.ImprovementStatus = altered
	case "novelty":
		result.Novelty = altered
	case "gap_code":
		result.GapCode = altered
	case "error_code":
		result.ErrorCode = altered
	case "canonical_sha256":
		result.CanonicalSHA256 = "sha256:" + variantID
	case "artifact_sha256":
		result.ArtifactSHA256 = "sha256:" + variantID
	case "replay_sha256":
		result.ReplaySHA256 = "sha256:" + variantID
	case "transfer_mode":
		result.TransferMode = altered
	case "adoption_disposition":
		result.AdoptionDisposition = altered
	case "run_id":
		result.RunID = altered
	case "lifecycle":
		result.Lifecycle = altered
	case "recovery_state":
		result.RecoveryState = altered
	case "store_identity":
		result.StoreIdentity = altered
	case "queue_identity":
		result.QueueIdentity = altered
	case "root_identity":
		result.RootIdentity = altered
	case "secret_identity":
		result.SecretIdentity = altered
	case "worker_identity":
		result.WorkerIdentity = altered
	case "target_state":
		result.TargetState = altered
	case "license":
		result.License = altered
	case "fit":
		result.Fit = altered
	case "owner":
		result.Owner = altered
	case "repo":
		result.Repo = altered
	case "commit":
		result.Commit = altered
	case "path":
		result.Path = altered
	case "language":
		result.Language = altered
	case "ecosystem":
		result.Ecosystem = altered
	case "binding_verified":
		result.BindingVerified = !result.BindingVerified
	case "provenance_complete":
		result.ProvenanceComplete = !result.ProvenanceComplete
	case "runtime_unique":
		result.RuntimeUnique = !result.RuntimeUnique
	case "basal_hash_verified":
		result.BasalHashVerified = !result.BasalHashVerified
	case "refinements_preserved":
		result.RefinementsPreserved = !result.RefinementsPreserved
	case "collision_detected":
		result.CollisionDetected = !result.CollisionDetected
	case "resolve_calls":
		result.ResolveCalls++
	case "reserve_calls":
		result.ReserveCalls++
	case "id_calls":
		result.IDCalls++
	case "create_calls":
		result.CreateCalls++
	case "completed_steps":
		result.CompletedSteps++
	case "network_calls":
		result.NetworkCalls++
	case "warning_codes":
		result.WarningCodes = appendCopy(result.WarningCodes, altered)
	case "integration_unknowns":
		result.IntegrationUnknowns = appendCopy(result.IntegrationUnknowns, altered)
	case "run_ids":
		result.RunIDs = appendCopy(result.RunIDs, altered)
	case "artifact_names":
		result.ArtifactNames = appendCopy(result.ArtifactNames, altered)
	case "source_refs":
		result.SourceRefs = appendCopy(result.SourceRefs, altered)
	case "counter_scope":
		result.CounterScope = appendCopy(result.CounterScope, altered)
	case "leaked_fields":
		result.LeakedFields = appendCopy(result.LeakedFields, altered)
	case "external_effects":
		result.ExternalEffects = appendCopy(result.ExternalEffects, altered)
	case "layer_findings":
		result.LayerFindings = appendCopy(result.LayerFindings, altered)
	case "validation_plan":
		result.ValidationPlan = appendCopy(result.ValidationPlan, altered)
	case "target_gap_refs":
		result.TargetGapRefs = appendCopy(result.TargetGapRefs, altered)
	case "mechanism_refs":
		result.MechanismRefs = appendCopy(result.MechanismRefs, altered)
	case "counterevidence_refs":
		result.CounterEvidenceRefs = appendCopy(result.CounterEvidenceRefs, altered)
	case "ordering":
		result.Ordering = appendCopy(result.Ordering, altered)
	case "artifact_hashes":
		result.ArtifactHashes = alteredMap(result.ArtifactHashes, variantID)
	case "index_hashes":
		result.IndexHashes = alteredMap(result.IndexHashes, variantID)
	case "tool_definitions":
		result.ToolDefinitions = alteredMap(result.ToolDefinitions, variantID)
	case "refinement_receipts":
		result.RefinementReceipts = alteredMap(result.RefinementReceipts, variantID)
	case "tools":
		result.Tools = append(append([]Tool(nil), result.Tools...), Tool{Name: altered, Definition: "sha256:" + variantID})
	default:
		return fmt.Errorf("variant %s names unsupported observable %q", variantID, field)
	}
	return nil
}

func appendCopy(values []string, value string) []string {
	return append(append([]string(nil), values...), value)
}

func alteredMap(values map[string]string, variantID string) map[string]string {
	result := cloneMap(values)
	if result == nil {
		result = map[string]string{}
	}
	result["variant:"+variantID] = "sha256:" + variantID
	return result
}
