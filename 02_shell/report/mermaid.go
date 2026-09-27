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
	"unicode"

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
	sb.WriteString(fmt.Sprintf("    Root[\"%s\"] --> Candidates\n", escapeMermaidLabel(report.Query.Query)))

	for i, res := range report.Results {
		repoID := fmt.Sprintf("Repo_%d", i+1)
		repoLabel := fmt.Sprintf("%s/%s", res.Candidate.Owner, res.Candidate.Repo)

		sb.WriteString(fmt.Sprintf("    Root --> %s[\"%s\"]\n", repoID, escapeMermaidLabel(repoLabel)))

		if res.Manifest != nil && len(res.Manifest.Dependencies) > 0 {
			depCount := 0
			for _, dep := range res.Manifest.Dependencies {
				if !dep.Indirect && depCount < 3 {
					depID := fmt.Sprintf("Dep_%d_%d", i+1, depCount+1)
					sb.WriteString(fmt.Sprintf("    %s --> %s[\"%s\"]\n", repoID, depID, escapeMermaidLabel(dep.Name)))
					depCount++
				}
			}
		}
	}

	sb.WriteString("```\n")
	return sb.String()
}

// Mermaid entities are not HTML escaping: syntax delimiters, entity markers,
// backticks and markup characters must not leave the quoted node label.
func escapeMermaidLabel(raw string) string {
	var sb strings.Builder
	for _, r := range raw {
		switch {
		case unicode.IsControl(r) || r == '\u2028' || r == '\u2029':
			sb.WriteByte(' ')
		case unicode.IsLetter(r) || unicode.IsNumber(r) || strings.ContainsRune(" /._:-", r):
			sb.WriteRune(r)
		default:
			fmt.Fprintf(&sb, "#%d;", r)
		}
	}
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
