/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/trace_files.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validateDecisionTraceOutput(output string, inputs ...string) error {
	if output == "" {
		return nil
	}
	absolute, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("resolver saída trace: %w", err)
	}
	parent := filepath.Dir(absolute)
	for path := absolute; ; path = filepath.Dir(path) {
		info, err := os.Lstat(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspecionar saída trace %s: %w", path, err)
		}
		if err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("saída trace não aceita links simbólicos: %s", path)
			}
			if path == absolute && !info.Mode().IsRegular() {
				return fmt.Errorf("saída trace existente deve ser arquivo regular")
			}
			if path != absolute && !info.IsDir() {
				return fmt.Errorf("diretório de saída trace inválido: %s", path)
			}
		}
		if path == filepath.Dir(path) {
			break
		}
	}
	if parent == absolute {
		return fmt.Errorf("saída trace inválida")
	}
	for _, input := range inputs {
		if input == "" {
			continue
		}
		if sameTechnicalOutputPath(input, output) {
			return fmt.Errorf("saída trace não pode identificar uma entrada")
		}
		resolved, err := resolveDecisionTraceProtectedPath(input)
		if err != nil {
			return fmt.Errorf("resolver entrada trace: %w", err)
		}
		if resolved == absolute {
			return fmt.Errorf("saída trace não pode identificar uma entrada")
		}
	}
	return nil
}

// Protected paths can also be future outputs with missing leaves or parents.
func resolveDecisionTraceProtectedPath(path string) (string, error) {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return filepath.Abs(resolved)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("destino protegido não pode ser link simbólico pendente")
	}
	// Missing paths with parent traversal are ambiguous under future directory
	// creation, especially when a symlink precedes '..'. Fail closed instead.
	for _, component := range strings.Split(filepath.ToSlash(path), "/") {
		if component == ".." {
			return "", fmt.Errorf("destino protegido ausente não aceita travessia de pai")
		}
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	var missing []string
	for prefix := absolute; ; prefix = filepath.Dir(prefix) {
		resolved, err := filepath.EvalSymlinks(prefix)
		if err == nil {
			for index := len(missing) - 1; index >= 0; index-- {
				resolved = filepath.Join(resolved, missing[index])
			}
			return filepath.Abs(resolved)
		}
		if !errors.Is(err, os.ErrNotExist) || prefix == filepath.Dir(prefix) {
			return "", err
		}
		info, statErr := os.Lstat(prefix)
		if statErr == nil && info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("destino protegido possui link simbólico pendente")
		}
		if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return "", statErr
		}
		missing = append(missing, filepath.Base(prefix))
	}
}

// Pin the output parent while writing and rename only a complete, synced result.
func (*Host) WriteTraceFile(encoded []byte, output string, inputs ...string) error {
	if len(encoded) > 16<<20 {
		return fmt.Errorf("resultado trace excede orçamento de serialização")
	}
	if err := validateDecisionTraceOutput(output, inputs...); err != nil {
		return err
	}
	parentInfo, err := os.Lstat(filepath.Dir(output))
	if err != nil {
		return err
	}
	if !parentInfo.IsDir() || parentInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("diretório de saída trace deve ser diretório real")
	}
	root, err := os.OpenRoot(filepath.Dir(output))
	if err != nil {
		return fmt.Errorf("abrir diretório trace: %w", err)
	}
	defer root.Close()
	openedParent, err := root.Stat(".")
	if err != nil || !os.SameFile(parentInfo, openedParent) {
		return fmt.Errorf("diretório de saída trace mudou durante abertura")
	}
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return err
	}
	temporary := ".bateia-trace-" + hex.EncodeToString(nonce[:])
	file, err := root.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("criar resultado trace temporário: %w", err)
	}
	defer func() { _ = file.Close(); _ = root.Remove(temporary) }()
	if _, err := file.Write(encoded); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := validateDecisionTraceOutput(output, inputs...); err != nil {
		return err
	}
	leaf := filepath.Base(output)
	info, err := root.Lstat(leaf)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("saída trace existente deve ser arquivo regular")
		}
		for _, input := range inputs {
			if input == "" {
				continue
			}
			inputInfo, err := os.Stat(input)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				return err
			}
			if os.SameFile(info, inputInfo) {
				return fmt.Errorf("saída trace não pode identificar uma entrada")
			}
		}
	}
	if err := root.Rename(temporary, leaf); err != nil {
		return fmt.Errorf("gravar resultado trace: %w", err)
	}
	return nil
}

func (*Host) ValidateTraceOutput(output string, inputs ...string) error {
	return validateDecisionTraceOutput(output, inputs...)
}
func (*Host) ResolveProtectedPath(path string) (string, error) {
	return resolveDecisionTraceProtectedPath(path)
}
func sameTechnicalOutputPath(a, b string) bool {
	aa, errA := filepath.Abs(a)
	bb, errB := filepath.Abs(b)
	if errA == nil && errB == nil && aa == bb {
		return true
	}
	sa, errA := os.Stat(a)
	sb, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(sa, sb)
}

func (*Host) SameOutputPath(a, b string) bool { return sameTechnicalOutputPath(a, b) }
