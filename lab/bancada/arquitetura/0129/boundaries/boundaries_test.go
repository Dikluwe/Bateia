// Architectural boundary verifier, step 0129. Lab only; never a product dependency.
package boundaries

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var layers = []string{"00_nucleo", "01_core", "02_shell", "03_infra", "04_wiring"}

func inspectSource(name string, layer int, source []byte) ([]string, error) {
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, name, source, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var issues []string
	report := func(node ast.Node, message string) {
		issues = append(issues, fmt.Sprintf("%s:%d: %s", name, set.Position(node.Pos()).Line, message))
	}
	imports := map[string]string{}
	for _, imp := range file.Imports {
		path, _ := strconv.Unquote(imp.Path.Value)
		alias := filepath.Base(path)
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		imports[alias] = path
		if alias == "." {
			report(imp, "dot import prevents explicit boundary attribution")
		}
		for target, folder := range layers {
			if !strings.HasPrefix(path, "bateia/"+folder+"/") {
				continue
			}
			forbidden := layer == 1 && (target == 2 || target == 3 || target == 4) || layer == 2 && (target == 3 || target == 4) || layer == 3 && (target == 2 || target == 4)
			if forbidden {
				report(imp, "gravity inversion: "+path)
			}
		}
		if strings.HasPrefix(path, "bateia/lab/") {
			report(imp, "product imports Lab")
		}
		if layer == 1 || layer == 2 {
			if path == "os" && layer == 1 || path == "syscall" || path == "io/ioutil" || path == "os/exec" || path == "os/signal" || path == "crypto/rand" || path == "net" || strings.HasPrefix(path, "net/") && path != "net/url" {
				report(imp, "concrete effect import: "+path)
			}
		}
	}
	ast.Inspect(file, func(n ast.Node) bool {
		selector, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}
		// Locally shadowed names are not package references.
		if ident.Obj != nil {
			return true
		}
		path := imports[ident.Name]
		member := selector.Sel.Name
		if layer == 1 || layer == 2 {
			prohibited := false
			switch path {
			case "time":
				prohibited = contains([]string{"Now", "Since", "Until", "After", "AfterFunc", "Tick", "NewTimer", "NewTicker", "Sleep"}, member)
			case "context":
				prohibited = strings.HasPrefix(member, "WithTimeout") || strings.HasPrefix(member, "WithDeadline")
			case "path/filepath":
				prohibited = contains([]string{"Abs", "Glob", "Walk", "WalkDir", "EvalSymlinks"}, member)
			case "os":
				prohibited = !(layer == 2 && name == "02_shell/cli/app.go" && contains([]string{"Stdin", "Stdout", "Stderr"}, member))
			case "fmt":
				prohibited = layer == 1 && contains([]string{"Print", "Println", "Printf", "Scan", "Scanf", "Scanln"}, member)
			case "log":
				prohibited = layer == 1
			}
			if prohibited {
				report(selector, "concrete effect: "+path+"."+member)
			}
		}
		return true
	})
	if layer == 4 {
		for _, decl := range file.Decls {
			if method, ok := decl.(*ast.FuncDecl); ok && method.Recv != nil {
				report(method, "method implementation belongs to a capability")
			}
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typ := spec.(*ast.TypeSpec)
				if _, ok := typ.Type.(*ast.InterfaceType); ok || typ.Assign.IsValid() {
					report(typ, "contract type belongs in L1/L2")
				}
			}
		}
	}
	return issues, nil
}
func contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

func TestProductBoundaries(t *testing.T) {
	root, err := filepath.Abs("../../../../..")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for layer, dir := range layers[1:] {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			name, _ := filepath.Rel(root, path)
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			issues, err := inspectSource(filepath.ToSlash(name), layer+1, data)
			if err != nil {
				return err
			}
			for _, issue := range issues {
				t.Error(issue)
			}
			count++
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("audited %d productive Go files", count)
	err = filepath.WalkDir(filepath.Join(root, layers[0]), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		lower := strings.ToLower(entry.Name())
		extension := filepath.Ext(lower)
		forbidden := contains([]string{".go", ".py", ".sh", ".js", ".ts", ".html", ".wasm"}, extension) || lower == "go.mod" || info.Mode()&0111 != 0
		for _, suffix := range []string{".go.txt", ".py.txt", ".sh.txt", ".js.txt"} {
			forbidden = forbidden || strings.HasSuffix(lower, suffix)
		}
		if forbidden {
			t.Errorf("executable material in L0: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
func TestVerifierRejectsConcreteEffectsAndGravity(t *testing.T) {
	cases := []struct {
		name   string
		layer  int
		source string
		bad    bool
	}{
		{"alias", 1, `package x; import clock "time"; func f(){_ = clock.Now()}`, true},
		{"value alias", 1, `package x; import "time"; var clock = time.Now`, true},
		{"deadline", 1, `package x; import "context"; var deadline = context.WithTimeoutCause`, true},
		{"disk", 2, `package x; import fs "os"; var read = fs.ReadFile`, true},
		{"upward", 2, `package x; import _ "bateia/03_infra/host"`, true},
		{"lab", 3, `package x; import _ "bateia/lab/arena"`, true},
		{"method", 4, `package x; type value struct{};func(value) Decide(){}`, true},
		{"contract", 4, `package x;type Policy interface{Decide()}`, true},
		{"factory", 4, `package x;func Build()any{return nil}`, false},
		{"data", 1, `package x;import "time";type Record struct{At time.Time}`, false},
		{"port", 1, `package x;import "io/fs";type Store interface{Stat(string)(fs.FileInfo,error)}`, false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			issues, err := inspectSource(test.name, test.layer, []byte(test.source))
			if err != nil {
				t.Fatal(err)
			}
			if (len(issues) > 0) != test.bad {
				t.Fatalf("issues=%v bad=%v", issues, test.bad)
			}
		})
	}
}
