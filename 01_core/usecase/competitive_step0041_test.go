/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/competitive_step0041_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"strings"
	"testing"
)

// These fixtures supply assessments to test selection, not behavioral verifiers.
func competitive41Candidate(repo string, caps ...string) domain.CandidateAssessment {
	ref := domain.CommitRef{Owner: "external", Repo: repo, Commit: strings.Repeat("a", 40)}
	source := domain.EvidenceSource{Kind: domain.EvidenceSymbol, Ref: ref, Path: "src/implementation.go", Line: 1, Snippet: "declared behavior fixture", RuleID: "supplied-assessment", ContentHash: contentHash("declared behavior fixture")}
	a := domain.CandidateAssessment{Repository: ref, ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{source}}, License: domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT"}, Maintenance: domain.MaintenanceAssessment{Status: domain.EvidenceProven}}
	license := source
	license.Kind = domain.EvidenceLicense
	license.Path = "LICENSE"
	license.Snippet = "MIT License"
	license.ContentHash = contentHash(license.Snippet)
	a.License.Sources = []domain.EvidenceSource{license}
	metadata := source
	metadata.Kind = domain.EvidenceRepository
	metadata.Path = "@repository-metadata"
	a.Maintenance.Sources = []domain.EvidenceSource{metadata}
	for _, c := range caps {
		a.Capabilities = append(a.Capabilities, domain.CapabilityEvidence{Capability: c, Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{source}})
	}
	return a
}
func competitive41Auto(t *testing.T) domain.AutophagyReport {
	t.Helper()
	b := autophagyBundle("local", "target", strings.Repeat("b", 40), "", true, "")
	auto, err := AnalyzeAutophagy(domain.IdeaSpec{Capabilities: []string{"dependency graph", "blast radius"}}, domain.ResolvedRepository{Identity: domain.RepositoryIdentity{Owner: "local", Repo: "target"}, Commit: b.Repository.Commit}, []domain.EvidenceBundle{b}, DefaultEvidenceRules())
	if err != nil {
		t.Fatal(err)
	}
	if len(auto.Gaps) != 2 {
		t.Fatalf("fixture gaps: %+v", auto)
	}
	return auto
}
func competitive41Policy() domain.KatamariPolicy {
	return domain.KatamariPolicy{SchemaVersion: "1", MaxProjects: 3, MaxComponents: 3, AllowedLicenses: []string{"MIT"}, RequireImmutableRef: true}
}
func competitive41Proposal(t *testing.T, auto domain.AutophagyReport, as ...domain.CandidateAssessment) domain.KatamariProposal {
	t.Helper()
	p, e := ProposeKatamari(auto, domain.AssessmentReport{Assessments: as}, competitive41Policy(), domain.ReuseDecisionPolicy{AllowedLicenses: []string{"MIT"}})
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func TestStep0041CoverageIsNotIntegration(t *testing.T) {
	auto := competitive41Auto(t)
	partial := competitive41Proposal(t, auto, competitive41Candidate("graph", "dependency graph"))
	if partial.Status != "PARTIAL_COVERAGE" || partial.Feasible || partial.CoverageComplete || len(partial.UncoveredGaps) != 1 {
		t.Fatalf("partial promoted: %+v", partial)
	}
	complete := competitive41Proposal(t, auto, competitive41Candidate("graph", "dependency graph"), competitive41Candidate("impact", "blast radius"))
	if complete.Status != "COVERED_AWAITING_VALIDATION" || !complete.CoverageComplete || complete.Feasible || !complete.Minimal {
		t.Fatalf("invalid coverage: %+v", complete)
	}
	for _, c := range complete.Components {
		if c.Decision != domain.DecisionUnknown || c.CostStatus != "NOT_ESTIMATED" || c.Cost != 0 || len(c.NextChecks) == 0 {
			t.Fatalf("invented integration: %+v", c)
		}
	}
}
func TestStep0041GatesCannotBeCompensatedByCoverage(t *testing.T) {
	auto := competitive41Auto(t)
	cases := map[string]func(*domain.CandidateAssessment){
		"constraint": func(a *domain.CandidateAssessment) {
			a.Constraints = []domain.ConstraintEvidence{{Constraint: "privacy", Evidence: domain.CapabilityEvidence{Status: domain.EvidenceConflict}}}
		},
		"maintenance":           func(a *domain.CandidateAssessment) { a.Maintenance.Status = domain.EvidenceInconclusive },
		"license-source":        func(a *domain.CandidateAssessment) { a.License.Sources[0].Path = "license_checker_test.go" },
		"capability-no-sources": func(a *domain.CandidateAssessment) { a.Capabilities[0].Sources = nil },
		"foreign-source":        func(a *domain.CandidateAssessment) { a.Capabilities[0].Sources[0].Ref.Repo = "another" },
		"bad-hash":              func(a *domain.CandidateAssessment) { a.Capabilities[0].Sources[0].ContentHash = "fake" },
		"archived":              func(a *domain.CandidateAssessment) { a.Maintenance.Archived = true },
		"failed-collection": func(a *domain.CandidateAssessment) {
			a.Failures = []domain.ExecutionFailure{{Stage: "content", Message: "timeout"}}
		},
		"forbidden-license": func(a *domain.CandidateAssessment) { a.License.Identifier = "GPL-3.0" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			a := competitive41Candidate("both", "dependency graph", "blast radius")
			mutate(&a)
			p := competitive41Proposal(t, auto, a)
			if len(p.Components) != 0 || len(p.Rejected) == 0 || p.Feasible {
				t.Fatalf("gate bypass: %+v", p)
			}
		})
	}
}
func TestStep0041PartialEvidenceCreatesAnInvestigation(t *testing.T) {
	a := competitive41Candidate("graph", "dependency graph")
	a.Capabilities[0].Status = domain.EvidencePartial
	a.Maintenance.Status = domain.EvidenceInconclusive
	p, e := ProposeKatamari(competitive41Auto(t), domain.AssessmentReport{Assessments: []domain.CandidateAssessment{a}}, domain.KatamariPolicy{SchemaVersion: "1", MaxProjects: 2, MaxComponents: 2, AllowedLicenses: []string{"MIT"}, AllowUnknown: true}, domain.ReuseDecisionPolicy{})
	if e != nil {
		t.Fatal(e)
	}
	found := false
	for _, o := range p.Options {
		found = found || o.Action == "INVESTIGATE_EXTERNAL" && len(o.NextChecks) > 0 && len(o.Sources) > 0
	}
	if !found || len(p.Components) > 0 || len(p.CoveredGaps) > 0 {
		t.Fatalf("partial lost or promoted: %+v", p)
	}
}
func TestStep0041LegacyGapListCannotAuthorizeCoverage(t *testing.T) {
	p := competitive41Proposal(t, domain.AutophagyReport{SchemaVersion: "1", Gaps: []string{"dependency graph"}}, competitive41Candidate("graph", "dependency graph"))
	if len(p.Components) > 0 || len(p.Warnings) == 0 || p.Status != "INVESTIGATE" {
		t.Fatalf("legacy gap trusted: %+v", p)
	}
}
func TestStep0041MinimumSubsetAndDeterministicDuplicates(t *testing.T) {
	auto := competitive41Auto(t)
	a, b, c := competitive41Candidate("graph", "dependency graph"), competitive41Candidate("impact", "blast radius"), competitive41Candidate("both", "dependency graph", "blast radius")
	left := competitive41Proposal(t, auto, a, b, c)
	right := competitive41Proposal(t, auto, c, b, a, a)
	x, _ := json.Marshal(left)
	y, _ := json.Marshal(right)
	if string(x) != string(y) {
		t.Fatalf("permutation or duplicate changed output:\n%s\n%s", x, y)
	}
	if len(left.Components) != 1 || left.Components[0].Repository.Repo != "both" {
		t.Fatalf("not minimum: %+v", left)
	}
}
func TestStep0041ConflictingAssessmentsMustBeReconciled(t *testing.T) {
	a, b := competitive41Candidate("graph", "dependency graph"), competitive41Candidate("graph", "dependency graph")
	b.Capabilities[0].Status = domain.EvidenceConflict
	p := competitive41Proposal(t, competitive41Auto(t), a, b)
	if len(p.Components) > 0 || len(p.Rejected) != 1 || !strings.Contains(p.Rejected[0].Message, "divergentes") {
		t.Fatalf("conflicting revision selected: %+v", p)
	}
}
func TestStep0041CannotComposeTwoRevisionsOfOneProject(t *testing.T) {
	a, b := competitive41Candidate("same", "dependency graph"), competitive41Candidate("same", "blast radius")
	b.Repository.Commit = strings.Repeat("c", 40)
	b.ProblemAlignment.Sources[0].Ref = b.Repository
	b.License.Sources[0].Ref = b.Repository
	b.Maintenance.Sources[0].Ref = b.Repository
	b.Capabilities[0].Sources[0].Ref = b.Repository
	p := competitive41Proposal(t, competitive41Auto(t), a, b)
	if len(p.Components) != 1 || p.CoverageComplete {
		t.Fatalf("incompatible revisions composed: %+v", p)
	}
}

