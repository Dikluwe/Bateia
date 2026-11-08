// Contract-specific ownership and artifact-regime checks. No product dependency.
package contracts

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func ownerIssues(name string, source []byte) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), name, source, 0)
	if err != nil {
		return nil, err
	}
	var issues []string
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		nameNode := call.Fun
		method := ""
		if selector, ok := nameNode.(*ast.SelectorExpr); ok {
			method = selector.Sel.Name
		} else if id, ok := nameNode.(*ast.Ident); ok {
			method = id.Name
		}
		denied := false
		switch name {
		case "02_shell/cli/idea_recommend.go":
			denied = strings.HasPrefix(method, "NewRecommendationService") || method == "RunIdeaPipeline" || method == "RankSoftwareAlternatives" || method == "StableRunID"
		case "02_shell/cli/idea_web.go":
			denied = method == "SearchOne"
		case "01_core/usecase/market_discovery.go":
			denied = method == "PlanMarketDiscovery" || method == "AllocateMarketDiscoveryBudget" || method == "Discover" || method == "marketStatusFor"
		}
		if strings.HasPrefix(name, "04_wiring/") && !strings.HasSuffix(name, "_test.go") {
			for _, business := range []string{"ExtractMechanisms", "RecommendWithRun", "RunIdeaPipeline", "BuildDecisionTrace", "PlanMarketDiscovery", "Transition", "PublishArtifact"} {
				if method == business {
					denied = true
				}
			}
		}
		if strings.HasPrefix(name, "04_wiring/") && strings.HasSuffix(name, "_test.go") {
			denied = denied || method == "ParseFile" || strings.HasPrefix(method, "ExportedValid")
		}
		if denied {
			issues = append(issues, fmt.Sprintf("%s retains operation %s outside its owner", name, method))
		}
		return true
	})
	return issues, nil
}
func TestReviewedOwnershipBoundaries(t *testing.T) {
	root, err := filepath.Abs("../../../../..")
	if err != nil {
		t.Fatal(err)
	}
	names := []string{"02_shell/cli/idea_recommend.go", "02_shell/cli/idea_web.go", "01_core/usecase/market_discovery.go"}
	entries, err := os.ReadDir(filepath.Join(root, "04_wiring"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".go") {
			names = append(names, "04_wiring/"+entry.Name())
		}
	}
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		issues, err := ownerIssues(name, data)
		if err != nil {
			t.Fatal(err)
		}
		for _, issue := range issues {
			t.Error(issue)
		}
	}
}
func TestOwnershipChecksRejectRegressions(t *testing.T) {
	for _, test := range []struct {
		name, source string
		bad          bool
	}{
		{"02_shell/cli/idea_recommend.go", `package x;func f(){usecase.NewRecommendationService(nil,nil,nil)}`, true},
		{"02_shell/cli/idea_web.go", `package x;func f(){provider.SearchOne(ctx,req)}`, true},
		{"01_core/usecase/market_discovery.go", `package x;func f(){connector.Discover(ctx,req)}`, true},
		{"04_wiring/example_test.go", `package x;func f(){parser.ParseFile(nil,"",nil,0)}`, true},
		{"02_shell/cli/idea_recommend.go", `package x;func f(){usecase.RunRecommendationDossier(ctx,input,files)}`, false},
	} {
		issues, err := ownerIssues(test.name, []byte(test.source))
		if err != nil {
			t.Fatal(err)
		}
		if (len(issues) > 0) != test.bad {
			t.Fatalf("incorrect guard: %s %v", test.name, issues)
		}
	}
}
func TestNucleusContainsDeclarationsNotRawBenchResults(t *testing.T) {
	root, err := filepath.Abs("../../../../..")
	if err != nil {
		t.Fatal(err)
	}
	// These are declarative parameters and a causal contract-change receipt, not measured results.
	allowed := map[string]bool{"modelos/screening-operacional-v1.json": true, "relatorios/evidence-contract-supersession.json": true}
	err = filepath.WalkDir(filepath.Join(root, "00_nucleo"), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		relative, _ := filepath.Rel(filepath.Join(root, "00_nucleo"), path)
		if !allowed[filepath.ToSlash(relative)] {
			t.Errorf("undeclared non-Markdown material in L0: %s", relative)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
