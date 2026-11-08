/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/execution_runtime_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"context"
	"errors"
	"testing"
	"time"
)

type controlledRuntime struct {
	now      time.Time
	duration time.Duration
	parent   context.Context
}

func (r *controlledRuntime) Now() time.Time { return r.now }
func (r *controlledRuntime) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	r.duration = d
	r.parent = ctx
	return context.WithCancel(ctx)
}
func TestRuntimeIsAnExplicitDependency(t *testing.T) {
	ctx := context.Background()
	if !RuntimeNow(ctx).IsZero() {
		t.Fatal("missing runtime invented a date")
	}
	if _, _, err := ExecutionTimeout(ctx, time.Second); !errors.Is(err, ErrRuntimeNotConfigured) {
		t.Fatal("missing deadline runtime accepted")
	}
	r := &controlledRuntime{now: time.Date(2026, 10, 28, 1, 0, 0, 0, time.UTC)}
	ctx = WithExecutionRuntime(ctx, r)
	if RuntimeNow(ctx) != r.now {
		t.Fatal("controlled clock ignored")
	}
	next, cancel, err := ExecutionTimeout(ctx, 17*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if r.duration != 17*time.Second || r.parent != ctx {
		t.Fatal("deadline port input changed")
	}
	cancel()
	if !errors.Is(next.Err(), context.Canceled) {
		t.Fatal("cancel not propagated")
	}
	if ElapsedMillis(r.now.Add(25*time.Millisecond), r.now) != 25 || ElapsedMillis(r.now.Add(-time.Second), r.now) != 0 || ElapsedMillis(r.now, time.Time{}) != 0 {
		t.Fatal("elapsed metadata was fabricated")
	}
}
