/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_authority_test.md
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

func TestStep0043AuthorityPurityIdempotenceAndPositive(t *testing.T) {
	ref := step0043Bundle("").Repository
	snippet := "func BuildGraph() {}"
	source := domain.EvidenceSource{Kind: domain.EvidenceSymbol, Ref: ref, Path: "graph.go", Line: 1, Symbol: "BuildGraph", Snippet: snippet, RuleID: "dependency-graph", ContentHash: step0043Hash(snippet), Strength: 1}
	if !usecase.EvidenceSourcesAdmissible(ref, []domain.EvidenceSource{source}, true) {
		t.Fatal("valid symbol witness rejected")
	}
	document := source
	document.Kind = domain.EvidenceDocumentation
	if usecase.EvidenceSourcesAdmissible(ref, []domain.EvidenceSource{document}, true) {
		t.Fatal("documentation authorized implementation")
	}
	foreign := source
	foreign.Ref.Repo = "other"
	if usecase.EvidenceSourcesAdmissible(ref, []domain.EvidenceSource{foreign}, true) {
		t.Fatal("foreign witness authorized implementation")
	}
	input := domain.CandidateAssessment{Repository: ref, Capabilities: []domain.CapabilityEvidence{{Capability: "dependency graph", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{source}}, {Capability: "unsupported", Status: domain.EvidenceProven}}}
	before, _ := json.Marshal(input)
	got := usecase.SanitizeAssessmentEvidence(input)
	after, _ := json.Marshal(input)
	if string(before) != string(after) {
		t.Fatal("sanitization mutated input")
	}
	if !reflect.DeepEqual(got, usecase.SanitizeAssessmentEvidence(got)) {
		t.Fatal("sanitization not idempotent")
	}
	for _, capability := range got.Capabilities {
		if capability.Capability == "unsupported" && capability.Status == domain.EvidenceProven {
			t.Fatal("sourceless status retained proof")
		}
		if capability.Capability == "dependency graph" && capability.Status != domain.EvidenceProven {
			t.Fatal("admissible positive lost authority")
		}
	}
}
