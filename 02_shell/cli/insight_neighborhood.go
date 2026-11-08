/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_neighborhood.md
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
	"strconv"
	"time"
)

func RunInsightNeighborhoodCommand(args []string, stdout io.Writer, factory domain.InsightSiteFactory) error {
	fs := flag.NewFlagSet("idea insights expand", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	request := domain.InsightNeighborhoodRequest{}
	input, output, memoryPath := "", "", ""
	seen := map[string]bool{}
	once := func(name string, assign func(string) error) func(string) error {
		return func(value string) error {
			if seen[name] {
				return fmt.Errorf("flag de expansão duplicada")
			}
			seen[name] = true
			return assign(value)
		}
	}
	fs.BoolFunc("all-trails", "Voltar à fronteira geral", once("all-trails", func(value string) error {
		parsed, err := strconv.ParseBool(value)
		request.AllTrails = parsed
		return err
	}))
	for _, option := range []struct {
		name string
		dest *string
	}{{"input", &input}, {"output", &output}, {"memory", &memoryPath}, {"focus", &request.Focus}} {
		fs.Func(option.name, option.name, once(option.name, func(value string) error {
			if value == "" {
				return fmt.Errorf("caminho ou foco vazio")
			}
			*option.dest = value
			return nil
		}))
	}
	fs.Func("follow-links", "Limite de leituras", once("follow-links", func(value string) error {
		n, err := strconv.Atoi(value)
		if err == nil {
			request.FollowLinks = &n
		}
		return err
	}))
	fs.Func("max-passages", "Limite de passagens", once("max-passages", func(value string) error {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return fmt.Errorf("limite de passagens inválido")
		}
		request.MaxPassages = n
		return nil
	}))
	for _, option := range []struct {
		name string
		dest *time.Duration
	}{{"timeout", &request.TotalTimeout}, {"attempt-timeout", &request.AttemptTimeout}} {
		fs.Func(option.name, "Deadline", once(option.name, func(value string) error {
			n, err := time.ParseDuration(value)
			if err != nil || n <= 0 {
				return fmt.Errorf("deadline inválido")
			}
			*option.dest = n
			return nil
		}))
	}
	if err := fs.Parse(args); err != nil || fs.NArg() != 0 {
		return fmt.Errorf("parâmetros de expansão de insights inválidos")
	}
	if output == "" {
		return fmt.Errorf("insights expand exige --output antes de acessar a rede")
	}
	if input == "" {
		return fmt.Errorf("insights expand exige --input com relatório anterior")
	}
	if err := validateDecisionTraceOutput(output, input, memoryPath); err != nil {
		return err
	}
	if err := readStrictIntakeJSON(input, &request.Input); err != nil {
		return err
	}
	if seen["memory"] {
		var err error
		request.Memory, err = readInsightMemory(memoryPath)
		if err != nil {
			return err
		}
	}
	request, err := usecase.ValidateInsightNeighborhood(request)
	if err != nil {
		return fmt.Errorf("expansão rejeitada: relatório, foco ou orçamento inválido")
	}
	if factory == nil {
		return fmt.Errorf("insights expand exige factory configurada pelo executável")
	}
	bindings, err := factory()
	if err != nil {
		return fmt.Errorf("configuração de fontes indisponível")
	}
	ctx, stop := commandHost().NotifyContext(commandContext(context.Background()))
	defer stop()
	report, expansionErr := usecase.ExpandInsightNeighborhood(ctx, request, bindings)
	if report.ID == "" {
		return fmt.Errorf("expansão não iniciada: %w", expansionErr)
	}
	if err := writeDecisionTraceJSON(report, output, stdout, input, memoryPath); err != nil {
		return err
	}
	if expansionErr != nil {
		return fmt.Errorf("expansão interrompida; resultados preservados: %w", expansionErr)
	}
	return nil
}
