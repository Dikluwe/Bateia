/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_run.md
 * @layer L1
 * @updated 2026-09-01
 */

package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

var runIDPattern = regexp.MustCompile(`^run_[0-9a-f]{32}$`)

var runArtifactOrder = []domain.RunArtifactName{
	domain.ArtifactRequest,
	domain.ArtifactRecommendation,
	domain.ArtifactDecision,
	domain.ArtifactAlternatives,
	domain.ArtifactFailures,
	domain.ArtifactProvenance,
}

type RecommendationRunService struct {
	store                  domain.RecommendationRunStore
	queue                  domain.RecommendationJobQueue
	ids                    domain.RecommendationRunIDGenerator
	executor               domain.ObservableRecommendationExecutor
	marketRegistry         *DiscoveryConnectorRegistry
	marketExecutor         domain.ObservableMarketDiscoveryExecutor
	mechanismExtractor     domain.MechanismExtractor
	mechanismCatalog       []byte
	mechanismCatalogSHA256 string
	mechanismBaseline      *domain.ConventionalBaseline
	improvementEnabled     bool
}

type cancellationViewMemo interface {
	CancellationView(domain.RunID) (domain.RecommendationRunView, bool)
	RememberCancellationView(domain.RunID, domain.RecommendationRunView)
	ForgetCancellationView(domain.RunID)
}

func NewRecommendationRunService(store domain.RecommendationRunStore, queue domain.RecommendationJobQueue, ids domain.RecommendationRunIDGenerator, executor domain.ObservableRecommendationExecutor) *RecommendationRunService {
	return &RecommendationRunService{store: store, queue: queue, ids: ids, executor: executor}
}

func CanonicalRecommendationRequestHash(request domain.RecommendationRequest) (domain.RecommendationRequest, string, error) {
	normalized := cloneRecommendationRequest(request)
	normalized.Idea.Users = canonicalSet(normalized.Idea.Users)
	normalized.Idea.Capabilities = canonicalSet(normalized.Idea.Capabilities)
	normalized.Idea.Constraints = canonicalSet(normalized.Idea.Constraints)
	normalized.Idea.NonGoals = canonicalSet(normalized.Idea.NonGoals)
	normalized.Idea.PreferredStack = canonicalSet(normalized.Idea.PreferredStack)
	normalized.Policy.AllowedLicenses = canonicalSet(normalized.Policy.AllowedLicenses)
	normalized.Search.Locales = canonicalSet(normalized.Search.Locales)
	validated, err := ValidateRecommendationRequest(normalized)
	if err != nil {
		return domain.RecommendationRequest{}, "", domain.ErrInvalidRunInput
	}
	encoded, err := json.Marshal(validated)
	if err != nil {
		return domain.RecommendationRequest{}, "", domain.ErrInvalidRunInput
	}
	digest := sha256.Sum256(encoded)
	return validated, "sha256:" + hex.EncodeToString(digest[:]), nil
}

func cloneRecommendationRequest(request domain.RecommendationRequest) domain.RecommendationRequest {
	request.Idea.Users = append([]string(nil), request.Idea.Users...)
	request.Idea.Capabilities = append([]string(nil), request.Idea.Capabilities...)
	request.Idea.Constraints = append([]string(nil), request.Idea.Constraints...)
	request.Idea.NonGoals = append([]string(nil), request.Idea.NonGoals...)
	request.Idea.PreferredStack = append([]string(nil), request.Idea.PreferredStack...)
	request.Policy.AllowedLicenses = append([]string(nil), request.Policy.AllowedLicenses...)
	request.Search.Locales = append([]string(nil), request.Search.Locales...)
	return request
}

func canonicalSet(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	if out == nil {
		return []string{}
	}
	return out
}

func validIdempotencyKey(key string) bool {
	if len(key) < 1 || len(key) > 128 || strings.Contains(key, "..") {
		return false
	}
	for _, value := range []byte(key) {
		if value < 0x20 || value > 0x7e || value == '/' || value == '\\' {
			return false
		}
	}
	return true
}

