/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/independent_step0043_test.md
 * @layer L3
 * @updated 2026-09-07
 */
package local_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/local"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func step0043Candidate(root string) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{Owner: "fixture", Repo: "candidate", URL: "file://" + root, CommitID: strings.Repeat("a", 40)}
}
func step0043Write(t *testing.T, path, value string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(value), 0600); err != nil {
		t.Fatal(err)
	}
}
func step0043Link(t *testing.T, source, target string) {
	t.Helper()
	if err := os.Symlink(source, target); err != nil {
		t.Fatal(err)
	}
}
func TestStep0043ExternalLocalBytesNeverReturned(t *testing.T) {
	for _, name := range []string{"go.mod", "README.md", "linked/secret.go"} {
		t.Run(name, func(t *testing.T) {
			root, outside := t.TempDir(), t.TempDir()
			sentinel := "STEP0043_OUTSIDE_SCOPE_SENTINEL"
			client := local.NewClient()
			candidate := step0043Candidate(root)
			if name == "linked/secret.go" {
				step0043Write(t, filepath.Join(outside, "secret.go"), "package secret\n// "+sentinel)
				step0043Link(t, outside, filepath.Join(root, "linked"))
			} else {
				content := sentinel
				if name == "go.mod" {
					content = "module " + sentinel + "\n\ngo 1.25\n"
				}
				step0043Write(t, filepath.Join(outside, name), content)
				step0043Link(t, filepath.Join(outside, name), filepath.Join(root, name))
			}
			extracted, _ := client.ExtractRepoFiles(context.Background(), candidate, []string{name})
			encoded, _ := json.Marshal(extracted)
			if strings.Contains(string(encoded), sentinel) {
				t.Fatalf("out-of-root bytes exposed through extraction: %s", encoded)
			}
			if name == "go.mod" {
				manifest, _ := client.FetchManifest(context.Background(), candidate)
				encoded, _ = json.Marshal(manifest)
				if strings.Contains(string(encoded), sentinel) {
					t.Fatalf("out-of-root manifest exposed: %s", encoded)
				}
			}
		})
	}
}
func TestStep0043LocalInternalPositiveAndByteIdentity(t *testing.T) {
	root := t.TempDir()
	client := local.NewClient()
	candidate := step0043Candidate(root)
	step0043Write(t, filepath.Join(root, "go.mod"), "module example.com/inside\n\ngo 1.25\n")
	step0043Write(t, filepath.Join(root, "inside.go"), "package inside\nconst Value = 1\n")
	step0043Link(t, filepath.Join(root, "inside.go"), filepath.Join(root, "alias.go"))
	manifest, err := client.FetchManifest(context.Background(), candidate)
	if err != nil || manifest == nil || manifest.ModuleName != "example.com/inside" {
		t.Fatalf("valid internal manifest rejected: %+v %v", manifest, err)
	}
	read := func() *domain.ExtractedRepository {
		t.Helper()
		got, err := client.ExtractRepoFiles(context.Background(), candidate, []string{"inside.go", "alias.go"})
		if err != nil || got == nil {
			t.Fatalf("valid extraction: %v", err)
		}
		if len(got.Files) != 2 {
			t.Fatalf("internal real/link positive: got %d files", len(got.Files))
		}
		return got
	}
	first, again := read(), read()
	identity := func(r *domain.ExtractedRepository) string {
		return r.Ref.Key() + "|" + r.SnapshotHash + "|" + r.Files[0].Evidence.SHA256
	}
	if identity(first) != identity(again) {
		t.Fatal("unchanged bytes changed identity")
	}
	step0043Write(t, filepath.Join(root, "inside.go"), "package inside\nconst Value = 2\n")
	next := read()
	if identity(first) == identity(next) {
		t.Fatal("changed worktree bytes kept identical extraction identity")
	}
	for _, file := range next.Files {
		if !strings.Contains(file.Content, "Value = 2") {
			t.Fatalf("stale content: %s", file.Path)
		}
	}
}
