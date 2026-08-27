/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_decide.md
 * @layer L2
 * @updated 2026-08-26
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

func runIdeaDecide(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea decide", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	ideaPath := fs.String("idea", "", "IdeaSpec JSON")
	assessmentPath := fs.String("assessment", "", "AssessmentReport JSON")
	policyPath := fs.String("policy", "", "ReuseDecisionPolicy JSON")
	outputPath := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea decide inválidos: %w", err)
	}
	if *ideaPath == "" || *assessmentPath == "" || *policyPath == "" {
		return fmt.Errorf("idea decide exige --idea, --assessment e --policy")
	}
	idea, err := readIdeaSpec(*ideaPath)
	if err != nil {
		return err
	}
	var assessment domain.AssessmentReport
	if err := readDecisionJSON(*assessmentPath, &assessment); err != nil {
		return err
	}
	var policy domain.ReuseDecisionPolicy
	if err := readDecisionJSON(*policyPath, &policy); err != nil {
		return err
	}
	report, err := usecase.DecideReuse(domain.ReuseDecisionRequest{Idea: idea, Assessments: assessment.Assessments, Policy: policy})
	if err != nil {
		return fmt.Errorf("não foi possível decidir reúso: %w", err)
	}
	return writeIdeaJSON(report, *outputPath, stdout)
}

func readDecisionJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("não foi possível ler %s: %w", path, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("JSON inválido em %s: %w", path, err)
	}
	return nil
}
