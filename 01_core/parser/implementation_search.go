/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/implementation_search.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"path"
	"sort"
	"strings"
	"unicode"

	"bateia/01_core/domain"
)

// ImplementationSearchPathAllowed applies the same scope before disk reads and ranking.
func ImplementationSearchPathAllowed(name, query string) bool {
	name = path.Clean(strings.ReplaceAll(name, "\\", "/"))
	if name == "." || name == ".." || strings.HasPrefix(name, "../") || strings.HasPrefix(name, "/") {
		return false
	}
	for _, part := range strings.Split(strings.ToLower(name), "/") {
		switch part {
		case ".git", ".hg", ".svn", ".agents", ".codex", ".cache", ".idea", ".vscode", "node_modules", "vendor", "dist", "build", "target", "tmp", "temp", "generated", "testdata", "fixtures", "fixture", "oracle", "oracles", "gates", "gate":
			return false
		}
		// Protected corpus names may carry a step number or descriptive suffix.
		for _, t := range implementationWords(part) {
			if t == "oracle" || t == "oracles" || t == "gate" || t == "gates" {
				return false
			}
		}
	}
	if strings.HasSuffix(strings.ToLower(name), "_test.go") {
		for _, t := range implementationWords(query) {
			if t == "test" || t == "tests" || t == "teste" || t == "testes" || t == "testing" {
				return true
			}
		}
		return false
	}
	return true
}

// implementationWords separates acronym boundaries as well as CamelCase and snake_case.
func implementationWords(s string) []string {
	r := []rune(s)
	var b strings.Builder
	for i, c := range r {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			b.WriteByte(' ')
			continue
		}
		if i > 0 && unicode.IsUpper(c) && (unicode.IsLower(r[i-1]) || unicode.IsDigit(r[i-1]) || (unicode.IsUpper(r[i-1]) && i+1 < len(r) && unicode.IsLower(r[i+1]))) {
			b.WriteByte(' ')
		}
		b.WriteRune(unicode.ToLower(c))
	}
	return strings.Fields(b.String())
}

// This finite glossary is lexical PT/EN expansion, not a semantic capability detector.
var implementationGlossary = map[string]string{
	"buscar": "search", "busca": "search", "pesquisa": "search", "encontrar": "find",
	"fila": "queue", "filas": "queue", "receita": "recipe", "receitas": "recipe",
	"arquivo": "file", "arquivos": "file", "ler": "read", "leitura": "read",
	"escrever": "write", "gravar": "write", "salvar": "save", "validar": "validate",
	"sanitizar": "sanitize",
	"evidencia": "evidence", "evidência": "evidence", "evidencias": "evidence", "evidências": "evidence",
	"avaliacao": "assessment", "avaliação": "assessment", "processar": "process",
	"mensagem": "message", "mensagens": "message", "consumir": "consume",
	"publicar": "publish", "autenticacao": "authentication", "autenticação": "authentication",
	"armazenar": "store", "armazenamento": "store", "storage": "store", "persistir": "persist", "persistencia": "persist", "persistência": "persist",
	"carregar": "load", "carregamento": "load", "escrita": "write", "writing": "write", "reading": "read", "saving": "save",
	"disco": "disk", "memoria": "memory", "memória": "memory", "chave": "key", "chaves": "key", "valor": "value", "valores": "value",
	"autenticar": "auth", "authentication": "auth", "authenticate": "auth", "authorization": "auth", "autorizar": "auth", "autorização": "auth",
	"usuario": "user", "usuário": "user", "usuarios": "user", "usuários": "user", "senha": "password", "credenciais": "credentials",
	"verificar": "verify", "verificação": "verify", "validacao": "validate", "validação": "validate", "validation": "validate",
	"assinatura": "signature", "token": "token", "tokens": "token", "expiração": "expiry", "expirar": "expire",
	"enfileirar": "enqueue", "desenfileirar": "dequeue", "consumo": "consume", "consumer": "consume", "consuming": "consume", "consumption": "consume",
	"producao": "produce", "produção": "produce", "producer": "produce", "produzir": "produce", "entrega": "deliver", "entregar": "deliver", "delivery": "deliver",
	"processamento": "process", "processing": "process", "mensageria": "message", "messages": "message", "messaging": "message",
	"enviar": "send", "envio": "send", "receber": "receive", "recebimento": "receive", "repetir": "retry", "retries": "retry",
	"bloquear": "lock", "bloqueio": "lock", "concorrente": "concurrent", "concorrencia": "concurrency", "concorrência": "concurrency",
	"remover": "remove", "excluir": "delete", "deletion": "delete", "cancelar": "cancel", "cancellation": "cancel",
}

