/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/security.md
 * @layer L1
 * @updated 2026-08-04
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
	{
		pkgSub:      "golang.org/x/net",
		versionSub:  "v0.16.0",
		severity:    "CRITICAL",
		advisoryID:  "CVE-2023-39325 / CVE-2023-44487",
		description: "Vulnerabilidade HTTP/2 Rapid Reset causando Negação de Serviço (DoS).",
	},
	{
		pkgSub:      "spring-core",
		versionSub:  "5.3.",
		severity:    "CRITICAL",
		advisoryID:  "CVE-2022-22965",
		description: "Spring4Shell RCE no DataBinder do Spring Framework.",
	},
	{
		pkgSub:      "requests",
		versionSub:  "2.30",
		severity:    "HIGH",
		advisoryID:  "CVE-2023-32681",
		description: "Vazamento não intencional de cabeçalho Proxy-Authorization em redirecionamentos HTTP.",
	},
	{
		pkgSub:      "pyyaml",
		versionSub:  "5.",
		severity:    "HIGH",
		advisoryID:  "CVE-2020-14343",
		description: "Execução remota de código via desserialização não confiável no PyYAML.",
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
