/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/recommendation_run_test.md
 * @layer L1
 * @updated 2026-09-01
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRecommendationRunJSONProjectsOnlyPublicView(t *testing.T) {
	run := RecommendationRun{
		SchemaVersion: "1", RunID: "run_0123456789abcdef0123456789abcdef",
		RequestHash: "sha256:internal", Lifecycle: RunQueued, Version: 1,
		Stage: RunStageQueued, Progress: RecommendationRunProgress{Planned: 8},
		Artifacts: []RunArtifactIndex{},
	}
	encoded, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "request_hash") || strings.Contains(string(encoded), `"request"`) {
		t.Fatalf("internal request material leaked: %s", encoded)
	}
	if !strings.Contains(string(encoded), `"artifacts":[]`) {
		t.Fatalf("public artifacts must be a non-null array: %s", encoded)
	}
}