func TestStep0041EnumerationLimitAndOperationPolicy(t *testing.T) {
	var as []domain.CandidateAssessment
	for i := 0; i < 21; i++ {
		as = append(as, competitive41Candidate(string(rune('a'+i)), "dependency graph"))
	}
	if _, err := ProposeKatamari(competitive41Auto(t), domain.AssessmentReport{Assessments: as}, competitive41Policy(), domain.ReuseDecisionPolicy{}); err == nil {
		t.Fatal("unbounded enumeration")
	}
	p := competitive41Policy()
	p.RequiredOperation = "remote-only"
	if _, err := ProposeKatamari(competitive41Auto(t), domain.AssessmentReport{}, p, domain.ReuseDecisionPolicy{}); err == nil {
		t.Fatal("unsupported policy silently ignored")
	}
}

func TestStep0041OperationAndBothLicensePoliciesAreEnforced(t *testing.T) {
	auto := competitive41Auto(t)
	a := competitive41Candidate("graph", "dependency graph")
	a.Operation.LocalStatus = domain.EvidenceAbsent
	p, err := ProposeKatamari(auto, domain.AssessmentReport{Assessments: []domain.CandidateAssessment{a}}, competitive41Policy(), domain.ReuseDecisionPolicy{RequireLocal: true})
	if err != nil || len(p.Components) > 0 || len(p.Rejected) == 0 {
		t.Fatalf("local gate bypass: %+v %v", p, err)
	}
	p, err = ProposeKatamari(auto, domain.AssessmentReport{Assessments: []domain.CandidateAssessment{a}}, competitive41Policy(), domain.ReuseDecisionPolicy{AllowedLicenses: []string{"Apache-2.0"}})
	if err != nil || len(p.Components) > 0 || len(p.Rejected) == 0 {
		t.Fatalf("reuse policy bypass: %+v %v", p, err)
	}
}

