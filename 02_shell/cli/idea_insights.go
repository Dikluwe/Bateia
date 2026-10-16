/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_insights.md
 * @layer L2
 * @updated 2026-10-06
 */

package cli

import (
	"context"
	"fmt"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaInsights(args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "analyze" {
		return fmt.Errorf("uso: bateia idea insights analyze --input pedido.json [--output resultado.json]")
	}
	input, output, err := parseDecisionTraceArgs(args[1:])
	if err != nil {
		return fmt.Errorf("parâmetros de insights inválidos: %w", err)
	}
	if err := validateDecisionTraceOutput(output, input); err != nil {
		return err
	}
	var request domain.InsightInvestigationRequest
	if err := readStrictIntakeJSON(input, &request); err != nil {
		return err
	}
	report, err := usecase.InvestigateInsights(context.Background(), request)
	if err != nil {
		return fmt.Errorf("investigação de insights rejeitada: %w", err)
	}
	return writeDecisionTraceJSON(report, output, stdout, input)
}
