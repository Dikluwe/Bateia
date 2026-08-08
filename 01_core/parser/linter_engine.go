/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/linter_engine.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"sync"

	"bateia/01_core/domain"
)

// LinterEngine orquestra a execução concorrente de regras arquiteturais em $L_1$.
type LinterEngine struct {
	workers int
}

// NewLinterEngine cria uma instância do motor concorrente com N goroutines.
func NewLinterEngine(workers int) *LinterEngine {
	if workers <= 0 {
		workers = 4
	}
	return &LinterEngine{workers: workers}
}

// RunConcurrent executa a análise arquitetural de forma paralela e deduplicada.
func (e *LinterEngine) RunConcurrent(manifest *domain.Manifest, tree []domain.DirectoryNode) []domain.ArchViolation {
	if manifest == nil {
		return nil
	}

	tasks := []func() []domain.ArchViolation{
		func() []domain.ArchViolation {
			return EvaluateArchViolations(manifest, tree)
		},
	}

	var wg sync.WaitGroup
	resultsChan := make(chan []domain.ArchViolation, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(t func() []domain.ArchViolation) {
			defer wg.Done()
			resultsChan <- t()
		}(task)
	}

	wg.Wait()
	close(resultsChan)

	var allViolations []domain.ArchViolation
	seen := make(map[string]bool)

	for res := range resultsChan {
		for _, v := range res {
			key := v.Rule + ":" + v.Target + ":" + v.Message
			if !seen[key] {
				seen[key] = true
				allViolations = append(allViolations, v)
			}
		}
	}

	return allViolations
}
