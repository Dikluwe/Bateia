/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/project_catalog.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

func DiscoverProjectCatalog(ctx context.Context, plan domain.ExpansionPlan, maxResults int, policy domain.DiscoveryExecutionPolicy, provider domain.ProjectCatalog) (domain.CatalogSnapshot, error) {
	out := domain.CatalogSnapshot{SchemaVersion: "1", Idea: plan.Idea, Plan: plan, PerExpansion: []domain.CatalogDiscovery{}, Candidates: []domain.CatalogCandidate{}}
	if provider == nil {
		return out, fmt.Errorf("catálogo ausente")
	}
	out.Provider = provider.Capabilities()
	if !domain.SupportsRole(out.Provider, domain.RoleProjectCatalog) {
		return out, fmt.Errorf("provedor %s não declara PROJECT_CATALOG", out.Provider.Provider)
	}
	if maxResults < 1 || maxResults > 25 {
		return out, fmt.Errorf("maxResults deve estar entre 1 e 25")
	}
	byRepo := map[string]int{}
	started, completed, timedOut, failed := 0, 0, 0, 0
	for _, expansion := range plan.Expansions {
		if ctx.Err() != nil {
			break
		}
		if len(expansion.Discriminators) == 0 || !containsAnyPhrase(expansion.Query, expansion.Discriminators) || isGenericOntologyQuery(expansion.Query) {
			return out, fmt.Errorf("query %s não está auditada para catálogo", expansion.ID)
		}
		started++
		discovery, err := provider.SearchCatalog(ctx, domain.CatalogSearchRequest{ExpansionID: expansion.ID, Query: expansion.Query, Discriminators: append([]string(nil), expansion.Discriminators...), MaxResults: maxResults, AttemptTimeout: policy.AttemptTimeout})
		completed++
		if discovery.Role != domain.RoleProjectCatalog {
			return out, fmt.Errorf("provedor retornou papel incompatível %s", discovery.Role)
		}
		if discovery.Attempt.Status == "TIMED_OUT" {
			timedOut++
		}
		if err != nil {
			failed++
			discovery.Attempt.SafeMessage = err.Error()
		}
		out.PerExpansion = append(out.PerExpansion, discovery)
		for _, candidate := range discovery.Candidates {
			key := strings.ToLower(candidate.Repository.Owner + "/" + candidate.Repository.Repo)
			origin := domain.CandidateOrigin{ExpansionID: expansion.ID, PlannedQuery: expansion.Query, EffectiveQuery: discovery.Attempt.Query}
			if index, ok := byRepo[key]; ok {
				out.Candidates[index].MatchedFields = uniqueSorted(append(out.Candidates[index].MatchedFields, candidate.MatchedFields...))
				out.Candidates[index].Origins = append(out.Candidates[index].Origins, origin)
				continue
			}
			candidate.Origins = []domain.CandidateOrigin{origin}
			byRepo[key] = len(out.Candidates)
			out.Candidates = append(out.Candidates, candidate)
		}
	}
	status, cause := domain.DiscoveryRunComplete, domain.CompletionFinished
	if completed < len(plan.Expansions) {
		status, cause = domain.DiscoveryRunPartial, domain.CompletionCancelled
	}
	out.Receipt = &domain.DiscoveryRunReceipt{SchemaVersion: "1", Status: status, CompletionCause: cause, Policy: policy, Planned: len(plan.Expansions), Started: started, Completed: completed, TimedOut: timedOut, Failed: failed}
	sort.Slice(out.Candidates, func(i, j int) bool {
		return strings.ToLower(out.Candidates[i].Repository.Owner+"/"+out.Candidates[i].Repository.Repo) < strings.ToLower(out.Candidates[j].Repository.Owner+"/"+out.Candidates[j].Repository.Repo)
	})
	if failed == completed && completed > 0 {
		return out, fmt.Errorf("todas as consultas de catálogo falharam")
	}
	return out, nil
}

func ScreenProjectCatalog(snapshot domain.CatalogSnapshot) domain.CatalogScreeningReport {
	out := domain.CatalogScreeningReport{SchemaVersion: "1", Decisions: []domain.CatalogDecision{}}
	for _, candidate := range snapshot.Candidates {
		signals := []domain.CatalogSignal{}
		score := 0.0
		add := func(id string, value float64, reason string) {
			score += value
			signals = append(signals, domain.CatalogSignal{ID: id, Contribution: value, Reason: reason})
		}
		if candidate.Repository.CommitID == "" || strings.EqualFold(candidate.Repository.CommitID, "HEAD") {
			add("IDENTITY_UNKNOWN", -.5, "revisão imutável ausente")
		} else {
			add("IDENTITY_RESOLVED", .35, "revisão imutável resolvida")
		}
		if candidate.Repository.Archived {
			add("ARCHIVED", -.5, "repositório arquivado")
		}
		strong := false
		for _, field := range candidate.MatchedFields {
			if field == "name" || field == "description" || field == "topics" {
				strong = true
			}
		}
		if strong {
			add("CATALOG_FIELD_MATCH", .35, "match em campo de catálogo")
		}
		if len(candidate.Origins) >= 2 {
			add("MULTI_QUERY_SUPPORT", .3, "duas ou mais origens independentes")
		}
		verdict := domain.CatalogReview
		if score >= .5 && strong && candidate.Repository.CommitID != "" && !candidate.Repository.Archived {
			verdict = domain.CatalogAdvance
		}
		if score < 0 || candidate.Repository.Archived {
			verdict = domain.CatalogReject
		}
		out.Decisions = append(out.Decisions, domain.CatalogDecision{Repository: candidate.Repository, Verdict: verdict, Score: score, Signals: signals, Origins: candidate.Origins})
	}
	sort.Slice(out.Decisions, func(i, j int) bool {
		if out.Decisions[i].Score != out.Decisions[j].Score {
			return out.Decisions[i].Score > out.Decisions[j].Score
		}
		return out.Decisions[i].Repository.Owner+"/"+out.Decisions[i].Repository.Repo < out.Decisions[j].Repository.Owner+"/"+out.Decisions[j].Repository.Repo
	})
	return out
}

func BuildCatalogProofPlan(report domain.CatalogScreeningReport, source domain.ExpansionPlan) domain.ExpansionPlan {
	byID := map[string]domain.SearchExpansion{}
	for _, expansion := range source.Expansions {
		byID[expansion.ID] = expansion
	}
	out := domain.ExpansionPlan{Idea: source.Idea, Expansions: []domain.SearchExpansion{}}
	for _, decision := range report.Decisions {
		if decision.Verdict != domain.CatalogAdvance {
			continue
		}
		scope := decision.Repository.Owner + "/" + decision.Repository.Repo
		for _, origin := range decision.Origins {
			expansion, ok := byID[origin.ExpansionID]
			if !ok {
				continue
			}
			expansion.ID = semanticExpansionID("proof:"+scope, expansion.Facet, expansion.Query)
			expansion.RepositoryScope = scope
			expansion.Rationale = "comprovar capacidade no candidato de catálogo"
			out.Expansions = append(out.Expansions, expansion)
		}
	}
	sort.Slice(out.Expansions, func(i, j int) bool { return out.Expansions[i].ID < out.Expansions[j].ID })
	return out
}

func containsAnyPhrase(query string, phrases []string) bool {
	n := normalizeOntologyText(query)
	for _, p := range phrases {
		if ontologyPhraseMatch(n, normalizeOntologyText(p)) {
			return true
		}
	}
	return false
}
