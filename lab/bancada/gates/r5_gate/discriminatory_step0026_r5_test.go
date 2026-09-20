/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r5_gate/discriminatory_step0026_r5_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package r5gate

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

type adversarialFile struct {
	Cases []attackCase `json:"cases"`
}

type attackCase struct {
	ID, Kind, Mutation, Stimulus, Expected, Witness string
}

type observation struct {
	Calls  []string `json:"calls"`
	Files  []string `json:"files"`
	Source []string `json:"source,omitempty"`
}

type runResult struct {
	Verdict, Witness string
}

func TestStep0026R5RealSubstituteMutationGateThreeByThree(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "lab", "bancada", "history", "benchmarks", "0026", "adversarial-cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var suite adversarialFile
	if err := json.Unmarshal(data, &suite); err != nil {
		t.Fatal(err)
	}
	if len(suite.Cases) != 70 {
		t.Fatalf("want 70 cases, got %d", len(suite.Cases))
	}
	byID := make(map[string]attackCase, len(suite.Cases))
	for _, c := range suite.Cases {
		if c.ID == "" || c.Witness == "" || c.Expected == "" || (c.Kind != "opaque_case" && c.Stimulus == "") {
			t.Fatalf("incomplete case %#v", c)
		}
		if c.Kind == "negative_mutation" && c.Mutation == "" {
			t.Fatalf("mutation missing for %s", c.ID)
		}
		byID[c.ID] = c
	}
	orders := [][]string{orderedIDs(false), orderedIDs(true), permutedIDs()}
	baseline := map[string]runResult{}
	for oi, order := range orders {
		for repetition := 0; repetition < 3; repetition++ {
			for _, id := range order {
				c := byID[id]
				got := executeCase(t, c)
				if got.Verdict != c.Expected {
					t.Fatalf("order=%d repetition=%d case=%s expected=%s got=%s witness=%s", oi, repetition, id, c.Expected, got.Verdict, got.Witness)
				}
				if prior, ok := baseline[id]; ok && prior != got {
					t.Fatalf("unstable case %s: %#v != %#v", id, prior, got)
				}
				baseline[id] = got
			}
		}
	}
	for i := 1; i <= 60; i++ {
		if baseline[fmt.Sprintf("M%02d", i)].Verdict != "Violated" {
			t.Fatalf("surviving mutation M%02d", i)
		}
	}
}

func executeCase(t *testing.T, c attackCase) runResult {
	t.Helper()
	caseBinding := sha256.Sum256([]byte(c.Mutation + "\x00" + c.Stimulus + "\x00" + c.Witness))
	if strings.HasPrefix(c.ID, "O") {
		return runResult{"Unknown", "environment-owned:" + hex.EncodeToString(caseBinding[:8])}
	}
	if strings.HasPrefix(c.ID, "P") {
		obs, ok := positiveControl(t, c.ID)
		if !ok {
			return runResult{"Violated", stableWitness(caseBinding[:], obs, obs)}
		}
		return runResult{"Preserved", stableWitness(caseBinding[:], obs, obs)}
	}
	ref := observeMutation(t, c.ID, "")
	mut := observeMutation(t, c.ID, c.ID)
	verdict := "Preserved"
	if !reflect.DeepEqual(ref, mut) {
		verdict = "Violated"
	}
	return runResult{verdict, stableWitness(caseBinding[:], ref, mut)}
}

func stableWitness(binding []byte, ref, mut observation) string {
	payload, _ := json.Marshal(struct {
		Binding           string
		Reference, Mutant observation
	}{hex.EncodeToString(binding), ref, mut})
	d := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(d[:])
}

func orderedIDs(reverse bool) []string {
	ids := make([]string, 0, 70)
	for i := 1; i <= 5; i++ {
		ids = append(ids, fmt.Sprintf("P%02d", i))
	}
	for i := 1; i <= 60; i++ {
		ids = append(ids, fmt.Sprintf("M%02d", i))
	}
	for i := 1; i <= 5; i++ {
		ids = append(ids, fmt.Sprintf("O%02d", i))
	}
	if reverse {
		for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
			ids[i], ids[j] = ids[j], ids[i]
		}
	}
	return ids
}

