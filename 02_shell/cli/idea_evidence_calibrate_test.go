/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_calibrate_test.md
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

func TestEvidenceCalibrateCLI(t *testing.T) {
	if err := runIdeaEvidenceCalibrate(nil, &bytes.Buffer{}); err == nil {
		t.Fatal("comando aceitou ausência de observações")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "observations.json")
	data := `[{"id":"opaque","split":"TEST","expected":"INCONCLUSIVE","observed":{"status":"INCONCLUSIVE"},"source_exists":false}]`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := runIdeaEvidenceCalibrate([]string{"--observations", path}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"verdict": "PROMOTE_EVIDENCE_CATALOG"`) {
		t.Fatalf("saída inesperada: %s", out.String())
	}
}
