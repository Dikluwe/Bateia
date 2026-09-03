/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runstore/file_store_step0026_contract_test.md
 * @layer L3
 * @updated 2026-09-01
 */
package runstore

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

var step0026CursorSecret = []byte("0123456789abcdef0123456789abcdef")

func step0026StoredRun(id domain.RunID, digest string) (domain.RecommendationRun, domain.IdempotencyRecord) {
	request := domain.RecommendationRequest{
		Idea:   domain.IdeaSpec{Problem: "offline", Capabilities: []string{"inventory"}},
		Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .3},
		Target: domain.SoftwareAlternativeTarget{Name: "offline"},
		Search: domain.RecommendationSearch{Strategy: "literal", Locales: []string{"en-US"}, MaxPerQuery: 1, MaxCandidates: 1},
	}
	run := domain.RecommendationRun{SchemaVersion: "1", RunID: id, Request: request, RequestHash: "sha256:" + strings.Repeat("a", 64), Lifecycle: domain.RunQueued, Version: 1, Stage: domain.RunStageQueued, Progress: domain.RecommendationRunProgress{Planned: 8}, Artifacts: []domain.RunArtifactIndex{}}
	record := domain.IdempotencyRecord{KeyDigest: digest, RequestHash: run.RequestHash, RunID: id}
	return run, record
}

