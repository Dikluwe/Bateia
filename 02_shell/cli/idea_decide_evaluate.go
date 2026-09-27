/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_decide_evaluate.md
 * @layer L2
 * @updated 2026-09-24
 */
package cli

import (
	"flag"
	"fmt"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaDecideEvaluate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea decide evaluate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "casos adjudicados de decisão final em JSON")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de decide evaluate inválidos: %w", err)
	}
	if *input == "" || fs.NArg() != 0 {
		return fmt.Errorf("decide evaluate exige --input e não aceita argumentos posicionais")
	}
	var request domain.SelectiveDecisionRequest
	if err := readStrictIntakeJSON(*input, &request); err != nil {
		return err
	}
	report, err := usecase.EvaluateSelectiveDecisions(request)
	if err != nil {
		return fmt.Errorf("avaliação seletiva rejeitada: %w", err)
	}
	return writeIdeaJSON(report, *output, stdout)
}
