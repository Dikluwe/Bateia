/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/effects/source_binding_test.md
 * @layer L3
 * @updated 2026-09-08
 */

package effects_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/03_infra/effects"
)

func binding63SHA(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func binding63Source(t *testing.T, source string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "candidate.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func binding63JSON(t *testing.T, report domain.EffectReport) map[string]any {
	t.Helper()
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		t.Fatal(err)
	}
	binding, ok := object["source_binding"].(map[string]any)
	if !ok {
		t.Fatalf("source_binding ausente ou inválido: %s", raw)
	}
	return binding
}

func binding63Verify(t *testing.T, source, symbol, expected string) domain.EffectReport {
	t.Helper()
	report, err := effects.NewRunner().Verify(context.Background(), domain.EffectRequest{SourcePath: binding63Source(t, source), Symbol: symbol, ExpectedSourceSHA256: expected, Timeout: 30 * time.Second, Seed: 63})
	if err != nil {
		t.Fatal(err)
	}
	return report
}

func TestSourceBinding63ExactPhysicalDeclaration(t *testing.T) {
	declaration := "func SortCopy(input []int64) []int64 {\n\t// coração λ — bytes brutos\n\toutput := clone(input)\n\tsort.Slice(output, func(i, j int) bool { return output[i] < output[j] })\n\treturn output\n}"
	source := "package candidate\nimport \"sort\"\n// func SortCopy(falso) {} — não é declaração\nfunc clone(input []int64) []int64 { return append([]int64(nil), input...) }\n//line inventado.go:700\n// SortCopy preserva entrada.\n" + declaration + "\n"
	for _, tc := range []struct{ name, newline, hashMode string }{{"compatibility_without_hash", "\n", "none"}, {"lowercase_crlf", "\r\n", "lower"}, {"uppercase_lf", "\n", "upper"}} {
		t.Run(tc.name, func(t *testing.T) {
			raw := strings.ReplaceAll(source, "\n", tc.newline)
			fragment := strings.ReplaceAll(declaration, "\n", tc.newline)
			expected := binding63SHA(raw)
			status := "matched"
			if tc.hashMode == "upper" {
				expected = strings.ToUpper(expected)
			}
			if tc.hashMode == "none" {
				expected = ""
				status = "not_requested"
			}
			report := binding63Verify(t, raw, "SortCopy", expected)
			if report.Status != "supported_in_scope" {
				t.Fatalf("fonte aceita perdeu veredito: %+v", report)
			}
			if len(report.Cases) == 0 {
				t.Fatal("nenhum caso registrado")
			}
			binding := binding63JSON(t, report)
			if binding["status"] != status || binding["symbol"] != "SortCopy" || binding["expected_sha256"] != strings.ToLower(expected) {
				t.Fatalf("binding: %#v", binding)
			}
			actual := binding63SHA(raw)
			if binding["actual_sha256"] != actual || report.SourceSHA256 != actual || report.Execution["source_sha256"] != actual {
				t.Fatalf("hash integral inconsistente: %+v", report)
			}
			decl, ok := binding["declaration"].(map[string]any)
			if !ok {
				t.Fatalf("declaração ausente: %#v", binding)
			}
			start := strings.Index(raw, fragment)
			end := start + len(fragment)
			want := map[string]any{"start_byte": float64(start), "end_byte": float64(end), "start_line": float64(strings.Count(raw[:start], "\n") + 1), "end_line": float64(strings.Count(raw[:end], "\n") + 1), "content": fragment, "sha256": binding63SHA(fragment)}
			for field, value := range want {
				if decl[field] != value {
					t.Errorf("%s: obtido %#v; esperado %#v", field, decl[field], value)
				}
			}
		})
	}
}

// A chamada de qualquer comando go vira evidência persistente, sem compilar a fonte.
func binding63NoGo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	marker := filepath.Join(dir, "called")
	script := "#!/bin/sh\nprintf invoked > '" + marker + "'\nexit 93\n"
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return marker
}

