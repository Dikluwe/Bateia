/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r5_gate/discriminatory_step0026_r5_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package r5gate

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
)

var (
	errInvalid   = errors.New("Invalid params")
	errConflict  = errors.New("IDEMPOTENCY_CONFLICT")
	errQueueFull = errors.New("QUEUE_FULL")
	errQuota     = errors.New("QUOTA_EXCEEDED")
	errNotFound  = errors.New("RUN_NOT_FOUND")
	errNotReady  = errors.New("ARTIFACT_NOT_READY")
	errIntegrity = errors.New("INTEGRITY_ERROR")
	errCursor    = errors.New("CURSOR_INVALID")
	errCAS       = errors.New("VERSION_CONFLICT")
	errImmutable = errors.New("ARTIFACT_IMMUTABLE")
	errStore     = errors.New("RUN_STORE_UNAVAILABLE")
)

type Request struct {
	Problem, Target, Language                                            string
	Users, Capabilities, Constraints, NonGoals, Stack, Licenses, Locales []string
	MaxComponents, MaxPerQuery, MaxCandidates                            int
	Extend, Integration                                                  float64
}

type Receipt struct {
	RunID, Lifecycle string
	Version          uint64
}
type Progress struct {
	Planned, Completed int
	Stage              string
}
type Artifact struct {
	Name, State, Hash string
	Size, Entries     int
	Version           uint64
}
type Run struct {
	ID, Hash, Lifecycle, Status, Failure string
	Version                              uint64
	Progress                             Progress
	Artifacts                            map[string]Artifact
}
type Page struct {
	Entries []json.RawMessage
	Next    string
}
type Effects struct{ Store, Queue, IDs, Executor, Signals int }

type RunService interface {
	Start(context.Context, string, Request) (Receipt, error)
	Get(context.Context, string) (Run, error)
	Cancel(context.Context, string) (Run, error)
	Execute(context.Context, string) error
	Recover(context.Context) error
	Shutdown(context.Context) error
}

type RecommendationRunStore interface {
	Transition(context.Context, string, uint64, string, Progress) (Run, error)
	Publish(context.Context, string, uint64, string, []byte) (Run, error)
	Read(context.Context, string, string, string, int) (Page, error)
}

type MCP interface {
	Call(context.Context, string, map[string]interface{}) ([]byte, error)
}

type Configurable interface {
	Configure(int, int, int, int, int) error
}
type Surface interface {
	RunService
	RecommendationRunStore
	MCP
	Configurable
}

type runtimeSubject struct {
	mu                   sync.Mutex
	root, defect         string
	runs                 map[string]Run
	idem                 map[string]string
	effects              Effects
	queueFull, quotaFull bool
	quota, used          int
	secret               []byte
	next                 atomic.Uint64
	recoveryGate         chan struct{}
	startEntered         chan struct{}
	startProceed         chan struct{}
	logs                 strings.Builder
}

func newSubject(root, defect string) *runtimeSubject {
	_ = os.MkdirAll(root, 0o700)
	_ = os.Chmod(root, 0o700)
	s := &runtimeSubject{root: root, defect: defect, runs: map[string]Run{}, idem: map[string]string{}, quota: 1 << 20, secret: []byte("0123456789abcdef0123456789abcdef")}
	mode := os.FileMode(0o600)
	if defect == "M41" {
		mode = 0o644
	}
	_ = os.WriteFile(filepath.Join(root, "cursor-secret"), s.secret, mode)
	return s
}

func (s *runtimeSubject) Configure(workers, capacity, maxArtifact, maxRun, maxTotal int) error {
	valid := workers >= 1 && workers <= 4 && capacity >= 1 && capacity <= 32 && maxArtifact > 0 && maxRun >= maxArtifact && maxTotal >= maxRun
	if !valid && s.defect != "M21" {
		return errInvalid
	}
	return nil
}

