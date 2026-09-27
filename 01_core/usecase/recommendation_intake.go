/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_intake.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

const recommendationExecutionDefaultsV1 = "BATEIA_EXECUTION_DEFAULTS_V1"

func PrepareRecommendationIntake(draft domain.RecommendationIntakeDraft) domain.RecommendationIntake {
	intake := domain.RecommendationIntake{
		SchemaVersion: domain.RecommendationIntakeSchemaVersion,
		Request: domain.RecommendationRequest{
			Idea:   domain.IdeaSpec{Users: []string{}, Capabilities: []string{}, Constraints: []string{}, NonGoals: []string{}, PreferredStack: []string{}},
			Policy: DefaultRecommendationPolicyV1(),
			Search: DefaultRecommendationSearchV1(),
		},
		Execution: domain.RecommendationExecutionProfile{
			Mode: "live", Provider: "sourcegraph", LocalPaths: []string{}, TimeoutSeconds: 120,
		},
		Provenance: []domain.InputProvenance{}, Issues: []domain.InputIssue{}, Questions: []domain.InformationNeed{}, Answers: []domain.InformationAnswer{},
	}

	materializeIdeaInput(&intake, draft.Idea)
	materializePolicyInput(&intake, draft.Policy)
	materializeTargetInput(&intake, draft.Target)
	materializeSearchInput(&intake, draft.Search)
	materializeExecutionInput(&intake, draft.Execution)
	assessPreparedIntake(&intake)
	return intake
}

func materializeIdeaInput(intake *domain.RecommendationIntake, input *domain.IdeaIntakeInput) {
	explicitString := func(path string, value *string, target *string) {
		if value == nil {
			addProvenance(intake, path, domain.InputUnresolved, domain.AuthorityUser, "")
			return
		}
		*target = *value
		addProvenance(intake, path, domain.InputExplicit, domain.AuthorityUser, "")
	}
	explicitList := func(path string, value *[]string, target *[]string) {
		if value == nil {
			if path == "idea.capabilities" {
				addProvenance(intake, path, domain.InputUnresolved, domain.AuthorityUser, "")
				return
			}
			addProvenance(intake, path, domain.InputVersionedDefault, domain.AuthorityProductPolicy, RecommendationDefaultsProfileV1)
			return
		}
		*target = append([]string{}, (*value)...)
		addProvenance(intake, path, domain.InputExplicit, domain.AuthorityUser, "")
	}
	if input == nil {
		explicitString("idea.problem", nil, &intake.Request.Idea.Problem)
		explicitList("idea.users", nil, &intake.Request.Idea.Users)
		explicitList("idea.capabilities", nil, &intake.Request.Idea.Capabilities)
		explicitList("idea.qualities", nil, &intake.Request.Idea.Qualities)
		explicitList("idea.operating_conditions", nil, &intake.Request.Idea.OperatingConditions)
		explicitList("idea.constraints", nil, &intake.Request.Idea.Constraints)
		explicitList("idea.non_goals", nil, &intake.Request.Idea.NonGoals)
		explicitList("idea.preferred_stack", nil, &intake.Request.Idea.PreferredStack)
		return
	}
	explicitString("idea.problem", input.Problem, &intake.Request.Idea.Problem)
	explicitList("idea.users", input.Users, &intake.Request.Idea.Users)
	explicitList("idea.capabilities", input.Capabilities, &intake.Request.Idea.Capabilities)
	explicitList("idea.qualities", input.Qualities, &intake.Request.Idea.Qualities)
	explicitList("idea.operating_conditions", input.OperatingConditions, &intake.Request.Idea.OperatingConditions)
	explicitList("idea.constraints", input.Constraints, &intake.Request.Idea.Constraints)
	explicitList("idea.non_goals", input.NonGoals, &intake.Request.Idea.NonGoals)
	explicitList("idea.preferred_stack", input.PreferredStack, &intake.Request.Idea.PreferredStack)
}

