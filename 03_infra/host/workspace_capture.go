/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/workspace_capture.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/03_infra/local"
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

func localHashCandidate(abs, digest string) domain.RepositoryCandidate {
	return usecase.LocalRepositoryCandidate(abs, filepath.Base(abs), digest)
}
func (s *capturedLocalRun) View() domain.CapturedWorkspaceView {
	return domain.CapturedWorkspaceView{Candidates: s.Candidates, Bindings: s.Bindings, Identities: s.Identities, ContainerPath: s.directory}
}
func (*Host) CaptureLocal(ctx context.Context, paths []string, excluded string, metadata bool) (domain.CapturedWorkspace, error) {
	s, err := captureLocalCandidates(ctx, paths, excluded, metadata)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (*Host) NormalizedLocalRoots(ctx context.Context, paths []string, excluded string) ([]string, error) {
	return normalizedLocalHashRoots(ctx, paths, excluded)
}
func (*Host) SnapshotProvider(fallback domain.ContentProvider, bindings []domain.SnapshotBinding) (domain.ContentProvider, error) {
	return local.NewSnapshotProvider(fallback, bindings)
}
