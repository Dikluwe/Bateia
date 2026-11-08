/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/storage_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestStoragePreservesBoundsAndDossierRules(t *testing.T) {
	h := New()
	dir := t.TempDir()
	input := filepath.Join(dir, "input")
	if err := os.WriteFile(input, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := h.ReadBounded(input, 3); err == nil {
		t.Fatal("oversized input accepted")
	}
	data, err := h.ReadBounded(input, 4)
	if err != nil || string(data) != "data" {
		t.Fatalf("read: %q %v", data, err)
	}
	fifo := filepath.Join(dir, "fifo")
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := h.ReadBounded(fifo, 16); err == nil {
		t.Fatal("nonregular input accepted")
	}
	if err := h.PrepareDossier(dir, false); err == nil {
		t.Fatal("nonempty dossier accepted")
	}
	if err := h.PrepareDossier(dir, true); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(dir, link); err != nil {
		t.Fatal(err)
	}
	if err := h.PrepareDossier(link, true); err == nil {
		t.Fatal("symlink dossier accepted")
	}
	output := filepath.Join(t.TempDir(), "result")
	if err := h.WriteAtomic(output, []byte("original"), 0640); err != nil {
		t.Fatal(err)
	}
	if err := h.WriteAtomic(output, []byte("replacement"), 0600); err != nil {
		t.Fatal(err)
	}
	result, err := os.ReadFile(output)
	if err != nil || string(result) != "replacement" {
		t.Fatalf("output: %q %v", result, err)
	}
	entries, err := os.ReadDir(filepath.Dir(output))
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary leaked: %v %v", entries, err)
	}
}
