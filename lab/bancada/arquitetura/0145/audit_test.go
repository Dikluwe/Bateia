// Ownership controls belong to the Lab, never to the composition root.
package architecture

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var owners = map[string]string{
	"ExpandInsightNeighborhood":         "01_core/usecase/insight_neighborhood.go",
	"ValidateInsightNeighborhood":       "01_core/usecase/insight_neighborhood_plan.go",
	"admitInsightNeighborhood":          "01_core/usecase/insight_neighborhood_plan.go",
	"neighborhoodSiteRequest":           "01_core/usecase/insight_neighborhood_plan.go",
	"insightNeighborhoodOrigins":        "01_core/usecase/insight_neighborhood_sources.go",
	"newInsightNeighborhoodReport":      "01_core/usecase/insight_neighborhood_result.go",
	"prepareInsightNeighborhoodContext": "01_core/usecase/insight_neighborhood_result.go",
	"finishInsightNeighborhood":         "01_core/usecase/insight_neighborhood_result.go",
	"completeInsightNeighborhood":       "01_core/usecase/insight_neighborhood_result.go",
	"prepareInsightNeighborhoodLeads":   "01_core/usecase/insight_neighborhood_frontier.go",
	"insightNeighborhoodDistances":      "01_core/usecase/insight_neighborhood_frontier.go",
	"insightNeighborhoodFrontier":       "01_core/usecase/insight_neighborhood_frontier.go",
	"beginInsightExplorationMemory":     "01_core/usecase/insight_exploration_memory.go",
	"insightMemoryReadState":            "01_core/usecase/insight_exploration_memory.go",
	"rememberInsightReferenceLeads":     "01_core/usecase/insight_exploration_memory.go",
	"sortInsightFollowupQueue":          "01_core/usecase/insight_followup_selection.go",
	"prepareInsightFollowupContext":     "01_core/usecase/insight_followup_context.go",
	"runInsightFollowupQueue":           "01_core/usecase/insight_followups.go",
	"runInsightFollowups":               "01_core/usecase/insight_followups.go",
	"WebTokenRequested":                 "02_shell/cli/web_input.go",
}

func ownershipIssues(name string, raw []byte) ([]string, map[string]bool, error) {
	file, err := parser.ParseFile(token.NewFileSet(), name, raw, 0)
	if err != nil {
		return nil, nil, err
	}
	issues := []string{}
	found := map[string]bool{}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if owner, known := owners[fn.Name.Name]; known {
			found[fn.Name.Name] = true
			if name != owner {
				issues = append(issues, "misplaced "+fn.Name.Name)
			}
		}
		if name == "01_core/usecase/insight_neighborhood.go" && fn.Name.Name != "ExpandInsightNeighborhood" {
			issues = append(issues, "facade owns an extra function")
		}
	}
	ast.Inspect(file, func(node ast.Node) bool {
		if literal, ok := node.(*ast.BasicLit); ok && literal.Kind == token.STRING && strings.HasPrefix(name, "04_wiring/") {
			text, _ := strconv.Unquote(literal.Value)
			if text == "--no-token" {
				issues = append(issues, "argument interpretation retained in L4")
			}
		}
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		method := ""
		switch x := call.Fun.(type) {
		case *ast.Ident:
			method = x.Name
		case *ast.SelectorExpr:
			method = x.Sel.Name
		}
		denied := false
		switch name {
		case "01_core/usecase/insight_neighborhood.go":
			denied = method == "insightNeighborhoodOrigins" || method == "Marshal" || method == "Unmarshal" || method == "finishInsightNeighborhood"
		case "01_core/usecase/insight_neighborhood_sources.go":
			denied = method == "insightNeighborhoodFrontier" || method == "runInsightFollowupQueue" || method == "Discover" || method == "FetchTree"
		case "01_core/usecase/insight_neighborhood_plan.go":
			denied = method == "runInsightFollowupQueue" || method == "ResolveRepository"
		case "01_core/usecase/insight_followup_selection.go", "01_core/usecase/insight_followup_context.go", "01_core/usecase/insight_neighborhood_result.go":
			denied = method == "Discover" || method == "ResolveRepository" || method == "FetchTree"
		}
		if denied {
			issues = append(issues, "operation outside owner: "+method)
		}
		return true
	})
	return issues, found, nil
}

