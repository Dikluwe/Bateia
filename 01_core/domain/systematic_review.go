/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/systematic_review.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

const SystematicReviewSchemaVersion = "1"

type ReviewProtocolDraft struct {
	Intake                  RecommendationIntake `json:"intake"`
	ScreeningModel          ScreeningModel       `json:"screening_model"`
	ScreeningModelReference string               `json:"screening_model_reference,omitempty"`
	EvidenceRulesVersion    string               `json:"evidence_rules_version"`
	EvidenceRulesHash       string               `json:"evidence_rules_hash"`
	UpdatePolicy            LivingUpdatePolicy   `json:"update_policy"`
}

type EligibilityCriterion struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Rule      string `json:"rule"`
	SourceRef string `json:"source_ref"`
}

type ReviewSearchProtocol struct {
	Strategy      string   `json:"strategy"`
	Language      string   `json:"language,omitempty"`
	Locales       []string `json:"locales"`
	MaxPerQuery   int      `json:"max_per_query"`
	MaxCandidates int      `json:"max_candidates"`
	Provider      string   `json:"provider"`
}

type ReviewScreeningProtocol struct {
	ModelReference   string  `json:"model_reference,omitempty"`
	ModelHash        string  `json:"model_hash"`
	AdvanceThreshold float64 `json:"advance_threshold"`
	DropThreshold    float64 `json:"drop_threshold"`
	AuditSampleRate  float64 `json:"audit_sample_rate"`
	Calibration      string  `json:"calibration_status"`
}

type ReviewEvidenceProtocol struct {
	RulesVersion string `json:"rules_version"`
	RulesHash    string `json:"rules_hash"`
}

type LivingUpdatePolicy struct {
	Mode                 string `json:"mode"`
	Trigger              string `json:"trigger"`
	RequireSameProtocol  bool   `json:"require_same_protocol"`
	PreservePriorVersion bool   `json:"preserve_prior_version"`
}

type ReviewProtocol struct {
	SchemaVersion string                         `json:"schema_version"`
	ProtocolHash  string                         `json:"protocol_hash"`
	IntakeHash    string                         `json:"intake_hash"`
	Question      IdeaSpec                       `json:"question"`
	Target        SoftwareAlternativeTarget      `json:"target"`
	Search        ReviewSearchProtocol           `json:"search"`
	Eligibility   []EligibilityCriterion         `json:"eligibility"`
	Screening     ReviewScreeningProtocol        `json:"screening"`
	Evidence      ReviewEvidenceProtocol         `json:"evidence"`
	Decision      ReuseDecisionPolicy            `json:"decision"`
	Execution     RecommendationExecutionProfile `json:"execution"`
	UpdatePolicy  LivingUpdatePolicy             `json:"update_policy"`
}

type SourceRecordDisposition string

const (
	SourceRecordCanonical SourceRecordDisposition = "CANONICAL"
	SourceRecordDuplicate SourceRecordDisposition = "DUPLICATE"
)

type CandidateTerminalDecision string

const (
	CandidateIncluded     CandidateTerminalDecision = "INCLUDED"
	CandidateExcluded     CandidateTerminalDecision = "EXCLUDED"
	CandidateInconclusive CandidateTerminalDecision = "INCONCLUSIVE"
)

type CandidateDispositionReason string

const (
	ReasonIncludedRecommended  CandidateDispositionReason = "INCLUDED_RECOMMENDED"
	ReasonIncludedPartial      CandidateDispositionReason = "INCLUDED_PARTIAL"
	ReasonIncludedReference    CandidateDispositionReason = "INCLUDED_REFERENCE"
	ReasonTargetRepository     CandidateDispositionReason = "TARGET_REPOSITORY"
	ReasonScreeningDrop        CandidateDispositionReason = "SCREENING_DROP"
	ReasonEvidenceNotAssessed  CandidateDispositionReason = "EVIDENCE_NOT_ASSESSED"
	ReasonEvidenceRetrieval    CandidateDispositionReason = "EVIDENCE_RETRIEVAL_FAILED"
	ReasonLicenseBlocked       CandidateDispositionReason = "LICENSE_BLOCKED"
	ReasonArchivedBlocked      CandidateDispositionReason = "ARCHIVED_BLOCKED"
	ReasonLocalRequired        CandidateDispositionReason = "LOCAL_OPERATION_REQUIRED"
	ReasonEvidenceConflict     CandidateDispositionReason = "EVIDENCE_CONFLICT"
	ReasonNotAlternative       CandidateDispositionReason = "NOT_AN_ALTERNATIVE"
	ReasonPolicyBlocked        CandidateDispositionReason = "POLICY_BLOCKED"
	ReasonEvidenceInconclusive CandidateDispositionReason = "EVIDENCE_INCONCLUSIVE"
)

