/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_memory.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
)

func readInsightMemory(path string) (*domain.InsightExplorationMemory, error) {
	var raw json.RawMessage
	if err := readStrictIntakeJSON(path, &raw); err != nil {
		return nil, err
	}
	var tag struct {
		SchemaVersion string `json:"schema_version"`
	}
	if json.Unmarshal(raw, &tag) != nil {
		return nil, fmt.Errorf("memória de exploração inválida")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var memory *domain.InsightExplorationMemory
	switch tag.SchemaVersion {
	case "insight-exploration-memory-v1":
		memory = &domain.InsightExplorationMemory{}
		if decoder.Decode(memory) != nil {
			return nil, fmt.Errorf("snapshot de memória inválido")
		}
	case "insight-site-discovery-v1":
		var report domain.InsightSiteReport
		if decoder.Decode(&report) != nil || report.ExplorationMemory == nil {
			return nil, fmt.Errorf("relatório sem memória de exploração válida")
		}
		memory = report.ExplorationMemory
	default:
		return nil, fmt.Errorf("formato de memória não suportado")
	}
	validated, err := usecase.ValidateInsightExplorationMemory(memory)
	if err != nil {
		return nil, fmt.Errorf("integridade/contrato da memória inválido")
	}
	return &validated, nil
}
func runInsightMemory(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea insights memory", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input, output, target, state, reason := "", "", "", "", ""
	seen := map[string]bool{}
	for _, option := range []struct {
		name string
		dest *string
	}{{"input", &input}, {"output", &output}, {"target", &target}, {"state", &state}, {"reason", &reason}} {
		fs.Func(option.name, option.name, func(value string) error {
			if seen[option.name] {
				return fmt.Errorf("flag de memória duplicada")
			}
			seen[option.name] = true
			*option.dest = value
			return nil
		})
	}
	if err := fs.Parse(args); err != nil || fs.NArg() != 0 || input == "" || output == "" {
		return fmt.Errorf("uso: insights memory --input arquivo --target URL --state pending|discarded --reason motivo --output memoria.json")
	}
	if err := validateDecisionTraceOutput(output, input); err != nil {
		return err
	}
	memory, err := readInsightMemory(input)
	if err != nil {
		return err
	}
	updated, err := usecase.UpdateInsightExplorationMemory(*memory, target, state, reason)
	if err != nil {
		return fmt.Errorf("alteração de memória rejeitada: alvo, estado ou motivo inválido")
	}
	return writeDecisionTraceJSON(updated, output, stdout, input)
}
