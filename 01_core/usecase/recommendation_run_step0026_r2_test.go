/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_run_step0026_r2_test.md
 * @layer L1
 * @updated 2026-09-01
 */
package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"bateia/01_core/domain"
)

func step0026SetVariants(request domain.RecommendationRequest) []struct {
	name   string
	mutate func(*domain.RecommendationRequest)
} {
	return []struct {
		name   string
		mutate func(*domain.RecommendationRequest)
	}{
		{"users", func(r *domain.RecommendationRequest) { r.Idea.Users = []string{"homes", "families", "homes"} }},
		{"capabilities", func(r *domain.RecommendationRequest) {
			r.Idea.Capabilities = []string{"inventory", "forecast", "inventory"}
		}},
		{"constraints", func(r *domain.RecommendationRequest) { r.Idea.Constraints = []string{"private", "private"} }},
		{"non-goals", func(r *domain.RecommendationRequest) { r.Idea.NonGoals = []string{"ads", "ads"} }},
		{"preferred-stack", func(r *domain.RecommendationRequest) { r.Idea.PreferredStack = []string{"Go", "Go"} }},
		{"licenses", func(r *domain.RecommendationRequest) { r.Policy.AllowedLicenses = []string{"MIT", "Apache-2.0", "MIT"} }},
		{"locales", func(r *domain.RecommendationRequest) { r.Search.Locales = []string{"pt-BR", "en-US", "pt-BR"} }},
	}
}