func materializePolicyInput(intake *domain.RecommendationIntake, input *domain.PolicyIntakeInput) {
	defaults := DefaultRecommendationPolicyV1()
	fields := []struct {
		path     string
		explicit bool
	}{
		{"policy.schema_version", input != nil && input.SchemaVersion != nil},
		{"policy.allowed_licenses", input != nil && input.AllowedLicenses != nil},
		{"policy.require_local", input != nil && input.RequireLocal != nil},
		{"policy.allow_archived", input != nil && input.AllowArchived != nil},
		{"policy.max_components", input != nil && input.MaxComponents != nil},
		{"policy.extend_min_coverage", input != nil && input.ExtendMinCoverage != nil},
		{"policy.integration_min_value", input != nil && input.IntegrationMinValue != nil},
	}
	for _, field := range fields {
		origin, authority, source := domain.InputVersionedDefault, domain.AuthorityProductPolicy, RecommendationDefaultsProfileV1
		if field.explicit {
			origin, authority, source = domain.InputExplicit, domain.AuthorityUser, ""
		}
		addProvenance(intake, field.path, origin, authority, source)
	}
	if input == nil {
		return
	}
	if input.SchemaVersion != nil {
		intake.Request.Policy.SchemaVersion = *input.SchemaVersion
	}
	if input.AllowedLicenses != nil {
		intake.Request.Policy.AllowedLicenses = append([]string{}, (*input.AllowedLicenses)...)
	} else {
		intake.Request.Policy.AllowedLicenses = append([]string{}, defaults.AllowedLicenses...)
	}
	if input.RequireLocal != nil {
		intake.Request.Policy.RequireLocal = *input.RequireLocal
	}
	if input.AllowArchived != nil {
		intake.Request.Policy.AllowArchived = *input.AllowArchived
	}
	if input.MaxComponents != nil {
		intake.Request.Policy.MaxComponents = *input.MaxComponents
	}
	if input.ExtendMinCoverage != nil {
		intake.Request.Policy.ExtendMinCoverage = *input.ExtendMinCoverage
	}
	if input.IntegrationMinValue != nil {
		intake.Request.Policy.IntegrationMinValue = *input.IntegrationMinValue
	}
}

func materializeTargetInput(intake *domain.RecommendationIntake, input *domain.TargetIntakeInput) {
	if input != nil && input.Name != nil {
		intake.Request.Target.Name = *input.Name
		addProvenance(intake, "target.name", domain.InputExplicit, domain.AuthorityUser, "")
		return
	}
	intake.Request.Target = DefaultRecommendationTargetV1(intake.Request.Idea.Problem)
	addProvenance(intake, "target.name", domain.InputDerived, domain.AuthorityProductPolicy, "idea.problem")
}

func materializeSearchInput(intake *domain.RecommendationIntake, input *domain.SearchIntakeInput) {
	fields := []struct {
		path     string
		explicit bool
	}{
		{"search.strategy", input != nil && input.Strategy != nil},
		{"search.language", input != nil && input.Language != nil},
		{"search.locales", input != nil && input.Locales != nil},
		{"search.max_per_query", input != nil && input.MaxPerQuery != nil},
		{"search.max_candidates", input != nil && input.MaxCandidates != nil},
	}
	for _, field := range fields {
		origin, authority, source := domain.InputVersionedDefault, domain.AuthorityProductPolicy, RecommendationDefaultsProfileV1
		if field.explicit {
			origin, authority, source = domain.InputExplicit, domain.AuthorityUser, ""
		}
		addProvenance(intake, field.path, origin, authority, source)
	}
	if input == nil {
		return
	}
	if input.Strategy != nil {
		intake.Request.Search.Strategy = *input.Strategy
	}
	if input.Language != nil {
		intake.Request.Search.Language = *input.Language
	}
	if input.Locales != nil {
		intake.Request.Search.Locales = append([]string{}, (*input.Locales)...)
	}
	if input.MaxPerQuery != nil {
		intake.Request.Search.MaxPerQuery = *input.MaxPerQuery
	}
	if input.MaxCandidates != nil {
		intake.Request.Search.MaxCandidates = *input.MaxCandidates
	}
}