type ReviewSourceRecord struct {
	RecordID          string                  `json:"record_id"`
	CandidateID       string                  `json:"candidate_id"`
	Repository        RepositoryCandidate     `json:"repository"`
	Provider          string                  `json:"provider"`
	ExpansionID       string                  `json:"expansion_id"`
	PlannedQuery      string                  `json:"planned_query"`
	Disposition       SourceRecordDisposition `json:"disposition"`
	CanonicalRecordID string                  `json:"canonical_record_id,omitempty"`
}

type CandidateDispositionEvent struct {
	Stage      string                     `json:"stage"`
	Outcome    string                     `json:"outcome"`
	ReasonCode CandidateDispositionReason `json:"reason_code,omitempty"`
	SourceRef  string                     `json:"source_ref,omitempty"`
}

type CandidateDisposition struct {
	CandidateID     string                      `json:"candidate_id"`
	Repository      CommitRef                   `json:"repository"`
	SourceRecordIDs []string                    `json:"source_record_ids"`
	Events          []CandidateDispositionEvent `json:"events"`
	Terminal        CandidateTerminalDecision   `json:"terminal_decision"`
	ReasonCode      CandidateDispositionReason  `json:"reason_code"`
}

type CandidateDispositionLedger struct {
	SchemaVersion string                 `json:"schema_version"`
	ProtocolHash  string                 `json:"protocol_hash"`
	SourceRecords []ReviewSourceRecord   `json:"source_records"`
	Candidates    []CandidateDisposition `json:"candidates"`
}

type ReasonCount struct {
	ReasonCode CandidateDispositionReason `json:"reason_code"`
	Count      int                        `json:"count"`
}

type FunnelFlowReport struct {
	SchemaVersion       string        `json:"schema_version"`
	ProtocolHash        string        `json:"protocol_hash"`
	IdentifiedRecords   int           `json:"identified_records"`
	ProviderRejected    int           `json:"provider_rejected_observations"`
	DuplicateRecords    int           `json:"duplicate_records"`
	UniqueCandidates    int           `json:"unique_candidates"`
	Screened            int           `json:"screened"`
	Unscreened          int           `json:"unscreened"`
	ScreeningAdvance    int           `json:"screening_advance"`
	ScreeningAudit      int           `json:"screening_audit"`
	ScreeningDrop       int           `json:"screening_drop"`
	Assessed            int           `json:"assessed"`
	EvidenceFailures    int           `json:"evidence_failures"`
	Included            int           `json:"included"`
	Excluded            int           `json:"excluded"`
	Inconclusive        int           `json:"inconclusive"`
	TerminalReasons     []ReasonCount `json:"terminal_reasons"`
	ConservationChecked bool          `json:"conservation_checked"`
}

type CertaintyLevel string

const (
	CertaintyRobust    CertaintyLevel = "ROBUST"
	CertaintySupported CertaintyLevel = "SUPPORTED"
	CertaintyFragile   CertaintyLevel = "FRAGILE"
	CertaintyUnknown   CertaintyLevel = "UNKNOWN"
)

type CertaintyConcern string

const (
	ConcernNone    CertaintyConcern = "NONE"
	ConcernSome    CertaintyConcern = "SOME"
	ConcernSerious CertaintyConcern = "SERIOUS"
	ConcernUnknown CertaintyConcern = "UNKNOWN"
)

type ClaimConcern struct {
	Dimension string           `json:"dimension"`
	Concern   CertaintyConcern `json:"concern"`
	Reason    string           `json:"reason"`
}

type ClaimCertainty struct {
	CandidateID    string         `json:"candidate_id"`
	ClaimID        string         `json:"claim_id"`
	ObservedStatus EvidenceStatus `json:"observed_status"`
	Level          CertaintyLevel `json:"certainty"`
	Concerns       []ClaimConcern `json:"concerns"`
	SourceCount    int            `json:"source_count"`
}

type BiasRisk string

const (
	BiasLow     BiasRisk = "LOW"
	BiasPresent BiasRisk = "PRESENT"
	BiasUnknown BiasRisk = "UNKNOWN"
)

type BiasAssessment struct {
	Bias       string   `json:"bias"`
	Risk       BiasRisk `json:"risk"`
	Reason     string   `json:"reason"`
	Mitigation string   `json:"mitigation"`
}

type ReviewCertaintyReport struct {
	SchemaVersion string           `json:"schema_version"`
	Claims        []ClaimCertainty `json:"claims"`
	Biases        []BiasAssessment `json:"biases"`
}

type SystematicReviewPackage struct {
	SchemaVersion   string                     `json:"schema_version"`
	PackageHash     string                     `json:"package_hash"`
	Protocol        ReviewProtocol             `json:"protocol"`
	SearchedThrough string                     `json:"searched_through"`
	Ledger          CandidateDispositionLedger `json:"ledger"`
	Flow            FunnelFlowReport           `json:"flow"`
	Certainty       ReviewCertaintyReport      `json:"certainty"`
	Decision        ReuseDecisionReport        `json:"decision"`
	Alternatives    SoftwareAlternativeReport  `json:"alternatives"`
}

