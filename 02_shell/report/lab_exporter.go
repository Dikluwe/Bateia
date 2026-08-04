/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/lab_exporter.md
 * @layer L2
 * @updated 2026-08-04
 */

package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

// ExportLabJSON serializa o dataset do Lab de Projetos para JSON.
func ExportLabJSON(dataset *domain.LabDataset) ([]byte, error) {
	return json.MarshalIndent(dataset, "", "  ")
}

// GenerateLabMarkdown compila o dataset do Lab num caderno de exercícios e testes em Markdown.
func GenerateLabMarkdown(dataset *domain.LabDataset) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# 🧪 Caderno de Exercícios & Benchmarks — Lab de Projetos Bateia\n\n"))
	sb.WriteString(fmt.Sprintf("**Query Alvo**: `%s`  \n", dataset.Query))
	sb.WriteString(fmt.Sprintf("**Desafios Minerados**: `%d`  \n", dataset.TotalChallenges))
	sb.WriteString(fmt.Sprintf("**Gerado Em**: %s  \n\n", dataset.GeneratedAt.Format("2006-01-02 15:04:05 MST")))

	sb.WriteString("---\n\n")

	if len(dataset.Challenges) == 0 {
		sb.WriteString("> Nenhum algoritmo ou suíte de testes identificada nos repositórios extraídos.\n\n")
		return sb.String()
	}

	for i, ch := range dataset.Challenges {
		sb.WriteString(fmt.Sprintf("## %d. %s (`%s`)\n\n", i+1, ch.Name, ch.Category))
		sb.WriteString(fmt.Sprintf("- **Repositório Origem**: `%s`\n", ch.Repository))
		sb.WriteString(fmt.Sprintf("- **Arquivo**: `%s`\n", ch.FilePath))
		sb.WriteString(fmt.Sprintf("- **Linguagem**: `%s`\n", ch.Language))
		sb.WriteString(fmt.Sprintf("- **Complexidade Teórica**: `%s`\n\n", ch.Complexity))

		if len(ch.TestCases) > 0 {
			sb.WriteString("### 🧪 Casos de Teste Vinculados\n")
			for _, tc := range ch.TestCases {
				sb.WriteString(fmt.Sprintf("- `%s`\n", tc))
			}
			sb.WriteString("\n")
		}

		if len(ch.Benchmarks) > 0 {
			sb.WriteString("### ⏱️ Suítes de Benchmark Vinculadas\n")
			for _, bm := range ch.Benchmarks {
				sb.WriteString(fmt.Sprintf("- `%s`\n", bm))
			}
			sb.WriteString("\n")
		}

		if ch.CodeSnippet != "" {
			sb.WriteString("### 📝 Trecho de Código de Referência\n\n")
			sb.WriteString(fmt.Sprintf("```%s\n", strings.ToLower(ch.Language)))
			sb.WriteString(ch.CodeSnippet)
			if !strings.HasSuffix(ch.CodeSnippet, "\n") {
				sb.WriteString("\n")
			}
			sb.WriteString("```\n\n")
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}
