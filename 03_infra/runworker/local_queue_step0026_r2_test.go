/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runworker/local_queue_step0026_r2_test.md
 * @layer L3
 * @updated 2026-09-01
 */
package runworker

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func step0026Commit(t *testing.T, queue *LocalQueue, id domain.RunID) {
	t.Helper()
	reservation, err := queue.Reserve(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if err := reservation.Commit(id); err != nil {
		t.Fatal(err)
	}
}

func TestStep0026R2QueuedCancellationWinsWithoutStartingHandler(t *testing.T) {
	blockerStarted := make(chan struct{})
	releaseBlocker := make(chan struct{})
	var targetCalls atomic.Int32
	blocker := domain.RunID("run_00000000000000000000000000000001")
	target := domain.RunID("run_00000000000000000000000000000002")
	queue, err := NewLocalQueue(1, 4, func(ctx context.Context, id domain.RunID) error {
		if id == target {
			targetCalls.Add(1)
			return nil
		}
		close(blockerStarted)
		select {
		case <-releaseBlocker:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Start(); err != nil {
		t.Fatal(err)
	}
	step0026Commit(t, queue, blocker)
	<-blockerStarted
	step0026Commit(t, queue, target)
	if err := queue.Cancel(context.Background(), target); err != nil {
		t.Fatal(err)
	}
	close(releaseBlocker)
	time.Sleep(20 * time.Millisecond)
	if targetCalls.Load() != 0 {
		t.Fatalf("queued-cancelled handler started %d times", targetCalls.Load())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := queue.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestStep0026R2RunningCancellationIsDeliveredExactlyOnce(t *testing.T) {
	id := domain.RunID("run_11111111111111111111111111111111")
	entered := make(chan struct{})
	observed := make(chan struct{})
	var observations atomic.Int32
	queue, err := NewLocalQueue(1, 2, func(ctx context.Context, got domain.RunID) error {
		if got != id {
			t.Fatalf("handler id = %s", got)
		}
		close(entered)
		<-ctx.Done()
		observations.Add(1)
		close(observed)
		return ctx.Err()
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Start(); err != nil {
		t.Fatal(err)
	}
	step0026Commit(t, queue, id)
	<-entered
	var wg sync.WaitGroup
	errs := make(chan error, 16)
	for index := 0; index < 16; index++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- queue.Cancel(context.Background(), id)
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("idempotent running cancel: %v", err)
		}
	}
	select {
	case <-observed:
	case <-time.After(time.Second):
		t.Fatal("worker did not observe cancellation")
	}
	if observations.Load() != 1 {
		t.Fatalf("cancellation observations=%d", observations.Load())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := queue.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestStep0026R2WorkerCountsBoundThirtyTwoJobsAndGoroutines(t *testing.T) {
	for _, workers := range []int{1, 4} {
		t.Run(string(rune('0'+workers)), func(t *testing.T) {
			baseline := runtime.NumGoroutine()
			release := make(chan struct{})
			entered := make(chan struct{}, 32)
			var active, maximum atomic.Int32
			queue, err := NewLocalQueue(workers, 32, func(ctx context.Context, _ domain.RunID) error {
				current := active.Add(1)
				defer active.Add(-1)
				for {
					old := maximum.Load()
					if current <= old || maximum.CompareAndSwap(old, current) {
						break
					}
				}
				entered <- struct{}{}
				select {
				case <-release:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			})
			if err != nil {
				t.Fatal(err)
			}
			if err := queue.Start(); err != nil {
				t.Fatal(err)
			}
			for index := 0; index < 32; index++ {
				step0026Commit(t, queue, domain.RunID("run_"+fmt32(index+1)))
			}
			for index := 0; index < workers; index++ {
				<-entered
			}
			if maximum.Load() != int32(workers) {
				t.Fatalf("active maximum=%d, workers=%d", maximum.Load(), workers)
			}
			if growth := runtime.NumGoroutine() - baseline; growth > workers+3 {
				t.Fatalf("goroutine-per-start growth=%d workers=%d", growth, workers)
			}
			close(release)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := queue.Shutdown(ctx); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStep0026R2BoundedShutdownDoesNotPresumeAcknowledgement(t *testing.T) {
	entered := make(chan struct{})
	ignoreCancellation := make(chan struct{})
	queue, err := NewLocalQueue(1, 1, func(context.Context, domain.RunID) error {
		close(entered)
		<-ignoreCancellation
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Start(); err != nil {
		t.Fatal(err)
	}
	step0026Commit(t, queue, "run_99999999999999999999999999999999")
	<-entered
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if err := queue.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("unacknowledged shutdown presumed success: %v", err)
	}
	close(ignoreCancellation)
}

func fmt32(value int) string {
	const hex = "0123456789abcdef"
	result := make([]byte, 32)
	for index := range result {
		result[index] = '0'
	}
	for index := 31; value > 0 && index >= 0; index-- {
		result[index] = hex[value&15]
		value >>= 4
	}
	return string(result)
}
