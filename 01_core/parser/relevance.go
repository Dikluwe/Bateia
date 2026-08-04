/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/relevance.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"fmt"
	"path"
	"strings"

	"bateia/01_core/domain"
)

const (
	ScoreDirectImport   = 100
	ScoreModuleName     = 85
	ScoreIndirectImport = 70
	ScoreReadmeKeyword  = 65
	ScoreStructural     = 50
	ScoreGenericText    = 20

	LabelDirectImport   = "DIRECT_IMPORT_MATCH"
	LabelModuleName     = "MODULE_NAME_MATCH"
	LabelIndirectImport = "INDIRECT_IMPORT_MATCH"
	LabelReadmeKeyword  = "README_KEYWORD_MATCH"
	LabelStructural     = "STRUCTURAL_MATCH"
	LabelGenericText    = "GENERIC_TEXT_MATCH"
)

type RelevanceResult struct {
	Score  int    `json:"score"`
	Label  string `json:"label"`
	Reason string `json:"reason"`
}

// EvaluateRelevance calcula a relevância do repositório em relação à query pesquisada.
func EvaluateRelevance(query string, manifest *domain.Manifest, cand domain.RepositoryCandidate, tree []domain.DirectoryNode) RelevanceResult {
	normQuery := NormalizeQuery(query)

	if manifest != nil {
		// 1. Check direct imports
		for _, dep := range manifest.Dependencies {
			if !dep.Indirect {
				if matchesDependency(normQuery, dep.Name) {
					return RelevanceResult{
						Score:  ScoreDirectImport,
						Label:  LabelDirectImport,
						Reason: fmt.Sprintf("Match exato de dependência direta no go.mod: '%s'", dep.Name),
					}
				}
			}
		}

		// 2. Check module name
		if manifest.ModuleName != "" {
			normModule := NormalizeQuery(manifest.ModuleName)
			if normModule == normQuery || path.Base(normModule) == normQuery {
				return RelevanceResult{
					Score:  ScoreModuleName,
					Label:  LabelModuleName,
					Reason: fmt.Sprintf("Match com o nome do módulo: '%s'", manifest.ModuleName),
				}
			}
		}

		// 3. Check indirect imports
		for _, dep := range manifest.Dependencies {
			if dep.Indirect {
				if matchesDependency(normQuery, dep.Name) {
					return RelevanceResult{
						Score:  ScoreIndirectImport,
						Label:  LabelIndirectImport,
						Reason: fmt.Sprintf("Match de dependência indireta no go.mod: '%s'", dep.Name),
					}
				}
			}
		}

		// 4. README content keyword scoring (README Score Booster)
		if manifest.Content != nil {
			readmeLower := strings.ToLower(string(manifest.Content))
			queryWords := strings.Fields(normQuery)
			matchedWords := 0
			for _, word := range queryWords {
				if len(word) >= 3 && strings.Contains(readmeLower, word) {
					matchedWords++
				}
			}
			if matchedWords > 0 && matchedWords >= len(queryWords)/2+1 {
				return RelevanceResult{
					Score:  ScoreReadmeKeyword,
					Label:  LabelReadmeKeyword,
					Reason: fmt.Sprintf("Match de %d/%d termos-chave no README do repositório", matchedWords, len(queryWords)),
				}
			}
		}
	}

	// 5. Check structural directory matches
	for _, node := range tree {
		parts := strings.Split(node.Path, "/")
		for _, p := range parts {
			if NormalizeQuery(p) == normQuery {
				return RelevanceResult{
					Score:  ScoreStructural,
					Label:  LabelStructural,
					Reason: fmt.Sprintf("Match na estrutura de diretórios: '%s'", node.Path),
				}
			}
		}
	}

	// 6. Generic text match (default)
	return RelevanceResult{
		Score:  ScoreGenericText,
		Label:  LabelGenericText,
		Reason: "Correspondência genérica de texto na descoberta",
	}
}

// NormalizeImportPath limpa o caminho de importação removendo sufixos de versão comuns (/v2, /v3) e espaços.
func NormalizeImportPath(q string) string {
	q = strings.ToLower(strings.TrimSpace(q))
	versionSuffixes := []string{"/v2", "/v3", "/v4", "/v5"}
	for _, s := range versionSuffixes {
		q = strings.TrimSuffix(q, s)
	}
	return q
}

func matchesDependency(normQuery, depName string) bool {
	normDep := NormalizeImportPath(depName)
	if normDep == normQuery {
		return true
	}
	if path.Base(normDep) == normQuery {
		return true
	}
	return false
}