func TestStep0026R2EverySetFieldCanonicalizesAndRetriesOneRun(t *testing.T) {
	base := step0026Request()
	_, baseHash, err := CanonicalRecommendationRequestHash(base)
	if err != nil {
		t.Fatal(err)
	}
	for _, variant := range step0026SetVariants(base) {
		t.Run(variant.name, func(t *testing.T) {
			changed := base
			variant.mutate(&changed)
			_, hash, err := CanonicalRecommendationRequestHash(changed)
			if err != nil {
				t.Fatal(err)
			}
			if hash != baseHash {
				t.Fatalf("set-like reordering/duplication changed hash: %s != %s", hash, baseHash)
			}
		})
	}
	combined := base
	for _, variant := range step0026SetVariants(base) {
		variant.mutate(&combined)
	}
	store, queue := newStep0026MemoryStore(), &step0026Queue{capacity: 8}
	service := NewRecommendationRunService(store, queue, &step0026IDGenerator{}, &step0026Executor{})
	barrier := make(chan struct{})
	receipts := make(chan domain.RecommendationStartReceipt, 2)
	errs := make(chan error, 2)
	for _, request := range []domain.RecommendationRequest{base, combined} {
		request := request
		go func() {
			<-barrier
			receipt, err := service.Start(context.Background(), "all-sets", request)
			receipts <- receipt
			errs <- err
		}()
	}
	close(barrier)
	first, second := <-receipts, <-receipts
	for range 2 {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
	if !reflect.DeepEqual(first, second) || store.creates != 1 || len(queue.committed) != 1 {
		t.Fatalf("equivalent aggregate request duplicated run: first=%#v second=%#v creates=%d commits=%d", first, second, store.creates, len(queue.committed))
	}
}

func TestStep0026R2EverySemanticFieldParticipatesInRequestHashButKeyDoesNot(t *testing.T) {
	base := step0026Request()
	_, baseline, err := CanonicalRecommendationRequestHash(base)
	if err != nil {
		t.Fatal(err)
	}
	mutations := []struct {
		name   string
		mutate func(*domain.RecommendationRequest)
	}{
		{"problem", func(r *domain.RecommendationRequest) { r.Idea.Problem = "another" }},
		{"users", func(r *domain.RecommendationRequest) { r.Idea.Users = append(r.Idea.Users, "schools") }},
		{"capabilities", func(r *domain.RecommendationRequest) { r.Idea.Capabilities = append(r.Idea.Capabilities, "alerts") }},
		{"constraints", func(r *domain.RecommendationRequest) { r.Idea.Constraints = append(r.Idea.Constraints, "offline") }},
		{"non-goals", func(r *domain.RecommendationRequest) { r.Idea.NonGoals = append(r.Idea.NonGoals, "tracking") }},
		{"preferred-stack", func(r *domain.RecommendationRequest) { r.Idea.PreferredStack = append(r.Idea.PreferredStack, "Rust") }},
		{"licenses", func(r *domain.RecommendationRequest) {
			r.Policy.AllowedLicenses = append(r.Policy.AllowedLicenses, "BSD-3-Clause")
		}},
		{"require-local", func(r *domain.RecommendationRequest) { r.Policy.RequireLocal = !r.Policy.RequireLocal }},
		{"allow-archived", func(r *domain.RecommendationRequest) { r.Policy.AllowArchived = !r.Policy.AllowArchived }},
		{"max-components", func(r *domain.RecommendationRequest) { r.Policy.MaxComponents++ }},
		{"extend-threshold", func(r *domain.RecommendationRequest) { r.Policy.ExtendMinCoverage = .8 }},
		{"integration-threshold", func(r *domain.RecommendationRequest) { r.Policy.IntegrationMinValue = .2 }},
		{"target", func(r *domain.RecommendationRequest) { r.Target.Name = "another target" }},
		{"strategy", func(r *domain.RecommendationRequest) { r.Search.Strategy = "lexical" }},
		{"language", func(r *domain.RecommendationRequest) { r.Search.Language = "go" }},
		{"locales", func(r *domain.RecommendationRequest) { r.Search.Locales = append(r.Search.Locales, "ja-JP") }},
		{"max-per-query", func(r *domain.RecommendationRequest) { r.Search.MaxPerQuery++ }},
		{"max-candidates", func(r *domain.RecommendationRequest) { r.Search.MaxCandidates++ }},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			changed := base
			mutation.mutate(&changed)
			_, got, err := CanonicalRecommendationRequestHash(changed)
			if err != nil {
				t.Fatal(err)
			}
			if got == baseline {
				t.Fatalf("semantic field %s omitted from hash", mutation.name)
			}
		})
	}
	store, queue := newStep0026MemoryStore(), &step0026Queue{capacity: 4}
	service := NewRecommendationRunService(store, queue, &step0026IDGenerator{}, &step0026Executor{})
	one, err := service.Start(context.Background(), "key-one", base)
	if err != nil {
		t.Fatal(err)
	}
	two, err := service.Start(context.Background(), "key-two", base)
	if err != nil {
		t.Fatal(err)
	}
	if one.RunID == two.RunID || len(store.runs) != 2 {
		t.Fatalf("raw idempotency key leaked into request hash or identity collapsed: %#v %#v", one, two)
	}
	for _, run := range store.runs {
		if run.RequestHash != baseline {
			t.Fatalf("same request under another key changed hash: %#v", run)
		}
	}
}

func TestStep0026R2RetryReceiptTracksCurrentLifecycleAndDoesNotMutate(t *testing.T) {
	store, queue := newStep0026MemoryStore(), &step0026Queue{capacity: 2}
	service := NewRecommendationRunService(store, queue, &step0026IDGenerator{fixed: step0026RunID}, &step0026Executor{})
	queued, err := service.Start(context.Background(), "retry-key", step0026Request())
	if err != nil {
		t.Fatal(err)
	}
	run, _ := store.Get(context.Background(), queued.RunID)
	run, err = store.Transition(context.Background(), domain.RunTransition{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: domain.RecommendationRunProgress{Planned: 8}})
	if err != nil {
		t.Fatal(err)
	}
	running, err := service.Start(context.Background(), "retry-key", step0026Request())
	if err != nil {
		t.Fatal(err)
	}
	if running.Lifecycle != domain.RunRunning || running.Version != run.Version {
		t.Fatalf("retry returned stale queued receipt: %#v current=%#v", running, run)
	}
	status := domain.IdeaRunPartial
	store.mu.Lock()
	run.Lifecycle, run.Stage, run.Version, run.RecommendationStatus = domain.RunSucceeded, domain.RunStageTerminal, 19, &status
	store.runs[run.RunID] = run
	store.mu.Unlock()
	succeeded, err := service.Start(context.Background(), "retry-key", step0026Request())
	if err != nil {
		t.Fatal(err)
	}
	if succeeded.Lifecycle != domain.RunSucceeded || succeeded.Version != 19 || len(queue.committed) != 1 || store.creates != 1 {
		t.Fatalf("terminal retry mutated/staled receipt: %#v commits=%d creates=%d", succeeded, len(queue.committed), store.creates)
	}
}

