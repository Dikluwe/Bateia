/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/observable_discovery.md
 * @layer L3
 * @updated 2026-08-24
 */

package sourcegraph

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"bateia/01_core/domain"
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

func observableQueries(query, language string, discriminators []string, repositoryScope string) []string {
	short := shortQuery(query)
	if len(discriminators) > 0 {
		short = strings.TrimSpace(discriminators[0])
	}
	prefix := ""
	if repositoryScope != "" {
		prefix = "repo:^" + strings.TrimSpace(repositoryScope) + "$ "
	}
	if language == "" {
		return []string{prefix + query, prefix + short}
	}
	lang := sourcegraphLanguage(language)
	return []string{prefix + fmt.Sprintf("lang:%s %s", lang, query), prefix + fmt.Sprintf("lang:%s %s", lang, short), prefix + short}
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

// Discover executa e torna observáveis as tentativas e fallbacks do Sourcegraph.
func (c *Client) Discover(ctx context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	result := domain.ProviderDiscovery{Provider: "sourcegraph", Attempts: make([]domain.ProviderAttempt, 0, 3)}
	var lastErr error
	requestSucceeded := false
	maxAttempts := request.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	for attemptIndex, query := range observableQueries(request.Query, request.Language, request.Discriminators, request.RepositoryScope) {
		if attemptIndex >= maxAttempts || ctx.Err() != nil {
			break
		}
		attemptCtx := ctx
		cancel := func() {}
		if request.AttemptTimeout > 0 {
			attemptCtx, cancel = context.WithTimeout(ctx, request.AttemptTimeout)
		}
		started := time.Now().UTC()
		candidates, matches, err := c.executeSGQueryDetailed(attemptCtx, query, request.MaxResults)
		attemptErr := attemptCtx.Err()
		cancel()
		status, class := "SUCCEEDED", ""
		if err == nil && len(candidates) == 0 {
			status = "EMPTY"
		}
		if err != nil {
			status, class = "FAILED", "PROVIDER"
			if errors.Is(attemptErr, context.DeadlineExceeded) {
				status, class = "TIMED_OUT", "TIMEOUT"
			}
			if errors.Is(attemptErr, context.Canceled) || errors.Is(ctx.Err(), context.Canceled) {
				status, class = "CANCELLED", "CANCELLED"
			}
		}
		attempt := domain.ProviderAttempt{Query: query, AttemptIndex: attemptIndex, StartedAt: started, Duration: time.Since(started), Status: status, ResultCount: len(candidates), ErrorClass: class, Error: c.safeError(err)}
		result.Attempts = append(result.Attempts, attempt)
		if err != nil {
			lastErr = err
			continue
		}
		requestSucceeded = true
		if len(candidates) > 0 {
			result.Candidates = candidates
			for index := range matches {
				matches[index].AttemptIndex = attemptIndex
				matches[index].EffectiveQuery = query
			}
			result.Matches = matches
			return result, nil
		}
	}
	if !requestSucceeded && lastErr != nil {
		return result, fmt.Errorf("todas as consultas ao Sourcegraph falharam: %s", c.safeError(lastErr))
	}
	return result, nil
}
