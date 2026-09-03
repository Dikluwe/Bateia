/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runworker/local_queue.md
 * @layer L3
 * @updated 2026-09-01
 */

package runworker

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"bateia/01_core/domain"
)

type CryptoRunIDGenerator struct{}

func (CryptoRunIDGenerator) NewRunID() (domain.RunID, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return domain.RunID("run_" + hex.EncodeToString(raw)), nil
}

type jobState uint8

const (
	jobQueued jobState = iota + 1
	jobRunning
	jobCancelled
)

type LocalQueue struct {
	workers  int
	capacity int
	handler  func(context.Context, domain.RunID) error

	mu       sync.Mutex
	reserved int
	states   map[domain.RunID]jobState
	cancels  map[domain.RunID]context.CancelFunc
	views    map[domain.RunID]domain.RecommendationRunView
	jobs     chan domain.RunID
	stop     chan struct{}
	started  bool
	stopped  bool
	wg       sync.WaitGroup
}

func NewLocalQueue(workers, capacity int, handler func(context.Context, domain.RunID) error) (*LocalQueue, error) {
	if workers < 1 || workers > 4 || capacity < 1 || capacity > 32 || handler == nil {
		return nil, fmt.Errorf("invalid local queue configuration")
	}
	return &LocalQueue{
		workers: workers, capacity: capacity, handler: handler,
		states: map[domain.RunID]jobState{}, cancels: map[domain.RunID]context.CancelFunc{}, views: map[domain.RunID]domain.RecommendationRunView{},
		jobs: make(chan domain.RunID, capacity), stop: make(chan struct{}),
	}, nil
}

type localReservation struct {
	queue *LocalQueue
	once  sync.Once
	used  bool
}

func (q *LocalQueue) Reserve(ctx context.Context) (domain.RecommendationQueueReservation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.stopped {
		return nil, domain.ErrQueueFull
	}
	if q.reserved+len(q.states) >= q.capacity {
		return nil, domain.ErrQueueFull
	}
	q.reserved++
	return &localReservation{queue: q}, nil
}

func (r *localReservation) Commit(id domain.RunID) error {
	var result error = fmt.Errorf("reservation already consumed")
	r.once.Do(func() {
		r.used = true
		q := r.queue
		q.mu.Lock()
		q.reserved--
		if q.stopped {
			q.mu.Unlock()
			result = domain.ErrQueueFull
			return
		}
		q.states[id] = jobQueued
		q.mu.Unlock()
		q.jobs <- id
		result = nil
	})
	return result
}

func (r *localReservation) Release() {
	r.once.Do(func() {
		r.used = true
		r.queue.mu.Lock()
		r.queue.reserved--
		r.queue.mu.Unlock()
	})
}

func (q *LocalQueue) Start() error {
	q.mu.Lock()
	if q.started || q.stopped {
		q.mu.Unlock()
		return fmt.Errorf("local queue already started or stopped")
	}
	q.started = true
	q.wg.Add(q.workers)
	q.mu.Unlock()
	for index := 0; index < q.workers; index++ {
		go q.worker()
	}
	return nil
}

func (q *LocalQueue) worker() {
	defer q.wg.Done()
	for {
		select {
		case <-q.stop:
			return
		case id := <-q.jobs:
			q.mu.Lock()
			state, exists := q.states[id]
			if !exists || state == jobCancelled {
				delete(q.states, id)
				q.mu.Unlock()
				continue
			}
			ctx, cancel := context.WithCancel(context.Background())
			q.states[id] = jobRunning
			q.cancels[id] = cancel
			q.mu.Unlock()

			_ = q.handler(ctx, id)
			cancel()
			q.mu.Lock()
			delete(q.cancels, id)
			delete(q.states, id)
			q.mu.Unlock()
		}
	}
}

func (q *LocalQueue) Cancel(ctx context.Context, id domain.RunID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	switch q.states[id] {
	case jobQueued:
		q.states[id] = jobCancelled
	case jobRunning:
		if cancel := q.cancels[id]; cancel != nil {
			cancel()
		}
	}
	return nil
}

func (q *LocalQueue) CancellationView(id domain.RunID) (domain.RecommendationRunView, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	view, ok := q.views[id]
	return view, ok
}

func (q *LocalQueue) RememberCancellationView(id domain.RunID, view domain.RecommendationRunView) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, exists := q.views[id]; !exists {
		q.views[id] = view
	}
}

func (q *LocalQueue) ForgetCancellationView(id domain.RunID) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.views, id)
}

func (q *LocalQueue) Shutdown(ctx context.Context) error {
	q.mu.Lock()
	if !q.stopped {
		q.stopped = true
		for _, cancel := range q.cancels {
			cancel()
		}
		close(q.stop)
	}
	q.mu.Unlock()
	done := make(chan struct{})
	go func() {
		q.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
