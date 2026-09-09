/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/effect_verification.md
 * @layer L1
 * @updated 2026-09-08
 */

package domain

import (
	"context"
	"time"
)

type EffectObservation struct {
	InputBefore []int64 `json:"input_before"`
	InputAfter  []int64 `json:"input_after"`
	Output      []int64 `json:"output"`
	Completed   bool    `json:"completed"`
	Failure     string  `json:"failure"`
}
type EffectVerdict struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}
type EffectRequest struct {
	ExpectedSourceSHA256 string        `json:"expected_source_sha256"`
	SourcePath           string        `json:"source_path"`
	Symbol               string        `json:"symbol"`
	Seed                 uint64        `json:"seed"`
	Timeout              time.Duration `json:"timeout"`
}
type EffectCase struct {
	Observation EffectObservation `json:"observation"`
	Verdict     EffectVerdict     `json:"verdict"`
}
type EffectSourceBinding struct {
	Status         string             `json:"status"`
	ExpectedSHA256 string             `json:"expected_sha256"`
	ActualSHA256   string             `json:"actual_sha256"`
	Symbol         string             `json:"symbol"`
	Declaration    *EffectDeclaration `json:"declaration"`
}
type EffectDeclaration struct {
	StartByte int    `json:"start_byte"`
	EndByte   int    `json:"end_byte"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Content   string `json:"content"`
	SHA256    string `json:"sha256"`
}
type EffectAnalyzerIdentity struct {
	Status        string `json:"status"`
	Reason        string `json:"reason"`
	BinarySHA256  string `json:"binary_sha256"`
	ModulePath    string `json:"module_path"`
	ModuleVersion string `json:"module_version"`
	GoVersion     string `json:"go_version"`
	VCSRevision   string `json:"vcs_revision"`
	VCSModified   *bool  `json:"vcs_modified"`
}
type EffectTraceInputs struct {
	ContractID           string `json:"contract_id"`
	AnalyzerBinarySHA256 string `json:"analyzer_binary_sha256"`
	SourceSHA256         string `json:"source_sha256"`
	Symbol               string `json:"symbol"`
	Seed                 uint64 `json:"seed"`
	ExecutionTimeoutNS   int64  `json:"execution_timeout_ns"`
	BuildTimeoutNS       int64  `json:"build_timeout_ns"`
	StdoutLimitBytes     int64  `json:"stdout_limit_bytes"`
	StderrLimitBytes     int64  `json:"stderr_limit_bytes"`
}
type EffectTrace struct {
	SchemaVersion string                 `json:"schema_version"`
	Status        string                 `json:"status"`
	Reason        string                 `json:"reason"`
	StartedAt     string                 `json:"started_at"`
	FinishedAt    string                 `json:"finished_at"`
	RequestSHA256 string                 `json:"request_sha256"`
	Analyzer      EffectAnalyzerIdentity `json:"analyzer"`
	Inputs        EffectTraceInputs      `json:"inputs"`
}
type EffectReport struct {
	Trace         EffectTrace         `json:"trace"`
	SourceBinding EffectSourceBinding `json:"source_binding"`
	SchemaVersion string              `json:"schema_version"`
	ContractID    string              `json:"contract_id"`
	Status        string              `json:"status"`
	Reason        string              `json:"reason"`
	SourceSHA256  string              `json:"source_sha256"`
	SourcePath    string              `json:"source_path"`
	Symbol        string              `json:"symbol"`
	Seed          uint64              `json:"seed"`
	Cases         []EffectCase        `json:"cases"`
	Limitations   []string            `json:"limitations"`
	Execution     map[string]any      `json:"execution"`
}
type EffectVerifier interface {
	Verify(context.Context, EffectRequest) (EffectReport, error)
}
