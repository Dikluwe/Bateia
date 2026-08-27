/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
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
	}
	catalog := domain.CapabilityRuleCatalog{SchemaVersion: "1", Version: "2026-08-26.2", Rules: rules}
	catalog.Hash = ruleCatalogHash(catalog)
	return catalog
}

func rule(id string, terms, patterns []string) domain.CapabilityRule {
	return domain.CapabilityRule{
		ID: id, CapabilityTerms: terms,
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

func ruleCatalogHash(c domain.CapabilityRuleCatalog) string {
	c.Hash = ""
	b, _ := json.Marshal(c)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func ValidateRuleCatalog(c domain.CapabilityRuleCatalog) error {
	if strings.TrimSpace(c.Version) == "" || len(c.Rules) == 0 {
		return fmt.Errorf("catálogo de regras exige version e rules")
	}
	ids := map[string]bool{}
	for _, r := range c.Rules {
		if strings.TrimSpace(r.ID) == "" || ids[r.ID] {
			return fmt.Errorf("ID de regra vazio ou duplicado: %q", r.ID)
		}
		ids[r.ID] = true
		for _, group := range append(append([]domain.CapabilityPattern{}, r.StrongEvidence...), r.SupportingEvidence...) {
			for _, pattern := range group.Patterns {
				if len(pattern) == 0 || len(pattern) > 120 {
					return fmt.Errorf("regra %s contém padrão vazio ou longo demais", r.ID)
				}
			}
		}
	}
	return nil
}

func PlanEvidence(report domain.ScreeningReport) ([]domain.EvidenceExtractionPlan, error) {
	return PlanEvidenceForIdea(domain.IdeaSpec{}, report, DefaultEvidenceRules())
}

func PlanEvidenceForIdea(idea domain.IdeaSpec, report domain.ScreeningReport, catalog domain.CapabilityRuleCatalog) ([]domain.EvidenceExtractionPlan, error) {
	return PlanEvidenceForIdeaAndRepositories(idea, report, catalog, nil)
}

func PlanEvidenceForIdeaAndRepositories(idea domain.IdeaSpec, report domain.ScreeningReport, catalog domain.CapabilityRuleCatalog, includes []string) ([]domain.EvidenceExtractionPlan, error) {
	if err := ValidateRuleCatalog(catalog); err != nil {
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
	tree, err := provider.FetchTree(ctx, plan.Candidate)
	if err != nil {
		bundle.Failures = append(bundle.Failures, domain.ExecutionFailure{Scope: plan.Repository.Key(), Stage: "tree", Message: err.Error()})
		bundle.Complete = false
		bundle.TreeComplete = false
	} else {
		tree = prioritizeEvidenceTree(tree, plan.Targets)
		for _, node := range tree {
			if depth := pathDepth(node.Path); depth > bundle.DepthReached {
				bundle.DepthReached = depth
			}
			if node.Type != "file" || pathDepth(node.Path) > plan.Budget.TreeDepth || ignoredEvidencePath(node.Path) {
				continue
			}
			kind := classifyEvidencePath(node.Path)
			if kind == domain.EvidenceFile && !isCodePath(node.Path) {
				continue
			}
			if withinKindBudget(paths, kinds, kind, plan.Budget) {
				if _, exists := targetByPath[node.Path]; exists {
					continue
				}
				target := domain.ExtractionTarget{Path: node.Path, Kind: kind, Capabilities: capabilitiesForPath(node.Path, plan.Targets), Reason: "árvore limitada no commit resolvido"}
				paths = append(paths, node.Path)
				kinds[node.Path] = kind
				targetByPath[node.Path] = target
			}
		}
	}
	sort.Strings(paths)
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
	if bundle.BytesRead >= plan.Budget.TotalBytes {
		bundle.Complete = false
	}
	return bundle, nil
}

func prioritizeEvidenceTree(tree []domain.DirectoryNode, targets []domain.ExtractionTarget) []domain.DirectoryNode {
	out := append([]domain.DirectoryNode(nil), tree...)
	sort.SliceStable(out, func(i, j int) bool {
		iRelevant := len(capabilitiesForPath(out[i].Path, targets)) > 0
		jRelevant := len(capabilitiesForPath(out[j].Path, targets)) > 0
		if iRelevant != jRelevant {
			return iRelevant
		}
		return out[i].Path < out[j].Path
	})
	return out
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
	return len(strings.Split(filepath.ToSlash(strings.Trim(path, "/")), "/"))
}
func ignoredEvidencePath(path string) bool {
	p := "/" + strings.ToLower(filepath.ToSlash(path)) + "/"
	return strings.Contains(p, "/vendor/") || strings.Contains(p, "/fixtures/") || strings.Contains(p, "/translations/") || strings.Contains(p, "/node_modules/")
}
func classifyEvidencePath(path string) domain.EvidenceKind {
	p := strings.ToLower(filepath.ToSlash(path))
	base := filepath.Base(p)
	if strings.HasPrefix(base, "readme") || strings.Contains(p, "/docs/") {
		return domain.EvidenceDocumentation
	}
	if strings.HasPrefix(base, "license") || strings.HasPrefix(base, "copying") {
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
	extension := strings.ToLower(filepath.Ext(path))
	switch extension {
	case ".go", ".rs", ".py", ".js", ".jsx", ".ts", ".tsx", ".java", ".kt", ".rb", ".cs", ".c", ".cc", ".cpp", ".h", ".hpp", ".ex", ".exs":
		return true
	default:
		return false
	}
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
	if err := ValidateRuleCatalog(catalog); err != nil {
		return domain.AssessmentReport{}, err
	}
	if catalog.Hash == "" {
		catalog.Hash = ruleCatalogHash(catalog)
	}
	decisionByRepo := map[string]domain.ScreeningDecision{}
	for _, d := range screening.Decisions {
		decisionByRepo[repositoryIdentityKey(d.Repository)] = d
	}
	report := domain.AssessmentReport{SchemaVersion: "1", RulesVersion: catalog.Version, RulesHash: catalog.Hash}
	for _, bundle := range bundles {
		if err := validateEvidenceBundle(bundle); err != nil {
			return domain.AssessmentReport{}, err
		}
		decision, ok := decisionByRepo[bundle.Repository.Owner+"/"+bundle.Repository.Repo]
		if !ok {
			decision = bundle.Plan.Decision
		}
		assessment := domain.CandidateAssessment{Repository: bundle.Repository, ScreeningDecision: decision, Failures: append([]domain.ExecutionFailure(nil), bundle.Failures...)}
		assessment.ProblemAlignment = assessCapability(idea.Problem, bundle, findRule(idea.Problem, catalog))
		for _, capability := range idea.Capabilities {
			e := assessCapability(capability, bundle, findRule(capability, catalog))
			assessment.Capabilities = append(assessment.Capabilities, e)
			if e.Status != domain.EvidenceProven {
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

func validateEvidenceBundle(bundle domain.EvidenceBundle) error {
	if bundle.Repository.Owner == "" || bundle.Repository.Repo == "" || bundle.Repository.Commit == "" || bundle.Repository.Commit == "HEAD" {
		return fmt.Errorf("bundle exige repository em commit imutável")
	}
	for _, item := range bundle.Items {
		if item.Ref != bundle.Repository {
			return fmt.Errorf("fonte %s não pertence ao commit do bundle", item.Path)
		}
		if item.ContentHash == "" || item.ContentHash != contentHash(item.Content) {
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
	n := normalizeEvidenceText(capability)
	for i := range catalog.Rules {
		for _, term := range catalog.Rules[i].CapabilityTerms {
			if strings.Contains(n, normalizeEvidenceText(term)) || strings.Contains(normalizeEvidenceText(term), n) {
				return &catalog.Rules[i]
			}
		}
	}
	return nil
}

func assessCapability(capability string, bundle domain.EvidenceBundle, rule *domain.CapabilityRule) domain.CapabilityEvidence {
	if rule == nil {
		return capabilityResult(capability, domain.EvidenceInconclusive, nil, 0, "nenhuma regra versionada cobre esta capacidade")
	}
	var sources []domain.EvidenceSource
	negative := false
	for _, item := range bundle.Items {
		norm := normalizeEvidenceText(item.Content + " " + item.Path)
		if explicitNegativeMatch(item.Content, *rule) {
			negative = true
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
	if negative {
		return capabilityResult(capability, domain.EvidenceConflict, sources, maxStrength(sources), "evidência negativa explícita conflita com a capacidade ou restrição")
	}
	if !bundleConclusive(bundle) {
		return capabilityResult(capability, domain.EvidenceInconclusive, sources, maxStrength(sources), "coleta incompleta ou com falha; ausência não pode ser concluída")
	}
	max, combined, independent := evidenceScores(sources)
	if max >= 0.85 || (max >= 0.65 && independent >= 2 && combined >= 0.90) {
		return capabilityResult(capability, domain.EvidenceProven, sources, combined, "capacidade comprovada por evidência concreta e rastreável")
	}
	if combined >= 0.50 {
		return capabilityResult(capability, domain.EvidencePartial, sources, combined, "há evidência concreta, mas insuficiente para comprovação")
	}
	if len(sources) > 0 {
		return capabilityResult(capability, domain.EvidenceInconclusive, sources, combined, "sinais fracos impedem concluir presença ou ausência")
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
	return domain.EvidenceSource{Kind: kind, Ref: item.Ref, Path: item.Path, Line: 1, Snippet: snippet, RuleID: rule, ContentHash: item.ContentHash, Strength: strength}
}
func strengthFor(kind domain.EvidenceKind, path string) float64 {
	if ignoredEvidencePath(path) {
		return .05
	}
	switch kind {
	case domain.EvidenceTest, domain.EvidenceTestImplementation:
		return 1
	case domain.EvidenceSymbol:
		return .9
	case domain.EvidenceExample:
		return .85
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
	return domain.CapabilityEvidence{Capability: c, Status: status, Sources: s, Confidence: confidence, Explanation: explanation}
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
		if matchesAny(normalizeEvidenceText(item.Content), []string{nonGoal}) {
			sources = append(sources, sourceFromItem(item, item.Kind, "non-goal-conflict", supportingStrength(item.Kind, item.Path), firstLine(item.Content, []string{nonGoal})))
		}
	}
	sources = dedupeSources(sources)
	if len(sources) > 0 {
		return capabilityResult(nonGoal, domain.EvidenceConflict, sources, maxStrength(sources), "evidência explícita coincide com um non-goal")
	}
	if !bundleConclusive(b) {
		return capabilityResult(nonGoal, domain.EvidenceInconclusive, nil, 0, "coleta incompleta")
	}
	return capabilityResult(nonGoal, domain.EvidenceAbsent, nil, 0, "non-goal não observado no escopo coletado")
}
func assessLicense(b domain.EvidenceBundle) domain.LicenseAssessment {
	for _, i := range b.Items {
		if i.Kind != domain.EvidenceLicense {
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
	var metadataSource *domain.EvidenceSource
	var metadataCandidate domain.RepositoryCandidate
	for _, i := range b.Items {
		n := normalizeEvidenceText(i.Content)
		if strings.Contains(n, "no longer maintained") || strings.Contains(n, "unmaintained") || strings.Contains(n, "maintenance ended") {
			s := sourceFromItem(i, i.Kind, "maintenance-ended", .7, firstLine(i.Content, []string{"maintained", "unmaintained"}))
			return domain.MaintenanceAssessment{Status: domain.EvidenceConflict, Sources: []domain.EvidenceSource{s}, Warnings: []string{"manutenção declarada como encerrada"}}
		}
		if i.Kind == domain.EvidenceRepository {
			if json.Unmarshal([]byte(i.Content), &metadataCandidate) == nil {
				s := sourceFromItem(i, domain.EvidenceRepository, "repository-maintenance", .7, "metadados imutáveis da coleta")
				metadataSource = &s
			}
		}
	}
	if metadataSource != nil {
		if metadataCandidate.Archived {
			return domain.MaintenanceAssessment{Status: domain.EvidenceConflict, Archived: true, Sources: []domain.EvidenceSource{*metadataSource}, Warnings: []string{"repositório arquivado"}}
		}
		if !metadataCandidate.UpdatedAt.IsZero() {
			return domain.MaintenanceAssessment{Status: domain.EvidencePartial, Sources: []domain.EvidenceSource{*metadataSource}, Warnings: []string{"atividade observada em " + metadataCandidate.UpdatedAt.UTC().Format("2006-01-02")}}
		}
	}
	if !b.Complete {
		return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive}
	}
	return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive, Warnings: []string{"metadados temporais não disponíveis no bundle"}}
}
func assessOperation(b domain.EvidenceBundle) domain.OperationAssessment {
	var sources []domain.EvidenceSource
	var localSources []domain.EvidenceSource
	requirementSet := map[string]bool{}
	for _, item := range b.Items {
		for lineNumber, raw := range strings.Split(item.Content, "\n") {
			line := normalizeEvidenceText(raw)
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
		if strings.Contains(normalized, "run locally") || strings.Contains(normalized, "local execution") || strings.Contains(normalized, "works offline") || strings.Contains(normalized, "offline mode") {
			localSources = append(localSources, sourceFromItem(item, item.Kind, "local-operation", strengthFor(item.Kind, item.Path), firstLine(item.Content, []string{"run locally", "local execution", "offline"})))
		}
	}
	if len(sources) > 0 {
		requirementSet["serviço externo ou envio de código"] = true
		return domain.OperationAssessment{LocalStatus: domain.EvidenceConflict, Sources: sources, Requirements: sortedRequirements(requirementSet)}
	}
	if !b.Complete || len(b.Failures) > 0 {
		return domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive}
	}
	if len(localSources) > 0 {
		return domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: dedupeSources(localSources), Requirements: sortedRequirements(requirementSet)}
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
	return transfer && code && destination
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
