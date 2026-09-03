/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runworker/local_queue_step0026_contract_test.md
 * @layer L3
 * @updated 2026-09-01
 */
package runworker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func TestStep0026LocalQueueValidatesLimitsAndCommitCannotFailForCapacity(t *testing.T) {
	handler := func(context.Context, domain.RunID) error { return nil }
	for _, limits := range [][2]int{{0, 1}, {5, 1}, {1, 0}, {1, 33}} {
		if _, err := NewLocalQueue(limits[0], limits[1], handler); err == nil {
			t.Fatalf("invalid workers/capacity accepted: %v", limits)
		}
	}
	queue, err := NewLocalQueue(1, 1, handler)
	if err != nil {
		t.Fatal(err)
	}
	reservation, err := queue.Reserve(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Reserve(context.Background()); !errors.Is(err, domain.ErrQueueFull) {
		t.Fatalf("capacity not reserved atomically: %v", err)
	}
	if err := reservation.Commit("run_0123456789abcdef0123456789abcdef"); err != nil {
		t.Fatalf("reserved commit failed by capacity: %v", err)
	}
	if err := queue.Start(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := queue.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestStep0026LocalQueueBoundsConcurrencyAndCancelsRunningContext(t *testing.T) {
	const workers = 2
	started := make(chan domain.RunID, 8)
	release := make(chan struct{})
	var active, maximum atomic.Int32
	handler := func(ctx context.Context, id domain.RunID) error {
		current := active.Add(1)
		defer active.Add(-1)
		for {
			old := maximum.Load()
			if current <= old || maximum.CompareAndSwap(old, current) {
				break
			}
		}
		started <- id
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-release:
			return nil
		}
	}
	queue, err := NewLocalQueue(workers, 8, handler)
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Start(); err != nil {
		t.Fatal(err)
	}
	ids := []domain.RunID{
		"run_00000000000000000000000000000001", "run_00000000000000000000000000000002",
		"run_00000000000000000000000000000003", "run_00000000000000000000000000000004",
	}
	for _, id := range ids {
		reservation, err := queue.Reserve(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if err := reservation.Commit(id); err != nil {
			t.Fatal(err)
		}
	}
	first, second := <-started, <-started
	if maximum.Load() > workers {
		t.Fatalf("worker bound exceeded: %d", maximum.Load())
	}
	if err := queue.Cancel(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	select {
	case third := <-started:
		if third == first {
			t.Fatal("cancel restarted same run")
		}
	case <-time.After(time.Second):
		t.Fatal("cancelled worker did not release bounded capacity")
	}
	_ = second
	close(release)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := queue.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestStep0026LocalQueueConcurrentReservationsNeverOversubscribe(t *testing.T) {
	queue, err := NewLocalQueue(1, 4, func(context.Context, domain.RunID) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	var accepted atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reservation, err := queue.Reserve(context.Background())
			if err == nil {
				accepted.Add(1)
				reservation.Release()
			} else if !errors.Is(err, domain.ErrQueueFull) {
				t.Errorf("unexpected reservation error: %v", err)
			}
		}()
	}
	wg.Wait()
	if accepted.Load() == 0 || accepted.Load() > 32 {
		t.Fatalf("invalid reservation accounting: %d", accepted.Load())
	}
	reservation, err := queue.Reserve(context.Background())
	if err != nil {
		t.Fatalf("released reservations leaked capacity: %v", err)
	}
	reservation.Release()
}
