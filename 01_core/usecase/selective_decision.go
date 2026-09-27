/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/selective_decision.md
 * @layer L1
 * @updated 2026-09-24
 */
package usecase

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

type scoredDecision struct {
	score float64
	error bool
}

// EvaluateSelectiveDecisions measures final-decision risk only among non-abstentions.
// Missing gold labels are reported but never enter coverage or risk denominators.
func EvaluateSelectiveDecisions(request domain.SelectiveDecisionRequest) (domain.SelectiveDecisionReport, error) {
	report := domain.SelectiveDecisionReport{
		SchemaVersion: "1", Split: request.Split, Total: len(request.Cases),
		RiskCoverage: []domain.SelectiveDecisionPoint{}, CurveStatus: "SCORE_UNAVAILABLE",
		RiskBoundMethod: "one-sided Wilson 95% (approximate)",
	}
	if request.SchemaVersion != "1" {
		return report, fmt.Errorf("schema_version deve ser 1")
	}
	if request.Split != domain.SplitContract && request.Split != domain.SplitTrain && request.Split != domain.SplitValidation && request.Split != domain.SplitTest {
		return report, fmt.Errorf("split inválido: %q", request.Split)
	}
	if len(request.Cases) == 0 {
		return report, fmt.Errorf("cases não pode ser vazio")
	}
	seen := make(map[string]bool, len(request.Cases))
	var scored []scoredDecision
	allScored := true
	for _, c := range request.Cases {
		id := strings.TrimSpace(c.ID)
		if id == "" || seen[id] {
			return report, fmt.Errorf("id de caso vazio ou duplicado: %q", c.ID)
		}
		seen[id] = true
		if !validSelectiveDecision(c.Observed) {
			return report, fmt.Errorf("%s: decisão observada inválida: %q", id, c.Observed)
		}
		if c.Reliability != nil && (math.IsNaN(*c.Reliability) || math.IsInf(*c.Reliability, 0) || *c.Reliability < 0 || *c.Reliability > 1) {
			return report, fmt.Errorf("%s: reliability deve estar entre 0 e 1", id)
		}
		if c.Expected == "" {
			if len(c.Acceptable) != 0 {
				return report, fmt.Errorf("%s: acceptable exige expected", id)
			}
			report.Unadjudicated++
			continue
		}
		if !validSelectiveDecision(c.Expected) {
			return report, fmt.Errorf("%s: decisão esperada inválida: %q", id, c.Expected)
		}
		for _, a := range c.Acceptable {
			if !validSelectiveDecision(a) {
				return report, fmt.Errorf("%s: decisão aceitável inválida: %q", id, a)
			}
		}
		report.Adjudicated++
		if c.Observed == domain.DecisionUnknown {
			report.Abstained++
			if acceptableSelectiveDecision(domain.DecisionUnknown, c.Expected, c.Acceptable) {
				report.CorrectAbstentions++
			} else {
				report.AvoidableAbstentions++
			}
			continue
		}
		report.Accepted++
		wrong := !acceptableSelectiveDecision(c.Observed, c.Expected, c.Acceptable)
		if wrong {
			report.Errors++
		}
		if c.Reliability == nil {
			allScored = false
		} else {
			scored = append(scored, scoredDecision{score: *c.Reliability, error: wrong})
		}
	}
	report.Current = selectivePoint(nil, report.Accepted, report.Errors, report.Adjudicated)
	if report.Unadjudicated > 0 {
		report.Warnings = append(report.Warnings, "casos sem julgamento foram excluídos dos denominadores")
	}
	if report.Accepted < 30 {
		report.Warnings = append(report.Warnings, "menos de 30 decisões conclusivas; a estimativa de risco é instável")
	}
	if request.Split == domain.SplitTrain || request.Split == domain.SplitContract {
		report.Warnings = append(report.Warnings, "split de treino/contrato não estima desempenho fora da amostra")
	}
	if report.Accepted == 0 {
		report.CurveStatus = "NO_ACCEPTED_DECISIONS"
		return report, nil
	}
	if !allScored {
		report.Warnings = append(report.Warnings, "curva indisponível: falta reliability para uma ou mais decisões conclusivas")
		return report, nil
	}
	sort.Slice(scored, func(i, j int) bool { return scored[i].score > scored[j].score })
	accepted, errors := 0, 0
	for i := 0; i < len(scored); {
		threshold := scored[i].score
		for i < len(scored) && scored[i].score == threshold {
			accepted++
			if scored[i].error {
				errors++
			}
			i++
		}
		t := threshold
		report.RiskCoverage = append(report.RiskCoverage, selectivePoint(&t, accepted, errors, report.Adjudicated))
	}
	report.CurveStatus = "AVAILABLE"
	return report, nil
}

func validSelectiveDecision(d domain.ReuseDecision) bool {
	switch d {
	case domain.DecisionFork, domain.DecisionExtend, domain.DecisionCompose, domain.DecisionIntegrate, domain.DecisionReference, domain.DecisionBuild, domain.DecisionUnknown:
		return true
	}
	return false
}

func acceptableSelectiveDecision(observed, expected domain.ReuseDecision, acceptable []domain.ReuseDecision) bool {
	if observed == expected {
		return true
	}
	for _, d := range acceptable {
		if observed == d {
			return true
		}
	}
	return false
}

func selectivePoint(threshold *float64, accepted, errors, adjudicated int) domain.SelectiveDecisionPoint {
	p := domain.SelectiveDecisionPoint{MinimumReliability: threshold, Accepted: accepted, Errors: errors}
	if adjudicated > 0 {
		p.Coverage = float64(accepted) / float64(adjudicated)
	}
	if accepted > 0 {
		risk := float64(errors) / float64(accepted)
		p.SelectiveRisk = &risk
		upper := wilsonRiskUpper95(errors, accepted)
		p.RiskUpper95 = &upper
	}
	return p
}

func wilsonRiskUpper95(errors, accepted int) float64 {
	const z = 1.6448536269514722
	n := float64(accepted)
	p := float64(errors) / n
	z2 := z * z
	return (p + z2/(2*n) + z*math.Sqrt(p*(1-p)/n+z2/(4*n*n))) / (1 + z2/n)
}
