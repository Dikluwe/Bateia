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
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaEvidence(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea evidence <plan|assess>")
	}
	switch args[0] {
	case "plan":
		return runIdeaEvidencePlan(args[1:], stdout)
	case "assess":
		return runIdeaEvidenceAssess(args[1:], stdout)
	case "collect":
		return fmt.Errorf("evidence collect requer o provedor de conteúdo configurado pelo executável")
	default:
		return fmt.Errorf("subcomando evidence desconhecido %q", args[0])
	}
}

func runIdeaEvidencePlan(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea evidence plan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	screeningPath := fs.String("screening", "", "Relatório de peneira")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence plan inválidos: %w", err)
	}
	if *screeningPath == "" {
		return fmt.Errorf("evidence plan exige --screening")
	}
	var report domain.ScreeningReport
	if err := readJSONFile(*screeningPath, &report); err != nil {
		return err
	}
	plans, err := usecase.PlanEvidence(report)
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
			return fmt.Errorf("coleta de %s falhou: %w", plan.Repository.Key(), collectErr)
		}
		bundles = append(bundles, bundle)
	}
	return writeIdeaJSON(bundles, *outputPath, stdout)
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
