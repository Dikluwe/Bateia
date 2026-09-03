/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/durable_config_test.md
 * @layer L4
 * @updated 2026-09-01
 */

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDurableCursorSecretRejectsHardlinkBeforeChmod(t *testing.T) {
	root := t.TempDir()
	external := filepath.Join(t.TempDir(), "external.key")
	if err := os.WriteFile(external, make([]byte, durableCursorSecretBytes), 0o644); err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(root, "cursor.key")
	if err := os.Link(external, linked); err != nil {
		t.Skipf("hardlinks unavailable: %v", err)
	}
	if _, err := loadOrCreateCursorSecret(linked); err == nil {
		t.Fatal("hardlinked cursor secret was accepted")
	}
	info, err := os.Stat(external)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o644 {
		t.Fatalf("external inode mode changed to %04o", got)
	}
}

func TestDurableConfigRejectsSymlinkRootWithoutExternalMutation(t *testing.T) {
	external := t.TempDir()
	if err := os.Chmod(external, 0o755); err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(t.TempDir(), "linked-root")
	if err := os.Symlink(external, linked); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	t.Setenv("BATEIA_RUN_ROOT", linked)
	if _, err := loadDurableLocalConfig(); err == nil {
		t.Fatal("symlinked durable root was accepted")
	}
	info, err := os.Stat(external)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o755 {
		t.Fatalf("external directory mode changed to %04o", got)
	}
	entries, err := os.ReadDir(external)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("external directory was mutated: %#v", entries)
	}
}
