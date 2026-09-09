/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/effect_binding_test.md
 * @layer L2
 * @updated 2026-09-08
 */

package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/02_shell/cli"
)

type binding63Verifier struct{ requests []domain.EffectRequest }

func (v *binding63Verifier) Verify(_ context.Context, request domain.EffectRequest) (domain.EffectReport, error) {
	v.requests = append(v.requests, request)
	return domain.EffectReport{SchemaVersion: "effect-run-v1", ContractID: "sort-preserve-v1", Status: "unknown", Reason: "source_identity_mismatch", SourceSHA256: strings.Repeat("b", 64), SourceBinding: domain.EffectSourceBinding{Status: "mismatch", ExpectedSHA256: request.ExpectedSourceSHA256, ActualSHA256: strings.Repeat("b", 64), Symbol: request.Symbol}, Execution: map[string]any{"stop_reason": "source_identity_mismatch"}}, nil
}

func TestEffectBinding63CLINormalizesAndPreservesReceipt(t *testing.T) {
	for _, tc := range []struct{ name, hash string }{{"omitted", ""}, {"lowercase", strings.Repeat("a1", 32)}, {"uppercase", strings.Repeat("AbCdEf01", 8)}} {
		t.Run(tc.name, func(t *testing.T) {
			verifier := &binding63Verifier{}
			var out bytes.Buffer
			args := []string{"verify", "--source", "reviewed.go", "--symbol", "SortCopy"}
			if tc.hash != "" {
				args = append(args, "--expect-sha256", tc.hash)
			}
			if err := cli.RunEffectsCommand(args, &out, verifier); err != nil {
				t.Fatal(err)
			}
			if len(verifier.requests) != 1 {
				t.Fatalf("chamadas: %d", len(verifier.requests))
			}
			request := verifier.requests[0]
			if request.ExpectedSourceSHA256 != strings.ToLower(tc.hash) || request.SourcePath != "reviewed.go" || request.Symbol != "SortCopy" {
				t.Fatalf("pedido: %+v", request)
			}
			var report map[string]any
			if err := json.Unmarshal(out.Bytes(), &report); err != nil {
				t.Fatalf("JSON inválido: %v: %s", err, out.String())
			}
			binding, ok := report["source_binding"].(map[string]any)
			if !ok || binding["status"] != "mismatch" || binding["expected_sha256"] != strings.ToLower(tc.hash) || binding["declaration"] != nil {
				t.Fatalf("CLI perdeu recibo: %#v", report)
			}
			if report["status"] != "unknown" || report["reason"] != "source_identity_mismatch" {
				t.Fatalf("CLI promoveu mismatch: %#v", report)
			}
		})
	}
}

func TestEffectBinding63CLIRejectsMalformedBeforeVerifier(t *testing.T) {
	for _, suffix := range [][]string{
		{"--expect-sha256", ""},
		{"--expect-sha256="},
		{"--expect-sha256"},
		{"--expect-sha256", " "},
		{"--expect-sha256", strings.Repeat("a", 63)},
		{"--expect-sha256", strings.Repeat("a", 65)},
		{"--expect-sha256", strings.Repeat("a", 63) + "g"},
		{"--expect-sha256", strings.Repeat("Ａ", 64)},
		{"--expect-sha256", strings.Repeat("a", 63) + "\n"},
		{"--expect-sha256", "0x" + strings.Repeat("a", 62)},
	} {
		t.Run(strings.Join(suffix, "="), func(t *testing.T) {
			verifier := &binding63Verifier{}
			var out bytes.Buffer
			args := append([]string{"verify", "--source", "reviewed.go", "--symbol", "SortCopy"}, suffix...)
			if err := cli.RunEffectsCommand(args, &out, verifier); err == nil {
				t.Fatal("argumento inválido aceito")
			}
			if len(verifier.requests) != 0 {
				t.Fatalf("argumento inválido chamou verifier: %+v", verifier.requests)
			}
		})
	}
}

func TestEffectBinding63CLIHelpExposesFlagWithoutExecution(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"verify", "--help"}} {
		verifier := &binding63Verifier{}
		var out bytes.Buffer
		if err := cli.RunEffectsCommand(args, &out, verifier); err != nil {
			t.Fatal(err)
		}
		if len(verifier.requests) != 0 {
			t.Fatal("ajuda executou verifier")
		}
		if !strings.Contains(out.String(), "--expect-sha256") {
			t.Fatalf("flag não documentada: %s", out.String())
		}
	}
}