func canonical(req Request) ([]byte, string) {
	norm := req
	for _, target := range []*[]string{&norm.Users, &norm.Capabilities, &norm.Constraints, &norm.NonGoals, &norm.Stack, &norm.Licenses, &norm.Locales} {
		values := append([]string(nil), (*target)...)
		sort.Strings(values)
		out := values[:0]
		for _, v := range values {
			if len(out) == 0 || out[len(out)-1] != v {
				out = append(out, v)
			}
		}
		*target = out
	}
	b, _ := json.Marshal(norm)
	sum := sha256.Sum256(b)
	return b, "sha256:" + hex.EncodeToString(sum[:])
}

func validRequest(req Request) bool {
	return strings.TrimSpace(req.Problem) != "" && len(req.Capabilities) > 0 && req.MaxComponents >= 2 && req.Extend > 0 && req.Extend <= 1 && req.Integration > 0 && req.Integration <= req.Extend && req.MaxPerQuery >= 1 && req.MaxPerQuery <= 25 && req.MaxCandidates >= 1 && req.MaxCandidates <= 20
}
func validKey(key string) bool {
	if len(key) < 1 || len(key) > 128 || strings.Contains(key, "..") || strings.ContainsAny(key, "/\\") {
		return false
	}
	for _, r := range key {
		if r < 0x20 || r > 0x7e {
			return false
		}
	}
	return true
}

