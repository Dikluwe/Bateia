/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/web_snapshot.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "context"

type ObservableWebRequest struct {
	RequestID  string `json:"request_id"`
	QueryID    string `json:"query_id"`
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
	Locale     string `json:"locale"`
}

type WebSearchCapabilities struct {
	Provider        string `json:"provider"`
	SupportsLocale  bool   `json:"supports_locale"`
	PreservesRank   bool   `json:"preserves_rank"`
	MaxResultsLimit int    `json:"max_results_limit"`
}

type WebReceipt struct {
	SanitizedURL string `json:"sanitized_url"`
	ResponseTime int64  `json:"response_time_ms"`
}

type RawWebResult struct {
	Rank         int    `json:"rank"`
	URL          string `json:"url"`
	CanonicalURL string `json:"canonical_url,omitempty"`
	Title        string `json:"title"`
	Snippet      string `json:"snippet"`
	IsAd         bool   `json:"is_ad"`
}

type RawWebResponse struct {
	RequestID   string         `json:"request_id"`
	QueryID     string         `json:"query_id"`
	Provider    string         `json:"provider"`
	RetrievedAt string         `json:"retrieved_at"`
	Status      string         `json:"status"`
	HTTPStatus  int            `json:"http_status,omitempty"`
	Results     []RawWebResult `json:"results"`
	Receipt     WebReceipt     `json:"receipt"`
}

type WebSearchSnapshot struct {
	SchemaVersion string           `json:"schema_version"`
	RunID         string           `json:"run_id"`
	Arm           string           `json:"arm"`
	GrammarHash   string           `json:"grammar_hash"`
	QuerySetHash  string           `json:"query_set_hash"`
	Provider      string           `json:"provider"`
	Responses     []RawWebResponse `json:"responses"`
	SnapshotHash  string           `json:"snapshot_hash"`
}

type WebReplayResult struct {
	SchemaVersion string         `json:"schema_version"`
	SnapshotHash  string         `json:"snapshot_hash"`
	TotalQueries  int            `json:"total_queries"`
	UniqueHits    int            `json:"unique_hits"`
	Deduplicated  []RawWebResult `json:"deduplicated"`
	MissingFields []string       `json:"missing_fields"`
	Replayable    bool           `json:"replayable"`
}

type BlindedJudgment struct {
	CanonicalURL string `json:"canonical_url"`
	Judgment     string `json:"judgment"`
	IsSoftware   bool   `json:"is_software"`
	IsRepoBacked bool   `json:"is_repo_backed"`
}

type BlindedJudgmentsFile struct {
	Judgments []BlindedJudgment `json:"judgments"`
}

type GrammarExperimentReport struct {
	SchemaVersion             string  `json:"schema_version"`
	JudgedPrecisionAt10Nominal float64 `json:"judged_precision_at_10_nominal"`
	JudgedPrecisionAt10F1      float64 `json:"judged_precision_at_10_f1"`
	JudgedPrecisionAt10F2      float64 `json:"judged_precision_at_10_f2"`
	RepoBackedRecallNominal   float64 `json:"repo_backed_recall_nominal"`
	RepoBackedRecallF1        float64 `json:"repo_backed_recall_f1"`
	RepoBackedRecallF2        float64 `json:"repo_backed_recall_f2"`
	RecallGainOverNominal     float64 `json:"recall_gain_over_nominal"`
	Verdict                   string  `json:"verdict"`
}

type ObservableWebSearch interface {
	Capabilities(context.Context) WebSearchCapabilities
	SearchOne(context.Context, ObservableWebRequest) (RawWebResponse, error)
}
