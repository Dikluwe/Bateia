/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaEvidence(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea evidence <plan|collect|assess|evaluate|diagnose|calibrate>")
	}
	switch args[0] {
	case "plan":
		return runIdeaEvidencePlan(args[1:], stdout)
	case "assess":
		return runIdeaEvidenceAssess(args[1:], stdout)
	case "collect":
		return fmt.Errorf("evidence collect requer o provedor de conteúdo configurado pelo executável")
	case "evaluate":
		return runIdeaEvidenceEvaluate(args[1:], stdout)
	case "diagnose":
		return runIdeaEvidenceDiagnose(args[1:], stdout)
	case "calibrate":
		return runIdeaEvidenceCalibrate(args[1:], stdout)
	default:
		return fmt.Errorf("subcomando evidence desconhecido %q", args[0])
	}
}

func runIdeaEvidencePlan(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea evidence plan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	screeningPath := fs.String("screening", "", "Relatório de peneira")
	ideaPath := fs.String("idea", "", "IdeaSpec usada para orientar os alvos")
	rulesPath := fs.String("rules", "", "Catálogo de regras; omita para catálogo embutido")
	includeRepositories := fs.String("include", "", "Controles owner/repo separados por vírgula")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence plan inválidos: %w", err)
	}
	if *screeningPath == "" || *ideaPath == "" {
		return fmt.Errorf("evidence plan exige --idea e --screening")
	}
	idea, err := readIdeaSpec(*ideaPath)
	if err != nil {
		return err
	}
	var report domain.ScreeningReport
	if err := readJSONFile(*screeningPath, &report); err != nil {
		return err
	}
	catalog := usecase.DefaultEvidenceRules()
	if *rulesPath != "" {
		if err := readJSONFile(*rulesPath, &catalog); err != nil {
			return err
		}
	}
	var includes []string
	for _, value := range strings.Split(*includeRepositories, ",") {
		if value = strings.TrimSpace(value); value != "" {
			includes = append(includes, value)
		}
	}
	plans, err := usecase.PlanEvidenceForIdeaAndRepositories(idea, report, catalog, includes)
	if err != nil {
		return fmt.Errorf("não foi possível planejar evidências: %w", err)
	}
	return writeIdeaJSON(plans, *outputPath, stdout)
}

func RunIdeaEvidenceCollectCommand(args []string, stdout io.Writer, provider domain.ContentProvider) error {
	fs := flag.NewFlagSet("idea evidence collect", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	planPath := fs.String("plan", "", "Plano de extração")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	timeout := fs.Duration("timeout", 2*time.Minute, "Tempo máximo da coleta")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence collect inválidos: %w", err)
	}
	if *planPath == "" {
		return fmt.Errorf("evidence collect exige --plan")
	}
	data, err := os.ReadFile(*planPath)
	if err != nil {
		return fmt.Errorf("não foi possível ler %s: %w", *planPath, err)
	}
	var plans []domain.EvidenceExtractionPlan
	if err = json.Unmarshal(data, &plans); err != nil {
		var one domain.EvidenceExtractionPlan
		if oneErr := json.Unmarshal(data, &one); oneErr != nil {
			return fmt.Errorf("plano JSON inválido: %w", err)
		}
		plans = []domain.EvidenceExtractionPlan{one}
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	bundles := make([]domain.EvidenceBundle, 0, len(plans))
	for _, plan := range plans {
		bundle, collectErr := usecase.CollectEvidence(ctx, plan, provider)
		if collectErr != nil {
			bundle = domain.EvidenceBundle{SchemaVersion: "1", Repository: plan.Repository, Plan: plan, Complete: false, TreeComplete: false, Failures: []domain.ExecutionFailure{{Scope: plan.Repository.Key(), Stage: "resolution", Message: collectErr.Error()}}}
		}
		bundles = append(bundles, bundle)
	}
	return writeIdeaJSON(bundles, *outputPath, stdout)
}

func runIdeaEvidenceEvaluate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea evidence evaluate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	assessmentPath := fs.String("assessment", "", "Relatório de assessments")
	evidencePath := fs.String("evidence", "", "Bundles congelados")
	judgmentsPath := fs.String("judgments", "", "Julgamentos manuais versionados")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence evaluate inválidos: %w", err)
	}
	if *assessmentPath == "" || *evidencePath == "" || *judgmentsPath == "" {
		return fmt.Errorf("evidence evaluate exige --assessment, --evidence e --judgments")
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
	return writeIdeaJSON(usecase.EvaluateEvidenceMetrics(report, bundles, judgments), *outputPath, stdout)
}

func runIdeaEvidenceAssess(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea evidence assess", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	ideaPath := fs.String("idea", "", "IdeaSpec")
	screeningPath := fs.String("screening", "", "Relatório de peneira")
	evidencePath := fs.String("evidence", "", "Bundle congelado de evidências")
	rulesPath := fs.String("rules", "", "Catálogo de regras; omita para catálogo embutido")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence assess inválidos: %w", err)
	}
	if *ideaPath == "" || *screeningPath == "" || *evidencePath == "" {
		return fmt.Errorf("evidence assess exige --idea, --screening e --evidence")
	}
	idea, err := readIdeaSpec(*ideaPath)
	if err != nil {
		return err
	}
	var screening domain.ScreeningReport
	if err := readJSONFile(*screeningPath, &screening); err != nil {
		return err
	}
	data, err := os.ReadFile(*evidencePath)
	if err != nil {
		return fmt.Errorf("não foi possível ler %s: %w", *evidencePath, err)
	}
	var bundles []domain.EvidenceBundle
	if err = json.Unmarshal(data, &bundles); err != nil {
		var one domain.EvidenceBundle
		if oneErr := json.Unmarshal(data, &one); oneErr != nil {
			return fmt.Errorf("bundle JSON inválido: %w", err)
		}
		bundles = []domain.EvidenceBundle{one}
	}
	catalog := usecase.DefaultEvidenceRules()
	if *rulesPath != "" {
		if err := readJSONFile(*rulesPath, &catalog); err != nil {
			return err
		}
	}
	report, err := usecase.AssessEvidence(idea, screening, bundles, catalog)
	if err != nil {
		return fmt.Errorf("avaliação de evidências inválida: %w", err)
	}
	return writeIdeaJSON(report, *outputPath, stdout)
}
