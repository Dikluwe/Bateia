/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence.md
 * @layer L1
 * @updated 2026-09-16
 */

package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	pathpkg "path"
	"sort"
	"strings"
	"time"
	"unicode"

	"bateia/01_core/domain"
)

var defaultEvidenceBudget = domain.ExtractionBudget{
	Documents: 8, Manifests: 4, CodeFiles: 20, Tests: 12,
	TotalBytes: 2 * 1024 * 1024, TreeDepth: 4,
}

func DefaultEvidenceRules() domain.CapabilityRuleCatalog {
	rules := []domain.CapabilityRule{
		rule("dependency-graph", []string{"grafo de dependencias", "dependency graph"}, []string{"dependencygraph", "buildgraph", "dependency graph"}),
		rule("transitive-dependents", []string{"dependentes transitivos", "transitive dependents"}, []string{"transitivedependents", "transitive dependents", "reverse dependencies"}),
		rule("blast-radius", []string{"impacto", "blast radius"}, []string{"blastradius", "impact analysis", "blast radius", "spillover when the change stays in scope", "derived targets ignore additions"}),
		rule("structural-change", []string{"comparacao estrutural de mudancas", "structural change"}, []string{"structuraldiff", "architecture diff", "structural change"}),
		rule("static-analysis", []string{"analise estatica", "static analysis", "analisar codigo sem executar o programa"}, []string{"staticanalyzer", "static analysis", "analyse code without actually running", "ast metrics", "abstract syntax tree metrics"}),
		rule("local-offline", []string{"execucao local", "offline", "local execution"}, []string{"offline", "run locally", "local execution"}),
		rule("agent-mcp", []string{"integracao mcp", "agentes", "model context protocol"}, []string{"mcpserver", "model context protocol", "mcp server"}),
		rule("async-events", []string{"processamento assincrono de eventos", "processar eventos assincronos de servicos", "asynchronous event processing"}, []string{"eventprocessor", "eventdispatcher", "eventhandler", "eventstore", "event queue", "async event", "event driven architecture", "event sourcing", "message bridge"}),
		rule("recipe-import", []string{"organizacao de receitas", "organizar receitas culinarias para uma familia", "importacao de receitas", "recipe import", "recipe organization"}, []string{"recipeimporter", "import recipe", "recipe organizer", "recipe manager", "family recipes", "mealie import", "tandoor import", "nextcloud import", "mela recipes", "recipe keeper", "recipe sage", "repository recipes", "recipe controller"}),
		independentRule("repository-discovery", []string{"descobrir repositorios candidatos", "repository discovery", "software discovery"}, []string{"repositorydiscovery", "candidatefinder", "discover repositories", "search candidates"}),
		independentRule("adoption-decision", []string{"emitir decisao de adocao de componente", "component adoption decision", "dependency adoption gate"}, []string{"adoptiondecision", "adoptiongate", "validate or drop", "go validate drop"}),
		independentRule("grounded-verdict", []string{"julgar candidato com evidencia rastreavel", "evidence grounded verdict", "traceable candidate judgment"}, []string{"groundedverdict", "evidencepointer", "traceable verdict", "evidence grounded judgment"}),
		independentRule("candidate-comparison", []string{"comparar candidatos por criterios", "criteria based candidate comparison", "compare software candidates"}, []string{"candidatecomparison", "comparisonmatrix", "compare candidates", "criteria matrix"}),
		independentRule("human-review-escalation", []string{"escalonar duvida para revisao humana", "human review escalation", "escalate uncertain judgment"}, []string{"humanreviewescalation", "reviewqueue", "escalate uncertainty", "mandatory review"}),
	}
	catalog := domain.CapabilityRuleCatalog{SchemaVersion: "2", Version: "2026-09-07.0040", Rules: rules}
	catalog.Hash = ruleCatalogHash(catalog)
	return catalog
}

func rule(id string, terms, patterns []string) domain.CapabilityRule {
	return domain.CapabilityRule{
		ID: id, Basis: domain.EvidenceRuleDeclared, Family: "software-capability", Provenance: "step-0022-contract", TargetHints: []string{"source symbol", "related test or example"}, CapabilityTerms: terms,
		StrongEvidence: []domain.CapabilityPattern{
			{Kind: domain.EvidenceSymbol, Patterns: patterns},
			{Kind: domain.EvidenceTest, Patterns: patterns},
		},
		SupportingEvidence: []domain.CapabilityPattern{
			{Kind: domain.EvidenceDocumentation, Patterns: patterns},
			{Kind: domain.EvidenceFile, Patterns: patterns},
		},
		NegativeEvidence: []string{"not implemented", "nao implementado", "planned", "roadmap only"},
	}
}

func independentRule(id string, terms, patterns []string) domain.CapabilityRule {
	r := rule(id, terms, patterns)
	r.RequiresIndependent = true
	return r
}

