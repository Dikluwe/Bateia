/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r6_gate/discriminatory_step0026_r6_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package variant

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"bateia/04_wiring/r6_gate/conformance"
)

type persisted struct {
	Runs map[string]conformance.Run `json:"runs"`
	Keys map[string]string          `json:"keys"`
}
type Subject struct {
	mu         sync.Mutex
	root       string
	data       persisted
	effects    conformance.Effects
	recovering bool
	secret     []byte
}

func New(root string) conformance.Surface { return NewSubject(root) }
func NewSubject(root string) *Subject {
	_ = os.MkdirAll(root, 0o700)
	_ = os.Chmod(root, 0o700)
	s := &Subject{root: root, data: persisted{Runs: map[string]conformance.Run{}, Keys: map[string]string{}}, secret: []byte("0123456789abcdef0123456789abcdef")}
	s.load()
	return s
}

func canonical(req conformance.Request) string {
	copy := req
	copy.Capabilities = sortedSet(copy.Capabilities)
	copy.Users = sortedSet(copy.Users)
	b, _ := json.Marshal(copy)
	d := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(d[:])
}
func sortedSet(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	j := 0
	for _, v := range out {
		if j == 0 || out[j-1] != v {
			out[j] = v
			j++
		}
	}
	return out[:j]
}
func valid(key string, req conformance.Request) bool {
	return key != "" && len(key) <= 128 && !strings.ContainsAny(key, "/\\") && !strings.Contains(key, "..") && req.Problem != "" && len(req.Capabilities) > 0
}

func (s *Subject) Start(ctx context.Context, key string, req conformance.Request) (conformance.Receipt, error) {
	for {
		s.mu.Lock()
		recovering := s.recovering
		s.mu.Unlock()
		if !recovering {
			break
		}
		select {
		case <-ctx.Done():
			return conformance.Receipt{}, ctx.Err()
		case <-time.After(time.Millisecond):
		}
	}
	return s.StartDirect(key, req, "")
}
func (s *Subject) StartBypassRecovery(_ context.Context, key string, req conformance.Request) (conformance.Receipt, error) {
	return s.StartDirect(key, req, "")
}
func (s *Subject) StartWithID(_ context.Context, key string, req conformance.Request, id string) (conformance.Receipt, error) {
	return s.StartDirect(key, req, id)
}
func (s *Subject) StartIgnoringQuota(ctx context.Context, key string, req conformance.Request) (conformance.Receipt, error) {
	path := filepath.Join(s.root, "quota-full")
	data, _ := os.ReadFile(path)
	_ = os.Remove(path)
	r, e := s.Start(ctx, key, req)
	if data != nil {
		_ = os.WriteFile(path, data, 0o600)
	}
	return r, e
}
func (s *Subject) StartDirect(key string, req conformance.Request, forced string) (conformance.Receipt, error) {
	if !valid(key, req) {
		return conformance.Receipt{}, conformance.ErrInvalid
	}
	if _, e := os.Stat(filepath.Join(s.root, "quota-full")); e == nil {
		return conformance.Receipt{}, conformance.ErrQuota
	}
	hash := canonical(req)
	keyID := digest([]byte(key))
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.data.Keys[keyID]; ok {
		run := s.data.Runs[id]
		if run.Status != hash {
			return conformance.Receipt{}, conformance.ErrConflict
		}
		return conformance.Receipt{RunID: id, Lifecycle: run.Lifecycle, Hash: hash, Version: run.Version}, nil
	}
	id := forced
	if id == "" {
		raw := make([]byte, 16)
		if _, e := rand.Read(raw); e != nil {
			return conformance.Receipt{}, e
		}
		id = "run_" + hex.EncodeToString(raw)
	}
	requestData, _ := json.Marshal(req)
	if e := atomicWrite(s.blob(id, "request"), requestData); e != nil {
		return conformance.Receipt{}, e
	}
	run := conformance.Run{ID: id, Lifecycle: "QUEUED", Status: hash, Version: 1, Progress: conformance.Progress{Planned: 8, Stage: "queued"}, Artifacts: map[string]conformance.Artifact{"request": {Name: "request", State: "PUBLISHED", Hash: digest(requestData), Version: 1}}, TerminalEvents: []string{}}
	s.data.Runs[id] = run
	s.data.Keys[keyID] = id
	if e := s.persist(); e != nil {
		return conformance.Receipt{}, e
	}
	s.effects.Store++
	s.effects.Queue++
	s.effects.IDs++
	return conformance.Receipt{RunID: id, Lifecycle: run.Lifecycle, Hash: hash, Version: run.Version}, nil
}

