/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/client_test.md
 * @layer L3
 * @updated 2026-10-09
 */
package clang_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/clang"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func TestSiblingMesonDatabaseWithGeneratedHeader(t *testing.T) {
	parent := t.TempDir()
	root, build := filepath.Join(parent, "source"), filepath.Join(parent, "build")
	for _, dir := range []string{root, build} {
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	contractWrite(t, root, "main.c", "#include \"generated.h\"\nconst char *entry(void) { return LABEL; }\n")
	contractWrite(t, build, "generated.h", "#ifndef LABEL\n#error missing compilation context\n#endif\n")
	marker := filepath.Join(build, "EXECUTED")
	compiler := filepath.Join(build, "compiler")
	contractWrite(t, build, "compiler", "#!/bin/sh\ntouch '"+marker+"'\n")
	if err := os.Chmod(compiler, 0700); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal([]map[string]any{{"directory": ".", "file": "../source/main.c", "command": compiler + " -I. '-DLABEL=\"literal with spaces\"' -mtls-dialect=gnu2 -flifetime-dse=1 -c ../source/main.c -o main.o"}})
	// Exercise a real database beyond the previous 4 MiB bound.
	data = append(data, []byte(strings.Repeat(" ", 4*1024*1024))...)
	path := filepath.Join(build, "compile_commands.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	report, err := contractAnalyze(t, root)
	if err != nil || report.Status != "complete" || len(report.Units) != 1 {
		t.Fatalf("sibling Meson context: %+v %v", report, err)
	}
	want := fmt.Sprintf("%x", sha256.Sum256(data))
	if report.ConfigurationHash != want {
		t.Fatal("database hash mismatch")
	}
	found := false
	for _, input := range report.Units[0].Inputs {
		if input.Path == path && input.SHA256 == want {
			found = true
		}
	}
	if !found {
		t.Fatal("sibling database provenance missing")
	}
	for _, name := range []string{"EXECUTED", "main.o"} {
		if _, err := os.Stat(filepath.Join(build, name)); !os.IsNotExist(err) {
			t.Fatalf("unexpected effect %s: %v", name, err)
		}
	}
}

func TestDatabasePrecedenceAndInvalidLocalDoNotFallBack(t *testing.T) {
	root, _ := contractFixture(t)
	build := filepath.Join(root, "build")
	if err := os.Mkdir(build, 0700); err != nil {
		t.Fatal(err)
	}
	contractWrite(t, build, "compile_commands.json", "invalid lower priority")
	report, err := contractAnalyze(t, root)
	if err != nil || report.Status != "complete" {
		t.Fatalf("local precedence: %+v %v", report, err)
	}
	contractWrite(t, root, "compile_commands.json", "invalid preferred database")
	report, err = contractAnalyze(t, root)
	contractIncomplete(t, report, err)
	if !strings.Contains(report.Reason, "inválido") {
		t.Fatalf("no invalid database diagnostic: %+v", report)
	}
}

func TestDatabaseLimitIsObservable(t *testing.T) {
	root := t.TempDir()
	contractWrite(t, root, "compile_commands.json", strings.Repeat(" ", 16*1024*1024+1))
	report, err := contractAnalyze(t, root)
	contractIncomplete(t, report, err)
	if !strings.Contains(report.Reason, "limite de bytes") {
		t.Fatalf("missing size diagnostic: %+v", report)
	}
}

func TestLiteralCommandRejectsShellAndCompilerExtensions(t *testing.T) {
	for _, command := range []string{
		"clang -c main.c; touch EXECUTED", "clang -c main.c $(touch EXECUTED)",
		"clang -c main.c `touch EXECUTED`", "clang -c main.c | cat",
		"clang -c main.c > EXECUTED", "clang -c main.c @flags.rsp",
		"clang -Xclang -load -Xclang plugin.so -c main.c", "clang 'unterminated",
	} {
		t.Run(command, func(t *testing.T) {
			root, _ := contractFixture(t)
			contractConfig(t, root, nil, command)
			report, err := contractAnalyze(t, root)
			contractIncomplete(t, report, err)
			if _, err := os.Stat(filepath.Join(root, "EXECUTED")); !os.IsNotExist(err) {
				t.Fatalf("command executed: %v", err)
			}
		})
	}
}
