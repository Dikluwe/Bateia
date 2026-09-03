/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/durable_mcp.md
 * @layer L2
 * @updated 2026-09-01
 */

package cli

import "bateia/02_shell/mcp"

// WithDurableRuns conecta à casca CLI uma porta já composta. A casca não
// conhece store, fila, workers, quotas nem qualquer adapter de infraestrutura.
func (a *CLIApp) WithDurableRuns(runs mcp.DurableRecommendationRuns) *CLIApp {
	if a != nil {
		a.durableRuns = runs
	}
	return a
}