func (s *Subject) GetRun(_ context.Context, id string) (conformance.Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Run{}, conformance.ErrNotFound
	}
	raw, e := os.ReadFile(filepath.Join(s.root, "run-"+id))
	if e != nil || json.Unmarshal(raw, &conformance.Run{}) != nil {
		return conformance.Run{}, conformance.ErrIntegrity
	}
	if terminal(run.Lifecycle) {
		for _, name := range requiredNames(run.Lifecycle) {
			a, ok := run.Artifacts[name]
			if !ok || a.State != "PUBLISHED" {
				return conformance.Run{}, conformance.ErrIntegrity
			}
			data, e := os.ReadFile(s.blob(id, name))
			if e != nil || digest(data) != a.Hash {
				return conformance.Run{}, conformance.ErrIntegrity
			}
		}
	}
	return cloneRun(run), nil
}
func (s *Subject) RawRun(id string) (conformance.Run, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	return cloneRun(run), ok
}

func (s *Subject) Transition(_ context.Context, id string, expected uint64, lifecycle string, p conformance.Progress) (conformance.Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Run{}, conformance.ErrNotFound
	}
	if run.Version != expected {
		return conformance.Run{}, conformance.ErrCAS
	}
	if p.Planned != 8 || p.Completed < run.Progress.Completed || p.Completed < 0 || p.Completed > 8 {
		return conformance.Run{}, conformance.ErrInvalid
	}
	if run.Lifecycle != "QUEUED" && run.Lifecycle != "RUNNING" && run.Lifecycle != "CANCELLING" {
		return conformance.Run{}, conformance.ErrImmutable
	}
	run.Lifecycle = lifecycle
	run.Progress = p
	run.Version++
	s.data.Runs[id] = run
	if e := s.persist(); e != nil {
		return conformance.Run{}, e
	}
	return cloneRun(run), nil
}
func (s *Subject) TransitionCurrent(ctx context.Context, id, lifecycle string, p conformance.Progress) (conformance.Run, error) {
	run, ok := s.RawRun(id)
	if !ok {
		return conformance.Run{}, conformance.ErrNotFound
	}
	return s.Transition(ctx, id, run.Version, lifecycle, p)
}
func (s *Subject) ForceTransition(id, lifecycle string, p conformance.Progress) (conformance.Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Run{}, conformance.ErrNotFound
	}
	run.Lifecycle = lifecycle
	run.Progress = p
	run.Version++
	s.data.Runs[id] = run
	_ = s.persist()
	return cloneRun(run), nil
}

func (s *Subject) Publish(_ context.Context, id string, expected uint64, name string, data []byte) (conformance.Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Run{}, conformance.ErrNotFound
	}
	if run.Version != expected {
		return conformance.Run{}, conformance.ErrCAS
	}
	if _, ok := run.Artifacts[name]; ok {
		return conformance.Run{}, conformance.ErrImmutable
	}
	if e := atomicWrite(s.blob(id, name), data); e != nil {
		return conformance.Run{}, e
	}
	run.Version++
	run.Artifacts[name] = conformance.Artifact{Name: name, State: "PUBLISHED", Hash: digest(data), Version: run.Version}
	s.data.Runs[id] = run
	if e := s.persist(); e != nil {
		return conformance.Run{}, e
	}
	return cloneRun(run), nil
}
func (s *Subject) ForgetArtifact(id, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run := s.data.Runs[id]
	delete(run.Artifacts, name)
	s.data.Runs[id] = run
}

