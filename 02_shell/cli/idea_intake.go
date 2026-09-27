/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_intake.md
 * @layer L2
 * @updated 2026-09-24
 */

package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaIntake(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea intake <prepare|apply|diagnose|compile> [opções]")
	}
	switch args[0] {
	case "prepare":
		return runIdeaIntakePrepare(args[1:], stdout)
	case "apply":
		return runIdeaIntakeApply(args[1:], stdout)
	case "diagnose":
		return runIdeaIntakeDiagnose(args[1:], stdout)
	case "compile":
		return runIdeaIntakeCompile(args[1:], stdout)
	default:
		return fmt.Errorf("subcomando idea intake desconhecido %q; use prepare, apply, diagnose ou compile", args[0])
	}
}

func runIdeaIntakePrepare(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea intake prepare", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "Rascunho JSON da entrada avançada")
	output := fs.String("output", "", "Arquivo JSON de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea intake prepare inválidos: %w", err)
	}
	if *input == "" {
		return fmt.Errorf("idea intake prepare exige --input")
	}
	var draft domain.RecommendationIntakeDraft
	if err := readStrictIntakeJSON(*input, &draft); err != nil {
		return err
	}
	return writeIdeaJSON(usecase.PrepareRecommendationIntake(draft), *output, stdout)
}

func runIdeaIntakeApply(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea intake apply", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "Entrada preparada por idea intake prepare")
	answersPath := fs.String("answers", "", "Lista JSON de respostas tipadas")
	output := fs.String("output", "", "Arquivo JSON de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea intake apply inválidos: %w", err)
	}
	if *input == "" || *answersPath == "" {
		return fmt.Errorf("idea intake apply exige --input e --answers")
	}
	var intake domain.RecommendationIntake
	if err := readStrictIntakeJSON(*input, &intake); err != nil {
		return err
	}
	var answers []domain.InformationAnswer
	if err := readStrictIntakeJSON(*answersPath, &answers); err != nil {
		return err
	}
	updated, err := usecase.ApplyInformationAnswers(intake, answers)
	if err != nil {
		return fmt.Errorf("respostas rejeitadas: %w", err)
	}
	return writeIdeaJSON(updated, *output, stdout)
}

func runIdeaIntakeDiagnose(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea intake diagnose", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "Entrada canônica preparada")
	diagnosisPath := fs.String("diagnosis", "", "Diagnóstico de inconclusão JSON")
	output := fs.String("output", "", "Arquivo JSON de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea intake diagnose inválidos: %w", err)
	}
	if *input == "" || *diagnosisPath == "" {
		return fmt.Errorf("idea intake diagnose exige --input e --diagnosis")
	}
	var intake domain.RecommendationIntake
	if err := readStrictIntakeJSON(*input, &intake); err != nil {
		return err
	}
	var diagnosis domain.InconclusiveDiagnosis
	if err := readStrictIntakeJSON(*diagnosisPath, &diagnosis); err != nil {
		return err
	}
	return writeIdeaJSON(usecase.CompileInconclusiveNeeds(diagnosis, intake), *output, stdout)
}

func readStrictIntakeJSON(path string, target any) error {
	data, err := readIdeaInput(path)
	if err != nil {
		return fmt.Errorf("não foi possível ler %s: %w", path, err)
	}
	if err := validateIdeaJSONStructure(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("JSON de intake inválido em %s: %w", path, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("JSON de intake inválido em %s: esperado um único valor", path)
	}
	return nil
}
