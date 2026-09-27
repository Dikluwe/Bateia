/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_value_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDecisionValueContractSeparatesHypothesisFromEvidence(t *testing.T) {
	plan := DecisionValuePlan{
		SchemaVersion: "1", Method: "COUNTERFACTUAL_DECISION_REACHABILITY_V1", BaselineDecision: DecisionUnknown,
		Options: []InvestigationValueOption{{
			ActionID: "a", Impact: DecisionPivotal, EpistemicBasis: "EXACT_DECISION_REPLAY_WITH_DECLARED_OUTCOMES",
			Outcomes:          []CounterfactualOutcome{{OutcomeCode: "PROVEN", Status: EvidenceProven, Decision: DecisionFork}},
			PossibleDecisions: []ReuseDecision{DecisionFork}, Cost: InvestigationCost{Tests: 1},
		}},
	}
	encoded, err := json.Marshal(plan)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	for _, expected := range []string{`"impact":"DECISION_PIVOTAL"`, `"epistemic_basis":"EXACT_DECISION_REPLAY_WITH_DECLARED_OUTCOMES"`, `"candidate_execution":false`} {
		if !strings.Contains(text, expected) {
			t.Fatalf("campo ausente %s: %s", expected, text)
		}
	}
	for _, forbidden := range []string{`"probability"`, `"evsi"`, `"sources"`} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("hipótese ganhou autoridade indevida %s: %s", forbidden, text)
		}
	}
}
