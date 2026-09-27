/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/composition_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func compositionTestSource(repo string) domain.EvidenceSource {
	s := decisionSource(repo, "integration-test")
	s.Kind = domain.EvidenceTestImplementation
	s.Path = "integration_test.go"
	s.Snippet = "func TestComposition() { verifyIntegration() }"
	s.ContentHash = fmt.Sprintf("%x", sha256.Sum256([]byte(s.Snippet)))
	return s
}

func compositionFixture() (domain.ReuseDecisionRequest, string, string) {
	a := decisionAssessment("a", domain.EvidenceProven, domain.EvidenceAbsent, domain.EvidenceAbsent)
	a.ProblemAlignment.Status = domain.EvidencePartial
	a.ReusableParts = []domain.ReusableComponent{{Name: "core", Kind: "package", Path: "core", Capabilities: []string{"core"}, Evidence: []domain.EvidenceSource{decisionSource("a", "cap-core")}}}
	b := decisionAssessment("b", domain.EvidenceAbsent, domain.EvidenceProven, domain.EvidenceProven)
	b.ProblemAlignment.Status = domain.EvidencePartial
	b.ReusableParts = []domain.ReusableComponent{{Name: "io", Kind: "package", Path: "io", Capabilities: []string{"export", "sync"}, Evidence: []domain.EvidenceSource{decisionSource("b", "cap-export"), decisionSource("b", "cap-sync")}}}
	return decisionRequest(a, b), compositionMemberID(a.Repository, a.ReusableParts[0]), compositionMemberID(b.Repository, b.ReusableParts[0])
}

func verifiedComposition(a, b string, repo ...string) *domain.CompositionInput {
	sourceRepo := "a"
	if len(repo) > 0 {
		sourceRepo = repo[0]
	}
	src := compositionTestSource(sourceRepo)
	return &domain.CompositionInput{
		Connections:   []domain.CompositionConnection{{From: a, To: b, Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{src}, Cost: &domain.IntegrationCost{Unit: "hours", Low: 2, High: 4}}},
		Verifications: []domain.CompositionVerification{{Members: []string{a, b}, Status: domain.EvidenceProven, Executed: true, Passed: true, Sources: []domain.EvidenceSource{src}}},
	}
}

func TestCompositionVerifiedNeedsConnectedInterfacesAndPassingRun(t *testing.T) {
	req, a, b := compositionFixture()
	req.Composition = verifiedComposition(a, b)
	before, _ := json.Marshal(req)
	report, err := DecideReuse(req)
	if err != nil || report.Decision != domain.DecisionCompose {
		t.Fatalf("decision=%s err=%v", report.Decision, err)
	}
	if report.Composition == nil || len(report.Composition.Candidates) != 1 {
		t.Fatal("missing first-class composition candidate")
	}
	candidate := report.Composition.Candidates[0]
	if candidate.Status != domain.CompositionVerified || candidate.IntegrationCost == nil || candidate.IntegrationCost.Low != 2 || candidate.IntegrationCost.High != 4 {
		t.Fatalf("candidate=%#v", candidate)
	}
	after, _ := json.Marshal(req)
	if string(before) != string(after) {
		t.Fatal("input mutated")
	}
	req.Assessments[0], req.Assessments[1] = req.Assessments[1], req.Assessments[0]
	reordered, err := DecideReuse(req)
	if err != nil || !reflect.DeepEqual(report, reordered) {
		t.Fatal("composition changed with assessment order")
	}
}

func TestCompositionConflictAndMissingDependencyBlock(t *testing.T) {
	req, a, b := compositionFixture()
	req.Composition = verifiedComposition(a, b)
	req.Composition.Exclusions = []domain.CompositionExclusion{{Left: a, Right: b, Reason: "incompatible licenses"}}
	report, err := DecideReuse(req)
	if err != nil || report.Decision == domain.DecisionCompose || report.Composition.Candidates[0].Status != domain.CompositionIncompatible {
		t.Fatalf("exclusion: %#v, %v", report.Composition, err)
	}
	req.Composition.Exclusions = nil
	req.Composition.Requirements = []domain.CompositionRequirement{{Component: a, Requires: b}}
	report, err = DecideReuse(req)
	if err != nil || report.Decision != domain.DecisionCompose {
		t.Fatalf("satisfied dependency: %s, %v", report.Decision, err)
	}
	req.Composition.Verifications[0].Passed = false
	report, err = DecideReuse(req)
	if err != nil || report.Composition.Candidates[0].Status != domain.CompositionIncompatible {
		t.Fatalf("failed run: %#v, %v", report.Composition, err)
	}
}

