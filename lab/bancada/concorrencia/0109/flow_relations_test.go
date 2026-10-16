/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/flow_relations_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"reflect"
	"strings"
	"testing"
)

const flowFixture = `package p
 type Gate chan struct{}
 type Runner struct{}
 func(r *Runner) Start(f func()){go f()}
 type P struct{tokens Gate;other Gate;runner Runner}
 func(p *P) Submit(f func()){
  select{case p.tokens<-struct{}{}:p.runner.Start(func(){p.work(f)})}
 }
 func(p *P) work(f func()){defer p.tokens.release();f()}
 func(l Gate) release(){if l!=nil{<-l}}
`

func relationProfile(r Report) map[string]int {
	p := map[string]int{}
	for _, e := range r.Relations {
		p[e.Kind]++
	}
	return p
}
func TestFlowRelationsContrasts(t *testing.T) {
	baseline := map[string]int{"STATIC_CALL_TARGET": 3, "FUNCTION_CONTAINS_GO": 1, "CALL_ARGUMENT_FUNCTION": 1, "SEND_CASE_CONTAINS_GO_CALL": 1, "DEFERRED_RECEIVER_RECEIVE": 1, "SAME_DECLARED_CHANNEL_FIELD": 1}
	for _, c := range []struct {
		name, code string
		missing    []string
	}{
		{"baseline", flowFixture, nil},
		{"renamed", strings.NewReplacer("Gate", "Permit", "Runner", "Launcher", "Start", "Launch", "tokens", "permits", "work", "execute", "release", "giveBack").Replace(flowFixture), nil},
		{"no-go", strings.Replace(flowFixture, "go f()", "f()", 1), []string{"FUNCTION_CONTAINS_GO", "SEND_CASE_CONTAINS_GO_CALL"}},
		{"other-receive", strings.Replace(flowFixture, "<-l", "var x Gate;<-x", 1), []string{"DEFERRED_RECEIVER_RECEIVE", "SAME_DECLARED_CHANNEL_FIELD"}},
		{"other-field", strings.Replace(flowFixture, "case p.tokens<-", "case p.other<-", 1), []string{"SAME_DECLARED_CHANNEL_FIELD"}},
		{"other-case", strings.Replace(flowFixture, "struct{}{}:p.runner", "struct{}{}:return;default:p.runner", 1), []string{"SEND_CASE_CONTAINS_GO_CALL"}},
		{"nested-launch", strings.Replace(flowFixture, "struct{}{}:p.runner.Start(func(){p.work(f)})", "struct{}{}:func(){p.runner.Start(func(){p.work(f)})}()", 1), []string{"SEND_CASE_CONTAINS_GO_CALL"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			r, err := Observe(request(c.code))
			if err != nil {
				t.Fatal(err)
			}
			expected := map[string]int{}
			for k, v := range baseline {
				expected[k] = v
			}
			for _, k := range c.missing {
				delete(expected, k)
			}
			if !reflect.DeepEqual(relationProfile(r), expected) {
				t.Fatalf("relations %+v expected %+v", relationProfile(r), expected)
			}
			if !hasGap(r, "WORKER_BOUND_UNRESOLVED") || !hasGap(r, "FLOW_PATH_FEASIBILITY_UNRESOLVED") {
				t.Fatal("relations promoted to bound/path proof")
			}
			for _, e := range r.Relations {
				if e.Authority != "STATIC_TYPES_AND_SYNTAX" {
					t.Fatal("excessive authority")
				}
				for _, loc := range []Location{e.From, e.To} {
					if loc.ContentSHA256 != Digest(c.code) || loc.Offset < 0 || loc.EndOffset > len(c.code) || loc.EndOffset <= loc.Offset {
						t.Fatal("unbound relation")
					}
				}
			}
		})
	}
}
func TestFlowConditionsAndCallModes(t *testing.T) {
	conditional := strings.Replace(flowFixture, "defer p.tokens.release()", "if f!=nil {defer p.tokens.release()}", 1)
	r, err := Observe(request(conditional))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range r.Relations {
		if e.Kind == "DEFERRED_RECEIVER_RECEIVE" {
			found = true
			if !reflect.DeepEqual(e.Conditions, []string{"if(f!=nil)"}) || !reflect.DeepEqual(e.TargetConditions, []string{"if(l!=nil)"}) {
				t.Fatal("conditional defer lost")
			}
		}
	}
	if !found {
		t.Fatal("defer relationship missing")
	}
	deferred := strings.Replace(flowFixture, "struct{}{}:p.runner.Start", "struct{}{}:defer p.runner.Start", 1)
	r, err = Observe(request(deferred))
	if err != nil {
		t.Fatal(err)
	}
	found = false
	for _, e := range r.Relations {
		if e.Kind == "SEND_CASE_CONTAINS_GO_CALL" {
			found = true
			if e.CallMode != "DEFERRED" {
				t.Fatal("deferred launch called direct")
			}
		}
	}
	if !found {
		t.Fatal("missing lexical deferred call")
	}
	unreachable := strings.Replace(flowFixture, "struct{}{}:p.runner.Start", "struct{}{}:return;p.runner.Start", 1)
	r, err = Observe(request(unreachable))
	if err != nil {
		t.Fatal(err)
	}
	if relationProfile(r)["SEND_CASE_CONTAINS_GO_CALL"] != 1 || !hasGap(r, "FLOW_PATH_FEASIBILITY_UNRESOLVED") {
		t.Fatal("containment mistaken for reachability")
	}
}
func TestFlowIndirectCallsErrorsAndBound(t *testing.T) {
	code := strings.Replace(flowFixture, "struct{}{}:p.runner.Start(func()", "struct{}{}:launch:=p.runner.Start;launch(func()", 1)
	r, err := Observe(request(code))
	if err != nil {
		t.Fatal(err)
	}
	if relationProfile(r)["STATIC_CALL_TARGET"] != 2 || relationProfile(r)["SEND_CASE_CONTAINS_GO_CALL"] != 0 {
		t.Fatal("function variable resolved as direct method")
	}
	code = flowFixture + "\nvar bad=missing\n"
	r, err = Observe(request(code))
	if err != nil || len(r.Relations) != 0 || !hasGap(r, "SUPPLIED_PACKAGE_TYPECHECK_FAILED") {
		t.Fatal("invalid package credited")
	}
	code = "package p;func f(){};func g(){" + strings.Repeat("f();", 10001) + "}"
	if _, err = Observe(request(code)); err == nil {
		t.Fatal("relation overflow not rejected")
	}
}

func TestFlowInterfaceAndOrder(t *testing.T) {
	code := `package p; type I interface{Start(func())};type R struct{};func(r R)Start(f func()){go f()};func submit(i I,c chan struct{}){select{case c<-struct{}{}:i.Start(func(){})}}`
	r, err := Observe(request(code))
	if err != nil || relationProfile(r)["STATIC_CALL_TARGET"] != 0 || relationProfile(r)["SEND_CASE_CONTAINS_GO_CALL"] != 0 {
		t.Fatal("interface implementation inferred")
	}
	split := strings.Split(flowFixture, "func(p *P) work")
	input := request(split[0])
	second := "package p;func(p *P) work" + split[1]
	input.Files = append(input.Files, SourceFile{Path: "other.go", Content: second, SHA256: Digest(second)})
	r, err = Observe(input)
	if err != nil {
		t.Fatal(err)
	}
	input.Files[0], input.Files[1] = input.Files[1], input.Files[0]
	reordered, err := Observe(input)
	if err != nil || !reflect.DeepEqual(r, reordered) {
		t.Fatal("relation order changed across source reordering")
	}
}
