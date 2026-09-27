/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recover_value.md
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

func runRecoverValue(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea recover value", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "requisição JSON de valor decisório")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de recover value inválidos: %w", err)
	}
	if *input == "" {
		return fmt.Errorf("recover value exige --input")
	}
	var request domain.DecisionValueRequest
	if err := readStrictIntakeJSON(*input, &request); err != nil {
		return err
	}
	plan, err := usecase.PlanDecisionValue(request)
	if err != nil {
		return fmt.Errorf("valor decisório rejeitado: %w", err)
	}
	return writeIdeaJSON(plan, *output, stdout)
}
