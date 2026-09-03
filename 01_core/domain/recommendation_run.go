/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/recommendation_run.md
 * @layer L1
 * @updated 2026-09-01
 */

package domain

import (
	"context"
	"encoding/json"
	"errors"
)

type RunID string

type RecommendationLifecycle string

const (
	RunQueued      RecommendationLifecycle = "QUEUED"
	RunRunning     RecommendationLifecycle = "RUNNING"
	RunSucceeded   RecommendationLifecycle = "SUCCEEDED"
	RunFailed      RecommendationLifecycle = "FAILED"
	RunCancelling  RecommendationLifecycle = "CANCELLING"
	RunCancelled   RecommendationLifecycle = "CANCELLED"
	RunInterrupted RecommendationLifecycle = "INTERRUPTED"
)

type RecommendationRunStage string

const (
	RunStageQueued                   RecommendationRunStage = "queued"
	RunStageExpand                   RecommendationRunStage = "expand"
	RunStageDiscover                 RecommendationRunStage = "discover"
	RunStageScreen                   RecommendationRunStage = "screen"
	RunStageCollect                  RecommendationRunStage = "collect"
	RunStageAssess                   RecommendationRunStage = "assess"
	RunStageDecide                   RecommendationRunStage = "decide"
	RunStageAlternatives             RecommendationRunStage = "alternatives"
	RunStagePublish                  RecommendationRunStage = "publish"
	RunStageMarketPlan               RecommendationRunStage = "market_plan"
	RunStageMarketDiscover           RecommendationRunStage = "market_discover"
	RunStageMarketResolve            RecommendationRunStage = "market_resolve"
	RunStageMarketPublish            RecommendationRunStage = "market_publish"
	RunStageSynthesisCompare         RecommendationRunStage = "synthesis_compare"
	RunStageSynthesisDerive          RecommendationRunStage = "synthesis_derive"
	RunStageSynthesisCounterEvidence RecommendationRunStage = "synthesis_counterevidence"
	RunStageSynthesisValidate        RecommendationRunStage = "synthesis_validate"
	RunStageSynthesisPublish         RecommendationRunStage = "synthesis_publish"
	RunStageTerminal                 RecommendationRunStage = "terminal"
)

type RecommendationRunProgress struct {
	Planned   int `json:"planned"`
	Completed int `json:"completed"`
}

type RunArtifactName string

const (
	ArtifactRequest                RunArtifactName = "request"
	ArtifactRecommendation         RunArtifactName = "recommendation"
	ArtifactDecision               RunArtifactName = "decision"
	ArtifactAlternatives           RunArtifactName = "alternatives"
	ArtifactFailures               RunArtifactName = "failures"
	ArtifactProvenance             RunArtifactName = "provenance"
	ArtifactConnectorCapabilities  RunArtifactName = "connector_capabilities"
	ArtifactMarketPlan             RunArtifactName = "market_plan"
	ArtifactMarketSnapshot         RunArtifactName = "market_snapshot"
	ArtifactMarketCandidates       RunArtifactName = "market_candidates"
	ArtifactMarketGaps             RunArtifactName = "market_gaps"
	ArtifactMechanismCatalog       RunArtifactName = "mechanism_catalog"
	ArtifactMechanismGaps          RunArtifactName = "mechanism_gaps"
	ArtifactArchitectureComparison RunArtifactName = "architecture_comparison"
	ArtifactImprovementHypotheses  RunArtifactName = "improvement_hypotheses"
	ArtifactHypothesisGaps         RunArtifactName = "hypothesis_gaps"
	ArtifactValidationPlans        RunArtifactName = "validation_plans"
	ArtifactSynthesisProvenance    RunArtifactName = "synthesis_provenance"
)

type RunArtifactState string

const (
	ArtifactPending   RunArtifactState = "PENDING"
	ArtifactPublished RunArtifactState = "PUBLISHED"
)

type RunArtifactIndex struct {
	Name          RunArtifactName  `json:"name"`
	SchemaVersion string           `json:"schema_version"`
	Hash          string           `json:"artifact_hash"`
	ArtifactHash  string           `json:"-"`
	SizeBytes     int64            `json:"size_bytes"`
	EntryCount    int              `json:"entry_count"`
	Version       uint64           `json:"version"`
	State         RunArtifactState `json:"state"`
}