func digestString(value string) string {
	digest := sha256.Sum256([]byte(value))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func (s *RecommendationRunService) Start(ctx context.Context, key string, request domain.RecommendationRequest) (domain.RecommendationStartReceipt, error) {
	if s == nil || s.store == nil || s.queue == nil || s.ids == nil || s.executor == nil || !validIdempotencyKey(key) {
		return domain.RecommendationStartReceipt{}, domain.ErrInvalidRunInput
	}
	// A API síncrona rejeita locale duplicado. Para o hash idempotente, porém,
	// uma requisição integralmente set-like pode chegar reordenada/deduplicada.
	if hasDuplicate(request.Search.Locales) && !hasDuplicateInOtherSemanticSets(request) {
		return domain.RecommendationStartReceipt{}, domain.ErrInvalidRunInput
	}
	normalized, requestHash, err := CanonicalRecommendationRequestHash(request)
	if err != nil {
		return domain.RecommendationStartReceipt{}, domain.ErrInvalidRunInput
	}
	keyDigest := digestString(key)
	if record, ok, resolveErr := s.store.ResolveIdempotency(ctx, keyDigest); resolveErr != nil {
		return domain.RecommendationStartReceipt{}, resolveErr
	} else if ok {
		if record.RequestHash != requestHash {
			return domain.RecommendationStartReceipt{}, domain.ErrIdempotencyConflict
		}
		run, getErr := s.store.Get(ctx, record.RunID)
		if getErr != nil {
			return domain.RecommendationStartReceipt{}, getErr
		}
		return receiptFor(run), nil
	}
	reservation, err := s.queue.Reserve(ctx)
	if err != nil {
		return domain.RecommendationStartReceipt{}, err
	}
	consumed := false
	defer func() {
		if !consumed {
			reservation.Release()
		}
	}()
	id, err := s.ids.NewRunID()
	if err != nil || !runIDPattern.MatchString(string(id)) {
		return domain.RecommendationStartReceipt{}, domain.ErrRunStoreUnavailable
	}
	run := domain.RecommendationRun{
		SchemaVersion: "1", RunID: id, Request: normalized, RequestHash: requestHash,
		Lifecycle: domain.RunQueued, Version: 1, Stage: domain.RunStageQueued,
		Progress: domain.RecommendationRunProgress{Planned: 8}, Artifacts: []domain.RunArtifactIndex{},
	}
	if s.mechanismExtractor != nil {
		run.RunKind = domain.RunKindRecommendation
	}
	record := domain.IdempotencyRecord{KeyDigest: keyDigest, RequestHash: requestHash, RunID: id}
	actual, created, err := s.store.CreateOrResolve(ctx, run, record)
	if err != nil {
		return domain.RecommendationStartReceipt{}, err
	}
	if !created {
		if actual.RequestHash != requestHash {
			return domain.RecommendationStartReceipt{}, domain.ErrIdempotencyConflict
		}
		return receiptFor(actual), nil
	}
	if err := reservation.Commit(actual.RunID); err != nil {
		return domain.RecommendationStartReceipt{}, domain.ErrRunStoreUnavailable
	}
	consumed = true
	return receiptFor(actual), nil
}

func hasDuplicate(values []string) bool {
	seen := map[string]bool{}
	for _, value := range values {
		if seen[value] {
			return true
		}
		seen[value] = true
	}
	return false
}

func hasDuplicateInOtherSemanticSets(request domain.RecommendationRequest) bool {
	return hasDuplicate(request.Idea.Users) || hasDuplicate(request.Idea.Capabilities) ||
		hasDuplicate(request.Idea.Constraints) || hasDuplicate(request.Idea.NonGoals) ||
		hasDuplicate(request.Idea.PreferredStack) || hasDuplicate(request.Policy.AllowedLicenses)
}

func receiptFor(run domain.RecommendationRun) domain.RecommendationStartReceipt {
	return domain.RecommendationStartReceipt{SchemaVersion: "1", RunID: run.RunID, Lifecycle: run.Lifecycle, Version: run.Version}
}

func (s *RecommendationRunService) GetRun(ctx context.Context, id domain.RunID) (domain.RecommendationRunView, error) {
	if !runIDPattern.MatchString(string(id)) {
		return domain.RecommendationRunView{}, domain.ErrInvalidRunInput
	}
	run, err := s.store.Get(ctx, id)
	if err != nil {
		return domain.RecommendationRunView{}, err
	}
	if terminalLifecycle(run.Lifecycle) {
		if memo, ok := s.queue.(cancellationViewMemo); ok {
			memo.ForgetCancellationView(id)
		}
	}
	return projectRunView(run), nil
}

func (s *RecommendationRunService) GetArtifact(ctx context.Context, request domain.ArtifactReadRequest) (domain.ArtifactPage, error) {
	if !validArtifactRequest(request) {
		return domain.ArtifactPage{}, domain.ErrInvalidRunInput
	}
	run, err := s.store.Get(ctx, request.RunID)
	if err != nil {
		return domain.ArtifactPage{}, err
	}
	allowed := false
	for _, name := range artifactOrderForRun(run) {
		allowed = allowed || name == request.Artifact
	}
	if !allowed {
		return domain.ArtifactPage{}, domain.ErrArtifactNotReady
	}
	return s.store.ReadArtifact(ctx, request)
}

func validArtifactRequest(request domain.ArtifactReadRequest) bool {
	if !runIDPattern.MatchString(string(request.RunID)) || request.Limit < 1 || request.Limit > 50 || len(request.Cursor) > 512 {
		return false
	}
	for _, name := range allRunArtifactOrder() {
		if request.Artifact == name {
			return true
		}
	}
	return false
}

func projectRunView(run domain.RecommendationRun) domain.RecommendationRunView {
	byName := make(map[domain.RunArtifactName]domain.RunArtifactIndex, len(run.Artifacts))
	for _, artifact := range run.Artifacts {
		if artifact.State == domain.ArtifactPublished {
			if artifact.Hash == "" {
				artifact.Hash = artifact.ArtifactHash
			}
			byName[artifact.Name] = artifact
		}
	}
	artifacts := make([]domain.RunArtifactIndex, 0, len(byName))
	for _, name := range artifactOrderForRun(run) {
		if artifact, ok := byName[name]; ok {
			artifacts = append(artifacts, artifact)
		}
	}
	return domain.RecommendationRunView{
		SchemaVersion: "1", RunID: run.RunID, RunKind: run.RunKind, Lifecycle: run.Lifecycle, Version: run.Version,
		Stage: run.Stage, Progress: run.Progress, RecommendationStatus: run.RecommendationStatus,
		MarketDiscoveryStatus: run.MarketDiscoveryStatus,
		ImprovementStatus:     run.ImprovementStatus,
		Artifacts:             artifacts, Failure: run.Failure,
	}
}

func allRunArtifactOrder() []domain.RunArtifactName {
	return []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance, domain.ArtifactConnectorCapabilities, domain.ArtifactMarketPlan, domain.ArtifactMarketSnapshot, domain.ArtifactMarketCandidates, domain.ArtifactMarketGaps, domain.ArtifactMechanismCatalog, domain.ArtifactMechanismGaps, domain.ArtifactArchitectureComparison, domain.ArtifactImprovementHypotheses, domain.ArtifactHypothesisGaps, domain.ArtifactValidationPlans, domain.ArtifactSynthesisProvenance}
}

