/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0121/declarations_test.md
 * @layer lab
 * @updated 2026-10-07
 */
package declarations

import (
	concurrency "bateia/lab/bancada/concorrencia/0109"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func fixture(code string) Request {
	return Request{Source: concurrency.Request{SchemaVersion: "1", Origin: "AUTHOR_DECLARATION_FIXTURE", Files: []concurrency.SourceFile{{Path: "source.go", Content: code, SHA256: concurrency.Digest(code)}}}}
}
func TestAttachedDocumentationAndEvidence(t *testing.T) {
	code := `package p
// Package comment only.
type P struct { ch chan int }
// Configure changes a field; this is a source declaration.
func (p *P) Configure(){p.ch=make(chan int,2);p.Run()}
// Run sends a value.
func (p *P) Run(){p.ch<-1}
type API interface {
 // Configure interface text differs from the implementation.
 Configure()
}
// Detached comment.

func Bare(){} // trailing comment
func hidden(){}
`
	request := fixture(code)
	report, err := Inspect(request)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Declarations) != 3 {
		t.Fatalf("declarations: %+v", report.Declarations)
	}
	observed, err := concurrency.Observe(request.Source)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(observed, report.ObserverReport) {
		t.Fatal("observer projection changed")
	}
	methods, interfaces := 0, 0
	for _, d := range report.Declarations {
		if d.Authority != "DOCUMENTATION_TEXT_ONLY" {
			t.Fatal(d)
		}
		if d.Owner.Kind == "INTERFACE_METHOD" {
			interfaces++
			if len(d.Evidence) != 0 {
				t.Fatal("interface associated by name")
			}
		}
		if d.Owner.Kind == "METHOD" {
			methods++
			if len(d.Evidence) == 0 {
				t.Fatal("missing body evidence")
			}
		}
		if strings.Contains(d.Text, "Detached") || strings.Contains(d.Text, "trailing") {
			t.Fatal("unattached comment accepted")
		}
	}
	if methods != 2 || interfaces != 1 {
		t.Fatal(methods, interfaces)
	}
}
func TestExternalBindingIntegrity(t *testing.T) {
	request := fixture("package p\n// Configure is documented.\nfunc Configure(){}\n")
	baseline, err := Inspect(request)
	if err != nil {
		t.Fatal(err)
	}
	doc := concurrency.SourceFile{Path: "README.md", Content: "heading\nExternal promise.\n", SHA256: concurrency.Digest("heading\nExternal promise.\n")}
	request.Documents = []concurrency.SourceFile{doc}
	request.Bindings = []Binding{{DocumentPath: doc.Path, FirstLine: 2, LastLine: 2, Owner: baseline.Declarations[0].Owner.Location}}
	report, err := Inspect(request)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Declarations) != 2 {
		t.Fatal(report)
	}
	for _, d := range report.Declarations {
		if d.Association == "CALLER_PROVIDED_BINDING" && (d.Text != "External promise." || d.Document.Offset != 8) {
			t.Fatal(d)
		}
	}
	mutations := []func(*Request){
		func(r *Request) { r.Bindings[0].Owner.ContentSHA256 = "wrong" },
		func(r *Request) { r.Bindings[0].Owner.EndOffset-- },
		func(r *Request) { r.Bindings[0].Owner.Line++ },
		func(r *Request) { r.Bindings[0].LastLine = 20 },
		func(r *Request) { r.Bindings = append(r.Bindings, r.Bindings[0]) },
		func(r *Request) { r.Documents[0].SHA256 = "wrong" },
		func(r *Request) { r.Documents[0].Path = "source.go" },
	}
	for i, mutate := range mutations {
		raw, _ := json.Marshal(request)
		var altered Request
		json.Unmarshal(raw, &altered)
		mutate(&altered)
		if _, err := Inspect(altered); err == nil {
			t.Fatalf("invalid binding accepted %d", i)
		}
	}
}
func TestPhysicalLocationsAndTypeFailure(t *testing.T) {
	code := "package p\nimport \"missing\"\n//line fake.go:999\n// New text.\nfunc New(){ missing.Call() }\n"
	report, err := Inspect(fixture(code))
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Declarations) != 1 || len(report.ObserverReport.TypeDiagnostics) == 0 {
		t.Fatal(report)
	}
	declaration := report.Declarations[0]
	if declaration.Document.Line != 3 || declaration.Owner.Location.Line != 5 || declaration.Document.Path != "source.go" {
		t.Fatal(declaration)
	}
	if code[declaration.Document.Offset:declaration.Document.EndOffset] != "//line fake.go:999\n// New text." {
		t.Fatal(declaration)
	}
	if len(declaration.Evidence) != 0 {
		t.Fatal("unresolved type call linked")
	}
	invalid := fixture("package p\n// New text.\nfunc New( {\n")
	partial, err := Inspect(invalid)
	if err != nil {
		t.Fatal(err)
	}
	if len(partial.Declarations) != 0 {
		t.Fatal("parse-invalid owner accepted")
	}
}
func TestOrderAndClosureScope(t *testing.T) {
	request := fixture("package p\n// New source.\nfunc New(){f:=func(){ch:=make(chan int,1);ch<-1};f()}\n")
	other := concurrency.SourceFile{Path: "other.go", Content: "package p\n// Other source.\nfunc Other(){}\n"}
	other.SHA256 = concurrency.Digest(other.Content)
	request.Source.Files = append(request.Source.Files, other)
	first, err := Inspect(request)
	if err != nil {
		t.Fatal(err)
	}
	request.Source.Files[0], request.Source.Files[1] = request.Source.Files[1], request.Source.Files[0]
	second, err := Inspect(request)
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(first)
	b, _ := json.Marshal(second)
	if string(a) != string(b) {
		t.Fatal("order changed projection")
	}
	found := false
	for _, d := range first.Declarations {
		for _, fact := range d.Evidence {
			if fact.Kind == "OBSERVATION:channel_send" {
				found = true
				if fact.Authority != "LEXICAL_CONTAINMENT_ONLY" {
					t.Fatal(fact)
				}
			}
		}
	}
	if !found {
		t.Fatal("closure evidence not inventoried lexically")
	}
}
