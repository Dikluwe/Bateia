/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_dossier_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestPrepareDossierBoundedAdmission(t *testing.T) {
	root := t.TempDir()
	empty := filepath.Join(root, "empty")
	if e := prepareDossier(empty, false); e != nil {
		t.Fatal(e)
	}
	if e := prepareDossier(empty, false); e != nil {
		t.Fatal("existing empty refused", e)
	}
	putHashFile(t, empty, "evidence.json", "fixture")
	if e := prepareDossier(empty, false); e == nil {
		t.Fatal("nonempty accepted without resume")
	}
	if e := prepareDossier(empty, true); e != nil {
		t.Fatal("resume refused", e)
	}
	if b, e := os.ReadFile(filepath.Join(empty, "evidence.json")); e != nil || string(b) != "fixture" {
		t.Fatal("admission changed evidence")
	}
	file := filepath.Join(root, "file")
	putHashFile(t, root, "file", "fixture")
	fifo := filepath.Join(root, "fifo")
	if e := syscall.Mkfifo(fifo, 0600); e != nil {
		t.Fatal(e)
	}
	alias := filepath.Join(root, "alias")
	if e := os.Symlink(empty, alias); e != nil {
		t.Fatal(e)
	}
	for _, path := range []string{file, fifo, alias} {
		for _, resume := range []bool{true, false} {
			if e := prepareDossier(path, resume); e == nil {
				t.Fatal("unsafe dossier admitted", path, resume)
			}
		}
	}
}