func artifactOrderForRun(run domain.RecommendationRun) []domain.RunArtifactName {
	if effectiveRunKind(run) == domain.RunKindMarketDiscovery {
		return []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactConnectorCapabilities, domain.ArtifactMarketPlan, domain.ArtifactMarketSnapshot, domain.ArtifactMarketCandidates, domain.ArtifactMarketGaps}
	}
	if effectiveRunKind(run) == domain.RunKindImprovementSynthesis {
		return []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactArchitectureComparison, domain.ArtifactImprovementHypotheses, domain.ArtifactHypothesisGaps, domain.ArtifactValidationPlans, domain.ArtifactSynthesisProvenance}
	}
	if run.RunKind == domain.RunKindRecommendation && sMechanismArtifactsExpected(run) {
		return append(append([]domain.RunArtifactName{}, runArtifactOrder...), domain.ArtifactMechanismCatalog, domain.ArtifactMechanismGaps)
	}
	return runArtifactOrder
}

func sMechanismArtifactsExpected(run domain.RecommendationRun) bool {
	return run.RunKind == domain.RunKindRecommendation
}

func (s *RecommendationRunService) Cancel(ctx context.Context, id domain.RunID) (domain.RecommendationRunView, error) {
	if !runIDPattern.MatchString(string(id)) {
		return domain.RecommendationRunView{}, domain.ErrInvalidRunInput
	}
	if memo, ok := s.queue.(cancellationViewMemo); ok {
		if view, found := memo.CancellationView(id); found {
			return view, nil
		}
	}
	for {
		run, err := s.store.Get(ctx, id)
		if err != nil {
			return domain.RecommendationRunView{}, err
		}
		switch run.Lifecycle {
		case domain.RunSucceeded, domain.RunFailed, domain.RunCancelled, domain.RunInterrupted, domain.RunCancelling:
			return projectRunView(run), nil
		case domain.RunQueued:
			if artifactPublished(run, domain.ArtifactRequest) {
				// A publicação é o CAS que elege exatamente um cancelador da
				// fila. Outro chamador aguarda o terminal durável desse líder.
				continue
			}
			entry := any(run.Request)
			if effectiveRunKind(run) == domain.RunKindMarketDiscovery && run.MarketRequest != nil {
				entry = *run.MarketRequest
			} else if effectiveRunKind(run) == domain.RunKindImprovementSynthesis && run.ImprovementRequest != nil {
				entry = *run.ImprovementRequest
			}
			run, err = s.publishDocument(ctx, run, domain.ArtifactRequest, []any{entry})
			if errors.Is(err, domain.ErrRunVersionConflict) {
				continue
			}
			if err != nil {
				return domain.RecommendationRunView{}, err
			}
			if err := s.queue.Cancel(ctx, id); err != nil {
				return domain.RecommendationRunView{}, err
			}
			updated, err := s.store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunCancelled, Stage: domain.RunStageTerminal, Progress: run.Progress})
			if errors.Is(err, domain.ErrRunVersionConflict) {
				continue
			}
			return projectRunView(updated), err
		case domain.RunRunning:
			updated, err := s.store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunCancelling, Stage: run.Stage, Progress: run.Progress})
			if errors.Is(err, domain.ErrRunVersionConflict) {
				continue
			}
			if err != nil {
				return domain.RecommendationRunView{}, err
			}
			view := projectRunView(updated)
			if memo, ok := s.queue.(cancellationViewMemo); ok {
				memo.RememberCancellationView(id, view)
			}
			if err := s.queue.Cancel(ctx, id); err != nil {
				if memo, ok := s.queue.(cancellationViewMemo); ok {
					memo.ForgetCancellationView(id)
				}
				return domain.RecommendationRunView{}, err
			}
			return view, nil
		default:
			return domain.RecommendationRunView{}, domain.ErrInvalidRunTransition
		}
	}
}

