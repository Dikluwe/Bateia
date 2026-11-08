/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/screening_model_policy.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"math"
)

func IsValidOptInScreeningModel(model domain.ScreeningModel) bool {
	if model.SchemaVersion != "1" || model.CalibrationStatus != "uncalibrated" || model.Calibration != nil || model.Weights == nil {
		return false
	}
	if !finiteOptInModelNumber(model.DropThreshold) || !finiteOptInModelNumber(model.AdvanceThreshold) ||
		model.DropThreshold < 0 || model.DropThreshold >= model.AdvanceThreshold || model.AdvanceThreshold > 1 {
		return false
	}
	if !finiteOptInModelNumber(model.AuditSampleRate) || model.AuditSampleRate < 0 || model.AuditSampleRate > 1 {
		return false
	}
	for _, weight := range model.Weights {
		if !finiteOptInModelNumber(weight) {
			return false
		}
	}
	for _, prior := range model.ExpansionPriors {
		if !finiteOptInModelNumber(prior) {
			return false
		}
	}
	return true
}

func finiteOptInModelNumber(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
