/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/funnel.md
 * @layer L1
 * @updated 2026-08-23
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestValidateQueryFilter(t *testing.T) {
	valid := domain.QueryFilter{Query: "architecture", Language: "go,rust", MaxResults: 20, Depth: 2, DirectRepos: []string{"owner/repo"}}
	if err := ValidateQueryFilter(valid); err != nil {
		t.Fatalf("filtro válido rejeitado: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*domain.QueryFilter)
	}{
		{"empty query", func(f *domain.QueryFilter) { f.Query = " " }},
		{"zero max", func(f *domain.QueryFilter) { f.MaxResults = 0 }},
		{"excessive max", func(f *domain.QueryFilter) { f.MaxResults = MaxDiscoveryResults + 1 }},
		{"zero depth", func(f *domain.QueryFilter) { f.Depth = 0 }},
		{"excessive depth", func(f *domain.QueryFilter) { f.Depth = MaxSearchDepth + 1 }},
		{"unknown language", func(f *domain.QueryFilter) { f.Language = "brainfuck" }},
		{"invalid repo", func(f *domain.QueryFilter) { f.DirectRepos = []string{"owner/repo/extra"} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := valid
			tt.mutate(&filter)
			if err := ValidateQueryFilter(filter); err == nil {
				t.Fatal("esperava erro de validação")
			}
		})
	}
}
