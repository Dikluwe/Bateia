/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/flow_relations.md
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
	"path"
	"sort"
	"strconv"
)

type Relation struct {
	ParameterDeclaration *Location `json:"parameter_declaration,omitempty"`
	ExecutionScope       *Location `json:"execution_scope,omitempty"`
	Kind                 string    `json:"kind"`
	From                 Location  `json:"from"`
	To                   Location  `json:"to"`
	Conditions           []string  `json:"conditions"`
	TargetConditions     []string  `json:"target_conditions,omitempty"`
	CallMode             string    `json:"call_mode,omitempty"`
	Authority            string    `json:"authority"`
}
type flowFunction struct {
	decl     *ast.FuncDecl
	obj      *types.Func
	w        walker
	hasGo    bool
	receives []flowReceive
}
type flowReceive struct {
	node       *ast.UnaryExpr
	conditions []string
}
type flowFieldUse struct {
	obj        *types.Var
	location   Location
	conditions []string
}
type flowCall struct {
	fn         *flowFunction
	call       *ast.CallExpr
	anchor     ast.Node
	conditions []string
	mode       string
}
type flowBuilder struct {
	fieldBases            map[Location]types.Object
	identityPairs         []channelPair
	calls                 []flowCall
	functions             map[*types.Func]*flowFunction
	relations             []Relation
	sends, deferredFields []flowFieldUse
	tooMany               bool
}

