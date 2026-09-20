/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/durableconfig/config_security_contract_test.md
 * @layer L3
 * @updated 2026-09-18
 */

package durableconfig_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type parsedSecuritySource struct {
	file          *ast.File
	importAliases map[string]string
}

type securitySourceAudit struct {
	hasNoFollowOpen bool
	chmod0600       []token.Position
	unorderedChmod  []token.Position
	uncheckedReads  []token.Position
}

func TestSecuritySourceContract(t *testing.T) {
	fileSet, sources := parseProductionSources(t)
	audit := auditSecuritySources(fileSet, sources)

	if !audit.hasNoFollowOpen {
		t.Error("cursor opening does not combine an open operation with syscall.O_NOFOLLOW")
	}
	if len(audit.chmod0600) == 0 {
		t.Error("no cursor mode normalization to 0600 was found")
	}
	for _, position := range audit.unorderedChmod {
		t.Errorf("Chmod(0600) at %s is not preceded in its function by an exact Nlink == 1/!= 1 validation", position)
	}
	for _, position := range audit.uncheckedReads {
		t.Errorf("plain Read at %s does not check both the returned byte count and error; use io.ReadFull or validate both results", position)
	}
}

func TestSecurityAnalyzerDiscriminatesSyntheticMutants(t *testing.T) {
	secure := `package durableconfig
import "syscall"
func load() {
	_, _ = syscall.Open(path, syscall.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if metadata.Nlink != 1 { return }
	_ = cursor.Chmod(0o600)
	count, err := cursor.Read(secret)
	if err != nil { return }
	if count != len(secret) { return }
}`
	withoutNoFollow := strings.Replace(secure, "syscall.O_RDONLY|syscall.O_NOFOLLOW", "syscall.O_RDONLY", 1)
	chmodBeforeLinkValidation := strings.Replace(
		secure,
		"if metadata.Nlink != 1 { return }\n\t_ = cursor.Chmod(0o600)",
		"_ = cursor.Chmod(0o600)\n\tif metadata.Nlink != 1 { return }",
		1,
	)
	ignoredReadCount := strings.Replace(
		secure,
		"count, err := cursor.Read(secret)\n\tif err != nil { return }\n\tif count != len(secret) { return }",
		"_, err := cursor.Read(secret)\n\tif err != nil { return }",
		1,
	)

	secureAudit := auditSyntheticSource(t, secure)
	if !secureAudit.hasNoFollowOpen || len(secureAudit.chmod0600) != 1 || len(secureAudit.unorderedChmod) != 0 || len(secureAudit.uncheckedReads) != 0 {
		t.Fatalf("secure control was not accepted: %+v", secureAudit)
	}

	testCases := []struct {
		name       string
		source     string
		rejected   func(securitySourceAudit) bool
		diagnostic string
	}{
		{
			name:       "removed-O_NOFOLLOW",
			source:     withoutNoFollow,
			rejected:   func(audit securitySourceAudit) bool { return !audit.hasNoFollowOpen },
			diagnostic: "missing O_NOFOLLOW was not rejected",
		},
		{
			name:       "chmod-before-Nlink",
			source:     chmodBeforeLinkValidation,
			rejected:   func(audit securitySourceAudit) bool { return len(audit.unorderedChmod) == 1 },
			diagnostic: "Chmod before Nlink validation was not rejected",
		},
		{
			name:       "ignored-Read-count",
			source:     ignoredReadCount,
			rejected:   func(audit securitySourceAudit) bool { return len(audit.uncheckedReads) == 1 },
			diagnostic: "ignored Read count was not rejected",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			audit := auditSyntheticSource(t, testCase.source)
			if !testCase.rejected(audit) {
				t.Fatal(testCase.diagnostic)
			}
		})
	}
}

func parseProductionSources(t *testing.T) (*token.FileSet, []parsedSecuritySource) {
	t.Helper()
	_, testPath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate security contract test source")
	}
	packageDirectory := filepath.Dir(testPath)
	entries, err := os.ReadDir(packageDirectory)
	if err != nil {
		t.Fatalf("read durableconfig package directory: %v", err)
	}

	fileSet := token.NewFileSet()
	var sources []parsedSecuritySource
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(packageDirectory, name)
		file, parseErr := parser.ParseFile(fileSet, path, nil, 0)
		if parseErr != nil {
			t.Fatalf("parse production source %s: %v", name, parseErr)
		}
		sources = append(sources, parsedSecuritySource{
			file:          file,
			importAliases: sourceImportAliases(file),
		})
	}
	if len(sources) == 0 {
		t.Fatal("no non-test Go source exists in the durableconfig package")
	}
	return fileSet, sources
}

func auditSyntheticSource(t *testing.T, source string) securitySourceAudit {
	t.Helper()
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "synthetic.go", source, 0)
	if err != nil {
		t.Fatalf("parse synthetic source: %v", err)
	}
	return auditSecuritySources(fileSet, []parsedSecuritySource{{
		file:          file,
		importAliases: sourceImportAliases(file),
	}})
}