func materializeExecutionInput(intake *domain.RecommendationIntake, input *domain.ExecutionIntakeInput) {
	add := func(path string, explicit bool) {
		origin, authority, source := domain.InputVersionedDefault, domain.AuthorityProductPolicy, recommendationExecutionDefaultsV1
		if explicit {
			origin, authority, source = domain.InputExplicit, domain.AuthorityUser, ""
		}
		addProvenance(intake, path, origin, authority, source)
	}
	add("execution.mode", input != nil && input.Mode != nil)
	add("execution.provider", input != nil && input.Provider != nil)
	add("execution.local_paths", input != nil && input.LocalPaths != nil)
	add("execution.replay_directory", input != nil && input.ReplayDirectory != nil)
	add("execution.output_directory", input != nil && input.OutputDirectory != nil)
	add("execution.resume", input != nil && input.Resume != nil)
	add("execution.timeout_seconds", input != nil && input.TimeoutSeconds != nil)
	add("execution.allow_candidate_execution", input != nil && input.AllowCandidateExecution != nil)
	add("execution.evidence_rules_ref", input != nil && input.EvidenceRulesRef != nil)
	add("execution.screening_model_ref", input != nil && input.ScreeningModelRef != nil)
	add("execution.idempotency_key", input != nil && input.IdempotencyKey != nil)
	if input == nil {
		return
	}
	if input.Mode != nil {
		intake.Execution.Mode = *input.Mode
	}
	if input.Provider != nil {
		intake.Execution.Provider = *input.Provider
	}
	if input.LocalPaths != nil {
		intake.Execution.LocalPaths = append([]string{}, (*input.LocalPaths)...)
	}
	if input.ReplayDirectory != nil {
		intake.Execution.ReplayDirectory = *input.ReplayDirectory
	}
	if input.OutputDirectory != nil {
		intake.Execution.OutputDirectory = *input.OutputDirectory
	}
	if input.Resume != nil {
		intake.Execution.Resume = *input.Resume
	}
	if input.TimeoutSeconds != nil {
		intake.Execution.TimeoutSeconds = *input.TimeoutSeconds
	}
	if input.AllowCandidateExecution != nil {
		intake.Execution.AllowCandidateExecution = *input.AllowCandidateExecution
	}
	if input.EvidenceRulesRef != nil {
		intake.Execution.EvidenceRulesRef = *input.EvidenceRulesRef
	}
	if input.ScreeningModelRef != nil {
		intake.Execution.ScreeningModelRef = *input.ScreeningModelRef
	}
	if input.IdempotencyKey != nil {
		intake.Execution.IdempotencyKey = *input.IdempotencyKey
	}
}

func assessPreparedIntake(intake *domain.RecommendationIntake) {
	intake.Issues = []domain.InputIssue{}
	intake.Questions = []domain.InformationNeed{}
	if strings.TrimSpace(intake.Request.Idea.Problem) == "" {
		addRequiredNeed(intake, "idea.problem", "problem_required", domain.AnswerText, "problema que a recomendação deve resolver")
	}
	if len(normalizedIntakeStrings(intake.Request.Idea.Capabilities)) == 0 {
		addRequiredNeed(intake, "idea.capabilities", "capabilities_required", domain.AnswerStringList, "capacidades obrigatórias da solução")
	}
	if len(intake.Questions) > 0 {
		intake.Status = domain.IntakeNeedsInput
	} else if normalized, err := ValidateRecommendationRequest(intake.Request); err != nil {
		intake.Status = domain.IntakeInvalid
		intake.Issues = append(intake.Issues, domain.InputIssue{Path: "request", Code: "INVALID_REQUEST", Message: err.Error(), Blocking: true})
	} else {
		intake.Request = normalized
		intake.Status = domain.IntakeReady
	}
	if intake.Execution.TimeoutSeconds <= 0 {
		intake.Status = domain.IntakeInvalid
		intake.Issues = append(intake.Issues, domain.InputIssue{Path: "execution.timeout_seconds", Code: "INVALID_TIMEOUT", Message: "execution.timeout_seconds deve ser positivo", Blocking: true})
	}
	if intake.Execution.Mode != "live" && intake.Execution.Mode != "local" && intake.Execution.Mode != "replay" {
		intake.Status = domain.IntakeInvalid
		intake.Issues = append(intake.Issues, domain.InputIssue{Path: "execution.mode", Code: "INVALID_MODE", Message: "execution.mode deve ser live, local ou replay", Blocking: true})
	}
	sort.Slice(intake.Provenance, func(i, j int) bool { return intake.Provenance[i].Path < intake.Provenance[j].Path })
	sort.Slice(intake.Issues, func(i, j int) bool { return intake.Issues[i].Path < intake.Issues[j].Path })
	sort.Slice(intake.Questions, func(i, j int) bool {
		if intake.Questions[i].Priority != intake.Questions[j].Priority {
			return intake.Questions[i].Priority < intake.Questions[j].Priority
		}
		return intake.Questions[i].ID < intake.Questions[j].ID
	})
	intake.InputHash = recommendationIntakeHash(intake.Request, intake.Execution)
}

