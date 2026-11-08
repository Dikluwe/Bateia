/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_web.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

func runIdeaWeb(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: idea web <discover-observable|experiment|replay|compare>")
	}
	switch args[0] {
	case "discover-observable":
		return runWebDiscoverObservable(args[1:], stdout, provider)
	case "experiment":
		return runWebExperiment(args[1:], stdout, provider)
	case "authenticated-experiment":
		return runWebAuthenticatedExperiment(args[1:], stdout, provider)
	case "replay":
		return runWebReplay(args[1:], stdout)
	case "compare":
		return runWebCompare(args[1:], stdout)
	default:
		return fmt.Errorf("operação web desconhecida")
	}
}

func runWebExperiment(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	return runWebExperimentMode(args, stdout, provider, false)
}

func runWebExperimentMode(args []string, stdout io.Writer, provider domain.ObservableWebSearch, authenticated bool) error {
	fs := flag.NewFlagSet("web experiment", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	nominalPath := fs.String("nominal", "", "queries N")
	f1Path := fs.String("functional-v1", "", "queries F1")
	f2Path := fs.String("functional-v2", "", "queries F2")
	outputDir := fs.String("output-dir", "", "diretório de snapshots")
	defaultMax := 5
	if authenticated {
		defaultMax = 10
	}
	maxPerQuery := fs.Int("max-per-query", defaultMax, "máximo por query")
	noToken := fs.Bool("no-token", false, "proíbe autenticação")
	authenticatedFlag := fs.Bool("authenticated", false, "autoriza autenticação por variável de ambiente")
	authorizationPath := fs.String("authorization", "", "autorização redigida do experimento")
	safeRemaining := fs.Int("safe-remaining", 10, "margem segura da cota")
	resumeDir := fs.String("resume-dir", "", "snapshots autenticados anteriores para retomada")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := usecase.ValidateWebExperimentAdmission(usecase.WebExperimentAdmission{Authenticated: authenticated, AuthenticatedFlag: *authenticatedFlag, NoToken: *noToken, HasAuthorization: *authorizationPath != "", SafeRemaining: *safeRemaining, MaxPerQuery: *maxPerQuery}); err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("experiment exige provedor injetado")
	}
	if *nominalPath == "" || *f1Path == "" || *f2Path == "" || *outputDir == "" {
		return fmt.Errorf("experiment exige nominal, functional-v1, functional-v2 e output-dir")
	}
	var authorization *domain.WebExperimentAuthorization
	if authenticated {
		var value domain.WebExperimentAuthorization
		if err := readOntologyJSON(*authorizationPath, &value); err != nil {
			return err
		}
		if err := validateStep0020Authorization(value); err != nil {
			return err
		}
		authorization = &value
	}
	arms := []struct {
		name, path string
		exp        domain.FunctionalGrammarExpansion
	}{{name: "N", path: *nominalPath}, {name: "F1", path: *f1Path}, {name: "F2", path: *f2Path}}
	for i := range arms {
		if err := readOntologyJSON(arms[i].path, &arms[i].exp); err != nil {
			return err
		}
	}
	input := usecase.WebExperimentInput{Authenticated: authenticated, Authorization: authorization, SafeRemaining: *safeRemaining, MaxPerQuery: *maxPerQuery, Resume: *resumeDir != ""}
	for _, arm := range arms {
		input.Arms = append(input.Arms, usecase.WebExperimentArm{Name: arm.name, Expansion: arm.exp})
	}
	execution, err := usecase.ExecuteWebExperiment(commandContext(context.Background()), input, provider, webExperimentResumeFiles{directory: *resumeDir})
	if err != nil {
		return err
	}
	snapshots, receipt := execution.Snapshots, execution.Receipt
	for i := range snapshots {
		if err := rejectSensitiveWebArtifact(snapshots[i]); err != nil {
			return err
		}
	}

	if err := rejectSensitiveWebArtifact(receipt); err != nil {
		return err
	}
	if err := commandHost().MkdirAll(*outputDir, 0755); err != nil {
		return err
	}
	for i := range snapshots {
		path := filepath.Join(*outputDir, "snapshot-"+strings.ToLower(snapshots[i].Arm)+".json")
		if err := writeIdeaJSON(snapshots[i], path, stdout); err != nil {
			return err
		}
	}
	if authenticated {
		if err := writeIdeaJSON(receipt, filepath.Join(*outputDir, "receipt.json"), stdout); err != nil {
			return err
		}
	}
	return nil
}

func runWebDiscoverObservable(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	fs := flag.NewFlagSet("web discover-observable", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	queriesPath := fs.String("queries", "", "FunctionalExpansion")
	arm := fs.String("arm", "F2", "braço N, F1 ou F2")
	maxPerQuery := fs.Int("max-per-query", 5, "máximo por query")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *queriesPath == "" {
		return fmt.Errorf("discover-observable exige queries")
	}
	if provider == nil {
		return fmt.Errorf("discover-observable exige provedor injetado por L4")
	}
	var exp domain.FunctionalGrammarExpansion
	if err := readOntologyJSON(*queriesPath, &exp); err != nil {
		return err
	}

	snapshot := usecase.DiscoverWebQueries(commandContext(context.Background()), exp, *arm, *maxPerQuery, provider)

	return writeIdeaJSON(snapshot, *output, stdout)
}

type webExperimentResumeFiles struct{ directory string }

func (r webExperimentResumeFiles) PreviousSnapshot(arm string) (domain.WebSearchSnapshot, error) {
	var result domain.WebSearchSnapshot
	err := readOntologyJSON(filepath.Join(r.directory, "snapshot-"+strings.ToLower(arm)+".json"), &result)
	return result, err
}
func querySetHash(queries []domain.RealizedFunctionalQuery) string {
	return usecase.WebExperimentQuerySetHash(queries)
}

func runWebReplay(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("web replay", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	snapshotPath := fs.String("snapshot", "", "WebSearchSnapshot")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *snapshotPath == "" {
		return fmt.Errorf("replay exige snapshot")
	}
	var snap domain.WebSearchSnapshot
	if err := readOntologyJSON(*snapshotPath, &snap); err != nil {
		return err
	}
	res, err := usecase.ReplayWebSnapshot(snap)
	if err != nil {
		return err
	}
	return writeIdeaJSON(res, *output, stdout)
}

func runWebCompare(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("web compare", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	nominalPath := fs.String("nominal", "", "Nominal Snapshot")
	f1Path := fs.String("functional-v1", "", "F1 Snapshot")
	f2Path := fs.String("functional-v2", "", "F2 Snapshot")
	judgmentsPath := fs.String("judgments", "", "Judgments JSON")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *nominalPath == "" || *f1Path == "" || *f2Path == "" || *judgmentsPath == "" {
		return fmt.Errorf("compare exige nominal, functional-v1, functional-v2 e judgments")
	}

	var nom, f1, f2 domain.WebSearchSnapshot
	var j domain.BlindedJudgmentsFile
	if err := readOntologyJSON(*nominalPath, &nom); err != nil {
		return err
	}
	if err := readOntologyJSON(*f1Path, &f1); err != nil {
		return err
	}
	if err := readOntologyJSON(*f2Path, &f2); err != nil {
		return err
	}
	if err := readOntologyJSON(*judgmentsPath, &j); err != nil {
		return err
	}

	report, err := usecase.CompareGrammarExperiments(nom, f1, f2, j)
	if err != nil {
		return err
	}
	return writeIdeaJSON(report, *output, stdout)
}
