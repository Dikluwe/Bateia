/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_investigation_test.md
 * @layer L1
 * @updated 2026-10-06
 */

package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func insightFixture() domain.InsightInvestigationRequest {
	ref := domain.CommitRef{Owner: "illustrative", Repo: "queue-library", Commit: strings.Repeat("b", 40)}
	implementation := "package queue\nfunc Reserve(capacity int) bool { return capacity > 0 }\n"
	independent := "package queue\nfunc TestCapacity() { if Reserve(0) { panic(\"admitted\") } }\n"
	items := []domain.CollectedEvidence{
		{Kind: domain.EvidenceFile, Ref: ref, Path: "queue.go", Content: implementation, ContentHash: "sha256:" + contentHash(implementation)},
		{Kind: domain.EvidenceTest, Ref: ref, Path: "queue_test.go", Content: independent, ContentHash: "sha256:" + contentHash(independent)},
	}
	sources := []domain.EvidenceSource{
		{Kind: domain.EvidenceFile, Ref: ref, Path: "queue.go", Line: 2, Snippet: "func Reserve(capacity int) bool { return capacity > 0 }", Symbol: "Reserve", RuleID: "bounded-queue-capacity", ContentHash: items[0].ContentHash},
		{Kind: domain.EvidenceTest, Ref: ref, Path: "queue_test.go", Line: 2, Snippet: "func TestCapacity() { if Reserve(0) { panic(\"admitted\") } }", Symbol: "TestCapacity", RuleID: "queue-capacity-test", ContentHash: items[1].ContentHash},
	}
	assessment := domain.CandidateAssessment{Repository: ref, Capabilities: []domain.CapabilityEvidence{{Capability: "bounded queue", Status: domain.EvidenceProven, Sources: sources}}, License: domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "GPL-3.0-only"}}
	return domain.InsightInvestigationRequest{SchemaVersion: "1", Question: "Como limitar trabalho admitido sem perder controle de cancelamento?", Target: domain.InsightTarget{Name: "Bateia", Context: "Execução de investigações com recursos limitados", Property: "trabalho admitido", DesiredEffect: "Evitar crescimento de trabalho pendente sob sobrecarga", Constraints: []string{"preservar cancelamento"}}, MechanismClasses: []string{"BOUNDED_QUEUE_BACKPRESSURE"}, Sources: []domain.InsightEvidenceInput{{Bundle: domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Items: items, Complete: true}, Assessments: []domain.CandidateAssessment{assessment}}}}
}

func TestInsightInvestigationIndependentOfReuseAndConservativeAboutTarget(t *testing.T) {
	request := insightFixture()
	before, _ := json.Marshal(request)
	result, err := InvestigateInsights(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Mechanisms) != 1 || result.Mechanisms[0].Status != domain.MechanismProven || len(result.Hypotheses) != 1 {
		t.Fatalf("source not investigated: %+v", result)
	}
	hypothesis := result.Hypotheses[0]
	if hypothesis.Status != domain.HypothesisProposed || hypothesis.TargetState != domain.TargetUnknown || hypothesis.TransferMode != domain.TransferReferenceOnly || hypothesis.AdoptionDisposition != "NOT_AN_ADOPTION_DECISION" || hypothesis.License != domain.ImprovementLicenseUnrecognized {
		t.Fatalf("target or adoption promoted: %+v", hypothesis)
	}
	if result.Status != domain.ImprovementPartial || len(hypothesis.ValidationPlan) != 3 || result.Provenance[0].AssessmentOrigin != "DECLARED_NOT_REASSESSED" {
		t.Fatal("missing validation limits")
	}
	after, _ := json.Marshal(request)
	if string(before) != string(after) {
		t.Fatal("input mutated")
	}
	again, err := InvestigateInsights(context.Background(), request)
	if err != nil || !reflect.DeepEqual(result, again) {
		t.Fatal("analysis is not deterministic")
	}
	request.Question = "Como reduzir a pressão de memória?"
	changed, err := InvestigateInsights(context.Background(), request)
	if err != nil || changed.ID == result.ID || changed.Hypotheses[0].ID == hypothesis.ID {
		t.Fatal("question not causal")
	}
	request = insightFixture()
	request.Target.Context = "Novo contexto"
	changed, err = InvestigateInsights(context.Background(), request)
	if err != nil || changed.ID == result.ID {
		t.Fatal("target context not causal")
	}
}

