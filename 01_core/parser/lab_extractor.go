/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/lab_extractor.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"bateia/01_core/domain"
)

// ExtractLabDataset analisa o relatório extraído em memória (Zero I/O) e compila o dataset do Lab de Projetos.
func ExtractLabDataset(report *domain.SearchReport) *domain.LabDataset {
	if report == nil {
		return &domain.LabDataset{Challenges: []domain.AlgorithmChallenge{}}
	}
	dataset := &domain.LabDataset{
		Query:       report.Query.Query,
		GeneratedAt: report.GeneratedAt,
		Challenges:  make([]domain.AlgorithmChallenge, 0),
	}

	if report == nil || len(report.Results) == 0 {
		return dataset
	}

	challengeID := 1

	for _, res := range report.Results {
		repoName := fmt.Sprintf("%s/%s", res.Candidate.Owner, res.Candidate.Repo)

		var testFiles []domain.ExtractedFile
		var codeFiles []domain.ExtractedFile

		for _, file := range res.Files {
			lowerPath := strings.ToLower(file.Path)
			if strings.Contains(lowerPath, "test") || strings.Contains(lowerPath, "spec") || strings.HasSuffix(lowerPath, "_test.go") {
				testFiles = append(testFiles, file)
			} else {
				codeFiles = append(codeFiles, file)
			}
		}

		for _, file := range codeFiles {
			for _, sym := range file.Symbols {
				if sym.Kind != "function" && sym.Kind != "method" {
					continue
				}

				category := classifyAlgorithmCategory(sym.Name, file.Path)
				if category == "" && !isAlgorithmLike(sym.Name) {
					continue
				}
				if category == "" {
					category = "General"
				}

				var testCases []string
				var benchmarks []string

				for _, tf := range testFiles {
					for _, tSym := range tf.Symbols {
						lowerName := strings.ToLower(tSym.Name)
						lowerSym := strings.ToLower(sym.Name)

						if strings.HasPrefix(tSym.Name, "Benchmark") || strings.Contains(lowerName, "benchmark") {
							if strings.Contains(lowerName, lowerSym) || len(benchmarks) == 0 {
								benchmarks = append(benchmarks, fmt.Sprintf("%s (em %s)", tSym.Name, filepath.Base(tf.Path)))
							}
						} else if strings.HasPrefix(tSym.Name, "Test") || strings.Contains(lowerName, "test") {
							if strings.Contains(lowerName, lowerSym) || len(testCases) == 0 {
								testCases = append(testCases, fmt.Sprintf("%s (em %s)", tSym.Name, filepath.Base(tf.Path)))
							}
						}
					}
				}

				challenge := domain.AlgorithmChallenge{
					ID:          fmt.Sprintf("lab-%03d", challengeID),
					Name:        sym.Name,
					Category:    category,
					Language:    res.Manifest.Language,
					Repository:  repoName,
					FilePath:    file.Path,
					CodeSnippet: ExtractSymbolSnippet(file.Content, sym.Name, sym.Line),
					TestCases:   limitList(testCases, 3),
					Benchmarks:  limitList(benchmarks, 3),
					Complexity:  estimateComplexity(category),
				}

				dataset.Challenges = append(dataset.Challenges, challenge)
				challengeID++
			}
		}
	}

	dataset.TotalChallenges = len(dataset.Challenges)
	return dataset
}

func classifyAlgorithmCategory(name string, path string) string {
	combined := strings.ToLower(name + " " + path)

	switch {
	case strings.Contains(combined, "sort") || strings.Contains(combined, "quicksort") || strings.Contains(combined, "mergesort") || strings.Contains(combined, "heap"):
		return "Sorting"
	case strings.Contains(combined, "graph") || strings.Contains(combined, "bfs") || strings.Contains(combined, "dfs") || strings.Contains(combined, "dijkstra"):
		return "Graph"
	case strings.Contains(combined, "tree") || strings.Contains(combined, "binarytree") || strings.Contains(combined, "trie") || strings.Contains(combined, "ast"):
		return "Tree"
	case strings.Contains(combined, "dp") || strings.Contains(combined, "dynamic") || strings.Contains(combined, "fibonacci") || strings.Contains(combined, "knapsack"):
		return "DynamicProgramming"
	case strings.Contains(combined, "hash") || strings.Contains(combined, "crypto") || strings.Contains(combined, "cipher") || strings.Contains(combined, "sha256"):
		return "Crypto"
	case strings.Contains(combined, "math") || strings.Contains(combined, "matrix") || strings.Contains(combined, "prime") || strings.Contains(combined, "gcd"):
		return "Math"
	case strings.Contains(combined, "search") || strings.Contains(combined, "binarysearch") || strings.Contains(combined, "find"):
		return "Search"
	default:
		return ""
	}
}

func isAlgorithmLike(name string) bool {
	lower := strings.ToLower(name)
	prefixes := []string{"solve", "compute", "parse", "filter", "evaluate", "calculate", "process", "reduce", "transform"}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}
	return false
}

func ExtractSymbolSnippet(fileContent string, symName string, startLine int) string {
	if fileContent == "" {
		return ""
	}
	lines := strings.Split(fileContent, "\n")
	if startLine <= 0 || startLine > len(lines) {
		res, _ := TruncateSmart(fileContent, 1000, 15)
		return res
	}

	startIdx := startLine - 1
	endIdx := startIdx + 20
	if endIdx > len(lines) {
		endIdx = len(lines)
	}

	snippet := strings.Join(lines[startIdx:endIdx], "\n")
	result, _ := TruncateSmart(snippet, 1000, 15)
	return result
}

func limitList(items []string, max int) []string {
	if len(items) <= max {
		return items
	}
	return items[:max]
}

func estimateComplexity(category string) string {
	switch category {
	case "Sorting":
		return "O(N log N)"
	case "Search":
		return "O(log N)"
	case "Graph":
		return "O(V + E)"
	case "DynamicProgramming":
		return "O(N^2)"
	case "Tree":
		return "O(log N)"
	default:
		return "O(N)"
	}
}
