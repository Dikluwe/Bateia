/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_diff.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestComputeArchDiff(t *testing.T) {
	oldRepo := &domain.ExtractedRepository{
		Ref: domain.CommitRef{Owner: "org", Repo: "app", Commit: "v1.0.0"},
		Manifest: &domain.Manifest{
			Dependencies: []domain.Dependency{
				{Name: "lib-a", Version: "1.0"},
			},
		},
	}
	newRepo := &domain.ExtractedRepository{
		Ref: domain.CommitRef{Owner: "org", Repo: "app", Commit: "v2.0.0"},
		Manifest: &domain.Manifest{
			Dependencies: []domain.Dependency{
				{Name: "lib-b", Version: "2.0"},
			},
		},
	}

	diff := ComputeArchDiff(oldRepo, newRepo)
	if len(diff.AddedDeps) != 1 || diff.AddedDeps[0].Name != "lib-b" {
		t.Errorf("Esperava dependência lib-b adicionada")
	}
	if len(diff.RemovedDeps) != 1 || diff.RemovedDeps[0].Name != "lib-a" {
		t.Errorf("Esperava dependência lib-a removida")
	}
}
