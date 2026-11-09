/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/calibration_group_admission_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestCalibrationRejectsMissingAndRepeatedIdeaIdentity(t *testing.T) {
	for _, ids := range [][]string{{"same", "same"}, {"same", "other", "same"}, {"", "other"}, {"  ", "other"}, {" one", "two"}, {"one\t", "two"}, {"one"}} {
		t.Run(fmt.Sprint(ids), func(t *testing.T) {
			groups := []domain.ScreeningDatasetGroup{}
			for _, id := range ids {
				groups = append(groups, calibrationGroup(id, true))
			}
			before, _ := json.Marshal(groups)
			got, err := CalibrateScreening(groups, DefaultScreeningModel())
			if err == nil || !reflect.DeepEqual(got, domain.CalibrationReport{}) {
				t.Errorf("invalid idea grouping admitted: %v %+v", err, got)
			}
			after, _ := json.Marshal(groups)
			if string(before) != string(after) {
				t.Fatal("groups mutated")
			}
		})
	}
}

func TestCalibrationDistinctIdeasKeepFrozenReport(t *testing.T) {
	groups := []domain.ScreeningDatasetGroup{calibrationGroup("one", true), calibrationGroup("two", true), calibrationGroup("three", false)}
	before, _ := json.Marshal(groups)
	got, err := CalibrateScreening(groups, DefaultScreeningModel())
	if err != nil || len(got.Folds) != 3 || got.ObservationCount != 6 {
		t.Fatal("valid groups changed", err, got)
	}
	for i, fold := range got.Folds {
		if fold.HeldOutIdea != groups[i].IdeaID || fold.TrainingCount != 4 || fold.TestCount != 2 {
			t.Fatal("fold or uncertain handling changed", fold)
		}
	}
	raw, _ := json.Marshal(got)
	if fmt.Sprintf("%x", sha256.Sum256(raw)) != "9758913ab5c81eeecef7bd1010aca097d3a16f2845193d2b20b6277baeed2551" {
		t.Fatal("valid calibration report changed")
	}
	after, _ := json.Marshal(groups)
	if string(before) != string(after) {
		t.Fatal("valid input mutated")
	}
}
