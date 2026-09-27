/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/snapshot_provider.md
 * @layer L3
 * @updated 2026-09-26
 */
package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bateia/01_core/domain"
)

const maxSnapshotBindings = 20

// SnapshotBinding associates a source identity with a caller-owned private capture.
// The caller must retain the directory unchanged until every provider call finishes.
// This adapter neither captures nor removes directories, and does not authenticate
// their bytes against Git or make a filesystem-wide atomic-snapshot claim.
type SnapshotBinding = domain.SnapshotBinding

type snapshotIdentity struct {
	owner, repo, url, commit string
}

type snapshotRepository struct {
	owner, repo string
}

type SnapshotProvider struct {
	fallback domain.ContentProvider
	bindings map[snapshotIdentity]SnapshotBinding
	client   *Client
}

var _ domain.ContentProvider = (*SnapshotProvider)(nil)
var _ domain.ImplementationSearcher = (*SnapshotProvider)(nil)

func snapshotKey(candidate domain.RepositoryCandidate) snapshotIdentity {
	return snapshotIdentity{candidate.Owner, candidate.Repo, candidate.URL, candidate.CommitID}
}

// NewSnapshotProvider accepts only unambiguous local identities. Source paths are
// not inspected: they may already have changed or disappeared after capture.
func NewSnapshotProvider(fallback domain.ContentProvider, bindings []SnapshotBinding) (*SnapshotProvider, error) {
	if len(bindings) > maxSnapshotBindings {
		return nil, fmt.Errorf("capturas locais excedem limite de %d vínculos", maxSnapshotBindings)
	}
	p := &SnapshotProvider{fallback: fallback, bindings: make(map[snapshotIdentity]SnapshotBinding), client: NewClient()}
	repositories := make(map[snapshotRepository]bool)
	urls := make(map[string]bool)
	var directories []os.FileInfo
	for _, binding := range bindings {
		candidate := binding.Candidate
		if candidate.Owner == "" || candidate.Repo == "" || strings.TrimSpace(candidate.CommitID) == "" || !strings.HasPrefix(candidate.URL, "file://") || !filepath.IsAbs(strings.TrimPrefix(candidate.URL, "file://")) {
			return nil, fmt.Errorf("identidade de captura local incompleta ou não local")
		}
		repository := snapshotRepository{candidate.Owner, candidate.Repo}
		if repositories[repository] || urls[candidate.URL] {
			return nil, fmt.Errorf("identidade de captura local duplicada ou ambígua")
		}
		if !filepath.IsAbs(binding.Directory) || filepath.Clean(binding.Directory) != binding.Directory {
			return nil, fmt.Errorf("diretório de captura local deve ser absoluto e canônico")
		}
		info, err := os.Lstat(binding.Directory)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("diretório de captura local indisponível ou não regular")
		}
		physical, err := filepath.EvalSymlinks(binding.Directory)
		if err != nil || physical != binding.Directory {
			return nil, fmt.Errorf("diretório de captura local com alias não permitido")
		}
		for _, prior := range directories {
			if os.SameFile(prior, info) {
				return nil, fmt.Errorf("diretório de captura local reutilizado por identidades distintas")
			}
		}
		directories = append(directories, info)
		repositories[repository], urls[candidate.URL] = true, true
		p.bindings[snapshotKey(candidate)] = binding
	}
	return p, nil
}

func (p *SnapshotProvider) resolve(candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, bool, error) {
	if !strings.HasPrefix(candidate.URL, "file://") {
		if p.fallback == nil {
			return candidate, false, fmt.Errorf("provedor remoto indisponível; nenhuma fonte alternativa configurada")
		}
		return candidate, false, nil
	}
	binding, ok := p.bindings[snapshotKey(candidate)]
	if !ok {
		return candidate, true, fmt.Errorf("captura local não vinculada à identidade e revisão solicitadas; fonte original não será consultada")
	}
	internal := binding.Candidate
	internal.URL = "file://" + binding.Directory
	return internal, true, nil
}

