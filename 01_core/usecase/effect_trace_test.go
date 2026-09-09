/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/effect_trace_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package usecase_test

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func trace64Input() domain.EffectTraceInputs {
	return domain.EffectTraceInputs{ContractID: "sort-preserve-v1", AnalyzerBinarySHA256: strings.Repeat("a", 64), SourceSHA256: strings.Repeat("b", 64), Symbol: "Sort", Seed: 1, ExecutionTimeoutNS: 5000000000, BuildTimeoutNS: 60000000000, StdoutLimitBytes: 1048576, StderrLimitBytes: 1048576}
}
func trace64Oracle(in domain.EffectTraceInputs) string {
	fields := []string{"effect-request-v1", in.ContractID, strings.ToLower(in.AnalyzerBinarySHA256), strings.ToLower(in.SourceSHA256), in.Symbol, strconv.FormatUint(in.Seed, 10), strconv.FormatInt(in.ExecutionTimeoutNS, 10), strconv.FormatInt(in.BuildTimeoutNS, 10), strconv.FormatInt(in.StdoutLimitBytes, 10), strconv.FormatInt(in.StderrLimitBytes, 10)}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(fields, "\n")+"\n")))
}
func TestTrace64KnownVectorNormalizationAndImmutability(t *testing.T) {
	in := trace64Input()
	const want = "c35a66c9f2f94248db649702db13107b8f2af85c8a20e4c701800ff7033590d9"
	if got := usecase.EffectTraceFingerprint(in); got != want {
		t.Fatalf("fingerprint=%q want %q", got, want)
	}
	in.AnalyzerBinarySHA256 = strings.ToUpper(in.AnalyzerBinarySHA256)
	in.SourceSHA256 = strings.ToUpper(in.SourceSHA256)
	before := in
	if got := usecase.EffectTraceFingerprint(in); got != want {
		t.Errorf("uppercase fingerprint=%q", got)
	}
	if in != before {
		t.Fatal("input mutated")
	}
}
func TestTrace64ValidSensitivityAndBoundaries(t *testing.T) {
	base := trace64Input()
	baseline := trace64Oracle(base)
	cases := map[string]func(*domain.EffectTraceInputs){
		"analyzer":       func(v *domain.EffectTraceInputs) { v.AnalyzerBinarySHA256 = strings.Repeat("c", 64) },
		"source":         func(v *domain.EffectTraceInputs) { v.SourceSHA256 = strings.Repeat("d", 64) },
		"symbol":         func(v *domain.EffectTraceInputs) { v.Symbol = "Other" },
		"unicode_export": func(v *domain.EffectTraceInputs) { v.Symbol = "Éxample２" },
		"seed_zero":      func(v *domain.EffectTraceInputs) { v.Seed = 0 },
		"seed_max":       func(v *domain.EffectTraceInputs) { v.Seed = ^uint64(0) },
		"execution_min":  func(v *domain.EffectTraceInputs) { v.ExecutionTimeoutNS = 1 },
		"execution_max":  func(v *domain.EffectTraceInputs) { v.ExecutionTimeoutNS = 30000000000 },
		"build_min":      func(v *domain.EffectTraceInputs) { v.BuildTimeoutNS = 1 },
		"stdout_min":     func(v *domain.EffectTraceInputs) { v.StdoutLimitBytes = 1 },
		"stderr_min":     func(v *domain.EffectTraceInputs) { v.StderrLimitBytes = 1 },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			in := base
			change(&in)
			want := trace64Oracle(in)
			if want == baseline {
				t.Fatal("non-discriminating oracle")
			}
			if got := usecase.EffectTraceFingerprint(in); got != want {
				t.Fatalf("got %q want %q", got, want)
			}
		})
	}
}
func TestTrace64InvalidInputs(t *testing.T) {
	var cases []struct {
		name   string
		change func(*domain.EffectTraceInputs)
	}
	for _, value := range []string{"", "other", "sort-preserve-v1\n"} {
		value := value
		cases = append(cases, struct {
			name   string
			change func(*domain.EffectTraceInputs)
		}{"contract_" + value, func(v *domain.EffectTraceInputs) { v.ContractID = value }})
	}
	for _, value := range []string{"", strings.Repeat("a", 63), strings.Repeat("a", 65), strings.Repeat("g", 64), strings.Repeat("Ａ", 64), strings.Repeat("a", 63) + "\n"} {
		value := value
		for _, field := range []string{"source", "analyzer"} {
			field := field
			cases = append(cases, struct {
				name   string
				change func(*domain.EffectTraceInputs)
			}{field + "_" + value, func(v *domain.EffectTraceInputs) {
				if field == "source" {
					v.SourceSHA256 = value
				} else {
					v.AnalyzerBinarySHA256 = value
				}
			}})
		}
	}
	for _, value := range []string{"", "sort", "_Sort", "1Sort", "Sort.Name", "Sort\n", "Sort-Name", "É!", "ǅName", "Sort\x00"} {
		value := value
		cases = append(cases, struct {
			name   string
			change func(*domain.EffectTraceInputs)
		}{"symbol_" + value, func(v *domain.EffectTraceInputs) { v.Symbol = value }})
	}
	for _, field := range []string{"execution", "build", "stdout", "stderr"} {
		field := field
		max := int64(1048576)
		if field == "execution" {
			max = 30000000000
		}
		if field == "build" {
			max = 60000000000
		}
		for _, value := range []int64{-1, 0, max + 1, 1<<63 - 1} {
			value := value
			cases = append(cases, struct {
				name   string
				change func(*domain.EffectTraceInputs)
			}{fmt.Sprintf("%s_%d", field, value), func(v *domain.EffectTraceInputs) {
				switch field {
				case "execution":
					v.ExecutionTimeoutNS = value
				case "build":
					v.BuildTimeoutNS = value
				case "stdout":
					v.StdoutLimitBytes = value
				case "stderr":
					v.StderrLimitBytes = value
				}
			}})
		}
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := trace64Input()
			tc.change(&in)
			if got := usecase.EffectTraceFingerprint(in); got != "" {
				t.Fatalf("invalid input accepted: %q", got)
			}
		})
	}
}
