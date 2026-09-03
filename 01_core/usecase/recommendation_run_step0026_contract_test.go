/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_run_step0026_contract_test.md
 * @layer L1
 * @updated 2026-09-01
 */
package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"bateia/01_core/domain"
)

const step0026RunID domain.RunID = "run_0123456789abcdef0123456789abcdef"

func step0026Request() domain.RecommendationRequest {
	return domain.RecommendationRequest{
		Idea:   domain.IdeaSpec{Problem: "food waste", Users: []string{"families", "homes"}, Capabilities: []string{"forecast", "inventory"}, Constraints: []string{"private"}, NonGoals: []string{"ads"}, PreferredStack: []string{"Go"}},
		Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"Apache-2.0", "MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .3},
		Target: domain.SoftwareAlternativeTarget{Name: "pantry"},
		Search: domain.RecommendationSearch{Strategy: "literal", Locales: []string{"en-US", "pt-BR"}, MaxPerQuery: 2, MaxCandidates: 4},
	}
}

type step0026MemoryStore struct {
	mu          sync.Mutex
	runs        map[domain.RunID]domain.RecommendationRun
	idem        map[string]domain.IdempotencyRecord
	artifacts   map[domain.RunID]map[domain.RunArtifactName][]byte
	creates     int
	transitions int
	createErr   error
}

func newStep0026MemoryStore() *step0026MemoryStore {
	return &step0026MemoryStore{runs: map[domain.RunID]domain.RecommendationRun{}, idem: map[string]domain.IdempotencyRecord{}, artifacts: map[domain.RunID]map[domain.RunArtifactName][]byte{}}
}

func (s *step0026MemoryStore) ResolveIdempotency(_ context.Context, digest string) (domain.IdempotencyRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.idem[digest]
	return record, ok, nil
}

func (s *step0026MemoryStore) CreateOrResolve(_ context.Context, run domain.RecommendationRun, record domain.IdempotencyRecord) (domain.RecommendationRun, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.createErr != nil {
		return domain.RecommendationRun{}, false, s.createErr
	}
	if prior, ok := s.idem[record.KeyDigest]; ok {
		return s.runs[prior.RunID], false, nil
	}
	s.creates++
	s.idem[record.KeyDigest] = record
	s.runs[run.RunID] = run
	return run, true, nil
}

func (s *step0026MemoryStore) Get(_ context.Context, id domain.RunID) (domain.RecommendationRun, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return domain.RecommendationRun{}, domain.ErrRunNotFound
	}
	return run, nil
}

func (s *step0026MemoryStore) Transition(_ context.Context, transition domain.RunTransition) (domain.RecommendationRun, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[transition.RunID]
	if !ok {
		return domain.RecommendationRun{}, domain.ErrRunNotFound
	}
	if run.Version != transition.ExpectedVersion {
		return domain.RecommendationRun{}, domain.ErrRunVersionConflict
	}
	if err := ValidateRecommendationRunTransition(run, transition); err != nil {
		return domain.RecommendationRun{}, err
	}
	run.Lifecycle, run.Stage, run.Progress = transition.Lifecycle, transition.Stage, transition.Progress
	run.RecommendationStatus, run.Failure = transition.RecommendationStatus, transition.Failure
	run.Version++
	s.transitions++
	s.runs[run.RunID] = run
	return run, nil
}

func (s *step0026MemoryStore) PublishArtifact(_ context.Context, publication domain.ArtifactPublication) (domain.RecommendationRun, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[publication.RunID]
	if !ok {
		return domain.RecommendationRun{}, domain.ErrRunNotFound
	}
	if run.Version != publication.ExpectedVersion {
		return domain.RecommendationRun{}, domain.ErrRunVersionConflict
	}
	if s.artifacts[run.RunID] == nil {
		s.artifacts[run.RunID] = map[domain.RunArtifactName][]byte{}
	}
	if _, exists := s.artifacts[run.RunID][publication.Name]; exists {
		return domain.RecommendationRun{}, domain.ErrArtifactImmutable
	}
	s.artifacts[run.RunID][publication.Name] = append([]byte(nil), publication.Document...)
	run.Version++
	run.Artifacts = append(run.Artifacts, domain.RunArtifactIndex{Name: publication.Name, SchemaVersion: "1", State: domain.ArtifactPublished, Version: run.Version})
	s.runs[run.RunID] = run
	return run, nil
}

