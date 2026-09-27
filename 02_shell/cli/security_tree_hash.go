/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_tree_hash.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

const (
	localTreeHashVersion = "tree-sha256-v2"
	maxHashLocalRoots    = 20
	maxHashTreeEntries   = 20000
	maxHashTreeBytes     = 64 << 20
	maxHashTreeDepth     = 64
	maxHashTreePath      = 4096
)

// One budget is retained across all --local roots. It is private host policy,
// not adjustable by repository input. Counts include skipped entry names, but
// ignored/excluded directory contents are never enumerated.
type treeHashBudget struct {
	entriesLeft int
	bytesLeft   int64
}

func newTreeHashBudget() *treeHashBudget {
	return &treeHashBudget{entriesLeft: maxHashTreeEntries, bytesLeft: maxHashTreeBytes}
}

func treeHashExcludingContext(ctx context.Context, rootPath, excludedPath string, budget *treeHashBudget) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if len(rootPath) > maxHashTreePath || len(excludedPath) > maxHashTreePath {
		return "", fmt.Errorf("caminho de hash excede limite de 4096 bytes")
	}
	abs, err := filepath.Abs(rootPath)
	if err != nil {
		return "", err
	}
	if len(abs) > maxHashTreePath {
		return "", fmt.Errorf("raiz de hash excede limite de caminho")
	}
	if err := validateHashLayout(abs, excludedPath); err != nil {
		return "", err
	}
	excluded := ""
	if strings.TrimSpace(excludedPath) != "" {
		excludedAbs, err := filepath.Abs(excludedPath)
		if err != nil {
			return "", fmt.Errorf("caminho de saída inválido: %w", err)
		}
		rel, err := filepath.Rel(abs, excludedAbs)
		if err != nil {
			return "", err
		}
		if rel == "." || filepath.IsLocal(rel) {
			excluded = filepath.ToSlash(rel)
		}
	}
	root, err := os.OpenRoot(abs)
	if err != nil {
		return "", err
	}
	defer root.Close()
	return hashRootedTree(ctx, root, excluded, budget)
}

// Admission rejects lexical aliases so output exclusion cannot silently miss
// an existing physical directory. Rooted reads separately confine races.
func validateHashLayout(rootPath, outputPath string) error {
	if len(rootPath) > maxHashTreePath || len(outputPath) > maxHashTreePath {
		return fmt.Errorf("caminho de hash excede limite de 4096 bytes")
	}
	if err := rejectHashPathAliases(rootPath, false); err != nil {
		return err
	}
	if strings.TrimSpace(outputPath) == "" {
		return nil
	}
	output, err := filepath.Abs(outputPath)
	if err != nil {
		return err
	}
	if len(output) > maxHashTreePath {
		return fmt.Errorf("caminho de saída excede limite")
	}
	if err := rejectHashPathAliases(output, true); err != nil {
		return err
	}
	rel, err := filepath.Rel(output, rootPath)
	if err != nil {
		return err
	}
	if rel == "." || filepath.IsLocal(rel) {
		return fmt.Errorf("saída não pode ser igual ou ancestral da raiz local")
	}
	return nil
}

func rejectHashPathAliases(path string, allowMissing bool) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// Inspect ancestors even when the leaf does not yet exist.
	for current := abs; ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil {
			if !allowMissing || !os.IsNotExist(err) {
				return err
			}
		} else if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return fmt.Errorf("raiz/saída de hash exige diretórios sem aliases symlink")
		}
		if filepath.Dir(current) == current {
			break
		}
	}
	return nil
}

// Preserve filepath.WalkDir's lexical DFS order, not global full-path order.
// A directory 'a' is visited before sibling 'a.go', even though 'a.go' sorts
// before 'a/x' as a complete string. V2 frames path and content with lengths;
// the unframed legacy preimage allowed structural collisions between trees.
func hashRootedTree(ctx context.Context, root *os.Root, excluded string, budget *treeHashBudget) (string, error) {
	return hashRootedTreeTo(ctx, root, excluded, budget, nil)
}

