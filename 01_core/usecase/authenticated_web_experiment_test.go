/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/authenticated_web_experiment_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import "testing"

func TestAuthenticatedWebExperimentStates(t *testing.T) {
	for _, state := range []string{"SUCCESS", "EMPTY", "UNKNOWN_RATE_LIMIT", "UNKNOWN_TIMEOUT", "UNKNOWN_PROVIDER_RESPONSE", "REJECTED_PRIVATE", "SKIPPED_BUDGET"} {
		if !validWebStates[state] {
			t.Fatalf("estado obrigatório ausente: %s", state)
		}
	}
	if validWebStates["NOT_A_STATE"] {
		t.Fatal("estado arbitrário aceito")
	}
}
