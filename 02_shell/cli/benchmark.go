/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/benchmark.md
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

// RunBenchmarkCommand exposes the frozen benchmark without loading credentials,
// caches, providers or any other network-capable dependency.
func RunBenchmarkCommand(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("benchmark exige validate, mutate ou run")
	}
	switch args[0] {
	case "validate":
		return runBenchmarkValidate(args[1:], stdout)
	case "mutate":
		return runBenchmarkMutate(args[1:], stdout)
	case "run":
		return runBenchmark(args[1:], stdout)
	default:
		return fmt.Errorf("subcomando benchmark desconhecido: %s", args[0])
	}
}

func benchmarkInputs(fs *flag.FlagSet) (*string, *string, *string) {
	manifest := fs.String("manifest", "", "BenchmarkManifest JSON congelado")
	corpus := fs.String("corpus", "", "BenchmarkCorpus JSON congelado")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	return manifest, corpus, output
}

func readBenchmarkInputs(manifestPath, corpusPath string) (domain.BenchmarkManifest, domain.BenchmarkCorpus, error) {
	var manifest domain.BenchmarkManifest
	var corpus domain.BenchmarkCorpus
	if manifestPath == "" || corpusPath == "" {
		return manifest, corpus, fmt.Errorf("--manifest e --corpus são obrigatórios")
	}
	if err := readDecisionJSON(manifestPath, &manifest); err != nil {
		return manifest, corpus, err
	}
	if err := readDecisionJSON(corpusPath, &corpus); err != nil {
		return manifest, corpus, err
	}
	return manifest, corpus, nil
}

func runBenchmarkValidate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("benchmark validate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	manifestPath, corpusPath, output := benchmarkInputs(fs)
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de benchmark validate inválidos: %w", err)
	}
	manifest, corpus, err := readBenchmarkInputs(*manifestPath, *corpusPath)
	if err != nil {
		return err
	}
	return writeIdeaJSON(usecase.ValidateBenchmark(manifest, corpus), *output, stdout)
}

func runBenchmarkMutate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("benchmark mutate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	manifestPath := fs.String("manifest", "", "BenchmarkManifest JSON congelado")
	mutationsPath := fs.String("mutations", "", "lista JSON de BenchmarkMutation")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de benchmark mutate inválidos: %w", err)
	}
	if *manifestPath == "" || *mutationsPath == "" {
		return fmt.Errorf("--manifest e --mutations são obrigatórios")
	}
	var manifest domain.BenchmarkManifest
	var mutations []domain.BenchmarkMutation
	if err := readDecisionJSON(*manifestPath, &manifest); err != nil {
		return err
	}
	if err := readDecisionJSON(*mutationsPath, &mutations); err != nil {
		return err
	}
	return writeIdeaJSON(usecase.EvaluateBenchmarkMutations(manifest, mutations), *output, stdout)
}

func runBenchmark(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("benchmark run", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	manifestPath, corpusPath, output := benchmarkInputs(fs)
	split := fs.String("split", "", "CONTRACT, TRAIN, VALIDATION ou TEST; vazio executa todos")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de benchmark run inválidos: %w", err)
	}
	manifest, corpus, err := readBenchmarkInputs(*manifestPath, *corpusPath)
	if err != nil {
		return err
	}
	report, err := usecase.RunEndToEndBenchmark(manifest, corpus, domain.BenchmarkSplit(*split))
	if err != nil {
		return err
	}
	return writeIdeaJSON(report, *output, stdout)
}
