/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/software_alternative_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type alternativeTestInput struct {
	target      domain.SoftwareAlternativeTarget
	idea        domain.IdeaSpec
	discovery   domain.DiscoverySnapshot
	assessments []domain.CandidateAssessment
	policy      domain.ReuseDecisionPolicy
}

func newAlternativeTestInput(t *testing.T) alternativeTestInput {
	t.Helper()
	idea := domain.IdeaSpec{Problem: "analyze dependency impact", Capabilities: []string{"dependency graph", "impact analysis"}}
	return alternativeTestInput{
		target:    domain.SoftwareAlternativeTarget{Name: "Known Tool"},
		idea:      idea,
		discovery: domain.DiscoverySnapshot{SchemaVersion: "2", Idea: idea, Strategy: "literal"},
		policy:    alternativeTestPolicy(t, false, false, "MIT", "Apache-2.0"),
	}
}

func alternativeTestPolicy(t *testing.T, requireLocal, allowArchived bool, licenses ...string) domain.ReuseDecisionPolicy {
	t.Helper()
	payload, err := json.Marshal(map[string]any{
		"schema_version":        "1",
		"allowed_licenses":      licenses,
		"max_components":        2,
		"extend_min_coverage":   0.8,
		"integration_min_value": 0.5,
		"require_local":         requireLocal,
		"allow_archived":        allowArchived,
	})
	if err != nil {
		t.Fatal(err)
	}
	var policy domain.ReuseDecisionPolicy
	if err := json.Unmarshal(payload, &policy); err != nil {
		t.Fatal(err)
	}
	return policy
}

func alternativeTestCandidate(owner, repo, commit string, stars int) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{
		Owner: owner, Repo: repo, CommitID: commit, Stars: stars,
		URL: "https://github.com/" + owner + "/" + repo, Description: "dependency graph and impact analysis",
	}
}

