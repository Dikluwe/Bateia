/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/mermaid.md
 * @layer L2
 * @updated 2026-08-03
 */

package report

import (
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

// GenerateMermaidGraph gera o diagrama de grafos em sintaxe Mermaid (graph TD).
func GenerateMermaidGraph(report *domain.SearchReport) string {
	if len(report.Results) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("```mermaid\ngraph TD\n")

	// Nó raiz
	queryNode := sanitizeMermaidID("Query: " + report.Query.Query)
	sb.WriteString(fmt.Sprintf("    Root[\"%s\"] --> Candidates\n", report.Query.Query))

	for i, res := range report.Results {
		repoID := fmt.Sprintf("Repo_%d", i+1)
		repoLabel := fmt.Sprintf("%s/%s", res.Candidate.Owner, res.Candidate.Repo)

		sb.WriteString(fmt.Sprintf("    Root --> %s[\"%s\"]\n", repoID, repoLabel))

		if res.Manifest != nil && len(res.Manifest.Dependencies) > 0 {
			depCount := 0
			for _, dep := range res.Manifest.Dependencies {
				if !dep.Indirect && depCount < 3 {
					depID := sanitizeMermaidID(dep.Name)
					sb.WriteString(fmt.Sprintf("    %s --> %s[\"%s\"]\n", repoID, depID, dep.Name))
					depCount++
				}
			}
		}
	}

	_ = queryNode
	sb.WriteString("```\n")
	return sb.String()
}

func sanitizeMermaidID(raw string) string {
	cleaned := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, raw)
	if cleaned == "" {
		return "Node"
	}
	return "N_" + cleaned
}
