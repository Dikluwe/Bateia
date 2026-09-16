/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/observable_discovery.md
 * @layer L3
 * @updated 2026-09-11
 */

package sourcegraph

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

func sourcegraphLanguage(language string) string {
	if mapped, ok := sgLangName[strings.ToLower(language)]; ok {
		return mapped
	}
	if language == "" {
		return ""
	}
	return strings.ToUpper(language[:1]) + language[1:]
}

// Discriminators remain compatible metadata; they never replace the user's query.
func observableQueries(query, language string, discriminators []string, repositoryScope string) []string {
	prefix := ""
	if scope := strings.TrimSpace(repositoryScope); scope != "" {
		prefix = "repo:^" + regexp.QuoteMeta(scope) + "$ "
	}
	if language = strings.TrimSpace(language); language != "" {
		prefix += "lang:" + sourcegraphLanguage(language) + " "
	}
	return []string{prefix + query}
}

func (c *Client) safeError(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	if c.apiToken != "" {
		message = strings.ReplaceAll(message, c.apiToken, "[redacted]")
	}
	return message
}

// redactedProviderError preserves error identity while keeping credentials out of messages.
type redactedProviderError struct {
	message string
	cause   error
}

func (e redactedProviderError) Error() string { return e.message }
func (e redactedProviderError) Unwrap() error { return e.cause }

func (c *Client) redactedError(err error) error {
	if err == nil {
		return nil
	}
	return redactedProviderError{message: c.safeError(err), cause: err}
}

// discoveryDerivedQuery interpolates only terms supplied by the finite catalog.
func discoveryDerivedQuery(plan domain.DiscoveryQueryPlan, language, scope string) string {
	expression := `patterntype:regexp (` + strings.Join(plan.ActionTerms, "|") + `) AND (` + strings.Join(plan.ObjectTerms, "|") + `)`
	return observableQueries(expression, language, nil, scope)[0]
}

func normalizeTechnicalValue(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}

func technicalDiscoveryPlan(request domain.DiscoveryRequest) (domain.DiscoveryQueryPlan, error) {
	query := strings.TrimSpace(request.Query)
	if query == "" || len(query) > 256 {
		return domain.DiscoveryQueryPlan{}, fmt.Errorf("consulta técnica deve conter de 1 a 256 bytes")
	}
	if request.Derivation == nil || len(request.Derivation.RuleIDs) != 1 || strings.TrimSpace(request.Derivation.RuleIDs[0]) == "" || len(request.Derivation.SourceIDs) != 1 || strings.TrimSpace(request.Derivation.SourceIDs[0]) == "" {
		return domain.DiscoveryQueryPlan{}, fmt.Errorf("consulta técnica exige exatamente uma regra e uma origem não vazias")
	}
	if len(request.Discriminators) != 2 {
		return domain.DiscoveryQueryPlan{}, fmt.Errorf("consulta técnica exige exatamente dois discriminadores")
	}
	discriminators := make([]string, 2)
	for index, value := range request.Discriminators {
		discriminators[index] = strings.TrimSpace(value)
		if discriminators[index] == "" || len(discriminators[index]) > 64 {
			return domain.DiscoveryQueryPlan{}, fmt.Errorf("discriminador técnico deve conter de 1 a 64 bytes")
		}
	}
	if normalizeTechnicalValue(discriminators[0]) == normalizeTechnicalValue(discriminators[1]) {
		return domain.DiscoveryQueryPlan{}, fmt.Errorf("discriminadores técnicos devem ser distintos")
	}
	return domain.DiscoveryQueryPlan{
		Original: query, Version: "0068-technical-v1", Status: "hypothesis", Reason: "technical_discriminators",
		ActionTerms: []string{discriminators[0]}, ObjectTerms: []string{discriminators[1]},
		RuleIDs: append([]string(nil), request.Derivation.RuleIDs...), Unverified: []string{query},
	}, nil
}

