/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runcontainer/container_test.md
 * @layer L3
 * @updated 2026-09-19
 */

package runcontainer_test

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/03_infra/runcontainer"
	"bateia/03_infra/runstore"
)

var step0079Limits = runstore.RunStoreLimits{
	MaxArtifactBytes: 1 << 20,
	MaxRunBytes:      4 << 20,
	MaxTotalBytes:    32 << 20,
}

type step0079ExecutorProbe struct {
	mu       sync.Mutex
	calls    int
	requests []domain.RecommendationRequest
	entered  chan struct{}
	release  chan struct{}
}

func (p *step0079ExecutorProbe) Recommend(context.Context, domain.RecommendationRequest) (domain.RecommendationResult, error) {
	return step0079RecommendationResult(), nil
}

func (p *step0079ExecutorProbe) RecommendObserved(ctx context.Context, request domain.RecommendationRequest, reporter domain.RecommendationProgressReporter) (domain.RecommendationResult, error) {
	p.mu.Lock()
	p.calls++
	p.requests = append(p.requests, request)
	p.mu.Unlock()
	if p.entered != nil {
		select {
		case p.entered <- struct{}{}:
		default:
		}
	}
	if p.release != nil {
		<-p.release
	}
	for index, stage := range []domain.RecommendationRunStage{
		domain.RunStageExpand,
		domain.RunStageDiscover,
		domain.RunStageScreen,
		domain.RunStageCollect,
		domain.RunStageAssess,
		domain.RunStageDecide,
		domain.RunStageAlternatives,
		domain.RunStagePublish,
	} {
		reporter.ReportRecommendationProgress(domain.RecommendationProgressEvent{
			Stage:    stage,
			Progress: domain.RecommendationRunProgress{Planned: 8, Completed: index + 1},
		})
	}
	return step0079RecommendationResult(), nil
}

func step0079RecommendationResult() domain.RecommendationResult {
	return domain.RecommendationResult{
		SchemaVersion: "1",
		Status:        domain.IdeaRunInconclusive,
		Decision: domain.ReuseDecisionReport{
			SchemaVersion: "1",
			Decision:      domain.DecisionUnknown,
			Matrix:        []domain.CapabilityComparison{},
			Gaps:          []string{},
			Reasons:       []domain.DecisionReason{},
			Alternatives:  []domain.DecisionAlternative{},
		},
		Alternatives: domain.SoftwareAlternativeReport{
			SchemaVersion: "1",
			Alternatives:  []domain.SoftwareAlternative{},
			Excluded:      []domain.SoftwareAlternative{},
		},
		Failures: []domain.ExecutionFailure{},
		Warnings: []string{},
	}
}

func step0079RecommendationRequest(problem string) domain.RecommendationRequest {
	return domain.RecommendationRequest{
		Idea: domain.IdeaSpec{Problem: problem, Capabilities: []string{"inventory"}},
		Policy: domain.ReuseDecisionPolicy{
			SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2,
			ExtendMinCoverage: 0.7, IntegrationMinValue: 0.3,
		},
		Target: domain.SoftwareAlternativeTarget{Name: "pantry"},
		Search: domain.RecommendationSearch{
			Strategy: "literal", Locales: []string{"en-US"}, MaxPerQuery: 1, MaxCandidates: 1,
		},
	}
}

func step0079NewContainer(t *testing.T, root string, executor domain.ObservableRecommendationExecutor) *runcontainer.Container {
	t.Helper()
	container, err := runcontainer.New(
		executor,
		root,
		[]byte("0123456789abcdef0123456789abcdef"),
		step0079Limits,
		1,
		8,
	)
	if err != nil {
		t.Fatalf("new runtime container: %v", err)
	}
	return container
}

func TestStep0079ContainerAPIValidationAndStableService(t *testing.T) {
	executor := &step0079ExecutorProbe{}
	if container, err := runcontainer.New(executor, filepath.Join(t.TempDir(), "bad-secret"), []byte("short"), step0079Limits, 1, 8); err == nil || container != nil {
		t.Fatalf("invalid secret produced container=%p err=%v", container, err)
	}

	queueFailureRoot := filepath.Join(t.TempDir(), "queue-error")
	if container, err := runcontainer.New(executor, queueFailureRoot, []byte("0123456789abcdef0123456789abcdef"), step0079Limits, 0, 8); err == nil || container != nil {
		t.Fatalf("invalid queue configuration produced container=%p err=%v", container, err)
	}
	// A queue-construction failure happens after the concrete store exists. The
	// publicly observable guarantee is that its directory remains reusable. A
	// concrete Close failure cannot be produced by FileStore and remains Unknown.
	reopened, err := runstore.NewFileStore(queueFailureRoot, []byte("0123456789abcdef0123456789abcdef"), step0079Limits)
	if err != nil {
		t.Fatalf("store leaked after queue construction failure: %v", err)
	}
	if err := reopened.Close(); err != nil {
		t.Fatal(err)
	}

	container := step0079NewContainer(t, filepath.Join(t.TempDir(), "valid"), executor)
	if container.Service() == nil || container.Service() != container.Service() {
		t.Fatal("Service did not expose one stable service instance")
	}
	if err := container.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown before start: %v", err)
	}
}