func ruleCatalogHash(c domain.CapabilityRuleCatalog) string {
	c.Hash = ""
	c.Rules = canonicalCapabilityRules(c.Rules)
	b, _ := json.Marshal(c)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func canonicalCapabilityRules(rules []domain.CapabilityRule) []domain.CapabilityRule {
	out := append([]domain.CapabilityRule(nil), rules...)
	for i := range out {
		if out[i].Basis == "" {
			out[i].Basis = domain.EvidenceRuleDeclared
		}
		out[i].TargetHints = canonicalRuleStrings(out[i].TargetHints)
		out[i].CapabilityTerms = canonicalRuleStrings(out[i].CapabilityTerms)
		out[i].NegativeEvidence = canonicalRuleStrings(out[i].NegativeEvidence)
		out[i].StrongEvidence = canonicalCapabilityPatterns(out[i].StrongEvidence)
		out[i].SupportingEvidence = canonicalCapabilityPatterns(out[i].SupportingEvidence)
	}
	sort.Slice(out, func(i, j int) bool {
		left, _ := json.Marshal(out[i])
		right, _ := json.Marshal(out[j])
		return string(left) < string(right)
	})
	return out
}

func canonicalRuleStrings(values []string) []string {
	out := append([]string(nil), values...)
	sort.Slice(out, func(i, j int) bool {
		left, right := normalizeEvidenceText(out[i]), normalizeEvidenceText(out[j])
		if left != right {
			return left < right
		}
		return out[i] < out[j]
	})
	return out
}

func canonicalCapabilityPatterns(values []domain.CapabilityPattern) []domain.CapabilityPattern {
	out := append([]domain.CapabilityPattern(nil), values...)
	for i := range out {
		out[i].Patterns = canonicalRuleStrings(out[i].Patterns)
	}
	sort.Slice(out, func(i, j int) bool {
		left, _ := json.Marshal(out[i])
		right, _ := json.Marshal(out[j])
		return string(left) < string(right)
	})
	return out
}

func ValidateRuleCatalog(c domain.CapabilityRuleCatalog) error {
	return validateRuleCatalog(c, false)
}

func validateRuleCatalog(c domain.CapabilityRuleCatalog, allowAmbiguousCoverage bool) error {
	if strings.TrimSpace(c.Version) == "" || len(c.Rules) == 0 {
		return fmt.Errorf("catálogo de regras exige version e rules")
	}
	ids := map[string]bool{}
	terms := map[string]string{}
	patterns := map[string]string{}
	for _, r := range c.Rules {
		if strings.TrimSpace(r.ID) == "" || ids[r.ID] {
			return fmt.Errorf("ID de regra vazio ou duplicado: %q", r.ID)
		}
		ids[r.ID] = true
		switch r.Basis {
		case "", domain.EvidenceRuleDeclared, domain.EvidenceRuleHypothesis:
		default:
			return fmt.Errorf("regra %s possui basis inválida: %q", r.ID, r.Basis)
		}
		if c.SchemaVersion == "2" && (strings.TrimSpace(r.Family) == "" || strings.TrimSpace(r.Provenance) == "" || len(r.TargetHints) == 0) {
			return fmt.Errorf("regra %s exige família, proveniência e target_hints", r.ID)
		}
		for _, term := range r.CapabilityTerms {
			n := normalizeEvidenceText(term)
			if n == "" || identityDependentRuleText(term) || overlyGenericRuleText(n) {
				return fmt.Errorf("regra %s contém termo inválido ou dependente de identidade: %q", r.ID, term)
			}
			if owner := terms[n]; !allowAmbiguousCoverage && owner != "" && owner != r.ID {
				return fmt.Errorf("termo duplicado entre %s e %s: %q", owner, r.ID, term)
			}
			terms[n] = r.ID
		}
		for _, group := range append(append([]domain.CapabilityPattern{}, r.StrongEvidence...), r.SupportingEvidence...) {
			for _, pattern := range group.Patterns {
				n := normalizeEvidenceText(pattern)
				if len(pattern) == 0 || len(pattern) > 120 || identityDependentRuleText(pattern) || overlyGenericRuleText(n) {
					return fmt.Errorf("regra %s contém padrão vazio ou longo demais", r.ID)
				}
				if owner := patterns[n]; !allowAmbiguousCoverage && owner != "" && owner != r.ID {
					return fmt.Errorf("padrão duplicado entre %s e %s: %q", owner, r.ID, pattern)
				}
				patterns[n] = r.ID
			}
		}
	}
	if c.Hash != "" && c.Hash != ruleCatalogHash(c) {
		return fmt.Errorf("hash do catálogo não corresponde ao conteúdo")
	}
	return nil
}

func AnalyzeEvidenceCatalogCoverage(idea domain.IdeaSpec, catalog domain.CapabilityRuleCatalog) (domain.EvidenceCatalogCoverage, error) {
	if err := validateRuleCatalog(catalog, true); err != nil {
		return domain.EvidenceCatalogCoverage{}, err
	}
	coverage := domain.EvidenceCatalogCoverage{
		SchemaVersion: "1",
		RulesVersion:  catalog.Version,
		RulesHash:     ruleCatalogHash(catalog),
		Entries:       []domain.CatalogCoverageEntry{},
	}
	type requirement struct {
		kind, original, normalized string
	}
	requirements := make([]requirement, 0, len(idea.Capabilities)+1)
	requirements = append(requirements, requirement{kind: "problem", original: idea.Problem, normalized: normalizeEvidenceText(idea.Problem)})
	for _, capability := range idea.Capabilities {
		requirements = append(requirements, requirement{kind: "capability", original: capability, normalized: normalizeEvidenceText(capability)})
	}
	for _, requirement := range requirements {
		entry := domain.CatalogCoverageEntry{
			RequirementKind: requirement.kind,
			Requirement:     requirement.original,
			Status:          domain.CatalogCoverageGap,
			ReasonCode:      "CATALOG_GAP",
		}
		matches := exactEvidenceRules(requirement.normalized, catalog)
		if len(matches) == 1 {
			basis := evidenceRuleBasis(matches[0])
			entry.RuleID = matches[0].ID
			entry.RuleBasis = basis
			switch basis {
			case domain.EvidenceRuleDeclared:
				entry.Status = domain.CatalogCoverageCovered
				entry.ReasonCode = "COVERED"
			case domain.EvidenceRuleHypothesis:
				entry.Status = domain.CatalogCoverageHypothesis
				entry.ReasonCode = "HYPOTHESIS_ONLY"
			}
		}
		coverage.Entries = append(coverage.Entries, entry)
	}
	sort.Slice(coverage.Entries, func(i, j int) bool {
		left, right := coverage.Entries[i], coverage.Entries[j]
		leftNormalized, rightNormalized := normalizeEvidenceText(left.Requirement), normalizeEvidenceText(right.Requirement)
		if left.RequirementKind != right.RequirementKind {
			return left.RequirementKind < right.RequirementKind
		}
		if leftNormalized != rightNormalized {
			return leftNormalized < rightNormalized
		}
		if left.Requirement != right.Requirement {
			return left.Requirement < right.Requirement
		}
		return left.RuleID < right.RuleID
	})
	return coverage, nil
}

func evidenceRuleBasis(rule domain.CapabilityRule) domain.EvidenceRuleBasis {
	if rule.Basis == "" {
		return domain.EvidenceRuleDeclared
	}
	return rule.Basis
}

func exactEvidenceRules(normalizedRequirement string, catalog domain.CapabilityRuleCatalog) []domain.CapabilityRule {
	if normalizedRequirement == "" {
		return nil
	}
	matches := []domain.CapabilityRule{}
	for _, rule := range catalog.Rules {
		matched := false
		for _, term := range rule.CapabilityTerms {
			if normalizeEvidenceText(term) == normalizedRequirement {
				matched = true
				break
			}
		}
		if matched {
			matches = append(matches, rule)
		}
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].ID < matches[j].ID })
	return matches
}

func identityDependentRuleText(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return strings.Contains(lower, "github.com/") || strings.Contains(lower, "gitlab.com/") || strings.Contains(lower, "http://") || strings.Contains(lower, "https://") || strings.Contains(lower, "@") || strings.Contains(lower, "owner/")
}

func overlyGenericRuleText(value string) bool {
	switch value {
	case "tool", "software", "repository", "repo", "manager", "controller", "service", "search", "evaluation", "avaliacao", "component", "componente":
		return true
	}
	return false
}

func PlanEvidence(report domain.ScreeningReport) ([]domain.EvidenceExtractionPlan, error) {
	return PlanEvidenceForIdea(domain.IdeaSpec{}, report, DefaultEvidenceRules())
}

func PlanEvidenceForIdea(idea domain.IdeaSpec, report domain.ScreeningReport, catalog domain.CapabilityRuleCatalog) ([]domain.EvidenceExtractionPlan, error) {
	return PlanEvidenceForIdeaAndRepositories(idea, report, catalog, nil)
}

func PlanEvidenceForIdeaAndRepositories(idea domain.IdeaSpec, report domain.ScreeningReport, catalog domain.CapabilityRuleCatalog, includes []string) ([]domain.EvidenceExtractionPlan, error) {
	if _, err := AnalyzeEvidenceCatalogCoverage(idea, catalog); err != nil {
		return nil, err
	}
	capabilities := append([]string(nil), idea.Capabilities...)
	if strings.TrimSpace(idea.Problem) != "" {
		capabilities = append(capabilities, idea.Problem)
	}
	capabilities = uniqueNormalizedStrings(capabilities)
	sort.Strings(capabilities)
	ruleIDs := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		if matched := findRule(capability, catalog); matched != nil {
			ruleIDs = append(ruleIDs, matched.ID)
		}
	}
	ruleReason := "amostra geral"
	if len(ruleIDs) > 0 {
		ruleReason = "regras: " + strings.Join(ruleIDs, ",")
	}
	selected := selectEvidenceCandidates(report.Decisions, 10, 2, 16)
	selected = includeEvidenceCandidates(selected, report.Decisions, includes)
	plans := make([]domain.EvidenceExtractionPlan, 0, len(selected))
	for _, decision := range selected {
		ref, err := parseRepositoryReference(decision.Repository)
		if err != nil {
			return nil, err
		}
		candidate := domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit, URL: "https://github.com/" + ref.Owner + "/" + ref.Repo}
		targets := []domain.ExtractionTarget{
			{Path: "README.md", Kind: domain.EvidenceDocumentation, Capabilities: capabilities, Reason: "documentação principal; " + ruleReason},
			{Path: "LICENSE", Kind: domain.EvidenceLicense, Reason: "licença do candidato"},
			{Path: "go.mod", Kind: domain.EvidenceManifest, Reason: "manifesto reconhecido"},
			{Path: "package.json", Kind: domain.EvidenceManifest, Reason: "manifesto reconhecido"},
			{Path: "Cargo.toml", Kind: domain.EvidenceManifest, Reason: "manifesto reconhecido"},
			{Path: "pyproject.toml", Kind: domain.EvidenceManifest, Reason: "manifesto reconhecido"},
		}
		plans = append(plans, domain.EvidenceExtractionPlan{SchemaVersion: "1", Repository: ref, Candidate: candidate, Decision: decision, Targets: targets, Budget: defaultEvidenceBudget})
	}
	return plans, nil
}

func uniqueNormalizedStrings(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		key := normalizeEvidenceText(value)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, value)
	}
	return out
}

