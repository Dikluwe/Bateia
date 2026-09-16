/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/algorithm_history.md
 * @layer L2
 * @updated 2026-09-12
 */

package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func openAlgorithmHistoryInput(inputPath, outputPath string) (*os.File, os.FileInfo, error) {
	if outputPath != "" {
		inputAbsolute, err := filepath.Abs(inputPath)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve algorithm history input path: %w", err)
		}
		outputAbsolute, err := filepath.Abs(outputPath)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve algorithm history output path: %w", err)
		}
		if filepath.Clean(inputAbsolute) == filepath.Clean(outputAbsolute) {
			return nil, nil, fmt.Errorf("algorithm history output must not identify its input")
		}
	}
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read algorithm history input %s: %w", inputPath, err)
	}
	fail := func(err error) (*os.File, os.FileInfo, error) { inputFile.Close(); return nil, nil, err }
	inputInfo, err := inputFile.Stat()
	if err != nil {
		return fail(fmt.Errorf("inspect algorithm history input %s: %w", inputPath, err))
	}
	if outputPath == "" {
		return inputFile, inputInfo, nil
	}
	if outputInfo, err := os.Stat(outputPath); err == nil && os.SameFile(inputInfo, outputInfo) {
		return fail(fmt.Errorf("algorithm history output must not identify its input"))
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fail(fmt.Errorf("inspect algorithm history output %s: %w", outputPath, err))
	}
	resolvedInput, err := filepath.EvalSymlinks(inputPath)
	if err != nil {
		return fail(fmt.Errorf("resolve algorithm history input %s: %w", inputPath, err))
	}
	resolvedParent, err := filepath.EvalSymlinks(filepath.Dir(outputPath))
	if err != nil {
		return fail(fmt.Errorf("resolve algorithm history output directory %s: %w", filepath.Dir(outputPath), err))
	}
	resolvedOutput := filepath.Join(resolvedParent, filepath.Base(outputPath))
	resolvedInput, err = filepath.Abs(resolvedInput)
	if err != nil {
		return fail(fmt.Errorf("resolve algorithm history input %s: %w", inputPath, err))
	}
	resolvedOutput, err = filepath.Abs(resolvedOutput)
	if err != nil {
		return fail(fmt.Errorf("resolve algorithm history output %s: %w", outputPath, err))
	}
	if filepath.Clean(resolvedInput) == filepath.Clean(resolvedOutput) {
		return fail(fmt.Errorf("algorithm history output must not identify its input"))
	}
	return inputFile, inputInfo, nil
}

func writeAlgorithmHistoryAtomic(outputPath string, data []byte, inputInfo os.FileInfo) error {
	directory := filepath.Dir(outputPath)
	temporary, err := os.CreateTemp(directory, ".bateia-algorithm-history-*")
	if err != nil {
		return fmt.Errorf("create temporary algorithm history report in %s: %w", directory, err)
	}
	temporaryPath := temporary.Name()
	defer func() { _ = temporary.Close(); _ = os.Remove(temporaryPath) }()
	if err := temporary.Chmod(0644); err != nil {
		return fmt.Errorf("set temporary algorithm history permissions: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		return fmt.Errorf("write temporary algorithm history report: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync temporary algorithm history report: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary algorithm history report: %w", err)
	}
	if outputInfo, err := os.Stat(outputPath); err == nil && os.SameFile(inputInfo, outputInfo) {
		return fmt.Errorf("algorithm history output must not identify its input")
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect algorithm history output %s: %w", outputPath, err)
	}
	if err := os.Rename(temporaryPath, outputPath); err != nil {
		return fmt.Errorf("atomically replace algorithm history output %s: %w", outputPath, err)
	}
	return nil
}
