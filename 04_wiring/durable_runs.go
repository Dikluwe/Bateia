/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/durable_runs.md
 * @layer L4
 * @updated 2026-09-01
 */

package main

import (
	"context"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/03_infra/runstore"
	"bateia/03_infra/runworker"
)

func NewDurableMCPServer(
	funnel *usecase.SearchFunnel,
	executor domain.ObservableRecommendationExecutor,
	root string,
	cursorSecret []byte,
	limits runstore.RunStoreLimits,
	workers int,
	capacity int,
	reader io.Reader,
	writer io.Writer,
) (*mcp.MCPServer, func(context.Context) error, error) {
	service, queue, store, err := newDurableComponents(executor, root, cursorSecret, limits, workers, capacity)
	if err != nil {
		return nil, nil, err
	}
	if err := service.Recover(context.Background()); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	if err := queue.Start(); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	server := mcp.NewMCPServer(funnel, reader, writer, executor).WithDurableRuns(service)
	return server, durableShutdown(queue, store), nil
}

func newDurableRecommendationRuns(
	executor domain.ObservableRecommendationExecutor,
	root string,
	cursorSecret []byte,
	limits runstore.RunStoreLimits,
	workers int,
	capacity int,
) (*usecase.RecommendationRunService, func(context.Context) error, error) {
	service, queue, store, err := newDurableComponents(executor, root, cursorSecret, limits, workers, capacity)
	if err != nil {
		return nil, nil, err
	}
	if err := service.Recover(context.Background()); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	if err := queue.Start(); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	return service, durableShutdown(queue, store), nil
}

func newDurableComponents(
	executor domain.ObservableRecommendationExecutor,
	root string,
	cursorSecret []byte,
	limits runstore.RunStoreLimits,
	workers int,
	capacity int,
) (*usecase.RecommendationRunService, *runworker.LocalQueue, *runstore.FileStore, error) {
	store, err := runstore.NewFileStore(root, cursorSecret, limits)
	if err != nil {
		return nil, nil, nil, err
	}
	var service *usecase.RecommendationRunService
	queue, err := runworker.NewLocalQueue(workers, capacity, func(ctx context.Context, id domain.RunID) error {
		return service.Execute(ctx, id)
	})
	if err != nil {
		_ = store.Close()
		return nil, nil, nil, err
	}
	service = usecase.NewRecommendationRunService(store, queue, runworker.CryptoRunIDGenerator{}, executor)
	return service, queue, store, nil
}

func durableShutdown(queue *runworker.LocalQueue, store *runstore.FileStore) func(context.Context) error {
	return func(ctx context.Context) error {
		queueErr := queue.Shutdown(ctx)
		storeErr := store.Close()
		if queueErr != nil {
			return queueErr
		}
		return storeErr
	}
}
