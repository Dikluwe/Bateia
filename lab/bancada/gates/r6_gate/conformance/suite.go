/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r6_gate/discriminatory_step0026_r6_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package conformance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

func RunSuite(t *testing.T, factory Factory) {
	t.Helper()
	ctx := context.Background()
	req := Request{Problem: "choose", Capabilities: []string{"rank", "search"}, Users: []string{"u2", "u1"}}
	newSurface := func(t *testing.T) (Surface, string) { t.Helper(); root := t.TempDir(); return factory(root), root }
	t.Run("positive-canonical-concurrent-start", func(t *testing.T) {
		s, _ := newSurface(t)
		a := req
		b := req
		b.Capabilities = []string{"search", "rank", "search"}
		b.Users = []string{"u1", "u2", "u1"}
		var wg sync.WaitGroup
		receipts := make(chan Receipt, 16)
		errs := make(chan error, 16)
		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				candidate := a
				if n%2 == 1 {
					candidate = b
				}
				r, e := s.Start(ctx, "positive-same", candidate)
				receipts <- r
				errs <- e
			}(i)
		}
		wg.Wait()
		close(receipts)
		close(errs)
		ids := map[string]bool{}
		for e := range errs {
			if e != nil {
				t.Fatal(e)
			}
		}
		for r := range receipts {
			ids[r.RunID] = true
		}
		if len(ids) != 1 || s.Effects().IDs != 1 {
			t.Fatalf("duplicate work ids=%v effects=%#v", ids, s.Effects())
		}
	})
	t.Run("positive-terminal-reopens", func(t *testing.T) {
		s, root := newSurface(t)
		r, e := s.Start(ctx, "positive-terminal", req)
		if e == nil {
			_, e = s.Complete(ctx, r.RunID)
		}
		reopened := factory(root)
		run, e2 := reopened.GetRun(ctx, r.RunID)
		if e != nil || e2 != nil || run.Lifecycle != "SUCCEEDED" || run.Status != "INCONCLUSIVE" {
			t.Fatalf("reopen %#v %v/%v", run, e, e2)
		}
	})
	t.Run("positive-cursor-reopens", func(t *testing.T) {
		s, root := newSurface(t)
		r, _ := s.Start(ctx, "positive-cursor", req)
		run, e := s.Publish(ctx, r.RunID, r.Version, "provenance", document(1, 2, 3))
		if e != nil {
			t.Fatal(e)
		}
		p1, e := s.Read(ctx, run.ID, "provenance", "", 1)
		if e != nil {
			t.Fatal(e)
		}
		reopened := factory(root)
		p2, e := reopened.Read(ctx, run.ID, "provenance", p1.Next, 2)
		if e != nil || len(p2.Entries) != 2 {
			t.Fatalf("cursor reopen %#v %v", p2, e)
		}
	})
	t.Run("positive-cancellation-is-idempotent", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "positive-cancel", req)
		first, e := s.Cancel(ctx, r.RunID)
		second, e2 := s.Cancel(ctx, r.RunID)
		if e != nil || e2 != nil || first.Lifecycle != "CANCELLED" || first.Version != second.Version {
			t.Fatalf("cancel %#v %#v %v/%v", first, second, e, e2)
		}
	})
	t.Run("positive-legacy-surface", func(t *testing.T) {
		s, _ := newSurface(t)
		b, e := s.Call(ctx, "legacy")
		if e != nil || string(b) != `{"tools":["search","benchmark","recommend"]}` {
			t.Fatalf("legacy %s %v", b, e)
		}
	})
	t.Run("identity-is-opaque", func(t *testing.T) {
		s, root := newSurface(t)
		seen := map[string]bool{}
		for i := 0; i < 6; i++ {
			key := "CANARY-" + string(rune('a'+i))
			r, e := s.Start(ctx, key, req)
			if e != nil {
				t.Fatal(e)
			}
			keyDigest := sha256.Sum256([]byte(key))
			correlated := "run_" + hex.EncodeToString(keyDigest[:16])
			if !regexp.MustCompile(`^run_[0-9a-f]{32}$`).MatchString(r.RunID) || r.RunID == correlated || strings.Contains(strings.Join(rootNames(root), " "), "CANARY") || seen[r.RunID] {
				t.Fatalf("correlated identity %q", r.RunID)
			}
			seen[r.RunID] = true
		}
	})
	t.Run("retry-returns-current-receipt", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "retry", req)
		run, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "provider"})
		if e != nil {
			t.Fatal(e)
		}
		again, e := s.Start(ctx, "retry", req)
		if e != nil || again.Lifecycle != run.Lifecycle || again.Version != run.Version {
			t.Fatalf("stale retry %#v %v", again, e)
		}
	})
	t.Run("invalid-input-has-zero-effects", func(t *testing.T) {
		s, _ := newSurface(t)
		_, e := s.Start(ctx, "../bad", Request{})
		if !errors.Is(e, ErrInvalid) || s.Effects() != (Effects{}) {
			t.Fatalf("invalid touched ports: %v %#v", e, s.Effects())
		}
	})
	t.Run("quota-full-is-nondestructive", func(t *testing.T) {
		s, root := newSurface(t)
		if e := os.WriteFile(filepath.Join(root, "quota-full"), []byte("1"), 0o600); e != nil {
			t.Fatal(e)
		}
		_, e := s.Start(ctx, "quota", req)
		if !errors.Is(e, ErrQuota) || s.Effects() != (Effects{}) || hasPrefix(root, "run-") {
			t.Fatalf("quota residue %v %#v", e, s.Effects())
		}
	})
	t.Run("artifact-is-immutable", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "immutable", req)
		run, e := s.Publish(ctx, r.RunID, r.Version, "provenance", document(1))
		if e != nil {
			t.Fatal(e)
		}
		_, e = s.Publish(ctx, r.RunID, run.Version, "provenance", document(2))
		if !errors.Is(e, ErrImmutable) {
			t.Fatalf("overwrite accepted: %v", e)
		}
	})
	t.Run("cas-has-one-winner", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "cas", req)
		if _, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "provider"}); e != nil {
			t.Fatal(e)
		}
		if _, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 2, "provider"}); !errors.Is(e, ErrCAS) {
			t.Fatalf("stale CAS accepted: %v", e)
		}
		s, _ = newSurface(t)
		r, _ = s.Start(ctx, "cas-concurrent", req)
		var wg sync.WaitGroup
		errs := make(chan error, 2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				_, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, n + 1, "provider"})
				errs <- e
			}(i)
		}
		wg.Wait()
		close(errs)
		wins, conflicts := 0, 0
		for e := range errs {
			if e == nil {
				wins++
			} else if errors.Is(e, ErrCAS) {
				conflicts++
			}
		}
		if wins != 1 || conflicts != 1 {
			t.Fatalf("wins=%d conflicts=%d", wins, conflicts)
		}
	})
	t.Run("progress-matrix-is-strict", func(t *testing.T) {
		for _, p := range []Progress{{7, 1, "bad"}, {8, 9, "bad"}, {8, -1, "bad"}} {
			s, _ := newSurface(t)
			r, _ := s.Start(ctx, "progress", req)
			if _, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", p); !errors.Is(e, ErrInvalid) {
				t.Fatalf("accepted %#v: %v", p, e)
			}
		}
	})
	t.Run("terminal-race-has-one-event", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "race", req)
		_, _ = s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "provider"})
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); _, _ = s.Cancel(ctx, r.RunID) }()
		go func() { defer wg.Done(); _, _ = s.Complete(ctx, r.RunID) }()
		wg.Wait()
		run, e := s.GetRun(ctx, r.RunID)
		if e != nil || len(run.TerminalEvents) != 1 {
			t.Fatalf("terminal events %#v %v", run.TerminalEvents, e)
		}
	})
	t.Run("recovery-interrupts-nonterminal", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "recover", req)
		_, _ = s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 3, "provider"})
		if e := s.Recover(ctx); e != nil {
			t.Fatal(e)
		}
		run, e := s.GetRun(ctx, r.RunID)
		if e != nil || run.Lifecycle != "INTERRUPTED" {
			t.Fatalf("recovery %#v %v", run, e)
		}
	})
	t.Run("start-waits-for-recovery", func(t *testing.T) {
		s, root := newSurface(t)
		block := filepath.Join(root, "recovery-block")
		_ = os.WriteFile(block, []byte("1"), 0o600)
		recovered := make(chan error, 1)
		go func() { recovered <- s.Recover(ctx) }()
		waitFile(t, filepath.Join(root, "recovery-entered"))
		deadline, cancel := context.WithTimeout(ctx, 25*time.Millisecond)
		_, startErr := s.Start(deadline, "during-recovery", req)
		cancel()
		if startErr == nil {
			t.Fatal("start escaped recovery")
		}
		_ = os.Remove(block)
		if e := <-recovered; e != nil {
			t.Fatal(e)
		}
	})
	t.Run("shutdown-does-not-guess-cancelled", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "shutdown", req)
		_, _ = s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "provider"})
		c, cancel := context.WithCancel(ctx)
		cancel()
		_ = s.Shutdown(c)
		_ = s.Recover(ctx)
		run, _ := s.GetRun(ctx, r.RunID)
		if run.Lifecycle != "INTERRUPTED" {
			t.Fatalf("shutdown guessed %s", run.Lifecycle)
		}
	})
	t.Run("unpublished-is-not-ready", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "ready", req)
		_, e := s.Read(ctx, r.RunID, "provenance", "", 1)
		if !errors.Is(e, ErrNotReady) {
			t.Fatalf("unpublished read %v", e)
		}
	})
	t.Run("pending-is-not-readable", func(t *testing.T) {
		s, root := newSurface(t)
		r, _ := s.Start(ctx, "pending", req)
		writeArtifactFixture(t, root, r.RunID, "provenance", "PENDING", document(1))
		_, e := s.Read(ctx, r.RunID, "provenance", "", 1)
		if !errors.Is(e, ErrNotReady) {
			t.Fatalf("pending read %v", e)
		}
	})
	t.Run("orphan-blob-is-not-readable", func(t *testing.T) {
		s, root := newSurface(t)
		r, _ := s.Start(ctx, "orphan", req)
		_ = os.WriteFile(blobPath(root, r.RunID, "failures"), document(1), 0o600)
		_, e := s.Read(ctx, r.RunID, "failures", "", 1)
		if !errors.Is(e, ErrNotReady) {
			t.Fatalf("orphan read %v", e)
		}
	})
	t.Run("terminal-requires-complete-artifacts", func(t *testing.T) {
		s, root := newSurface(t)
		r, _ := s.Start(ctx, "terminal", req)
		_, _ = s.Complete(ctx, r.RunID)
		_ = os.Remove(blobPath(root, r.RunID, "decision"))
		_, e := s.GetRun(ctx, r.RunID)
		if !errors.Is(e, ErrIntegrity) {
			t.Fatalf("missing terminal artifact %v", e)
		}
	})
	t.Run("corrupt-terminal-is-integrity", func(t *testing.T) {
		s, root := newSurface(t)
		r, _ := s.Start(ctx, "corrupt", req)
		_, _ = s.Complete(ctx, r.RunID)
		_ = os.WriteFile(blobPath(root, r.RunID, "decision"), []byte("corrupt"), 0o600)
		_, e := s.GetRun(ctx, r.RunID)
		if !errors.Is(e, ErrIntegrity) {
			t.Fatalf("corrupt terminal %v", e)
		}
	})
	t.Run("cursor-hmac-rejects-tamper", func(t *testing.T) {
		s, _ := newSurface(t)
		r, _ := s.Start(ctx, "cursor", req)
		run, _ := s.Publish(ctx, r.RunID, r.Version, "provenance", document(1, 2))
		p, e := s.Read(ctx, r.RunID, "provenance", "", 1)
		if e != nil || p.Next == "" {
			t.Fatalf("cursor setup %v", e)
		}
		parts := strings.Split(p.Next, ".")
		replacement := "A"
		if strings.HasPrefix(parts[1], "A") {
			replacement = "B"
		}
		tampered := parts[0] + "." + replacement + parts[1][1:]
		_, e = s.Read(ctx, run.ID, "provenance", tampered, 1)
		if !errors.Is(e, ErrCursor) {
			t.Fatalf("tamper accepted %v", e)
		}
	})
	t.Run("filesystem-rejects-symlink", func(t *testing.T) {
		s, root := newSurface(t)
		outside := filepath.Join(t.TempDir(), "outside")
		_ = os.WriteFile(outside, []byte("x"), 0o600)
		_ = os.Symlink(outside, filepath.Join(root, "link"))
		if !errors.Is(s.ValidateRoot(), ErrIntegrity) {
			t.Fatal("symlink accepted")
		}
	})
	t.Run("filesystem-rejects-hardlink", func(t *testing.T) {
		s, root := newSurface(t)
		a := filepath.Join(root, "a")
		_ = os.WriteFile(a, []byte("x"), 0o600)
		_ = os.Link(a, filepath.Join(root, "b"))
		if !errors.Is(s.ValidateRoot(), ErrIntegrity) {
			t.Fatal("hardlink accepted")
		}
	})
	t.Run("create-crash-prefixes-are-atomic", func(t *testing.T) {
		s, _ := newSurface(t)
		observations, err := s.CrashCreate()
		assertCrashMatrix(t, observations, err)
	})
	t.Run("transition-crash-prefixes-are-atomic", func(t *testing.T) {
		s, _ := newSurface(t)
		observations, err := s.CrashTransition()
		assertCrashMatrix(t, observations, err)
	})
	t.Run("publish-crash-prefixes-are-atomic", func(t *testing.T) {
		s, _ := newSurface(t)
		observations, err := s.CrashPublish()
		assertCrashMatrix(t, observations, err)
	})
	t.Run("disk-and-log-privacy", func(t *testing.T) {
		s, root := newSurface(t)
		_, _ = s.Start(ctx, "RAW-KEY-CANARY", req)
		b, e := s.Call(ctx, "logs")
		if e != nil {
			t.Fatal(e)
		}
		for _, name := range rootNames(root) {
			data, _ := os.ReadFile(filepath.Join(root, name))
			b = append(b, data...)
		}
		if strings.Contains(string(b), "RAW-KEY-CANARY") || strings.Contains(string(b), "Bearer TOKEN") {
			t.Fatalf("canary leaked: %s", b)
		}
	})
	t.Run("projection-is-minimal", func(t *testing.T) {
		s, _ := newSurface(t)
		b, e := s.Call(ctx, "projection")
		if e != nil || string(b) != `{"entries":[],"failure":"recommendation execution failed"}` {
			t.Fatalf("projection %s %v", b, e)
		}
	})
	t.Run("errors-are-sanitized", func(t *testing.T) {
		s, _ := newSurface(t)
		_, e := s.Call(ctx, "error")
		if e == nil || e.Error() != "RUN_STORE_UNAVAILABLE" {
			t.Fatalf("error leaked %v", e)
		}
	})
	t.Run("json-closure-is-exact", func(t *testing.T) {
		s, _ := newSurface(t)
		b, e := s.Call(ctx, "empty")
		if e != nil || string(b) != `{"artifacts":[],"failure":null}` {
			t.Fatalf("closure %s %v", b, e)
		}
	})
	t.Run("l1-imports-remain-pure", func(t *testing.T) {
		imports := packageImports(t)
		for _, bad := range []string{"net/http", "os/exec", "database/sql"} {
			if imports[bad] {
				t.Fatalf("forbidden import %s", bad)
			}
		}
	})
	t.Run("l2-does-not-import-concrete-store", func(t *testing.T) {
		imports := packageImports(t)
		if imports["bateia/lab/bancada/gates/r6_gate/storefixture"] {
			t.Fatal("L2 imported concrete store")
		}
		s, _ := newSurface(t)
		b, e := s.Call(ctx, "legacy")
		if e != nil || string(b) != `{"tools":["search","benchmark","recommend"]}` {
			t.Fatalf("legacy %s %v", b, e)
		}
	})
	t.Run("filesystem-rejects-public-and-nonregular", func(t *testing.T) {
		s, root := newSurface(t)
		_ = os.WriteFile(filepath.Join(root, "public"), []byte("x"), 0o644)
		_ = syscall.Mkfifo(filepath.Join(root, "fifo"), 0o600)
		if !errors.Is(s.ValidateRoot(), ErrIntegrity) {
			t.Fatal("unsafe kind/mode accepted")
		}
	})
	t.Run("not-found-and-integrity-are-distinct", func(t *testing.T) {
		s, root := newSurface(t)
		_, e := s.GetRun(ctx, "run_00000000000000000000000000000000")
		if !errors.Is(e, ErrNotFound) {
			t.Fatalf("missing %v", e)
		}
		r, _ := s.Start(ctx, "corrupt-record", req)
		_ = os.WriteFile(filepath.Join(root, "run-"+r.RunID), []byte("{"), 0o600)
		_, e = s.GetRun(ctx, r.RunID)
		if !errors.Is(e, ErrIntegrity) {
			t.Fatalf("corrupt %v", e)
		}
	})
}