func alternativeTestSource(ref domain.CommitRef, kind domain.EvidenceKind, path, rule, snippet string) domain.EvidenceSource {
	source := decisionSource(ref.Repo, rule)
	source.Ref = ref
	source.Kind = kind
	source.Path = path
	source.Snippet = snippet
	source.ContentHash = fmt.Sprintf("%x", sha256.Sum256([]byte(snippet)))
	return source
}
func alternativeTestAssessment(owner, repo, commit string, statuses ...domain.EvidenceStatus) domain.CandidateAssessment {
	ref := domain.CommitRef{Owner: owner, Repo: repo, Commit: commit}
	if len(statuses) == 0 {
		statuses = []domain.EvidenceStatus{domain.EvidenceProven, domain.EvidenceProven}
	}
	capabilities := []domain.CapabilityEvidence{}
	names := []string{"dependency graph", "impact analysis"}
	for i, status := range statuses {
		name := names[i%len(names)]
		capabilities = append(capabilities, domain.CapabilityEvidence{Capability: name, Status: status, Sources: []domain.EvidenceSource{alternativeTestSource(ref, domain.EvidenceSymbol, name+".go", "CAP-001", "func Analyze() {} // "+name)}})
	}
	return domain.CandidateAssessment{Repository: ref, ProblemAlignment: domain.CapabilityEvidence{Capability: "analyze dependency impact", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{alternativeTestSource(ref, domain.EvidenceSymbol, "analysis.go", "alignment", "func AnalyzeDependencyImpact() {}")}}, Capabilities: capabilities, License: domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT", Sources: []domain.EvidenceSource{alternativeTestSource(ref, domain.EvidenceLicense, "LICENSE", "license", "MIT License")}}, Maintenance: domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{alternativeTestSource(ref, domain.EvidenceRepository, "@repository-metadata", "maintenance", `{"archived":false}`)}}, Operation: domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{alternativeTestSource(ref, domain.EvidenceSymbol, "main.go", "local", "func RunLocally() {}")}}}
}

func (input *alternativeTestInput) add(candidate domain.RepositoryCandidate, assessment domain.CandidateAssessment) {
	input.discovery.Candidates = append(input.discovery.Candidates, domain.DiscoveredCandidate{Repository: candidate})
	input.assessments = append(input.assessments, assessment)
}

func rankAlternativeTestInput(t *testing.T, input alternativeTestInput) domain.SoftwareAlternativeReport {
	t.Helper()
	report, err := RankSoftwareAlternatives(input.target, input.idea, input.discovery, input.assessments, input.policy)
	if err != nil {
		t.Fatal(err)
	}
	return report
}

func alternativeIdentity(ref domain.CommitRef) string {
	return strings.ToLower(ref.Owner + "/" + ref.Repo)
}

func requireAlternativeFit(t *testing.T, report domain.SoftwareAlternativeReport, repository string, fit domain.SoftwareAlternativeFit, excluded bool) domain.SoftwareAlternative {
	t.Helper()
	values := report.Alternatives
	if excluded {
		values = report.Excluded
	}
	for _, alternative := range values {
		if alternativeIdentity(alternative.Repository) == strings.ToLower(repository) {
			if alternative.Fit != fit {
				t.Fatalf("%s fit = %s, want %s", repository, alternative.Fit, fit)
			}
			return alternative
		}
	}
	t.Fatalf("%s absent from excluded=%v: alternatives=%+v excluded=%+v", repository, excluded, report.Alternatives, report.Excluded)
	return domain.SoftwareAlternative{}
}

func requireTopChoice(t *testing.T, report domain.SoftwareAlternativeReport, repository string) {
	t.Helper()
	if repository == "" {
		if report.TopChoice != nil {
			t.Fatalf("unexpected top choice: %+v", *report.TopChoice)
		}
		return
	}
	if report.TopChoice == nil || alternativeIdentity(*report.TopChoice) != strings.ToLower(repository) {
		t.Fatalf("top choice = %+v, want %s", report.TopChoice, repository)
	}
}

func TestRankSoftwareAlternativesRejectsSealedMutations(t *testing.T) {
	cases := []struct {
		id  string
		run func(*testing.T)
	}{
		{"M01-stars-before-coverage", testMutationStarsBeforeCoverage},
		{"M02-stars-final-tiebreak", testMutationStarsFinalTiebreak},
		{"M03-license-bypass", testMutationLicenseBypass},
		{"M04a-maintenance-bypass", testMutationMaintenanceBypass},
		{"M04b-operation-bypass", testMutationOperationBypass},
		{"M05-target-self", testMutationTargetSelf},
		{"M06a-target-case", testMutationTargetCase},
		{"M07a-empty-commit", func(t *testing.T) { testMutationOpaqueCommit(t, "") }},
		{"M07b-head-commit", func(t *testing.T) { testMutationOpaqueCommit(t, "HEAD") }},
		{"M07c-divergent-commit", testMutationDivergentCommit},
		{"M07d-assessment-orphan", testMutationAssessmentMissing},
		{"M07e-snapshot-orphan", testMutationSnapshotMissing},
		{"M08a-capability-unknown", testMutationCapabilityUnknown},
		{"M08b-capability-conflict", testMutationCapabilityConflict},
		{"M08b2-alignment-conflict", testMutationAlignmentConflict},
		{"M08c-license-partial", testMutationLicenseOpaque},
		{"M08d-license-conflict", testMutationLicenseConflict},
		{"M08e-maintenance-partial", testMutationMaintenanceOpaque},
		{"M08f-operation-partial", testMutationOperationOpaque},
		{"M09-failure-ignored", testMutationFailureIgnored},
		{"M10-fallback-top-choice", testMutationFallbackTopChoice},
		{"M11a-coverage-duplicate", testMutationCoverageDuplicate},
		{"M11b-coverage-extra", testMutationCoverageExtra},
		{"M11c-coverage-omitted", testMutationCoverageOmitted},
		{"M12a-duplicate-assessment", testMutationDuplicateAssessment},
		{"M12b-duplicate-snapshot", testMutationDuplicateSnapshot},
		{"M12c-join-orphan", testMutationJoinOrphans},
		{"M13-nested-order-instability", testMutationNestedOrder},
		{"M14-identity-tiebreak", testMutationIdentityTiebreak},
	}
	killed := 0
	for _, mutation := range cases {
		mutation := mutation
		if t.Run(mutation.id, mutation.run) {
			killed++
		}
	}
	score := float64(killed) / float64(len(cases))
	if score != 1 {
		t.Fatalf("core mutation score = %d/%d = %.4f, want 1.0", killed, len(cases), score)
	}
}

func testMutationStarsBeforeCoverage(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.add(alternativeTestCandidate("org", "proven", "aaa", 1), alternativeTestAssessment("org", "proven", "aaa"))
	input.add(alternativeTestCandidate("org", "popular", "bbb", 1_000_000), alternativeTestAssessment("org", "popular", "bbb", domain.EvidenceProven, domain.EvidenceAbsent))
	report := rankAlternativeTestInput(t, input)
	if len(report.Alternatives) < 2 || alternativeIdentity(report.Alternatives[0].Repository) != "org/proven" {
		t.Fatalf("stars outranked coverage: %+v", report.Alternatives)
	}
	requireTopChoice(t, report, "org/proven")
}

func testMutationStarsFinalTiebreak(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.add(alternativeTestCandidate("org", "low-stars", "aaa", 1), alternativeTestAssessment("org", "low-stars", "aaa"))
	input.add(alternativeTestCandidate("org", "high-stars", "bbb", 50), alternativeTestAssessment("org", "high-stars", "bbb"))
	report := rankAlternativeTestInput(t, input)
	if len(report.Alternatives) < 2 || alternativeIdentity(report.Alternatives[0].Repository) != "org/high-stars" {
		t.Fatalf("stars did not break exact evidence tie: %+v", report.Alternatives)
	}
}

func testMutationLicenseBypass(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "gpl", "aaa")
	assessment.License.Identifier = "GPL-3.0"
	input.add(alternativeTestCandidate("org", "gpl", "aaa", 1_000_000), assessment)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/gpl", domain.AlternativeBlocked, true)
	requireTopChoice(t, report, "")
}

func testMutationMaintenanceBypass(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "archived", "aaa")
	assessment.Maintenance.Archived = true
	input.add(alternativeTestCandidate("org", "archived", "aaa", 500), assessment)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/archived", domain.AlternativeBlocked, true)
	requireTopChoice(t, report, "")
}

func testMutationOperationBypass(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.policy = alternativeTestPolicy(t, true, false, "MIT")
	assessment := alternativeTestAssessment("org", "cloud-only", "aaa")
	assessment.Operation.LocalStatus = domain.EvidenceAbsent
	input.add(alternativeTestCandidate("org", "cloud-only", "aaa", 500), assessment)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/cloud-only", domain.AlternativeBlocked, true)
}

func testMutationTargetSelf(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.target.Repository = &domain.RepositoryIdentity{Owner: "org", Repo: "self"}
	input.add(alternativeTestCandidate("org", "self", "new-commit", 1_000_000), alternativeTestAssessment("org", "self", "new-commit"))
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/self", domain.AlternativeNotAlternative, true)
	requireTopChoice(t, report, "")
}