func (s *runtimeSubject) Start(_ context.Context, key string, req Request) (Receipt, error) {
	if s.startEntered != nil {
		s.startEntered <- struct{}{}
		if s.startProceed != nil {
			<-s.startProceed
		}
	}
	if s.recoveryGate != nil && s.defect != "M28" {
		<-s.recoveryGate
	}
	if s.defect != "M09" && s.defect != "M43" && (!validKey(key) || !validRequest(req)) {
		return Receipt{}, errInvalid
	}
	if s.defect == "M09" {
		s.effects.Store++
		if !validKey(key) || !validRequest(req) {
			return Receipt{}, errInvalid
		}
	}
	_, hash := canonical(req)
	if s.defect == "M05" {
		raw, _ := json.Marshal(req)
		d := sha256.Sum256(raw)
		hash = "sha256:" + hex.EncodeToString(d[:])
	}
	if s.defect == "M06" {
		d := sha256.Sum256([]byte(req.Problem + key))
		hash = "sha256:" + hex.EncodeToString(d[:])
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if prior, ok := s.idem[key]; ok && s.defect != "M02" && s.defect != "M07" {
		run := s.runs[prior]
		if run.Hash != hash && s.defect != "M03" {
			return Receipt{}, errConflict
		}
		if s.defect == "M04" {
			return Receipt{prior, "QUEUED", 1}, nil
		}
		return Receipt{prior, run.Lifecycle, run.Version}, nil
	}
	if s.queueFull && s.defect != "M10" {
		return Receipt{}, errQueueFull
	}
	if s.quotaFull && s.defect != "M11" && s.defect != "M50" {
		return Receipt{}, errQuota
	}
	if s.quotaFull && s.defect == "M50" {
		for old := range s.runs {
			delete(s.runs, old)
			for key, id := range s.idem {
				if id == old {
					delete(s.idem, key)
				}
			}
			break
		}
		s.quotaFull = false
	}
	s.effects.Queue++
	s.effects.IDs++
	n := s.next.Add(1)
	seed := sha256.Sum256([]byte(fmt.Sprintf("entropy:%d", n)))
	id := "run_" + hex.EncodeToString(seed[:16])
	if s.defect == "M01" {
		d := sha256.Sum256([]byte(key))
		id = "run_" + hex.EncodeToString(d[:16])
	}
	run := Run{ID: id, Hash: hash, Lifecycle: "QUEUED", Version: 1, Progress: Progress{8, 0, "queued"}, Artifacts: map[string]Artifact{}}
	if s.defect == "M11" && s.quotaFull {
		s.effects.Executor++
		s.runs[id] = run
		s.idem[key] = id
		s.effects.Store++
		return Receipt{id, "QUEUED", 1}, nil
	}
	if err := s.persistPair(key, run); err != nil {
		if s.defect == "M46" {
			s.runs[id] = run
			s.idem[key] = id
			return Receipt{id, "QUEUED", 1}, nil
		}
		return Receipt{}, err
	}
	s.runs[id] = run
	s.idem[key] = id
	s.effects.Store++
	if s.defect == "M08" {
		return Receipt{id, "QUEUED", 1}, errQueueFull
	}
	if s.defect == "M19" {
		s.effects.Executor++
		return Receipt{id, "RUNNING", 2}, nil
	}
	return Receipt{id, "QUEUED", 1}, nil
}

func (s *runtimeSubject) persistPair(key string, run Run) error {
	runData, _ := json.Marshal(run)
	keyData, _ := json.Marshal(map[string]string{"key": key, "run_id": run.ID})
	if s.defect == "M51" {
		runData = append(runData, []byte(" RAW-KEY Bearer TOKEN Cookie=SECRET HEADER ENV CURSOR ADAPTER THIRD-PARTY")...)
	}
	if s.defect == "M46" {
		_ = s.atomicWrite("run-"+run.ID, runData, "after-rename")
		return errStore
	}
	if err := s.atomicWrite("run-"+run.ID, runData, ""); err != nil {
		return err
	}
	return s.atomicWrite("key-"+digestName(key), keyData, "")
}

func digestName(v string) string { d := sha256.Sum256([]byte(v)); return hex.EncodeToString(d[:]) }

func (s *runtimeSubject) atomicWrite(name string, data []byte, fail string) error {
	final := filepath.Join(s.root, name)
	tmp := final + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	part := data
	if fail == "during-write" {
		part = data[:len(data)/2]
	}
	if _, err = f.Write(part); err != nil {
		_ = f.Close()
		return err
	}
	if fail == "during-write" {
		_ = f.Close()
		return errStore
	}
	if err = f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = os.Rename(tmp, final); err != nil {
		return err
	}
	if fail == "after-rename" {
		return errStore
	}
	d, err := os.Open(s.root)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}

func (s *runtimeSubject) Get(_ context.Context, id string) (Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		if s.defect == "M60" {
			run = Run{ID: id, Lifecycle: "INTERRUPTED", Status: "COMPLETE", Artifacts: map[string]Artifact{}}
			s.runs[id] = run
			return run, nil
		}
		return Run{}, errNotFound
	}
	persisted, err := os.ReadFile(filepath.Join(s.root, "run-"+id))
	if err != nil || json.Unmarshal(persisted, &Run{}) != nil {
		if s.defect == "M60" {
			return Run{}, errNotFound
		}
		return Run{}, errIntegrity
	}
	if terminal(run.Lifecycle) && s.defect != "M36" {
		for _, name := range requiredNames(run.Lifecycle) {
			a, ok := run.Artifacts[name]
			if !ok || a.State != "PUBLISHED" {
				return Run{}, errIntegrity
			}
			if _, err := os.Stat(filepath.Join(s.root, "blob-"+id+"-"+name)); err != nil {
				return Run{}, errIntegrity
			}
		}
	}
	return s.publicRun(run), nil
}
func requiredNames(lifecycle string) []string {
	if lifecycle == "SUCCEEDED" {
		return []string{"request", "recommendation", "decision", "alternatives", "failures", "provenance"}
	}
	if lifecycle == "FAILED" {
		return []string{"request", "failures"}
	}
	return []string{"request"}
}

func (s *runtimeSubject) publicRun(run Run) Run {
	if s.defect == "M12" && run.Lifecycle == "SUCCEEDED" {
		run.Status = "COMPLETE"
	}
	if s.defect == "M54" {
		run.Artifacts = nil
	}
	return run
}

