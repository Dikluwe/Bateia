/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_run.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

type IdeaRunDependencies struct {
	Semantic        SemanticExpander
	Discovery       domain.ObservableDiscoveryProvider
	Content         domain.ContentProvider
	ReplaySnapshot  *domain.DiscoverySnapshot
	ReplayBundles   []domain.EvidenceBundle
	LocalCandidates []domain.RepositoryCandidate
	ScreeningModel  domain.ScreeningModel
	EvidenceRules   domain.CapabilityRuleCatalog
}

func RunIdeaPipeline(ctx context.Context, req domain.IdeaRunRequest, deps IdeaRunDependencies) (domain.IdeaRunResult, error) {
	if err := validateIdeaRun(req, deps); err != nil {
		return domain.IdeaRunResult{}, err
	}
	durations := map[string]int64{}
	started := time.Now()
	literal, err := ExpandIdea(req.Idea)
	if err != nil {
		return domain.IdeaRunResult{}, err
	}
	plan := *literal
	if req.Strategy != "" && req.Strategy != "literal" {
		enriched, e := EnrichExpansionPlan(plan, deps.Semantic)
		if e != nil {
			return domain.IdeaRunResult{}, e
		}
		plan = enriched.Plan
	}
	durations["expand"] = time.Since(started).Milliseconds()
	result := domain.IdeaRunResult{Expansion: plan, StageDurations: durations}
	started = time.Now()
	if req.Resume && deps.ReplaySnapshot != nil {
		result.Discovery = *deps.ReplaySnapshot
	} else {
		switch req.Mode {
		case domain.IdeaRunReplay:
			result.Discovery = *deps.ReplaySnapshot
		case domain.IdeaRunLive:
			maxPerQuery := req.MaxPerQuery
			if maxPerQuery == 0 {
				maxPerQuery = 5
			}
			snapshot, e := DiscoverExpansionPlan(ctx, plan, req.Strategy, req.Language, maxPerQuery, 2, deps.Discovery)
			if e != nil {
				return result, e
			}
			result.Discovery = snapshot
		case domain.IdeaRunLocal:
			result.Discovery = localDiscoverySnapshot(req, plan, deps.LocalCandidates)
		}
	}
	durations["discover"] = time.Since(started).Milliseconds()
	started = time.Now()
	model := deps.ScreeningModel
	if model.SchemaVersion == "" {
		model = DefaultScreeningModel()
	}
	result.Screening, err = ScreenDiscovery(result.Discovery, model)
	if err != nil {
		return result, err
	}
	durations["screen"] = time.Since(started).Milliseconds()
	started = time.Now()
	maxCandidates := req.MaxCandidates
	if maxCandidates == 0 {
		maxCandidates = 5
	}
	if len(result.Plans) > maxCandidates {
		result.Plans = append([]domain.EvidenceExtractionPlan(nil), result.Plans[:maxCandidates]...)
	}
	rules := deps.EvidenceRules
	if rules.SchemaVersion == "" {
		rules = DefaultEvidenceRules()
	}
	includes := []string{}
	if req.Mode == domain.IdeaRunLocal {
		for _, c := range deps.LocalCandidates {
			includes = append(includes, c.Owner+"/"+c.Repo)
		}
	}
	result.Plans, err = PlanEvidenceForIdeaAndRepositories(req.Idea, result.Screening, rules, includes)
	if err != nil {
		return result, err
	}
	if req.Mode == domain.IdeaRunLocal {
		byKey := map[string]domain.RepositoryCandidate{}
		for _, c := range deps.LocalCandidates {
			byKey[c.Owner+"/"+c.Repo] = c
		}
		for i := range result.Plans {
			if c, ok := byKey[result.Plans[i].Repository.Owner+"/"+result.Plans[i].Repository.Repo]; ok {
				result.Plans[i].Candidate = c
				result.Plans[i].Repository.Commit = c.CommitID
			}
		}
	}
	if req.Mode == domain.IdeaRunReplay || (req.Resume && len(deps.ReplayBundles) > 0) {
		result.Bundles = append([]domain.EvidenceBundle(nil), deps.ReplayBundles...)
	} else {
		for _, p := range result.Plans {
			b, e := CollectEvidence(ctx, p, deps.Content)
			if e != nil {
				b = domain.EvidenceBundle{SchemaVersion: "1", Repository: p.Repository, Plan: p, Complete: false, Failures: []domain.ExecutionFailure{{Scope: p.Repository.Key(), Stage: "collection", Message: e.Error()}}}
				result.Failures = append(result.Failures, b.Failures...)
			}
			result.Bundles = append(result.Bundles, b)
		}
	}
	durations["collect"] = time.Since(started).Milliseconds()
	started = time.Now()
	validBundles := make([]domain.EvidenceBundle, 0, len(result.Bundles))
	failedBundles := make([]domain.EvidenceBundle, 0)
	for _, bundle := range result.Bundles {
		if bundle.Repository.Commit == "" || strings.EqualFold(bundle.Repository.Commit, "HEAD") {
			failedBundles = append(failedBundles, bundle)
		} else {
			validBundles = append(validBundles, bundle)
		}
	}
	result.Assessment, err = AssessEvidence(req.Idea, result.Screening, validBundles, rules)
	if err != nil {
		return result, err
	}
	for _, bundle := range failedBundles {
		result.Assessment.Assessments = append(result.Assessment.Assessments, inconclusiveRunAssessment(req.Idea, bundle))
	}
	durations["assess"] = time.Since(started).Milliseconds()
	started = time.Now()
	result.Decision, err = DecideReuse(domain.ReuseDecisionRequest{Idea: req.Idea, Assessments: result.Assessment.Assessments, Policy: req.Policy})
	if err != nil {
		return result, err
	}
	durations["decide"] = time.Since(started).Milliseconds()
	return result, nil
}

