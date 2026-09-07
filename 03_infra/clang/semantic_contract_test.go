/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/semantic_contract_test.md
 * @layer L3
 * @updated 2026-09-07
 */
package clang_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/03_infra/clang"
)

func contractBinary(t *testing.T) string {
	t.Helper()
	binary, err := exec.LookPath("clang")
	if err != nil {
		t.Skip("Clang real indisponível: " + err.Error())
	}
	return binary
}

func contractWrite(t *testing.T, root, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
}

func contractConfig(t *testing.T, root string, arguments []string, command string) {
	t.Helper()
	entry := map[string]any{"directory": root, "file": "main.c"}
	if arguments != nil {
		entry["arguments"] = arguments
	}
	if command != "" {
		entry["command"] = command
	}
	b, err := json.Marshal([]any{entry})
	if err != nil {
		t.Fatal(err)
	}
	contractWrite(t, root, "compile_commands.json", string(b))
}

func contractFixture(t *testing.T) (string, []string) {
	t.Helper()
	root := t.TempDir()
	contractWrite(t, root, "helper.h", "#ifndef ENABLE_HELPER\n#error missing compile macro\n#endif\nstatic inline int helper(int x) { return x + OFFSET; }\n")
	contractWrite(t, root, "main.c", "#include \"helper.h\"\nint entry(int x);\nint entry(int x) { return helper(x); }\n")
	args := []string{"clang", "-std=c11", "-DENABLE_HELPER", "-DOFFSET=2", "-I", root, "-c", "main.c", "-o", "main.o"}
	contractConfig(t, root, args, "")
	return root, args
}

func contractAnalyze(t *testing.T, root string) (domain.SemanticReport, error) {
	t.Helper()
	return clang.NewClient(contractBinary(t)).Analyze(context.Background(), domain.SemanticRequest{Root: root, Query: "entry", Paths: []string{"main.c"}})
}

func contractIncomplete(t *testing.T, report domain.SemanticReport, err error) {
	t.Helper()
	if report.Status == "complete" {
		t.Fatal("impedimento foi promovido a complete")
	}
	if report.Status != "partial" && report.Status != "not_evaluated" {
		t.Errorf("estado inválido: %q", report.Status)
	}
	explanation := report.Reason
	for _, u := range report.Units {
		explanation += u.Diagnostics
	}
	if err == nil && strings.TrimSpace(explanation) == "" {
		t.Error("impedimento sem erro, motivo ou diagnóstico")
	}
}

func contractHashes(t *testing.T, root string, report domain.SemanticReport) map[string]string {
	t.Helper()
	got := map[string]string{}
	for _, unit := range report.Units {
		for _, input := range unit.Inputs {
			name := filepath.Base(input.Path)
			if name != "main.c" && name != "helper.h" && name != "compile_commands.json" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(root, name))
			if err != nil {
				t.Fatal(err)
			}
			sum := sha256.Sum256(data)
			want := hex.EncodeToString(sum[:])
			if input.SHA256 != want {
				t.Errorf("hash %s = %q, esperado %s", name, input.SHA256, want)
			}
			got[name] = input.SHA256
		}
	}
	for _, name := range []string{"main.c", "helper.h", "compile_commands.json"} {
		if got[name] == "" {
			t.Errorf("entrada examinada sem hash: %s", name)
		}
	}
	return got
}

