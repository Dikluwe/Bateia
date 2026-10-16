/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0121/declarations.md
 * @layer lab
 * @updated 2026-10-07
 */
package declarations

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	concurrency "bateia/lab/bancada/concorrencia/0109"
)

type Binding struct {
	DocumentPath string               `json:"document_path"`
	FirstLine    int                  `json:"first_line"`
	LastLine     int                  `json:"last_line"`
	Owner        concurrency.Location `json:"owner"`
}
type Request struct {
	Source      concurrency.Request      `json:"source"`
	Environment *concurrency.Environment `json:"environment,omitempty"`
	Documents   []concurrency.SourceFile `json:"documents,omitempty"`
	Bindings    []Binding                `json:"bindings,omitempty"`
}
type Owner struct {
	Kind     string               `json:"kind"`
	Name     string               `json:"name"`
	Location concurrency.Location `json:"location"`
}
type Evidence struct {
	Kind      string               `json:"kind"`
	Location  concurrency.Location `json:"location"`
	Authority string               `json:"authority"`
}
type Declaration struct {
	ID          string               `json:"id"`
	Text        string               `json:"text"`
	Document    concurrency.Location `json:"document"`
	Owner       Owner                `json:"owner"`
	Association string               `json:"association"`
	Authority   string               `json:"authority"`
	Evidence    []Evidence           `json:"evidence"`
}
type Report struct {
	SchemaVersion  string             `json:"schema_version"`
	ObserverReport concurrency.Report `json:"observer_report"`
	Declarations   []Declaration      `json:"declarations"`
	APIOwners      []Owner            `json:"api_owners"`
	Scope          string             `json:"scope"`
	Limitations    []string           `json:"limitations"`
}
type ownerRecord struct {
	owner Owner
	body  *ast.BlockStmt
	doc   *ast.CommentGroup
	file  concurrency.SourceFile
}

