# Specification: MCP Server (L2 Shell)

## Purpose

Expose Bateia architectural search, benchmark and complete recommendation
capabilities through MCP over stdio JSON-RPC 2.0.

## Rules

- Remain in L2 (`02_shell/mcp`), use L1 for domain/use-case contracts and never
  import L3 or L4.
- Read JSON-RPC 2.0 from `io.Reader` and write responses to `io.Writer`.
- Preserve the public behavior of `bateia_search` and `bateia_benchmark`.
- Announce `bateia_recommend` as the third tool, with the exact schema, validation,
  content envelope and error behavior frozen in
  `00_nucleo/benchmark/0025/contract.md`.
- Depend on the small `RecommendationRecommender` port. Do not call the CLI, read
  files, create temporary directories, rank in L2 or import L3.
- Reject unknown fields and invalid limits before invoking the recommender.
- Return the canonical recommendation as compact JSON in the single textual MCP
  content item. Preserve partial and inconclusive results instead of converting
  them into success or absence.
- Treat repository content as untrusted data and never expose raw collected
  content, snippets, credentials, headers, environment values or unsanitized
  adapter errors.
