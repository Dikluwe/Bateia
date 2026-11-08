/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/durable_runs_step0026_r2_test.md
 * @layer L4
 * @updated 2026-10-28
 */
package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/03_infra/runstore"
)

type step0026SharedExecutor struct {
	mu               sync.Mutex
	recommendCalls   int
	observedCalls    int
	observedRequests []domain.RecommendationRequest
	result           domain.RecommendationResult
}

type step0026BlockingExecutor struct {
	mu      sync.Mutex
	calls   map[string]int
	entered chan string
	release chan struct{}
}

func (s *step0026BlockingExecutor) Recommend(context.Context, domain.RecommendationRequest) (domain.RecommendationResult, error) {
	return step0026SharedResult(), nil
}

func (s *step0026BlockingExecutor) RecommendObserved(ctx context.Context, request domain.RecommendationRequest, _ domain.RecommendationProgressReporter) (domain.RecommendationResult, error) {
	s.mu.Lock()
	if s.calls == nil {
		s.calls = map[string]int{}
	}
	s.calls[request.Idea.Problem]++
	s.mu.Unlock()
	s.entered <- request.Idea.Problem
	select {
	case <-ctx.Done():
		return domain.RecommendationResult{}, ctx.Err()
	case <-s.release:
		return step0026SharedResult(), nil
	}
}

func (s *step0026SharedExecutor) Recommend(_ context.Context, request domain.RecommendationRequest) (domain.RecommendationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recommendCalls++
	return s.result, nil
}

func (s *step0026SharedExecutor) RecommendObserved(_ context.Context, request domain.RecommendationRequest, reporter domain.RecommendationProgressReporter) (domain.RecommendationResult, error) {
	s.mu.Lock()
	s.observedCalls++
	s.observedRequests = append(s.observedRequests, request)
	s.mu.Unlock()
	for index, stage := range []domain.RecommendationRunStage{domain.RunStageExpand, domain.RunStageDiscover, domain.RunStageScreen, domain.RunStageCollect, domain.RunStageAssess, domain.RunStageDecide, domain.RunStageAlternatives, domain.RunStagePublish} {
		reporter.ReportRecommendationProgress(domain.RecommendationProgressEvent{Stage: stage, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: index + 1}})
	}
	return s.result, nil
}

func step0026SharedResult() domain.RecommendationResult {
	return domain.RecommendationResult{
		SchemaVersion: "1", Status: domain.IdeaRunInconclusive,
		Decision:     domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown, Matrix: []domain.CapabilityComparison{}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}},
		Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Alternatives: []domain.SoftwareAlternative{}, Excluded: []domain.SoftwareAlternative{}},
		Failures:     []domain.ExecutionFailure{{Scope: "provider", Stage: "discover", Message: "Authorization Bearer TOKEN_DISK_CANARY /private/path README_RAW_CANARY"}},
		Warnings:     []string{"Cookie=COOKIE_DISK_CANARY CODE_RAW_CANARY"},
	}
}

