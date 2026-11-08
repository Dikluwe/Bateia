/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_collection_test.md
 * @layer L2
 * @updated 2026-11-08
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInsightCollectionCLIOptInAndLegacyShape(t *testing.T) {
	pin := strings.Repeat("a", 40)
	plan := domain.EvidenceExtractionPlan{SchemaVersion: "1", Repository: domain.CommitRef{Owner: "org", Repo: "repo", Commit: pin}, Candidate: domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: pin}, Targets: []domain.ExtractionTarget{{Path: "README.md", Kind: domain.EvidenceDocumentation}}, Budget: domain.ExtractionBudget{Documents: 2, Manifests: 2, CodeFiles: 2, Tests: 2, TotalBytes: 32768, TreeDepth: 4}}
	path := filepath.Join(t.TempDir(), "plan.json")
	data, _ := json.Marshal([]domain.EvidenceExtractionPlan{plan})
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := RunIdeaEvidenceCollectCommand([]string{"--plan", path}, &out, fixtureContent{}); err != nil {
		t.Fatal(err)
	}
	var legacy []domain.EvidenceBundle
	if err := json.Unmarshal(out.Bytes(), &legacy); err != nil || len(legacy) != 1 {
		t.Fatal("legacy array changed")
	}
	out.Reset()
	if err := RunIdeaEvidenceCollectCommand([]string{"--plan", path, "--question", "How to inspect dependencies?", "--max-passages", "4"}, &out, fixtureContent{}); err != nil {
		t.Fatal(err)
	}
	var report domain.InsightCollectionReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != "insight-collection-v1" || len(report.Bundles) != 1 || len(report.Exploration.Cards) == 0 || len(report.Exploration.Tasks) != 7 {
		t.Fatal("new bytes not explored")
	}
	for _, args := range [][]string{{"--plan", path, "--familiar", "approach"}, {"--plan", path, "--max-passages", "4"}, {"--plan", path, "--question", ""}, {"--plan", path, "--question", "question", "--output", path}} {
		out.Reset()
		if err := RunIdeaEvidenceCollectCommand(args, &out, fixtureContent{}); err == nil || out.Len() != 0 {
			t.Fatal("invalid input published")
		}
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(after, data) {
		t.Fatal("plan overwritten")
	}
}
