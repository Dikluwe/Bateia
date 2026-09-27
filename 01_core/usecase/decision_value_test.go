/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_value_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func decisionValueFixture() domain.DecisionValueRequest {
	pivotal := decisionAssessment("pivotal", domain.EvidenceProven, domain.EvidenceProven, domain.EvidenceInconclusive)
	neutral := decisionAssessment("neutral", domain.EvidenceAbsent, domain.EvidenceAbsent, domain.EvidenceAbsent)
	neutral.License = domain.LicenseAssessment{Status: domain.EvidenceInconclusive}
	policy := decisionPolicy()
	return domain.DecisionValueRequest{
		Idea:           domain.IdeaSpec{Problem: "x", Capabilities: []string{"core", "export", "sync"}},
		Assessment:     domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{pivotal, neutral}},
		DecisionPolicy: policy, RecoveryPolicy: DefaultRecoveryPolicy(), Round: 1,
		Diagnosis: domain.InconclusiveDiagnosis{SchemaVersion: "1", Findings: []domain.InconclusiveFinding{
			{ID: "missing-test", Cause: domain.CauseMissingTest, Candidate: pivotal.Repository.Key(), Capability: "sync", RuleID: "cap-sync", Recoverable: true},
			{ID: "license", Cause: domain.CauseUnknownLicense, Candidate: neutral.Repository.Key(), RuleID: "GATE-LICENSE-001", Recoverable: true},
		}},
	}
}

func TestDecisionValuePrioritizesPivotalActionOverNeutralGate(t *testing.T) {
	request := decisionValueFixture()
	plan, err := PlanDecisionValue(request)
	if err != nil {
		t.Fatal(err)
	}
	if plan.StopReason != "ACTION_READY" || plan.RecommendedActionID == "" || len(plan.ExecutablePlan.Targets) != 1 {
		t.Fatalf("plano não ficou operacional: %+v", plan)
	}
	selected := plan.ExecutablePlan.Targets[0]
	if selected.FindingID != "missing-test" || selected.Path == "@compatibility-test" {
		t.Fatalf("ação errada selecionada: %+v", selected)
	}
	seenNeutral := false
	for _, option := range plan.Options {
		if option.FindingID == "license" {
			seenNeutral = true
			if option.Impact != domain.DecisionNeutral || !reflect.DeepEqual(option.PossibleDecisions, []domain.ReuseDecision{domain.DecisionUnknown}) {
				t.Fatalf("gate irrelevante deveria ser neutro: %+v", option)
			}
		}
	}
	if !seenNeutral || plan.Options[0].Impact != domain.DecisionPivotal || plan.Options[0].ResolvedOutcomeCount == 0 {
		t.Fatalf("classificação inesperada: %+v", plan.Options)
	}
}

func TestDecisionValueIsInvariantToFindingOrder(t *testing.T) {
	request := decisionValueFixture()
	first, err := PlanDecisionValue(request)
	if err != nil {
		t.Fatal(err)
	}
	request.Diagnosis.Findings[0], request.Diagnosis.Findings[1] = request.Diagnosis.Findings[1], request.Diagnosis.Findings[0]
	second, err := PlanDecisionValue(request)
	if err != nil {
		t.Fatal(err)
	}
	if first.RecommendedActionID != second.RecommendedActionID || !reflect.DeepEqual(first.Options, second.Options) {
		t.Fatal("ordem dos findings alterou valor decisório")
	}
}

func TestDecisionValueAsksBeforeCompatibilityExecution(t *testing.T) {
	request := decisionValueFixture()
	request.Diagnosis.Findings = request.Diagnosis.Findings[:1]
	finding := request.Diagnosis.Findings[0]
	ref, err := parseRepositoryReference(finding.Candidate)
	if err != nil {
		t.Fatal(err)
	}
	request.PreviousTargets = targetsForFinding(finding, ref, request.RecoveryPolicy)
	plan, err := PlanDecisionValue(request)
	if err != nil {
		t.Fatal(err)
	}
	if plan.StopReason != "AWAITING_CANDIDATE_EXECUTION_AUTHORIZATION" || len(plan.Questions) != 1 || len(plan.ExecutablePlan.Targets) != 0 {
		t.Fatalf("execução sem autorização não foi bloqueada: %+v", plan)
	}
	if len(plan.Options) != 1 || plan.Options[0].Method != domain.InvestigationRunCompatibilityTest || !plan.Options[0].RequiresAuthorization {
		t.Fatalf("opção de teste não foi explicitada: %+v", plan.Options)
	}
}

func TestDecisionValueDoesNotInventImpactForUnresolvedRevision(t *testing.T) {
	request := decisionValueFixture()
	request.Assessment.Assessments[0].Repository.Commit = "HEAD"
	request.Diagnosis.Findings = []domain.InconclusiveFinding{{ID: "revision", Cause: domain.CauseUnresolvedRevision, Candidate: request.Assessment.Assessments[0].Repository.Key(), RuleID: "GATE-REVISION-001", Recoverable: true}}
	plan, err := PlanDecisionValue(request)
	if err != nil {
		t.Fatal(err)
	}
	if plan.StopReason != "NO_DECISION_PIVOTAL_ACTION" || len(plan.Options) != 1 || plan.Options[0].Impact != domain.DecisionCounterfactualUnsupported {
		t.Fatalf("revisão opaca recebeu valor fabricado: %+v", plan)
	}
}

func TestProjectedDecisionDoesNotMutateObservedAssessment(t *testing.T) {
	request := decisionValueFixture()
	before := request.Assessment.Assessments[0].Capabilities[2].Status
	if _, err := PlanDecisionValue(request); err != nil {
		t.Fatal(err)
	}
	if request.Assessment.Assessments[0].Capabilities[2].Status != before {
		t.Fatal("simulação alterou assessment observado")
	}
}

func TestDecisionValueBranchesLicenseMetadataIntoAllowedAndDisallowedOutcomes(t *testing.T) {
	finding := domain.InconclusiveFinding{Cause: domain.CauseUnknownLicense}
	target := domain.RecoveryTarget{Kind: domain.EvidenceRepository, Expected: []domain.EvidenceStatus{domain.EvidenceProven}}
	hypotheses := decisionHypotheses(finding, target, decisionPolicy())
	if len(hypotheses) != 2 || hypotheses[0].code != "PROVEN_ALLOWED" || hypotheses[1].code != "PROVEN_NOT_ALLOWED" {
		t.Fatalf("metadados de repositório não ramificaram a licença: %+v", hypotheses)
	}
}
