# Specification: MCP Server (L2 Shell)

## Purpose
Expose Bateia architectural search and benchmark capabilities via the Model Context Protocol (MCP) using stdio JSON-RPC 2.0.

## Rules
- L2 package (`02_shell/mcp`).
- Reads JSON-RPC 2.0 requests from `io.Reader` (stdio stdin).
- Writes JSON-RPC 2.0 responses to `io.Writer` (stdio stdout).
- Invokes `usecase.SearchFunnel` ($L_1$) to perform queries.
- Exposes tools: `bateia_search` and `bateia_benchmark`.
