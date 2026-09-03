/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runstore/file_store_step0026_r3_test.md
 * @layer L3
 * @updated 2026-09-01
 */
package runstore

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestStep0026R3EveryConfirmedMutationIsImmediatelyReadableAndSurvivesReopen(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	id := domain.RunID("run_90909090909090909090909090909090")
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("9", 64))
	run, won, err := store.CreateOrResolve(ctx, run, record)
	if err != nil || !won {
		t.Fatalf("create: won=%v err=%v", won, err)
	}
	assertRead := func(want domain.RecommendationRun) {
		t.Helper()
		got, err := store.Get(ctx, id)
		if err != nil || got.Version != want.Version || got.Lifecycle != want.Lifecycle || got.Stage != want.Stage || got.Progress != want.Progress {
			t.Fatalf("read-after-confirmation: got=%#v want=%#v err=%v", got, want, err)
		}
	}
	assertRead(run)
	run, err = store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: run.Version, Name: domain.ArtifactRequest, SchemaVersion: "1", Document: []byte(`{"schema_version":"1","artifact":"request","entries":[{"safe":true}]}`)})
	if err != nil {
		t.Fatal(err)
	}
	assertRead(run)
	for index, stage := range []domain.RecommendationRunStage{domain.RunStageExpand, domain.RunStageDiscover, domain.RunStageScreen, domain.RunStageCollect, domain.RunStageAssess, domain.RunStageDecide, domain.RunStageAlternatives, domain.RunStagePublish} {
		completed := index
		run, err = store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: stage, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: completed}})
		if err != nil {
			t.Fatalf("stage %s: %v", stage, err)
		}
		assertRead(run)
	}
	for _, name := range []domain.RunArtifactName{domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance} {
		run, err = store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: run.Version, Name: name, SchemaVersion: "1", Document: []byte(`{"schema_version":"1","artifact":"` + string(name) + `","entries":[]}`)})
		if err != nil {
			t.Fatalf("publish %s: %v", name, err)
		}
		assertRead(run)
		if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: name, Limit: 1}); err != nil {
			t.Fatalf("published %s not immediately readable: %v", name, err)
		}
	}
	status := domain.IdeaRunComplete
	run, err = store.Transition(ctx, domain.RunTransition{RunID: id, ExpectedVersion: run.Version, Lifecycle: domain.RunSucceeded, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 8}, RecommendationStatus: &status})
	if err != nil {
		t.Fatal(err)
	}
	assertRead(run)
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened := step0026Store(t, root, step0026Limits())
	defer reopened.Close()
	got, err := reopened.Get(ctx, id)
	if err != nil || got.Version != run.Version || got.Lifecycle != run.Lifecycle {
		t.Fatalf("reopen lost final confirmation: %#v err=%v", got, err)
	}
}

func TestStep0026R3ValidMACCursorRejectsResignedVersionHashAndOffsetChanges(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	store := step0026Store(t, root, step0026Limits())
	defer store.Close()
	id := domain.RunID("run_40404040404040404040404040404040")
	run, record := step0026StoredRun(id, "sha256:"+strings.Repeat("4", 64))
	run, _, _ = store.CreateOrResolve(ctx, run, record)
	run, err := store.PublishArtifact(ctx, domain.ArtifactPublication{RunID: id, ExpectedVersion: run.Version, Name: domain.ArtifactProvenance, SchemaVersion: "1", Document: []byte(`{"schema_version":"1","artifact":"provenance","entries":[{"n":0},{"n":1}]}`)})
	if err != nil {
		t.Fatal(err)
	}
	page, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Limit: 1})
	if err != nil || page.NextCursor == nil {
		t.Fatalf("cursor fixture: %#v err=%v", page, err)
	}
	parts := strings.Split(*page.NextCursor, ".")
	if len(parts) != 2 {
		t.Fatalf("cursor format: %q", *page.NextCursor)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatal(err)
	}
	sign := func(raw []byte, overEncoded bool) string {
		encoded := base64.RawURLEncoding.EncodeToString(raw)
		input := raw
		if overEncoded {
			input = []byte(encoded)
		}
		mac := hmac.New(sha256.New, step0026CursorSecret)
		_, _ = mac.Write(input)
		return encoded + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	}
	overEncoded := false
	if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Cursor: sign(payload, false), Limit: 1}); err != nil {
		overEncoded = true
		if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Cursor: sign(payload, true), Limit: 1}); err != nil {
			t.Fatalf("neither public HMAC interpretation reproduced the valid cursor: %v", err)
		}
	}
	var fields map[string]interface{}
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatal(err)
	}
	mutations := map[string]func(map[string]interface{}){
		"version": func(value map[string]interface{}) { value["artifact_version"] = float64(1 << 20) },
		"hash":    func(value map[string]interface{}) { value["hash"] = "sha256:" + strings.Repeat("0", 64) },
		"offset":  func(value map[string]interface{}) { value["offset"] = float64(1 << 20) },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			copy := map[string]interface{}{}
			for key, value := range fields {
				copy[key] = value
			}
			mutate(copy)
			changed, err := json.Marshal(copy)
			if err != nil {
				t.Fatal(err)
			}
			cursor := sign(changed, overEncoded)
			if _, err := store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: id, Artifact: domain.ArtifactProvenance, Cursor: cursor, Limit: 1}); !errors.Is(err, domain.ErrCursorInvalid) {
				t.Fatalf("valid-MAC %s mutation = %v", name, err)
			}
		})
	}
	_ = run
}

func TestStep0026R3FullStoreRemainsFullAfterReopenAndNeverEvictsOnCreate(t *testing.T) {
	ctx, root := context.Background(), t.TempDir()
	limits := RunStoreLimits{MaxArtifactBytes: 512, MaxRunBytes: 2048, MaxTotalBytes: 4096}
	store := step0026Store(t, root, limits)
	var admitted []domain.RunID
	for index := 1; index <= 128; index++ {
		id := domain.RunID(fmt.Sprintf("run_%032x", index))
		run, record := step0026StoredRun(id, "sha256:"+fmt.Sprintf("%064x", index))
		_, won, err := store.CreateOrResolve(ctx, run, record)
		if errors.Is(err, domain.ErrQuotaExceeded) {
			break
		}
		if err != nil || !won {
			t.Fatalf("fill create %d: won=%v err=%v", index, won, err)
		}
		admitted = append(admitted, id)
	}
	if len(admitted) == 0 || len(admitted) == 128 {
		t.Fatalf("fixture did not reach quota after real creates: admitted=%d", len(admitted))
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened := step0026Store(t, root, limits)
	defer reopened.Close()
	newID := domain.RunID("run_ffffffffffffffffffffffffffffffff")
	newRun, newRecord := step0026StoredRun(newID, "sha256:"+strings.Repeat("f", 64))
	if _, _, err := reopened.CreateOrResolve(ctx, newRun, newRecord); !errors.Is(err, domain.ErrQuotaExceeded) {
		t.Fatalf("full reopened store admitted create: %v", err)
	}
	for _, id := range admitted {
		if _, err := reopened.Get(ctx, id); err != nil {
			t.Fatalf("create-time quota evicted %s: %v", id, err)
		}
	}
}
