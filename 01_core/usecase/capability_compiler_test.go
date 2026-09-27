/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/capability_compiler_test.md
 * @layer L1
 * @updated 2026-09-24
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"encoding/json"
	"strings"
	"testing"
)

func TestCompileCapabilityPlanQuestionsAndAuthority(t *testing.T) {
	problem := "buscar documentos privados"
	capabilities := []string{"buscar trechos"}
	qualities := []string{"resposta em até 2 segundos"}
	conditions := []string{"sem acesso à rede"}
	nonGoals := []string{"sincronização em nuvem"}
	intake := usecase.PrepareRecommendationIntake(domain.RecommendationIntakeDraft{
		Idea: &domain.IdeaIntakeInput{
			Problem: &problem, Capabilities: &capabilities, Qualities: &qualities,
			OperatingConditions: &conditions, NonGoals: &nonGoals,
		},
	})
	before, err := json.Marshal(intake)
	if err != nil {
		t.Fatal(err)
	}
	out := usecase.CompileCapabilityPlan(intake)
	after, _ := json.Marshal(intake)
	if string(before) != string(after) {
		t.Fatal("compiler mutated intake")
	}
	if out.IntakeHash != intake.InputHash || out.IntakeStatus != domain.IntakeReady {
		t.Fatalf("lost intake identity: %+v", out)
	}
	if len(out.Requirements) != 6 || len(out.Questions) != 16 {
		t.Fatalf("unexpected compilation size: %d requirements, %d questions", len(out.Requirements), len(out.Questions))
	}
	kinds := map[string]domain.CompiledRequirement{}
	for _, r := range out.Requirements {
		kinds[r.Kind] = r
	}
	for _, kind := range []string{"capability", "quality", "operating_condition", "non_goal", "license_policy", "archive_policy"} {
		if kinds[kind].ID == "" {
			t.Fatalf("missing requirement kind %s: %+v", kind, out.Requirements)
		}
	}
	if kinds["capability"].Authority != domain.AuthorityUser || kinds["quality"].Authority != domain.AuthorityUser {
		t.Fatalf("user authority lost: %+v", kinds)
	}
	if kinds["license_policy"].Authority != domain.AuthorityProductPolicy {
		t.Fatalf("default license policy misattributed to user: %+v", kinds["license_policy"])
	}
	byKind := map[string]map[string]bool{}
	byID := map[string]string{}
	for _, r := range out.Requirements {
		byID[r.ID] = r.Kind
	}
	for _, q := range out.Questions {
		if q.Status != "PENDING" || q.Question == "" || q.Support == "" || q.Refute == "" || q.Limit == "" {
			t.Fatalf("question masquerades as evidence or lacks limits: %+v", q)
		}
		kind := byID[q.RequirementID]
		if kind == "" {
			t.Fatalf("orphan question: %+v", q)
		}
		if byKind[kind] == nil {
			byKind[kind] = map[string]bool{}
		}
		byKind[kind][q.Channel] = true
	}
	if !byKind["capability"]["documentation"] || !byKind["capability"]["code"] || !byKind["capability"]["tests"] || !byKind["capability"]["execution"] {
		t.Fatalf("capability channel gap: %+v", byKind["capability"])
	}
	if byKind["quality"]["code"] || !byKind["quality"]["execution"] {
		t.Fatalf("quality should require observation, not static proof: %+v", byKind["quality"])
	}
	if !byKind["license_policy"]["license"] || len(byKind["license_policy"]) != 1 {
		t.Fatalf("license is not independent: %+v", byKind["license_policy"])
	}
	if !byKind["archive_policy"]["maintenance"] {
		t.Fatal("maintenance policy not compiled")
	}
	for _, q := range out.Questions {
		if byID[q.RequirementID] == "non_goal" && !strings.Contains(q.Limit, "não prova ausência") {
			t.Fatalf("non-goal absence caveat missing: %+v", q)
		}
	}
}

func TestCompileCapabilityPlanStableIDsAndExplicitPolicy(t *testing.T) {
	problem := "analisar texto"
	capabilities := []string{"  Buscar  trechos "}
	requireLocal := true
	intake := usecase.PrepareRecommendationIntake(domain.RecommendationIntakeDraft{
		Idea:   &domain.IdeaIntakeInput{Problem: &problem, Capabilities: &capabilities},
		Policy: &domain.PolicyIntakeInput{RequireLocal: &requireLocal},
	})
	a := usecase.CompileCapabilityPlan(intake)
	b := usecase.CompileCapabilityPlan(intake)
	encodedA, _ := json.Marshal(a)
	encodedB, _ := json.Marshal(b)
	if string(encodedA) != string(encodedB) {
		t.Fatal("non-deterministic compilation")
	}
	found := false
	for _, r := range a.Requirements {
		if r.Kind == "local_policy" {
			found = true
			if r.Authority != domain.AuthorityUser {
				t.Fatalf("explicit policy authority lost: %+v", r)
			}
		}
	}
	if !found {
		t.Fatal("require_local did not generate a verification obligation")
	}
}

func TestAdvancedQualityCannotBeSilentlyIgnoredByDecision(t *testing.T) {
	idea := domain.IdeaSpec{
		Problem: "buscar documentos", Capabilities: []string{"buscar trechos"},
		Qualities:           []string{"resposta em até 2 segundos"},
		OperatingConditions: []string{"sem rede"},
	}
	report, err := usecase.DecideReuse(domain.ReuseDecisionRequest{Idea: idea, Policy: usecase.DefaultRecommendationPolicyV1()})
	if err != nil {
		t.Fatal(err)
	}
	if report.Decision != domain.DecisionUnknown || len(report.Gaps) != 2 {
		t.Fatalf("unverified requirements were ignored: %+v", report)
	}
	foundQuality, foundCondition, foundReason := false, false, false
	for _, r := range report.Investigation.Requirements {
		foundQuality = foundQuality || r.Kind == "quality"
		foundCondition = foundCondition || r.Kind == "operating_condition"
	}
	for _, r := range report.Reasons {
		foundReason = foundReason || r.RuleID == "DEC-REQUIREMENTS-0087"
	}
	if !foundQuality || !foundCondition || !foundReason {
		t.Fatalf("missing explicit abstention path: %+v", report)
	}
}

func TestCompileCapabilityPlanDoesNotInventAuthority(t *testing.T) {
	intake := domain.RecommendationIntake{Request: domain.RecommendationRequest{
		Idea: domain.IdeaSpec{Capabilities: []string{"buscar trechos"}},
	}}
	out := usecase.CompileCapabilityPlan(intake)
	if len(out.Requirements) != 1 || out.Requirements[0].Authority != domain.AuthorityUnresolved {
		t.Fatalf("invented requirement authority: %+v", out.Requirements)
	}
}