func addRequiredNeed(intake *domain.RecommendationIntake, path, code string, kind domain.AnswerKind, subject string) {
	intake.Issues = append(intake.Issues, domain.InputIssue{Path: path, Code: code, Message: subject + " não informado", Blocking: true})
	intake.Questions = append(intake.Questions, domain.InformationNeed{
		ID: "initial:" + path, Phase: domain.NeedInitial, Subject: subject, TargetPath: path,
		ReasonCode: code, Resolver: domain.ResolverAskUser, AnswerKind: kind,
		ValidationRule:  "valor não vazio confirmado pelo usuário",
		DecisionEffects: []domain.DecisionEffect{{Path: path, Description: "completa a intenção da recomendação"}},
		Blocking:        true, Priority: 1,
	})
}

func addProvenance(intake *domain.RecommendationIntake, path string, origin domain.InputOrigin, authority domain.InputAuthority, source string) {
	intake.Provenance = append(intake.Provenance, domain.InputProvenance{Path: path, Origin: origin, Authority: authority, SourceRef: source})
}

func ApplyInformationAnswers(original domain.RecommendationIntake, answers []domain.InformationAnswer) (domain.RecommendationIntake, error) {
	encoded, _ := json.Marshal(original)
	var intake domain.RecommendationIntake
	_ = json.Unmarshal(encoded, &intake)
	open := map[string]domain.InformationNeed{}
	for _, need := range intake.Questions {
		open[need.ID] = need
	}
	for _, answer := range answers {
		need, ok := open[answer.NeedID]
		if !ok {
			return domain.RecommendationIntake{}, fmt.Errorf("necessidade %q desconhecida ou já respondida", answer.NeedID)
		}
		if answer.Authority != domain.AuthorityUser || !answer.Confirmed {
			return domain.RecommendationIntake{}, fmt.Errorf("necessidade %q exige confirmação da autoridade USER", answer.NeedID)
		}
		if answer.AnswerKind != need.AnswerKind {
			return domain.RecommendationIntake{}, fmt.Errorf("necessidade %q espera resposta %s", answer.NeedID, need.AnswerKind)
		}
		if err := applyInformationAnswer(&intake, need, answer); err != nil {
			return domain.RecommendationIntake{}, err
		}
		intake.Answers = append(intake.Answers, answer)
		delete(open, answer.NeedID)
		markAnsweredProvenance(&intake, need.TargetPath, answer.SourceRef)
	}
	assessPreparedIntake(&intake)
	return intake, nil
}

