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
	"bateia/01_core/parser"
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
					githubURL := parser.GenerateGitHubFileURL(res.Candidate.Owner, res.Candidate.Repo, res.Ref.Commit, f.Path)
					content, isTruncated := parser.TruncateSmart(f.Content, parser.DefaultMaxChars, parser.DefaultMaxLines)

					sb.WriteString(fmt.Sprintf("<details><summary><b>%s</b> (%d bytes) — <a href=\"%s\" target=\"_blank\">Ver no GitHub ↗</a></summary>\n\n", f.Path, f.Size, githubURL))
					sb.WriteString("```go\n")
					sb.WriteString(content)
					if isTruncated {
						sb.WriteString(fmt.Sprintf("\n... 📄 [Conteúdo truncado para visualização no relatório. Ver arquivo completo em: %s]\n", githubURL))
					} else if !strings.HasSuffix(content, "\n") {
						sb.WriteString("\n")
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
	sb.WriteString("Análise do ecossistema técnico e padrões arquiteturais presentes nos repositórios extraídos:\n\n")

	if report.Insights != nil {
		ins := report.Insights

		// 1. Estratégia de Testes
		sb.WriteString("### 🧪 1. Estratégia de Testes\n\n")
		if len(ins.Testing) > 0 {
			sb.WriteString("| Ferramenta | Repositórios (%) | Repositórios Detectados | Versões Comuns |\n")
			sb.WriteString("| :--- | :--- | :--- | :--- |\n")
			for _, m := range ins.Testing {
				vStr := strings.Join(m.Versions, ", ")
				if vStr == "" {
					vStr = "N/A"
				}
				sb.WriteString(fmt.Sprintf("| **%s** | %.0f%% (%d/%d) | %s | %s |\n",
					m.Name, m.Percentage, m.Count, len(report.Results), strings.Join(m.FoundIn, ", "), vStr))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("> Nenhuma ferramenta de teste específica detectada nos manifestos.\n\n")
		}

		// 2. Configuração e Ambiente
		sb.WriteString("### ⚙️ 2. Configuração e Ambiente\n\n")
		if len(ins.Configuration) > 0 {
			sb.WriteString("| Ferramenta | Repositórios (%) | Repositórios Detectados | Versões Comuns |\n")
			sb.WriteString("| :--- | :--- | :--- | :--- |\n")
			for _, m := range ins.Configuration {
				vStr := strings.Join(m.Versions, ", ")
				if vStr == "" {
					vStr = "N/A"
				}
				sb.WriteString(fmt.Sprintf("| **%s** | %.0f%% (%d/%d) | %s | %s |\n",
					m.Name, m.Percentage, m.Count, len(report.Results), strings.Join(m.FoundIn, ", "), vStr))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("> Nenhuma biblioteca padrão de configuração detectada nos manifestos.\n\n")
		}

		// 3. Banco de Dados e Migrações
		sb.WriteString("### 🗄️ 3. Banco de Dados e Migrações\n\n")
		if len(ins.Database) > 0 {
			sb.WriteString("| Ferramenta / Driver | Repositórios (%) | Repositórios Detectados | Versões Comuns |\n")
			sb.WriteString("| :--- | :--- | :--- | :--- |\n")
			for _, m := range ins.Database {
				vStr := strings.Join(m.Versions, ", ")
				if vStr == "" {
					vStr = "N/A"
				}
				sb.WriteString(fmt.Sprintf("| **%s** | %.0f%% (%d/%d) | %s | %s |\n",
					m.Name, m.Percentage, m.Count, len(report.Results), strings.Join(m.FoundIn, ", "), vStr))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("> Nenhum driver ou ORM de banco de dados detectado nos manifestos.\n\n")
		}

		// 4. Observabilidade e Comunicação
		sb.WriteString("### 📊 4. Observabilidade e Comunicação\n\n")
		if len(ins.Observability) > 0 {
			sb.WriteString("| Ferramenta | Repositórios (%) | Repositórios Detectados | Versões Comuns |\n")
			sb.WriteString("| :--- | :--- | :--- | :--- |\n")
			for _, m := range ins.Observability {
				vStr := strings.Join(m.Versions, ", ")
				if vStr == "" {
					vStr = "N/A"
				}
				sb.WriteString(fmt.Sprintf("| **%s** | %.0f%% (%d/%d) | %s | %s |\n",
					m.Name, m.Percentage, m.Count, len(report.Results), strings.Join(m.FoundIn, ", "), vStr))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("> Nenhuma biblioteca de observabilidade detectada nos manifestos.\n\n")
		}

		// 5. Automação e DevOps
		sb.WriteString("### 🔄 5. Automação e DevOps\n\n")
		if len(ins.DevOps) > 0 {
			sb.WriteString("| Arquivo / Recurso | Ocorrências | Repositórios Detectados |\n")
			sb.WriteString("| :--- | :--- | :--- |\n")
			for _, m := range ins.DevOps {
				sb.WriteString(fmt.Sprintf("| **%s** | %d repos | %s |\n",
					m.Name, m.Count, strings.Join(m.FoundIn, ", ")))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("> Nenhum arquivo padrão de automação/DevOps identificado na estrutura.\n\n")
		}
	} else {
		sb.WriteString("- Presença frequente de `go.mod`, arquivos em `cmd/` e `internal/` para separação limpa de domínios.\n")
	}

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