func TestStep0079ContainerRecoverBeforeStartExecutionAndReopen(t *testing.T) {
	root := t.TempDir()
	executor := &step0079ExecutorProbe{}
	container := step0079NewContainer(t, root, executor)
	service := container.Service()

	queued, err := service.Start(context.Background(), "queued-before-recovery", step0079RecommendationRequest("must be interrupted"))
	if err != nil {
		t.Fatalf("queue before recovery: %v", err)
	}
	if err := container.Recover(context.Background()); err != nil {
		t.Fatalf("recover: %v", err)
	}
	interrupted, err := service.GetRun(context.Background(), queued.RunID)
	if err != nil || interrupted.Lifecycle != domain.RunInterrupted {
		t.Fatalf("pre-start recovery = %+v, %v; want INTERRUPTED", interrupted, err)
	}
	if err := container.Start(); err != nil {
		t.Fatalf("start queue: %v", err)
	}

	started, err := service.Start(context.Background(), "executes-through-container", step0079RecommendationRequest("execute me"))
	if err != nil {
		t.Fatalf("start durable run: %v", err)
	}
	view := step0079WaitForLifecycle(t, service.GetRun, started.RunID, domain.RunSucceeded)
	if view.Lifecycle != domain.RunSucceeded {
		t.Fatalf("run lifecycle = %s, want SUCCEEDED", view.Lifecycle)
	}
	executor.mu.Lock()
	calls := executor.calls
	requests := append([]domain.RecommendationRequest(nil), executor.requests...)
	executor.mu.Unlock()
	if calls != 1 || len(requests) != 1 || requests[0].Idea.Problem != "execute me" {
		t.Fatalf("queue did not execute through the injected service/executor: calls=%d requests=%+v", calls, requests)
	}
	if err := container.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	reopened := step0079NewContainer(t, root, &step0079ExecutorProbe{})
	if err := reopened.Recover(context.Background()); err != nil {
		t.Fatalf("reopen recovery: %v", err)
	}
	if err := reopened.Start(); err != nil {
		t.Fatalf("reopen start: %v", err)
	}
	persisted, err := reopened.Service().GetRun(context.Background(), started.RunID)
	if err != nil || persisted.Lifecycle != domain.RunSucceeded {
		t.Fatalf("reopened run = %+v, %v", persisted, err)
	}
	if err := reopened.Shutdown(context.Background()); err != nil {
		t.Fatalf("reopened shutdown: %v", err)
	}
}

func TestStep0079ContainerStartFailureClosesStore(t *testing.T) {
	container := step0079NewContainer(t, t.TempDir(), &step0079ExecutorProbe{})
	if err := container.Start(); err != nil {
		t.Fatalf("first start: %v", err)
	}
	if err := container.Start(); err == nil {
		t.Fatal("second start unexpectedly succeeded")
	}
	if err := container.Recover(context.Background()); !errors.Is(err, domain.ErrRunStoreUnavailable) {
		t.Fatalf("store remained usable after Start failure: %v", err)
	}
	if err := container.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown after failed start: %v", err)
	}
}

func TestStep0079ContainerShutdownQueueErrorWinsAndStoreStillCloses(t *testing.T) {
	executor := &step0079ExecutorProbe{entered: make(chan struct{}, 1), release: make(chan struct{})}
	container := step0079NewContainer(t, t.TempDir(), executor)
	if err := container.Recover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := container.Start(); err != nil {
		t.Fatal(err)
	}
	if _, err := container.Service().Start(context.Background(), "blocked-shutdown", step0079RecommendationRequest("block shutdown")); err != nil {
		t.Fatal(err)
	}
	select {
	case <-executor.entered:
	case <-time.After(time.Second):
		t.Fatal("worker did not enter injected executor")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := container.Shutdown(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("shutdown error = %v, want queue context cancellation", err)
	}
	if err := container.Recover(context.Background()); !errors.Is(err, domain.ErrRunStoreUnavailable) {
		t.Fatalf("store was not closed after queue shutdown error: %v", err)
	}
	close(executor.release)
	if err := container.Shutdown(context.Background()); err != nil {
		t.Fatalf("final worker drain: %v", err)
	}
}

func step0079WaitForLifecycle(
	t *testing.T,
	get func(context.Context, domain.RunID) (domain.RecommendationRunView, error),
	id domain.RunID,
	want domain.RecommendationLifecycle,
) domain.RecommendationRunView {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	var last domain.RecommendationRunView
	for time.Now().Before(deadline) {
		view, err := get(context.Background(), id)
		if err != nil {
			t.Fatalf("get run: %v", err)
		}
		last = view
		if view.Lifecycle == want {
			return view
		}
		time.Sleep(time.Millisecond)
	}
	return last
}
