/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_git_metadata_test.md
 * @layer L2
 * @updated 2026-09-26
 */

package cli

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"bateia/01_core/domain"
)

const internalRemoteConfig = "[remote \"origin\"]\n url = https://github.com/inside/repo\n"
const externalRemoteConfig = "[remote \"origin\"]\n url = https://github.com/external/repo\n"

func writeGitMetadataFixture(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func localMetadataIdentities(root string) []domain.RepositoryIdentity {
	return localGitHubRepositoryIdentities((&url.URL{Scheme: "file", Path: root}).String())
}

func TestLocalGitMetadataPreservesInternalLayouts(t *testing.T) {
	for _, layout := range []string{"normal", "relative_gitdir", "relative_common"} {
		t.Run(layout, func(t *testing.T) {
			root := t.TempDir()
			switch layout {
			case "normal":
				writeGitMetadataFixture(t, root, ".git/config", internalRemoteConfig)
			case "relative_gitdir":
				writeGitMetadataFixture(t, root, ".git", "gitdir: metadata/state\n")
				writeGitMetadataFixture(t, root, "metadata/state/config", internalRemoteConfig)
			case "relative_common":
				writeGitMetadataFixture(t, root, ".git", "gitdir: metadata/state\n")
				writeGitMetadataFixture(t, root, "metadata/state/commondir", "../common\n")
				writeGitMetadataFixture(t, root, "metadata/common/config", internalRemoteConfig)
			}
			got := localMetadataIdentities(root)
			if len(got) != 1 || got[0].Owner != "inside" || got[0].Repo != "repo" {
				t.Fatalf("internal layout lost identity: %+v", got)
			}
		})
	}
}

func TestLocalGitMetadataRejectsExternalPathsAndAliases(t *testing.T) {
	for _, kind := range []string{"gitdir_absolute", "gitdir_relative_escape", "commondir_absolute", "commondir_relative_escape", "dotgit_symlink", "config_symlink", "config_hardlink", "commondir_symlink", "commondir_hardlink", "gitdir_parent_symlink", "config_fifo"} {
		t.Run(kind, func(t *testing.T) {
			parent := t.TempDir()
			root, external := filepath.Join(parent, "worktree"), filepath.Join(parent, "external")
			if err := os.MkdirAll(root, 0700); err != nil {
				t.Fatal(err)
			}
			writeGitMetadataFixture(t, external, "config", externalRemoteConfig)
			writeGitMetadataFixture(t, external, "common-pointer", "../internal\n")
			switch kind {
			case "gitdir_absolute":
				writeGitMetadataFixture(t, root, ".git", "gitdir: "+external)
			case "gitdir_relative_escape":
				writeGitMetadataFixture(t, root, ".git", "gitdir: ../external")
			case "commondir_absolute":
				writeGitMetadataFixture(t, root, ".git/config", internalRemoteConfig)
				writeGitMetadataFixture(t, root, ".git/commondir", external)
			case "commondir_relative_escape":
				writeGitMetadataFixture(t, root, ".git/config", internalRemoteConfig)
				writeGitMetadataFixture(t, root, ".git/commondir", "../../external")
			case "dotgit_symlink":
				if err := os.Symlink(external, filepath.Join(root, ".git")); err != nil {
					t.Fatal(err)
				}
			case "config_symlink", "config_hardlink", "config_fifo":
				if err := os.Mkdir(filepath.Join(root, ".git"), 0700); err != nil {
					t.Fatal(err)
				}
				var err error
				if kind == "config_symlink" {
					err = os.Symlink(filepath.Join(external, "config"), filepath.Join(root, ".git/config"))
				}
				if kind == "config_hardlink" {
					err = os.Link(filepath.Join(external, "config"), filepath.Join(root, ".git/config"))
				}
				if kind == "config_fifo" {
					err = syscall.Mkfifo(filepath.Join(root, ".git/config"), 0600)
				}
				if err != nil {
					t.Fatal(err)
				}
			case "commondir_symlink", "commondir_hardlink":
				writeGitMetadataFixture(t, root, ".git/config", internalRemoteConfig)
				writeGitMetadataFixture(t, root, "internal/config", internalRemoteConfig)
				var err error
				if kind == "commondir_symlink" {
					err = os.Symlink(filepath.Join(external, "common-pointer"), filepath.Join(root, ".git/commondir"))
				} else {
					err = os.Link(filepath.Join(external, "common-pointer"), filepath.Join(root, ".git/commondir"))
				}
				if err != nil {
					t.Fatal(err)
				}
			case "gitdir_parent_symlink":
				writeGitMetadataFixture(t, root, ".git", "gitdir: metadata")
				if err := os.Symlink(external, filepath.Join(root, "metadata")); err != nil {
					t.Fatal(err)
				}
			}
			if got := localMetadataIdentities(root); got != nil {
				t.Fatalf("external metadata remapped identity: %+v", got)
			}
		})
	}
}

func TestGitMetadataReadRetainsRootDescriptorAfterPathReplacement(t *testing.T) {
	parent := t.TempDir()
	rootPath := filepath.Join(parent, "worktree")
	writeGitMetadataFixture(t, rootPath, ".git/config", internalRemoteConfig)
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if err := os.Rename(rootPath, rootPath+".original"); err != nil {
		t.Fatal(err)
	}
	writeGitMetadataFixture(t, rootPath, ".git/config", externalRemoteConfig)
	data, err := readRootGitMetadata(root, ".git/config", maxIdeaInputBytes)
	if err != nil || string(data) != internalRemoteConfig {
		t.Fatalf("read escaped held root descriptor: %q %v", data, err)
	}
}

func TestLocalGitMetadataRejectsOversizedPointersAndNonlocalFileAuthority(t *testing.T) {
	root := t.TempDir()
	writeGitMetadataFixture(t, root, ".git", "gitdir: "+strings.Repeat("x", 4096))
	if got := localMetadataIdentities(root); got != nil {
		t.Fatal("oversized gitdir admitted")
	}
	root = t.TempDir()
	writeGitMetadataFixture(t, root, ".git/config", internalRemoteConfig)
	if got := localGitHubRepositoryIdentities((&url.URL{Scheme: "file", Host: "remote-server", Path: root}).String()); got != nil {
		t.Fatal("remote file authority treated as local")
	}
}