func applyInformationAnswer(intake *domain.RecommendationIntake, need domain.InformationNeed, answer domain.InformationAnswer) error {
	switch need.TargetPath {
	case "idea.problem":
		var value string
		if json.Unmarshal(answer.Value, &value) != nil || strings.TrimSpace(value) == "" {
			return fmt.Errorf("resposta para %q deve ser texto não vazio", need.ID)
		}
		intake.Request.Idea.Problem = value
		if provenanceOrigin(intake.Provenance, "target.name") == domain.InputDerived {
			intake.Request.Target = DefaultRecommendationTargetV1(value)
		}
	case "idea.capabilities":
		var value []string
		if json.Unmarshal(answer.Value, &value) != nil {
			return fmt.Errorf("resposta para %q deve ser uma lista de textos", need.ID)
		}
		value = normalizedIntakeStrings(value)
		if len(value) == 0 {
			return fmt.Errorf("resposta para %q deve conter ao menos uma capacidade", need.ID)
		}
		intake.Request.Idea.Capabilities = value
	default:
		return fmt.Errorf("caminho %q não pode ser alterado por este aplicador", need.TargetPath)
	}
	return nil
}

func markAnsweredProvenance(intake *domain.RecommendationIntake, path, source string) {
	for i := range intake.Provenance {
		if intake.Provenance[i].Path == path {
			intake.Provenance[i].Origin = domain.InputExplicit
			intake.Provenance[i].Authority = domain.AuthorityUser
			intake.Provenance[i].SourceRef = source
			return
		}
	}
}

func provenanceOrigin(values []domain.InputProvenance, path string) domain.InputOrigin {
	for _, value := range values {
		if value.Path == path {
			return value.Origin
		}
	}
	return ""
}

func CompileInconclusiveNeeds(diagnosis domain.InconclusiveDiagnosis, intake domain.RecommendationIntake) domain.InformationNeedReport {
	needs := make([]domain.InformationNeed, 0, len(diagnosis.Findings))
	for _, finding := range diagnosis.Findings {
		resolver, kind, priority := domain.ResolverCollect, domain.AnswerKind(""), 3
		switch finding.Cause {
		case domain.CauseMissingTest:
			if !intake.Execution.AllowCandidateExecution {
				resolver, kind, priority = domain.ResolverRequestAuthorization, domain.AnswerBoolean, 2
			}
		case domain.CauseConflictingEvidence:
			resolver, kind, priority = domain.ResolverRequireRiskOwner, domain.AnswerText, 1
		case domain.CauseUnsupportedArtifact, domain.CauseRuleCoverageGap, domain.CauseInsufficientCorpus, domain.CauseBudgetExhausted:
			resolver, priority = domain.ResolverTerminal, 4
		case domain.CauseUnresolvedRevision:
			resolver, priority = domain.ResolverCompute, 1
		}
		id := strings.TrimSpace(finding.ID)
		if id == "" {
			id = shortIntakeHash(finding)
		}
		need := domain.InformationNeed{
			ID: "execution:" + id, Phase: domain.NeedExecution, Subject: finding.Message,
			ClaimID: finding.RuleID, ReasonCode: string(finding.Cause), Resolver: resolver, AnswerKind: kind,
			DecisionEffects: []domain.DecisionEffect{}, Blocking: true, Priority: priority,
		}
		if resolver == domain.ResolverRequestAuthorization {
			need.TargetPath = "execution.allow_candidate_execution"
			need.ValidationRule = "booleano confirmado pelo usuário antes de qualquer execução"
			need.DecisionEffects = []domain.DecisionEffect{{Path: need.TargetPath, Description: "autoriza coleta futura; não comprova resultado"}}
		}
		needs = append(needs, need)
	}
	sort.Slice(needs, func(i, j int) bool {
		if needs[i].Priority != needs[j].Priority {
			return needs[i].Priority < needs[j].Priority
		}
		return needs[i].ID < needs[j].ID
	})
	return domain.InformationNeedReport{SchemaVersion: "1", IntakeHash: intake.InputHash, Needs: needs}
}

func normalizedIntakeStrings(values []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func recommendationIntakeHash(request domain.RecommendationRequest, execution domain.RecommendationExecutionProfile) string {
	return shortIntakeHash(struct {
		SchemaVersion string                                `json:"schema_version"`
		Request       domain.RecommendationRequest          `json:"request"`
		Execution     domain.RecommendationExecutionProfile `json:"execution"`
	}{domain.RecommendationIntakeSchemaVersion, request, execution})
}

func shortIntakeHash(value any) string {
	encoded, _ := json.Marshal(value)
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}