func knownLegacyDerivation(strategy string) bool {
	return strategy == "lexical" || strategy == "controlled-catalog-exact-v1" || (strings.HasPrefix(strategy, "ontology:") && len(strategy) > len("ontology:"))
}

func discoverySearchAttempts(original string, attempts []domain.ProviderAttempt, matches []domain.ProviderCandidateMatch) []domain.SearchAttempt {
	if len(attempts) == 0 {
		return nil
	}
	traces := make([]domain.SearchAttempt, 0, len(attempts))
	for _, attempt := range attempts {
		trace := domain.SearchAttempt{OriginalQuery: original, Provider: "Sourcegraph", Query: attempt.Query, Candidates: attempt.ResultCount, Error: attempt.Error, QueryPlan: attempt.QueryPlan, RawCandidates: attempt.RawCandidates, RejectedCandidates: attempt.RejectedCandidates, Coverage: attempt.Coverage, CandidateDiagnostics: attempt.CandidateDiagnostics}
		if attempt.QueryPlan != nil {
			for _, match := range matches {
				if match.AttemptIndex == attempt.AttemptIndex {
					trace.AdmissionEvidence = append(trace.AdmissionEvidence, match)
				}
			}
			trace.Relaxation = "hipótese lexical de ação e objeto; implementação e restrições da consulta original não verificadas"
		}
		traces = append(traces, trace)
	}
	return traces
}

