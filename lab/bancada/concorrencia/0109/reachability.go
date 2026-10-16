/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/reachability.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/types"
	"sort"
)

type CallSiteReachability struct {
	Site       Location  `json:"site"`
	Scope      Location  `json:"scope"`
	State      string    `json:"state"`
	BlockedBy  *Location `json:"blocked_by,omitempty"`
	CallMode   string    `json:"call_mode"`
	Conditions []string  `json:"conditions"`
	Authority  string    `json:"authority"`
}
type DeferOrder struct {
	Kind       string   `json:"kind"`
	Defer      Location `json:"defer"`
	Site       Location `json:"site"`
	Scope      Location `json:"scope"`
	Conditions []string `json:"conditions"`
	Authority  string   `json:"authority"`
}
type siteState struct{ blocker ast.Node }
type prefixDefer struct {
	index int
	node  *ast.DeferStmt
}
type pathAnalyzer struct {
	prefix    []prefixDefer
	w         walker
	body      *ast.BlockStmt
	unknown   bool
	states    map[*ast.CallExpr]siteState
	rootIndex map[*ast.CallExpr]int
	root      int
}

func deriveReachability(functions []*flowFunction, calls []flowCall, r *Report) error {
	scopes := map[ast.Node]*pathAnalyzer{}
	for _, fn := range functions {
		addScope := func(anchor ast.Node, body *ast.BlockStmt) {
			p := &pathAnalyzer{w: fn.w, body: body, states: map[*ast.CallExpr]siteState{}, rootIndex: map[*ast.CallExpr]int{}}
			ast.Inspect(body, func(n ast.Node) bool {
				if _, ok := n.(*ast.FuncLit); ok {
					return false
				}
				switch n.(type) {
				case *ast.BranchStmt, *ast.LabeledStmt:
					p.unknown = true
				}
				return true
			})
			alive := true
			var blocker ast.Node
			for i, s := range body.List {
				if d, ok := s.(*ast.DeferStmt); ok {
					p.prefix = append(p.prefix, prefixDefer{i, d})
				}
				p.root = i
				if p.unknown {
					p.mark(s, nil)
				} else {
					alive, blocker = p.statement(s, alive, blocker)
				}
			}
			scopes[anchor] = p
		}
		addScope(fn.decl, fn.decl.Body)
		ast.Inspect(fn.decl.Body, func(n ast.Node) bool {
			if f, ok := n.(*ast.FuncLit); ok {
				addScope(f, f.Body)
			}
			return true
		})
	}
	r.CallSiteReachability = []CallSiteReachability{}
	r.DeferOrder = []DeferOrder{}
	for _, call := range calls {
		p := scopes[call.anchor]
		if p == nil {
			continue
		}
		state, ok := p.states[call.call]
		if !ok {
			return fmt.Errorf("call site missing from control projection")
		}
		fact := CallSiteReachability{Site: p.w.location(call.call), Scope: p.w.location(call.anchor), State: "POSSIBLY_REACHABLE", CallMode: call.mode, Conditions: append([]string{}, call.conditions...), Authority: "LIMITED_STRUCTURED_CONTROL"}
		if p.unknown {
			fact.State = "UNSUPPORTED_CONTROL_FLOW"
		} else if state.blocker != nil {
			fact.State = "UNREACHABLE_AFTER_TERMINATOR"
			loc := p.w.location(state.blocker)
			fact.BlockedBy = &loc
		}
		r.CallSiteReachability = append(r.CallSiteReachability, fact)
		if !p.unknown && state.blocker == nil && call.mode == "DIRECT" {
			for _, prior := range p.prefix {
				if prior.index >= p.rootIndex[call.call] {
					break
				}
				d := prior.node
				{
					r.DeferOrder = append(r.DeferOrder, DeferOrder{Kind: "PREFIX_DEFER_BEFORE_CALL_SITE", Defer: p.w.location(d), Site: fact.Site, Scope: fact.Scope, Conditions: fact.Conditions, Authority: "LIMITED_STRUCTURED_CONTROL"})
					if len(r.CallSiteReachability)+len(r.DeferOrder) > 10000 {
						return fmt.Errorf("control fact limit exceeded")
					}
				}
			}
		}
		if len(r.CallSiteReachability)+len(r.DeferOrder) > 10000 {
			return fmt.Errorf("control fact limit exceeded")
		}
	}
	sort.Slice(r.CallSiteReachability, func(i, j int) bool {
		a, _ := json.Marshal(r.CallSiteReachability[i])
		b, _ := json.Marshal(r.CallSiteReachability[j])
		return string(a) < string(b)
	})
	sort.Slice(r.DeferOrder, func(i, j int) bool {
		a, _ := json.Marshal(r.DeferOrder[i])
		b, _ := json.Marshal(r.DeferOrder[j])
		return string(a) < string(b)
	})
	return nil
}
func (p *pathAnalyzer) mark(n ast.Node, blocker ast.Node) {
	if n == nil {
		return
	}
	ast.Inspect(n, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncLit); ok {
			return false
		}
		if c, ok := node.(*ast.CallExpr); ok {
			p.states[c] = siteState{blocker}
			p.rootIndex[c] = p.root
		}
		return true
	})
}
func (p *pathAnalyzer) sequence(list []ast.Stmt, alive bool, blocker ast.Node) (bool, ast.Node) {
	for _, s := range list {
		alive, blocker = p.statement(s, alive, blocker)
	}
	return alive, blocker
}
func (p *pathAnalyzer) statement(s ast.Stmt, alive bool, blocker ast.Node) (bool, ast.Node) {
	if s == nil {
		return alive, blocker
	}
	if !alive {
		p.mark(s, blocker)
		return false, blocker
	}
	switch x := s.(type) {
	case *ast.ReturnStmt:
		p.mark(x, nil)
		return false, x
	case *ast.BlockStmt:
		return p.sequence(x.List, true, nil)
	case *ast.IfStmt:
		p.mark(x.Init, nil)
		p.mark(x.Cond, nil)
		thenAlive, _ := p.statement(x.Body, true, nil)
		elseAlive := true
		if x.Else != nil {
			elseAlive, _ = p.statement(x.Else, true, nil)
		}
		if !thenAlive && !elseAlive {
			return false, x
		}
		return true, nil
	case *ast.ForStmt:
		p.mark(x.Init, nil)
		p.mark(x.Cond, nil)
		p.mark(x.Post, nil)
		p.statement(x.Body, true, nil)
		return true, nil
	case *ast.RangeStmt:
		p.mark(x.X, nil)
		p.statement(x.Body, true, nil)
		return true, nil
	case *ast.SelectStmt:
		for _, node := range x.Body.List {
			clause := node.(*ast.CommClause)
			p.mark(clause.Comm, nil)
			p.sequence(clause.Body, true, nil)
		}
		return true, nil
	case *ast.SwitchStmt:
		p.mark(x.Init, nil)
		p.mark(x.Tag, nil)
		for _, node := range x.Body.List {
			clause := node.(*ast.CaseClause)
			for _, e := range clause.List {
				p.mark(e, nil)
			}
			p.sequence(clause.Body, true, nil)
		}
		return true, nil
	case *ast.TypeSwitchStmt:
		p.mark(x.Init, nil)
		p.mark(x.Assign, nil)
		for _, node := range x.Body.List {
			clause := node.(*ast.CaseClause)
			for _, e := range clause.List {
				p.mark(e, nil)
			}
			p.sequence(clause.Body, true, nil)
		}
		return true, nil
	case *ast.ExprStmt:
		p.mark(x, nil)
		e := x.X
		for {
			par, ok := e.(*ast.ParenExpr)
			if !ok {
				break
			}
			e = par.X
		}
		if call, ok := e.(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok {
				if obj, ok := p.w.typing.info.Uses[id].(*types.Builtin); ok && obj.Name() == "panic" {
					return false, x
				}
			}
		}
		return true, nil
	default:
		p.mark(s, nil)
		return true, nil
	}
}
