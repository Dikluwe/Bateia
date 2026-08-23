/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/app.md
 * @layer L2
 * @updated 2026-08-23
 */

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIRejectsInvalidSearchInputsBeforeExecution(t *testing.T) {
	tests := [][]string{
		{"search", "-query", "test", "-max", "0"},
		{"search", "-query", "test", "-depth", "99"},
		{"search", "-query", "test", "-lang", "unknown"},
		{"search", "-query", "test", "-repos", "invalid"},
	}
	for _, args := range tests {
		err := NewCLIApp(nil).Run(args)
		if err == nil || !strings.Contains(err.Error(), "parâmetros inválidos") {
			t.Fatalf("esperava erro de validação para %v, recebeu %v", args, err)
		}
	}
}

func TestCLIRejectsMissingLocalWorktree(t *testing.T) {
	err := NewCLIApp(nil).Run([]string{"analyze", "-path", filepath.Join(t.TempDir(), "missing")})
	if err == nil || !strings.Contains(err.Error(), "worktree local") {
		t.Fatalf("esperava erro de caminho local, recebeu %v", err)
	}
}
