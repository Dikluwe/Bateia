/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/insights.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestExtractArchitecturalInsights_Detection(t *testing.T) {
	repos := []domain.ExtractedRepository{
		{
			Candidate: domain.RepositoryCandidate{Owner: "werbot", Repo: "werbot"},
			Manifest: &domain.Manifest{
				Dependencies: []domain.Dependency{
					{Name: "github.com/stretchr/testify", Version: "v1.11.1"},
					{Name: "github.com/pressly/goose/v3", Version: "v3.26.0"},
					{Name: "github.com/rs/zerolog", Version: "v1.34.0"},
				},
			},
			Files: []domain.ExtractedFile{
				{Path: "Dockerfile"},
			},
		},
	}

	insights := ExtractArchitecturalInsights(repos)

	if len(insights.Testing) == 0 || insights.Testing[0].Name != "testify" {
		t.Errorf("esperava detectar a ferramenta 'testify' em Testing")
	}

	if len(insights.Database) == 0 || insights.Database[0].Name != "goose" {
		t.Errorf("esperava detectar a ferramenta 'goose' em Database")
	}

	if len(insights.Observability) == 0 || insights.Observability[0].Name != "zerolog" {
		t.Errorf("esperava detectar a ferramenta 'zerolog' em Observability")
	}

	if len(insights.DevOps) == 0 || insights.DevOps[0].Name != "Dockerfile" {
		t.Errorf("esperava detectar 'Dockerfile' em DevOps")
	}
}

func TestExtractArchitecturalInsights_Empty(t *testing.T) {
	insights := ExtractArchitecturalInsights(nil)

	if len(insights.Testing) != 0 || len(insights.Database) != 0 {
		t.Errorf("esperava listas vazias quando o slice de repositórios for nil ou vazio")
	}
}
