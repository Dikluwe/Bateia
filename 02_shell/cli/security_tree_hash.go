/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_tree_hash.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"context"
	"hash"
	"io"
)

const (
	localTreeHashVersion = "tree-sha256-v2"
	maxHashLocalRoots    = 20
	maxHashTreeEntries   = 20000
	maxHashTreeBytes     = 64 << 20
	maxHashTreeDepth     = 64
	maxHashTreePath      = 4096
)

type treeHashBudget = domain.TreeHashBudget

func newTreeHashBudget() *treeHashBudget {
	return &treeHashBudget{EntriesLeft: maxHashTreeEntries, BytesLeft: maxHashTreeBytes}
}
func treeHashExcludingContext(ctx context.Context, path, excluded string, budget *treeHashBudget) (string, error) {
	return commandHost().TreeHash(ctx, path, excluded, budget)
}
func validateHashLayout(root, output string) error {
	return commandHost().ValidateHashLayout(root, output)
}
func hashRootedTree(ctx context.Context, root domain.CommandRoot, excluded string, budget *treeHashBudget) (string, error) {
	return commandHost().HashRootedTree(ctx, root, excluded, budget, nil)
}
func hashRootedTreeTo(ctx context.Context, root domain.CommandRoot, excluded string, budget *treeHashBudget, destination domain.CommandRoot) (string, error) {
	return commandHost().HashRootedTree(ctx, root, excluded, budget, destination)
}
func hashRootedFileTo(ctx context.Context, root domain.CommandRoot, name string, h hash.Hash, buffer []byte, budget *treeHashBudget, destination domain.CommandRoot) error {
	return commandHost().HashRootedFile(ctx, root, name, h, buffer, budget, destination)
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
