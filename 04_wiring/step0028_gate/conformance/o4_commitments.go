/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o4_commitments.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

type O4Predicate struct {
	ExpectedOutcome         string `json:"expected_outcome"`
	SubjectKind             string `json:"subject_kind"`
	RequiredRelation        string `json:"required_relation"`
	SemanticVariantID       string `json:"semantic_variant_id,omitempty"`
	ChangedFieldCardinality int    `json:"changed_field_cardinality,omitempty"`
	Allowed                 *bool  `json:"allowed,omitempty"`
	Executable              *bool  `json:"executable,omitempty"`
}

type O4Commitment struct {
	ID                         string      `json:"id"`
	CanonicalAttackEntrySHA256 string      `json:"canonical_attack_entry_sha256"`
	PublicContractBinding      string      `json:"public_contract_binding"`
	RequiredPort               Operation   `json:"required_port"`
	InputFixtureFamily         string      `json:"input_fixture_family"`
	SemanticVariantID          string      `json:"semantic_variant_id"`
	BaselinePredicate          O4Predicate `json:"baseline_predicate"`
	ViolatedPredicate          O4Predicate `json:"violated_predicate"`
	RequiredObservableFields   []string    `json:"required_observable_fields"`
}

type O4CommitmentDocument struct {
	SchemaVersion        string          `json:"schema_version"`
	Visibility           string          `json:"visibility"`
	Step                 string          `json:"step"`
	SourceOracleRevision string          `json:"source_oracle_revision"`
	TargetOracleRevision string          `json:"target_oracle_revision"`
	ContractRevision     string          `json:"contract_revision"`
	Attestation          string          `json:"attestation"`
	PrivacyContract      json.RawMessage `json:"privacy_contract"`
	Canonicalization     json.RawMessage `json:"canonicalization"`
	GlobalAttackBinding  struct {
		SourcePath                   string          `json:"source_path"`
		SourceRawSHA256              string          `json:"source_raw_sha256"`
		SourceEntryCount             int             `json:"source_entry_count"`
		SourceCaseOrder              string          `json:"source_case_order"`
		CanonicalEnvelopeSHA256      string          `json:"canonical_envelope_without_cases_sha256"`
		CanonicalDocumentSHA256      string          `json:"canonical_document_sha256"`
		EntryCommitmentRootAlgorithm string          `json:"entry_commitment_root_algorithm"`
		EntryCommitmentRootSHA256    string          `json:"entry_commitment_root_sha256"`
		ReconstructionValidation     json.RawMessage `json:"reconstruction_validation"`
	} `json:"global_attack_binding"`
	Taxonomy json.RawMessage `json:"taxonomy"`
	Counts   struct {
		Negative int `json:"negative"`
		Positive int `json:"positive"`
		Opaque   int `json:"opaque"`
		Total    int `json:"total"`
	} `json:"counts"`
	Entries []O4Commitment `json:"entries"`
}

func DecodeO4Commitments(document []byte) (O4CommitmentDocument, error) {
	var envelope O4CommitmentDocument
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return envelope, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return envelope, fmt.Errorf("trailing O4 commitment document")
	}
	if envelope.SchemaVersion != "tekt.step0028.mutation-commitments-public.v4" || envelope.Visibility != "PUBLIC_ALLOWLISTED_ORACLE_FEEDBACK" || envelope.Step != "0028" || envelope.SourceOracleRevision != "O3" || envelope.TargetOracleRevision != "O4" || envelope.ContractRevision != "R3" || envelope.Counts.Negative != 160 || envelope.Counts.Positive != 12 || envelope.Counts.Opaque != 10 || envelope.Counts.Total != 182 || len(envelope.Entries) != 182 {
		return envelope, fmt.Errorf("invalid O4 commitment envelope")
	}
	if envelope.GlobalAttackBinding.SourceEntryCount != 182 || envelope.GlobalAttackBinding.SourceCaseOrder != "M001-M160,P01-P12,O01-O10" || envelope.GlobalAttackBinding.EntryCommitmentRootSHA256 != "a486a67b330946daea5d9bfb41f7764ce35a1e300f94bb77210341d3893a4a3a" {
		return envelope, fmt.Errorf("invalid O4 global attack binding")
	}
	seenIDs, seenVariants := map[string]bool{}, map[string]bool{}
	for index, entry := range envelope.Entries {
		wantID := o4IDAt(index)
		if entry.ID != wantID || seenIDs[entry.ID] || seenVariants[entry.SemanticVariantID] || len(entry.CanonicalAttackEntrySHA256) != 64 || entry.RequiredPort == "" || entry.InputFixtureFamily == "" || entry.PublicContractBinding == "" || len(entry.RequiredObservableFields) == 0 {
			return envelope, fmt.Errorf("invalid O4 commitment at %d/%s", index, entry.ID)
		}
		if _, err := hex.DecodeString(entry.CanonicalAttackEntrySHA256); err != nil {
			return envelope, fmt.Errorf("invalid entry hash for %s", entry.ID)
		}
		seenIDs[entry.ID], seenVariants[entry.SemanticVariantID] = true, true
	}
	root, err := O4EntryCommitmentRoot(envelope.Entries)
	if err != nil || root != envelope.GlobalAttackBinding.EntryCommitmentRootSHA256 {
		return envelope, fmt.Errorf("entry commitment root mismatch: %s", root)
	}
	return envelope, nil
}

func O4EntryCommitmentRoot(entries []O4Commitment) (string, error) {
	digest := sha256.New()
	for _, entry := range entries {
		raw, err := hex.DecodeString(entry.CanonicalAttackEntrySHA256)
		if err != nil || len(raw) != sha256.Size {
			return "", fmt.Errorf("invalid canonical entry hash for %s", entry.ID)
		}
		_, _ = digest.Write([]byte(entry.ID))
		_, _ = digest.Write([]byte{0})
		_, _ = digest.Write(raw)
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func o4IDAt(index int) string {
	switch {
	case index < 160:
		return fmt.Sprintf("M%03d", index+1)
	case index < 172:
		return fmt.Sprintf("P%02d", index-159)
	default:
		return fmt.Sprintf("O%02d", index-171)
	}
}
