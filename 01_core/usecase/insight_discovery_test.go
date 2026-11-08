/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_discovery_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

type insightTestRuntime struct{ durations []time.Duration }

func (r *insightTestRuntime) Now() time.Time { return time.Time{} }
func (r *insightTestRuntime) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	r.durations = append(r.durations, d)
	return context.WithCancel(ctx)
}

type insightTestCatalog struct {
	calls   []domain.CatalogSearchRequest
	respond func(context.Context, domain.CatalogSearchRequest) (domain.CatalogDiscovery, error)
}

func (*insightTestCatalog) Capabilities() domain.ProviderCapabilities {
	return domain.ProviderCapabilities{Provider: "fixture", Roles: []domain.DiscoveryProviderRole{domain.RoleProjectCatalog}}
}
func (p *insightTestCatalog) SearchCatalog(ctx context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
	p.calls = append(p.calls, r)
	if p.respond != nil {
		return p.respond(ctx, r)
	}
	return insightCatalogResult(r, insightCatalogCandidate("org", "queue")), nil
}
func insightCatalogCandidate(owner, repo string) domain.CatalogCandidate {
	return domain.CatalogCandidate{Repository: domain.RepositoryCandidate{Owner: owner, Repo: repo, CommitID: strings.Repeat("a", 40), Description: "backpressure queue", URL: "https://example.invalid/" + owner + "/" + repo}, Topics: []string{"concurrency"}}
}
func insightCatalogResult(r domain.CatalogSearchRequest, candidates ...domain.CatalogCandidate) domain.CatalogDiscovery {
	return domain.CatalogDiscovery{Provider: "fixture", Role: domain.RoleProjectCatalog, Candidates: candidates, Attempt: domain.CatalogAttempt{ExpansionID: r.ExpansionID, Query: r.Query, Status: "SUCCEEDED", ResultCount: len(candidates)}}
}

type insightTestReader struct {
	treeCalls, extractCalls int
	tree                    []domain.DirectoryNode
	content                 string
	mutate                  func(*domain.ExtractedRepository)
}

func (p *insightTestReader) FetchTree(_ context.Context, _ domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	p.treeCalls++
	if p.tree != nil {
		return p.tree, nil
	}
	return []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, nil
}
func (p *insightTestReader) ExtractRepoFiles(_ context.Context, r domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.extractCalls++
	content := p.content
	if content == "" {
		content = "# Queue\nA bounded mechanism with cooperative cancellation.\n"
	}
	out := &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: r.Owner, Repo: r.Repo, Commit: r.CommitID}, Files: []domain.ExtractedFile{{Path: paths[0], Content: content}}}
	if p.mutate != nil {
		p.mutate(out)
	}
	return out, nil
}

func TestInsightDiscoveryRunsEveryRouteAndFeedsExploration(t *testing.T) {
	runtime := &insightTestRuntime{}
	catalog := &insightTestCatalog{}
	reader := &insightTestReader{}
	got, err := DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), runtime), insightDiscoveryFixture(), catalog, reader)
	if err != nil || len(catalog.calls) != 5 || len(got.Candidates) != 1 || len(got.Candidates[0].QueryIDs) != 5 || len(got.Bundles) != 1 || len(got.Exploration.Cards) != 1 || got.Status != "REFERENCES_COLLECTED" {
		t.Fatalf("workflow incomplete: %v %+v", err, got)
	}
	if got.NoveltyState != "NOT_ASSESSED" || got.EffectState != "UNVERIFIED" || got.Bundles[0].Complete || got.Readmes[0].Scope != "ROOT_README_ONLY" {
		t.Fatal("discovery promoted semantic evidence")
	}
	if len(runtime.durations) != 7 || runtime.durations[0] != 90*time.Second {
		t.Fatal("timeouts not acquired through runtime")
	}
	for _, call := range catalog.calls {
		if call.AttemptTimeout != 10*time.Second || call.MaxResults != 2 {
			t.Fatal("query budgets not forwarded")
		}
	}
	if reader.treeCalls != 1 || reader.extractCalls != 1 {
		t.Fatal("readme budget exceeded")
	}
}

func TestInsightDiscoveryPreservesPartialAndSanitizesProviderFailures(t *testing.T) {
	const private = "provider-private-message"
	catalog := &insightTestCatalog{}
	catalog.respond = func(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
		if len(catalog.calls) == 1 {
			return insightCatalogResult(r, insightCatalogCandidate("org", "kept")), nil
		}
		out := insightCatalogResult(r)
		out.Attempt.Status = "RATE_LIMITED"
		return out, errors.New(private)
	}
	got, err := DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), insightDiscoveryFixture(), catalog, &insightTestReader{})
	data, _ := json.Marshal(got)
	if err != nil || got.Status != "PARTIAL" || len(got.Candidates) != 1 || len(catalog.calls) != 5 || strings.Contains(string(data), private) || got.Attempts[1].ErrorCode != "RATE_LIMIT" {
		t.Fatal("partial observations discarded or unsafe error exposed")
	}
	catalog.respond = func(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
		return insightCatalogResult(r), errors.New(private)
	}
	got, err = DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), insightDiscoveryFixture(), catalog, &insightTestReader{})
	if !errors.Is(err, domain.ErrInsightDiscoveryUnavailable) || got.ID == "" || got.Status != "FAILED" {
		t.Fatal("total failure not reviewable")
	}
}

