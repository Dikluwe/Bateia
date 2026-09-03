/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runstore/file_store_step0026_r2_test.md
 * @layer L3
 * @updated 2026-09-01
 */
package runstore

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"testing"

	"bateia/01_core/domain"
)

func TestStep0026R2FileStoreCreateOrResolveAndCASAreActuallyConcurrent(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	defer store.Close()
	const contenders = 32
	start := make(chan struct{})
	results := make(chan domain.RecommendationRun, contenders)
	wins := make(chan bool, contenders)
	errs := make(chan error, contenders)
	var wg sync.WaitGroup
	for index := 0; index < contenders; index++ {
		index := index
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			id := domain.RunID(fmt.Sprintf("run_%032x", index+1))
			run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("f", 64))
			got, won, err := store.CreateOrResolve(ctx, run, record)
			results <- got
			wins <- won
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	close(wins)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("CreateOrResolve race: %v", err)
		}
	}
	winnerIDs := map[domain.RunID]bool{}
	for result := range results {
		winnerIDs[result.RunID] = true
	}
	winnerCount := 0
	for won := range wins {
		if won {
			winnerCount++
		}
	}
	if winnerCount != 1 || len(winnerIDs) != 1 {
		t.Fatalf("atomic create winners=%d ids=%v", winnerCount, winnerIDs)
	}
	var winner domain.RunID
	for id := range winnerIDs {
		winner = id
	}
	run, err := store.Transition(ctx, domain.RunTransition{RunID: winner, ExpectedVersion: 1, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: domain.RecommendationRunProgress{Planned: 8}})
	if err != nil {
		t.Fatal(err)
	}
	race := make(chan struct{})
	transitionErrs := make(chan error, 4)
	for index, stage := range []domain.RecommendationRunStage{domain.RunStageDiscover, domain.RunStageScreen, domain.RunStageCollect, domain.RunStageAssess} {
		index, stage := index, stage
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-race
			_, err := store.Transition(ctx, domain.RunTransition{RunID: winner, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: stage, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: index + 1}})
			transitionErrs <- err
		}()
	}
	close(race)
	wg.Wait()
	close(transitionErrs)
	success, conflict := 0, 0
	for err := range transitionErrs {
		if err == nil {
			success++
		} else if errors.Is(err, domain.ErrRunVersionConflict) {
			conflict++
		} else {
			t.Fatalf("unexpected CAS race error: %v", err)
		}
	}
	after, _ := store.Get(ctx, winner)
	if success != 1 || conflict != 3 || after.Version != run.Version+1 {
		t.Fatalf("CAS race success=%d conflicts=%d run=%#v", success, conflict, after)
	}
}

func TestStep0026R2StoreLimitValidationIsExhaustive(t *testing.T) {
	valid := step0026Limits()
	cases := []RunStoreLimits{
		{MaxArtifactBytes: 0, MaxRunBytes: valid.MaxRunBytes, MaxTotalBytes: valid.MaxTotalBytes},
		{MaxArtifactBytes: -1, MaxRunBytes: valid.MaxRunBytes, MaxTotalBytes: valid.MaxTotalBytes},
		{MaxArtifactBytes: 100, MaxRunBytes: 99, MaxTotalBytes: valid.MaxTotalBytes},
		{MaxArtifactBytes: 100, MaxRunBytes: 200, MaxTotalBytes: 199},
	}
	for _, limits := range cases {
		if _, err := NewFileStore(t.TempDir(), step0026CursorSecret, limits); err == nil {
			t.Fatalf("invalid limits accepted: %#v", limits)
		}
	}
	for _, secret := range [][]byte{nil, []byte("short"), make([]byte, 31), make([]byte, 33)} {
		if _, err := NewFileStore(t.TempDir(), secret, valid); err == nil {
			t.Fatalf("non-256-bit cursor secret accepted: %d bytes", len(secret))
		}
	}
}

func step0026PublishAll(t *testing.T, store *FileStore, run domain.RecommendationRun) domain.RecommendationRun {
	t.Helper()
	for _, name := range []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance} {
		entries := `[]`
		if name == domain.ArtifactRequest || name == domain.ArtifactRecommendation || name == domain.ArtifactDecision || name == domain.ArtifactAlternatives {
			entries = `[{"kind":"` + string(name) + `"}]`
		}
		document := []byte(`{"schema_version":"1","artifact":"` + string(name) + `","entries":` + entries + `}`)
		var err error
		run, err = store.PublishArtifact(context.Background(), domain.ArtifactPublication{RunID: run.RunID, ExpectedVersion: run.Version, Name: name, SchemaVersion: "1", Document: document})
		if err != nil {
			t.Fatalf("publish %s: %v", name, err)
		}
	}
	return run
}

