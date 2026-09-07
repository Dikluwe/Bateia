/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query_evidence.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// BuildQueryEvidence selects verified source intervals, then follows direct static calls.
// Lexical anchors and graph proximity are hypotheses of pertinence, not behavioral proof.
func buildUnitEvidence(query string, semantic *domain.SemanticReport, files []domain.ExtractedFile, root string) *domain.QueryEvidence {
	out := &domain.QueryEvidence{Status: "insufficient", Reason: "Nenhuma definição com corpo e hash confirmado ancorada na consulta."}
	if semantic == nil {
		return out
	}
	bodies := map[string]domain.ExtractedFile{}
	bodyHashes := map[string]string{}
	for _, f := range files {
		key := strings.TrimSuffix(root, "/") + "/" + f.Path
		bodies[key] = f
		bodyHashes[key] = fmt.Sprintf("%x", sha256.Sum256([]byte(f.Content)))
	}
	type ranked struct {
		ev    domain.ImplementationEvidence
		score int
	}
	candidates := map[string]ranked{}
	calls := map[string]domain.SemanticCall{}
	for _, u := range semantic.Units {
		if u.Status != "complete" {
			continue
		}
		hashes := map[string]string{}
		for _, in := range u.Inputs {
			hashes[in.Path] = in.SHA256
		}
		for _, s := range u.Symbols {
			if !s.Definition || !(s.Kind == "FunctionDecl" || s.Kind == "CXXMethodDecl") || strings.HasPrefix(s.Name, "operator") {
				continue
			}
			f, ok := bodies[s.File]
			if !ok {
				continue
			}
			hash := bodyHashes[s.File]
			if hashes[s.File] != hash {
				continue
			}
			lines := strings.Split(f.Content, "\n")
			if s.Line < 1 || s.EndLine < s.Line || s.EndLine > len(lines) {
				continue
			}
			activeLines := append([]string(nil), lines...)
			if u.Flow != nil && u.Flow.Status == "complete" {
				for _, r := range u.Flow.Skipped {
					if r.File == s.File {
						lo, hi := r.Line, r.EndLine
						if lo < 1 {
							lo = 1
						}
						if hi > len(activeLines) {
							hi = len(activeLines)
						}
						for n := lo; n <= hi; n++ {
							activeLines[n-1] = ""
						}
					}
				}
			}
			body := strings.Join(activeLines[s.Line-1:s.EndLine], "\n")
			terms := MatchingQueryTerms(query, body)
			score := len(terms) * 1000
			if score > 0 {
				penalty := s.EndLine - s.Line
				if penalty > 900 {
					penalty = 900
				}
				score -= penalty
			}
			nearCall := false
			for _, call := range u.Calls {
				if call.CallerID != s.ID || call.File != s.File {
					continue
				}
				lo, hi := call.Line-8, call.Line+2
				if lo < s.Line {
					lo = s.Line
				}
				if hi > s.EndLine {
					hi = s.EndLine
				}
				if lo >= 1 && hi >= lo && len(MatchingQueryTerms(query, strings.Join(activeLines[lo-1:hi], "\n"))) > 0 {
					nearCall = true
				}
			}
			if !nearCall {
				score = 0
			}
			for _, term := range QueryTerms(query) {
				if strings.Contains(strings.ToLower(s.Name), term) {
					if nearCall {
						score += 30
					}
				}
			}
			// Display a bounded window around the first body match; retain full examined range.
			first := s.Line
			for i := s.Line - 1; i < s.EndLine; i++ {
				if len(MatchingQueryTerms(query, activeLines[i])) > 0 {
					first = i + 1
					break
				}
			}
			if first > s.Line+2 {
				first -= 2
			} else {
				first = s.Line
			}
			last := first + 11
			if last > s.EndLine {
				last = s.EndLine
			}
			ev := domain.ImplementationEvidence{Symbol: s, SHA256: hash, MatchedTerms: terms, Selection: "âncora lexical no corpo/nome (inferência)", Snippet: strings.Join(lines[first-1:last], "\n"), SnippetFirstLine: first, SnippetLastLine: last}
			prev, exists := candidates[s.ID]
			if !exists || score > prev.score || (score == prev.score && s.File < prev.ev.Symbol.File) {
				candidates[s.ID] = ranked{ev, score}
			}
		}
		for _, c := range u.Calls {
			hash, ok := bodyHashes[c.File]
			if !ok || hashes[c.File] != hash {
				continue
			}
			calls[fmt.Sprintf("%s:%d:%s:%s", c.File, c.Line, c.CallerID, c.CalleeID)] = c
		}
	}
	var ordered []ranked
	for _, c := range candidates {
		if c.score > 0 {
			ordered = append(ordered, c)
		}
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].score != ordered[j].score {
			return ordered[i].score > ordered[j].score
		}
		a, b := ordered[i].ev.Symbol, ordered[j].ev.Symbol
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		return a.ID < b.ID
	})
	selected := map[string]bool{}
	for i, c := range ordered {
		if i == 6 {
			break
		}
		selected[c.ev.Symbol.ID] = true
		out.Definitions = append(out.Definitions, c.ev)
	}
	// Traverse each anchor's call sites in source order, preserving a mechanism's
	// descendants before unrelated anchors consume the global display budget.
	adjacency := map[string][]domain.SemanticCall{}
	for _, c := range calls {
		adjacency[c.CallerID] = append(adjacency[c.CallerID], c)
	}
	for id := range adjacency {
		sort.Slice(adjacency[id], func(i, j int) bool {
			a, b := adjacency[id][i], adjacency[id][j]
			if a.File != b.File {
				return a.File < b.File
			}
			if a.Line != b.Line {
				return a.Line < b.Line
			}
			return a.CalleeID < b.CalleeID
		})
	}
	var visit func(string, int)
	visit = func(id string, depth int) {
		if depth == 0 {
			return
		}
		for _, call := range adjacency[id] {
			c, ok := candidates[call.CalleeID]
			if !ok {
				continue
			}
			if !selected[call.CalleeID] {
				if len(out.Definitions) >= 24 {
					return
				}
				c.ev.Selection = "chamada direta de uma definição selecionada (inferência de pertinência)"
				selected[call.CalleeID] = true
				out.Definitions = append(out.Definitions, c.ev)
			}
			visit(call.CalleeID, depth-1)
		}
	}
	seeds := append([]domain.ImplementationEvidence(nil), out.Definitions...)
	for _, c := range seeds {
		visit(c.Symbol.ID, 2)
	}
	for _, c := range calls {
		if selected[c.CallerID] && selected[c.CalleeID] {
			out.Calls = append(out.Calls, c)
		}
	}
	sort.Slice(out.Calls, func(i, j int) bool {
		a, b := out.Calls[i], out.Calls[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		return a.CallerID+a.CalleeID < b.CallerID+b.CalleeID
	})
	if len(out.Calls) > 64 {
		out.Calls = out.Calls[:64]
	}
	for i := range out.Definitions {
		e := &out.Definitions[i]
		f := bodies[e.Symbol.File]
		e.Excerpts = mechanismExcerpts(*e, f.Content, out.Calls)
	}
	if len(out.Definitions) > 0 {
		out.Status = "implementation_evidence"
		out.Reason = "Fato observado: definições e chamadas estáticas com hashes confirmados. Inferência: seleção lexical (até 6 âncoras) e até 2 saltos diretos, máximo 24 definições e 64 chamadas. Hipótese: pertinência à consulta. Não verificado: adequação completa, execução, GPU e desempenho."
	}
	return out
}

