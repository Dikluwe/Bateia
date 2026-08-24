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
		rule("blast-radius", []string{"impacto", "blast radius"}, []string{"blastradius", "impact analysis", "blast radius"}),
		rule("structural-change", []string{"comparacao estrutural de mudancas", "structural change"}, []string{"structuraldiff", "architecture diff", "structural change"}),
		rule("static-analysis", []string{"analise estatica", "static analysis"}, []string{"staticanalyzer", "static analysis", "lint"}),
		rule("local-offline", []string{"execucao local", "offline", "local execution"}, []string{"offline", "run locally", "local execution"}),
		rule("agent-mcp", []string{"integracao mcp", "agentes", "model context protocol"}, []string{"mcpserver", "model context protocol", "mcp server"}),
		rule("async-events", []string{"processamento assincrono de eventos", "asynchronous event processing"}, []string{"eventprocessor", "event queue", "async event"}),
		rule("recipe-import", []string{"organizacao de receitas", "importacao de receitas", "recipe import"}, []string{"recipeimporter", "import recipe", "recipe organizer"}),
	}
	catalog := domain.CapabilityRuleCatalog{SchemaVersion: "1", Version: "2026-08-24", Rules: rules}
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
	selected := selectEvidenceCandidates(report.Decisions, 10, 2, 16)
	plans := make([]domain.EvidenceExtractionPlan, 0, len(selected))
	for _, decision := range selected {
		ref, err := parseImmutableRepository(decision.Repository)
		if err != nil {
			return nil, err
		}
		candidate := domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit, URL: "https://github.com/" + ref.Owner + "/" + ref.Repo}
		targets := []domain.ExtractionTarget{
			{Path: "README.md", Kind: domain.EvidenceDocumentation, Reason: "documentação principal"},
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

func parseImmutableRepository(value string) (domain.CommitRef, error) {
	parts := strings.Split(value, "@")
	if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" || strings.EqualFold(parts[1], "HEAD") {
		return domain.CommitRef{}, fmt.Errorf("candidato %q deve usar referência imutável owner/repo@commit", value)
	}
	repo := strings.Split(strings.Trim(parts[0], "/"), "/")
	if len(repo) != 2 || repo[0] == "" || repo[1] == "" {
		return domain.CommitRef{}, fmt.Errorf("repositório inválido %q; use owner/repo@commit", value)
	}
	return domain.CommitRef{Owner: repo[0], Repo: repo[1], Commit: parts[1]}, nil
}