type RecommendationRunFailure struct {
	Code    string                 `json:"code"`
	Stage   RecommendationRunStage `json:"stage"`
	Message string                 `json:"message"`
}

type RecommendationRun struct {
	SchemaVersion         string                       `json:"schema_version"`
	RunID                 RunID                        `json:"run_id"`
	RunKind               RunKind                      `json:"run_kind,omitempty"`
	Request               RecommendationRequest        `json:"request"`
	MarketRequest         *MarketDiscoveryStartRequest `json:"market_request,omitempty"`
	ImprovementRequest    *ImprovementSynthesisRequest `json:"improvement_request,omitempty"`
	RequestHash           string                       `json:"request_hash"`
	Lifecycle             RecommendationLifecycle      `json:"lifecycle"`
	Version               uint64                       `json:"version"`
	Stage                 RecommendationRunStage       `json:"stage"`
	Progress              RecommendationRunProgress    `json:"progress"`
	RecommendationStatus  *IdeaRunStatus               `json:"recommendation_status"`
	MarketDiscoveryStatus *MarketDiscoveryStatus       `json:"market_discovery_status,omitempty"`
	ImprovementStatus     *ImprovementSynthesisStatus  `json:"improvement_status,omitempty"`
	Artifacts             []RunArtifactIndex           `json:"artifacts"`
	Failure               *RecommendationRunFailure    `json:"failure"`
}

// MarshalJSON evita que a representação pública acidental de um registro
// persistido exponha o pedido ou seu hash operacional.
func (r RecommendationRun) MarshalJSON() ([]byte, error) {
	if r.RunKind == "" && r.MarketDiscoveryStatus == nil {
		return json.Marshal(struct {
			SchemaVersion        string                    `json:"schema_version"`
			RunID                RunID                     `json:"run_id"`
			Lifecycle            RecommendationLifecycle   `json:"lifecycle"`
			Version              uint64                    `json:"version"`
			Stage                RecommendationRunStage    `json:"stage"`
			Progress             RecommendationRunProgress `json:"progress"`
			RecommendationStatus *IdeaRunStatus            `json:"recommendation_status"`
			Artifacts            []RunArtifactIndex        `json:"artifacts"`
			Failure              *RecommendationRunFailure `json:"failure"`
		}{r.SchemaVersion, r.RunID, r.Lifecycle, r.Version, r.Stage, r.Progress, r.RecommendationStatus, r.Artifacts, r.Failure})
	}
	return json.Marshal(RecommendationRunView{
		SchemaVersion: r.SchemaVersion, RunID: r.RunID, RunKind: r.RunKind, Lifecycle: r.Lifecycle,
		Version: r.Version, Stage: r.Stage, Progress: r.Progress,
		RecommendationStatus: r.RecommendationStatus, MarketDiscoveryStatus: r.MarketDiscoveryStatus, ImprovementStatus: r.ImprovementStatus, Artifacts: r.Artifacts, Failure: r.Failure,
	})
}

type IdempotencyRecord struct {
	KeyDigest   string `json:"key_digest"`
	RequestHash string `json:"request_hash"`
	RunID       RunID  `json:"run_id"`
}

type RecommendationStartReceipt struct {
	SchemaVersion string                  `json:"schema_version"`
	RunID         RunID                   `json:"run_id"`
	Lifecycle     RecommendationLifecycle `json:"lifecycle"`
	Version       uint64                  `json:"version"`
}

type RecommendationRunView struct {
	SchemaVersion         string                      `json:"schema_version"`
	RunID                 RunID                       `json:"run_id"`
	RunKind               RunKind                     `json:"run_kind,omitempty"`
	Lifecycle             RecommendationLifecycle     `json:"lifecycle"`
	Version               uint64                      `json:"version"`
	Stage                 RecommendationRunStage      `json:"stage"`
	Progress              RecommendationRunProgress   `json:"progress"`
	RecommendationStatus  *IdeaRunStatus              `json:"recommendation_status"`
	MarketDiscoveryStatus *MarketDiscoveryStatus      `json:"market_discovery_status"`
	ImprovementStatus     *ImprovementSynthesisStatus `json:"improvement_status"`
	Artifacts             []RunArtifactIndex          `json:"artifacts"`
	Failure               *RecommendationRunFailure   `json:"failure"`
}

