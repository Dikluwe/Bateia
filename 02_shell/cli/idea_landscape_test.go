/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_landscape_test.md
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

func TestIdeaLandscapeAnalyzeCLI(t *testing.T) {
	path := filepath.Join(t.TempDir(), "landscape.json")
	in := domain.TechnologyLandscapeRequest{SchemaVersion: "1", Current: domain.LandscapeSnapshot{AsOf: "2026-09-24", Need: []string{"search"}, Candidates: []domain.LandscapeCandidate{{ID: "org/tool", Kind: domain.LandscapeRepository, Name: "Tool", Capabilities: []domain.LandscapeCapability{{ID: "search", Status: domain.EvidenceProven, Sources: []string{"assessment:1"}}}, Standalone: domain.LandscapeAssertion{Status: domain.EvidenceProven, Sources: []string{"test:1"}}}}}}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := RunIdeaCommand([]string{"landscape", "analyze", "--input", path}, &out); err != nil {
		t.Fatal(err)
	}
	var report domain.TechnologyLandscapeReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Assessments) != 1 || report.Assessments[0].Role != domain.LandscapeSubstitute {
		t.Fatalf("relatório inesperado: %+v", report)
	}
	if err := RunIdeaCommand([]string{"landscape", "analyze"}, &out); err == nil || !strings.Contains(err.Error(), "--input") {
		t.Fatalf("erro de input ausente esperado: %v", err)
	}
}