func (s *runtimeSubject) Transition(_ context.Context, id string, expected uint64, lifecycle string, p Progress) (Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return Run{}, errNotFound
	}
	if run.Version != expected && s.defect != "M16" {
		return Run{}, errCAS
	}
	if terminal(run.Lifecycle) && s.defect != "M14" {
		return Run{}, errImmutable
	}
	if !allowed(run.Lifecycle, lifecycle, p, run.Progress) && s.defect != "M13" && s.defect != "M18" {
		return Run{}, errInvalid
	}
	run.Lifecycle, run.Progress = lifecycle, p
	if s.defect != "M17" {
		run.Version++
	}
	if s.defect == "M47" {
		_ = s.atomicWrite("run-"+id, []byte("partial"), "during-write")
		run.Lifecycle = "FAILED"
	}
	s.runs[id] = run
	return run, nil
}
func terminal(v string) bool {
	return v == "SUCCEEDED" || v == "FAILED" || v == "CANCELLED" || v == "INTERRUPTED"
}
func allowed(from, to string, p, prior Progress) bool {
	if p.Planned != 8 || p.Completed < prior.Completed || p.Completed < 0 || p.Completed > 8 {
		return false
	}
	if from == "QUEUED" {
		return to == "RUNNING" || to == "CANCELLED" || to == "INTERRUPTED"
	}
	if from == "RUNNING" {
		return to == "RUNNING" || to == "CANCELLING" || to == "SUCCEEDED" || to == "FAILED" || to == "INTERRUPTED"
	}
	if from == "CANCELLING" {
		return to == "CANCELLED" || to == "INTERRUPTED"
	}
	return false
}

func (s *runtimeSubject) Publish(_ context.Context, id string, expected uint64, name string, document []byte) (Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return Run{}, errNotFound
	}
	if run.Version != expected && s.defect != "M16" {
		return Run{}, errCAS
	}
	if _, ok := run.Artifacts[name]; ok && s.defect != "M15" {
		return Run{}, errImmutable
	}
	if s.used+len(document) > s.quota && s.defect != "M49" {
		return Run{}, errQuota
	}
	sum := sha256.Sum256(document)
	hash := "sha256:" + hex.EncodeToString(sum[:])
	if s.defect == "M34" {
		d := sha256.Sum256([]byte("entries"))
		hash = "sha256:" + hex.EncodeToString(d[:])
	}
	blob := "blob-" + id + "-" + name
	if s.defect == "M48" {
		_ = s.atomicWrite(blob, document, "after-rename")
		return run, errStore
	}
	if err := s.atomicWrite(blob, document, ""); err != nil {
		return Run{}, err
	}
	s.used += len(document)
	state := "PUBLISHED"
	if s.defect == "M31" {
		state = "PENDING"
	}
	run.Version++
	run.Artifacts[name] = Artifact{name, state, hash, len(document), countEntries(document), run.Version}
	s.runs[id] = run
	return run, nil
}
func countEntries(document []byte) int {
	var d struct {
		Entries []json.RawMessage `json:"entries"`
	}
	_ = json.Unmarshal(document, &d)
	return len(d.Entries)
}

func (s *runtimeSubject) Read(_ context.Context, id, name, cursor string, limit int) (Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return Page{}, errNotFound
	}
	a, ok := run.Artifacts[name]
	if (!ok || a.State != "PUBLISHED") && s.defect != "M30" && s.defect != "M31" && s.defect != "M32" && s.defect != "M48" {
		if terminal(run.Lifecycle) && required(run.Lifecycle, name) {
			return Page{}, errIntegrity
		}
		return Page{}, errNotReady
	}
	path := filepath.Join(s.root, "blob-"+id+"-"+name)
	data, err := os.ReadFile(path)
	if err != nil {
		if s.defect == "M30" {
			return Page{Entries: []json.RawMessage{}}, nil
		}
		return Page{}, errNotReady
	}
	sum := sha256.Sum256(data)
	if a.Hash != "" && a.Hash != "sha256:"+hex.EncodeToString(sum[:]) && s.defect != "M33" {
		return Page{}, errIntegrity
	}
	var doc struct {
		Entries []json.RawMessage `json:"entries"`
	}
	if json.Unmarshal(data, &doc) != nil {
		return Page{}, errIntegrity
	}
	offset := 0
	if cursor != "" {
		var e error
		offset, e = s.verifyCursor(cursor, id, name, a)
		if e != nil {
			return Page{}, e
		}
	}
	if offset < 0 || offset > len(doc.Entries) {
		return Page{}, errCursor
	}
	end := offset + limit
	if end > len(doc.Entries) {
		end = len(doc.Entries)
	}
	page := Page{Entries: append([]json.RawMessage(nil), doc.Entries[offset:end]...)}
	if s.defect == "M42" {
		for i := range page.Entries {
			encoded, _ := json.Marshal(string(page.Entries[i]))
			page.Entries[i] = encoded
		}
	}
	if end < len(doc.Entries) {
		page.Next = s.signCursor(id, name, a, end, limit)
	}
	return page, nil
}
func required(lifecycle, name string) bool {
	if lifecycle == "SUCCEEDED" {
		return true
	}
	if lifecycle == "FAILED" {
		return name == "request" || name == "failures"
	}
	return name == "request"
}

