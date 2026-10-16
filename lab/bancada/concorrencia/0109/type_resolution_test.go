/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/type_resolution_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"reflect"
	"testing"
)

func TestNamedChannelTypesAndContrasts(t *testing.T) {
	for _, c := range []struct {
		name, code string
		want       bool
	}{
		{"named", "package p; type C chan int; func f(){_=make(C,3)}", true},
		{"alias-chain", "package p; type C = chan int; type D = C; func f(){_=make(D,3)}", true},
		{"local", "package p; func f(){type C chan int;_=make(C,3)}", true},
		{"rename", "package p; type Renamed chan int; func other(){_=make(Renamed,3)}", true},
		{"slice", "package p; type C []int; func f(){_=make(C,3)}", false},
		{"map", "package p; type C map[int]int; func f(){_=make(C,3)}", false},
		{"shadow-make", "package p; type C chan int; func f(){make:=func(int,int)int{return 0};_=make(1,3)}", false},
		{"shadow-type", "package p; type C chan int; func f(){C:=0;_=make(C,3)}", false},
		{"cycle", "package p; type C = D; type D = C; func f(){_=make(C,3)}", false},
		{"missing-import", "package p; import \"missing.invalid/package\"; type C chan int; func f(){_=make(C,3)}", false},
		{"unrelated-error", "package p; type C chan int; func f(){_=make(C,3);_=missing}", false},
		{"invalid-capacity", "package p; type C chan int; func f(){_=make(C,-3)}", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			r, err := Observe(request(c.code))
			if err != nil {
				t.Fatal(err)
			}
			got := profile(r)["channel_type_identity=RESOLVED_CHANNEL_TYPE"] == 1
			if got != c.want {
				t.Fatalf("resolution mismatch: %+v", r)
			}
			for _, o := range r.Observations {
				if o.Property == "channel_type_identity" && o.Authority != "GO_TYPES_SUPPLIED_PACKAGE" {
					t.Fatal("wrong authority")
				}
				if o.Location.ContentSHA256 != Digest(c.code) {
					t.Fatal("unbound source")
				}
			}
			if got && (!hasGap(r, "WORKER_BOUND_UNRESOLVED") || !hasGap(r, "ADMISSION_BOUND_UNRESOLVED")) {
				t.Fatal("type promoted to bound")
			}
			if (c.name == "slice" || c.name == "map") && hasGap(r, "NAMED_CHANNEL_TYPE_UNRESOLVED") {
				t.Fatal("non-channel falsely marked unresolved")
			}
			if c.name == "missing-import" && !hasGap(r, "SUPPLIED_PACKAGE_TYPECHECK_FAILED") {
				t.Fatal("missing failure gap")
			}
		})
	}
}

func TestTypeResolutionGroupsAndDeterminism(t *testing.T) {
	input := request("package p;func f(){_=make(C,3)}")
	extra := "package p;type C chan int"
	input.Files = append(input.Files, SourceFile{Path: "other.go", Content: extra, SHA256: Digest(extra)})
	r, err := Observe(input)
	if err != nil || profile(r)["channel_type_identity=RESOLVED_CHANNEL_TYPE"] != 1 {
		t.Fatalf("cross-file %+v %v", r, err)
	}
	input.Files[0], input.Files[1] = input.Files[1], input.Files[0]
	reordered, err := Observe(input)
	if err != nil || !reflect.DeepEqual(r, reordered) {
		t.Fatal("order changed typing")
	}
	input.Files[0].Path = "different/other.go"
	r, err = Observe(input)
	if err != nil || profile(r)["channel_type_identity=RESOLVED_CHANNEL_TYPE"] != 0 {
		t.Fatal("type leaked across package group")
	}
	input = request("package p;type C chan int;func f(){_=make(C,3)}")
	extra = "package p;func broken( {"
	input.Files = append(input.Files, SourceFile{Path: "broken.go", Content: extra, SHA256: Digest(extra)})
	r, err = Observe(input)
	if err != nil || r.Status != "PARTIAL_PARSE" || profile(r)["channel_type_identity=RESOLVED_CHANNEL_TYPE"] != 0 {
		t.Fatal("incomplete parsed package credited")
	}
}

func TestNamedChannelCapacityShapes(t *testing.T) {
	for _, c := range []struct{ code, value string }{
		{"package p;type C chan int;func f(n int){_=make(C,n)}", "CAPACITY_EXPRESSION"},
		{"package p;type C chan int;func f(){_=make(C)}", "UNBUFFERED_FORM"},
		{"package p;type C chan int;func f(){_=make(C,0)}", "CAPACITY_LITERAL:0"},
		{"package p;type C chan<- int;func f(){_=make(C,3)}", "CAPACITY_LITERAL:3"},
	} {
		r, err := Observe(request(c.code))
		if err != nil || profile(r)["channel_buffer_syntax="+c.value] != 1 || profile(r)["channel_type_identity=RESOLVED_CHANNEL_TYPE"] != 1 {
			t.Fatalf("capacity %+v %v", r, err)
		}
	}
	input := request("package p;type C chan int;func f(){_=make(C,3)}")
	bad := "func broken( {"
	input.Files = append(input.Files, SourceFile{Path: "broken.go", Content: bad, SHA256: Digest(bad)})
	r, err := Observe(input)
	if err != nil || r.Status != "PARTIAL_PARSE" || profile(r)["channel_type_identity=RESOLVED_CHANNEL_TYPE"] != 0 {
		t.Fatalf("missing package not contained: %+v %v", r, err)
	}
}
