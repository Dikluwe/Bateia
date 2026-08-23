/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/cache/cache_store.md
 * @layer L3
 * @updated 2026-08-23
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