func TestStep0026R2SucceededPreservesAllThreeRecommendationStatuses(t *testing.T) {
	for index, status := range []domain.IdeaRunStatus{domain.IdeaRunComplete, domain.IdeaRunPartial, domain.IdeaRunInconclusive} {
		t.Run(string(status), func(t *testing.T) {
			id := domain.RunID("run_" + strings.Repeat(string("abc"[index]), 32))
			decision := domain.DecisionBuild
			failures := []domain.ExecutionFailure{}
			if status == domain.IdeaRunInconclusive {
				decision = domain.DecisionUnknown
			}
			if status == domain.IdeaRunPartial {
				failures = []domain.ExecutionFailure{{Scope: "run", Stage: "collect", Message: "collection failed"}}
			}
			result := domain.RecommendationResult{SchemaVersion: "1", Status: status, Decision: domain.ReuseDecisionReport{SchemaVersion: "1", Decision: decision, Matrix: []domain.CapabilityComparison{}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}}, Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Alternatives: []domain.SoftwareAlternative{}, Excluded: []domain.SoftwareAlternative{}}, Failures: failures, Warnings: []string{}}
			store, queue := newStep0026MemoryStore(), &step0026Queue{capacity: 1}
			service := NewRecommendationRunService(store, queue, &step0026IDGenerator{fixed: id}, &step0026Executor{result: result})
			if _, err := service.Start(context.Background(), "status-"+string(status), step0026Request()); err != nil {
				t.Fatal(err)
			}
			if err := service.Execute(context.Background(), id); err != nil {
				t.Fatal(err)
			}
			run, _ := store.Get(context.Background(), id)
			if run.Lifecycle != domain.RunSucceeded || run.RecommendationStatus == nil || *run.RecommendationStatus != status {
				t.Fatalf("status %s projected as %#v", status, run)
			}
		})
	}
}