// Do not interpolate provider errors: they can contain private capture paths.
// Preserve only safe, actionable sentinels used by the pipeline.
func snapshotError(operation string, err error) error {
	if err == nil || errors.Is(err, domain.ErrManifestNotFound) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return fmt.Errorf("captura local: %s indisponível; fonte original não foi consultada", operation)
}

func (p *SnapshotProvider) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	internal, local, err := p.resolve(candidate)
	if err != nil {
		return nil, err
	}
	if !local {
		return p.fallback.FetchManifest(ctx, candidate)
	}
	manifest, err := p.client.FetchManifest(ctx, internal)
	return manifest, snapshotError("manifesto", err)
}

func (p *SnapshotProvider) FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	internal, local, err := p.resolve(candidate)
	if err != nil {
		return nil, err
	}
	if !local {
		return p.fallback.FetchTree(ctx, candidate)
	}
	tree, err := p.client.FetchTree(ctx, internal)
	return tree, snapshotError("árvore", err)
}

func (p *SnapshotProvider) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	internal, local, err := p.resolve(candidate)
	if err != nil {
		return nil, err
	}
	if !local {
		return p.fallback.ExtractRepoFiles(ctx, candidate, paths)
	}
	extracted, err := p.client.ExtractRepoFiles(ctx, internal, paths)
	if err != nil {
		return nil, snapshotError("extração", err)
	}
	extracted.Candidate = p.bindings[snapshotKey(candidate)].Candidate
	extracted.Candidate.Description = snapshotDescription
	for i := range extracted.Files {
		extracted.Files[i].Evidence.Method = "leitura confinada de bytes capturados; sem autenticação Git ou verificação semântica"
	}
	extracted.BaseCommit = ""
	extracted.Coverage.Reason = "Captura local de bytes; metadados Git não capturados e indisponíveis. Bytes não autenticados por commit Git."
	found := false
	for _, check := range extracted.Coverage.ChecksNotRun {
		found = found || check == "git_base_metadata"
	}
	if !found {
		extracted.Coverage.ChecksNotRun = append(extracted.Coverage.ChecksNotRun, "git_base_metadata")
	}
	return extracted, nil
}

const snapshotDescription = "Captura local de bytes para análise consistente; sem autenticação por Git e sem metadata Git capturada"

func (p *SnapshotProvider) FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	_, local, err := p.resolve(candidate)
	if err != nil {
		return candidate, err
	}
	if !local {
		return p.fallback.FetchRepoMetadata(ctx, candidate)
	}
	if err := ctx.Err(); err != nil {
		return candidate, err
	}
	metadata := p.bindings[snapshotKey(candidate)].Candidate
	metadata.Description = snapshotDescription
	metadata.ConsultadoEm = time.Now().UTC()
	return metadata, nil
}

func (p *SnapshotProvider) SearchImplementations(ctx context.Context, candidate domain.RepositoryCandidate, query string, budget domain.ImplementationSearchBudget) (domain.ImplementationSearchReport, error) {
	internal, local, err := p.resolve(candidate)
	if err != nil {
		return domain.ImplementationSearchReport{}, err
	}
	if !local {
		searcher, ok := p.fallback.(domain.ImplementationSearcher)
		if !ok {
			return domain.ImplementationSearchReport{}, fmt.Errorf("provedor remoto não oferece busca de implementação")
		}
		return searcher.SearchImplementations(ctx, candidate, query, budget)
	}
	report, err := p.client.SearchImplementations(ctx, internal, query, budget)
	directory := p.bindings[snapshotKey(candidate)].Directory
	// Diagnostics are the only path-bearing fields generated by the local reader.
	for i := range report.Issues {
		report.Issues[i] = strings.ReplaceAll(report.Issues[i], directory, "[captura local]")
	}
	report.Reason = strings.ReplaceAll(report.Reason, directory, "[captura local]")
	return report, snapshotError("busca de implementação", err)
}
