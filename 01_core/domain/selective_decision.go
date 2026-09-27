/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/selective_decision.md
 * @layer L1
 * @updated 2026-09-24
 */
package domain

// SelectiveDecisionCase records one independently adjudicated final decision.
// Reliability is an optional ordering score, never a calibrated probability.
type SelectiveDecisionCase struct {
	ID          string          `json:"id"`
	Group       string          `json:"group,omitempty"`
	Observed    ReuseDecision   `json:"observed"`
	Expected    ReuseDecision   `json:"expected,omitempty"`
	Acceptable  []ReuseDecision `json:"acceptable,omitempty"`
	Reliability *float64        `json:"reliability,omitempty"`
}

type SelectiveDecisionRequest struct {
	SchemaVersion string                  `json:"schema_version"`
	Split         BenchmarkSplit          `json:"split"`
	Cases         []SelectiveDecisionCase `json:"cases"`
}

type SelectiveDecisionPoint struct {
	MinimumReliability *float64 `json:"minimum_reliability,omitempty"`
	Accepted           int      `json:"accepted"`
	Errors             int      `json:"errors"`
	Coverage           float64  `json:"coverage"`
	SelectiveRisk      *float64 `json:"selective_risk"`
	RiskUpper95        *float64 `json:"risk_upper_95"`
}

type SelectiveDecisionReport struct {
	SchemaVersion        string                   `json:"schema_version"`
	Split                BenchmarkSplit           `json:"split"`
	Total                int                      `json:"total"`
	Adjudicated          int                      `json:"adjudicated"`
	Unadjudicated        int                      `json:"unadjudicated"`
	Accepted             int                      `json:"accepted"`
	Errors               int                      `json:"errors"`
	Abstained            int                      `json:"abstained"`
	CorrectAbstentions   int                      `json:"correct_abstentions"`
	AvoidableAbstentions int                      `json:"avoidable_abstentions"`
	Current              SelectiveDecisionPoint   `json:"current"`
	RiskCoverage         []SelectiveDecisionPoint `json:"risk_coverage"`
	CurveStatus          string                   `json:"curve_status"`
	RiskBoundMethod      string                   `json:"risk_bound_method"`
	Warnings             []string                 `json:"warnings,omitempty"`
}