func TestStep0041AutophagyScopeHashAndInputAreStable(t *testing.T) {
	a := autophagyBundle("local", "one", "a1", "DependencyGraph", true, "")
	b := autophagyBundle("local", "two", "b1", "BlastRadius", true, "")
	target := domain.ResolvedRepository{Identity: domain.RepositoryIdentity{Owner: "local", Repo: "one"}, Commit: "a1"}
	input := []domain.EvidenceBundle{a, b}
	before, _ := json.Marshal(input)
	idea := domain.IdeaSpec{Capabilities: []string{"dependency graph", "blast radius"}}
	left, err := AnalyzeAutophagy(idea, target, input, DefaultEvidenceRules())
	if err != nil {
		t.Fatal(err)
	}
	right, err := AnalyzeAutophagy(idea, target, []domain.EvidenceBundle{b, a, a}, DefaultEvidenceRules())
	if err != nil {
		t.Fatal(err)
	}
	x, _ := json.Marshal(left)
	y, _ := json.Marshal(right)
	after, _ := json.Marshal(input)
	if string(x) != string(y) || string(before) != string(after) {
		t.Fatal("order, duplication or mutation changed inventory")
	}
	if len(left.Scope) != 2 || len(left.Options) == 0 {
		t.Fatal("missing scope or useful next steps")
	}
}