func testMutationTargetCase(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.target.Repository = &domain.RepositoryIdentity{Owner: "ORG", Repo: "SELF"}
	input.add(alternativeTestCandidate("org", "self", "aaa", 10), alternativeTestAssessment("org", "self", "aaa"))
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/self", domain.AlternativeNotAlternative, true)
}

func testMutationOpaqueCommit(t *testing.T, commit string) {
	input := newAlternativeTestInput(t)
	input.add(alternativeTestCandidate("org", "opaque", commit, 10), alternativeTestAssessment("org", "opaque", commit))
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/opaque", domain.AlternativeInconclusive, true)
	requireTopChoice(t, report, "")
}

func testMutationDivergentCommit(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.add(alternativeTestCandidate("org", "divergent", "snapshot-sha", 10), alternativeTestAssessment("org", "divergent", "assessment-sha"))
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/divergent", domain.AlternativeInconclusive, true)
}

func testMutationAssessmentMissing(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.discovery.Candidates = append(input.discovery.Candidates, domain.DiscoveredCandidate{Repository: alternativeTestCandidate("org", "no-assessment", "aaa", 10)})
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/no-assessment", domain.AlternativeInconclusive, true)
}

func testMutationSnapshotMissing(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.assessments = append(input.assessments, alternativeTestAssessment("org", "no-snapshot", "aaa"))
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/no-snapshot", domain.AlternativeInconclusive, true)
}

func testMutationCapabilityUnknown(t *testing.T) {
	variants := map[string][]domain.CapabilityEvidence{
		"omitted":      alternativeTestAssessment("org", "opaque", "aaa").Capabilities[:1],
		"duplicated":   append(alternativeTestAssessment("org", "opaque", "aaa").Capabilities, alternativeTestAssessment("org", "opaque", "aaa").Capabilities[0]),
		"inconclusive": alternativeTestAssessment("org", "opaque", "aaa", domain.EvidenceProven, domain.EvidenceInconclusive).Capabilities,
	}
	for name, capabilities := range variants {
		t.Run(name, func(t *testing.T) {
			input := newAlternativeTestInput(t)
			assessment := alternativeTestAssessment("org", "opaque", "aaa")
			assessment.Capabilities = capabilities
			input.add(alternativeTestCandidate("org", "opaque", "aaa", 10), assessment)
			report := rankAlternativeTestInput(t, input)
			requireAlternativeFit(t, report, "org/opaque", domain.AlternativeInconclusive, true)
		})
	}
}

func testMutationCapabilityConflict(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.add(alternativeTestCandidate("org", "conflict", "aaa", 10), alternativeTestAssessment("org", "conflict", "aaa", domain.EvidenceProven, domain.EvidenceConflict))
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/conflict", domain.AlternativeBlocked, true)
}

func testMutationAlignmentConflict(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "conflict", "aaa")
	assessment.ProblemAlignment.Status = domain.EvidenceConflict
	input.add(alternativeTestCandidate("org", "conflict", "aaa", 10), assessment)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/conflict", domain.AlternativeBlocked, true)
}

func testMutationLicenseOpaque(t *testing.T) {
	for _, status := range []domain.EvidenceStatus{domain.EvidencePartial, domain.EvidenceAbsent, domain.EvidenceInconclusive} {
		t.Run(string(status), func(t *testing.T) {
			input := newAlternativeTestInput(t)
			assessment := alternativeTestAssessment("org", "opaque", "aaa")
			assessment.License.Status = status
			assessment.License.Identifier = ""
			input.add(alternativeTestCandidate("org", "opaque", "aaa", 10), assessment)
			alternative := requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/opaque", domain.AlternativeRecommended, false)
			if alternative.License.Status != status || alternative.License.Identifier != "" {
				t.Fatalf("license changed from unknown status %s: %+v", status, alternative.License)
			}
		})
	}
}

func testMutationLicenseConflict(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "blocked", "aaa")
	assessment.License.Status = domain.EvidenceConflict
	assessment.License.Identifier = ""
	input.add(alternativeTestCandidate("org", "blocked", "aaa", 10), assessment)
	requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/blocked", domain.AlternativeBlocked, true)
}

func testMutationMaintenanceOpaque(t *testing.T) {
	for _, status := range []domain.EvidenceStatus{domain.EvidencePartial, domain.EvidenceInconclusive} {
		t.Run(string(status), func(t *testing.T) {
			input := newAlternativeTestInput(t)
			assessment := alternativeTestAssessment("org", "opaque", "aaa")
			assessment.Maintenance.Status = status
			input.add(alternativeTestCandidate("org", "opaque", "aaa", 10), assessment)
			requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/opaque", domain.AlternativeInconclusive, true)
		})
	}
}

func testMutationOperationOpaque(t *testing.T) {
	for _, status := range []domain.EvidenceStatus{domain.EvidencePartial, domain.EvidenceInconclusive} {
		t.Run(string(status), func(t *testing.T) {
			input := newAlternativeTestInput(t)
			input.policy = alternativeTestPolicy(t, true, false, "MIT")
			assessment := alternativeTestAssessment("org", "opaque", "aaa")
			assessment.Operation.LocalStatus = status
			input.add(alternativeTestCandidate("org", "opaque", "aaa", 10), assessment)
			requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/opaque", domain.AlternativeInconclusive, true)
		})
	}
}