func CollectEvidence(ctx context.Context, plan domain.EvidenceExtractionPlan, provider domain.ContentProvider) (domain.EvidenceBundle, error) {
	if provider == nil {
		return domain.EvidenceBundle{}, fmt.Errorf("provedor de conteúdo não configurado")
	}
	if plan.Repository.Commit == "" || strings.EqualFold(plan.Repository.Commit, "HEAD") {
		return domain.EvidenceBundle{}, fmt.Errorf("coleta exige commit imutável")
	}
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: plan.Repository, Plan: plan, Complete: true}
	paths := make([]string, 0)
	kinds := map[string]domain.EvidenceKind{}
	for _, target := range plan.Targets {
		paths = append(paths, target.Path)
		kinds[target.Path] = target.Kind
	}
	tree, err := provider.FetchTree(ctx, plan.Candidate)
	if err != nil {
		bundle.Failures = append(bundle.Failures, domain.ExecutionFailure{Scope: plan.Repository.Key(), Stage: "tree", Message: err.Error()})
		bundle.Complete = false
	} else {
		for _, node := range tree {
			if node.Type != "file" || pathDepth(node.Path) > plan.Budget.TreeDepth || ignoredEvidencePath(node.Path) {
				continue
			}
			kind := classifyEvidencePath(node.Path)
			if kind == domain.EvidenceFile {
				continue
			}
			if withinKindBudget(paths, kinds, kind, plan.Budget) {
				paths = append(paths, node.Path)
				kinds[node.Path] = kind
			}
		}
	}
	sort.Strings(paths)
	extracted, err := provider.ExtractRepoFiles(ctx, plan.Candidate, paths)
	if err != nil {
		bundle.Failures = append(bundle.Failures, domain.ExecutionFailure{Scope: plan.Repository.Key(), Stage: "content", Message: err.Error()})
		bundle.Complete = false
		return bundle, nil
	}
	seenHash := map[string]bool{}
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
	}
	return bundle, nil
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
		decisionByRepo[d.Repository] = d
	}
	report := domain.AssessmentReport{SchemaVersion: "1", RulesVersion: catalog.Version, RulesHash: catalog.Hash}
	for _, bundle := range bundles {
		assessment := domain.CandidateAssessment{Repository: bundle.Repository, ScreeningDecision: decisionByRepo[bundle.Repository.Key()], Failures: append([]domain.ExecutionFailure(nil), bundle.Failures...)}
		assessment.ProblemAlignment = assessCapability(idea.Problem, bundle, findRule(idea.Problem, catalog))
		for _, capability := range idea.Capabilities {
			e := assessCapability(capability, bundle, findRule(capability, catalog))
			assessment.Capabilities = append(assessment.Capabilities, e)
			if e.Status != domain.EvidenceProven {
				assessment.Gaps = append(assessment.Gaps, capability)
			}
			if e.Status == domain.EvidenceProven && len(e.Sources) > 0 {
				s := e.Sources[0]
				assessment.ReusableParts = append(assessment.ReusableParts, domain.ReusableComponent{Name: filepath.Base(s.Path), Kind: "module_or_pattern", Path: s.Path, Symbols: nonEmpty(s.Symbol), Capabilities: []string{capability}, Evidence: []domain.EvidenceSource{s}})
			}
		}
		for _, constraint := range idea.Constraints {
			assessment.Constraints = append(assessment.Constraints, domain.ConstraintEvidence{Constraint: constraint, Evidence: assessConstraint(constraint, bundle)})
		}
		assessment.License = assessLicense(bundle)
		assessment.Maintenance = assessMaintenance(bundle)
		assessment.Operation = assessOperation(bundle)
		report.Assessments = append(report.Assessments, assessment)
	}
	return report, nil
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
		for _, p := range rule.NegativeEvidence {
			if phraseMatch(norm, p) && matchesAny(norm, rule.CapabilityTerms) {
				negative = true
			}
		}
		for _, group := range rule.StrongEvidence {
			if strongKindMatches(item, group, norm) {
				kind := group.Kind
				if item.Kind == domain.EvidenceTest {
					kind = domain.EvidenceTest
				}
				sources = append(sources, sourceFromItem(item, kind, rule.ID, strengthFor(kind, item.Path), firstLine(item.Content, group.Patterns)))
			}
		}
		for _, group := range rule.SupportingEvidence {
			if matchesAny(norm, group.Patterns) {
				kind := item.Kind
				strength := strengthFor(kind, item.Path)
				sources = append(sources, sourceFromItem(item, kind, rule.ID, strength, firstLine(item.Content, group.Patterns)))
			}
		}
	}
	sources = dedupeSources(sources)
	if negative {
		return capabilityResult(capability, domain.EvidenceConflict, sources, maxStrength(sources), "evidência negativa explícita conflita com a capacidade ou restrição")
	}
	if len(bundle.Failures) > 0 || !bundle.Complete {
		return capabilityResult(capability, domain.EvidenceInconclusive, sources, maxStrength(sources), "coleta incompleta ou com falha; ausência não pode ser concluída")
	}
	max, combined, independent := evidenceScores(sources)
	if max >= 0.85 || (independent >= 2 && combined >= 0.90) {
		return capabilityResult(capability, domain.EvidenceProven, sources, combined, "capacidade comprovada por evidência concreta e rastreável")
	}
	if combined >= 0.50 {
		return capabilityResult(capability, domain.EvidencePartial, sources, combined, "há evidência concreta, mas insuficiente para comprovação")
	}
	return capabilityResult(capability, domain.EvidenceAbsent, sources, combined, "busca concluída no orçamento declarado sem evidência suficiente")
}

