/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_value.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"fmt"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

const counterfactualDecisionMethod = "COUNTERFACTUAL_DECISION_REACHABILITY_V1"

// PlanDecisionValue determines which recovery action can cross a real decision
// boundary. It deliberately does not estimate probabilities, utilities or EVSI.
func PlanDecisionValue(request domain.DecisionValueRequest) (domain.DecisionValuePlan, error) {
	if err := validateDecisionValueRequest(request); err != nil {
		return domain.DecisionValuePlan{}, err
	}
	baseline, err := DecideReuse(domain.ReuseDecisionRequest{Idea: request.Idea, Assessments: request.Assessment.Assessments, Policy: request.DecisionPolicy})
	if err != nil {
		return domain.DecisionValuePlan{}, err
	}
	if baseline.Decision != domain.DecisionUnknown {
		return domain.DecisionValuePlan{}, fmt.Errorf("valor de investigação exige decisão INCONCLUSIVE; atual é %s", baseline.Decision)
	}
	plan := domain.DecisionValuePlan{
		SchemaVersion: domain.DecisionValueSchemaVersion, Method: counterfactualDecisionMethod,
		BaselineDecision: baseline.Decision, Options: []domain.InvestigationValueOption{}, Questions: []domain.InformationNeed{},
		ExecutablePlan: domain.EvidenceRecoveryPlan{
			SchemaVersion: "1", Round: request.Round, DiagnosisHash: canonicalDiagnosisHash(request.Diagnosis),
			Policy: request.RecoveryPolicy, PolicyHash: recoveryHash(request.RecoveryPolicy), Targets: []domain.RecoveryTarget{}, Skipped: []domain.InconclusiveFinding{},
		},
	}
	findings := map[string]domain.InconclusiveFinding{}
	for _, finding := range request.Diagnosis.Findings {
		findings[finding.ID] = finding
	}
	for _, target := range decisionValueTargets(request) {
		finding, ok := findings[target.FindingID]
		if !ok {
			return domain.DecisionValuePlan{}, fmt.Errorf("target sem finding correspondente: %s", target.FindingID)
		}
		option := evaluateInvestigationOption(request, baseline.Decision, finding, target)
		plan.Options = append(plan.Options, option)
	}
	markCostDominance(plan.Options)
	sort.Slice(plan.Options, func(i, j int) bool { return investigationOptionLess(plan.Options[i], plan.Options[j]) })
	selected := firstPivotalOption(plan.Options)
	if selected == nil {
		if len(plan.Options) == 0 {
			plan.StopReason = "NO_RECOVERABLE_ACTION"
		} else {
			plan.StopReason = "NO_DECISION_PIVOTAL_ACTION"
		}
		return plan, nil
	}
	plan.RecommendedActionID = selected.ActionID
	if selected.RequiresAuthorization {
		plan.StopReason = "AWAITING_CANDIDATE_EXECUTION_AUTHORIZATION"
		plan.Questions = append(plan.Questions, candidateExecutionQuestion(*selected))
		return plan, nil
	}
	if !selected.Executable {
		plan.StopReason = "ACTION_REQUIRES_EXECUTOR"
		return plan, nil
	}
	plan.ExecutablePlan.Targets = []domain.RecoveryTarget{selected.Target}
	plan.StopReason = "ACTION_READY"
	return plan, nil
}

func validateDecisionValueRequest(request domain.DecisionValueRequest) error {
	if request.Round < 1 || request.Round > request.RecoveryPolicy.MaxRounds {
		return fmt.Errorf("rodada %d excede max_rounds", request.Round)
	}
	if request.RecoveryPolicy.SchemaVersion == "" || request.RecoveryPolicy.MaxCandidates < 1 || request.RecoveryPolicy.AdditionalBytes < 1 {
		return fmt.Errorf("política de recuperação inválida")
	}
	if request.Diagnosis.SchemaVersion == "" {
		return fmt.Errorf("diagnosis.schema_version é obrigatório")
	}
	return validateDecisionRequest(domain.ReuseDecisionRequest{Idea: request.Idea, Assessments: request.Assessment.Assessments, Policy: request.DecisionPolicy})
}

