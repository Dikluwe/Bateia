/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/lab_exporter_test.md
 * @layer L2
 * @updated 2026-08-04
 */

package report_test

import (
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/02_shell/report"
)

func TestLabExporter(t *testing.T) {
	dataset := &domain.LabDataset{
		Query:           "quicksort",
		GeneratedAt:     time.Now().UTC(),
		TotalChallenges: 1,
		Challenges: []domain.AlgorithmChallenge{
			{
				ID:          "lab-001",
				Name:        "QuickSort",
				Category:    "Sorting",
				Language:    "go",
				Repository:  "algorithms/sorting",
				FilePath:    "sort.go",
				CodeSnippet: "func QuickSort() {}",
				TestCases:   []string{"TestQuickSort"},
				Benchmarks:  []string{"BenchmarkQuickSort"},
				Complexity:  "O(N log N)",
			},
		},
	}

	jsonBytes, err := report.ExportLabJSON(dataset)
	if err != nil {
		t.Fatalf("Erro ao exportar Lab JSON: %v", err)
	}

	if len(jsonBytes) == 0 {
		t.Errorf("JSON do Lab veio vazio")
	}

	mdContent := report.GenerateLabMarkdown(dataset)
	if len(mdContent) == 0 {
		t.Errorf("Markdown do Lab veio vazio")
	}
}
