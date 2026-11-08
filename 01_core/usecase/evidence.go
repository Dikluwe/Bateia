/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence.md
 * @layer L1
 * @updated 2026-11-08
 */

package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	pathpkg "path"
	"sort"
	"strings"
	"time"
	"unicode"

	"bateia/01_core/domain"
)

var defaultEvidenceBudget = domain.ExtractionBudget{
	Documents: 8, Manifests: 4, CodeFiles: 20, RequirementFiles: 2, Tests: 12,
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
	coverage, err := AnalyzeEvidenceCatalogCoverage(idea, catalog)
	if err != nil {
		return nil, err
	}
	requirementQueries := compileEvidenceRequirementQueries(coverage, catalog)
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
		plans = append(plans, domain.EvidenceExtractionPlan{SchemaVersion: "1", Repository: ref, Candidate: candidate, Decision: decision, Targets: targets, Budget: defaultEvidenceBudget, RequirementQueries: cloneRequirementQueries(requirementQueries)})
	}
	return plans, nil
}

func compileEvidenceRequirementQueries(coverage domain.EvidenceCatalogCoverage, catalog domain.CapabilityRuleCatalog) []domain.EvidenceRequirementQuery {
	rules := make(map[string]domain.CapabilityRule, len(catalog.Rules))
	for _, rule := range catalog.Rules {
		rules[rule.ID] = rule
	}
	queries := []domain.EvidenceRequirementQuery{}
	for _, entry := range coverage.Entries {
		if entry.Status != domain.CatalogCoverageCovered || entry.RuleBasis != domain.EvidenceRuleDeclared {
			continue
		}
		rule, ok := rules[entry.RuleID]
		if !ok || evidenceRuleBasis(rule) != domain.EvidenceRuleDeclared {
			continue
		}
		terms := evidenceRequirementTerms(rule)
		if len(terms) == 0 {
			continue
		}
		queries = append(queries, domain.EvidenceRequirementQuery{RequirementKind: entry.RequirementKind, Requirement: entry.Requirement, RuleID: entry.RuleID, Terms: terms})
	}
	sortRequirementQueries(queries)
	return queries
}

func evidenceRequirementTerms(rule domain.CapabilityRule) []string {
	set := map[string]string{}
	add := func(raw string, hint bool) {
		for _, term := range technicalQueryTerms(raw, hint) {
			key := strings.ToLower(term)
			if _, exists := set[key]; !exists {
				set[key] = term
			}
		}
	}
	for _, value := range rule.TargetHints {
		add(value, true)
		// Preserve the step-0072 C++ header convention without borrowing a
		// proof pattern: an identifier-like header hint may also name its .h
		// spelling. Both terms are still derived solely from TargetHints.
		for _, term := range technicalQueryTerms(value, true) {
			if pathpkg.Ext(term) == "" && strings.Contains(term, "_") {
				add(term+".h", true)
			}
		}
	}
	out := make([]string, 0, len(set))
	for _, value := range set {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool {
		left, right := strings.ToLower(out[i]), strings.ToLower(out[j])
		if left != right {
			return left < right
		}
		return out[i] < out[j]
	})
	return out
}

func technicalQueryTerms(raw string, hint bool) []string {
	value := strings.TrimSpace(strings.ToLower(slashPath(raw)))
	if value == "" || identityDependentRuleText(value) {
		return nil
	}
	if strings.Contains(value, "/") {
		value = pathpkg.Base(value)
	}
	if strings.ContainsAny(value, " \t\n") {
		return nil
	}
	base := strings.Trim(value, "<>[](){}\"'`,;:*=+&|!?")
	if extension := pathpkg.Ext(base); recognizedCPPQueryExtension(extension) {
		return []string{base}
	}
	if index := strings.LastIndex(base, "::"); index >= 0 {
		base = base[index+2:]
	}
	if index := strings.Index(base, "<"); index >= 0 {
		base = base[:index]
	}
	base = strings.Trim(base, "<>[](){}\"'`,;:*=+&|!?")
	if base == "" || overlyGenericRuleText(normalizeEvidenceText(base)) {
		return nil
	}
	technical := hint || strings.ContainsAny(raw, "_</.:")
	if !technical && len(base) < 4 {
		return nil
	}
	for _, r := range base {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '.' && r != '-' {
			return nil
		}
	}
	return []string{base}
}

func recognizedCPPQueryExtension(extension string) bool {
	switch strings.ToLower(extension) {
	case ".h", ".hh", ".hpp", ".hxx", ".c", ".cc", ".cpp", ".cxx":
		return true
	default:
		return false
	}
}

func sortRequirementQueries(queries []domain.EvidenceRequirementQuery) {
	for i := range queries {
		queries[i].Terms = uniqueSortedFolded(queries[i].Terms)
	}
	sort.Slice(queries, func(i, j int) bool {
		left, right := queries[i], queries[j]
		if left.RequirementKind != right.RequirementKind {
			return left.RequirementKind < right.RequirementKind
		}
		ln, rn := normalizeEvidenceText(left.Requirement), normalizeEvidenceText(right.Requirement)
		if ln != rn {
			return ln < rn
		}
		if left.Requirement != right.Requirement {
			return left.Requirement < right.Requirement
		}
		return left.RuleID < right.RuleID
	})
}

