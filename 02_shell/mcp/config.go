/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/config.md
 * @layer L2
 * @updated 2026-10-28
 */
package mcp

import (
	"fmt"
	"strconv"
)

func ParseMaxToolCalls(raw string, present bool) (*int, error) {
	if !present {
		return nil, nil
	}
	if raw == "" || raw[0] < '1' || raw[0] > '9' {
		return nil, fmt.Errorf("BATEIA_MCP_MAX_TOOL_CALLS must be a canonical positive decimal integer")
	}
	for index := 1; index < len(raw); index++ {
		if raw[index] < '0' || raw[index] > '9' {
			return nil, fmt.Errorf("BATEIA_MCP_MAX_TOOL_CALLS must be a canonical positive decimal integer")
		}
	}
	parsed, err := strconv.ParseInt(raw, 10, strconv.IntSize)
	if err != nil || parsed <= 0 {
		return nil, fmt.Errorf("BATEIA_MCP_MAX_TOOL_CALLS must be a canonical positive decimal integer")
	}
	limit := int(parsed)
	return &limit, nil
}