func TestStep0026R2ValidationMatrixTouchesNoPort(t *testing.T) {
	base := step0026Request()
	invalid := []struct {
		name   string
		key    string
		mutate func(*domain.RecommendationRequest)
	}{
		{"empty-key", "", func(*domain.RecommendationRequest) {}},
		{"non-ascii-key", "olá", func(*domain.RecommendationRequest) {}},
		{"slash-key", "a/b", func(*domain.RecommendationRequest) {}},
		{"backslash-key", `a\b`, func(*domain.RecommendationRequest) {}},
		{"dotdot-key", "a..b", func(*domain.RecommendationRequest) {}},
		{"long-key", strings.Repeat("a", 129), func(*domain.RecommendationRequest) {}},
		{"problem", "key", func(r *domain.RecommendationRequest) { r.Idea.Problem = "" }},
		{"capabilities", "key", func(r *domain.RecommendationRequest) { r.Idea.Capabilities = nil }},
		{"policy-version", "key", func(r *domain.RecommendationRequest) { r.Policy.SchemaVersion = "2" }},
		{"max-components", "key", func(r *domain.RecommendationRequest) { r.Policy.MaxComponents = 1 }},
		{"extend-zero", "key", func(r *domain.RecommendationRequest) { r.Policy.ExtendMinCoverage = 0 }},
		{"extend-high", "key", func(r *domain.RecommendationRequest) { r.Policy.ExtendMinCoverage = 1.1 }},
		{"integration-zero", "key", func(r *domain.RecommendationRequest) { r.Policy.IntegrationMinValue = 0 }},
		{"integration-above-extend", "key", func(r *domain.RecommendationRequest) { r.Policy.IntegrationMinValue = .9 }},
		{"strategy", "key", func(r *domain.RecommendationRequest) { r.Search.Strategy = "opaque" }},
		{"duplicate-locale", "key", func(r *domain.RecommendationRequest) { r.Search.Locales = []string{"en-US", "en-US"} }},
		{"unknown-locale", "key", func(r *domain.RecommendationRequest) { r.Search.Locales = []string{"pt-PT"} }},
		{"max-per-query-low", "key", func(r *domain.RecommendationRequest) { r.Search.MaxPerQuery = 0 }},
		{"max-per-query-high", "key", func(r *domain.RecommendationRequest) { r.Search.MaxPerQuery = 26 }},
		{"max-candidates-low", "key", func(r *domain.RecommendationRequest) { r.Search.MaxCandidates = 0 }},
		{"max-candidates-high", "key", func(r *domain.RecommendationRequest) { r.Search.MaxCandidates = 21 }},
	}
	for _, test := range invalid {
		t.Run(test.name, func(t *testing.T) {
			request := base
			test.mutate(&request)
			store, queue, ids, executor := newStep0026MemoryStore(), &step0026Queue{capacity: 1}, &step0026IDGenerator{}, &step0026Executor{}
			service := NewRecommendationRunService(store, queue, ids, executor)
			if _, err := service.Start(context.Background(), test.key, request); !errors.Is(err, domain.ErrInvalidRunInput) {
				t.Fatalf("invalid %s = %v", test.name, err)
			}
			if store.creates != 0 || len(store.idem) != 0 || queue.reserved != 0 || len(queue.committed) != 0 || ids.next != 0 || executor.calls != 0 {
				t.Fatalf("invalid input acquired effects: store=%#v queue=%#v ids=%d executor=%d", store, queue, ids.next, executor.calls)
			}
		})
	}
}

