/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/client.md
 * @layer L3
 * @updated 2026-09-07
 */

package local

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

type Client struct{}

func NewClient() *Client { return &Client{} }

var ignoredDirectories = map[string]bool{
	".git": true, ".hg": true, ".svn": true, "node_modules": true, "vendor": true,
	"dist": true, "build": true, "target": true, ".cache": true, ".idea": true, ".vscode": true,
}

func candidateRoot(candidate domain.RepositoryCandidate) (string, error) {
	if !strings.HasPrefix(candidate.URL, "file://") {
		return "", fmt.Errorf("candidato não é local: %s", candidate.URL)
	}
	root := strings.TrimPrefix(candidate.URL, "file://")
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("worktree local inacessível: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("worktree local não é diretório: %s", abs)
	}
	return abs, nil
}

func safePath(root, relative string) (string, error) {
	if filepath.IsAbs(relative) {
		return "", fmt.Errorf("caminho absoluto não permitido: %s", relative)
	}
	clean := filepath.Clean(relative)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("caminho fora da worktree: %s", relative)
	}
	path := filepath.Join(root, clean)
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("caminho fora da worktree: %s", relative)
	}
	return path, nil
}

// readLocalFile bounds the read itself, and validates the opened descriptor.
// Root.OpenFile enforces confinement during path resolution, including symlink races.
const maxLocalFileSize = 2 * 1024 * 1024

// resolveRootAlias translates absolute internal symlinks into relative names.
// Every metadata lookup stays rooted; the resulting name is still opened with
// Root.OpenFile, so this translation is not the security boundary.
func resolveRootAlias(root *os.Root, name string) (string, error) {
	for hops := 0; hops < 40; hops++ {
		if !filepath.IsLocal(name) {
			return "", fmt.Errorf("alias fora da worktree: %s", name)
		}
		parts := strings.Split(filepath.Clean(name), string(filepath.Separator))
		prefix := ""
		replaced := false
		for i, part := range parts {
			prefix = filepath.Join(prefix, part)
			info, err := root.Lstat(prefix)
			if err != nil {
				return "", err
			}
			if info.Mode()&os.ModeSymlink == 0 {
				continue
			}
			target, err := root.Readlink(prefix)
			if err != nil {
				return "", err
			}
			if filepath.IsAbs(target) {
				target, err = filepath.Rel(root.Name(), target)
				if err != nil {
					return "", err
				}
			} else {
				target = filepath.Join(filepath.Dir(prefix), target)
			}
			if !filepath.IsLocal(target) {
				return "", fmt.Errorf("alias fora da worktree: %s", prefix)
			}
			name = filepath.Join(append([]string{target}, parts[i+1:]...)...)
			replaced = true
			break
		}
		if !replaced {
			return name, nil
		}
	}
	return "", fmt.Errorf("limite de resolução de aliases locais excedido")
}

func readLocalFile(root *os.Root, name string) ([]byte, error) {
	f, err := root.OpenFile(name, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		resolved, resolveErr := resolveRootAlias(root, name)
		if resolveErr != nil {
			return nil, resolveErr
		}
		f, err = root.OpenFile(resolved, os.O_RDONLY|syscall.O_NONBLOCK, 0)
		if err != nil {
			return nil, err
		}
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > maxLocalFileSize {
		return nil, fmt.Errorf("arquivo local não regular ou acima de %d bytes: %s", maxLocalFileSize, name)
	}
	data, err := io.ReadAll(io.LimitReader(f, maxLocalFileSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxLocalFileSize {
		return nil, fmt.Errorf("arquivo local excedeu limite durante leitura: %s", name)
	}
	return data, nil
}

func (c *Client) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	root, err := candidateRoot(candidate)
	if err != nil {
		return nil, err
	}
	scoped, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	defer scoped.Close()
	entries, err := fs.ReadDir(scoped.FS(), ".")
	if err != nil {
		return nil, err
	}
	var manifests []string
	for _, entry := range entries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if !entry.IsDir() && parser.DetectManifestType(entry.Name()) != "" {
			manifests = append(manifests, entry.Name())
		}
	}
	sort.Strings(manifests)
	for _, name := range manifests {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		data, readErr := readLocalFile(scoped, name)
		if readErr != nil {
			continue
		}
		manifest, parseErr := parser.ParseManifest(name, domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID}, string(data))
		if parseErr != nil {
			continue
		}
		if readme, readErr := readLocalFile(scoped, "README.md"); readErr == nil {
			manifest.Content = readme
			manifest.ContentPath = "README.md"
		}
		return manifest, nil
	}
	return nil, domain.ErrManifestNotFound
}

