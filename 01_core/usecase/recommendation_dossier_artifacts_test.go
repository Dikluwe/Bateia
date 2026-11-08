/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_dossier_artifacts_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"errors"
	"testing"
)

func TestDossierArtifactFailureCannotAdvertiseCompleteTrace(t *testing.T) {
	files := &memoryDossierFiles{docs: map[string][]byte{}}
	if _, err := BuildRecommendationTraceStage(files, domain.DecisionTraceGraph{}, domain.DecisionTraceExplanation{}, nil); err == nil {
		t.Fatal("missing artifacts accepted")
	}
	manifest := domain.IdeaRunManifest{Status: domain.IdeaRunComplete}
	MarkRecommendationTraceFailed(&manifest, errors.New("unavailable"))
	if manifest.Status != domain.IdeaRunFailed || len(manifest.Stages) != 1 || manifest.Stages[0].Status != domain.IdeaStageFailed || len(manifest.Stages[0].Outputs) != 0 {
		t.Fatal("failed sidecar advertised")
	}
}
