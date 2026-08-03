/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/mermaid.md
 * @layer L2
 * @updated 2026-08-03
 */

package report

import (
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestGenerateMermaidGraph(t *testing.T) {
	rep := &domain.SearchReport{
		Query: domain.QueryFilter{
			Query: "architecture",
		},
		Results: []domain.ExtractedRepository{
			{
				Candidate: domain.RepositoryCandidate{
					Owner: "test",
					Repo:  "repo",
				},
				Manifest: &domain.Manifest{
					Dependencies: []domain.Dependency{
						{Name: "express", Version: "1.0.0"},
					},
				},
			},
		},
	}

	graph := GenerateMermaidGraph(rep)
	if !strings.Contains(graph, "graph TD") {
		t.Errorf("Esperava graph TD no diagrama Mermaid")
	}
	if !strings.Contains(graph, "test/repo") {
		t.Errorf("Esperava test/repo no diagrama")
	}
	if !strings.Contains(graph, "express") {
		t.Errorf("Esperava dependência express no diagrama")
	}
}
