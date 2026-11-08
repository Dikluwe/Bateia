/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/exports.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"bateia/01_core/domain"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func openAlgorithmHistoryInput(inputPath, outputPath string) (domain.CommandFile, domain.FileInfo, error) {
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
	fail := func(err error) (domain.CommandFile, domain.FileInfo, error) { inputFile.Close(); return nil, nil, err }
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

func writeAlgorithmHistoryAtomic(outputPath string, data []byte, inputInfo domain.FileInfo) error {
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

func openIdeaRadarExportInput(inputPath, outputPath string) (domain.CommandFile, domain.FileInfo, error) {
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
	fail := func(err error) (domain.CommandFile, domain.FileInfo, error) {
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

func rejectIdeaRadarExportAlias(inputPath, outputPath string, inputInfo domain.FileInfo) error {
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

func writeIdeaRadarExportAtomic(outputPath string, data []byte, inputInfo domain.FileInfo) error {
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

func (*Host) OpenAlgorithmHistoryInput(a, b string) (domain.CommandFile, domain.FileInfo, error) {
	return openAlgorithmHistoryInput(a, b)
}
func (*Host) WriteAlgorithmHistoryAtomic(p string, data []byte, info domain.FileInfo) error {
	return writeAlgorithmHistoryAtomic(p, data, info)
}
func (*Host) OpenIdeaRadarExportInput(a, b string) (domain.CommandFile, domain.FileInfo, error) {
	return openIdeaRadarExportInput(a, b)
}
func (*Host) RejectIdeaRadarExportAlias(a, b string, info domain.FileInfo) error {
	return rejectIdeaRadarExportAlias(a, b, info)
}
func (*Host) WriteIdeaRadarExportAtomic(p string, data []byte, info domain.FileInfo) error {
	return writeIdeaRadarExportAtomic(p, data, info)
}
