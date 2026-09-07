/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/client.md
 * @layer L3
 * @updated 2026-09-07
 */

package local

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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

func (c *Client) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	root, err := candidateRoot(candidate)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(root)
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
		data, readErr := os.ReadFile(filepath.Join(root, name))
		if readErr != nil {
			continue
		}
		manifest, parseErr := parser.ParseManifest(name, domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID}, string(data))
		if parseErr != nil {
			continue
		}
		if readme, readErr := os.ReadFile(filepath.Join(root, "README.md")); readErr == nil {
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
	nodes := make([]domain.DirectoryNode, 0, 512)
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if path == root {
			return nil
		}
		if entry.IsDir() && ignoredDirectories[entry.Name()] {
			return filepath.SkipDir
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return nil
		}
		node := domain.DirectoryNode{Path: filepath.ToSlash(relative), Type: "file"}
		if entry.IsDir() {
			node.Type = "dir"
		} else if info, infoErr := entry.Info(); infoErr == nil {
			node.Size = info.Size()
		}
		nodes = append(nodes, node)
		if len(nodes) >= 20_000 {
			return filepath.SkipAll
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
	for _, relative := range paths {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		path, pathErr := safePath(root, relative)
		if pathErr != nil || seen[path] {
			continue
		}
		seen[path] = true
		info, statErr := os.Lstat(path)
		if statErr != nil || !info.Mode().IsRegular() || info.Size() > 2*1024*1024 {
			continue
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			continue
		}
		extracted.Files = append(extracted.Files, domain.ExtractedFile{Path: filepath.ToSlash(relative), Content: string(data), Size: info.Size()})
	}
	if len(extracted.Files) == 0 {
		return nil, fmt.Errorf("nenhum arquivo local solicitado pôde ser extraído")
	}
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