func permutedIDs() []string {
	ids := orderedIDs(false)
	sort.Slice(ids, func(i, j int) bool {
		a := sha256.Sum256([]byte("781fd3b7" + "05140182" + ids[i]))
		b := sha256.Sum256([]byte("781fd3b7" + "05140182" + ids[j]))
		return string(a[:]) < string(b[:])
	})
	return ids
}

func validRequestFixture() Request {
	return Request{Problem: "choose", Target: "go", Language: "pt-BR", Users: []string{"u2", "u1"}, Capabilities: []string{"search", "rank"}, Constraints: []string{"offline"}, NonGoals: []string{"cloud"}, Stack: []string{"go"}, Licenses: []string{"MIT"}, Locales: []string{"pt-BR"}, MaxComponents: 2, MaxPerQuery: 3, MaxCandidates: 4, Extend: .8, Integration: .5}
}

func errText(err error) string {
	if err == nil {
		return "<nil>"
	}
	return err.Error()
}
func value(v interface{}, err error) string {
	b, _ := json.Marshal(v)
	return string(b) + "|" + errText(err)
}
func files(root string) []string {
	entries, _ := os.ReadDir(root)
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		info, _ := os.Lstat(filepath.Join(root, e.Name()))
		out = append(out, fmt.Sprintf("%s:%s:%04o", e.Name(), info.Mode().Type(), info.Mode().Perm()))
	}
	sort.Strings(out)
	return out
}
func doc(values ...int) []byte {
	b, _ := json.Marshal(map[string]interface{}{"entries": values})
	return b
}

func start(s Surface, key string) (Receipt, error) {
	return s.Start(context.Background(), key, validRequestFixture())
}

