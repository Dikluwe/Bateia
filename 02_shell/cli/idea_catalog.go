/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_catalog.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func RunIdeaCatalogCommand(args []string, stdout io.Writer, provider domain.ProjectCatalog, surfaces ...domain.CatalogSurface) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea catalog <discover|screen>")
	}
	switch args[0] {
	case "discover":
		return runIdeaCatalogDiscover(args[1:], stdout, provider)
	case "screen":
		return runIdeaCatalogScreen(args[1:], stdout)
	case "discover-surfaces":
		return runIdeaDiscoverSurfaces(args[1:], stdout, surfaces)
	case "screen-multiview":
		return runIdeaScreenMultiView(args[1:], stdout)
	default:
		return fmt.Errorf("operação de catálogo desconhecida %q", args[0])
	}
}

func runIdeaDiscoverSurfaces(args []string, stdout io.Writer, providers []domain.CatalogSurface) error {
	fs := flag.NewFlagSet("idea catalog discover-surfaces", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	planPath := fs.String("plan", "", "ExpansionPlan")
	max := fs.Int("max-per-query", 3, "Máximo por consulta")
	timeout := fs.Duration("timeout", 3*time.Minute, "Deadline total")
	output := fs.String("output", "", "Snapshot obrigatório")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *planPath == "" || *output == "" {
		return fmt.Errorf("discover-surfaces exige --plan e --output")
	}
	var plan domain.ExpansionPlan
	if err := readOntologyJSON(*planPath, &plan); err != nil {
		return err
	}
	caps := []domain.CatalogSurfaceCapabilities{}
	for _, p := range providers {
		caps = append(caps, p.SurfaceCapabilities())
	}
	plans := usecase.PlanCatalogSurfaces(plan, caps)
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	snapshot, discoverErr := usecase.DiscoverCatalogSurfaces(ctx, plans, providers, *max)
	if err := writeJSONAtomic(*output, snapshot); err != nil {
		return err
	}
	return discoverErr
}

func runIdeaScreenMultiView(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea catalog screen-multiview", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	snapshotPath := fs.String("snapshot", "", "MultiViewSnapshot")
	output := fs.String("output", "", "Relatório")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *snapshotPath == "" {
		return fmt.Errorf("screen-multiview exige --snapshot")
	}
	var snapshot domain.MultiViewSnapshot
	if err := readOntologyJSON(*snapshotPath, &snapshot); err != nil {
		return err
	}
	return writeIdeaJSON(usecase.ScreenMultiView(snapshot), *output, stdout)
}

func runIdeaCatalogDiscover(args []string, stdout io.Writer, provider domain.ProjectCatalog) error {
	fs := flag.NewFlagSet("idea catalog discover", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	planPath := fs.String("plan", "", "ExpansionPlan ontológico auditado")
	providerName := fs.String("provider", "github", "Provedor PROJECT_CATALOG")
	max := fs.Int("max-per-query", 3, "Máximo por query")
	concurrency := fs.Int("concurrency", 2, "Concorrência registrada")
	timeout := fs.Duration("timeout", 3*time.Minute, "Deadline total")
	attempt := fs.Duration("attempt-timeout", 20*time.Second, "Timeout por consulta")
	output := fs.String("output", "", "Snapshot obrigatório")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *planPath == "" || *output == "" {
		return fmt.Errorf("catalog discover exige --plan e --output")
	}
	if *providerName != "github" {
		return fmt.Errorf("provedor de catálogo %q desconhecido", *providerName)
	}
	if *timeout <= 0 || *attempt <= 0 || *attempt > *timeout {
		return fmt.Errorf("timeouts inválidos")
	}
	var plan domain.ExpansionPlan
	if err := readOntologyJSON(*planPath, &plan); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	policy := domain.DiscoveryExecutionPolicy{SchemaVersion: "1", TotalTimeout: *timeout, AttemptTimeout: *attempt, Concurrency: *concurrency, MaxAttempts: 1}
	snapshot, discoverErr := usecase.DiscoverProjectCatalog(ctx, plan, *max, policy, provider)
	if err := writeJSONAtomic(*output, snapshot); err != nil {
		return err
	}
	if discoverErr != nil {
		return fmt.Errorf("descoberta de catálogo: %w", discoverErr)
	}
	return nil
}

func runIdeaCatalogScreen(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea catalog screen", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	snapshotPath := fs.String("snapshot", "", "CatalogSnapshot")
	output := fs.String("output", "", "Saída; omita para stdout")
	proofPlan := fs.String("proof-plan", "", "Plano opcional, restrito aos repositórios ADVANCE")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *snapshotPath == "" {
		return fmt.Errorf("catalog screen exige --snapshot")
	}
	b, err := readIdeaJSONInput(*snapshotPath)
	if err != nil {
		return err
	}
	var snapshot domain.CatalogSnapshot
	if err := json.Unmarshal(b, &snapshot); err != nil {
		return err
	}
	report := usecase.ScreenProjectCatalog(snapshot)
	if *proofPlan != "" {
		if err := writeIdeaJSON(usecase.BuildCatalogProofPlan(report, snapshot.Plan), *proofPlan, stdout); err != nil {
			return err
		}
	}
	return writeIdeaJSON(report, *output, stdout)
}