func decisionValueTargets(request domain.DecisionValueRequest) []domain.RecoveryTarget {
	seen := map[string]bool{}
	for _, target := range request.PreviousTargets {
		seen[targetIdentity(target)] = true
	}
	findings := append([]domain.InconclusiveFinding(nil), request.Diagnosis.Findings...)
	sort.Slice(findings, func(i, j int) bool { return findingSortKey(findings[i]) < findingSortKey(findings[j]) })
	result := []domain.RecoveryTarget{}
	for _, finding := range findings {
		if !finding.Recoverable {
			continue
		}
		ref, err := parseRepositoryReference(finding.Candidate)
		if err != nil {
			continue
		}
		options := targetsForFinding(finding, ref, request.RecoveryPolicy)
		if finding.Cause == domain.CauseMissingTest {
			options = append(options, compatibilityTestTarget(finding, ref))
		}
		for _, target := range options {
			identity := targetIdentity(target)
			if !seen[identity] {
				seen[identity] = true
				result = append(result, target)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return targetIdentity(result[i]) < targetIdentity(result[j]) })
	return result
}

func compatibilityTestTarget(finding domain.InconclusiveFinding, ref domain.CommitRef) domain.RecoveryTarget {
	return domain.RecoveryTarget{
		FindingID: finding.ID, Candidate: ref, Capability: finding.Capability,
		Path: "@compatibility-test", Kind: domain.EvidenceTestImplementation, RuleID: finding.RuleID,
		Expected: []domain.EvidenceStatus{domain.EvidenceProven, domain.EvidenceAbsent, domain.EvidenceConflict, domain.EvidenceInconclusive},
		Reason:   "executar teste delimitado que confronte a capacidade com o alvo",
	}
}

func evaluateInvestigationOption(request domain.DecisionValueRequest, baseline domain.ReuseDecision, finding domain.InconclusiveFinding, target domain.RecoveryTarget) domain.InvestigationValueOption {
	method := investigationMethod(finding, target)
	option := domain.InvestigationValueOption{
		ActionID: recoveryHash([]string{targetIdentity(target), string(method)})[:16], FindingID: finding.ID,
		Target: target, Method: method, EpistemicBasis: "EXACT_DECISION_REPLAY_WITH_DECLARED_OUTCOMES",
		Outcomes: []domain.CounterfactualOutcome{}, PossibleDecisions: []domain.ReuseDecision{}, Cost: investigationCost(target),
		Executable: targetExecutable(target, request.RecoveryPolicy),
	}
	option.RequiresAuthorization = option.Cost.CandidateExecution && !request.AllowCandidateExecution
	option.SecurityRisk = domain.InvestigationBoundedRead
	option.SecurityRationale = "ler artefato não confiável com limites; conteúdo não concede autorização"
	if option.Cost.CandidateExecution {
		option.SecurityRisk = domain.InvestigationHostExecution
		option.SecurityRationale = "executar código de terceiros exige autorização e isolamento; executor seguro indisponível"
	}
	if !counterfactualCauseSupported(finding) {
		option.Impact = domain.DecisionCounterfactualUnsupported
		option.EpistemicBasis = "NO_SAFE_COUNTERFACTUAL_TRANSFORM"
		option.Reason = "o resultado exige reconstruir identidade ou coleta, não apenas alterar um estado decisório"
		return option
	}
	hypotheses := decisionHypotheses(finding, target, request.DecisionPolicy)
	decisions := map[domain.ReuseDecision]bool{}
	for _, hypothesis := range hypotheses {
		decision, err := projectedDecision(request, finding, target, hypothesis)
		if err != nil {
			option.Impact = domain.DecisionCounterfactualUnsupported
			option.EpistemicBasis = "COUNTERFACTUAL_REPLAY_FAILED"
			option.Reason = err.Error()
			option.Outcomes = []domain.CounterfactualOutcome{}
			return option
		}
		option.Outcomes = append(option.Outcomes, domain.CounterfactualOutcome{OutcomeCode: hypothesis.code, Status: hypothesis.status, Decision: decision})
		decisions[decision] = true
		if decision != domain.DecisionUnknown {
			option.ResolvedOutcomeCount++
		}
	}
	for decision := range decisions {
		option.PossibleDecisions = append(option.PossibleDecisions, decision)
	}
	sort.Slice(option.PossibleDecisions, func(i, j int) bool { return option.PossibleDecisions[i] < option.PossibleDecisions[j] })
	option.Impact = domain.DecisionNeutral
	for _, outcome := range option.Outcomes {
		if outcome.Decision != baseline {
			option.Impact = domain.DecisionPivotal
			break
		}
	}
	if option.Impact == domain.DecisionPivotal {
		option.Reason = fmt.Sprintf("resultados declarados alcançam %s", joinDecisions(option.PossibleDecisions))
	} else {
		option.Reason = fmt.Sprintf("todos os resultados mantêm %s", baseline)
	}
	if method == domain.InvestigationRunCompatibilityTest {
		option.Executable = false
		option.Reason += "; o executor de testes de terceiros não faz parte deste passo"
	}
	return option
}

type decisionHypothesis struct {
	code      string
	status    domain.EvidenceStatus
	licenseID string
}

func decisionHypotheses(finding domain.InconclusiveFinding, target domain.RecoveryTarget, policy domain.ReuseDecisionPolicy) []decisionHypothesis {
	seen := map[domain.EvidenceStatus]bool{}
	statuses := []domain.EvidenceStatus{}
	for _, status := range target.Expected {
		if !seen[status] {
			seen[status] = true
			statuses = append(statuses, status)
		}
	}
	sort.Slice(statuses, func(i, j int) bool { return statuses[i] < statuses[j] })
	result := []decisionHypothesis{}
	for _, status := range statuses {
		if finding.Cause == domain.CauseUnknownLicense && status == domain.EvidenceProven {
			allowed := "COUNTERFACTUAL_ALLOWED"
			if len(policy.AllowedLicenses) > 0 {
				allowed = normalizedLicenses(policy.AllowedLicenses)[0]
			}
			result = append(result,
				decisionHypothesis{code: "PROVEN_ALLOWED", status: status, licenseID: allowed},
				decisionHypothesis{code: "PROVEN_NOT_ALLOWED", status: status, licenseID: "COUNTERFACTUAL_NOT_ALLOWED"},
			)
			continue
		}
		result = append(result, decisionHypothesis{code: string(status), status: status})
	}
	return result
}

func projectedDecision(request domain.DecisionValueRequest, finding domain.InconclusiveFinding, target domain.RecoveryTarget, hypothesis decisionHypothesis) (domain.ReuseDecision, error) {
	matched := false
	applied := false
	projection := func(assessment domain.CandidateAssessment) domain.CandidateAssessment {
		if reviewRepositoryKey(assessment.Repository.Owner, assessment.Repository.Repo) != reviewRepositoryKey(target.Candidate.Owner, target.Candidate.Repo) {
			return assessment
		}
		matched = true
		switch finding.Cause {
		case domain.CauseUnknownLicense:
			assessment.License.Status = hypothesis.status
			assessment.License.Identifier = hypothesis.licenseID
			applied = true
		case domain.CauseUnknownOperation:
			assessment.Operation.LocalStatus = hypothesis.status
			applied = true
		case domain.CauseMissingImplementation, domain.CauseMissingTest, domain.CauseIncompleteCollection, domain.CauseBudgetExhausted, domain.CauseProviderFailure:
			for index := range assessment.Capabilities {
				if normalizeDecisionText(assessment.Capabilities[index].Capability) == normalizeDecisionText(finding.Capability) {
					assessment.Capabilities[index].Status = hypothesis.status
					applied = true
				}
			}
			if finding.Cause == domain.CauseProviderFailure && hypothesis.status != domain.EvidenceInconclusive {
				assessment.Failures = []domain.ExecutionFailure{}
			}
		}
		return assessment
	}
	report, err := decideReuseProjected(domain.ReuseDecisionRequest{Idea: request.Idea, Assessments: request.Assessment.Assessments, Policy: request.DecisionPolicy}, projection)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("candidato do target não existe no assessment")
	}
	if !applied {
		return "", fmt.Errorf("claim do finding não existe no assessment")
	}
	return report.Decision, nil
}

