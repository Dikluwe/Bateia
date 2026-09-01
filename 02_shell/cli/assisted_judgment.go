/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/assisted_judgment.md
 * @layer L2
 * @updated 2026-08-31
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"flag"
	"fmt"
	"io"
)

type repeatedSnapshotFlags []string

func (f *repeatedSnapshotFlags) String() string         { return fmt.Sprint([]string(*f)) }
func (f *repeatedSnapshotFlags) Set(value string) error { *f = append(*f, value); return nil }

func RunAssistedJudgmentCommand(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: idea judge <prepare|seal>")
	}
	switch args[0] {
	case "prepare":
		return runAssistedPrepare(args[1:], stdout)
	case "seal":
		return runAssistedSeal(args[1:], stdout)
	default:
		return fmt.Errorf("operação de julgamento desconhecida")
	}
}

func runAssistedPrepare(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea judge prepare", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var paths repeatedSnapshotFlags
	fs.Var(&paths, "snapshot", "snapshot de entrada; pode ser repetido")
	output := fs.String("output", "", "pacote cego")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("prepare exige ao menos um snapshot")
	}
	var snapshots []domain.WebSearchSnapshot
	for _, path := range paths {
		var snapshot domain.WebSearchSnapshot
		if err := readOntologyJSON(path, &snapshot); err != nil {
			return err
		}
		snapshots = append(snapshots, snapshot)
	}
	pkg, err := usecase.PrepareBlindCandidateUnion(snapshots)
	if err != nil {
		return err
	}
	return writeIdeaJSON(pkg, *output, stdout)
}

func runAssistedSeal(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea judge seal", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	packagePath := fs.String("package", "", "pacote cego")
	judgmentsPath := fs.String("judgments", "", "julgamentos da IA")
	reviewsPath := fs.String("reviews", "", "revisões humanas")
	output := fs.String("output", "", "oráculo final")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *packagePath == "" || *judgmentsPath == "" {
		return fmt.Errorf("seal exige package e judgments")
	}
	var pkg domain.BlindCandidatePackage
	var judgments []domain.AssistedJudgment
	var reviews []domain.HumanReview
	if err := readOntologyJSON(*packagePath, &pkg); err != nil {
		return err
	}
	if err := readOntologyJSON(*judgmentsPath, &judgments); err != nil {
		return err
	}
	if *reviewsPath != "" {
		if err := readOntologyJSON(*reviewsPath, &reviews); err != nil {
			return err
		}
	}
	return writeIdeaJSON(usecase.ValidateAssistedJudgments(pkg, judgments, reviews), *output, stdout)
}