func TestStep0026R2QueueFullDoesNotReserveIdentityAndRetryCanCreate(t *testing.T) {
	store, queue, ids, executor := newStep0026MemoryStore(), &step0026Queue{capacity: 1}, &step0026IDGenerator{}, &step0026Executor{}
	held, err := queue.Reserve(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	service := NewRecommendationRunService(store, queue, ids, executor)
	if _, err := service.Start(context.Background(), "queue-retry-key", step0026Request()); !errors.Is(err, domain.ErrQueueFull) {
		t.Fatalf("full queue = %v", err)
	}
	if ids.next != 0 || store.creates != 0 || len(store.idem) != 0 || executor.calls != 0 {
		t.Fatalf("full queue acquired identity/store/provider: ids=%d creates=%d idem=%d calls=%d", ids.next, store.creates, len(store.idem), executor.calls)
	}
	held.Release()
	receipt, err := service.Start(context.Background(), "queue-retry-key", step0026Request())
	if err != nil || receipt.Lifecycle != domain.RunQueued || store.creates != 1 || len(store.idem) != 1 {
		t.Fatalf("retry after capacity release failed: %#v err=%v creates=%d idem=%d", receipt, err, store.creates, len(store.idem))
	}
}

func TestStep0026R2TerminalOperationsAndViewsAreImmutableNoOps(t *testing.T) {
	lifecycles := []domain.RecommendationLifecycle{domain.RunSucceeded, domain.RunFailed, domain.RunCancelled, domain.RunInterrupted}
	for index, lifecycle := range lifecycles {
		t.Run(string(lifecycle), func(t *testing.T) {
			store, queue := newStep0026MemoryStore(), &step0026Queue{capacity: 2}
			id := domain.RunID("run_" + strings.Repeat(string("abcdef01"[index]), 32))
			published := func(name domain.RunArtifactName, version uint64) domain.RunArtifactIndex {
				return domain.RunArtifactIndex{Name: name, SchemaVersion: "1", ArtifactHash: "sha256:" + strings.Repeat(string("123456"[version%6]), 64), SizeBytes: 64, EntryCount: 1, Version: version, State: domain.ArtifactPublished}
			}
			run := domain.RecommendationRun{SchemaVersion: "1", RunID: id, Request: step0026Request(), RequestHash: "sha256:" + strings.Repeat("a", 64), Lifecycle: lifecycle, Version: 11, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 3}, Artifacts: []domain.RunArtifactIndex{published(domain.ArtifactRequest, 2)}}
			if lifecycle == domain.RunSucceeded {
				status := domain.IdeaRunInconclusive
				run.RecommendationStatus = &status
				run.Progress.Completed = 8
				run.Artifacts = []domain.RunArtifactIndex{
					published(domain.ArtifactRequest, 2),
					published(domain.ArtifactRecommendation, 3),
					published(domain.ArtifactDecision, 4),
					published(domain.ArtifactAlternatives, 5),
					published(domain.ArtifactFailures, 6),
					published(domain.ArtifactProvenance, 7),
				}
			}
			if lifecycle == domain.RunFailed {
				run.Failure = &domain.RecommendationRunFailure{Code: "RECOMMENDATION_FAILED", Stage: domain.RunStageCollect, Message: "recommendation execution failed"}
				run.Artifacts = append(run.Artifacts, published(domain.ArtifactFailures, 3))
			}
			store.runs[id] = run
			service := NewRecommendationRunService(store, queue, &step0026IDGenerator{}, &step0026Executor{})
			before, _ := json.Marshal(run)
			view, err := service.GetRun(context.Background(), id)
			if err != nil {
				t.Fatal(err)
			}
			if len(view.Artifacts) != len(run.Artifacts) || view.Artifacts == nil {
				t.Fatalf("view leaked PENDING or null array: %#v", view)
			}
			for _, artifact := range view.Artifacts {
				if artifact.State != domain.ArtifactPublished {
					t.Fatalf("terminal view exposed non-published artifact: %#v", view)
				}
			}
			cancelled, err := service.Cancel(context.Background(), id)
			if err != nil || cancelled.Version != 11 || cancelled.Lifecycle != lifecycle {
				t.Fatalf("terminal cancel changed run: %#v err=%v", cancelled, err)
			}
			if err := service.Execute(context.Background(), id); err != nil {
				t.Fatalf("terminal execute should be no-op: %v", err)
			}
			if err := service.Recover(context.Background()); err != nil {
				t.Fatal(err)
			}
			after, _ := store.Get(context.Background(), id)
			afterBytes, _ := json.Marshal(after)
			if string(before) != string(afterBytes) || len(queue.cancelled) != 0 {
				t.Fatalf("terminal operation mutated bytes/signalled queue:\n%s\n%s signals=%d", before, afterBytes, len(queue.cancelled))
			}
		})
	}
}

func TestStep0026R2ProgressRejectsEveryInvalidVariantAndAcceptsMonotonicStages(t *testing.T) {
	stages := []domain.RecommendationRunStage{domain.RunStageExpand, domain.RunStageDiscover, domain.RunStageScreen, domain.RunStageCollect, domain.RunStageAssess, domain.RunStageDecide, domain.RunStageAlternatives, domain.RunStagePublish}
	run := domain.RecommendationRun{RunID: step0026RunID, Lifecycle: domain.RunRunning, Version: 2, Stage: domain.RunStageExpand, Progress: domain.RecommendationRunProgress{Planned: 8}}
	for index := 1; index < len(stages); index++ {
		transition := domain.RunTransition{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: stages[index], Progress: domain.RecommendationRunProgress{Planned: 8, Completed: index}}
		if err := ValidateRecommendationRunTransition(run, transition); err != nil {
			t.Fatalf("monotonic stage %d rejected: %v", index, err)
		}
		run.Stage, run.Progress, run.Version = transition.Stage, transition.Progress, run.Version+1
	}
	bad := []domain.RunTransition{
		{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: run.Stage, Progress: run.Progress},
		{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStageCollect, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 2}},
		{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStagePublish, Progress: domain.RecommendationRunProgress{Planned: 7, Completed: 8}},
		{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStagePublish, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 9}},
		{RunID: run.RunID, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: "TOKEN provider said ready", Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 8}},
	}
	for _, transition := range bad {
		if err := ValidateRecommendationRunTransition(run, transition); !errors.Is(err, domain.ErrInvalidRunTransition) {
			t.Errorf("invalid progress accepted: %#v err=%v", transition, err)
		}
	}
}

