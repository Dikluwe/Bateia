/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/functional_retrieval.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"sort"

	"bateia/01_core/domain"
)

func functionalBodyTerms(node ast.Node) (map[string]bool, bool) {
	terms := map[string]bool{}
	fn, ok := node.(*ast.FuncDecl)
	if !ok || fn.Body == nil {
		return terms, false
	}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if _, closure := n.(*ast.FuncLit); closure {
			return false
		}
		if id, ok := n.(*ast.Ident); ok {
			for term := range implementationTerms(id.Name) {
				terms[term] = true
			}
		}
		return true
	})
	return terms, true
}

func functionalLexicalScore(required, name, body map[string]bool, functional bool) (int, []string) {
	score := 0
	matched := []string{}
	for term := range required {
		if body[term] {
			score += 90
		} else if name[term] {
			score += 60
		} else {
			continue
		}
		matched = append(matched, term)
	}
	if score > 0 {
		if score > 900 {
			score = 900
		}
		if functional {
			score += 1000
		}
	}
	return score, matched
}

type functionalCall struct {
	name string
	line int
}
type functionalFunction struct {
	hit   domain.ImplementationHit
	scope string
	calls []functionalCall
}
type functionalRelations struct {
	functions    map[string]functionalFunction
	targets      map[string][]string
	declarations map[string]int
	blocked      map[string]bool
}

func newFunctionalRelations() *functionalRelations {
	return &functionalRelations{functions: map[string]functionalFunction{}, targets: map[string][]string{}, declarations: map[string]int{}, blocked: map[string]bool{}}
}
func functionalHitKey(hit domain.ImplementationHit) string {
	return fmt.Sprintf("%s\x00%d\x00%s", hit.Path, hit.Line, hit.Symbol)
}
func (r *functionalRelations) collect(filePath, hash string, tree *ast.File, set *token.FileSet) {
	scope := path.Dir(filePath) + "\x00" + tree.Name.Name + "\x00"
	for _, decl := range tree.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					r.declarations[scope+s.Name.Name]++
				case *ast.ValueSpec:
					for _, n := range s.Names {
						r.declarations[scope+n.Name]++
					}
				}
			}
		case *ast.FuncDecl:
			if d.Recv == nil {
				r.declarations[scope+d.Name.Name]++
			}
			if d.Body == nil {
				continue
			}
			hit := domain.ImplementationHit{Path: filePath, Symbol: d.Name.Name, Kind: "definition", SHA256: hash, Line: set.Position(d.Pos()).Line, EndLine: set.Position(d.End()).Line, Retrieval: "direct_call", Reason: "function: chamada direta AST a destino único no mesmo diretório/pacote; um salto estático, sem prova de execução", MatchedTerms: []string{}}
			key := functionalHitKey(hit)
			fn := functionalFunction{hit: hit, scope: scope}
			ast.Inspect(d.Body, func(n ast.Node) bool {
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
				fn.calls = append(fn.calls, functionalCall{name: id.Name, line: set.Position(call.Pos()).Line})
				return true
			})
			r.functions[key] = fn
			if d.Recv == nil {
				r.targets[scope+d.Name.Name] = append(r.targets[scope+d.Name.Name], key)
			}
		}
	}
}

// expand never walks the newly discovered targets. Only retained lexical callers
// create references, so extraction can retain both ends within its hit budget.
func (r *functionalRelations) expand(lexical []domain.ImplementationHit, limit int) []domain.ImplementationHit {
	capacity := limit / 3
	if capacity == 0 || len(lexical) == 0 {
		return lexical
	}
	seeds := len(lexical)
	if seeds > limit-capacity {
		seeds = limit - capacity
	}
	// Exact identifiers are never displaced by expansion.
	for seeds < len(lexical) && seeds < limit && lexical[seeds].Score >= 100000 {
		seeds++
	}
	if capacity > limit-seeds {
		capacity = limit - seeds
	}
	refs := map[string][]domain.ImplementationReference{}
	for _, hit := range lexical[:seeds] {
		fn, ok := r.functions[functionalHitKey(hit)]
		if !ok || r.blocked[path.Dir(hit.Path)] {
			continue
		}
		for _, call := range fn.calls {
			targetName := fn.scope + call.name
			if r.declarations[targetName] != 1 || len(r.targets[targetName]) != 1 {
				continue
			}
			target := r.targets[targetName][0]
			refs[target] = append(refs[target], domain.ImplementationReference{Path: hit.Path, Symbol: hit.Symbol, Line: call.line, SHA256: hit.SHA256})
		}
	}
	lexicalByKey := map[string]domain.ImplementationHit{}
	seedKeys := map[string]bool{}
	for i := range lexical {
		key := functionalHitKey(lexical[i])
		lexical[i].RelatedTo = canonicalFunctionalReferences(refs[key])
		lexicalByKey[key] = lexical[i]
		if i < seeds {
			seedKeys[key] = true
		}
	}
	extraKeys := []string{}
	for key := range refs {
		if !seedKeys[key] {
			extraKeys = append(extraKeys, key)
		}
	}
	sort.Strings(extraKeys)
	if len(extraKeys) > capacity {
		extraKeys = extraKeys[:capacity]
	}
	if len(extraKeys) == 0 {
		return lexical
	}
	promoted := map[string]bool{}
	for _, key := range extraKeys {
		promoted[key] = true
	}
	retained := limit - len(extraKeys)
	out := make([]domain.ImplementationHit, 0, len(lexical)+len(extraKeys))
	tail := []domain.ImplementationHit{}
	for _, hit := range lexical {
		if promoted[functionalHitKey(hit)] {
			continue
		}
		if len(out) < retained {
			out = append(out, hit)
		} else {
			tail = append(tail, hit)
		}
	}
	for _, key := range extraKeys {
		hit, alreadyLexical := lexicalByKey[key]
		if !alreadyLexical {
			hit = r.functions[key].hit
		}
		hit.RelatedTo = canonicalFunctionalReferences(refs[key])
		out = append(out, hit)
	}
	// Preserve the unselected tail for the caller's truncation reporting.
	return append(out, tail...)
}
func canonicalFunctionalReferences(refs []domain.ImplementationReference) []domain.ImplementationReference {
	sort.Slice(refs, func(i, j int) bool {
		a, b := refs[i], refs[j]
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
	out := refs[:0]
	for _, ref := range refs {
		if len(out) == 0 || out[len(out)-1] != ref {
			out = append(out, ref)
		}
	}
	return out
}