func TestStep0026R2TerminalArtifactMatrixReadinessIntegrityAndReopen(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	id := domain.RunID("run_22222222222222222222222222222222")
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("2", 64))
	run, _, _ = store.CreateOrResolve(ctx, run, record)
	for _, name := range []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance} {
		if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: name, Limit: 1}); !errors.Is(err, domain.ErrArtifactNotReady) {
			t.Fatalf("active absent %s = %v", name, err)
		}
	}
	run, err := store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: domain.RecommendationRunProgress{Planned: 8}})
	if err != nil {
		t.Fatal(err)
	}
	run = step0026PublishAll(t, store, run)
	status := domain.IdeaRunInconclusive
	run, err = store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunSucceeded, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 8}, RecommendationStatus: &status})
	if err != nil {
		t.Fatal(err)
	}
	before, err := store.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened := step0026Store(t, root, step0026Limits())
	defer reopened.Close()
	after, err := reopened.Get(ctx, id)
	if err != nil || !reflect.DeepEqual(before, after) || !reflect.DeepEqual(after, run) {
		t.Fatalf("terminal changed on reopen:\n%#v\n%#v err=%v", before, after, err)
	}
	for _, name := range []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance} {
		if _, err := reopened.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: name, Limit: 50}); err != nil {
			t.Fatalf("required terminal artifact %s unavailable: %v", name, err)
		}
	}
}

func TestStep0026R2CancelledAndInterruptedOnlyExposePublishedRequest(t *testing.T) {
	for _, lifecycle := range []domain.RecommendationLifecycle{domain.RunCancelled, domain.RunInterrupted} {
		t.Run(string(lifecycle), func(t *testing.T) {
			ctx, root := context.Background(), t.TempDir()
			store := step0026Store(t, root, step0026Limits())
			defer store.Close()
			id := domain.RunID("run_" + strings.Repeat(string("cd"[len(string(lifecycle))%2]), 32))
			run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("9", 64))
			run, _, _ = store.CreateOrResolve(ctx, run, record)
			document := []byte(`{"schema_version":"1","artifact":"request","entries":[{"safe":true}]}`)
			run, err := store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: run.Version, Name: domain.ArtifactRequest, SchemaVersion: "1", Document: document})
			if err != nil {
				t.Fatal(err)
			}
			if lifecycle == domain.RunCancelled {
				run, err = store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunCancelled, Stage: domain.RunStageTerminal, Progress: run.Progress})
			} else {
				var changed []domain.RecommendationRun
				changed, err = store.InterruptNonTerminal(ctx)
				if len(changed) == 1 {
					run = changed[0]
				}
			}
			if err != nil || run.Lifecycle != lifecycle {
				t.Fatalf("terminal setup: %#v err=%v", run, err)
			}
			if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactRequest, Limit: 1}); err != nil {
				t.Fatalf("published request hidden: %v", err)
			}
			for _, name := range []domain.RunArtifactName{domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance} {
				if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: name, Limit: 1}); !errors.Is(err, domain.ErrArtifactNotReady) {
					t.Fatalf("%s missing %s = %v, want NOT_READY", lifecycle, name, err)
				}
			}
		})
	}
}

func TestStep0026R2PublishedBlobCorruptionAndMissingTerminalArtifactAreIntegrity(t *testing.T) {
	t.Run("blob-flip-truncate-append", func(t *testing.T) {
		for _, mutation := range []string{"flip", "truncate", "append"} {
			t.Run(mutation, func(t *testing.T) {
				ctx, root := context.Background(), t.TempDir()
				store := step0026Store(t, root, step0026Limits())
				id := domain.RunID("run_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
				run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("a", 64))
				run, _, _ = store.CreateOrResolve(ctx, run, record)
				document := []byte(`{"schema_version":"1","artifact":"failures","entries":[{"safe":true}]}`)
				if _, err := store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: run.Version, Name: domain.ArtifactFailures, SchemaVersion: "1", Document: document}); err != nil {
					t.Fatal(err)
				}
				_ = store.Close()
				blob := step0026FindFileContaining(t, root, `"artifact":"failures"`)
				data, _ := os.ReadFile(blob)
				switch mutation {
				case "flip":
					data[len(data)/2] ^= 1
				case "truncate":
					data = data[:len(data)-1]
				case "append":
					data = append(data, 'x')
				}
				if err := os.WriteFile(blob, data, 0o600); err != nil {
					t.Fatal(err)
				}
				reopened, err := NewFileStore(root, step0026CursorSecret, step0026Limits())
				if err == nil {
					defer reopened.Close()
					_, err = reopened.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactFailures, Limit: 1})
				}
				if !errors.Is(err, domain.ErrIntegrity) {
					t.Fatalf("%s corruption = %v", mutation, err)
				}
			})
		}
	})
	t.Run("missing-required-terminal", func(t *testing.T) {
		ctx, root := context.Background(), t.TempDir()
		store := step0026Store(t, root, step0026Limits())
		id := domain.RunID("run_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
		run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("b", 64))
		run, _, _ = store.CreateOrResolve(ctx, run, record)
		run, _ = store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: domain.RecommendationRunProgress{Planned: 8}})
		run = step0026PublishAll(t, store, run)
		status := domain.IdeaRunComplete
		if _, err := store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunSucceeded, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 8}, RecommendationStatus: &status}); err != nil {
			t.Fatal(err)
		}
		_ = store.Close()
		decisionBlob := step0026FindFileContaining(t, root, `"artifact":"decision"`)
		if err := os.Remove(decisionBlob); err != nil {
			t.Fatal(err)
		}
		reopened, err := NewFileStore(root, step0026CursorSecret, step0026Limits())
		if err == nil {
			defer reopened.Close()
			_, err = reopened.Get(ctx, id)
		}
		if !errors.Is(err, domain.ErrIntegrity) {
			t.Fatalf("missing required terminal artifact = %v", err)
		}
	})
}

