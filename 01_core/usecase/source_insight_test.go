/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/source_insight_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func sourceInsightFixture() domain.InsightInvestigationRequest {
	r := insightFixture()
	r.SchemaVersion = "2"
	for _, entry := range []struct {
		kind domain.InsightObservationKind
		text string
	}{
		{domain.InsightStrategy, "Reserve retorna false para capacidade não positiva"},
		{domain.InsightPrecondition, "Configurar capacidade positiva para admitir"},
		{domain.InsightTradeoff, "O chamador precisa tratar a recusa"},
		{domain.InsightLimitation, "O trecho não observa o estado do alvo"},
	} {
		r.Sources[0].Observations = append(r.Sources[0].Observations, domain.SourceInsightObservation{
			MechanismClass: r.MechanismClasses[0], Kind: entry.kind, Statement: entry.text,
			References: []domain.InsightObservationReference{{Path: "queue.go", ContentHash: r.Sources[0].Bundle.Items[0].ContentHash, StartLine: 2, EndLine: 2}},
		})
	}
	return r
}

func TestSourceInsightCarriesLocatedConditionsWithoutPromotingJudgments(t *testing.T) {
	request := sourceInsightFixture()
	before, _ := json.Marshal(request)
	got, err := InvestigateSourceInsights(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if got.SchemaVersion != sourceInsightContract || got.Status != domain.ImprovementPartial || len(got.SourceAnalyses) != 1 || len(got.Hypotheses) != 1 {
		t.Fatalf("unexpected report: %+v", got)
	}
	analysis, hypothesis := got.SourceAnalyses[0], got.Hypotheses[0]
	if analysis.CatalogDefinition.Origin != "CATALOG_TEMPLATE" || !reflect.DeepEqual(analysis.ExtractionScope.Languages, []string{"go", "java", "python", "rust", "typescript"}) {
		t.Fatal("catalog definition and extraction scope not separated")
	}
	if analysis.AssessmentOrigin != "DECLARED_NOT_REASSESSED" || analysis.MechanismStatus != domain.MechanismProven || hypothesis.Status != domain.HypothesisProposed || hypothesis.TargetState != domain.TargetUnknown || hypothesis.AdoptionDisposition != "NOT_AN_ADOPTION_DECISION" {
		t.Fatal("declared interpretation promoted")
	}
	if !reflect.DeepEqual(hypothesis.Preconditions, []string{"Configurar capacidade positiva para admitir"}) || !reflect.DeepEqual(hypothesis.Tradeoffs, []string{"O chamador precisa tratar a recusa"}) {
		t.Fatalf("source conditions lost or replaced by languages: %+v", hypothesis)
	}
	for _, observation := range analysis.Observations {
		if observation.Origin != "DECLARED_NOT_REASSESSED" || observation.ID == "" || observation.References[0].Excerpt != "func Reserve(capacity int) bool { return capacity > 0 }" {
			t.Fatalf("unbound or promoted observation: %+v", observation)
		}
		if !containsText(hypothesis.EvidenceRefs, observation.ID) {
			t.Fatal("hypothesis lost observation provenance")
		}
	}
	after, _ := json.Marshal(request)
	if !bytesEqual(before, after) {
		t.Fatal("input mutated")
	}
	baseRequest := request
	baseRequest.SchemaVersion = "1"
	baseRequest.Sources = append([]domain.InsightEvidenceInput(nil), request.Sources...)
	baseRequest.Sources[0].Observations = nil
	base, err := InvestigateInsights(context.Background(), baseRequest)
	if err != nil || base.Mechanisms[0].ID != analysis.MechanismRef || base.ID == got.ID || base.Hypotheses[0].ID == hypothesis.ID {
		t.Fatal("classification changed or v2 reused a legacy identity")
	}
}

func TestSourceInsightRejectsInvalidObservationsAndRanges(t *testing.T) {
	mutations := map[string]func(*domain.InsightInvestigationRequest){
		"schema":        func(r *domain.InsightInvestigationRequest) { r.SchemaVersion = "1" },
		"unknown-kind":  func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].Kind = "CONFIRMED" },
		"unknown-class": func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].MechanismClass = "INVENTED" },
		"unrequested-class": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].MechanismClass = "APPEND_ONLY_LOG_RECOVERY"
		},
		"blank-statement":   func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].Statement = "  " },
		"statement-control": func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].Statement = "bad\ntext" },
		"long-statement": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].Statement = strings.Repeat("a", 1001)
		},
		"no-references": func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].References = nil },
		"many-references": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].References = make([]domain.InsightObservationReference, 9)
		},
		"path": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].References[0].Path = "../private"
		},
		"hash": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].References[0].ContentHash = strings.Repeat("0", 64)
		},
		"missing-hash": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].References[0].ContentHash = ""
		},
		"zero-line":             func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].References[0].StartLine = 0 },
		"reversed-lines":        func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].References[0].EndLine = 1 },
		"trailing-phantom-line": func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].References[0].EndLine = 3 },
		"duplicate-reference": func(r *domain.InsightInvestigationRequest) {
			o := &r.Sources[0].Observations[0]
			o.References = append(o.References, o.References[0])
		},
		"duplicate-observation": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations = append(r.Sources[0].Observations, r.Sources[0].Observations[0])
		},
		"source-budget": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations = make([]domain.SourceInsightObservation, 33)
		},
		"tampered-bundle": func(r *domain.InsightInvestigationRequest) { r.Sources[0].Bundle.Items[0].Content += "tampered" },
	}
	for name, change := range mutations {
		t.Run(name, func(t *testing.T) {
			r := sourceInsightFixture()
			change(&r)
			got, err := InvestigateSourceInsights(context.Background(), r)
			if err == nil || got.ID != "" {
				t.Fatalf("invalid observation accepted: %+v %v", got, err)
			}
		})
	}
	r := sourceInsightFixture()
	r.SchemaVersion = "1"
	if _, err := InvestigateInsights(context.Background(), r); err == nil {
		t.Fatal("v1 accepted observations")
	}
	r.SchemaVersion = "2"
	if _, err := InvestigateInsights(context.Background(), r); err == nil {
		t.Fatal("v1 facade accepted v2")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := InvestigateSourceInsights(ctx, r); !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation ignored")
	}
	if _, err := InvestigateSourceInsights(nil, r); !errors.Is(err, domain.ErrInvalidRunInput) {
		t.Fatal("nil context accepted")
	}
}

