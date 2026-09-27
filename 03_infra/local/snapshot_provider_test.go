/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/snapshot_provider_test.md
 * @layer L3
 * @updated 2026-09-26
 */
package local

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func capturedCandidate(source string) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{Owner: "local", Repo: "captured", URL: "file://" + source, CommitID: "tree-v2-original-digest"}
}

func bindCapture(t *testing.T, candidate domain.RepositoryCandidate, directory string, fallback domain.ContentProvider) *SnapshotProvider {
	t.Helper()
	p, err := NewSnapshotProvider(fallback, []SnapshotBinding{{Candidate: candidate, Directory: directory}})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSnapshotProviderNeverReopensChangedOrDeletedSource(t *testing.T) {
	source, directory := t.TempDir(), t.TempDir()
	for _, root := range []string{source, directory} {
		writeFixture(t, root, "go.mod", "module example.test/captured\n\ngo 1.25\n")
		writeFixture(t, root, "README.md", "Captured readme")
		writeFixture(t, root, "src/save.go", "package captured\nfunc SaveRecord() {}\n")
	}
	candidate := capturedCandidate(source)
	p := bindCapture(t, candidate, directory, nil)
	ctx := context.Background()
	manifestBefore, err := p.FetchManifest(ctx, candidate)
	if err != nil {
		t.Fatal(err)
	}
	treeBefore, err := p.FetchTree(ctx, candidate)
	if err != nil {
		t.Fatal(err)
	}
	extractBefore, err := p.ExtractRepoFiles(ctx, candidate, []string{"src/save.go", "README.md"})
	if err != nil {
		t.Fatal(err)
	}
	searchBefore, err := p.SearchImplementations(ctx, candidate, "SaveRecord", domain.ImplementationSearchBudget{})
	if err != nil || len(searchBefore.Hits) == 0 {
		t.Fatalf("busca inicial: %+v, %v", searchBefore, err)
	}
	writeFixture(t, source, "go.mod", "module malicious.test/replaced\n\ngo 1.25\n")
	writeFixture(t, source, "src/save.go", "package replacement\nfunc Other() {}\n")
	writeFixture(t, source, "extra.txt", "New source content")
	check := func(t *testing.T) {
		t.Helper()
		manifest, err := p.FetchManifest(ctx, candidate)
		if err != nil || manifest.ModuleName != manifestBefore.ModuleName || string(manifest.Content) != string(manifestBefore.Content) || manifest.Ref != manifestBefore.Ref {
			t.Fatalf("manifesto relido da origem: %+v, %v", manifest, err)
		}
		tree, err := p.FetchTree(ctx, candidate)
		if err != nil || !reflect.DeepEqual(tree, treeBefore) {
			t.Fatalf("árvore relida da origem: %+v, %v", tree, err)
		}
		extracted, err := p.ExtractRepoFiles(ctx, candidate, []string{"src/save.go", "README.md"})
		if err != nil || !reflect.DeepEqual(extracted.Files, extractBefore.Files) || extracted.SnapshotHash != extractBefore.SnapshotHash {
			t.Fatalf("extração relida da origem: %+v, %v", extracted, err)
		}
		if extracted.Candidate.URL != candidate.URL || extracted.Candidate.Owner != candidate.Owner || extracted.Candidate.Repo != candidate.Repo || extracted.Candidate.CommitID != candidate.CommitID || extracted.Ref != (domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID}) {
			t.Fatalf("identidade foi alterada: %+v", extracted)
		}
		if extracted.BaseCommit != "" || !strings.Contains(extracted.Coverage.Reason, "Git não capturados") || !reflect.DeepEqual(extracted.Coverage.ChecksNotRun, []string{"git_base_metadata"}) {
			t.Fatalf("metadados não comprovados: %+v", extracted)
		}
		metadata, err := p.FetchRepoMetadata(ctx, candidate)
		if err != nil || metadata.URL != candidate.URL || metadata.CommitID != candidate.CommitID || metadata.Description != snapshotDescription {
			t.Fatalf("metadata indevida: %+v, %v", metadata, err)
		}
		search, err := p.SearchImplementations(ctx, candidate, "SaveRecord", domain.ImplementationSearchBudget{})
		if err != nil || !reflect.DeepEqual(search, searchBefore) {
			t.Fatalf("busca relida da origem: %+v, %v", search, err)
		}
		payload, _ := json.Marshal([]any{manifest, tree, extracted, metadata, search})
		if strings.Contains(string(payload), directory) {
			t.Fatalf("caminho privado vazou: %s", payload)
		}
	}
	t.Run("source changed", check)
	// Move the fixture out of its original location without deleting user data.
	if err := os.Rename(source, filepath.Join(t.TempDir(), "removed-source")); err != nil {
		t.Fatal(err)
	}
	t.Run("source no longer exists", check)
	if _, err := NewSnapshotProvider(nil, []SnapshotBinding{{Candidate: candidate, Directory: directory}}); err != nil {
		t.Fatalf("constructor consulted absent origin: %v", err)
	}
}