func (s *step0026MemoryStore) ReadArtifact(_ context.Context, request domain.ArtifactReadRequest) (domain.ArtifactPage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	document, ok := s.artifacts[request.RunID][request.Artifact]
	if !ok {
		return domain.ArtifactPage{}, domain.ErrArtifactNotReady
	}
	return domain.ArtifactPage{SchemaVersion: "1", RunID: request.RunID, Artifact: request.Artifact, Entries: [][]byte{append([]byte(nil), document...)}}, nil
}

func (s *step0026MemoryStore) InterruptNonTerminal(_ context.Context) ([]domain.RecommendationRun, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var changed []domain.RecommendationRun
	for id, run := range s.runs {
		if run.Lifecycle == domain.RunQueued || run.Lifecycle == domain.RunRunning || run.Lifecycle == domain.RunCancelling {
			run.Lifecycle, run.Stage, run.Version = domain.RunInterrupted, domain.RunStageTerminal, run.Version+1
			run.RecommendationStatus, run.Failure = nil, nil
			s.runs[id] = run
			changed = append(changed, run)
		}
	}
	return changed, nil
}

type step0026Reservation struct {
	queue    *step0026Queue
	released bool
}

func (r *step0026Reservation) Commit(id domain.RunID) error {
	r.queue.mu.Lock()
	defer r.queue.mu.Unlock()
	if r.released {
		return errors.New("reservation already consumed")
	}
	r.released = true
	r.queue.reserved--
	r.queue.committed = append(r.queue.committed, id)
	return nil
}
func (r *step0026Reservation) Release() {
	r.queue.mu.Lock()
	defer r.queue.mu.Unlock()
	if !r.released {
		r.released = true
		r.queue.reserved--
	}
}

type step0026Queue struct {
	mu        sync.Mutex
	capacity  int
	reserved  int
	committed []domain.RunID
	cancelled []domain.RunID
}

func (q *step0026Queue) Reserve(context.Context) (domain.RecommendationQueueReservation, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.reserved+len(q.committed) >= q.capacity {
		return nil, domain.ErrQueueFull
	}
	q.reserved++
	return &step0026Reservation{queue: q}, nil
}
func (q *step0026Queue) Cancel(_ context.Context, id domain.RunID) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.cancelled = append(q.cancelled, id)
	return nil
}
func (*step0026Queue) Shutdown(context.Context) error { return nil }

type step0026IDGenerator struct {
	mu    sync.Mutex
	next  int
	fixed domain.RunID
}

func (g *step0026IDGenerator) NewRunID() (domain.RunID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	if g.fixed != "" {
		return g.fixed, nil
	}
	return domain.RunID(fmt.Sprintf("run_%032x", g.next)), nil
}

type step0026Executor struct {
	mu     sync.Mutex
	calls  int
	result domain.RecommendationResult
	err    error
}

func (e *step0026Executor) Recommend(ctx context.Context, request domain.RecommendationRequest) (domain.RecommendationResult, error) {
	return e.RecommendObserved(ctx, request, nil)
}
func (e *step0026Executor) RecommendObserved(_ context.Context, _ domain.RecommendationRequest, reporter domain.RecommendationProgressReporter) (domain.RecommendationResult, error) {
	e.mu.Lock()
	e.calls++
	e.mu.Unlock()
	for index, stage := range []domain.RecommendationRunStage{domain.RunStageExpand, domain.RunStageDiscover, domain.RunStageScreen, domain.RunStageCollect, domain.RunStageAssess, domain.RunStageDecide, domain.RunStageAlternatives, domain.RunStagePublish} {
		if reporter != nil {
			reporter.ReportRecommendationProgress(domain.RecommendationProgressEvent{Stage: stage, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: index + 1}})
		}
	}
	return e.result, e.err
}

