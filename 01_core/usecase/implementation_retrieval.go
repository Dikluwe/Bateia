/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/implementation_retrieval.md
 * @layer L1
 * @updated 2026-09-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"context"
	"crypto/sha256"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"strings"
)

func RetrieveImplementationEvidence(ctx context.Context, provider domain.ContentProvider, candidate domain.RepositoryCandidate, query string) domain.ImplementationSearchReport {
	budget := domain.ImplementationSearchBudget{MaxFiles: 2048, MaxBytes: 32 * 1024 * 1024, MaxFileBytes: 512 * 1024, MaxResults: 24}
	searcher, ok := provider.(domain.ImplementationSearcher)
	if !ok {
		return domain.ImplementationSearchReport{Status: "not_evaluated", Method: "implementation-search", Reason: "Provedor sem busca de implementação; seleção limitada a caminhos.", Budget: budget}
	}
	report, err := searcher.SearchImplementations(ctx, candidate, query, budget)
	if err != nil {
		report.Status = "partial"
		report.Issues = append(report.Issues, err.Error())
		report.Reason = "Busca de implementação incompleta; não demonstra ausência."
	}
	return report
}

func SelectImplementationPaths(query string, tree []domain.DirectoryNode, report domain.ImplementationSearchReport, limit int) []string {
	base := parser.SelectEvidenceFiles(query, tree, limit)
	if limit < 0 {
		limit = 0
	}
	out := []string{}
	seen := map[string]bool{}
	present := map[string]bool{}
	for _, n := range tree {
		if n.Type == "file" || n.Type == "blob" {
			present[n.Path] = true
		}
	}
	add := func(p string) {
		if !seen[p] {
			out = append(out, p)
			seen[p] = true
		}
	}
	for _, p := range base {
		if !strings.Contains(p, "/") && (parser.DetectManifestType(p) != "" || strings.EqualFold(p, "README.md")) {
			add(p)
		}
	}
	contextCount := len(out)
	for _, h := range report.Hits {
		if len(out)-contextCount >= limit {
			break
		}
		if present[h.Path] && validMechanismPath(h.Path) {
			add(h.Path)
		}
	}
	for _, p := range base {
		if len(out)-contextCount >= limit {
			break
		}
		add(p)
	}
	return out
}

// ValidateImplementationHits checks the second read against the indexed bytes.
func ValidateImplementationHits(report domain.ImplementationSearchReport, files []domain.ExtractedFile) domain.ImplementationSearchReport {
	out := report
	out.Hits = []domain.ImplementationHit{}
	out.Issues = append([]string(nil), report.Issues...)
	contents := map[string]domain.ExtractedFile{}
	hashes := map[string]string{}
	calls := map[string]map[string]bool{}
	docs := map[string]map[string]bool{}
	for _, f := range files {
		contents[f.Path] = f
		hashes[f.Path] = fmt.Sprintf("%x", sha256.Sum256([]byte(f.Content)))
	}
	for _, h := range report.Hits {
		f, ok := contents[h.Path]
		if !ok {
			out.Status = "partial"
			out.Issues = append(out.Issues, "hit não extraído: "+h.Path)
			continue
		}
		hash := hashes[h.Path]
		if hash != h.SHA256 {
			out.Status = "partial"
			out.Issues = append(out.Issues, "índice obsoleto: bytes divergentes em "+h.Path)
			continue
		}
		if h.Line < 1 || h.EndLine < h.Line || h.EndLine > len(strings.Split(f.Content, "\n")) {
			out.Status = "partial"
			out.Issues = append(out.Issues, "intervalo de definição inválido: "+h.Path)
			continue
		}
		if h.Retrieval == "documentation" {
			if docs[h.Path] == nil {
				docs[h.Path] = implementationDocumentationDefinitions(f.Content, h.Path)
			}
			if !docs[h.Path][implementationDocumentationKey(h.Symbol, h.Line, h.EndLine, h.DocumentationLine, h.DocumentationEndLine)] {
				out.Status = "partial"
				out.Issues = append(out.Issues, "documentação não vinculada à declaração: "+h.Path)
				continue
			}
		}
		references := h.RelatedTo
		h.RelatedTo = nil
		for _, ref := range references {
			caller, exists := contents[ref.Path]
			if exists && hashes[ref.Path] == ref.SHA256 && calls[ref.Path] == nil {
				calls[ref.Path] = implementationCallReferences(caller.Content, ref.Path)
			}
			if !exists || hashes[ref.Path] != ref.SHA256 || !calls[ref.Path][implementationCallKey(ref.Symbol, h.Symbol, ref.Line)] {
				out.Status = "partial"
				out.Issues = append(out.Issues, "referência de chamada não confirmada: "+ref.Path)
				continue
			}
			h.RelatedTo = append(h.RelatedTo, ref)
		}
		if h.Retrieval == "direct_call" && len(h.RelatedTo) == 0 {
			out.Status = "partial"
			out.Issues = append(out.Issues, "hit de chamada sem referência confirmada: "+h.Path)
			continue
		}
		out.Hits = append(out.Hits, h)
	}
	if out.Status == "partial" {
		out.Reason = "Busca ou releitura parcial; somente hits vinculados aos bytes extraídos são apresentados. Ausência não demonstrada."
	}
	return out
}

// A line in the file is insufficient: it must be the claimed direct call in the caller.
// Parse each confirmed caller only once, even when it supports several destinations.
func implementationCallReferences(content, filePath string) map[string]bool {
	out := map[string]bool{}
	set := token.NewFileSet()
	tree, err := goparser.ParseFile(set, filePath, content, goparser.AllErrors)
	if err != nil {
		return out
	}
	for _, decl := range tree.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if _, closure := n.(*ast.FuncLit); closure {
				return false
			}
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			id, ok := call.Fun.(*ast.Ident)
			if !ok {
				return true
			}
			if id.Obj != nil {
				target, function := id.Obj.Decl.(*ast.FuncDecl)
				if id.Obj.Kind != ast.Fun || !function || target.Recv != nil {
					return true
				}
			}
			out[implementationCallKey(fn.Name.Name, id.Name, set.Position(call.Pos()).Line)] = true
			return true
		})
	}
	return out
}
func implementationCallKey(caller, callee string, line int) string {
	return fmt.Sprintf("%s\x00%s\x00%d", caller, callee, line)
}

// Documentation must be the AST Doc group of the claimed function, not arbitrary nearby text.
func implementationDocumentationDefinitions(content, filePath string) map[string]bool {
	out := map[string]bool{}
	set := token.NewFileSet()
	tree, err := goparser.ParseFile(set, filePath, content, goparser.ParseComments|goparser.AllErrors)
	if err != nil || ast.IsGenerated(tree) {
		return out
	}
	for _, decl := range tree.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Doc == nil {
			continue
		}
		out[implementationDocumentationKey(fn.Name.Name, set.Position(fn.Pos()).Line, set.Position(fn.End()).Line, set.Position(fn.Doc.Pos()).Line, set.Position(fn.Doc.End()).Line)] = true
	}
	return out
}
func implementationDocumentationKey(symbol string, line, end, docLine, docEnd int) string {
	return fmt.Sprintf("%s\x00%d\x00%d\x00%d\x00%d", symbol, line, end, docLine, docEnd)
}
