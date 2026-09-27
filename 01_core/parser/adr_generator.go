/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/adr_generator.md
 * @layer L1
 * @updated 2026-09-26
 */

package parser

import (
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

// GenerateADRDocument compila o SearchReport em uma proposta de ADR, sem inferir aprovação.
func GenerateADRDocument(report *domain.SearchReport) string {
	var sb strings.Builder

	sb.WriteString("# ADR-001: Decisão de Arquitetura de Software\n\n")
	sb.WriteString("**Status**: PROPOSTO (gerado a partir de pesquisa; não constitui aprovação)  \n")
	sb.WriteString(fmt.Sprintf("**Data**: %s  \n", report.GeneratedAt.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Contexto de Pesquisa**: `%s` (Linguagem: %s)  \n\n", report.Query.Query, report.Query.Language))

	sb.WriteString("## 1. Contexto\n\n")
	sb.WriteString(fmt.Sprintf("Foi realizada uma busca e triagem de padrões arquiteturais no ecossistema open-source para a solução de '%s'. ", report.Query.Query))
	sb.WriteString(fmt.Sprintf("O funil de redução do Bateia analisou %d candidatos no topo, selecionou %d na triagem de manifestos e extraiu %d repositórios relevantes.\n\n", report.TopCandidates, report.TriagedCount, report.ExtractedCount))

	sb.WriteString("## 2. Proposta Arquitetural\n\n")
	if len(report.Results) > 0 {
		topRepo := report.Results[0]
		sb.WriteString(fmt.Sprintf("Propõe-se avaliar o modelo arquitetural do candidato **[%s/%s](%s)** (Score de relevância da pesquisa: %d/100; não é garantia de adequação).\n\n", topRepo.Candidate.Owner, topRepo.Candidate.Repo, topRepo.Candidate.URL, topRepo.RelevanceScore))
		sb.WriteString(fmt.Sprintf("- **Justificativa da relevância**: %s\n", topRepo.RelevanceReason))
		if topRepo.Ref.Commit != "" && topRepo.Ref.Commit != "HEAD" {
			sb.WriteString(fmt.Sprintf("- **Rastreabilidade**: Commit `%s`\n\n", topRepo.Ref.Key()))
		} else if topRepo.Ref.Owner != "" || topRepo.Ref.Repo != "" {
			sb.WriteString(fmt.Sprintf("- **Rastreabilidade**: Referência de repositório `%s`; commit não fixado.\n\n", topRepo.Ref.Key()))
		} else {
			sb.WriteString("- **Rastreabilidade**: Commit não informado neste relatório.\n\n")
		}

		if topRepo.Manifest != nil && len(topRepo.Manifest.Dependencies) > 0 {
			depCount := 0
			for _, dep := range topRepo.Manifest.Dependencies {
				if !dep.Indirect && depCount < 5 {
					if depCount == 0 {
						sb.WriteString("### Dependências diretas observadas no candidato (não aprovadas):\n\n")
					}
					sb.WriteString(fmt.Sprintf("- `%s` (%s)\n", dep.Name, dep.Version))
					depCount++
				}
			}
			if depCount > 0 {
				sb.WriteString("\n")
			}
		}
	} else {
		sb.WriteString("Este relatório não contém candidatos extraídos para propor uma referência arquitetural. Isso não demonstra a inexistência de alternativas adequadas.\n\n")
	}

	sb.WriteString("## 3. Consequências possíveis e validações pendentes\n\n")
	if len(report.Results) > 0 {
		sb.WriteString("- **Hipótese de benefício**: O reúso pode diminuir esforço de implementação; a pesquisa não comprova esse efeito nem a adequação ao contexto de destino.\n")
		sb.WriteString("- **Validações pendentes**: Avaliar compatibilidade funcional, versões, licença, custo de integração e riscos de segurança antes de decidir pela adoção.\n")
	} else {
		sb.WriteString("- **Próximo passo**: Revisar a pergunta e os critérios de busca antes de formular uma proposta de adoção.\n")
	}

	return sb.String()
}