func auditSecuritySources(fileSet *token.FileSet, sources []parsedSecuritySource) securitySourceAudit {
	var audit securitySourceAudit
	for _, source := range sources {
		modeConstants := integerConstants(source.file)
		for _, declaration := range source.file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Body == nil {
				continue
			}
			noFollowNames := assignedNoFollowNames(function.Body, source.importAliases)
			linkCountNames := assignedLinkCountNames(function.Body)
			linkValidations := exactLinkValidationPositions(function.Body, linkCountNames)
			parents := parentNodes(function.Body)

			ast.Inspect(function.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				if isOpenCall(call) && callUsesNoFollow(call, source.importAliases, noFollowNames) {
					audit.hasNoFollowOpen = true
				}
				if isChmod0600(call, modeConstants) {
					position := fileSet.Position(call.Pos())
					audit.chmod0600 = append(audit.chmod0600, position)
					if !hasEarlierPosition(linkValidations, call.Pos()) {
						audit.unorderedChmod = append(audit.unorderedChmod, position)
					}
				}
				if isPlainMethodRead(call, source.importAliases) && !plainReadChecksCountAndError(function, call, parents) {
					audit.uncheckedReads = append(audit.uncheckedReads, fileSet.Position(call.Pos()))
				}
				return true
			})
		}
	}
	return audit
}

func sourceImportAliases(file *ast.File) map[string]string {
	aliases := make(map[string]string)
	for _, imported := range file.Imports {
		path, err := strconv.Unquote(imported.Path.Value)
		if err != nil {
			continue
		}
		alias := filepath.Base(path)
		if imported.Name != nil {
			alias = imported.Name.Name
		}
		aliases[alias] = path
	}
	return aliases
}

func integerConstants(file *ast.File) map[string]int64 {
	constants := make(map[string]int64)
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.CONST {
			continue
		}
		for _, specification := range general.Specs {
			values, ok := specification.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for index, name := range values.Names {
				if index >= len(values.Values) {
					continue
				}
				if value, ok := integerExpression(values.Values[index], constants); ok {
					constants[name.Name] = value
				}
			}
		}
	}
	return constants
}

func assignedNoFollowNames(body *ast.BlockStmt, importAliases map[string]string) map[string]bool {
	names := make(map[string]bool)
	ast.Inspect(body, func(node ast.Node) bool {
		assignment, ok := node.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for index, right := range assignment.Rhs {
			if index >= len(assignment.Lhs) || !expressionUsesNoFollow(right, importAliases, names) {
				continue
			}
			if identifier, ok := assignment.Lhs[index].(*ast.Ident); ok {
				names[identifier.Name] = true
			}
		}
		return true
	})
	return names
}

func assignedLinkCountNames(body *ast.BlockStmt) map[string]bool {
	names := make(map[string]bool)
	ast.Inspect(body, func(node ast.Node) bool {
		assignment, ok := node.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for index, right := range assignment.Rhs {
			if index >= len(assignment.Lhs) || !expressionContainsSelector(right, "Nlink") {
				continue
			}
			if identifier, ok := assignment.Lhs[index].(*ast.Ident); ok {
				names[identifier.Name] = true
			}
		}
		return true
	})
	return names
}

func exactLinkValidationPositions(body *ast.BlockStmt, linkCountNames map[string]bool) []token.Pos {
	var positions []token.Pos
	ast.Inspect(body, func(node ast.Node) bool {
		conditional, ok := node.(*ast.IfStmt)
		if ok && conditionValidatesOneLink(conditional.Cond, linkCountNames) {
			positions = append(positions, conditional.Cond.Pos())
		}
		return true
	})
	return positions
}

func conditionValidatesOneLink(expression ast.Expr, linkCountNames map[string]bool) bool {
	validated := false
	ast.Inspect(expression, func(node ast.Node) bool {
		binary, ok := node.(*ast.BinaryExpr)
		if !ok || (binary.Op != token.EQL && binary.Op != token.NEQ) {
			return true
		}
		leftIsLink := expressionContainsSelector(binary.X, "Nlink") || expressionContainsNamedIdentifier(binary.X, linkCountNames)
		rightIsLink := expressionContainsSelector(binary.Y, "Nlink") || expressionContainsNamedIdentifier(binary.Y, linkCountNames)
		left, leftIsOne := integerExpression(binary.X, nil)
		right, rightIsOne := integerExpression(binary.Y, nil)
		if (leftIsLink && rightIsOne && right == 1) || (rightIsLink && leftIsOne && left == 1) {
			validated = true
			return false
		}
		return true
	})
	return validated
}

func isOpenCall(call *ast.CallExpr) bool {
	name := calledName(call)
	return name == "Open" || name == "OpenFile" || name == "Openat" || name == "Openat2"
}

func callUsesNoFollow(call *ast.CallExpr, importAliases map[string]string, assignedNames map[string]bool) bool {
	for _, argument := range call.Args {
		if expressionUsesNoFollow(argument, importAliases, assignedNames) {
			return true
		}
	}
	return false
}

