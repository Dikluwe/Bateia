/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation.md
 * @layer L1
 * @updated 2026-09-01
 */

package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

type RecommendationService struct {
	dependencies IdeaRunDependencies
}

type RecommendationExecution struct {
	Run            domain.IdeaRunResult
	Recommendation domain.RecommendationResult
}

func NewRecommendationService(semantic SemanticExpander, discovery domain.ObservableDiscoveryProvider, content domain.ContentProvider) *RecommendationService {
	return &RecommendationService{dependencies: IdeaRunDependencies{Semantic: semantic, Discovery: discovery, Content: content}}
}

func ValidateRecommendationRequest(request domain.RecommendationRequest) (domain.RecommendationRequest, error) {
	request.Idea = normalizeIdea(request.Idea)
	if request.Idea.Problem == "" {
		return domain.RecommendationRequest{}, fmt.Errorf("idea.problem é obrigatório")
	}
	if len(request.Idea.Capabilities) == 0 {
		return domain.RecommendationRequest{}, fmt.Errorf("idea exige ao menos uma capacidade obrigatória")
	}
	if request.Policy.SchemaVersion != "1" {
		return domain.RecommendationRequest{}, fmt.Errorf("policy.schema_version deve ser 1")
	}
	if request.Policy.MaxComponents < 2 {
		return domain.RecommendationRequest{}, fmt.Errorf("policy.max_components deve ser ao menos 2")
	}
	if invalidCoverageThreshold(request.Policy.ExtendMinCoverage) || invalidCoverageThreshold(request.Policy.IntegrationMinValue) {
		return domain.RecommendationRequest{}, fmt.Errorf("limiares da política devem estar em (0,1]")
	}
	if request.Policy.IntegrationMinValue > request.Policy.ExtendMinCoverage {
		return domain.RecommendationRequest{}, fmt.Errorf("policy.integration_min_value não pode exceder extend_min_coverage")
	}
	request.Policy.AllowedLicenses = normalizedLicenses(request.Policy.AllowedLicenses)
	request.Target.Name = normalizeIdeaText(request.Target.Name)
	if request.Target.Name == "" {
		request.Target.Name = request.Idea.Problem
	}
	if request.Target.Repository != nil {
		return domain.RecommendationRequest{}, fmt.Errorf("target.repository não pertence ao gateway síncrono")
	}
	switch request.Search.Strategy {
	case "literal":
	case "lexical":
	default:
		return domain.RecommendationRequest{}, fmt.Errorf("strategy %q desconhecida", request.Search.Strategy)
	}
	seenLocales := map[string]bool{}
	for _, locale := range request.Search.Locales {
		if seenLocales[locale] {
			return domain.RecommendationRequest{}, fmt.Errorf("search locale duplicado %q", locale)
		}
		seenLocales[locale] = true
	}
	locales, err := canonicalizeSearchLocales(request.Search.Locales)
	if err != nil {
		return domain.RecommendationRequest{}, err
	}
	if len(locales) != len(request.Search.Locales) {
		return domain.RecommendationRequest{}, fmt.Errorf("search locales inválidos")
	}
	request.Search.Locales = locales
	if request.Search.MaxPerQuery < 1 || request.Search.MaxPerQuery > 25 {
		return domain.RecommendationRequest{}, fmt.Errorf("search.max_per_query deve estar em [1,25]")
	}
	if request.Search.MaxCandidates < 1 || request.Search.MaxCandidates > 20 {
		return domain.RecommendationRequest{}, fmt.Errorf("search.max_candidates deve estar em [1,20]")
	}
	return request, nil
}

func (s *RecommendationService) Recommend(ctx context.Context, request domain.RecommendationRequest) (domain.RecommendationResult, error) {
	execution, err := s.RecommendWithRun(ctx, request)
	return execution.Recommendation, err
}

