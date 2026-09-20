/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runcontainer/container.md
 * @layer L3
 * @updated 2026-09-19
 */

package runcontainer

import (
	"context"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/03_infra/runstore"
	"bateia/03_infra/runworker"
)

// Container owns the concrete process-local resources for durable runs.
type Container struct {
	service *usecase.RecommendationRunService
	queue   *runworker.LocalQueue
	store   *runstore.FileStore
}

// New constructs the durable resources without recovering runs or starting workers.
func New(
	executor domain.ObservableRecommendationExecutor,
	root string,
	cursorSecret []byte,
	limits runstore.RunStoreLimits,
	workers int,
	capacity int,
) (*Container, error) {
	store, err := runstore.NewFileStore(root, cursorSecret, limits)
	if err != nil {
		return nil, err
	}
	var service *usecase.RecommendationRunService
	queue, err := runworker.NewLocalQueue(workers, capacity, func(ctx context.Context, id domain.RunID) error {
		return service.Execute(ctx, id)
	})
	if err != nil {
		_ = store.Close()
		return nil, err
	}
	service = usecase.NewRecommendationRunService(store, queue, runworker.CryptoRunIDGenerator{}, executor)
	return &Container{service: service, queue: queue, store: store}, nil
}

// Service returns the single service instance wired to this container.
func (c *Container) Service() *usecase.RecommendationRunService {
	return c.service
}

// Recover interrupts non-terminal runs before workers are started.
func (c *Container) Recover(ctx context.Context) error {
	if err := c.service.Recover(ctx); err != nil {
		_ = c.store.Close()
		return err
	}
	return nil
}

// Start starts the local queue after recovery and feature configuration.
func (c *Container) Start() error {
	if err := c.queue.Start(); err != nil {
		_ = c.store.Close()
		return err
	}
	return nil
}

// Shutdown drains the queue before closing the store and preserves queue error precedence.
func (c *Container) Shutdown(ctx context.Context) error {
	queueErr := c.queue.Shutdown(ctx)
	storeErr := c.store.Close()
	if queueErr != nil {
		return queueErr
	}
	return storeErr
}
