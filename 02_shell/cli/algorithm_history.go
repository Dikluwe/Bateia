/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/algorithm_history.md
 * @layer L2
 * @updated 2026-10-28
 */

package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

// RunAlgorithmHistoryCommand executes the strictly offline historical analysis.
func RunAlgorithmHistoryCommand(args []string, stdout io.Writer) error {
	inputPath, outputPath, err := parseAlgorithmHistoryArgs(args)
	if err != nil {
		return err
	}
	inputFile, inputInfo, err := openAlgorithmHistoryInput(inputPath, outputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	var input domain.AlgorithmHistoryInput
	decoder := json.NewDecoder(inputFile)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return fmt.Errorf("invalid algorithm history JSON in %s: %w", inputPath, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("invalid algorithm history JSON in %s: trailing JSON value", inputPath)
		}
		return fmt.Errorf("invalid algorithm history JSON in %s: trailing data: %w", inputPath, err)
	}

	report := usecase.AnalyzeAlgorithmHistory(input)
	encoded, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize algorithm history report: %w", err)
	}
	encoded = append(encoded, '\n')
	if len(encoded) > 4*1024*1024 {
		return fmt.Errorf("algorithm history report exceeds serialized budget")
	}
	if outputPath == "" {
		if stdout == nil {
			return fmt.Errorf("algorithm history stdout is unavailable")
		}
		if _, err := stdout.Write(encoded); err != nil {
			return fmt.Errorf("write algorithm history report to stdout: %w", err)
		}
		return nil
	}
	return writeAlgorithmHistoryAtomic(outputPath, encoded, inputInfo)
}

func parseAlgorithmHistoryArgs(args []string) (string, string, error) {
	if len(args) < 2 || args[0] != "history" || args[1] != "analyze" {
		return "", "", fmt.Errorf("usage: bateia algorithm history analyze --input FILE [--output FILE]")
	}
	inputPath, outputPath := "", ""
	seenInput, seenOutput := false, false
	for index := 2; index < len(args); index++ {
		name, value, hasEquals := strings.Cut(args[index], "=")
		if name != "--input" && name != "--output" {
			return "", "", fmt.Errorf("invalid algorithm history argument %q", args[index])
		}
		if !hasEquals {
			index++
			if index >= len(args) || strings.HasPrefix(args[index], "--") {
				return "", "", fmt.Errorf("algorithm history flag %s requires a value", name)
			}
			value = args[index]
		}
		if strings.TrimSpace(value) == "" {
			return "", "", fmt.Errorf("algorithm history flag %s requires a non-empty value", name)
		}
		switch name {
		case "--input":
			if seenInput {
				return "", "", fmt.Errorf("algorithm history accepts --input exactly once")
			}
			seenInput, inputPath = true, value
		case "--output":
			if seenOutput {
				return "", "", fmt.Errorf("algorithm history accepts --output at most once")
			}
			seenOutput, outputPath = true, value
		}
	}
	if !seenInput {
		return "", "", fmt.Errorf("algorithm history requires --input FILE")
	}
	return inputPath, outputPath, nil
}

func openAlgorithmHistoryInput(a, b string) (domain.CommandFile, domain.FileInfo, error) {
	return commandHost().OpenAlgorithmHistoryInput(a, b)
}
func writeAlgorithmHistoryAtomic(p string, data []byte, info domain.FileInfo) error {
	return commandHost().WriteAlgorithmHistoryAtomic(p, data, info)
}