func TestStep0026R2CursorTamperMatrixAndSecretSurviveReopen(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	id := domain.RunID("run_33333333333333333333333333333333")
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("3", 64))
	run, _, _ = store.CreateOrResolve(ctx, run, record)
	document := []byte(`{"schema_version":"1","artifact":"provenance","entries":[{"n":0},{"n":1},{"n":2},{"n":3}]}`)
	run, err := store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: run.Version, Name: domain.ArtifactProvenance, SchemaVersion: "1", Document: document})
	if err != nil {
		t.Fatal(err)
	}
	page, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Limit: 1})
	if err != nil || page.NextCursor == nil {
		t.Fatalf("cursor fixture: %#v err=%v", page, err)
	}
	cursor := *page.NextCursor
	parts := strings.Split(cursor, ".")
	if len(parts) != 2 {
		t.Fatalf("cursor format: %q", cursor)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatal(err)
	}
	for name, value := range map[string]string{
		"payload-first-byte": base64.RawURLEncoding.EncodeToString(append([]byte{payload[0] ^ 1}, payload[1:]...)) + "." + parts[1],
		"mac-last-byte":      parts[0] + "." + parts[1][:len(parts[1])-1] + "A",
		"truncated":          cursor[:len(cursor)-1],
		"extra-padding":      parts[0] + "=." + parts[1],
		"empty-mac":          parts[0] + ".",
		"three-segments":     cursor + ".extra",
		"forged-offset":      base64.RawURLEncoding.EncodeToString([]byte(`{"schema_version":"1","run_id":"`+string(id)+`","artifact":"provenance","offset":999}`)) + "." + parts[1],
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Cursor: value, Limit: 1}); !errors.Is(err, domain.ErrCursorInvalid) {
				t.Fatalf("tampered cursor %q = %v", value, err)
			}
		})
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened := step0026Store(t, root, step0026Limits())
	continued, err := reopened.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Cursor: cursor, Limit: 2})
	if err != nil || len(continued.Entries) != 2 || string(continued.Entries[0]) != `{"n":1}` {
		t.Fatalf("cursor/secret did not survive reopen: %#v err=%v", continued, err)
	}
	_ = reopened.Close()
	wrongSecret := append([]byte(nil), step0026CursorSecret...)
	wrongSecret[0] ^= 1
	wrong, err := NewFileStore(root, wrongSecret, step0026Limits())
	if err == nil {
		defer wrong.Close()
		_, err = wrong.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Cursor: cursor, Limit: 1})
	}
	if err == nil {
		t.Fatal("different reopen secret accepted an existing cursor")
	}
	_ = run
}