func (r RecommendationRunView) MarshalJSON() ([]byte, error) {
	if r.RunKind == "" && r.MarketDiscoveryStatus == nil && r.ImprovementStatus == nil {
		return json.Marshal(struct {
			SchemaVersion        string                    `json:"schema_version"`
			RunID                RunID                     `json:"run_id"`
			Lifecycle            RecommendationLifecycle   `json:"lifecycle"`
			Version              uint64                    `json:"version"`
			Stage                RecommendationRunStage    `json:"stage"`
			Progress             RecommendationRunProgress `json:"progress"`
			RecommendationStatus *IdeaRunStatus            `json:"recommendation_status"`
			Artifacts            []RunArtifactIndex        `json:"artifacts"`
			Failure              *RecommendationRunFailure `json:"failure"`
		}{r.SchemaVersion, r.RunID, r.Lifecycle, r.Version, r.Stage, r.Progress, r.RecommendationStatus, r.Artifacts, r.Failure})
	}
	type publicView RecommendationRunView
	return json.Marshal(publicView(r))
}

type ArtifactReadRequest struct {
	RunID    RunID
	Artifact RunArtifactName
	Cursor   string
	Limit    int
}

type ArtifactPage struct {
	SchemaVersion string          `json:"schema_version"`
	RunID         RunID           `json:"run_id"`
	Artifact      RunArtifactName `json:"artifact"`
	ArtifactHash  string          `json:"artifact_hash"`
	Entries       [][]byte        `json:"-"`
	NextCursor    *string         `json:"next_cursor"`
}

type RunTransition struct {
	RunID                 RunID
	RunKind               RunKind
	ExpectedVersion       uint64
	Lifecycle             RecommendationLifecycle
	Stage                 RecommendationRunStage
	Progress              RecommendationRunProgress
	RecommendationStatus  *IdeaRunStatus
	MarketDiscoveryStatus *MarketDiscoveryStatus
	ImprovementStatus     *ImprovementSynthesisStatus
	Failure               *RecommendationRunFailure
}

type ArtifactPublication struct {
	RunID           RunID
	ExpectedVersion uint64
	Name            RunArtifactName
	SchemaVersion   string
	Document        []byte
}

type RecommendationProgressEvent struct {
	Stage    RecommendationRunStage
	Progress RecommendationRunProgress
}

type RecommendationRunStore interface {
	ResolveIdempotency(context.Context, string) (IdempotencyRecord, bool, error)
	CreateOrResolve(context.Context, RecommendationRun, IdempotencyRecord) (RecommendationRun, bool, error)
	Get(context.Context, RunID) (RecommendationRun, error)
	Transition(context.Context, RunTransition) (RecommendationRun, error)
	PublishArtifact(context.Context, ArtifactPublication) (RecommendationRun, error)
	ReadArtifact(context.Context, ArtifactReadRequest) (ArtifactPage, error)
	InterruptNonTerminal(context.Context) ([]RecommendationRun, error)
}

type RecommendationQueueReservation interface {
	Commit(RunID) error
	Release()
}

type RecommendationJobQueue interface {
	Reserve(context.Context) (RecommendationQueueReservation, error)
	Cancel(context.Context, RunID) error
	Shutdown(context.Context) error
}

type RecommendationRunIDGenerator interface {
	NewRunID() (RunID, error)
}

type ObservableRecommendationExecutor interface {
	Recommend(context.Context, RecommendationRequest) (RecommendationResult, error)
	RecommendObserved(context.Context, RecommendationRequest, RecommendationProgressReporter) (RecommendationResult, error)
}

type RecommendationProgressReporter interface {
	ReportRecommendationProgress(RecommendationProgressEvent)
}

var (
	ErrInvalidRunInput        = errors.New("invalid run input")
	ErrRunNotFound            = errors.New("run not found")
	ErrArtifactNotReady       = errors.New("artifact not ready")
	ErrIdempotencyConflict    = errors.New("idempotency conflict")
	ErrQueueFull              = errors.New("queue full")
	ErrQuotaExceeded          = errors.New("quota exceeded")
	ErrIntegrity              = errors.New("integrity error")
	ErrCursorInvalid          = errors.New("cursor invalid")
	ErrRunStoreUnavailable    = errors.New("run store unavailable")
	ErrRunVersionConflict     = errors.New("run version conflict")
	ErrInvalidRunTransition   = errors.New("invalid run transition")
	ErrArtifactImmutable      = errors.New("artifact immutable")
	ErrSourceBasisUnavailable = errors.New("source basis unavailable")
)
