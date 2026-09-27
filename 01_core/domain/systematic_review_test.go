/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/systematic_review_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSystematicReviewVocabularyKeepsEvidenceCertaintyAndDecisionSeparate(t *testing.T) {
	claim := ClaimCertainty{CandidateID: "acme/tool", ClaimID: "capability:index", ObservedStatus: EvidenceProven, Level: CertaintyFragile, Concerns: []ClaimConcern{{Dimension: "completeness", Concern: ConcernSerious, Reason: "partial collection"}}, SourceCount: 1}
	encoded, err := json.Marshal(claim)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	for _, expected := range []string{`"observed_status":"PROVEN"`, `"certainty":"FRAGILE"`, `"dimension":"completeness"`} {
		if !strings.Contains(text, expected) {
			t.Fatalf("serialização não preservou %s: %s", expected, text)
		}
	}
}

func TestActiveLearningContractCannotAuthorizeExclusion(t *testing.T) {
	plan := ActiveLearningPlan{SchemaVersion: "1", Authority: "PRIORITIZATION_ONLY", CanExclude: false, Priorities: []ActiveLearningPriority{}}
	encoded, err := json.Marshal(plan)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), `"can_exclude":true`) {
		t.Fatalf("plano ganhou autoridade de exclusão: %s", encoded)
	}
}
