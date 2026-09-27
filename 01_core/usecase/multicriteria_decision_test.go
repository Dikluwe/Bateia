/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/multicriteria_decision_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func multicriteriaSource(ref domain.CommitRef, kind domain.EvidenceKind, path, snippet string) domain.EvidenceSource {
	source := domain.EvidenceSource{Kind: kind, Ref: ref, Path: path, RuleID: "MCDA-FIXTURE", ContentHash: contentHash(snippet), Strength: 1}
	if kind != domain.EvidenceRepository {
		source.Line = 1
		source.Snippet = snippet
	}
	return source
}

func multicriteriaAssessment(owner, repo, license string) domain.CandidateAssessment {
	ref := domain.CommitRef{Owner: owner, Repo: repo, Commit: "abc123"}
	codeSource := multicriteriaSource(ref, domain.EvidenceFile, "evidence.go", "func Core() {}")
	licenseSource := multicriteriaSource(ref, domain.EvidenceLicense, "LICENSE", license+" license")
	maintenanceSource := multicriteriaSource(ref, domain.EvidenceRepository, "@repository-metadata", "active repository")
	return domain.CandidateAssessment{
		Repository:       ref,
		ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource}},
		Capabilities:     []domain.CapabilityEvidence{{Capability: "core", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource}}},
		License:          domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: license, Sources: []domain.EvidenceSource{licenseSource}},
		Maintenance:      domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{maintenanceSource}},
		Operation:        domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{codeSource}},
	}
}

func multicriteriaFixture() domain.MulticriteriaDecisionRequest {
	a := multicriteriaAssessment("acme", "alpha", "MIT")
	b := multicriteriaAssessment("acme", "beta", "MIT")
	c := multicriteriaAssessment("acme", "gamma", "MIT")
	criteria := []domain.MulticriteriaCriterion{
		{ID: "coverage", Label: "Cobertura", Direction: domain.CriterionMaximize, ScaleMin: 0, ScaleMax: 1},
		{ID: "adaptation_cost", Label: "Custo de adaptação", Direction: domain.CriterionMinimize, ScaleMin: 0, ScaleMax: 1},
	}
	measurement := func(assessment domain.CandidateAssessment, criterion string, value float64) domain.CriterionMeasurement {
		return domain.CriterionMeasurement{Candidate: assessment.Repository.Key(), CriterionID: criterion, Status: domain.EvidenceProven, Interval: domain.CriterionInterval{Low: value, High: value}, Sources: []domain.EvidenceSource{multicriteriaSource(assessment.Repository, domain.EvidenceFile, "measurement.go", criterion)}}
	}
	return domain.MulticriteriaDecisionRequest{
		Idea:       domain.IdeaSpec{Problem: "selecionar componente", Capabilities: []string{"core"}},
		Assessment: domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{c, a, b}},
		Policy:     domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, RequireLocal: true, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3},
		Criteria:   criteria,
		Measurements: []domain.CriterionMeasurement{
			measurement(a, "coverage", .9), measurement(a, "adaptation_cost", .6),
			measurement(b, "coverage", .8), measurement(b, "adaptation_cost", .4),
			measurement(c, "coverage", .7), measurement(c, "adaptation_cost", .7),
		},
		Scenarios: []domain.PreferenceScenario{},
	}
}

func TestMulticriteriaDecisionExposesParetoFrontierAndTradeoff(t *testing.T) {
	report, err := AnalyzeMulticriteriaDecision(multicriteriaFixture())
	if err != nil {
		t.Fatal(err)
	}
	wantFrontier := []string{"acme/alpha@abc123", "acme/beta@abc123"}
	if !reflect.DeepEqual(report.Frontier, wantFrontier) || !reflect.DeepEqual(report.RobustFrontier, wantFrontier) {
		t.Fatalf("fronteira inesperada: %+v / %+v", report.Frontier, report.RobustFrontier)
	}
	if len(report.Dominated) != 1 || report.Dominated[0].Candidate != "acme/gamma@abc123" {
		t.Fatalf("dominado não explicado: %+v", report.Dominated)
	}
	if report.StopReason != "PREFERENCES_REQUIRED" || len(report.Questions) == 0 || report.Questions[0].Kind != "PREFERENCE_TRADEOFF" {
		t.Fatalf("abstenção ou pergunta ausente: %+v", report)
	}
}