func step0026Store(t *testing.T, root string, limits RunStoreLimits) *FileStore {
	t.Helper()
	store, err := NewFileStore(root, step0026CursorSecret, limits)
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func step0026Limits() RunStoreLimits {
	return RunStoreLimits{MaxArtifactBytes: 1 << 20, MaxRunBytes: 4 << 20, MaxTotalBytes: 32 << 20}
}

func TestStep0026FileStoreCreateCASPermissionsAndReopen(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	id := domain.RunID("run_0123456789abcdef0123456789abcdef")
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("b", 64))
	created, won, err := store.CreateOrResolve(ctx, run, record)
	if err != nil || !won || !reflect.DeepEqual(created, run) {
		t.Fatalf("create = %#v won=%v err=%v", created, won, err)
	}
	if again, won, err := store.CreateOrResolve(ctx, run, record); err != nil || won || again.RunID != id {
		t.Fatalf("atomic resolve = %#v won=%v err=%v", again, won, err)
	}
	if _, err := store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: 99, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: domain.RecommendationRunProgress{Planned: 8}}); !errors.Is(err, domain.ErrRunVersionConflict) {
		t.Fatalf("stale CAS = %v", err)
	}
	running, err := store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: 1, Lifecycle: domain.RunRunning, Stage: domain.RunStageExpand, Progress: domain.RecommendationRunProgress{Planned: 8}})
	if err != nil || running.Version != 2 || running.Lifecycle != domain.RunRunning {
		t.Fatalf("valid CAS = %#v err=%v", running, err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened := step0026Store(t, root, step0026Limits())
	t.Cleanup(func() { _ = reopened.Close() })
	got, err := reopened.Get(ctx, id)
	if err != nil || !reflect.DeepEqual(got, running) {
		t.Fatalf("reopen lost run: %#v err=%v", got, err)
	}
	resolved, ok, err := reopened.ResolveIdempotency(ctx, record.KeyDigest)
	if err != nil || !ok || !reflect.DeepEqual(resolved, record) {
		t.Fatalf("reopen lost idempotency: %#v ok=%v err=%v", resolved, ok, err)
	}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() && info.Mode().Perm()&0o077 != 0 {
			t.Errorf("public directory permissions at %s: %o", path, info.Mode().Perm())
		}
		if info.Mode().IsRegular() && info.Mode().Perm()&0o077 != 0 {
			t.Errorf("public file permissions at %s: %o", path, info.Mode().Perm())
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestStep0026FileStorePublishesFullDocumentHashAndPaginatesBoundCursor(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	t.Cleanup(func() { _ = store.Close() })
	id := domain.RunID("run_0123456789abcdef0123456789abcdef")
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("b", 64))
	run, _, _ = store.CreateOrResolve(ctx, run, record)
	document := []byte(`{"schema_version":"1","artifact":"failures","entries":[{"code":"A"},{"code":"B"},{"code":"C"}]}`)
	published, err := store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: run.Version, Name: domain.ArtifactFailures, SchemaVersion: "1", Document: document})
	if err != nil {
		t.Fatal(err)
	}
	if published.Version != run.Version+1 || len(published.Artifacts) != 1 || published.Artifacts[0].State != domain.ArtifactPublished {
		t.Fatalf("publication was not one confirmed mutation: %#v", published)
	}
	digest := sha256.Sum256(document)
	wantHash := "sha256:" + hex.EncodeToString(digest[:])
	index := published.Artifacts[0]
	if index.Hash != wantHash || index.SizeBytes != int64(len(document)) || index.EntryCount != 3 {
		t.Fatalf("index does not hash full canonical document: %#v want hash=%s size=%d", index, wantHash, len(document))
	}
	if _, err := store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: published.Version, Name: domain.ArtifactFailures, SchemaVersion: "1", Document: []byte(`{"schema_version":"1","artifact":"failures","entries":[]}`)}); !errors.Is(err, domain.ErrArtifactImmutable) {
		t.Fatalf("published artifact was overwriteable: %v", err)
	}
	unchanged, _ := store.Get(ctx, id)
	if unchanged.Version != published.Version || unchanged.Artifacts[0].Hash != wantHash {
		t.Fatalf("overwrite attempt changed terminal artifact/index: %#v", unchanged)
	}
	first, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactFailures, Limit: 1})
	if err != nil || len(first.Entries) != 1 || first.NextCursor == nil || string(first.Entries[0]) != `{"code":"A"}` {
		t.Fatalf("first page = %#v err=%v", first, err)
	}
	second, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactFailures, Cursor: *first.NextCursor, Limit: 2})
	if err != nil || len(second.Entries) != 2 || second.NextCursor != nil || string(second.Entries[1]) != `{"code":"C"}` {
		t.Fatalf("cursor bound to prior limit or pagination changed: %#v err=%v", second, err)
	}
	repeated, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactFailures, Limit: 1})
	if err != nil || !reflect.DeepEqual(first, repeated) {
		t.Fatalf("pagination is not repeatable: %#v vs %#v err=%v", first, repeated, err)
	}
	tampered := []byte(*first.NextCursor)
	if tampered[len(tampered)-1] == 'A' {
		tampered[len(tampered)-1] = 'B'
	} else {
		tampered[len(tampered)-1] = 'A'
	}
	if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactFailures, Cursor: string(tampered), Limit: 1}); !errors.Is(err, domain.ErrCursorInvalid) {
		t.Fatalf("tampered HMAC cursor = %v", err)
	}
	otherID := domain.RunID("run_11111111111111111111111111111111")
	other, otherRecord := step0026StoredRun(otherID, "sha256:"+strings.Repeat("c", 64))
	if _, _, err := store.CreateOrResolve(ctx, other, otherRecord); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: otherID, Artifact: domain.ArtifactFailures, Cursor: *first.NextCursor, Limit: 1}); !errors.Is(err, domain.ErrCursorInvalid) {
		t.Fatalf("cross-run cursor = %v", err)
	}
	if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Cursor: *first.NextCursor, Limit: 1}); !errors.Is(err, domain.ErrCursorInvalid) {
		t.Fatalf("cross-artifact cursor = %v", err)
	}
}

