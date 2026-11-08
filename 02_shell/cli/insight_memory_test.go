/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_memory_test.md
 * @layer L2
 * @updated 2026-10-29
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInsightMemoryCLIReportFeedbackAndReuse(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.json")
	feedback := filepath.Join(dir, "memory.json")
	next := filepath.Join(dir, "next.json")
	factory := func() ([]domain.InsightSiteBinding, error) {
		return []domain.InsightSiteBinding{{Connector: &insightSiteCLIFixture{}}}, nil
	}
	if err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--output", first}, &bytes.Buffer{}, nil, nil, factory); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(first)
	var report domain.InsightSiteReport
	_ = json.Unmarshal(raw, &report)
	if report.ExplorationMemory == nil || len(report.ExplorationMemory.Entries) != 1 {
		t.Fatal("result did not persist memory")
	}
	command := []string{"memory", "--input", first, "--target", "https://gitee.com/org/repo", "--state", "discarded", "--reason", "Already investigated", "--output", feedback}
	if err := runIdeaInsights(command, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--memory", feedback, "--output", next}, &bytes.Buffer{}, nil, nil, factory); err != nil {
		t.Fatal(err)
	}
	raw, _ = os.ReadFile(next)
	_ = json.Unmarshal(raw, &report)
	if report.Readmes[0].Status != "MEMORY_DISCARDED" || report.ExplorationMemory.Entries[0].Reason != "Already investigated" {
		t.Fatal("discarded history ignored")
	}
	if err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--memory", feedback, "--output", feedback}, &bytes.Buffer{}, nil, nil, factory); err == nil {
		t.Fatal("memory input overwritten")
	}
}
func TestInsightMemoryCLIRejectsBeforeFactoryAndProtectsOutput(t *testing.T) {
	dir := t.TempDir()
	invalid := filepath.Join(dir, "invalid.json")
	_ = os.WriteFile(invalid, []byte(`{"schema_version":"insight-exploration-memory-v1","id":"tampered","generation":0,"entries":[]}`), 0600)
	for _, flags := range [][]string{{"--memory", invalid}, {"--memory", ""}, {"--follow-links", "-1"}, {"--follow-links", "5"}} {
		called := false
		factory := func() ([]domain.InsightSiteBinding, error) { called = true; return nil, nil }
		args := append([]string{"--question", "queue", "--sources", "gitee", "--output", filepath.Join(dir, "out.json")}, flags...)
		if err := RunInsightDiscoveryCommand(args, &bytes.Buffer{}, nil, nil, factory); err == nil || called {
			t.Fatal("invalid memory/options reached factory")
		}
	}
	first := filepath.Join(dir, "first.json")
	factory := func() ([]domain.InsightSiteBinding, error) {
		return []domain.InsightSiteBinding{{Connector: &insightSiteCLIFixture{}}}, nil
	}
	if err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--output", first}, &bytes.Buffer{}, nil, nil, factory); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(dir, "keep.json")
	_ = os.WriteFile(output, []byte("keep"), 0600)
	secret := "fixture-secret-token<>&"
	t.Setenv("GITHUB_TOKEN", secret)
	if err := runIdeaInsights([]string{"memory", "--input", first, "--target", "https://gitee.com/org/repo", "--state", "discarded", "--reason", secret, "--output", output}, &bytes.Buffer{}); err == nil {
		t.Fatal("secret persisted in memory")
	}
	data, _ := os.ReadFile(output)
	if string(data) != "keep" {
		t.Fatal("rejected update changed output")
	}
}

func TestInsightMemoryCLIGenerationLimitDoesNotPublishInvalidSnapshot(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.json")
	factory := func() ([]domain.InsightSiteBinding, error) {
		return []domain.InsightSiteBinding{{Connector: &insightSiteCLIFixture{}}}, nil
	}
	if err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--output", first}, &bytes.Buffer{}, nil, nil, factory); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(first)
	if err != nil {
		t.Fatal(err)
	}
	var report domain.InsightSiteReport
	if err := json.Unmarshal(raw, &report); err != nil || report.ExplorationMemory == nil {
		t.Fatal("fixture lacks memory", err)
	}
	for _, tc := range []struct {
		name       string
		generation int
		wantError  bool
	}{{"last-valid", 999998, false}, {"exhausted", 999999, true}} {
		t.Run(tc.name, func(t *testing.T) {
			memory := *report.ExplorationMemory
			memory.Generation, memory.ID = tc.generation, ""
			var canonical bytes.Buffer
			encoder := json.NewEncoder(&canonical)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(memory); err != nil {
				t.Fatal(err)
			}
			digest := sha256.Sum256(bytes.TrimSuffix(canonical.Bytes(), []byte{'\n'}))
			memory.ID = "exploration_memory_" + hex.EncodeToString(digest[:])
			if _, err := usecase.ValidateInsightExplorationMemory(&memory); err != nil {
				t.Fatal("boundary history should remain readable", err)
			}
			input, output := filepath.Join(dir, tc.name+"-input.json"), filepath.Join(dir, tc.name+"-output.json")
			data, _ := json.Marshal(memory)
			if err := os.WriteFile(input, data, 0600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			err := runIdeaInsights([]string{"memory", "--input", input, "--target", "https://gitee.com/org/repo", "--state", "pending", "--reason", "Boundary feedback", "--output", output}, &stdout)
			after, readErr := os.ReadFile(input)
			if readErr != nil || !bytes.Equal(data, after) {
				t.Fatal("feedback changed the input file", readErr)
			}
			if tc.wantError {
				_, statErr := os.Stat(output)
				if err == nil || !os.IsNotExist(statErr) || stdout.Len() != 0 {
					t.Fatal("invalid next generation was published", err, statErr, stdout.String())
				}
				called := false
				rejectedFactory := func() ([]domain.InsightSiteBinding, error) { called = true; return nil, nil }
				err = RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--memory", input, "--output", output}, &stdout, nil, nil, rejectedFactory)
				_, statErr = os.Stat(output)
				if err == nil || called || !os.IsNotExist(statErr) || stdout.Len() != 0 {
					t.Fatal("exhausted generation reached the factory or publication", err, called)
				}
				return
			}
			if err != nil {
				t.Fatal("last valid feedback failed", err)
			}
			updated, err := readInsightMemory(output)
			if err != nil || updated.Generation != 999999 {
				t.Fatal("published memory cannot be read", err)
			}
		})
	}
}