func counterfactualCauseSupported(finding domain.InconclusiveFinding) bool {
	switch finding.Cause {
	case domain.CauseUnknownLicense, domain.CauseUnknownOperation, domain.CauseMissingImplementation,
		domain.CauseMissingTest, domain.CauseIncompleteCollection, domain.CauseBudgetExhausted, domain.CauseProviderFailure:
		return finding.Cause == domain.CauseUnknownLicense || finding.Cause == domain.CauseUnknownOperation || strings.TrimSpace(finding.Capability) != ""
	default:
		return false
	}
}

func investigationMethod(finding domain.InconclusiveFinding, target domain.RecoveryTarget) domain.InvestigationMethod {
	if target.Path == "@compatibility-test" {
		return domain.InvestigationRunCompatibilityTest
	}
	if finding.Cause == domain.CauseProviderFailure {
		return domain.InvestigationRetryProvider
	}
	return domain.InvestigationInspectArtifact
}

func investigationCost(target domain.RecoveryTarget) domain.InvestigationCost {
	cost := domain.InvestigationCost{}
	switch target.Kind {
	case domain.EvidenceRepository:
		cost.ProviderCalls = 1
	case domain.EvidenceDocumentation, domain.EvidenceManifest, domain.EvidenceLicense:
		cost.Documents = 1
	case domain.EvidenceFile, domain.EvidenceSymbol:
		cost.CodeFiles = 1
	case domain.EvidenceTest:
		cost.CodeFiles, cost.Tests = 1, 1
	case domain.EvidenceTestImplementation:
		cost.CodeFiles, cost.Tests, cost.CandidateExecution = 1, 1, true
	}
	return cost
}

