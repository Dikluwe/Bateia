/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_radar_export.md
 * @layer L2
 * @updated 2026-10-28
 */

package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

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

func openIdeaRadarExportInput(a, b string) (domain.CommandFile, domain.FileInfo, error) {
	return commandHost().OpenIdeaRadarExportInput(a, b)
}
func rejectIdeaRadarExportAlias(a, b string, info domain.FileInfo) error {
	return commandHost().RejectIdeaRadarExportAlias(a, b, info)
}
func writeIdeaRadarExportAtomic(p string, data []byte, info domain.FileInfo) error {
	return commandHost().WriteIdeaRadarExportAtomic(p, data, info)
}
