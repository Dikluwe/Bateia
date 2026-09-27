/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/security_scan.md
 * @layer L1
 * @updated 2026-09-26
 */
package parser

import "bateia/01_core/domain"

func EvaluateSecurityScan(manifest *domain.Manifest) (domain.SecurityScanReport, []domain.SecurityAdvisory) {
	report := domain.SecurityScanReport{
		Status: domain.SecurityNotScanned, Scanner: "LOCAL_HEURISTIC_V1", Scope: "declared_direct_manifest_dependencies", DatabaseUpdatedAt: advisoryDatabaseUpdatedAt,
		Limitations: []string{"A base local é parcial e heurística, não é SCA completa.", "Ausência de correspondência não prova segurança.", "Não resolve lockfiles, dependências transitivas, alcance no binário ou explorabilidade."},
	}
	if manifest == nil {
		return report, nil
	}
	advisories := EvaluateSecurityAdvisories(manifest)
	report.Status = domain.SecurityScannedNoMatch
	if len(advisories) > 0 {
		report.Status = domain.SecurityVulnerabilityFound
	}
	return report, advisories
}
