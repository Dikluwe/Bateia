/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/security.md
 * @layer L1
 * @updated 2026-08-03
 */

package parser

import (
	"strings"

	"bateia/01_core/domain"
)

type knownAdvisory struct {
	pkgSub      string
	versionSub  string
	severity    string
	advisoryID  string
	description string
}

var knownAdvisories = []knownAdvisory{
	{
		pkgSub:      "log4j",
		versionSub:  "2.",
		severity:    "CRITICAL",
		advisoryID:  "GHSA-j2ge-4423-1574 / CVE-2021-44228",
		description: "RCE em Log4j2 (Log4Shell) via interpolação de JNDI.",
	},
	{
		pkgSub:      "event-stream",
		versionSub:  "3.3.6",
		severity:    "HIGH",
		advisoryID:  "GHSA-8g9b-32a2-3f11",
		description: "Injeção maliciosa em dependência transitiva do npm.",
	},
	{
		pkgSub:      "jwt-go",
		versionSub:  "v3.2.0",
		severity:    "HIGH",
		advisoryID:  "CVE-2020-26160",
		description: "Vulnerabilidade em validação de token JWT descontinuado.",
	},
}

// EvaluateSecurityAdvisories inspeciona as dependências do manifesto em busca de CVEs conhecidos.
func EvaluateSecurityAdvisories(manifest *domain.Manifest) []domain.SecurityAdvisory {
	var results []domain.SecurityAdvisory

	if manifest == nil {
		return results
	}

	for _, dep := range manifest.Dependencies {
		lowerDep := strings.ToLower(dep.Name)
		lowerVer := strings.ToLower(dep.Version)

		for _, adv := range knownAdvisories {
			if strings.Contains(lowerDep, adv.pkgSub) {
				if adv.versionSub == "" || strings.Contains(lowerVer, adv.versionSub) {
					results = append(results, domain.SecurityAdvisory{
						Package:           dep.Name,
						VulnerableVersion: dep.Version,
						Severity:          adv.severity,
						AdvisoryID:        adv.advisoryID,
						Message:           adv.description,
					})
				}
			}
		}
	}

	return results
}
