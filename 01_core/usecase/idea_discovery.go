/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_discovery.md
 * @layer L1
 * @updated 2026-11-08
 */

package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

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
	policy := domain.DiscoveryExecutionPolicy{SchemaVersion: "1", Concurrency: concurrency, MaxAttempts: 3}
	return DiscoverExpansionPlanWithPolicy(ctx, plan, strategy, language, maxResults, policy, provider, nil)
}

func DiscoverExpansionPlanWithPolicy(ctx context.Context, plan domain.ExpansionPlan, strategy, language string, maxResults int, policy domain.DiscoveryExecutionPolicy, provider domain.ObservableDiscoveryProvider, reporter domain.DiscoveryProgressReporter) (domain.DiscoverySnapshot, error) {
	if provider == nil {
		return domain.DiscoverySnapshot{}, fmt.Errorf("provedor de descoberta é obrigatório")
	}
	if len(plan.Expansions) == 0 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("plano não possui expansões")
	}
	if maxResults < 1 || maxResults > 25 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("maxResults deve estar entre 1 e 25")
	}
	if policy.SchemaVersion != "1" || policy.Concurrency < 1 || policy.Concurrency > 4 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("concorrência deve estar entre 1 e 4")
	}
	if policy.TotalTimeout < 0 || policy.AttemptTimeout < 0 || (policy.TotalTimeout > 0 && policy.AttemptTimeout > policy.TotalTimeout) {
		return domain.DiscoverySnapshot{}, fmt.Errorf("timeouts de descoberta inválidos")
	}
	if policy.MaxAttempts < 1 || policy.MaxAttempts > 3 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("maxAttempts deve estar entre 1 e 3")
	}
	if policy.TotalTimeout > 0 {
		var cancel context.CancelFunc
		var deadlineErr error
		ctx, cancel, deadlineErr = domain.ExecutionTimeout(ctx, policy.TotalTimeout)
		if deadlineErr != nil {
			return domain.DiscoverySnapshot{}, deadlineErr
		}
		defer cancel()
	}

	plan = clonePlan(plan)
	requirements := investigationRequirements(plan.Idea)
	for i := range plan.Expansions {
		plan.Expansions[i].RequirementIDs = expansionRequirements(plan.Idea, plan.Expansions[i], requirements)
	}
	outcomes := make([]discoveryOutcome, len(plan.Expansions))
	for index, expansion := range plan.Expansions {
		if strings.TrimSpace(expansion.ID) == "" || strings.TrimSpace(expansion.Query) == "" {
			return domain.DiscoverySnapshot{}, fmt.Errorf("expansão %d possui id ou query vazia", index)
		}
	}
	// A porta não exige que o provider seja thread-safe. Mantemos a fronteira
	// serial para que cancelamento observado impeça o início da próxima chamada;
	// o paralelismo limitado permanece responsabilidade do adapter concreto.
	sem := make(chan struct{}, 1)
	var wg sync.WaitGroup
	started := 0
	var completed atomic.Int64
	for index, expansion := range plan.Expansions {
		if ctx.Err() != nil {
			break
		}
		sem <- struct{}{}
		if ctx.Err() != nil {
			<-sem
			break
		}
		started++
		currentIndex, currentExpansion := index, expansion
		wg.Go(func() {
			defer func() { <-sem }()
			if ctx.Err() != nil {
				return
			}
			request := domain.DiscoveryRequest{RequirementIDs: append([]string(nil), currentExpansion.RequirementIDs...), Derivation: cloneDerivation(currentExpansion.Derivation), ExpansionID: currentExpansion.ID, Query: currentExpansion.Query, Discriminators: append([]string(nil), currentExpansion.Discriminators...), RepositoryScope: currentExpansion.RepositoryScope, Language: language, MaxResults: maxResults, AttemptTimeout: policy.AttemptTimeout, MaxAttempts: policy.MaxAttempts}
			if currentExpansion.Regional != nil {
				request.SearchLocale = currentExpansion.Regional.Locale
				request.RegionalProposalID = currentExpansion.Regional.ProposalID
			}
			result, err := provider.Discover(ctx, request)
			if result.Attempts != nil {
				result.Attempts = append([]domain.ProviderAttempt{}, result.Attempts...)
			}
			for attemptIndex := range result.Attempts {
				if result.Attempts[attemptIndex].SearchLocale == "" {
					result.Attempts[attemptIndex].SearchLocale = request.SearchLocale
				}
				if result.Attempts[attemptIndex].RegionalProposalID == "" {
					result.Attempts[attemptIndex].RegionalProposalID = request.RegionalProposalID
				}
			}
			outcomes[currentIndex].discovery = domain.ExpansionDiscovery{Expansion: currentExpansion, Provider: result}
			if err != nil {
				outcomes[currentIndex].failure = &domain.ExecutionFailure{Scope: currentExpansion.ID, Stage: "discovery", Message: err.Error()}
			}
			done := int(completed.Add(1))
			if reporter != nil {
				status := "SUCCEEDED"
				if err != nil {
					status = "FAILED"
				}
				reporter.ReportDiscoveryProgress(domain.DiscoveryProgressEvent{ExpansionID: currentExpansion.ID, AttemptIndex: -1, Status: status, Completed: done, Planned: len(plan.Expansions)})
			}
		})
	}
	wg.Wait()
	if started == 0 {
		return domain.DiscoverySnapshot{}, fmt.Errorf("descoberta cancelada antes de iniciar: %w", ctx.Err())
	}

	snapshot := domain.DiscoverySnapshot{
		SchemaVersion: "2", Idea: normalizeIdea(plan.Idea), Strategy: strings.ToLower(strings.TrimSpace(strategy)), Language: language,
		PerExpansion: make([]domain.ExpansionDiscovery, 0, started), Candidates: make([]domain.DiscoveredCandidate, 0),
	}
	if locales := searchLocalesForExpansionPlan(plan); len(locales) > 0 {
		snapshot.SearchLocales = locales
	}
	candidateIndexes := make(map[string]int)
	failures := 0
	for index := 0; index < len(outcomes); index++ {
		outcome := outcomes[index]
		if outcome.discovery.Expansion.ID == "" {
			continue
		}
		canonicalizeDiscoveryProvider(&outcome.discovery.Provider)
		snapshot.PerExpansion = append(snapshot.PerExpansion, outcome.discovery)
		if outcome.failure != nil {
			failures++
			snapshot.Failures = append(snapshot.Failures, *outcome.failure)
		}
		for _, candidate := range outcome.discovery.Provider.Candidates {
			key := repositoryKey(candidate)
			if scope := strings.ToLower(strings.TrimSpace(outcome.discovery.Expansion.RepositoryScope)); scope != "" && key != scope {
				continue
			}
			if key == "/" {
				continue
			}
			origins := candidateOrigins(outcome.discovery, candidate)
			for _, origin := range origins {
				if err := domain.ValidateRegionalCandidateOrigin(origin); err != nil {
					return domain.DiscoverySnapshot{}, err
				}
			}
			if existingIndex, exists := candidateIndexes[key]; exists {
				existing := &snapshot.Candidates[existingIndex]
				existing.ExpansionIDs = appendUnique(existing.ExpansionIDs, outcome.discovery.Expansion.ID)
				existing.Queries = appendUnique(existing.Queries, outcome.discovery.Expansion.Query)
				existing.Origins = append(existing.Origins, origins...)
				if directCandidateMetadataLess(candidate, existing.Repository) {
					existing.Repository = candidate
				}
				continue
			}
			candidateIndexes[key] = len(snapshot.Candidates)
			snapshot.Candidates = append(snapshot.Candidates, domain.DiscoveredCandidate{
				Repository: candidate, ExpansionIDs: []string{outcome.discovery.Expansion.ID}, Queries: []string{outcome.discovery.Expansion.Query}, Origins: origins,
			})
		}
	}
	for index := range snapshot.Candidates {
		snapshot.Candidates[index].Origins = canonicalRegionalOrigins(snapshot.Candidates[index].Origins)
	}
	canonicalizeDirectDiscoverySnapshot(&snapshot)
	receipt := buildDiscoveryReceipt(policy, plan, outcomes, started, int(completed.Load()), failures, ctx.Err())
	snapshot.Receipt = &receipt
	if ctx.Err() != nil {
		return snapshot, fmt.Errorf("descoberta %s: %w", receipt.CompletionCause, ctx.Err())
	}
	if failures == started {
		message := "erro não informado"
		if len(snapshot.Failures) > 0 {
			message = snapshot.Failures[0].Message
		}
		return snapshot, fmt.Errorf("todas as %d expansões falharam: %s", failures, message)
	}
	return snapshot, nil
}

