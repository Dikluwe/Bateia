/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/configuration_guard_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"reflect"
	"strings"
	"testing"
)

const guardFixture = `package p
 type Gate chan int
 type P struct{tokens Gate;tasks Gate}
 func(l Gate)release(){<-l}
 func(p *P)run(){p.tokens<-1;defer p.tokens.release()}
 func(p *P)guard(){if p.tasks!=nil{panic(1)}}
 func(p *P)configure(){p.guard();p.tokens=make(Gate,2)}
 func(p *P)init(){p.tasks=make(Gate)}
`

func TestConfigurationGuardContrasts(t *testing.T) {
	for _, c := range []struct {
		name, code string
		want       int
	}{
		{"baseline", guardFixture, 1},
		{"rename", strings.NewReplacer("guard", "check", "tasks", "state", "tokens", "permits").Replace(guardFixture), 1},
		{"reversed", strings.Replace(guardFixture, "p.tasks!=nil", "nil!=p.tasks", 1), 1},
		{"absent", strings.Replace(guardFixture, "p.guard();", "", 1), 0},
		{"conditional", strings.Replace(guardFixture, "p.guard();", "if p.tasks==nil{p.guard()};", 1), 0},
		{"deferred", strings.Replace(guardFixture, "p.guard();", "defer p.guard();", 1), 0},
		{"go-call", strings.Replace(guardFixture, "p.guard();", "go p.guard();", 1), 0},
		{"later", strings.Replace(guardFixture, "p.guard();p.tokens=make(Gate,2)", "p.tokens=make(Gate,2);p.guard()", 1), 0},
		{"different-base", strings.Replace(guardFixture, "configure(){p.guard();", "configure(q *P){q.guard();", 1), 0},
		{"rebind", strings.Replace(guardFixture, "p.guard();", "p.guard();p=&P{};", 1), 0},
		{"address", strings.Replace(guardFixture, "p.guard();", "p.guard();_= &p;", 1), 0},
		{"capture", strings.Replace(guardFixture, "p.guard();", "p.guard();_=func(){_=p};", 1), 0},
		{"shadow-panic", guardFixture + "\nfunc panic(int){}", 0},
		{"shadow-nil", guardFixture + "\nvar nil Gate", 0},
		{"recover", strings.Replace(guardFixture, "guard(){", "guard(){defer func(){recover()}();", 1), 0},
		{"wrong-predicate", strings.Replace(guardFixture, "p.tasks!=nil", "p.tasks==nil", 1), 0},
		{"closure", strings.Replace(guardFixture, "p.guard();p.tokens=make(Gate,2)", "func(){p.guard();p.tokens=make(Gate,2)}()", 1), 0},
		{"jump", strings.Replace(guardFixture, "configure(){", "configure(){goto check;check:", 1), 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			r, err := Observe(request(c.code))
			if err != nil {
				t.Fatal(err)
			}
			if len(r.ConfigurationGuards) != c.want {
				t.Fatalf("%+v", r.ConfigurationGuards)
			}
			for _, g := range r.ConfigurationGuards {
				if g.Atomicity != "UNRESOLVED" || g.State != "PRIOR_GUARD_SAME_BASE_OBJECT" || len(g.CheckedFieldWrites) != 1 {
					t.Fatal("guard promoted to lifecycle proof")
				}
				for _, loc := range []Location{g.Write, g.GuardCall, g.GuardDefinition, g.Predicate, g.Panic, g.CheckedFieldWrites[0]} {
					if loc.ContentSHA256 != Digest(c.code) || loc.EndOffset <= loc.Offset {
						t.Fatal("unbound guard")
					}
				}
			}
		})
	}
}
func TestConfigurationGuardOrderAndInterveningEffects(t *testing.T) {
	// The projection records a prior check, not the field value at the later write.
	code := strings.Replace(guardFixture, "p.guard();", "p.guard();p.tasks=nil;", 1)
	r, err := Observe(request(code))
	if err != nil || len(r.ConfigurationGuards) != 1 || r.ConfigurationGuards[0].Atomicity != "UNRESOLVED" {
		t.Fatal("prior check confused with atomic protection")
	}
	input := request(guardFixture)
	extra := "package p;func other(p *P){p.tasks=make(Gate)}"
	input.Files = append(input.Files, SourceFile{Path: "other.go", Content: extra, SHA256: Digest(extra)})
	r, err = Observe(input)
	if err != nil || len(r.ConfigurationGuards[0].CheckedFieldWrites) != 2 {
		t.Fatal("declaration-scoped sentinel inventory wrong")
	}
	input.Files[0], input.Files[1] = input.Files[1], input.Files[0]
	again, err := Observe(input)
	if err != nil || !reflect.DeepEqual(r, again) {
		t.Fatal("guard order changed")
	}
	r, err = Observe(request(guardFixture + ";var bad=missing"))
	if err != nil || len(r.ConfigurationGuards) != 0 {
		t.Fatal("invalid package guarded")
	}
}

func TestConfigurationGuardLimit(t *testing.T) {
	code := strings.Replace(guardFixture, "p.guard();p.tokens=make(Gate,2)", strings.Repeat("p.guard();", 101)+strings.Repeat("p.tokens=nil;", 100), 1)
	if _, err := Observe(request(code)); err == nil {
		t.Fatal("configuration projection overflow accepted")
	}
}
