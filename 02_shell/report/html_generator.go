/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/html_generator.md
 * @layer L2
 * @updated 2026-08-03
 */

package report

import (
	"fmt"
	"html"
	"path/filepath"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

// GenerateHTMLReport compila o domain.SearchReport em um documento HTML standalone.
func GenerateHTMLReport(report *domain.SearchReport) string {
	var sb strings.Builder

	sb.WriteString("<!DOCTYPE html>\n<html lang=\"pt-BR\">\n<head>\n")
	sb.WriteString("  <meta charset=\"UTF-8\">\n")
	sb.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	sb.WriteString(fmt.Sprintf("  <title>Bateia v1.0 — Relatório: %s</title>\n", html.EscapeString(report.Query.Query)))
	sb.WriteString("  <style>\n")
	sb.WriteString("    :root { --bg: #0d1117; --card: #161b22; --border: #30363d; --text: #c9d1d9; --accent: #58a6ff; --green: #3fb950; --gold: #d29922; }\n")
	sb.WriteString("    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background: var(--bg); color: var(--text); padding: 2rem; margin: 0; line-height: 1.6; }\n")
	sb.WriteString("    .container { max-width: 1000px; margin: 0 auto; }\n")
	sb.WriteString("    header { border-bottom: 1px solid var(--border); padding-bottom: 1rem; margin-bottom: 2rem; }\n")
	sb.WriteString("    h1 { color: var(--accent); margin: 0 0 0.5rem 0; font-size: 1.8rem; }\n")
	sb.WriteString("    .meta { color: #8b949e; font-size: 0.9rem; }\n")
	sb.WriteString("    .badge { display: inline-block; padding: 0.2rem 0.6rem; border-radius: 12px; font-weight: bold; font-size: 0.8rem; background: var(--border); color: #fff; }\n")
	sb.WriteString("    .badge-green { background: #238636; }\n")
	sb.WriteString("    .badge-gold { background: #9e6a03; }\n")
	sb.WriteString("    .card { background: var(--card); border: 1px solid var(--border); border-radius: 8px; padding: 1.5rem; margin-bottom: 1.5rem; }\n")
	sb.WriteString("    .card h3 { margin-top: 0; color: var(--accent); }\n")
	sb.WriteString("    table { width: 100%; border-collapse: collapse; margin: 1rem 0; }\n")
	sb.WriteString("    th, td { border: 1px solid var(--border); padding: 0.6rem 0.8rem; text-align: left; }\n")
	sb.WriteString("    th { background: #21262d; color: var(--accent); }\n")
	sb.WriteString("    pre { background: #010409; padding: 1rem; border-radius: 6px; overflow-x: auto; font-family: SFMono-Regular, Consolas, 'Liberation Mono', Menlo, monospace; font-size: 0.85rem; border: 1px solid var(--border); }\n")
	sb.WriteString("    details { margin-top: 1rem; background: #0d1117; border: 1px solid var(--border); border-radius: 6px; padding: 0.5rem 1rem; }\n")
	sb.WriteString("    summary { cursor: pointer; font-weight: bold; color: var(--accent); }\n")
	sb.WriteString("  </style>\n")
	sb.WriteString("  <script type=\"module\">\n")
	sb.WriteString("    import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';\n")
	sb.WriteString("    mermaid.initialize({ startOnLoad: true, theme: 'dark' });\n")
	sb.WriteString("  </script>\n</head>\n<body>\n<div class=\"container\">\n")

	sb.WriteString("  <header>\n")
	sb.WriteString(fmt.Sprintf("    <h1>🔍 Bateia v1.0 — Relatório Arquitetural</h1>\n"))
	sb.WriteString(fmt.Sprintf("    <div class=\"meta\">Query: <strong>%s</strong> | Linguagem: <strong>%s</strong> | Gerado em: %s</div>\n", html.EscapeString(report.Query.Query), html.EscapeString(report.Query.Language), report.GeneratedAt.Format("2006-01-02 15:04:05 MST")))
	sb.WriteString("  </header>\n")

	sb.WriteString("  <section class=\"card\">\n")
	sb.WriteString("    <h2>Resumo do Funil de Redução</h2>\n")
	sb.WriteString("    <table>\n")
	sb.WriteString("      <thead><tr><th>Estágio</th><th>Quantidade</th><th>Observação</th></tr></thead>\n")
	sb.WriteString("      <tbody>\n")
	sb.WriteString(fmt.Sprintf("        <tr><td><strong>1. Descoberta</strong></td><td>%d candidatos</td><td>Fonte: Sourcegraph</td></tr>\n", report.TopCandidates))
	sb.WriteString(fmt.Sprintf("        <tr><td><strong>2. Triagem</strong></td><td>%d aprovados</td><td>Leitura de manifestos (%d descartados)</td></tr>\n", report.TriagedCount, report.DiscardedCount))
	sb.WriteString(fmt.Sprintf("        <tr><td><strong>3. Extração</strong></td><td>%d entregues</td><td>Extração idempotente por commit</td></tr>\n", report.ExtractedCount))
	sb.WriteString("      </tbody>\n    </table>\n")

	mermaidGraph := GenerateMermaidGraph(report)
	if mermaidGraph != "" {
		cleanMermaid := strings.TrimPrefix(strings.TrimSuffix(mermaidGraph, "```\n"), "```mermaid\n")
		sb.WriteString("    <h3>📊 Grafo de Dependências</h3>\n")
		sb.WriteString(fmt.Sprintf("    <pre class=\"mermaid\">\n%s    </pre>\n", cleanMermaid))
	}

	if len(report.PolyglotMatrix) > 0 {
		sb.WriteString("    <h3>🌐 Matriz Polyglot de Comunicação</h3>\n")
		sb.WriteString("    <table>\n      <thead><tr><th>Protocolo</th><th>Tecnologia</th><th>Repositórios</th></tr></thead>\n      <tbody>\n")
		for _, poly := range report.PolyglotMatrix {
			sb.WriteString(fmt.Sprintf("        <tr><td><strong>%s</strong></td><td><code>%s</code></td><td>%s</td></tr>\n", html.EscapeString(poly.Protocol), html.EscapeString(poly.Technology), html.EscapeString(strings.Join(poly.Repositories, ", "))))
		}
		sb.WriteString("      </tbody>\n    </table>\n")
	}

	sb.WriteString("  </section>\n")

	sb.WriteString("  <h2>Repositórios Extraídos</h2>\n")
	for i, res := range report.Results {
		badgeClass := "badge-gold"
		if res.RelevanceScore >= 80 {
			badgeClass = "badge-green"
		}

		sb.WriteString("  <article class=\"card\">\n")
		sb.WriteString(fmt.Sprintf("    <h3>%d. <a href=\"%s\" target=\"_blank\" style=\"color:var(--accent);text-decoration:none;\">%s/%s</a> <span class=\"badge %s\">Score %d/100</span></h3>\n", i+1, res.Candidate.URL, html.EscapeString(res.Candidate.Owner), html.EscapeString(res.Candidate.Repo), badgeClass, res.RelevanceScore))
		sb.WriteString(fmt.Sprintf("    <p><strong>Justificativa:</strong> %s</p>\n", html.EscapeString(res.RelevanceReason)))
		sb.WriteString(fmt.Sprintf("    <p class=\"meta\">Commit: <code>%s</code></p>\n", html.EscapeString(res.Ref.Key())))

		if len(res.Symbols) > 0 {
			sb.WriteString("    <h4>Símbolos Estruturais Extraídos (AST)</h4>\n    <ul>\n")
			symCount := 0
			for _, sym := range res.Symbols {
				if symCount < 8 {
					sb.WriteString(fmt.Sprintf("      <li><strong><code>%s</code></strong> (%s) — em <code>%s</code> (linha %d)</li>\n", html.EscapeString(sym.Name), html.EscapeString(sym.Kind), html.EscapeString(filepath.Base(sym.File)), sym.Line))
					symCount++
				}
			}
			sb.WriteString("    </ul>\n")
		}

		if len(res.ArchViolations) > 0 {
			sb.WriteString("    <h4>⚠️ Violações Arquiteturais Detectadas (Arch-Linter)</h4>\n    <ul>\n")
			for _, v := range res.ArchViolations {
				sb.WriteString(fmt.Sprintf("      <li><strong style=\"color:#f85149;\">[%s] %s</strong>: %s (alvo: <code>%s</code>)</li>\n", html.EscapeString(v.Severity), html.EscapeString(v.Rule), html.EscapeString(v.Message), html.EscapeString(v.Target)))
			}
			sb.WriteString("    </ul>\n")
		}

		if len(res.SecurityAdvisories) > 0 {
			sb.WriteString("    <h4>🛡️ Alertas de Segurança & CVEs (Security Scanner)</h4>\n    <ul>\n")
			for _, adv := range res.SecurityAdvisories {
				sb.WriteString(fmt.Sprintf("      <li><strong style=\"color:#da3633;\">[%s] %s</strong>: %s em <code>%s</code> (%s)</li>\n", html.EscapeString(adv.Severity), html.EscapeString(adv.AdvisoryID), html.EscapeString(adv.Message), html.EscapeString(adv.Package), html.EscapeString(adv.VulnerableVersion)))
			}
			sb.WriteString("    </ul>\n")
		}

		if len(res.Files) > 0 {
			sb.WriteString("    <h4>Arquivos Extraídos</h4>\n")
			for _, f := range res.Files {
				githubURL := parser.GenerateGitHubFileURL(res.Candidate.Owner, res.Candidate.Repo, res.Ref.Commit, f.Path)
				content, _ := parser.TruncateSmart(f.Content, parser.DefaultMaxChars, parser.DefaultMaxLines)

				sb.WriteString(fmt.Sprintf("    <details><summary>%s (%d bytes)</summary>\n", html.EscapeString(f.Path), f.Size))
				sb.WriteString(fmt.Sprintf("      <p><a href=\"%s\" target=\"_blank\" style=\"color:var(--accent);\">Ver no GitHub ↗</a></p>\n", githubURL))
				sb.WriteString(fmt.Sprintf("      <pre><code>%s</code></pre>\n", html.EscapeString(content)))
				sb.WriteString("    </details>\n")
			}
		}

		sb.WriteString("  </article>\n")
	}

	sb.WriteString("</div>\n</body>\n</html>\n")
	return sb.String()
}