func (s *Subject) Read(_ context.Context, id, name, cursor string, limit int) (conformance.Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Page{}, conformance.ErrNotFound
	}
	a, ok := s.artifact(run, id, name)
	if !ok || a.State != "PUBLISHED" {
		return conformance.Page{}, conformance.ErrNotReady
	}
	return s.readArtifact(id, name, a, cursor, limit, true)
}
func (s *Subject) ReadUnsigned(id, name, cursor string, limit int) (conformance.Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Page{}, conformance.ErrNotFound
	}
	a, ok := s.artifact(run, id, name)
	if !ok {
		return conformance.Page{}, conformance.ErrNotReady
	}
	return s.readArtifact(id, name, a, cursor, limit, false)
}
func (s *Subject) ReadOrphan(id, name string, limit int) (conformance.Page, error) {
	data, e := os.ReadFile(s.blob(id, name))
	if e != nil {
		return conformance.Page{}, conformance.ErrNotReady
	}
	return decodePage(data, 0, limit), nil
}
func (s *Subject) artifact(run conformance.Run, id, name string) (conformance.Artifact, bool) {
	if a, ok := run.Artifacts[name]; ok {
		return a, true
	}
	raw, e := os.ReadFile(filepath.Join(s.root, "artifact-"+id+"-"+name))
	if e != nil {
		return conformance.Artifact{}, false
	}
	var a conformance.Artifact
	if json.Unmarshal(raw, &a) != nil {
		return conformance.Artifact{}, false
	}
	return a, true
}
func (s *Subject) readArtifact(id, name string, a conformance.Artifact, cursor string, limit int, verify bool) (conformance.Page, error) {
	data, e := os.ReadFile(s.blob(id, name))
	if e != nil {
		return conformance.Page{}, conformance.ErrNotReady
	}
	if digest(data) != a.Hash {
		return conformance.Page{}, conformance.ErrIntegrity
	}
	offset := 0
	if cursor != "" {
		var ok bool
		offset, ok = s.cursorOffset(cursor, id, name, a, verify)
		if !ok {
			return conformance.Page{}, conformance.ErrCursor
		}
	}
	page := decodePage(data, offset, limit)
	var doc struct {
		Entries []json.RawMessage `json:"entries"`
	}
	_ = json.Unmarshal(data, &doc)
	if offset+len(page.Entries) < len(doc.Entries) {
		page.Next = s.sign(id, name, a, offset+len(page.Entries))
	}
	return page, nil
}

func (s *Subject) Cancel(_ context.Context, id string) (conformance.Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Run{}, conformance.ErrNotFound
	}
	if terminal(run.Lifecycle) {
		return cloneRun(run), nil
	}
	run.Lifecycle = "CANCELLED"
	run.Version++
	run.TerminalEvents = append(run.TerminalEvents, "CANCELLED")
	s.data.Runs[id] = run
	_ = s.persist()
	return cloneRun(run), nil
}
func (s *Subject) Complete(_ context.Context, id string) (conformance.Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Run{}, conformance.ErrNotFound
	}
	if terminal(run.Lifecycle) {
		return cloneRun(run), nil
	}
	for _, name := range []string{"request", "recommendation", "decision", "alternatives", "failures", "provenance"} {
		data := []byte(`{"entries":[]}`)
		_ = atomicWrite(s.blob(id, name), data)
		run.Version++
		run.Artifacts[name] = conformance.Artifact{Name: name, State: "PUBLISHED", Hash: digest(data), Version: run.Version}
	}
	run.Lifecycle = "SUCCEEDED"
	run.Status = "INCONCLUSIVE"
	run.Version++
	run.TerminalEvents = append(run.TerminalEvents, "SUCCEEDED")
	s.data.Runs[id] = run
	_ = s.persist()
	return cloneRun(run), nil
}
func (s *Subject) CompleteTwice(id string) (conformance.Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.data.Runs[id]
	if !ok {
		return conformance.Run{}, conformance.ErrNotFound
	}
	run.TerminalEvents = append(run.TerminalEvents, "SUCCEEDED")
	run.Lifecycle = "SUCCEEDED"
	run.Version++
	s.data.Runs[id] = run
	_ = s.persist()
	return cloneRun(run), nil
}

