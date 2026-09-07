/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/mechanism_test.md
 * @layer L3
 * @updated 2026-09-07
 */
package clang_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"bateia/03_infra/clang"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDisabledImplementationDoesNotAnswerQuery0037(t *testing.T) {
	binary, err := exec.LookPath("clang")
	if err != nil {
		t.Skip("clang unavailable")
	}
	root := t.TempDir()
	file := filepath.Join(root, "sample.cpp")
	code := "void helper() {}\n#if 0\nvoid packet_decoding() { helper(); }\n#endif\nvoid unrelated() { helper(); }\n"
	if err = os.WriteFile(file, []byte(code), 0600); err != nil {
		t.Fatal(err)
	}
	db, _ := json.Marshal([]map[string]any{{"directory": root, "file": file, "arguments": []string{"clang++", "-c", file}}})
	if err = os.WriteFile(filepath.Join(root, "compile_commands.json"), db, 0600); err != nil {
		t.Fatal(err)
	}
	sem, err := clang.NewClient(binary).Analyze(context.Background(), domain.SemanticRequest{Root: root, Query: "packet decoding"})
	if err != nil || sem.Status != "complete" {
		t.Fatalf("compiler failed: %v %s", err, sem.Status)
	}
	evidence := parser.BuildQueryEvidence("packet decoding", &sem, []domain.ExtractedFile{{Path: "sample.cpp", Content: code}}, root)
	if len(evidence.Definitions) != 0 {
		t.Fatal("disabled implementation promoted")
	}
}
