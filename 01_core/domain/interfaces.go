/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/types.md
 * @layer L1
 * @updated 2026-07-25
 */

package domain

import (
	"context"
	"time"
)

// SearchProvider abstrai o provedor de busca de candidatos no topo do funil (Estágio 1 - Descoberta).
type SearchProvider interface {
	SearchCandidates(ctx context.Context, filter QueryFilter) ([]RepositoryCandidate, error)
}

// ContentProvider abstrai a leitura de manifestos, árvores e arquivos nos estágios 2 (Triagem) e 3 (Extração).
type ContentProvider interface {
	FetchManifest(ctx context.Context, candidate RepositoryCandidate) (*Manifest, error)
	FetchTree(ctx context.Context, candidate RepositoryCandidate) ([]DirectoryNode, error)
	ExtractRepoFiles(ctx context.Context, candidate RepositoryCandidate, paths []string) (*ExtractedRepository, error)
	FetchRepoMetadata(ctx context.Context, candidate RepositoryCandidate) (RepositoryCandidate, error)
}

// CacheStore abstrai a persistência e busca de registros em cache no File System.
type CacheStore interface {
	GetDiscovery(queryKey string, ttl time.Duration) ([]RepositoryCandidate, bool, error)
	SaveDiscovery(queryKey string, candidates []RepositoryCandidate) error

	GetManifest(ref CommitRef) (*Manifest, bool, error)
	SaveManifest(ref CommitRef, manifest *Manifest) error

	GetExtracted(ref CommitRef) (*ExtractedRepository, bool, error)
	SaveExtracted(ref CommitRef, extracted *ExtractedRepository) error
}

// ProgressReporter recebe eventos operacionais sem acoplar o núcleo a log, terminal ou arquivo.
type ProgressReporter interface {
	ReportProgress(event ProgressEvent)
}
