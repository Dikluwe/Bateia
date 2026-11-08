/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/execution_boundaries_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"
)

type advancingRuntime struct{ ticks atomic.Int64 }

func (r *advancingRuntime) Now() time.Time {
	return time.Date(2042, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(r.ticks.Add(1)) * time.Millisecond)
}
func (*advancingRuntime) WithTimeout(ctx context.Context, _ time.Duration) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}
func TestFunnelUsesInjectedClockForReportAndTelemetry(t *testing.T) {
	runtime := &advancingRuntime{}
	funnel := NewSearchFunnel(&MockSearchProvider{}, nil, NewMockCacheStore()).WithExecutionRuntime(runtime)
	report, err := funnel.Execute(context.Background(), domain.QueryFilter{Query: "queue", Language: "go", Depth: 1, MaxResults: 5})
	if err != nil {
		t.Fatal(err)
	}
	if report.GeneratedAt.Year() != 2042 || report.Telemetry == nil || report.Telemetry.TotalDurationMs <= 0 {
		t.Fatalf("runtime ignored: %+v", report)
	}
}
func TestIdeaRunTemporalMetadataDoesNotChangeBusinessProjection(t *testing.T) {
	request, deps := replayRunFixture()
	without, err := RunIdeaPipeline(context.Background(), request, deps)
	if err != nil {
		t.Fatal(err)
	}
	deps.Runtime = &advancingRuntime{}
	with, err := RunIdeaPipeline(context.Background(), request, deps)
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(without)
	b, _ := json.Marshal(with)
	if string(a) != string(b) {
		t.Fatal("temporal runtime changed business projection")
	}
	if len(with.StageDurations) == 0 {
		t.Fatal("missing stage metadata")
	}
	for stage, duration := range with.StageDurations {
		if duration <= 0 {
			t.Fatalf("runtime ignored for %s: %d", stage, duration)
		}
	}
	for stage, duration := range without.StageDurations {
		if duration != 0 {
			t.Fatalf("missing runtime invented duration for %s: %d", stage, duration)
		}
	}
}
