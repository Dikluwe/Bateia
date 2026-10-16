/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/configuration_guard.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"sort"
)

type ConfigurationGuard struct {
	Write              Location   `json:"write"`
	GuardCall          Location   `json:"guard_call"`
	GuardDefinition    Location   `json:"guard_definition"`
	Predicate          Location   `json:"predicate"`
	PredicateText      string     `json:"predicate_text"`
	Panic              Location   `json:"panic"`
	CheckedFieldWrites []Location `json:"checked_field_writes"`
	State              string     `json:"state"`
	Authority          string     `json:"authority"`
	Atomicity          string     `json:"atomicity"`
	Limitations        []string   `json:"limitations"`
}
type guardPattern struct {
	fn          *flowFunction
	conditional *ast.IfStmt
	panic       *ast.ExprStmt
	field       *types.Var
}

func matchGuard(fn *flowFunction) *guardPattern {
	sig := fn.obj.Type().(*types.Signature)
	if sig.Recv() == nil || sig.Params().Len() != 0 || sig.Results().Len() != 0 || len(fn.decl.Body.List) != 1 {
		return nil
	}
	cond, ok := fn.decl.Body.List[0].(*ast.IfStmt)
	if !ok || cond.Init != nil || cond.Else != nil || len(cond.Body.List) != 1 {
		return nil
	}
	binary, ok := unparen(cond.Cond).(*ast.BinaryExpr)
	if !ok || binary.Op != token.NEQ {
		return nil
	}
	fieldExpr, nilExpr := unparen(binary.X), unparen(binary.Y)
	if id, ok := fieldExpr.(*ast.Ident); ok && fn.w.typing.info.Uses[id] == types.Universe.Lookup("nil") {
		fieldExpr, nilExpr = nilExpr, fieldExpr
	}
	id, ok := nilExpr.(*ast.Ident)
	if !ok || fn.w.typing.info.Uses[id] != types.Universe.Lookup("nil") {
		return nil
	}
	selector, ok := fieldExpr.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	base, ok := unparen(selector.X).(*ast.Ident)
	if !ok || fn.w.typing.info.Uses[base] != sig.Recv() {
		return nil
	}
	selection := fn.w.typing.info.Selections[selector]
	if selection == nil || len(selection.Index()) != 1 {
		return nil
	}
	field, ok := selection.Obj().(*types.Var)
	if !ok || !field.IsField() {
		return nil
	}
	statement, ok := cond.Body.List[0].(*ast.ExprStmt)
	if !ok {
		return nil
	}
	call, ok := unparen(statement.X).(*ast.CallExpr)
	if !ok {
		return nil
	}
	panicID, ok := unparen(call.Fun).(*ast.Ident)
	if !ok || fn.w.typing.info.Uses[panicID] != types.Universe.Lookup("panic") {
		return nil
	}
	return &guardPattern{fn, cond, statement, field}
}
func stableBaseSyntax(fn *flowFunction, base types.Object) bool {
	bad := false
	info := fn.w.typing.info
	same := func(e ast.Expr) bool {
		if id, ok := unparen(e).(*ast.Ident); ok {
			return info.Uses[id] == base
		}
		return false
	}
	ast.Inspect(fn.decl.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BranchStmt, *ast.LabeledStmt:
			bad = true
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				if same(lhs) {
					bad = true
				}
				if star, ok := unparen(lhs).(*ast.StarExpr); ok && same(star.X) {
					bad = true
				}
			}
		case *ast.UnaryExpr:
			if x.Op == token.AND && same(x.X) {
				bad = true
			}
		case *ast.FuncLit:
			ast.Inspect(x.Body, func(child ast.Node) bool {
				if id, ok := child.(*ast.Ident); ok && info.Uses[id] == base {
					bad = true
				}
				return true
			})
		}
		return !bad
	})
	return !bad
}
func deriveConfigurationGuards(functions []*flowFunction, index map[*types.Func]*flowFunction, pairs []channelPair) ([]ConfigurationGuard, error) {
	if len(pairs) == 0 {
		return []ConfigurationGuard{}, nil
	}
	candidates := map[*types.Var]bool{}
	for _, p := range pairs {
		candidates[p.field] = true
	}
	patterns := map[*types.Func]*guardPattern{}
	for _, fn := range functions {
		if p := matchGuard(fn); p != nil {
			patterns[fn.obj] = p
		}
	}
	sentinelWrites := map[*types.Var][]Location{}
	for _, fn := range functions {
		ast.Inspect(fn.decl.Body, func(n ast.Node) bool {
			if a, ok := n.(*ast.AssignStmt); ok {
				for _, lhs := range a.Lhs {
					selector, ok := unparen(lhs).(*ast.SelectorExpr)
					if !ok {
						return true
					}
					if sel := fn.w.typing.info.Selections[selector]; sel != nil {
						if field, ok := sel.Obj().(*types.Var); ok && field.IsField() {
							sentinelWrites[field] = append(sentinelWrites[field], fn.w.location(lhs))
						}
					}
				}
			}
			return true
		})
	}
	result := []ConfigurationGuard{}
	facts := 0
	for _, fn := range functions {
		for i, stmt := range fn.decl.Body.List {
			assignment, ok := stmt.(*ast.AssignStmt)
			if !ok {
				continue
			}
			for _, lhs := range assignment.Lhs {
				field := channelField(unparen(lhs), fn.w.typing.info)
				if !candidates[field] {
					continue
				}
				base := channelBase(lhs, fn.w.typing.info)
				if base == nil || !stableBaseSyntax(fn, base) {
					continue
				}
				for j := 0; j < i; j++ {
					expr, ok := fn.decl.Body.List[j].(*ast.ExprStmt)
					if !ok {
						continue
					}
					call, ok := unparen(expr.X).(*ast.CallExpr)
					if !ok || len(call.Args) != 0 {
						continue
					}
					pattern := patterns[staticTarget(call, fn.w.typing.info)]
					if pattern == nil {
						continue
					}
					selector, ok := unparen(call.Fun).(*ast.SelectorExpr)
					if !ok {
						continue
					}
					receiver, ok := unparen(selector.X).(*ast.Ident)
					if !ok || fn.w.typing.info.Uses[receiver] != base {
						continue
					}
					// Prior lexical guard is not a proof of atomicity or of the field value at the write.
					writes := append([]Location{}, sentinelWrites[pattern.field]...)
					facts += 1 + len(writes)
					if facts > 10000 {
						return nil, fmt.Errorf("configuration guard limit exceeded")
					}
					sort.Slice(writes, func(i, j int) bool {
						a, _ := json.Marshal(writes[i])
						b, _ := json.Marshal(writes[j])
						return string(a) < string(b)
					})
					result = append(result, ConfigurationGuard{fn.w.location(lhs), fn.w.location(call), pattern.fn.w.location(pattern.fn.decl), pattern.fn.w.location(pattern.conditional.Cond), pattern.fn.w.text(pattern.conditional.Cond), pattern.fn.w.location(pattern.panic), writes, "PRIOR_GUARD_SAME_BASE_OBJECT", "LIMITED_GUARD_PATTERN", "UNRESOLVED", []string{"Guard describes a prior check, not the checked field value at the write.", "Sentinel writes are declaration scoped; initialization of this instance is unresolved.", "Recovery, aliases, intervening effects, external mutation and concurrent atomicity are unresolved."}})
				}
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		a, _ := json.Marshal(result[i])
		b, _ := json.Marshal(result[j])
		return string(a) < string(b)
	})
	return result, nil
}
