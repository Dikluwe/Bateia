/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_plan.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

// RunIdeaCommand executa comandos offline relacionados à materialização de ideias.
func RunIdeaCommand(args []string, stdout io.Writer, semantic ...usecase.SemanticExpander) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea <plan|expand|evaluate> [opções]")
	}
	switch args[0] {
	case "plan":
		return runIdeaPlan(args[1:], stdout)
	case "expand":
		return runIdeaExpand(args[1:], stdout, firstSemanticExpander(semantic))
	case "evaluate":
		return runIdeaEvaluate(args[1:], stdout, firstSemanticExpander(semantic))
	case "ontology":
		return runIdeaOntology(args[1:], stdout)
	case "functional":
		return runIdeaFunctional(args[1:], stdout)
	case "discovery":
		return runIdeaDiscoveryEvaluate(args[1:], stdout)
	case "screen":
		return runIdeaScreen(args[1:], stdout)
	case "screening":
		return runIdeaScreening(args[1:], stdout)
	case "evidence":
		return runIdeaEvidence(args[1:], stdout)
	case "decide":
		return runIdeaDecide(args[1:], stdout)
	case "recommend":
		return RunIdeaRecommendCommand(args[1:], stdout, firstSemanticExpander(semantic), nil, nil)
	case "recover":
		return RunIdeaRecoverCommand(args[1:], stdout, nil)
	case "autofagia":
		return runIdeaAutophagy(args[1:], stdout)
	case "katamari":
		return runIdeaKatamari(args[1:], stdout)
	default:
		return fmt.Errorf("subcomando idea desconhecido %q; use plan, expand, evaluate, ontology, discovery, screen, screening, evidence, decide, recommend, recover, autofagia ou katamari", args[0])
	}
}

func firstSemanticExpander(expanders []usecase.SemanticExpander) usecase.SemanticExpander {
	if len(expanders) == 0 {
		return nil
	}
	return expanders[0]
}

func writeIdeaJSON(value any, outputPath string, stdout io.Writer) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("não foi possível serializar resultado: %w", err)
	}
	encoded = append(encoded, '\n')
	if outputPath != "" {
		if err := os.WriteFile(outputPath, encoded, 0644); err != nil {
			return fmt.Errorf("não foi possível gravar resultado %s: %w", outputPath, err)
		}
		return nil
	}
	if _, err := stdout.Write(encoded); err != nil {
		return fmt.Errorf("não foi possível escrever resultado em stdout: %w", err)
	}
	return nil
}

func readIdeaSpec(path string) (domain.IdeaSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.IdeaSpec{}, fmt.Errorf("não foi possível ler IdeaSpec %s: %w", path, err)
	}
	var spec domain.IdeaSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return domain.IdeaSpec{}, fmt.Errorf("IdeaSpec JSON inválida em %s: %w", path, err)
	}
	return spec, nil
}

func runIdeaPlan(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea plan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inputPath := fs.String("input", "", "Arquivo JSON contendo IdeaSpec")
	outputPath := fs.String("output", "", "Arquivo JSON de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea plan inválidos: %w", err)
	}
	if *inputPath == "" {
		return fmt.Errorf("idea plan exige --input com um arquivo IdeaSpec JSON")
	}

	spec, err := readIdeaSpec(*inputPath)
	if err != nil {
		return err
	}
	plan, err := usecase.ExpandIdea(spec)
	if err != nil {
		return fmt.Errorf("IdeaSpec inválida: %w", err)
	}
	return writeIdeaJSON(plan, *outputPath, stdout)
}
