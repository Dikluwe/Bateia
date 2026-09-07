/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_competitive_step0041_test.md
 * @layer L2
 * @updated 2026-09-07
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestStep0041CompetitiveCLIRetainsInvestigativeValue(t *testing.T) {
	dir := t.TempDir()
	files := map[string]any{
		"idea.json":    domain.IdeaSpec{Capabilities: []string{"dependency graph"}},
		"bundles.json": []domain.EvidenceBundle{{SchemaVersion: "1", Repository: domain.CommitRef{Owner: "local", Repo: "target", Commit: "abc"}, Complete: true, TreeComplete: true}},
		"rules.json":   usecase.DefaultEvidenceRules(),
		"policy.json":  domain.KatamariPolicy{SchemaVersion: "1", MaxProjects: 2, MaxComponents: 2, AllowedLicenses: []string{"MIT"}},
		"reuse.json":   domain.ReuseDecisionPolicy{SchemaVersion: "1"},
	}
	ref := domain.CommitRef{Owner: "external", Repo: "candidate", Commit: "def"}
	files["assessments.json"] = domain.AssessmentReport{Assessments: []domain.CandidateAssessment{{Repository: ref, Capabilities: []domain.CapabilityEvidence{{Capability: "dependency graph", Status: domain.EvidencePartial, Sources: []domain.EvidenceSource{{Ref: ref, Kind: domain.EvidenceSymbol, Path: "graph.go", Line: 1, Snippet: "type DependencyGraph struct{}", RuleID: "lexical", ContentHash: strings.Repeat("a", 64)}}}}}}}
	for name, value := range files {
		if err := writeJSONAtomic(filepath.Join(dir, name), value); err != nil {
			t.Fatal(err)
		}
	}
	if err := RunIdeaCommand([]string{"autofagia", "--idea", filepath.Join(dir, "idea.json"), "--bundles", filepath.Join(dir, "bundles.json"), "--rules", filepath.Join(dir, "rules.json"), "--target", "local/target@abc", "--output", filepath.Join(dir, "auto.json")}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := RunIdeaCommand([]string{"katamari", "--autofagia", filepath.Join(dir, "auto.json"), "--assessments", filepath.Join(dir, "assessments.json"), "--policy", filepath.Join(dir, "policy.json"), "--reuse-policy", filepath.Join(dir, "reuse.json")}, &output); err != nil {
		t.Fatal(err)
	}
	var p domain.KatamariProposal
	if err := json.Unmarshal(output.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, o := range p.Options {
		found = found || o.Action == "INVESTIGATE_EXTERNAL" && len(o.NextChecks) > 0
	}
	if p.SchemaVersion != "2" || p.Feasible || len(p.Components) > 0 || !found {
		t.Fatalf("CLI lost useful conservative plan: %s", output.String())
	}
}
