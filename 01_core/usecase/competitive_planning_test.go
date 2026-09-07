/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/competitive_planning_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func TestStep0041AutophagyDoesNotInventConsolidationOrHideConflicts(t *testing.T) {
	a, b := competitive41Candidate("one", "dependency graph"), competitive41Candidate("two", "dependency graph")
	f := aggregateInternalCapability("dependency graph", []domain.CandidateAssessment{a, a})
	if f.Disposition != domain.InternalKeep {
		t.Fatalf("duplicate created overlap: %+v", f)
	}
	f = aggregateInternalCapability("dependency graph", []domain.CandidateAssessment{a, b})
	if f.Disposition != domain.InternalReviewOverlap || len(f.NextChecks) == 0 {
		t.Fatalf("overlap not reviewed: %+v", f)
	}
	b.Capabilities[0].Status = domain.EvidenceConflict
	f = aggregateInternalCapability("dependency graph", []domain.CandidateAssessment{a, b})
	if f.Status != domain.EvidenceConflict || f.Disposition != domain.InternalUnknown {
		t.Fatalf("conflict hidden: %+v", f)
	}
	b.Capabilities = nil
	f = aggregateInternalCapability("dependency graph", []domain.CandidateAssessment{a, b})
	if f.Disposition != domain.InternalUnknown {
		t.Fatalf("missing assessment hidden: %+v", f)
	}
}
func TestStep0041NoChangeRequiresKnownInventory(t *testing.T) {
	a := competitive41Candidate("internal", "dependency graph")
	req := domain.CompetitiveAutophagyRequest{Idea: domain.IdeaSpec{Capabilities: []string{"dependency graph"}}, Target: domain.ResolvedRepository{Identity: domain.RepositoryIdentity{Owner: a.Repository.Owner, Repo: a.Repository.Repo}, Commit: a.Repository.Commit}}
	auto := buildAutophagyReport(req, domain.AssessmentReport{Assessments: []domain.CandidateAssessment{a}})
	p := competitive41Proposal(t, auto)
	if p.Status != "NO_CHANGE_REQUIRED" || p.Minimal || len(p.Components) > 0 {
		t.Fatalf("no-change wrong: %+v", p)
	}
	empty := competitive41Proposal(t, domain.AutophagyReport{})
	if empty.Status == "NO_CHANGE_REQUIRED" {
		t.Fatal("empty inventory invented success")
	}
}
func TestStep0041ReplacementRemainsAnOptionForExistingCapability(t *testing.T) {
	a := competitive41Candidate("internal", "dependency graph")
	req := domain.CompetitiveAutophagyRequest{Idea: domain.IdeaSpec{Capabilities: []string{"dependency graph"}}}
	auto := buildAutophagyReport(req, domain.AssessmentReport{Assessments: []domain.CandidateAssessment{a}})
	p := competitive41Proposal(t, auto, competitive41Candidate("alternative", "dependency graph"))
	found := false
	for _, o := range p.Options {
		found = found || o.Action == "COMPARE_REPLACEMENT"
	}
	if !found || len(p.Components) > 0 {
		t.Fatalf("replacement disappeared or became automatic: %+v", p)
	}
}
