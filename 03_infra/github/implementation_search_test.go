/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/implementation_search_test.md
 * @layer L3
 * @updated 2026-09-08
 */
package github_test

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/03_infra/content"
	"bateia/03_infra/github"
	"bateia/03_infra/local"
)

const step50Commit = "1234567890abcdef1234567890abcdef12345678"
const step50Source = "package mechanism\n\n// ExpireSession removes a session after its deadline.\nfunc ExpireSession(deadline, now int) bool {\n return now > deadline\n}\n"
const step50Other = "package mechanism\n\nfunc RenderPalette() string { return \"blue\" }\n"

func step50Candidate() domain.RepositoryCandidate {
	return domain.RepositoryCandidate{Owner: "fixture", Repo: "mechanism", CommitID: step50Commit, URL: "https://github.com/fixture/mechanism"}
}
func step50SHA(body string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("blob %d\x00%s", len(body), body))))
}
func step50Hash(body string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(body))) }
func step50Node(path, body string) map[string]any {
	return map[string]any{"path": path, "type": "blob", "mode": "100644", "sha": step50SHA(body), "size": len(body)}
}
func step50Blob(body string) map[string]any {
	return map[string]any{"sha": step50SHA(body), "size": len(body), "encoding": "base64", "content": base64.StdEncoding.EncodeToString([]byte(body))}
}

type step50Fixture struct {
	mu     sync.Mutex
	tree   any
	blobs  map[string]any
	status int
	paths  []string
}

func step50Server(t *testing.T, f *step50Fixture) *github.Client {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.paths = append(f.paths, r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/repos/fixture/mechanism/git/trees/"+step50Commit {
			if r.URL.Query().Get("recursive") != "1" {
				t.Errorf("nonrecursive tree request: %s", r.URL)
			}
			if raw, ok := f.tree.(string); ok {
				fmt.Fprint(w, raw)
			} else {
				json.NewEncoder(w).Encode(f.tree)
			}
			return
		}
		prefix := "/repos/fixture/mechanism/git/blobs/"
		if strings.HasPrefix(r.URL.Path, prefix) {
			if f.status != 0 {
				w.WriteHeader(f.status)
				fmt.Fprint(w, `{"message":"blocked"}`)
				return
			}
			if blob, ok := f.blobs[strings.TrimPrefix(r.URL.Path, prefix)]; ok {
				if raw, yes := blob.(string); yes {
					fmt.Fprint(w, raw)
				} else {
					json.NewEncoder(w).Encode(blob)
				}
				return
			}
		}
		t.Errorf("unexpected request, snapshot must remain pinned: %s", r.URL)
		http.NotFound(w, r)
	}))
	t.Cleanup(s.Close)
	return github.NewClientWithBaseURL("", s.URL)
}
func step50Setup(nodes ...map[string]any) *step50Fixture {
	return &step50Fixture{tree: map[string]any{"sha": step50Commit, "truncated": false, "tree": nodes}, blobs: map[string]any{step50SHA(step50Source): step50Blob(step50Source), step50SHA(step50Other): step50Blob(step50Other)}}
}
func (f *step50Fixture) requests() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.paths...)
}
func step50Search(t *testing.T, c *github.Client, b domain.ImplementationSearchBudget) domain.ImplementationSearchReport {
	t.Helper()
	r, e := c.SearchImplementations(context.Background(), step50Candidate(), "expire session", b)
	if e != nil {
		t.Fatalf("search: %v", e)
	}
	return r
}
func step50Partial(t *testing.T, r domain.ImplementationSearchReport) {
	t.Helper()
	if r.Status != "partial" {
		t.Errorf("opaque/incomplete retrieval must be partial: %+v", r)
	}
}

