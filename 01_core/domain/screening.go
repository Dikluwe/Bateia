/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/screening.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

type MatchEvidence struct {
	MatchType string `json:"match_type"`
	Path      string `json:"path,omitempty"`
	Snippet   string `json:"snippet,omitempty"`
}

type CandidateOrigin struct {
	ExpansionID        string          `json:"expansion_id"`
	PlannedQuery       string          `json:"planned_query"`
	AttemptIndex       int             `json:"attempt_index"`
	EffectiveQuery     string          `json:"effective_query"`
	Evidence           []MatchEvidence `json:"evidence,omitempty"`
	SearchLocale       string          `json:"search_locale,omitempty"`
	RegionalProposalID string          `json:"regional_proposal_id,omitempty"`
}

type ProviderCandidateMatch struct {
	Repository     string          `json:"repository"`
	AttemptIndex   int             `json:"attempt_index"`
	EffectiveQuery string          `json:"effective_query"`
	Evidence       []MatchEvidence `json:"evidence,omitempty"`
}

type ScreeningSignals struct {
	ExpansionPrior       float64 `json:"expansion_prior"`
	QueryRetention       float64 `json:"query_retention"`
	FallbackDepth        int     `json:"fallback_depth"`
	MatchSpecificity     float64 `json:"match_specificity"`
	DescriptionAlignment float64 `json:"description_alignment"`
	OriginConvergence    int     `json:"origin_convergence"`
	CapabilityCoverage   float64 `json:"capability_coverage"`
	NonGoalConflict      float64 `json:"non_goal_conflict"`
	EvidenceCompleteness float64 `json:"evidence_completeness"`
}

type SignalContribution struct {
	Signal string  `json:"signal"`
	Value  float64 `json:"value"`
	Weight float64 `json:"weight"`
	Effect float64 `json:"effect"`
	Reason string  `json:"reason"`
}

type CalibrationParameters struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type ScreeningModel struct {
	SchemaVersion     string                 `json:"schema_version"`
	Weights           map[string]float64     `json:"weights"`
	ExpansionPriors   map[string]float64     `json:"expansion_priors,omitempty"`
	AdvanceThreshold  float64                `json:"advance_threshold"`
	DropThreshold     float64                `json:"drop_threshold"`
	AuditSampleRate   float64                `json:"audit_sample_rate"`
	AuditSeed         int64                  `json:"audit_seed"`
	CalibrationStatus string                 `json:"calibration_status"`
	Calibration       *CalibrationParameters `json:"calibration,omitempty"`
}

type ScreeningDecision struct {
	Repository         string               `json:"repository"`
	RankingScore       float64              `json:"ranking_score"`
	Probability        *float64             `json:"probability"`
	CalibrationStatus  string               `json:"calibration_status"`
	Decision           string               `json:"decision"`
	Signals            ScreeningSignals     `json:"signals"`
	Contributions      []SignalContribution `json:"contributions"`
	UncertaintyReasons []string             `json:"uncertainty_reasons"`
	Origins            []CandidateOrigin    `json:"origins"`
}

type ScreeningReport struct {
	SchemaVersion string              `json:"schema_version"`
	Strategy      string              `json:"strategy"`
	Decisions     []ScreeningDecision `json:"decisions"`
}

type CalibrationFold struct {
	HeldOutIdea     string  `json:"held_out_idea"`
	TrainingCount   int     `json:"training_count"`
	TestCount       int     `json:"test_count"`
	ModelBrier      float64 `json:"model_brier"`
	PrevalenceBrier float64 `json:"prevalence_brier"`
	ModelLogLoss    float64 `json:"model_log_loss"`
}

type CalibrationBin struct {
	LowerBound         float64 `json:"lower_bound"`
	UpperBound         float64 `json:"upper_bound"`
	Count              int     `json:"count"`
	AverageProbability float64 `json:"average_probability"`
	ObservedRate       float64 `json:"observed_rate"`
	Sufficient         bool    `json:"sufficient"`
}

type CalibrationReport struct {
	SchemaVersion    string            `json:"schema_version"`
	Status           string            `json:"status"`
	ObservationCount int               `json:"observation_count"`
	ModelBrier       float64           `json:"model_brier"`
	PrevalenceBrier  float64           `json:"prevalence_brier"`
	ModelLogLoss     float64           `json:"model_log_loss"`
	PrecisionAtK     float64           `json:"precision_at_k"`
	Folds            []CalibrationFold `json:"folds"`
	Bins             []CalibrationBin  `json:"bins"`
	Model            ScreeningModel    `json:"model"`
}

type ScreeningDatasetGroup struct {
	IdeaID    string             `json:"idea_id"`
	Snapshot  DiscoverySnapshot  `json:"snapshot"`
	Judgments DiscoveryJudgments `json:"judgments"`
}

type AuditReport struct {
	SchemaVersion      string             `json:"schema_version"`
	ReviewedCount      int                `json:"reviewed_count"`
	FalseNegativeCount int                `json:"false_negative_count"`
	Revisions          []JudgmentRevision `json:"revisions"`
}

type AuditCandidate struct {
	Repository      string         `json:"repository"`
	CurrentLabel    RelevanceLabel `json:"current_label"`
	RankingScore    float64        `json:"ranking_score"`
	SelectionReason string         `json:"selection_reason"`
}

type AuditSelection struct {
	SchemaVersion string           `json:"schema_version"`
	Candidates    []AuditCandidate `json:"candidates"`
}

type JudgmentRevision struct {
	Repository    string         `json:"repository"`
	PreviousLabel RelevanceLabel `json:"previous_label"`
	NewLabel      RelevanceLabel `json:"new_label"`
	Evidence      string         `json:"evidence"`
	Reason        string         `json:"reason"`
}
