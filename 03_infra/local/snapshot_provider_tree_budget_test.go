/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/snapshot_provider_tree_budget_test.md
 * @layer L3
 * @updated 2026-09-26
 */
package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshotProviderPreservesTreeNodeBudget(t *testing.T) {
	directory := t.TempDir()
	candidate := capturedCandidate(filepath.Join(t.TempDir(), "absent-source"))
	p := bindCapture(t, candidate, directory, nil)
	for i := 0; i < maxLocalTreeNodes; i++ {
		if err := os.WriteFile(filepath.Join(directory, fmt.Sprintf("node-%05d", i)), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if tree, err := p.FetchTree(context.Background(), candidate); err != nil || len(tree) != maxLocalTreeNodes {
		t.Fatalf("exact node budget failed: len=%d, err=%v", len(tree), err)
	}
	if err := os.WriteFile(filepath.Join(directory, "one-too-many"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if tree, err := p.FetchTree(context.Background(), candidate); err == nil || tree != nil || strings.Contains(err.Error(), directory) {
		t.Fatalf("node budget bypassed or private path leaked: len=%d, err=%v", len(tree), err)
	}
}
