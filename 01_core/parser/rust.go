/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/rust.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"
	"strings"

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

		// Dependências de Workspace
		if currentSection == "[workspace.dependencies]" {
			dep := parseCargoDepLine(trimmed, "workspace")
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

var (
	rustFnRegex     = regexp.MustCompile(`(?m)^\s*(?:pub(?:\([\w:]+\))?\s+)?fn\s+([a-zA-Z0-9_]+)`)
	rustStructRegex = regexp.MustCompile(`(?m)^\s*(?:pub(?:\([\w:]+\))?\s+)?struct\s+([a-zA-Z0-9_]+)`)
	rustEnumRegex   = regexp.MustCompile(`(?m)^\s*(?:pub(?:\([\w:]+\))?\s+)?enum\s+([a-zA-Z0-9_]+)`)
	rustTraitRegex  = regexp.MustCompile(`(?m)^\s*(?:pub(?:\([\w:]+\))?\s+)?trait\s+([a-zA-Z0-9_]+)`)
	rustImplRegex   = regexp.MustCompile(`(?m)^\s*impl(?:<[^>]+>)?\s+([a-zA-Z0-9_]+)`)
)

// ExtractRustSymbols extrai os símbolos de um arquivo fonte Rust (.rs) em memória.
func ExtractRustSymbols(path string, content []byte) []domain.ASTSymbol {
	lines := strings.Split(string(content), "\n")
	var symbols []domain.ASTSymbol

	for lineIdx, line := range lines {
		lineNo := lineIdx + 1
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		if matches := rustStructRegex.FindStringSubmatch(line); len(matches) >= 2 {
			symbols = append(symbols, domain.ASTSymbol{Name: matches[1], Kind: "struct", File: path, Line: lineNo})
		} else if matches := rustFnRegex.FindStringSubmatch(line); len(matches) >= 2 {
			symbols = append(symbols, domain.ASTSymbol{Name: matches[1], Kind: "function", File: path, Line: lineNo})
		} else if matches := rustEnumRegex.FindStringSubmatch(line); len(matches) >= 2 {
			symbols = append(symbols, domain.ASTSymbol{Name: matches[1], Kind: "enum", File: path, Line: lineNo})
		} else if matches := rustTraitRegex.FindStringSubmatch(line); len(matches) >= 2 {
			symbols = append(symbols, domain.ASTSymbol{Name: matches[1], Kind: "interface", File: path, Line: lineNo})
		} else if matches := rustImplRegex.FindStringSubmatch(line); len(matches) >= 2 {
			symbols = append(symbols, domain.ASTSymbol{Name: matches[1], Kind: "impl", File: path, Line: lineNo})
		}
	}

	return symbols
}
