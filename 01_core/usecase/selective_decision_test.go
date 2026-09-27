/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/selective_decision_test.md
 * @layer L1
 * @updated 2026-09-24
 */
package usecase

import (
	"math"
	"testing"

	"bateia/01_core/domain"
)

func selectiveScore(v float64) *float64 { return &v }

func TestSelectiveDecisionRiskCoverageAndTies(t *testing.T) {
	request := domain.SelectiveDecisionRequest{SchemaVersion: "1", Split: domain.SplitTest, Cases: []domain.SelectiveDecisionCase{
		{ID: "correct-high", Observed: domain.DecisionBuild, Expected: domain.DecisionBuild, Reliability: selectiveScore(.9)},
		{ID: "wrong-mid", Observed: domain.DecisionBuild, Expected: domain.DecisionFork, Reliability: selectiveScore(.7)},
		{ID: "correct-mid", Observed: domain.DecisionFork, Expected: domain.DecisionFork, Reliability: selectiveScore(.7)},
		{ID: "abstain-correct", Observed: domain.DecisionUnknown, Expected: domain.DecisionUnknown},
		{ID: "abstain-avoidable", Observed: domain.DecisionUnknown, Expected: domain.DecisionBuild},
		{ID: "unlabeled", Observed: domain.DecisionBuild},
	}}
	report, err := EvaluateSelectiveDecisions(request)
	if err != nil {
		t.Fatal(err)
	}
	if report.Adjudicated != 5 || report.Unadjudicated != 1 || report.Accepted != 3 || report.Errors != 1 || report.Abstained != 2 || report.CorrectAbstentions != 1 || report.AvoidableAbstentions != 1 {
		t.Fatalf("contagens incorretas: %+v", report)
	}
	if report.Current.Coverage != .6 || *report.Current.SelectiveRisk != 1.0/3.0 {
		t.Fatalf("denominadores incorretos: %+v", report.Current)
	}
	if report.CurveStatus != "AVAILABLE" || len(report.RiskCoverage) != 2 || report.RiskCoverage[0].Accepted != 1 || report.RiskCoverage[0].Errors != 0 || report.RiskCoverage[1].Accepted != 3 || report.RiskCoverage[1].Errors != 1 {
		t.Fatalf("curva/ties incorretos: %+v", report)
	}
	if report.RiskCoverage[0].RiskUpper95 == nil || *report.RiskCoverage[0].RiskUpper95 <= 0 {
		t.Fatal("limite superior não pode ser zero após 1 acerto")
	}
}

func TestSelectiveDecisionNoScoreAndNoAccepted(t *testing.T) {
	for _, tc := range []struct {
		name, status string
		cases        []domain.SelectiveDecisionCase
	}{
		{"missing", "SCORE_UNAVAILABLE", []domain.SelectiveDecisionCase{{ID: "a", Observed: domain.DecisionBuild, Expected: domain.DecisionBuild}}},
		{"abstention", "NO_ACCEPTED_DECISIONS", []domain.SelectiveDecisionCase{{ID: "a", Observed: domain.DecisionUnknown, Expected: domain.DecisionBuild}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			report, err := EvaluateSelectiveDecisions(domain.SelectiveDecisionRequest{SchemaVersion: "1", Split: domain.SplitTest, Cases: tc.cases})
			if err != nil {
				t.Fatal(err)
			}
			if report.CurveStatus != tc.status || len(report.RiskCoverage) != 0 {
				t.Fatalf("curva fabricada: %+v", report)
			}
			if tc.name == "abstention" && (report.Current.SelectiveRisk != nil || report.Current.Coverage != 0) {
				t.Fatalf("risco deveria ser indefinido: %+v", report.Current)
			}
		})
	}
}

func TestSelectiveDecisionRejectsInvalidInput(t *testing.T) {
	base := domain.SelectiveDecisionRequest{SchemaVersion: "1", Split: domain.SplitTest, Cases: []domain.SelectiveDecisionCase{{ID: "a", Observed: domain.DecisionBuild, Expected: domain.DecisionBuild}}}
	for _, mutate := range []func(*domain.SelectiveDecisionRequest){
		func(r *domain.SelectiveDecisionRequest) { r.Split = "" },
		func(r *domain.SelectiveDecisionRequest) { r.Cases[0].Observed = "MAGIC" },
		func(r *domain.SelectiveDecisionRequest) { r.Cases[0].Reliability = selectiveScore(math.NaN()) },
		func(r *domain.SelectiveDecisionRequest) { r.Cases[0].Reliability = selectiveScore(1.1) },
		func(r *domain.SelectiveDecisionRequest) { r.Cases = append(r.Cases, r.Cases[0]) },
	} {
		r := base
		r.Cases = append([]domain.SelectiveDecisionCase(nil), base.Cases...)
		mutate(&r)
		if _, err := EvaluateSelectiveDecisions(r); err == nil {
			t.Fatalf("entrada inválida aceita: %+v", r)
		}
	}
}
