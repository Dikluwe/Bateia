/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/funnel_test.md
 * @layer L1
 * @updated 2026-10-28
 */

package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"bateia/01_core/domain"
)

// MockSearchProvider simula o provedor de busca de descoberta.
type MockSearchProvider struct {
	Candidates []domain.RepositoryCandidate
	Err        error
}

func (m *MockSearchProvider) SearchCandidates(ctx context.Context, filter domain.QueryFilter) ([]domain.RepositoryCandidate, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Candidates, nil
}

// MockContentProvider simula a leitura de manifestos e extração.
type MockContentProvider struct {
	ManifestErr error
	Delay       time.Duration
	active      int32
	maxActive   int32
}

func (m *MockContentProvider) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	active := atomic.AddInt32(&m.active, 1)
	defer atomic.AddInt32(&m.active, -1)
	for {
		maximum := atomic.LoadInt32(&m.maxActive)
		if active <= maximum || atomic.CompareAndSwapInt32(&m.maxActive, maximum, active) {
			break
		}
	}
	if m.Delay > 0 {
		time.Sleep(m.Delay)
	}
	if m.ManifestErr != nil {
		return nil, m.ManifestErr
	}
	return &domain.Manifest{
		Ref:          domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID},
		ModuleName:   "github.com/" + candidate.Owner + "/" + candidate.Repo,
		GoVersion:    "1.22",
		ConsultadoEm: time.Now().UTC(),
	}, nil
}

func (m *MockContentProvider) FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{
		{Path: "cmd", Type: "dir"},
		{Path: "internal", Type: "dir"},
	}, nil
}

func (m *MockContentProvider) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	return &domain.ExtractedRepository{
		Ref:       domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID},
		Candidate: candidate,
		Files: []domain.ExtractedFile{
			{Path: "go.mod", Content: "module " + candidate.Repo, Size: 20},
		},
		ConsultadoEm: time.Now().UTC(),
	}, nil
}

func (m *MockContentProvider) FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	candidate.Stars = 1234
	return candidate, nil
}

// MockCacheStore simula um cache em memória.
type MockCacheStore struct {
	discovery map[string][]domain.RepositoryCandidate
	Err       error
}

func NewMockCacheStore() *MockCacheStore {
	return &MockCacheStore{
		discovery: make(map[string][]domain.RepositoryCandidate),
	}
}

func (m *MockCacheStore) GetDiscovery(queryKey string, ttl time.Duration) ([]domain.RepositoryCandidate, bool, error) {
	if m.Err != nil {
		return nil, false, m.Err
	}
	c, ok := m.discovery[queryKey]
	return c, ok, nil
}

func (m *MockCacheStore) SaveDiscovery(queryKey string, candidates []domain.RepositoryCandidate) error {
	if m.Err != nil {
		return m.Err
	}
	m.discovery[queryKey] = candidates
	return nil
}

func (m *MockCacheStore) GetManifest(ref domain.CommitRef) (*domain.Manifest, bool, error) {
	return nil, false, m.Err
}

func (m *MockCacheStore) SaveManifest(ref domain.CommitRef, manifest *domain.Manifest) error {
	return m.Err
}

func (m *MockCacheStore) GetExtracted(ref domain.CommitRef) (*domain.ExtractedRepository, bool, error) {
	return nil, false, m.Err
}

func (m *MockCacheStore) SaveExtracted(ref domain.CommitRef, extracted *domain.ExtractedRepository) error {
	return m.Err
}

type MockProgressReporter struct {
	Events []domain.ProgressEvent
}

func (m *MockProgressReporter) ReportProgress(event domain.ProgressEvent) {
	m.Events = append(m.Events, event)
}

