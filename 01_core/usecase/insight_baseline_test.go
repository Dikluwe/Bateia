/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_baseline_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestInsightComparisonRejectsInvalidBaselineWithoutPartialReport(t *testing.T) {
	mutations := map[string]func(*domain.InsightComparisonRequest){
		"schema":           func(r *domain.InsightComparisonRequest) { r.SchemaVersion = "2" },
		"source-schema":    func(r *domain.InsightComparisonRequest) { r.Investigation.SchemaVersion = "1" },
		"repository":       func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Repository = "../org/repo" },
		"empty-repository": func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Repository = "/repo" },
		"snapshot":         func(r *domain.InsightComparisonRequest) { r.TargetBaseline.SnapshotSHA256 = strings.Repeat("0", 64) },
		"hash":             func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Files[0].ContentHash = "" },
		"content":          func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Files[0].Content += "altered" },
		"path":             func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Files[0].Path = "../private" },
		"duplicate-file": func(r *domain.InsightComparisonRequest) {
			r.TargetBaseline.Files = append(r.TargetBaseline.Files, r.TargetBaseline.Files[0])
		},
		"empty-files": func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Files = nil },
		"utf8":        func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Files[0].Content = string([]byte{255}) },
		"range":       func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Observations[0].References[0].EndLine = 3 },
		"foreign-reference": func(r *domain.InsightComparisonRequest) {
			r.TargetBaseline.Observations[0].References[0].Path = "queue.go"
		},
		"observation-duplicate": func(r *domain.InsightComparisonRequest) {
			r.TargetBaseline.Observations = append(r.TargetBaseline.Observations, r.TargetBaseline.Observations[0])
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			r := comparisonFixture()
			mutate(&r)
			got, err := CompareInsightBaseline(context.Background(), r)
			if err == nil || got.ID != "" {
				t.Fatal("invalid baseline accepted or partial report leaked")
			}
		})
	}
	r := comparisonFixture()
	if _, err := CompareInsightBaseline(nil, r); !errors.Is(err, domain.ErrInvalidRunInput) {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := CompareInsightBaseline(ctx, r); !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation ignored")
	}
}

func TestInsightComparisonBaselineContentBudgets(t *testing.T) {
	for _, test := range []struct {
		name         string
		files, bytes int
		valid        bool
	}{
		{"file-limit", 1, 256 << 10, true}, {"file-overflow", 1, (256 << 10) + 1, false},
		{"total-limit", 4, 256 << 10, true}, {"total-overflow", 5, 210 << 10, false},
		{"file-count-limit", 64, 1, true}, {"file-count-overflow", 65, 1, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := comparisonFixture()
			r.TargetBaseline.Observations = nil
			r.TargetBaseline.Files = nil
			for i := 0; i < test.files; i++ {
				r.TargetBaseline.Files = append(r.TargetBaseline.Files, domain.InsightBaselineFile{Path: fmt.Sprintf("file%02d.go", i), Content: strings.Repeat("a", test.bytes)})
			}
			sealBaselineFixture(r.TargetBaseline)
			_, err := CompareInsightBaseline(context.Background(), r)
			if (err == nil) != test.valid {
				t.Fatalf("budget validity %v: %v", test.valid, err)
			}
		})
	}
}
