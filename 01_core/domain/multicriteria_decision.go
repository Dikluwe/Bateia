/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/multicriteria_decision.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

const MulticriteriaDecisionSchemaVersion = "1"

type CriterionDirection string

const (
	CriterionMaximize CriterionDirection = "MAXIMIZE"
	CriterionMinimize CriterionDirection = "MINIMIZE"
)

type MulticriteriaCriterion struct {
	ID        string             `json:"id"`
	Label     string             `json:"label"`
	Direction CriterionDirection `json:"direction"`
	ScaleMin  float64            `json:"scale_min"`
	ScaleMax  float64            `json:"scale_max"`
	Unit      string             `json:"unit,omitempty"`
}

type CriterionInterval struct {
	Low  float64 `json:"low"`
	High float64 `json:"high"`
}

type CriterionMeasurement struct {
	Candidate   string            `json:"candidate"`
	CriterionID string            `json:"criterion_id"`
	Status      EvidenceStatus    `json:"status"`
	Interval    CriterionInterval `json:"interval"`
	Sources     []EvidenceSource  `json:"sources,omitempty"`
}

type CriterionWeight struct {
	CriterionID string  `json:"criterion_id"`
	Weight      float64 `json:"weight"`
}

type PreferenceScenario struct {
	ID      string            `json:"id"`
	Label   string            `json:"label"`
	Weights []CriterionWeight `json:"weights"`
}

type MulticriteriaDecisionRequest struct {
	Idea         IdeaSpec                 `json:"idea"`
	Assessment   AssessmentReport         `json:"assessment"`
	Policy       ReuseDecisionPolicy      `json:"policy"`
	Composition  *CompositionInput        `json:"composition,omitempty"`
	Criteria     []MulticriteriaCriterion `json:"criteria"`
	Measurements []CriterionMeasurement   `json:"measurements"`
	Scenarios    []PreferenceScenario     `json:"preference_scenarios"`
}

type CandidateEligibilityStatus string

const (
	CandidateEligible                CandidateEligibilityStatus = "ELIGIBLE"
	CandidateIneligible              CandidateEligibilityStatus = "INELIGIBLE"
	CandidateEligibilityInconclusive CandidateEligibilityStatus = "INCONCLUSIVE"
)

type CandidateEligibility struct {
	Candidate string                     `json:"candidate"`
	Status    CandidateEligibilityStatus `json:"status"`
	Reasons   []DecisionReason           `json:"reasons"`
}

type PairwiseRelationKind string

const (
	PairwiseDominates    PairwiseRelationKind = "DOMINATES"
	PairwiseDominatedBy  PairwiseRelationKind = "DOMINATED_BY"
	PairwiseEquivalent   PairwiseRelationKind = "EQUIVALENT"
	PairwiseIncomparable PairwiseRelationKind = "INCOMPARABLE"
	PairwiseUncertain    PairwiseRelationKind = "EVIDENCE_UNCERTAIN"
)

type PairwiseCriterionComparison struct {
	CriterionID string            `json:"criterion_id"`
	Left        CriterionInterval `json:"left"`
	Right       CriterionInterval `json:"right"`
}

type PairwiseRelation struct {
	Left               string                        `json:"left"`
	Right              string                        `json:"right"`
	Relation           PairwiseRelationKind          `json:"relation"`
	LeftBetter         []string                      `json:"left_better"`
	RightBetter        []string                      `json:"right_better"`
	UnresolvedCriteria []string                      `json:"unresolved_criteria"`
	Comparisons        []PairwiseCriterionComparison `json:"comparisons"`
}

type DominatedCandidate struct {
	Candidate   string   `json:"candidate"`
	DominatedBy []string `json:"dominated_by"`
}

type CriterionContribution struct {
	CriterionID string            `json:"criterion_id"`
	Weight      float64           `json:"weight"`
	Utility     CriterionInterval `json:"utility"`
	Weighted    CriterionInterval `json:"weighted"`
}

type ScenarioCandidateScore struct {
	Candidate     string                  `json:"candidate"`
	Score         CriterionInterval       `json:"score"`
	Contributions []CriterionContribution `json:"contributions"`
}

type PreferenceScenarioResult struct {
	ScenarioID       string                   `json:"scenario_id"`
	Scores           []ScenarioCandidateScore `json:"scores"`
	NecessaryWinners []string                 `json:"necessary_winners"`
	PossibleWinners  []string                 `json:"possible_winners"`
}

type RankingSensitivity struct {
	RobustWinners       []string `json:"robust_winners"`
	PossibleWinners     []string `json:"possible_winners"`
	PreferenceSensitive bool     `json:"preference_sensitive"`
}

type MulticriteriaQuestion struct {
	Kind       string   `json:"kind"`
	Candidates []string `json:"candidates"`
	Criteria   []string `json:"criteria"`
	Question   string   `json:"question"`
}

type MulticriteriaDecisionReport struct {
	SchemaVersion  string                     `json:"schema_version"`
	Method         string                     `json:"method"`
	Composition    *CompositionReport         `json:"composition,omitempty"`
	Eligibility    []CandidateEligibility     `json:"eligibility"`
	Frontier       []string                   `json:"frontier"`
	RobustFrontier []string                   `json:"robust_frontier"`
	Dominated      []DominatedCandidate       `json:"dominated"`
	Pairwise       []PairwiseRelation         `json:"pairwise"`
	Scenarios      []PreferenceScenarioResult `json:"preference_scenarios"`
	Sensitivity    RankingSensitivity         `json:"sensitivity"`
	Recommended    []string                   `json:"recommended_candidates"`
	Questions      []MulticriteriaQuestion    `json:"questions"`
	StopReason     string                     `json:"stop_reason"`
}