func testMutationFailureIgnored(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "failed", "aaa")
	assessment.Failures = []domain.ExecutionFailure{{Scope: "org/failed", Stage: "collection", Message: "timeout"}}
	input.add(alternativeTestCandidate("org", "failed", "aaa", 10), assessment)
	requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/failed", domain.AlternativeInconclusive, true)

	blocked := newAlternativeTestInput(t)
	assessment = alternativeTestAssessment("org", "blocked", "bbb")
	assessment.License.Status = domain.EvidenceConflict
	assessment.Failures = []domain.ExecutionFailure{{Scope: "org/blocked", Stage: "collection", Message: "timeout"}}
	blocked.add(alternativeTestCandidate("org", "blocked", "bbb", 10), assessment)
	requireAlternativeFit(t, rankAlternativeTestInput(t, blocked), "org/blocked", domain.AlternativeBlocked, true)
}

func testMutationFallbackTopChoice(t *testing.T) {
	input := newAlternativeTestInput(t)
	reference := alternativeTestAssessment("org", "reference", "aaa", domain.EvidenceProven, domain.EvidenceAbsent)
	reference.ProblemAlignment.Status = domain.EvidencePartial
	input.add(alternativeTestCandidate("org", "reference", "aaa", 10), reference)
	blocked := alternativeTestAssessment("org", "blocked", "bbb")
	blocked.License.Status = domain.EvidenceConflict
	input.add(alternativeTestCandidate("org", "blocked", "bbb", 20), blocked)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/reference", domain.AlternativeReference, false)
	requireAlternativeFit(t, report, "org/blocked", domain.AlternativeBlocked, true)
	requireTopChoice(t, report, "")
}

func testMutationCoverageDuplicate(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "duplicate", "aaa")
	assessment.Capabilities = append(assessment.Capabilities, assessment.Capabilities[0])
	input.add(alternativeTestCandidate("org", "duplicate", "aaa", 10), assessment)
	alternative := requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/duplicate", domain.AlternativeInconclusive, true)
	if alternative.Coverage > 1 {
		t.Fatalf("coverage exceeded 1: %v", alternative.Coverage)
	}
}

func testMutationCoverageExtra(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "extra", "aaa")
	assessment.Capabilities = append(assessment.Capabilities, domain.CapabilityEvidence{Capability: "unrequested feature", Status: domain.EvidenceProven})
	input.add(alternativeTestCandidate("org", "extra", "aaa", 10), assessment)
	alternative := requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/extra", domain.AlternativeRecommended, false)
	if alternative.Coverage != 1 {
		t.Fatalf("extra capability changed coverage: %v", alternative.Coverage)
	}
}

func testMutationCoverageOmitted(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "omitted", "aaa")
	assessment.Capabilities = assessment.Capabilities[:1]
	input.add(alternativeTestCandidate("org", "omitted", "aaa", 10), assessment)
	requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/omitted", domain.AlternativeInconclusive, true)
}

func testMutationDuplicateAssessment(t *testing.T) {
	input := newAlternativeTestInput(t)
	candidate := alternativeTestCandidate("org", "duplicate", "aaa", 10)
	assessment := alternativeTestAssessment("org", "duplicate", "aaa")
	input.add(candidate, assessment)
	input.assessments = append(input.assessments, assessment)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/duplicate", domain.AlternativeInconclusive, true)
	if len(report.Excluded) != 1 {
		t.Fatalf("duplicate assessment emitted %d records", len(report.Excluded))
	}
}

func testMutationDuplicateSnapshot(t *testing.T) {
	input := newAlternativeTestInput(t)
	candidate := alternativeTestCandidate("org", "duplicate", "aaa", 10)
	assessment := alternativeTestAssessment("org", "duplicate", "aaa")
	input.add(candidate, assessment)
	input.discovery.Candidates = append(input.discovery.Candidates, domain.DiscoveredCandidate{Repository: candidate})
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/duplicate", domain.AlternativeInconclusive, true)
	if len(report.Excluded) != 1 {
		t.Fatalf("duplicate snapshot emitted %d records", len(report.Excluded))
	}
}

func testMutationJoinOrphans(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.discovery.Candidates = append(input.discovery.Candidates, domain.DiscoveredCandidate{Repository: alternativeTestCandidate("org", "snapshot-only", "aaa", 10)})
	input.assessments = append(input.assessments, alternativeTestAssessment("org", "assessment-only", "bbb"))
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/snapshot-only", domain.AlternativeInconclusive, true)
	requireAlternativeFit(t, report, "org/assessment-only", domain.AlternativeInconclusive, true)
	if len(report.Excluded) != 2 {
		t.Fatalf("orphans emitted %d records", len(report.Excluded))
	}
}

func testMutationNestedOrder(t *testing.T) {
	input := newAlternativeTestInput(t)
	first := alternativeTestAssessment("org", "alpha", "aaa")
	first.Capabilities[0].Sources = []domain.EvidenceSource{
		{Kind: domain.EvidenceTest, Ref: first.Repository, Path: "z_test.go", RuleID: "CAP-Z", Strength: .9},
		{Kind: domain.EvidenceFile, Ref: first.Repository, Path: "a.go", RuleID: "CAP-A", Strength: 1},
	}
	input.add(alternativeTestCandidate("org", "zulu", "zzz", 10), alternativeTestAssessment("org", "zulu", "zzz"))
	input.add(alternativeTestCandidate("org", "alpha", "aaa", 10), first)
	reportA := rankAlternativeTestInput(t, input)

	input.discovery.Candidates[0], input.discovery.Candidates[1] = input.discovery.Candidates[1], input.discovery.Candidates[0]
	input.assessments[0], input.assessments[1] = input.assessments[1], input.assessments[0]
	for index := range input.assessments {
		capabilities := input.assessments[index].Capabilities
		for left, right := 0, len(capabilities)-1; left < right; left, right = left+1, right-1 {
			capabilities[left], capabilities[right] = capabilities[right], capabilities[left]
		}
		for capabilityIndex := range capabilities {
			sources := capabilities[capabilityIndex].Sources
			for left, right := 0, len(sources)-1; left < right; left, right = left+1, right-1 {
				sources[left], sources[right] = sources[right], sources[left]
			}
		}
	}
	reportB := rankAlternativeTestInput(t, input)
	a, _ := json.Marshal(reportA)
	b, _ := json.Marshal(reportB)
	if string(a) != string(b) {
		t.Fatalf("input permutation changed canonical report:\n%s\n%s", a, b)
	}
}

