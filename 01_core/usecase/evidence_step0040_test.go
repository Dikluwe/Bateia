/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_step0040_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"strings"
	"testing"
)

func TestStep0040EchoDoesNotProveAsynchronousBehavior(t *testing.T) {
	b := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "events.go", Content: "package events\nfunc EventProcessor(value string) string { return value }"},
		domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "events_test.go", Content: "func TestEventProcessor(t *testing.T) { if EventProcessor(\"x\") != \"x\" {t.Fatal(\"echo failed\")} }"},
		domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Asynchronous event processing."})
	got := assessCapability("asynchronous event processing", b, findRule("asynchronous event processing", DefaultEvidenceRules()))
	if got.Status != domain.EvidencePartial || !strings.Contains(got.Explanation, "comportamento não verificado") || got.Confidence > .65 {
		t.Fatalf("lexical evidence promoted: %+v", got)
	}
	if len(got.Sources) < 2 {
		t.Fatal("lost useful locations")
	}
}

func TestStep0040LicenseAuthorityAndScope(t *testing.T) {
	content := step0025MITLicense
	for _, path := range []string{"license_checker_test.go", "src/LICENSE", "vendor/LICENSE", "examples/LICENSE", "docs/LICENSE"} {
		// Even an incorrectly tagged item from a saved bundle must not establish authority.
		b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceLicense, Path: path, Content: content})
		if got := assessLicense(b); got.Status == domain.EvidenceProven {
			t.Fatalf("unauthorized path %s: %+v", path, got)
		}
		if classifyEvidencePath(path) == domain.EvidenceLicense {
			t.Fatalf("misclassified %s", path)
		}
	}
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceLicense, Path: "LICENSE.md", Content: content})
	if got := assessLicense(b); got.Status != domain.EvidenceProven || got.Identifier != "MIT" {
		t.Fatalf("root license lost: %+v", got)
	}
}

func TestStep0040NonGoalMentionsAndRealLine(t *testing.T) {
	for _, sentence := range []string{"This tool does not upload code.", "Esta ferramenta não faz upload code.", "This tool can upload code."} {
		b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "# Project\n\n" + sentence})
		got := assessNonGoal("upload code", b)
		if got.Status != domain.EvidenceInconclusive || len(got.Sources) != 1 || got.Sources[0].Line != 3 || got.Sources[0].Snippet != sentence {
			t.Fatalf("mention misrepresented: %+v", got)
		}
	}
	item := domain.CollectedEvidence{Content: "hello", Path: "README.md"}
	if s := sourceFromItem(item, domain.EvidenceDocumentation, "test", .4, "not present"); s.Line != 0 {
		t.Fatal("invented line")
	}
}

func TestStep0040ProductionBeforeValidationArtifacts(t *testing.T) {
	targets := []domain.ExtractionTarget{{Capabilities: []string{"event processing"}}}
	tree := []domain.DirectoryNode{{Path: "tmp/events.go"}, {Path: "benchmark/events.py"}, {Path: "generated/events.cpp"}, {Path: "src/events.cpp"}}
	got := prioritizeEvidenceTree(tree, targets)
	if got[0].Path != "src/events.cpp" {
		t.Fatalf("wrong priority: %+v", got)
	}
	explicit := []domain.ExtractionTarget{{Capabilities: []string{"benchmark event processing"}}}
	if evidencePathPriority("benchmark/events.py", explicit) < 0 {
		t.Fatal("explicit benchmark target penalized")
	}
	if tree[0].Path != "tmp/events.go" {
		t.Fatal("mutated input")
	}
}

func TestStep0040UnknownGateIsNotCapabilityGap(t *testing.T) {
	c := decisionCandidate{key: "org/repo@abc", a: domain.CandidateAssessment{
		License:          domain.LicenseAssessment{Status: domain.EvidenceProven},
		Maintenance:      domain.MaintenanceAssessment{Status: domain.EvidenceInconclusive},
		ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceProven},
		Capabilities:     []domain.CapabilityEvidence{{Capability: "export", Status: domain.EvidenceProven}},
	}}
	reasons := unresolvedCandidateReasons(c, domain.ReuseDecisionPolicy{})
	if len(reasons) != 1 || !strings.Contains(reasons[0].Message, "manutenção") || reasons[0].Candidate != c.key {
		t.Fatalf("not actionable: %+v", reasons)
	}
}

func TestStep0040MaintenanceUnknownDoesNotEraseProvenCoverage(t *testing.T) {
	a := decisionAssessment("repo", domain.EvidenceProven, domain.EvidenceProven, domain.EvidenceProven)
	a.Maintenance.Status = domain.EvidenceInconclusive
	got, err := DecideReuse(decisionRequest(a))
	if err != nil {
		t.Fatal(err)
	}
	if got.Decision != domain.DecisionUnknown || len(got.Gaps) != 0 {
		t.Fatalf("gate became functional gap: %+v", got)
	}
	found := false
	for _, reason := range got.Reasons {
		found = found || strings.Contains(reason.Message, "manutenção não verificada")
	}
	if !found {
		t.Fatal("missing resolution action")
	}
}

func TestStep0040QueriesDoNotBecomeCandidateCoverage(t *testing.T) {
	snapshot := domain.DiscoverySnapshot{Idea: domain.IdeaSpec{Capabilities: []string{"export markdown"}, NonGoals: []string{"upload code"}}}
	candidate := domain.DiscoveredCandidate{Queries: []string{"export markdown upload code"}}
	got := ExtractScreeningSignals(snapshot, candidate, DefaultScreeningModel())
	if got.CapabilityCoverage != 0 || got.NonGoalConflict != 0 {
		t.Fatalf("queries used as content: %+v", got)
	}
	candidate.Repository.Description = "export markdown"
	got = ExtractScreeningSignals(snapshot, candidate, DefaultScreeningModel())
	if got.CapabilityCoverage != 1 {
		t.Fatal("lost observed description signal")
	}
}

func TestStep0040BudgetOmissionIsNotCompleteAnalysis(t *testing.T) {
	plan := domain.EvidenceExtractionPlan{Repository: domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc"}, Candidate: domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "abc"}, Budget: domain.ExtractionBudget{TreeDepth: 4, TotalBytes: 4096}}
	bundle, err := CollectEvidence(context.Background(), plan, evidenceProvider{})
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Complete || len(bundle.Targets) != 2 {
		t.Fatalf("omission hidden: %+v", bundle)
	}
	for _, target := range bundle.Targets {
		if target.Status != domain.TargetSkippedBudget {
			t.Fatalf("wrong state: %+v", target)
		}
	}
	got := assessCapability("dependency graph", bundle, findRule("dependency graph", DefaultEvidenceRules()))
	if got.Status != domain.EvidenceInconclusive {
		t.Fatalf("budget became absence: %+v", got)
	}
}