func TestSnapshotProviderBindingAdmission(t *testing.T) {
	directory := t.TempDir()
	valid := SnapshotBinding{Candidate: capturedCandidate(filepath.Join(t.TempDir(), "absent-source")), Directory: directory}
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(directory, alias); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, []byte("not directory"), 0600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		edit func(*SnapshotBinding)
	}{
		{"empty owner", func(b *SnapshotBinding) { b.Candidate.Owner = "" }},
		{"empty repo", func(b *SnapshotBinding) { b.Candidate.Repo = "" }},
		{"empty revision", func(b *SnapshotBinding) { b.Candidate.CommitID = "" }},
		{"whitespace revision", func(b *SnapshotBinding) { b.Candidate.CommitID = " " }},
		{"remote", func(b *SnapshotBinding) { b.Candidate.URL = "https://example.invalid/repo" }},
		{"relative source", func(b *SnapshotBinding) { b.Candidate.URL = "file://relative" }},
		{"relative directory", func(b *SnapshotBinding) { b.Directory = "relative" }},
		{"unclean directory", func(b *SnapshotBinding) { b.Directory += "/." }},
		{"absent directory", func(b *SnapshotBinding) { b.Directory = filepath.Join(directory, "absent") }},
		{"file directory", func(b *SnapshotBinding) { b.Directory = file }},
		{"alias", func(b *SnapshotBinding) { b.Directory = alias }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			binding := valid
			test.edit(&binding)
			if _, err := NewSnapshotProvider(nil, []SnapshotBinding{binding}); err == nil {
				t.Fatal("invalid binding admitted")
			} else if strings.Contains(err.Error(), directory) || strings.Contains(err.Error(), alias) || strings.Contains(err.Error(), file) {
				t.Fatalf("private path exposed: %v", err)
			}
		})
	}
	for _, kind := range []string{"exact duplicate", "same repo new URL", "same URL new repo", "same directory"} {
		t.Run(kind, func(t *testing.T) {
			second := valid
			if kind != "exact duplicate" {
				second.Directory = t.TempDir()
			}
			switch kind {
			case "same repo new URL":
				second.Candidate.URL += "/other"
				second.Candidate.CommitID += "-other"
			case "same URL new repo":
				second.Candidate.Repo += "-other"
			case "same directory":
				second.Directory = valid.Directory
				second.Candidate.Repo += "-other"
				second.Candidate.URL += "/other"
			}
			if _, err := NewSnapshotProvider(nil, []SnapshotBinding{valid, second}); err == nil {
				t.Fatal("ambiguous bindings admitted")
			}
		})
	}
	if _, err := NewSnapshotProvider(nil, make([]SnapshotBinding, maxSnapshotBindings+1)); err == nil || !strings.Contains(err.Error(), "limite") {
		t.Fatalf("binding budget not applied before validation: %v", err)
	}
}

