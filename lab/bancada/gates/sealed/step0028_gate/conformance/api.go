/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/api.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import "context"

type Obligation string

const (
	P01ExtractionBinding   Obligation = "P01_EXTRACTION_BINDING"
	P02EpistemicSeparation Obligation = "P02_EPISTEMIC_SEPARATION"
	P03ProductOnly         Obligation = "P03_SOURCE_PREVALIDATION"
	P04UnknownTarget       Obligation = "P04_CANONICAL_IDEMPOTENCY"
	P05License             Obligation = "P05_LIFECYCLE_STATUS"
	P06Composition         Obligation = "P06_CANCEL_RECOVERY"
	P07CounterEvidence     Obligation = "P07_OFFLINE_REPLAY"
	P08PopularityNovelty   Obligation = "P08_RUNTIME_IDENTITY"
	P09Idempotency         Obligation = "P09_MCP_REFINEMENT"
	P10CancelRecovery      Obligation = "P10_PREDECESSOR_REFINEMENT"
	P11OfflineReplay       Obligation = "P11_HYPOTHESIS_SAFETY"
	P12RefinementRuntime   Obligation = "P12_PROJECTION_SAFETY"
)

type Operation string

const (
	OperationExtract    Operation = "EXTRACT"
	OperationStart      Operation = "START"
	OperationLifecycle  Operation = "LIFECYCLE"
	OperationReplay     Operation = "REPLAY"
	OperationRuntime    Operation = "RUNTIME"
	OperationTools      Operation = "TOOLS"
	OperationRefinement Operation = "REFINEMENT"
	OperationHypothesis Operation = "HYPOTHESIS"
	OperationSafety     Operation = "SAFETY"
	OperationArtifact   Operation = "ARTIFACT"
	OperationPagination Operation = "PAGINATION"
	OperationSchema     Operation = "SCHEMA"
	OperationFull       Operation = "FULL_REFINEMENT"
)

