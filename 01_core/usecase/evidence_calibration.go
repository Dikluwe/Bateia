/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_calibration.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import (
	"bateia/01_core/domain"
	"sort"
	"strings"
)

func EvaluateEvidenceCalibration(observations []domain.EvidenceCalibrationObservation) domain.EvidenceCalibrationReport {
	report := domain.EvidenceCalibrationReport{SchemaVersion: "1", Verdict: "PROMOTE_EVIDENCE_CATALOG", BySplit: map[string]domain.EvidenceCalibrationMetrics{}}
	ordered := append([]domain.EvidenceCalibrationObservation(nil), observations...)
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Split == ordered[j].Split {
			return ordered[i].ID < ordered[j].ID
		}
		return ordered[i].Split < ordered[j].Split
	})
	for _, observation := range ordered {
		metrics := report.BySplit[observation.Split]
		accumulateCalibration(&metrics, observation)
		report.BySplit[observation.Split] = metrics
		accumulateCalibration(&report.Aggregate, observation)
	}
	for split, metrics := range report.BySplit {
		metrics = finishCalibration(metrics)
		report.BySplit[split] = metrics
	}
	report.Aggregate = finishCalibration(report.Aggregate)
	if report.Aggregate.InvalidSources > 0 || report.Aggregate.OpaquePromotions > 0 {
		report.Verdict = "FAIL_FALSE_PROOF"
	} else if report.Aggregate.UndiagnosedMisses > 0 || rateBelow(report.Aggregate.ProvenRecall, .85) {
		report.Verdict = "KEEP_INCONCLUSIVE_AND_DIAGNOSE"
	} else if rateBelow(report.Aggregate.ProvenPrecision, .90) {
		report.Verdict = "FAIL_FALSE_PROOF"
	}
	if test, ok := report.BySplit["TEST"]; ok && (rateBelow(test.ProvenPrecision, .90) || rateBelow(test.ProvenRecall, .85) || test.InvalidSources > 0 || test.OpaquePromotions > 0 || test.UndiagnosedMisses > 0) {
		report.Verdict = "FAIL_GENERALIZATION"
	}
	return report
}

func accumulateCalibration(metrics *domain.EvidenceCalibrationMetrics, observation domain.EvidenceCalibrationObservation) {
	metrics.Cases++
	if observation.Expected == domain.EvidenceProven {
		metrics.ExpectedProven++
	}
	if observation.Observed.Status == domain.EvidenceProven {
		metrics.ObservedProven++
		if observation.Expected == domain.EvidenceProven {
			metrics.TrueProven++
		}
		if observation.Expected == domain.EvidenceInconclusive {
			metrics.OpaquePromotions++
		}
		if !validImmutableCalibrationSources(observation.Observed.Sources, observation.Evidence) {
			metrics.InvalidSources++
		}
	}
	if observation.Expected == domain.EvidenceProven && observation.Observed.Status != domain.EvidenceProven && observation.Diagnosis == "" {
		metrics.UndiagnosedMisses++
	}
}

func validImmutableCalibrationSources(sources []domain.EvidenceSource, evidence []domain.CollectedEvidence) bool {
	if len(sources) == 0 {
		return false
	}
	for _, source := range sources {
		if strings.TrimSpace(source.Path) == "" || strings.TrimSpace(source.Ref.Commit) == "" || strings.EqualFold(source.Ref.Commit, "HEAD") {
			return false
		}
		found := false
		for _, item := range evidence {
			if item.Ref == source.Ref && item.Path == source.Path && strings.TrimSpace(item.Content) != "" && (source.ContentHash == "" || source.ContentHash == item.ContentHash) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func finishCalibration(metrics domain.EvidenceCalibrationMetrics) domain.EvidenceCalibrationMetrics {
	if metrics.ObservedProven > 0 {
		value := float64(metrics.TrueProven) / float64(metrics.ObservedProven)
		metrics.ProvenPrecision = &value
	}
	if metrics.ExpectedProven > 0 {
		value := float64(metrics.TrueProven) / float64(metrics.ExpectedProven)
		metrics.ProvenRecall = &value
	}
	return metrics
}

func rateBelow(value *float64, threshold float64) bool {
	return value != nil && *value < threshold
}