func TestInsightInvestigationPartialSourcesNeverValidateEffect(t *testing.T) {
	for _, condition := range []string{"incomplete", "documentation", "unknown-class-observation", "no-assessments"} {
		t.Run(condition, func(t *testing.T) {
			request := insightFixture()
			switch condition {
			case "incomplete":
				request.Sources[0].Bundle.Complete = false
			case "documentation":
				ref := request.Sources[0].Bundle.Repository
				text := "bounded queue with finite capacity"
				item := domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Ref: ref, Path: "README.md", Content: text, ContentHash: "sha256:" + contentHash(text)}
				source := domain.EvidenceSource{Kind: domain.EvidenceDocumentation, Ref: ref, Path: item.Path, Line: 1, Snippet: text, RuleID: "bounded-queue-architecture-documentation", ContentHash: item.ContentHash}
				request.Sources[0].Bundle.Items = []domain.CollectedEvidence{item}
				request.Sources[0].Assessments[0].Capabilities = []domain.CapabilityEvidence{{Capability: "bounded queue", Status: domain.EvidencePartial, Sources: []domain.EvidenceSource{source}}}
			case "unknown-class-observation":
				request.MechanismClasses = []string{"APPEND_ONLY_LOG_RECOVERY"}
			case "no-assessments":
				request.Sources[0].Assessments = nil
			}
			result, err := InvestigateInsights(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}
			if condition == "unknown-class-observation" || condition == "no-assessments" {
				if result.Status != domain.ImprovementInconclusive || len(result.Hypotheses) != 0 {
					t.Fatal("absence promoted")
				}
			} else if len(result.Hypotheses) != 1 || result.Hypotheses[0].Status != domain.HypothesisInconclusive {
				t.Fatal("limited evidence promoted")
			}
		})
	}
}

func TestInsightInvestigationRejectsIntegrityAndAdmissionViolations(t *testing.T) {
	cases := map[string]func(*domain.InsightInvestigationRequest){
		"tampered-bytes":   func(r *domain.InsightInvestigationRequest) { r.Sources[0].Bundle.Items[0].Content += "changed" },
		"cross-repository": func(r *domain.InsightInvestigationRequest) { r.Sources[0].Assessments[0].Repository.Repo = "other" },
		"head":             func(r *domain.InsightInvestigationRequest) { r.Sources[0].Bundle.Repository.Commit = "HEAD" },
		"duplicate-source": func(r *domain.InsightInvestigationRequest) { r.Sources = append(r.Sources, r.Sources[0]) },
		"duplicate-path": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Bundle.Items = append(r.Sources[0].Bundle.Items, r.Sources[0].Bundle.Items[0])
		},
		"unknown-class": func(r *domain.InsightInvestigationRequest) { r.MechanismClasses = []string{"INVENTED"} },
		"duplicate-class": func(r *domain.InsightInvestigationRequest) {
			r.MechanismClasses = append(r.MechanismClasses, r.MechanismClasses[0])
		},
		"hypothesis-budget": func(r *domain.InsightInvestigationRequest) { r.MaxHypotheses = 13 },
		"control":           func(r *domain.InsightInvestigationRequest) { r.Question = "bad\nquestion" },
		"byte-budget": func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Bundle.Items[0].Content = strings.Repeat("x", 2<<20)
		},
		"source-budget": func(r *domain.InsightInvestigationRequest) { r.Sources = make([]domain.InsightEvidenceInput, 9) },
		"whitespace":    func(r *domain.InsightInvestigationRequest) { r.Target.Name = " " },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			r := insightFixture()
			mutate(&r)
			result, err := InvestigateInsights(context.Background(), r)
			if err == nil || result.ID != "" {
				t.Fatalf("invalid input accepted: %v", err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := InvestigateInsights(ctx, insightFixture()); !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation ignored")
	}
}

func TestInsightInvestigationUnboundEvaluationsCannotCreateMechanism(t *testing.T) {
	request := insightFixture()
	for i := range request.Sources[0].Assessments[0].Capabilities[0].Sources {
		request.Sources[0].Assessments[0].Capabilities[0].Sources[i].ContentHash = "sha256:" + strings.Repeat("0", 64)
	}
	result, err := InvestigateInsights(context.Background(), request)
	if err != nil || len(result.Mechanisms) != 0 || len(result.Hypotheses) != 0 || len(result.ExtractionGaps) == 0 {
		t.Fatal("unbound assessments accepted")
	}
}
