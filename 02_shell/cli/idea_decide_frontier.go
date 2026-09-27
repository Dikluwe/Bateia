/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_decide_frontier.md
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

func runIdeaDecideFrontier(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea decide frontier", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "requisição JSON multicritério")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de decide frontier inválidos: %w", err)
	}
	if *input == "" {
		return fmt.Errorf("decide frontier exige --input")
	}
	var request domain.MulticriteriaDecisionRequest
	if err := readStrictIntakeJSON(*input, &request); err != nil {
		return err
	}
	report, err := usecase.AnalyzeMulticriteriaDecision(request)
	if err != nil {
		return fmt.Errorf("análise multicritério rejeitada: %w", err)
	}
	return writeIdeaJSON(report, *output, stdout)
}
