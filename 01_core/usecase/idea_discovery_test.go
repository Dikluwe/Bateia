/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_discovery_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"bateia/01_core/domain"
)

type fakeObservableProvider struct {
	mu                  sync.Mutex
	active, peak, calls int
	results             map[string][]domain.RepositoryCandidate
	failures            map[string]bool
}

type deadlineProvider struct{}

func (deadlineProvider) Discover(ctx context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	if request.ExpansionID == "exp-a" {
		return domain.ProviderDiscovery{Provider: "deadline", Attempts: []domain.ProviderAttempt{{Query: request.Query, Status: "SUCCEEDED", ResultCount: 1}}, Candidates: []domain.RepositoryCandidate{candidate("org", "kept")}}, nil
	}
	<-ctx.Done()
	return domain.ProviderDiscovery{Provider: "deadline", Attempts: []domain.ProviderAttempt{{Query: request.Query, Status: "TIMED_OUT", ErrorClass: "TIMEOUT", Error: ctx.Err().Error()}}}, ctx.Err()
}

func (provider *fakeObservableProvider) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	provider.mu.Lock()
	provider.active++
	provider.calls++
	if provider.active > provider.peak {
		provider.peak = provider.active
	}
	provider.mu.Unlock()
	defer func() { provider.mu.Lock(); provider.active--; provider.mu.Unlock() }()
	result := domain.ProviderDiscovery{Provider: "fake", Attempts: []domain.ProviderAttempt{{Query: request.Query}}, Candidates: provider.results[request.ExpansionID]}
	result.Attempts[0].ResultCount = len(result.Candidates)
	if provider.failures[request.ExpansionID] {
		result.Attempts[0].Error = "failed"
		return result, errors.New("provider failed")
	}
	return result, nil
}

func discoveryPlan() domain.ExpansionPlan {
	return domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "target"}, Expansions: []domain.SearchExpansion{
		{ID: "exp-a", Facet: domain.FacetProblem, Query: "query a", Covers: []string{"problem"}},
		{ID: "exp-b", Facet: domain.FacetCapability, Query: "query b", Covers: []string{"capabilities[0]"}},
		{ID: "exp-c", Facet: domain.FacetCapability, Query: "query c", Covers: []string{"capabilities[1]"}},
	}}
}

func candidate(owner, repo string) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{Owner: owner, Repo: repo}
}

func TestDiscoverExpansionPlanDeduplicatesOriginsAndKeepsPartialFailure(t *testing.T) {
	provider := &fakeObservableProvider{results: map[string][]domain.RepositoryCandidate{
		"exp-a": {candidate("Org", "A"), candidate("Org", "B")},
		"exp-b": {candidate("org", "b"), candidate("Org", "C")},
	}, failures: map[string]bool{"exp-c": true}}
	snapshot, err := DiscoverExpansionPlan(context.Background(), discoveryPlan(), "lexical", "go", 10, 2, provider)
	if err != nil {
		t.Fatal(err)
	}
	if provider.calls != 3 || provider.peak > 2 {
		t.Fatalf("chamadas/concorrência incorretas: calls=%d peak=%d", provider.calls, provider.peak)
	}
	if len(snapshot.Candidates) != 3 || snapshot.Candidates[0].Repository.Repo != "A" || snapshot.Candidates[1].Repository.Repo != "B" || snapshot.Candidates[2].Repository.Repo != "C" {
		t.Fatalf("ordem/deduplicação incorreta: %+v", snapshot.Candidates)
	}
	if len(snapshot.Candidates[1].ExpansionIDs) != 2 || len(snapshot.Candidates[1].Queries) != 2 {
		t.Fatalf("origens de B perdidas: %+v", snapshot.Candidates[1])
	}
	if len(snapshot.Failures) != 1 || snapshot.Failures[0].Scope != "exp-c" {
		t.Fatalf("falha parcial perdida: %+v", snapshot.Failures)
	}
}

func TestDiscoverExpansionPlanFailsWhenAllExpansionsFail(t *testing.T) {
	provider := &fakeObservableProvider{results: map[string][]domain.RepositoryCandidate{}, failures: map[string]bool{"exp-a": true, "exp-b": true, "exp-c": true}}
	if _, err := DiscoverExpansionPlan(context.Background(), discoveryPlan(), "literal", "go", 10, 1, provider); err == nil {
		t.Fatal("falha total foi aceita")
	}
}

func TestDiscoverExpansionPlanDoesNotMutateInput(t *testing.T) {
	plan := discoveryPlan()
	original := clonePlan(plan)
	provider := &fakeObservableProvider{results: map[string][]domain.RepositoryCandidate{}, failures: map[string]bool{}}
	if _, err := DiscoverExpansionPlan(context.Background(), plan, "literal", "go", 10, 1, provider); err != nil {
		t.Fatal(err)
	}
	if plan.Expansions[0].Query != original.Expansions[0].Query || plan.Idea.Problem != original.Idea.Problem {
		t.Fatal("plano foi modificado")
	}
}

func TestBoundedDiscoveryPreservesCandidateAndMarksNotStarted(t *testing.T) {
	policy := domain.DiscoveryExecutionPolicy{SchemaVersion: "1", TotalTimeout: 40 * time.Millisecond, AttemptTimeout: 20 * time.Millisecond, Concurrency: 1, MaxAttempts: 3}
	snapshot, err := DiscoverExpansionPlanWithPolicy(context.Background(), discoveryPlan(), "literal", "", 3, policy, deadlineProvider{}, nil)
	if err == nil {
		t.Fatal("deadline must be observable")
	}
	if len(snapshot.Candidates) != 1 || snapshot.Candidates[0].Repository.Repo != "kept" {
		t.Fatalf("completed candidate was lost: %#v", snapshot.Candidates)
	}
	if snapshot.Receipt == nil || snapshot.Receipt.Status != domain.DiscoveryRunPartial || snapshot.Receipt.CompletionCause != domain.CompletionTotalTimeout {
		t.Fatalf("unexpected receipt: %#v", snapshot.Receipt)
	}
	notStarted := 0
	for _, a := range snapshot.Receipt.Attempts {
		if a.Status == "NOT_STARTED" {
			notStarted++
		}
	}
	if notStarted != 1 || snapshot.Receipt.TimedOut != 1 {
		t.Fatalf("timeout/not-started misclassified: %#v", snapshot.Receipt)
	}
}

func TestDiscoveryPolicyRejectsInvalidTimeoutsBeforeProvider(t *testing.T) {
	provider := &fakeObservableProvider{}
	policy := domain.DiscoveryExecutionPolicy{SchemaVersion: "1", TotalTimeout: time.Second, AttemptTimeout: 2 * time.Second, Concurrency: 1, MaxAttempts: 3}
	if _, err := DiscoverExpansionPlanWithPolicy(context.Background(), discoveryPlan(), "literal", "", 3, policy, provider, nil); err == nil {
		t.Fatal("attempt timeout greater than total survived")
	}
	if provider.calls != 0 {
		t.Fatal("provider was touched before policy validation")
	}
}
