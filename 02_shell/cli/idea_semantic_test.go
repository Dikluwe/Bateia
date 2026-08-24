/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_semantic.md
 * @layer L2
 * @updated 2026-08-24
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

type cliSemanticExpander struct{}

func (cliSemanticExpander) Expand(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error) {
	source := plan.Expansions[0]
	return []domain.SearchExpansion{{
		Facet: source.Facet, Query: "semantic alternative", Rationale: "fixture",
		Covers:     append([]string(nil), source.Covers...),
		Derivation: &domain.ExpansionDerivation{Strategy: "lexical", RuleIDs: []string{"fixture"}, SourceIDs: []string{source.ID}},
	}}, nil
}

func TestRunIdeaExpandProducesEnrichedPlan(t *testing.T) {
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{"expand", "--input", filepath.Join("testdata", "idea_spec.json"), "--strategy", "lexical"}, &stdout, cliSemanticExpander{})
	if err != nil {
		t.Fatal(err)
	}
	var plan domain.ExpansionPlan
	if err := json.Unmarshal(stdout.Bytes(), &plan); err != nil {
		t.Fatalf("saída inválida: %v", err)
	}
	last := plan.Expansions[len(plan.Expansions)-1]
	if last.Query != "semantic alternative" || last.Derivation == nil {
		t.Fatalf("plano não enriquecido: %+v", last)
	}
}

func TestRunIdeaEvaluateComparesCorpus(t *testing.T) {
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{"evaluate", "--corpus", filepath.Join("testdata", "evaluation_corpus.json"), "--strategy", "lexical"}, &stdout, cliSemanticExpander{})
	if err != nil {
		t.Fatal(err)
	}
	var report domain.EvaluationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("saída inválida: %v", err)
	}
	if report.Strategy != "lexical" || len(report.Cases) != 8 {
		t.Fatalf("relatório inesperado: %+v", report)
	}
}

func TestRunIdeaExpandWritesOutputWithoutStdout(t *testing.T) {
	output := filepath.Join(t.TempDir(), "expanded.json")
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{
		"expand", "--input", filepath.Join("testdata", "idea_spec.json"),
		"--strategy", "lexical", "--output", output,
	}, &stdout, cliSemanticExpander{})
	if err != nil {
		t.Fatal(err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout deveria estar vazio: %q", stdout.String())
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	var plan domain.ExpansionPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		t.Fatalf("arquivo inválido: %v", err)
	}
}

func TestRunIdeaEvaluateRejectsInvalidCorpus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.json")
	if err := os.WriteFile(path, []byte(`{"cases":`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := RunIdeaCommand([]string{"evaluate", "--corpus", path}, &bytes.Buffer{}); err == nil {
		t.Fatal("corpus inválido foi aceito")
	}
}

func TestRunIdeaSemanticRejectsUnknownOrMissingStrategy(t *testing.T) {
	input := filepath.Join("testdata", "idea_spec.json")
	if err := RunIdeaCommand([]string{"expand", "--input", input, "--strategy", "magic"}, &bytes.Buffer{}); err == nil {
		t.Fatal("estratégia desconhecida foi aceita")
	}
	if err := RunIdeaCommand([]string{"expand", "--input", input, "--strategy", "lexical"}, &bytes.Buffer{}); err == nil {
		t.Fatal("estratégia lexical sem adaptador foi aceita")
	}
}