func TestSearchFunnelExecution(t *testing.T) {
	searcher := &MockSearchProvider{
		Candidates: []domain.RepositoryCandidate{
			{Owner: "testowner", Repo: "testrepo", CommitID: "c1", Stars: 50},
		},
	}
	fetcher := &MockContentProvider{}
	cacheStore := NewMockCacheStore()

	funnel := NewSearchFunnel(searcher, fetcher, cacheStore)

	filter := domain.QueryFilter{
		Query:      "hexagonal",
		Language:   "go",
		MaxResults: 5,
		MustHave:   []string{"cmd", "internal"},
		Depth:      1,
	}

	report, err := funnel.Execute(context.Background(), filter)
	if err != nil {
		t.Fatalf("Erro inesperado na execução do funil: %v", err)
	}

	if report.TopCandidates != 1 {
		t.Errorf("Esperado 1 candidato no topo, obtido: %d", report.TopCandidates)
	}

	if report.TriagedCount != 1 {
		t.Errorf("Esperado 1 candidato triado, obtido: %d", report.TriagedCount)
	}

	if report.ExtractedCount != 1 {
		t.Errorf("Esperado 1 candidato extraído, obtido: %d", report.ExtractedCount)
	}
}

func TestSearchFunnelPolyglotPreservesPartialFailures(t *testing.T) {
	providerErr := errors.New("provider unavailable")
	funnel := NewSearchFunnel(
		&MockSearchProvider{Err: providerErr},
		&MockContentProvider{},
		NewMockCacheStore(),
	)

	report, err := funnel.Execute(context.Background(), domain.QueryFilter{
		Query:      "architecture",
		Language:   "go,rust",
		MaxResults: 5,
		Depth:      1,
	})
	if err != nil {
		t.Fatalf("o sweep deve preservar falhas parciais no relatório: %v", err)
	}
	if len(report.Failures) != 2 {
		t.Fatalf("esperava uma falha para cada linguagem, recebeu: %+v", report.Failures)
	}
	for _, failure := range report.Failures {
		if failure.Message != providerErr.Error() && !strings.Contains(failure.Message, providerErr.Error()) {
			t.Errorf("falha não preservou erro do provedor: %+v", failure)
		}
	}
}