func TestStep0026CanonicalRequestHashIsSetLikeButSemantic(t *testing.T) {
	left := step0026Request()
	right := step0026Request()
	right.Idea.Users = []string{"homes", "families", "families"}
	right.Idea.Capabilities = []string{"inventory", "forecast"}
	right.Policy.AllowedLicenses = []string{"MIT", "Apache-2.0", "MIT"}
	right.Search.Locales = []string{"pt-BR", "en-US"}
	normalizedLeft, leftHash, err := CanonicalRecommendationRequestHash(left)
	if err != nil {
		t.Fatal(err)
	}
	normalizedRight, rightHash, err := CanonicalRecommendationRequestHash(right)
	if err != nil {
		t.Fatal(err)
	}
	if leftHash != rightHash || !reflect.DeepEqual(normalizedLeft, normalizedRight) {
		t.Fatalf("equivalent sets differ: %s != %s\n%#v\n%#v", leftHash, rightHash, normalizedLeft, normalizedRight)
	}
	changed := right
	changed.Search.MaxCandidates++
	_, changedHash, err := CanonicalRecommendationRequestHash(changed)
	if err != nil {
		t.Fatal(err)
	}
	if changedHash == leftHash {
		t.Fatal("semantic change preserved request hash")
	}
	if got, want := normalizedLeft.Idea.Capabilities, []string{"forecast", "inventory"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("canonical capabilities = %#v, want %#v", got, want)
	}
}

func TestStep0026ValidationPrecedesStoreQueueAndIDGeneration(t *testing.T) {
	for _, key := range []string{"", "../escape", `bad\\key`, "bad/key", string([]byte{0x1f}), string(make([]byte, 129))} {
		store, queue, ids := newStep0026MemoryStore(), &step0026Queue{capacity: 1}, &step0026IDGenerator{}
		service := NewRecommendationRunService(store, queue, ids, &step0026Executor{})
		if _, err := service.Start(context.Background(), key, step0026Request()); !errors.Is(err, domain.ErrInvalidRunInput) {
			t.Fatalf("key %q: err = %v, want invalid", key, err)
		}
		if store.creates != 0 || queue.reserved != 0 || len(queue.committed) != 0 || ids.next != 0 {
			t.Fatalf("invalid key acquired capacity: store=%d reserved=%d commit=%d ids=%d", store.creates, queue.reserved, len(queue.committed), ids.next)
		}
	}
}

func TestStep0026ConcurrentStartCreatesAndCommitsExactlyOnce(t *testing.T) {
	store, queue, ids := newStep0026MemoryStore(), &step0026Queue{capacity: 64}, &step0026IDGenerator{fixed: step0026RunID}
	executor := &step0026Executor{}
	service := NewRecommendationRunService(store, queue, ids, executor)
	const callers = 32
	receipts := make(chan domain.RecommendationStartReceipt, callers)
	errorsSeen := make(chan error, callers)
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			receipt, err := service.Start(context.Background(), "same-key", step0026Request())
			receipts <- receipt
			errorsSeen <- err
		}()
	}
	wg.Wait()
	close(receipts)
	close(errorsSeen)
	for err := range errorsSeen {
		if err != nil {
			t.Fatalf("concurrent start failed: %v", err)
		}
	}
	for receipt := range receipts {
		if receipt.RunID != step0026RunID || receipt.Lifecycle != domain.RunQueued || receipt.Version != 1 {
			t.Fatalf("unstable receipt: %#v", receipt)
		}
	}
	if store.creates != 1 || len(queue.committed) != 1 || queue.reserved != 0 {
		t.Fatalf("TOCTOU/idempotency failure: creates=%d committed=%d residual=%d", store.creates, len(queue.committed), queue.reserved)
	}
	for digest := range store.idem {
		if strings.Contains(digest, "same-key") || !strings.HasPrefix(digest, "sha256:") {
			t.Fatalf("clear idempotency key persisted: %q", digest)
		}
	}
	if executor.calls != 0 {
		t.Fatalf("Start waited for or invoked executor: calls=%d", executor.calls)
	}
}

