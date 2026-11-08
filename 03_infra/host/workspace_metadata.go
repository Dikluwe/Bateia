/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/workspace_metadata.md
 * @layer L3
 * @updated 2026-10-28
 */

package host

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unicode"
)

// All Git metadata is resolved/read through one retained worktree descriptor.
// Worktrees whose Git administration lives elsewhere are deliberately unknown.
func readLocalGitConfig(worktree string) ([]byte, error) {
	return readLocalGitConfigLimit(worktree, maxIdeaInputBytes)
}

func readLocalGitConfigLimit(worktree string, limit int64) ([]byte, error) {
	if limit <= 0 || limit > maxIdeaInputBytes {
		return nil, fmt.Errorf("invalid Git metadata budget")
	}
	root, err := os.OpenRoot(worktree)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	gitdir := ".git"
	info, err := root.Lstat(gitdir)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("Git metadata symlink rejected")
	}
	if !info.IsDir() {
		data, err := readRootGitMetadata(root, gitdir, 4096)
		if err != nil {
			return nil, err
		}
		pointer := strings.TrimSpace(string(data))
		if !strings.HasPrefix(strings.ToLower(pointer), "gitdir:") {
			return nil, fmt.Errorf("invalid gitdir pointer")
		}
		gitdir, err = confinedGitMetadataPath(".", strings.TrimSpace(pointer[len("gitdir:"):]))
		if err != nil {
			return nil, err
		}
	}
	if err := rejectGitMetadataAliases(root, gitdir); err != nil {
		return nil, err
	}
	gitInfo, err := root.Lstat(gitdir)
	if err != nil || !gitInfo.IsDir() {
		return nil, fmt.Errorf("gitdir is not an internal directory")
	}
	commonPath := filepath.Join(gitdir, "commondir")
	commonInfo, err := root.Lstat(commonPath)
	if err == nil {
		if commonInfo.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("commondir symlink rejected")
		}
		common, err := readRootGitMetadata(root, commonPath, 4096)
		if err != nil {
			return nil, err
		}
		gitdir, err = confinedGitMetadataPath(gitdir, strings.TrimSpace(string(common)))
		if err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return readRootGitMetadata(root, filepath.Join(gitdir, "config"), int(limit))
}

func confinedGitMetadataPath(base, pointer string) (string, error) {
	if pointer == "" || filepath.IsAbs(pointer) || strings.Contains(pointer, "\\") || strings.IndexFunc(pointer, unicode.IsControl) >= 0 {
		return "", fmt.Errorf("external or invalid Git metadata pointer")
	}
	joined := filepath.Clean(filepath.Join(base, pointer))
	if joined == "." || !filepath.IsLocal(joined) {
		return "", fmt.Errorf("Git metadata escaped worktree")
	}
	return joined, nil
}

func rejectGitMetadataAliases(root *os.Root, name string) error {
	if !filepath.IsLocal(name) {
		return fmt.Errorf("Git metadata escaped worktree")
	}
	prefix := ""
	for _, component := range strings.Split(filepath.Clean(name), string(filepath.Separator)) {
		prefix = filepath.Join(prefix, component)
		info, err := root.Lstat(prefix)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("Git metadata symlink rejected")
		}
	}
	return nil
}

func readRootGitMetadata(root *os.Root, name string, limit int) ([]byte, error) {
	if err := rejectGitMetadataAliases(root, name); err != nil {
		return nil, err
	}
	f, err := root.OpenFile(name, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.Mode().IsRegular() || !ok || stat.Nlink != 1 || info.Size() < 0 || info.Size() > int64(limit) {
		return nil, fmt.Errorf("Git metadata requires a bounded single-link regular file")
	}
	data, err := io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > limit {
		return nil, fmt.Errorf("Git metadata exceeded byte budget")
	}
	after, err := f.Stat()
	if err != nil {
		return nil, err
	}
	afterStat, ok := after.Sys().(*syscall.Stat_t)
	if !after.Mode().IsRegular() || !ok || afterStat.Nlink != 1 || !os.SameFile(info, after) || after.Size() != info.Size() || after.Size() != int64(len(data)) || !after.ModTime().Equal(info.ModTime()) {
		return nil, fmt.Errorf("Git metadata changed or acquired another link during read")
	}
	return data, nil
}

func localGitHubRepositoryIdentitiesLimit(candidateURL string, limit int64) []domain.RepositoryIdentity {
	parsed, err := url.Parse(candidateURL)
	if err != nil || !strings.EqualFold(parsed.Scheme, "file") || !filepath.IsAbs(parsed.Path) || (parsed.Host != "" && parsed.Host != "localhost") || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil
	}
	data, err := readLocalGitConfigLimit(filepath.FromSlash(parsed.Path), limit)
	if err != nil {
		return nil
	}
	return usecase.ParseGitRepositoryHints(data)
}

func (*Host) LocalGitRepositoryIdentities(candidateURL string, limit int64) []domain.RepositoryIdentity {
	return localGitHubRepositoryIdentitiesLimit(candidateURL, limit)
}
func (*Host) ReadLocalGitConfig(root string, limit int64) ([]byte, error) {
	return readLocalGitConfigLimit(root, limit)
}
func (*Host) ReadRootGitMetadata(source domain.CommandRoot, name string, limit int) ([]byte, error) {
	r, ok := source.(*root)
	if !ok || r == nil {
		return nil, fmt.Errorf("root was not supplied by the host")
	}
	return readRootGitMetadata(r.Root, name, limit)
}
