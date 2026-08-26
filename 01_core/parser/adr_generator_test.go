/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/adr_generator_test.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestGenerateADRDocument(t *testing.T) {
	rep := &domain.SearchReport{
		Query: domain.QueryFilter{
			Query:    "time tracking",
			Language: "go",
		},
		Results: []domain.ExtractedRepository{
			{
				Candidate: domain.RepositoryCandidate{
					Owner: "solidtime",
					Repo:  "solidtime",
				},
				RelevanceScore:  85,
				RelevanceReason: "Excellent match",
			},
		},
	}

	adr := GenerateADRDocument(rep)
	if !strings.Contains(adr, "# ADR-001") {
		t.Errorf("Esperava cabeçalho ADR-001 no documento")
	}
	if !strings.Contains(adr, "solidtime/solidtime") {
		t.Errorf("Esperava nome do repositório recomendado no ADR")
	}
}
