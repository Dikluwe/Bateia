/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_diagnosis_test.md
 * @layer L1
 * @updated 2026-08-26
 */

package usecase

import (
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func diagnosisFixture(status domain.EvidenceStatus, sources []domain.EvidenceSource, targets []domain.CollectionTargetResult) (domain.AssessmentReport, []domain.EvidenceBundle, []domain.EvidenceJudgment) {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	report := domain.AssessmentReport{SchemaVersion: "1", RulesVersion: "test", Assessments: []domain.CandidateAssessment{{Repository: ref, ProblemAlignment: domain.CapabilityEvidence{Capability: "analisar código sem executar o programa", Status: status, Sources: sources}, Capabilities: []domain.CapabilityEvidence{{Capability: "static analysis", Status: status, Sources: sources}}}}}
	bundles := []domain.EvidenceBundle{{SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true, Targets: targets}}
	judgments := []domain.EvidenceJudgment{{Repository: "org/repo", Capability: "analisar código sem executar o programa", Expected: domain.EvidenceProven}, {Repository: "org/repo", Capability: "static analysis", Expected: domain.EvidenceProven}}
	return report, bundles, judgments
}

func coveredTarget(status domain.CollectionTargetStatus) domain.CollectionTargetResult {
	return domain.CollectionTargetResult{Target: domain.ExtractionTarget{Path: "analyzer.go", Capabilities: []string{"static analysis", "analisar código sem executar o programa"}}, Status: status}
}

func TestDiagnosisDeduplicatesAliasesForRecall(t *testing.T) {
	report, bundles, judgments := diagnosisFixture(domain.EvidenceProven, nil, []domain.CollectionTargetResult{coveredTarget(domain.TargetCollected)})
	got := DiagnoseEvidenceCoverage(report, bundles, judgments, DefaultEvidenceRules())
	if got.Metrics.ExpectedProvenCount != 1 || got.Metrics.EmittedProvenCount != 1 || got.Metrics.ProvenRecall == nil || *got.Metrics.ProvenRecall != 1 {
		t.Fatalf("aliases inflaram métricas: %+v", got.Metrics)
	}
	if len(got.Diagnoses) != 0 {
		t.Fatalf("diagnósticos inesperados: %+v", got.Diagnoses)
	}
}

func TestDiagnosisOperationalCausesPrecedeInterpretation(t *testing.T) {
	tests := []struct {
		name    string
		targets []domain.CollectionTargetResult
		cause   domain.EvidenceMissCause
	}{
		{"not planned", nil, domain.MissTargetNotPlanned},
		{"not collected", []domain.CollectionTargetResult{coveredTarget(domain.TargetFailed)}, domain.MissTargetNotCollected},
		{"budget", []domain.CollectionTargetResult{coveredTarget(domain.TargetSkippedBudget)}, domain.MissBudgetExhausted},
		{"unsupported", []domain.CollectionTargetResult{coveredTarget(domain.TargetUnsupported)}, domain.MissLanguageUnsupported},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, bundles, judgments := diagnosisFixture(domain.EvidenceInconclusive, nil, tt.targets)
			got := DiagnoseEvidenceCoverage(report, bundles, judgments, DefaultEvidenceRules())
			if len(got.Diagnoses) != 1 || got.Diagnoses[0].PrimaryCause != tt.cause {
				t.Fatalf("diagnóstico = %+v", got.Diagnoses)
			}
		})
	}
}

func TestDiagnosisDistinguishesPatternLinkageStrengthAndIndependence(t *testing.T) {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	tests := []struct {
		name    string
		sources []domain.EvidenceSource
		cause   domain.EvidenceMissCause
	}{
		{"pattern", nil, domain.MissPattern},
		{"linkage", []domain.EvidenceSource{{Kind: domain.EvidenceFile, Ref: ref, Path: "analyzer.go", RuleID: "static-analysis", Strength: .2}}, domain.MissSymbolLinkage},
		{"strength", []domain.EvidenceSource{{Kind: domain.EvidenceSymbol, Ref: ref, Path: "analyzer.go", Symbol: "Analyzer", RuleID: "static-analysis", Strength: .4}}, domain.MissStrengthUnderestimated},
		{"independence", []domain.EvidenceSource{{Kind: domain.EvidenceDocumentation, Ref: ref, Path: "README.md", RuleID: "static-analysis", Strength: .4}}, domain.MissInsufficientIndependence},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, bundles, judgments := diagnosisFixture(domain.EvidencePartial, tt.sources, []domain.CollectionTargetResult{coveredTarget(domain.TargetCollected)})
			got := DiagnoseEvidenceCoverage(report, bundles, judgments, DefaultEvidenceRules())
			if got.Diagnoses[0].PrimaryCause != tt.cause {
				t.Fatalf("causa = %s; diagnóstico=%+v", got.Diagnoses[0].PrimaryCause, got.Diagnoses[0])
			}
		})
	}
}

func TestDiagnosisRecognizesNegativeConflict(t *testing.T) {
	report, bundles, judgments := diagnosisFixture(domain.EvidenceConflict, nil, []domain.CollectionTargetResult{coveredTarget(domain.TargetCollected)})
	got := DiagnoseEvidenceCoverage(report, bundles, judgments, DefaultEvidenceRules())
	if got.Diagnoses[0].PrimaryCause != domain.MissNegativeConflict {
		t.Fatalf("diagnóstico = %+v", got.Diagnoses[0])
	}
}

func TestDiagnosisIsOrderInvariantAndDoesNotMutateInputs(t *testing.T) {
	report, bundles, judgments := diagnosisFixture(domain.EvidencePartial, []domain.EvidenceSource{{Kind: domain.EvidenceDocumentation, Path: "README.md", Strength: .4}}, []domain.CollectionTargetResult{coveredTarget(domain.TargetCollected)})
	reportBefore, bundlesBefore, judgmentsBefore := report, append([]domain.EvidenceBundle(nil), bundles...), append([]domain.EvidenceJudgment(nil), judgments...)
	a := DiagnoseEvidenceCoverage(report, bundles, judgments, DefaultEvidenceRules())
	judgments[0], judgments[1] = judgments[1], judgments[0]
	b := DiagnoseEvidenceCoverage(report, bundles, judgments, DefaultEvidenceRules())
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("ordem alterou diagnóstico\n%+v\n%+v", a, b)
	}
	if !reflect.DeepEqual(report, reportBefore) || !reflect.DeepEqual(bundles, bundlesBefore) || !reflect.DeepEqual(judgmentsBefore[0], judgments[1]) {
		t.Fatal("entrada foi modificada")
	}
}