func (s *Subject) Recover(_ context.Context) error {
	s.mu.Lock()
	s.recovering = true
	s.mu.Unlock()
	block := filepath.Join(s.root, "recovery-block")
	if _, e := os.Stat(block); e == nil {
		_ = os.WriteFile(filepath.Join(s.root, "recovery-entered"), []byte("1"), 0o600)
		for {
			if _, e := os.Stat(block); errors.Is(e, os.ErrNotExist) {
				break
			}
			time.Sleep(time.Millisecond)
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, run := range s.data.Runs {
		if !terminal(run.Lifecycle) {
			run.Lifecycle = "INTERRUPTED"
			run.Version++
			s.data.Runs[id] = run
		}
	}
	s.recovering = false
	return s.persist()
}
func (s *Subject) RecoverAs(lifecycle string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, run := range s.data.Runs {
		if !terminal(run.Lifecycle) {
			run.Lifecycle = lifecycle
			run.Version++
			s.data.Runs[id] = run
		}
	}
	return s.persist()
}
func (s *Subject) Shutdown(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, run := range s.data.Runs {
		if !terminal(run.Lifecycle) {
			run.Lifecycle = "CANCELLING"
			run.Version++
			s.data.Runs[id] = run
		}
	}
	return s.persist()
}
func (s *Subject) ShutdownAs(lifecycle string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, run := range s.data.Runs {
		if !terminal(run.Lifecycle) {
			run.Lifecycle = lifecycle
			run.Version++
			s.data.Runs[id] = run
		}
	}
	return s.persist()
}

func (s *Subject) Call(_ context.Context, tool string) ([]byte, error) {
	switch tool {
	case "logs":
		return []byte("run service event"), nil
	case "projection":
		return []byte(`{"entries":[],"failure":"recommendation execution failed"}`), nil
	case "error":
		return nil, errors.New("RUN_STORE_UNAVAILABLE")
	case "empty":
		return []byte(`{"artifacts":[],"failure":null}`), nil
	case "legacy":
		return []byte(`{"tools":["search","benchmark","recommend"]}`), nil
	}
	return nil, conformance.ErrInvalid
}
func (s *Subject) ValidateRoot() error {
	return filepath.WalkDir(s.root, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == s.root {
			return nil
		}
		info, statErr := os.Lstat(path)
		if statErr != nil {
			return statErr
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return conformance.ErrIntegrity
		}
		if st, ok := info.Sys().(*syscall.Stat_t); ok && st.Nlink != 1 {
			return conformance.ErrIntegrity
		}
		if info.Mode().Perm()&0o077 != 0 {
			return conformance.ErrIntegrity
		}
		if !info.Mode().IsRegular() && !info.IsDir() {
			return conformance.ErrIntegrity
		}
		return nil
	})
}
func (s *Subject) ValidateRootAllowSymlink() error {
	return filepath.WalkDir(s.root, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == s.root {
			return nil
		}
		info, e := os.Lstat(path)
		if e != nil {
			return e
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		if st, ok := info.Sys().(*syscall.Stat_t); ok && st.Nlink != 1 {
			return conformance.ErrIntegrity
		}
		if info.Mode().Perm()&0o077 != 0 || (!info.Mode().IsRegular() && !info.IsDir()) {
			return conformance.ErrIntegrity
		}
		return nil
	})
}
func (s *Subject) ValidateRootAllowHardlink() error {
	return filepath.WalkDir(s.root, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == s.root {
			return nil
		}
		info, e := os.Lstat(path)
		if e != nil {
			return e
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return conformance.ErrIntegrity
		}
		if info.Mode().Perm()&0o077 != 0 || (!info.Mode().IsRegular() && !info.IsDir()) {
			return conformance.ErrIntegrity
		}
		return nil
	})
}
func (s *Subject) ValidateRootAllowUnsafeKind() error {
	return filepath.WalkDir(s.root, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == s.root {
			return nil
		}
		info, e := os.Lstat(path)
		if e != nil {
			return e
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return conformance.ErrIntegrity
		}
		if st, ok := info.Sys().(*syscall.Stat_t); ok && st.Nlink != 1 {
			return conformance.ErrIntegrity
		}
		return nil
	})
}

func (s *Subject) CrashCreate() ([]conformance.CrashObservation, error) {
	return crashMatrix(filepath.Join(s.root, "crash-create"), false)
}
func (s *Subject) CrashTransition() ([]conformance.CrashObservation, error) {
	return crashMatrix(filepath.Join(s.root, "crash-transition"), false)
}
func (s *Subject) CrashPublish() ([]conformance.CrashObservation, error) {
	return crashMatrix(filepath.Join(s.root, "crash-publish"), false)
}
func (s *Subject) UnsafeCrashCreate() ([]conformance.CrashObservation, error) {
	return crashMatrix(filepath.Join(s.root, "unsafe-create"), true)
}
func (s *Subject) UnsafeCrashTransition() ([]conformance.CrashObservation, error) {
	return crashMatrix(filepath.Join(s.root, "unsafe-transition"), true)
}
func (s *Subject) UnsafeCrashPublish() ([]conformance.CrashObservation, error) {
	return crashMatrix(filepath.Join(s.root, "unsafe-publish"), true)
}
func crashMatrix(base string, unsafe bool) ([]conformance.CrashObservation, error) {
	boundaries := []string{"open", "write", "file-sync", "rename", "dir-sync"}
	out := make([]conformance.CrashObservation, 0, len(boundaries))
	for i, b := range boundaries {
		dir := base + "-" + b
		_ = os.MkdirAll(dir, 0o700)
		final := filepath.Join(dir, "record")
		tmp := final + ".tmp"
		f, e := os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if e != nil {
			return nil, e
		}
		confirmed := false
		if i > 0 {
			_, _ = f.Write([]byte(`{"ok":true}`))
		}
		if i > 1 {
			_ = f.Sync()
		}
		_ = f.Close()
		if i > 2 {
			_ = os.Rename(tmp, final)
		}
		if i > 3 {
			d, _ := os.Open(dir)
			_ = d.Sync()
			_ = d.Close()
			confirmed = true
		}
		if unsafe && i >= 1 {
			if _, e := os.Stat(tmp); e == nil {
				_ = os.Rename(tmp, final)
			}
			confirmed = false
		} else if !confirmed {
			_ = os.Remove(tmp)
			_ = os.Remove(final)
		}
		_, finalErr := os.Stat(final)
		_, tempErr := os.Stat(tmp)
		reopen := "absent"
		if finalErr == nil {
			raw, _ := os.ReadFile(final)
			if json.Valid(raw) {
				reopen = "valid"
			} else {
				reopen = "ambiguous"
			}
		}
		out = append(out, conformance.CrashObservation{Boundary: b, Confirmed: confirmed, Final: finalErr == nil, Temp: tempErr == nil, Reopen: reopen})
	}
	return out, nil
}