func TestStep0050RemoteLocalParityAndUsecase(t *testing.T) {
	f := step50Setup(step50Node("internal/zeta.go", step50Other), step50Node("internal/alpha.go", step50Source))
	remote := step50Server(t, f)
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "internal"), 0700); err != nil {
		t.Fatal(err)
	}
	for p, b := range map[string]string{"internal/alpha.go": step50Source, "internal/zeta.go": step50Other} {
		if err := os.WriteFile(filepath.Join(root, p), []byte(b), 0600); err != nil {
			t.Fatal(err)
		}
	}
	localCandidate := step50Candidate()
	localCandidate.URL = "file://" + root
	router := content.NewRouter(remote, local.NewClient())
	budget := domain.ImplementationSearchBudget{MaxFiles: 10, MaxBytes: 4096, MaxFileBytes: 2048, MaxResults: 24}
	r, e := router.SearchImplementations(context.Background(), step50Candidate(), "expire session", budget)
	if e != nil {
		t.Fatal(e)
	}
	l, e := router.SearchImplementations(context.Background(), localCandidate, "expire session", budget)
	if e != nil {
		t.Fatal(e)
	}
	if r.Status != "complete" || len(r.Hits) == 0 {
		t.Fatalf("functional query needs remote implementation at neutral path: %+v", r)
	}
	if !reflect.DeepEqual(r.Hits, l.Hits) || !reflect.DeepEqual(r.Examined, l.Examined) {
		t.Fatalf("local/remote hits and hashes differ\nremote=%+v\nlocal=%+v", r, l)
	}
	if len(r.Examined) != 2 || r.BytesRead != int64(len(step50Source)+len(step50Other)) {
		t.Fatalf("incorrect accounting: %+v", r)
	}
	if r.Examined[0].Path != "internal/alpha.go" || r.Examined[0].SHA256 != step50Hash(step50Source) || r.Examined[0].Bytes != int64(len(step50Source)) {
		t.Fatalf("missing byte identity: %+v", r.Examined)
	}
	if len(f.requests()) != 3 {
		t.Fatalf("expected one tree and two blobs: %v", f.requests())
	}
	journey := usecase.RetrieveImplementationEvidence(context.Background(), router, step50Candidate(), "expire session")
	paths := usecase.SelectImplementationPaths("expire session", []domain.DirectoryNode{{Path: "internal/alpha.go", Type: "file"}, {Path: "internal/zeta.go", Type: "file"}}, journey, 1)
	if len(paths) != 1 || paths[0] != "internal/alpha.go" {
		t.Fatalf("hit must win selection: %v", paths)
	}
	intact := usecase.ValidateImplementationHits(journey, []domain.ExtractedFile{{Path: "internal/alpha.go", Content: step50Source}, {Path: "internal/zeta.go", Content: step50Other}})
	if len(intact.Hits) != len(journey.Hits) {
		t.Fatal("valid reread lost hits")
	}
	for _, files := range [][]domain.ExtractedFile{nil, {{Path: "internal/alpha.go", Content: step50Source + "// changed\n"}}} {
		changed := usecase.ValidateImplementationHits(journey, files)
		step50Partial(t, changed)
		if len(changed.Hits) != 0 {
			t.Fatal("unverified reread retained hit")
		}
	}
}

func TestStep0050RejectIdentityBeforeIO(t *testing.T) {
	for _, bad := range []domain.RepositoryCandidate{{Owner: "fixture", Repo: "mechanism", CommitID: "HEAD"}, {Owner: "fixture", Repo: "mechanism", CommitID: "1234567"}, {Owner: "fixture", Repo: "mechanism", CommitID: strings.Repeat("g", 40)}, {Owner: "../fixture", Repo: "mechanism", CommitID: step50Commit}, {Owner: "fixture", Repo: "../mechanism", CommitID: step50Commit}, {Owner: "", Repo: "mechanism", CommitID: step50Commit}} {
		f := step50Setup()
		c := step50Server(t, f)
		r, e := c.SearchImplementations(context.Background(), bad, "expire session", domain.ImplementationSearchBudget{})
		if e == nil && r.Status != "partial" {
			t.Errorf("invalid identity accepted: %+v", bad)
		}
		if len(f.requests()) != 0 {
			t.Fatalf("invalid identity performed I/O: %v", f.requests())
		}
	}
}

