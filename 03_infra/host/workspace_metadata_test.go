/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/workspace_metadata_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkspaceMetadataUsesInternalConfigOnly(t *testing.T) {
	h := New()
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0700); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, ".git", "config")
	if err := os.WriteFile(config, []byte("[remote \"origin\"]\nurl = https://github.com/org/repo\n"), 0600); err != nil {
		t.Fatal(err)
	}
	hints := h.LocalGitRepositoryIdentities("file://"+dir, 1024)
	if len(hints) != 1 || hints[0].Repo != "repo" {
		t.Fatal(hints)
	}
	if err := os.Remove(config); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(outside, []byte("outside"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, config); err != nil {
		t.Fatal(err)
	}
	if _, err := h.ReadLocalGitConfig(dir, 1024); err == nil {
		t.Fatal("external config followed")
	}
}
