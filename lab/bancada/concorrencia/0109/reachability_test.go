/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/reachability_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"reflect"
	"strings"
	"testing"
)

func callbackSite(t *testing.T, code string) CallSiteReachability {
	t.Helper()
	r, err := Observe(request(code))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range r.CallSiteReachability {
		if code[s.Site.Offset:s.Site.EndOffset] == "cb()" {
			return s
		}
	}
	t.Fatal("site missing")
	return CallSiteReachability{}
}
func TestReachabilityContrasts(t *testing.T) {
	for _, c := range []struct{ name, code, state string }{
		{"return-before", "package p;func f(cb func()){return;cb()}", "UNREACHABLE_AFTER_TERMINATOR"},
		{"both-branches", "package p;func f(cb func(),ok bool){if ok{return}else{return};cb()}", "UNREACHABLE_AFTER_TERMINATOR"},
		{"one-branch", "package p;func f(cb func(),ok bool){if ok{return};cb()}", "POSSIBLY_REACHABLE"},
		{"panic", "package p;func f(cb func()){panic(1);cb()}", "UNREACHABLE_AFTER_TERMINATOR"},
		{"shadow-panic", "package p;func f(cb func()){panic:=func(int){};panic(1);cb()}", "POSSIBLY_REACHABLE"},
		{"defer-panic", "package p;func f(cb func()){defer panic(1);cb()}", "POSSIBLY_REACHABLE"},
		{"go-panic", "package p;func f(cb func()){go panic(1);cb()}", "POSSIBLY_REACHABLE"},
		{"return-expression", "package p;func f(cb func()int)int{return cb()}", "POSSIBLY_REACHABLE"},
		{"case-return", "package p;func f(cb func(),c chan int){select{case c<-1:return;cb()}}", "UNREACHABLE_AFTER_TERMINATOR"},
		{"closure-after-return", "package p;func f(cb func()){return;func(){cb()}()}", "POSSIBLY_REACHABLE"},
		{"goto", "package p;func f(cb func()){goto target;return;target:cb()}", "UNSUPPORTED_CONTROL_FLOW"},
		{"break", "package p;func f(cb func()){for{break};cb()}", "UNSUPPORTED_CONTROL_FLOW"},
		{"loop-exit-unproved", "package p;func f(cb func()){for{};cb()}", "POSSIBLY_REACHABLE"},
	} {
		t.Run(c.name, func(t *testing.T) {
			s := callbackSite(t, c.code)
			if s.State != c.state {
				t.Fatalf("%+v", s)
			}
			if (s.BlockedBy != nil) != (c.state == "UNREACHABLE_AFTER_TERMINATOR") {
				t.Fatal("missing/false blocker")
			}
			if s.Site.ContentSHA256 != Digest(c.code) || s.Scope.ContentSHA256 != Digest(c.code) || s.Authority != "LIMITED_STRUCTURED_CONTROL" {
				t.Fatal("unbound control site")
			}
		})
	}
}
func TestDeferOrderContrasts(t *testing.T) {
	for _, c := range []struct {
		code  string
		count int
	}{
		{"package p;func done(){};func f(cb func()){defer done();cb()}", 1},
		{"package p;func done(){};func f(cb func(),ok bool){if ok{defer done()};cb()}", 0},
		{"package p;func done(){};func f(cb func()){cb();defer done()}", 0},
		{"package p;func done(){};func f(cb func()){defer done();return;cb()}", 0},
		{"package p;func done(){};func f(cb func()){defer done();go func(){cb()}()}", 0},
		{"package p;func done(){};func f(cb func(),ok bool){defer done();if ok{cb()}}", 1},
	} {
		r, err := Observe(request(c.code))
		if err != nil {
			t.Fatal(err)
		}
		if len(r.DeferOrder) != c.count {
			t.Fatalf("%d expected %d %+v", len(r.DeferOrder), c.count, r.DeferOrder)
		}
		for _, e := range r.DeferOrder {
			if e.Defer.Offset >= e.Site.Offset || e.Defer.ContentSHA256 != Digest(c.code) || e.Authority != "LIMITED_STRUCTURED_CONTROL" {
				t.Fatal("invalid prefix order")
			}
		}
	}
}
func TestControlOrderLimitsAndHistoricalRelations(t *testing.T) {
	code := "package p;func f(cb func()){return;cb()}"
	r, err := Observe(request(code))
	if err != nil || len(r.CallbackRelations) != 1 || r.CallSiteReachability[0].State != "UNREACHABLE_AFTER_TERMINATOR" {
		t.Fatal("lexical link removed/promoted")
	}
	code = "package p;func f(){" + strings.Repeat("println();", 10001) + "}"
	if _, err = Observe(request(code)); err == nil {
		t.Fatal("control overflow accepted")
	}
	input := request("package p;func a(cb func()){defer b();cb()}")
	extra := "package p;func b(){}"
	input.Files = append(input.Files, SourceFile{Path: "other.go", Content: extra, SHA256: Digest(extra)})
	r, err = Observe(input)
	if err != nil {
		t.Fatal(err)
	}
	input.Files[0], input.Files[1] = input.Files[1], input.Files[0]
	reordered, err := Observe(input)
	if err != nil || !reflect.DeepEqual(r, reordered) {
		t.Fatal("control order changed")
	}
}

func TestTypeCaseExpressionsRemainAnOverapproximation(t *testing.T) {
	code := `package p;func f(x any){switch x.(type){case [len([1]int{})]int:}}`
	r, err := Observe(request(code))
	if err != nil || len(r.CallSiteReachability) != 1 {
		t.Fatalf("type case projection %+v %v", r, err)
	}
	// These are lexical sites; this builtin len is a constant expression, not a runtime invocation.
	for _, s := range r.CallSiteReachability {
		if s.State != "POSSIBLY_REACHABLE" {
			t.Fatal("constant evaluation claimed")
		}
	}
}
