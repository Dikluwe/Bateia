/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/recommendation_run_step0026_contract_test.md
 * @layer L1
 * @updated 2026-09-01
 */
package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStep0026RunDomainConstantsAndPublicJSONAreFrozen(t *testing.T) {
	if got, want := []RecommendationLifecycle{RunQueued, RunRunning, RunSucceeded, RunFailed, RunCancelling, RunCancelled, RunInterrupted}, []RecommendationLifecycle{"QUEUED", "RUNNING", "SUCCEEDED", "FAILED", "CANCELLING", "CANCELLED", "INTERRUPTED"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("lifecycles = %#v, want %#v", got, want)
	}
	if got, want := []RecommendationRunStage{RunStageQueued, RunStageExpand, RunStageDiscover, RunStageScreen, RunStageCollect, RunStageAssess, RunStageDecide, RunStageAlternatives, RunStagePublish, RunStageTerminal}, []RecommendationRunStage{"queued", "expand", "discover", "screen", "collect", "assess", "decide", "alternatives", "publish", "terminal"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("stages = %#v, want %#v", got, want)
	}
	if got, want := []RunArtifactName{ArtifactRequest, ArtifactRecommendation, ArtifactDecision, ArtifactAlternatives, ArtifactFailures, ArtifactProvenance}, []RunArtifactName{"request", "recommendation", "decision", "alternatives", "failures", "provenance"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("artifacts = %#v, want %#v", got, want)
	}
	if ArtifactPending != "PENDING" || ArtifactPublished != "PUBLISHED" {
		t.Fatalf("artifact states changed: %q %q", ArtifactPending, ArtifactPublished)
	}

	status := IdeaRunInconclusive
	run := RecommendationRunView{
		SchemaVersion: "1", RunID: "run_0123456789abcdef0123456789abcdef",
		Lifecycle: RunSucceeded, Version: 9, Stage: RunStageTerminal,
		Progress:             RecommendationRunProgress{Planned: 8, Completed: 8},
		RecommendationStatus: &status,
		Artifacts:            []RunArtifactIndex{}, Failure: nil,
	}
	got, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","lifecycle":"SUCCEEDED","version":9,"stage":"terminal","progress":{"planned":8,"completed":8},"recommendation_status":"INCONCLUSIVE","artifacts":[],"failure":null}`
	if string(got) != want {
		t.Fatalf("run JSON changed:\ngot  %s\nwant %s", got, want)
	}

	receipt, err := json.Marshal(RecommendationStartReceipt{SchemaVersion: "1", RunID: "run_0123456789abcdef0123456789abcdef", Lifecycle: RunQueued, Version: 1})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","lifecycle":"QUEUED","version":1}`; string(receipt) != want {
		t.Fatalf("receipt JSON changed: %s", receipt)
	}

	page := ArtifactPage{SchemaVersion: "1", RunID: "run_0123456789abcdef0123456789abcdef", Artifact: ArtifactAlternatives, ArtifactHash: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Entries: [][]byte{[]byte(`{"rank":1}`)}, NextCursor: nil}
	if page.Entries == nil {
		t.Fatal("fixture error: entries must be materially empty/non-null")
	}
}

func TestStep0026InitialRunSeparatesLifecycleFromRecommendationStatus(t *testing.T) {
	run := RecommendationRun{
		SchemaVersion: "1", RunID: "run_0123456789abcdef0123456789abcdef",
		Lifecycle: RunQueued, Version: 1, Stage: RunStageQueued,
		Progress:  RecommendationRunProgress{Planned: 8, Completed: 0},
		Artifacts: []RunArtifactIndex{},
	}
	encoded, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]interface{}
	if err := json.Unmarshal(encoded, &object); err != nil {
		t.Fatal(err)
	}
	if object["recommendation_status"] != nil || object["failure"] != nil {
		t.Fatalf("initial run invented terminal meaning: %s", encoded)
	}
	if artifacts, ok := object["artifacts"].([]interface{}); !ok || len(artifacts) != 0 {
		t.Fatalf("artifacts must encode as [], got %s", encoded)
	}
	if _, leaked := object["request_hash"]; leaked {
		t.Fatalf("public view leaks request hash: %s", encoded)
	}
}