func TestInsightDiscoveryInvalidResponseAndCancellationDoNotInventResults(t *testing.T) {
	catalog := &insightTestCatalog{}
	catalog.respond = func(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
		out := insightCatalogResult(r, insightCatalogCandidate("org", "bad"))
		out.Role = domain.RoleCodeProof
		return out, nil
	}
	got, err := DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), insightDiscoveryFixture(), catalog, &insightTestReader{})
	if !errors.Is(err, domain.ErrInsightDiscoveryUnavailable) || len(got.Candidates) != 0 || got.Attempts[0].Status != "INTEGRITY_ERROR" {
		t.Fatal("invalid provider role admitted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	catalog = &insightTestCatalog{}
	catalog.respond = func(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
		if len(catalog.calls) == 2 {
			cancel()
			return insightCatalogResult(r), context.Canceled
		}
		return insightCatalogResult(r, insightCatalogCandidate("org", "kept")), nil
	}
	reader := &insightTestReader{}
	got, err = DiscoverInsightReferences(domain.WithExecutionRuntime(ctx, &insightTestRuntime{}), insightDiscoveryFixture(), catalog, reader)
	if !errors.Is(err, context.Canceled) || got.Status != "CANCELLED" || len(got.Candidates) != 1 || len(catalog.calls) != 2 || reader.treeCalls != 0 || got.Attempts[2].Status != "NOT_STARTED" {
		t.Fatal("cancelled workflow continued or lost prior candidates")
	}
}

func TestInsightDiscoveryEmptyAndConfigurationFailures(t *testing.T) {
	catalog := &insightTestCatalog{respond: func(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
		return insightCatalogResult(r), nil
	}}
	got, err := DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), insightDiscoveryFixture(), catalog, &insightTestReader{})
	if err != nil || got.Status != "NO_CANDIDATES" || got.Exploration.Status != "SEARCH_NEEDED" || len(catalog.calls) != 5 {
		t.Fatal("valid emptiness conflated with failed search")
	}
	if _, err := DiscoverInsightReferences(nil, insightDiscoveryFixture(), catalog, &insightTestReader{}); !errors.Is(err, domain.ErrInvalidRunInput) {
		t.Fatal("nil context accepted")
	}
	if _, err := DiscoverInsightReferences(context.Background(), insightDiscoveryFixture(), catalog, &insightTestReader{}); !errors.Is(err, domain.ErrRuntimeNotConfigured) {
		t.Fatal("implicit runtime created")
	}
	if _, err := DiscoverInsightReferences(context.Background(), insightDiscoveryFixture(), nil, &insightTestReader{}); !errors.Is(err, domain.ErrInvalidRunInput) {
		t.Fatal("nil catalog accepted")
	}
}

type insightExpiredContext struct {
	context.Context
	done chan struct{}
}

func (ctx insightExpiredContext) Done() <-chan struct{} { return ctx.done }
func (insightExpiredContext) Err() error                { return context.DeadlineExceeded }

type insightImmediateDeadlineRuntime struct{}

func (insightImmediateDeadlineRuntime) Now() time.Time { return time.Time{} }
func (insightImmediateDeadlineRuntime) WithTimeout(ctx context.Context, _ time.Duration) (context.Context, context.CancelFunc) {
	done := make(chan struct{})
	close(done)
	return insightExpiredContext{Context: ctx, done: done}, func() {}
}

func TestInsightDiscoveryExhaustedDeadlineStartsNoCallsAndKeepsReceipt(t *testing.T) {
	catalog, reader := &insightTestCatalog{}, &insightTestReader{}
	ctx := domain.WithExecutionRuntime(context.Background(), insightImmediateDeadlineRuntime{})
	report, err := DiscoverInsightReferences(ctx, insightDiscoveryFixture(), catalog, reader)
	if !errors.Is(err, context.DeadlineExceeded) || report.Status != "TIMED_OUT" || report.ID == "" || len(catalog.calls) != 0 || reader.treeCalls != 0 || len(report.Attempts) != 5 {
		t.Fatal("exhausted deadline discarded receipt or started work")
	}
}

func TestInsightDiscoveryCollectedBundlesReplayExplorationAfterJSON(t *testing.T) {
	request := insightDiscoveryFixture()
	report, err := DiscoverInsightReferences(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), request, &insightTestCatalog{}, &insightTestReader{content: "# Reference\n<img src=\"x\" /> queue cancellation & conditions\n"})
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.MarshalIndent(report.Bundles, "", "  ")
	var decoded []domain.EvidenceBundle
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	replay, err := ExploreInsights(context.Background(), domain.InsightExplorationRequest{SchemaVersion: "1", Question: report.Plan.Request.Question, FamiliarApproaches: report.Plan.Request.FamiliarApproaches, MaxPassages: report.Plan.Request.MaxPassages, Bundles: decoded})
	if err != nil {
		t.Fatal(err)
	}
	if replay.ID != report.Exploration.ID {
		before, _ := normalizeInsightExploration(domain.InsightExplorationRequest{SchemaVersion: "1", Question: report.Plan.Request.Question, FamiliarApproaches: report.Plan.Request.FamiliarApproaches, MaxPassages: report.Plan.Request.MaxPassages, Bundles: report.Bundles})
		after, _ := normalizeInsightExploration(domain.InsightExplorationRequest{SchemaVersion: "1", Question: report.Plan.Request.Question, FamiliarApproaches: report.Plan.Request.FamiliarApproaches, MaxPassages: report.Plan.Request.MaxPassages, Bundles: decoded})
		left, right := canonicalJSON(before), canonicalJSON(after)
		for i := 0; i < min(len(left), len(right)); i++ {
			if left[i] != right[i] {
				t.Fatalf("preimage differs at %d: %s / %s", i, string(left[max(0, i-70):min(len(left), i+90)]), string(right[max(0, i-70):min(len(right), i+90)]))
			}
		}
		t.Fatal("context identity lost after JSON")
	}
}