func testMutationIdentityTiebreak(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.add(alternativeTestCandidate("org", "zulu", "zzz", 10), alternativeTestAssessment("org", "zulu", "zzz"))
	input.add(alternativeTestCandidate("org", "alpha", "aaa", 10), alternativeTestAssessment("org", "alpha", "aaa"))
	report := rankAlternativeTestInput(t, input)
	if len(report.Alternatives) != 2 || alternativeIdentity(report.Alternatives[0].Repository) != "org/alpha" {
		t.Fatalf("identity tie-break is unstable: %+v", report.Alternatives)
	}
}

func TestRankSoftwareAlternativesIgnoresOperationOpacityWhenLocalIsNotRequired(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "remote-ok", "aaa")
	assessment.Operation.LocalStatus = domain.EvidenceInconclusive
	input.add(alternativeTestCandidate("org", "remote-ok", "aaa", 10), assessment)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/remote-ok", domain.AlternativeRecommended, false)
	requireTopChoice(t, report, "org/remote-ok")
}

func TestRankSoftwareAlternativesProjectsNormalizedPolicyFilters(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.policy = alternativeTestPolicy(t, true, true, "mit", "Apache-2.0", "MIT")
	report := rankAlternativeTestInput(t, input)
	if !report.Filters.RequireLocal || !report.Filters.AllowArchived || report.Filters.RecommendedMinCoverage != .8 || report.Filters.PartialMinCoverage != .5 {
		t.Fatalf("policy projection = %+v", report.Filters)
	}
	want := []string{"Apache-2.0", "MIT"}
	if strings.Join(report.Filters.AllowedLicenses, ",") != strings.Join(want, ",") {
		t.Fatalf("normalized licenses = %v, want %v", report.Filters.AllowedLicenses, want)
	}
}

func TestRankSoftwareAlternativesValidatesSealedInputs(t *testing.T) {
	base := newAlternativeTestInput(t)
	cases := []struct {
		name   string
		mutate func(*alternativeTestInput)
	}{
		{"empty target", func(input *alternativeTestInput) { input.target.Name = "" }},
		{"no capabilities", func(input *alternativeTestInput) { input.idea.Capabilities = nil }},
		{"empty policy schema", func(input *alternativeTestInput) { input.policy.SchemaVersion = "" }},
		{"zero recommended", func(input *alternativeTestInput) { input.policy.ExtendMinCoverage = 0 }},
		{"nan partial", func(input *alternativeTestInput) { input.policy.IntegrationMinValue = math.NaN() }},
		{"partial above recommended", func(input *alternativeTestInput) {
			input.policy.ExtendMinCoverage = .5
			input.policy.IntegrationMinValue = .8
		}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			input := base
			testCase.mutate(&input)
			if _, err := RankSoftwareAlternatives(input.target, input.idea, input.discovery, input.assessments, input.policy); err == nil {
				t.Fatal("invalid input accepted")
			}
		})
	}
}

const alternativeM20Warning = "licença não reconhecida; valide os termos antes da adoção"

func TestRankSoftwareAlternativesM20RejectsLicenseMutations(t *testing.T) {
	mutations := []struct {
		id  string
		run func(*testing.T)
	}{
		{"M20a-partial-warning-only", func(t *testing.T) { testM20UnknownLicense(t, domain.EvidencePartial) }},
		{"M20b-absent-warning-only", func(t *testing.T) { testM20UnknownLicense(t, domain.EvidenceAbsent) }},
		{"M20c-inconclusive-warning-only", func(t *testing.T) { testM20UnknownLicense(t, domain.EvidenceInconclusive) }},
		{"M20d-warning-canonical-and-allowlist-safe", testM20WarningCanonical},
		{"M20e-proven-outside-allowlist-blocked", testM20OutsideAllowlistBlocked},
		{"M20f-conflict-blocked", testM20ConflictBlocked},
		{"M20g-proven-allowed-has-no-warning", testM20AllowedNoWarning},
		{"M20h-invalid-status-remains-opaque", testM20InvalidStatusOpaque},
		{"M20i-recipesage-and-json-reorder", testM20RecipeSageAndReorder},
	}
	killed := 0
	for _, mutation := range mutations {
		mutation := mutation
		if t.Run(mutation.id, mutation.run) {
			killed++
		}
	}
	if score := float64(killed) / float64(len(mutations)); score != 1 {
		t.Fatalf("M20 mutation score = %d/%d = %.4f, want 1.0", killed, len(mutations), score)
	}
}

func testM20UnknownLicense(t *testing.T, status domain.EvidenceStatus) {
	input := newAlternativeTestInput(t)
	repo := strings.ToLower(string(status))
	assessment := alternativeTestAssessment("org", repo, "abc123")
	assessment.License = domain.LicenseAssessment{Status: status}
	input.add(alternativeTestCandidate("org", repo, "abc123", 10), assessment)
	alternative := requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/"+repo, domain.AlternativeRecommended, false)
	if alternative.License.Status != status || alternative.License.Identifier != "" {
		t.Fatalf("unknown license changed: got %+v, want status %s and empty identifier", alternative.License, status)
	}
	requireM20Warning(t, alternative)
	if len(alternative.Blockers) != 0 {
		t.Fatalf("unknown license added blockers: %+v", alternative.Blockers)
	}
}

