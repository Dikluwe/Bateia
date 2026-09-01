/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_run.md
 * @layer L1
 * @updated 2026-08-26
 */

package domain

type IdeaRunMode string

const (
	IdeaRunLive   IdeaRunMode = "live"
	IdeaRunReplay IdeaRunMode = "replay"
	IdeaRunLocal  IdeaRunMode = "local"
)

type IdeaRunStatus string

const (
	IdeaRunRunning      IdeaRunStatus = "RUNNING"
	IdeaRunComplete     IdeaRunStatus = "COMPLETE"
	IdeaRunPartial      IdeaRunStatus = "PARTIAL"
	IdeaRunInconclusive IdeaRunStatus = "INCONCLUSIVE"
	IdeaRunFailed       IdeaRunStatus = "FAILED"
)

type IdeaRunStageStatus string

const (
	IdeaStagePending  IdeaRunStageStatus = "PENDING"
	IdeaStageRunning  IdeaRunStageStatus = "RUNNING"
	IdeaStageComplete IdeaRunStageStatus = "COMPLETE"
	IdeaStageSkipped  IdeaRunStageStatus = "SKIPPED"
	IdeaStageFailed   IdeaRunStageStatus = "FAILED"
)

type IdeaRunRequest struct {
	Idea            IdeaSpec            `json:"idea"`
	Policy          ReuseDecisionPolicy `json:"policy"`
	Mode            IdeaRunMode         `json:"mode"`
	Strategy        string              `json:"strategy"`
	Provider        string              `json:"provider,omitempty"`
	Language        string              `json:"language,omitempty"`
	SearchLocales   []string            `json:"search_locales,omitempty"`
	LocalPaths      []string            `json:"local_paths,omitempty"`
	ReplayDirectory string              `json:"replay_directory,omitempty"`
	OutputDirectory string              `json:"output_directory"`
	Resume          bool                `json:"resume"`
	MaxPerQuery     int                 `json:"max_per_query"`
	MaxCandidates   int                 `json:"max_candidates"`
}
type IdeaRunArtifact struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}
type IdeaRunStage struct {
	Name            string             `json:"name"`
	ContractVersion string             `json:"contract_version"`
	Status          IdeaRunStageStatus `json:"status"`
	InputHashes     map[string]string  `json:"input_hashes"`
	Outputs         []IdeaRunArtifact  `json:"outputs"`
	Counts          map[string]int64   `json:"counts,omitempty"`
	Message         string             `json:"message,omitempty"`
	DurationMillis  int64              `json:"duration_ms"`
}
type IdeaRunManifest struct {
	SchemaVersion string             `json:"schema_version"`
	RunID         string             `json:"run_id"`
	Mode          IdeaRunMode        `json:"mode"`
	Status        IdeaRunStatus      `json:"status"`
	InputHashes   map[string]string  `json:"input_hashes"`
	RuleVersions  map[string]string  `json:"rule_versions"`
	Stages        []IdeaRunStage     `json:"stages"`
	DecisionPath  string             `json:"decision_path,omitempty"`
	Failures      []ExecutionFailure `json:"failures,omitempty"`
}
type IdeaRunResult struct {
	Expansion      ExpansionPlan            `json:"expansion"`
	Discovery      DiscoverySnapshot        `json:"discovery"`
	Screening      ScreeningReport          `json:"screening"`
	Plans          []EvidenceExtractionPlan `json:"plans"`
	Bundles        []EvidenceBundle         `json:"bundles"`
	Assessment     AssessmentReport         `json:"assessment"`
	Decision       ReuseDecisionReport      `json:"decision"`
	Failures       []ExecutionFailure       `json:"failures,omitempty"`
	StageDurations map[string]int64         `json:"-"`
}
