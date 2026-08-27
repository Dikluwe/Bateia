/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_ontology_test.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func TestIdeaOntologyEvaluateAndMaterializeAuditedPlan(t *testing.T) {
	ontology := filepath.Join("..", "..", "03_infra", "lexical", "discovery_ontology.json")
	var evaluation bytes.Buffer
	if err := RunIdeaCommand([]string{"ontology", "evaluate", "--corpus", filepath.Join("testdata", "discovery_ontology_corpus.json"), "--ontology", ontology}, &evaluation); err != nil {
		t.Fatal(err)
	}
	var report domain.OntologyEvaluationReport
	if json.Unmarshal(evaluation.Bytes(), &report) != nil || !report.Passed {
		t.Fatalf("gate offline falhou: %s", evaluation.Bytes())
	}
	dir := t.TempDir()
	input := filepath.Join(dir, "idea.json")
	planPath := filepath.Join(dir, "plan.json")
	if err := os.WriteFile(input, []byte(`{"problem":"não começar do zero","capabilities":["buscar repositórios"]}`), 0600); err != nil {
		t.Fatal(err)
	}
	var result bytes.Buffer
	if err := RunIdeaCommand([]string{"ontology", "expand", "--input", input, "--ontology", ontology, "--plan-output", planPath}, &result); err != nil {
		t.Fatal(err)
	}
	var plan domain.ExpansionPlan
	data, _ := os.ReadFile(planPath)
	if json.Unmarshal(data, &plan) != nil || len(plan.Expansions) != 4 {
		t.Fatalf("plano inválido: %s", data)
	}
}

func TestIdeaDiscoverConsumesPlanWithoutSemanticExpander(t *testing.T) {
	dir := t.TempDir()
	planPath := filepath.Join(dir, "plan.json")
	output := filepath.Join(dir, "snapshot.json")
	plan := domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "x"}, Expansions: []domain.SearchExpansion{{ID: "q1", Query: "software reuse discovery"}}}
	b, _ := json.Marshal(plan)
	_ = os.WriteFile(planPath, b, 0600)
	p := &cliDiscoveryProvider{}
	if err := RunIdeaDiscoverCommand([]string{"--plan", planPath, "--output", output, "--concurrency", "1"}, &bytes.Buffer{}, &bytes.Buffer{}, nil, p); err != nil {
		t.Fatal(err)
	}
	if p.calls != 1 {
		t.Fatalf("calls=%d", p.calls)
	}
}
