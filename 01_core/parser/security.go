/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/security.md
 * @layer L1
 * @updated 2026-08-23
 */

package parser

import (
	"regexp"
	"strconv"
	"strings"

	"bateia/01_core/domain"
)

const advisoryDatabaseUpdatedAt = "2026-08-23"

type knownAdvisory struct {
	packages    []string
	minVersion  string
	maxVersion  string
	includeMax  bool
	severity    string
	advisoryID  string
	description string
	source      string
}

var knownAdvisories = []knownAdvisory{
	{[]string{"log4j-core"}, "2.0.0", "2.15.0", false, "CRITICAL", "GHSA-j2ge-4423-1574", "Possível Log4Shell em versão afetada do Log4j Core.", "https://github.com/advisories/GHSA-j2ge-4423-1574"},
	{[]string{"event-stream"}, "3.3.6", "3.3.6", true, "HIGH", "GHSA-8g9b-32a2-3f11", "Versão associada ao comprometimento da dependência event-stream.", "https://github.com/advisories/GHSA-8g9b-32a2-3f11"},
	{[]string{"github.com/dgrijalva/jwt-go", "jwt-go"}, "3.2.0", "3.2.0", true, "HIGH", "CVE-2020-26160", "Possível validação incorreta de tokens no jwt-go descontinuado.", "https://nvd.nist.gov/vuln/detail/CVE-2020-26160"},
	{[]string{"golang.org/x/net"}, "0.16.0", "0.16.0", true, "CRITICAL", "CVE-2023-39325", "Versão localmente sinalizada para revisão de vulnerabilidades HTTP/2.", "https://nvd.nist.gov/vuln/detail/CVE-2023-39325"},
	{[]string{"spring-core"}, "5.3.0", "5.3.18", false, "CRITICAL", "CVE-2022-22965", "Possível Spring4Shell em faixa afetada do Spring Framework.", "https://nvd.nist.gov/vuln/detail/CVE-2022-22965"},
	{[]string{"requests"}, "2.3.0", "2.31.0", false, "HIGH", "CVE-2023-32681", "Possível vazamento de Proxy-Authorization em redirecionamentos.", "https://nvd.nist.gov/vuln/detail/CVE-2023-32681"},
	{[]string{"pyyaml"}, "5.1.0", "5.4.0", false, "HIGH", "CVE-2020-14343", "Possível desserialização insegura em faixa afetada do PyYAML.", "https://nvd.nist.gov/vuln/detail/CVE-2020-14343"},
}

var versionNumberPattern = regexp.MustCompile(`\d+`)

func normalizedVersion(version string) [3]int {
	var result [3]int
	numbers := versionNumberPattern.FindAllString(version, 3)
	for i, number := range numbers {
		result[i], _ = strconv.Atoi(number)
	}
	return result
}

func compareVersion(a, b [3]int) int {
	for i := range a {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

func packageMatches(name string, aliases []string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, alias := range aliases {
		alias = strings.ToLower(alias)
		if name == alias || strings.HasSuffix(name, "/"+alias) || strings.HasSuffix(name, ":"+alias) {
			return true
		}
	}
	return false
}

func versionMatches(version string, advisory knownAdvisory) bool {
	if len(versionNumberPattern.FindAllString(version, 1)) == 0 {
		return false
	}
	actual := normalizedVersion(version)
	if compareVersion(actual, normalizedVersion(advisory.minVersion)) < 0 {
		return false
	}
	maxComparison := compareVersion(actual, normalizedVersion(advisory.maxVersion))
	return maxComparison < 0 || (advisory.includeMax && maxComparison == 0)
}

// EvaluateSecurityAdvisories gera correspondências heurísticas de uma base local.
// Os resultados não substituem um scanner que resolva versões e dependências transitivas.
func EvaluateSecurityAdvisories(manifest *domain.Manifest) []domain.SecurityAdvisory {
	var results []domain.SecurityAdvisory
	if manifest == nil {
		return results
	}
	for _, dep := range manifest.Dependencies {
		for _, advisory := range knownAdvisories {
			if packageMatches(dep.Name, advisory.packages) && versionMatches(dep.Version, advisory) {
				results = append(results, domain.SecurityAdvisory{
					Package: dep.Name, VulnerableVersion: dep.Version, Severity: advisory.severity,
					AdvisoryID: advisory.advisoryID, Message: advisory.description,
					Confidence: "MEDIUM", MatchType: "local_package_and_version_range",
					Source: advisory.source, DatabaseUpdatedAt: advisoryDatabaseUpdatedAt,
				})
			}
		}
	}
	return results
}
