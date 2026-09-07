/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/authority_absence_step0043_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"testing"
)

func TestStep0043SourcelessAbsenceCannotAuthorizeBuild(t *testing.T) {
	ref := step0043Bundle("").Repository
	source := func(kind domain.EvidenceKind, path, rule, snippet string) domain.EvidenceSource {
		return domain.EvidenceSource{Kind: kind, Ref: ref, Path: path, Line: 1, RuleID: rule, Snippet: snippet, ContentHash: step0043Hash(snippet), Strength: 1}
	}
	absent := source(domain.EvidenceSymbol, "queue.go", "queue-unsupported", `func StartQueue() error { return errors.New("unsupported") }`)
	a := domain.CandidateAssessment{Repository: ref, ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceAbsent, Sources: []domain.EvidenceSource{absent}}, Capabilities: []domain.CapabilityEvidence{{Capability: "bounded queue", Status: domain.EvidenceAbsent}}, License: domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT", Sources: []domain.EvidenceSource{source(domain.EvidenceLicense, "LICENSE", "license", "MIT License")}}, Maintenance: domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{source(domain.EvidenceRepository, "@repository-metadata", "maintenance", `{"archived":false}`)}}, Operation: domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{source(domain.EvidenceSymbol, "main.go", "local", "func RunLocally() {}")}}}
	got := usecase.SanitizeAssessmentEvidence(a)
	if got.Capabilities[0].Status != domain.EvidenceInconclusive {
		t.Fatalf("sourceless absence remained conclusive: %s", got.Capabilities[0].Status)
	}
	if a.Capabilities[0].Status != domain.EvidenceAbsent {
		t.Fatal("sanitizer mutated caller")
	}
	decision, err := usecase.DecideReuse(domain.ReuseDecisionRequest{Idea: domain.IdeaSpec{Problem: "bounded queue", Capabilities: []string{"bounded queue"}}, Assessments: []domain.CandidateAssessment{a}, Policy: usecase.DefaultRecommendationPolicyV1()})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Decision == domain.DecisionBuild {
		t.Fatal("sourceless absence authorized BUILD")
	}
	a.Capabilities[0].Sources = []domain.EvidenceSource{absent}
	positive := usecase.SanitizeAssessmentEvidence(a)
	if positive.Capabilities[0].Status != domain.EvidenceAbsent {
		t.Fatal("bounded witnessed absence lost admissibility")
	}
}
