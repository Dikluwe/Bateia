/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/funnel.md
 * @layer L1
 * @updated 2026-07-25
 */

package usecase

import (
	"context"
	"testing"
	"time"

	"bateia/01_core/domain"
)

// MockSearchProvider simula o provedor de busca de descoberta.
type MockSearchProvider struct {
	Candidates []domain.RepositoryCandidate
}

func (m *MockSearchProvider) SearchCandidates(ctx context.Context, filter domain.QueryFilter) ([]domain.RepositoryCandidate, error) {
	return m.Candidates, nil
}

// MockContentProvider simula a leitura de manifestos e extração.
type MockContentProvider struct{}

func (m *MockContentProvider) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
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
}

func NewMockCacheStore() *MockCacheStore {
	return &MockCacheStore{
		discovery: make(map[string][]domain.RepositoryCandidate),
	}
}

func (m *MockCacheStore) GetDiscovery(queryKey string, ttl time.Duration) ([]domain.RepositoryCandidate, bool, error) {
	c, ok := m.discovery[queryKey]
	return c, ok, nil
}

func (m *MockCacheStore) SaveDiscovery(queryKey string, candidates []domain.RepositoryCandidate) error {
	m.discovery[queryKey] = candidates
	return nil
}

func (m *MockCacheStore) GetManifest(ref domain.CommitRef) (*domain.Manifest, bool, error) {
	return nil, false, nil
}

func (m *MockCacheStore) SaveManifest(ref domain.CommitRef, manifest *domain.Manifest) error {
	return nil
}

func (m *MockCacheStore) GetExtracted(ref domain.CommitRef) (*domain.ExtractedRepository, bool, error) {
	return nil, false, nil
}

func (m *MockCacheStore) SaveExtracted(ref domain.CommitRef, extracted *domain.ExtractedRepository) error {
	return nil
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
