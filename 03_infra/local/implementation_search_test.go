/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/implementation_search_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package local_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/local"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func implementationFixture(t *testing.T, root, path, body string) {
	t.Helper()
	full := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}
func implementationCandidate(root string) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{Owner: "fixture", Repo: "sample", URL: "file://" + root, CommitID: strings.Repeat("a", 40)}
}
func implementationBudget() domain.ImplementationSearchBudget {
	return domain.ImplementationSearchBudget{MaxFiles: 50, MaxBytes: 100000, MaxFileBytes: 10000, MaxResults: 10}
}

func TestLocalImplementationSearchContainmentAndAccounting(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	source := "package sample\nfunc RetrieveWidget() {}\n"
	implementationFixture(t, root, "neutral.go", source)
	for _, path := range []string{"testdata/fixture.go", "oracle/hidden.go", "oracles/hidden.go", "gates/hidden.go", ".git/hidden.go", ".agents/hidden.go", ".codex/hidden.go", "vendor/hidden.go", "neutral_test.go"} {
		implementationFixture(t, root, path, "package sample\nfunc RetrieveWidget() {}\n")
	}
	implementationFixture(t, outside, "outside.go", "package sample\nfunc ExternalSecret() {}\n")
	if err := os.Symlink(filepath.Join(outside, "outside.go"), filepath.Join(root, "escape.go")); err != nil {
		t.Fatal(err)
	}
	r, err := local.NewClient().SearchImplementations(context.Background(), implementationCandidate(root), "RetrieveWidget", implementationBudget())
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Hits) != 1 || r.Hits[0].Path != "neutral.go" {
		t.Fatalf("protected/ordinary test input leaked or target lost: %+v", r)
	}
	if r.FilesRead != 1 || r.BytesRead != int64(len(source)) || len(r.Examined) != 1 || r.Examined[0].SHA256 != fmt.Sprintf("%x", sha256.Sum256([]byte(source))) {
		t.Fatalf("unexpected bytes read: %+v", r)
	}
}

func TestLocalImplementationSearchBudgets(t *testing.T) {
	root := t.TempDir()
	source := "package sample\nfunc QueueItems() {}\n"
	implementationFixture(t, root, "a.go", source)
	implementationFixture(t, root, "b.go", source)
	for _, tc := range []struct {
		name   string
		budget domain.ImplementationSearchBudget
	}{
		{"files", domain.ImplementationSearchBudget{MaxFiles: 1, MaxBytes: 10000, MaxFileBytes: 1000, MaxResults: 10}},
		{"bytes", domain.ImplementationSearchBudget{MaxFiles: 10, MaxBytes: int64(len(source)), MaxFileBytes: 1000, MaxResults: 10}},
		{"per_file", domain.ImplementationSearchBudget{MaxFiles: 10, MaxBytes: 10000, MaxFileBytes: 5, MaxResults: 10}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, err := local.NewClient().SearchImplementations(context.Background(), implementationCandidate(root), "QueueItems", tc.budget)
			if err != nil {
				t.Fatal(err)
			}
			if !r.Limited || r.Status == "complete" || r.FilesRead > tc.budget.MaxFiles || r.BytesRead > tc.budget.MaxBytes {
				t.Fatalf("budget not explicit/enforced: %+v", r)
			}
			if tc.name == "per_file" && len(r.Hits) > 0 {
				t.Fatalf("oversize file produced definition: %+v", r)
			}
		})
	}
	b := implementationBudget()
	b.MaxResults = 1
	r, err := local.NewClient().SearchImplementations(context.Background(), implementationCandidate(root), "QueueItems", b)
	if err != nil || len(r.Hits) != 1 {
		t.Fatalf("result cap: %+v %v", r, err)
	}
}

func TestLocalImplementationSearchCancelAndMissingRoot(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, err := local.NewClient().SearchImplementations(ctx, implementationCandidate(t.TempDir()), "Anything", implementationBudget())
	if err == nil {
		t.Fatalf("cancellation not returned: %+v", r)
	}
	if r.FilesRead != 0 || r.BytesRead != 0 {
		t.Fatalf("read after cancellation: %+v", r)
	}
	r, err = local.NewClient().SearchImplementations(context.Background(), implementationCandidate(filepath.Join(t.TempDir(), "missing")), "Anything", implementationBudget())
	if err == nil {
		t.Fatalf("missing root not returned as error: %+v", r)
	}
}

func TestLocalImplementationExaminedOrderIsCanonical(t *testing.T) {
	root := t.TempDir()
	source := "package sample\nfunc QueueItems() {}\n"
	for _, path := range []string{"zeta.go", "alpha.go"} {
		implementationFixture(t, root, path, source)
	}
	report, err := local.NewClient().SearchImplementations(context.Background(), implementationCandidate(root), "QueueItems", implementationBudget())
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Examined) != 2 || report.Examined[0].Path != "alpha.go" || report.Examined[1].Path != "zeta.go" {
		t.Fatalf("noncanonical metadata: %+v", report.Examined)
	}
	for _, input := range report.Examined {
		if input.SHA256 != fmt.Sprintf("%x", sha256.Sum256([]byte(source))) {
			t.Fatal("byte identity changed")
		}
	}
}