func (s *RecommendationService) RecommendWithRun(ctx context.Context, request domain.RecommendationRequest) (RecommendationExecution, error) {
	if s == nil {
		return RecommendationExecution{}, fmt.Errorf("serviço de recomendação não configurado")
	}
	normalized, err := ValidateRecommendationRequest(request)
	if err != nil {
		return RecommendationExecution{}, err
	}
	if normalized.Search.Strategy == "lexical" && s.dependencies.Semantic == nil {
		return RecommendationExecution{}, fmt.Errorf("estratégia lexical não foi configurada")
	}
	if ctx.Err() != nil {
		return RecommendationExecution{Recommendation: ProjectRecommendationResult(inconclusiveRecommendation("run", "cancelled", "execution cancelled"))}, nil
	}
	runRequest := domain.IdeaRunRequest{
		Idea: normalized.Idea, Policy: normalized.Policy, Mode: domain.IdeaRunLive,
		Strategy: normalized.Search.Strategy, Language: normalized.Search.Language,
		SearchLocales: normalized.Search.Locales, MaxPerQuery: normalized.Search.MaxPerQuery,
		MaxCandidates: normalized.Search.MaxCandidates,
	}
	run, runErr := RunIdeaPipeline(ctx, runRequest, s.dependencies)
	if runErr != nil {
		result := recommendationFromPartialRun(run, normalized.Target, normalized.Policy, runErr)
		return RecommendationExecution{Run: run, Recommendation: ProjectRecommendationResult(result)}, nil
	}
	alternatives, err := RankSoftwareAlternatives(normalized.Target, normalized.Idea, run.Discovery, run.Assessment.Assessments, normalized.Policy)
	if err != nil {
		return RecommendationExecution{Run: run}, err
	}
	result := domain.RecommendationResult{
		SchemaVersion: "1",
		Status:        IdeaRunStatusFor(run),
		Summary: domain.RecommendationSummary{
			ExpansionCount: len(run.Expansion.Expansions),
			CandidateCount: len(run.Discovery.Candidates),
			AssessedCount:  len(run.Assessment.Assessments),
		},
		Decision:     run.Decision,
		Alternatives: alternatives,
		Failures:     recommendationFailures(run),
		Warnings:     recommendationWarnings(run, alternatives),
	}
	return RecommendationExecution{Run: run, Recommendation: ProjectRecommendationResult(result)}, nil
}

func recommendationFromPartialRun(run domain.IdeaRunResult, target domain.SoftwareAlternativeTarget, policy domain.ReuseDecisionPolicy, runErr error) domain.RecommendationResult {
	failures := recommendationFailures(run)
	if len(failures) == 0 {
		stage := "pipeline"
		if len(run.Discovery.PerExpansion) > 0 || len(run.Discovery.Failures) > 0 {
			stage = "discovery"
		}
		failures = []domain.ExecutionFailure{{Scope: "run", Stage: stage, Message: publicFailureMessage(stage, runErr.Error())}}
	}
	decision := run.Decision
	if decision.SchemaVersion == "" {
		decision = domain.ReuseDecisionReport{
			SchemaVersion: "1", Decision: domain.DecisionUnknown,
			Policy: policy,
			Matrix: []domain.CapabilityComparison{}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{},
		}
	}
	return domain.RecommendationResult{
		SchemaVersion: "1", Status: domain.IdeaRunInconclusive,
		Summary:      domain.RecommendationSummary{ExpansionCount: len(run.Expansion.Expansions), CandidateCount: len(run.Discovery.Candidates), AssessedCount: len(run.Assessment.Assessments)},
		Decision:     decision,
		Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Target: target, Alternatives: []domain.SoftwareAlternative{}, Excluded: []domain.SoftwareAlternative{}},
		Failures:     failures, Warnings: []string{},
	}
}

func inconclusiveRecommendation(scope, stage, message string) domain.RecommendationResult {
	return domain.RecommendationResult{
		SchemaVersion: "1", Status: domain.IdeaRunInconclusive,
		Decision:     domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown, Matrix: []domain.CapabilityComparison{}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}},
		Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Alternatives: []domain.SoftwareAlternative{}, Excluded: []domain.SoftwareAlternative{}},
		Failures:     []domain.ExecutionFailure{{Scope: scope, Stage: stage, Message: message}}, Warnings: []string{},
	}
}

func recommendationFailures(run domain.IdeaRunResult) []domain.ExecutionFailure {
	values := append([]domain.ExecutionFailure(nil), run.Failures...)
	values = append(values, run.Discovery.Failures...)
	for _, assessment := range run.Assessment.Assessments {
		values = append(values, assessment.Failures...)
	}
	return canonicalRecommendationFailures(values)
}

