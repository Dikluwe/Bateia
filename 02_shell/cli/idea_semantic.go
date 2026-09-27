/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_semantic.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func selectIdeaStrategy(name string, lexical usecase.SemanticExpander) (usecase.SemanticExpander, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "literal":
		return nil, nil
	case "lexical":
		if lexical == nil {
			return nil, fmt.Errorf("estratégia lexical não foi configurada")
		}
		return lexical, nil
	default:
		return nil, fmt.Errorf("estratégia %q desconhecida; use literal ou lexical", name)
	}
}

func runIdeaExpand(args []string, stdout io.Writer, lexical usecase.SemanticExpander) error {
	fs := flag.NewFlagSet("idea expand", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inputPath := fs.String("input", "", "Arquivo JSON contendo IdeaSpec")
	strategyName := fs.String("strategy", "lexical", "Estratégia literal ou lexical")
	outputPath := fs.String("output", "", "Arquivo JSON de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea expand inválidos: %w", err)
	}
	if *inputPath == "" {
		return fmt.Errorf("idea expand exige --input com um arquivo IdeaSpec JSON")
	}
	strategy, err := selectIdeaStrategy(*strategyName, lexical)
	if err != nil {
		return err
	}
	spec, err := readIdeaSpec(*inputPath)
	if err != nil {
		return err
	}
	literal, err := usecase.ExpandIdea(spec)
	if err != nil {
		return fmt.Errorf("IdeaSpec inválida: %w", err)
	}
	result, err := usecase.EnrichExpansionPlan(*literal, strategy)
	if err != nil {
		return fmt.Errorf("falha na expansão %s: %w", *strategyName, err)
	}
	return writeIdeaJSON(result.Plan, *outputPath, stdout)
}

func runIdeaEvaluate(args []string, stdout io.Writer, lexical usecase.SemanticExpander) error {
	fs := flag.NewFlagSet("idea evaluate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	corpusPath := fs.String("corpus", "", "Arquivo JSON contendo o corpus")
	strategyName := fs.String("strategy", "literal", "Estratégia literal ou lexical")
	outputPath := fs.String("output", "", "Arquivo JSON de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea evaluate inválidos: %w", err)
	}
	if *corpusPath == "" {
		return fmt.Errorf("idea evaluate exige --corpus com um arquivo JSON")
	}
	strategy, err := selectIdeaStrategy(*strategyName, lexical)
	if err != nil {
		return err
	}
	data, err := readIdeaJSONInput(*corpusPath)
	if err != nil {
		return fmt.Errorf("não foi possível ler corpus %s: %w", *corpusPath, err)
	}
	var corpus domain.EvaluationCorpus
	if err := json.Unmarshal(data, &corpus); err != nil {
		return fmt.Errorf("corpus JSON inválido em %s: %w", *corpusPath, err)
	}
	report, err := usecase.EvaluateExpansionStrategy(corpus, strings.ToLower(*strategyName), strategy)
	if err != nil {
		return fmt.Errorf("não foi possível avaliar corpus: %w", err)
	}
	return writeIdeaJSON(report, *outputPath, stdout)
}
