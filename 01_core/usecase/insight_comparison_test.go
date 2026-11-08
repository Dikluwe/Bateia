/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_comparison_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func comparisonFixture() domain.InsightComparisonRequest {
	source := sourceInsightFixture()
	baseline := &domain.InsightTargetBaseline{Repository: "org/target", Files: []domain.InsightBaselineFile{
		{Path: "target.go", Content: "package target\nfunc Reserve() bool { return false }\n"},
	}}
	sealBaselineFixture(baseline)
	for _, kind := range []domain.InsightObservationKind{domain.InsightStrategy, domain.InsightPrecondition, domain.InsightTradeoff} {
		baseline.Observations = append(baseline.Observations, domain.SourceInsightObservation{MechanismClass: source.MechanismClasses[0], Kind: kind, Statement: "O alvo recusa a admissão neste trecho", References: []domain.InsightObservationReference{{Path: "target.go", ContentHash: baseline.Files[0].ContentHash, StartLine: 2, EndLine: 2}}})
	}
	return domain.InsightComparisonRequest{SchemaVersion: "1", Investigation: source, TargetBaseline: baseline}
}

func sealBaselineFixture(b *domain.InsightTargetBaseline) {
	files := append([]domain.InsightBaselineFile(nil), b.Files...)
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	parts := []string{}
	for i := range b.Files {
		b.Files[i].ContentHash = "sha256:" + contentHash(b.Files[i].Content)
	}
	for _, file := range files {
		path, _ := json.Marshal(file.Path)
		parts = append(parts, fmt.Sprintf(`{"path":%s,"content_hash":"sha256:%s"}`, path, contentHash(file.Content)))
	}
	b.SnapshotSHA256 = "sha256:" + contentHash("["+strings.Join(parts, ",")+"]")
}

