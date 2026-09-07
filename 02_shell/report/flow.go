/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/flow.md
 * @layer L2
 * @updated 2026-09-07
 */
package report

import (
	"bateia/01_core/domain"
	"fmt"
	"strings"
)

func writeFlowEvidence(sb *strings.Builder, repo domain.ExtractedRepository, ids []string) {
	if repo.Semantic == nil || repo.QueryEvidence == nil {
		return
	}
	for _, u := range repo.Semantic.Units {
		if u.Path != repo.QueryEvidence.Context {
			continue
		}
		if u.Status != "complete" || u.Flow == nil || u.Flow.Status != "complete" {
			reason := "auxiliar ausente ou unidade incompleta"
			if u.Flow != nil {
				reason = u.Flow.Reason
			}
			sb.WriteString("**Fluxo de controle**: não avaliado — " + reason + ".\n\n")
			return
		}
		flow := u.Flow
		sb.WriteString(fmt.Sprintf("**CFG observado**: `%s`, Clang `%s`, auxiliar SHA-256 `%s`; escopo `%s`; poda de arestas trivialmente falsas: %t. %d intervalos ignorados pelo pré-processador. Arestas não comprovam viabilidade nem ordem de execução; polaridade não atribuída.\n\n", flow.Schema, flow.Version, flow.ToolSHA256, flow.Scope, flow.PruneTriviallyFalseEdges, len(flow.Skipped)))
		for _, f := range flow.Functions {
			wanted := false
			for _, d := range repo.QueryEvidence.Definitions {
				for _, id := range ids {
					if d.Symbol.ID == id && d.Symbol.File == f.File && d.Symbol.Line == f.Line {
						wanted = true
					}
				}
			}
			if !wanted {
				continue
			}
			sb.WriteString(fmt.Sprintf("Função `%s`, `%s:%d–%d`: estado %s, entrada B%d, saída B%d. Exibição limitada a 16/%d blocos; grafo integral no JSON.\n\n", f.Name, f.File, f.Line, f.EndLine, f.Status, f.Entry, f.Exit, len(f.Blocks)))
			sb.WriteString("| Bloco | Sucessores | Terminador AST | Fatos no bloco |\n| --- | --- | --- | --- |\n")
			for i, b := range f.Blocks {
				if i >= 16 {
					break
				}
				var next, facts []string
				for _, n := range b.Successors {
					if n == nil {
						next = append(next, "indisponível")
					} else {
						next = append(next, fmt.Sprintf("B%d", *n))
					}
				}
				term := "—"
				if b.Terminator != nil {
					term = fmt.Sprintf("%s (%s:%d)", b.Terminator.Kind, b.Terminator.File, b.Terminator.Line)
				}
				for _, s := range b.Statements {
					if s.Resolution != "" {
						name := s.Callee
						if name == "" {
							name = "alvo não resolvido"
						}
						facts = append(facts, fmt.Sprintf("%s (%s:%d)", name, s.File, s.Line))
					} else if s.Kind == "ReturnStmt" || s.Kind == "BreakStmt" || s.Kind == "ContinueStmt" {
						facts = append(facts, fmt.Sprintf("%s (%s:%d)", s.Kind, s.File, s.Line))
					}
				}
				sb.WriteString(fmt.Sprintf("| B%d | %s | %s | %s |\n", b.ID, strings.Join(next, ", "), term, strings.Join(facts, "; ")))
			}
			sb.WriteString("\n")
		}
		return
	}
}