func includeEvidenceCandidates(selected, available []domain.ScreeningDecision, includes []string) []domain.ScreeningDecision {
	byRepository := map[string]domain.ScreeningDecision{}
	for _, decision := range available {
		byRepository[repositoryIdentityKey(decision.Repository)] = decision
	}
	for _, decision := range selected {
		byRepository[repositoryIdentityKey(decision.Repository)] = decision
	}
	seen := map[string]bool{}
	out := append([]domain.ScreeningDecision(nil), selected...)
	for _, decision := range selected {
		seen[repositoryIdentityKey(decision.Repository)] = true
	}
	normalized := make([]string, 0, len(includes))
	for _, repository := range includes {
		repository = repositoryIdentityKey(strings.TrimSpace(repository))
		if repository != "" {
			normalized = append(normalized, repository)
		}
	}
	sort.Strings(normalized)
	for _, repository := range normalized {
		if seen[repository] {
			continue
		}
		decision, ok := byRepository[repository]
		if !ok {
			decision = domain.ScreeningDecision{Repository: repository, Decision: "AUDIT", UncertaintyReasons: []string{"controle incluído explicitamente para calibração"}}
		}
		out = append(out, decision)
		seen[repository] = true
	}
	return out
}

func selectEvidenceCandidates(decisions []domain.ScreeningDecision, auditMax, dropMax, totalMax int) []domain.ScreeningDecision {
	copyDecisions := append([]domain.ScreeningDecision(nil), decisions...)
	sort.SliceStable(copyDecisions, func(i, j int) bool {
		if copyDecisions[i].RankingScore == copyDecisions[j].RankingScore {
			return copyDecisions[i].Repository < copyDecisions[j].Repository
		}
		return copyDecisions[i].RankingScore > copyDecisions[j].RankingScore
	})
	selected := make([]domain.ScreeningDecision, 0, totalMax)
	audits := 0
	for _, d := range copyDecisions {
		if strings.EqualFold(d.Decision, "ADVANCE") || (strings.EqualFold(d.Decision, "AUDIT") && audits < auditMax) {
			selected = append(selected, d)
			if strings.EqualFold(d.Decision, "AUDIT") {
				audits++
			}
		}
	}
	drops := make([]domain.ScreeningDecision, 0)
	for _, d := range copyDecisions {
		if strings.EqualFold(d.Decision, "DROP") {
			drops = append(drops, d)
		}
	}
	sort.Slice(drops, func(i, j int) bool { return drops[i].Repository < drops[j].Repository })
	if len(drops) > dropMax {
		drops = drops[:dropMax]
	}
	selected = append(selected, drops...)
	if len(selected) > totalMax {
		selected = selected[:totalMax]
	}
	return selected
}

func parseRepositoryReference(value string) (domain.CommitRef, error) {
	parts := strings.Split(value, "@")
	if len(parts) > 2 {
		return domain.CommitRef{}, fmt.Errorf("referência de candidato inválida %q", value)
	}
	repo := strings.Split(strings.Trim(parts[0], "/"), "/")
	if len(repo) != 2 || repo[0] == "" || repo[1] == "" {
		return domain.CommitRef{}, fmt.Errorf("repositório inválido %q; use owner/repo@commit", value)
	}
	commit := "HEAD"
	if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
		commit = strings.TrimSpace(parts[1])
	}
	return domain.CommitRef{Owner: repo[0], Repo: repo[1], Commit: commit}, nil
}

func CollectEvidence(ctx context.Context, plan domain.EvidenceExtractionPlan, provider domain.ContentProvider) (domain.EvidenceBundle, error) {
	if provider == nil {
		return domain.EvidenceBundle{}, fmt.Errorf("provedor de conteúdo não configurado")
	}
	if err := ctx.Err(); err != nil {
		return domain.EvidenceBundle{}, err
	}
	if plan.Repository.Commit == "" || strings.EqualFold(plan.Repository.Commit, "HEAD") {
		resolved, err := provider.FetchRepoMetadata(ctx, plan.Candidate)
		if err != nil || resolved.CommitID == "" || strings.EqualFold(resolved.CommitID, "HEAD") {
			return domain.EvidenceBundle{}, fmt.Errorf("não foi possível resolver %s para commit imutável", plan.Repository.Key())
		}
		plan.Candidate = resolved
		plan.Repository.Commit = resolved.CommitID
		plan.Decision.Repository = plan.Repository.Key()
	} else if metadata, err := provider.FetchRepoMetadata(ctx, plan.Candidate); err == nil {
		plan.Candidate = metadata
	}
	if err := ctx.Err(); err != nil {
		return domain.EvidenceBundle{}, err
	}
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: plan.Repository, Plan: plan, Complete: true, TreeComplete: true}
	if metadata, err := json.Marshal(plan.Candidate); err == nil {
		content := string(metadata)
		bundle.Items = append(bundle.Items, domain.CollectedEvidence{Kind: domain.EvidenceRepository, Ref: plan.Repository, Path: "@repository-metadata", Content: content, ContentHash: contentHash(content)})
		bundle.BytesRead += int64(len(content))
	}
	paths := make([]string, 0)
	kinds := map[string]domain.EvidenceKind{}
	targetByPath := map[string]domain.ExtractionTarget{}
	for _, target := range plan.Targets {
		if _, exists := targetByPath[target.Path]; exists {
			continue
		}
		paths = append(paths, target.Path)
		kinds[target.Path] = target.Kind
		targetByPath[target.Path] = target
	}
	implementationQuery := ""
	tree, err := provider.FetchTree(ctx, plan.Candidate)
	if err != nil {
		bundle.Failures = append(bundle.Failures, domain.ExecutionFailure{Scope: plan.Repository.Key(), Stage: "tree", Message: err.Error()})
		bundle.Complete = false
		bundle.TreeComplete = false
	} else {
		tree = prioritizeEvidenceTree(tree, plan.Targets)
		if _, ok := provider.(domain.ImplementationSearcher); ok {
			requirements := []string{}
			for _, target := range plan.Targets {
				requirements = append(requirements, target.Capabilities...)
			}
			implementationQuery = strings.Join(uniqueNormalizedStrings(requirements), "\n")
			retrieval := RetrieveImplementationEvidence(ctx, provider, plan.Candidate, implementationQuery)
			bundle.ImplementationSearch = &retrieval
			rank := map[string]int{}
			for _, h := range retrieval.Hits {
				if _, ok := rank[h.Path]; !ok {
					rank[h.Path] = len(rank) + 1
				}
			}
			sort.SliceStable(tree, func(i, j int) bool {
				a, ao := rank[tree[i].Path]
				b, bo := rank[tree[j].Path]
				if ao != bo {
					return ao
				}
				if ao && bo {
					return a < b
				}
				return false
			})
		}
		for _, node := range tree {
			if depth := pathDepth(node.Path); depth > bundle.DepthReached {
				bundle.DepthReached = depth
			}
			if node.Type != "file" || (ignoredEvidencePath(node.Path) && !explicitEvidenceDirectory(node.Path, plan.Targets)) {
				continue
			}
			kind := classifyEvidencePath(node.Path)
			if kind == domain.EvidenceFile && !isCodePath(node.Path) {
				continue
			}
			if _, exists := targetByPath[node.Path]; exists {
				continue
			}
			target := domain.ExtractionTarget{Path: node.Path, Kind: kind, Capabilities: capabilitiesForPath(node.Path, plan.Targets), Reason: "árvore limitada no commit resolvido"}
			if pathDepth(node.Path) > plan.Budget.TreeDepth || !withinKindBudget(paths, kinds, kind, plan.Budget) {
				bundle.Complete = false
				bundle.Targets = append(bundle.Targets, domain.CollectionTargetResult{Target: target, Status: domain.TargetSkippedBudget, Message: "arquivo encontrado, não examinado: limite de profundidade ou quantidade de arquivos"})
				continue
			}
			paths = append(paths, node.Path)
			kinds[node.Path] = kind
			targetByPath[node.Path] = target
		}
	}
	if err := ctx.Err(); err != nil {
		bundle.Complete = false
		return bundle, err
	}
	if bundle.ImplementationSearch == nil {
		sort.Strings(paths)
	}
	extracted, err := provider.ExtractRepoFiles(ctx, plan.Candidate, paths)
	if err != nil {
		bundle.Failures = append(bundle.Failures, domain.ExecutionFailure{Scope: plan.Repository.Key(), Stage: "content", Message: err.Error()})
		bundle.Complete = false
		for _, path := range paths {
			bundle.Targets = append(bundle.Targets, domain.CollectionTargetResult{Target: targetByPath[path], Status: domain.TargetFailed, Message: err.Error()})
		}
		return bundle, nil
	}
	seenHash := map[string]bool{}
	collectedPaths := map[string]int64{}
	for _, file := range extracted.Files {
		content := file.Content
		remaining := plan.Budget.TotalBytes - bundle.BytesRead
		if remaining <= 0 {
			bundle.Complete = false
			break
		}
		if int64(len(content)) > remaining {
			content = content[:remaining]
			bundle.Complete = false
		}
		hash := contentHash(content)
		if seenHash[hash] {
			continue
		}
		seenHash[hash] = true
		kind := kinds[file.Path]
		if kind == "" {
			kind = classifyEvidencePath(file.Path)
		}
		bundle.Items = append(bundle.Items, domain.CollectedEvidence{Kind: kind, Ref: plan.Repository, Path: file.Path, Content: content, ContentHash: hash})
		bundle.BytesRead += int64(len(content))
		collectedPaths[file.Path] = int64(len(content))
	}
	for _, path := range paths {
		result := domain.CollectionTargetResult{Target: targetByPath[path], Status: domain.TargetCollected, Bytes: collectedPaths[path]}
		if _, ok := collectedPaths[path]; !ok {
			result.Status = domain.TargetNotFound
			result.Message = "provedor não devolveu o alvo solicitado"
			if len(result.Target.Capabilities) > 0 {
				bundle.Complete = false
			}
		}
		bundle.Targets = append(bundle.Targets, result)
	}
	if bundle.ImplementationSearch != nil {
		readFiles := []domain.ExtractedFile{}
		for _, item := range bundle.Items {
			readFiles = append(readFiles, domain.ExtractedFile{Path: item.Path, Content: item.Content})
		}
		validated := ValidateImplementationRequirements(implementationQuery, *bundle.ImplementationSearch, readFiles)
		bundle.ImplementationSearch = &validated
	}
	if bundle.BytesRead >= plan.Budget.TotalBytes {
		bundle.Complete = false
	}
	return bundle, nil
}

