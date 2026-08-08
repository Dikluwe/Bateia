/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/generator.md
 * @layer L2
 * @updated 2026-07-25
 */

package report

import (
	"fmt"
	"path/filepath"
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

	if report.Telemetry != nil && report.Telemetry.TotalDurationMs > 0 {
		sb.WriteString("#### ⏱️ Telemetria de Desempenho do Funil\n\n")
		sb.WriteString(fmt.Sprintf("- **Tempo Total de Execução**: `%d ms`  \n", report.Telemetry.TotalDurationMs))
		if report.Telemetry.DiscoveryDurationMs > 0 {
			sb.WriteString(fmt.Sprintf("- **Tempo de Descoberta (Estágio 1)**: `%d ms`  \n", report.Telemetry.DiscoveryDurationMs))
		}
		if report.Telemetry.TriageDurationMs > 0 {
			sb.WriteString(fmt.Sprintf("- **Tempo de Triagem (Estágio 2)**: `%d ms`  \n", report.Telemetry.TriageDurationMs))
		}
		if report.Telemetry.ExtractionDurationMs > 0 {
			sb.WriteString(fmt.Sprintf("- **Tempo de Extração (Estágio 3)**: `%d ms`  \n", report.Telemetry.ExtractionDurationMs))
		}
		sb.WriteString(fmt.Sprintf("- **Métricas de Cache**: `%d Hits` / `%d Misses`  \n\n", report.Telemetry.CacheHits, report.Telemetry.CacheMisses))
	}

	mermaidGraph := GenerateMermaidGraph(report)
	if mermaidGraph != "" {
		sb.WriteString("### 📊 Grafo de Dependências (Mermaid.js)\n\n")
		sb.WriteString(mermaidGraph)
		sb.WriteString("\n")
	}

	if len(report.PolyglotMatrix) > 0 {
		sb.WriteString("### 🌐 Matriz Polyglot de Comunicação & Protocolos\n\n")
		sb.WriteString("| Protocolo / Categoria | Tecnologia | Repositories Encontrados |\n")
		sb.WriteString("| :--- | :--- | :--- |\n")
		for _, poly := range report.PolyglotMatrix {
			sb.WriteString(fmt.Sprintf("| **%s** | `%s` | `%s` |\n", poly.Protocol, poly.Technology, strings.Join(poly.Repositories, ", ")))
		}
		sb.WriteString("\n")
	}

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
				langName := res.Manifest.Language
				if langName == "" {
					langName = "Desconhecida"
				} else {
					langName = strings.ToUpper(langName[:1]) + langName[1:]
				}

				manifestName := "Manifesto"
				switch res.Manifest.Language {
				case "go":
					manifestName = "go.mod"
				case "typescript", "javascript", "node":
					manifestName = "package.json"
				case "python":
					manifestName = "pyproject.toml/requirements.txt"
				case "rust":
					manifestName = "Cargo.toml"
				case "ruby":
					manifestName = "Gemfile"
				case "java", "kotlin":
					manifestName = "pom.xml/build.gradle"
				}

				sb.WriteString(fmt.Sprintf("#### Módulo & Dependências (`%s`)\n", manifestName))
				
				versionStr := ""
				if res.Manifest.GoVersion != "" {
					versionStr = " " + res.Manifest.GoVersion
				}
				sb.WriteString(fmt.Sprintf("- **Módulo**: `%s` (%s%s)\n", res.Manifest.ModuleName, langName, versionStr))
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

			if len(res.Symbols) > 0 {
				sb.WriteString("#### Símbolos Estruturais Extraídos (AST)\n")
				symCount := 0
				for _, sym := range res.Symbols {
					if symCount < 8 {
						sb.WriteString(fmt.Sprintf("- **`%s`** (%s) — em `%s` (linha %d)\n", sym.Name, sym.Kind, filepath.Base(sym.File), sym.Line))
						symCount++
					}
				}
				sb.WriteString("\n")
			}

			if len(res.ArchViolations) > 0 {
				sb.WriteString("#### ⚠️ Violações Arquiteturais Detectadas (Arch-Linter)\n\n")
				for _, v := range res.ArchViolations {
					sb.WriteString(fmt.Sprintf("- **[%s] %s**: %s (alvo: `%s`)\n", v.Severity, v.Rule, v.Message, v.Target))
					if v.CodeSnippet != "" || v.Suggestion != "" {
						sb.WriteString("```text\n")
						sb.WriteString(fmt.Sprintf("  Error: [%s]\n", v.Rule))
						lineNo := v.Line
						if lineNo <= 0 {
							lineNo = 1
						}
						sb.WriteString(fmt.Sprintf("     ╭─[%s:%d:1]\n", v.Target, lineNo))
						if v.CodeSnippet != "" {
							sb.WriteString(fmt.Sprintf("  %d │ %s\n", lineNo, v.CodeSnippet))
							sb.WriteString("     · ^^^^^^^^^^^^^^\n")
						}
						if v.Suggestion != "" {
							sb.WriteString(fmt.Sprintf("     ╰─[Dica: %s]\n", v.Suggestion))
						}
						sb.WriteString("```\n")
					}
				}
				sb.WriteString("\n")
			}

			if len(res.SecurityAdvisories) > 0 {
				sb.WriteString("#### 🛡️ Alertas de Segurança & CVEs (Security Scanner)\n")
				for _, adv := range res.SecurityAdvisories {
					sb.WriteString(fmt.Sprintf("- **[%s] %s** (%s): %s em `%s` (%s)\n", adv.Severity, adv.AdvisoryID, adv.Package, adv.Message, adv.Package, adv.VulnerableVersion))
				}
				sb.WriteString("\n")
			}

			if len(res.Files) > 0 {
				sb.WriteString("#### Arquivos Extraídos\n\n")
				for _, f := range res.Files {
					githubURL := parser.GenerateGitHubFileURL(res.Candidate.Owner, res.Candidate.Repo, res.Ref.Commit, f.Path)
					content, isTruncated := parser.TruncateSmart(f.Content, parser.DefaultMaxChars, parser.DefaultMaxLines)

					sb.WriteString(fmt.Sprintf("<details><summary><b>%s</b> (%d bytes) — <a href=\"%s\" target=\"_blank\">Ver no GitHub ↗</a></summary>\n\n", f.Path, f.Size, githubURL))
					
					ext := strings.TrimPrefix(filepath.Ext(f.Path), ".")
					if ext == "" {
						ext = "text"
					}
					sb.WriteString(fmt.Sprintf("```%s\n", ext))
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
		if res.Manifest != nil && len(res.Manifest.Dependencies) > 0 {
			deps := make([]string, 0)
			for _, dep := range res.Manifest.Dependencies {
				if !dep.Indirect {
					deps = append(deps, dep.Name)
					if len(deps) >= 6 {
						break
					}
				}
			}
			if len(deps) > 0 {
				sb.WriteString(fmt.Sprintf("- **%s/%s**: Principais pacotes/ferramentas: `%s`\n", res.Candidate.Owner, res.Candidate.Repo, strings.Join(deps, ", ")))
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
