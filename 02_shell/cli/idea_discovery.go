/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_discovery.md
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
	"os/signal"
	"strings"
	"syscall"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

// RunIdeaDiscoverCommand executa a única operação remota da família idea.
func RunIdeaDiscoverCommand(args []string, stdout, stderr io.Writer, lexical usecase.SemanticExpander, provider domain.ObservableDiscoveryProvider) error {
	fs := flag.NewFlagSet("idea discover", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inputPath := fs.String("input", "", "Arquivo JSON contendo IdeaSpec")
	planPath := fs.String("plan", "", "ExpansionPlan auditado; mutuamente exclusivo com --input")
	strategyName := fs.String("strategy", "lexical", "Estratégia literal ou lexical")
	providerName := fs.String("provider", "sourcegraph", "Provedor de descoberta")
	language := fs.String("language", "", "Linguagem opcional")
	maxPerQuery := fs.Int("max-per-query", 10, "Máximo de candidatos por expansão")
	concurrency := fs.Int("concurrency", 2, "Concorrência entre expansões")
	totalTimeout := fs.Duration("timeout", 3*time.Minute, "Deadline total da descoberta")
	attemptTimeout := fs.Duration("attempt-timeout", 20*time.Second, "Timeout por tentativa do provedor")
	outputPath := fs.String("output", "", "Arquivo obrigatório para o snapshot")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea discover inválidos: %w", err)
	}
	if (*inputPath == "") == (*planPath == "") {
		return fmt.Errorf("idea discover exige exatamente um de --input ou --plan")
	}
	if *outputPath == "" {
		return fmt.Errorf("idea discover exige --output antes de acessar a rede")
	}
	if strings.ToLower(*providerName) != "sourcegraph" {
		return fmt.Errorf("provedor %q desconhecido; use sourcegraph", *providerName)
	}
	if *maxPerQuery < 1 || *maxPerQuery > 25 {
		return fmt.Errorf("--max-per-query deve estar entre 1 e 25")
	}
	if *concurrency < 1 || *concurrency > 4 {
		return fmt.Errorf("--concurrency deve estar entre 1 e 4")
	}
	if *totalTimeout <= 0 || *attemptTimeout <= 0 || *attemptTimeout > *totalTimeout {
		return fmt.Errorf("--timeout e --attempt-timeout devem ser positivos e attempt não pode exceder total")
	}
	var plan domain.ExpansionPlan
	proposedCount := 0
	if *planPath != "" {
		data, err := readIdeaJSONInput(*planPath)
		if err != nil {
			return fmt.Errorf("não foi possível ler plano %s: %w", *planPath, err)
		}
		if err := json.Unmarshal(data, &plan); err != nil {
			return fmt.Errorf("ExpansionPlan JSON inválido em %s: %w", *planPath, err)
		}
		if len(plan.Expansions) == 0 {
			return fmt.Errorf("ExpansionPlan não contém expansões")
		}
		*strategyName = "ontology"
	} else {
		strategy, err := selectIdeaStrategy(*strategyName, lexical)
		if err != nil {
			return err
		}
		spec, err := readIdeaSpec(*inputPath)
		if err != nil {
			return err
		}
		literal, err := usecase.ExpandIdea(spec)
		if err != nil {
			return fmt.Errorf("IdeaSpec inválida: %w", err)
		}
		enriched, err := usecase.EnrichExpansionPlan(*literal, strategy)
		if err != nil {
			return fmt.Errorf("falha na expansão %s: %w", *strategyName, err)
		}
		plan, proposedCount = enriched.Plan, enriched.ProposedCount
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	policy := domain.DiscoveryExecutionPolicy{SchemaVersion: "1", TotalTimeout: *totalTimeout, AttemptTimeout: *attemptTimeout, Concurrency: *concurrency, MaxAttempts: 3}
	reporter := discoveryCLIReporter{writer: stderr}
	snapshot, discoverErr := usecase.DiscoverExpansionPlanWithPolicy(ctx, plan, *strategyName, *language, *maxPerQuery, policy, provider, reporter)
	if snapshot.Receipt != nil && *planPath == "" && proposedCount == 0 {
		snapshot.Receipt.Warnings = append(snapshot.Receipt.Warnings, "NO_SEMANTIC_EXPANSION")
	}
	if err := writeJSONAtomic(*outputPath, snapshot); err != nil {
		return err
	}
	if stderr != nil {
		_, _ = fmt.Fprintf(stderr, "snapshot gravado em %s\n", *outputPath)
	}
	if discoverErr != nil {
		return fmt.Errorf("falha na descoberta: %w", discoverErr)
	}
	return nil
}

type discoveryCLIReporter struct{ writer io.Writer }

func (r discoveryCLIReporter) ReportDiscoveryProgress(event domain.DiscoveryProgressEvent) {
	if r.writer != nil {
		_, _ = fmt.Fprintf(r.writer, "descoberta %d/%d %s %s\n", event.Completed, event.Planned, event.ExpansionID, event.Status)
	}
}

func runIdeaDiscoveryEvaluate(args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "evaluate" {
		return fmt.Errorf("uso: bateia idea discovery evaluate --snapshot arquivo --judgments arquivo")
	}
	fs := flag.NewFlagSet("idea discovery evaluate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	snapshotPath := fs.String("snapshot", "", "Snapshot de descoberta")
	judgmentsPath := fs.String("judgments", "", "Julgamentos de relevância")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("parâmetros de replay inválidos: %w", err)
	}
	if *snapshotPath == "" || *judgmentsPath == "" {
		return fmt.Errorf("replay exige --snapshot e --judgments")
	}
	readJSON := func(path string, destination any) error {
		data, err := readIdeaJSONInput(path)
		if err != nil {
			return fmt.Errorf("não foi possível ler %s: %w", path, err)
		}
		if err := json.Unmarshal(data, destination); err != nil {
			return fmt.Errorf("JSON inválido em %s: %w", path, err)
		}
		return nil
	}
	var snapshot domain.DiscoverySnapshot
	if err := readJSON(*snapshotPath, &snapshot); err != nil {
		return err
	}
	var judgments domain.DiscoveryJudgments
	if err := readJSON(*judgmentsPath, &judgments); err != nil {
		return err
	}
	evaluation, err := usecase.EvaluateDiscovery(snapshot, judgments)
	if err != nil {
		return fmt.Errorf("avaliação inválida: %w", err)
	}
	return writeIdeaJSON(evaluation, *outputPath, stdout)
}