func TestSourceInsightDeterminismAndCausality(t *testing.T) {
	r := sourceInsightFixture()
	r.Sources[0].Observations[0].References = append(r.Sources[0].Observations[0].References, domain.InsightObservationReference{Path: "queue_test.go", ContentHash: strings.TrimPrefix(r.Sources[0].Bundle.Items[1].ContentHash, "sha256:"), StartLine: 2, EndLine: 2})
	first, err := InvestigateSourceInsights(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	for i, j := 0, len(r.Sources[0].Observations)-1; i < j; i, j = i+1, j-1 {
		r.Sources[0].Observations[i], r.Sources[0].Observations[j] = r.Sources[0].Observations[j], r.Sources[0].Observations[i]
	}
	for i := range r.Sources[0].Observations {
		refs := r.Sources[0].Observations[i].References
		for a, b := 0, len(refs)-1; a < b; a, b = a+1, b-1 {
			refs[a], refs[b] = refs[b], refs[a]
		}
		for j := range refs {
			refs[j].ContentHash = strings.TrimPrefix(refs[j].ContentHash, "sha256:")
		}
	}
	again, err := InvestigateSourceInsights(context.Background(), r)
	if err != nil || !reflect.DeepEqual(first, again) {
		t.Fatal("ordering or optional hash prefix changed result")
	}
	for _, change := range []func(*domain.InsightInvestigationRequest){
		func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].Statement += " outra condição"
		},
		func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].Kind = domain.InsightFailureMode
		},
		func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].References[0].StartLine = 1 },
	} {
		input := sourceInsightFixture()
		original, _ := InvestigateSourceInsights(context.Background(), input)
		change(&input)
		changed, err := InvestigateSourceInsights(context.Background(), input)
		if err != nil || changed.ID == original.ID || changed.Hypotheses[0].ID == original.Hypotheses[0].ID {
			t.Fatalf("observation is not causal: %v", err)
		}
	}
}

func renamedSourceFixture(owner string) domain.InsightEvidenceInput {
	source := sourceInsightFixture().Sources[0]
	source.Bundle.Repository.Owner = owner
	for i := range source.Bundle.Items {
		source.Bundle.Items[i].Ref.Owner = owner
	}
	for i := range source.Assessments {
		source.Assessments[i].Repository.Owner = owner
		for j := range source.Assessments[i].Capabilities {
			for k := range source.Assessments[i].Capabilities[j].Sources {
				source.Assessments[i].Capabilities[j].Sources[k].Ref.Owner = owner
			}
		}
	}
	return source
}

