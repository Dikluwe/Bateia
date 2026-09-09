/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/effect_trace.md
 * @layer L1
 * @updated 2026-09-08
 */

package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"time"

	"bateia/01_core/domain"
)

// EffectTraceFingerprint identifies the effective request, not the entire receipt.
func EffectTraceFingerprint(in domain.EffectTraceInputs) string {
	validHash := func(value string) bool {
		if len(value) != 64 {
			return false
		}
		_, err := hex.DecodeString(value)
		return err == nil
	}
	if in.ContractID != "sort-preserve-v1" || !validHash(in.AnalyzerBinarySHA256) || !validHash(in.SourceSHA256) ||
		!token.IsIdentifier(in.Symbol) || !ast.IsExported(in.Symbol) ||
		in.ExecutionTimeoutNS <= 0 || in.ExecutionTimeoutNS > int64(30*time.Second) ||
		in.BuildTimeoutNS <= 0 || in.BuildTimeoutNS > int64(60*time.Second) ||
		in.StdoutLimitBytes <= 0 || in.StdoutLimitBytes > 1<<20 ||
		in.StderrLimitBytes <= 0 || in.StderrLimitBytes > 1<<20 {
		return ""
	}
	canonical := strings.Join([]string{
		"effect-request-v1", in.ContractID, strings.ToLower(in.AnalyzerBinarySHA256),
		strings.ToLower(in.SourceSHA256), in.Symbol, strconv.FormatUint(in.Seed, 10),
		strconv.FormatInt(in.ExecutionTimeoutNS, 10), strconv.FormatInt(in.BuildTimeoutNS, 10),
		strconv.FormatInt(in.StdoutLimitBytes, 10), strconv.FormatInt(in.StderrLimitBytes, 10),
	}, "\n") + "\n"
	sum := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(sum[:])
}