type snapshotFallback struct{ calls int }

type snapshotContentOnly struct{ domain.ContentProvider }

func (f *snapshotFallback) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	f.calls++
	return &domain.Manifest{ModuleName: "remote"}, nil
}
func (f *snapshotFallback) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	f.calls++
	return []domain.DirectoryNode{{Path: "remote"}}, nil
}
func (f *snapshotFallback) ExtractRepoFiles(_ context.Context, candidate domain.RepositoryCandidate, _ []string) (*domain.ExtractedRepository, error) {
	f.calls++
	return &domain.ExtractedRepository{Candidate: candidate}, nil
}
func (f *snapshotFallback) FetchRepoMetadata(_ context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	f.calls++
	return candidate, nil
}
func (f *snapshotFallback) SearchImplementations(context.Context, domain.RepositoryCandidate, string, domain.ImplementationSearchBudget) (domain.ImplementationSearchReport, error) {
	f.calls++
	return domain.ImplementationSearchReport{Status: "remote"}, nil
}

func snapshotOperations(p *SnapshotProvider, ctx context.Context, candidate domain.RepositoryCandidate) []error {
	_, manifestErr := p.FetchManifest(ctx, candidate)
	_, treeErr := p.FetchTree(ctx, candidate)
	_, extractErr := p.ExtractRepoFiles(ctx, candidate, []string{"go.mod"})
	_, metadataErr := p.FetchRepoMetadata(ctx, candidate)
	_, searchErr := p.SearchImplementations(ctx, candidate, "SaveRecord", domain.ImplementationSearchBudget{})
	return []error{manifestErr, treeErr, extractErr, metadataErr, searchErr}
}