func (c *Client) FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	root, err := candidateRoot(candidate)
	if err != nil {
		return nil, err
	}
	scoped, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	defer scoped.Close()
	nodes := make([]domain.DirectoryNode, 0, 512)
	err = fs.WalkDir(scoped.FS(), ".", func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if path == "." {
			return nil
		}
		if entry.IsDir() && ignoredDirectories[entry.Name()] {
			return fs.SkipDir
		}
		node := domain.DirectoryNode{Path: path, Type: "file"}
		if entry.IsDir() {
			node.Type = "dir"
		} else if info, infoErr := entry.Info(); infoErr == nil {
			node.Size = info.Size()
		}
		nodes = append(nodes, node)
		if len(nodes) >= 20_000 {
			return fs.SkipAll
		}
		return nil
	})
	return nodes, err
}

func (c *Client) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	root, err := candidateRoot(candidate)
	if err != nil {
		return nil, err
	}
	scoped, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	defer scoped.Close()
	extracted := &domain.ExtractedRepository{
		Ref:       domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID},
		Candidate: candidate, Files: make([]domain.ExtractedFile, 0), ConsultadoEm: time.Now().UTC(),
	}
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Dir = root
	if base, err := cmd.Output(); err == nil {
		extracted.BaseCommit = strings.TrimSpace(string(base))
	}
	seen := make(map[string]bool)
	var hashes []string
	for _, relative := range paths {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		_, pathErr := safePath(root, relative)
		name := filepath.Clean(relative)
		if pathErr != nil || seen[name] {
			continue
		}
		seen[name] = true
		data, readErr := readLocalFile(scoped, name)
		if readErr != nil {
			continue
		}
		normalized := filepath.ToSlash(name)
		hash := fmt.Sprintf("%x", sha256.Sum256(data))
		evidence := domain.FileEvidence{SHA256: hash, State: "content_inspected", Method: "leitura local confinada; sem verificação semântica"}
		if len(data) > 0 {
			evidence.FirstLine = 1
			evidence.LastLine = strings.Count(strings.TrimSuffix(string(data), "\n"), "\n") + 1
		}
		extracted.Files = append(extracted.Files, domain.ExtractedFile{Path: normalized, Content: string(data), Size: int64(len(data)), Evidence: evidence})
		hashes = append(hashes, normalized+"\x00"+hash)
	}
	if len(extracted.Files) == 0 {
		return nil, fmt.Errorf("nenhum arquivo local solicitado pôde ser extraído")
	}
	sort.Strings(hashes)
	extracted.SnapshotHash = fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(hashes, "\n"))))
	return extracted, nil
}

func (c *Client) FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	root, err := candidateRoot(candidate)
	if err != nil {
		return candidate, err
	}
	candidate.Description = "Worktree local em " + root
	if candidate.CommitID == "" || strings.EqualFold(candidate.CommitID, "HEAD") || candidate.CommitID == "WORKTREE" {
		cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
		cmd.Dir = root
		output, resolveErr := cmd.Output()
		if resolveErr != nil {
			return candidate, fmt.Errorf("não foi possível resolver commit da worktree local: %w", resolveErr)
		}
		candidate.CommitID = strings.TrimSpace(string(output))
	}
	candidate.ConsultadoEm = time.Now().UTC()
	return candidate, nil
}