func testM20WarningCanonical(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "manual-review", "abc123")
	assessment.License = domain.LicenseAssessment{Status: domain.EvidenceInconclusive}
	input.add(alternativeTestCandidate("org", "manual-review", "abc123", 10), assessment)
	alternative := requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/manual-review", domain.AlternativeRecommended, false)
	requireM20Warning(t, alternative)
	encoded, err := json.Marshal(alternative)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(encoded), alternativeM20Warning) != 1 {
		t.Fatalf("canonical warning missing or duplicated in JSON: %s", encoded)
	}
}

func testM20OutsideAllowlistBlocked(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "gpl", "abc123")
	assessment.License = domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "GPL-3.0", Sources: []domain.EvidenceSource{alternativeTestSource(assessment.Repository, domain.EvidenceLicense, "LICENSE", "license", "GNU GENERAL PUBLIC LICENSE Version 3")}}
	input.add(alternativeTestCandidate("org", "gpl", "abc123", 100), assessment)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/gpl", domain.AlternativeBlocked, true)
	requireTopChoice(t, report, "")
}

func testM20ConflictBlocked(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "conflict", "abc123")
	assessment.License = domain.LicenseAssessment{Status: domain.EvidenceConflict}
	input.add(alternativeTestCandidate("org", "conflict", "abc123", 100), assessment)
	report := rankAlternativeTestInput(t, input)
	requireAlternativeFit(t, report, "org/conflict", domain.AlternativeBlocked, true)
	requireTopChoice(t, report, "")
}

func testM20AllowedNoWarning(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.add(alternativeTestCandidate("org", "mit", "abc123", 10), alternativeTestAssessment("org", "mit", "abc123"))
	alternative := requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/mit", domain.AlternativeRecommended, false)
	if alternative.License.Status != domain.EvidenceProven || alternative.License.Identifier != "MIT" {
		t.Fatalf("allowed PROVEN license changed: %+v", alternative.License)
	}
	for _, reason := range alternative.Reasons {
		if reason == alternativeM20Warning {
			t.Fatalf("allowed license received unknown warning: %+v", alternative.Reasons)
		}
	}
}

func testM20InvalidStatusOpaque(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "invalid", "abc123")
	assessment.License = domain.LicenseAssessment{Status: domain.EvidenceStatus("INVALID")}
	input.add(alternativeTestCandidate("org", "invalid", "abc123", 10), assessment)
	alternative := requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/invalid", domain.AlternativeInconclusive, true)
	if alternative.License.Status != domain.EvidenceStatus("INVALID") || alternative.License.Identifier != "" {
		t.Fatalf("invalid license status was rewritten: %+v", alternative.License)
	}
	for _, reason := range alternative.Reasons {
		if reason == alternativeM20Warning {
			t.Fatalf("invalid status was silently accepted as recognized unknown: %+v", alternative.Reasons)
		}
	}
}

func testM20RecipeSageAndReorder(t *testing.T) {
	input := m20RecipeSageInput(t)
	first := rankAlternativeTestInput(t, input)
	recipeSage := requireAlternativeFit(t, first, "julianpoy/RecipeSage", domain.AlternativeReference, false)
	requireM20Warning(t, recipeSage)
	requireAlternativeFit(t, first, "norish-recipes/norish", domain.AlternativeRecommended, false)
	requireTopChoice(t, first, "norish-recipes/norish")
	if len(first.Alternatives) != 2 || alternativeIdentity(first.Alternatives[0].Repository) != "norish-recipes/norish" || alternativeIdentity(first.Alternatives[1].Repository) != "julianpoy/recipesage" {
		t.Fatalf("witness order changed: %+v", first.Alternatives)
	}

	for left, right := 0, len(input.discovery.Candidates)-1; left < right; left, right = left+1, right-1 {
		input.discovery.Candidates[left], input.discovery.Candidates[right] = input.discovery.Candidates[right], input.discovery.Candidates[left]
	}
	for left, right := 0, len(input.assessments)-1; left < right; left, right = left+1, right-1 {
		input.assessments[left], input.assessments[right] = input.assessments[right], input.assessments[left]
	}
	second := rankAlternativeTestInput(t, input)
	a, err := json.Marshal(first)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(second)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatalf("M20 input reorder changed canonical JSON:\n%s\n%s", a, b)
	}
}

func m20RecipeSageInput(t *testing.T) alternativeTestInput {
	t.Helper()
	input := newAlternativeTestInput(t)
	const norishCommit = "e3691fe08addabd42feb68e8cc942fcd6468634b"
	input.add(alternativeTestCandidate("norish-recipes", "norish", norishCommit, 10), alternativeTestAssessment("norish-recipes", "norish", norishCommit))
	const recipeSageCommit = "120aa5c8816f9a17e043d48c653a6ce200937ea9"
	recipeSage := alternativeTestAssessment("julianpoy", "RecipeSage", recipeSageCommit, domain.EvidenceProven, domain.EvidenceAbsent)
	recipeSage.ProblemAlignment.Status = domain.EvidencePartial
	recipeSage.License = domain.LicenseAssessment{Status: domain.EvidenceInconclusive}
	input.add(alternativeTestCandidate("julianpoy", "RecipeSage", recipeSageCommit, 948), recipeSage)
	return input
}