func TestSnapshotProviderNeverFallsBackForLocalMismatch(t *testing.T) {
	directory, source := t.TempDir(), t.TempDir()
	writeFixture(t, source, "go.mod", "module example.test/live\n\ngo 1.25\n")
	candidate := capturedCandidate(source)
	fallback := &snapshotFallback{}
	p := bindCapture(t, candidate, directory, fallback)
	for _, field := range []string{"owner", "repo", "url", "revision", "HEAD", "empty"} {
		t.Run(field, func(t *testing.T) {
			changed := candidate
			switch field {
			case "owner":
				changed.Owner += "-changed"
			case "repo":
				changed.Repo += "-changed"
			case "url":
				changed.URL += "/changed"
			case "revision":
				changed.CommitID += "-changed"
			case "HEAD":
				changed.CommitID = "HEAD"
			case "empty":
				changed.CommitID = ""
			}
			for _, err := range snapshotOperations(p, context.Background(), changed) {
				if err == nil || !strings.Contains(err.Error(), "não vinculada") {
					t.Fatalf("local mismatch did not fail closed: %v", err)
				}
			}
		})
	}
	if fallback.calls != 0 {
		t.Fatalf("local mismatch reached fallback %d times", fallback.calls)
	}
	remote := domain.RepositoryCandidate{URL: "https://example.invalid/repo", CommitID: "remote-commit"}
	for _, err := range snapshotOperations(p, context.Background(), remote) {
		if err != nil {
			t.Fatal(err)
		}
	}
	if fallback.calls != 5 {
		t.Fatalf("remote calls not delegated: %d", fallback.calls)
	}
	withoutSearch, err := NewSnapshotProvider(snapshotContentOnly{fallback}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := withoutSearch.SearchImplementations(context.Background(), remote, "save", domain.ImplementationSearchBudget{}); err == nil || !strings.Contains(err.Error(), "não oferece") {
		t.Fatalf("missing search interface not explicit: %v", err)
	}
	withoutFallback := bindCapture(t, candidate, directory, nil)
	for _, err := range snapshotOperations(withoutFallback, context.Background(), remote) {
		if err == nil || !strings.Contains(err.Error(), "remoto indisponível") {
			t.Fatalf("missing fallback not explicit: %v", err)
		}
	}
}

func TestSnapshotProviderSanitizesErrorsAndDoesNotInventGit(t *testing.T) {
	directory := t.TempDir()
	candidate := capturedCandidate(filepath.Join(t.TempDir(), "absent-source"))
	p := bindCapture(t, candidate, directory, nil)
	writeFixture(t, directory, "go.mod", "module example.test/captured\n\ngo 1.25\n")
	writeFixture(t, directory, ".git/HEAD", strings.Repeat("a", 40))
	extracted, err := p.ExtractRepoFiles(context.Background(), candidate, []string{"go.mod"})
	if err != nil || extracted.BaseCommit != "" || !strings.Contains(extracted.Coverage.Reason, "indisponíveis") {
		t.Fatalf("Git metadata was transferred: %+v, %v", extracted, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, err := range snapshotOperations(p, ctx, candidate) {
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation identity lost: %v", err)
		}
	}
	// Remove access by moving only this synthetic capture; errors must not expose its path.
	if err := os.Rename(directory, filepath.Join(t.TempDir(), "removed-capture")); err != nil {
		t.Fatal(err)
	}
	for _, err := range snapshotOperations(p, context.Background(), candidate) {
		// Metadata is the recorded binding, and does not require reading capture bytes.
		if err != nil && strings.Contains(err.Error(), directory) {
			t.Fatalf("private path leaked: %v", err)
		}
	}
	report, err := p.SearchImplementations(context.Background(), candidate, "save", domain.ImplementationSearchBudget{})
	if err == nil || strings.Contains(strings.Join(report.Issues, " "), directory) {
		t.Fatalf("search error/report leaks capture: %+v, %v", report, err)
	}
}

func TestSnapshotProviderPreservesClientBudgets(t *testing.T) {
	directory := t.TempDir()
	candidate := capturedCandidate(filepath.Join(t.TempDir(), "absent-source"))
	p := bindCapture(t, candidate, directory, nil)
	paths := make([]string, maxLocalExtractionFiles+1)
	if _, err := p.ExtractRepoFiles(context.Background(), candidate, paths); err == nil {
		t.Fatal("extraction probe budget bypassed")
	}
	paths = nil
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("chunk-%d.txt", i)
		writeFixture(t, directory, name, strings.Repeat("x", maxLocalFileSize))
		paths = append(paths, name)
	}
	if extracted, err := p.ExtractRepoFiles(context.Background(), candidate, paths[:4]); err != nil || len(extracted.Files) != 4 {
		t.Fatalf("exact byte budget failed: %+v, %v", extracted, err)
	}
	if _, err := p.ExtractRepoFiles(context.Background(), candidate, paths); err == nil {
		t.Fatal("aggregate extraction budget bypassed")
	}
	for i := 0; i < 3; i++ {
		writeFixture(t, directory, fmt.Sprintf("src/file%d.go", i), "package captured\nfunc SaveRecord() {}\n")
	}
	budget := domain.ImplementationSearchBudget{MaxFiles: 1, MaxBytes: 128, MaxFileBytes: 64, MaxResults: 1}
	report, err := p.SearchImplementations(context.Background(), candidate, "SaveRecord", budget)
	if err != nil || report.FilesRead > 1 || report.BytesRead > 128 || len(report.Hits) > 1 || !report.Limited || report.Budget != budget {
		t.Fatalf("search budget bypassed: %+v, %v", report, err)
	}
}

func TestSnapshotProviderPreservesManifestNotFound(t *testing.T) {
	directory := t.TempDir()
	candidate := capturedCandidate(filepath.Join(t.TempDir(), "absent-source"))
	p := bindCapture(t, candidate, directory, nil)
	if _, err := p.FetchManifest(context.Background(), candidate); !errors.Is(err, domain.ErrManifestNotFound) {
		t.Fatalf("manifest sentinel lost: %v", err)
	}
}
