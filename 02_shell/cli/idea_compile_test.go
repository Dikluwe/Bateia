/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_compile_test.md
 * @layer L2
 * @updated 2026-09-24
 */
package cli_test

import (
	"bateia/01_core/domain"
	"bateia/02_shell/cli"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestIdeaIntakeCompileEndToEnd(t *testing.T) {
	dir := t.TempDir()
	draftPath := filepath.Join(dir, "draft.json")
	preparedPath := filepath.Join(dir, "prepared.json")
	draft := []byte(`{"idea":{"problem":"buscar arquivos privados","capabilities":["buscar trechos"],"qualities":["responder em até 2 segundos"],"operating_conditions":["sem rede"]}}`)
	if err := os.WriteFile(draftPath, draft, 0600); err != nil {
		t.Fatal(err)
	}
	if err := cli.RunIdeaCommand([]string{"intake", "prepare", "--input", draftPath, "--output", preparedPath}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := cli.RunIdeaCommand([]string{"intake", "compile", "--input", preparedPath}, &output); err != nil {
		t.Fatal(err)
	}
	var plan domain.CapabilityCompilation
	if err := json.Unmarshal(output.Bytes(), &plan); err != nil {
		t.Fatalf("invalid compilation JSON: %v", err)
	}
	if plan.Goal != "buscar arquivos privados" || plan.IntakeStatus != domain.IntakeReady {
		t.Fatalf("unexpected compilation: %+v", plan)
	}
	foundQuality, foundCondition := false, false
	for _, r := range plan.Requirements {
		switch r.Kind {
		case "quality":
			foundQuality = r.Authority == domain.AuthorityUser
		case "operating_condition":
			foundCondition = r.Authority == domain.AuthorityUser
		}
	}
	if !foundQuality || !foundCondition || len(plan.Questions) == 0 {
		t.Fatalf("advanced fields lost: %+v", plan)
	}
}