func (s *Subject) Effects() conformance.Effects { s.mu.Lock(); defer s.mu.Unlock(); return s.effects }
func (s *Subject) blob(id, name string) string  { return filepath.Join(s.root, "blob-"+id+"-"+name) }
func (s *Subject) persist() error {
	data, _ := json.Marshal(s.data)
	if e := atomicWrite(filepath.Join(s.root, "state"), data); e != nil {
		return e
	}
	for id, run := range s.data.Runs {
		b, _ := json.Marshal(run)
		if e := atomicWrite(filepath.Join(s.root, "run-"+id), b); e != nil {
			return e
		}
	}
	return nil
}
func (s *Subject) load() {
	raw, e := os.ReadFile(filepath.Join(s.root, "state"))
	if e == nil {
		_ = json.Unmarshal(raw, &s.data)
	}
}
func atomicWrite(path string, data []byte) error {
	tmp := path + ".tmp"
	_ = os.Remove(tmp)
	f, e := os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if e != nil {
		return e
	}
	if _, e = f.Write(data); e != nil {
		_ = f.Close()
		return e
	}
	if e = f.Sync(); e != nil {
		_ = f.Close()
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	if e = os.Rename(tmp, path); e != nil {
		return e
	}
	d, e := os.Open(filepath.Dir(path))
	if e != nil {
		return e
	}
	defer d.Close()
	return d.Sync()
}
func digest(data []byte) string {
	d := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(d[:])
}
func terminal(v string) bool {
	return v == "SUCCEEDED" || v == "FAILED" || v == "CANCELLED" || v == "INTERRUPTED"
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
func cloneRun(run conformance.Run) conformance.Run {
	artifacts := run.Artifacts
	run.Artifacts = map[string]conformance.Artifact{}
	for k, v := range artifacts {
		run.Artifacts[k] = v
	}
	run.TerminalEvents = append([]string(nil), run.TerminalEvents...)
	return run
}
func decodePage(data []byte, offset, limit int) conformance.Page {
	var doc struct {
		Entries []json.RawMessage `json:"entries"`
	}
	_ = json.Unmarshal(data, &doc)
	end := offset + limit
	if end > len(doc.Entries) {
		end = len(doc.Entries)
	}
	return conformance.Page{Entries: append([]json.RawMessage(nil), doc.Entries[offset:end]...)}
}
func (s *Subject) sign(id, name string, a conformance.Artifact, offset int) string {
	raw, _ := json.Marshal(map[string]interface{}{"run_id": id, "artifact": name, "version": a.Version, "hash": a.Hash, "offset": offset})
	enc := base64.RawURLEncoding.EncodeToString(raw)
	m := hmac.New(sha256.New, s.secret)
	_, _ = m.Write([]byte(enc))
	return enc + "." + base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}
func (s *Subject) cursorOffset(cursor, id, name string, a conformance.Artifact, verify bool) (int, bool) {
	parts := strings.Split(cursor, ".")
	if len(parts) != 2 {
		return 0, false
	}
	if verify {
		sig, e := base64.RawURLEncoding.DecodeString(parts[1])
		if e != nil {
			return 0, false
		}
		m := hmac.New(sha256.New, s.secret)
		_, _ = m.Write([]byte(parts[0]))
		if !hmac.Equal(sig, m.Sum(nil)) {
			return 0, false
		}
	}
	raw, e := base64.RawURLEncoding.DecodeString(parts[0])
	if e != nil {
		return 0, false
	}
	var p struct {
		RunID          string `json:"run_id"`
		Artifact, Hash string
		Version        uint64
		Offset         int
	}
	if json.Unmarshal(raw, &p) != nil || p.RunID != id || p.Artifact != name || p.Hash != a.Hash || p.Version != a.Version {
		return 0, false
	}
	return p.Offset, true
}