func document(values ...int) []byte {
	b, _ := json.Marshal(map[string]interface{}{"entries": values})
	return b
}
func blobPath(root, id, name string) string { return filepath.Join(root, "blob-"+id+"-"+name) }
func rootNames(root string) []string {
	entries, _ := os.ReadDir(root)
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out
}
func hasPrefix(root, prefix string) bool {
	for _, name := range rootNames(root) {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
func waitFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if _, e := os.Stat(path); e == nil {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("file did not appear: %s", path)
}
func writeArtifactFixture(t *testing.T, root, id, name, state string, data []byte) {
	t.Helper()
	sum := sha256.Sum256(data)
	meta := Artifact{Name: name, State: state, Hash: "sha256:" + hex.EncodeToString(sum[:]), Version: 2}
	b, _ := json.Marshal(meta)
	_ = os.WriteFile(filepath.Join(root, "artifact-"+id+"-"+name), b, 0o600)
	_ = os.WriteFile(blobPath(root, id, name), data, 0o600)
}
func assertCrashMatrix(t *testing.T, observations []CrashObservation, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"open", "write", "file-sync", "rename", "dir-sync"}
	if len(observations) != len(want) {
		t.Fatalf("boundaries=%d", len(observations))
	}
	for i, o := range observations {
		if o.Boundary != want[i] || o.Temp || (!o.Confirmed && o.Final) || o.Reopen == "ambiguous" {
			t.Fatalf("unsafe prefix %#v", o)
		}
	}
}
func packageImports(t *testing.T) map[string]bool {
	t.Helper()
	cwd, e := os.Getwd()
	if e != nil {
		t.Fatal(e)
	}
	set := map[string]bool{}
	fset := token.NewFileSet()
	entries, _ := os.ReadDir(cwd)
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".go" {
			continue
		}
		f, e := parser.ParseFile(fset, filepath.Join(cwd, entry.Name()), nil, parser.ImportsOnly)
		if e != nil {
			t.Fatal(e)
		}
		for _, imp := range f.Imports {
			set[strings.Trim(imp.Path.Value, "\"")] = true
		}
	}
	return set
}

var _ = reflect.DeepEqual
