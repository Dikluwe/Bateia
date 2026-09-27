/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_recheck.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func RunIdeaEvidenceRecheckCommand(args []string, stdout io.Writer, provider domain.ContentProvider) error {
	fs := flag.NewFlagSet("idea evidence recheck", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "TechnicalSourceRecheckRequest JSON")
	asOfText := fs.String("as-of", "", "Data RFC3339 de corte")
	output := fs.String("output", "", "Relatório; omita para stdout")
	bundlesOutput := fs.String("bundles-output", "", "Bundles relidos opcionais para idea decide")
	timeout := fs.Duration("timeout", time.Minute, "Prazo da releitura pelo provedor")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *input == "" || *timeout <= 0 || *timeout > 5*time.Minute {
		return fmt.Errorf("evidence recheck exige --input e timeout entre 0 e 5 minutos")
	}
	paths := []string{*input, *output, *bundlesOutput}
	for i, path := range paths {
		for _, other := range paths[i+1:] {
			if path != "" && other != "" && sameTechnicalOutputPath(path, other) {
				return fmt.Errorf("--input, --output e --bundles-output devem ser distintos")
			}
		}
	}
	var req domain.TechnicalSourceRecheckRequest
	if err := readDecisionJSON(*input, &req); err != nil {
		return err
	}
	asOf, err := parseTechnicalEvidenceAsOf(*asOfText)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	report, err := usecase.RecheckTechnicalSources(ctx, req, provider, asOf)
	if err != nil {
		return fmt.Errorf("não foi possível reler fontes: %w", err)
	}
	if *bundlesOutput != "" {
		if err := writeIdeaJSON(report.Bundles, *bundlesOutput, io.Discard); err != nil {
			return err
		}
	}
	return writeIdeaJSON(report, *output, stdout)
}

func sameTechnicalOutputPath(a, b string) bool {
	aa, errA := filepath.Abs(a)
	bb, errB := filepath.Abs(b)
	if errA == nil && errB == nil && aa == bb {
		return true
	}
	sa, errA := os.Stat(a)
	sb, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(sa, sb)
}
