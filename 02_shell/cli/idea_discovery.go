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
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

// RunIdeaDiscoverCommand executa a única operação remota da família idea.
func RunIdeaDiscoverCommand(args []string, stdout, stderr io.Writer, lexical usecase.SemanticExpander, provider domain.ObservableDiscoveryProvider) error {
	fs := flag.NewFlagSet("idea discover", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inputPath := fs.String("input", "", "Arquivo JSON contendo IdeaSpec")
	strategyName := fs.String("strategy", "lexical", "Estratégia literal ou lexical")
	providerName := fs.String("provider", "sourcegraph", "Provedor de descoberta")
	language := fs.String("language", "", "Linguagem opcional")
	maxPerQuery := fs.Int("max-per-query", 10, "Máximo de candidatos por expansão")
	concurrency := fs.Int("concurrency", 2, "Concorrência entre expansões")
	outputPath := fs.String("output", "", "Arquivo obrigatório para o snapshot")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea discover inválidos: %w", err)
	}
	if *inputPath == "" {
		return fmt.Errorf("idea discover exige --input")
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
	snapshot, err := usecase.DiscoverExpansionPlan(context.Background(), enriched.Plan, *strategyName, *language, *maxPerQuery, *concurrency, provider)
	if err != nil {
		return fmt.Errorf("falha na descoberta: %w", err)
	}
	if err := writeIdeaJSON(snapshot, *outputPath, stdout); err != nil {
		return err
	}
	if stderr != nil {
		_, _ = fmt.Fprintf(stderr, "snapshot gravado em %s\n", *outputPath)
	}
	return nil
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
		data, err := os.ReadFile(path)
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
