/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/lab_extractor_test.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser_test

import (
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

func TestExtractLabDataset(t *testing.T) {
	report := &domain.SearchReport{
		Query: domain.QueryFilter{Query: "quicksort"},
		Results: []domain.ExtractedRepository{
			{
				Candidate: domain.RepositoryCandidate{Owner: "algorithms", Repo: "sorting"},
				Manifest:  &domain.Manifest{Language: "go"},
				Files: []domain.ExtractedFile{
					{
						Path:    "sort.go",
						Size:    200,
						Content: "package sort\n\nfunc QuickSort(arr []int) {\n\t// impl\n}",
						Symbols: []domain.ASTSymbol{
							{Name: "QuickSort", Kind: "function", Line: 3},
						},
					},
					{
						Path:    "sort_test.go",
						Size:    300,
						Content: "package sort\n\nfunc TestQuickSort(t *testing.T) {}\nfunc BenchmarkQuickSort(b *testing.B) {}",
						Symbols: []domain.ASTSymbol{
							{Name: "TestQuickSort", Kind: "function", Line: 3},
							{Name: "BenchmarkQuickSort", Kind: "function", Line: 4},
						},
					},
				},
			},
		},
	}

	labDataset := parser.ExtractLabDataset(report)

	if labDataset.TotalChallenges != 1 {
		t.Fatalf("Esperava 1 desafio no lab dataset, obteve %d", labDataset.TotalChallenges)
	}

	ch := labDataset.Challenges[0]
	if ch.Name != "QuickSort" {
		t.Errorf("Esperava nome QuickSort, obteve %s", ch.Name)
	}

	if ch.Category != "Sorting" {
		t.Errorf("Esperava categoria Sorting, obteve %s", ch.Category)
	}

	if len(ch.TestCases) == 0 {
		t.Errorf("Esperava encontrar casos de teste vinculados")
	}

	if len(ch.Benchmarks) == 0 {
		t.Errorf("Esperava encontrar benchmarks vinculados")
	}
}
