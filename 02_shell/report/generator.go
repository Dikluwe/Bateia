/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/generator.md
 * @layer L2
 * @updated 2026-07-25
 */

package report

import (
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

// GenerateMarkdownReport compila o objeto domain.SearchReport no documento Markdown de apresentação L2.
func GenerateMarkdownReport(report *domain.SearchReport) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Relatório de Análise Arquitetural — Bateia\n\n"))
	sb.WriteString(fmt.Sprintf("**Query de Busca**: `%s`  \n", report.Query.Query))
	sb.WriteString(fmt.Sprintf("**Linguagem**: `%s`  \n", report.Query.Language))
	sb.WriteString(fmt.Sprintf("**Gerado em**: %s  \n\n", report.GeneratedAt.Format("2006-01-02 15:04:05 MST")))

	sb.WriteString("## Resumo do Funil de Redução\n\n")
	sb.WriteString(fmt.Sprintf("| Estágio | Quantidade | Observação |\n"))
	sb.WriteString(fmt.Sprintf("| :--- | :--- | :--- |\n"))
	sb.WriteString(fmt.Sprintf("| **1. Descoberta (Topo)** | %d candidatos | Fonte: Sourcegraph (sem consumo de Search API) |\n", report.TopCandidates))
	sb.WriteString(fmt.Sprintf("| **2. Triagem (Meio)** | %d aprovados | Leitura de manifestos e estrutura (%d descartados) |\n", report.TriagedCount, report.DiscardedCount))
	sb.WriteString(fmt.Sprintf("| **3. Extração (Fundo)** | %d entregues | Extração idempotente e cacheados por `repo@commit` |\n\n", report.ExtractedCount))

	sb.WriteString("---\n\n")
	sb.WriteString("## 1. Soluções Exatas Encontradas (Com Rastreabilidade `repo@commit`)\n\n")

	if len(report.Results) == 0 {
		sb.WriteString("> Nenhuma solução encontrada para os critérios informados. Considere alargar os termos da query ou remover filtros restritivos.\n\n")
	} else {
		for i, res := range report.Results {
			sb.WriteString(fmt.Sprintf("### %d. [%s/%s](%s)\n\n", i+1, res.Candidate.Owner, res.Candidate.Repo, res.Candidate.URL))
			sb.WriteString(fmt.Sprintf("- **Score de Relevância**: `%d/100` (`%s`)\n", res.RelevanceScore, res.RelevanceLabel))
			sb.WriteString(fmt.Sprintf("- **Justificativa**: %s\n", res.RelevanceReason))
			sb.WriteString(fmt.Sprintf("- **Commit Rastreável**: `%s`\n", res.Ref.Key()))
			sb.WriteString(fmt.Sprintf("- **Estrelas**: %s\n", formatStars(res.Candidate.Stars)))
			sb.WriteString(fmt.Sprintf("- **Descrição**: %s\n", res.Candidate.Description))
			sb.WriteString(fmt.Sprintf("- **Consultado Em**: %s\n\n", res.ConsultadoEm.Format("2006-01-02 15:04:05")))

			if res.Manifest != nil {
				sb.WriteString(fmt.Sprintf("#### Módulo & Dependências (`go.mod`)\n"))
				sb.WriteString(fmt.Sprintf("- **Módulo**: `%s` (Go %s)\n", res.Manifest.ModuleName, res.Manifest.GoVersion))
				if len(res.Manifest.Dependencies) > 0 {
					sb.WriteString("- **Principais Dependências Directas**:\n")
					depCount := 0
					for _, dep := range res.Manifest.Dependencies {
						if !dep.Indirect && depCount < 5 {
							sb.WriteString(fmt.Sprintf("  - `%s` (%s)\n", dep.Name, dep.Version))
							depCount++
						}
					}
				}
				sb.WriteString("\n")
			}

			if len(res.Files) > 0 {
				sb.WriteString("#### Arquivos Extraídos\n\n")
				for _, f := range res.Files {
					sb.WriteString(fmt.Sprintf("<details><summary><b>%s</b> (%d bytes)</summary>\n\n", f.Path, f.Size))
					sb.WriteString("```go\n")
					if len(f.Content) > 1500 {
						sb.WriteString(f.Content[:1500] + "\n... (conteúdo truncado para visualização no relatório)\n")
					} else {
						sb.WriteString(f.Content + "\n")
					}
					sb.WriteString("```\n</details>\n\n")
				}
			}
			sb.WriteString("---\n\n")
		}
	}

	sb.WriteString("## 2. Ideias Alternativas & Abordagens Encontradas\n\n")
	sb.WriteString("Agrupamento das técnicas identificadas nos repositórios analisados:\n\n")
	for _, res := range report.Results {
		if res.Manifest != nil {
			frameworks := make([]string, 0)
			for _, dep := range res.Manifest.Dependencies {
				if strings.Contains(dep.Name, "gin") || strings.Contains(dep.Name, "fiber") || strings.Contains(dep.Name, "chi") || strings.Contains(dep.Name, "echo") || strings.Contains(dep.Name, "grpc") {
					frameworks = append(frameworks, dep.Name)
				}
			}
			if len(frameworks) > 0 {
				sb.WriteString(fmt.Sprintf("- **%s/%s**: Utiliza os frameworks `%s`\n", res.Candidate.Owner, res.Candidate.Repo, strings.Join(frameworks, ", ")))
			}
		}
	}

	sb.WriteString("\n## 3. Descobrimento de Necessidades Ocultas\n\n")
	sb.WriteString("Padrões estruturais e arquivos de configuração detectados que podem ser relevantes:\n\n")
	sb.WriteString("- Presença frequente de `go.mod`, arquivos em `cmd/` e `internal/` para separação limpa de domínios.\n")

	return sb.String()
}

func formatStars(stars int) string {
	switch {
	case stars >= 1_000_000:
		return fmt.Sprintf("⭐ %.1fM", float64(stars)/1_000_000)
	case stars >= 1_000:
		return fmt.Sprintf("⭐ %.1fK (%d)", float64(stars)/1_000, stars)
	case stars > 0:
		return fmt.Sprintf("⭐ %d", stars)
	default:
		return "⭐ N/A (metadado indisponível)"
	}
}