func TestStep0026FileStoreRejectsPendingCorruptionSymlinkHardlinkAndQuota(t *testing.T) {
	ctx := context.Background()
	t.Run("invalid-identifiers-before-disk", func(t *testing.T) {
		root := t.TempDir()
		store := step0026Store(t, root, step0026Limits())
		defer store.Close()
		before := step0026TreeNames(t, root)
		for _, id := range []domain.RunID{"../escape", "run_abc", domain.RunID("run_" + strings.Repeat("a", 1024))} {
			if _, err := store.Get(ctx, id); !errors.Is(err, domain.ErrInvalidRunInput) {
				t.Fatalf("invalid id %q = %v", id, err)
			}
		}
		if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: "run_0123456789abcdef0123456789abcdef", Artifact: "../../raw", Limit: 1}); !errors.Is(err, domain.ErrInvalidRunInput) {
			t.Fatalf("invalid artifact name = %v", err)
		}
		if after := step0026TreeNames(t, root); !reflect.DeepEqual(after, before) {
			t.Fatalf("invalid identifiers touched disk:\nbefore=%#v\nafter=%#v", before, after)
		}
	})
	t.Run("artifact-not-ready-is-not-empty", func(t *testing.T) {
		root := t.TempDir()
		store := step0026Store(t, root, step0026Limits())
		defer store.Close()
		id := domain.RunID("run_0123456789abcdef0123456789abcdef")
		run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("b", 64))
		_, _, _ = store.CreateOrResolve(ctx, run, record)
		if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactAlternatives, Limit: 1}); !errors.Is(err, domain.ErrArtifactNotReady) {
			t.Fatalf("pending/unpublished became empty success: %v", err)
		}
	})
	t.Run("artifact-quota-before-write", func(t *testing.T) {
		root := t.TempDir()
		limits := RunStoreLimits{MaxArtifactBytes: 64, MaxRunBytes: 1 << 20, MaxTotalBytes: 2 << 20}
		store := step0026Store(t, root, limits)
		defer store.Close()
		id := domain.RunID("run_0123456789abcdef0123456789abcdef")
		run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("b", 64))
		run, _, _ = store.CreateOrResolve(ctx, run, record)
		document := []byte(`{"schema_version":"1","artifact":"failures","entries":["` + strings.Repeat("x", 256) + `"]}`)
		if _, err := store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: run.Version, Name: domain.ArtifactFailures, SchemaVersion: "1", Document: document}); !errors.Is(err, domain.ErrQuotaExceeded) {
			t.Fatalf("oversized publication = %v", err)
		}
		after, _ := store.Get(ctx, id)
		if after.Version != run.Version || len(after.Artifacts) != 0 {
			t.Fatalf("quota overshoot left pending/index state: %#v", after)
		}
	})
	t.Run("symlink-root", func(t *testing.T) {
		target := t.TempDir()
		link := filepath.Join(t.TempDir(), "store-link")
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}
		if _, err := NewFileStore(link, step0026CursorSecret, step0026Limits()); err == nil {
			t.Fatal("symlink root accepted")
		}
	})
	t.Run("hardlink-on-reopen", func(t *testing.T) {
		root := t.TempDir()
		store := step0026Store(t, root, step0026Limits())
		id := domain.RunID("run_0123456789abcdef0123456789abcdef")
		run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("b", 64))
		_, _, _ = store.CreateOrResolve(ctx, run, record)
		_ = store.Close()
		var regular string
		_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err == nil && !entry.IsDir() && entry.Type().IsRegular() && regular == "" {
				regular = path
			}
			return nil
		})
		if regular == "" {
			t.Fatal("store created no regular durable file")
		}
		if err := os.Link(regular, filepath.Join(root, "unexpected-hardlink")); err != nil {
			t.Skipf("hardlinks unavailable: %v", err)
		}
		if _, err := NewFileStore(root, step0026CursorSecret, step0026Limits()); err == nil {
			t.Fatal("hardlinked store file accepted on reopen")
		}
	})
}

func step0026TreeNames(t *testing.T, root string) []string {
	t.Helper()
	var names []string
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		names = append(names, relative)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return names
}

func TestStep0026FileStoreCorruptionIsIntegrityAndOrphansInterruptHonestly(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	id := domain.RunID("run_0123456789abcdef0123456789abcdef")
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("b", 64))
	if _, _, err := store.CreateOrResolve(ctx, run, record); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened := step0026Store(t, root, step0026Limits())
	changed, err := reopened.InterruptNonTerminal(ctx)
	if err != nil || len(changed) != 1 || changed[0].Lifecycle != domain.RunInterrupted || changed[0].RecommendationStatus != nil || changed[0].Version != 2 {
		t.Fatalf("honest reopen recovery = %#v err=%v", changed, err)
	}
	_ = reopened.Close()

	var recordPath string
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !entry.Type().IsRegular() {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr == nil && strings.Contains(string(data), string(id)) {
			recordPath = path
		}
		return nil
	})
	if recordPath == "" {
		t.Fatal("could not locate opaque record for corruption witness")
	}
	if err := os.WriteFile(recordPath, []byte(`{"corrupt":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	corrupt, err := NewFileStore(root, step0026CursorSecret, step0026Limits())
	if err == nil {
		defer corrupt.Close()
		_, err = corrupt.Get(ctx, id)
	}
	if !errors.Is(err, domain.ErrIntegrity) {
		t.Fatalf("corruption became absence/success: %v", err)
	}
}
