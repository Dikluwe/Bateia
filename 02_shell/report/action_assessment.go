/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/action_assessment.md
 * @layer L2
 * @updated 2026-09-08
 */
package report

import (
	"bateia/01_core/domain"
	"fmt"
	"html"
	"sort"
	"strconv"
	"strings"
)

// actionDisplay keeps untrusted text inert, including HTML, Markdown and newlines.
func actionDisplay(value string) string {
	quoted := html.EscapeString(strconv.Quote(value))
	// Numeric entities also neutralize Markdown when the renderer drops HTML tags.
	quoted = strings.NewReplacer("`", "&#96;", "*", "&#42;", "_", "&#95;", "[", "&#91;", "]", "&#93;", "\\", "&#92;").Replace(quoted)
	return "<code>" + quoted + "</code>"
}

// WriteActionAssessments presents bounded hints while JSON retains every candidate.
func WriteActionAssessments(sb *strings.Builder, assessments []domain.RequirementActionAssessment) {
	if len(assessments) == 0 {
		return
	}
	sb.WriteString("#### Hipóteses de ação por candidato\n\nRótulos e chamadas não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade. A ordem dos candidatos prioriza leitura, não confiança semântica.\n\n")
	for _, a := range assessments {
		fmt.Fprintf(sb, "- Consulta %s; ação %s; estado %s; motivo %s.\n", actionDisplay(a.Query), actionDisplay(a.RequiredAction), actionDisplay(a.Status), actionDisplay(a.Reason))
		candidates := append([]domain.ActionCandidate(nil), a.Candidates...)
		priority := func(relation string) int {
			switch relation {
			case "aligned_hint":
				return 0
			case "conflicting_hints":
				return 1
			case "opposed_hint":
				return 2
			default:
				return 3
			}
		}
		sort.SliceStable(candidates, func(i, j int) bool {
			return priority(candidates[i].Relation) < priority(candidates[j].Relation)
		})
		for i, c := range candidates {
			if i == 12 {
				fmt.Fprintf(sb, "  Mais %d candidatos disponíveis no JSON.\n", len(a.Candidates)-i)
				break
			}
			fmt.Fprintf(sb, "  - Relação %s; função %s; arquivo %s; SHA-256 %s; linhas %d–%d; motivo %s.\n", actionDisplay(c.Relation), actionDisplay(c.Source.Symbol), actionDisplay(c.Source.Path), actionDisplay(c.Source.SHA256), c.Source.Line, c.EndLine, actionDisplay(c.Reason))
			for j, e := range c.Evidence {
				if j == 4 {
					fmt.Fprintf(sb, "    Mais %d evidências disponíveis no JSON.\n", len(c.Evidence)-j)
					break
				}
				fmt.Fprintf(sb, "    - Evidência %s; ação %s; termo %s; linhas %d–%d; trecho %s.\n", actionDisplay(e.Kind), actionDisplay(e.Action), actionDisplay(e.Term), e.Line, e.EndLine, actionDisplay(e.Snippet))
			}
		}
	}
	sb.WriteString("\n")
}
