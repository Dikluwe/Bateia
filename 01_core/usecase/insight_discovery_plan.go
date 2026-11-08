/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_discovery_plan.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"strings"
	"time"
)

func PlanInsightDiscovery(request domain.InsightDiscoveryRequest) (domain.InsightDiscoveryPlan, error) {
	if request.SchemaVersion != "1" || (request.QueryPolicy != "" && request.QueryPolicy != "lexical-v1" && request.QueryPolicy != "context-v2") {
		return domain.InsightDiscoveryPlan{}, domain.ErrInvalidRunInput
	}
	base, err := normalizeInsightExploration(domain.InsightExplorationRequest{SchemaVersion: "1", AlgorithmIdeas: request.AlgorithmIdeas, Question: request.Question, FamiliarApproaches: request.FamiliarApproaches, MaxPassages: request.MaxPassages})
	if err != nil {
		return domain.InsightDiscoveryPlan{}, err
	}
	request.Question, request.FamiliarApproaches, request.MaxPassages = base.Question, base.FamiliarApproaches, base.MaxPassages
	if request.MaxPerQuery == 0 {
		request.MaxPerQuery = 2
		if request.QueryPolicy == "context-v2" {
			request.MaxPerQuery = 5
		}
	}
	if request.MaxReferences == 0 {
		request.MaxReferences = 4
	}
	if request.TotalTimeout == 0 {
		request.TotalTimeout = 90 * time.Second
	}
	if request.AttemptTimeout == 0 {
		request.AttemptTimeout = 10 * time.Second
	}
	if request.MaxPerQuery < 1 || request.MaxPerQuery > 5 || request.MaxReferences < 1 || request.MaxReferences > 8 || request.TotalTimeout < time.Second || request.TotalTimeout > 5*time.Minute || request.AttemptTimeout < time.Second || request.AttemptTimeout > 30*time.Second || request.AttemptTimeout > request.TotalTimeout {
		return domain.InsightDiscoveryPlan{}, domain.ErrInvalidRunInput
	}
	terms := explorationTerms(request.Question)
	if len(terms) == 0 {
		return domain.InsightDiscoveryPlan{}, domain.ErrInvalidRunInput
	}
	familiar := map[string]bool{}
	for _, label := range request.FamiliarApproaches {
		for _, token := range explorationWords(label) {
			familiar[token] = true
		}
	}
	anchors := []string{}
	for _, token := range terms {
		if !familiar[token] {
			anchors = append(anchors, token)
		}
	}
	if len(anchors) == 0 {
		anchors = append(anchors, terms...)
	}
	inputHash := insightHash(request)
	plan := domain.InsightDiscoveryPlan{SchemaVersion: "insight-discovery-plan-v1", ID: "discovery_plan_" + insightHash(struct{ Contract, Input string }{"insight-discovery-v1", inputHash}), RequestSHA256: inputHash, Request: request, Queries: []domain.InsightDiscoveryQuery{}}
	if request.QueryPolicy == "context-v2" {
		pairs := [][]string{terms[:min(3, len(terms))], anchors[:min(2, len(anchors))], {anchors[0], anchors[len(anchors)-1]}, {anchors[min(1, len(anchors)-1)], anchors[min(2, len(anchors)-1)]}, {anchors[len(anchors)/2], anchors[len(anchors)-1]}}
		for i, perspective := range []string{"DIRECT", "BROADENED", "ALTERNATIVES", "TRADEOFFS", "COUNTEREXAMPLES"} {
			seen := map[string]bool{}
			tokens := []string{}
			for _, term := range pairs[i] {
				if !seen[term] {
					seen[term] = true
					tokens = append(tokens, term)
				}
			}
			q := domain.InsightDiscoveryQuery{Perspective: perspective, Query: strings.Join(tokens, " "), AnchorTerms: append([]string{}, tokens...), Derivation: "CONTEXTUAL_TERM_CONJUNCTION_V2"}
			q.ID = "discovery_query_" + insightHash(struct{ Plan, Perspective, Query string }{plan.ID, q.Perspective, q.Query})
			plan.Queries = append(plan.Queries, q)
		}
		return plan, nil
	}
	for _, route := range []struct {
		name, query string
		anchors     []string
	}{
		{"DIRECT", strings.Join(terms[:min(3, len(terms))], " "), terms[:min(3, len(terms))]},
		{"BROADENED", anchors[0], anchors[:1]},
		{"ALTERNATIVES", anchors[0] + " alternative", anchors[:1]},
		{"TRADEOFFS", anchors[len(anchors)-1] + " tradeoffs", anchors[len(anchors)-1:]},
		{"COUNTEREXAMPLES", anchors[0] + " failure", anchors[:1]},
	} {
		query := domain.InsightDiscoveryQuery{Perspective: route.name, Query: route.query, AnchorTerms: append([]string{}, route.anchors...), Derivation: "LEXICAL_EXPLORATION_HEURISTIC"}
		query.ID = "discovery_query_" + insightHash(struct{ Plan, Perspective, Query string }{plan.ID, route.name, route.query})
		plan.Queries = append(plan.Queries, query)
	}
	return plan, nil
}
