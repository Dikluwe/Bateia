// Structural oracle relocated from L4, step 0131; assertions preserved.
package contracts

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStep0026R2LayeringAndRecoveryOrderASTAudit(t *testing.T) {
	t.Chdir("../../../../../04_wiring")
	roots := []string{"../01_core/domain", "../01_core/usecase", "../02_shell/mcp", "../03_infra/runstore", "../03_infra/runworker", "."}
	parsedFiles := 0
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
			if err != nil {
				return err
			}
			parsedFiles++
			forbidden := []string{}
			switch {
			case strings.Contains(path, "/01_core/"):
				forbidden = []string{"os", "io", "path/filepath", "net/http", "database/sql", "bateia/02_shell", "bateia/03_infra", "bateia/04_wiring"}
			case strings.Contains(path, "/02_shell/"):
				forbidden = []string{"bateia/03_infra", "bateia/04_wiring", "os", "path/filepath"}
			case strings.Contains(path, "/03_infra/"):
				forbidden = []string{"bateia/02_shell", "bateia/04_wiring"}
			}
			for _, imported := range parsed.Imports {
				name := strings.Trim(imported.Path.Value, `"`)
				if (name == "io" || name == "io/fs") && declarativePortTypesOnly(parsed) {
					continue
				}
				for _, denied := range forbidden {
					if name == denied || strings.HasPrefix(name, denied+"/") {
						t.Errorf("%s imports forbidden %s", path, name)
					}
				}
			}
			if strings.Contains(path, "/01_core/") {
				ast.Inspect(parsed, func(node ast.Node) bool {
					if _, ok := node.(*ast.GoStmt); ok {
						t.Errorf("L1 starts goroutine in %s", path)
					}
					return true
				})
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if parsedFiles == 0 {
		t.Fatal("AST audit parsed no production files")
	}

	matches, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var factory *ast.FuncDecl
	for _, path := range matches {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		source, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), path, source, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, declaration := range parsed.Decls {
			if fn, ok := declaration.(*ast.FuncDecl); ok && fn.Name.Name == "NewDurableMCPServer" {
				factory = fn
			}
		}
	}
	if factory == nil {
		t.Fatal("exported observable durable factory missing")
	}
	positions := map[string]token.Pos{}
	forbiddenCalls := map[string]bool{"CanonicalRecommendationRequestHash": true, "ValidateRecommendationRunTransition": true, "Transition": true, "PublishArtifact": true}
	ast.Inspect(factory.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		name := step0026ASTName(call.Fun)
		if name == "Recover" || name == "Start" {
			positions[name] = call.Pos()
		}
		if forbiddenCalls[name] {
			t.Errorf("L4 owns business operation %s", name)
		}
		return true
	})
	if positions["Recover"] == 0 || positions["Start"] == 0 || positions["Recover"] > positions["Start"] {
		t.Fatalf("recovery must precede worker Start in exported factory: %v", positions)
	}
}

// Standard Reader/FileInfo types do not perform I/O. This exception is structural:
// executable functions, variable initializers and calls remain forbidden.
func declarativePortTypesOnly(file *ast.File) bool {
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.IMPORT && general.Tok != token.TYPE && general.Tok != token.CONST {
			return false
		}
	}
	valid := true
	ast.Inspect(file, func(node ast.Node) bool {
		if _, call := node.(*ast.CallExpr); call {
			valid = false
		}
		return valid
	})
	return valid
}

func TestDeclarativePortAuditRejectsExecutableIO(t *testing.T) {
	for _, test := range []struct {
		source  string
		allowed bool
	}{
		{`package p; import "io"; type Reader interface { Read([]byte) (int,error) }`, true},
		{`package p; import "io"; func read(){ io.ReadAll(nil) }`, false},
		{`package p; import "io"; var data = io.ReadAll(nil)`, false},
	} {
		file, err := parser.ParseFile(token.NewFileSet(), "fixture.go", test.source, 0)
		if err != nil {
			t.Fatal(err)
		}
		if got := declarativePortTypesOnly(file); got != test.allowed {
			t.Fatalf("incorrect port classification %t for %s", got, test.source)
		}
	}
}

func step0026ASTName(expression ast.Expr) string {
	switch value := expression.(type) {
	case *ast.SelectorExpr:
		return value.Sel.Name
	case *ast.Ident:
		return value.Name
	default:
		return ""
	}
}