func TestMulticriteriaDecisionNeverCompensatesFailedGate(t *testing.T) {
	request := multicriteriaFixture()
	request.Assessment.Assessments = request.Assessment.Assessments[:2]
	for index := range request.Assessment.Assessments {
		if request.Assessment.Assessments[index].Repository.Repo == "alpha" {
			request.Assessment.Assessments[index].License.Identifier = "GPL-3.0"
		}
	}
	allowed := map[string]bool{}
	for _, assessment := range request.Assessment.Assessments {
		allowed[assessment.Repository.Key()] = true
	}
	filtered := []domain.CriterionMeasurement{}
	for _, measurement := range request.Measurements {
		if allowed[measurement.Candidate] {
			filtered = append(filtered, measurement)
		}
	}
	request.Measurements = filtered
	report, err := AnalyzeMulticriteriaDecision(request)
	if err != nil {
		t.Fatal(err)
	}
	for _, candidate := range report.Frontier {
		if candidate == "acme/alpha@abc123" {
			t.Fatal("candidato com gate de licença falho entrou na fronteira")
		}
	}
	if report.Eligibility[0].Candidate != "acme/alpha@abc123" || report.Eligibility[0].Status != domain.CandidateIneligible {
		t.Fatalf("gate não preservado: %+v", report.Eligibility)
	}
}

func TestMulticriteriaDecisionTreatsMissingMeasurementAsUncertainty(t *testing.T) {
	request := multicriteriaFixture()
	request.Assessment.Assessments = request.Assessment.Assessments[1:]
	filtered := []domain.CriterionMeasurement{}
	for _, measurement := range request.Measurements {
		if measurement.Candidate == "acme/alpha@abc123" || (measurement.Candidate == "acme/beta@abc123" && measurement.CriterionID == "coverage") {
			filtered = append(filtered, measurement)
		}
	}
	request.Measurements = filtered
	request.Scenarios = []domain.PreferenceScenario{{ID: "balanced", Weights: []domain.CriterionWeight{{CriterionID: "coverage", Weight: .5}, {CriterionID: "adaptation_cost", Weight: .5}}}}
	report, err := AnalyzeMulticriteriaDecision(request)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Dominated) != 0 || len(report.Frontier) != 2 || report.Pairwise[0].Relation != domain.PairwiseUncertain || report.StopReason != "EVIDENCE_INCONCLUSIVE" {
		t.Fatalf("ausência virou derrota ou valor pontual: %+v", report)
	}
}

func TestMulticriteriaDecisionRejectsUnknownMeasurementStatus(t *testing.T) {
	request := multicriteriaFixture()
	request.Measurements[0].Status = domain.EvidenceStatus("CERTAIN")
	if _, err := AnalyzeMulticriteriaDecision(request); err == nil {
		t.Fatal("status desconhecido deveria ser rejeitado")
	}
}

func TestMulticriteriaDecisionShowsWhichDeclaredScenarioMakesEachCandidateWin(t *testing.T) {
	request := multicriteriaFixture()
	request.Scenarios = []domain.PreferenceScenario{
		{ID: "1-coverage-first", Weights: []domain.CriterionWeight{{CriterionID: "coverage", Weight: .8}, {CriterionID: "adaptation_cost", Weight: .2}}},
		{ID: "2-cost-first", Weights: []domain.CriterionWeight{{CriterionID: "coverage", Weight: .2}, {CriterionID: "adaptation_cost", Weight: .8}}},
	}
	report, err := AnalyzeMulticriteriaDecision(request)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(report.Scenarios[0].NecessaryWinners, []string{"acme/alpha@abc123"}) || !reflect.DeepEqual(report.Scenarios[1].NecessaryWinners, []string{"acme/beta@abc123"}) {
		t.Fatalf("vencedores por cenário inesperados: %+v", report.Scenarios)
	}
	if report.StopReason != "PREFERENCE_SENSITIVE" || !report.Sensitivity.PreferenceSensitive || len(report.Sensitivity.RobustWinners) != 0 {
		t.Fatalf("sensibilidade não explicitada: %+v", report.Sensitivity)
	}
}

func TestMulticriteriaDecisionIsInvariantToInputOrder(t *testing.T) {
	request := multicriteriaFixture()
	request.Scenarios = []domain.PreferenceScenario{{ID: "balanced", Weights: []domain.CriterionWeight{{CriterionID: "coverage", Weight: .5}, {CriterionID: "adaptation_cost", Weight: .5}}}}
	first, err := AnalyzeMulticriteriaDecision(request)
	if err != nil {
		t.Fatal(err)
	}
	for left, right := 0, len(request.Assessment.Assessments)-1; left < right; left, right = left+1, right-1 {
		request.Assessment.Assessments[left], request.Assessment.Assessments[right] = request.Assessment.Assessments[right], request.Assessment.Assessments[left]
	}
	for left, right := 0, len(request.Measurements)-1; left < right; left, right = left+1, right-1 {
		request.Measurements[left], request.Measurements[right] = request.Measurements[right], request.Measurements[left]
	}
	request.Criteria[0], request.Criteria[1] = request.Criteria[1], request.Criteria[0]
	request.Scenarios[0].Weights[0], request.Scenarios[0].Weights[1] = request.Scenarios[0].Weights[1], request.Scenarios[0].Weights[0]
	second, err := AnalyzeMulticriteriaDecision(request)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("ordem alterou resultado\nfirst=%+v\nsecond=%+v", first, second)
	}
}
