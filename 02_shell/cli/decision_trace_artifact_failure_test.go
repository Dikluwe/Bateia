/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/decision_trace_artifact_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func TestDecisionTraceStageRequiresIndexAndFailureAdvertisesNoOutputs(t *testing.T) {
	dir := t.TempDir()
	for _, path := range []string{"idea.json", "reuse-decision.json", "trace-graph.json"} {
		if err := os.WriteFile(filepath.Join(dir, path), []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	stage, err := buildDecisionTraceStage(dir, domain.DecisionTraceGraph{}, domain.DecisionTraceExplanation{}, nil)
	if err == nil || len(stage.Outputs) != 0 {
		t.Fatal("missing index was announced as a completed trace")
	}
	manifest := domain.IdeaRunManifest{Status: domain.IdeaRunRunning, Stages: []domain.IdeaRunStage{}}
	markDecisionTraceFailed(&manifest, errors.New("index unavailable"))
	if manifest.Status != domain.IdeaRunFailed || len(manifest.Failures) != 1 || len(manifest.Stages) != 1 {
		t.Fatalf("failure not recorded: %+v", manifest)
	}
	failed := manifest.Stages[0]
	if failed.Status != domain.IdeaStageFailed || failed.ContractVersion != decisionTraceArtifactContractVersion || len(failed.Outputs) != 0 {
		t.Fatalf("partial index advertised: %+v", failed)
	}
}
