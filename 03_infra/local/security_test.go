/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/security_test.md
 * @layer L3
 * @updated 2026-09-26
 */
package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"bateia/01_core/domain"
)

func secureLocalCandidate(root string) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{Owner: "local", Repo: "fixture", URL: "file://" + root, CommitID: "WORKTREE"}
}

func TestLocalHardlinksCannotExposeOutsideBytes(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.go")
	sentinel := "synthetic_outside_secret"
	if err := os.WriteFile(outside, []byte("package hidden\n// "+sentinel), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(outside, filepath.Join(root, "linked.go")); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "inside.go", "package inside\n")
	candidate := secureLocalCandidate(root)
	got, err := NewClient().ExtractRepoFiles(context.Background(), candidate, []string{"inside.go", "linked.go"})
	if err != nil || len(got.Files) != 1 || got.Files[0].Path != "inside.go" || strings.Contains(got.Files[0].Content, sentinel) {
		t.Fatalf("hardlink admitted: %+v %v", got, err)
	}
	search, err := NewClient().SearchImplementations(context.Background(), candidate, sentinel, domain.ImplementationSearchBudget{})
	if err != nil || search.FilesRead != 1 || len(search.Hits) != 0 {
		t.Fatalf("hardlink entered implementation parser: %+v %v", search, err)
	}
	for _, input := range search.Examined {
		if input.Path == "linked.go" {
			t.Fatal("hardlink bytes counted as examined")
		}
	}
}

func TestLocalTreeExposesSymlinkHardlinkAndNonregularHonestly(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "inside.go", "package inside\n")
	writeFixture(t, root, "base", "base")
	if err := os.Symlink(filepath.Join(root, "inside.go"), filepath.Join(root, "alias.go")); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(filepath.Join(root, "base"), filepath.Join(root, "hard")); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(root, "pipe"), 0600); err != nil {
		t.Fatal(err)
	}
	nodes, err := NewClient().FetchTree(context.Background(), secureLocalCandidate(root))
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]string{}
	for _, node := range nodes {
		seen[node.Path] = node.Type
	}
	for path, want := range map[string]string{"inside.go": "file", "alias.go": "symlink", "base": "hardlink", "hard": "hardlink", "pipe": "nonregular"} {
		if seen[path] != want {
			t.Fatalf("%s reported as %s, want %s", path, seen[path], want)
		}
	}
}

func TestLocalExtractionFileAndTotalBudgets(t *testing.T) {
	root := t.TempDir()
	candidate := secureLocalCandidate(root)
	if result, err := NewClient().ExtractRepoFiles(context.Background(), candidate, make([]string, maxLocalExtractionFiles+1)); result != nil || err == nil {
		t.Fatal("file count budget ignored")
	}
	paths := []string{}
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("file%d", i)
		writeFixture(t, root, name, strings.Repeat("x", maxLocalFileSize))
		paths = append(paths, name)
	}
	if result, err := NewClient().ExtractRepoFiles(context.Background(), candidate, paths); result != nil || err == nil {
		t.Fatal("total budget ignored")
	}
	if result, err := NewClient().ExtractRepoFiles(context.Background(), candidate, paths[:4]); result == nil || err != nil || len(result.Files) != 4 {
		t.Fatalf("exact budget positive: %+v %v", result, err)
	}
}

func TestLocalBaseMetadataConfinedAndNoGitSubprocess(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "inside.go", "package inside\n")
	writeFixture(t, root, ".git/HEAD", "ref: refs/heads/main\n")
	sha := strings.Repeat("a", 40)
	writeFixture(t, root, ".git/refs/heads/main", sha+"\n")
	shimDir := t.TempDir()
	marker := filepath.Join(shimDir, "called")
	writeFixture(t, shimDir, "git", "#!/bin/sh\nprintf called > '"+marker+"'\nexit 1\n")
	if err := os.Chmod(filepath.Join(shimDir, "git"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", shimDir)
	t.Setenv("GITHUB_TOKEN", "synthetic-only-token")
	t.Setenv("GIT_CONFIG_COUNT", "1")
	t.Setenv("GIT_CONFIG_KEY_0", "alias.rev-parse")
	t.Setenv("GIT_CONFIG_VALUE_0", "!touch "+marker)
	candidate := secureLocalCandidate(root)
	got, err := NewClient().ExtractRepoFiles(context.Background(), candidate, []string{"inside.go"})
	if err != nil || got.BaseCommit != sha {
		t.Fatalf("confined loose refs positive: %+v %v", got, err)
	}
	metadata, err := NewClient().FetchRepoMetadata(context.Background(), candidate)
	if err != nil || metadata.CommitID != sha || !strings.Contains(metadata.Description, "consultiva") {
		t.Fatalf("metadata positive: %+v %v", metadata, err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("Git subprocess was executed")
	}
}

func TestLocalMetadataDetachedAndPackedRefs(t *testing.T) {
	for _, packed := range []bool{false, true} {
		root := t.TempDir()
		sha := strings.Repeat("b", 40)
		if packed {
			writeFixture(t, root, ".git/HEAD", "ref: refs/heads/main\n")
			writeFixture(t, root, ".git/packed-refs", "# packed refs\n"+sha+" refs/heads/main\n")
		} else {
			writeFixture(t, root, ".git/HEAD", sha+"\n")
		}
		got, err := NewClient().FetchRepoMetadata(context.Background(), secureLocalCandidate(root))
		if err != nil || got.CommitID != sha {
			t.Fatalf("metadata positive: %+v %v", got, err)
		}
	}
}

func TestExternalGitdirPointerAndMetadataHardlinkStayUnavailable(t *testing.T) {
	for _, hardlink := range []bool{false, true} {
		root, outside := t.TempDir(), t.TempDir()
		writeFixture(t, root, "inside.go", "package inside\n")
		sha := strings.Repeat("c", 40)
		writeFixture(t, outside, "HEAD", sha+"\n")
		if hardlink {
			if err := os.Mkdir(filepath.Join(root, ".git"), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.Link(filepath.Join(outside, "HEAD"), filepath.Join(root, ".git/HEAD")); err != nil {
				t.Fatal(err)
			}
		} else {
			writeFixture(t, root, ".git", "gitdir: "+outside+"\n")
		}
		candidate := secureLocalCandidate(root)
		got, err := NewClient().ExtractRepoFiles(context.Background(), candidate, []string{"inside.go"})
		if err != nil || got.BaseCommit != "" || len(got.Coverage.ChecksNotRun) != 1 || !strings.Contains(got.Coverage.Reason, "indisponíveis") {
			t.Fatalf("external metadata admitted/hidden: %+v %v", got, err)
		}
		metadata, err := NewClient().FetchRepoMetadata(context.Background(), candidate)
		if err == nil || metadata.CommitID != "WORKTREE" {
			t.Fatalf("external metadata resolved: %+v %v", metadata, err)
		}
	}
}