func TestSearchFunnelDiscardsRepositoryWithoutSupportedManifest(t *testing.T) {
	funnel := NewSearchFunnel(
		&MockSearchProvider{Candidates: []domain.RepositoryCandidate{{Owner: "empty", Repo: "repo", CommitID: "HEAD"}}},
		&MockContentProvider{ManifestErr: domain.ErrManifestNotFound},
		NewMockCacheStore(),
	)

	report, err := funnel.Execute(context.Background(), domain.QueryFilter{
		Query:      "architecture",
		Language:   "go",
		MaxResults: 5,
		Depth:      1,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if report.TriagedCount != 0 || report.ExtractedCount != 0 || report.DiscardedCount != 1 {
		t.Fatalf("funil inesperado: %d triados, %d extraídos, %d descartados", report.TriagedCount, report.ExtractedCount, report.DiscardedCount)
	}
}

func TestSearchFunnelReportsCacheDegradationAndContinues(t *testing.T) {
	cacheErr := errors.New("cache unavailable")
	cacheStore := NewMockCacheStore()
	cacheStore.Err = cacheErr
	reporter := &MockProgressReporter{}
	funnel := NewSearchFunnel(
		&MockSearchProvider{Candidates: []domain.RepositoryCandidate{{Owner: "test", Repo: "repo", CommitID: "c1"}}},
		&MockContentProvider{},
		cacheStore,
	).WithProgressReporter(reporter)

	report, err := funnel.Execute(context.Background(), domain.QueryFilter{Query: "test", Language: "go", MaxResults: 1, Depth: 1})
	if err != nil {
		t.Fatalf("falha de cache deve ser recuperável: %v", err)
	}
	if report.ExtractedCount != 1 {
		t.Fatalf("esperava resultado apesar do cache indisponível, recebeu %d", report.ExtractedCount)
	}
	if len(report.Warnings) == 0 {
		t.Fatal("esperava avisos estruturados de cache")
	}
	if len(reporter.Events) == 0 {
		t.Fatal("esperava eventos de progresso pela abstração injetada")
	}
}

func TestSearchFunnelLimitsConcurrentTriageWorkers(t *testing.T) {
	candidates := make([]domain.RepositoryCandidate, 8)
	for i := range candidates {
		candidates[i] = domain.RepositoryCandidate{Owner: "test", Repo: fmt.Sprintf("repo-%d", i), CommitID: "c1"}
	}
	provider := &MockContentProvider{Delay: 10 * time.Millisecond}
	funnel := NewSearchFunnel(&MockSearchProvider{Candidates: candidates}, provider, NewMockCacheStore())
	_, err := funnel.Execute(context.Background(), domain.QueryFilter{Query: "test", Language: "go", MaxResults: 8, Depth: 1})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	maximum := atomic.LoadInt32(&provider.maxActive)
	if maximum < 2 || maximum > 4 {
		t.Fatalf("concorrência esperada entre 2 e 4 workers, observada: %d", maximum)
	}
}

// Local static analysis must be injected through the port and remain distinct from relevance.
type semanticSpy0035 struct{ called bool }

func (s *semanticSpy0035) Analyze(ctx context.Context, req domain.SemanticRequest) (domain.SemanticReport, error) {
	s.called = true
	return domain.SemanticReport{Status: "partial", Reason: "fixture: budget limitado", Units: []domain.SemanticUnit{{Path: "render.cpp", Status: "complete", Symbols: []domain.SemanticSymbol{{Name: "render", File: "render.cpp", Line: 2, EndLine: 4, Definition: true}}}}}, nil
}
func TestFunnelInjectsStaticAnalysisWithoutSolutionPromotion0035(t *testing.T) {
	spy := &semanticSpy0035{}
	f := NewSearchFunnel(&MockSearchProvider{}, &MockContentProvider{}, &MockCacheStore{}).WithSemanticProvider(spy)
	result, err := f.Execute(context.Background(), domain.QueryFilter{Query: "transparency", Language: "cpp", MaxResults: 1, Depth: 1, LocalPath: "/fixture"})
	if err != nil || !spy.called || len(result.Results) != 1 {
		t.Fatalf("semantic port not invoked: %v %+v", err, result)
	}
	repo := result.Results[0]
	if repo.Semantic == nil || repo.Semantic.Status != "partial" || repo.HealthScore != -1 || repo.Coverage.Status != "insufficient" {
		t.Fatalf("static evidence promoted: %+v", repo)
	}
}

type manifestTimestampProvider struct {
	MockContentProvider
	collectedAt time.Time
}

func (p *manifestTimestampProvider) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	extracted, err := p.MockContentProvider.ExtractRepoFiles(ctx, candidate, paths)
	if err == nil {
		extracted.ConsultadoEm = p.collectedAt
	}
	return extracted, err
}

func TestFunnelReparsePreservesCollectionInstant(t *testing.T) {
	for _, collectedAt := range []time.Time{time.Date(2026, 10, 28, 10, 0, 0, 321, time.UTC), {}} {
		provider := &manifestTimestampProvider{collectedAt: collectedAt}
		funnel := NewSearchFunnel(&MockSearchProvider{Candidates: []domain.RepositoryCandidate{{Owner: "example", Repo: "app", CommitID: "abc123"}}}, provider, NewMockCacheStore())
		report, err := funnel.Execute(context.Background(), domain.QueryFilter{Query: "app", Language: "go", MaxResults: 1, Depth: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(report.Results) != 1 || report.Results[0].Manifest == nil {
			t.Fatalf("missing parsed extraction: %+v", report)
		}
		if !report.Results[0].Manifest.ConsultadoEm.Equal(collectedAt) {
			t.Fatalf("reparse replaced collection instant: %v", report.Results[0].Manifest.ConsultadoEm)
		}
	}
}
