/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/reuse_decision_test.md
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

func decisionPolicy() domain.ReuseDecisionPolicy {
	return domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, RequireLocal: true, MaxComponents: 3, ExtendMinCoverage: .6, IntegrationMinValue: .3}
}
func decisionSource(repo, rule string) domain.EvidenceSource {
	kind, path, snippet := domain.EvidenceSymbol, "main.go", "func Execute() { return } // "+rule
	if rule == "license" {
		kind, path, snippet = domain.EvidenceLicense, "LICENSE", "MIT License"
	}
	if rule == "maintenance" {
		kind, path, snippet = domain.EvidenceRepository, "@repository-metadata", `{"archived":false,"updated_at":"2026-09-01T00:00:00Z"}`
	}
	return domain.EvidenceSource{Kind: kind, Ref: domain.CommitRef{Owner: "o", Repo: repo, Commit: "abc"}, Path: path, RuleID: rule, Line: 1, Snippet: snippet, ContentHash: fmt.Sprintf("%x", sha256.Sum256([]byte(snippet))), Strength: 1}
}
func decisionAssessment(repo string, statuses ...domain.EvidenceStatus) domain.CandidateAssessment {
	a := domain.CandidateAssessment{Repository: domain.CommitRef{Owner: "o", Repo: repo, Commit: "abc"}, ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{decisionSource(repo, "alignment")}}, License: domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT", Sources: []domain.EvidenceSource{decisionSource(repo, "license")}}, Maintenance: domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{decisionSource(repo, "maintenance")}}, Operation: domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{decisionSource(repo, "local")}}}
	for i, status := range statuses {
		capability := []string{"core", "export", "sync"}[i]
		source := decisionSource(repo, "cap-"+capability)
		test := source
		test.Kind = domain.EvidenceTestImplementation
		test.Path = "main_test.go"
		test.Snippet = "func TestExecute() { Execute() } // " + capability
		test.ContentHash = fmt.Sprintf("%x", sha256.Sum256([]byte(test.Snippet)))
		a.Capabilities = append(a.Capabilities, domain.CapabilityEvidence{Capability: capability, Status: status, Sources: []domain.EvidenceSource{source, test}})
	}
	return a
}
func decisionRequest(as ...domain.CandidateAssessment) domain.ReuseDecisionRequest {
	return domain.ReuseDecisionRequest{Idea: domain.IdeaSpec{Problem: "x", Capabilities: []string{"core", "export", "sync"}}, Assessments: as, Policy: decisionPolicy()}
}

func TestDecideReuseForkAndIgnoresScreeningScore(t *testing.T) {
	a := decisionAssessment("full", domain.EvidenceProven, domain.EvidenceProven, domain.EvidenceProven)
	a.ScreeningDecision.RankingScore = .01
	r, err := DecideReuse(decisionRequest(a))
	if err != nil || r.Decision != domain.DecisionFork {
		t.Fatalf("got %s, %v", r.Decision, err)
	}
	if r.Primary == nil || len(r.Reasons[0].Evidence) == 0 {
		t.Fatal("fork must be traceable")
	}
}
func TestDecideReuseConflictPreventsForkAndKeepsReference(t *testing.T) {
	a := decisionAssessment("conflict", domain.EvidenceProven, domain.EvidenceProven, domain.EvidenceProven)
	a.Constraints = []domain.ConstraintEvidence{{Constraint: "offline", Evidence: domain.CapabilityEvidence{Status: domain.EvidenceConflict, Sources: []domain.EvidenceSource{decisionSource("conflict", "constraint")}}}}
	r, _ := DecideReuse(decisionRequest(a))
	if r.Decision != domain.DecisionReference || len(r.Blockers) == 0 {
		t.Fatalf("got %#v", r)
	}
}
func TestDecideReuseInconclusiveAndFailureNeverBuild(t *testing.T) {
	cases := []domain.CandidateAssessment{decisionAssessment("unknown", domain.EvidenceInconclusive, domain.EvidenceAbsent, domain.EvidenceAbsent), decisionAssessment("failed", domain.EvidenceAbsent, domain.EvidenceAbsent, domain.EvidenceAbsent)}
	cases[1].Failures = []domain.ExecutionFailure{{Scope: "repository", Stage: "read", Message: "failed"}}
	for _, a := range cases {
		r, _ := DecideReuse(decisionRequest(a))
		if r.Decision != domain.DecisionUnknown {
			t.Fatalf("%s became %s", a.Repository.Key(), r.Decision)
		}
	}
}