func TestStep0050RejectOpaqueBlobs(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"sha", func(b map[string]any) { b["sha"] = strings.Repeat("f", 40) }},
		{"content", func(b map[string]any) {
			b["content"] = base64.StdEncoding.EncodeToString([]byte(strings.Repeat("x", len(step50Source))))
		}},
		{"size", func(b map[string]any) { b["size"] = len(step50Source) + 1 }},
		{"encoding", func(b map[string]any) { b["encoding"] = "utf-8" }},
		{"base64", func(b map[string]any) { b["content"] = "!!!" }},
		{"missing_content", func(b map[string]any) { delete(b, "content") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := step50Setup(step50Node("internal/alpha.go", step50Source))
			b := step50Blob(step50Source)
			tc.mutate(b)
			f.blobs[step50SHA(step50Source)] = b
			c := step50Server(t, f)
			r, _ := c.SearchImplementations(context.Background(), step50Candidate(), "expire session", domain.ImplementationSearchBudget{})
			step50Partial(t, r)
			if len(r.Examined) != 0 || len(r.Hits) != 0 {
				t.Fatal("opaque blob admitted")
			}
			if tc.name == "sha" || tc.name == "content" || tc.name == "size" {
				if r.BytesRead != int64(len(step50Source)) {
					t.Errorf("discarded decoded bytes unaccounted: %d", r.BytesRead)
				}
			}
		})
	}
}

func TestStep0050OpaqueTreesAndPaths(t *testing.T) {
	for _, tc := range []struct {
		name     string
		tree     any
		maxBlobs int
	}{
		{"truncated", map[string]any{"sha": step50Commit, "truncated": true, "tree": []any{step50Node("internal/alpha.go", step50Source)}}, 1},
		{"invalid_json", `{"tree":`, 0},
		{"missing_tree", map[string]any{"sha": step50Commit}, 0},
		{"duplicate_ambiguous", map[string]any{"sha": step50Commit, "truncated": false, "tree": []any{step50Node("internal/alpha.go", step50Source), step50Node("internal/alpha.go", step50Other)}}, 0},
		{"unsafe_path", map[string]any{"sha": step50Commit, "truncated": false, "tree": []any{step50Node("../alpha.go", step50Source)}}, 0},
		{"oversized_envelope", `{"padding":"` + strings.Repeat("x", 8*1024*1024) + `","tree":[]}`, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := step50Setup()
			f.tree = tc.tree
			c := step50Server(t, f)
			r, _ := c.SearchImplementations(context.Background(), step50Candidate(), "expire session", domain.ImplementationSearchBudget{})
			step50Partial(t, r)
			if len(f.requests()) > 1+tc.maxBlobs {
				t.Fatalf("unexpected blob reads: %v", f.requests())
			}
			if tc.maxBlobs == 0 && len(r.Examined) != 0 {
				t.Fatal("invalid tree admitted sources")
			}
		})
	}
	t.Run("symlink", func(t *testing.T) {
		node := step50Node("internal/alpha.go", step50Source)
		node["mode"] = "120000"
		f := step50Setup(node)
		c := step50Server(t, f)
		r, _ := c.SearchImplementations(context.Background(), step50Candidate(), "expire session", domain.ImplementationSearchBudget{})
		if len(r.Examined) != 0 || len(f.requests()) != 1 {
			t.Fatal("symlink read as regular blob")
		}
	})
}

