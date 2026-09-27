/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/recommendation_intake_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

import (
	"encoding/json"
	"testing"
)

func TestRecommendationIntakeDraftPreservesOmittedAndExplicitFalse(t *testing.T) {
	var omitted RecommendationIntakeDraft
	encoded, err := json.Marshal(omitted)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != `{}` {
		t.Fatalf("rascunho omitido = %s", encoded)
	}

	value := false
	explicit := RecommendationIntakeDraft{Execution: &ExecutionIntakeInput{AllowCandidateExecution: &value}}
	encoded, err = json.Marshal(explicit)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != `{"execution":{"allow_candidate_execution":false}}` {
		t.Fatalf("false explícito = %s", encoded)
	}
}

func TestRecommendationIntakeAuthorityVocabularyIsStable(t *testing.T) {
	if AuthorityUser != "USER" || AuthorityAIProposal != "AI_PROPOSAL" || ResolverAskUser != "ASK_USER" || ResolverTerminal != "TERMINAL" {
		t.Fatal("vocabulário de autoridade ou resolução mudou")
	}
}
