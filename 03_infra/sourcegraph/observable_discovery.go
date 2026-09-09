/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/observable_discovery.md
 * @layer L3
 * @updated 2026-09-08
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
