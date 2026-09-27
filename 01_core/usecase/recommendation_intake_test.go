/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_intake_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"encoding/json"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func TestPrepareRecommendationIntakeAsksOnlyForMissingHumanIntent(t *testing.T) {
	first := PrepareRecommendationIntake(domain.RecommendationIntakeDraft{})
	second := PrepareRecommendationIntake(domain.RecommendationIntakeDraft{})
	if first.Status != domain.IntakeNeedsInput {
		t.Fatalf("status = %s", first.Status)
	}
	if len(first.Questions) != 2 || first.Questions[0].ID != "initial:idea.capabilities" || first.Questions[1].ID != "initial:idea.problem" {
		t.Fatalf("questions = %#v", first.Questions)
	}
	for _, need := range first.Questions {
		if need.Resolver != domain.ResolverAskUser {
			t.Fatalf("resolver = %s", need.Resolver)
		}
	}
	if originAt(first.Provenance, "idea.capabilities") != domain.InputUnresolved {
		t.Fatal("capacidade obrigatória ausente foi rotulada como default")
	}
	if first.InputHash != second.InputHash || !reflect.DeepEqual(first.Provenance, second.Provenance) {
		t.Fatal("preparo idêntico não foi determinístico")
	}
}

func TestPrepareRecommendationIntakeMaterializesDefaultsAndPreservesExplicitValues(t *testing.T) {
	problem, strategy, mode := "Selecionar componentes", "literal", "live"
	capabilities := []string{"buscar", "comparar"}
	emptyLicenses := []string{}
	requireLocal, resume := false, false
	zero := 0
	draft := domain.RecommendationIntakeDraft{
		Idea:      &domain.IdeaIntakeInput{Problem: &problem, Capabilities: &capabilities},
		Policy:    &domain.PolicyIntakeInput{AllowedLicenses: &emptyLicenses, RequireLocal: &requireLocal},
		Search:    &domain.SearchIntakeInput{Strategy: &strategy, MaxCandidates: &zero},
		Execution: &domain.ExecutionIntakeInput{Mode: &mode, Resume: &resume},
	}
	intake := PrepareRecommendationIntake(draft)
	if intake.Status != domain.IntakeInvalid {
		t.Fatalf("status = %s", intake.Status)
	}
	if intake.Request.Search.MaxCandidates != 0 || intake.Request.Policy.RequireLocal || intake.Execution.Resume {
		t.Fatal("zero ou false explícito foi substituído por default")
	}
	if originAt(intake.Provenance, "search.max_candidates") != domain.InputExplicit || originAt(intake.Provenance, "search.max_per_query") != domain.InputVersionedDefault {
		t.Fatalf("proveniência incorreta: %#v", intake.Provenance)
	}
}

func TestApplyInformationAnswersRequiresUserConfirmationAndProducesReadyIntake(t *testing.T) {
	initial := PrepareRecommendationIntake(domain.RecommendationIntakeDraft{})
	problem, _ := json.Marshal("Descobrir software correlato")
	capabilities, _ := json.Marshal([]string{"buscar", "comparar"})
	answers := []domain.InformationAnswer{
		{NeedID: "initial:idea.problem", AnswerKind: domain.AnswerText, Value: problem, Authority: domain.AuthorityUser, Confirmed: true},
		{NeedID: "initial:idea.capabilities", AnswerKind: domain.AnswerStringList, Value: capabilities, Authority: domain.AuthorityUser, Confirmed: true},
	}
	updated, err := ApplyInformationAnswers(initial, answers)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != domain.IntakeReady || len(updated.Questions) != 0 || len(updated.Answers) != 2 {
		t.Fatalf("updated = %#v", updated)
	}
	if updated.InputHash == initial.InputHash || updated.Request.Target.Name != "Descobrir software correlato" {
		t.Fatal("resposta não gerou nova entrada canônica")
	}

	answers[0].Authority = domain.AuthorityAIProposal
	if _, err := ApplyInformationAnswers(initial, answers[:1]); err == nil {
		t.Fatal("proposta de IA sem autoridade foi aceita")
	}
}

func TestCompileInconclusiveNeedsSeparatesMachineAuthorizationRiskAndTerminal(t *testing.T) {
	diagnosis := domain.InconclusiveDiagnosis{Findings: []domain.InconclusiveFinding{
		{ID: "test", Cause: domain.CauseMissingTest, Message: "teste ausente"},
		{ID: "license", Cause: domain.CauseUnknownLicense, Message: "licença desconhecida"},
		{ID: "conflict", Cause: domain.CauseConflictingEvidence, Message: "evidência conflitante"},
		{ID: "rule", Cause: domain.CauseRuleCoverageGap, Message: "regra ausente"},
	}}
	report := CompileInconclusiveNeeds(diagnosis, domain.RecommendationIntake{})
	resolvers := map[string]domain.ResolverKind{}
	for _, need := range report.Needs {
		resolvers[need.ID] = need.Resolver
	}
	if resolvers["execution:test"] != domain.ResolverRequestAuthorization || resolvers["execution:license"] != domain.ResolverCollect || resolvers["execution:conflict"] != domain.ResolverRequireRiskOwner || resolvers["execution:rule"] != domain.ResolverTerminal {
		t.Fatalf("resolvers = %#v", resolvers)
	}
	report = CompileInconclusiveNeeds(diagnosis, domain.RecommendationIntake{Execution: domain.RecommendationExecutionProfile{AllowCandidateExecution: true}})
	for _, need := range report.Needs {
		if need.ID == "execution:test" && need.Resolver != domain.ResolverCollect {
			t.Fatalf("teste autorizado deveria ser coleta, recebeu %s", need.Resolver)
		}
	}
}

func originAt(values []domain.InputProvenance, path string) domain.InputOrigin {
	for _, value := range values {
		if value.Path == path {
			return value.Origin
		}
	}
	return ""
}
