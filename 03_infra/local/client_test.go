/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/client_test.md
 * @layer L3
 * @updated 2026-09-07
 */

package local

import (
	"bateia/01_core/parser"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func writeFixture(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestClientAnalyzesLocalWorktree(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "go.mod", "module example.local/project\n\ngo 1.25\n")
	writeFixture(t, root, "README.md", "# Local Project")
	writeFixture(t, root, "cmd/main.go", "package main\nfunc main() {}\n")
	writeFixture(t, root, "node_modules/ignored.js", "ignored")
	candidate := domain.RepositoryCandidate{Owner: "local", Repo: "project", URL: "file://" + root, CommitID: "WORKTREE"}
	client := NewClient()

	manifest, err := client.FetchManifest(context.Background(), candidate)
	if err != nil {
		t.Fatalf("manifesto local: %v", err)
	}
	if manifest.ModuleName != "example.local/project" || string(manifest.Content) != "# Local Project" {
		t.Fatalf("manifesto inesperado: %+v", manifest)
	}
	tree, err := client.FetchTree(context.Background(), candidate)
	if err != nil {
		t.Fatalf("árvore local: %v", err)
	}
	for _, node := range tree {
		if node.Path == "node_modules/ignored.js" {
			t.Fatal("node_modules não deveria ser percorrido")
		}
	}
	extracted, err := client.ExtractRepoFiles(context.Background(), candidate, []string{"go.mod", "cmd/main.go"})
	if err != nil || len(extracted.Files) != 2 {
		t.Fatalf("extração local inesperada: files=%d err=%v", len(extracted.Files), err)
	}
}

func TestClientRejectsPathTraversal(t *testing.T) {
	root := t.TempDir()
	candidate := domain.RepositoryCandidate{Owner: "local", Repo: "project", URL: "file://" + root, CommitID: "WORKTREE"}
	_, err := NewClient().ExtractRepoFiles(context.Background(), candidate, []string{"../secret"})
	if err == nil {
		t.Fatal("esperava rejeição de traversal")
	}
}

func TestWorktreeRecordsBaseAndChangedUntrackedBytes(t *testing.T) {
	root := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s %v", args, out, err)
		}
	}
	run("init")
	writeFixture(t, root, "CMakeLists.txt", "project(Coin CXX)\n")
	writeFixture(t, root, "src/render.cpp", "void render() {}\n")
	run("add", ".")
	run("-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-m", "fixture")
	candidate := domain.RepositoryCandidate{Owner: "local", Repo: "coin", URL: "file://" + root, CommitID: "WORKTREE"}
	client := NewClient()
	read := func() *domain.ExtractedRepository {
		t.Helper()
		r, err := client.ExtractRepoFiles(context.Background(), candidate, []string{"src/render.cpp", "new.cpp"})
		if err != nil {
			t.Fatal(err)
		}
		parser.AnnotateEvidence(r, "render", []string{"src/render.cpp", "new.cpp"})
		return r
	}
	first := read()
	writeFixture(t, root, "src/render.cpp", "void render() { changed(); }\n")
	writeFixture(t, root, "new.cpp", "void untracked() {}\n")
	second := read()
	if first.BaseCommit == "" || first.BaseCommit != second.BaseCommit || first.SnapshotHash == second.SnapshotHash || len(second.Files) != 2 || first.Files[0].Evidence.SHA256 == second.Files[0].Evidence.SHA256 {
		t.Fatalf("snapshot did not capture changes: %+v %+v", first, second)
	}
}
