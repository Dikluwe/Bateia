/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/workspace_capture_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWorkspaceCaptureOwnsPhysicalLifecycle(t *testing.T) {
	h := New()
	dir := t.TempDir()
	file := filepath.Join(dir, "a")
	if err := os.WriteFile(file, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	capture, err := h.CaptureLocal(context.Background(), []string{dir}, "", false)
	if err != nil {
		t.Fatal(err)
	}
	view := capture.View()
	if len(view.Bindings) != 1 {
		t.Fatal("binding missing")
	}
	if err := os.WriteFile(file, []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(view.Bindings[0].Directory, "a"))
	if err != nil || string(data) != "original" {
		t.Fatal("capture reread origin")
	}
	if err := capture.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(view.ContainerPath); !os.IsNotExist(err) {
		t.Fatal("capture leaked")
	}
}