// Discover tries the intact query, then at most one lexical hypothesis after valid emptiness.
func (c *Client) Discover(ctx context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	result := domain.ProviderDiscovery{Provider: "sourcegraph", Attempts: make([]domain.ProviderAttempt, 0, 2)}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if request.Derivation != nil && !knownLegacyDerivation(request.Derivation.Strategy) && request.Derivation.Strategy != "technical-lexicon-v1" {
		return result, fmt.Errorf("estratégia de derivação Sourcegraph desconhecida %q", request.Derivation.Strategy)
	}
	if request.Derivation != nil && request.Derivation.Strategy == "technical-lexicon-v1" {
		plan, err := technicalDiscoveryPlan(request)
		if err != nil {
			return result, err
		}
		query := observableQueries(plan.Original, request.Language, nil, request.RepositoryScope)[0] + fmt.Sprintf(" count:%d", discoveryRequestedMatches(request.MaxResults))
		attemptCtx := ctx
		cancel := func() {}
		if request.AttemptTimeout > 0 {
			attemptCtx, cancel = context.WithTimeout(ctx, request.AttemptTimeout)
		}
		started := time.Now().UTC()
		attempt := domain.ProviderAttempt{Query: query, AttemptIndex: 0, StartedAt: started, QueryPlan: &plan}
		candidates, matches, coverage, diagnostics, discoveryErr := c.executeDerivedDiscovery(attemptCtx, query, request.MaxResults, plan)
		attempt.Coverage, attempt.CandidateDiagnostics = coverage, diagnostics
		if coverage != nil {
			attempt.RawCandidates = coverage.ExaminedRepositories
		}
		for _, diagnostic := range diagnostics {
			if diagnostic.Status == "insufficient_evidence" {
				attempt.RejectedCandidates++
			}
		}
		if discoveryErr == nil {
			discoveryErr = attemptCtx.Err()
		}
		attempt.Status = "SUCCEEDED"
		if discoveryErr == nil && len(candidates) == 0 {
			attempt.Status = "EMPTY"
		}
		if discoveryErr != nil {
			attempt.Status, attempt.ErrorClass = "FAILED", "PROVIDER"
			if errors.Is(discoveryErr, context.DeadlineExceeded) || errors.Is(attemptCtx.Err(), context.DeadlineExceeded) {
				attempt.Status, attempt.ErrorClass = "TIMED_OUT", "TIMEOUT"
			}
			if errors.Is(discoveryErr, context.Canceled) || errors.Is(attemptCtx.Err(), context.Canceled) {
				attempt.Status, attempt.ErrorClass = "CANCELLED", "CANCELLED"
			}
		}
		attempt.ResultCount, attempt.Duration, attempt.Error = len(candidates), time.Since(started), c.safeError(discoveryErr)
		cancel()
		result.Attempts = append(result.Attempts, attempt)
		if discoveryErr != nil {
			return result, c.redactedError(discoveryErr)
		}
		result.Candidates, result.Matches = candidates, matches
		for index := range result.Matches {
			result.Matches[index].AttemptIndex, result.Matches[index].EffectiveQuery = 0, query
		}
		for index := range result.Candidates {
			result.Candidates[index].SearchAttempts = discoverySearchAttempts(request.Query, result.Attempts, result.Matches)
		}
		return result, nil
	}
	if strings.TrimSpace(request.Query) == "" {
		return result, fmt.Errorf("consulta Sourcegraph vazia")
	}
	plan := parser.PlanDiscoveryQuery(request.Query)
	queries := observableQueries(request.Query, request.Language, request.Discriminators, request.RepositoryScope)
	if plan.Status == "hypothesis" && (request.MaxAttempts <= 0 || request.MaxAttempts > 1) {
		queries = append(queries, discoveryDerivedQuery(plan, request.Language, request.RepositoryScope)+fmt.Sprintf(" count:%d", discoveryRequestedMatches(request.MaxResults)))
	}
	for index, query := range queries {
		if err := ctx.Err(); err != nil {
			return result, c.redactedError(err)
		}
		attemptCtx := ctx
		cancel := func() {}
		if request.AttemptTimeout > 0 {
			attemptCtx, cancel = context.WithTimeout(ctx, request.AttemptTimeout)
		}
		started := time.Now().UTC()
		var candidates []domain.RepositoryCandidate
		var matches []domain.ProviderCandidateMatch
		var err error
		attempt := domain.ProviderAttempt{Query: query, AttemptIndex: index, StartedAt: started}
		if index > 0 {
			attempt.QueryPlan = &plan
			candidates, matches, attempt.Coverage, attempt.CandidateDiagnostics, err = c.executeDerivedDiscovery(attemptCtx, query, request.MaxResults, plan)
			if attempt.Coverage != nil {
				attempt.RawCandidates = attempt.Coverage.ExaminedRepositories
			}
			for _, diagnostic := range attempt.CandidateDiagnostics {
				if diagnostic.Status == "insufficient_evidence" {
					attempt.RejectedCandidates++
				}
			}
		} else {
			candidates, matches, err = c.executeSGQueryDetailed(attemptCtx, query, request.MaxResults)
		}

		if err == nil {
			err = attemptCtx.Err()
		}
		attempt.Status = "SUCCEEDED"
		if err == nil && len(candidates) == 0 {
			attempt.Status = "EMPTY"
		}
		if err != nil {
			attempt.Status, attempt.ErrorClass = "FAILED", "PROVIDER"
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(attemptCtx.Err(), context.DeadlineExceeded) {
				attempt.Status, attempt.ErrorClass = "TIMED_OUT", "TIMEOUT"
			}
			if errors.Is(err, context.Canceled) || errors.Is(attemptCtx.Err(), context.Canceled) {
				attempt.Status, attempt.ErrorClass = "CANCELLED", "CANCELLED"
			}
		}
		attempt.ResultCount, attempt.Duration, attempt.Error = len(candidates), time.Since(started), c.safeError(err)
		cancel()
		result.Attempts = append(result.Attempts, attempt)
		if err != nil {
			return result, c.redactedError(err)
		}
		result.Candidates, result.Matches = candidates, matches
		for i := range result.Matches {
			result.Matches[i].AttemptIndex, result.Matches[i].EffectiveQuery = index, query
		}
		for i := range result.Candidates {
			result.Candidates[i].SearchAttempts = discoverySearchAttempts(request.Query, result.Attempts, result.Matches)
		}
		if len(candidates) > 0 {
			return result, nil
		}
	}
	return result, nil
}
