/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/content/router.md
 * @layer L3
 * @updated 2026-09-08
 */

package content

import (
	"context"
	"strings"

	"bateia/01_core/domain"
)

type Router struct {
	remote domain.ContentProvider
	local  domain.ContentProvider
}

func NewRouter(remote, local domain.ContentProvider) *Router {
	return &Router{remote: remote, local: local}
}

func (r *Router) provider(candidate domain.RepositoryCandidate) domain.ContentProvider {
	if strings.HasPrefix(candidate.URL, "file://") {
		return r.local
	}
	return r.remote
}

func (r *Router) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	return r.provider(candidate).FetchManifest(ctx, candidate)
}

func (r *Router) FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return r.provider(candidate).FetchTree(ctx, candidate)
}

func (r *Router) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	return r.provider(candidate).ExtractRepoFiles(ctx, candidate, paths)
}

func (r *Router) FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	return r.provider(candidate).FetchRepoMetadata(ctx, candidate)
}

// SearchImplementations keeps unsupported providers explicit without widening ContentProvider.
func (r *Router) SearchImplementations(ctx context.Context, candidate domain.RepositoryCandidate, query string, budget domain.ImplementationSearchBudget) (domain.ImplementationSearchReport, error) {
	if searcher, ok := r.provider(candidate).(domain.ImplementationSearcher); ok {
		return searcher.SearchImplementations(ctx, candidate, query, budget)
	}
	return domain.ImplementationSearchReport{Status: "not_evaluated", Method: "go_ast_query_v1", Reason: "provedor sem busca de implementações", Budget: budget, Hits: []domain.ImplementationHit{}, Examined: []domain.ImplementationInput{}}, nil
}
