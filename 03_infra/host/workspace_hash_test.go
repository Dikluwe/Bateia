/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/workspace_hash_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"bateia/01_core/domain"
	"context"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestWorkspaceHashPreservesSharedBudget(t *testing.T) {
	h := New()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a"), []byte("abc"), 0600); err != nil {
		t.Fatal(err)
	}
	budget := &domain.TreeHashBudget{EntriesLeft: 10, BytesLeft: 3}
	digest, err := h.TreeHash(context.Background(), dir, "", budget)
	if err != nil || digest == "" || budget.BytesLeft != 0 {
		t.Fatalf("hash/budget: %s %+v %v", digest, budget, err)
	}
	if _, err := h.TreeHash(context.Background(), dir, "", budget); err == nil {
		t.Fatal("aggregate quota bypassed")
	}
	if err := syscall.Mkfifo(filepath.Join(dir, "fifo"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := h.TreeHash(context.Background(), dir, "", &domain.TreeHashBudget{EntriesLeft: 10, BytesLeft: 10}); err == nil {
		t.Fatal("FIFO accepted")
	}
}