func prioritizeEvidenceTree(tree []domain.DirectoryNode, targets []domain.ExtractionTarget) []domain.DirectoryNode {
	out := append([]domain.DirectoryNode(nil), tree...)
	sort.SliceStable(out, func(i, j int) bool {
		iScore := evidencePathPriority(out[i].Path, targets)
		jScore := evidencePathPriority(out[j].Path, targets)
		if iScore != jScore {
			return iScore > jScore
		}
		return out[i].Path < out[j].Path
	})
	return out
}

// Path priority is a collection heuristic, never evidence of implemented behavior.
func evidencePathPriority(path string, targets []domain.ExtractionTarget) int {
	p := strings.ToLower(slashPath(path))
	score := 0
	if isCodePath(p) && classifyEvidencePath(p) == domain.EvidenceFile {
		score += 20
	}
	score += 5 * len(capabilitiesForPath(p, targets))
	for _, segment := range strings.Split(p, "/") {
		switch segment {
		case "tmp", "temp", "generated", "dist", "build", "benchmark", "benchmarks", "fixtures", "prompts", "vendor", "node_modules":
			if !explicitEvidenceDirectory(p, targets) {
				score -= 100
			}
		}
	}
	return score
}

func explicitEvidenceDirectory(path string, targets []domain.ExtractionTarget) bool {
	directory := pathpkg.Dir(strings.ToLower(slashPath(path)))
	for _, target := range targets {
		if slashPath(target.Path) == slashPath(path) {
			return true
		}
		for _, segment := range strings.Split(directory, "/") {
			switch segment {
			case "tmp", "temp", "generated", "dist", "build", "benchmark", "benchmarks", "fixtures", "prompts", "vendor", "node_modules", "translations":
				for _, capability := range target.Capabilities {
					if phraseMatch(normalizeEvidenceText(capability), segment) {
						return true
					}
				}
			}
		}
	}
	return false
}

func capabilitiesForPath(path string, targets []domain.ExtractionTarget) []string {
	pathNorm := normalizeEvidenceText(path)
	set := map[string]bool{}
	for _, target := range targets {
		for _, capability := range target.Capabilities {
			terms := strings.Fields(normalizeEvidenceText(capability))
			for _, term := range terms {
				if len(term) >= 4 && strings.Contains(pathNorm, term) {
					set[capability] = true
				}
			}
		}
	}
	out := make([]string, 0, len(set))
	for capability := range set {
		out = append(out, capability)
	}
	sort.Strings(out)
	return out
}

func pathDepth(path string) int {
	return len(strings.Split(slashPath(strings.Trim(path, "/")), "/"))
}
func ignoredEvidencePath(path string) bool {
	p := "/" + strings.ToLower(slashPath(path)) + "/"
	return strings.Contains(p, "/vendor/") || strings.Contains(p, "/fixtures/") || strings.Contains(p, "/translations/") || strings.Contains(p, "/node_modules/")
}
func classifyEvidencePath(path string) domain.EvidenceKind {
	p := strings.ToLower(slashPath(path))
	base := pathpkg.Base(p)
	if strings.HasPrefix(base, "readme") || strings.Contains(p, "/docs/") {
		return domain.EvidenceDocumentation
	}
	if authoritativeLicensePath(p) {
		return domain.EvidenceLicense
	}
	if base == "go.mod" || base == "package.json" || base == "cargo.toml" || base == "pyproject.toml" {
		return domain.EvidenceManifest
	}
	if strings.Contains(base, "_test.") || strings.Contains(p, "/test/") || strings.Contains(p, "/tests/") {
		return domain.EvidenceTest
	}
	if strings.Contains(p, "/example") || strings.Contains(p, "/examples/") {
		return domain.EvidenceExample
	}
	return domain.EvidenceFile
}

func isCodePath(path string) bool {
	extension := strings.ToLower(pathpkg.Ext(slashPath(path)))
	switch extension {
	case ".go", ".rs", ".py", ".js", ".jsx", ".ts", ".tsx", ".java", ".kt", ".rb", ".cs", ".c", ".cc", ".cpp", ".h", ".hpp", ".ex", ".exs":
		return true
	default:
		return false
	}
}

func slashPath(value string) string {
	return strings.ReplaceAll(value, "\\", "/")
}
func withinKindBudget(paths []string, kinds map[string]domain.EvidenceKind, kind domain.EvidenceKind, b domain.ExtractionBudget) bool {
	count := 0
	for _, p := range paths {
		if kinds[p] == kind {
			count++
		}
	}
	limit := b.CodeFiles
	switch kind {
	case domain.EvidenceDocumentation:
		limit = b.Documents
	case domain.EvidenceManifest:
		limit = b.Manifests
	case domain.EvidenceTest, domain.EvidenceExample:
		limit = b.Tests
	}
	return count < limit
}
func contentHash(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])
}

