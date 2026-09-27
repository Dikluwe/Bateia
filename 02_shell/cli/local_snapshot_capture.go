/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/local_snapshot_capture.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const localCollectionVersion = "frozen-local-v1"
const maxCapturedGitConfigBytes = 64 << 10

// Private copies are immutable by workflow, not against another process with
// the same OS identity. They are not a point-in-time filesystem snapshot.
type capturedLocalRun struct {
	Candidates    []domain.RepositoryCandidate
	Bindings      []domain.SnapshotBinding
	Identities    map[string][]domain.RepositoryIdentity
	directory     string
	directoryInfo os.FileInfo
}

func (s *capturedLocalRun) Close() error {
	if s == nil || s.directory == "" {
		return nil
	}
	info, err := os.Lstat(s.directory)
	if os.IsNotExist(err) {
		s.directory = ""
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !os.SameFile(info, s.directoryInfo) {
		return fmt.Errorf("diretório privado de captura mudou; limpeza recusada")
	}
	if err := os.RemoveAll(s.directory); err != nil {
		return fmt.Errorf("não foi possível limpar captura local privada")
	}
	s.directory = ""
	return nil
}

func normalizedLocalHashRoots(ctx context.Context, paths []string, excluded string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(paths) > maxHashLocalRoots {
		return nil, fmt.Errorf("hash admite no máximo 20 raízes locais")
	}
	roots := make([]string, 0, len(paths))
	seen := map[string]bool{}
	identities := map[string]bool{}
	for _, path := range paths {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if len(path) > maxHashTreePath {
			return nil, fmt.Errorf("caminho local excede limite de 4096 bytes")
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		if seen[abs] {
			return nil, fmt.Errorf("raiz local duplicada")
		}
		seen[abs] = true
		if err := validateHashLayout(abs, excluded); err != nil {
			return nil, err
		}
		identity := localHashCandidate(abs, "").Repo
		if identities[identity] {
			return nil, fmt.Errorf("identidade local ambígua entre raízes")
		}
		identities[identity] = true
		roots = append(roots, abs)
	}
	return roots, nil
}

func captureLocalCandidates(ctx context.Context, paths []string, excluded string, captureMetadata bool) (result *capturedLocalRun, returnErr error) {
	roots, err := normalizedLocalHashRoots(ctx, paths, excluded)
	if err != nil {
		return nil, err
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("captura exige raiz local")
	}
	directory, err := os.MkdirTemp("", "bateia-local-snapshot-")
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar captura local privada")
	}
	info, err := os.Lstat(directory)
	if err != nil {
		os.RemoveAll(directory)
		return nil, fmt.Errorf("não foi possível validar captura privada")
	}
	snapshot := &capturedLocalRun{directory: directory, directoryInfo: info, Identities: map[string][]domain.RepositoryIdentity{}}
	defer func() {
		if returnErr != nil {
			if cleanupErr := snapshot.Close(); cleanupErr != nil {
				returnErr = fmt.Errorf("%w; %v", returnErr, cleanupErr)
			}
		}
	}()
	// Never let traversal ingest its own copy. Broad roots like /tmp or / are
	// refused when the host temporary directory is nested under them.
	for _, root := range roots {
		rel, err := filepath.Rel(root, directory)
		if err != nil {
			return nil, err
		}
		if rel == "." || filepath.IsLocal(rel) {
			return nil, fmt.Errorf("armazenamento temporário sobrepõe raiz local")
		}
	}
	budget := newTreeHashBudget()
	for i, rootPath := range roots {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		candidate := localHashCandidate(rootPath, "")
		if captureMetadata {
			snapshot.Identities[candidate.URL] = localGitHubRepositoryIdentitiesLimit(candidate.URL, maxCapturedGitConfigBytes)
		}
		source, err := os.OpenRoot(rootPath)
		if err != nil {
			return nil, err
		}
		copyPath := filepath.Join(directory, fmt.Sprintf("root-%02d", i))
		if err := os.Mkdir(copyPath, 0700); err != nil {
			source.Close()
			return nil, fmt.Errorf("não foi possível criar raiz capturada")
		}
		destination, err := os.OpenRoot(copyPath)
		if err != nil {
			source.Close()
			return nil, fmt.Errorf("não foi possível abrir raiz capturada")
		}
		excludedRelative := ""
		if excluded != "" {
			output, err := filepath.Abs(excluded)
			if err != nil {
				source.Close()
				destination.Close()
				return nil, err
			}
			rel, err := filepath.Rel(rootPath, output)
			if err != nil {
				source.Close()
				destination.Close()
				return nil, err
			}
			if filepath.IsLocal(rel) {
				excludedRelative = filepath.ToSlash(rel)
			}
		}
		digest, hashErr := hashRootedTreeTo(ctx, source, excludedRelative, budget, destination)
		source.Close()
		destination.Close()
		if hashErr != nil {
			return nil, hashErr
		}
		candidate.CommitID = digest
		snapshot.Candidates = append(snapshot.Candidates, candidate)
		snapshot.Bindings = append(snapshot.Bindings, domain.SnapshotBinding{Candidate: candidate, Directory: copyPath})
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	sort.Slice(snapshot.Candidates, func(i, j int) bool { return snapshot.Candidates[i].Repo < snapshot.Candidates[j].Repo })
	return snapshot, nil
}

// SnapshotProviderFactory creates a ContentProvider backed by captured snapshot bindings.
type SnapshotProviderFactory func(fallback domain.ContentProvider, bindings []domain.SnapshotBinding) (domain.ContentProvider, error)

var snapshotProviderFactory SnapshotProviderFactory

// SetSnapshotProviderFactory registers the snapshot provider builder injected by L4 wiring.
func SetSnapshotProviderFactory(factory SnapshotProviderFactory) {
	snapshotProviderFactory = factory
}

func buildSnapshotProvider(fallback domain.ContentProvider, bindings []domain.SnapshotBinding) (domain.ContentProvider, error) {
	if snapshotProviderFactory != nil {
		return snapshotProviderFactory(fallback, bindings)
	}
	return newCLISnapshotProvider(fallback, bindings)
}

type cliSnapshotProvider struct {
	fallback domain.ContentProvider
	bindings map[string]domain.SnapshotBinding
}

func newCLISnapshotProvider(fallback domain.ContentProvider, bindings []domain.SnapshotBinding) (domain.ContentProvider, error) {
	m := make(map[string]domain.SnapshotBinding, len(bindings))
	for _, b := range bindings {
		m[b.Candidate.URL] = b
		m["file://"+b.Directory] = b
	}
	return &cliSnapshotProvider{fallback: fallback, bindings: m}, nil
}

func (p *cliSnapshotProvider) resolve(candidate domain.RepositoryCandidate) (domain.SnapshotBinding, bool) {
	b, ok := p.bindings[candidate.URL]
	if ok {
		return b, true
	}
	b, ok = p.bindings["file://"+candidate.URL]
	return b, ok
}

func (p *cliSnapshotProvider) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	b, ok := p.resolve(candidate)
	if !ok {
		if p.fallback != nil {
			return p.fallback.FetchManifest(ctx, candidate)
		}
		return nil, domain.ErrManifestNotFound
	}
	data, err := os.ReadFile(filepath.Join(b.Directory, "go.mod"))
	if err != nil {
		return nil, domain.ErrManifestNotFound
	}
	parsed, err := parser.ParseGoMod(domain.CommitRef{Commit: candidate.CommitID}, string(data))
	if err != nil {
		return &domain.Manifest{ModuleName: candidate.Repo, Content: data}, nil
	}
	return parsed, nil
}

func (p *cliSnapshotProvider) FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	b, ok := p.resolve(candidate)
	if !ok {
		if p.fallback != nil {
			return p.fallback.FetchTree(ctx, candidate)
		}
		return nil, nil
	}
	var nodes []domain.DirectoryNode
	_ = filepath.Walk(b.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == b.Directory {
			return nil
		}
		rel, _ := filepath.Rel(b.Directory, path)
		nodeType := "file"
		if info.IsDir() {
			nodeType = "directory"
		}
		nodes = append(nodes, domain.DirectoryNode{Path: rel, Type: nodeType, Size: info.Size()})
		return nil
	})
	return nodes, nil
}

