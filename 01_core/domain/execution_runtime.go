/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/execution_runtime.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"context"
	"errors"
	"time"
)

type ExecutionRuntime interface {
	Now() time.Time
	WithTimeout(context.Context, time.Duration) (context.Context, context.CancelFunc)
}
type executionRuntimeKey struct{}

var ErrRuntimeNotConfigured = errors.New("execution runtime is not configured")

func WithExecutionRuntime(ctx context.Context, runtime ExecutionRuntime) context.Context {
	if runtime == nil {
		return ctx
	}
	return context.WithValue(ctx, executionRuntimeKey{}, runtime)
}
func ExecutionRuntimeFromContext(ctx context.Context) ExecutionRuntime {
	runtime, _ := ctx.Value(executionRuntimeKey{}).(ExecutionRuntime)
	return runtime
}
func RuntimeNow(ctx context.Context) time.Time {
	if runtime := ExecutionRuntimeFromContext(ctx); runtime != nil {
		return runtime.Now()
	}
	return time.Time{}
}
func ExecutionTimeout(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc, error) {
	runtime := ExecutionRuntimeFromContext(ctx)
	if runtime == nil {
		return ctx, func() {}, ErrRuntimeNotConfigured
	}
	next, cancel := runtime.WithTimeout(ctx, duration)
	return next, cancel, nil
}
func ElapsedMillis(now, started time.Time) int64 {
	if now.IsZero() || started.IsZero() || now.Before(started) {
		return 0
	}
	return now.Sub(started).Milliseconds()
}