func TestCompositionRejectsWrongEvidenceAndUnknownCost(t *testing.T) {
	req, a, b := compositionFixture()
	req.Composition = verifiedComposition(a, b)
	req.Composition.Connections[0].Sources = []domain.EvidenceSource{compositionTestSource("outsider")}
	report, err := DecideReuse(req)
	if err != nil || report.Decision != domain.DecisionUnknown || report.Composition.Candidates[0].Status != domain.CompositionProposed {
		t.Fatalf("wrong source: %#v, %v", report.Composition, err)
	}
	req.Composition = verifiedComposition(a, b)
	req.Composition.Connections[0].Cost = nil
	report, err = DecideReuse(req)
	if err != nil || report.Decision != domain.DecisionCompose || report.Composition.Candidates[0].IntegrationCost != nil {
		t.Fatalf("unknown cost must not become zero: %#v, %v", report.Composition, err)
	}
}

func TestCompositionRejectsUnknownMemberReference(t *testing.T) {
	req, a, _ := compositionFixture()
	req.Composition = &domain.CompositionInput{Requirements: []domain.CompositionRequirement{{Component: a, Requires: "missing"}}}
	if _, err := DecideReuse(req); err == nil {
		t.Fatal("unknown member must fail validation")
	}
}

func TestCompositionOverlapAndBoundedSearch(t *testing.T) {
	req, aID, bID := compositionFixture()
	a := &req.Assessments[0]
	a.Capabilities[1].Status = domain.EvidenceProven
	a.ReusableParts[0].Capabilities = append(a.ReusableParts[0].Capabilities, "export")
	a.ReusableParts[0].Evidence = append(a.ReusableParts[0].Evidence, decisionSource("a", "cap-export"))
	req.Composition = verifiedComposition(aID, bID)
	report, err := DecideReuse(req)
	if err != nil || report.Decision != domain.DecisionCompose || !contains(report.Composition.Candidates[0].Overlap, "export") {
		t.Fatalf("overlap not represented: %+v, %v", report.Composition, err)
	}
	many := decisionAssessment("many", domain.EvidenceProven, domain.EvidenceAbsent, domain.EvidenceAbsent)
	many.ProblemAlignment.Status = domain.EvidencePartial
	for i := 0; i < compositionPartLimit+1; i++ {
		name := fmt.Sprintf("part-%02d", i)
		many.ReusableParts = append(many.ReusableParts, domain.ReusableComponent{Name: name, Kind: "package", Path: name, Capabilities: []string{"core"}, Evidence: []domain.EvidenceSource{decisionSource("many", "cap-core")}})
	}
	limited, err := analyzeCompositions([]decisionCandidate{inspectCandidate(many, []string{"core", "export", "sync"}, decisionPolicy())}, []string{"core", "export", "sync"}, decisionPolicy(), nil)
	if err != nil || !limited.Truncated || limited.StopReason != "PART_LIMIT" {
		t.Fatalf("search bound not exposed: %+v, %v", limited, err)
	}
}

func TestCompositionDoesNotBorrowRepositoryCapabilityForWrongPart(t *testing.T) {
	req, _, _ := compositionFixture()
	req.Assessments[0].ReusableParts[0].Evidence = []domain.EvidenceSource{decisionSource("a", "unrelated")}
	report, err := DecideReuse(req)
	if err != nil {
		t.Fatal(err)
	}
	if report.Composition == nil || len(report.Composition.Candidates) != 0 || report.Decision == domain.DecisionCompose {
		t.Fatalf("capability borrowed by unrelated part: %+v", report.Composition)
	}
}
