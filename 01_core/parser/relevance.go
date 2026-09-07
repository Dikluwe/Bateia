/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/relevance.md
 * @layer L1
 * @updated 2026-09-07
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
	ScoreGenericText    = 0

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
	normQuery := strings.ToLower(NormalizeQuery(query))
	queryWords := strings.Fields(normQuery)

	// 1. Check repository name match (ex: "Bateia" na query "Bateia Architectural Engine")
	if cand.Repo != "" {
		repoLower := strings.ToLower(cand.Repo)
		for _, w := range queryWords {
			if strings.ToLower(w) == repoLower {
				return RelevanceResult{
					Score:  ScoreModuleName,
					Label:  LabelModuleName,
					Reason: fmt.Sprintf("Campo candidato.repo: '%s'; termo da consulta: '%s'", cand.Repo, w),
				}
			}
		}
	}

	if manifest != nil {
		// 2. Check direct imports
		for _, dep := range manifest.Dependencies {
			if !dep.Indirect {
				if matchesDependency(normQuery, dep.Name) {
					return RelevanceResult{
						Score:  ScoreDirectImport,
						Label:  LabelDirectImport,
						Reason: fmt.Sprintf("Match exato de dependência direta no manifesto: '%s'", dep.Name),
					}
				}
			}
		}

		// 3. Check module name
		if manifest.ModuleName != "" {
			normModule := strings.ToLower(NormalizeQuery(manifest.ModuleName))
			if normModule == normQuery || path.Base(normModule) == normQuery {
				return RelevanceResult{
					Score:  ScoreModuleName,
					Label:  LabelModuleName,
					Reason: fmt.Sprintf("Campo manifesto.module_name: '%s'; consulta normalizada: '%s'", manifest.ModuleName, normQuery),
				}
			}
		}

		// 4. Check indirect imports
		for _, dep := range manifest.Dependencies {
			if dep.Indirect {
				if matchesDependency(normQuery, dep.Name) {
					return RelevanceResult{
						Score:  ScoreIndirectImport,
						Label:  LabelIndirectImport,
						Reason: fmt.Sprintf("Match de dependência indireta no manifesto: '%s'", dep.Name),
					}
				}
			}
		}

		// Content is legacy manifest data unless its source is explicitly README.
		if manifest.ContentPath != "" && strings.EqualFold(path.Base(manifest.ContentPath), "README.md") {
			matched := MatchingQueryTerms(query, string(manifest.Content))
			if len(matched) > 0 {
				return RelevanceResult{
					Score:  ScoreReadmeKeyword * len(matched) / len(QueryTerms(query)),
					Label:  LabelReadmeKeyword,
					Reason: fmt.Sprintf("Campo %s; termos literais %q (%d/%d); indício textual, implementação não verificada", manifest.ContentPath, matched, len(matched), len(QueryTerms(query))),
				}
			}
		}
	}

	// 5. Check structural directory matches
	for _, node := range tree {
		parts := strings.Split(node.Path, "/")
		for _, p := range parts {
			if strings.ToLower(NormalizeQuery(p)) == normQuery {
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
		Reason: "Candidato descoberto; nenhum termo pertinente verificado nos campos avaliados",
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