func recommendationWarnings(run domain.IdeaRunResult, alternatives domain.SoftwareAlternativeReport) []string {
	values := []string{}
	if run.Expansion.Regional != nil {
		for _, gap := range run.Expansion.Regional.Gaps {
			values = append(values, fmt.Sprintf("locale %s: %s", gap.Locale, gap.Reason))
		}
	}
	collect := func(items []domain.SoftwareAlternative) {
		for _, alternative := range items {
			for _, reason := range alternative.Reasons {
				if strings.Contains(strings.ToLower(reason), "licença não reconhecida") {
					values = append(values, unknownLicenseWarning)
				}
			}
			values = append(values, alternative.Maintenance.Warnings...)
		}
	}
	collect(alternatives.Alternatives)
	collect(alternatives.Excluded)
	return canonicalRecommendationStrings(values, "warning")
}

func MarshalRecommendationResult(result domain.RecommendationResult) ([]byte, error) {
	return json.Marshal(ProjectRecommendationResult(result))
}

func ProjectRecommendationResult(result domain.RecommendationResult) domain.RecommendationResult {
	result.SchemaVersion = "1"
	result.Decision = projectDecision(result.Decision)
	result.Alternatives = projectAlternatives(result.Alternatives)
	result.Failures = canonicalRecommendationFailures(result.Failures)
	result.Warnings = canonicalRecommendationStrings(result.Warnings, "warning")
	if result.Decision.Decision == domain.DecisionUnknown {
		result.Status = domain.IdeaRunInconclusive
		result.Alternatives.TopChoice = nil
	} else if len(result.Failures) > 0 {
		result.Status = domain.IdeaRunPartial
	} else if result.Status != domain.IdeaRunPartial && result.Status != domain.IdeaRunInconclusive {
		result.Status = domain.IdeaRunComplete
	}
	if !validRecommendationTopChoice(result.Alternatives) {
		result.Alternatives.TopChoice = nil
	}
	return result
}

func projectDecision(value domain.ReuseDecisionReport) domain.ReuseDecisionReport {
	value.Matrix = append([]domain.CapabilityComparison(nil), value.Matrix...)
	for i := range value.Matrix {
		value.Matrix[i].Candidates = append([]domain.CandidateCapabilityStatus(nil), value.Matrix[i].Candidates...)
		for j := range value.Matrix[i].Candidates {
			value.Matrix[i].Candidates[j].Evidence = projectEvidence(value.Matrix[i].Candidates[j].Evidence)
		}
	}
	value.Components = append([]domain.ReusableComponent(nil), value.Components...)
	for i := range value.Components {
		value.Components[i].Evidence = projectEvidence(value.Components[i].Evidence)
		value.Components[i].Capabilities = canonicalRecommendationStrings(value.Components[i].Capabilities, "value")
		value.Components[i].Symbols = canonicalRecommendationStrings(value.Components[i].Symbols, "value")
		value.Components[i].Caveats = canonicalRecommendationStrings(value.Components[i].Caveats, "caveat")
	}
	value.Blockers = projectDecisionReasons(value.Blockers)
	value.Reasons = projectDecisionReasons(value.Reasons)
	if value.Matrix == nil {
		value.Matrix = []domain.CapabilityComparison{}
	}
	if value.Gaps == nil {
		value.Gaps = []string{}
	}
	if value.Reasons == nil {
		value.Reasons = []domain.DecisionReason{}
	}
	if value.Alternatives == nil {
		value.Alternatives = []domain.DecisionAlternative{}
	}
	return value
}

func projectDecisionReasons(values []domain.DecisionReason) []domain.DecisionReason {
	out := append([]domain.DecisionReason(nil), values...)
	for i := range out {
		out[i].Message = publicText(out[i].Message, "decision reason")
		out[i].Evidence = projectEvidence(out[i].Evidence)
	}
	return out
}

func projectAlternatives(value domain.SoftwareAlternativeReport) domain.SoftwareAlternativeReport {
	value.Alternatives = append([]domain.SoftwareAlternative(nil), value.Alternatives...)
	value.Excluded = append([]domain.SoftwareAlternative(nil), value.Excluded...)
	for i := range value.Alternatives {
		projectAlternative(&value.Alternatives[i])
	}
	for i := range value.Excluded {
		projectAlternative(&value.Excluded[i])
	}
	if value.Alternatives == nil {
		value.Alternatives = []domain.SoftwareAlternative{}
	}
	if value.Excluded == nil {
		value.Excluded = []domain.SoftwareAlternative{}
	}
	return value
}