func TestSourceInsightMultipleSourcesDoNotMixConditions(t *testing.T) {
	r := sourceInsightFixture()
	second := renamedSourceFixture("other")
	second.Observations = second.Observations[:1]
	second.Observations[0].Kind = domain.InsightPrecondition
	second.Observations[0].Statement = "Condição exclusiva da outra fonte"
	r.Sources = append(r.Sources, second)
	first, err := InvestigateSourceInsights(context.Background(), r)
	if err != nil || len(first.Hypotheses) != 2 || len(first.SourceAnalyses) != 2 {
		t.Fatalf("multiple sources not retained: %v", err)
	}
	for _, hypothesis := range first.Hypotheses {
		if len(hypothesis.Preconditions) != 1 {
			t.Fatalf("source conditions mixed: %+v", hypothesis)
		}
	}
	r.Sources[0], r.Sources[1] = r.Sources[1], r.Sources[0]
	again, err := InvestigateSourceInsights(context.Background(), r)
	if err != nil || !reflect.DeepEqual(first, again) {
		t.Fatal("source order changed result")
	}
}

func TestSourceInsightMissingAndLimitedEvidenceRemainsExplicit(t *testing.T) {
	for _, condition := range []string{"no-observations", "no-assessments", "partial", "conflict", "incomplete", "other-class"} {
		t.Run(condition, func(t *testing.T) {
			r := sourceInsightFixture()
			switch condition {
			case "no-observations":
				r.Sources[0].Observations = nil
			case "no-assessments":
				r.Sources[0].Assessments = nil
			case "partial":
				r.Sources[0].Assessments[0].Capabilities[0].Status = domain.EvidencePartial
			case "conflict":
				witness := r.Sources[0].Assessments[0].Capabilities[0].Sources[0]
				witness.RuleID = "unbounded-queue"
				r.Sources[0].Assessments[0].Capabilities = append(r.Sources[0].Assessments[0].Capabilities, domain.CapabilityEvidence{
					Capability: "conflicting declaration", Status: domain.EvidenceConflict, Sources: []domain.EvidenceSource{witness},
				})
			case "incomplete":
				r.Sources[0].Bundle.Complete = false
			case "other-class":
				r.MechanismClasses = append(r.MechanismClasses, "APPEND_ONLY_LOG_RECOVERY")
				r.Sources[0].Observations[0].MechanismClass = "APPEND_ONLY_LOG_RECOVERY"
			}
			got, err := InvestigateSourceInsights(context.Background(), r)
			if err != nil || len(got.SourceAnalyses) == 0 {
				t.Fatalf("observations lost: %v", err)
			}
			if condition == "no-assessments" && (len(got.Hypotheses) != 0 || got.SourceAnalyses[0].MechanismStatus != domain.MechanismInconclusive || len(got.SourceAnalyses[0].Observations) != 4) {
				t.Fatal("observations alone promoted classification")
			}
			if condition == "no-observations" && (len(got.Hypotheses[0].Preconditions) != 0 || !containsText(got.SourceAnalyses[0].Gaps, "OPERATIONAL_PRECONDITIONS_NOT_DECLARED")) {
				t.Fatal("unknown conditions replaced by languages")
			}
			if condition == "conflict" && (got.SourceAnalyses[0].MechanismStatus != domain.MechanismInconclusive || len(got.SourceAnalyses[0].CounterEvidence) == 0) {
				t.Fatal("conflicting classification promoted or counterevidence lost")
			}
			for _, hypothesis := range got.Hypotheses {
				if hypothesis.TargetState != domain.TargetUnknown || hypothesis.Status == domain.HypothesisSupported || ((condition == "partial" || condition == "conflict" || condition == "incomplete") && hypothesis.Status != domain.HypothesisInconclusive) {
					t.Fatal("source or target promoted")
				}
			}
			if condition == "other-class" {
				for _, analysis := range got.SourceAnalyses {
					if analysis.CatalogDefinition.MechanismClass == "APPEND_ONLY_LOG_RECOVERY" && (analysis.MechanismRef != "" || analysis.MechanismStatus != domain.MechanismInconclusive || len(analysis.Observations) != 1) {
						t.Fatal("observation crossed mechanism classes")
					}
				}
			}
		})
	}
}