type Tool struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type Input struct {
	ScenarioVersion          string            `json:"scenario_version"`
	SemanticVariantID        string            `json:"semantic_variant_id"`
	EntryCommitmentSHA256    string            `json:"entry_commitment_sha256"`
	SemanticFault            string            `json:"semantic_fault"`
	Opaque                   bool              `json:"opaque"`
	OpacityReason            string            `json:"opacity_reason"`
	CatalogState             string            `json:"catalog_state"`
	RuleState                string            `json:"rule_state"`
	RuleFromContent          bool              `json:"rule_from_content"`
	SourceKindState          string            `json:"source_kind_state"`
	Implementation           bool              `json:"implementation"`
	Independent              bool              `json:"independent"`
	DocumentationOnly        bool              `json:"documentation_only"`
	ProductOnly              bool              `json:"product_only"`
	SameImmutableSource      bool              `json:"same_immutable_source"`
	ExactBundleBinding       bool              `json:"exact_bundle_binding"`
	ContentHashValid         bool              `json:"content_hash_valid"`
	ParentStatus             string            `json:"parent_status"`
	CentralConflict          bool              `json:"central_conflict"`
	ClassInBaseline          bool              `json:"class_in_baseline"`
	SourceExists             bool              `json:"source_exists"`
	SourceTerminal           bool              `json:"source_terminal"`
	SourcePublished          bool              `json:"source_published"`
	SourceHashMatches        bool              `json:"source_hash_matches"`
	CanonicalDocument        bool              `json:"canonical_document"`
	PriorRequestHash         string            `json:"prior_request_hash"`
	RequestHash              string            `json:"request_hash"`
	PriorRunID               string            `json:"prior_run_id"`
	CanonicalItems           []string          `json:"canonical_items"`
	ConcurrentStarts         int               `json:"concurrent_starts"`
	Lifecycle                string            `json:"lifecycle"`
	TargetState              string            `json:"target_state"`
	MechanismStatus          string            `json:"mechanism_status"`
	CounterSearch            string            `json:"counter_search"`
	Cancelled                bool              `json:"cancelled"`
	RecoveryRequested        bool              `json:"recovery_requested"`
	ReplayItems              []string          `json:"replay_items"`
	CatalogHashMatches       bool              `json:"catalog_hash_matches"`
	CaptureHashMatches       bool              `json:"capture_hash_matches"`
	NetworkAvailable         bool              `json:"network_available"`
	StoreIDs                 []string          `json:"store_ids"`
	QueueIDs                 []string          `json:"queue_ids"`
	RootIDs                  []string          `json:"root_ids"`
	SecretIDs                []string          `json:"secret_ids"`
	WorkerIDs                []string          `json:"worker_ids"`
	BasalTools               []Tool            `json:"basal_tools"`
	BasalToolsSHA256         string            `json:"basal_tools_sha256"`
	ActualBasalToolsSHA256   string            `json:"actual_basal_tools_sha256"`
	AppendTool10             bool              `json:"append_tool_10"`
	RefinementHashes         map[string]string `json:"refinement_hashes"`
	RequiredRefinementHashes map[string]string `json:"required_refinement_hashes"`
	License                  string            `json:"license"`
	RequestedTransfer        string            `json:"requested_transfer"`
	MechanismCount           int               `json:"mechanism_count"`
	CompatibilityWitnessed   bool              `json:"compatibility_witnessed"`
	Popularity               int               `json:"popularity"`
	AdoptionRate             float64           `json:"adoption_rate"`
	ClaimedNovelty           string            `json:"claimed_novelty"`
	RawFields                []string          `json:"raw_fields"`
	ExternalEffects          []string          `json:"external_effects"`
	ArtifactEntries          []string          `json:"artifact_entries"`
	ArtifactNames            []string          `json:"artifact_names"`
	ArtifactsCanonical       bool              `json:"artifacts_canonical"`
	ArtifactsPublished       bool              `json:"artifacts_published"`
	SourceRefs               []string          `json:"source_refs"`
	CounterScope             []string          `json:"counter_scope"`
	Owner                    string            `json:"owner"`
	Repo                     string            `json:"repo"`
	Commit                   string            `json:"commit"`
	Path                     string            `json:"path"`
	Language                 string            `json:"language"`
	Ecosystem                string            `json:"ecosystem"`
	SourceHashes             []string          `json:"source_hashes"`
	SourceKinds              []string          `json:"source_kinds"`
	RuleIDs                  []string          `json:"rule_ids"`
	WitnessRoles             []string          `json:"witness_roles"`
	AppliedLimits            map[string]int    `json:"applied_limits"`
	MechanismIDs             []string          `json:"mechanism_ids"`
	HypothesisIDs            []string          `json:"hypothesis_ids"`
	ArtifactIDs              []string          `json:"artifact_ids"`
	ForceCollision           bool              `json:"force_collision"`
	ValidationPlan           []string          `json:"validation_plan"`
	TargetGapRefs            []string          `json:"target_gap_refs"`
	MechanismRefs            []string          `json:"mechanism_refs"`
	CounterEvidenceRefs      []string          `json:"counterevidence_refs"`
	Assumptions              []string          `json:"assumptions"`
	Preconditions            []string          `json:"preconditions"`
	Risks                    []string          `json:"risks"`
	Tradeoffs                []string          `json:"tradeoffs"`
	Ordering                 []string          `json:"ordering"`
	Fit                      string            `json:"fit"`
	ImprovementStatus        string            `json:"improvement_status"`
	ArtifactState            string            `json:"artifact_state"`
	ArtifactTransitions      []string          `json:"artifact_transitions"`
	ArtifactHashes           map[string]string `json:"artifact_hashes"`
	IndexHashes              map[string]string `json:"index_hashes"`
	SnapshotSHA256           string            `json:"snapshot_sha256"`
	Cursor                   string            `json:"cursor"`
	NextCursor               string            `json:"next_cursor"`
	PageItems                []string          `json:"page_items"`
	ToolDefinitions          map[string]string `json:"tool_definitions"`
	SchemaDefinitions        map[string]string `json:"schema_definitions"`
	RefinementReceipts       map[string]string `json:"refinement_receipts"`
	LayerFindings            []string          `json:"layer_findings"`
}