func TestStep0026ConflictAndQueueFailureDoNotMutateExistingRun(t *testing.T) {
	store, queue, ids := newStep0026MemoryStore(), &step0026Queue{capacity: 1}, &step0026IDGenerator{fixed: step0026RunID}
	service := NewRecommendationRunService(store, queue, ids, &step0026Executor{})
	if _, err := service.Start(context.Background(), "stable-key", step0026Request()); err != nil {
		t.Fatal(err)
	}
	changed := step0026Request()
	changed.Idea.Problem = "different problem"
	if _, err := service.Start(context.Background(), "stable-key", changed); !errors.Is(err, domain.ErrIdempotencyConflict) {
		t.Fatalf("different request = %v, want conflict", err)
	}
	original, err := store.Get(context.Background(), step0026RunID)
	if err != nil || original.Version != 1 || original.Request.Idea.Problem != "food waste" {
		t.Fatalf("conflict changed original: %#v err=%v", original, err)
	}
	if _, err := service.Start(context.Background(), "another-key", step0026Request()); !errors.Is(err, domain.ErrQueueFull) {
		t.Fatalf("full queue = %v", err)
	}
	if store.creates != 1 || len(store.idem) != 1 || queue.reserved != 0 || ids.next != 1 {
		t.Fatalf("queue failure left durable/reservation state: creates=%d idem=%d reserved=%d ids=%d", store.creates, len(store.idem), queue.reserved, ids.next)
	}

	quotaStore := newStep0026MemoryStore()
	quotaStore.createErr = domain.ErrQuotaExceeded
	quotaQueue := &step0026Queue{capacity: 1}
	quotaService := NewRecommendationRunService(quotaStore, quotaQueue, &step0026IDGenerator{fixed: "run_99999999999999999999999999999999"}, &step0026Executor{})
	if _, err := quotaService.Start(context.Background(), "quota-key", step0026Request()); !errors.Is(err, domain.ErrQuotaExceeded) {
		t.Fatalf("quota start = %v", err)
	}
	if quotaStore.creates != 0 || len(quotaStore.idem) != 0 || quotaQueue.reserved != 0 || len(quotaQueue.committed) != 0 {
		t.Fatalf("quota failure left durable/reservation state: store=%#v queue=%#v", quotaStore, quotaQueue)
	}
}

