/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/client_test.md
 * @layer L3
 * @updated 2026-09-07
 */
package clang_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/clang"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCompilationBudgetKeepsGlobalCoveragePartial0035(t *testing.T) {
	binary, err := exec.LookPath("clang")
	if err != nil {
		t.Skip("clang unavailable")
	}
	root := t.TempDir()
	var entries []map[string]any
	for i := 0; i < 13; i++ {
		path := filepath.Join(root, fmt.Sprintf("render%02d.cpp", i))
		if err := os.WriteFile(path, []byte("int render() {return 1;}\n"), 0600); err != nil {
			t.Fatal(err)
		}
		entries = append(entries, map[string]any{"directory": root, "file": path, "arguments": []string{"clang++", "-c", path}})
	}
	data, _ := json.Marshal(entries)
	if err := os.WriteFile(filepath.Join(root, "compile_commands.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	report, err := clang.NewClient(binary).Analyze(context.Background(), domain.SemanticRequest{Root: root, Query: "render"})
	if err != nil || report.Status != "partial" || report.TotalUnits != 13 || report.SelectedUnits != 12 || len(report.Units) != 12 {
		t.Fatalf("budget not observable: %+v %v", report, err)
	}
}
