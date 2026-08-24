/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/observable_discovery.md
 * @layer L3
 * @updated 2026-08-24
 */

package sourcegraph

import (
	"context"
	"fmt"
	"strings"

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

func observableQueries(query, language string) []string {
	short := shortQuery(query)
	if language == "" {
		return []string{query, short}
	}
	lang := sourcegraphLanguage(language)
	return []string{fmt.Sprintf("lang:%s %s", lang, query), fmt.Sprintf("lang:%s %s", lang, short), short}
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
	for attemptIndex, query := range observableQueries(request.Query, request.Language) {
		candidates, matches, err := c.executeSGQueryDetailed(ctx, query, request.MaxResults)
		attempt := domain.ProviderAttempt{Query: query, ResultCount: len(candidates), Error: c.safeError(err)}
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
