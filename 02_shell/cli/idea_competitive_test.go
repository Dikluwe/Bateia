/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_competitive_test.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestIdeaAutophagyRunsFromFrozenLocalFiles(t *testing.T) {
	dir := t.TempDir()
	ideaPath := filepath.Join(dir, "idea.json")
	bundlesPath := filepath.Join(dir, "bundles.json")
	rulesPath := filepath.Join(dir, "rules.json")
	inputs := []struct {
		path  string
		value any
	}{
		{ideaPath, domain.IdeaSpec{Problem: "analisar", Capabilities: []string{"dependency graph"}}},
		{bundlesPath, []domain.EvidenceBundle{{SchemaVersion: "1", Repository: domain.CommitRef{Owner: "local", Repo: "bateia", Commit: "abc"}, Complete: false}}},
		{rulesPath, usecase.DefaultEvidenceRules()},
	}
	for _, input := range inputs {
		if err := writeJSONAtomic(input.path, input.value); err != nil {
			t.Fatal(err)
		}
	}
	var out bytes.Buffer
	err := RunIdeaCommand([]string{"autofagia", "--idea", ideaPath, "--bundles", bundlesPath, "--rules", rulesPath, "--target", "local/bateia@abc"}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"disposition": "UNKNOWN"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestIdeaKatamariPreservesEmptyDenominatorAsInfeasible(t *testing.T) {
	dir := t.TempDir()
	autoPath := filepath.Join(dir, "auto.json")
	assessmentsPath := filepath.Join(dir, "assessments.json")
	policyPath := filepath.Join(dir, "policy.json")
	reusePath := filepath.Join(dir, "reuse.json")
	inputs := []struct {
		path  string
		value any
	}{
		{autoPath, domain.AutophagyReport{SchemaVersion: "1", Gaps: []string{"missing"}}},
		{assessmentsPath, domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{}}},
		{policyPath, domain.KatamariPolicy{SchemaVersion: "1", MaxProjects: 2, MaxComponents: 2}},
		{reusePath, domain.ReuseDecisionPolicy{SchemaVersion: "1"}},
	}
	for _, input := range inputs {
		if err := writeJSONAtomic(input.path, input.value); err != nil {
			t.Fatal(err)
		}
	}
	var out bytes.Buffer
	err := RunIdeaCommand([]string{"katamari", "--autofagia", autoPath, "--assessments", assessmentsPath, "--policy", policyPath, "--reuse-policy", reusePath}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"feasible": false`) || !strings.Contains(out.String(), `"missing"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestCompetitiveCommandsRejectMissingInputsAndHead(t *testing.T) {
	if err := RunIdeaCommand([]string{"autofagia"}, &bytes.Buffer{}); err == nil {
		t.Fatal("missing inputs survived")
	}
	if _, err := parseResolvedTarget("local/bateia@HEAD"); err == nil {
		t.Fatal("HEAD survived")
	}
	if err := RunIdeaCommand([]string{"katamari"}, &bytes.Buffer{}); err == nil {
		t.Fatal("missing Katamari inputs survived")
	}
}