func TestSemanticContractConfiguredDefinitionCallAndEvidence(t *testing.T) {
	root, _ := contractFixture(t)
	report, err := contractAnalyze(t, root)
	if err != nil {
		t.Fatal(err)
	}
	if report.Status != "complete" {
		t.Fatalf("unidade válida não completa: %+v", report)
	}
	if report.Analyzer == "" || report.Version == "" {
		t.Error("executável/versão ausentes")
	}
	if len(report.Units) != 1 {
		t.Fatalf("unidades = %d, esperado 1", len(report.Units))
	}
	u := report.Units[0]
	if filepath.Base(u.Path) != "main.c" || u.Status != "complete" {
		t.Errorf("unidade inválida: %+v", u)
	}
	if len(u.Arguments) == 0 {
		t.Error("argumentos efetivos ausentes")
	}
	ids := map[string]string{}
	definition, declaration := false, false
	for _, s := range u.Symbols {
		ids[s.ID] = s.Name
		if s.Name == "entry" {
			if filepath.Base(s.File) != "main.c" || s.Line < 2 || s.EndLine < s.Line {
				t.Errorf("intervalo inválido: %+v", s)
			}
			if s.Definition {
				definition = true
			} else {
				declaration = true
			}
		}
	}
	if !definition || !declaration {
		t.Errorf("entry: definição=%v declaração=%v", definition, declaration)
	}
	resolved := false
	for _, c := range u.Calls {
		caller, callee := c.Caller, c.Callee
		if ids[caller] != "" {
			caller = ids[caller]
		}
		if ids[callee] != "" {
			callee = ids[callee]
		}
		if caller == "entry" && callee == "helper" && filepath.Base(c.File) == "main.c" && c.Line == 3 {
			resolved = true
		}
	}
	if !resolved {
		t.Errorf("chamada entry -> helper resolvida ausente: %+v", u.Calls)
	}
	contractHashes(t, root, report)
}

func TestSemanticContractInputHashesTrackChanges(t *testing.T) {
	for _, name := range []string{"main.c", "helper.h", "compile_commands.json"} {
		t.Run(name, func(t *testing.T) {
			root, _ := contractFixture(t)
			before, err := contractAnalyze(t, root)
			if err != nil || before.Status != "complete" {
				t.Fatalf("baseline: %+v %v", before, err)
			}
			old := contractHashes(t, root, before)
			data, err := os.ReadFile(filepath.Join(root, name))
			if err != nil {
				t.Fatal(err)
			}
			contractWrite(t, root, name, string(data)+"\n")
			after, err := contractAnalyze(t, root)
			if err != nil || after.Status != "complete" {
				t.Fatalf("após edição: %+v %v", after, err)
			}
			newHashes := contractHashes(t, root, after)
			if old[name] == newHashes[name] {
				t.Errorf("hash não mudou para %s", name)
			}
		})
	}
}

func TestSemanticContractMissingContextAndInclude(t *testing.T) {
	for _, missing := range []string{"compile_commands.json", "helper.h"} {
		t.Run(missing, func(t *testing.T) {
			root, _ := contractFixture(t)
			if err := os.Remove(filepath.Join(root, missing)); err != nil {
				t.Fatal(err)
			}
			report, err := contractAnalyze(t, root)
			contractIncomplete(t, report, err)
		})
	}
}

func TestSemanticContractRejectsExecutableConfiguration(t *testing.T) {
	for _, mode := range []string{"command", "response", "plugin"} {
		t.Run(mode, func(t *testing.T) {
			root, args := contractFixture(t)
			marker := filepath.Join(root, "EXECUTED")
			switch mode {
			case "command":
				contractConfig(t, root, nil, "touch "+marker+"; clang -c main.c")
			case "response":
				contractWrite(t, root, "flags.rsp", "-DENABLE_HELPER -DOFFSET=2")
				contractConfig(t, root, append(args, "@"+filepath.Join(root, "flags.rsp")), "")
			case "plugin":
				contractConfig(t, root, append(args, "-Xclang", "-load", "-Xclang", filepath.Join(root, "untrusted.so")), "")
			}
			report, err := contractAnalyze(t, root)
			contractIncomplete(t, report, err)
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatalf("sentinela de execução: %v", err)
			}
		})
	}
}

func TestSemanticContractDoesNotExecuteRepositoryCompiler(t *testing.T) {
	root, args := contractFixture(t)
	marker := filepath.Join(root, "EXECUTED")
	compiler := filepath.Join(root, "compiler")
	contractWrite(t, root, "compiler", "#!/bin/sh\ntouch '"+marker+"'\nexit 0\n")
	if err := os.Chmod(compiler, 0700); err != nil {
		t.Fatal(err)
	}
	args[0] = compiler
	contractConfig(t, root, args, "")
	_, _ = contractAnalyze(t, root)
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("compilador do repositório executado: %v", err)
	}
}

func TestSemanticContractCanceledContextCannotComplete(t *testing.T) {
	root, _ := contractFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	report, err := clang.NewClient(contractBinary(t)).Analyze(ctx, domain.SemanticRequest{Root: root, Query: "entry", Paths: []string{"main.c"}})
	contractIncomplete(t, report, err)
}