func projectAlternative(value *domain.SoftwareAlternative) {
	value.URL = publicHTTPSURL(value.URL)
	value.Evidence = projectEvidence(value.Evidence)
	value.License.Sources = projectEvidence(value.License.Sources)
	value.License.Caveat = publicText(value.License.Caveat, "license caveat")
	value.Maintenance.Sources = projectEvidence(value.Maintenance.Sources)
	value.Maintenance.Warnings = canonicalRecommendationStrings(value.Maintenance.Warnings, "maintenance warning")
	value.Operation.Sources = projectEvidence(value.Operation.Sources)
	value.Operation.Requirements = canonicalRecommendationStrings(value.Operation.Requirements, "operation requirement")
	value.Capabilities = append([]domain.SoftwareAlternativeCapability(nil), value.Capabilities...)
	for i := range value.Capabilities {
		value.Capabilities[i].Evidence = projectEvidence(value.Capabilities[i].Evidence)
	}
	value.Reasons = canonicalRecommendationStrings(value.Reasons, "reason")
	value.Blockers = canonicalRecommendationStrings(value.Blockers, "blocker")
	value.Gaps = canonicalRecommendationStrings(value.Gaps, "gap")
	value.DiscoveryLocales, _ = canonicalizeSearchLocales(value.DiscoveryLocales)
	if value.Capabilities == nil {
		value.Capabilities = []domain.SoftwareAlternativeCapability{}
	}
	if value.Gaps == nil {
		value.Gaps = []string{}
	}
	if value.Reasons == nil {
		value.Reasons = []string{}
	}
	if value.Blockers == nil {
		value.Blockers = []string{}
	}
	if value.Evidence == nil {
		value.Evidence = []domain.EvidenceSource{}
	}
}

func projectEvidence(values []domain.EvidenceSource) []domain.EvidenceSource {
	out := append([]domain.EvidenceSource(nil), values...)
	for i := range out {
		out[i].Snippet = ""
	}
	return canonicalEvidenceSources(out)
}

func canonicalRecommendationFailures(values []domain.ExecutionFailure) []domain.ExecutionFailure {
	byKey := map[string]domain.ExecutionFailure{}
	for _, value := range values {
		value.Scope = publicText(value.Scope, "scope")
		value.Stage = publicText(value.Stage, "stage")
		value.Message = publicFailureMessage(value.Stage, value.Message)
		key := value.Stage + "\x00" + value.Scope + "\x00" + value.Message
		byKey[key] = value
	}
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]domain.ExecutionFailure, 0, len(keys))
	for _, key := range keys {
		out = append(out, byKey[key])
	}
	return out
}

func canonicalRecommendationStrings(values []string, fallback string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		value = publicText(value, fallback)
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func publicFailureMessage(stage, value string) string {
	stage = normalizeIdeaText(stage)
	switch stage {
	case "discover", "discovery", "provider":
		return "provider failure"
	case "":
		return "execution failure"
	default:
		return stage + " failure"
	}
}

func publicText(value, fallback string) string {
	value = strings.Join(strings.Fields(value), " ")
	if containsSensitiveRecommendationText(value) || strings.HasPrefix(strings.ToLower(value), "falha de coleta ") {
		return fallback
	}
	return value
}

func containsSensitiveRecommendationText(value string) bool {
	upper := strings.ToUpper(value)
	for _, marker := range []string{"STEP0025_", "AUTHORIZATION", "BEARER ", "COOKIE=", "TOKEN=", "SECRET", "PROVIDER_CONFIG", "ADAPTER_METADATA", "ENV_VALUE", "ENV_NAME", "HEADER_CANARY"} {
		if strings.Contains(upper, marker) {
			return true
		}
	}
	return false
}

func publicHTTPSURL(raw string) string {
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return ""
	}
	return parsed.String()
}

func validRecommendationTopChoice(report domain.SoftwareAlternativeReport) bool {
	if report.TopChoice == nil || report.TopChoice.Commit == "" || strings.EqualFold(report.TopChoice.Commit, "HEAD") {
		return false
	}
	for _, alternative := range report.Alternatives {
		if alternative.Repository == *report.TopChoice && (alternative.Fit == domain.AlternativeRecommended || alternative.Fit == domain.AlternativePartial) && len(alternative.Blockers) == 0 {
			return true
		}
	}
	return false
}
