/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_run.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type IdeaRunDependencies struct {
	Runtime         domain.ExecutionRuntime
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
	ctx = domain.WithExecutionRuntime(ctx, deps.Runtime)
	canonicalLocales, localeErr := canonicalizeSearchLocales(req.SearchLocales)
	if localeErr != nil {
		return domain.IdeaRunResult{}, localeErr
	}
	req.SearchLocales = canonicalLocales
	if err := validateIdeaRun(req, deps); err != nil {
		return domain.IdeaRunResult{}, err
	}
	rules := deps.EvidenceRules
	if !evidenceRuleCatalogProvided(rules) {
		rules = DefaultEvidenceRules()
	}
	catalogCoverage, err := AnalyzeEvidenceCatalogCoverage(req.Idea, rules)
	if err != nil {
		return domain.IdeaRunResult{}, err
	}
	durations := map[string]int64{}
	started := domain.RuntimeNow(ctx)
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
	plan, err = addRegionalDiscovery(plan, req.SearchLocales)
	if err != nil {
		return domain.IdeaRunResult{}, err
	}
	durations["expand"] = domain.ElapsedMillis(domain.RuntimeNow(ctx), started)
	result := domain.IdeaRunResult{Expansion: plan, CatalogCoverage: catalogCoverage, StageDurations: durations}
	started = domain.RuntimeNow(ctx)
	if req.Resume && deps.ReplaySnapshot != nil {
		if !equalSearchLocales(req.SearchLocales, deps.ReplaySnapshot.SearchLocales) {
			return result, fmt.Errorf("snapshot de retomada pertence a outro conjunto de search locales")
		}
		result.Discovery = *deps.ReplaySnapshot
	} else {
		switch req.Mode {
		case domain.IdeaRunReplay:
			if !equalSearchLocales(req.SearchLocales, deps.ReplaySnapshot.SearchLocales) {
				return result, fmt.Errorf("snapshot de replay pertence a outro conjunto de search locales")
			}
			result.Discovery = *deps.ReplaySnapshot
		case domain.IdeaRunLive:
			maxPerQuery := req.MaxPerQuery
			if maxPerQuery == 0 {
				maxPerQuery = 5
			}
			snapshot, e := DiscoverExpansionPlan(ctx, plan, req.Strategy, req.Language, maxPerQuery, 2, deps.Discovery)
			result.Discovery = snapshot
			if e != nil {
				return result, e
			}
		case domain.IdeaRunLocal:
			result.Discovery = localDiscoverySnapshot(req, plan, deps.LocalCandidates)
		}
	}
	durations["discover"] = domain.ElapsedMillis(domain.RuntimeNow(ctx), started)
	started = domain.RuntimeNow(ctx)
	model := deps.ScreeningModel
	if model.SchemaVersion == "" {
		model = DefaultScreeningModel()
	}
	result.Screening, err = ScreenDiscovery(result.Discovery, model)
	if err != nil {
		return result, err
	}
	durations["screen"] = domain.ElapsedMillis(domain.RuntimeNow(ctx), started)
	started = domain.RuntimeNow(ctx)
	maxCandidates := req.MaxCandidates
	if maxCandidates == 0 {
		maxCandidates = 5
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
	if len(result.Plans) > maxCandidates {
		result.Plans = append([]domain.EvidenceExtractionPlan(nil), result.Plans[:maxCandidates]...)
	}
	byKey := map[string]domain.RepositoryCandidate{}
	for _, discovered := range result.Discovery.Candidates {
		candidate := discovered.Repository
		byKey[candidate.Owner+"/"+candidate.Repo] = candidate
	}
	for i := range result.Plans {
		if candidate, ok := byKey[result.Plans[i].Repository.Owner+"/"+result.Plans[i].Repository.Repo]; ok {
			result.Plans[i].Candidate = candidate
			result.Plans[i].Repository.Commit = candidate.CommitID
		}
	}
	if req.Mode == domain.IdeaRunReplay || (req.Resume && len(deps.ReplayBundles) > 0) {
		result.Bundles = replayBundlesForPlans(result.Plans, deps.ReplayBundles)
	} else {
		for _, p := range result.Plans {
			if err := ctx.Err(); err != nil {
				return result, err
			}
			b, e := CollectEvidence(ctx, p, deps.Content)
			if e != nil {
				b = domain.EvidenceBundle{SchemaVersion: "1", Repository: p.Repository, Plan: p, Complete: false, Failures: []domain.ExecutionFailure{{Scope: p.Repository.Key(), Stage: "collection", Message: e.Error()}}}
			}
			result.Bundles = append(result.Bundles, b)
			if err := ctx.Err(); err != nil {
				return result, err
			}
		}
	}
	for _, bundle := range result.Bundles {
		result.Failures = append(result.Failures, bundle.Failures...)
	}
	propagateResolvedDiscoveryCommits(&result.Discovery, result.Bundles)
	canonicalizeDiscoveryCandidateOrder(&result.Discovery)
	durations["collect"] = domain.ElapsedMillis(domain.RuntimeNow(ctx), started)
	started = domain.RuntimeNow(ctx)
	validBundles := make([]domain.EvidenceBundle, 0, len(result.Bundles))
	for _, bundle := range result.Bundles {
		if bundle.Repository.Commit != "" && !strings.EqualFold(bundle.Repository.Commit, "HEAD") {
			validBundles = append(validBundles, bundle)
		}
	}
	result.Assessment, err = AssessEvidence(req.Idea, result.Screening, validBundles, rules)
	if err != nil {
		return result, err
	}
	validAssessments := result.Assessment.Assessments
	result.Assessment.Assessments = make([]domain.CandidateAssessment, 0, len(result.Bundles))
	validIndex := 0
	for _, bundle := range result.Bundles {
		if bundle.Repository.Commit == "" || strings.EqualFold(bundle.Repository.Commit, "HEAD") {
			result.Assessment.Assessments = append(result.Assessment.Assessments, inconclusiveRunAssessment(req.Idea, bundle, rules, catalogCoverage))
			continue
		}
		result.Assessment.Assessments = append(result.Assessment.Assessments, validAssessments[validIndex])
		validIndex++
	}
	durations["assess"] = domain.ElapsedMillis(domain.RuntimeNow(ctx), started)
	started = domain.RuntimeNow(ctx)
	result.Decision, err = DecideReuse(domain.ReuseDecisionRequest{Idea: req.Idea, Assessments: result.Assessment.Assessments, Policy: req.Policy})
	if err != nil {
		return result, err
	}
	investigation := BuildInvestigationPlanWithCatalogCoverage(req.Idea, result.Discovery, result.Assessment.Assessments, catalogCoverage)
	result.Decision.Investigation = &investigation
	durations["decide"] = domain.ElapsedMillis(domain.RuntimeNow(ctx), started)
	return result, nil
}

func replayBundlesForPlans(plans []domain.EvidenceExtractionPlan, bundles []domain.EvidenceBundle) []domain.EvidenceBundle {
	byRepository := make(map[string]map[string][]domain.EvidenceBundle, len(bundles))
	for _, bundle := range bundles {
		key := ideaRunRepositoryKey(bundle.Repository.Owner, bundle.Repository.Repo)
		if byRepository[key] == nil {
			byRepository[key] = map[string][]domain.EvidenceBundle{}
		}
		byRepository[key][bundle.Repository.Commit] = append(byRepository[key][bundle.Repository.Commit], bundle)
	}
	out := make([]domain.EvidenceBundle, 0, len(plans))
	for _, plan := range plans {
		key := ideaRunRepositoryKey(plan.Repository.Owner, plan.Repository.Repo)
		commits := sortedReplayCommits(byRepository[key])
		if mutableIdeaRunCommit(plan.Repository.Commit) {
			if len(commits) == 1 && !mutableIdeaRunCommit(commits[0]) {
				out = append(out, canonicalReplayBundle(byRepository[key][commits[0]]))
				continue
			}
			if len(commits) == 1 && commits[0] == plan.Repository.Commit {
				out = append(out, canonicalReplayBundle(byRepository[key][commits[0]]))
				continue
			}
			if len(commits) > 0 {
				message := fmt.Sprintf("resolução de replay não é única e imutável para %s: commits %s", key, quotedCommits(commits))
				out = append(out, incompleteReplayBundle(plan, message))
				continue
			}
			out = append(out, incompleteReplayBundle(plan, "bundle de replay ausente para candidato planejado"))
			continue
		}
		if exact := byRepository[key][plan.Repository.Commit]; len(exact) > 0 {
			out = append(out, canonicalReplayBundle(exact))
			continue
		}
		if len(commits) > 0 {
			message := fmt.Sprintf("bundle de replay divergente para %s: commit planejado %q, recebidos %s", key, plan.Repository.Commit, quotedCommits(commits))
			out = append(out, incompleteReplayBundle(plan, message))
			continue
		}
		out = append(out, incompleteReplayBundle(plan, "bundle de replay ausente para candidato planejado"))
	}
	return out
}

func propagateResolvedDiscoveryCommits(snapshot *domain.DiscoverySnapshot, bundles []domain.EvidenceBundle) {
	resolved := map[string]map[string]bool{}
	unresolved := map[string]bool{}
	for _, bundle := range bundles {
		key := ideaRunRepositoryKey(bundle.Repository.Owner, bundle.Repository.Repo)
		if mutableIdeaRunCommit(bundle.Repository.Commit) {
			unresolved[key] = true
			continue
		}
		if resolved[key] == nil {
			resolved[key] = map[string]bool{}
		}
		resolved[key][bundle.Repository.Commit] = true
	}
	existing := map[string]map[string]bool{}
	recordExisting := func(candidate domain.RepositoryCandidate) {
		if mutableIdeaRunCommit(candidate.CommitID) {
			return
		}
		key := ideaRunRepositoryKey(candidate.Owner, candidate.Repo)
		if existing[key] == nil {
			existing[key] = map[string]bool{}
		}
		existing[key][candidate.CommitID] = true
	}
	for _, discovered := range snapshot.Candidates {
		recordExisting(discovered.Repository)
	}
	for _, expansion := range snapshot.PerExpansion {
		for _, candidate := range expansion.Provider.Candidates {
			recordExisting(candidate)
		}
	}
	refinements := map[string]string{}
	for key, commits := range resolved {
		if unresolved[key] || len(commits) != 1 {
			continue
		}
		commit := onlyIdeaRunCommit(commits)
		if immutable := existing[key]; len(immutable) > 0 && (len(immutable) != 1 || !immutable[commit]) {
			continue
		}
		refinements[key] = commit
	}
	refine := func(candidate *domain.RepositoryCandidate) {
		if !mutableIdeaRunCommit(candidate.CommitID) {
			return
		}
		if commit, ok := refinements[ideaRunRepositoryKey(candidate.Owner, candidate.Repo)]; ok {
			candidate.CommitID = commit
		}
	}
	for i := range snapshot.Candidates {
		refine(&snapshot.Candidates[i].Repository)
	}
	for i := range snapshot.PerExpansion {
		for j := range snapshot.PerExpansion[i].Provider.Candidates {
			refine(&snapshot.PerExpansion[i].Provider.Candidates[j])
		}
	}
}

func canonicalizeDiscoveryCandidateOrder(snapshot *domain.DiscoverySnapshot) {
	sort.Slice(snapshot.Candidates, func(i, j int) bool {
		return discoveredCandidateLess(snapshot.Candidates[i], snapshot.Candidates[j])
	})
	for i := range snapshot.PerExpansion {
		sort.Slice(snapshot.PerExpansion[i].Provider.Candidates, func(left, right int) bool {
			leftCandidate := domain.DiscoveredCandidate{Repository: snapshot.PerExpansion[i].Provider.Candidates[left]}
			rightCandidate := domain.DiscoveredCandidate{Repository: snapshot.PerExpansion[i].Provider.Candidates[right]}
			return discoveredCandidateLess(leftCandidate, rightCandidate)
		})
	}
}

func discoveredCandidateLess(left, right domain.DiscoveredCandidate) bool {
	leftKey := ideaRunRepositoryKey(left.Repository.Owner, left.Repository.Repo) + "@" + left.Repository.CommitID
	rightKey := ideaRunRepositoryKey(right.Repository.Owner, right.Repository.Repo) + "@" + right.Repository.CommitID
	if leftKey != rightKey {
		return leftKey < rightKey
	}
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return string(leftJSON) < string(rightJSON)
}

func ideaRunRepositoryKey(owner, repo string) string {
	return strings.ToLower(strings.TrimSpace(owner)) + "/" + strings.ToLower(strings.TrimSpace(repo))
}

func mutableIdeaRunCommit(commit string) bool {
	return strings.TrimSpace(commit) == "" || strings.EqualFold(strings.TrimSpace(commit), "HEAD")
}

func onlyIdeaRunCommit(commits map[string]bool) string {
	for commit := range commits {
		return commit
	}
	return ""
}

func sortedReplayCommits(bundles map[string][]domain.EvidenceBundle) []string {
	commits := make([]string, 0, len(bundles))
	for commit := range bundles {
		commits = append(commits, commit)
	}
	sort.Strings(commits)
	return commits
}

func quotedCommits(commits []string) string {
	quoted := make([]string, 0, len(commits))
	for _, commit := range commits {
		quoted = append(quoted, fmt.Sprintf("%q", commit))
	}
	return strings.Join(quoted, ", ")
}

func canonicalReplayBundle(bundles []domain.EvidenceBundle) domain.EvidenceBundle {
	best := bundles[0]
	bestJSON, _ := json.Marshal(best)
	for _, bundle := range bundles[1:] {
		encoded, _ := json.Marshal(bundle)
		if string(encoded) < string(bestJSON) {
			best = bundle
			bestJSON = encoded
		}
	}
	return best
}

func incompleteReplayBundle(plan domain.EvidenceExtractionPlan, message string) domain.EvidenceBundle {
	return domain.EvidenceBundle{
		SchemaVersion: "1",
		Repository:    plan.Repository,
		Plan:          plan,
		Complete:      false,
		Failures: []domain.ExecutionFailure{{
			Scope:   plan.Repository.Key(),
			Stage:   "replay",
			Message: message,
		}},
	}
}

func inconclusiveRunAssessment(idea domain.IdeaSpec, bundle domain.EvidenceBundle, catalog domain.CapabilityRuleCatalog, coverage domain.EvidenceCatalogCoverage) domain.CandidateAssessment {
	incomplete := bundle
	incomplete.Complete = false
	incomplete.TreeComplete = false
	a := domain.CandidateAssessment{
		Repository: bundle.Repository, ScreeningDecision: bundle.Plan.Decision,
		ProblemAlignment: assessCapability(idea.Problem, incomplete, findRule(idea.Problem, catalog)),
		Capabilities:     []domain.CapabilityEvidence{}, Constraints: []domain.ConstraintEvidence{},
		License:       domain.LicenseAssessment{Status: domain.EvidenceInconclusive, Caveat: "coleta incompleta"},
		Maintenance:   domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive},
		Operation:     domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive},
		ReusableParts: []domain.ReusableComponent{}, Gaps: []string{},
		Failures: append([]domain.ExecutionFailure(nil), bundle.Failures...),
	}
	capabilities := append([]string(nil), idea.Capabilities...)
	sort.Slice(capabilities, func(i, j int) bool {
		left, right := normalizeEvidenceText(capabilities[i]), normalizeEvidenceText(capabilities[j])
		if left != right {
			return left < right
		}
		return capabilities[i] < capabilities[j]
	})
	for _, capability := range capabilities {
		evidence := assessCapability(capability, incomplete, findRule(capability, catalog))
		a.Capabilities = append(a.Capabilities, evidence)
		if catalogCoverageStatus(coverage, "capability", capability) == domain.CatalogCoverageCovered {
			a.Gaps = append(a.Gaps, capability)
		}
	}
	return a
}