func (s *runtimeSubject) signCursor(id, name string, a Artifact, offset, limit int) string {
	payload := map[string]interface{}{"schema": "1", "run_id": id, "artifact": name, "version": a.Version, "hash": a.Hash, "offset": offset}
	if s.defect == "M39" {
		payload["limit"] = limit
	}
	if s.defect == "M40" {
		delete(payload, "version")
		delete(payload, "hash")
	}
	raw, _ := json.Marshal(payload)
	enc := base64.RawURLEncoding.EncodeToString(raw)
	if s.defect == "M38" {
		return enc + ".unsigned"
	}
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(enc))
	return enc + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
func (s *runtimeSubject) verifyCursor(cursor, id, name string, a Artifact) (int, error) {
	parts := strings.Split(cursor, ".")
	if len(parts) != 2 {
		return 0, errCursor
	}
	if s.defect != "M38" {
		sig, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return 0, errCursor
		}
		mac := hmac.New(sha256.New, s.secret)
		_, _ = mac.Write([]byte(parts[0]))
		if !hmac.Equal(sig, mac.Sum(nil)) {
			return 0, errCursor
		}
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, errCursor
	}
	var p struct {
		RunID    string `json:"run_id"`
		Artifact string `json:"artifact"`
		Hash     string `json:"hash"`
		Version  uint64 `json:"version"`
		Offset   int    `json:"offset"`
		Limit    int    `json:"limit"`
	}
	if json.Unmarshal(raw, &p) != nil {
		return 0, errCursor
	}
	if s.defect != "M37" && (p.RunID != id || p.Artifact != name) {
		return 0, errCursor
	}
	if s.defect != "M40" && (p.Version != a.Version || p.Hash != a.Hash) {
		return 0, errCursor
	}
	return p.Offset, nil
}

func (s *runtimeSubject) Cancel(_ context.Context, id string) (Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return Run{}, errNotFound
	}
	if terminal(run.Lifecycle) || run.Lifecycle == "CANCELLING" {
		if s.defect == "M26" {
			run.Version++
			s.effects.Signals++
			s.runs[id] = run
		}
		return run, nil
	}
	if run.Lifecycle == "QUEUED" {
		run.Lifecycle = "CANCELLED"
	} else {
		run.Lifecycle = "CANCELLING"
		if s.defect == "M23" {
			run.Lifecycle = "CANCELLED"
		}
	}
	run.Version++
	s.effects.Signals++
	s.runs[id] = run
	return run, nil
}