func TestStep0026StateMachineAllowsOnlyFrozenArcsAndMonotonicProgress(t *testing.T) {
	base := domain.RecommendationRun{RunID: step0026RunID, Lifecycle: domain.RunQueued, Version: 1, Stage: domain.RunStageQueued, Progress: domain.RecommendationRunProgress{Planned: 8}}
	allowed := []struct {
		from, to domain.RecommendationLifecycle
	}{
		{domain.RunQueued, domain.RunRunning}, {domain.RunQueued, domain.RunCancelled}, {domain.RunQueued, domain.RunInterrupted},
		{domain.RunRunning, domain.RunRunning}, {domain.RunRunning, domain.RunCancelling}, {domain.RunRunning, domain.RunSucceeded}, {domain.RunRunning, domain.RunFailed}, {domain.RunRunning, domain.RunInterrupted},
		{domain.RunCancelling, domain.RunCancelled}, {domain.RunCancelling, domain.RunInterrupted},
	}
	for _, item := range allowed {
		run := base
		run.Lifecycle = item.from
		if item.from == domain.RunRunning || item.from == domain.RunCancelling {
			run.Stage = domain.RunStageExpand
		}
		transition := domain.RunTransition{RunID: step0026RunID, ExpectedVersion: run.Version, Lifecycle: item.to, Stage: domain.RunStageTerminal, Progress: run.Progress}
		if item.to == domain.RunRunning || item.to == domain.RunCancelling {
			transition.Stage = domain.RunStageExpand
		}
		if item.from == domain.RunRunning && item.to == domain.RunRunning {
			transition.Stage = domain.RunStageDiscover
			transition.Progress.Completed = 1
		}
		if item.to == domain.RunSucceeded {
			status := domain.IdeaRunComplete
			transition.RecommendationStatus = &status
			transition.Progress.Completed = 8
			run.Artifacts = []domain.RunArtifactIndex{
				{Name: domain.ArtifactRequest, State: domain.ArtifactPublished}, {Name: domain.ArtifactRecommendation, State: domain.ArtifactPublished},
				{Name: domain.ArtifactDecision, State: domain.ArtifactPublished}, {Name: domain.ArtifactAlternatives, State: domain.ArtifactPublished},
				{Name: domain.ArtifactFailures, State: domain.ArtifactPublished}, {Name: domain.ArtifactProvenance, State: domain.ArtifactPublished},
			}
		}
		if item.to == domain.RunFailed {
			transition.Failure = &domain.RecommendationRunFailure{Code: "RECOMMENDATION_FAILED", Stage: domain.RunStageExpand, Message: "recommendation execution failed"}
		}
		if err := ValidateRecommendationRunTransition(run, transition); err != nil {
			t.Errorf("allowed %s -> %s rejected: %v", item.from, item.to, err)
		}
	}
	terminals := []domain.RecommendationLifecycle{domain.RunSucceeded, domain.RunFailed, domain.RunCancelled, domain.RunInterrupted}
	for _, from := range terminals {
		run := base
		run.Lifecycle = from
		if err := ValidateRecommendationRunTransition(run, domain.RunTransition{RunID: step0026RunID, ExpectedVersion: 1, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: run.Progress}); !errors.Is(err, domain.ErrInvalidRunTransition) {
			t.Errorf("terminal %s rewound: %v", from, err)
		}
	}
	all := []domain.RecommendationLifecycle{domain.RunQueued, domain.RunRunning, domain.RunSucceeded, domain.RunFailed, domain.RunCancelling, domain.RunCancelled, domain.RunInterrupted}
	allowedKey := map[string]bool{}
	for _, item := range allowed {
		allowedKey[string(item.from)+"->"+string(item.to)] = true
	}
	for _, from := range all {
		for _, to := range all {
			if allowedKey[string(from)+"->"+string(to)] {
				continue
			}
			run := base
			run.Lifecycle = from
			transition := domain.RunTransition{RunID: step0026RunID, ExpectedVersion: 1, Lifecycle: to, Stage: domain.RunStageTerminal, Progress: run.Progress}
			if err := ValidateRecommendationRunTransition(run, transition); !errors.Is(err, domain.ErrInvalidRunTransition) {
				t.Errorf("forbidden arc %s -> %s accepted: %v", from, to, err)
			}
		}
	}
	running := base
	running.Lifecycle, running.Stage = domain.RunRunning, domain.RunStageCollect
	running.Progress = domain.RecommendationRunProgress{Planned: 8, Completed: 4}
	for _, bad := range []domain.RunTransition{
		{RunID: step0026RunID, ExpectedVersion: 1, Lifecycle: domain.RunRunning, Stage: domain.RunStageScreen, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 3}},
		{RunID: step0026RunID, ExpectedVersion: 1, Lifecycle: domain.RunRunning, Stage: "provider said TOKEN_CANARY", Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 5}},
		{RunID: step0026RunID, ExpectedVersion: 1, Lifecycle: domain.RunRunning, Stage: domain.RunStageAssess, Progress: domain.RecommendationRunProgress{Planned: 9, Completed: 5}},
	} {
		if err := ValidateRecommendationRunTransition(running, bad); !errors.Is(err, domain.ErrInvalidRunTransition) {
			t.Errorf("bad progress accepted: %#v err=%v", bad, err)
		}
	}
}