func canonicalizeDirectDiscoverySnapshot(snapshot *domain.DiscoverySnapshot) {
	for index := range snapshot.PerExpansion {
		canonicalizeDiscoveryProvider(&snapshot.PerExpansion[index].Provider)
	}
	sort.Slice(snapshot.Candidates, func(left, right int) bool {
		return discoveredCandidateLess(snapshot.Candidates[left], snapshot.Candidates[right])
	})
}

func searchLocalesForExpansionPlan(plan domain.ExpansionPlan) []string {
	values := []string{}
	if plan.Regional != nil {
		values = append(values, plan.Regional.Requested...)
	}
	for _, expansion := range plan.Expansions {
		if expansion.Regional != nil {
			values = append(values, expansion.Regional.Locale)
		}
	}
	locales, err := canonicalizeSearchLocales(values)
	if err != nil {
		return nil
	}
	return locales
}

func buildDiscoveryReceipt(policy domain.DiscoveryExecutionPolicy, plan domain.ExpansionPlan, outcomes []discoveryOutcome, started, completed, expansionFailures int, runErr error) domain.DiscoveryRunReceipt {
	r := domain.DiscoveryRunReceipt{SchemaVersion: "1", Status: domain.DiscoveryRunComplete, CompletionCause: domain.CompletionFinished, Policy: policy, Planned: len(plan.Expansions), Started: started, Completed: completed, Attempts: []domain.ProviderAttemptReceipt{}}
	for i, outcome := range outcomes {
		if i >= started || outcome.discovery.Expansion.ID == "" {
			r.Attempts = append(r.Attempts, domain.ProviderAttemptReceipt{ExpansionID: plan.Expansions[i].ID, AttemptIndex: -1, PlannedQuery: plan.Expansions[i].Query, Status: "NOT_STARTED"})
			continue
		}
		for index, attempt := range outcome.discovery.Provider.Attempts {
			status := attempt.Status
			if status == "" {
				if attempt.Error != "" {
					status = "FAILED"
				} else if attempt.ResultCount == 0 {
					status = "EMPTY"
				} else {
					status = "SUCCEEDED"
				}
			}
			receipt := domain.ProviderAttemptReceipt{ExpansionID: outcome.discovery.Expansion.ID, AttemptIndex: index, PlannedQuery: outcome.discovery.Expansion.Query, EffectiveQuery: attempt.Query, StartedAt: attempt.StartedAt, Duration: attempt.Duration, Status: status, ResultCount: attempt.ResultCount, ErrorClass: attempt.ErrorClass, SafeMessage: attempt.Error, SearchLocale: attempt.SearchLocale, RegionalProposalID: attempt.RegionalProposalID}
			if outcome.discovery.Expansion.Regional != nil {
				if receipt.SearchLocale == "" {
					receipt.SearchLocale = outcome.discovery.Expansion.Regional.Locale
				}
				if receipt.RegionalProposalID == "" {
					receipt.RegionalProposalID = outcome.discovery.Expansion.Regional.ProposalID
				}
			}
			r.Attempts = append(r.Attempts, receipt)
			switch status {
			case "TIMED_OUT":
				r.TimedOut++
			case "CANCELLED":
				r.Cancelled++
			case "FAILED":
				r.Failed++
			}
		}
	}
	if runErr != nil {
		r.Status = domain.DiscoveryRunPartial
		if errors.Is(runErr, context.DeadlineExceeded) {
			r.CompletionCause = domain.CompletionTotalTimeout
		} else {
			r.CompletionCause = domain.CompletionCancelled
		}
	} else if started > 0 && expansionFailures == started {
		r.Status, r.CompletionCause = domain.DiscoveryRunFailed, domain.CompletionProviderFailed
	}
	return r
}

