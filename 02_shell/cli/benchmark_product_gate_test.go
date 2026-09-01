/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/benchmark_product_gate_test.md
 * @layer L2
 * @updated 2026-08-31
 */
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBenchmarkProductGateCLI(t *testing.T) {
	if err := runBenchmarkProductGate(nil, &bytes.Buffer{}); err == nil {
		t.Fatal("comando aceitou ausência de input")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "input.json")
	data := `{"schema_version":"1","domain":"software-search-and-evaluation","case_ids":["1","2","3","4","5","6","7","8","9","10"],"discovery_arms":[{"arm":"N","planned":10,"completed":5}]}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := runBenchmarkProductGate([]string{"--input", path}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"verdict": "FAIL_DISCOVERY"`) {
		t.Fatalf("saída inesperada: %s", out.String())
	}
}