func expressionUsesNoFollow(expression ast.Expr, importAliases map[string]string, assignedNames map[string]bool) bool {
	found := false
	ast.Inspect(expression, func(node ast.Node) bool {
		switch current := node.(type) {
		case *ast.SelectorExpr:
			identifier, ok := current.X.(*ast.Ident)
			if ok && current.Sel.Name == "O_NOFOLLOW" && importAliases[identifier.Name] == "syscall" {
				found = true
				return false
			}
		case *ast.Ident:
			if assignedNames[current.Name] {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func isChmod0600(call *ast.CallExpr, constants map[string]int64) bool {
	name := calledName(call)
	if name != "Chmod" && name != "Fchmod" && name != "Fchmodat" {
		return false
	}
	for _, argument := range call.Args {
		if value, ok := integerExpression(argument, constants); ok && value == 0o600 {
			return true
		}
	}
	return false
}

func integerExpression(expression ast.Expr, constants map[string]int64) (int64, bool) {
	switch current := expression.(type) {
	case *ast.BasicLit:
		if current.Kind != token.INT {
			return 0, false
		}
		value, err := strconv.ParseInt(current.Value, 0, 64)
		return value, err == nil
	case *ast.ParenExpr:
		return integerExpression(current.X, constants)
	case *ast.UnaryExpr:
		value, ok := integerExpression(current.X, constants)
		if !ok {
			return 0, false
		}
		if current.Op == token.SUB {
			return -value, true
		}
		if current.Op == token.ADD {
			return value, true
		}
	case *ast.CallExpr:
		if len(current.Args) == 1 {
			return integerExpression(current.Args[0], constants)
		}
	case *ast.Ident:
		if constants != nil {
			value, ok := constants[current.Name]
			return value, ok
		}
	}
	return 0, false
}

func isPlainMethodRead(call *ast.CallExpr, importAliases map[string]string) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Read" {
		return false
	}
	if identifier, ok := selector.X.(*ast.Ident); ok {
		_, importedPackage := importAliases[identifier.Name]
		return !importedPackage
	}
	return true
}

func plainReadChecksCountAndError(function *ast.FuncDecl, call *ast.CallExpr, parents map[ast.Node]ast.Node) bool {
	assignment, ok := parents[call].(*ast.AssignStmt)
	if !ok || len(assignment.Rhs) != 1 || assignment.Rhs[0] != call || len(assignment.Lhs) < 2 {
		return false
	}
	count, countOK := assignment.Lhs[0].(*ast.Ident)
	readError, errorOK := assignment.Lhs[1].(*ast.Ident)
	if !countOK || !errorOK || count.Name == "_" || readError.Name == "_" {
		return false
	}
	return identifierUsedInLaterCondition(function.Body, count.Name, call.End()) &&
		identifierUsedInLaterCondition(function.Body, readError.Name, call.End())
}

func identifierUsedInLaterCondition(body *ast.BlockStmt, name string, after token.Pos) bool {
	used := false
	ast.Inspect(body, func(node ast.Node) bool {
		conditional, ok := node.(*ast.IfStmt)
		if !ok || conditional.Cond.Pos() <= after {
			return true
		}
		if expressionContainsIdentifier(conditional.Cond, name) {
			used = true
			return false
		}
		return true
	})
	return used
}

func expressionContainsSelector(expression ast.Expr, selectorName string) bool {
	found := false
	ast.Inspect(expression, func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if ok && selector.Sel.Name == selectorName {
			found = true
			return false
		}
		return true
	})
	return found
}

func expressionContainsNamedIdentifier(expression ast.Expr, names map[string]bool) bool {
	found := false
	ast.Inspect(expression, func(node ast.Node) bool {
		identifier, ok := node.(*ast.Ident)
		if ok && names[identifier.Name] {
			found = true
			return false
		}
		return true
	})
	return found
}

func expressionContainsIdentifier(expression ast.Expr, name string) bool {
	found := false
	ast.Inspect(expression, func(node ast.Node) bool {
		identifier, ok := node.(*ast.Ident)
		if ok && identifier.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}

func calledName(call *ast.CallExpr) string {
	switch function := call.Fun.(type) {
	case *ast.Ident:
		return function.Name
	case *ast.SelectorExpr:
		return function.Sel.Name
	default:
		return ""
	}
}

func hasEarlierPosition(positions []token.Pos, target token.Pos) bool {
	for _, position := range positions {
		if position < target {
			return true
		}
	}
	return false
}

func parentNodes(root ast.Node) map[ast.Node]ast.Node {
	parents := make(map[ast.Node]ast.Node)
	var stack []ast.Node
	ast.Inspect(root, func(node ast.Node) bool {
		if node == nil {
			stack = stack[:len(stack)-1]
			return false
		}
		if len(stack) > 0 {
			parents[node] = stack[len(stack)-1]
		}
		stack = append(stack, node)
		return true
	})
	return parents
}
