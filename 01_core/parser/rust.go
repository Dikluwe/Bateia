/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/rust.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"regexp"
	"strings"
	"time"

	"bateia/01_core/domain"
)

var (
	cargoPackageNameRegex = regexp.MustCompile(`(?m)^\s*name\s*=\s*["']([^"']+)["']`)
	cargoSimpleDepRegex   = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-\.]+)\s*=\s*["']([^"']+)["']`)
	cargoInlineTableRegex = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-\.]+)\s*=\s*\{.*version\s*=\s*["']([^"']+)["'].*\}`)
)

// ParseCargoToml analisa o conteúdo de um manifesto Cargo.toml de Rust.
func ParseCargoToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "rust",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	lines := strings.Split(content, "\n")
	currentSection := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			currentSection = strings.ToLower(trimmed)
			continue
		}

		// Extrair nome do pacote de [package]
		if manifest.ModuleName == "" && currentSection == "[package]" {
			if matches := cargoPackageNameRegex.FindStringSubmatch(trimmed); len(matches) >= 2 {
				manifest.ModuleName = matches[1]
			}
		}

		// Dependências de produção
		if currentSection == "[dependencies]" || strings.HasPrefix(currentSection, "[target.") && strings.HasSuffix(currentSection, ".dependencies]") {
			dep := parseCargoDepLine(trimmed, "direct")
			if dep != nil {
				manifest.Dependencies = append(manifest.Dependencies, *dep)
			}
		}

		// Dependências de desenvolvimento
		if currentSection == "[dev-dependencies]" || strings.Contains(currentSection, "dev-dependencies") {
			dep := parseCargoDepLine(trimmed, "dev")
			if dep != nil {
				manifest.Dependencies = append(manifest.Dependencies, *dep)
			}
		}

		// Dependências de build
		if currentSection == "[build-dependencies]" {
			dep := parseCargoDepLine(trimmed, "build")
			if dep != nil {
				manifest.Dependencies = append(manifest.Dependencies, *dep)
			}
		}
	}

	if manifest.ModuleName == "" {
		manifest.ModuleName = ref.Repo
	}

	return manifest, nil
}

func parseCargoDepLine(line, depType string) *domain.Dependency {
	// Tentar inline table primeiro ex: `serde = { version = "1.0", features = [...] }`
	if matches := cargoInlineTableRegex.FindStringSubmatch(line); len(matches) >= 3 {
		return &domain.Dependency{
			Name:     strings.ToLower(matches[1]),
			Version:  matches[2],
			Indirect: false,
			Type:     depType,
		}
	}

	// Tentar versão string simples ex: `tokio = "1.35.0"`
	if matches := cargoSimpleDepRegex.FindStringSubmatch(line); len(matches) >= 3 {
		return &domain.Dependency{
			Name:     strings.ToLower(matches[1]),
			Version:  matches[2],
			Indirect: false,
			Type:     depType,
		}
	}

	return nil
}
