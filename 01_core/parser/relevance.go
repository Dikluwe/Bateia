/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/relevance.md
 * @layer L1
 * @updated 2026-07-31
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
	ScoreStructural     = 50
	ScoreGenericText    = 20

	LabelDirectImport   = "DIRECT_IMPORT_MATCH"
	LabelModuleName     = "MODULE_NAME_MATCH"
	LabelIndirectImport = "INDIRECT_IMPORT_MATCH"
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
	}

	// 4. Check structural directory matches
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

	// 5. Generic text match (default)
	return RelevanceResult{
		Score:  ScoreGenericText,
		Label:  LabelGenericText,
		Reason: "Correspondência genérica de texto na descoberta",
	}
}

// NormalizeQuery limpa a query removendo sufixos de versão comuns (/v2, /v3) e espaços.
func NormalizeQuery(q string) string {
	q = strings.ToLower(strings.TrimSpace(q))
	versionSuffixes := []string{"/v2", "/v3", "/v4", "/v5"}
	for _, s := range versionSuffixes {
		q = strings.TrimSuffix(q, s)
	}
	return q
}

func matchesDependency(normQuery, depName string) bool {
	normDep := NormalizeQuery(depName)
	if normDep == normQuery {
		return true
	}
	if path.Base(normDep) == normQuery {
		return true
	}
	return false
}
