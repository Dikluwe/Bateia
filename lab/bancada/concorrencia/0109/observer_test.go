/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/observer_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func request(code string) Request {
	return Request{SchemaVersion: "1", Origin: "AUTHOR_FIXTURE", Files: []SourceFile{{Path: "source.go", Content: code, SHA256: Digest(code)}}}
}
func profile(r Report) map[string]int {
	p := map[string]int{}
	for _, o := range r.Observations {
		p[o.Property+"="+o.Value]++
	}
	return p
}
func hasGap(r Report, code string) bool {
	for _, g := range r.Gaps {
		if g.Code == code {
			return true
		}
	}
	return false
}

func TestObserverContrastAndAuthority(t *testing.T) {
	for _, c := range []struct {
		code string
		want map[string]int
	}{
		{"package p;func f(){ c:=make(chan int,3);_=c }", map[string]int{"channel_buffer_syntax=CAPACITY_LITERAL:3": 1}},
		{"package p;func f(){ c:=make(chan int);_=c }", map[string]int{"channel_buffer_syntax=UNBUFFERED_FORM": 1}},
		{"package p;func f(c chan int){select{case c<-1:default:return}}", map[string]int{"select_shape=SELECT_WITH_DEFAULT": 1, "channel_send=SELECT_WITH_DEFAULT": 1}},
		{"package p;func g(other chan int){select{case other<-1:default:return}}", map[string]int{"select_shape=SELECT_WITH_DEFAULT": 1, "channel_send=SELECT_WITH_DEFAULT": 1}},
		{"package p;func f(c chan int){select{case c<-1:}}", map[string]int{"select_shape=WAITING_SELECT": 1, "channel_send=WAITING_SELECT": 1}},
		{"package p;func f(){_=`go make(chan int,3); ctx.Done()`}", map[string]int{}},
	} {
		r, err := Observe(request(c.code))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(profile(r), c.want) {
			t.Fatalf("unexpected profile: %+v", r)
		}
		for _, o := range r.Observations {
			if o.Authority != "SYNTAX_ONLY" {
				t.Fatal("excessive authority")
			}
			if o.Location.Offset < 0 || o.Location.EndOffset > len(c.code) || o.Location.ContentSHA256 != Digest(c.code) {
				t.Fatal("unbound location")
			}
		}
		if len(r.Observations) > 0 && (!hasGap(r, "ADMISSION_BOUND_UNRESOLVED") || !hasGap(r, "WORKER_BOUND_UNRESOLVED")) {
			t.Fatal("buffer promoted to workload bound")
		}
	}
}

func TestObserverConditionsClosuresAndPhysicalLocations(t *testing.T) {
	code := "package p\n//line fictitious.go:900\nfunc f(c chan int,ok bool){if ok { go func(){ c<-1 }() }}"
	r, err := Observe(request(code))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Observations) != 2 {
		t.Fatalf("observations %+v", r)
	}
	for _, o := range r.Observations {
		if o.Location.Path != "source.go" || o.Location.Line != 3 {
			t.Fatal("line directive corrupted provenance")
		}
		if o.Property == "goroutine_launch" && !reflect.DeepEqual(o.Conditions, []string{"if(ok)"}) {
			t.Fatal("missing condition")
		}
		if o.Property == "channel_send" && (len(o.Conditions) != 0 || !strings.Contains(o.Function, "/literal@")) {
			t.Fatal("closure attributed to outer execution scope")
		}
	}
}

func TestObserverShadowingAndCancellationAreNotSemanticProof(t *testing.T) {
	input := request("package p; func f(){ _=make(chan int,2) }")
	extra := "package p;var make func()"
	input.Files = append(input.Files, SourceFile{Path: "other.go", Content: extra, SHA256: Digest(extra)})
	r, err := Observe(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Observations) != 0 || !hasGap(r, "MAKE_NAME_SHADOWED") {
		t.Fatalf("shadowed make: %+v", r)
	}
	r, err = Observe(request("package p;func f(make func(), c chan int){ _=make(chan int);c<-1 }"))
	if err != nil {
		t.Fatal(err)
	}
	if !hasGap(r, "MAKE_NAME_SHADOWED") || profile(r)["channel_buffer_syntax=UNBUFFERED_FORM"] != 0 {
		t.Fatal("local shadow inferred as builtin")
	}
	r, err = Observe(request("package p;func f(notContext any){<-notContext.Done()}"))
	if err != nil {
		t.Fatal(err)
	}
	if !hasGap(r, "DONE_RECEIVER_TYPE_UNRESOLVED") || !hasGap(r, "CANCELLATION_COMPLETION_UNRESOLVED") {
		t.Fatal("Done name credited as verified context")
	}
}

func TestObserverIntegrityLimitsPartialParseAndNoMutation(t *testing.T) {
	input := request("package p;func f(){_ = make(chan int,2)}")
	extra := "package p;func broken( {"
	input.Files = append(input.Files, SourceFile{Path: "broken.go", Content: extra, SHA256: Digest(extra)})
	before, _ := json.Marshal(input)
	r, err := Observe(input)
	if err != nil || r.Status != "PARTIAL_PARSE" || len(r.Observations) != 1 {
		t.Fatalf("partial: %+v %v", r, err)
	}
	after, _ := json.Marshal(input)
	if string(before) != string(after) {
		t.Fatal("request mutated")
	}
	reversed := input
	reversed.Files = append([]SourceFile(nil), input.Files...)
	reversed.Files[0], reversed.Files[1] = reversed.Files[1], reversed.Files[0]
	r2, err := Observe(reversed)
	if err != nil || !reflect.DeepEqual(r, r2) {
		t.Fatal("file order changed report")
	}
	for _, mutate := range []func(*Request){
		func(r *Request) { r.Files[0].Content += " changed" },
		func(r *Request) { r.Files[0].Path = "../outside.go" },
		func(r *Request) { r.Files = append(r.Files, r.Files[0]) },
		func(r *Request) {
			r.Files[0].Content = strings.Repeat("x", (256<<10)+1)
			r.Files[0].SHA256 = Digest(r.Files[0].Content)
		},
		func(r *Request) { r.Build.Tags = make([]string, 21) },
		func(r *Request) { r.Origin = "bad\norigin" },
	} {
		candidate := request("package p;func f(){}")
		mutate(&candidate)
		if _, err := Observe(candidate); err == nil {
			t.Fatal("invalid input accepted")
		}
	}
}
