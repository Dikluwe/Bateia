/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/callback_relations.md
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

type callbackParameter struct {
	fn       *flowFunction
	id       *ast.Ident
	obj      *types.Var
	index    int
	modified bool
}

func deriveCallbacks(functions []*flowFunction, calls []flowCall, report *Report, limit int) ([]Relation, error) {
	params := map[*types.Var]*callbackParameter{}
	byFunction := map[*types.Func]map[int]*callbackParameter{}
	ordered := []*callbackParameter{}
	for _, fn := range functions {
		sig := fn.obj.Type().(*types.Signature)
		if sig.Variadic() {
			continue
		}
		for _, field := range fn.decl.Type.Params.List {
			for _, id := range field.Names {
				obj, ok := fn.w.typing.info.Defs[id].(*types.Var)
				if !ok {
					continue
				}
				if _, ok := obj.Type().Underlying().(*types.Signature); !ok {
					continue
				}
				for i := 0; i < sig.Params().Len(); i++ {
					if sig.Params().At(i) == obj {
						p := &callbackParameter{fn: fn, id: id, obj: obj, index: i}
						params[obj] = p
						ordered = append(ordered, p)
						if byFunction[fn.obj] == nil {
							byFunction[fn.obj] = map[int]*callbackParameter{}
						}
						byFunction[fn.obj][i] = p
					}
				}
			}
		}
	}
	// Conservative whole-body mutation check, including captured parameters in closures.
	for _, fn := range functions {
		mark := func(e ast.Expr) {
			for {
				p, ok := e.(*ast.ParenExpr)
				if !ok {
					break
				}
				e = p.X
			}
			if id, ok := e.(*ast.Ident); ok {
				if obj, ok := fn.w.typing.info.Uses[id].(*types.Var); ok {
					if p := params[obj]; p != nil {
						p.modified = true
					}
				}
			}
		}
		ast.Inspect(fn.decl.Body, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.AssignStmt:
				for _, lhs := range x.Lhs {
					mark(lhs)
				}
			case *ast.UnaryExpr:
				if x.Op == token.AND {
					mark(x.X)
				}
			}
			return true
		})
	}
	result := []Relation{}
	add := func(r Relation) error {
		if len(result) >= limit {
			return fmt.Errorf("combined flow relation limit exceeded")
		}
		r.Authority = "STATIC_TYPES_AND_SYNTAX"
		r.Conditions = append([]string{}, r.Conditions...)
		result = append(result, r)
		return nil
	}
	for _, site := range calls {
		w := site.fn.w
		info := w.typing.info
		scope := w.location(site.anchor)
		parameterAt := func(e ast.Expr) *callbackParameter {
			id, ok := e.(*ast.Ident)
			if !ok {
				return nil
			}
			obj, ok := info.Uses[id].(*types.Var)
			if !ok {
				return nil
			}
			return params[obj]
		}
		if p := parameterAt(site.call.Fun); p != nil && !p.modified {
			if err := add(Relation{Kind: "CALLBACK_PARAMETER_INVOCATION", From: p.fn.w.location(p.id), To: w.location(site.call), Conditions: site.conditions, CallMode: site.mode, ExecutionScope: &scope}); err != nil {
				return nil, err
			}
		}
		target := staticTarget(site.call, info)
		if target == nil || target.Type().(*types.Signature).Variadic() || site.call.Ellipsis.IsValid() || len(site.call.Args) != target.Type().(*types.Signature).Params().Len() {
			continue
		}
		// Method expressions include an explicit receiver argument; unsupported arity is skipped.
		for i, arg := range site.call.Args {
			destination := byFunction[target][i]
			if destination == nil || destination.modified {
				continue
			}
			r := Relation{From: w.location(arg), To: destination.fn.w.location(destination.id), Conditions: site.conditions, CallMode: site.mode, ExecutionScope: &scope}
			if _, ok := arg.(*ast.FuncLit); ok {
				r.Kind = "CALLBACK_ARGUMENT_BINDING"
			} else if p := parameterAt(arg); p != nil && !p.modified {
				r.Kind = "CALLBACK_PARAMETER_FORWARDING"
				declaration := p.fn.w.location(p.id)
				r.ParameterDeclaration = &declaration
			} else {
				continue
			}
			if err := add(r); err != nil {
				return nil, err
			}
		}
	}
	for _, p := range ordered {
		if p.modified {
			report.Gaps = append(report.Gaps, Gap{Code: "CALLBACK_VALUE_FLOW_UNRESOLVED", Location: p.fn.w.location(p.id)})
		}
	}
	seen := map[string]bool{}
	for _, edge := range result {
		for _, code := range []string{"CALLBACK_EXECUTION_UNRESOLVED", "FLOW_PATH_FEASIBILITY_UNRESOLVED"} {
			gap := Gap{Code: code, Location: edge.From}
			key, _ := json.Marshal(gap)
			if !seen[string(key)] {
				seen[string(key)] = true
				report.Gaps = append(report.Gaps, gap)
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
