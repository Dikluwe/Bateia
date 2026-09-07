/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/flow_test.md
 * @layer L3
 * @updated 2026-09-07
 */
package clang

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestStructuredFlowAndPreprocessor0038(t *testing.T) {
	tool := os.Getenv("BATEIA_TEST_FLOW_TOOL")
	if tool == "" {
		t.Skip("set BATEIA_TEST_FLOW_TOOL to compiled helper")
	}
	root := t.TempDir()
	file := filepath.Join(root, "sample.cpp")
	code := "#include <stddef.h>\n#define INVOKE work(3)\nvoid work(int); void done(); void fallback();\nvoid render(bool ok,int n,void(*fn)()) {\n if (!ok) {fallback(); return;}\n for(int i=0;i<n;++i) {if(i==2)continue;work(i);if(i==4)break;}\n while(n>0 && ok) {--n;}\n fn(); done(); INVOKE;\n#if 0\n // packet decoding\n work(99);\n#endif\n}\n"
	if err := os.WriteFile(file, []byte(code), 0600); err != nil {
		t.Fatal(err)
	}
	db, _ := json.Marshal([]compilation{{Directory: root, File: file, Arguments: []string{"clang++", "-std=c++17", "-c", file}}})
	if err := os.WriteFile(filepath.Join(root, "compile_commands.json"), db, 0600); err != nil {
		t.Fatal(err)
	}
	r, err := NewClient("clang").WithFlowTool(tool).Analyze(context.Background(), domain.SemanticRequest{Root: root, Query: "render"})
	if err != nil || len(r.Units) != 1 || r.Units[0].Status != "complete" {
		t.Fatalf("AST: %+v %v", r, err)
	}
	f := r.Units[0].Flow
	if f == nil || f.Status != "complete" || len(f.Skipped) == 0 || f.ToolSHA256 == "" {
		t.Fatalf("flow: %+v", f)
	}
	seen := map[string]bool{}
	unresolved := false
	macroLocation := false
	for _, fn := range f.Functions {
		for _, b := range fn.Blocks {
			if b.Terminator != nil {
				seen[b.Terminator.Kind] = true
			}
			for _, s := range b.Statements {
				seen[s.Kind] = true
				if s.SpellingLine > 0 && s.SpellingLine != s.Line {
					macroLocation = true
				}
				if s.Resolution == "unresolved" {
					unresolved = true
				}
			}
		}
	}
	for _, kind := range []string{"IfStmt", "ForStmt", "WhileStmt", "ReturnStmt", "BinaryOperator", "BreakStmt", "ContinueStmt"} {
		if !seen[kind] {
			t.Errorf("missing %s", kind)
		}
	}
	for _, fn := range f.Functions {
		if fn.Name == "render" {
			edges := map[int][]*int{}
			fallbackBlock, doneBlock := -1, -1
			for _, b := range fn.Blocks {
				edges[b.ID] = b.Successors
				for _, s := range b.Statements {
					if s.Callee == "fallback" {
						fallbackBlock = b.ID
					}
					if s.Callee == "done" {
						doneBlock = b.ID
					}
				}
			}
			if fallbackBlock < 0 || doneBlock < 0 {
				t.Fatal("missing call sites")
			}
			queue := []int{fallbackBlock}
			visited := map[int]bool{}
			for len(queue) > 0 {
				n := queue[0]
				queue = queue[1:]
				if visited[n] {
					continue
				}
				visited[n] = true
				if n == doneBlock {
					t.Fatal("early return can reach done in CFG")
				}
				for _, next := range edges[n] {
					if next != nil {
						queue = append(queue, *next)
					}
				}
			}
		}
	}
	if !macroLocation {
		t.Fatal("macro origin/expansion not recorded")
	}
	if !unresolved {
		t.Fatal("indirect call not marked unknown")
	}
	e := parser.BuildQueryEvidence("packet decoding", &r, []domain.ExtractedFile{{Path: "sample.cpp", Content: code}}, root)
	if len(e.Definitions) > 0 {
		t.Fatal("inactive text anchored evidence")
	}
}
func TestFlowToolAbsent0038(t *testing.T) {
	f := NewClient("clang").analyzeFlow(context.Background(), compilation{}, nil)
	if f.Status != "not_evaluated" {
		t.Fatal("missing flow became evaluated")
	}
}
