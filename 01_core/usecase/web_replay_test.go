/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/web_replay_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func TestReplayAndCompareWebSnapshot(t *testing.T) {
	snap := domain.WebSearchSnapshot{
		SchemaVersion: "1",
		RunID:         "run1",
		Arm:           "F2",
		Responses: []domain.RawWebResponse{
			{
				Results: []domain.RawWebResult{
					{Rank: 1, URL: "https://example.com/repo1/", Title: "Repo 1"},
					{Rank: 2, URL: "https://example.com/repo1", Title: "Repo 1 Dup"},
				},
			},
		},
	}

	res, err := ReplayWebSnapshot(snap)
	if err != nil {
		t.Fatalf("replay error: %v", err)
	}
	if len(res.Deduplicated) != 1 {
		t.Errorf("expected 1 deduplicated hit, got %d", len(res.Deduplicated))
	}
}