func (s *RecommendationRunService) Execute(ctx context.Context, id domain.RunID) error {
	if !runIDPattern.MatchString(string(id)) {
		return domain.ErrInvalidRunInput
	}
	run, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if terminalLifecycle(run.Lifecycle) {
		return nil
	}
	if effectiveRunKind(run) == domain.RunKindMarketDiscovery {
		return s.executeMarketDiscovery(ctx, run)
	}
	if effectiveRunKind(run) == domain.RunKindImprovementSynthesis {
		return s.ExecuteImprovementSynthesis(ctx, id)
	}
	if run.Lifecycle == domain.RunQueued {
		run, err = s.store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: domain.RecommendationRunProgress{Planned: 8}})
		if err != nil {
			return err
		}
	}
	if run.Lifecycle != domain.RunRunning {
		return nil
	}
	run, err = s.ensureRequestArtifact(context.Background(), id)
	if err != nil {
		return err
	}
	if run.Lifecycle == domain.RunCancelling {
		return s.confirmCancellation(context.Background(), id)
	}
	if run.Lifecycle != domain.RunRunning {
		return nil
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	reporter := &runProgressReporter{ctx: ctx, service: s, id: id, current: run}
	var execution RecommendationExecution
	var executionErr error
	if s.mechanismExtractor != nil {
		if detailed, ok := s.executor.(interface {
			RecommendObservedWithRun(context.Context, domain.RecommendationRequest, domain.RecommendationProgressReporter) (RecommendationExecution, error)
		}); ok {
			execution, executionErr = detailed.RecommendObservedWithRun(ctx, run.Request, reporter)
		} else {
			executionErr = domain.ErrInvalidRunInput
		}
	} else {
		execution.Recommendation, executionErr = s.executor.RecommendObserved(ctx, run.Request, reporter)
	}
	result := execution.Recommendation
	if reporter.err != nil {
		return reporter.err
	}
	run = reporter.current
	if executionErr != nil {
		if ctx.Err() != nil {
			return s.confirmCancellation(context.Background(), id)
		}
		return s.failRun(ctx, run, "RECOMMENDATION_FAILED", "recommendation execution failed")
	}
	result = ProjectRecommendationResult(result)
	if result.SchemaVersion != "1" || (result.Status != domain.IdeaRunComplete && result.Status != domain.IdeaRunPartial && result.Status != domain.IdeaRunInconclusive) {
		return s.failRun(ctx, run, "RESULT_INVALID", "recommendation result was invalid")
	}
	documents := []struct {
		name    domain.RunArtifactName
		entries []any
	}{
		{domain.ArtifactRecommendation, []any{result}},
		{domain.ArtifactDecision, []any{result.Decision}},
		{domain.ArtifactAlternatives, []any{result.Alternatives}},
		{domain.ArtifactFailures, interfacesForFailures(result.Failures)},
		{domain.ArtifactProvenance, []any{}},
	}
	for _, document := range documents {
		run, err = s.publishDocument(ctx, run, document.name, document.entries)
		if err != nil {
			return err
		}
	}
	if s.mechanismExtractor != nil {
		mechanisms := []domain.ArchitecturalMechanism{}
		gaps := []domain.MechanismGap{}
		for _, bundle := range execution.Run.Bundles {
			extracted, extractErr := s.mechanismExtractor.ExtractMechanisms(ctx, domain.MechanismExtractionInput{Bundle: bundle, Assessments: execution.Run.Assessment.Assessments, ConventionalBaseline: s.mechanismBaseline, RuleCatalogJSON: s.mechanismCatalog, RuleCatalogSHA256: s.mechanismCatalogSHA256})
			if extractErr != nil {
				return s.failRun(ctx, run, "MECHANISM_EXTRACTION_FAILED", "mechanism extraction failed")
			}
			mechanisms = append(mechanisms, extracted.Mechanisms...)
			gaps = append(gaps, extracted.Gaps...)
		}
		sort.Slice(mechanisms, func(i, j int) bool { return mechanisms[i].ID < mechanisms[j].ID })
		sortMechanismGaps(gaps)
		run, err = s.publishDocument(ctx, run, domain.ArtifactMechanismCatalog, interfacesForSlice(mechanisms))
		if err != nil {
			return err
		}
		run, err = s.publishDocument(ctx, run, domain.ArtifactMechanismGaps, interfacesForSlice(gaps))
		if err != nil {
			return err
		}
	}
	status := result.Status
	_, err = s.store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunSucceeded, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 8}, RecommendationStatus: &status})
	if errors.Is(err, domain.ErrRunVersionConflict) {
		winner, getErr := s.store.Get(ctx, id)
		if getErr == nil && winner.Lifecycle == domain.RunCancelling {
			return s.confirmCancellation(context.Background(), id)
		}
	}
	return err
}

