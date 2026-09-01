/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/authenticated_web_experiment.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

type WebExperimentAuthorization struct {
	Provider                   string `json:"provider"`
	Purpose                    string `json:"purpose"`
	Authentication             string `json:"authentication"`
	PrivateRepositoriesAllowed bool   `json:"private_repositories_allowed"`
	PlannedQueries             int    `json:"planned_queries"`
	Authorized                 bool   `json:"authorized"`
}

type WebExperimentBudget struct {
	QueriesPerArm      int `json:"queries_per_arm"`
	PrimaryRequests    int `json:"primary_requests"`
	MaxResultsPerQuery int `json:"max_results_per_query"`
	AttemptTimeoutSec  int `json:"attempt_timeout_seconds"`
	TotalTimeoutSec    int `json:"total_timeout_seconds"`
	MaxConcurrency     int `json:"max_concurrency"`
	MaxRetries         int `json:"max_retries"`
	SafeRemaining      int `json:"safe_remaining"`
}

type WebExperimentReceipt struct {
	SchemaVersion string `json:"schema_version"`
	RunID         string `json:"run_id"`
	Planned       int    `json:"planned"`
	Started       int    `json:"started"`
	Retried       int    `json:"retried"`
	Completed     int    `json:"completed"`
	Unknown       int    `json:"unknown"`
	SkippedBudget int    `json:"skipped_budget"`
	StopReason    string `json:"stop_reason,omitempty"`
	Verdict       string `json:"verdict"`
}
