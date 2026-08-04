/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_diff.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"bateia/01_core/domain"
)

// ComputeArchDiff calcula o delta arquitetural entre dois snapshots de um repositório.
func ComputeArchDiff(oldRepo, newRepo *domain.ExtractedRepository) domain.ArchDiff {
	diff := domain.ArchDiff{
		OldRef: oldRepo.Ref,
		NewRef: newRepo.Ref,
	}

	oldDeps := make(map[string]domain.Dependency)
	if oldRepo.Manifest != nil {
		for _, d := range oldRepo.Manifest.Dependencies {
			oldDeps[d.Name] = d
		}
	}

	newDeps := make(map[string]domain.Dependency)
	if newRepo.Manifest != nil {
		for _, d := range newRepo.Manifest.Dependencies {
			newDeps[d.Name] = d
		}
	}

	// Dependências adicionadas
	for name, dep := range newDeps {
		if _, exists := oldDeps[name]; !exists {
			diff.AddedDeps = append(diff.AddedDeps, dep)
		}
	}

	// Dependências removidas
	for name, dep := range oldDeps {
		if _, exists := newDeps[name]; !exists {
			diff.RemovedDeps = append(diff.RemovedDeps, dep)
		}
	}

	// Símbolos AST
	oldSyms := make(map[string]bool)
	for _, s := range oldRepo.Symbols {
		oldSyms[s.Name] = true
	}
	newSyms := make(map[string]bool)
	for _, s := range newRepo.Symbols {
		newSyms[s.Name] = true
	}

	for _, s := range newRepo.Symbols {
		if !oldSyms[s.Name] {
			diff.AddedSymbols = append(diff.AddedSymbols, s)
		}
	}
	for _, s := range oldRepo.Symbols {
		if !newSyms[s.Name] {
			diff.RemovedSymbols = append(diff.RemovedSymbols, s)
		}
	}

	return diff
}