func TestSourceInsightObservationCountLimits(t *testing.T) {
	for _, count := range []int{32, 33, 128, 129} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			r := sourceInsightFixture()
			r.Sources = nil
			remaining := count
			for remaining > 0 {
				source := renamedSourceFixture(fmt.Sprintf("owner%d", len(r.Sources)))
				base := source.Observations[0]
				source.Observations = nil
				n := remaining
				if count != 33 && n > 32 {
					n = 32
				}
				for j := 0; j < n; j++ {
					observation := base
					observation.Statement = fmt.Sprintf("Declaração %d", j)
					source.Observations = append(source.Observations, observation)
				}
				r.Sources = append(r.Sources, source)
				remaining -= n
			}
			_, err := InvestigateSourceInsights(context.Background(), r)
			if (err != nil) != (count == 33 || count == 129) {
				t.Fatalf("count=%d err=%v", count, err)
			}
		})
	}
}

func extendedSourceInsightFixture(extra string) domain.InsightInvestigationRequest {
	r := sourceInsightFixture()
	item := &r.Sources[0].Bundle.Items[0]
	item.Content += extra
	item.ContentHash = "sha256:" + contentHash(item.Content)
	r.Sources[0].Assessments[0].Capabilities[0].Sources[0].ContentHash = item.ContentHash
	for i := range r.Sources[0].Observations {
		r.Sources[0].Observations[i].References[0].ContentHash = item.ContentHash
	}
	return r
}

func TestSourceInsightExcerptAndReferenceBudgetBoundaries(t *testing.T) {
	for _, size := range []int{4096, 4097} {
		r := extendedSourceInsightFixture(strings.Repeat("a", size) + "\n")
		r.Sources[0].Observations = r.Sources[0].Observations[:1]
		r.Sources[0].Observations[0].References[0].StartLine = 3
		r.Sources[0].Observations[0].References[0].EndLine = 3
		_, err := InvestigateSourceInsights(context.Background(), r)
		if (err != nil) != (size > 4096) {
			t.Fatalf("excerpt size %d: %v", size, err)
		}
	}
	for _, lines := range []int{20, 21} {
		r := extendedSourceInsightFixture(strings.Repeat("// line\n", 25))
		r.Sources[0].Observations = r.Sources[0].Observations[:1]
		r.Sources[0].Observations[0].References[0].StartLine = 2
		r.Sources[0].Observations[0].References[0].EndLine = 2 + lines - 1
		_, err := InvestigateSourceInsights(context.Background(), r)
		if (err != nil) != (lines > 20) {
			t.Fatalf("excerpt lines %d: %v", lines, err)
		}
	}
	r := extendedSourceInsightFixture(strings.Repeat("// line\n", 10))
	base := r.Sources[0].Observations[0]
	base.References = nil
	for line := 2; line < 10; line++ {
		base.References = append(base.References, domain.InsightObservationReference{Path: "queue.go", ContentHash: r.Sources[0].Bundle.Items[0].ContentHash, StartLine: line, EndLine: line})
	}
	r.Sources[0].Observations = nil
	for i := 0; i < 32; i++ {
		observation := base
		observation.Statement = fmt.Sprintf("Declaração %d", i)
		r.Sources[0].Observations = append(r.Sources[0].Observations, observation)
	}
	if _, err := InvestigateSourceInsights(context.Background(), r); err != nil {
		t.Fatalf("256 references rejected: %v", err)
	}
	other := renamedSourceFixture("other")
	other.Observations = other.Observations[:1]
	r.Sources = append(r.Sources, other)
	if _, err := InvestigateSourceInsights(context.Background(), r); err == nil {
		t.Fatal("257 references accepted")
	}
}

func TestSourceInsightContradictoryDeclarationsRemainDeclared(t *testing.T) {
	r := sourceInsightFixture()
	contradiction := r.Sources[0].Observations[0]
	contradiction.Statement = "A fonte garante admissão ilimitada em qualquer condição"
	r.Sources[0].Observations = append(r.Sources[0].Observations, contradiction)
	got, err := InvestigateSourceInsights(context.Background(), r)
	if err != nil || len(got.SourceAnalyses[0].Observations) != 5 || len(got.SourceAnalyses[0].CounterEvidence) != 0 {
		t.Fatal("declaration interpreted as formal evidence or discarded")
	}
	for _, observation := range got.SourceAnalyses[0].Observations {
		if observation.Origin != "DECLARED_NOT_REASSESSED" {
			t.Fatal("interpretation promoted")
		}
	}
}

func containsText(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func bytesEqual(left, right []byte) bool { return string(left) == string(right) }