func (s *RecommendationRunService) ensureRequestArtifact(ctx context.Context, id domain.RunID) (domain.RecommendationRun, error) {
	for {
		run, err := s.store.Get(ctx, id)
		if err != nil {
			return domain.RecommendationRun{}, err
		}
		if artifactPublished(run, domain.ArtifactRequest) {
			return run, nil
		}
		entry := any(run.Request)
		if effectiveRunKind(run) == domain.RunKindMarketDiscovery && run.MarketRequest != nil {
			entry = *run.MarketRequest
		} else if effectiveRunKind(run) == domain.RunKindImprovementSynthesis && run.ImprovementRequest != nil {
			entry = *run.ImprovementRequest
		}
		updated, err := s.publishDocument(ctx, run, domain.ArtifactRequest, []any{entry})
		if errors.Is(err, domain.ErrRunVersionConflict) {
			continue
		}
		return updated, err
	}
}

type artifactDocument struct {
	SchemaVersion string                 `json:"schema_version"`
	Artifact      domain.RunArtifactName `json:"artifact"`
	Entries       []any                  `json:"entries"`
}

func (s *RecommendationRunService) publishDocument(ctx context.Context, run domain.RecommendationRun, name domain.RunArtifactName, entries []any) (domain.RecommendationRun, error) {
	if entries == nil {
		entries = []any{}
	}
	document, err := json.Marshal(artifactDocument{SchemaVersion: "1", Artifact: name, Entries: entries})
	if err != nil {
		return run, domain.ErrRunStoreUnavailable
	}
	return s.store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: run.RunID, ExpectedVersion: run.Version, Name: name, SchemaVersion: "1", Document: document})
}