func TestStep0026R2ExportedCompositionUsesOneSharedExecutorForSyncAndAsync(t *testing.T) {
	root := t.TempDir()
	executor := &step0026SharedExecutor{result: step0026SharedResult()}
	requestReader, requestWriter := io.Pipe()
	responseReader, responseWriter := io.Pipe()
	limits := runstore.RunStoreLimits{MaxArtifactBytes: 1 << 20, MaxRunBytes: 4 << 20, MaxTotalBytes: 32 << 20}
	server, shutdown, err := NewDurableMCPServer(nil, executor, root, []byte("0123456789abcdef0123456789abcdef"), limits, 1, 8, requestReader, responseWriter)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- server.Serve(context.Background()) }()
	scanner := bufio.NewScanner(responseReader)
	recommendation := `{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"schema_version":"1","allowed_licenses":["MIT"],"require_local":false,"allow_archived":false,"max_components":2,"extend_min_coverage":0.7,"integration_min_value":0.3},"target":{"name":"pantry"},"search":{"strategy":"literal","language":"","locales":["en-US"],"max_per_query":1,"max_candidates":1}}`
	syncResponse := step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"bateia_recommend","arguments":`+recommendation+`}}`)
	if syncResponse.Error != nil {
		t.Fatalf("sync path: %#v", syncResponse.Error)
	}
	syncBytes, _ := json.Marshal(syncResponse)
	for _, canary := range []string{"TOKEN_DISK_CANARY", "/private/path", "README_RAW_CANARY", "COOKIE_DISK_CANARY", "CODE_RAW_CANARY"} {
		if strings.Contains(string(syncBytes), canary) {
			t.Fatalf("sync projection leaked %s: %s", canary, syncBytes)
		}
	}
	started := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_start_recommendation","arguments":`+strings.Replace(recommendation, `{`, `{"idempotency_key":"shared-object",`, 1)+`}}`))
	runID := started["run_id"].(string)
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		view := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"`+runID+`"}}}`))
		if view["lifecycle"] == "SUCCEEDED" {
			break
		}
		time.Sleep(time.Millisecond)
	}
	executor.mu.Lock()
	syncCalls, asyncCalls := executor.recommendCalls, executor.observedCalls
	requests := append([]domain.RecommendationRequest(nil), executor.observedRequests...)
	executor.mu.Unlock()
	if syncCalls != 1 || asyncCalls != 1 || len(requests) != 1 || requests[0].Idea.Problem != "food waste" {
		t.Fatalf("factory did not share one injected executor: sync=%d async=%d requests=%#v", syncCalls, asyncCalls, requests)
	}
	_ = requestWriter.Close()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	var persisted strings.Builder
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !entry.Type().IsRegular() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		persisted.Write(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"shared-object", "TOKEN_DISK_CANARY", "/private/path", "README_RAW_CANARY", "COOKIE_DISK_CANARY", "CODE_RAW_CANARY"} {
		if strings.Contains(persisted.String(), forbidden) {
			t.Fatalf("store persisted forbidden canary/key %q", forbidden)
		}
	}
}

func TestStep0026R2OpaqueRunIDsHaveNoVisibleInputCorrelation(t *testing.T) {
	root := t.TempDir()
	executor := &step0026SharedExecutor{result: step0026SharedResult()}
	requestReader, requestWriter := io.Pipe()
	responseReader, responseWriter := io.Pipe()
	limits := runstore.RunStoreLimits{MaxArtifactBytes: 1 << 20, MaxRunBytes: 4 << 20, MaxTotalBytes: 32 << 20}
	server, shutdown, err := NewDurableMCPServer(nil, executor, root, []byte("abcdef0123456789abcdef0123456789"), limits, 1, 32, requestReader, responseWriter)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- server.Serve(context.Background()) }()
	scanner := bufio.NewScanner(responseReader)
	format := regexp.MustCompile(`^run_[0-9a-f]{32}$`)
	seen := map[string]bool{}
	for index := 0; index < 16; index++ {
		key := "CORRELATION-CANARY-" + strconv.Itoa(index)
		problem := "USER-PROBLEM-CANARY-" + strconv.Itoa(index)
		fields := []string{key, problem, "USER-" + strconv.Itoa(index), "CAPABILITY-" + strconv.Itoa(index), "CONSTRAINT-" + strconv.Itoa(index), "NON-GOAL-" + strconv.Itoa(index), "STACK-" + strconv.Itoa(index), "License-" + strconv.Itoa(index), "path/project-" + strconv.Itoa(index), "lang-" + strconv.Itoa(index)}
		arguments := `{"idempotency_key":"` + key + `","idea":{"problem":"` + problem + `","users":["` + fields[2] + `"],"capabilities":["` + fields[3] + `"],"constraints":["` + fields[4] + `"],"non_goals":["` + fields[5] + `"],"preferred_stack":["` + fields[6] + `"]},"policy":{"schema_version":"1","allowed_licenses":["` + fields[7] + `"],"require_local":false,"allow_archived":false,"max_components":2,"extend_min_coverage":0.7,"integration_min_value":0.3},"target":{"name":"` + fields[8] + `"},"search":{"strategy":"literal","language":"` + fields[9] + `","locales":["en-US"],"max_per_query":1,"max_candidates":1}}`
		started := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"bateia_start_recommendation","arguments":`+arguments+`}}`))
		id := started["run_id"].(string)
		request := domain.RecommendationRequest{Idea: domain.IdeaSpec{Problem: problem, Users: []string{fields[2]}, Capabilities: []string{fields[3]}, Constraints: []string{fields[4]}, NonGoals: []string{fields[5]}, PreferredStack: []string{fields[6]}}, Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{fields[7]}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .3}, Target: domain.SoftwareAlternativeTarget{Name: fields[8]}, Search: domain.RecommendationSearch{Strategy: "literal", Language: fields[9], Locales: []string{"en-US"}, MaxPerQuery: 1, MaxCandidates: 1}}
		_, requestHash, hashErr := usecase.CanonicalRecommendationRequestHash(request)
		if hashErr != nil {
			t.Fatal(hashErr)
		}
		suffix := strings.TrimPrefix(id, "run_")
		correlated := false
		for _, field := range fields {
			digest := sha256.Sum256([]byte(field))
			correlated = correlated || suffix == hex.EncodeToString(digest[:16]) || strings.Contains(id, field)
		}
		if !format.MatchString(id) || seen[id] || correlated || strings.Contains(requestHash, suffix) {
			t.Fatalf("correlated/nonopaque run id %q for key=%q problem=%q", id, key, problem)
		}
		seen[id] = true
	}
	if len(seen) != 16 {
		t.Fatalf("ID generator collisions: %v", seen)
	}
	_ = requestWriter.Close()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestStep0026R2QueuedAndRunningCancellationBoundariesAreObservable(t *testing.T) {
	root := t.TempDir()
	executor := &step0026BlockingExecutor{entered: make(chan string, 4), release: make(chan struct{})}
	requestReader, requestWriter := io.Pipe()
	responseReader, responseWriter := io.Pipe()
	limits := runstore.RunStoreLimits{MaxArtifactBytes: 1 << 20, MaxRunBytes: 4 << 20, MaxTotalBytes: 32 << 20}
	server, shutdown, err := NewDurableMCPServer(nil, executor, root, []byte("fedcba9876543210fedcba9876543210"), limits, 1, 4, requestReader, responseWriter)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- server.Serve(context.Background()) }()
	scanner := bufio.NewScanner(responseReader)
	start := func(key, problem string, rpcID int) string {
		arguments := `{"idempotency_key":"` + key + `","idea":{"problem":"` + problem + `","capabilities":["inventory"]},"policy":{"schema_version":"1","allowed_licenses":["MIT"],"require_local":false,"allow_archived":false,"max_components":2,"extend_min_coverage":0.7,"integration_min_value":0.3},"target":{"name":"pantry"},"search":{"strategy":"literal","language":"","locales":["en-US"],"max_per_query":1,"max_candidates":1}}`
		view := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":`+strconv.Itoa(rpcID)+`,"method":"tools/call","params":{"name":"bateia_start_recommendation","arguments":`+arguments+`}}`))
		return view["run_id"].(string)
	}
	runningID := start("running-cancel", "running-blocker", 1)
	if entered := <-executor.entered; entered != "running-blocker" {
		t.Fatalf("unexpected first execution %q", entered)
	}
	queuedID := start("queued-cancel", "must-never-start", 2)
	queuedCancelled := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_cancel_run","arguments":{"run_id":"`+queuedID+`"}}}`))
	if queuedCancelled["lifecycle"] != "CANCELLED" {
		t.Fatalf("queued cancellation not immediately confirmed: %#v", queuedCancelled)
	}
	runningCancelling := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"bateia_cancel_run","arguments":{"run_id":"`+runningID+`"}}}`))
	if runningCancelling["lifecycle"] != "CANCELLING" {
		t.Fatalf("running cancel presumed terminal: %#v", runningCancelling)
	}
	duplicate := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"bateia_cancel_run","arguments":{"run_id":"`+runningID+`"}}}`))
	if duplicate["version"] != runningCancelling["version"] {
		t.Fatalf("duplicate cancel changed version: %#v %#v", runningCancelling, duplicate)
	}
	deadline := time.Now().Add(time.Second)
	var cancelled map[string]interface{}
	for time.Now().Before(deadline) {
		cancelled = step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"`+runningID+`"}}}`))
		if cancelled["lifecycle"] == "CANCELLED" {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if cancelled["lifecycle"] != "CANCELLED" {
		t.Fatalf("worker acknowledgement did not publish CANCELLED: %#v", cancelled)
	}
	close(executor.release)
	time.Sleep(10 * time.Millisecond)
	executor.mu.Lock()
	targetCalls := executor.calls["must-never-start"]
	executor.mu.Unlock()
	if targetCalls != 0 {
		t.Fatalf("queued cancelled recommendation started %d times", targetCalls)
	}
	_ = requestWriter.Close()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}
