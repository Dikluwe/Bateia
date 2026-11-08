/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_trace.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

// Trace commands only describe declared dependencies; they never rerun a decision.
func runIdeaTrace(args []string, stdout io.Writer) error {
	if len(args) == 0 || (args[0] != "build" && args[0] != "impact" && args[0] != "explain" && args[0] != "compare" && args[0] != "review" && args[0] != "review-context") {
		return fmt.Errorf("uso: bateia idea trace <build|impact|explain|compare|review|review-context> --input arquivo.json [--output arquivo.json]")
	}
	input, output, err := parseDecisionTraceArgs(args[1:])
	if err != nil {
		return err
	}
	if err := validateDecisionTraceOutput(output, input); err != nil {
		return err
	}
	var result any
	if args[0] == "build" {
		var request domain.DecisionTraceBuildInput
		if err := readStrictIntakeJSON(input, &request); err != nil {
			return err
		}
		result, err = usecase.BuildDecisionTrace(request)
	} else if args[0] == "impact" {
		var request domain.DecisionTraceImpactInput
		if err := readStrictIntakeJSON(input, &request); err != nil {
			return err
		}
		result, err = usecase.PlanDecisionTraceImpact(request)
	} else if args[0] == "explain" {
		var request domain.DecisionTraceExplainInput
		if err := readStrictIntakeJSON(input, &request); err != nil {
			return err
		}
		result, err = usecase.ExplainDecisionTrace(request)
	} else if args[0] == "compare" {
		var request domain.DecisionTraceCompareInput
		if err := readStrictIntakeJSON(input, &request); err != nil {
			return err
		}
		result, err = usecase.CompareDecisionTraces(request)
	} else if args[0] == "review" {
		var request domain.DecisionTraceReviewInput
		if err := readStrictIntakeJSON(input, &request); err != nil {
			return err
		}
		result, err = usecase.TrackDecisionTraceReview(request)
	} else {
		var request domain.DecisionTraceReviewContextInput
		if err := readStrictIntakeJSON(input, &request); err != nil {
			return err
		}
		result, err = usecase.BuildDecisionTraceReviewContext(request)
	}
	if err != nil {
		return fmt.Errorf("rastreabilidade rejeitada: %w", err)
	}
	return writeDecisionTraceJSON(result, output, stdout, input)
}

func parseDecisionTraceArgs(args []string) (string, string, error) {
	input, output := "", ""
	seen := map[string]bool{}
	for index := 0; index < len(args); index++ {
		name, value, equal := strings.Cut(args[index], "=")
		if name != "--input" && name != "--output" {
			return "", "", fmt.Errorf("parâmetro idea trace inválido %q", args[index])
		}
		if seen[name] {
			return "", "", fmt.Errorf("idea trace não aceita repetição de %s", name)
		}
		seen[name] = true
		if !equal {
			index++
			if index == len(args) || strings.HasPrefix(args[index], "--") {
				return "", "", fmt.Errorf("%s exige um valor", name)
			}
			value = args[index]
		}
		if strings.TrimSpace(value) == "" {
			return "", "", fmt.Errorf("%s exige um valor não vazio", name)
		}
		if name == "--input" {
			input = value
		} else {
			output = value
		}
	}
	if input == "" {
		return "", "", fmt.Errorf("idea trace exige --input")
	}
	return input, output, nil
}

// Output paths may replace an explicitly selected regular result, never a source.
// Symlink components are conservatively refused, including missing output leaves.
func validateDecisionTraceOutput(output string, inputs ...string) error {
	return commandHost().ValidateTraceOutput(output, inputs...)
}
func resolveDecisionTraceProtectedPath(path string) (string, error) {
	return commandHost().ResolveProtectedPath(path)
}
func writeDecisionTraceJSON(value any, output string, stdout io.Writer, inputs ...string) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("serializar trace: %w", err)
	}
	encoded = append(encoded, '\n')
	if len(encoded) > maxIdeaInputBytes {
		return fmt.Errorf("resultado trace excede orçamento de serialização")
	}
	if err := rejectSerializedSecrets(encoded); err != nil {
		return err
	}
	if output == "" {
		if stdout == nil {
			return fmt.Errorf("stdout trace indisponível")
		}
		written, err := stdout.Write(encoded)
		if err == nil && written != len(encoded) {
			return io.ErrShortWrite
		}
		return err
	}
	return commandHost().WriteTraceFile(encoded, output, inputs...)
}