func hashRootedTreeTo(ctx context.Context, root *os.Root, excluded string, budget *treeHashBudget, destination *os.Root) (string, error) {
	if budget == nil {
		return "", fmt.Errorf("hash exige orçamento compartilhado")
	}
	h := sha256.New()
	_, _ = h.Write([]byte("bateia-local-tree-v2\x00"))
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if excluded == "." {
		_, _ = h.Write([]byte{0})
		return hex.EncodeToString(h.Sum(nil)), nil
	}
	buffer := make([]byte, 32<<10)
	var walk func(string, int) error
	walk = func(directory string, depth int) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		entries, err := readHashDirectory(ctx, root, directory, budget)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := ctx.Err(); err != nil {
				return err
			}
			name := filepath.Join(directory, entry.Name())
			if len(name) > maxHashTreePath {
				return fmt.Errorf("caminho da árvore excede limite de 4096 bytes")
			}
			slash := filepath.ToSlash(name)
			if slash == excluded {
				continue
			}
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}
			if entry.IsDir() {
				if entry.Name() == ".git" || entry.Name() == "target" || entry.Name() == "node_modules" {
					continue
				}
				if depth >= maxHashTreeDepth {
					return fmt.Errorf("árvore de hash excede profundidade de 64 diretórios")
				}
				if destination != nil {
					if err := destination.MkdirAll(name, 0700); err != nil {
						return fmt.Errorf("não foi possível capturar diretório local")
					}
				}
				if err := walk(name, depth+1); err != nil {
					return err
				}
				continue
			}
			if err := hashRootedFileTo(ctx, root, name, h, buffer, budget, destination); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk(".", 0); err != nil {
		return "", err
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	_, _ = h.Write([]byte{0})
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Read batches before sorting, so a hostile flat directory cannot allocate an
// unbounded entry list as filepath.WalkDir would. The total retained entry
// count across DFS frames is bounded by the same shared budget.
func readHashDirectory(ctx context.Context, root *os.Root, name string, budget *treeHashBudget) ([]os.DirEntry, error) {
	dir, err := root.OpenFile(name, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	var entries []os.DirEntry
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		batch, readErr := dir.ReadDir(128)
		if len(batch) > budget.entriesLeft {
			return nil, fmt.Errorf("hash excede orçamento compartilhado de 20000 entradas")
		}
		budget.entriesLeft -= len(batch)
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

func hashRootedFile(ctx context.Context, root *os.Root, name string, h hash.Hash, buffer []byte, budget *treeHashBudget) error {
	return hashRootedFileTo(ctx, root, name, h, buffer, budget, nil)
}

func hashRootedFileTo(ctx context.Context, root *os.Root, name string, h hash.Hash, buffer []byte, budget *treeHashBudget, destination *os.Root) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f, err := root.OpenFile(name, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.Mode().IsRegular() || !ok || stat.Nlink != 1 || info.Size() < 0 {
		return fmt.Errorf("hash exige arquivo regular sem hardlinks")
	}
	limit := int64(maxIdeaInputBytes)
	if budget.bytesLeft < limit {
		limit = budget.bytesLeft
	}
	if info.Size() > limit {
		return fmt.Errorf("arquivo excede limite de 16 MiB ou orçamento compartilhado de 64 MiB")
	}
	var copyFile *os.File
	var contentWriter io.Writer = h
	if destination != nil {
		copyFile, err = destination.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0600)
		if err != nil {
			return fmt.Errorf("não foi possível criar arquivo de captura local")
		}
		defer copyFile.Close()
		contentWriter = io.MultiWriter(h, copyFile)
	}
	path := []byte(filepath.ToSlash(name))
	_, _ = h.Write([]byte{1})
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(path)))
	_, _ = h.Write(length[:])
	_, _ = h.Write(path)
	binary.BigEndian.PutUint64(length[:], uint64(info.Size()))
	_, _ = h.Write(length[:])
	reader := &hashContextReader{ctx: ctx, reader: io.LimitReader(f, limit+1)}
	n, err := io.CopyBuffer(contentWriter, reader, buffer)
	budget.bytesLeft -= n
	if err != nil {
		return err
	}
	if n > limit {
		return fmt.Errorf("conteúdo cresceu além do orçamento de hash")
	}
	after, err := f.Stat()
	if err != nil {
		return err
	}
	afterStat, ok := after.Sys().(*syscall.Stat_t)
	if !after.Mode().IsRegular() || !ok || afterStat.Nlink != 1 || !os.SameFile(info, after) || after.Size() != info.Size() || after.Size() != n || !after.ModTime().Equal(info.ModTime()) {
		return fmt.Errorf("arquivo alterado durante cálculo de hash")
	}
	if copyFile != nil {
		if err := copyFile.Chmod(0400); err != nil {
			return fmt.Errorf("não foi possível proteger arquivo capturado")
		}
		if err := copyFile.Close(); err != nil {
			return fmt.Errorf("não foi possível concluir captura local")
		}
	}
	return ctx.Err()
}

type hashContextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *hashContextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}
