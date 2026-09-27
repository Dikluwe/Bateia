/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_reconcile.md
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

func runIdeaEvidenceReconcile(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea evidence reconcile", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inputPath := fs.String("input", "", "TechnicalEvidenceBatch JSON")
	asOfText := fs.String("as-of", "", "Data RFC3339; omita para agora")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence reconcile inválidos: %w", err)
	}
	if *inputPath == "" {
		return fmt.Errorf("evidence reconcile exige --input")
	}
	var batch domain.TechnicalEvidenceBatch
	if err := readDecisionJSON(*inputPath, &batch); err != nil {
		return err
	}
	asOf, err := parseTechnicalEvidenceAsOf(*asOfText)
	if err != nil {
		return err
	}
	review, err := usecase.ReviewTechnicalEvidence(batch, asOf)
	if err != nil {
		return fmt.Errorf("achados técnicos inválidos: %w", err)
	}
	return writeIdeaJSON(review, *outputPath, stdout)
}