func evidenceRuleCatalogProvided(catalog domain.CapabilityRuleCatalog) bool {
	return catalog.SchemaVersion != "" || catalog.Version != "" || catalog.Hash != "" || len(catalog.Rules) != 0
}

func validateIdeaRun(req domain.IdeaRunRequest, deps IdeaRunDependencies) error {
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
	out := domain.DiscoverySnapshot{SchemaVersion: "1", Idea: req.Idea, Strategy: req.Strategy, Language: req.Language, SearchLocales: append([]string(nil), req.SearchLocales...)}
	for _, c := range sorted {
		origins := []domain.CandidateOrigin{}
		ids := []string{}
		queries := []string{}
		for _, x := range plan.Expansions {
			ids = append(ids, x.ID)
			queries = append(queries, x.Query)
			origin := domain.CandidateOrigin{ExpansionID: x.ID, PlannedQuery: x.Query, EffectiveQuery: x.Query}
			if x.Regional != nil {
				origin.SearchLocale = x.Regional.Locale
				origin.RegionalProposalID = x.Regional.ProposalID
			}
			origins = append(origins, origin)
		}
		out.Candidates = append(out.Candidates, domain.DiscoveredCandidate{Repository: c, ExpansionIDs: ids, Queries: queries, Origins: canonicalRegionalOrigins(origins)})
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