func AssessEvidence(idea domain.IdeaSpec, screening domain.ScreeningReport, bundles []domain.EvidenceBundle, catalog domain.CapabilityRuleCatalog) (domain.AssessmentReport, error) {
	coverage, err := AnalyzeEvidenceCatalogCoverage(idea, catalog)
	if err != nil {
		return domain.AssessmentReport{}, err
	}
	catalog.Hash = coverage.RulesHash
	decisionByRepo := map[string]domain.ScreeningDecision{}
	for _, d := range screening.Decisions {
		decisionByRepo[repositoryIdentityKey(d.Repository)] = d
	}
	report := domain.AssessmentReport{SchemaVersion: "1", RulesVersion: catalog.Version, RulesHash: catalog.Hash, CatalogCoverage: coverage, Assessments: []domain.CandidateAssessment{}}
	for _, bundle := range bundles {
		if err := validateEvidenceBundle(bundle); err != nil {
			return domain.AssessmentReport{}, err
		}
		decision, ok := decisionByRepo[bundle.Repository.Owner+"/"+bundle.Repository.Repo]
		if !ok {
			decision = bundle.Plan.Decision
		}
		assessment := domain.CandidateAssessment{
			Repository: bundle.Repository, ScreeningDecision: decision,
			Capabilities: []domain.CapabilityEvidence{}, Constraints: []domain.ConstraintEvidence{},
			ReusableParts: []domain.ReusableComponent{}, Gaps: []string{},
			Failures: append([]domain.ExecutionFailure(nil), bundle.Failures...),
		}
		assessment.ProblemAlignment = assessCapability(idea.Problem, bundle, findRule(idea.Problem, catalog))
		capabilities := append([]string(nil), idea.Capabilities...)
		sort.Slice(capabilities, func(i, j int) bool {
			left, right := normalizeEvidenceText(capabilities[i]), normalizeEvidenceText(capabilities[j])
			if left != right {
				return left < right
			}
			return capabilities[i] < capabilities[j]
		})
		for _, capability := range capabilities {
			e := assessCapability(capability, bundle, findRule(capability, catalog))
			assessment.Capabilities = append(assessment.Capabilities, e)
			if catalogCoverageStatus(coverage, "capability", capability) == domain.CatalogCoverageCovered && e.Status != domain.EvidenceProven {
				assessment.Gaps = append(assessment.Gaps, capability)
			}
			if e.Status == domain.EvidenceProven {
				if s, ok := strongestReusableSource(e.Sources); ok {
					assessment.ReusableParts = append(assessment.ReusableParts, domain.ReusableComponent{Name: s.Symbol, Kind: "symbol", Path: s.Path, Symbols: []string{s.Symbol}, Capabilities: []string{capability}, Evidence: []domain.EvidenceSource{s}})
				}
			}
		}
		for _, constraint := range idea.Constraints {
			assessment.Constraints = append(assessment.Constraints, domain.ConstraintEvidence{Constraint: constraint, Evidence: assessConstraint(constraint, bundle)})
		}
		for _, nonGoal := range idea.NonGoals {
			assessment.Constraints = append(assessment.Constraints, domain.ConstraintEvidence{Constraint: "non_goal: " + nonGoal, Evidence: assessNonGoal(nonGoal, bundle)})
		}
		assessment.License = assessLicense(bundle)
		assessment.Maintenance = assessMaintenance(bundle)
		assessment.Operation = assessOperation(bundle)
		report.Assessments = append(report.Assessments, assessment)
	}
	return report, nil
}

func catalogCoverageStatus(coverage domain.EvidenceCatalogCoverage, kind, requirement string) domain.CatalogCoverageStatus {
	normalized := normalizeEvidenceText(requirement)
	for _, entry := range coverage.Entries {
		if entry.RequirementKind == kind && normalizeEvidenceText(entry.Requirement) == normalized {
			return entry.Status
		}
	}
	return domain.CatalogCoverageGap
}

func catalogLimitedEvidence(evidence domain.CapabilityEvidence) bool {
	return evidence.ReasonCode == "CATALOG_GAP" || evidence.ReasonCode == "HYPOTHESIS_ONLY"
}

func validateEvidenceBundle(bundle domain.EvidenceBundle) error {
	if bundle.Repository.Owner == "" || bundle.Repository.Repo == "" || bundle.Repository.Commit == "" || bundle.Repository.Commit == "HEAD" {
		return fmt.Errorf("bundle exige repository em commit imutável")
	}
	for _, item := range bundle.Items {
		if item.Ref != bundle.Repository {
			return fmt.Errorf("fonte %s não pertence ao commit do bundle", item.Path)
		}
		if digest, ok := evidenceDigest(item.ContentHash); !ok || digest != contentHash(item.Content) {
			return fmt.Errorf("hash de conteúdo inválido em %s", item.Path)
		}
	}
	return nil
}

func repositoryIdentityKey(value string) string {
	return strings.SplitN(value, "@", 2)[0]
}

func strongestReusableSource(sources []domain.EvidenceSource) (domain.EvidenceSource, bool) {
	var best domain.EvidenceSource
	found := false
	for _, source := range sources {
		if source.Kind != domain.EvidenceSymbol || source.Symbol == "" || source.Path == "" {
			continue
		}
		if !found || source.Strength > best.Strength || (source.Strength == best.Strength && source.Path < best.Path) {
			best, found = source, true
		}
	}
	return best, found
}

func findRule(capability string, catalog domain.CapabilityRuleCatalog) *domain.CapabilityRule {
	matches := exactEvidenceRules(normalizeEvidenceText(capability), catalog)
	if len(matches) != 1 {
		return nil
	}
	matched := matches[0]
	return &matched
}

func assessCapability(capability string, bundle domain.EvidenceBundle, rule *domain.CapabilityRule) domain.CapabilityEvidence {
	if rule == nil {
		result := capabilityResult(capability, domain.EvidenceInconclusive, nil, 0, "o catálogo não possui regra declarada única e segura para este requisito")
		result.ReasonCode = "CATALOG_GAP"
		return result
	}
	var sources []domain.EvidenceSource
	negative := false
	for _, item := range bundle.Items {
		norm := normalizeEvidenceText(item.Content + " " + item.Path)
		if explicitNegativeMatch(item.Content, *rule) {
			negative = true
			sources = append(sources, sourceFromItem(item, item.Kind, rule.ID, .5, firstLine(item.Content, rule.CapabilityTerms)))
		}
		for _, group := range rule.StrongEvidence {
			if strongKindMatches(item, group, norm) {
				kind := group.Kind
				if item.Kind == domain.EvidenceTest {
					kind = domain.EvidenceTestImplementation
					symbol, line, snippet, _ := extractTestMatch(item.Content, group.Patterns)
					source := sourceFromItem(item, kind, rule.ID, strengthFor(kind, item.Path), snippet)
					source.Symbol, source.Line = symbol, line
					sources = append(sources, source)
					continue
				}
				if symbol, line, snippet, ok := extractSymbolMatch(item.Content, group.Patterns); ok {
					source := sourceFromItem(item, kind, rule.ID, strengthFor(kind, item.Path), snippet)
					source.Symbol, source.Line = symbol, line
					sources = append(sources, source)
				}
			}
		}
		for _, group := range rule.SupportingEvidence {
			if matchesAny(norm, group.Patterns) {
				kind := item.Kind
				strength := supportingStrengthForMatch(item, group.Patterns)
				sources = append(sources, sourceFromItem(item, kind, rule.ID, strength, firstLine(item.Content, group.Patterns)))
			}
		}
	}
	sources = dedupeSources(sources)
	if evidenceRuleBasis(*rule) == domain.EvidenceRuleHypothesis {
		status := domain.EvidenceInconclusive
		if len(sources) > 0 {
			status = domain.EvidencePartial
		}
		result := capabilityResult(capability, status, sources, maxStrength(sources), "regra hipotética localizou somente indícios sem autoridade para concluir presença ou ausência")
		result.ReasonCode = "HYPOTHESIS_ONLY"
		return result
	}
	if negative {
		return capabilityResult(capability, domain.EvidenceConflict, sources, maxStrength(sources), "evidência negativa explícita conflita com a capacidade ou restrição")
	}
	if !bundleConclusive(bundle) {
		result := capabilityResult(capability, domain.EvidenceInconclusive, sources, maxStrength(sources), "coleta incompleta ou com falha; ausência não pode ser concluída")
		result.ReasonCode = "COLLECTION_INCOMPLETE"
		return result
	}
	// These sources locate declarations, not verified behavior. Independent lexical
	// matches cannot promote a capability to PROVEN, regardless of their sum.
	_, combined, _ := evidenceScores(sources)
	if combined > .65 {
		combined = .65
	}
	if combined >= 0.50 {
		return capabilityResult(capability, domain.EvidencePartial, sources, combined, "símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)")
	}
	if len(sources) > 0 {
		result := capabilityResult(capability, domain.EvidenceInconclusive, sources, combined, "sinais fracos impedem concluir presença ou ausência")
		result.ReasonCode = "INSUFFICIENT_EVIDENCE"
		return result
	}
	return capabilityResult(capability, domain.EvidenceAbsent, sources, combined, "busca concluída no orçamento declarado sem evidência suficiente")
}