func cloneRequirementQueries(input []domain.EvidenceRequirementQuery) []domain.EvidenceRequirementQuery {
	out := make([]domain.EvidenceRequirementQuery, len(input))
	for i, query := range input {
		out[i] = query
		out[i].Terms = append([]string(nil), query.Terms...)
	}
	return out
}

func uniqueSortedFolded(values []string) []string {
	set := map[string]string{}
	for _, value := range values {
		value = strings.TrimSpace(strings.ToLower(value))
		if value != "" {
			set[value] = value
		}
	}
	out := make([]string, 0, len(set))
	for _, value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
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

func validateEvidenceCollectionPlan(plan domain.EvidenceExtractionPlan, requireSchema bool) error {
	if plan.SchemaVersion != "1" && (requireSchema || plan.SchemaVersion != "") {
		return domain.ErrInvalidRunInput
	}
	ref := plan.Repository
	for _, part := range []string{ref.Owner, ref.Repo} {
		if !validSynthesisText(part, 100) || part != normalizeMechanismText(part) || strings.ContainsAny(part, "/\\@ :?#") || part == "." || part == ".." {
			return domain.ErrInvalidRunInput
		}
	}
	if plan.Candidate.Owner != ref.Owner || plan.Candidate.Repo != ref.Repo || plan.Budget.TotalBytes <= 0 {
		return domain.ErrInvalidRunInput
	}
	for _, quota := range []int{plan.Budget.Documents, plan.Budget.Manifests, plan.Budget.CodeFiles, plan.Budget.RequirementFiles, plan.Budget.Tests, plan.Budget.TreeDepth} {
		if quota < 0 {
			return domain.ErrInvalidRunInput
		}
	}
	return nil
}

func CollectEvidence(ctx context.Context, plan domain.EvidenceExtractionPlan, provider domain.ContentProvider) (domain.EvidenceBundle, error) {
	if ctx == nil {
		return domain.EvidenceBundle{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.EvidenceBundle{}, err
	}
	if err := validateEvidenceCollectionPlan(plan, false); err != nil {
		return domain.EvidenceBundle{}, err
	}
	if provider == nil {
		return domain.EvidenceBundle{}, fmt.Errorf("provedor de conteúdo não configurado")
	}
	resolving := plan.Repository.Commit == "" || strings.EqualFold(plan.Repository.Commit, "HEAD")
	if !resolving {
		plan.Candidate.CommitID = plan.Repository.Commit
	}
	metadata, metadataErr := provider.FetchRepoMetadata(ctx, plan.Candidate)
	if err := ctx.Err(); err != nil {
		return domain.EvidenceBundle{}, err
	}
	if errors.Is(metadataErr, domain.ErrIntegrity) {
		return domain.EvidenceBundle{}, domain.ErrIntegrity
	}
	if metadataErr == nil {
		if metadata.Owner != plan.Repository.Owner || metadata.Repo != plan.Repository.Repo {
			return domain.EvidenceBundle{}, domain.ErrIntegrity
		}
		if resolving {
			if metadata.CommitID == "" || strings.EqualFold(metadata.CommitID, "HEAD") {
				return domain.EvidenceBundle{}, fmt.Errorf("não foi possível resolver %s para commit imutável", plan.Repository.Key())
			}
			plan.Repository.Commit = metadata.CommitID
			plan.Decision.Repository = plan.Repository.Key()
		}
		metadata.CommitID = plan.Repository.Commit
		plan.Candidate = metadata
	} else if resolving {
		return domain.EvidenceBundle{}, fmt.Errorf("não foi possível resolver %s para commit imutável", plan.Repository.Key())
	}
	plan.RequirementQueries = cloneRequirementQueries(plan.RequirementQueries)
	sortRequirementQueries(plan.RequirementQueries)
	bundle := domain.EvidenceBundle{
		SchemaVersion: "1", Repository: plan.Repository, Plan: plan, Complete: true, TreeComplete: true,
		Items: []domain.CollectedEvidence{}, Targets: []domain.CollectionTargetResult{}, RequirementCoverage: []domain.RequirementCollectionCoverage{},
	}
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
	selectedPaths := map[string]bool{}
	for _, path := range paths {
		selectedPaths[path] = true
	}
	implementationQuery := ""
	var requirementSelections []requirementPathSelection
	tree, err := provider.FetchTree(ctx, plan.Candidate)
	if err != nil {
		bundle.Failures = append(bundle.Failures, domain.ExecutionFailure{Scope: plan.Repository.Key(), Stage: "tree", Message: err.Error()})
		bundle.Complete = false
		bundle.TreeComplete = false
		requirementSelections = buildRequirementPathSelections(nil, plan.RequirementQueries)
	} else {
		tree = canonicalEvidenceTree(tree)
		requirementSelections = buildRequirementPathSelections(tree, plan.RequirementQueries)
		reserveRequirementPaths(requirementSelections, plan.Budget, &paths, kinds, targetByPath, selectedPaths)
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
			if selectedPaths[node.Path] {
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
			selectedPaths[node.Path] = true
		}
	}
	if err == nil {
		for _, selection := range requirementSelections {
			for _, path := range selection.candidates {
				if !selectedPaths[path] {
					bundle.Complete = false
				}
			}
		}
	}
	if err := ctx.Err(); err != nil {
		bundle.Complete = false
		return bundle, err
	}
	extracted, err := provider.ExtractRepoFiles(ctx, plan.Candidate, paths)
	if contextErr := ctx.Err(); contextErr != nil {
		return domain.EvidenceBundle{}, contextErr
	}
	if errors.Is(err, domain.ErrIntegrity) {
		return domain.EvidenceBundle{}, domain.ErrIntegrity
	}
	if err == nil && (extracted == nil || extracted.Ref != plan.Repository) {
		return domain.EvidenceBundle{}, domain.ErrIntegrity
	}
	if err != nil {
		bundle.Failures = append(bundle.Failures, domain.ExecutionFailure{Scope: plan.Repository.Key(), Stage: "content", Message: err.Error()})
		bundle.Complete = false
		for _, path := range paths {
			bundle.Targets = append(bundle.Targets, domain.CollectionTargetResult{Target: targetByPath[path], Status: domain.TargetFailed, Message: err.Error()})
		}
		bundle.RequirementCoverage = finalizeRequirementCoverage(requirementSelections, bundle.TreeComplete, nil)
		return bundle, nil
	}
	filesByPath := map[string]domain.ExtractedFile{}
	for _, file := range extracted.Files {
		path := slashPath(file.Path)
		if _, exists := filesByPath[path]; !exists {
			file.Path = path
			filesByPath[path] = file
		}
	}
	extractionPaths := append([]string(nil), paths...)
	extraPaths := []string{}
	for path := range filesByPath {
		if !selectedPaths[path] {
			extraPaths = append(extraPaths, path)
		}
	}
	sort.Strings(extraPaths)
	extractionPaths = append(extractionPaths, extraPaths...)
	collectedPaths := map[string]int64{}
	fullyExamined := map[string]bool{}
	seenContent := map[string]bool{}
	for _, path := range extractionPaths {
		file, ok := filesByPath[slashPath(path)]
		if !ok {
			continue
		}
		content := file.Content
		remaining := plan.Budget.TotalBytes - bundle.BytesRead
		if remaining <= 0 {
			bundle.Complete = false
			continue
		}
		completeContent := true
		if int64(len(content)) > remaining {
			content = content[:remaining]
			bundle.Complete = false
			completeContent = false
		}
		hash := contentHash(content)
		if seenContent[hash] {
			fullyExamined[path] = completeContent
			collectedPaths[path] = 0
			continue
		}
		seenContent[hash] = true
		kind := kinds[path]
		if kind == "" {
			kind = classifyEvidencePath(path)
		}
		bundle.Items = append(bundle.Items, domain.CollectedEvidence{Kind: kind, Ref: plan.Repository, Path: path, Content: content, ContentHash: hash})
		bundle.BytesRead += int64(len(content))
		collectedPaths[path] = int64(len(content))
		fullyExamined[path] = completeContent
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
	bundle.RequirementCoverage = finalizeRequirementCoverage(requirementSelections, bundle.TreeComplete, fullyExamined)
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

type requirementPathSelection struct {
	query      domain.EvidenceRequirementQuery
	candidates []string
	selected   map[string]bool
}

func canonicalEvidenceTree(tree []domain.DirectoryNode) []domain.DirectoryNode {
	byPath := map[string]domain.DirectoryNode{}
	for _, node := range tree {
		node.Path = slashPath(strings.TrimPrefix(node.Path, "./"))
		if node.Path == "" {
			continue
		}
		if previous, exists := byPath[node.Path]; !exists || (previous.Type != "file" && node.Type == "file") {
			byPath[node.Path] = node
		}
	}
	out := make([]domain.DirectoryNode, 0, len(byPath))
	for _, node := range byPath {
		out = append(out, node)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Type < out[j].Type
	})
	return out
}

func buildRequirementPathSelections(tree []domain.DirectoryNode, queries []domain.EvidenceRequirementQuery) []requirementPathSelection {
	selections := make([]requirementPathSelection, 0, len(queries))
	for _, query := range queries {
		selection := requirementPathSelection{query: query, candidates: []string{}, selected: map[string]bool{}}
		type rankedPath struct {
			path string
			rank int
		}
		ranked := []rankedPath{}
		for _, node := range tree {
			if node.Type != "file" {
				continue
			}
			rank, matched := requirementPathRank(node.Path, query.Terms)
			if !matched || !isRequirementCodePath(node.Path, query.Terms) {
				continue
			}
			ranked = append(ranked, rankedPath{path: node.Path, rank: rank})
		}
		sort.Slice(ranked, func(i, j int) bool {
			if ranked[i].rank != ranked[j].rank {
				return ranked[i].rank < ranked[j].rank
			}
			leftDepth, rightDepth := pathDepth(ranked[i].path), pathDepth(ranked[j].path)
			if leftDepth != rightDepth {
				return leftDepth < rightDepth
			}
			return ranked[i].path < ranked[j].path
		})
		for _, candidate := range ranked {
			selection.candidates = append(selection.candidates, candidate.path)
		}
		selections = append(selections, selection)
	}
	return selections
}

func reserveRequirementPaths(selections []requirementPathSelection, budget domain.ExtractionBudget, paths *[]string, kinds map[string]domain.EvidenceKind, targets map[string]domain.ExtractionTarget, selected map[string]bool) {
	quota := budget.RequirementFiles
	if quota <= 0 {
		quota = 2
	}
	codeCount := 0
	for _, path := range *paths {
		if kinds[path] == domain.EvidenceFile && isCodePath(path) {
			codeCount++
		}
	}
	for round := 0; round < quota; round++ {
		for index := range selections {
			selection := &selections[index]
			for _, path := range selection.candidates {
				if selection.selected[path] {
					continue
				}
				if pathDepth(path) > budget.TreeDepth {
					continue
				}
				if selected[path] {
					selection.selected[path] = true
					break
				}
				if codeCount >= budget.CodeFiles {
					break
				}
				kind := classifyEvidencePath(path)
				*paths = append(*paths, path)
				kinds[path] = kind
				targets[path] = domain.ExtractionTarget{Path: path, Kind: kind, Capabilities: []string{selection.query.Requirement}, Reason: "recuperação dirigida pela regra " + selection.query.RuleID}
				selected[path] = true
				selection.selected[path] = true
				codeCount++
				break
			}
		}
	}
}

func requirementPathRank(path string, terms []string) (int, bool) {
	base := strings.ToLower(pathpkg.Base(slashPath(path)))
	stem := strings.TrimSuffix(base, pathpkg.Ext(base))
	tokens := identifierTokens(stem)
	for _, term := range terms {
		term = strings.ToLower(strings.TrimSpace(pathpkg.Base(slashPath(term))))
		if term != "" && base == term {
			return 0, true
		}
	}
	for _, term := range terms {
		term = strings.ToLower(strings.TrimSpace(pathpkg.Base(slashPath(term))))
		termStem := strings.TrimSuffix(term, pathpkg.Ext(term))
		if termStem != "" && stem == termStem {
			return 1, true
		}
	}
	for _, term := range terms {
		term = strings.ToLower(strings.TrimSpace(pathpkg.Base(slashPath(term))))
		term = strings.TrimSuffix(term, pathpkg.Ext(term))
		for _, token := range tokens {
			if token == term && term != "" {
				return 2, true
			}
		}
	}
	return 0, false
}

func identifierTokens(value string) []string {
	return strings.FieldsFunc(strings.ToLower(value), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
	})
}

func isRequirementCodePath(path string, terms []string) bool {
	if isCodePath(path) {
		return true
	}
	base := strings.ToLower(pathpkg.Base(slashPath(path)))
	if pathpkg.Ext(base) != "" {
		return false
	}
	for _, term := range terms {
		termBase := strings.ToLower(strings.TrimSpace(pathpkg.Base(slashPath(term))))
		if pathpkg.Ext(termBase) == "" && base == termBase {
			return true
		}
	}
	return false
}

func finalizeRequirementCoverage(selections []requirementPathSelection, searchComplete bool, examined map[string]bool) []domain.RequirementCollectionCoverage {
	out := make([]domain.RequirementCollectionCoverage, 0, len(selections))
	for _, selection := range selections {
		coverage := domain.RequirementCollectionCoverage{
			RequirementKind: selection.query.RequirementKind, Requirement: selection.query.Requirement, RuleID: selection.query.RuleID,
			SearchComplete: searchComplete, CandidatePaths: append([]string{}, selection.candidates...), ExaminedPaths: []string{}, SkippedPaths: []string{},
		}
		for _, path := range selection.candidates {
			if examined[path] {
				coverage.ExaminedPaths = append(coverage.ExaminedPaths, path)
			} else {
				coverage.SkippedPaths = append(coverage.SkippedPaths, path)
			}
		}
		sort.Strings(coverage.CandidatePaths)
		sort.Strings(coverage.ExaminedPaths)
		sort.Strings(coverage.SkippedPaths)
		switch {
		case !searchComplete || len(coverage.SkippedPaths) > 0:
			coverage.Status = domain.RequirementCollectionIncomplete
			coverage.ReasonCode = "COLLECTION_INCOMPLETE"
		case len(coverage.CandidatePaths) == 0:
			coverage.Status = domain.RequirementNoPlausiblePath
			coverage.ReasonCode = "NO_PLAUSIBLE_PATH"
		default:
			coverage.Status = domain.RequirementTargetedExamined
			coverage.ReasonCode = "TARGETED_EXAMINED"
		}
		out = append(out, coverage)
	}
	return out
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
	case ".go", ".rs", ".py", ".js", ".jsx", ".ts", ".tsx", ".java", ".kt", ".rb", ".cs", ".c", ".cc", ".cpp", ".cxx", ".h", ".hh", ".hpp", ".hxx", ".ex", ".exs":
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
			if catalogCoverageStatus(coverage, "capability", capability) == domain.CatalogCoverageCovered && e.Status != domain.EvidenceProven && e.ReasonCode != "COLLECTION_INCOMPLETE" {
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
	lexicalOpaque := false
	for _, item := range bundle.Items {
		itemSourceCount := len(sources)
		if explicitNegativeMatch(item.Content, *rule) {
			negative = true
			if match, ok := firstAdmissibleContentMatch(item, rule.CapabilityTerms, domain.EvidenceMatchSupportingContent, false); ok {
				sources = append(sources, sourceFromMatch(item, item.Kind, rule.ID, .5, match))
			}
		}
		strongSourceCount := len(sources)
		for _, group := range rule.StrongEvidence {
			for _, match := range locateRuleContentMatches(item, group, domain.EvidenceMatchStrongContent) {
				kind := group.Kind
				if item.Kind == domain.EvidenceTest {
					kind = domain.EvidenceTestImplementation
				}
				sources = append(sources, sourceFromMatch(item, kind, rule.ID, strengthFor(kind, item.Path), match))
			}
		}
		strongLocated := len(sources) > strongSourceCount
		pendingSupporting := []domain.EvidenceSource{}
		for _, group := range rule.SupportingEvidence {
			for _, match := range locateRuleContentMatches(item, group, domain.EvidenceMatchSupportingContent) {
				strength := supportingStrength(item.Kind, item.Path)
				pendingSupporting = append(pendingSupporting, sourceFromMatch(item, item.Kind, rule.ID, strength, match))
			}
		}
		// Supporting tokens in source code (for example a common type name)
		// only corroborate a substantive match in that file. Documentation is
		// independently useful, but retains its lower authority class.
		if strongLocated || item.Kind == domain.EvidenceDocumentation {
			sources = append(sources, pendingSupporting...)
		}
		if len(sources) == itemSourceCount && opaqueRuleSignal(item, *rule) {
			lexicalOpaque = true
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
	if requirementCollectionIncomplete(bundle, capability, rule.ID) && len(sources) == 0 {
		result := capabilityResult(capability, domain.EvidenceInconclusive, nil, 0, "coleta dirigida incompleta; caminhos plausíveis permaneceram sem exame")
		result.ReasonCode = "COLLECTION_INCOMPLETE"
		return result
	}
	targetedExamined := requirementTargetedExamined(bundle, capability, rule.ID)
	if !bundleConclusive(bundle) && (!targetedExamined || len(sources) == 0) {
		result := capabilityResult(capability, domain.EvidenceInconclusive, sources, maxStrength(sources), "coleta incompleta ou com falha; ausência não pode ser concluída")
		result.ReasonCode = "COLLECTION_INCOMPLETE"
		return result
	}
	if len(sources) == 0 && lexicalOpaque {
		result := capabilityResult(capability, domain.EvidenceInconclusive, nil, 0, "o arquivo ou bytes recuperados contêm apenas sinais que o localizador lexical não pode atribuir com segurança")
		result.ReasonCode = "LEXICAL_MATCH_OPAQUE"
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

func requirementCollectionIncomplete(bundle domain.EvidenceBundle, requirement, ruleID string) bool {
	normalized := normalizeEvidenceText(requirement)
	for _, coverage := range bundle.RequirementCoverage {
		if coverage.RuleID == ruleID && normalizeEvidenceText(coverage.Requirement) == normalized {
			return coverage.Status == domain.RequirementCollectionIncomplete || !coverage.SearchComplete || len(coverage.SkippedPaths) > 0
		}
	}
	return false
}

func requirementTargetedExamined(bundle domain.EvidenceBundle, requirement, ruleID string) bool {
	normalized := normalizeEvidenceText(requirement)
	for _, coverage := range bundle.RequirementCoverage {
		if coverage.RuleID == ruleID && normalizeEvidenceText(coverage.Requirement) == normalized {
			return coverage.Status == domain.RequirementTargetedExamined && coverage.SearchComplete && len(coverage.SkippedPaths) == 0
		}
	}
	return false
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

type locatedEvidenceMatch struct {
	Pattern string
	Symbol  string
	Line    int
	Snippet string
	Class   domain.EvidenceMatchClass
}

func locateRuleContentMatches(item domain.CollectedEvidence, group domain.CapabilityPattern, class domain.EvidenceMatchClass) []locatedEvidenceMatch {
	switch group.Kind {
	case domain.EvidenceTest:
		if item.Kind != domain.EvidenceTest {
			return nil
		}
	case domain.EvidenceSymbol:
		if item.Kind != domain.EvidenceFile && item.Kind != domain.EvidenceExample && item.Kind != domain.EvidenceTest {
			return nil
		}
	default:
		if item.Kind != group.Kind {
			return nil
		}
	}
	if item.Kind == domain.EvidenceDocumentation {
		class = domain.EvidenceMatchDocumentation
	}

	rawLines := strings.Split(item.Content, "\n")
	codeLines := lexicalCodeLines(item.Content)
	matches := make([]locatedEvidenceMatch, 0)
	patterns := canonicalRuleStrings(group.Patterns)
	for index, raw := range rawLines {
		searchable := codeLines[index]
		if class == domain.EvidenceMatchDocumentation {
			searchable = raw
		}
		if strings.HasPrefix(strings.TrimSpace(searchable), "#") && class != domain.EvidenceMatchDocumentation {
			continue
		}
		for _, pattern := range patterns {
			if !exactLexicalPatternMatch(searchable, pattern) {
				continue
			}
			symbol := ""
			if group.Kind == domain.EvidenceSymbol || group.Kind == domain.EvidenceTest {
				if !supportedSymbolOccurrence(searchable, pattern, group.Kind == domain.EvidenceTest) {
					continue
				}
				symbol = declaredOrMatchedSymbol(searchable, pattern)
			}
			matches = append(matches, locatedEvidenceMatch{
				Pattern: pattern, Symbol: symbol, Line: index + 1,
				Snippet: truncateSnippet(raw), Class: class,
			})
		}
	}
	return matches
}

func firstAdmissibleContentMatch(item domain.CollectedEvidence, patterns []string, class domain.EvidenceMatchClass, symbol bool) (locatedEvidenceMatch, bool) {
	kind := domain.EvidenceFile
	if symbol {
		kind = domain.EvidenceSymbol
	}
	if item.Kind == domain.EvidenceDocumentation {
		kind = domain.EvidenceDocumentation
	}
	matches := locateRuleContentMatches(item, domain.CapabilityPattern{Kind: kind, Patterns: patterns}, class)
	if len(matches) == 0 {
		return locatedEvidenceMatch{}, false
	}
	return matches[0], true
}

// lexicalCodeLines replaces comments and quoted literals with spaces while
// preserving line boundaries and byte positions. It deliberately does not
// interpret macros, token pasting or declarations split across lines.
func lexicalCodeLines(content string) []string {
	lines := strings.Split(content, "\n")
	out := make([]string, len(lines))
	inBlockComment := false
	for lineIndex, raw := range lines {
		bytes := []byte(raw)
		clean := make([]byte, len(bytes))
		for i := range clean {
			clean[i] = ' '
		}
		for i := 0; i < len(bytes); {
			if inBlockComment {
				if i+1 < len(bytes) && bytes[i] == '*' && bytes[i+1] == '/' {
					inBlockComment = false
					i += 2
				} else {
					i++
				}
				continue
			}
			if i+1 < len(bytes) && bytes[i] == '/' && bytes[i+1] == '/' {
				break
			}
			if i+1 < len(bytes) && bytes[i] == '/' && bytes[i+1] == '*' {
				inBlockComment = true
				i += 2
				continue
			}
			if bytes[i] == '\'' || bytes[i] == '"' {
				quote := bytes[i]
				i++
				for i < len(bytes) {
					if bytes[i] == '\\' {
						i += 2
						continue
					}
					if bytes[i] == quote {
						i++
						break
					}
					i++
				}
				continue
			}
			clean[i] = bytes[i]
			i++
		}
		out[lineIndex] = string(clean)
	}
	return out
}

func exactLexicalPatternMatch(line, pattern string) bool {
	line = strings.ToLower(strings.Join(strings.Fields(line), " "))
	pattern = strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(pattern)), " "))
	if line == "" || pattern == "" {
		return false
	}
	if lexicalSubstringWithBoundaries(line, pattern) {
		return true
	}
	// Preserve the legacy ability to match a multiword catalog term against a
	// CamelCase identifier, while retaining whole-identifier boundaries.
	if strings.Contains(pattern, " ") {
		joined := strings.ReplaceAll(pattern, " ", "")
		return strings.Contains(strings.ReplaceAll(line, " ", ""), joined)
	}
	return false
}

func lexicalSubstringWithBoundaries(line, pattern string) bool {
	for start := 0; ; {
		position := strings.Index(line[start:], pattern)
		if position < 0 {
			return false
		}
		position += start
		end := position + len(pattern)
		leftOK := !isIdentifierByte(pattern[0]) || position == 0 || !isIdentifierByte(line[position-1])
		rightOK := !isIdentifierByte(pattern[len(pattern)-1]) || end == len(line) || !isIdentifierByte(line[end])
		if leftOK && rightOK {
			return true
		}
		start = position + 1
		if start >= len(line) {
			return false
		}
	}
}

func isIdentifierByte(value byte) bool {
	return value == '_' || value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z' || value >= '0' && value <= '9'
}

func supportedSymbolOccurrence(line, pattern string, test bool) bool {
	compact := strings.Join(strings.Fields(line), " ")
	if strings.HasPrefix(strings.TrimSpace(compact), "#") {
		return false
	}
	if test {
		normalized := normalizeEvidenceText(compact)
		return strings.HasPrefix(normalized, "func test") || strings.HasPrefix(normalized, "def test") ||
			strings.HasPrefix(normalized, "fn test") || strings.HasPrefix(normalized, "test ") || strings.HasPrefix(normalized, "it ")
	}
	if _, ok := declaredSymbol(compact); ok {
		return true
	}
	// Some catalog symbols are deliberately operator-rich code idioms rather
	// than a single declaration or call (for example, a self-assignment guard).
	// The caller has already established an exact lexical match on comment- and
	// literal-free code, so this branch only classifies the declared pattern.
	if supportedCompoundCodeExpression(pattern) {
		return true
	}
	identifier := symbolForPattern(pattern)
	if identifier == "" || !exactLexicalPatternMatch(compact, identifier) {
		return false
	}
	searchCompact := strings.ToLower(compact)
	searchIdentifier := strings.ToLower(identifier)
	for start := 0; ; {
		position := strings.Index(searchCompact[start:], searchIdentifier)
		if position < 0 {
			break
		}
		position += start
		end := position + len(searchIdentifier)
		if (position == 0 || !isIdentifierByte(searchCompact[position-1])) && (end == len(searchCompact) || !isIdentifierByte(searchCompact[end])) {
			rest := strings.TrimSpace(compact[end:])
			if strings.HasPrefix(rest, "(") || strings.HasPrefix(rest, "<") {
				return true
			}
			prefix := normalizeEvidenceText(compact[:position])
			for _, keyword := range []string{"func", "type", "class", "struct", "def", "fn"} {
				if phraseMatch(prefix, keyword) {
					return true
				}
			}
		}
		start = position + 1
	}
	return false
}

func supportedCompoundCodeExpression(pattern string) bool {
	if strings.ContainsAny(pattern, "\r\n") {
		return false
	}
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return false
	}
	identifiers := 0
	for index := 0; index < len(pattern); {
		if !isIdentifierByte(pattern[index]) {
			index++
			continue
		}
		identifiers++
		for index < len(pattern) && isIdentifierByte(pattern[index]) {
			index++
		}
	}
	if identifiers < 2 {
		return false
	}
	for _, operator := range []string{"==", "!=", "<=", ">=", "&&", "||", "->"} {
		if strings.Contains(pattern, operator) {
			return true
		}
	}
	for index := 0; index < len(pattern); index++ {
		switch pattern[index] {
		case '=', '%', '+', '*', '/', '&', '|', '^', '<', '>':
			return true
		case '-':
			leftIdentifier := index > 0 && isIdentifierByte(pattern[index-1])
			rightIdentifier := index+1 < len(pattern) && isIdentifierByte(pattern[index+1])
			if !leftIdentifier || !rightIdentifier {
				return true
			}
		}
	}
	return false
}

func declaredOrMatchedSymbol(line, pattern string) string {
	if symbol, ok := declaredSymbol(line); ok {
		return symbol
	}
	identifier := symbolForPattern(pattern)
	lowerLine, lowerIdentifier := strings.ToLower(line), strings.ToLower(identifier)
	if position := strings.Index(lowerLine, lowerIdentifier); position >= 0 {
		return line[position : position+len(identifier)]
	}
	return identifier
}

func declaredSymbol(line string) (string, bool) {
	fields := strings.Fields(strings.TrimSpace(line))
	for index, field := range fields {
		keyword := strings.Trim(field, "(){}:")
		if keyword != "func" && keyword != "type" && keyword != "class" && keyword != "struct" && keyword != "def" && keyword != "fn" && keyword != "test" && keyword != "it" {
			continue
		}
		for candidate := index + 1; candidate < len(fields); candidate++ {
			name := strings.Trim(fields[candidate], "(){}:,[<")
			if name == "" || strings.HasPrefix(name, "(") {
				continue
			}
			if cut := strings.IndexAny(name, "([<{:"); cut >= 0 {
				name = name[:cut]
			}
			if name != "" {
				return name, true
			}
		}
	}
	return "", false
}

func symbolForPattern(pattern string) string {
	pattern = strings.TrimSpace(pattern)
	if fields := strings.Fields(pattern); len(fields) > 1 {
		joined := strings.Join(fields, "")
		for index := range joined {
			if !isIdentifierByte(joined[index]) {
				joined = ""
				break
			}
		}
		if joined != "" {
			return joined
		}
	}
	start := 0
	for start < len(pattern) && !isIdentifierByte(pattern[start]) {
		start++
	}
	end := start
	for end < len(pattern) && isIdentifierByte(pattern[end]) {
		end++
	}
	if start == end {
		return ""
	}
	return pattern[start:end]
}

func opaqueRuleSignal(item domain.CollectedEvidence, rule domain.CapabilityRule) bool {
	path := normalizeEvidenceText(item.Path)
	for _, hint := range rule.TargetHints {
		if exactLexicalPatternMatch(path, normalizeEvidenceText(hint)) {
			return true
		}
	}
	patterns := []string{}
	for _, group := range append(append([]domain.CapabilityPattern{}, rule.StrongEvidence...), rule.SupportingEvidence...) {
		patterns = append(patterns, group.Patterns...)
	}
	raw := strings.Join(strings.Fields(item.Content), " ")
	for _, pattern := range patterns {
		if exactLexicalPatternMatch(raw, pattern) {
			return true
		}
	}
	return false
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

func sourceFromMatch(item domain.CollectedEvidence, kind domain.EvidenceKind, rule string, strength float64, match locatedEvidenceMatch) domain.EvidenceSource {
	return domain.EvidenceSource{
		Kind: kind, Ref: item.Ref, Path: item.Path, Symbol: match.Symbol,
		Line: match.Line, Snippet: match.Snippet, RuleID: rule,
		MatchedPattern: match.Pattern, MatchClass: match.Class,
		ContentHash: item.ContentHash, Strength: strength,
	}
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
func dedupeSources(in []domain.EvidenceSource) []domain.EvidenceSource {
	seen := map[string]bool{}
	out := make([]domain.EvidenceSource, 0, len(in))
	sort.SliceStable(in, func(i, j int) bool {
		leftClass, rightClass := evidenceMatchClassRank(in[i].MatchClass), evidenceMatchClassRank(in[j].MatchClass)
		if leftClass != rightClass {
			return leftClass < rightClass
		}
		if in[i].Path != in[j].Path {
			return in[i].Path < in[j].Path
		}
		if in[i].Line != in[j].Line {
			return in[i].Line < in[j].Line
		}
		if in[i].MatchedPattern != in[j].MatchedPattern {
			return in[i].MatchedPattern < in[j].MatchedPattern
		}
		if in[i].Snippet != in[j].Snippet {
			return in[i].Snippet < in[j].Snippet
		}
		if in[i].Strength != in[j].Strength {
			return in[i].Strength > in[j].Strength
		}
		return in[i].Kind < in[j].Kind
	})
	for _, s := range in {
		key := fmt.Sprintf("%s|%d|%s|%s|%s|%s", s.Path, s.Line, s.MatchedPattern, s.MatchClass, s.Snippet, s.Kind)
		if !seen[key] {
			seen[key] = true
			out = append(out, s)
		}
	}
	return out
}

func evidenceMatchClassRank(class domain.EvidenceMatchClass) int {
	switch class {
	case domain.EvidenceMatchStrongContent:
		return 0
	case domain.EvidenceMatchSupportingContent:
		return 1
	case domain.EvidenceMatchDocumentation:
		return 2
	case domain.EvidenceMatchPathOnly:
		return 3
	default:
		// Legacy non-rule sources retain their historical relative ordering.
		return 4
	}
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
	if !b.Complete || len(b.Failures) > 0 {
		return domain.LicenseAssessment{Status: domain.EvidenceInconclusive, Caveat: "coleta incompleta"}
	}
	var roots []domain.CollectedEvidence
	for _, i := range b.Items {
		// Inventory authority is path-based. A caller-supplied Kind must not
		// hide another root document from the ambiguity check.
		if authoritativeLicensePath(i.Path) {
			roots = append(roots, i)
		}
	}
	if len(roots) != 1 {
		return domain.LicenseAssessment{Status: domain.EvidenceInconclusive, Caveat: "licença ausente ou múltiplos documentos raiz; revisão de escopo necessária"}
	}
	i := roots[0]
	if i.Kind != domain.EvidenceLicense {
		return domain.LicenseAssessment{Status: domain.EvidenceInconclusive, Caveat: "classificação do documento raiz de licença incoerente"}
	}
	id := recognizedLicenseText(i.Content)
	s := sourceFromItem(i, domain.EvidenceLicense, "license-full-text", .7, firstLine(i.Content, []string{"license", "copyright", "permission"}))
	if id != "" {
		return domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: id, RecognitionBasis: "FULL_TEXT_TEMPLATE", Sources: []domain.EvidenceSource{s}, Caveat: "texto MIT completo reconhecido; compatibilidade, titularidade e escopo jurídico não inferidos"}
	}
	return domain.LicenseAssessment{Status: domain.EvidenceInconclusive, RecognitionBasis: "UNVERIFIED_TEXT", DetectedIdentifiers: detectedLicenseMarkers(i.Content), Sources: []domain.EvidenceSource{s}, Caveat: "marcadores são consultivos; texto incompleto, alterado ou licença sem verificador de texto completo"}
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