func requireM20Warning(t *testing.T, alternative domain.SoftwareAlternative) {
	t.Helper()
	count := 0
	for _, reason := range alternative.Reasons {
		if reason == alternativeM20Warning {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("warning count = %d, reasons=%+v", count, alternative.Reasons)
	}
}

func TestM21RegionalOriginsAreProjectedButRankingNeutral(t *testing.T) {
	input := newAlternativeTestInput(t)
	input.add(alternativeTestCandidate("org", "beta", "bbb", 10), alternativeTestAssessment("org", "beta", "bbb"))
	input.add(alternativeTestCandidate("org", "alpha", "aaa", 10), alternativeTestAssessment("org", "alpha", "aaa"))
	input.discovery.Candidates[0].Origins = []domain.CandidateOrigin{
		{ExpansionID: "regional-ja-JP", PlannedQuery: "レシピ管理", EffectiveQuery: "レシピ管理", SearchLocale: "ja-JP", RegionalProposalID: "regional-ja-JP-recipe-manager"},
		{ExpansionID: "regional-pt-BR", PlannedQuery: "gerenciador de receitas", EffectiveQuery: "gerenciador de receitas", SearchLocale: "pt-BR", RegionalProposalID: "regional-pt-BR-recipe-manager"},
		{ExpansionID: "regional-ja-JP", PlannedQuery: "レシピ管理", EffectiveQuery: "レシピ管理", SearchLocale: "ja-JP", RegionalProposalID: "regional-ja-JP-recipe-manager"},
	}
	input.discovery.Candidates[1].Origins = []domain.CandidateOrigin{
		{ExpansionID: "regional-de-DE", PlannedQuery: "Rezeptmanager", EffectiveQuery: "Rezeptmanager", SearchLocale: "de-DE", RegionalProposalID: "regional-de-DE-recipe-manager"},
	}

	regional := rankAlternativeTestInput(t, input)
	if len(regional.Alternatives) != 2 || alternativeIdentity(regional.Alternatives[0].Repository) != "org/alpha" || alternativeIdentity(regional.Alternatives[1].Repository) != "org/beta" {
		t.Fatalf("locale changed identity tie-break: %+v", regional.Alternatives)
	}
	if !reflect.DeepEqual(regional.Alternatives[0].DiscoveryLocales, []string{"de-DE"}) || !reflect.DeepEqual(regional.Alternatives[1].DiscoveryLocales, []string{"pt-BR", "ja-JP"}) {
		t.Fatalf("regional locale projection is not canonical: %+v", regional.Alternatives)
	}

	without := input
	without.discovery.Candidates = append([]domain.DiscoveredCandidate(nil), input.discovery.Candidates...)
	for index := range without.discovery.Candidates {
		without.discovery.Candidates[index].Origins = nil
	}
	legacy := rankAlternativeTestInput(t, without)
	for index := range regional.Alternatives {
		got := regional.Alternatives[index]
		want := legacy.Alternatives[index]
		got.DiscoveryLocales = nil
		want.DiscoveryLocales = nil
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("locale changed non-provenance ranking fields:\nregional=%+v\nlegacy=%+v", got, want)
		}
	}

	reordered := input
	reordered.discovery.Candidates = append([]domain.DiscoveredCandidate(nil), input.discovery.Candidates...)
	for index := range reordered.discovery.Candidates {
		origins := append([]domain.CandidateOrigin(nil), reordered.discovery.Candidates[index].Origins...)
		for left, right := 0, len(origins)-1; left < right; left, right = left+1, right-1 {
			origins[left], origins[right] = origins[right], origins[left]
		}
		reordered.discovery.Candidates[index].Origins = origins
	}
	reordered.discovery.Candidates[0], reordered.discovery.Candidates[1] = reordered.discovery.Candidates[1], reordered.discovery.Candidates[0]
	reordered.assessments = append([]domain.CandidateAssessment(nil), input.assessments...)
	reordered.assessments[0], reordered.assessments[1] = reordered.assessments[1], reordered.assessments[0]
	a, _ := json.Marshal(regional)
	b, _ := json.Marshal(rankAlternativeTestInput(t, reordered))
	if string(a) != string(b) {
		t.Fatalf("regional origin reorder changed alternatives JSON:\n%s\n%s", a, b)
	}
}

func TestM21LocalizedDiscoveryNeverPromotesFunctionalUnknown(t *testing.T) {
	input := newAlternativeTestInput(t)
	assessment := alternativeTestAssessment("org", "localized-only", "abc123", domain.EvidenceInconclusive, domain.EvidenceInconclusive)
	assessment.ProblemAlignment.Status = domain.EvidenceInconclusive
	input.add(alternativeTestCandidate("org", "localized-only", "abc123", 10), assessment)
	input.discovery.Candidates[0].Origins = []domain.CandidateOrigin{{ExpansionID: "regional-ja-JP", PlannedQuery: "レシピ管理", EffectiveQuery: "レシピ管理", SearchLocale: "ja-JP", RegionalProposalID: "regional-ja-JP-recipe-manager"}}
	alternative := requireAlternativeFit(t, rankAlternativeTestInput(t, input), "org/localized-only", domain.AlternativeInconclusive, true)
	if alternative.Coverage != 0 || alternative.ProblemAlignment != domain.EvidenceInconclusive || len(alternative.Blockers) != 0 || !reflect.DeepEqual(alternative.DiscoveryLocales, []string{"ja-JP"}) {
		t.Fatalf("localized text was promoted to functional proof: %+v", alternative)
	}
}

type alternativeCapabilityJSON struct {
	Capability string                  `json:"capability"`
	Status     domain.EvidenceStatus   `json:"status"`
	Evidence   []domain.EvidenceSource `json:"evidence"`
}

type alternativeMatrixCandidateJSON struct {
	Repository   domain.CommitRef            `json:"repository"`
	Capabilities []alternativeCapabilityJSON `json:"capabilities"`
}

