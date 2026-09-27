/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_radar_export.md
 * @layer L2
 * @updated 2026-09-12
 */

package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"bateia/01_core/domain"
	"bateia/02_shell/report"
)

// RunIdeaRadarExportCommand handles the strictly offline idea radar export.
func RunIdeaRadarExportCommand(args []string, stdout io.Writer) error {
	inputPath, outputPath, err := parseIdeaRadarExportArgs(args)
	if err != nil {
		return err
	}
	inputFile, inputInfo, err := openIdeaRadarExportInput(inputPath, outputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	data, err := readIdeaDescriptor(inputFile)
	if err != nil {
		return fmt.Errorf("read idea radar input %s: %w", inputPath, err)
	}
	if err := validateIdeaJSONStructure(data); err != nil {
		return fmt.Errorf("invalid idea radar JSON in %s: %w", inputPath, err)
	}

	var input domain.IdeaRadarInput
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return fmt.Errorf("invalid idea radar JSON in %s: %w", inputPath, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("invalid idea radar JSON in %s: trailing JSON value", inputPath)
		}
		return fmt.Errorf("invalid idea radar JSON in %s: trailing data: %w", inputPath, err)
	}

	encoded, err := report.ExportIdeaRadarArchify(input)
	if err != nil {
		return err
	}
	if outputPath == "" {
		if stdout == nil {
			return fmt.Errorf("idea radar export stdout is unavailable")
		}
		if _, err := stdout.Write(encoded); err != nil {
			return fmt.Errorf("write idea radar export to stdout: %w", err)
		}
		return nil
	}
	return writeIdeaRadarExportAtomic(outputPath, encoded, inputInfo)
}

func parseIdeaRadarExportArgs(args []string) (string, string, error) {
	if len(args) == 0 || args[0] != "export" {
		return "", "", fmt.Errorf("usage: bateia idea radar export --input FILE [--output FILE]")
	}
	inputPath, outputPath := "", ""
	seenInput, seenOutput := false, false
	for index := 1; index < len(args); index++ {
		argument := args[index]
		name, value, hasEquals := strings.Cut(argument, "=")
		if name != "--input" && name != "--output" {
			return "", "", fmt.Errorf("invalid idea radar export argument %q", argument)
		}
		if !hasEquals {
			index++
			if index >= len(args) || strings.HasPrefix(args[index], "--") {
				return "", "", fmt.Errorf("idea radar export flag %s requires a value", name)
			}
			value = args[index]
		}
		if strings.TrimSpace(value) == "" {
			return "", "", fmt.Errorf("idea radar export flag %s requires a non-empty value", name)
		}
		switch name {
		case "--input":
			if seenInput {
				return "", "", fmt.Errorf("idea radar export accepts --input exactly once")
			}
			seenInput, inputPath = true, value
		case "--output":
			if seenOutput {
				return "", "", fmt.Errorf("idea radar export accepts --output at most once")
			}
			seenOutput, outputPath = true, value
		}
	}
	if !seenInput {
		return "", "", fmt.Errorf("idea radar export requires --input FILE")
	}
	return inputPath, outputPath, nil
}

func openIdeaRadarExportInput(inputPath, outputPath string) (*os.File, os.FileInfo, error) {
	if outputPath != "" {
		inputAbsolute, err := filepath.Abs(inputPath)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve idea radar input path: %w", err)
		}
		outputAbsolute, err := filepath.Abs(outputPath)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve idea radar output path: %w", err)
		}
		if filepath.Clean(inputAbsolute) == filepath.Clean(outputAbsolute) {
			return nil, nil, fmt.Errorf("idea radar output must not identify its input")
		}
	}

	inputFile, err := os.OpenFile(inputPath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("read idea radar input %s: %w", inputPath, err)
	}
	fail := func(err error) (*os.File, os.FileInfo, error) {
		inputFile.Close()
		return nil, nil, err
	}
	inputInfo, err := inputFile.Stat()
	if err != nil {
		return fail(fmt.Errorf("inspect idea radar input %s: %w", inputPath, err))
	}
	if outputPath == "" {
		return inputFile, inputInfo, nil
	}
	if err := rejectIdeaRadarExportAlias(inputPath, outputPath, inputInfo); err != nil {
		return fail(err)
	}
	return inputFile, inputInfo, nil
}

func rejectIdeaRadarExportAlias(inputPath, outputPath string, inputInfo os.FileInfo) error {
	outputInfo, err := os.Stat(outputPath)
	if err == nil && os.SameFile(inputInfo, outputInfo) {
		return fmt.Errorf("idea radar output must not identify its input")
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect idea radar output %s: %w", outputPath, err)
	}

	resolvedInput, err := filepath.EvalSymlinks(inputPath)
	if err != nil {
		return fmt.Errorf("resolve idea radar input %s: %w", inputPath, err)
	}
	resolvedInput, err = filepath.Abs(resolvedInput)
	if err != nil {
		return fmt.Errorf("resolve idea radar input %s: %w", inputPath, err)
	}
	resolvedOutputParent, err := filepath.EvalSymlinks(filepath.Dir(outputPath))
	if err != nil {
		return fmt.Errorf("resolve idea radar output directory %s: %w", filepath.Dir(outputPath), err)
	}
	resolvedOutput := filepath.Join(resolvedOutputParent, filepath.Base(outputPath))
	resolvedOutput, err = filepath.Abs(resolvedOutput)
	if err != nil {
		return fmt.Errorf("resolve idea radar output %s: %w", outputPath, err)
	}
	if filepath.Clean(resolvedInput) == filepath.Clean(resolvedOutput) {
		return fmt.Errorf("idea radar output must not identify its input")
	}
	return nil
}

func writeIdeaRadarExportAtomic(outputPath string, data []byte, inputInfo os.FileInfo) error {
	directory := filepath.Dir(outputPath)
	temporary, err := os.CreateTemp(directory, ".bateia-idea-radar-export-*")
	if err != nil {
		return fmt.Errorf("create temporary idea radar export in %s: %w", directory, err)
	}
	temporaryPath := temporary.Name()
	defer func() {
		_ = temporary.Close()
		_ = os.Remove(temporaryPath)
	}()
	if err := temporary.Chmod(0644); err != nil {
		return fmt.Errorf("set temporary idea radar export permissions: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		return fmt.Errorf("write temporary idea radar export: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync temporary idea radar export: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary idea radar export: %w", err)
	}
	if outputInfo, err := os.Stat(outputPath); err == nil && os.SameFile(inputInfo, outputInfo) {
		return fmt.Errorf("idea radar output must not identify its input")
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect idea radar output %s: %w", outputPath, err)
	}
	if err := os.Rename(temporaryPath, outputPath); err != nil {
		return fmt.Errorf("atomically replace idea radar output %s: %w", outputPath, err)
	}
	return nil
}
