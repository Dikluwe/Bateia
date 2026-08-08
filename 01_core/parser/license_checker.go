/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/license_checker.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "strings"

type LicenseCategory string

const (
	Permissive LicenseCategory = "PERMISSIVE"
	Copyleft   LicenseCategory = "COPYLEFT"
	Unknown    LicenseCategory = "UNKNOWN"
)

// CategorizeLicense identifica a licença de código a partir do conteúdo do arquivo LICENSE ou cabeçalho.
func CategorizeLicense(licenseText string) LicenseCategory {
	if licenseText == "" {
		return Unknown
	}

	lower := strings.ToLower(licenseText)

	// Licenças Copyleft (Infectivas)
	if strings.Contains(lower, "general public license") || strings.Contains(lower, "gpl") || strings.Contains(lower, "agpl") || strings.Contains(lower, "mozilla public license") || strings.Contains(lower, "mpl") {
		return Copyleft
	}

	// Licenças Permissivas
	if strings.Contains(lower, "mit license") || strings.Contains(lower, "apache license") || strings.Contains(lower, "bsd") || strings.Contains(lower, "isc license") || strings.Contains(lower, "unlicense") {
		return Permissive
	}

	return Unknown
}