func (p *cliSnapshotProvider) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	b, ok := p.resolve(candidate)
	if !ok {
		if p.fallback != nil {
			return p.fallback.ExtractRepoFiles(ctx, candidate, paths)
		}
		return nil, fmt.Errorf("candidato não vinculado à captura local")
	}
	extracted := &domain.ExtractedRepository{
		Candidate: candidate,
		Ref:       domain.CommitRef{Commit: candidate.CommitID},
		Coverage:  domain.AnalysisCoverage{Reason: "Captura local de bytes; metadados Git não capturados e indisponíveis. Bytes não autenticados por commit Git."},
	}
	for _, rel := range paths {
		target := filepath.Join(b.Directory, rel)
		data, err := os.ReadFile(target)
		if err == nil {
			extracted.Files = append(extracted.Files, domain.ExtractedFile{
				Path:     rel,
				Content:  string(data),
				Evidence: domain.FileEvidence{Method: "leitura confinada de bytes capturados; sem autenticação Git ou verificação semântica"},
			})
		}
	}
	return extracted, nil
}

func (p *cliSnapshotProvider) FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	candidate.Description = "Captura local de bytes para análise consistente; sem autenticação por Git e sem metadata Git capturada"
	return candidate, nil
}

func (p *cliSnapshotProvider) SearchImplementations(ctx context.Context, candidate domain.RepositoryCandidate, query string, budget domain.ImplementationSearchBudget) (domain.ImplementationSearchReport, error) {
	if p.fallback != nil {
		if searcher, ok := p.fallback.(domain.ImplementationSearcher); ok {
			return searcher.SearchImplementations(ctx, candidate, query, budget)
		}
	}
	return domain.ImplementationSearchReport{}, nil
}
