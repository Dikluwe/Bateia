/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/documentation_retrieval.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"

	"bateia/01_core/domain"
)

func documentationHit(path, hash string, decl *ast.FuncDecl, set *token.FileSet, terms map[string]bool) (domain.ImplementationHit, bool) {
	if decl.Doc == nil || len(terms) == 0 {
		return domain.ImplementationHit{}, false
	}
	words := implementationTerms(decl.Doc.Text())
	matched := []string{}
	for term := range terms {
		if words[term] {
			matched = append(matched, term)
		}
	}
	if len(matched) == 0 {
		return domain.ImplementationHit{}, false
	}
	sort.Strings(matched)
	score := len(matched)
	if score > 99 {
		score = 99
	}
	return domain.ImplementationHit{
		Retrieval: "documentation", Kind: "definition", Path: path, SHA256: hash,
		Symbol: decl.Name.Name, Line: set.Position(decl.Pos()).Line, EndLine: set.Position(decl.End()).Line,
		DocumentationLine: set.Position(decl.Doc.Pos()).Line, DocumentationEndLine: set.Position(decl.Doc.End()).Line,
		MatchedTerms: matched, Score: score,
		Reason: "indício documental em comentário Doc vinculado à declaração AST; declaração do autor, sem verificar implementação, execução ou cumprimento",
	}, true
}

// documentationBudget excludes every admitted lexical identity before the code cut.
// The returned candidates are also used for filling, so displaced code cannot
// reappear as documentary evidence and consume the reserved slots.
func documentationBudget(lexical, candidates []domain.ImplementationHit, limit int) ([]domain.ImplementationHit, int) {
	seen := map[string]bool{}
	exacts := 0
	for _, hit := range lexical {
		seen[documentationIdentity(hit)] = true
		if hit.Score >= 100000 {
			exacts++
		}
	}
	sortDocumentationHits(candidates)
	unique := make([]domain.ImplementationHit, 0, len(candidates))
	for _, hit := range candidates {
		identity := documentationIdentity(hit)
		if !seen[identity] {
			seen[identity] = true
			unique = append(unique, hit)
		}
	}
	reserve := limit / 4
	if reserve < 1 {
		reserve = 1
	}
	if reserve > len(unique) {
		reserve = len(unique)
	}
	codeMinimum := exacts
	if len(lexical) > 0 && codeMinimum < 1 {
		codeMinimum = 1
	}
	available := limit - codeMinimum
	if available < 0 {
		available = 0
	}
	if reserve > available {
		reserve = available
	}
	return unique, reserve
}

func documentationIdentity(hit domain.ImplementationHit) string {
	return fmt.Sprintf("%s\x00%d\x00%s", hit.Path, hit.Line, hit.Symbol)
}

func sortDocumentationHits(candidates []domain.ImplementationHit) {
	sort.Slice(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.Score != b.Score {
			return a.Score > b.Score
		}
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		return a.Symbol < b.Symbol
	})
}

func appendDocumentationHits(report *domain.ImplementationSearchReport, candidates []domain.ImplementationHit, limit int) {
	sortDocumentationHits(candidates)
	seen := map[string]bool{}
	for _, hit := range report.Hits {
		seen[documentationIdentity(hit)] = true
	}
	for _, hit := range candidates {
		identity := documentationIdentity(hit)
		if seen[identity] {
			continue
		}
		seen[identity] = true
		if len(report.Hits) >= limit {
			report.Limited = true
			report.Issues = append(report.Issues, "limite de resultados omitiu candidatos documentais")
			continue
		}
		report.Hits = append(report.Hits, hit)
	}
}