func observeMutation(t *testing.T, id, defect string) observation {
	t.Helper()
	root := t.TempDir()
	s := newSubject(root, defect)
	ctx := context.Background()
	o := observation{}
	record := func(v interface{}, err error) { o.Calls = append(o.Calls, value(v, err)) }
	switch id {
	case "M01":
		r, e := start(s, "CANARY-key")
		record(r, e)
		o.Files = files(root)
	case "M02", "M07":
		var wg sync.WaitGroup
		out := make(chan string, 32)
		for i := 0; i < 32; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); r, e := start(s, "same"); out <- value(r, e) }()
		}
		wg.Wait()
		close(out)
		for v := range out {
			o.Calls = append(o.Calls, v)
		}
		sort.Strings(o.Calls)
		record(s.effectsSnapshot(), nil)
		o.Files = files(root)
	case "M03":
		r, e := start(s, "same")
		record(r, e)
		changed := validRequestFixture()
		changed.Capabilities = []string{"different"}
		r, e = s.Start(ctx, "same", changed)
		record(r, e)
		record(s.effectsSnapshot(), nil)
	case "M04":
		r, e := start(s, "retry")
		record(r, e)
		run, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "provider"})
		record(run, e)
		r, e = start(s, "retry")
		record(r, e)
	case "M05":
		a := validRequestFixture()
		b := a
		b.Users = []string{"u1", "u2", "u1"}
		r, e := s.Start(ctx, "canon", a)
		record(r, e)
		r, e = s.Start(ctx, "canon", b)
		record(r, e)
		record(s.effectsSnapshot(), nil)
	case "M06":
		r1, e1 := start(s, "key-a")
		r2, e2 := start(s, "key-b")
		record(r1, e1)
		record(r2, e2)
		g1, _ := s.Get(ctx, r1.RunID)
		g2, _ := s.Get(ctx, r2.RunID)
		record([]string{g1.Hash, g2.Hash}, nil)
	case "M08":
		r, e := start(s, "reservation")
		record(r, e)
		record(s.effectsSnapshot(), nil)
		o.Files = files(root)
	case "M09":
		r, e := s.Start(ctx, "../bad", Request{})
		record(r, e)
		record(s.effectsSnapshot(), nil)
		o.Files = files(root)
	case "M10":
		s.setQueueFull(true)
		r, e := start(s, "full")
		record(r, e)
		record(s.effectsSnapshot(), nil)
		o.Files = files(root)
	case "M11":
		s.setQuotaFull(true)
		r, e := start(s, "quota")
		record(r, e)
		record(s.effectsSnapshot(), nil)
		o.Files = files(root)
	case "M12":
		r, _ := start(s, "status")
		e := s.Execute(ctx, r.RunID)
		record(e == nil, e)
		run, e := s.Get(ctx, r.RunID)
		record(run, e)
	case "M13":
		r, _ := start(s, "state")
		run, e := s.Transition(ctx, r.RunID, r.Version, "SUCCEEDED", Progress{8, 8, "terminal"})
		record(run, e)
	case "M14":
		r, _ := start(s, "terminal")
		_ = s.Execute(ctx, r.RunID)
		run, _ := s.runSnapshot(r.RunID)
		changed, e := s.Transition(ctx, r.RunID, run.Version, "RUNNING", Progress{8, 8, "again"})
		record(changed, e)
	case "M15":
		r, _ := start(s, "immutable")
		run, e := s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1))
		record(run, e)
		run, e = s.Publish(ctx, r.RunID, run.Version, "provenance", doc(2))
		record(run, e)
		page, e := s.Read(ctx, r.RunID, "provenance", "", 10)
		record(page, e)
	case "M16":
		r, _ := start(s, "cas")
		a, e1 := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "a"})
		b, e2 := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 2, "b"})
		record(a, e1)
		record(b, e2)
	case "M17":
		r, _ := start(s, "version")
		run, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "x"})
		record(run, e)
	case "M18":
		r, _ := start(s, "progress")
		for _, p := range []Progress{{7, 1, "bad"}, {8, 9, "bad"}, {8, -1, "bad"}} {
			run, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", p)
			record(run, e)
		}
	case "M19":
		r, e := start(s, "nonblocking")
		record(r, e)
		record(s.effectsSnapshot(), nil)
	case "M20":
		r, _ := start(s, "bounded")
		e := s.Execute(ctx, r.RunID)
		record(e == nil, e)
		record(s.effectsSnapshot(), nil)
	case "M21":
		for _, v := range [][5]int{{0, 1, 1, 1, 1}, {5, 1, 1, 1, 1}, {1, 0, 1, 1, 1}, {1, 33, 1, 1, 1}, {1, 1, 2, 1, 2}} {
			e := s.Configure(v[0], v[1], v[2], v[3], v[4])
			record(v, e)
		}
	case "M22":
		r, _ := start(s, "queued-cancel")
		run, e := s.Cancel(ctx, r.RunID)
		record(run, e)
		e = s.Execute(ctx, r.RunID)
		record(s.effectsSnapshot(), e)
	case "M23":
		r, _ := start(s, "running-cancel")
		run, _ := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "provider"})
		run, e := s.Cancel(ctx, run.ID)
		record(run, e)
	case "M24":
		r, _ := start(s, "cancel-boundary")
		run, _ := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "provider"})
		_, _ = s.Cancel(ctx, run.ID)
		e := s.Execute(ctx, run.ID)
		record(s.effectsSnapshot(), e)
	case "M25":
		r, _ := start(s, "terminal-race")
		e := s.Execute(ctx, r.RunID)
		record(e == nil, e)
		run, _ := s.runSnapshot(r.RunID)
		record(run, nil)
		o.Files = files(root)
	case "M26":
		r, _ := start(s, "signal")
		for i := 0; i < 4; i++ {
			run, e := s.Cancel(ctx, r.RunID)
			record(run, e)
		}
		record(s.effectsSnapshot(), nil)
	case "M27":
		r, _ := start(s, "recover")
		run, _ := s.Publish(ctx, r.RunID, r.Version, "request", doc())
		_, _ = s.Transition(ctx, r.RunID, run.Version, "RUNNING", Progress{8, 1, "provider"})
		e := s.Recover(ctx)
		record(e == nil, e)
		got, e := s.Get(ctx, r.RunID)
		record(got, e)
	case "M28":
		s.recoveryGate = make(chan struct{})
		s.startEntered = make(chan struct{})
		s.startProceed = make(chan struct{})
		recovered := make(chan error, 1)
		started := make(chan string, 1)
		go func() { recovered <- s.Recover(ctx) }()
		go func() { r, e := start(s, "during-recovery"); started <- value(r, e) }()
		<-s.startEntered
		close(s.startProceed)
		early := false
		select {
		case <-started:
			early = true
		case <-time.After(100 * time.Millisecond):
		}
		o.Calls = append(o.Calls, fmt.Sprintf("accepted-before-recovery=%t", early))
		close(s.recoveryGate)
		if !early {
			select {
			case <-started:
			case <-time.After(time.Second):
				t.Fatal("Start did not resume after recovery")
			}
		}
		record(nil, <-recovered)
	case "M29":
		r, _ := start(s, "shutdown")
		run, _ := s.Publish(ctx, r.RunID, r.Version, "request", doc())
		_, _ = s.Transition(ctx, r.RunID, run.Version, "RUNNING", Progress{8, 1, "provider"})
		cctx, cancel := context.WithCancel(ctx)
		cancel()
		record(nil, s.Shutdown(cctx))
		record(nil, s.Recover(ctx))
		got, e := s.Get(ctx, r.RunID)
		record(got, e)
	case "M30":
		r, _ := start(s, "not-ready")
		page, e := s.Read(ctx, r.RunID, "provenance", "", 1)
		record(page, e)
	case "M31":
		r, _ := start(s, "pending")
		run, e := s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1))
		record(run, e)
		page, e := s.Read(ctx, r.RunID, "provenance", "", 1)
		record(page, e)
	case "M32":
		r, _ := start(s, "missing-index")
		_ = os.WriteFile(filepath.Join(root, "blob-"+r.RunID+"-failures"), doc(1), 0o600)
		page, e := s.Read(ctx, r.RunID, "failures", "", 1)
		record(page, e)
	case "M33":
		r, _ := start(s, "corrupt")
		run, _ := s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1))
		_ = os.WriteFile(filepath.Join(root, "blob-"+r.RunID+"-provenance"), doc(9), 0o600)
		page, e := s.Read(ctx, r.RunID, "provenance", "", 1)
		record(page, e)
		record(run, nil)
	case "M34":
		r, _ := start(s, "digest")
		run, e := s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1, 2))
		record(run, e)
	case "M35":
		r, _ := start(s, "terminal-integrity")
		record(nil, s.Execute(ctx, r.RunID))
		run, e := s.Get(ctx, r.RunID)
		record(run, e)
		o.Files = files(root)
	case "M36":
		r, _ := start(s, "missing-terminal")
		_ = s.Execute(ctx, r.RunID)
		_ = os.Remove(filepath.Join(root, "blob-"+r.RunID+"-decision"))
		run, e := s.Get(ctx, r.RunID)
		record(run, e)
	case "M37":
		r1, _ := start(s, "run-a")
		a, _ := s.Publish(ctx, r1.RunID, r1.Version, "provenance", doc(1, 2))
		p, _ := s.Read(ctx, r1.RunID, "provenance", "", 1)
		r2, _ := start(s, "run-b")
		_, _ = s.Publish(ctx, r2.RunID, r2.Version, "provenance", doc(1, 2))
		page, e := s.Read(ctx, r2.RunID, "provenance", p.Next, 1)
		record(page, e)
		record(a, nil)
	case "M38":
		r, _ := start(s, "tamper")
		_, _ = s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1, 2))
		p, _ := s.Read(ctx, r.RunID, "provenance", "", 1)
		tampered := p.Next[:len(p.Next)-1] + "A"
		page, e := s.Read(ctx, r.RunID, "provenance", tampered, 1)
		record(page, e)
	case "M39", "M40":
		r, _ := start(s, "cursor-fields")
		_, _ = s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1, 2))
		p, e := s.Read(ctx, r.RunID, "provenance", "", 1)
		record(p, e)
		parts := strings.Split(p.Next, ".")
		raw, _ := base64.RawURLEncoding.DecodeString(parts[0])
		o.Calls = append(o.Calls, string(raw))
	case "M41":
		record(nil, s.Recover(ctx))
		o.Files = files(root)
	case "M42":
		r, _ := start(s, "json-page")
		_, _ = s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1, 2))
		p, e := s.Read(ctx, r.RunID, "provenance", "", 2)
		record(p, e)
	case "M43":
		r, e := s.Start(ctx, "../CANARY", validRequestFixture())
		record(r, e)
		record(s.effectsSnapshot(), nil)
		o.Files = files(root)
	case "M44":
		outside := filepath.Join(t.TempDir(), "outside")
		_ = os.WriteFile(outside, []byte("sentinel"), 0o600)
		_ = os.Symlink(outside, filepath.Join(root, "blob-link"))
		record(nil, s.Recover(ctx))
		b, _ := os.ReadFile(outside)
		o.Calls = append(o.Calls, string(b))
	case "M45":
		p := filepath.Join(root, "owned")
		_ = os.WriteFile(p, []byte("x"), 0o600)
		_ = os.Link(p, filepath.Join(root, "alias"))
		record(nil, s.Recover(ctx))
		o.Files = files(root)
	case "M46":
		r, e := start(s, "after-rename")
		record(r, e)
		o.Files = files(root)
	case "M47":
		r, _ := start(s, "during-write")
		run, e := s.Transition(ctx, r.RunID, r.Version, "RUNNING", Progress{8, 1, "provider"})
		record(run, e)
		o.Files = files(root)
	case "M48":
		r, _ := start(s, "blob-after-rename")
		run, e := s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1))
		record(run, e)
		page, e := s.Read(ctx, r.RunID, "provenance", "", 1)
		record(page, e)
		o.Files = files(root)
	case "M49":
		s.quota = 8
		r, _ := start(s, "quota-publish")
		run, e := s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1, 2, 3, 4))
		record(run, e)
		o.Files = files(root)
	case "M50":
		old, _ := start(s, "old")
		s.setQuotaFull(true)
		fresh, e := start(s, "new")
		record(fresh, e)
		got, e := s.Get(ctx, old.RunID)
		record(got, e)
		o.Files = files(root)
	case "M51":
		_, e := start(s, "RAW-KEY-CANARY")
		record(nil, e)
		b, _ := os.ReadFile(filepath.Join(root, filesWithPrefix(root, "run-")[0]))
		o.Calls = append(o.Calls, string(b))
		o.Files = files(root)
	case "M52":
		b, e := s.Call(ctx, "projection", nil)
		record(json.RawMessage(b), e)
	case "M53":
		b, e := s.Call(ctx, "error", nil)
		record(string(b), e)
	case "M54":
		b, e := s.Call(ctx, "unknown", nil)
		record(json.RawMessage(b), e)
	case "M55":
		o.Source = astObservation(t, defect, "l1")
	case "M56":
		o.Source = astObservation(t, defect, "l2")
		record(nil, certifiedPredecessor(t))
	case "M57":
		b, e := s.Call(ctx, "legacy", nil)
		record(json.RawMessage(b), e)
		record(nil, certifiedPredecessor(t))
	case "M58":
		r, _ := start(s, "projection-state")
		record(nil, s.Execute(ctx, r.RunID))
		run, e := s.Get(ctx, r.RunID)
		record(run, e)
		record(s.effectsSnapshot(), nil)
	case "M59":
		p := filepath.Join(root, "public")
		_ = os.WriteFile(p, []byte("x"), 0o644)
		fifo := filepath.Join(root, "fifo")
		_ = syscall.Mkfifo(fifo, 0o600)
		record(nil, s.Recover(ctx))
		o.Files = files(root)
	case "M60":
		r, _ := start(s, "corrupt-run")
		_ = os.WriteFile(filepath.Join(root, "run-"+r.RunID), []byte("{"), 0o600)
		run, e := s.Get(ctx, r.RunID)
		record(run, e)
	default:
		t.Fatalf("unbound mutation %s", id)
	}
	return o
}