type Observation struct {
	Classification       string            `json:"classification"`
	Decision             string            `json:"decision"`
	MechanismStatus      string            `json:"mechanism_status"`
	Novelty              string            `json:"novelty"`
	HypothesisStatus     string            `json:"hypothesis_status"`
	TransferMode         string            `json:"transfer_mode"`
	WarningCodes         []string          `json:"warning_codes"`
	IntegrationUnknowns  []string          `json:"integration_unknowns"`
	ErrorCode            string            `json:"error_code"`
	GapCode              string            `json:"gap_code"`
	ResolveCalls         int               `json:"resolve_calls"`
	ReserveCalls         int               `json:"reserve_calls"`
	IDCalls              int               `json:"id_calls"`
	CreateCalls          int               `json:"create_calls"`
	RunID                string            `json:"run_id"`
	RunIDs               []string          `json:"run_ids"`
	CanonicalSHA256      string            `json:"canonical_sha256"`
	Lifecycle            string            `json:"lifecycle"`
	CompletedSteps       int               `json:"completed_steps"`
	RecoveryState        string            `json:"recovery_state"`
	ReplaySHA256         string            `json:"replay_sha256"`
	NetworkCalls         int               `json:"network_calls"`
	RuntimeUnique        bool              `json:"runtime_unique"`
	Tools                []Tool            `json:"tools"`
	BasalHashVerified    bool              `json:"basal_hash_verified"`
	RefinementsPreserved bool              `json:"refinements_preserved"`
	LeakedFields         []string          `json:"leaked_fields"`
	ExternalEffects      []string          `json:"external_effects"`
	AdoptionDisposition  string            `json:"adoption_disposition"`
	BindingVerified      bool              `json:"binding_verified"`
	ProvenanceComplete   bool              `json:"provenance_complete"`
	ArtifactSHA256       string            `json:"artifact_sha256"`
	ArtifactNames        []string          `json:"artifact_names"`
	SourceRefs           []string          `json:"source_refs"`
	CounterScope         []string          `json:"counter_scope"`
	Owner                string            `json:"owner"`
	Repo                 string            `json:"repo"`
	Commit               string            `json:"commit"`
	Path                 string            `json:"path"`
	Language             string            `json:"language"`
	Ecosystem            string            `json:"ecosystem"`
	SourceHashes         []string          `json:"source_hashes"`
	SourceKinds          []string          `json:"source_kinds"`
	RuleIDs              []string          `json:"rule_ids"`
	WitnessRoles         []string          `json:"witness_roles"`
	AppliedLimits        map[string]int    `json:"applied_limits"`
	MechanismIDs         []string          `json:"mechanism_ids"`
	HypothesisIDs        []string          `json:"hypothesis_ids"`
	ArtifactIDs          []string          `json:"artifact_ids"`
	CollisionDetected    bool              `json:"collision_detected"`
	ValidationPlan       []string          `json:"validation_plan"`
	TargetState          string            `json:"target_state"`
	TargetGapRefs        []string          `json:"target_gap_refs"`
	MechanismRefs        []string          `json:"mechanism_refs"`
	CounterEvidenceRefs  []string          `json:"counterevidence_refs"`
	Assumptions          []string          `json:"assumptions"`
	Preconditions        []string          `json:"preconditions"`
	Risks                []string          `json:"risks"`
	Tradeoffs            []string          `json:"tradeoffs"`
	Ordering             []string          `json:"ordering"`
	License              string            `json:"license"`
	Fit                  string            `json:"fit"`
	ImprovementStatus    string            `json:"improvement_status"`
	ArtifactState        string            `json:"artifact_state"`
	ArtifactTransitions  []string          `json:"artifact_transitions"`
	ArtifactHashes       map[string]string `json:"artifact_hashes"`
	IndexHashes          map[string]string `json:"index_hashes"`
	SnapshotSHA256       string            `json:"snapshot_sha256"`
	Cursor               string            `json:"cursor"`
	NextCursor           string            `json:"next_cursor"`
	PageItems            []string          `json:"page_items"`
	StoreIdentity        string            `json:"store_identity"`
	QueueIdentity        string            `json:"queue_identity"`
	RootIdentity         string            `json:"root_identity"`
	SecretIdentity       string            `json:"secret_identity"`
	WorkerIdentity       string            `json:"worker_identity"`
	ToolDefinitions      map[string]string `json:"tool_definitions"`
	SchemaDefinitions    map[string]string `json:"schema_definitions"`
	RefinementReceipts   map[string]string `json:"refinement_receipts"`
	LayerFindings        []string          `json:"layer_findings"`
}

type Case struct {
	ID         string      `json:"id"`
	Obligation Obligation  `json:"obligation"`
	Operation  Operation   `json:"operation"`
	Input      Input       `json:"input"`
	Expected   Observation `json:"expected"`
}

type Subject interface {
	Observe(context.Context, Case) (Observation, error)
}

type Finding struct {
	CaseID     string     `json:"case_id"`
	Obligation Obligation `json:"obligation"`
	Outcome    string     `json:"outcome"`
	Detail     string     `json:"detail"`
}
