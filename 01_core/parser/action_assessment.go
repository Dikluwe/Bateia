/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/action_assessment.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"sort"
	"strings"
)

const actionHintLimits = "Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade."

// AssessImplementationActions rebuilds hypotheses without trusting prior projections.
func AssessImplementationActions(query string, report domain.ImplementationSearchReport, files []domain.ExtractedFile) []domain.RequirementActionAssessment {
	byPath := map[string][]domain.ExtractedFile{}
	for _, file := range files {
		byPath[file.Path] = append(byPath[file.Path], file)
	}
	hits := append([]domain.ImplementationHit(nil), report.Hits...)
	sort.Slice(hits, func(i, j int) bool {
		a, b := hits[i], hits[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Symbol != b.Symbol {
			return a.Symbol < b.Symbol
		}
		if a.SHA256 != b.SHA256 {
			return a.SHA256 < b.SHA256
		}
		if a.EndLine != b.EndLine {
			return a.EndLine < b.EndLine
		}
		if a.Retrieval != b.Retrieval {
			return a.Retrieval < b.Retrieval
		}
		return a.Kind < b.Kind
	})
	out := []domain.RequirementActionAssessment{}
	for _, q := range InterpretImplementationQuery(query) {
		a := domain.RequirementActionAssessment{Query: q.Original, Status: "not_evaluated", Candidates: []domain.ActionCandidate{}}
		action, reason := requiredAction(q)
		a.Reason = reason + " " + actionHintLimits
		if action != "" {
			a.RequiredAction, a.Status = action, "hypothesis"
			for _, hit := range hits {
				c := domain.ActionCandidate{Source: domain.ImplementationReference{Path: hit.Path, Symbol: hit.Symbol, Line: hit.Line, SHA256: hit.SHA256}, EndLine: hit.EndLine, Relation: "unknown", Evidence: []domain.ActionEvidence{}}
				c.Reason = assessActionHit(hit, byPath[hit.Path], action, &c) + " " + actionHintLimits
				a.Candidates = append(a.Candidates, c)
			}
		}
		out = append(out, a)
	}
	return out
}

func requiredAction(q domain.ImplementationQuery) (string, string) {
	if q.Mode != "natural" {
		return "", "Modo da consulta não natural; ação não avaliada."
	}
	if len(q.Constraints) > 0 {
		return "", "Restrições ou negação presentes; ação não avaliada."
	}
	start, stop := false, false
	for _, word := range strings.Fields(discoveryNormalized(q.Original)) {
		switch word {
		case "sem", "nao", "nunca", "not", "never", "without":
			return "", "Negação presente; ação não avaliada."
		case "iniciar", "comecar", "start", "subscribe":
			start = true
		case "parar", "interromper", "cancelar", "cessar", "stop", "cancel", "unsubscribe":
			stop = true
		}
	}
	if start && stop {
		return "", "Famílias start e stop ambíguas na consulta."
	}
	if start {
		return "start", "Hipótese lexical da ação solicitada: start."
	}
	if stop {
		return "stop", "Hipótese lexical da ação solicitada: stop."
	}
	return "", "Nenhuma família de ação reconhecida na consulta."
}

func assessActionHit(hit domain.ImplementationHit, files []domain.ExtractedFile, required string, candidate *domain.ActionCandidate) string {
	if hit.Retrieval == "documentation" {
		return "Origem documental não autoriza indícios da implementação."
	}
	if !strings.HasSuffix(hit.Path, ".go") {
		return "Linguagem não suportada; somente Go."
	}
	if hit.Kind != "definition" {
		return "Candidato não é definição de função."
	}
	if len(files) == 0 {
		return "Bytes completos ausentes do arquivo."
	}
	if len(files) != 1 {
		return "Arquivo duplicado; identidade ambígua."
	}
	file := files[0]
	if file.Size > 0 && file.Size != int64(len(file.Content)) {
		return "Tamanho diverge dos bytes completos."
	}
	if fmt.Sprintf("%x", sha256.Sum256([]byte(file.Content))) != hit.SHA256 {
		return "SHA256 diverge dos bytes fornecidos."
	}
	set := token.NewFileSet()
	tree, err := goparser.ParseFile(set, file.Path, file.Content, goparser.AllErrors)
	if err != nil {
		return "AST Go inválido; erro de parsing."
	}
	var matched *ast.FuncDecl
	for _, decl := range tree.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		if fn.Name.Name == hit.Symbol && set.PositionFor(fn.Pos(), false).Line == hit.Line && set.PositionFor(fn.End(), false).Line == hit.EndLine {
			if matched != nil {
				return "Declaração ambígua para símbolo e intervalo físico."
			}
			matched = fn
		}
	}
	if matched == nil {
		return "Símbolo ou intervalo físico não corresponde a função completa com corpo."
	}
	type observed struct {
		evidence domain.ActionEvidence
		offset   int
	}
	observations := []observed{}
	seen := map[string]bool{}
	overflow := false
	add := func(kind string, id *ast.Ident) {
		if overflow {
			return
		}
		words := implementationWords(id.Name)
		for _, word := range words {
			switch word {
			case "not", "never", "without", "no":
				return
			}
		}
		pos, end := set.PositionFor(id.Pos(), false), set.PositionFor(id.End(), false)
		for _, word := range words {
			action := ""
			switch word {
			case "add", "start", "subscribe", "register", "watch":
				action = "start"
			case "remove", "stop", "unsubscribe", "unregister", "cancel", "close", "unwatch":
				action = "stop"
			}
			if action == "" {
				continue
			}
			key := fmt.Sprintf("%s/%s/%s/%d/%d", kind, action, word, pos.Offset, end.Offset)
			if seen[key] {
				continue
			}
			seen[key] = true
			if len(observations) == 32 {
				overflow = true
				return
			}
			observations = append(observations, observed{domain.ActionEvidence{Kind: kind, Action: action, Term: word, Line: pos.Line, EndLine: end.Line, Snippet: file.Content[pos.Offset:end.Offset]}, pos.Offset})
		}
	}
	add("declaration_name", matched.Name)
	ast.Inspect(matched.Body, func(node ast.Node) bool {
		if overflow {
			return false
		}
		switch n := node.(type) {
		case *ast.FuncLit:
			return false
		case *ast.CallExpr:
			switch callee := n.Fun.(type) {
			case *ast.Ident:
				add("call_name", callee)
			case *ast.SelectorExpr:
				add("call_name", callee.Sel)
			}
		}
		return true
	})
	if overflow {
		return "Orçamento de 32 evidências excedido; prefixo parcial descartado, ação desconhecida."
	}
	sort.Slice(observations, func(i, j int) bool {
		a, b := observations[i], observations[j]
		if a.offset != b.offset {
			return a.offset < b.offset
		}
		if a.evidence.Kind != b.evidence.Kind {
			return a.evidence.Kind < b.evidence.Kind
		}
		if a.evidence.Action != b.evidence.Action {
			return a.evidence.Action < b.evidence.Action
		}
		return a.evidence.Term < b.evidence.Term
	})
	aligned, opposed := false, false
	for _, o := range observations {
		candidate.Evidence = append(candidate.Evidence, o.evidence)
		if o.evidence.Action == required {
			aligned = true
		} else {
			opposed = true
		}
	}
	switch {
	case aligned && opposed:
		candidate.Relation = "conflicting_hints"
	case aligned:
		candidate.Relation = "aligned_hint"
	case opposed:
		candidate.Relation = "opposed_hint"
	}
	return "Relação por indícios de nomes de declaração e chamadas: " + candidate.Relation + "."
}
