/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/channel_identity.md
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

type ChannelIdentityEvent struct {
	Kind     string   `json:"kind"`
	Location Location `json:"location"`
}
type ChannelIdentityAudit struct {
	Reservation     Location               `json:"reservation"`
	DeferredReceive Location               `json:"deferred_receive"`
	BaseLink        string                 `json:"base_link"`
	State           string                 `json:"state"`
	Authority       string                 `json:"authority"`
	Scope           string                 `json:"scope"`
	Events          []ChannelIdentityEvent `json:"events"`
	Limitations     []string               `json:"limitations"`
}
type channelPair struct {
	field                  *types.Var
	send, deferred         Location
	sendBase, deferredBase types.Object
}

func unparen(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}
func channelBase(e ast.Expr, info *types.Info) types.Object {
	s, ok := unparen(e).(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	id, ok := unparen(s.X).(*ast.Ident)
	if !ok {
		return nil
	}
	return info.Uses[id]
}
func auditChannelIdentity(functions []*flowFunction, pairs []channelPair) ([]ChannelIdentityAudit, error) {
	if len(pairs) == 0 {
		return []ChannelIdentityAudit{}, nil
	}
	candidates := map[*types.Var]bool{}
	bases := map[*types.Var]map[types.Object]bool{}
	for _, p := range pairs {
		candidates[p.field] = true
		if bases[p.field] == nil {
			bases[p.field] = map[types.Object]bool{}
		}
		if p.sendBase != nil {
			bases[p.field][p.sendBase] = true
		}
		if p.deferredBase != nil {
			bases[p.field][p.deferredBase] = true
		}
	}
	events := map[*types.Var][]ChannelIdentityEvent{}
	seen := map[string]bool{}
	total := 0
	overflow := false
	add := func(field *types.Var, kind string, w walker, n ast.Node) {
		if !candidates[field] || overflow {
			return
		}
		e := ChannelIdentityEvent{kind, w.location(n)}
		encoded, _ := json.Marshal(e)
		key := fmt.Sprintf("%p:%s", field, encoded)
		if !seen[key] {
			seen[key] = true
			events[field] = append(events[field], e)
			total++
			if total > 10000 {
				overflow = true
			}
		}
	}
	aggregateFields := func(t types.Type) []*types.Var {
		for i := 0; t != nil && i < 8; i++ {
			if ptr, ok := t.Underlying().(*types.Pointer); ok {
				t = ptr.Elem()
				continue
			}
			break
		}
		found := []*types.Var{}
		if t != nil {
			if s, ok := t.Underlying().(*types.Struct); ok {
				for i := 0; i < s.NumFields(); i++ {
					if candidates[s.Field(i)] {
						found = append(found, s.Field(i))
					}
				}
			}
		}
		return found
	}
	for _, fn := range functions {
		w := fn.w
		info := w.typing.info
		ast.Inspect(fn.decl.Body, func(n ast.Node) bool {
			if overflow {
				return false
			}
			switch x := n.(type) {
			case *ast.AssignStmt:
				for _, lhs := range x.Lhs {
					add(channelField(unparen(lhs), info), "FIELD_WRITE", w, lhs)
					for _, field := range aggregateFields(info.Types[lhs].Type) {
						add(field, "AGGREGATE_OR_BASE_WRITE", w, lhs)
					}
				}
			case *ast.UnaryExpr:
				if x.Op == token.AND {
					add(channelField(unparen(x.X), info), "FIELD_ADDRESS_EXPOSED", w, x)
				}
			case *ast.CallExpr:
				for _, arg := range x.Args {
					for _, field := range aggregateFields(info.Types[arg].Type) {
						add(field, "RECEIVER_ARGUMENT_EXPOSURE", w, arg)
					}
				}
			case *ast.FuncLit:
				ast.Inspect(x.Body, func(child ast.Node) bool {
					if id, ok := child.(*ast.Ident); ok {
						obj := info.Uses[id]
						if obj != nil {
							for field, objects := range bases {
								if objects[obj] {
									add(field, "RECEIVER_CAPTURE", w, x)
								}
							}
						}
					}
					return !overflow
				})
			}
			return !overflow
		})
	}
	if overflow {
		return nil, fmt.Errorf("channel audit limit exceeded")
	}
	result := []ChannelIdentityAudit{}
	facts := 0
	for _, p := range pairs {
		link := "DISTINCT_BASE_OBJECTS"
		if p.sendBase == nil || p.deferredBase == nil {
			link = "UNRESOLVED_BASE_EXPRESSION"
		} else if p.sendBase == p.deferredBase {
			link = "SAME_BASE_OBJECT"
		}
		es := append([]ChannelIdentityEvent{}, events[p.field]...)
		sort.Slice(es, func(i, j int) bool {
			a, _ := json.Marshal(es[i])
			b, _ := json.Marshal(es[j])
			return string(a) < string(b)
		})
		facts += 1 + len(es)
		if facts > 10000 {
			return nil, fmt.Errorf("channel audit limit exceeded")
		}
		result = append(result, ChannelIdentityAudit{p.send, p.deferred, link, "INSTANCE_IDENTITY_UNRESOLVED", "DECLARATION_SCOPED_AUDIT", "SUPPLIED_FUNCTION_BODIES_ALL_INSTANCES_OF_FIELD_DECLARATION", es, []string{"Events may involve other instances or unreachable code; no defect is established.", "Aliases, external mutation, unsafe/reflection and concurrent stability are unresolved.", "Same base object does not certify the same channel value; distinct bases do not certify different instances."}})
	}
	sort.Slice(result, func(i, j int) bool {
		a, _ := json.Marshal(result[i])
		b, _ := json.Marshal(result[j])
		return string(a) < string(b)
	})
	return result, nil
}