func explicitNegativeMatch(content string, rule domain.CapabilityRule) bool {
	patterns := append([]string(nil), rule.CapabilityTerms...)
	for _, group := range append(append([]domain.CapabilityPattern{}, rule.StrongEvidence...), rule.SupportingEvidence...) {
		patterns = append(patterns, group.Patterns...)
	}
	for _, line := range strings.Split(content, "\n") {
		normalized := normalizeEvidenceText(line)
		if matchesAny(normalized, patterns) && negativeLocalOperation(normalized) {
			return true
		}
		for _, negative := range rule.NegativeEvidence {
			if phraseMatch(normalized, negative) && matchesAny(normalized, patterns) {
				return true
			}
		}
	}
	return false
}

func bundleConclusive(bundle domain.EvidenceBundle) bool {
	if !bundle.Complete || !bundle.TreeComplete || len(bundle.Failures) > 0 {
		return false
	}
	for _, target := range bundle.Targets {
		if len(target.Target.Capabilities) == 0 {
			continue
		}
		if target.Status != domain.TargetCollected && target.Status != domain.TargetNotFound {
			return false
		}
	}
	return true
}

func strongKindMatches(item domain.CollectedEvidence, group domain.CapabilityPattern, normalized string) bool {
	if !matchesAny(normalized, group.Patterns) {
		return false
	}
	switch group.Kind {
	case domain.EvidenceTest:
		if item.Kind != domain.EvidenceTest {
			return false
		}
		_, _, _, ok := extractTestMatch(item.Content, group.Patterns)
		return ok
	case domain.EvidenceSymbol:
		if item.Kind != domain.EvidenceFile && item.Kind != domain.EvidenceExample {
			return false
		}
		_, _, _, ok := extractSymbolMatch(item.Content, group.Patterns)
		return ok
	default:
		return item.Kind == group.Kind
	}
}

func extractTestMatch(content string, patterns []string) (string, int, string, bool) {
	for index, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		if !matchesAny(normalizeEvidenceText(line), patterns) {
			continue
		}
		normalized := normalizeEvidenceText(line)
		declaration := strings.HasPrefix(normalized, "func test") || strings.HasPrefix(normalized, "def test") ||
			strings.HasPrefix(normalized, "fn test") || strings.HasPrefix(normalized, "test ") || strings.HasPrefix(normalized, "it ")
		if !declaration {
			continue
		}
		fields := strings.Fields(line)
		for _, field := range fields {
			name := strings.Trim(field, "(){}:,\"")
			lower := strings.ToLower(name)
			if strings.HasPrefix(name, "Test") || strings.HasPrefix(lower, "test_") || lower == "it" {
				return name, index + 1, truncateSnippet(line), true
			}
		}
	}
	return "", 0, "", false
}

func extractSymbolMatch(content string, patterns []string) (string, int, string, bool) {
	for index, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		if !matchesAny(normalizeEvidenceText(line), patterns) {
			continue
		}
		fields := strings.Fields(line)
		for position, field := range fields {
			keyword := strings.Trim(field, "(){}:")
			if keyword != "func" && keyword != "type" && keyword != "class" && keyword != "def" && keyword != "fn" {
				continue
			}
			if position+1 >= len(fields) {
				continue
			}
			name := strings.Trim(fields[position+1], "(){}:,[<")
			if cut := strings.IndexAny(name, "([<{:"); cut >= 0 {
				name = name[:cut]
			}
			if name == "" {
				continue
			}
			return name, index + 1, truncateSnippet(line), true
		}
	}
	return "", 0, "", false
}

func truncateSnippet(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if len(value) > 240 {
		return value[:240]
	}
	return value
}

func sourceFromItem(item domain.CollectedEvidence, kind domain.EvidenceKind, rule string, strength float64, snippet string) domain.EvidenceSource {
	line := 0
	if snippet != "" {
		for index, text := range strings.Split(item.Content, "\n") {
			if truncateSnippet(text) == snippet {
				line = index + 1
				break
			}
		}
	}
	return domain.EvidenceSource{Kind: kind, Ref: item.Ref, Path: item.Path, Line: line, Snippet: snippet, RuleID: rule, ContentHash: item.ContentHash, Strength: strength}
}
func strengthFor(kind domain.EvidenceKind, path string) float64 {
	if ignoredEvidencePath(path) {
		return .05
	}
	switch kind {
	case domain.EvidenceTest, domain.EvidenceTestImplementation:
		return .6
	case domain.EvidenceSymbol:
		return .6
	case domain.EvidenceExample:
		return .5
	case domain.EvidenceManifest:
		return .7
	case domain.EvidenceDocumentation:
		return .4
	case domain.EvidenceRepository:
		return .1
	default:
		return .2
	}
}
func supportingStrength(kind domain.EvidenceKind, path string) float64 {
	strength := strengthFor(kind, path)
	if strength > .4 {
		return .4
	}
	return strength
}
func supportingStrengthForMatch(item domain.CollectedEvidence, patterns []string) float64 {
	if item.Kind == domain.EvidenceFile || item.Kind == domain.EvidenceExample {
		if matchesAny(normalizeEvidenceText(item.Path), patterns) {
			return .65
		}
	}
	return supportingStrength(item.Kind, item.Path)
}
func dedupeSources(in []domain.EvidenceSource) []domain.EvidenceSource {
	seen := map[string]bool{}
	out := make([]domain.EvidenceSource, 0, len(in))
	sort.SliceStable(in, func(i, j int) bool {
		if in[i].Path == in[j].Path {
			return in[i].Strength > in[j].Strength
		}
		return in[i].Path < in[j].Path
	})
	for _, s := range in {
		key := s.ContentHash
		if key == "" {
			key = s.Path + "|" + s.Snippet
		}
		if !seen[key] {
			seen[key] = true
			out = append(out, s)
		}
	}
	return out
}
func evidenceScores(s []domain.EvidenceSource) (float64, float64, int) {
	max, total := 0.0, 0.0
	files := map[string]bool{}
	for _, x := range s {
		if x.Strength > max {
			max = x.Strength
		}
		if !files[x.Path] {
			total = 1 - (1-total)*(1-x.Strength)
			files[x.Path] = true
		}
	}
	return max, total, len(files)
}
func maxStrength(s []domain.EvidenceSource) float64 { m, _, _ := evidenceScores(s); return m }
func capabilityResult(c string, status domain.EvidenceStatus, s []domain.EvidenceSource, confidence float64, explanation string) domain.CapabilityEvidence {
	sources := append([]domain.EvidenceSource(nil), s...)
	if sources == nil {
		sources = []domain.EvidenceSource{}
	}
	return domain.CapabilityEvidence{Capability: c, Status: status, Sources: sources, Confidence: confidence, Explanation: explanation}
}
func normalizeEvidenceText(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		r = foldEvidenceRune(r)
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}
func foldEvidenceRune(r rune) rune {
	switch r {
	case 'à', 'á', 'â', 'ã', 'ä', 'å':
		return 'a'
	case 'ç':
		return 'c'
	case 'è', 'é', 'ê', 'ë':
		return 'e'
	case 'ì', 'í', 'î', 'ï':
		return 'i'
	case 'ñ':
		return 'n'
	case 'ò', 'ó', 'ô', 'õ', 'ö':
		return 'o'
	case 'ù', 'ú', 'û', 'ü':
		return 'u'
	case 'ý', 'ÿ':
		return 'y'
	}
	return r
}
func phraseMatch(normalized, pattern string) bool {
	return strings.Contains(" "+normalized+" ", " "+normalizeEvidenceText(pattern)+" ")
}
func matchesAny(normalized string, patterns []string) bool {
	for _, p := range patterns {
		n := normalizeEvidenceText(p)
		if phraseMatch(normalized, n) || strings.Contains(strings.ReplaceAll(normalized, " ", ""), strings.ReplaceAll(n, " ", "")) {
			return true
		}
	}
	return false
}
func firstLine(content string, patterns []string) string {
	for _, line := range strings.Split(content, "\n") {
		if matchesAny(normalizeEvidenceText(line), patterns) {
			line = strings.Join(strings.Fields(line), " ")
			if len(line) > 240 {
				line = line[:240]
			}
			return line
		}
	}
	return ""
}
func assessConstraint(c string, b domain.EvidenceBundle) domain.CapabilityEvidence {
	n := normalizeEvidenceText(c)
	if strings.Contains(n, "local") || strings.Contains(n, "offline") {
		o := assessOperation(b)
		return capabilityResult(c, o.LocalStatus, o.Sources, maxStrength(o.Sources), "restrição operacional avaliada por regras locais")
	}
	return capabilityResult(c, domain.EvidenceInconclusive, nil, 0, "restrição sem regra determinística específica")
}
func assessNonGoal(nonGoal string, b domain.EvidenceBundle) domain.CapabilityEvidence {
	var sources []domain.EvidenceSource
	for _, item := range b.Items {
		for _, line := range strings.Split(item.Content, "\n") {
			if !phraseMatch(normalizeEvidenceText(line), normalizeEvidenceText(nonGoal)) {
				continue
			}
			// A mention (including a negated one) is not proof of prohibited behavior.
			sources = append(sources, sourceFromItem(item, item.Kind, "non-goal-mention", supportingStrength(item.Kind, item.Path), truncateSnippet(line)))
		}
	}
	sources = dedupeSources(sources)
	if len(sources) > 0 {
		return capabilityResult(nonGoal, domain.EvidenceInconclusive, sources, maxStrength(sources), "menção ao não objetivo, possivelmente negada; verificar comportamento antes de concluir conflito")
	}
	if !bundleConclusive(b) {
		return capabilityResult(nonGoal, domain.EvidenceInconclusive, nil, 0, "coleta incompleta")
	}
	return capabilityResult(nonGoal, domain.EvidenceAbsent, nil, 0, "non-goal não observado no escopo coletado")
}

