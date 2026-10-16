/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/observer.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const Version = "go-concurrency-syntax-0116-v8"

type SourceFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	SHA256  string `json:"sha256"`
}

type BuildContext struct {
	GoVersion string   `json:"go_version,omitempty"`
	GOOS      string   `json:"goos,omitempty"`
	GOARCH    string   `json:"goarch,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type Request struct {
	SchemaVersion string       `json:"schema_version"`
	Origin        string       `json:"origin"`
	Build         BuildContext `json:"build"`
	Files         []SourceFile `json:"files"`
}

type Location struct {
	Path          string `json:"path"`
	Line          int    `json:"line"`
	Column        int    `json:"column"`
	Offset        int    `json:"offset"`
	EndOffset     int    `json:"end_offset"`
	ContentSHA256 string `json:"content_sha256"`
}

type Observation struct {
	Property   string   `json:"property"`
	Value      string   `json:"value"`
	Resource   string   `json:"resource,omitempty"`
	Function   string   `json:"function"`
	Conditions []string `json:"conditions"`
	Authority  string   `json:"authority"`
	Location   Location `json:"location"`
}

type Gap struct {
	Code     string   `json:"code"`
	Function string   `json:"function,omitempty"`
	Location Location `json:"location"`
}

type TypeDiagnostic struct {
	Message  string   `json:"message"`
	Location Location `json:"location"`
}
type Report struct {
	ConfigurationGuards  []ConfigurationGuard   `json:"configuration_guards,omitempty"`
	ChannelIdentity      []ChannelIdentityAudit `json:"channel_identity,omitempty"`
	CallSiteReachability []CallSiteReachability `json:"call_site_reachability,omitempty"`
	DeferOrder           []DeferOrder           `json:"defer_order,omitempty"`
	CallbackRelations    []Relation             `json:"callback_relations,omitempty"`
	Relations            []Relation             `json:"relations,omitempty"`
	TypeDiagnostics      []TypeDiagnostic       `json:"type_diagnostics,omitempty"`
	EnvironmentSHA256    string                 `json:"environment_sha256,omitempty"`
	SchemaVersion        string                 `json:"schema_version"`
	Analyzer             string                 `json:"analyzer"`
	InputSHA256          string                 `json:"input_sha256"`
	Origin               string                 `json:"origin"`
	Build                BuildContext           `json:"declared_build"`
	Scope                string                 `json:"scope"`
	Status               string                 `json:"status"`
	Observations         []Observation          `json:"observations"`
	Gaps                 []Gap                  `json:"gaps"`
	Limitations          []string               `json:"limitations"`
}

func Digest(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func Observe(input Request) (Report, error) {
	return observe(input, nil, "")
}
func ObserveWithEnvironment(input Request, supplied Environment) (Report, error) {
	env, hash, err := canonicalEnvironment(supplied)
	if err != nil {
		return Report{}, err
	}
	return observe(input, &env, hash)
}
func observe(input Request, env *Environment, environmentHash string) (Report, error) {
	if input.SchemaVersion != "1" || !safeText(input.Origin, 500) || len(input.Files) == 0 || len(input.Files) > 32 {
		return Report{}, fmt.Errorf("invalid observation request")
	}
	for _, s := range []string{input.Build.GoVersion, input.Build.GOOS, input.Build.GOARCH} {
		if s != "" && !safeText(s, 100) {
			return Report{}, fmt.Errorf("invalid build declaration")
		}
	}
	if len(input.Build.Tags) > 20 {
		return Report{}, fmt.Errorf("too many build tags")
	}
	for _, tag := range input.Build.Tags {
		if !safeText(tag, 100) {
			return Report{}, fmt.Errorf("invalid build tag")
		}
	}
	input.Files = append([]SourceFile(nil), input.Files...)
	input.Build.Tags = append([]string(nil), input.Build.Tags...)
	sort.Strings(input.Build.Tags)
	seen, total := map[string]bool{}, 0
	for _, file := range input.Files {
		if !safeText(file.Path, 300) || path.IsAbs(file.Path) || path.Clean(file.Path) != file.Path || strings.Contains(file.Path, "\\") || strings.HasPrefix(file.Path, "../") || !strings.HasSuffix(file.Path, ".go") || seen[file.Path] {
			return Report{}, fmt.Errorf("invalid or duplicate source path")
		}
		seen[file.Path] = true
		total += len(file.Content)
		if len(file.Content) > 256<<10 || total > 2<<20 {
			return Report{}, fmt.Errorf("source byte limit exceeded")
		}
		if !utf8.ValidString(file.Content) || file.SHA256 != Digest(file.Content) {
			return Report{}, fmt.Errorf("source integrity mismatch: %s", file.Path)
		}
	}
	sort.Slice(input.Files, func(i, j int) bool { return input.Files[i].Path < input.Files[j].Path })
	encoded, err := json.Marshal(input)
	if err != nil || len(encoded) > 2<<20 {
		return Report{}, fmt.Errorf("request byte limit exceeded")
	}
	r := Report{EnvironmentSHA256: environmentHash, SchemaVersion: "concurrency-observation-v1", Analyzer: Version, InputSHA256: Digest(string(encoded)), Origin: input.Origin, Build: input.Build, Scope: "FILES_AS_SUPPLIED_NOT_BUILD_SELECTED", Status: "INCONCLUSIVE", Observations: []Observation{}, Gaps: []Gap{}, Limitations: []string{
		"Syntax and closed supplied-package types are observed; sources are not built or executed.",
		"Build context and origin are declarations, not independently verified receipts.",
		"Type resolution requires error-free supplied packages; imports need supplied sources/export data. Sizes are fixed gc/amd64, not a verified build.",
		"Buffer capacity does not establish a bound on workers or admitted work.",
		"A Done selector does not establish context identity or completion of cancellation.",
		"No mechanism, absence of defects or benefit to a target is certified.",
	}}
	parsedFiles := []parsedSource{}
	fs := token.NewFileSet()
	brokenGroups := map[string]bool{}
	// Package-level shadowing is checked across supplied files of the same directory/package.
	shadowed := map[string]bool{}
	for _, file := range input.Files {
		tree, parseErr := parser.ParseFile(fs, file.Path, file.Content, 0)
		base := Location{Path: file.Path, Line: 1, Column: 1, ContentSHA256: file.SHA256}
		r.Gaps = append(r.Gaps, Gap{Code: "BUILD_SELECTION_NOT_VERIFIED", Location: base}, Gap{Code: "SYNTAX_PROPERTIES_NOT_TYPECHECKED", Location: base})
		if parseErr != nil {
			brokenGroups[path.Dir(file.Path)+"/*"] = true
			if tree != nil && tree.Name != nil {
				brokenGroups[path.Dir(file.Path)+"/"+tree.Name.Name] = true
			}
			r.Gaps = append(r.Gaps, Gap{Code: "GO_PARSE_FAILED", Location: base})
			continue
		}
		key := path.Dir(file.Path) + "/" + tree.Name.Name
		for _, declaration := range tree.Decls {
			switch d := declaration.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil && d.Name.Name == "make" {
					shadowed[key] = true
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if name.Name == "make" {
								shadowed[key] = true
							}
						}
					case *ast.TypeSpec:
						if s.Name.Name == "make" {
							shadowed[key] = true
						}
					}
				}
			}
		}
		parsedFiles = append(parsedFiles, parsedSource{file, fs, tree})
	}
	resolutions := resolveTypes(parsedFiles, fs, brokenGroups, env)
	for _, p := range parsedFiles {
		typing := resolutions[path.Dir(p.source.Path)+"/"+p.tree.Name.Name]
		for _, diagnostic := range typing.diagnostics {
			pos := fs.PositionFor(diagnostic.Pos, false)
			if pos.Filename == p.source.Path {
				msg := []rune(strings.Join(strings.Fields(diagnostic.Msg), " "))
				if len(msg) > 300 {
					msg = msg[:300]
				}
				r.TypeDiagnostics = append(r.TypeDiagnostics, TypeDiagnostic{Message: string(msg), Location: Location{Path: p.source.Path, Line: pos.Line, Column: pos.Column, Offset: pos.Offset, EndOffset: pos.Offset, ContentSHA256: p.source.SHA256}})
			}
		}
		if !typing.complete {
			r.Gaps = append(r.Gaps, Gap{Code: "SUPPLIED_PACKAGE_TYPECHECK_FAILED", Location: Location{Path: p.source.Path, Line: 1, Column: 1, ContentSHA256: p.source.SHA256}})
		}
		o := walker{environmentSupplied: env != nil, typing: typing, source: p.source, fs: p.fset, report: &r, packageMakeShadowed: shadowed[path.Dir(p.source.Path)+"/"+p.tree.Name.Name]}
		ast.Inspect(p.tree, func(n ast.Node) bool {
			switch f := n.(type) {
			case *ast.FuncDecl:
				if f.Body != nil {
					name := f.Name.Name
					if f.Recv != nil {
						name = o.text(f.Recv.List[0].Type) + "." + name
					}
					o.function(f.Body, name+"@"+strconv.Itoa(o.location(f).Offset))
				}
				return false
			case *ast.FuncLit:
				o.function(f.Body, "literal@"+strconv.Itoa(o.location(f).Offset))
				return false
			}
			return true
		})
	}
	var relationErr error
	r.Relations, relationErr = deriveRelations(parsedFiles, resolutions, &r)
	if relationErr != nil {
		return Report{}, relationErr
	}
	r.Limitations = append(r.Limitations, "Relations describe typed targets and lexical containment, not feasible paths, callback execution or shared channel instances.")
	sort.Slice(r.Observations, func(i, j int) bool {
		a, b := r.Observations[i], r.Observations[j]
		if a.Location.Path != b.Location.Path {
			return a.Location.Path < b.Location.Path
		}
		if a.Location.Offset != b.Location.Offset {
			return a.Location.Offset < b.Location.Offset
		}
		return a.Property < b.Property
	})
	sort.Slice(r.Gaps, func(i, j int) bool {
		a, _ := json.Marshal(r.Gaps[i])
		b, _ := json.Marshal(r.Gaps[j])
		return string(a) < string(b)
	})
	if len(r.Observations) > 0 {
		r.Status = "OBSERVED_SYNTAX"
	}
	for _, gap := range r.Gaps {
		if gap.Code == "GO_PARSE_FAILED" {
			r.Status = "PARTIAL_PARSE"
		}
	}
	return r, nil
}

func safeText(s string, max int) bool {
	return utf8.ValidString(s) && s != "" && len(s) <= max && strings.TrimSpace(s) == s && !strings.ContainsFunc(s, unicode.IsControl)
}

type walker struct {
	environmentSupplied bool
	typing              typeResolution
	source              SourceFile
	fs                  *token.FileSet
	report              *Report
	packageMakeShadowed bool
}

func (w walker) location(n ast.Node) Location {
	a, b := w.fs.PositionFor(n.Pos(), false), w.fs.PositionFor(n.End(), false)
	return Location{Path: w.source.Path, Line: a.Line, Column: a.Column, Offset: a.Offset, EndOffset: b.Offset, ContentSHA256: w.source.SHA256}
}

func (w walker) text(n ast.Node) string {
	l := w.location(n)
	s := strings.Join(strings.Fields(w.source.Content[l.Offset:l.EndOffset]), " ")
	if len([]rune(s)) > 200 {
		return string([]rune(s)[:200]) + "…"
	}
	return s
}

func (w walker) add(n ast.Node, fn, property, value, resource string, conditions []string) {
	cs := append([]string{}, conditions...)
	w.report.Observations = append(w.report.Observations, Observation{Property: property, Value: value, Resource: resource, Function: fn, Conditions: cs, Authority: "SYNTAX_ONLY", Location: w.location(n)})
}

func (w walker) gap(n ast.Node, fn, code string) {
	w.report.Gaps = append(w.report.Gaps, Gap{Code: code, Function: fn, Location: w.location(n)})
}

func (w walker) function(body *ast.BlockStmt, fn string) {
	before := len(w.report.Observations)
	w.walk(body, fn, nil, "")
	// Nested function observations do not imply a property of the enclosing function.
	for _, observation := range w.report.Observations[before:] {
		if observation.Function == fn {
			for _, code := range []string{"WORKER_BOUND_UNRESOLVED", "ADMISSION_BOUND_UNRESOLVED", "CANCELLATION_COMPLETION_UNRESOLVED"} {
				w.gap(body, fn, code)
			}
			break
		}
	}
}

func withCondition(cs []string, condition string) []string {
	return append(append([]string{}, cs...), condition)
}

func (w walker) walk(n ast.Node, fn string, cs []string, communication string) {
	if n == nil {
		return
	}
	switch x := n.(type) {
	case *ast.FuncLit:
		w.function(x.Body, fn+"/literal@"+strconv.Itoa(w.location(x).Offset))
		return
	case *ast.IfStmt:
		w.walk(x.Init, fn, cs, "")
		w.walk(x.Cond, fn, cs, "")
		w.walk(x.Body, fn, withCondition(cs, "if("+w.text(x.Cond)+")"), "")
		w.walk(x.Else, fn, withCondition(cs, "else("+w.text(x.Cond)+")"), "")
		return
	case *ast.ForStmt:
		w.walk(x.Init, fn, cs, "")
		label := "for(;;)"
		if x.Cond != nil {
			label = "for(" + w.text(x.Cond) + ")"
		}
		next := withCondition(cs, label)
		w.walk(x.Cond, fn, next, "")
		w.walk(x.Post, fn, next, "")
		w.walk(x.Body, fn, next, "")
		return
	case *ast.RangeStmt:
		w.walk(x.X, fn, cs, "")
		w.walk(x.Body, fn, withCondition(cs, "range("+w.text(x.X)+")"), "")
		return
	case *ast.SelectStmt:
		kind := "WAITING_SELECT"
		hasDefault := false
		for _, node := range x.Body.List {
			if node.(*ast.CommClause).Comm == nil {
				hasDefault = true
			}
		}
		if hasDefault {
			kind = "SELECT_WITH_DEFAULT"
		}
		w.add(x, fn, "select_shape", kind, "", cs)
		for i, node := range x.Body.List {
			clause := node.(*ast.CommClause)
			next := withCondition(cs, "select.case:"+strconv.Itoa(i))
			w.walk(clause.Comm, fn, next, kind)
			for _, statement := range clause.Body {
				w.walk(statement, fn, next, "")
			}
		}
		return
	case *ast.CaseClause:
		for _, e := range x.List {
			w.walk(e, fn, cs, "")
		}
		label := "case:default"
		if len(x.List) > 0 {
			label = "case(" + w.text(x.List[0]) + ")"
		}
		for _, s := range x.Body {
			w.walk(s, fn, withCondition(cs, label), "")
		}
		return
	case *ast.GoStmt:
		w.add(x, fn, "goroutine_launch", "GO_STATEMENT", w.text(x.Call.Fun), cs)
	case *ast.SendStmt:
		kind := "DIRECT_SEND"
		if communication != "" {
			kind = communication
		}
		w.add(x, fn, "channel_send", kind, w.text(x.Chan), cs)
	case *ast.UnaryExpr:
		if x.Op == token.ARROW {
			kind := "DIRECT_RECEIVE"
			if communication != "" {
				kind = communication
			}
			w.add(x, fn, "channel_receive", kind, w.text(x.X), cs)
			if call, ok := x.X.(*ast.CallExpr); ok {
				if selector, ok := call.Fun.(*ast.SelectorExpr); ok && selector.Sel.Name == "Done" && len(call.Args) == 0 {
					w.add(x, fn, "cancellation_candidate", "DONE_SELECTOR_RECEIVE", w.text(selector.X), cs)
					w.gap(x, fn, "DONE_RECEIVER_TYPE_UNRESOLVED")
				}
			}
		}
	case *ast.CallExpr:
		if id, ok := x.Fun.(*ast.Ident); ok && id.Name == "make" && len(x.Args) > 0 {
			_, literalChannel := x.Args[0].(*ast.ChanType)
			resolved, knownOther := w.typing.channelMake(x, id)
			if !literalChannel && !resolved && !knownOther {
				w.gap(x, fn, "NAMED_CHANNEL_TYPE_UNRESOLVED")
			}
			if literalChannel || resolved {
				if id.Obj != nil || w.packageMakeShadowed {
					w.gap(x, fn, "MAKE_NAME_SHADOWED")
				} else {
					value := "UNBUFFERED_FORM"
					if len(x.Args) > 2 {
						w.gap(x, fn, "INVALID_CHANNEL_MAKE_ARITY")
						break
					}
					if len(x.Args) == 2 {
						value = "CAPACITY_EXPRESSION"
						if literal, ok := x.Args[1].(*ast.BasicLit); ok && literal.Kind == token.INT {
							if capacity, err := strconv.ParseUint(strings.ReplaceAll(literal.Value, "_", ""), 0, 64); err == nil {
								value = "CAPACITY_LITERAL:" + strconv.FormatUint(capacity, 10)
							}
						}
						w.gap(x, fn, "CAPACITY_VALIDITY_NOT_TYPECHECKED")
					}
					w.add(x, fn, "channel_buffer_syntax", value, w.text(x.Args[0]), cs)
					if !literalChannel && resolved {
						w.add(x, fn, "channel_type_identity", "RESOLVED_CHANNEL_TYPE", w.text(x.Args[0]), cs)
						w.report.Observations[len(w.report.Observations)-1].Authority = "GO_TYPES_SUPPLIED_PACKAGE"
						if w.environmentSupplied {
							w.report.Observations[len(w.report.Observations)-1].Authority = "GO_TYPES_SUPPLIED_ENVIRONMENT"
						}
					}
				}
			}
		}
	}
	// Walk only direct children; custom branches above supply their own conditions.
	first := true
	ast.Inspect(n, func(child ast.Node) bool {
		if first {
			first = false
			return true
		}
		if child != nil {
			w.walk(child, fn, cs, communication)
		}
		return false
	})
}