func targetExecutable(target domain.RecoveryTarget, policy domain.RecoveryPolicy) bool {
	switch target.Kind {
	case domain.EvidenceDocumentation, domain.EvidenceManifest, domain.EvidenceLicense:
		return policy.AdditionalDocuments > 0
	case domain.EvidenceFile, domain.EvidenceSymbol:
		return policy.AdditionalCodeFiles > 0
	case domain.EvidenceTest:
		return policy.AdditionalTests > 0
	case domain.EvidenceRepository:
		return true
	case domain.EvidenceTestImplementation:
		return false
	default:
		return false
	}
}

func investigationOptionLess(left, right domain.InvestigationValueOption) bool {
	impactRank := map[domain.DecisionImpactClass]int{domain.DecisionPivotal: 0, domain.DecisionNeutral: 1, domain.DecisionCounterfactualUnsupported: 2}
	if impactRank[left.Impact] != impactRank[right.Impact] {
		return impactRank[left.Impact] < impactRank[right.Impact]
	}
	if left.SecurityRisk != right.SecurityRisk {
		return left.SecurityRisk == domain.InvestigationBoundedRead
	}
	if left.ResolvedOutcomeCount != right.ResolvedOutcomeCount {
		return left.ResolvedOutcomeCount > right.ResolvedOutcomeCount
	}
	if len(left.PossibleDecisions) != len(right.PossibleDecisions) {
		return len(left.PossibleDecisions) > len(right.PossibleDecisions)
	}
	if left.RequiresAuthorization != right.RequiresAuthorization {
		return !left.RequiresAuthorization
	}
	if left.Executable != right.Executable {
		return left.Executable
	}
	if left.CostParetoDominated != right.CostParetoDominated {
		return !left.CostParetoDominated
	}
	return left.ActionID < right.ActionID
}

func markCostDominance(options []domain.InvestigationValueOption) {
	for left := range options {
		for right := range options {
			if left == right || options[left].Impact != options[right].Impact || options[left].ResolvedOutcomeCount != options[right].ResolvedOutcomeCount || len(options[left].PossibleDecisions) != len(options[right].PossibleDecisions) || options[left].RequiresAuthorization != options[right].RequiresAuthorization || options[left].Executable != options[right].Executable {
				continue
			}
			if costDominates(options[right].Cost, options[left].Cost) {
				options[left].CostParetoDominated = true
				break
			}
		}
	}
}

func costDominates(left, right domain.InvestigationCost) bool {
	leftExecution, rightExecution := 0, 0
	if left.CandidateExecution {
		leftExecution = 1
	}
	if right.CandidateExecution {
		rightExecution = 1
	}
	leftValues := []int{leftExecution, left.ProviderCalls, left.Documents, left.CodeFiles, left.Tests}
	rightValues := []int{rightExecution, right.ProviderCalls, right.Documents, right.CodeFiles, right.Tests}
	strict := false
	for index := range leftValues {
		if leftValues[index] > rightValues[index] {
			return false
		}
		strict = strict || leftValues[index] < rightValues[index]
	}
	return strict
}

func firstPivotalOption(options []domain.InvestigationValueOption) *domain.InvestigationValueOption {
	for index := range options {
		if options[index].Impact == domain.DecisionPivotal {
			return &options[index]
		}
	}
	return nil
}

func joinDecisions(decisions []domain.ReuseDecision) string {
	values := make([]string, len(decisions))
	for index, decision := range decisions {
		values[index] = string(decision)
	}
	return strings.Join(values, ", ")
}

func candidateExecutionQuestion(option domain.InvestigationValueOption) domain.InformationNeed {
	return domain.InformationNeed{
		ID: "decision-value-execution-" + option.ActionID, Phase: domain.NeedExecution,
		Subject:    "Autorizar teste de compatibilidade no candidato " + option.Target.Candidate.Key(),
		TargetPath: "allow_candidate_execution", ReasonCode: "CANDIDATE_EXECUTION_AUTHORIZATION_REQUIRED",
		Resolver: domain.ResolverRequestAuthorization, AnswerKind: domain.AnswerBoolean,
		DecisionEffects: []domain.DecisionEffect{{Path: "recommended_action_id", Description: "a autorização permite encaminhar a investigação pivotal para um executor isolado"}},
		Blocking:        true, Priority: 1,
	}
}