func inconclusiveRunAssessment(idea domain.IdeaSpec, bundle domain.EvidenceBundle) domain.CandidateAssessment {
	a := domain.CandidateAssessment{Repository: bundle.Repository, ScreeningDecision: bundle.Plan.Decision, ProblemAlignment: domain.CapabilityEvidence{Capability: idea.Problem, Status: domain.EvidenceInconclusive, Explanation: "repositório não resolvido em commit imutável"}, License: domain.LicenseAssessment{Status: domain.EvidenceInconclusive, Caveat: "coleta incompleta"}, Maintenance: domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive}, Operation: domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive}, Failures: append([]domain.ExecutionFailure(nil), bundle.Failures...)}
	for _, capability := range idea.Capabilities {
		a.Capabilities = append(a.Capabilities, domain.CapabilityEvidence{Capability: capability, Status: domain.EvidenceInconclusive, Explanation: "coleta incompleta"})
	}
	return a
}

func validateIdeaRun(req domain.IdeaRunRequest, deps IdeaRunDependencies) error {
	if req.OutputDirectory == "" {
		return fmt.Errorf("output_directory é obrigatório")
	}
	if req.Strategy == "" {
		return fmt.Errorf("strategy é obrigatória")
	}
	switch req.Mode {
	case domain.IdeaRunReplay:
		if deps.ReplaySnapshot == nil {
			return fmt.Errorf("replay exige snapshot")
		}
		if len(deps.ReplayBundles) == 0 {
			return fmt.Errorf("replay exige bundles")
		}
	case domain.IdeaRunLive:
		if (!req.Resume || deps.ReplaySnapshot == nil) && (deps.Discovery == nil || deps.Content == nil) {
			return fmt.Errorf("live exige provedores de descoberta e conteúdo")
		}
	case domain.IdeaRunLocal:
		if len(deps.LocalCandidates) == 0 || deps.Content == nil {
			return fmt.Errorf("local exige candidatos e provedor de conteúdo")
		}
	default:
		return fmt.Errorf("modo %q desconhecido", req.Mode)
	}
	return nil
}
func localDiscoverySnapshot(req domain.IdeaRunRequest, plan domain.ExpansionPlan, candidates []domain.RepositoryCandidate) domain.DiscoverySnapshot {
	sorted := append([]domain.RepositoryCandidate(nil), candidates...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Owner+"/"+sorted[i].Repo < sorted[j].Owner+"/"+sorted[j].Repo })
	out := domain.DiscoverySnapshot{SchemaVersion: "1", Idea: req.Idea, Strategy: req.Strategy, Language: req.Language}
	for _, c := range sorted {
		origins := []domain.CandidateOrigin{}
		ids := []string{}
		queries := []string{}
		for _, x := range plan.Expansions {
			ids = append(ids, x.ID)
			queries = append(queries, x.Query)
			origins = append(origins, domain.CandidateOrigin{ExpansionID: x.ID, PlannedQuery: x.Query, EffectiveQuery: x.Query})
		}
		out.Candidates = append(out.Candidates, domain.DiscoveredCandidate{Repository: c, ExpansionIDs: ids, Queries: queries, Origins: origins})
	}
	return out
}
func IdeaRunStatusFor(result domain.IdeaRunResult) domain.IdeaRunStatus {
	if result.Decision.Decision == domain.DecisionUnknown {
		return domain.IdeaRunInconclusive
	}
	if len(result.Failures) > 0 || len(result.Discovery.Failures) > 0 {
		return domain.IdeaRunPartial
	}
	return domain.IdeaRunComplete
}
func StableRunID(inputHash string) string {
	value := strings.TrimSpace(inputHash)
	if len(value) > 16 {
		return value[:16]
	}
	return value
}
