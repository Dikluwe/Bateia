/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recovery_original_targets_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func TestOriginalRecoveryTargetsKeepRulesAndOrder(t *testing.T) {
	bundles := []domain.EvidenceBundle{{Repository: domain.CommitRef{Owner: "org", Repo: "tool", Commit: "frozen"}, Plan: domain.EvidenceExtractionPlan{Targets: []domain.ExtractionTarget{{Kind: domain.EvidenceLicense, Path: "LICENSE", Capabilities: []string{"first", "second"}}, {Kind: domain.EvidenceManifest, Path: "go.mod"}}}, Targets: []domain.CollectionTargetResult{{Target: domain.ExtractionTarget{Kind: domain.EvidenceFile, Path: "main.go"}}}}}
	out := OriginalRecoveryTargets(bundles)
	if len(out) != 3 || out[0].RuleID != "GATE-LICENSE-001" || out[0].Capability != "first" || out[1].RuleID != "GATE-OPERATION-001" || out[2].RuleID != "original" || out[2].Path != "main.go" {
		t.Fatalf("history projection changed: %+v", out)
	}
	if bundles[0].Plan.Targets[0].Capabilities[0] != "first" {
		t.Fatal("input mutated")
	}
}
