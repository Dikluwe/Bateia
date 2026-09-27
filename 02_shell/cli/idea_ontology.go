/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_ontology.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaOntology(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea ontology <expand|evaluate|audit>")
	}
	switch args[0] {
	case "expand":
		return runIdeaOntologyExpand(args[1:], stdout)
	case "evaluate":
		return runIdeaOntologyEvaluate(args[1:], stdout)
	case "audit":
		return runIdeaOntologyAudit(args[1:], stdout)
	default:
		return fmt.Errorf("operação ontológica desconhecida %q", args[0])
	}
}

func readOntologyJSON(path string, destination any) error {
	data, err := readIdeaJSONInput(path)
	if err != nil {
		return fmt.Errorf("não foi possível ler %s: %w", path, err)
	}
	if err := json.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("JSON inválido em %s: %w", path, err)
	}
	return nil
}

func runIdeaOntologyExpand(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea ontology expand", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "IdeaSpec JSON")
	ontologyPath := fs.String("ontology", "", "Ontologia JSON")
	output := fs.String("output", "", "Resultado JSON; omita para stdout")
	planOutput := fs.String("plan-output", "", "Plano auditado aceito para descoberta")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros inválidos: %w", err)
	}
	if *input == "" || *ontologyPath == "" {
		return fmt.Errorf("ontology expand exige --input e --ontology")
	}
	spec, err := readIdeaSpec(*input)
	if err != nil {
		return err
	}
	var ontology domain.DiscoveryOntology
	if err := readOntologyJSON(*ontologyPath, &ontology); err != nil {
		return err
	}
	literal, err := usecase.ExpandIdea(spec)
	if err != nil {
		return fmt.Errorf("IdeaSpec inválida: %w", err)
	}
	result, err := usecase.ExpandWithOntology(*literal, ontology)
	if err != nil {
		return fmt.Errorf("expansão ontológica inválida: %w", err)
	}
	quality, err := usecase.AuditOntologyQueries(result.Expansions, ontology)
	if err != nil {
		return err
	}
	if quality.Rejected > 0 {
		return fmt.Errorf("plano ontológico contém %d queries rejeitadas", quality.Rejected)
	}
	if *planOutput != "" {
		if err := writeIdeaJSON(domain.ExpansionPlan{Idea: spec, Expansions: quality.Accepted}, *planOutput, stdout); err != nil {
			return err
		}
	}
	return writeIdeaJSON(result, *output, stdout)
}

func runIdeaOntologyEvaluate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea ontology evaluate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	corpusPath := fs.String("corpus", "", "Corpus JSON")
	ontologyPath := fs.String("ontology", "", "Ontologia JSON")
	output := fs.String("output", "", "Relatório JSON; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros inválidos: %w", err)
	}
	if *corpusPath == "" || *ontologyPath == "" {
		return fmt.Errorf("ontology evaluate exige --corpus e --ontology")
	}
	var corpus domain.OntologyEvaluationCorpus
	var ontology domain.DiscoveryOntology
	if err := readOntologyJSON(*corpusPath, &corpus); err != nil {
		return err
	}
	if err := readOntologyJSON(*ontologyPath, &ontology); err != nil {
		return err
	}
	report, err := usecase.EvaluateDiscoveryOntology(corpus, ontology)
	if err != nil {
		return err
	}
	return writeIdeaJSON(report, *output, stdout)
}

func runIdeaOntologyAudit(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea ontology audit", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	planPath := fs.String("plan", "", "ExpansionPlan JSON")
	ontologyPath := fs.String("ontology", "", "Ontologia JSON")
	output := fs.String("output", "", "Relatório JSON; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros inválidos: %w", err)
	}
	if *planPath == "" || *ontologyPath == "" {
		return fmt.Errorf("ontology audit exige --plan e --ontology")
	}
	var plan domain.ExpansionPlan
	var ontology domain.DiscoveryOntology
	if err := readOntologyJSON(*planPath, &plan); err != nil {
		return err
	}
	if err := readOntologyJSON(*ontologyPath, &ontology); err != nil {
		return err
	}
	report, err := usecase.AuditOntologyQueries(plan.Expansions, ontology)
	if err != nil {
		return err
	}
	return writeIdeaJSON(report, *output, stdout)
}