func repository(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("../../../..")
	if err != nil {
		t.Fatal(err)
	}
	return root
}
func TestInvestigationOwnershipBoundaries(t *testing.T) {
	root := repository(t)
	seen := map[string]bool{}
	for _, layer := range []string{"01_core", "02_shell", "03_infra", "04_wiring"} {
		err := filepath.WalkDir(filepath.Join(root, layer), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			relative, _ := filepath.Rel(root, path)
			issues, found, err := ownershipIssues(filepath.ToSlash(relative), raw)
			if err != nil {
				return err
			}
			for key := range found {
				if seen[key] {
					t.Error("duplicate owner", key)
				}
				seen[key] = true
			}
			for _, issue := range issues {
				t.Error(relative, issue)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	for name := range owners {
		if !seen[name] {
			t.Error("missing owner", name)
		}
	}
}
func TestOwnershipControlsRejectRegressions(t *testing.T) {
	for _, test := range []struct {
		name, source string
		bad          bool
	}{
		{"04_wiring/example.go", `package x;func f(){_="--no-token"}`, true},
		{"01_core/usecase/insight_neighborhood.go", `package x;func finishInsightNeighborhood(){}`, true},
		{"01_core/usecase/insight_neighborhood.go", `package x;func ExpandInsightNeighborhood(){insightNeighborhoodOrigins(input)}`, true},
		{"01_core/usecase/insight_neighborhood_sources.go", `package x;func insightNeighborhoodOrigins(){provider.FetchTree(ctx,repo)}`, true},
		{"01_core/usecase/insight_neighborhood_result.go", `package x;func finishInsightNeighborhood(){reader.ResolveRepository(ctx,id)}`, true},
		{"04_wiring/example.go", `package x;func f(){cli.WebTokenRequested(args)}`, false},
		{"01_core/usecase/insight_neighborhood_plan.go", `package x;func admitInsightNeighborhood(){insightNeighborhoodOrigins(input)}`, false},
	} {
		issues, _, err := ownershipIssues(test.name, []byte(test.source))
		if err != nil || (len(issues) > 0) != test.bad {
			t.Fatalf("incorrect control %s: %v %v", test.name, issues, err)
		}
	}
}

func declarations(raw []byte) (map[string]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), "source.go", raw, 0)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil {
			continue
		}
		var out bytes.Buffer
		if err := format.Node(&out, token.NewFileSet(), fn.Body); err != nil {
			return nil, err
		}
		result[fn.Name.Name] = fmt.Sprintf("%x", sha256.Sum256(out.Bytes()))
	}
	return result, nil
}
func currentDeclarations(t *testing.T, root string, tests bool) map[string]string {
	t.Helper()
	result := map[string]string{}
	files, err := filepath.Glob(filepath.Join(root, "01_core/usecase/*.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") != tests {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		items, err := declarations(raw)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range items {
			result[k] = v
		}
	}
	return result
}
func TestExistingConsumersAndMovedRulesRemainUnchanged(t *testing.T) {
	root := repository(t)
	currentTests := currentDeclarations(t, root, true)
	currentRules := currentDeclarations(t, root, false)
	for _, name := range []string{"insight_neighborhood_test", "insight_followups_test"} {
		raw, err := os.ReadFile(filepath.Join(root, "lab/bancada/arquitetura/0145/before/01_core__usecase__"+name+".go.txt"))
		if err != nil {
			t.Fatal(err)
		}
		old, err := declarations(raw)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range old {
			if strings.HasPrefix(k, "Test") && currentTests[k] != v {
				t.Error("consumer changed or disappeared", k)
			}
		}
	}
	for _, test := range []struct {
		file  string
		names []string
	}{
		{"insight_neighborhood", []string{"neighborhoodSiteRequest", "finishInsightNeighborhood"}},
		{"insight_neighborhood_frontier", []string{"insightNeighborhoodOrigins"}},
		{"insight_followups", []string{"beginInsightExplorationMemory", "insightMemoryReadState", "rememberInsightReferenceLeads", "sortInsightFollowupQueue"}},
	} {
		raw, err := os.ReadFile(filepath.Join(root, "lab/bancada/arquitetura/0145/before/01_core__usecase__"+test.file+".go.txt"))
		if err != nil {
			t.Fatal(err)
		}
		old, err := declarations(raw)
		if err != nil {
			t.Fatal(err)
		}
		for _, name := range test.names {
			if old[name] == "" || currentRules[name] != old[name] {
				t.Error("moved rule changed", name)
			}
		}
	}
}
