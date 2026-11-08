/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_exploration.md
 * @layer L2
 * @updated 2026-11-08
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"flag"
	"fmt"
	"io"
	"strconv"
)

func runInsightExploration(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea insights explore", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	question, bundles, output, maximum := "", "", "", 8
	seen := map[string]bool{}
	once := func(name string, assign func(string) error) func(string) error {
		return func(value string) error {
			if seen[name] {
				return fmt.Errorf("flag de exploração duplicada")
			}
			seen[name] = true
			return assign(value)
		}
	}
	fs.Func("question", "Pergunta a explorar", once("question", func(value string) error { question = value; return nil }))
	fs.Func("bundles", "Lista JSON emitida por evidence collect", once("bundles", func(value string) error { bundles = value; return nil }))
	fs.Func("output", "Saída JSON; omita para stdout", once("output", func(value string) error { output = value; return nil }))
	fs.Func("max-passages", "Máximo de trechos de leitura", once("max-passages", func(value string) error {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			maximum = parsed
		}
		return err
	}))
	algorithmIdeas, review := false, false
	for _, option := range []struct {
		name, usage string
		dest        *bool
	}{{"algorithm-ideas", "Investigar princípios, condições e decisões algorítmicas", &algorithmIdeas}, {"review-alternatives", "Comparar decisões na mesma fonte e revisar cobertura", &review}} {
		fs.BoolFunc(option.name, option.usage, once(option.name, func(value string) error {
			parsed, err := strconv.ParseBool(value)
			if err == nil {
				*option.dest = parsed
			}
			return err
		}))
	}
	familiar := []string{}
	fs.Func("familiar", "Abordagem habitual; flag repetível", func(value string) error { familiar = append(familiar, value); return nil })
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de exploração inválidos: %w", err)
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("explore não aceita argumentos posicionais")
	}
	if err := validateDecisionTraceOutput(output, bundles); err != nil {
		return err
	}
	request := domain.InsightExplorationRequest{SchemaVersion: "1", Question: question, FamiliarApproaches: familiar, MaxPassages: maximum, ReviewAlternatives: review, AlgorithmIdeas: algorithmIdeas}
	if bundles != "" {
		if err := readStrictIntakeJSON(bundles, &request.Bundles); err != nil {
			return err
		}
	}
	report, err := usecase.ExploreInsights(context.Background(), request)
	if err != nil {
		return fmt.Errorf("exploração de insights rejeitada: %w", err)
	}
	return writeDecisionTraceJSON(report, output, stdout, bundles)
}
