/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/effect_verification_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package domain_test

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type effectVerifierContractStub struct{}

func (effectVerifierContractStub) Verify(_ context.Context, request domain.EffectRequest) (domain.EffectReport, error) {
	return domain.EffectReport{SourcePath: request.SourcePath, Symbol: request.Symbol, Seed: request.Seed}, nil
}

var _ domain.EffectVerifier = effectVerifierContractStub{}

func TestEffectVerifierPublicInterface(t *testing.T) {
	request := domain.EffectRequest{SourcePath: "candidate.go", Symbol: "Transform", Seed: ^uint64(0), Timeout: time.Second}
	var verifier domain.EffectVerifier = effectVerifierContractStub{}
	got, err := verifier.Verify(context.Background(), request)
	if err != nil || got.SourcePath != request.SourcePath || got.Symbol != request.Symbol || got.Seed != request.Seed {
		t.Fatalf("unexpected interface result: %#v %v", got, err)
	}
}

func TestEffectPublicJSONAndInt64RoundTrip(t *testing.T) {
	observation := domain.EffectObservation{InputBefore: []int64{9223372036854775807, -9223372036854775808}, InputAfter: []int64{9223372036854775807, -9223372036854775808}, Output: []int64{-9223372036854775808, 9223372036854775807}, Completed: true, Failure: ""}
	verdict := domain.EffectVerdict{Status: "supported_in_scope", Reason: "checked values"}
	report := domain.EffectReport{SchemaVersion: "effect-run-v1", ContractID: "sort-preserve-v1", Status: verdict.Status, Reason: verdict.Reason, SourceSHA256: "digest", SourcePath: "candidate.go", Symbol: "Transform", Seed: ^uint64(0), Cases: []domain.EffectCase{{Observation: observation, Verdict: verdict}}, Limitations: []string{"finite scope"}, Execution: map[string]any{"kind": "tool telemetry"}}
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"schema_version", "contract_id", "status", "reason", "source_sha256", "source_path", "symbol", "seed", "cases", "limitations", "execution"} {
		if _, ok := fields[key]; !ok {
			t.Errorf("missing public JSON key %s", key)
		}
	}
	var cases []map[string]json.RawMessage
	if err := json.Unmarshal(fields["cases"], &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 1 {
		t.Fatalf("case count=%d", len(cases))
	}
	var o map[string]json.RawMessage
	if err := json.Unmarshal(cases[0]["observation"], &o); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"input_before", "input_after", "output", "completed"} {
		if _, ok := o[key]; !ok {
			t.Errorf("missing observation key %s", key)
		}
	}
	var v map[string]json.RawMessage
	if err := json.Unmarshal(cases[0]["verdict"], &v); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"status", "reason"} {
		if _, ok := v[key]; !ok {
			t.Errorf("missing verdict key %s", key)
		}
	}
	var roundTrip domain.EffectReport
	if err := json.Unmarshal(raw, &roundTrip); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(report, roundTrip) {
		t.Fatalf("public JSON roundtrip changed report: %#v", roundTrip)
	}
	observation.Failure = "panic marker"
	failureRaw, err := json.Marshal(observation)
	if err != nil {
		t.Fatal(err)
	}
	var failureFields map[string]json.RawMessage
	if err := json.Unmarshal(failureRaw, &failureFields); err != nil {
		t.Fatal(err)
	}
	if string(failureFields["failure"]) != `"panic marker"` {
		t.Fatalf("failure field not snake_case: %s", failureRaw)
	}
}