func strongKindMatches(item domain.CollectedEvidence, group domain.CapabilityPattern, normalized string) bool {
	if !matchesAny(normalized, group.Patterns) {
		return false
	}
	switch group.Kind {
	case domain.EvidenceTest:
		return item.Kind == domain.EvidenceTest
	case domain.EvidenceSymbol:
		if item.Kind != domain.EvidenceFile && item.Kind != domain.EvidenceExample {
			return false
		}
		content := strings.ToLower(item.Content)
		return strings.Contains(content, "func ") || strings.Contains(content, "type ") ||
			strings.Contains(content, "class ") || strings.Contains(content, "def ") ||
			strings.Contains(content, "fn ")
	default:
		return item.Kind == group.Kind
	}
}

func sourceFromItem(item domain.CollectedEvidence, kind domain.EvidenceKind, rule string, strength float64, snippet string) domain.EvidenceSource {
	return domain.EvidenceSource{Kind: kind, Ref: item.Ref, Path: item.Path, Line: 1, Snippet: snippet, RuleID: rule, ContentHash: item.ContentHash, Strength: strength}
}
func strengthFor(kind domain.EvidenceKind, path string) float64 {
	if ignoredEvidencePath(path) {
		return .05
	}
	switch kind {
	case domain.EvidenceTest:
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
		key := s.Path + "|" + s.ContentHash
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
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
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
func nonEmpty(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
}

func assessConstraint(c string, b domain.EvidenceBundle) domain.CapabilityEvidence {
	n := normalizeEvidenceText(c)
	if strings.Contains(n, "local") || strings.Contains(n, "offline") {
		o := assessOperation(b)
		return capabilityResult(c, o.LocalStatus, o.Sources, maxStrength(o.Sources), "restrição operacional avaliada por regras locais")
	}
	return capabilityResult(c, domain.EvidenceInconclusive, nil, 0, "restrição sem regra determinística específica")
}
func assessLicense(b domain.EvidenceBundle) domain.LicenseAssessment {
	for _, i := range b.Items {
		if i.Kind != domain.EvidenceLicense {
			continue
		}
		n := strings.ToLower(i.Content)
		id := ""
		for _, x := range []string{"MIT", "Apache-2.0", "MPL-2.0", "GPL-3.0"} {
			if strings.Contains(n, strings.ToLower(strings.Split(x, "-")[0])) {
				id = x
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
func assessMaintenance(b domain.EvidenceBundle) domain.MaintenanceAssessment {
	for _, i := range b.Items {
		n := normalizeEvidenceText(i.Content)
		if strings.Contains(n, "no longer maintained") || strings.Contains(n, "unmaintained") || strings.Contains(n, "maintenance ended") {
			s := sourceFromItem(i, i.Kind, "maintenance-ended", .7, firstLine(i.Content, []string{"maintained", "unmaintained"}))
			return domain.MaintenanceAssessment{Status: domain.EvidenceConflict, Sources: []domain.EvidenceSource{s}, Warnings: []string{"manutenção declarada como encerrada"}}
		}
	}
	if !b.Complete {
		return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive}
	}
	return domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive, Warnings: []string{"metadados temporais não disponíveis no bundle"}}
}
func assessOperation(b domain.EvidenceBundle) domain.OperationAssessment {
	var sources []domain.EvidenceSource
	for _, item := range b.Items {
		normalized := normalizeEvidenceText(item.Content)
		if (strings.Contains(normalized, "upload") || strings.Contains(normalized, "send")) &&
			(strings.Contains(normalized, "api") || strings.Contains(normalized, "cloud")) {
			sources = append(sources, sourceFromItem(item, item.Kind, "external-code-upload", .7, firstLine(item.Content, []string{"upload", "send"})))
		}
	}
	if len(sources) > 0 {
		return domain.OperationAssessment{LocalStatus: domain.EvidenceConflict, Sources: sources, Requirements: []string{"serviço externo ou envio de código"}}
	}
	if !b.Complete || len(b.Failures) > 0 {
		return domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive}
	}
	return domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive}
}
