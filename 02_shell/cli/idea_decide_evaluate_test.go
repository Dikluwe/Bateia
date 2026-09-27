/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_decide_evaluate_test.md
 * @layer L2
 * @updated 2026-09-24
 */
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestIdeaDecideEvaluateCLI(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cases.json")
	input := domain.SelectiveDecisionRequest{SchemaVersion: "1", Split: domain.SplitTest, Cases: []domain.SelectiveDecisionCase{{ID: "a", Observed: domain.DecisionBuild, Expected: domain.DecisionBuild}}}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := runIdeaDecide([]string{"evaluate", "--input", path}, &out); err != nil {
		t.Fatal(err)
	}
	var report domain.SelectiveDecisionReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Accepted != 1 || report.CurveStatus != "SCORE_UNAVAILABLE" {
		t.Fatalf("relatório inesperado: %+v", report)
	}
	if err := runIdeaDecide([]string{"evaluate"}, &out); err == nil || !strings.Contains(err.Error(), "--input") {
		t.Fatalf("erro esperado para input ausente: %v", err)
	}
}
