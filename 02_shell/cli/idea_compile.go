/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_compile.md
 * @layer L2
 * @updated 2026-09-24
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"flag"
	"fmt"
	"io"
)

func runIdeaIntakeCompile(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea intake compile", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "Entrada preparada por idea intake prepare")
	output := fs.String("output", "", "Arquivo JSON de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea intake compile inválidos: %w", err)
	}
	if *input == "" {
		return fmt.Errorf("idea intake compile exige --input")
	}
	var intake domain.RecommendationIntake
	if err := readStrictIntakeJSON(*input, &intake); err != nil {
		return err
	}
	return writeIdeaJSON(usecase.CompileCapabilityPlan(intake), *output, stdout)
}
