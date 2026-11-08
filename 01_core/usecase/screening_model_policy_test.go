/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/screening_model_policy_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"math"
	"testing"
)

func TestOptInScreeningProfileKeepsItsExistingAdmissionRules(t *testing.T) {
	valid := func() domain.ScreeningModel {
		return domain.ScreeningModel{SchemaVersion: "1", CalibrationStatus: "uncalibrated", Weights: map[string]float64{"negative-is-allowed": -2}, ExpansionPriors: map[string]float64{}, DropThreshold: 0, AdvanceThreshold: 1, AuditSampleRate: 0}
	}
	cases := []struct {
		name   string
		mutate func(*domain.ScreeningModel)
		want   bool
	}{
		{"finite-negative-weight", func(*domain.ScreeningModel) {}, true},
		{"empty-present-weights", func(m *domain.ScreeningModel) { m.Weights = map[string]float64{} }, true},
		{"full-audit", func(m *domain.ScreeningModel) { m.AuditSampleRate = 1 }, true},
		{"schema", func(m *domain.ScreeningModel) { m.SchemaVersion = "2" }, false},
		{"calibrated", func(m *domain.ScreeningModel) { m.CalibrationStatus = "calibrated" }, false},
		{"opaque-status", func(m *domain.ScreeningModel) { m.CalibrationStatus = "" }, false},
		{"calibrator-present", func(m *domain.ScreeningModel) { m.Calibration = &domain.CalibrationParameters{} }, false},
		{"weights-absent", func(m *domain.ScreeningModel) { m.Weights = nil }, false},
		{"negative-drop", func(m *domain.ScreeningModel) { m.DropThreshold = -.01 }, false},
		{"equal-thresholds", func(m *domain.ScreeningModel) { m.DropThreshold = m.AdvanceThreshold }, false},
		{"advance-over-one", func(m *domain.ScreeningModel) { m.AdvanceThreshold = 1.01 }, false},
		{"negative-audit", func(m *domain.ScreeningModel) { m.AuditSampleRate = -.01 }, false},
		{"audit-over-one", func(m *domain.ScreeningModel) { m.AuditSampleRate = 1.01 }, false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			m := valid()
			test.mutate(&m)
			if got := usecase.IsValidOptInScreeningModel(m); got != test.want {
				t.Fatalf("admitted=%v want=%v", got, test.want)
			}
		})
	}
	for _, number := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		for _, field := range []string{"drop", "advance", "audit", "weight", "prior"} {
			m := valid()
			switch field {
			case "drop":
				m.DropThreshold = number
			case "advance":
				m.AdvanceThreshold = number
			case "audit":
				m.AuditSampleRate = number
			case "weight":
				m.Weights["x"] = number
			case "prior":
				m.ExpansionPriors["x"] = number
			}
			if usecase.IsValidOptInScreeningModel(m) {
				t.Errorf("nonfinite %s admitted", field)
			}
		}
	}
}