func interfacesForFailures(values []domain.ExecutionFailure) []any {
	out := make([]any, len(values))
	for index := range values {
		out[index] = values[index]
	}
	return out
}

func (s *RecommendationRunService) failRun(ctx context.Context, run domain.RecommendationRun, code, message string) error {
	if !artifactPublished(run, domain.ArtifactFailures) {
		var err error
		run, err = s.publishDocument(ctx, run, domain.ArtifactFailures, []any{})
		if err != nil {
			return err
		}
	}
	failure := &domain.RecommendationRunFailure{Code: code, Stage: run.Stage, Message: message}
	_, err := s.store.Transition(ctx, domain.RunTransition{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunFailed, Stage: domain.RunStageTerminal, Progress: run.Progress, Failure: failure})
	return err
}

func (s *RecommendationRunService) confirmCancellation(ctx context.Context, id domain.RunID) error {
	for {
		run, err := s.store.Get(ctx, id)
		if err != nil {
			return err
		}
		if run.Lifecycle == domain.RunCancelled || terminalLifecycle(run.Lifecycle) {
			return nil
		}
		if run.Lifecycle != domain.RunCancelling {
			return nil
		}
		_, err = s.store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunCancelled, Stage: domain.RunStageTerminal, Progress: run.Progress})
		if errors.Is(err, domain.ErrRunVersionConflict) {
			continue
		}
		return err
	}
}

type runProgressReporter struct {
	ctx     context.Context
	service *RecommendationRunService
	id      domain.RunID
	current domain.RecommendationRun
	err     error
}

func (r *runProgressReporter) ReportRecommendationProgress(event domain.RecommendationProgressEvent) {
	if r.err != nil || r.ctx.Err() != nil {
		return
	}
	if event.Stage == r.current.Stage && event.Progress == r.current.Progress {
		return
	}
	transition := domain.RunTransition{RunID: r.id, ExpectedVersion: r.current.Version, Lifecycle: domain.RunRunning, Stage: event.Stage, Progress: event.Progress}
	updated, err := r.service.store.Transition(r.ctx, transition)
	if err != nil {
		r.err = err
		return
	}
	r.current = updated
}

func (s *RecommendationRunService) Recover(ctx context.Context) error {
	_, err := s.store.InterruptNonTerminal(ctx)
	return err
}

func terminalLifecycle(value domain.RecommendationLifecycle) bool {
	return value == domain.RunSucceeded || value == domain.RunFailed || value == domain.RunCancelled || value == domain.RunInterrupted
}

func effectiveRunKind(run domain.RecommendationRun) domain.RunKind {
	if run.RunKind == "" {
		return domain.RunKindRecommendation
	}
	return run.RunKind
}

