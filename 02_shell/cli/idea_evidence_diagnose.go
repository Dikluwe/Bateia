/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_diagnose.md
 * @layer L2
 * @updated 2026-08-26
 */

package cli

import (
	"flag"
	"fmt"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaEvidenceDiagnose(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea evidence diagnose", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	assessmentPath := fs.String("assessment", "", "Relatório de assessments")
	evidencePath := fs.String("evidence", "", "Bundles congelados")
	judgmentsPath := fs.String("judgments", "", "Julgamentos manuais versionados")
	rulesPath := fs.String("rules", "", "Catálogo de regras; omita para catálogo embutido")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence diagnose inválidos: %w", err)
	}
	if *assessmentPath == "" || *evidencePath == "" || *judgmentsPath == "" {
		return fmt.Errorf("evidence diagnose exige --assessment, --evidence e --judgments")
	}
	var report domain.AssessmentReport
	if err := readJSONFile(*assessmentPath, &report); err != nil {
		return err
	}
	var bundles []domain.EvidenceBundle
	if err := readJSONFile(*evidencePath, &bundles); err != nil {
		return err
	}
	var judgments []domain.EvidenceJudgment
	if err := readJSONFile(*judgmentsPath, &judgments); err != nil {
		return err
	}
	if err := usecase.ValidateEvidenceCoverageInputs(report, judgments); err != nil {
		return fmt.Errorf("entradas de evidence diagnose inválidas: %w", err)
	}
	catalog := usecase.DefaultEvidenceRules()
	if *rulesPath != "" {
		if err := readJSONFile(*rulesPath, &catalog); err != nil {
			return err
		}
	}
	return writeIdeaJSON(usecase.DiagnoseEvidenceCoverage(report, bundles, judgments, catalog), *outputPath, stdout)
}