func TestStep0026R2ExternalIdentifiersNeverBecomePathComponents(t *testing.T) {
	root := t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	defer store.Close()
	id := domain.RunID("run_deadbeefdeadbeefdeadbeefdeadbeef")
	rawKey := "RAW-IDEMPOTENCY-KEY-CANARY"
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("4", 64))
	if _, _, err := store.CreateOrResolve(context.Background(), run, record); err != nil {
		t.Fatal(err)
	}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		if strings.Contains(path, string(id)) || strings.Contains(path, rawKey) || strings.Contains(path, "deadbeefdeadbeefdeadbeefdeadbeef") {
			t.Errorf("external identifier became path: %s", path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	bytes := step0026ReadRoot(t, root)
	if strings.Contains(bytes, rawKey) {
		t.Fatal("raw idempotency key persisted")
	}
}

func TestStep0026R2StoreRejectsSymlinkHardlinkFIFOAndIgnoresUnconfirmedTopLevelFiles(t *testing.T) {
	ctx := context.Background()
	for _, kind := range []string{"symlink", "hardlink", "fifo"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			store := step0026Store(t, root, step0026Limits())
			id := domain.RunID("run_55555555555555555555555555555555")
			run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("5", 64))
			_, _, _ = store.CreateOrResolve(ctx, run, record)
			_ = store.Close()
			regular := step0026FirstRegular(t, root)
			external := filepath.Join(t.TempDir(), "external")
			if err := os.WriteFile(external, []byte("SENTINEL"), 0o600); err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "symlink":
				if err := os.Remove(regular); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(external, regular); err != nil {
					t.Skip(err)
				}
			case "hardlink":
				if err := os.Link(regular, external+"-link"); err != nil {
					t.Skip(err)
				}
			case "fifo":
				if err := os.Remove(regular); err != nil {
					t.Fatal(err)
				}
				if err := syscall.Mkfifo(regular, 0o600); err != nil {
					t.Skip(err)
				}
			}
			reopened, err := NewFileStore(root, step0026CursorSecret, step0026Limits())
			if err == nil {
				defer reopened.Close()
				_, err = reopened.Get(ctx, id)
			}
			if err == nil {
				t.Fatalf("%s component accepted", kind)
			}
			data, _ := os.ReadFile(external)
			if string(data) != "SENTINEL" {
				t.Fatalf("external sentinel changed through %s: %q", kind, data)
			}
		})
	}
	root := t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	id := domain.RunID("run_66666666666666666666666666666666")
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("6", 64))
	_, _, _ = store.CreateOrResolve(ctx, run, record)
	_ = store.Close()
	for _, suffix := range []string{"pending", "partial", "fsynced", "renamed", "orphan"} {
		if err := os.WriteFile(filepath.Join(root, ".crash-"+suffix+".tmp"), []byte(`{"schema_version":"1","artifact":"alternatives","entries":[{"CANARY":"ORPHAN"}]}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	reopened := step0026Store(t, root, step0026Limits())
	defer reopened.Close()
	if _, err := reopened.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactAlternatives, Limit: 1}); !errors.Is(err, domain.ErrArtifactNotReady) {
		t.Fatalf("orphan crash bytes became visible: %v", err)
	}
}

func TestStep0026R2ConcurrentQuotaNeverOvershootsOrEvicts(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	limits := RunStoreLimits{MaxArtifactBytes: 3500, MaxRunBytes: 4096, MaxTotalBytes: 16384}
	store := step0026Store(t, root, limits)
	defer store.Close()
	ids := []domain.RunID{"run_77777777777777777777777777777777", "run_88888888888888888888888888888888"}
	runs := make([]domain.RecommendationRun, len(ids))
	for index, id := range ids {
		run, record := step0026StoredRun(id, "sha256:"+strings.Repeat(fmt.Sprintf("%x", index+7), 64))
		var err error
		runs[index], _, err = store.CreateOrResolve(ctx, run, record)
		if err != nil {
			t.Fatal(err)
		}
	}
	document := []byte(`{"schema_version":"1","artifact":"failures","entries":["` + strings.Repeat("x", 3000) + `"]}`)
	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for index := range runs {
		index := index
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: runs[index].RunID, ExpectedVersion: runs[index].Version, Name: domain.ArtifactFailures, SchemaVersion: "1", Document: document})
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	quotaErrors := 0
	for err := range errs {
		if errors.Is(err, domain.ErrQuotaExceeded) {
			quotaErrors++
		} else if err != nil {
			t.Fatalf("unexpected quota race error: %v", err)
		}
	}
	if quotaErrors == 0 {
		t.Fatal("concurrent quota race admitted every large publication")
	}
	for _, id := range ids {
		if _, err := store.Get(ctx, id); err != nil {
			t.Fatalf("quota evicted prior run %s: %v", id, err)
		}
	}
}

func step0026FirstRegular(t *testing.T, root string) string {
	t.Helper()
	var found string
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err == nil && found == "" && !entry.IsDir() && entry.Type().IsRegular() {
			found = path
		}
		return nil
	})
	if found == "" {
		t.Fatal("no regular store file")
	}
	return found
}

func step0026FindFileContaining(t *testing.T, root, fragment string) string {
	t.Helper()
	var found string
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !entry.Type().IsRegular() || found != "" {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr == nil && strings.Contains(string(data), fragment) {
			found = path
		}
		return nil
	})
	if found == "" {
		t.Fatalf("no file contains %q", fragment)
	}
	return found
}

func step0026ReadRoot(t *testing.T, root string) string {
	t.Helper()
	var out strings.Builder
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !entry.Type().IsRegular() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out.Write(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return out.String()
}
