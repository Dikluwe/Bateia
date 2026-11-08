/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/host_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"bateia/01_core/domain"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHostConfinesRootsAndPreservesFileOptions(t *testing.T) {
	h := New()
	dir := t.TempDir()
	root, err := h.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	f, err := root.OpenFile("a", domain.OpenWriteOnly|domain.OpenCreate|domain.OpenExclusive|domain.OpenNoFollow|domain.OpenNonblock, 0600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("safe")); err != nil {
		t.Fatal(err)
	}
	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	if n, ok := h.LinkCount(info); !ok || n != 1 {
		t.Fatal("link identity lost")
	}
	if _, err := root.OpenFile("a", domain.OpenWriteOnly|domain.OpenCreate|domain.OpenExclusive, 0600); err == nil {
		t.Fatal("exclusive open overwritten file")
	}
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("outside"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	if f, err := root.OpenFile("link", domain.OpenReadOnly|domain.OpenNoFollow|domain.OpenNonblock, 0); err == nil {
		f.Close()
		t.Fatal("nofollow root followed link")
	}
	if f, err := root.OpenFile("../outside", domain.OpenReadOnly, 0); err == nil {
		f.Close()
		t.Fatal("root escaped")
	}
}
func TestHostOwnsTimeAndDeadlines(t *testing.T) {
	h := New()
	if h.Now().IsZero() || h.Now().Location() != time.UTC {
		t.Fatal("clock not supplied")
	}
	ctx, cancel := h.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, ok := ctx.Deadline(); !ok {
		t.Fatal("deadline missing")
	}
	cancel()
	if ctx.Err() == nil {
		t.Fatal("cancel missing")
	}
}
