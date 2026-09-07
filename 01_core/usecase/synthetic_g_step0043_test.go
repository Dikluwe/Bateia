/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/synthetic_g_step0043_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"encoding/json"
	"reflect"
	"testing"
)

func TestStep0043SyntheticGPublicCrossDomain(t *testing.T) {
	type outcome struct {
		Domain            string `json:"domain"`
		Condition         string `json:"condition"`
		FalseProof        bool   `json:"false_proof"`
		CorrectAbstention bool   `json:"correct_abstention"`
		Actionable        bool   `json:"actionable"`
		Deterministic     bool   `json:"deterministic"`
	}
	outcomes := []outcome{}
	for _, domainCase := range []struct{ name, capability string }{{"rendering", "render OpenGL order independent transparency"}, {"recipe", "execute recipe substitutions with allergy constraints"}, {"async", "enforce bounded async queue backpressure"}} {
		for _, condition := range []string{"documentation_only", "irrelevant_candidate", "zero_candidates_recorded_search"} {
			idea := domain.IdeaSpec{Problem: domainCase.capability, Capabilities: []string{domainCase.capability}}
			text := "This software can " + domainCase.capability + "."
			if condition == "irrelevant_candidate" {
				text = "An astronomy catalogue of stellar spectra and galaxy distances."
			}
			bundle := step0043Bundle(text)
			snapshot := domain.DiscoverySnapshot{SchemaVersion: "1", Idea: idea, PerExpansion: []domain.ExpansionDiscovery{{Expansion: domain.SearchExpansion{ID: "public-q", Query: domainCase.capability, Covers: idea.Capabilities}, Provider: domain.ProviderDiscovery{Provider: "synthetic-public", Attempts: []domain.ProviderAttempt{{Query: domainCase.capability, Status: "COMPLETE"}}}}}}
			screening := domain.ScreeningReport{SchemaVersion: "1"}
			bundles := []domain.EvidenceBundle{}
			if condition != "zero_candidates_recorded_search" {
				bundles = append(bundles, bundle)
				screening.Decisions = []domain.ScreeningDecision{{Repository: "fixture/candidate", Decision: "ADVANCE"}}
				snapshot.Candidates = []domain.DiscoveredCandidate{{Repository: domain.RepositoryCandidate{Owner: bundle.Repository.Owner, Repo: bundle.Repository.Repo, CommitID: bundle.Repository.Commit}}}
				snapshot.PerExpansion[0].Provider.Attempts[0].ResultCount = 1
			}
			evaluate := func() (domain.AssessmentReport, domain.ReuseDecisionReport, domain.InvestigationPlan) {
				t.Helper()
				assessment, err := usecase.AssessEvidence(idea, screening, bundles, usecase.DefaultEvidenceRules())
				if err != nil {
					t.Fatal(err)
				}
				decision, err := usecase.DecideReuse(domain.ReuseDecisionRequest{Idea: idea, Assessments: assessment.Assessments, Policy: usecase.DefaultRecommendationPolicyV1()})
				if err != nil {
					t.Fatal(err)
				}
				plan := usecase.BuildInvestigationPlan(idea, snapshot, assessment.Assessments)
				return assessment, decision, plan
			}
			assessment, decision, plan := evaluate()
			result := outcome{Domain: domainCase.name, Condition: condition}
			result.FalseProof = decision.Decision == domain.DecisionFork
			for _, a := range assessment.Assessments {
				if a.Operation.LocalStatus == domain.EvidenceProven {
					result.FalseProof = true
				}
				for _, capability := range a.Capabilities {
					if capability.Status == domain.EvidenceProven {
						result.FalseProof = true
					}
				}
			}
			result.Actionable = plan.Scope != "" && len(plan.Actions) > 0
			requirements := map[string]bool{}
			for _, r := range plan.Requirements {
				requirements[r.ID] = true
			}
			for _, action := range plan.Actions {
				if action.Status == domain.EvidenceProven {
					result.FalseProof = true
				}
				if !requirements[action.RequirementID] || action.NextAction == "" || action.Acceptance == "" {
					result.Actionable = false
				}
			}
			result.CorrectAbstention = !result.FalseProof
			a2, d2, p2 := evaluate()
			result.Deterministic = reflect.DeepEqual(assessment, a2) && reflect.DeepEqual(decision, d2) && reflect.DeepEqual(plan, p2)
			outcomes = append(outcomes, result)
			if result.FalseProof || !result.CorrectAbstention || !result.Actionable || !result.Deterministic {
				t.Errorf("synthetic G contract: %+v", result)
			}
		}
	}
	encoded, _ := json.Marshal(outcomes)
	t.Logf("STEP0043_SYNTHETIC_G_RESULTS=%s", encoded)
}