func location(fset *token.FileSet, source concurrency.SourceFile, start, end token.Pos) concurrency.Location {
	p := fset.PositionFor(start, false)
	return concurrency.Location{Path: source.Path, Line: p.Line, Column: p.Column,
		Offset: p.Offset, EndOffset: fset.PositionFor(end, false).Offset, ContentSHA256: source.SHA256}
}
func Inspect(request Request) (Report, error) {
	var observed concurrency.Report
	var err error
	if request.Environment == nil {
		observed, err = concurrency.Observe(request.Source)
	} else {
		observed, err = concurrency.ObserveWithEnvironment(request.Source, *request.Environment)
	}
	if err != nil {
		return Report{}, err
	}
	if len(request.Bindings) > 1000 {
		return Report{}, fmt.Errorf("binding limit exceeded")
	}
	if len(request.Documents) > 32 {
		return Report{}, fmt.Errorf("document count limit exceeded")
	}
	docs := map[string]concurrency.SourceFile{}
	total := 0
	for _, doc := range request.Documents {
		if doc.Path == "" || doc.Path == "." || len(doc.Path) > 300 || !utf8.ValidString(doc.Path) || strings.TrimSpace(doc.Path) != doc.Path || strings.IndexFunc(doc.Path, unicode.IsControl) >= 0 || path.IsAbs(doc.Path) || path.Clean(doc.Path) != doc.Path || strings.HasPrefix(doc.Path, "../") || strings.Contains(doc.Path, "\\") {
			return Report{}, fmt.Errorf("invalid document path")
		}
		if _, found := docs[doc.Path]; found {
			return Report{}, fmt.Errorf("duplicate document path")
		}
		total += len(doc.Content)
		if len(doc.Content) > 256<<10 || total > 2<<20 {
			return Report{}, fmt.Errorf("document byte limit exceeded")
		}
		if !utf8.ValidString(doc.Content) || concurrency.Digest(doc.Content) != doc.SHA256 {
			return Report{}, fmt.Errorf("document integrity mismatch")
		}
		docs[doc.Path] = doc
	}
	encoded, _ := json.Marshal(request.Documents)
	if len(encoded) > 2<<20 {
		return Report{}, fmt.Errorf("serialized document limit exceeded")
	}
	fset := token.NewFileSet()
	owners := map[concurrency.Location]ownerRecord{}
	for _, source := range request.Source.Files {
		if _, found := docs[source.Path]; found {
			return Report{}, fmt.Errorf("document/source path collision")
		}
		file, parseErr := parser.ParseFile(fset, source.Path, source.Content, parser.ParseComments|parser.AllErrors|parser.SkipObjectResolution)
		if parseErr != nil {
			continue
		}
		add := func(kind, name string, node ast.Node, body *ast.BlockStmt, doc *ast.CommentGroup) {
			loc := location(fset, source, node.Pos(), node.End())
			owners[loc] = ownerRecord{Owner{kind, name, loc}, body, doc, source}
		}
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if !ast.IsExported(d.Name.Name) {
					continue
				}
				name, kind := d.Name.Name, "FUNCTION"
				if d.Recv != nil {
					receiver := d.Recv.List[0].Type
					a, b := location(fset, source, receiver.Pos(), receiver.End()).Offset, location(fset, source, receiver.Pos(), receiver.End()).EndOffset
					name, kind = source.Content[a:b]+"."+name, "METHOD"
				}
				add(kind, name, d, d.Body, d.Doc)
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					typed, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					iface, ok := typed.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}
					for _, field := range iface.Methods.List {
						if len(field.Names) != 1 || !ast.IsExported(field.Names[0].Name) {
							continue
						}
						if _, ok := field.Type.(*ast.FuncType); !ok {
							continue
						}
						add("INTERFACE_METHOD", typed.Name.Name+"."+field.Names[0].Name, field, nil, field.Doc)
					}
				}
			}
		}
	}
	facts := []Evidence{}
	for _, o := range observed.Observations {
		facts = append(facts, Evidence{"OBSERVATION:" + o.Property, o.Location, "LEXICAL_CONTAINMENT_ONLY"})
	}
	for _, r := range observed.Relations {
		facts = append(facts, Evidence{r.Kind, r.From, "LEXICAL_CONTAINMENT_ONLY"})
	}
	for _, g := range observed.ConfigurationGuards {
		facts = append(facts, Evidence{"CONFIGURATION_WRITE", g.Write, "LEXICAL_CONTAINMENT_ONLY"}, Evidence{"CONFIGURATION_GUARD_CALL", g.GuardCall, "LEXICAL_CONTAINMENT_ONLY"})
	}
	sort.Slice(facts, func(i, j int) bool {
		a, _ := json.Marshal(facts[i])
		b, _ := json.Marshal(facts[j])
		return string(a) < string(b)
	})
	if len(owners) > 10000 {
		return Report{}, fmt.Errorf("API owner limit exceeded")
	}
	ownerList := []Owner{}
	for _, record := range owners {
		ownerList = append(ownerList, record.owner)
	}
	sort.Slice(ownerList, func(i, j int) bool {
		a, b := ownerList[i].Location, ownerList[j].Location
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Offset < b.Offset
	})
	result := Report{"api-declarations-v1", observed, []Declaration{}, ownerList, "SUPPLIED_DOCUMENTATION_AND_LEXICAL_SOURCE_LINKS", []string{
		"Documentation text is a declaration, not a verified behavior or semantic verdict.",
		"External associations are caller supplied; valid anchors do not authenticate interpretation.",
		"Lexical containment includes closures and does not prove execution, agreement or applicability.",
		"Interface implementation, concurrent atomicity and repository build selection are unresolved.",
		"Source files with parse errors have no declaration inventory; type failures do not erase syntax documentation.",
	}}
	count := 0
	appendDeclaration := func(record ownerRecord, doc concurrency.Location, text, association string) error {
		evidence := []Evidence{}
		if record.body != nil {
			body := location(fset, record.file, record.body.Pos(), record.body.End())
			for _, fact := range facts {
				loc := fact.Location
				if loc.Path == body.Path && loc.ContentSHA256 == body.ContentSHA256 && loc.Offset >= body.Offset && loc.EndOffset <= body.EndOffset {
					if len(evidence) == 0 || evidence[len(evidence)-1] != fact {
						evidence = append(evidence, fact)
					}
				}
			}
		}
		count += 1 + len(evidence)
		if count > 10000 {
			return fmt.Errorf("declaration/link limit exceeded")
		}
		identity, _ := json.Marshal(struct {
			Owner       Owner
			Document    concurrency.Location
			Association string
		}{record.owner, doc, association})
		result.Declarations = append(result.Declarations, Declaration{concurrency.Digest(string(identity)), text, doc, record.owner, association, "DOCUMENTATION_TEXT_ONLY", evidence})
		return nil
	}
	for _, record := range owners {
		if record.doc == nil {
			continue
		}
		doc := location(fset, record.file, record.doc.Pos(), record.doc.End())
		if err := appendDeclaration(record, doc, record.file.Content[doc.Offset:doc.EndOffset], "AST_ATTACHED_GO_DOC"); err != nil {
			return Report{}, err
		}
	}
	seen := map[Binding]bool{}
	for _, binding := range request.Bindings {
		if seen[binding] {
			return Report{}, fmt.Errorf("duplicate binding")
		}
		seen[binding] = true
		record, ok := owners[binding.Owner]
		if !ok {
			return Report{}, fmt.Errorf("binding owner is not an exact supplied API declaration")
		}
		source, ok := docs[binding.DocumentPath]
		if !ok {
			return Report{}, fmt.Errorf("binding document missing")
		}
		lines := strings.Split(source.Content, "\n")
		if binding.FirstLine < 1 || binding.LastLine < binding.FirstLine || binding.LastLine > len(lines) {
			return Report{}, fmt.Errorf("invalid document span")
		}
		offset := 0
		for _, line := range lines[:binding.FirstLine-1] {
			offset += len(line) + 1
		}
		text := strings.Join(lines[binding.FirstLine-1:binding.LastLine], "\n")
		if strings.TrimSpace(text) == "" {
			return Report{}, fmt.Errorf("empty document span")
		}
		doc := concurrency.Location{Path: source.Path, Line: binding.FirstLine, Column: 1, Offset: offset, EndOffset: offset + len(text), ContentSHA256: source.SHA256}
		if err := appendDeclaration(record, doc, text, "CALLER_PROVIDED_BINDING"); err != nil {
			return Report{}, err
		}
	}
	sort.Slice(result.Declarations, func(i, j int) bool {
		a, b := result.Declarations[i], result.Declarations[j]
		if a.Document.Path != b.Document.Path {
			return a.Document.Path < b.Document.Path
		}
		if a.Document.Offset != b.Document.Offset {
			return a.Document.Offset < b.Document.Offset
		}
		return a.ID < b.ID
	})
	return result, nil
}