// Only root license documents establish repository-wide authority. Nested licenses
// need a scope model before they can be used for the whole project's adoption gate.
func authoritativeLicensePath(path string) bool {
	p := strings.ToLower(slashPath(path))
	if strings.Contains(p, "/") {
		return false
	}
	switch p {
	case "license", "license.txt", "license.md", "license.rst", "copying", "copying.txt", "copying.md":
		return true
	}
	return false
}

func assessLicense(b domain.EvidenceBundle) domain.LicenseAssessment {
	for _, i := range b.Items {
		if i.Kind != domain.EvidenceLicense || !authoritativeLicensePath(i.Path) {
			continue
		}
		n := strings.ToLower(i.Content)
		id := ""
		for _, candidate := range []struct {
			id      string
			markers []string
		}{
			{"MIT", []string{"mit license", "permission is hereby granted, free of charge"}},
			{"Apache-2.0", []string{"apache license", "version 2.0"}},
			{"MPL-2.0", []string{"mozilla public license version 2.0", "mpl-2.0"}},
			{"GPL-3.0", []string{"gnu general public license", "version 3"}},
		} {
			if allMarkers(n, candidate.markers) {
				id = candidate.id
				break
			}
		}
		s := sourceFromItem(i, domain.EvidenceLicense, "license-detect", .7, firstLine(i.Content, []string{"license", "copyright"}))
		if id != "" {
			return domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: id, Sources: []domain.EvidenceSource{s}, Caveat: "identificação automática; não é aconselhamento jurídico"}
		}
	}
	if !b.Complete || len(b.Failures) > 0 {
		return domain.LicenseAssessment{Status: domain.EvidenceInconclusive, Caveat: "coleta incompleta"}
	}
	return domain.LicenseAssessment{Status: domain.EvidenceInconclusive, Caveat: "licença ausente ou ambígua"}
}
func allMarkers(content string, markers []string) bool {
	for _, marker := range markers {
		if !strings.Contains(content, marker) {
			return false
		}
	}
	return true
}
func assessMaintenance(b domain.EvidenceBundle) domain.MaintenanceAssessment {
	type maintenanceMetadata struct {
		candidate domain.RepositoryCandidate
		source    domain.EvidenceSource
	}
	recordsByKey := map[string]maintenanceMetadata{}
	var sources []domain.EvidenceSource
	var endedSources []domain.EvidenceSource
	malformed := false
	for _, i := range b.Items {
		normalized := normalizeEvidenceText(i.Content)
		if strings.Contains(normalized, "no longer maintained") || strings.Contains(normalized, "unmaintained") || strings.Contains(normalized, "maintenance ended") {
			endedSources = append(endedSources, sourceFromItem(i, i.Kind, "maintenance-ended", .7, firstLine(i.Content, []string{"maintained", "unmaintained"})))
		}
		if i.Kind != domain.EvidenceRepository {
			continue
		}
		source := sourceFromItem(i, domain.EvidenceRepository, "repository-maintenance", .7, "metadados imutáveis da coleta")
		sources = append(sources, source)
		var candidate domain.RepositoryCandidate
		if json.Unmarshal([]byte(i.Content), &candidate) != nil {
			malformed = true
			continue
		}
		key := maintenanceMetadataKey(candidate, i.Ref)
		recordsByKey[key] = maintenanceMetadata{candidate: candidate, source: source}
	}
	if len(endedSources) > 0 {
		return domain.MaintenanceAssessment{Status: domain.EvidenceConflict, Sources: canonicalMaintenanceSources(endedSources), Warnings: []string{"manutenção declarada como encerrada"}}
	}
	sources = canonicalMaintenanceSources(sources)
	if malformed {
		return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive, Sources: sources, Warnings: []string{"metadados de repositório inválidos"}}
	}
	keys := make([]string, 0, len(recordsByKey))
	for key := range recordsByKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) == 0 {
		return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive, Sources: sources, Warnings: []string{"metadados temporais não disponíveis no bundle"}}
	}
	if len(keys) > 1 {
		return domain.MaintenanceAssessment{Status: domain.EvidenceConflict, Sources: sources, Warnings: []string{"metadados de manutenção conflitantes"}}
	}
	metadata := recordsByKey[keys[0]]
	candidate := metadata.candidate
	if mutableIdeaRunCommit(b.Repository.Commit) || mutableIdeaRunCommit(metadata.source.Ref.Commit) || mutableIdeaRunCommit(candidate.CommitID) {
		return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive, Sources: sources, Warnings: []string{"identidade de manutenção não está em commit imutável"}}
	}
	if candidate.CommitID != b.Repository.Commit || metadata.source.Ref != b.Repository {
		return domain.MaintenanceAssessment{Status: domain.EvidenceConflict, Sources: sources, Warnings: []string{"metadados pertencem a commit divergente"}}
	}
	if ideaRunRepositoryKey(candidate.Owner, candidate.Repo) != ideaRunRepositoryKey(b.Repository.Owner, b.Repository.Repo) {
		return domain.MaintenanceAssessment{Status: domain.EvidenceConflict, Sources: sources, Warnings: []string{"metadados pertencem a repositório divergente"}}
	}
	if candidate.Archived {
		return domain.MaintenanceAssessment{Status: domain.EvidenceConflict, Archived: true, Sources: sources, Warnings: []string{"repositório arquivado"}}
	}
	if !b.Complete || len(b.Failures) > 0 {
		return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive, Sources: sources, Warnings: []string{"coleta incompleta"}}
	}
	if candidate.UpdatedAt.IsZero() || candidate.ConsultadoEm.IsZero() {
		return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive, Sources: sources, Warnings: []string{"metadados temporais incompletos"}}
	}
	if candidate.UpdatedAt.After(candidate.ConsultadoEm) {
		return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive, Sources: sources, Warnings: []string{"updated_at posterior ao instante de consulta"}}
	}
	age := candidate.ConsultadoEm.Sub(candidate.UpdatedAt)
	warning := "atividade observada em " + candidate.UpdatedAt.UTC().Format("2006-01-02")
	if age <= 365*24*time.Hour {
		return domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: sources, Warnings: []string{warning}}
	}
	return domain.MaintenanceAssessment{Status: domain.EvidencePartial, Sources: sources, Warnings: []string{warning, "atividade anterior à janela de 365 dias"}}
}