func artifactPublished(run domain.RecommendationRun, name domain.RunArtifactName) bool {
	for _, artifact := range run.Artifacts {
		if artifact.Name == name && artifact.State == domain.ArtifactPublished {
			return true
		}
	}
	return false
}

func ValidateRecommendationRunTransition(run domain.RecommendationRun, transition domain.RunTransition) error {
	kind := effectiveRunKind(run)
	planned := 8
	if kind == domain.RunKindMarketDiscovery {
		planned = 4
	} else if kind == domain.RunKindImprovementSynthesis {
		planned = 5
	}
	if transition.RunKind != "" && transition.RunKind != kind || transition.RunID != run.RunID || transition.ExpectedVersion != run.Version || transition.Progress.Planned != planned || transition.Progress.Completed < 0 || transition.Progress.Completed > planned {
		return domain.ErrInvalidRunTransition
	}
	allowed := false
	switch run.Lifecycle {
	case domain.RunQueued:
		allowed = transition.Lifecycle == domain.RunRunning || transition.Lifecycle == domain.RunCancelled || transition.Lifecycle == domain.RunInterrupted
	case domain.RunRunning:
		allowed = transition.Lifecycle == domain.RunRunning || transition.Lifecycle == domain.RunCancelling || transition.Lifecycle == domain.RunSucceeded || transition.Lifecycle == domain.RunFailed || transition.Lifecycle == domain.RunInterrupted
	case domain.RunCancelling:
		allowed = transition.Lifecycle == domain.RunCancelled || transition.Lifecycle == domain.RunInterrupted
	}
	if !allowed {
		return domain.ErrInvalidRunTransition
	}
	if transition.Progress.Completed < run.Progress.Completed {
		return domain.ErrInvalidRunTransition
	}
	if transition.Lifecycle == domain.RunRunning {
		if !validRunStageForKind(transition.Stage, kind) || transition.Stage == domain.RunStageQueued || transition.Stage == domain.RunStageTerminal || transition.RecommendationStatus != nil || transition.MarketDiscoveryStatus != nil || transition.ImprovementStatus != nil || transition.Failure != nil {
			return domain.ErrInvalidRunTransition
		}
		if run.Lifecycle == domain.RunRunning && transition.Progress.Completed <= run.Progress.Completed {
			return domain.ErrInvalidRunTransition
		}
		if stagePositionForKind(transition.Stage, kind) < stagePositionForKind(run.Stage, kind) {
			return domain.ErrInvalidRunTransition
		}
		return nil
	}
	if transition.Lifecycle == domain.RunCancelling {
		if transition.Stage != run.Stage || transition.Progress != run.Progress || transition.RecommendationStatus != nil || transition.MarketDiscoveryStatus != nil || transition.ImprovementStatus != nil || transition.Failure != nil {
			return domain.ErrInvalidRunTransition
		}
		return nil
	}
	if transition.Stage != domain.RunStageTerminal || transition.Progress.Completed != run.Progress.Completed && transition.Lifecycle != domain.RunSucceeded {
		return domain.ErrInvalidRunTransition
	}
	switch transition.Lifecycle {
	case domain.RunSucceeded:
		if transition.Progress.Completed != planned || transition.Failure != nil || !allArtifactsPublished(run) {
			return domain.ErrInvalidRunTransition
		}
		if kind == domain.RunKindMarketDiscovery {
			if transition.RecommendationStatus != nil || transition.MarketDiscoveryStatus == nil || transition.ImprovementStatus != nil {
				return domain.ErrInvalidRunTransition
			}
			status := *transition.MarketDiscoveryStatus
			if status != domain.MarketDiscoveryComplete && status != domain.MarketDiscoveryPartial && status != domain.MarketDiscoveryInconclusive {
				return domain.ErrInvalidRunTransition
			}
		} else if kind == domain.RunKindImprovementSynthesis {
			if transition.RecommendationStatus != nil || transition.MarketDiscoveryStatus != nil || transition.ImprovementStatus == nil {
				return domain.ErrInvalidRunTransition
			}
			status := *transition.ImprovementStatus
			if status != domain.ImprovementComplete && status != domain.ImprovementPartial && status != domain.ImprovementInconclusive {
				return domain.ErrInvalidRunTransition
			}
		} else {
			if transition.RecommendationStatus == nil || transition.MarketDiscoveryStatus != nil || transition.ImprovementStatus != nil {
				return domain.ErrInvalidRunTransition
			}
			status := *transition.RecommendationStatus
			if status != domain.IdeaRunComplete && status != domain.IdeaRunPartial && status != domain.IdeaRunInconclusive {
				return domain.ErrInvalidRunTransition
			}
		}
	case domain.RunFailed:
		if transition.RecommendationStatus != nil || transition.MarketDiscoveryStatus != nil || transition.ImprovementStatus != nil || transition.Failure == nil {
			return domain.ErrInvalidRunTransition
		}
	case domain.RunCancelled, domain.RunInterrupted:
		if transition.RecommendationStatus != nil || transition.MarketDiscoveryStatus != nil || transition.ImprovementStatus != nil || transition.Failure != nil {
			return domain.ErrInvalidRunTransition
		}
	}
	return nil
}

