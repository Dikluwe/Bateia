/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/html_generator_test.md
 * @layer L2
 * @updated 2026-08-03
 */

package report

import (
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func TestGenerateHTMLReport(t *testing.T) {
	rep := &domain.SearchReport{
		Query: domain.QueryFilter{
			Query:    "test query",
			Language: "go",
		},
		GeneratedAt:    time.Now().UTC(),
		TopCandidates:  5,
		TriagedCount:   3,
		ExtractedCount: 2,
		Results: []domain.ExtractedRepository{
			{
				Candidate: domain.RepositoryCandidate{
					Owner: "testowner",
					Repo:  "testrepo",
					URL:   "https://github.com/testowner/testrepo",
				},
				RelevanceScore:  90,
				RelevanceReason: "Direct match",
			},
		},
	}

	htmlOut := GenerateHTMLReport(rep)

	if !strings.Contains(htmlOut, "<!DOCTYPE html>") {
		t.Errorf("Esperava DOCTYPE html na saída")
	}
	if !strings.Contains(htmlOut, "testowner/testrepo") {
		t.Errorf("Esperava nome do repositório no HTML")
	}
	if !strings.Contains(htmlOut, "Score 90/100") {
		t.Errorf("Esperava badge de score 90/100")
	}
}