func maintenanceMetadataKey(candidate domain.RepositoryCandidate, ref domain.CommitRef) string {
	return strings.Join([]string{
		ideaRunRepositoryKey(ref.Owner, ref.Repo),
		ref.Commit,
		ideaRunRepositoryKey(candidate.Owner, candidate.Repo),
		candidate.CommitID,
		fmt.Sprintf("%t", candidate.Archived),
		candidate.UpdatedAt.UTC().Format(time.RFC3339Nano),
		candidate.ConsultadoEm.UTC().Format(time.RFC3339Nano),
	}, "\x00")
}

func canonicalMaintenanceSources(sources []domain.EvidenceSource) []domain.EvidenceSource {
	byKey := map[string]domain.EvidenceSource{}
	for _, source := range sources {
		encoded, _ := json.Marshal(source)
		byKey[string(encoded)] = source
	}
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]domain.EvidenceSource, 0, len(keys))
	for _, key := range keys {
		out = append(out, byKey[key])
	}
	return out
}
func assessOperation(b domain.EvidenceBundle) domain.OperationAssessment {
	var sources []domain.EvidenceSource
	var localSources []domain.EvidenceSource
	var negativeLocalSources []domain.EvidenceSource
	requirementSet := map[string]bool{}
	for _, item := range b.Items {
		for lineNumber, raw := range strings.Split(item.Content, "\n") {
			line := normalizeEvidenceText(raw)
			if negativeLocalOperation(line) {
				source := sourceFromItem(item, item.Kind, "local-operation-negative", .5, truncateSnippet(raw))
				source.Line = lineNumber + 1
				negativeLocalSources = append(negativeLocalSources, source)
			} else if matchesAny(line, []string{"run locally", "local execution", "works offline", "offline mode"}) {
				source := sourceFromItem(item, item.Kind, "local-operation", .5, truncateSnippet(raw))
				source.Line = lineNumber + 1
				localSources = append(localSources, source)
			}
			if matchesExternalCodeTransfer(line) {
				source := sourceFromItem(item, item.Kind, "external-code-upload", .7, truncateSnippet(raw))
				source.Line = lineNumber + 1
				sources = append(sources, source)
			}
		}
		normalized := normalizeEvidenceText(item.Content)
		if strings.Contains(normalized, "api token") || strings.Contains(normalized, "api key") {
			requirementSet["token ou chave de API"] = true
		}
		if strings.Contains(normalized, "postgres") || strings.Contains(normalized, "mysql") || strings.Contains(normalized, "mongodb") {
			requirementSet["banco de dados"] = true
		}
		if item.Kind == domain.EvidenceManifest {
			requirementSet["runtime declarado em manifesto"] = true
		}

	}
	if len(negativeLocalSources) > 0 {
		return domain.OperationAssessment{LocalStatus: domain.EvidenceConflict, Sources: dedupeSources(negativeLocalSources), Requirements: sortedRequirements(requirementSet)}
	}
	if len(sources) > 0 {
		requirementSet["serviço externo ou envio de código"] = true
		return domain.OperationAssessment{LocalStatus: domain.EvidenceConflict, Sources: sources, Requirements: sortedRequirements(requirementSet)}
	}
	if !b.Complete || len(b.Failures) > 0 {
		return domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive}
	}
	if len(localSources) > 0 {
		return domain.OperationAssessment{LocalStatus: domain.EvidencePartial, Sources: dedupeSources(localSources), Requirements: sortedRequirements(requirementSet)}
	}
	return domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive, Requirements: sortedRequirements(requirementSet)}
}

func sortedRequirements(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func matchesExternalCodeTransfer(line string) bool {
	transfer := strings.Contains(line, "upload") || strings.Contains(line, "send code") || strings.Contains(line, "sends code")
	code := strings.Contains(line, "repository") || strings.Contains(line, "source code") || strings.Contains(line, "codebase") || strings.Contains(line, "your code")
	destination := strings.Contains(line, " api") || strings.Contains(line, "cloud") || strings.Contains(line, "server") || strings.Contains(line, "service")
	return transfer && code && destination && !negatedExternalTransfer(line)
}

func EvaluateEvidenceMetrics(report domain.AssessmentReport, bundles []domain.EvidenceBundle, judgments []domain.EvidenceJudgment) domain.EvidenceMetrics {
	metrics := domain.EvidenceMetrics{}
	judgmentByKey := map[string]domain.EvidenceJudgment{}
	for _, judgment := range judgments {
		judgmentByKey[repositoryIdentityKey(judgment.Repository)+"|"+normalizeEvidenceText(judgment.Capability)] = judgment
	}
	provenReviewed, provenCorrect := 0, 0
	absentReviewed, falseAbsent := 0, 0
	total, inconclusive, sources, traced := 0, 0, 0, 0
	componentsReviewed, componentsCorrect := 0, 0
	for _, assessment := range report.Assessments {
		repo := assessment.Repository.Owner + "/" + assessment.Repository.Repo
		evidenceSet := make([]domain.CapabilityEvidence, 0, len(assessment.Capabilities)+1)
		if assessment.ProblemAlignment.Capability != "" {
			evidenceSet = append(evidenceSet, assessment.ProblemAlignment)
		}
		evidenceSet = append(evidenceSet, assessment.Capabilities...)
		for _, evidence := range evidenceSet {
			total++
			if evidence.Status == domain.EvidenceInconclusive {
				inconclusive++
			}
			if evidence.Status == domain.EvidenceProven {
				metrics.ProvenCount++
			}
			for _, source := range evidence.Sources {
				sources++
				if source.Ref.Owner != "" && source.Ref.Repo != "" && source.Ref.Commit != "" && source.Ref.Commit != "HEAD" && source.Path != "" && source.ContentHash != "" {
					traced++
				}
			}
			if judgment, ok := judgmentByKey[repo+"|"+normalizeEvidenceText(evidence.Capability)]; ok {
				if evidence.Status == domain.EvidenceProven {
					provenReviewed++
					if judgment.Expected == domain.EvidenceProven {
						provenCorrect++
					}
				}
				if evidence.Status == domain.EvidenceAbsent {
					absentReviewed++
					if judgment.Expected != domain.EvidenceAbsent {
						falseAbsent++
					}
				}
				if judgment.ComponentValid != nil && assessmentHasComponent(assessment, evidence.Capability) {
					componentsReviewed++
					if *judgment.ComponentValid {
						componentsCorrect++
					}
				}
			}
		}
	}
	if provenReviewed > 0 {
		value := float64(provenCorrect) / float64(provenReviewed)
		metrics.EvidencePrecision = &value
	}
	if absentReviewed > 0 {
		value := float64(falseAbsent) / float64(absentReviewed)
		metrics.FalseAbsenceRate = &value
	}
	if total > 0 {
		metrics.InconclusiveRate = float64(inconclusive) / float64(total)
	}
	if sources > 0 {
		metrics.SourceTraceRate = float64(traced) / float64(sources)
	}
	planned, completed := 0, 0
	for _, bundle := range bundles {
		for _, target := range bundle.Targets {
			planned++
			if target.Status == domain.TargetCollected || target.Status == domain.TargetNotFound {
				completed++
			}
		}
	}
	if planned > 0 {
		metrics.BudgetCompletion = float64(completed) / float64(planned)
	}
	if componentsReviewed > 0 {
		value := float64(componentsCorrect) / float64(componentsReviewed)
		metrics.ComponentPrecision = &value
	}
	metrics.SampleSufficient = metrics.ProvenCount >= 10
	return metrics
}

func assessmentHasComponent(assessment domain.CandidateAssessment, capability string) bool {
	for _, component := range assessment.ReusableParts {
		for _, covered := range component.Capabilities {
			if normalizeEvidenceText(covered) == normalizeEvidenceText(capability) {
				return true
			}
		}
	}
	return false
}