func candidateOrigin(discovery domain.ExpansionDiscovery, candidate domain.RepositoryCandidate) domain.CandidateOrigin {
	origin := domain.CandidateOrigin{ExpansionID: discovery.Expansion.ID, PlannedQuery: discovery.Expansion.Query, AttemptIndex: -1}
	if discovery.Expansion.Regional != nil {
		origin.SearchLocale = discovery.Expansion.Regional.Locale
		origin.RegionalProposalID = discovery.Expansion.Regional.ProposalID
	}
	key := repositoryKey(candidate)
	for _, match := range discovery.Provider.Matches {
		if strings.ToLower(strings.TrimSpace(match.Repository)) == key {
			origin.AttemptIndex = match.AttemptIndex
			origin.EffectiveQuery = match.EffectiveQuery
			origin.Evidence = append([]domain.MatchEvidence(nil), match.Evidence...)
			if origin.EffectiveQuery == "" {
				if match.AttemptIndex >= 0 && match.AttemptIndex < len(discovery.Provider.Attempts) {
					origin.EffectiveQuery = discovery.Provider.Attempts[match.AttemptIndex].Query
				}
				if origin.EffectiveQuery == "" {
					for _, attempt := range discovery.Provider.Attempts {
						if attempt.AttemptIndex == match.AttemptIndex && attempt.Query != "" {
							origin.EffectiveQuery = attempt.Query
							break
						}
					}
				}
			}
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

// Preserve the established metadata order, with a complete declaration tie-breaker.
func directCandidateMetadataLess(left, right domain.RepositoryCandidate) bool {
	if candidateMetadataLess(left, right) {
		return true
	}
	if candidateMetadataLess(right, left) {
		return false
	}
	a, _ := json.Marshal(left)
	b, _ := json.Marshal(right)
	return string(a) < string(b)
}
func canonicalizeDiscoveryProvider(provider *domain.ProviderDiscovery) {
	if provider.Candidates != nil {
		provider.Candidates = append([]domain.RepositoryCandidate{}, provider.Candidates...)
	}
	sort.Slice(provider.Candidates, func(i, j int) bool {
		return directCandidateMetadataLess(provider.Candidates[i], provider.Candidates[j])
	})
	if provider.Matches != nil {
		provider.Matches = append([]domain.ProviderCandidateMatch{}, provider.Matches...)
	}
	for i := range provider.Matches {
		evidence := provider.Matches[i].Evidence
		if evidence != nil {
			evidence = append([]domain.MatchEvidence{}, evidence...)
		}
		sort.Slice(evidence, func(i, j int) bool {
			a, _ := json.Marshal(evidence[i])
			b, _ := json.Marshal(evidence[j])
			return string(a) < string(b)
		})
		provider.Matches[i].Evidence = evidence
	}
	sort.Slice(provider.Matches, func(i, j int) bool {
		a, _ := json.Marshal(provider.Matches[i])
		b, _ := json.Marshal(provider.Matches[j])
		return string(a) < string(b)
	})
}
func candidateOrigins(discovery domain.ExpansionDiscovery, candidate domain.RepositoryCandidate) []domain.CandidateOrigin {
	var origins []domain.CandidateOrigin
	for _, match := range discovery.Provider.Matches {
		if strings.ToLower(strings.TrimSpace(match.Repository)) != repositoryKey(candidate) {
			continue
		}
		single := discovery
		single.Provider.Matches = []domain.ProviderCandidateMatch{match}
		origins = append(origins, candidateOrigin(single, candidate))
	}
	if len(origins) == 0 {
		origins = append(origins, candidateOrigin(discovery, candidate))
	}
	return canonicalRegionalOrigins(origins)
}
