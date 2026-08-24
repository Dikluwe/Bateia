/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_discovery.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"bateia/01_core/domain"
)

type discoveryOutcome struct {
	discovery domain.ExpansionDiscovery
	failure   *domain.ExecutionFailure
}

func repositoryKey(candidate domain.RepositoryCandidate) string {
	return strings.ToLower(strings.TrimSpace(candidate.Owner) + "/" + strings.TrimSpace(candidate.Repo))
}

// DiscoverExpansionPlan executa as expansões com concorrência limitada e saída canônica.
func DiscoverExpansionPlan(ctx context.Context, plan domain.ExpansionPlan, strategy, language string, maxResults, concurrency int, provider domain.ObservableDiscoveryProvider) (domain.DiscoverySnapshot, error) {
	if provider == nil {
		return domain.DiscoverySnapshot{}, fmt.Errorf("provedor de descoberta é obrigatório")
	}
	if len(plan.Expansions) == 0 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("plano não possui expansões")
	}
	if maxResults < 1 || maxResults > 25 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("maxResults deve estar entre 1 e 25")
	}
	if concurrency < 1 || concurrency > 4 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("concorrência deve estar entre 1 e 4")
	}

	outcomes := make([]discoveryOutcome, len(plan.Expansions))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	started := 0
discoveryLoop:
	for index, expansion := range plan.Expansions {
		if strings.TrimSpace(expansion.ID) == "" || strings.TrimSpace(expansion.Query) == "" {
			return domain.DiscoverySnapshot{}, fmt.Errorf("expansão %d possui id ou query vazia", index)
		}
		select {
		case <-ctx.Done():
			break discoveryLoop
		case sem <- struct{}{}:
			if ctx.Err() != nil {
				<-sem
				break discoveryLoop
			}
			started++
			wg.Add(1)
			go func(index int, expansion domain.SearchExpansion) {
				defer wg.Done()
				defer func() { <-sem }()
				request := domain.DiscoveryRequest{ExpansionID: expansion.ID, Query: expansion.Query, Language: language, MaxResults: maxResults}
				result, err := provider.Discover(ctx, request)
				outcomes[index].discovery = domain.ExpansionDiscovery{Expansion: expansion, Provider: result}
				if err != nil {
					outcomes[index].failure = &domain.ExecutionFailure{Scope: expansion.ID, Stage: "discovery", Message: err.Error()}
				}
			}(index, expansion)
		}
	}
	wg.Wait()
	if started == 0 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("descoberta cancelada antes de iniciar: %w", ctx.Err())
	}

	snapshot := domain.DiscoverySnapshot{
		SchemaVersion: "2", Idea: normalizeIdea(plan.Idea), Strategy: strings.ToLower(strings.TrimSpace(strategy)), Language: language,
		PerExpansion: make([]domain.ExpansionDiscovery, 0, started), Candidates: make([]domain.DiscoveredCandidate, 0),
	}
	candidateIndexes := make(map[string]int)
	failures := 0
	for index := 0; index < len(outcomes); index++ {
		outcome := outcomes[index]
		if outcome.discovery.Expansion.ID == "" {
			continue
		}
		snapshot.PerExpansion = append(snapshot.PerExpansion, outcome.discovery)
		if outcome.failure != nil {
			failures++
			snapshot.Failures = append(snapshot.Failures, *outcome.failure)
		}
		for _, candidate := range outcome.discovery.Provider.Candidates {
			key := repositoryKey(candidate)
			if key == "/" {
				continue
			}
			origin := candidateOrigin(outcome.discovery, candidate)
			if existingIndex, exists := candidateIndexes[key]; exists {
				existing := &snapshot.Candidates[existingIndex]
				existing.ExpansionIDs = appendUnique(existing.ExpansionIDs, outcome.discovery.Expansion.ID)
				existing.Queries = appendUnique(existing.Queries, outcome.discovery.Expansion.Query)
				existing.Origins = append(existing.Origins, origin)
				continue
			}
			candidateIndexes[key] = len(snapshot.Candidates)
			snapshot.Candidates = append(snapshot.Candidates, domain.DiscoveredCandidate{
				Repository: candidate, ExpansionIDs: []string{outcome.discovery.Expansion.ID}, Queries: []string{outcome.discovery.Expansion.Query}, Origins: []domain.CandidateOrigin{origin},
			})
		}
	}
	if failures == started {
		message := "erro não informado"
		if len(snapshot.Failures) > 0 {
			message = snapshot.Failures[0].Message
		}
		return domain.DiscoverySnapshot{}, fmt.Errorf("todas as %d expansões falharam: %s", failures, message)
	}
	return snapshot, nil
}

func candidateOrigin(discovery domain.ExpansionDiscovery, candidate domain.RepositoryCandidate) domain.CandidateOrigin {
	origin := domain.CandidateOrigin{ExpansionID: discovery.Expansion.ID, PlannedQuery: discovery.Expansion.Query, AttemptIndex: -1}
	key := repositoryKey(candidate)
	for _, match := range discovery.Provider.Matches {
		if strings.ToLower(strings.TrimSpace(match.Repository)) == key {
			origin.AttemptIndex = match.AttemptIndex
			origin.EffectiveQuery = match.EffectiveQuery
			origin.Evidence = append([]domain.MatchEvidence(nil), match.Evidence...)
			return origin
		}
	}
	for index, attempt := range discovery.Provider.Attempts {
		if attempt.ResultCount > 0 {
			origin.AttemptIndex = index
			origin.EffectiveQuery = attempt.Query
			break
		}
	}
	return origin
}
