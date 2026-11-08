/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/trace_files_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTraceFilesProtectSources(t *testing.T) {
	h := New()
	dir := t.TempDir()
	input := filepath.Join(dir, "input")
	output := filepath.Join(dir, "output")
	if err := os.WriteFile(input, []byte("source"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := h.WriteTraceFile([]byte("changed"), input, input); err == nil {
		t.Fatal("source alias overwritten")
	}
	if err := h.WriteTraceFile([]byte("result"), output, input); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 2 {
		t.Fatalf("temporary leaked: %v %v", entries, err)
	}
	data, err := os.ReadFile(input)
	if err != nil || string(data) != "source" {
		t.Fatal("input changed")
	}
}