func implementationTerms(s string) map[string]bool {
	stop := " a an and are as at be by can could does how where which what when who have has using use used uses support supports system systems service services mecanismo mecanismos sistema sistemas implementar implementado implementada onde qual quais fazer faz feita feito pelo pela pelos pelas esse essa este esta ele ela su sua seu that this these those code com como da das de do dos e em for from function functions func implementation implement implements implementacao implementação in is it metodo método na nas no nos o of on or os para por que s se the to um uma with codigo código "
	terms := map[string]bool{}
	for _, t := range implementationWords(s) {
		if strings.Contains(stop, " "+t+" ") {
			continue
		}
		terms[t] = true
		if en, ok := implementationGlossary[t]; ok {
			terms[en] = true
			if canonical, ok := implementationGlossary[en]; ok {
				terms[canonical] = true
			}
		}
	}
	return terms
}

// RankGoImplementations records static declarations only; scores are not proof of behavior.
func RankGoImplementations(query string, files []domain.ExtractedFile, limit int) domain.ImplementationSearchReport {
	if limit <= 0 {
		limit = 24
	}
	report := domain.ImplementationSearchReport{Status: "complete", Method: "go_ast_query_v1", Reason: "ponderação lexical IDF por declarações AST elegíveis do recorte lido, normalização de comprimento e chamadas diretas estáticas v2 limitadas; não representa ecossistema ou comprova comportamento", Hits: []domain.ImplementationHit{}, Examined: []domain.ImplementationInput{}, Budget: domain.ImplementationSearchBudget{MaxResults: limit}}
	inputs := append([]domain.ExtractedFile(nil), files...)
	for i := range inputs {
		inputs[i].Path = path.Clean(strings.ReplaceAll(inputs[i].Path, "\\", "/"))
	}
	sort.Slice(inputs, func(i, j int) bool {
		if inputs[i].Path != inputs[j].Path {
			return inputs[i].Path < inputs[j].Path
		}
		return inputs[i].Content < inputs[j].Content
	})
	relations := newFunctionalRelations()
	report.QueryInterpretation = InterpretImplementationQuery(query)
	exact := map[string]bool{}
	requirementTerms := []map[string]bool{}
	actionFiltered := false
	for _, interpretation := range report.QueryInterpretation {
		_, filtered := implementationQueryActionAdmission(interpretation, nil)
		actionFiltered = actionFiltered || filtered
		required := map[string]bool{}
		for _, term := range interpretation.Terms {
			required[term] = true
		}
		for _, symbol := range interpretation.Symbols {
			exact[strings.ToLower(symbol)] = true
		}
		requirementTerms = append(requirementTerms, required)
	}
	const actionReason = "; filtro lexical de admissão por ação e contexto da consulta aplicado por linha e campo, sem comprovar comportamento"
	if actionFiltered {
		report.Reason += actionReason
	}
	coverage := map[string][]int{}
	hitKey := func(hit domain.ImplementationHit) string {
		return fmt.Sprintf("%s\x00%d\x00%s", hit.Path, hit.Line, hit.Symbol)
	}

	documents := []corpusDocument{}
	documentationHits := []domain.ImplementationHit{}
	seen := map[string]string{}
	for _, file := range inputs {
		if !ImplementationSearchPathAllowed(file.Path, query) || !strings.HasSuffix(file.Path, ".go") {
			continue
		}
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(file.Content)))
		if previous, ok := seen[file.Path]; ok {
			if previous != hash {
				report.Issues = append(report.Issues, "conteúdos conflitantes para "+file.Path)
				relations.blocked[path.Dir(file.Path)] = true
			}
			continue
		}
		seen[file.Path] = hash
		report.Examined = append(report.Examined, domain.ImplementationInput{Path: file.Path, SHA256: hash, Bytes: int64(len(file.Content))})
		report.FilesRead++
		report.BytesRead += int64(len(file.Content))
		set := token.NewFileSet()
		tree, err := goparser.ParseFile(set, file.Path, file.Content, goparser.ParseComments|goparser.AllErrors)
		if err != nil {
			relations.blocked[path.Dir(file.Path)] = true
			report.Issues = append(report.Issues, fmt.Sprintf("parse %s: %v", file.Path, err))
			continue
		}
		if ast.IsGenerated(tree) {
			continue
		}
		relations.collect(file.Path, hash, tree, set)
		add := func(name, kind string, node ast.Node, doc *ast.CommentGroup) {
			nameTerms := implementationTerms(name)
			start, end := set.Position(node.Pos()), set.Position(node.End())
			bodyTerms, functional := functionalBodyTerms(node)
			documents = append(documents, corpusDocument{
				hit:  domain.ImplementationHit{Retrieval: "lexical", Path: file.Path, Symbol: name, Kind: "definition", SHA256: hash, Reason: kind, Line: start.Line, EndLine: end.Line},
				name: nameTerms, body: bodyTerms, functional: functional, exact: exact[strings.ToLower(name)],
			})
		}
		for _, decl := range tree.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				kind := "function"
				if d.Recv != nil {
					kind = "method"
				}
				add(d.Name.Name, kind, d, d.Doc)
				if file.Size <= int64(len(file.Content)) {
					documentationTerms := map[string]bool{}
					if d.Doc != nil {
						words := implementationTerms(d.Doc.Text())
						for i, interpretation := range report.QueryInterpretation {
							admitted, _ := implementationQueryActionAdmission(interpretation, words)
							if admitted && interpretation.Mode != "symbol" {
								for term := range requirementTerms[i] {
									documentationTerms[term] = true
								}
							}
						}
					}
					if hit, ok := documentationHit(file.Path, hash, d, set, documentationTerms); ok {
						if actionFiltered {
							hit.Reason += actionReason
						}
						documentationHits = append(documentationHits, hit)
					}
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						doc := s.Doc
						if doc == nil {
							doc = d.Doc
						}
						add(s.Name.Name, "type", s, doc)
					case *ast.ValueSpec:
						doc := s.Doc
						if doc == nil {
							doc = d.Doc
						}
						for _, name := range s.Names {
							add(name.Name, d.Tok.String(), s, doc)
						}
					}
				}
			}
		}
	}
	stats := collectCorpusStatistics(documents)
	for _, doc := range documents {
		field := map[string]bool{}
		for term := range doc.name {
			field[term] = true
		}
		for term := range doc.body {
			field[term] = true
		}
		terms := map[string]bool{}
		admissible := make([]bool, len(report.QueryInterpretation))
		for i, interpretation := range report.QueryInterpretation {
			admissible[i], _ = implementationQueryActionAdmission(interpretation, field)
			if admissible[i] {
				for term := range requirementTerms[i] {
					terms[term] = true
				}
			}
		}
		score, matched, weights := stats.score(terms, doc)
		if doc.exact {
			score += 100000
			matched = append(matched, doc.hit.Symbol)
		}
		if score == 0 {
			continue
		}
		sort.Strings(matched)
		compact := matched[:0]
		for _, term := range matched {
			if len(compact) == 0 || compact[len(compact)-1] != term {
				compact = append(compact, term)
			}
		}
		hit := doc.hit
		hit.Score, hit.MatchedTerms, hit.TermWeights = score, compact, weights
		hit.Reason += ": " + stats.reason(doc)
		if actionFiltered {
			hit.Reason += actionReason
		}
		if doc.exact {
			hit.Reason += "; identificador exato em declaração AST"
		}
		perRequirement := make([]int, len(requirementTerms))
		for i, required := range requirementTerms {
			if admissible[i] {
				perRequirement[i], _, _ = stats.score(required, doc)
			}
		}
		coverage[hitKey(hit)] = perRequirement
		report.Hits = append(report.Hits, hit)
	}

	sort.Slice(report.Hits, func(i, j int) bool {
		a, b := report.Hits[i], report.Hits[j]
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

	// Exact declarations remain globally strongest; one best hit per canonical
	// requirement follows, then the remaining global score order.
	if len(requirementTerms) > 1 {
		ordered := make([]domain.ImplementationHit, 0, len(report.Hits))
		used := make([]bool, len(report.Hits))
		for i, hit := range report.Hits {
			if hit.Score >= 100000 {
				ordered = append(ordered, hit)
				used[i] = true
			}
		}
		for line := range requirementTerms {
			best, bestScore := -1, 0
			for i, hit := range report.Hits {
				if used[i] {
					continue
				}
				score := coverage[hitKey(hit)][line]
				if score > bestScore {
					best, bestScore = i, score
				}
			}
			if best >= 0 {
				ordered = append(ordered, report.Hits[best])
				used[best] = true
			}
		}
		for i, hit := range report.Hits {
			if !used[i] {
				ordered = append(ordered, hit)
			}
		}
		report.Hits = ordered
	}
	documentationHits, documentationReserve := documentationBudget(report.Hits, documentationHits, limit)
	codeLimit := limit - documentationReserve
	if documentationReserve > 0 {
		report.Reason += "; reserva documental limitada para indícios do autor, sem comprovação de comportamento"
	}
	report.Hits = relations.expand(report.Hits, codeLimit)
	if len(report.Hits) > codeLimit {
		report.Hits = report.Hits[:codeLimit]
		report.Limited = true
		report.Issues = append(report.Issues, "limite de resultados atingido")
	}
	appendDocumentationHits(&report, documentationHits, limit)
	sort.Strings(report.Issues)
	uniqueIssues := report.Issues[:0]
	for _, issue := range report.Issues {
		if len(uniqueIssues) == 0 || uniqueIssues[len(uniqueIssues)-1] != issue {
			uniqueIssues = append(uniqueIssues, issue)
		}
	}
	report.Issues = uniqueIssues
	if report.FilesRead == 0 {
		report.Status = "not_evaluated"
		report.Reason = "nenhum arquivo Go elegível no escopo fornecido"
	}
	if len(report.Issues) > 0 {
		report.Status = "partial"
	}
	return report
}