func validRunStage(stage domain.RecommendationRunStage) bool {
	return stagePosition(stage) >= 0
}

func validRunStageForKind(stage domain.RecommendationRunStage, kind domain.RunKind) bool {
	return stagePositionForKind(stage, kind) >= 0
}

func stagePositionForKind(stage domain.RecommendationRunStage, kind domain.RunKind) int {
	if kind == domain.RunKindMarketDiscovery {
		stages := []domain.RecommendationRunStage{domain.RunStageQueued, domain.RunStageMarketPlan, domain.RunStageMarketDiscover, domain.RunStageMarketResolve, domain.RunStageMarketPublish, domain.RunStageTerminal}
		for index, value := range stages {
			if value == stage {
				return index
			}
		}
		return -1
	}
	if kind == domain.RunKindImprovementSynthesis {
		stages := []domain.RecommendationRunStage{domain.RunStageQueued, domain.RunStageSynthesisCompare, domain.RunStageSynthesisDerive, domain.RunStageSynthesisCounterEvidence, domain.RunStageSynthesisValidate, domain.RunStageSynthesisPublish, domain.RunStageTerminal}
		for index, value := range stages {
			if value == stage {
				return index
			}
		}
		return -1
	}
	return stagePosition(stage)
}

func stagePosition(stage domain.RecommendationRunStage) int {
	stages := []domain.RecommendationRunStage{domain.RunStageQueued, domain.RunStageExpand, domain.RunStageDiscover, domain.RunStageScreen, domain.RunStageCollect, domain.RunStageAssess, domain.RunStageDecide, domain.RunStageAlternatives, domain.RunStagePublish, domain.RunStageTerminal}
	for index, value := range stages {
		if value == stage {
			return index
		}
	}
	return -1
}

func allArtifactsPublished(run domain.RecommendationRun) bool {
	requiredArtifacts := runArtifactOrder
	if effectiveRunKind(run) == domain.RunKindMarketDiscovery {
		requiredArtifacts = []domain.RunArtifactName{domain.ArtifactConnectorCapabilities, domain.ArtifactMarketPlan, domain.ArtifactMarketSnapshot, domain.ArtifactMarketCandidates, domain.ArtifactMarketGaps}
	}
	if effectiveRunKind(run) == domain.RunKindImprovementSynthesis {
		requiredArtifacts = []domain.RunArtifactName{domain.ArtifactArchitectureComparison, domain.ArtifactImprovementHypotheses, domain.ArtifactHypothesisGaps, domain.ArtifactValidationPlans, domain.ArtifactSynthesisProvenance}
	} else if run.RunKind == domain.RunKindRecommendation {
		requiredArtifacts = append(append([]domain.RunArtifactName{}, runArtifactOrder...), domain.ArtifactMechanismCatalog, domain.ArtifactMechanismGaps)
	}
	for _, required := range requiredArtifacts {
		if !artifactPublished(run, required) {
			return false
		}
	}
	return true
}
