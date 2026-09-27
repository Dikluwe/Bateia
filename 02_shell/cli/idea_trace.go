/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_trace.md
 * @layer L2
 * @updated 2026-09-27
 */
package cli

import (
	"crypto/rand"
	"encoding/hex"
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