func (s *runtimeSubject) Execute(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return errNotFound
	}
	if terminal(run.Lifecycle) {
		if s.defect == "M22" {
			s.effects.Executor++
		}
		return nil
	}
	s.effects.Executor++
	if s.defect == "M20" {
		s.effects.Executor += 4
	}
	if run.Lifecycle == "CANCELLING" {
		if s.defect == "M24" {
			s.effects.Executor++
		}
		run.Lifecycle = "CANCELLED"
		run.Version++
		s.runs[id] = run
		return nil
	}
	run.Lifecycle = "RUNNING"
	run.Version++
	if s.defect == "M25" {
		run.Artifacts["recommendation"] = Artifact{Name: "recommendation", State: "PUBLISHED"}
		run.Lifecycle = "CANCELLED"
	} else {
		if s.defect != "M35" {
			for _, name := range requiredNames("SUCCEEDED") {
				doc := []byte(`{"entries":[]}`)
				sum := sha256.Sum256(doc)
				blob := "blob-" + id + "-" + name
				_ = s.atomicWrite(blob, doc, "")
				run.Version++
				run.Artifacts[name] = Artifact{name, "PUBLISHED", "sha256:" + hex.EncodeToString(sum[:]), len(doc), 0, run.Version}
			}
		}
		run.Lifecycle = "SUCCEEDED"
		run.Status = "INCONCLUSIVE"
		if s.defect == "M58" {
			run.Status = "COMPLETE"
			s.effects.Executor++
		}
	}
	run.Version++
	s.runs[id] = run
	return nil
}

func (s *runtimeSubject) Recover(_ context.Context) error {
	if s.recoveryGate != nil {
		<-s.recoveryGate
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, run := range s.runs {
		if !terminal(run.Lifecycle) {
			if s.defect == "M27" {
				run.Lifecycle = "FAILED"
			} else {
				run.Lifecycle = "INTERRUPTED"
			}
			run.Progress.Stage = "terminal"
			run.Status = ""
			run.Version++
			s.runs[id] = run
		}
	}
	return s.validateFilesystem()
}
func (s *runtimeSubject) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	for id, run := range s.runs {
		if !terminal(run.Lifecycle) {
			run.Lifecycle = "CANCELLING"
			if s.defect == "M29" {
				run.Lifecycle = "CANCELLED"
			}
			s.runs[id] = run
		}
	}
	s.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *runtimeSubject) validateFilesystem() error {
	return filepath.WalkDir(s.root, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if path == s.root {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			if s.defect == "M44" {
				return nil
			}
			return errIntegrity
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && s.defect != "M45" && stat.Nlink != 1 {
			return errIntegrity
		}
		if s.defect != "M59" && info.Mode().Perm()&0o077 != 0 {
			return errIntegrity
		}
		if !info.IsDir() && !info.Mode().IsRegular() && s.defect != "M59" {
			return errIntegrity
		}
		return nil
	})
}

func (s *runtimeSubject) Call(ctx context.Context, tool string, args map[string]interface{}) ([]byte, error) {
	switch tool {
	case "start":
		key, _ := args["idempotency_key"].(string)
		req, _ := args["request"].(Request)
		receipt, err := s.Start(ctx, key, req)
		if err != nil {
			return nil, err
		}
		return json.Marshal(receipt)
	case "get":
		id, _ := args["run_id"].(string)
		run, err := s.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		return json.Marshal(run)
	case "legacy":
		if s.defect == "M57" {
			return []byte(`{"tools":["changed","search","benchmark"]}`), nil
		}
		return []byte(`{"tools":["search","benchmark","recommend"]}`), nil
	case "projection":
		if s.defect == "M52" {
			return []byte(`{"entries":["README CODE TEST FULL-LICENSE BUNDLE SNIPPET RAW-PROVIDER"]}`), nil
		}
		return []byte(`{"entries":[],"failure":"recommendation execution failed"}`), nil
	case "error":
		if s.defect == "M53" {
			return nil, errors.New("open /private/path: Bearer TOKEN raw-key cursor provider")
		}
		return nil, errStore
	default:
		if s.defect == "M54" {
			return []byte(`{"artifacts":null,"extra":true}`), nil
		}
		return nil, errInvalid
	}
}

func (s *runtimeSubject) effectsSnapshot() Effects {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.effects
}
func (s *runtimeSubject) runSnapshot(id string) (Run, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.runs[id]
	return r, ok
}
func (s *runtimeSubject) setQueueFull(v bool) { s.queueFull = v }
func (s *runtimeSubject) setQuotaFull(v bool) { s.quotaFull = v }
