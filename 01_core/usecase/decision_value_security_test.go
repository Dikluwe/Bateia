/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_value_security_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func TestInvestigationPrefersBoundedReadBeforeHostExecution(t *testing.T) {
	read := domain.InvestigationValueOption{Impact: domain.DecisionPivotal, SecurityRisk: domain.InvestigationBoundedRead, ResolvedOutcomeCount: 1}
	run := domain.InvestigationValueOption{Impact: domain.DecisionPivotal, SecurityRisk: domain.InvestigationHostExecution, ResolvedOutcomeCount: 10}
	if !investigationOptionLess(read, run) || investigationOptionLess(run, read) {
		t.Fatal("more information authorized risk escalation")
	}
	req := decisionValueFixture()
	plan, err := PlanDecisionValue(req)
	if err != nil {
		t.Fatal(err)
	}
	for _, option := range plan.Options {
		if option.SecurityRisk == "" || option.SecurityRationale == "" {
			t.Fatalf("risk absent: %+v", option)
		}
		if option.Cost.CandidateExecution && option.Executable {
			t.Fatal("third-party execution became executable")
		}
	}
}
