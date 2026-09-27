/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_review.md
 * @layer L2
 * @updated 2026-09-24
 */

package cli

import (
	"flag"
	"fmt"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaReview(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea review <protocol|synthesize|update|prioritize|evaluate> --input arquivo.json")
	}
	switch args[0] {
	case "protocol":
		return runIdeaReviewProtocol(args[1:], stdout)
	case "synthesize":
		return runIdeaReviewSynthesize(args[1:], stdout)
	case "update":
		return runIdeaReviewUpdate(args[1:], stdout)
	case "prioritize":
		return runIdeaReviewPrioritize(args[1:], stdout)
	case "evaluate":
		return runIdeaReviewEvaluate(args[1:], stdout)
	default:
		return fmt.Errorf("subcomando idea review desconhecido %q; use protocol, synthesize, update, prioritize ou evaluate", args[0])
	}
}

func reviewIO(args []string, name string) (string, string, error) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "Arquivo JSON de entrada")
	output := fs.String("output", "", "Arquivo JSON de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return "", "", fmt.Errorf("parâmetros de %s inválidos: %w", name, err)
	}
	if *input == "" {
		return "", "", fmt.Errorf("%s exige --input", name)
	}
	return *input, *output, nil
}

func runIdeaReviewProtocol(args []string, stdout io.Writer) error {
	input, output, err := reviewIO(args, "idea review protocol")
	if err != nil {
		return err
	}
	var draft domain.ReviewProtocolDraft
	if err := readStrictIntakeJSON(input, &draft); err != nil {
		return err
	}
	protocol, err := usecase.FreezeReviewProtocol(draft)
	if err != nil {
		return fmt.Errorf("protocolo rejeitado: %w", err)
	}
	return writeIdeaJSON(protocol, output, stdout)
}

func runIdeaReviewSynthesize(args []string, stdout io.Writer) error {
	input, output, err := reviewIO(args, "idea review synthesize")
	if err != nil {
		return err
	}
	var synthesis domain.ReviewSynthesisInput
	if err := readStrictIntakeJSON(input, &synthesis); err != nil {
		return err
	}
	pack, err := usecase.SynthesizeSystematicReview(synthesis)
	if err != nil {
		return fmt.Errorf("síntese rejeitada: %w", err)
	}
	return writeIdeaJSON(pack, output, stdout)
}

func runIdeaReviewUpdate(args []string, stdout io.Writer) error {
	input, output, err := reviewIO(args, "idea review update")
	if err != nil {
		return err
	}
	var update domain.LivingReviewUpdateInput
	if err := readStrictIntakeJSON(input, &update); err != nil {
		return err
	}
	index, err := usecase.UpdateLivingReview(update)
	if err != nil {
		return fmt.Errorf("atualização viva rejeitada: %w", err)
	}
	return writeIdeaJSON(index, output, stdout)
}

func runIdeaReviewPrioritize(args []string, stdout io.Writer) error {
	input, output, err := reviewIO(args, "idea review prioritize")
	if err != nil {
		return err
	}
	var request domain.ActiveLearningRequest
	if err := readStrictIntakeJSON(input, &request); err != nil {
		return err
	}
	plan, err := usecase.PrioritizeReviewCandidates(request)
	if err != nil {
		return fmt.Errorf("priorização rejeitada: %w", err)
	}
	return writeIdeaJSON(plan, output, stdout)
}

func runIdeaReviewEvaluate(args []string, stdout io.Writer) error {
	input, output, err := reviewIO(args, "idea review evaluate")
	if err != nil {
		return err
	}
	var corpus domain.ActiveLearningCorpus
	if err := readStrictIntakeJSON(input, &corpus); err != nil {
		return err
	}
	report, err := usecase.EvaluateActiveLearning(corpus)
	if err != nil {
		return fmt.Errorf("avaliação rejeitada: %w", err)
	}
	return writeIdeaJSON(report, output, stdout)
}