func binding63NoExecution(t *testing.T, report domain.EffectReport, marker string) {
	t.Helper()
	if report.Status != "unknown" {
		t.Fatalf("rejeição promovida: %+v", report)
	}
	for _, c := range report.Cases {
		if c.Observation.Completed {
			t.Fatal("rejeição executou caso")
		}
	}
	if binding63JSON(t, report)["declaration"] != nil {
		t.Fatal("rejeição registrou declaração")
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("go foi executado ou marcador ilegível: %v", err)
	}
}

func TestSourceBinding63MismatchBeforeExecution(t *testing.T) {
	original := "package candidate\nfunc SortCopy(x []int64) []int64 { for {} }\nfunc helper() int { return 1 }\n"
	changed := strings.Replace(original, "return 1", "return 2", 1)
	marker := binding63NoGo(t)
	report := binding63Verify(t, changed, "SortCopy", strings.ToUpper(binding63SHA(original)))
	binding63NoExecution(t, report, marker)
	binding := binding63JSON(t, report)
	if binding["status"] != "mismatch" || binding["expected_sha256"] != binding63SHA(original) || binding["actual_sha256"] != binding63SHA(changed) {
		t.Fatalf("hash do helper não vinculou arquivo: %#v", binding)
	}
	if report.Reason != "source_identity_mismatch" || report.Execution["stop_reason"] != "source_identity_mismatch" {
		t.Fatalf("motivo incorreto: %+v", report)
	}
	if report.SourceSHA256 != binding63SHA(changed) || report.Execution["source_sha256"] != binding63SHA(changed) {
		t.Fatal("hash real inconsistente")
	}
}

func TestSourceBinding63MatchingHashCannotPromoteRejectedSource(t *testing.T) {
	for _, source := range []string{
		"package candidate\nvar global = 1\nfunc SortCopy(x []int64) []int64 { return x }\n",
		"package candidate\nfunc Another(x []int64) []int64 { return x }\n",
		"package candidate\nfunc SortCopy(x []int64) []int64 { return x }\nfunc SortCopy(x []int64) []int64 { return x }\n",
		"package candidate\nfunc SortCopy( {\n",
		"package candidate\nimport \"os\"\nfunc SortCopy(x []int64) []int64 { os.Exit(0); return x }\n",
	} {
		t.Run(binding63SHA(source)[:8], func(t *testing.T) {
			marker := binding63NoGo(t)
			report := binding63Verify(t, source, "SortCopy", binding63SHA(source))
			binding63NoExecution(t, report, marker)
			if binding63JSON(t, report)["status"] != "matched" {
				t.Fatal("preflight alterou identidade do arquivo")
			}
		})
	}
}

func TestSourceBinding63MalformedHashPrecedesRead(t *testing.T) {
	for _, expected := range []string{" ", strings.Repeat("a", 63), strings.Repeat("a", 65), strings.Repeat("g", 64), strings.Repeat("Ａ", 64), strings.Repeat("a", 63) + "\n"} {
		t.Run(binding63SHA(expected)[:8], func(t *testing.T) {
			marker := binding63NoGo(t)
			_, err := effects.NewRunner().Verify(context.Background(), domain.EffectRequest{SourcePath: filepath.Join(t.TempDir(), "missing.go"), Symbol: "SortCopy", ExpectedSourceSHA256: expected, Timeout: time.Second})
			if err == nil {
				t.Fatal("hash inválido aceito")
			}
			if errors.Is(err, os.ErrNotExist) {
				t.Fatalf("leu fonte antes de validar hash: %v", err)
			}
			if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("hash inválido executou go")
			}
		})
	}
}

func TestSourceBinding63ReadErrorPreserved(t *testing.T) {
	_, err := effects.NewRunner().Verify(context.Background(), domain.EffectRequest{SourcePath: filepath.Join(t.TempDir(), "missing.go"), Symbol: "SortCopy", ExpectedSourceSHA256: strings.Repeat("a", 64), Timeout: time.Second})
	if err == nil {
		t.Fatal("erro de leitura suprimido")
	}
}
