/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/requirement_retrieval.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser

import (
	"encoding/hex"
	"fmt"
	"path"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

// DescribeImplementationRequirements projects already revalidated hits. Token
// overlap identifies possible reading material, not a satisfied capability.
func DescribeImplementationRequirements(query string, report domain.ImplementationSearchReport) []domain.RequirementRetrieval {
	queries := InterpretImplementationQuery(query)
	out := make([]domain.RequirementRetrieval, 0, len(queries))
	for _, interpretation := range queries {
		line := interpretation.Original
		r := domain.RequirementRetrieval{Requirement: line, Status: "not_verified", MatchedTerms: []string{}, MissingTerms: []string{}, Sources: []domain.ImplementationReference{}, Constraints: append([]domain.ImplementationQueryConstraint(nil), interpretation.Constraints...)}
		terms := map[string]bool{}
		for _, term := range interpretation.Terms {
			terms[term] = true
		}
		matched := map[string]bool{}
		seen := map[domain.ImplementationReference]bool{}
		documentationSeen := map[domain.ImplementationReference]bool{}
		for _, hit := range report.Hits {
			if (hit.Retrieval != "lexical" && hit.Retrieval != "documentation") || !requirementSourceIdentity(hit) {
				continue
			}
			if hit.Retrieval == "documentation" && (interpretation.Mode == "symbol" || hit.DocumentationLine < 1 || hit.DocumentationEndLine < hit.DocumentationLine || hit.DocumentationEndLine >= hit.Line || hit.EndLine < hit.Line) {
				continue
			}
			field := map[string]bool{}
			for _, term := range hit.MatchedTerms {
				field[strings.ToLower(strings.TrimSpace(term))] = true
			}
			if admitted, _ := implementationQueryActionAdmission(interpretation, field); !admitted {
				continue
			}
			intersection := false
			for _, term := range hit.MatchedTerms {
				term = strings.ToLower(strings.TrimSpace(term))
				if terms[term] {
					matched[term] = true
					intersection = true
				}
			}
			if !intersection {
				continue
			}
			source := domain.ImplementationReference{Path: hit.Path, Symbol: hit.Symbol, Line: hit.Line, SHA256: hit.SHA256}
			if hit.Retrieval == "documentation" {
				if !documentationSeen[source] {
					r.DocumentationSources = append(r.DocumentationSources, source)
					documentationSeen[source] = true
				}
			} else if !seen[source] {
				r.Sources = append(r.Sources, source)
				seen[source] = true
			}
		}
		for term := range terms {
			if matched[term] {
				r.MatchedTerms = append(r.MatchedTerms, term)
			} else {
				r.MissingTerms = append(r.MissingTerms, term)
			}
		}
		sort.Strings(r.MatchedTerms)
		sort.Strings(r.MissingTerms)
		for _, sources := range [][]domain.ImplementationReference{r.Sources, r.DocumentationSources} {
			sort.Slice(sources, func(i, j int) bool {
				a, b := sources[i], sources[j]
				if a.Path != b.Path {
					return a.Path < b.Path
				}
				if a.Line != b.Line {
					return a.Line < b.Line
				}
				if a.Symbol != b.Symbol {
					return a.Symbol < b.Symbol
				}
				return a.SHA256 < b.SHA256
			})
		}
		if len(r.Sources) > 0 {
			r.Status = "evidence_located"
			r.Reason = fmt.Sprintf("Inferência lexical com fontes de código: %d de %d termos associados aos canais localizados; requisito e comportamento não verificados.", len(r.MatchedTerms), len(terms))
			if len(r.DocumentationSources) > 0 {
				r.Reason += " Fontes documentais separadas registram declarações do autor, não prova de implementação."
			}
		} else if len(r.DocumentationSources) > 0 {
			r.Status = "documentation_located"
			r.Reason = fmt.Sprintf("Indício documental: %d de %d termos em comentários Doc; declaração do autor, sem fonte lexical de código; requisito e comportamento não verificados.", len(r.MatchedTerms), len(terms))
		} else {
			r.Reason = "Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade."
		}
		if len(r.Constraints) > 0 {
			r.Reason += " Restrições preservadas; escopo e cumprimento da negação não foram verificados."
		}
		if _, filtered := implementationQueryActionAdmission(interpretation, nil); filtered {
			r.Reason += " Filtro lexical de admissão por ação e contexto aplicado separadamente a cada fonte."
		}
		if report.Status != "complete" || report.Limited {
			r.Reason += " Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência."
		}
		out = append(out, r)
	}
	return out
}

func requirementSourceIdentity(hit domain.ImplementationHit) bool {
	if hit.Line < 1 || strings.TrimSpace(hit.Symbol) == "" || hit.Path == "" || strings.Contains(hit.Path, "\\") || path.Clean(hit.Path) != hit.Path || !ImplementationSearchPathAllowed(hit.Path, "tests") {
		return false
	}
	bytes, err := hex.DecodeString(hit.SHA256)
	return err == nil && len(bytes) == 32
}
