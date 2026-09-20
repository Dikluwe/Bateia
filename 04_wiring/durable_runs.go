/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/durable_runs.md
 * @layer L4
 * @updated 2026-09-19
 */

package main

import (
	"context"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/03_infra/runcontainer"
	"bateia/03_infra/runstore"
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
	durable, err := runcontainer.New(executor, root, cursorSecret, limits, workers, capacity)
	if err != nil {
		return nil, nil, err
	}
	service := durable.Service()
	if err := durable.Recover(context.Background()); err != nil {
		return nil, nil, err
	}
	if err := durable.Start(); err != nil {
		return nil, nil, err
	}
	server := mcp.NewMCPServer(funnel, reader, writer, executor).WithDurableRuns(service)
	return server, durable.Shutdown, nil
}

func newDurableRecommendationRuns(
	executor domain.ObservableRecommendationExecutor,
	root string,
	cursorSecret []byte,
	limits runstore.RunStoreLimits,
	workers int,
	capacity int,
) (*usecase.RecommendationRunService, func(context.Context) error, error) {
	durable, err := runcontainer.New(executor, root, cursorSecret, limits, workers, capacity)
	if err != nil {
		return nil, nil, err
	}
	service := durable.Service()
	if err := durable.Recover(context.Background()); err != nil {
		return nil, nil, err
	}
	if err := durable.Start(); err != nil {
		return nil, nil, err
	}
	return service, durable.Shutdown, nil
}