type alternativeMatrixReportJSON struct {
	Alternatives []alternativeMatrixCandidateJSON `json:"alternatives"`
	Excluded     []alternativeMatrixCandidateJSON `json:"excluded"`
}

func TestRankSoftwareAlternativesM17PreservesCanonicalCapabilityMatrix(t *testing.T) {
	base := alternativeM17Input(t, false, false)
	report := rankAlternativeTestInput(t, base)
	assertAlternativeM17Matrix(t, report, map[string]struct {
		status domain.EvidenceStatus
		path   string
	}{
		"dependency graph": {status: domain.EvidenceProven, path: "graph.go"},
		"impact analysis":  {status: domain.EvidencePartial, path: "impact.go"},
	})

	reordered := alternativeM17Input(t, true, false)
	a, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(rankAlternativeTestInput(t, reordered))
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatalf("capability/input reorder changed canonical JSON:\n%s\n%s", a, b)
	}

	for _, mutation := range []struct {
		name       string
		duplicate  bool
		capability string
	}{
		{name: "omitted", capability: "impact analysis"},
		{name: "duplicated", duplicate: true, capability: "dependency graph"},
	} {
		mutation := mutation
		t.Run(mutation.name, func(t *testing.T) {
			input := alternativeM17Input(t, false, mutation.duplicate)
			if !mutation.duplicate {
				input.assessments[0].Capabilities = input.assessments[0].Capabilities[:2]
			}
			candidate := alternativeM17Candidate(t, rankAlternativeTestInput(t, input))
			matrix := map[string]alternativeCapabilityJSON{}
			for _, capability := range candidate.Capabilities {
				matrix[capability.Capability] = capability
			}
			if len(matrix) != 2 {
				t.Fatalf("required capability matrix = %+v", candidate.Capabilities)
			}
			opaque, ok := matrix[mutation.capability]
			if !ok || opaque.Status != domain.EvidenceInconclusive {
				t.Fatalf("%s did not remain identifiable as INCONCLUSIVE: %+v", mutation.capability, candidate.Capabilities)
			}
		})
	}
}

func alternativeM17Input(t *testing.T, reordered, duplicate bool) alternativeTestInput {
	t.Helper()
	input := newAlternativeTestInput(t)
	input.idea.Capabilities = []string{" Dependency   Graph ", "dependency graph", "Impact Analysis"}
	input.discovery.Idea = input.idea
	assessment := alternativeTestAssessment("org", "matrix", "abc123")
	ref := assessment.Repository
	assessment.Capabilities = []domain.CapabilityEvidence{
		{Capability: "dependency graph", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{alternativeTestSource(ref, domain.EvidenceSymbol, "graph.go", "CAP-GRAPH", "func BuildGraph() {}")}},
		{Capability: "unrequested feature", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{alternativeTestSource(ref, domain.EvidenceSymbol, "extra.go", "CAP-EXTRA", "func ExtraFeature() {}")}},
		{Capability: "impact analysis", Status: domain.EvidencePartial, Sources: []domain.EvidenceSource{alternativeTestSource(ref, domain.EvidenceDocumentation, "impact.go", "CAP-IMPACT", "This tool provides impact analysis.")}},
	}
	if duplicate {
		assessment.Capabilities = append(assessment.Capabilities, assessment.Capabilities[0])
	}
	input.add(alternativeTestCandidate("org", "matrix", "abc123", 10), assessment)
	if reordered {
		for left, right := 0, len(input.idea.Capabilities)-1; left < right; left, right = left+1, right-1 {
			input.idea.Capabilities[left], input.idea.Capabilities[right] = input.idea.Capabilities[right], input.idea.Capabilities[left]
		}
		input.discovery.Idea = input.idea
		capabilities := input.assessments[0].Capabilities
		for left, right := 0, len(capabilities)-1; left < right; left, right = left+1, right-1 {
			capabilities[left], capabilities[right] = capabilities[right], capabilities[left]
		}
	}
	return input
}

func alternativeM17Candidate(t *testing.T, report domain.SoftwareAlternativeReport) alternativeMatrixCandidateJSON {
	t.Helper()
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var decoded alternativeMatrixReportJSON
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, candidate := range append(decoded.Alternatives, decoded.Excluded...) {
		if alternativeIdentity(candidate.Repository) == "org/matrix" {
			return candidate
		}
	}
	t.Fatalf("org/matrix absent from report JSON: %s", encoded)
	return alternativeMatrixCandidateJSON{}
}

func assertAlternativeM17Matrix(t *testing.T, report domain.SoftwareAlternativeReport, expected map[string]struct {
	status domain.EvidenceStatus
	path   string
}) {
	t.Helper()
	candidate := alternativeM17Candidate(t, report)
	if len(candidate.Capabilities) != len(expected) {
		t.Fatalf("capability matrix = %+v, want %d required rows", candidate.Capabilities, len(expected))
	}
	seen := map[string]bool{}
	for _, capability := range candidate.Capabilities {
		want, ok := expected[capability.Capability]
		if !ok {
			t.Fatalf("extra or unnamed capability leaked: %+v", capability)
		}
		if seen[capability.Capability] {
			t.Fatalf("duplicate canonical capability: %+v", candidate.Capabilities)
		}
		seen[capability.Capability] = true
		if capability.Status != want.status {
			t.Fatalf("%s status = %s, want %s", capability.Capability, capability.Status, want.status)
		}
		if len(capability.Evidence) != 1 || capability.Evidence[0].Path != want.path {
			t.Fatalf("%s evidence = %+v, want path %s", capability.Capability, capability.Evidence, want.path)
		}
	}
}
