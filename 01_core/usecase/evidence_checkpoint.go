/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_checkpoint.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import "bateia/01_core/domain"

type EvidenceCheckpointCriteria struct {
	EvidenceHash       string
	EvidenceRulesHash  string
	RuleVersions       map[string]string
	ScreeningModelHash string
}

// EvidenceCheckpointReader reads hashes of logical artifacts selected by the policy.
type EvidenceCheckpointReader interface{ ArtifactSHA256(string) (string, error) }

type EvidenceCheckpointReason string

const (
	CheckpointReusable            EvidenceCheckpointReason = "CHECKPOINT_REUSABLE"
	CheckpointManifestUnavailable EvidenceCheckpointReason = "MANIFEST_UNAVAILABLE"
	CheckpointStatusIneligible    EvidenceCheckpointReason = "RUN_STATUS_NOT_ELIGIBLE"
	CheckpointSchemaUnsupported   EvidenceCheckpointReason = "SCHEMA_UNSUPPORTED"
	CheckpointEvidenceChanged     EvidenceCheckpointReason = "EVIDENCE_INPUT_CHANGED"
	CheckpointRulesChanged        EvidenceCheckpointReason = "EVIDENCE_RULES_CHANGED"
	CheckpointModelChanged        EvidenceCheckpointReason = "SCREENING_MODEL_CHANGED"
	CheckpointVersionChanged      EvidenceCheckpointReason = "RULE_VERSION_CHANGED"
	CheckpointReaderUnavailable   EvidenceCheckpointReason = "ARTIFACT_READER_UNAVAILABLE"
	CheckpointArtifactUnreadable  EvidenceCheckpointReason = "ARTIFACT_UNREADABLE"
	CheckpointArtifactMismatch    EvidenceCheckpointReason = "ARTIFACT_HASH_MISMATCH"
	CheckpointArtifactUnrecorded  EvidenceCheckpointReason = "ARTIFACT_NOT_RECORDED"
)

// EvidenceCheckpointAssessment describes eligibility, not the strength of evidence.
// Subject contains only a known logical artifact or field; never reader errors or input values.
type EvidenceCheckpointAssessment struct {
	Reusable bool
	Reason   EvidenceCheckpointReason
	Subject  string
}

func CanReuseEvidenceCheckpoint(previous domain.IdeaRunManifest, criteria EvidenceCheckpointCriteria, reader EvidenceCheckpointReader) bool {
	return EvaluateEvidenceCheckpoint(previous, criteria, reader).Reusable
}

// EvaluateEvidenceCheckpoint reports the first impediment without extra artifact reads.
func EvaluateEvidenceCheckpoint(previous domain.IdeaRunManifest, criteria EvidenceCheckpointCriteria, reader EvidenceCheckpointReader) EvidenceCheckpointAssessment {
	switch previous.Status {
	case domain.IdeaRunComplete, domain.IdeaRunPartial, domain.IdeaRunInconclusive:
	default:
		return EvidenceCheckpointAssessment{Reason: CheckpointStatusIneligible, Subject: "status"}
	}
	if previous.SchemaVersion != "1" {
		return EvidenceCheckpointAssessment{Reason: CheckpointSchemaUnsupported, Subject: "schema_version"}
	}
	if previous.InputHashes["evidence"] != criteria.EvidenceHash {
		return EvidenceCheckpointAssessment{Reason: CheckpointEvidenceChanged, Subject: "evidence"}
	}
	if previous.InputHashes["evidence_rules"] != criteria.EvidenceRulesHash {
		return EvidenceCheckpointAssessment{Reason: CheckpointRulesChanged, Subject: "evidence_rules"}
	}
	screeningModelHash := criteria.ScreeningModelHash
	previousScreeningModelHash, previousHasScreeningModel := previous.InputHashes["screening_model"]
	if screeningModelHash == "" {
		if previousHasScreeningModel {
			return EvidenceCheckpointAssessment{Reason: CheckpointModelChanged, Subject: "screening_model"}
		}
	} else if !previousHasScreeningModel || previousScreeningModelHash != screeningModelHash {
		return EvidenceCheckpointAssessment{Reason: CheckpointModelChanged, Subject: "screening_model"}
	}
	for _, key := range []string{"expansion", "screening", "evidence", "analyzer", "implementation_search", "discovery_query", "local_snapshot", "local_collection"} {
		if previous.RuleVersions[key] != criteria.RuleVersions[key] {
			return EvidenceCheckpointAssessment{Reason: CheckpointVersionChanged, Subject: key}
		}
	}
	if reader == nil {
		return EvidenceCheckpointAssessment{Reason: CheckpointReaderUnavailable}
	}
	required := map[string]bool{"discovery-snapshot.json": false, "evidence/bundles.json": false}
	for _, stage := range previous.Stages {
		for _, artifact := range stage.Outputs {
			expected, ok := required[artifact.Path]
			if !ok || expected || artifact.Hash == "" {
				continue
			}
			actual, err := reader.ArtifactSHA256(artifact.Path)
			if err != nil {
				return EvidenceCheckpointAssessment{Reason: CheckpointArtifactUnreadable, Subject: artifact.Path}
			}
			if actual != artifact.Hash {
				return EvidenceCheckpointAssessment{Reason: CheckpointArtifactMismatch, Subject: artifact.Path}
			}
			required[artifact.Path] = true
		}
	}
	for _, path := range []string{"discovery-snapshot.json", "evidence/bundles.json"} {
		if !required[path] {
			return EvidenceCheckpointAssessment{Reason: CheckpointArtifactUnrecorded, Subject: path}
		}
	}
	return EvidenceCheckpointAssessment{Reusable: true, Reason: CheckpointReusable}
}
