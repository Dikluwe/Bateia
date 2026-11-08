/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/exports_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"bateia/01_core/domain"
	"os"
	"path/filepath"
	"testing"
)

func TestExportStoresPreserveInputAndPublishAtomically(t *testing.T) {
	h := New()
	for _, op := range []struct {
		open  func(string, string) (domain.CommandFile, domain.FileInfo, error)
		write func(string, []byte, domain.FileInfo) error
	}{{h.OpenAlgorithmHistoryInput, h.WriteAlgorithmHistoryAtomic}, {h.OpenIdeaRadarExportInput, h.WriteIdeaRadarExportAtomic}} {
		dir := t.TempDir()
		input := filepath.Join(dir, "input")
		output := filepath.Join(dir, "output")
		if err := os.WriteFile(input, []byte("source"), 0600); err != nil {
			t.Fatal(err)
		}
		if f, _, err := op.open(input, input); err == nil {
			f.Close()
			t.Fatal("input/output alias accepted")
		}
		f, info, err := op.open(input, output)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := op.write(output, []byte("result"), info); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(input)
		if err != nil || string(data) != "source" {
			t.Fatal("source overwritten")
		}
		data, err = os.ReadFile(output)
		if err != nil || string(data) != "result" {
			t.Fatal("result missing")
		}
	}
}