type ReviewSynthesisInput struct {
	Protocol        ReviewProtocol            `json:"protocol"`
	Discovery       DiscoverySnapshot         `json:"discovery"`
	Screening       ScreeningReport           `json:"screening"`
	Assessment      AssessmentReport          `json:"assessment"`
	EvidenceBundles []EvidenceBundle          `json:"evidence_bundles"`
	Decision        ReuseDecisionReport       `json:"decision"`
	Alternatives    SoftwareAlternativeReport `json:"alternatives"`
	SearchedThrough string                    `json:"searched_through"`
}

type LivingDispositionChange struct {
	CandidateID  string                     `json:"candidate_id"`
	Before       CandidateTerminalDecision  `json:"before"`
	After        CandidateTerminalDecision  `json:"after"`
	BeforeReason CandidateDispositionReason `json:"before_reason"`
	AfterReason  CandidateDispositionReason `json:"after_reason"`
}

type LivingCertaintyChange struct {
	CandidateID string         `json:"candidate_id"`
	ClaimID     string         `json:"claim_id"`
	Before      CertaintyLevel `json:"before"`
	After       CertaintyLevel `json:"after"`
}

type LivingReviewDelta struct {
	FromPackageHash    string                    `json:"from_package_hash,omitempty"`
	ToPackageHash      string                    `json:"to_package_hash"`
	AddedCandidates    []string                  `json:"added_candidates"`
	RemovedCandidates  []string                  `json:"removed_candidates"`
	DispositionChanges []LivingDispositionChange `json:"disposition_changes"`
	CertaintyChanges   []LivingCertaintyChange   `json:"certainty_changes"`
	DecisionChanged    bool                      `json:"decision_changed"`
}

type LivingReviewVersion struct {
	Version         int               `json:"version"`
	PackageHash     string            `json:"package_hash"`
	SearchedThrough string            `json:"searched_through"`
	Delta           LivingReviewDelta `json:"delta"`
}

type LivingReviewLineage struct {
	LineageID    string                `json:"lineage_id"`
	ProtocolHash string                `json:"protocol_hash"`
	Versions     []LivingReviewVersion `json:"versions"`
}

type LivingReviewIndex struct {
	SchemaVersion string                `json:"schema_version"`
	Lineages      []LivingReviewLineage `json:"lineages"`
}

type LivingReviewUpdateInput struct {
	Index           LivingReviewIndex        `json:"index"`
	PreviousPackage *SystematicReviewPackage `json:"previous_package,omitempty"`
	CurrentPackage  SystematicReviewPackage  `json:"current_package"`
}

type ActiveLearningRequest struct {
	Screening ScreeningReport     `json:"screening"`
	Reviewed  []CandidateJudgment `json:"reviewed"`
	BatchSize int                 `json:"batch_size"`
}

type ActiveLearningPriority struct {
	Repository         string  `json:"repository"`
	Priority           float64 `json:"priority"`
	EstimatedRelevance float64 `json:"estimated_relevance"`
	Uncertainty        float64 `json:"uncertainty"`
	Novelty            float64 `json:"novelty"`
	Reason             string  `json:"reason"`
}

type ActiveLearningPlan struct {
	SchemaVersion string                   `json:"schema_version"`
	Authority     string                   `json:"authority"`
	CanExclude    bool                     `json:"can_exclude"`
	Priorities    []ActiveLearningPriority `json:"priorities"`
}

type ActiveLearningCase struct {
	ID        string              `json:"id"`
	Screening ScreeningReport     `json:"screening"`
	Seeds     []CandidateJudgment `json:"seeds"`
	Gold      []CandidateJudgment `json:"gold"`
}

type ActiveLearningCorpus struct {
	SchemaVersion string               `json:"schema_version"`
	Cases         []ActiveLearningCase `json:"cases"`
}

type ActiveLearningCaseMetrics struct {
	CaseID              string  `json:"case_id"`
	Candidates          int     `json:"candidates"`
	Relevant            int     `json:"relevant"`
	ScreenedFor95Recall int     `json:"screened_for_95_recall"`
	WorkSavedAt95       float64 `json:"work_saved_at_95_recall"`
	RecallAtHalf        float64 `json:"recall_at_half"`
}

type ActiveLearningEvaluation struct {
	SchemaVersion      string                      `json:"schema_version"`
	Cases              []ActiveLearningCaseMetrics `json:"cases"`
	AverageWorkSaved95 float64                     `json:"average_work_saved_at_95_recall"`
	AverageRecallHalf  float64                     `json:"average_recall_at_half"`
	StoppingAuthority  bool                        `json:"stopping_authority"`
}