func TestDecideReuseEmptyCorpusIsInconclusive(t *testing.T) {
	r, err := DecideReuse(decisionRequest())
	if err != nil || r.Decision != domain.DecisionUnknown {
		t.Fatalf("got %s, %v", r.Decision, err)
	}
}
func TestDecideReuseExtend(t *testing.T) {
	a := decisionAssessment("extend", domain.EvidenceProven, domain.EvidenceProven, domain.EvidenceAbsent)
	a.ReusableParts = []domain.ReusableComponent{{Name: "engine", Kind: "package", Path: "engine", Capabilities: []string{"core", "export"}, Evidence: []domain.EvidenceSource{decisionSource("extend", "cap-core"), decisionSource("extend", "cap-export")}}}
	r, _ := DecideReuse(decisionRequest(a))
	if r.Decision != domain.DecisionExtend {
		t.Fatalf("got %s", r.Decision)
	}
}
func TestDecideReuseCompose(t *testing.T) {
	a := decisionAssessment("a", domain.EvidenceProven, domain.EvidenceAbsent, domain.EvidenceAbsent)
	a.ReusableParts = []domain.ReusableComponent{{Name: "core", Kind: "package", Path: "core", Capabilities: []string{"core"}, Evidence: []domain.EvidenceSource{decisionSource("a", "cap-core")}}}
	b := decisionAssessment("b", domain.EvidenceAbsent, domain.EvidenceProven, domain.EvidenceProven)
	b.ProblemAlignment.Status = domain.EvidencePartial
	b.ReusableParts = []domain.ReusableComponent{{Name: "io", Kind: "package", Path: "io", Capabilities: []string{"export", "sync"}, Evidence: []domain.EvidenceSource{decisionSource("b", "cap-export"), decisionSource("b", "cap-sync")}}}
	r, _ := DecideReuse(decisionRequest(b, a))
	if r.Decision != domain.DecisionUnknown || r.Composition == nil || len(r.Composition.Candidates) != 1 || r.Composition.Candidates[0].Status != domain.CompositionProposed {
		t.Fatalf("got %#v", r)
	}
}
func TestDecideReuseIntegrate(t *testing.T) {
	a := decisionAssessment("lib", domain.EvidenceProven, domain.EvidenceAbsent, domain.EvidenceAbsent)
	a.ProblemAlignment.Status = domain.EvidencePartial
	a.ReusableParts = []domain.ReusableComponent{{Name: "parser", Kind: "library", Path: "parser", Capabilities: []string{"core"}, Evidence: []domain.EvidenceSource{decisionSource("lib", "cap-core")}}}
	r, _ := DecideReuse(decisionRequest(a))
	if r.Decision != domain.DecisionIntegrate {
		t.Fatalf("got %s", r.Decision)
	}
}
func TestDecideReuseBuildOnlyOnConclusiveAbsence(t *testing.T) {
	a := decisionAssessment("none", domain.EvidenceAbsent, domain.EvidenceAbsent, domain.EvidenceAbsent)
	a.ProblemAlignment.Status = domain.EvidenceAbsent
	r, _ := DecideReuse(decisionRequest(a))
	if r.Decision != domain.DecisionBuild {
		t.Fatalf("got %s", r.Decision)
	}
}
func TestDecideReuseOrderPolicyAndInputImmutability(t *testing.T) {
	a := decisionAssessment("z", domain.EvidenceAbsent, domain.EvidenceAbsent, domain.EvidenceAbsent)
	b := decisionAssessment("a", domain.EvidenceProven, domain.EvidenceProven, domain.EvidenceProven)
	req := decisionRequest(a, b)
	before, _ := json.Marshal(req)
	r1, _ := DecideReuse(req)
	req2 := decisionRequest(b, a)
	r2, _ := DecideReuse(req2)
	after, _ := json.Marshal(req)
	if !reflect.DeepEqual(r1, r2) {
		t.Fatal("candidate order changed report")
	}
	if string(before) != string(after) {
		t.Fatal("input mutated")
	}
	req.Policy.AllowedLicenses = []string{"Apache-2.0"}
	r3, _ := DecideReuse(req)
	if r3.Decision != domain.DecisionReference {
		t.Fatalf("policy did not alter expected result: %s", r3.Decision)
	}
}
