/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/local_snapshot_capture_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestLocalCaptureMatchesV2AndIgnoresLaterOriginChanges(t *testing.T) {
	source := t.TempDir()
	output := filepath.Join(source, "dossier")
	putHashFile(t, source, "a/z", "binary\x00content")
	putHashFile(t, source, "a.txt", "A")
	putHashFile(t, source, "README.md", "original documentation")
	putHashFile(t, source, "go.mod", "module captured\n\ngo 1.25\n")
	putHashFile(t, source, ".git", "gitdir: unsafe-external")
	putHashFile(t, source, "target/ignored", "ignored")
	putHashFile(t, source, "dossier/ignored", "ignored")
	if e := os.Mkdir(filepath.Join(source, "empty"), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.Symlink("README.md", filepath.Join(source, "alias")); e != nil {
		t.Fatal(e)
	}
	before, e := treeHashExcluding(source, output)
	if e != nil {
		t.Fatal(e)
	}
	capture, e := captureLocalCandidates(context.Background(), []string{source}, output, false)
	if e != nil {
		t.Fatal(e)
	}
	container := capture.directory
	defer func() {
		if e := capture.Close(); e != nil {
			t.Error(e)
		}
	}()
	if len(capture.Candidates) != 1 || capture.Candidates[0].CommitID != before {
		t.Fatal("capture changed v2 identity")
	}
	copied := capture.Bindings[0].Directory
	if h, e := treeHash(copied); e != nil || h != before {
		t.Fatal("captured bytes differ from their digest", e)
	}
	if info, e := os.Stat(filepath.Join(copied, "README.md")); e != nil || info.Mode().Perm() != 0400 {
		t.Fatal("copy not protected", e)
	}
	for _, name := range []string{"alias", "target", "dossier"} {
		if _, e := os.Lstat(filepath.Join(copied, name)); !os.IsNotExist(e) {
			t.Fatal("excluded name captured", name, e)
		}
	}
	if info, e := os.Stat(filepath.Join(copied, "empty")); e != nil || !info.IsDir() {
		t.Fatal("empty directory not retained")
	}
	putHashFile(t, source, "README.md", "mutated documentation")
	renamed := source + "-renamed"
	if e := os.Rename(source, renamed); e != nil {
		t.Fatal(e)
	}
	// Only this synthetic fixture is removed; the provider must not need it.
	if e := os.RemoveAll(renamed); e != nil {
		t.Fatal(e)
	}
	provider, e := buildSnapshotProvider(nil, capture.Bindings)
	if e != nil {
		t.Fatal(e)
	}
	candidate := capture.Candidates[0]
	extracted, e := provider.ExtractRepoFiles(context.Background(), candidate, []string{"README.md", "a/z"})
	if e != nil {
		t.Fatal(e)
	}
	if extracted.Candidate.URL != candidate.URL || extracted.Ref.Commit != before || extracted.BaseCommit != "" {
		t.Fatal("source identity not preserved")
	}
	if len(extracted.Files) != 2 || extracted.Files[0].Content != "original documentation" {
		t.Fatal("mutated origin was consulted")
	}
	if strings.Contains(extracted.Coverage.Reason, container) || !strings.Contains(extracted.Files[0].Evidence.Method, "captura") {
		t.Fatal("private path leaked or frozen method missing")
	}
	if e := capture.Close(); e != nil {
		t.Fatal(e)
	}
	if _, e := os.Lstat(container); !os.IsNotExist(e) {
		t.Fatal("private capture leaked after close")
	}
}

func TestLocalCaptureFrozenGitHintsAndBoundedConfig(t *testing.T) {
	source := t.TempDir()
	putHashFile(t, source, ".git/config", "[remote \"origin\"]\nurl = https://github.com/owner/old\n")
	putHashFile(t, source, "README.md", "fixture")
	capture, e := captureLocalCandidates(context.Background(), []string{source}, "", true)
	if e != nil {
		t.Fatal(e)
	}
	defer capture.Close()
	candidate := capture.Candidates[0]
	putHashFile(t, source, ".git/config", "[remote \"origin\"]\nurl = https://github.com/owner/new\n")
	if identities := capture.Identities[candidate.URL]; len(identities) != 1 || identities[0].Repo != "old" {
		t.Fatal("captured hints changed")
	}
	lookup := func(url string) []domain.RepositoryIdentity { return capture.Identities[url] }
	discovery := domain.DiscoverySnapshot{Candidates: []domain.DiscoveredCandidate{{Repository: candidate}}}
	rebound, _ := alternativeInputsForTargetUsing(domain.IdeaRunLocal, discovery, nil, &domain.RepositoryIdentity{Owner: "owner", Repo: "old"}, lookup)
	if rebound.Candidates[0].Repository.Owner != "owner" || rebound.Candidates[0].Repository.Repo != "old" {
		t.Fatal("target lookup reread source")
	}
	again, e := captureLocalCandidates(context.Background(), []string{source}, "", true)
	if e != nil {
		t.Fatal(e)
	}
	defer again.Close()
	if again.Candidates[0].CommitID != candidate.CommitID || hashJSON(again.Identities) == hashJSON(capture.Identities) {
		t.Fatal("Git hints conflated with content digest")
	}
	config := filepath.Join(source, ".git", "config")
	if e := os.Truncate(config, maxCapturedGitConfigBytes+1); e != nil {
		t.Fatal(e)
	}
	bounded, e := captureLocalCandidates(context.Background(), []string{source}, "", true)
	if e != nil {
		t.Fatal(e)
	}
	defer bounded.Close()
	if len(bounded.Identities[candidate.URL]) != 0 {
		t.Fatal("oversized config used for identity")
	}
}

func TestLocalCaptureCleansFailures(t *testing.T) {
	for _, kind := range []string{"fifo", "quota", "cancelled", "overlap"} {
		t.Run(kind, func(t *testing.T) {
			source := t.TempDir()
			storage := t.TempDir()
			t.Setenv("TMPDIR", storage)
			ctx := context.Background()
			switch kind {
			case "fifo":
				if e := syscall.Mkfifo(filepath.Join(source, "fifo"), 0600); e != nil {
					t.Fatal(e)
				}
			case "quota":
				for _, name := range []string{"a", "b", "c", "d", "e"} {
					f, e := os.Create(filepath.Join(source, name))
					if e != nil {
						t.Fatal(e)
					}
					if e = f.Truncate(maxIdeaInputBytes); e != nil {
						t.Fatal(e)
					}
					if e = f.Close(); e != nil {
						t.Fatal(e)
					}
				}
			case "cancelled":
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			case "overlap":
				source = storage
			}
			capture, e := captureLocalCandidates(ctx, []string{source}, "", false)
			if e == nil || capture != nil {
				t.Fatal("invalid capture returned")
			}
			entries, e := os.ReadDir(storage)
			if e != nil || len(entries) != 0 {
				t.Fatal("failed capture left private files", e)
			}
		})
	}
}

func TestLocalCaptureSinkFailureNeverReturnsDigest(t *testing.T) {
	source, destination := t.TempDir(), t.TempDir()
	putHashFile(t, source, "a", "A")
	putHashFile(t, destination, "a", "pre-existing")
	src, e := commandHost().OpenRoot(source)
	if e != nil {
		t.Fatal(e)
	}
	defer src.Close()
	dst, e := commandHost().OpenRoot(destination)
	if e != nil {
		t.Fatal(e)
	}
	defer dst.Close()
	if digest, e := hashRootedTreeTo(context.Background(), src, "", newTreeHashBudget(), dst); e == nil || digest != "" {
		t.Fatal("sink failure returned partial hash")
	}
	data, e := os.ReadFile(filepath.Join(destination, "a"))
	if e != nil || string(data) != "pre-existing" {
		t.Fatal("capture overwrote an existing file")
	}
	if e := dst.Close(); e != nil {
		t.Fatal(e)
	}
	if digest, e := hashRootedTreeTo(context.Background(), src, "", newTreeHashBudget(), dst); e == nil || digest != "" {
		t.Fatal("closed destination produced digest")
	}
	// A cancelled stream must fail before touching the destination.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := hashRootedFileTo(ctx, src, "a", sha256.New(), make([]byte, 32<<10), newTreeHashBudget(), dst); !errors.Is(e, context.Canceled) {
		t.Fatal("cancelled capture performed I/O", e)
	}
}

func TestLocalCaptureMetadataNotRequestedDoesNotReadGitConfig(t *testing.T) {
	source := t.TempDir()
	putHashFile(t, source, "README.md", "fixture")
	if e := os.Mkdir(filepath.Join(source, ".git"), 0700); e != nil {
		t.Fatal(e)
	}
	if e := syscall.Mkfifo(filepath.Join(source, ".git", "config"), 0600); e != nil {
		t.Fatal(e)
	}
	capture, e := captureLocalCandidates(context.Background(), []string{source}, "", false)
	if e != nil {
		t.Fatal(e)
	}
	defer capture.Close()
	if len(capture.Identities) != 0 {
		t.Fatal("Git hints collected without target")
	}
}
