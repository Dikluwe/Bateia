/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_screening.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func benchmarkPath(name string) string { return filepath.Join("testdata", "discovery_benchmark", name) }

func TestRunIdeaScreenLoadsModelFromCalibrationArtifact(t *testing.T) {
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{"screen", "--snapshot", benchmarkPath("structural-impact.lexical.snapshot.json"), "--model", benchmarkPath("screening.calibration.json")}, &stdout)
	if err != nil {
		t.Fatal(err)
	}
	var report domain.ScreeningReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Decisions) == 0 || report.Decisions[0].CalibrationStatus != "calibrated" || report.Decisions[0].Probability == nil {
		t.Fatalf("calibrador não carregado: %+v", report.Decisions)
	}
}

func TestRunIdeaScreeningCalibrateUsesIdeaGroups(t *testing.T) {
	directory := benchmarkPath("")
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{"screening", "calibrate", "--snapshots", directory, "--judgments", directory, "--strategy", "lexical"}, &stdout)
	if err != nil {
		t.Fatal(err)
	}
	var report domain.CalibrationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Folds) != 4 || report.ObservationCount != 57 {
		t.Fatalf("calibração inesperada: %+v", report)
	}
}

func TestRunIdeaScreeningAuditValidatesRevisions(t *testing.T) {
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{"screening", "audit", "--revisions", benchmarkPath("screening.audit-revisions.json")}, &stdout)
	if err != nil {
		t.Fatal(err)
	}
	var report domain.AuditReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.ReviewedCount != 33 || report.FalseNegativeCount != 0 {
		t.Fatalf("auditoria inesperada: %+v", report)
	}
}