func TestStep0026ExecutePublishesBeforeSucceededAndPreservesInconclusive(t *testing.T) {
	store, queue, ids := newStep0026MemoryStore(), &step0026Queue{capacity: 2}, &step0026IDGenerator{fixed: step0026RunID}
	result := domain.RecommendationResult{
		SchemaVersion: "1", Status: domain.IdeaRunInconclusive,
		Decision:     domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown, Matrix: []domain.CapabilityComparison{}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}},
		Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Alternatives: []domain.SoftwareAlternative{}, Excluded: []domain.SoftwareAlternative{}},
		Failures:     []domain.ExecutionFailure{{Scope: "provider", Stage: "discover", Message: "Authorization: Bearer TOKEN_CANARY /private/path README_CANARY"}},
		Warnings:     []string{"license remains unknown", "CODE_CANARY Cookie=SECRET"},
	}
	service := NewRecommendationRunService(store, queue, ids, &step0026Executor{result: result})
	if _, err := service.Start(context.Background(), "execute-key", step0026Request()); err != nil {
		t.Fatal(err)
	}
	if err := service.Execute(context.Background(), step0026RunID); err != nil {
		t.Fatal(err)
	}
	run, err := store.Get(context.Background(), step0026RunID)
	if err != nil {
		t.Fatal(err)
	}
	if run.Lifecycle != domain.RunSucceeded || run.Stage != domain.RunStageTerminal || run.RecommendationStatus == nil || *run.RecommendationStatus != domain.IdeaRunInconclusive {
		t.Fatalf("lifecycle/status conflated: %#v", run)
	}
	wantNames := []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance}
	gotNames := make([]domain.RunArtifactName, 0, len(run.Artifacts))
	for _, artifact := range run.Artifacts {
		if artifact.State != domain.ArtifactPublished {
			t.Fatalf("terminal exposed non-published artifact: %#v", artifact)
		}
		gotNames = append(gotNames, artifact.Name)
	}
	sort.Slice(gotNames, func(i, j int) bool { return gotNames[i] < gotNames[j] })
	sort.Slice(wantNames, func(i, j int) bool { return wantNames[i] < wantNames[j] })
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("terminal artifact matrix = %#v, want %#v", gotNames, wantNames)
	}
	for _, document := range store.artifacts[step0026RunID] {
		for _, canary := range []string{"TOKEN_CANARY", "/private/path", "README_CANARY", "CODE_CANARY", "SECRET"} {
			if strings.Contains(string(document), canary) {
				t.Fatalf("canary %q persisted in public artifact: %s", canary, document)
			}
		}
	}
}

