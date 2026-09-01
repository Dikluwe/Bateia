/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_calibration_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import (
	"bateia/01_core/domain"
	"reflect"
	"testing"
)

func evidenceCalibrationObservation(id, split string, expected, observed domain.EvidenceStatus) domain.EvidenceCalibrationObservation {
	source := domain.EvidenceSource{Kind: domain.EvidenceSymbol, Ref: domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}, Path: "feature.go", Symbol: "Feature", RuleID: "rule", Strength: .9}
	item := domain.CollectedEvidence{Kind: domain.EvidenceFile, Ref: source.Ref, Path: source.Path, Content: "type Feature struct{}"}
	return domain.EvidenceCalibrationObservation{ID: id, Split: split, Expected: expected, Observed: domain.CapabilityEvidence{Status: observed, Sources: []domain.EvidenceSource{source}}, Evidence: []domain.CollectedEvidence{item}}
}

func TestEvidenceCalibrationPromotesPassingThreeSplitsDeterministically(t *testing.T) {
	observations := []domain.EvidenceCalibrationObservation{
		evidenceCalibrationObservation("c", "TEST", domain.EvidenceProven, domain.EvidenceProven),
		evidenceCalibrationObservation("a", "CONTRACT", domain.EvidenceProven, domain.EvidenceProven),
		evidenceCalibrationObservation("b", "VALIDATION", domain.EvidenceProven, domain.EvidenceProven),
	}
	left := EvaluateEvidenceCalibration(observations)
	right := EvaluateEvidenceCalibration([]domain.EvidenceCalibrationObservation{observations[1], observations[2], observations[0]})
	if left.Verdict != "PROMOTE_EVIDENCE_CATALOG" || !reflect.DeepEqual(left, right) {
		t.Fatalf("calibração não promoveu deterministicamente: %#v %#v", left, right)
	}
}

func TestEvidenceCalibrationRejectsOpaqueOrMutableProof(t *testing.T) {
	opaque := evidenceCalibrationObservation("opaque", "TEST", domain.EvidenceInconclusive, domain.EvidenceProven)
	opaque.Observed.Sources[0].Ref.Commit = "HEAD"
	report := EvaluateEvidenceCalibration([]domain.EvidenceCalibrationObservation{opaque})
	if report.Verdict != "FAIL_GENERALIZATION" || report.Aggregate.OpaquePromotions != 1 || report.Aggregate.InvalidSources != 1 {
		t.Fatalf("prova opaca aceita: %#v", report)
	}
}

func TestEvidenceCalibrationRequiresDiagnosisForMiss(t *testing.T) {
	miss := evidenceCalibrationObservation("miss", "VALIDATION", domain.EvidenceProven, domain.EvidencePartial)
	if got := EvaluateEvidenceCalibration([]domain.EvidenceCalibrationObservation{miss}); got.Verdict != "KEEP_INCONCLUSIVE_AND_DIAGNOSE" {
		t.Fatalf("perda sem diagnóstico aceita: %#v", got)
	}
	miss.Diagnosis = domain.MissPattern
	if got := EvaluateEvidenceCalibration([]domain.EvidenceCalibrationObservation{miss}); got.Aggregate.UndiagnosedMisses != 0 {
		t.Fatalf("diagnóstico não reconhecido: %#v", got)
	}
}