func TestStep0026R2SixteenConcurrentCancelsProduceOneSignalAndVersion(t *testing.T) {
	store, queue := newStep0026MemoryStore(), &step0026Queue{capacity: 2}
	run, record := step0026StoredRunForRace()
	run.Lifecycle, run.Stage, run.Version = domain.RunRunning, domain.RunStageCollect, 4
	run.Progress = domain.RecommendationRunProgress{Planned: 8, Completed: 3}
	store.runs[run.RunID], store.idem[record.KeyDigest] = run, record
	service := NewRecommendationRunService(store, queue, &step0026IDGenerator{}, &step0026Executor{})
	var wg sync.WaitGroup
	views := make(chan domain.RecommendationRunView, 16)
	errs := make(chan error, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			view, err := service.Cancel(context.Background(), run.RunID)
			views <- view
			errs <- err
		}()
	}
	wg.Wait()
	close(views)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent cancel: %v", err)
		}
	}
	versions := map[uint64]bool{}
	for view := range views {
		if view.Lifecycle != domain.RunCancelling {
			t.Fatalf("concurrent cancel invented terminal: %#v", view)
		}
		versions[view.Version] = true
	}
	if len(versions) != 1 || !versions[5] || len(queue.cancelled) != 1 {
		t.Fatalf("cancel idempotency: versions=%v signals=%d", versions, len(queue.cancelled))
	}
}

func TestStep0026R2RecoveryInterruptsEveryNonterminalAtEveryProgressCount(t *testing.T) {
	store := newStep0026MemoryStore()
	index := 0
	for _, lifecycle := range []domain.RecommendationLifecycle{domain.RunQueued, domain.RunRunning, domain.RunCancelling} {
		for completed := 0; completed <= 8; completed++ {
			index++
			id := domain.RunID(fmt.Sprintf("run_%032x", index+100))
			stage := domain.RunStageQueued
			if lifecycle != domain.RunQueued {
				stage = domain.RunStageCollect
			}
			store.runs[id] = domain.RecommendationRun{SchemaVersion: "1", RunID: id, Request: step0026Request(), RequestHash: "sha256:" + strings.Repeat("a", 64), Lifecycle: lifecycle, Version: uint64(completed + 1), Stage: stage, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: completed}, Artifacts: []domain.RunArtifactIndex{{Name: domain.ArtifactRequest, State: domain.ArtifactPublished}}}
		}
	}
	service := NewRecommendationRunService(store, &step0026Queue{capacity: 32}, &step0026IDGenerator{}, &step0026Executor{})
	if err := service.Recover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(store.runs) != 27 {
		t.Fatalf("orphan fixture changed cardinality: %d", len(store.runs))
	}
	for _, run := range store.runs {
		if run.Lifecycle != domain.RunInterrupted || run.Stage != domain.RunStageTerminal || run.RecommendationStatus != nil || run.Failure != nil || len(run.Artifacts) != 1 {
			t.Fatalf("orphan recovered dishonestly: %#v", run)
		}
	}
}

func TestStep0026R2ArtifactAllowlistOrderIsCanonical(t *testing.T) {
	names := []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance}
	got := append([]domain.RunArtifactName(nil), names...)
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	if reflect.DeepEqual(got, names) {
		t.Fatal("fixture must prove canonical allowlist order is semantic, not lexical")
	}
}