func TestStep0050BudgetsAndStops(t *testing.T) {
	for _, tc := range []struct {
		name        string
		b           domain.ImplementationSearchBudget
		maxRequests int
		maxExamined int
	}{
		{"attempts", domain.ImplementationSearchBudget{MaxFiles: 1}, 2, 1},
		{"file_bytes", domain.ImplementationSearchBudget{MaxFileBytes: 8}, 1, 0},
		{"total_bytes", domain.ImplementationSearchBudget{MaxBytes: 8}, 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := step50Setup(step50Node("a.go", step50Source), step50Node("b.go", step50Other))
			c := step50Server(t, f)
			r, _ := c.SearchImplementations(context.Background(), step50Candidate(), "expire session", tc.b)
			step50Partial(t, r)
			if !r.Limited || len(f.requests()) > tc.maxRequests || len(r.Examined) > tc.maxExamined {
				t.Fatalf("budget not enforced: %+v requests=%v", r, f.requests())
			}
		})
	}
	t.Run("remote_ceiling", func(t *testing.T) {
		nodes := []map[string]any{}
		for i := 0; i < 66; i++ {
			nodes = append(nodes, step50Node(fmt.Sprintf("p%03d.go", i), step50Source))
		}
		f := step50Setup(nodes...)
		c := step50Server(t, f)
		r, _ := c.SearchImplementations(context.Background(), step50Candidate(), "expire session", domain.ImplementationSearchBudget{MaxFiles: 1000, MaxBytes: 1 << 30, MaxFileBytes: 1 << 30, MaxResults: 1000})
		if len(f.requests()) > 65 || len(r.Examined) > 64 || len(r.Hits) > 24 || !r.Limited {
			t.Fatalf("ceilings breached: %+v requests=%d", r, len(f.requests()))
		}
		if r.Budget.MaxFiles > 64 || r.Budget.MaxBytes > 4*1024*1024 || r.Budget.MaxFileBytes > 512*1024 || r.Budget.MaxResults > 24 {
			t.Fatalf("effective budget exceeds ceiling: %+v", r.Budget)
		}
	})
	for _, status := range []int{403, 429} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			f := step50Setup(step50Node("a.go", step50Source), step50Node("b.go", step50Other))
			f.status = status
			c := step50Server(t, f)
			r, _ := c.SearchImplementations(context.Background(), step50Candidate(), "expire session", domain.ImplementationSearchBudget{})
			step50Partial(t, r)
			if len(f.requests()) != 2 {
				t.Fatalf("rate limit did not stop reads: %v", f.requests())
			}
		})
	}
	t.Run("cancelled", func(t *testing.T) {
		f := step50Setup(step50Node("a.go", step50Source))
		c := step50Server(t, f)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		r, e := c.SearchImplementations(ctx, step50Candidate(), "expire session", domain.ImplementationSearchBudget{})
		if e == nil && r.Status != "partial" {
			t.Fatal("cancelled search considered complete")
		}
		if len(f.requests()) != 0 {
			t.Fatal("cancelled context performed I/O")
		}
	})
	t.Run("blob_envelope", func(t *testing.T) {
		f := step50Setup(step50Node("a.go", step50Source))
		f.blobs[step50SHA(step50Source)] = `{"padding":"` + strings.Repeat("x", 2*1024*1024) + `","sha":"` + step50SHA(step50Source) + `","encoding":"base64","size":` + fmt.Sprint(len(step50Source)) + `,"content":"` + base64.StdEncoding.EncodeToString([]byte(step50Source)) + `"}`
		c := step50Server(t, f)
		r, _ := c.SearchImplementations(context.Background(), step50Candidate(), "expire session", domain.ImplementationSearchBudget{MaxFileBytes: 1024})
		step50Partial(t, r)
		if len(r.Examined) != 0 {
			t.Fatal("unbounded blob envelope admitted")
		}
	})
}

// This wrapper intentionally exposes only the original ContentProvider port.
type step50LegacyProvider struct{ domain.ContentProvider }

func TestStep0050RouterUnsupportedProvider(t *testing.T) {
	f := step50Setup()
	remote := step50LegacyProvider{step50Server(t, f)}
	router := content.NewRouter(remote, local.NewClient())
	r := usecase.RetrieveImplementationEvidence(context.Background(), router, step50Candidate(), "expire session")
	if r.Status != "not_evaluated" || len(f.requests()) != 0 {
		t.Fatalf("unsupported remote must stay not_evaluated: %+v", r)
	}
}
