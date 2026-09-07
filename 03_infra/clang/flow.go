/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/flow.md
 * @layer L3
 * @updated 2026-09-07
 */
package clang

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"os/exec"
)

func (c *Client) WithFlowTool(binary string) *Client { c.flowTool = binary; return c }
func (c *Client) analyzeFlow(ctx context.Context, entry compilation, flags []string) *domain.StaticFlow {
	result := &domain.StaticFlow{Status: "not_evaluated", Reason: "Auxiliar CFG não configurado."}
	if c.flowTool == "" {
		return result
	}
	binary, err := exec.LookPath(c.flowTool)
	if err != nil {
		result.Reason = err.Error()
		return result
	}
	before, err := readLimited(binary, 128*1024*1024)
	if err != nil {
		result.Reason = err.Error()
		return result
	}
	args := append([]string{entry.File}, flags...)
	raw, diagnostics, err := run(ctx, binary, entry.Directory, args, maxOutput)
	if err != nil {
		result.Reason = string(diagnostics) + "; " + err.Error()
		return result
	}
	after, err := readLimited(binary, 128*1024*1024)
	if err != nil || digest(before) != digest(after) {
		result.Reason = "Auxiliar mudou durante análise"
		return result
	}
	if err = json.Unmarshal(raw, result); err != nil || result.Schema != "bateia-flow-v1" || result.Status != "complete" {
		return &domain.StaticFlow{Status: "not_evaluated", Reason: "Contrato JSON do auxiliar inválido"}
	}
	result.Reason = ""
	result.ToolSHA256 = digest(before)
	result.Arguments = args
	return result
}
