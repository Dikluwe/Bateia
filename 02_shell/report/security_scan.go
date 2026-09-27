/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/security_scan.md
 * @layer L2
 * @updated 2026-09-26
 */

package report

import (
	"html"
	"strings"

	"bateia/01_core/domain"
)

func securityScanText(scan *domain.SecurityScanReport) string {
	if scan == nil {
		return "Não avaliado: nenhum resultado de scanner registrado. Ausência de avisos não demonstra segurança."
	}
	var text string
	switch scan.Status {
	case domain.SecurityNotScanned:
		text = "NOT_SCANNED — não avaliado."
	case domain.SecurityScannedNoMatch:
		text = "SCANNED_NO_MATCH — nenhuma correspondência encontrada no escopo e na base consultados; não demonstra ausência de vulnerabilidades."
	case domain.SecurityVulnerabilityFound:
		text = "VULNERABILITY_FOUND — correspondência sinalizada pelo scanner; alcance e impacto precisam de validação."
	default:
		text = "Não avaliado: estado de scanner não reconhecido."
	}
	return text + " Scanner: " + scan.Scanner + ". Escopo: " + scan.Scope + ". Base: " + scan.DatabaseUpdatedAt + ". Limitações: " + strings.Join(scan.Limitations, "; ")
}

func writeSecurityScanHTML(sb *strings.Builder, scan *domain.SecurityScanReport) {
	sb.WriteString("    <h4>Estado da análise de segurança</h4>\n    <p>" + html.EscapeString(securityScanText(scan)) + "</p>\n")
}

func writeSecurityScanMarkdown(sb *strings.Builder, scan *domain.SecurityScanReport) {
	text := html.EscapeString(strings.ReplaceAll(strings.ReplaceAll(securityScanText(scan), "\n", " "), "\r", " "))
	// Entities prevent link/fence syntax from being parsed as Markdown.
	text = strings.NewReplacer("[", "&#91;", "]", "&#93;", "`", "&#96;", "*", "&#42;", "\\", "&#92;").Replace(text)
	sb.WriteString("- **Estado da análise de segurança**: " + text + "\n")
}