func filesWithPrefix(root, prefix string) []string {
	entries, _ := os.ReadDir(root)
	var out []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

func astObservation(t *testing.T, defect, layer string) []string {
	t.Helper()
	variant := "reference"
	if defect != "" {
		variant = "mutant"
	}
	dir := filepath.Join("fixture", layer+variant)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var out []string
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".go" {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, imp := range f.Imports {
			out = append(out, strings.Trim(imp.Path.Value, "\""))
		}
		ast.Inspect(f, func(n ast.Node) bool { return true })
	}
	sort.Strings(out)
	return out
}

var predecessorOnce sync.Once
var predecessorErr error

func certifiedPredecessor(t *testing.T) error {
	t.Helper()
	predecessorOnce.Do(func() {
		repo, err := filepath.Abs(filepath.Join("..", "..", "..", ".."))
		if err != nil {
			predecessorErr = err
			return
		}
		args := []string{"test", "02_shell/mcp/server.go", "02_shell/mcp/server_test.go", "02_shell/mcp/recommendation_step0025_contract_test.go", "02_shell/mcp/durable_runs_step0026_r3_test.go", "-run", "TestStep0026R3LegacyBenchmarkAndRecommendationUseIndependentByteGoldens|TestStep0026R3PredecessorServerHashesRemainCertified|TestStep0025LegacySuccessfulSearchMatchesFrozenBaselineBytes", "-count=1"}
		cmd := exec.Command("go", args...)
		cmd.Dir = repo
		cmd.Env = append(os.Environ(), "GOCACHE=/tmp/bateia-go-cache-0026-r5")
		out, err := cmd.CombinedOutput()
		if err != nil {
			predecessorErr = fmt.Errorf("certified predecessor: %w: %s", err, out)
		}
	})
	return predecessorErr
}

func positiveControl(t *testing.T, id string) (observation, bool) {
	t.Helper()
	s := newSubject(t.TempDir(), "")
	ctx := context.Background()
	o := observation{}
	switch id {
	case "P01":
		a := validRequestFixture()
		b := a
		b.Users = []string{"u1", "u2", "u1"}
		r1, e1 := s.Start(ctx, "same", a)
		r2, e2 := s.Start(ctx, "same", b)
		o.Calls = []string{value(r1, e1), value(r2, e2), value(s.effectsSnapshot(), nil)}
		return o, e1 == nil && e2 == nil && r1.RunID == r2.RunID && s.effectsSnapshot().IDs == 1
	case "P02":
		r, e := start(s, "lifecycle")
		if e == nil {
			e = s.Execute(ctx, r.RunID)
		}
		run, e2 := s.Get(ctx, r.RunID)
		o.Calls = []string{value(run, e2)}
		return o, e == nil && e2 == nil && run.Lifecycle == "SUCCEEDED" && run.Status == "INCONCLUSIVE"
	case "P03":
		r, e := start(s, "cursor")
		if e == nil {
			_, e = s.Publish(ctx, r.RunID, r.Version, "provenance", doc(1, 2, 3))
		}
		p1, e1 := s.Read(ctx, r.RunID, "provenance", "", 1)
		p2, e2 := s.Read(ctx, r.RunID, "provenance", p1.Next, 2)
		p3, e3 := s.Read(ctx, r.RunID, "provenance", p1.Next, 2)
		o.Calls = []string{value(p1, e1), value(p2, e2), value(p3, e3)}
		return o, e == nil && e1 == nil && e2 == nil && e3 == nil && reflect.DeepEqual(p2, p3) && len(p1.Entries)+len(p2.Entries) == 3
	case "P04":
		queued, _ := start(s, "queued")
		queuedRun, _ := s.Publish(ctx, queued.RunID, queued.Version, "request", doc())
		queued.Version = queuedRun.Version
		q, _ := s.Cancel(ctx, queued.RunID)
		_ = s.Execute(ctx, queued.RunID)
		running, _ := start(s, "running")
		runningRun, _ := s.Publish(ctx, running.RunID, running.Version, "request", doc())
		rr, _ := s.Transition(ctx, running.RunID, runningRun.Version, "RUNNING", Progress{8, 1, "provider"})
		c1, _ := s.Cancel(ctx, rr.ID)
		before := s.effectsSnapshot()
		c2, _ := s.Cancel(ctx, rr.ID)
		after := s.effectsSnapshot()
		_ = s.Execute(ctx, rr.ID)
		done, _ := s.Get(ctx, rr.ID)
		o.Calls = []string{value(q, nil), value(c1, nil), value(c2, nil), value(done, nil), value(after, nil)}
		return o, q.Lifecycle == "CANCELLED" && c1.Lifecycle == "CANCELLING" && c1.Version == c2.Version && before.Signals == after.Signals && done.Lifecycle == "CANCELLED"
	case "P05":
		legacy, e := s.Call(ctx, "legacy", nil)
		pred := certifiedPredecessor(t)
		r, e2 := start(s, "async")
		if e2 == nil {
			e2 = s.Execute(ctx, r.RunID)
		}
		o.Calls = []string{value(json.RawMessage(legacy), e), value(s.effectsSnapshot(), e2), errText(pred)}
		return o, e == nil && e2 == nil && pred == nil && string(legacy) == `{"tools":["search","benchmark","recommend"]}` && s.effectsSnapshot().Executor == 1
	default:
		t.Fatalf("unbound positive %s", id)
		return o, false
	}
}
