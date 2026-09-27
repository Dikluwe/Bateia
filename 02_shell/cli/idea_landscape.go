/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_landscape.md
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

func runIdeaLandscape(args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "analyze" {
		return fmt.Errorf("uso: bateia idea landscape analyze --input arquivo.json [--output arquivo.json]")
	}
	fs := flag.NewFlagSet("idea landscape analyze", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "snapshot de landscape em JSON")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("parâmetros de landscape analyze inválidos: %w", err)
	}
	if *input == "" || fs.NArg() != 0 {
		return fmt.Errorf("landscape analyze exige --input e não aceita argumentos posicionais")
	}
	var request domain.TechnologyLandscapeRequest
	if err := readStrictIntakeJSON(*input, &request); err != nil {
		return err
	}
	report, err := usecase.AnalyzeTechnologyLandscape(request)
	if err != nil {
		return fmt.Errorf("landscape rejeitado: %w", err)
	}
	return writeIdeaJSON(report, *output, stdout)
}
