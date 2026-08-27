/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_functional.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"flag"
	"fmt"
	"io"
)

func runIdeaFunctional(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: idea functional <expand|audit>")
	}
	switch args[0] {
	case "expand":
		return runFunctionalExpand(args[1:], stdout)
	case "audit":
		return runFunctionalAudit(args[1:], stdout)
	default:
		return fmt.Errorf("operação funcional desconhecida")
	}
}
func runFunctionalExpand(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("functional expand", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	ontologyPath := fs.String("ontology-expansion", "", "OntologyExpansionResult")
	lexiconPath := fs.String("lexicon", "", "FunctionalLexicon")
	templatesPath := fs.String("templates", "", "templates")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *ontologyPath == "" || *lexiconPath == "" || *templatesPath == "" {
		return fmt.Errorf("expand exige ontology-expansion, lexicon e templates")
	}
	var o domain.OntologyExpansionResult
	var l domain.FunctionalLexicon
	var t []domain.FunctionalQueryTemplate
	if err := readOntologyJSON(*ontologyPath, &o); err != nil {
		return err
	}
	if err := readOntologyJSON(*lexiconPath, &l); err != nil {
		return err
	}
	if err := readOntologyJSON(*templatesPath, &t); err != nil {
		return err
	}
	result, err := usecase.ExpandFunctionalQueries(o, l, t)
	if err != nil {
		return err
	}
	return writeIdeaJSON(result, *output, stdout)
}
func runFunctionalAudit(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("functional audit", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	queriesPath := fs.String("queries", "", "FunctionalExpansion")
	lexiconPath := fs.String("lexicon", "", "FunctionalLexicon")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *queriesPath == "" || *lexiconPath == "" {
		return fmt.Errorf("audit exige queries e lexicon")
	}
	var e domain.FunctionalExpansion
	var l domain.FunctionalLexicon
	if err := readOntologyJSON(*queriesPath, &e); err != nil {
		return err
	}
	if err := readOntologyJSON(*lexiconPath, &l); err != nil {
		return err
	}
	return writeIdeaJSON(usecase.AuditFunctionalQueries(e, l), *output, stdout)
}