func TestStep0026CancelIsIdempotentAndRecoveryInterruptsOnlyOrphans(t *testing.T) {
	store, queue, ids := newStep0026MemoryStore(), &step0026Queue{capacity: 4}, &step0026IDGenerator{fixed: step0026RunID}
	service := NewRecommendationRunService(store, queue, ids, &step0026Executor{})
	if _, err := service.Start(context.Background(), "cancel-key", step0026Request()); err != nil {
		t.Fatal(err)
	}
	first, err := service.Cancel(context.Background(), step0026RunID)
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Cancel(context.Background(), step0026RunID)
	if err != nil {
		t.Fatal(err)
	}
	if first.Lifecycle != domain.RunCancelled || second.Lifecycle != domain.RunCancelled || first.Version != second.Version || len(queue.cancelled) != 1 {
		t.Fatalf("queued cancel not idempotent: first=%#v second=%#v signals=%d", first, second, len(queue.cancelled))
	}

	orphanID := domain.RunID("run_11111111111111111111111111111111")
	terminalID := domain.RunID("run_22222222222222222222222222222222")
	store.runs[orphanID] = domain.RecommendationRun{SchemaVersion: "1", RunID: orphanID, Lifecycle: domain.RunRunning, Version: 4, Stage: domain.RunStageCollect, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 3}, Artifacts: []domain.RunArtifactIndex{{Name: domain.ArtifactRequest, State: domain.ArtifactPublished}}}
	status := domain.IdeaRunPartial
	store.runs[terminalID] = domain.RecommendationRun{SchemaVersion: "1", RunID: terminalID, Lifecycle: domain.RunSucceeded, Version: 9, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 8}, RecommendationStatus: &status, Artifacts: []domain.RunArtifactIndex{}}
	before := store.runs[terminalID]
	if err := service.Recover(context.Background()); err != nil {
		t.Fatal(err)
	}
	orphan := store.runs[orphanID]
	if orphan.Lifecycle != domain.RunInterrupted || orphan.Progress.Completed != 3 || orphan.RecommendationStatus != nil || len(orphan.Artifacts) != 1 {
		t.Fatalf("orphan recovery invented/lost state: %#v", orphan)
	}
	if after := store.runs[terminalID]; !reflect.DeepEqual(after, before) {
		t.Fatalf("recovery rewrote terminal:\nbefore %#v\nafter  %#v", before, after)
	}
}

func TestStep0026CompletionCancellationRaceHasOneImmutableWinner(t *testing.T) {
	store := newStep0026MemoryStore()
	status := domain.IdeaRunPartial
	run, record := step0026StoredRunForRace()
	store.runs[run.RunID], store.idem[record.KeyDigest] = run, record
	transitions := []domain.RunTransition{
		{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunSucceeded, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 8}, RecommendationStatus: &status},
		{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunCancelling, Stage: domain.RunStagePublish, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 7}},
	}
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for _, transition := range transitions {
		transition := transition
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Transition(context.Background(), transition)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	successes, conflicts := 0, 0
	for err := range errs {
		if err == nil {
			successes++
		} else if errors.Is(err, domain.ErrRunVersionConflict) {
			conflicts++
		} else {
			t.Fatalf("unexpected race error: %v", err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("CAS race winners=%d conflicts=%d", successes, conflicts)
	}
	winner, _ := store.Get(context.Background(), run.RunID)
	if winner.Version != run.Version+1 {
		t.Fatalf("race mutated version more than once: %#v", winner)
	}
	if winner.Lifecycle == domain.RunSucceeded {
		if _, err := store.Transition(context.Background(), domain.RunTransition{RunID: run.RunID, ExpectedVersion: winner.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: winner.Progress}); !errors.Is(err, domain.ErrInvalidRunTransition) {
			t.Fatalf("terminal winner was overwritten: %v", err)
		}
	}
}

func step0026StoredRunForRace() (domain.RecommendationRun, domain.IdempotencyRecord) {
	id := domain.RunID("run_33333333333333333333333333333333")
	run := domain.RecommendationRun{
		SchemaVersion: "1", RunID: id, Request: step0026Request(), RequestHash: "sha256:" + strings.Repeat("d", 64),
		Lifecycle: domain.RunRunning, Version: 8, Stage: domain.RunStagePublish,
		Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 7},
		Artifacts: []domain.RunArtifactIndex{
			{Name: domain.ArtifactRequest, State: domain.ArtifactPublished}, {Name: domain.ArtifactRecommendation, State: domain.ArtifactPublished},
			{Name: domain.ArtifactDecision, State: domain.ArtifactPublished}, {Name: domain.ArtifactAlternatives, State: domain.ArtifactPublished},
			{Name: domain.ArtifactFailures, State: domain.ArtifactPublished}, {Name: domain.ArtifactProvenance, State: domain.ArtifactPublished},
		},
	}
	return run, domain.IdempotencyRecord{KeyDigest: "sha256:" + strings.Repeat("e", 64), RequestHash: run.RequestHash, RunID: id}
}
