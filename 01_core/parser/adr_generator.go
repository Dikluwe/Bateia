/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/adr_generator.md
 * @layer L1
 * @updated 2026-08-26
 */

package parser

import (
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

// GenerateADRDocument compila o SearchReport em um documento formal de ADR (Architecture Decision Record).
func GenerateADRDocument(report *domain.SearchReport) string {
	var sb strings.Builder

	sb.WriteString("# ADR-001: Decisão de Arquitetura de Software\n\n")
	sb.WriteString("**Status**: Aceito (Proposto via Bateia Engine)  \n")
	sb.WriteString(fmt.Sprintf("**Data**: %s  \n", report.GeneratedAt.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Contexto de Pesquisa**: `%s` (Linguagem: %s)  \n\n", report.Query.Query, report.Query.Language))

	sb.WriteString("## 1. Contexto\n\n")
	sb.WriteString(fmt.Sprintf("Foi realizada uma busca e triagem de padrões arquiteturais no ecossistema open-source para a solução de '%s'. ", report.Query.Query))
	sb.WriteString(fmt.Sprintf("O funil de redução do Bateia analisou %d candidatos no topo, aprovou %d na triagem de manifestos e extraiu %d repositórios relevantes.\n\n", report.TopCandidates, report.TriagedCount, report.ExtractedCount))

	sb.WriteString("## 2. Decisão Arquitetural\n\n")
	if len(report.Results) > 0 {
		topRepo := report.Results[0]
		sb.WriteString(fmt.Sprintf("Recomenda-se a adoção/referência do modelo arquitetural do repositório **[%s/%s](%s)** (Score: %d/100).\n\n", topRepo.Candidate.Owner, topRepo.Candidate.Repo, topRepo.Candidate.URL, topRepo.RelevanceScore))
		sb.WriteString(fmt.Sprintf("- **Justificativa**: %s\n", topRepo.RelevanceReason))
		sb.WriteString(fmt.Sprintf("- **Rastreabilidade**: Commit `%s`\n\n", topRepo.Ref.Key()))

		if topRepo.Manifest != nil && len(topRepo.Manifest.Dependencies) > 0 {
			sb.WriteString("### Dependências de Referência Recomendadas:\n\n")
			depCount := 0
			for _, dep := range topRepo.Manifest.Dependencies {
				if !dep.Indirect && depCount < 5 {
					sb.WriteString(fmt.Sprintf("- `%s` (%s)\n", dep.Name, dep.Version))
					depCount++
				}
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("Nenhuma solução de referência atendeu aos critérios estritos de triagem.\n\n")
	}

	sb.WriteString("## 3. Consequências\n\n")
	sb.WriteString("- **Positivas**: Redução de dívida técnica ao reutilizar padrões testados em produção no ecossistema open-source.\n")
	sb.WriteString("- **Negativas/Riscos**: Necessidade de alinhamento de versões de dependências e monitoramento de alertas de segurança (CVEs).\n")

	return sb.String()
}
