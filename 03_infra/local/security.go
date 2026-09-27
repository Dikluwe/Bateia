/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/security.md
 * @layer L3
 * @updated 2026-09-26
 */
package local

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

const (
	maxLocalTreeNodes       = 20_000
	maxLocalExtractionFiles = 100
	maxLocalExtractionBytes = 8 << 20
)

func singleLinkRegular(info os.FileInfo) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	return info.Mode().IsRegular() && ok && stat.Nlink == 1
}

// Stream directory entries instead of fs.WalkDir/fs.ReadDir, which allocate an
// entire attacker-sized directory before a callback can enforce its budget.
func walkLocalTree(ctx context.Context, root *os.Root, visit func(string, fs.DirEntry) error) error {
	dirs := []string{"."}
	count := 0
	for cursor := 0; cursor < len(dirs); cursor++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		name := dirs[cursor]
		dir, err := root.OpenFile(name, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
		if err != nil {
			return err
		}
		err = func() error {
			defer dir.Close()
			for {
				if err := ctx.Err(); err != nil {
					return err
				}
				entries, readErr := dir.ReadDir(128)
				for _, entry := range entries {
					count++
					if count > maxLocalTreeNodes {
						return fmt.Errorf("árvore local excede limite de %d entradas; resultado incompleto", maxLocalTreeNodes)
					}
					path := filepath.Join(name, entry.Name())
					if len(path) > 4096 {
						return fmt.Errorf("caminho de árvore local excede limite de 4096 bytes")
					}
					visitErr := visit(filepath.ToSlash(path), entry)
					if visitErr == fs.SkipAll {
						return fs.SkipAll
					}
					if visitErr == fs.SkipDir {
						continue
					}
					if visitErr != nil {
						return visitErr
					}
					if entry.IsDir() {
						dirs = append(dirs, path)
					}
				}
				if readErr == io.EOF {
					return nil
				}
				if readErr != nil {
					return readErr
				}
			}
		}()
		if err == fs.SkipAll {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func readLocalRootEntries(ctx context.Context, root *os.Root) ([]fs.DirEntry, error) {
	dir, err := root.Open(".")
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	var entries []fs.DirEntry
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		batch, readErr := dir.ReadDir(128)
		if len(entries)+len(batch) > maxLocalTreeNodes {
			return nil, fmt.Errorf("raiz local excede limite de entradas")
		}
		entries = append(entries, batch...)
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, readErr
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, nil
}

func localSHA(value string) bool {
	decoded, err := hex.DecodeString(value)
	return len(value) == 40 && err == nil && len(decoded) == 20
}

func metadataRelative(path string) bool {
	return filepath.IsLocal(path) && filepath.Clean(path) == path && !strings.ContainsAny(path, "\\\r\n\x00")
}

// Metadata is advisory only, not proof that worktree bytes match this commit.
// A worktree .git pointer is accepted only if its gitdir is inside this Root.
func readLocalBaseCommit(root *os.Root) (string, error) {
	gitdir := ".git"
	info, err := root.Lstat(gitdir)
	if err != nil {
		return "", fmt.Errorf("metadados Git locais indisponíveis: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("metadados Git por symlink indisponíveis")
	}
	if !info.IsDir() {
		data, err := readLocalFileLimit(root, gitdir, 4096)
		if err != nil {
			return "", err
		}
		if len(data) > 4096 {
			return "", fmt.Errorf("ponteiro gitdir excede limite")
		}
		pointer := strings.TrimSpace(string(data))
		if !strings.HasPrefix(pointer, "gitdir: ") {
			return "", fmt.Errorf("ponteiro gitdir inválido")
		}
		gitdir = strings.TrimPrefix(pointer, "gitdir: ")
		if !metadataRelative(gitdir) {
			return "", fmt.Errorf("gitdir externo ou inválido: metadados indisponíveis sem leitura externa")
		}
	}
	head, err := readLocalFileLimit(root, filepath.Join(gitdir, "HEAD"), 4096)
	if err != nil {
		return "", fmt.Errorf("HEAD local indisponível: %w", err)
	}
	if len(head) > 4096 {
		return "", fmt.Errorf("HEAD excede limite de metadados")
	}
	value := strings.TrimSpace(string(head))
	if localSHA(value) {
		return value, nil
	}
	if !strings.HasPrefix(value, "ref: ") {
		return "", fmt.Errorf("HEAD local inválido")
	}
	ref := strings.TrimPrefix(value, "ref: ")
	if !strings.HasPrefix(ref, "refs/") || !metadataRelative(ref) {
		return "", fmt.Errorf("referência HEAD local insegura")
	}
	loose, looseErr := readLocalFileLimit(root, filepath.Join(gitdir, ref), 4096)
	if looseErr == nil {
		value := strings.TrimSpace(string(loose))
		if len(loose) <= 4096 && localSHA(value) {
			return value, nil
		}
		return "", fmt.Errorf("referência solta local inválida")
	}
	if !os.IsNotExist(looseErr) {
		return "", fmt.Errorf("referência solta local indisponível: %w", looseErr)
	}
	packed, err := readLocalFile(root, filepath.Join(gitdir, "packed-refs"))
	if err != nil {
		return "", fmt.Errorf("packed-refs local indisponível: %w", err)
	}
	if strings.Count(string(packed), "\n") >= maxLocalTreeNodes {
		return "", fmt.Errorf("packed-refs excede limite de entradas")
	}
	lines := strings.Split(string(packed), "\n")
	if len(lines) > maxLocalTreeNodes {
		return "", fmt.Errorf("packed-refs excede limite de entradas")
	}
	resolved := ""
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == ref {
			if resolved != "" || !localSHA(fields[0]) {
				return "", fmt.Errorf("packed-refs ambíguo ou inválido")
			}
			resolved = fields[0]
		}
	}
	if resolved == "" {
		return "", fmt.Errorf("referência local não encontrada")
	}
	return resolved, nil
}
