/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_calibrate.md
 * @layer L2
 * @updated 2026-08-31
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"flag"
	"fmt"
	"io"
)

func runIdeaEvidenceCalibrate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea evidence calibrate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	observationsPath := fs.String("observations", "", "observações congeladas dos três splits")
	outputPath := fs.String("output", "", "relatório de calibração")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence calibrate inválidos: %w", err)
	}
	if *observationsPath == "" {
		return fmt.Errorf("evidence calibrate exige --observations")
	}
	var observations []domain.EvidenceCalibrationObservation
	if err := readJSONFile(*observationsPath, &observations); err != nil {
		return err
	}
	return writeIdeaJSON(usecase.EvaluateEvidenceCalibration(observations), *outputPath, stdout)
}