func deriveRelations(files []parsedSource, resolutions map[string]typeResolution, report *Report) ([]Relation, error) {
	b := flowBuilder{functions: map[*types.Func]*flowFunction{}, relations: []Relation{}, fieldBases: map[Location]types.Object{}}
	ordered := []*flowFunction{}
	for _, p := range files {
		typing := resolutions[path.Dir(p.source.Path)+"/"+p.tree.Name.Name]
		if !typing.complete {
			continue
		}
		w := walker{source: p.source, fs: p.fset, typing: typing, report: report}
		for _, d := range p.tree.Decls {
			f, ok := d.(*ast.FuncDecl)
			if !ok || f.Body == nil {
				continue
			}
			obj, ok := typing.info.Defs[f.Name].(*types.Func)
			if !ok {
				continue
			}
			fn := &flowFunction{decl: f, obj: obj, w: w}
			ast.Inspect(f.Body, func(n ast.Node) bool {
				if _, ok := n.(*ast.FuncLit); ok {
					return false
				}
				if _, ok := n.(*ast.GoStmt); ok {
					fn.hasGo = true
				}
				return true
			})
			b.functions[obj] = fn
			ordered = append(ordered, fn)
		}
	}
	// Index receives before traversing calls, so file order does not affect deferred relations.
	for _, fn := range ordered {
		b.visit(fn, fn.decl.Body, fn.decl, nil, nil, "DIRECT", true)
	}
	for _, fn := range ordered {
		b.visit(fn, fn.decl.Body, fn.decl, nil, nil, "DIRECT", false)
	}
	// Matching a declaration of a field is weaker than matching a channel instance.
	deferredByField := map[*types.Var][]flowFieldUse{}
	for _, deferred := range b.deferredFields {
		deferredByField[deferred.obj] = append(deferredByField[deferred.obj], deferred)
	}
	for _, send := range b.sends {
		for _, deferred := range deferredByField[send.obj] {
			b.identityPairs = append(b.identityPairs, channelPair{send.obj, send.location, deferred.location, b.fieldBases[send.location], b.fieldBases[deferred.location]})
			b.add(Relation{Kind: "SAME_DECLARED_CHANNEL_FIELD", From: send.location, To: deferred.location, Conditions: send.conditions, TargetConditions: deferred.conditions})
			if b.tooMany {
				return nil, fmt.Errorf("flow relation limit exceeded")
			}
		}
	}
	if b.tooMany {
		return nil, fmt.Errorf("flow relation limit exceeded")
	}
	callbacks, err := deriveCallbacks(ordered, b.calls, report, 10000-len(b.relations))
	if err != nil {
		return nil, err
	}
	report.CallbackRelations = callbacks
	identity, err := auditChannelIdentity(ordered, b.identityPairs)
	if err != nil {
		return nil, err
	}
	report.ChannelIdentity = identity
	guards, err := deriveConfigurationGuards(ordered, b.functions, b.identityPairs)
	if err != nil {
		return nil, err
	}
	report.ConfigurationGuards = guards
	if err := deriveReachability(ordered, b.calls, report); err != nil {
		return nil, err
	}
	sort.Slice(b.relations, func(i, j int) bool {
		a, _ := json.Marshal(b.relations[i])
		c, _ := json.Marshal(b.relations[j])
		return string(a) < string(c)
	})
	// Gaps are linked to each relation source; duplicate gap sites are suppressed.
	seen := map[string]bool{}
	for _, r := range b.relations {
		codes := []string{"FLOW_PATH_FEASIBILITY_UNRESOLVED"}
		if r.Kind == "CALL_ARGUMENT_FUNCTION" {
			codes = append(codes, "CALLBACK_EXECUTION_UNRESOLVED")
		}
		if r.Kind == "SAME_DECLARED_CHANNEL_FIELD" {
			codes = append(codes, "CHANNEL_INSTANCE_IDENTITY_UNRESOLVED")
		}
		for _, code := range codes {
			gap := Gap{Code: code, Location: r.From}
			key, _ := json.Marshal(gap)
			if !seen[string(key)] {
				seen[string(key)] = true
				report.Gaps = append(report.Gaps, gap)
			}
		}
	}
	return b.relations, nil
}
func (b *flowBuilder) add(r Relation) {
	if len(b.relations) >= 10000 {
		b.tooMany = true
		return
	}
	r.Authority = "STATIC_TYPES_AND_SYNTAX"
	r.Conditions = append([]string{}, r.Conditions...)
	r.TargetConditions = append([]string(nil), r.TargetConditions...)
	b.relations = append(b.relations, r)
}
func staticTarget(call *ast.CallExpr, info *types.Info) *types.Func {
	switch x := call.Fun.(type) {
	case *ast.Ident:
		f, _ := info.Uses[x].(*types.Func)
		return f
	case *ast.SelectorExpr:
		// Interface dispatch is not a static implementation target.
		if sel := info.Selections[x]; sel != nil {
			t := sel.Recv()
			if p, ok := t.(*types.Pointer); ok {
				t = p.Elem()
			}
			if _, ok := t.Underlying().(*types.Interface); ok {
				return nil
			}
		}
		f, _ := info.Uses[x.Sel].(*types.Func)
		return f
	}
	return nil
}
func channelField(e ast.Expr, info *types.Info) *types.Var {
	s, ok := e.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	selection := info.Selections[s]
	if selection == nil {
		return nil
	}
	v, ok := selection.Obj().(*types.Var)
	if !ok || !v.IsField() {
		return nil
	}
	if _, ok := v.Type().Underlying().(*types.Chan); !ok {
		return nil
	}
	return v
}
func (b *flowBuilder) visit(fn *flowFunction, n ast.Node, anchor ast.Node, cs []string, send *ast.SendStmt, mode string, indexReceives bool) {
	if n == nil || b.tooMany {
		return
	}
	w := fn.w
	switch x := n.(type) {
	case *ast.FuncLit:
		// A lexical closure is an independent execution scope.
		b.visit(fn, x.Body, x, nil, nil, "DIRECT", indexReceives)
		return
	case *ast.IfStmt:
		b.visit(fn, x.Init, anchor, cs, send, mode, indexReceives)
		b.visit(fn, x.Cond, anchor, cs, send, mode, indexReceives)
		b.visit(fn, x.Body, anchor, withCondition(cs, "if("+w.text(x.Cond)+")"), send, mode, indexReceives)
		b.visit(fn, x.Else, anchor, withCondition(cs, "else("+w.text(x.Cond)+")"), send, mode, indexReceives)
		return
	case *ast.ForStmt:
		b.visit(fn, x.Init, anchor, cs, send, mode, indexReceives)
		label := "for(;;)"
		if x.Cond != nil {
			label = "for(" + w.text(x.Cond) + ")"
		}
		next := withCondition(cs, label)
		b.visit(fn, x.Cond, anchor, next, send, mode, indexReceives)
		b.visit(fn, x.Post, anchor, next, send, mode, indexReceives)
		b.visit(fn, x.Body, anchor, next, send, mode, indexReceives)
		return
	case *ast.RangeStmt:
		b.visit(fn, x.X, anchor, cs, send, mode, indexReceives)
		b.visit(fn, x.Body, anchor, withCondition(cs, "range("+w.text(x.X)+")"), send, mode, indexReceives)
		return
	case *ast.CaseClause:
		for _, e := range x.List {
			b.visit(fn, e, anchor, cs, send, mode, indexReceives)
		}
		label := "case:default"
		if len(x.List) > 0 {
			label = "case(" + w.text(x.List[0]) + ")"
		}
		for _, s := range x.Body {
			b.visit(fn, s, anchor, withCondition(cs, label), send, mode, indexReceives)
		}
		return
	case *ast.SelectStmt:
		for i, node := range x.Body.List {
			clause := node.(*ast.CommClause)
			next := withCondition(cs, "select.case:"+strconv.Itoa(i))
			b.visit(fn, clause.Comm, anchor, next, nil, mode, indexReceives)
			reserved, _ := clause.Comm.(*ast.SendStmt)
			if reserved != nil {
				if t := w.typing.info.Types[reserved.Chan].Type; t == nil {
					reserved = nil
				} else if _, ok := t.Underlying().(*types.Chan); !ok {
					reserved = nil
				}
			}
			for _, s := range clause.Body {
				b.visit(fn, s, anchor, next, reserved, mode, indexReceives)
			}
		}
		return
	case *ast.GoStmt:
		if !indexReceives {
			b.add(Relation{Kind: "FUNCTION_CONTAINS_GO", From: w.location(anchor), To: w.location(x), Conditions: cs})
		}
		b.visit(fn, x.Call, anchor, cs, send, "GO", indexReceives)
		return
	case *ast.DeferStmt:
		if !indexReceives {
			if target := b.functions[staticTarget(x.Call, w.typing.info)]; target != nil {
				for _, receive := range target.receives {
					b.add(Relation{Kind: "DEFERRED_RECEIVER_RECEIVE", From: w.location(x), To: target.w.location(receive.node), Conditions: cs, TargetConditions: receive.conditions, CallMode: "DEFERRED"})
				}
				if len(target.receives) > 0 {
					if selector, ok := x.Call.Fun.(*ast.SelectorExpr); ok {
						if field := channelField(selector.X, w.typing.info); field != nil {
							b.deferredFields = append(b.deferredFields, flowFieldUse{field, w.location(x), cs})
							b.fieldBases[w.location(x)] = channelBase(selector.X, w.typing.info)
						}
					}
				}
			}
		}
		b.visit(fn, x.Call, anchor, cs, send, "DEFERRED", indexReceives)
		return
	case *ast.UnaryExpr:
		if indexReceives && anchor == fn.decl && x.Op == token.ARROW {
			if id, ok := x.X.(*ast.Ident); ok {
				sig := fn.obj.Type().(*types.Signature)
				if recv := sig.Recv(); recv != nil && w.typing.info.Uses[id] == recv {
					if _, ok := recv.Type().Underlying().(*types.Chan); ok {
						fn.receives = append(fn.receives, flowReceive{x, append([]string{}, cs...)})
					}
				}
			}
		}
	case *ast.SendStmt:
		if !indexReceives {
			if field := channelField(x.Chan, w.typing.info); field != nil {
				b.sends = append(b.sends, flowFieldUse{field, w.location(x), cs})
				b.fieldBases[w.location(x)] = channelBase(x.Chan, w.typing.info)
			}
		}
	case *ast.CallExpr:
		if !indexReceives {
			b.calls = append(b.calls, flowCall{fn, x, anchor, append([]string{}, cs...), mode})
			target := b.functions[staticTarget(x, w.typing.info)]
			if target != nil {
				b.add(Relation{Kind: "STATIC_CALL_TARGET", From: w.location(x), To: target.w.location(target.decl), Conditions: cs, CallMode: mode})
				if send != nil && target.hasGo {
					b.add(Relation{Kind: "SEND_CASE_CONTAINS_GO_CALL", From: w.location(send), To: w.location(x), Conditions: cs, CallMode: mode})
				}
			}
			for _, arg := range x.Args {
				if literal, ok := arg.(*ast.FuncLit); ok {
					b.add(Relation{Kind: "CALL_ARGUMENT_FUNCTION", From: w.location(x), To: w.location(literal), Conditions: cs, CallMode: mode})
				}
			}
		}
	}
	first := true
	ast.Inspect(n, func(child ast.Node) bool {
		if first {
			first = false
			return true
		}
		if child != nil {
			b.visit(fn, child, anchor, cs, send, "DIRECT", indexReceives)
		}
		return false
	})
}
