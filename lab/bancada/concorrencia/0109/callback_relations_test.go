/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/callback_relations_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"reflect"
	"strings"
	"testing"
)

const callbackFixture = `package p
func Start(cb func()){go func(){Forward(cb)}()}
func Forward(task func()){task()}
func Entry(){Start(func(){})}
`

func callbackProfile(r Report) map[string]int {
	p := map[string]int{}
	for _, e := range r.CallbackRelations {
		p[e.Kind]++
	}
	return p
}
func TestCallbackContrasts(t *testing.T) {
	baseline := map[string]int{"CALLBACK_ARGUMENT_BINDING": 1, "CALLBACK_PARAMETER_FORWARDING": 1, "CALLBACK_PARAMETER_INVOCATION": 1}
	for _, c := range []struct {
		name, code string
		missing    []string
	}{
		{"baseline", callbackFixture, nil},
		{"rename", strings.NewReplacer("Start", "Launch", "Forward", "Pass", "cb", "job", "task", "work").Replace(callbackFixture), nil},
		{"no-forward", strings.Replace(callbackFixture, "Forward(cb)", "_=cb", 1), []string{"CALLBACK_PARAMETER_FORWARDING"}},
		{"no-invocation", strings.Replace(callbackFixture, "{task()}", "{_=task}", 1), []string{"CALLBACK_PARAMETER_INVOCATION"}},
		{"alias", strings.Replace(callbackFixture, "{task()}", "{alias:=task;alias()}", 1), []string{"CALLBACK_PARAMETER_INVOCATION"}},
		{"overwrite", strings.Replace(callbackFixture, "{task()}", "{task=func(){};task()}", 1), []string{"CALLBACK_PARAMETER_INVOCATION", "CALLBACK_PARAMETER_FORWARDING"}},
		{"parenthesized-overwrite", strings.Replace(callbackFixture, "{task()}", "{(task)=func(){};task()}", 1), []string{"CALLBACK_PARAMETER_INVOCATION", "CALLBACK_PARAMETER_FORWARDING"}},
		{"address", strings.Replace(callbackFixture, "{task()}", "{_= &(task);task()}", 1), []string{"CALLBACK_PARAMETER_INVOCATION", "CALLBACK_PARAMETER_FORWARDING"}},
		{"shadow", strings.Replace(callbackFixture, "{task()}", "{_=task;{task:=func(){};task()}}", 1), []string{"CALLBACK_PARAMETER_INVOCATION"}},
		{"indirect-launch", strings.Replace(callbackFixture, "Start(func(){})", "launch:=Start;launch(func(){})", 1), []string{"CALLBACK_ARGUMENT_BINDING"}},
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
			if !reflect.DeepEqual(callbackProfile(r), expected) {
				t.Fatalf("%+v expected %+v", callbackProfile(r), expected)
			}
			if !hasGap(r, "CALLBACK_EXECUTION_UNRESOLVED") {
				t.Fatal("invocation mistaken for execution")
			}
			for _, e := range r.CallbackRelations {
				if e.Authority != "STATIC_TYPES_AND_SYNTAX" || e.ExecutionScope == nil {
					t.Fatal("wrong callback authority/scope")
				}
				for _, loc := range []Location{e.From, e.To, *e.ExecutionScope} {
					if loc.ContentSHA256 != Digest(c.code) || loc.EndOffset <= loc.Offset || loc.EndOffset > len(c.code) {
						t.Fatal("unbound callback")
					}
				}
				if e.Kind == "CALLBACK_PARAMETER_FORWARDING" && e.ParameterDeclaration == nil {
					t.Fatal("forwarding loses source declaration")
				}
			}
			if (c.name == "overwrite" || c.name == "address") && !hasGap(r, "CALLBACK_VALUE_FLOW_UNRESOLVED") {
				t.Fatal("mutation gap absent")
			}
		})
	}
}
func TestCallbackScopesAndLimits(t *testing.T) {
	code := `package p;func Start(cb func(),ok bool){if ok {go func(){Forward(cb)}()}};func Forward(task func()){if task!=nil{task()}};func Entry(){Start(func(){},true)}`
	r, err := Observe(request(code))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range r.CallbackRelations {
		if e.Kind == "CALLBACK_PARAMETER_FORWARDING" && len(e.Conditions) != 0 {
			t.Fatal("closure inherited creation condition")
		}
		if e.Kind == "CALLBACK_PARAMETER_INVOCATION" && !reflect.DeepEqual(e.Conditions, []string{"if(task!=nil)"}) {
			t.Fatal("invocation guard lost")
		}
	}
	code = `package p;func Start(cb func(),values ...int){cb()};func Entry(){Start(func(){},1)}`
	r, err = Observe(request(code))
	if err != nil || len(r.CallbackRelations) != 0 {
		t.Fatal("variadic callback invented")
	}
	code = "package p;func f(cb func()){" + strings.Repeat("cb();", 10001) + "}"
	if _, err = Observe(request(code)); err == nil {
		t.Fatal("combined callback limit not rejected")
	}
}
func TestCallbackOrderAndMutationInClosure(t *testing.T) {
	input := request("package p;func Entry(){Start(func(){})}")
	extra := "package p;func Start(cb func()){go func(){Forward(cb)}()};func Forward(task func()){task()}"
	input.Files = append(input.Files, SourceFile{Path: "other.go", Content: extra, SHA256: Digest(extra)})
	r, err := Observe(input)
	if err != nil {
		t.Fatal(err)
	}
	input.Files[0], input.Files[1] = input.Files[1], input.Files[0]
	again, err := Observe(input)
	if err != nil || !reflect.DeepEqual(r, again) {
		t.Fatal("callback order changed")
	}
	code := strings.Replace(callbackFixture, "Forward(cb)", "cb=func(){};Forward(cb)", 1)
	r, err = Observe(request(code))
	if err != nil || callbackProfile(r)["CALLBACK_ARGUMENT_BINDING"] != 0 || callbackProfile(r)["CALLBACK_PARAMETER_FORWARDING"] != 0 {
		t.Fatal("capture overwrite preserved original callback")
	}
}
