/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/investigation_plan.md
 * @layer L1
 * @updated 2026-09-16
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func requirementID(kind, text string) string {
	sum := sha256.Sum256([]byte(kind + "\x00" + strings.ToLower(normalizeIdeaText(text))))
	return fmt.Sprintf("req-%s-%x", kind, sum[:8])
}

func investigationRequirements(idea domain.IdeaSpec) []domain.InvestigationRequirement {
	out := []domain.InvestigationRequirement{}
	seen := map[string]bool{}
	add := func(kind, text string) {
		text = normalizeIdeaText(text)
		if text == "" {
			return
		}
		id := requirementID(kind, text)
		if seen[id] {
			return
		}
		seen[id] = true
		acceptance := "Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas."
		if kind == "non_goal" {
			acceptance = "Examinar presença e negação do comportamento excluído; ausência de menção não demonstra ausência do comportamento."
		}
		if kind == "preference" {
			acceptance = "Registrar compatibilidade observada; preferência não é gate eliminatório."
		}
		if kind == "quality" {
			acceptance = "Definir e observar critério mensurável sob condições declaradas; não inferir limiar da formulação textual."
		}
		if kind == "operating_condition" {
			acceptance = "Verificar a condição no ambiente de execução declarado; documentação e código isolados não comprovam operação."
		}
		out = append(out, domain.InvestigationRequirement{ID: id, Kind: kind, Text: text, Question: "Que evidência responde a: " + text + "?", Acceptance: acceptance})
	}
	add("problem", idea.Problem)
	for _, v := range idea.Capabilities {
		add("capability", v)
	}
	for _, v := range idea.Qualities {
		add("quality", v)
	}
	for _, v := range idea.OperatingConditions {
		add("operating_condition", v)
	}
	for _, v := range idea.Constraints {
		add("constraint", v)
	}
	for _, v := range idea.NonGoals {
		add("non_goal", v)
	}
	for _, v := range idea.PreferredStack {
		add("preference", v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// BuildInvestigationPlan never treats query terms as proof about a candidate.
func BuildInvestigationPlan(idea domain.IdeaSpec, snapshot domain.DiscoverySnapshot, assessments []domain.CandidateAssessment) domain.InvestigationPlan {
	return BuildInvestigationPlanWithCatalogCoverage(idea, snapshot, assessments, domain.EvidenceCatalogCoverage{})
}

func BuildInvestigationPlanWithCatalogCoverage(idea domain.IdeaSpec, snapshot domain.DiscoverySnapshot, assessments []domain.CandidateAssessment, coverage domain.EvidenceCatalogCoverage) domain.InvestigationPlan {
	out := domain.InvestigationPlan{SchemaVersion: "1", Requirements: investigationRequirements(idea), Searches: []domain.InvestigationSearch{}, Actions: []domain.InvestigationAction{}, StopReason: "RECORDED_SEARCH_FINISHED", Scope: "Somente candidatos, tentativas e evidências registrados; não demonstra busca global completa nem integração verificada.", TelemetryLabel: "Tempos de busca são telemetria do Bateia, não benchmarks dos candidatos."}
	if snapshot.SchemaVersion == "" && len(snapshot.Candidates) == 0 && len(snapshot.PerExpansion) == 0 {
		out.StopReason = "NO_CANDIDATES_IN_SCOPE"
		out.Scope = "Descoberta não registrada. " + out.Scope
	} else if len(snapshot.Candidates) == 0 {
		out.StopReason = "NO_CANDIDATES_IN_SCOPE"
	}
	if len(snapshot.Failures) > 0 {
		out.StopReason = "DISCOVERY_INCOMPLETE"
	}
	if snapshot.Receipt != nil && snapshot.Receipt.CompletionCause != domain.CompletionFinished {
		out.StopReason = string(snapshot.Receipt.CompletionCause)
	}
	for _, exp := range snapshot.PerExpansion {
		ids := expansionRequirements(snapshot.Idea, exp.Expansion, out.Requirements)
		if len(exp.Provider.Attempts) == 0 {
			out.Searches = append(out.Searches, domain.InvestigationSearch{RequirementIDs: ids, ExpansionID: exp.Expansion.ID, Purpose: "discovery", OriginalQuery: exp.Expansion.Query, Relaxation: "EFFECTIVE_QUERY_NOT_RECORDED", Provider: exp.Provider.Provider})
		}
		for _, attempt := range exp.Provider.Attempts {
			relaxation := "NONE_RECORDED"
			if attempt.Query == "" {
				relaxation = "EFFECTIVE_QUERY_NOT_RECORDED"
			} else if attempt.Query != exp.Expansion.Query {
				relaxation = "QUERY_CHANGED_REEVALUATE_ORIGINAL_INTENT"
			}
			out.Searches = append(out.Searches, domain.InvestigationSearch{RequirementIDs: ids, ExpansionID: exp.Expansion.ID, Purpose: "discovery", OriginalQuery: exp.Expansion.Query, EffectiveQuery: attempt.Query, Relaxation: relaxation, Provider: exp.Provider.Provider, Error: attempt.Error})
		}
	}
	sort.SliceStable(out.Searches, func(i, j int) bool {
		a, b := out.Searches[i], out.Searches[j]
		return a.ExpansionID+"\x00"+a.Provider+"\x00"+a.EffectiveQuery+"\x00"+a.Error < b.ExpansionID+"\x00"+b.Provider+"\x00"+b.EffectiveQuery+"\x00"+b.Error
	})
	byKey := map[string]domain.CandidateAssessment{}
	ambiguous := map[string]bool{}
	for _, raw := range assessments {
		a := SanitizeAssessmentEvidence(raw)
		key := a.Repository.Key()
		if existing, ok := byKey[key]; ok {
			left, _ := json.Marshal(existing)
			right, _ := json.Marshal(a)
			if string(left) != string(right) {
				ambiguous[key] = true
			}
		}
		byKey[key] = a
	}
	for key := range ambiguous {
		byKey[key] = domain.CandidateAssessment{Repository: byKey[key].Repository}
	}
	// Candidates lacking collection still require an explicit investigation action.
	for _, c := range snapshot.Candidates {
		ref := domain.CommitRef{Owner: c.Repository.Owner, Repo: c.Repository.Repo, Commit: c.Repository.CommitID}
		if _, ok := byKey[ref.Key()]; !ok {
			byKey[ref.Key()] = domain.CandidateAssessment{Repository: ref}
		}
	}
	keys := make([]string, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, r := range out.Requirements {
		if len(keys) == 0 {
			classification := "NOT_VERIFIED"
			gap := "Nenhum candidato foi avaliado para este requisito."
			nextAction := "Buscar candidatos para o requisito e registrar origem, filtros e orçamento; não concluir inexistência."
			switch investigationCoverageStatus(coverage, r.Kind, r.Text) {
			case domain.CatalogCoverageGap:
				classification = "CATALOG_GAP"
				gap = "O catálogo não possui regra declarada única e segura para este requisito."
				nextAction = "Definir uma verificação delimitada e segura antes de atribuir presença ou ausência a qualquer candidato."
			case domain.CatalogCoverageHypothesis:
				classification = "HYPOTHESIS"
				gap = "O catálogo possui apenas uma hipótese sem autoridade declarada para este requisito."
				nextAction = "Validar a hipótese e seus testemunhos antes de qualquer decisão conclusiva sobre candidatos."
			}
			out.Actions = append(out.Actions, domain.InvestigationAction{RequirementID: r.ID, Question: r.Question, Gap: gap, Classification: classification, Status: domain.EvidenceInconclusive, Evidence: []domain.EvidenceSource{}, NextAction: nextAction, Acceptance: r.Acceptance})
			continue
		}
		for _, key := range keys {
			a := byKey[key]
			e := domain.CapabilityEvidence{Status: domain.EvidenceInconclusive}
			if r.Kind == "capability" {
				for _, v := range a.Capabilities {
					if strings.EqualFold(normalizeIdeaText(v.Capability), r.Text) {
						e = v
						break
					}
				}
			}
			if r.Kind == "problem" {
				e = a.ProblemAlignment
			}
			if r.Kind == "constraint" {
				for _, v := range a.Constraints {
					if strings.EqualFold(normalizeIdeaText(v.Constraint), r.Text) {
						e = v.Evidence
						break
					}
				}
			}
			action := domain.InvestigationAction{RequirementID: r.ID, Question: r.Question, Gap: e.Explanation, Candidate: key, Classification: "NOT_VERIFIED", Status: e.Status, Evidence: append([]domain.EvidenceSource{}, e.Sources...), NextAction: "Localizar implementação, limitações e ensaio pertinente; registrar arquivos e intervalos efetivamente examinados.", Acceptance: r.Acceptance}
			if action.Status == "" {
				action.Status = domain.EvidenceInconclusive
			}
			if !EvidenceSourcesAdmissible(a.Repository, e.Sources, false) {
				action.Evidence = []domain.EvidenceSource{}
			}
			if EvidenceSourcesAdmissible(a.Repository, e.Sources, false) {
				action.Classification = "INFERENCE"
				action.NextAction = "Verificar se os indícios atendem ao requisito original e procurar contrapontos antes de recomendar reúso."
				if e.Status == domain.EvidenceProven {
					action.Classification = "OBSERVED"
					action.NextAction = "Validar interfaces, invariantes e integração no alvo; comprovação delimitada não garante substituição geral."
				}
				if e.Status == domain.EvidenceConflict {
					action.NextAction = "Examinar o conflito e suas condições; não aprovar o requisito até resolvê-lo."
				}
			}
			if e.ReasonCode == "CATALOG_GAP" {
				action.Classification = "CATALOG_GAP"
				action.Status = domain.EvidenceInconclusive
				action.Evidence = []domain.EvidenceSource{}
				action.NextAction = "Definir uma verificação delimitada e segura; a cobertura do catálogo não descreve ausência no candidato."
			}
			if e.ReasonCode == "HYPOTHESIS_ONLY" {
				action.Classification = "HYPOTHESIS"
				action.NextAction = "Validar a hipótese e seus testemunhos antes de atribuir presença, ausência ou componente ao candidato."
			}
			if r.Kind == "preference" {
				action.NextAction = "Comparar a preferência de stack com manifesto e interfaces; não excluir por preferência desconhecida."
			}
			if r.Kind == "non_goal" {
				action.NextAction = "Investigar explicitamente o comportamento excluído e sua negação; não inferir ausência por silêncio."
			}
			if action.Gap == "" && action.Status != domain.EvidenceProven {
				action.Gap = "Evidência admissível insuficiente para responder à pergunta."
			}
			if strings.Contains(e.Explanation, "nenhuma regra") || strings.Contains(e.Explanation, "sem regra") {
				action.NextAction = "Definir uma verificação delimitada para este requisito e os testemunhos necessários; o catálogo atual não possui detector aplicável."
			}
			if strings.Contains(e.Explanation, "coleta incompleta") {
				action.NextAction = "Completar a coleta indicada nas falhas do assessment antes de avaliar a capacidade; não inferir ausência."
			}
			out.Actions = append(out.Actions, action)
		}
	}
	return publicInvestigationPlan(out)
}

func investigationCoverageStatus(coverage domain.EvidenceCatalogCoverage, kind, requirement string) domain.CatalogCoverageStatus {
	normalized := normalizeEvidenceText(requirement)
	for _, entry := range coverage.Entries {
		if entry.RequirementKind == kind && normalizeEvidenceText(entry.Requirement) == normalized {
			return entry.Status
		}
	}
	return ""
}

func expansionRequirements(idea domain.IdeaSpec, exp domain.SearchExpansion, requirements []domain.InvestigationRequirement) []string {
	ids := []string{}
	valid := map[string]bool{}
	for _, r := range requirements {
		valid[r.ID] = true
	}
	for _, cover := range exp.Covers {
		for _, r := range requirements {
			if strings.EqualFold(normalizeIdeaText(cover), normalizeIdeaText(r.Text)) {
				ids = appendUnique(ids, r.ID)
			}
		}
		kind, text := "", ""
		var i int
		switch {
		case cover == "problem":
			kind, text = "problem", idea.Problem
		case strings.HasPrefix(cover, "capabilities["):
			if _, err := fmt.Sscanf(cover, "capabilities[%d]", &i); err == nil && i >= 0 && i < len(idea.Capabilities) {
				kind, text = "capability", idea.Capabilities[i]
			}
		case strings.HasPrefix(cover, "constraints["):
			if _, err := fmt.Sscanf(cover, "constraints[%d]", &i); err == nil && i >= 0 && i < len(idea.Constraints) {
				kind, text = "constraint", idea.Constraints[i]
			}
		}
		id := requirementID(kind, text)
		if kind != "" && valid[id] {
			ids = appendUnique(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

// Keep public projection at creation so the existing certified gateway stays unchanged.
func publicInvestigationPlan(plan domain.InvestigationPlan) domain.InvestigationPlan {
	for i := range plan.Requirements {
		r := &plan.Requirements[i]
		r.Text = publicText(r.Text, "requirement")
		r.Question = publicText(r.Question, "question")
		r.Acceptance = publicText(r.Acceptance, "acceptance")
	}
	for i := range plan.Searches {
		s := &plan.Searches[i]
		s.OriginalQuery = publicText(s.OriginalQuery, "query")
		s.EffectiveQuery = publicText(s.EffectiveQuery, "query")
		if s.Error != "" {
			s.Error = publicText(s.Error, "search failure")
		}
	}
	for i := range plan.Actions {
		a := &plan.Actions[i]
		a.Question = publicText(a.Question, "question")
		a.Gap = publicText(a.Gap, "evidence gap")
		a.Evidence = projectEvidence(a.Evidence)
		a.NextAction = publicText(a.NextAction, "next action")
		a.Acceptance = publicText(a.Acceptance, "acceptance")
	}
	return plan
}
