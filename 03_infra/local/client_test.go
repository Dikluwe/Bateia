/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/client_test.md
 * @layer L3
 * @updated 2026-08-23
 */

package local

import (
	"context"
	"os"
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
