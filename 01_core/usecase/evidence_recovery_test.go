/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_recovery_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"reflect"
	"testing"
)

func recoveryDiagnosisFixture() (domain.ReuseDecisionReport, domain.AssessmentReport, []domain.EvidenceBundle) {
	ref := domain.CommitRef{Owner: "o", Repo: "r", Commit: "abc"}
	doc := decisionSource("r", "cap")
	doc.Kind = domain.EvidenceDocumentation
	a := domain.CandidateAssessment{Repository: ref, Capabilities: []domain.CapabilityEvidence{{Capability: "dependency graph", Status: domain.EvidenceInconclusive, Sources: []domain.EvidenceSource{doc}}}, License: domain.LicenseAssessment{Status: domain.EvidenceInconclusive}, Operation: domain.OperationAssessment{LocalStatus: domain.EvidenceInconclusive}}
	b := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: false}
	return domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{a}}, []domain.EvidenceBundle{b}
}
func TestDiagnoseInconclusiveSeparatesCapabilityLicenseAndOperation(t *testing.T) {
	decision, assessment, bundles := recoveryDiagnosisFixture()
	d := DiagnoseInconclusive("run", "hash", decision, assessment, bundles, DefaultEvidenceRules())
	causes := map[domain.InconclusiveCause]bool{}
	for _, f := range d.Findings {
		causes[f.Cause] = true
		if f.ID == "" || f.RuleID == "" {
			t.Fatalf("untraceable %#v", f)
		}
	}
	for _, c := range []domain.InconclusiveCause{domain.CauseMissingImplementation, domain.CauseUnknownLicense, domain.CauseUnknownOperation} {
		if !causes[c] {
			t.Fatalf("missing %s: %#v", c, d)
		}
	}
}
func TestConclusiveDecisionProducesNoRecoveryTargets(t *testing.T) {
	_, assessment, bundles := recoveryDiagnosisFixture()
	d := DiagnoseInconclusive("run", "hash", domain.ReuseDecisionReport{Decision: domain.DecisionBuild}, assessment, bundles, DefaultEvidenceRules())
	p, err := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Targets) != 0 {
		t.Fatalf("got %#v", p.Targets)
	}
}
func TestRecoveryPlanIsOrderInvariantAndEveryTargetHasFinding(t *testing.T) {
	decision, assessment, bundles := recoveryDiagnosisFixture()
	d := DiagnoseInconclusive("run", "hash", decision, assessment, bundles, DefaultEvidenceRules())
	p1, err := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i, j := 0, len(d.Findings)-1; i < j; i, j = i+1, j-1 {
		d.Findings[i], d.Findings[j] = d.Findings[j], d.Findings[i]
	}
	p2, _ := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 1, nil)
	if !reflect.DeepEqual(p1, p2) {
		t.Fatal("finding order changed plan")
	}
	for _, x := range p1.Targets {
		if x.FindingID == "" {
			t.Fatal("target without finding")
		}
	}
}
func TestRecoveryPlanDoesNotRepeatTargetAndHonorsRounds(t *testing.T) {
	decision, assessment, bundles := recoveryDiagnosisFixture()
	d := DiagnoseInconclusive("run", "hash", decision, assessment, bundles, DefaultEvidenceRules())
	p, _ := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 1, nil)
	again, _ := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 2, p.Targets)
	seen := map[string]bool{}
	for _, x := range p.Targets {
		seen[targetIdentity(x)] = true
	}
	for _, x := range again.Targets {
		if seen[targetIdentity(x)] {
			t.Fatalf("repeated %#v", x)
		}
	}
	if _, err := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 3, nil); err == nil {
		t.Fatal("third round accepted")
	}
}
func TestRuleGapIsTerminalAndCreatesNoTarget(t *testing.T) {
	decision, assessment, bundles := recoveryDiagnosisFixture()
	assessment.Assessments[0].Capabilities[0].Capability = "capability with no catalog mapping whatsoever"
	d := DiagnoseInconclusive("run", "hash", decision, assessment, bundles, DefaultEvidenceRules())
	p, _ := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 1, nil)
	for _, x := range p.Targets {
		if x.Capability == assessment.Assessments[0].Capabilities[0].Capability {
			t.Fatal("rule gap created target")
		}
	}
}

func TestMergeRecoveryNeverDegradesOriginalCompleteness(t *testing.T) {
	ref := domain.CommitRef{Owner: "o", Repo: "r", Commit: "abc"}
	original := domain.EvidenceBundle{Repository: ref, Complete: true, TreeComplete: true}
	attempt := domain.EvidenceBundle{Repository: ref, Complete: false, TreeComplete: false, Failures: []domain.ExecutionFailure{{Stage: "recovery", Message: "missing"}}}
	merged, _ := mergeRecoveryBundles([]domain.EvidenceBundle{original}, []domain.EvidenceBundle{attempt})
	if !merged[0].Complete || !merged[0].TreeComplete {
		t.Fatal("additional failure degraded original conclusion")
	}
}

func TestRecoveryUsefulEvidenceRateIsBounded(t *testing.T) {
	plan := domain.EvidenceRecoveryPlan{Round: 1, Targets: []domain.RecoveryTarget{{FindingID: "a"}, {FindingID: "b"}}}
	d := domain.InconclusiveDiagnosis{Findings: []domain.InconclusiveFinding{{ID: "b"}}}
	m := recoveryMetrics(plan, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, nil, nil, 0, d, 0)
	if m.UsefulEvidenceRate == nil || *m.UsefulEvidenceRate < 0 || *m.UsefulEvidenceRate > 1 {
		t.Fatalf("rate %#v", m.UsefulEvidenceRate)
	}
}