// BuildQueryEvidence keeps all projected calls in one compiler context.
func BuildQueryEvidence(query string, semantic *domain.SemanticReport, files []domain.ExtractedFile, root string) *domain.QueryEvidence {
	best := &domain.QueryEvidence{Status: "insufficient", Reason: "Nenhuma definição ancorada em chamadas estáticas com bytes confirmados."}
	bestScore := -1
	if semantic == nil {
		return best
	}
	for _, unit := range semantic.Units {
		if unit.Status != "complete" {
			continue
		}
		one := *semantic
		one.Units = []domain.SemanticUnit{unit}
		ev := buildUnitEvidence(query, &one, files, root)
		if len(ev.Definitions) == 0 {
			continue
		}
		first := ev.Definitions[0]
		score := len(first.MatchedTerms)*1000 - (first.Symbol.EndLine - first.Symbol.Line)
		if score > bestScore || (score == bestScore && unit.Path < best.Context) {
			bestScore = score
			best = ev
			best.Context = unit.Path
			best.ContextSHA256 = fmt.Sprintf("%x", sha256.Sum256([]byte(unit.Path+"\x00"+strings.Join(unit.Arguments, "\x00"))))
		}
	}
	return best
}

func mechanismExcerpts(e domain.ImplementationEvidence, content string, calls []domain.SemanticCall) []domain.EvidenceExcerpt {
	lines := strings.Split(content, "\n")
	points := map[int]string{e.SnippetFirstLine: "âncora textual"}
	for n := e.Symbol.Line; n <= e.Symbol.EndLine; n++ {
		t := strings.TrimSpace(lines[n-1])
		if strings.HasPrefix(t, "if (") || strings.HasPrefix(t, "if(") || strings.HasPrefix(t, "for(") || strings.HasPrefix(t, "for (") || strings.HasPrefix(t, "else") || strings.HasPrefix(t, "switch (") {
			points[n] = "condição/laço textual; fluxo não verificado"
		}
	}
	for _, c := range calls {
		if c.CallerID == e.Symbol.ID && c.File == e.Symbol.File {
			points[c.Line] = "chamada resolvida neste contexto"
		}
	}
	var ns []int
	for n := range points {
		ns = append(ns, n)
	}
	sort.Ints(ns)
	var out []domain.EvidenceExcerpt
	for _, n := range ns {
		if len(out) >= 8 {
			break
		}
		lo, hi := n-1, n+3
		if lo < e.Symbol.Line {
			lo = e.Symbol.Line
		}
		if hi > e.Symbol.EndLine {
			hi = e.Symbol.EndLine
		}
		if len(out) > 0 && lo <= out[len(out)-1].LastLine {
			continue
		}
		out = append(out, domain.EvidenceExcerpt{Role: points[n], FirstLine: lo, LastLine: hi, Content: strings.Join(lines[lo-1:hi], "\n")})
	}
	return out
}