func TestInsightComparisonLocatesExistingTargetWithoutInferringBenefit(t *testing.T) {
	r := comparisonFixture()
	before, _ := json.Marshal(r)
	report, err := CompareInsightBaseline(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	base, err := InvestigateSourceInsights(context.Background(), r.Investigation)
	if err != nil || report.SourceReportID != base.ID || !reflect.DeepEqual(report.SourceAnalyses, base.SourceAnalyses) {
		t.Fatal("source investigation rewritten")
	}
	if report.Status != domain.ImprovementPartial || report.EffectState != "UNVERIFIED" || report.AdoptionDisposition != "NOT_AN_ADOPTION_DECISION" || len(report.Comparisons) != 5 {
		t.Fatalf("availability promoted to benefit: %+v", report)
	}
	baseline := report.TargetBaseline
	if baseline.Scope != "SUPPLIED_FILES_ONLY" || baseline.Origin != "DECLARED_NOT_REASSESSED" || baseline.SnapshotSHA256 != r.TargetBaseline.SnapshotSHA256 {
		t.Fatal("baseline presented with unsupported scope")
	}
	refs := map[string]bool{}
	for _, observation := range baseline.Observations {
		refs[observation.ID] = true
		if !strings.HasPrefix(observation.ID, "target_obs_") || observation.References[0].Excerpt != "func Reserve() bool { return false }" {
			t.Fatal("unlocated target declaration")
		}
	}
	for _, row := range report.Comparisons {
		for _, ref := range row.TargetObservationRefs {
			if !refs[ref] {
				t.Fatal("dangling target reference")
			}
		}
		if !containsText(row.Gaps, "TARGET_EFFECT_NOT_VERIFIED") || !strings.Contains(row.NextStep.SafeProcedure, "política já declarada") {
			t.Fatal("comparison bypasses existing behavior")
		}
	}
	after, _ := json.Marshal(r)
	if string(before) != string(after) {
		t.Fatal("caller input mutated")
	}
}

func TestInsightComparisonUnknownAndMismatchedDeclarationsStayScoped(t *testing.T) {
	for _, mode := range []string{"nil", "empty", "source-empty", "no-assessments", "other-class", "other-source"} {
		t.Run(mode, func(t *testing.T) {
			r := comparisonFixture()
			switch mode {
			case "nil":
				r.TargetBaseline = nil
			case "empty":
				r.TargetBaseline.Observations = nil
			case "source-empty":
				r.Investigation.Sources[0].Observations = nil
			case "no-assessments":
				r.Investigation.Sources[0].Assessments = nil
			case "other-class":
				r.Investigation.MechanismClasses = append(r.Investigation.MechanismClasses, "APPEND_ONLY_LOG_RECOVERY")
				for i := range r.TargetBaseline.Observations {
					r.TargetBaseline.Observations[i].MechanismClass = "APPEND_ONLY_LOG_RECOVERY"
				}
			case "other-source":
				source := renamedSourceFixture("other")
				source.Observations = nil
				r.Investigation.Sources = append(r.Investigation.Sources, source)
			}
			got, err := CompareInsightBaseline(context.Background(), r)
			if err != nil || got.EffectState != "UNVERIFIED" {
				t.Fatalf("missing declarations lost: %v", err)
			}
			for _, row := range got.Comparisons {
				if (mode == "other-class" || mode == "nil" || mode == "empty" || mode == "source-empty") && row.DeclarationStatus == "BOTH_DECLARED" {
					t.Fatal("declarations crossed class or missing baseline")
				}
				if mode == "other-source" && row.Source.Repository == "other/queue-library" && len(row.SourceObservationRefs) != 0 {
					t.Fatal("observations crossed source")
				}
			}
			if mode == "no-assessments" && (got.SourceAnalyses[0].MechanismStatus != domain.MechanismInconclusive || got.Status != domain.ImprovementPartial) {
				t.Fatal("availability confused with classification")
			}
			if (mode == "nil" || mode == "empty" || mode == "source-empty" || mode == "other-class") && got.Status != domain.ImprovementInconclusive {
				t.Fatal("unpaired declarations promoted")
			}
		})
	}
}

func TestInsightComparisonCanonicalIdentityAndCallerIsolation(t *testing.T) {
	r := comparisonFixture()
	r.TargetBaseline.Files = append(r.TargetBaseline.Files, domain.InsightBaselineFile{Path: "a.go", Content: "package a\n"})
	sealBaselineFixture(r.TargetBaseline)
	first, err := CompareInsightBaseline(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	r.TargetBaseline.Files[0], r.TargetBaseline.Files[1] = r.TargetBaseline.Files[1], r.TargetBaseline.Files[0]
	r.TargetBaseline.Observations[0], r.TargetBaseline.Observations[2] = r.TargetBaseline.Observations[2], r.TargetBaseline.Observations[0]
	r.TargetBaseline.SnapshotSHA256 = strings.TrimPrefix(r.TargetBaseline.SnapshotSHA256, "sha256:")
	for i := range r.TargetBaseline.Files {
		r.TargetBaseline.Files[i].ContentHash = strings.TrimPrefix(r.TargetBaseline.Files[i].ContentHash, "sha256:")
	}
	for i := range r.TargetBaseline.Observations {
		r.TargetBaseline.Observations[i].References[0].ContentHash = strings.TrimPrefix(r.TargetBaseline.Observations[i].References[0].ContentHash, "sha256:")
	}
	before, _ := json.Marshal(r)
	again, err := CompareInsightBaseline(context.Background(), r)
	after, _ := json.Marshal(r)
	if err != nil || !reflect.DeepEqual(first, again) || string(before) != string(after) {
		t.Fatal("ordering, hash prefix or normalization mutated identity/input")
	}
	for _, mutate := range []func(*domain.InsightComparisonRequest){
		func(r *domain.InsightComparisonRequest) { r.TargetBaseline.Repository = "org/another" },
		func(r *domain.InsightComparisonRequest) {
			r.TargetBaseline.Observations[0].Statement += " outra interpretação"
		},
		func(r *domain.InsightComparisonRequest) {
			r.TargetBaseline.Files = append(r.TargetBaseline.Files, domain.InsightBaselineFile{Path: "extra.go", Content: "package extra\n"})
			sealBaselineFixture(r.TargetBaseline)
		},
	} {
		r := comparisonFixture()
		original, _ := CompareInsightBaseline(context.Background(), r)
		mutate(&r)
		changed, err := CompareInsightBaseline(context.Background(), r)
		if err != nil || changed.ID == original.ID || changed.Comparisons[0].ID == original.Comparisons[0].ID {
			t.Fatalf("target input not causal: %v", err)
		}
	}
}
