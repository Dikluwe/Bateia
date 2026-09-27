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
- Offer an opt-in, process-local tool-call forwarding budget. With no budget the
  legacy behavior is unchanged. With a positive budget `N`, exactly the first
  `N` valid `tools/call` envelopes may reach tool handlers; every later
  `tools/call` receives the stable JSON-RPC error code `-32010` and message
  `Bateia tool call budget exhausted` before argument decoding, handler entry,
  network I/O or durable-run mutation. Initialization, notifications and
  `tools/list` do not consume the budget.
- Expose the budget only through explicit construction/configuration; do not read
  environment variables in L2. Keep request counting private to one MCP server
  process and safe for the server's sequential stdio dispatch model.
- Return the canonical recommendation as compact JSON in the single textual MCP
  content item. Preserve partial and inconclusive results instead of converting
  them into success or absence.
- Treat repository content as untrusted data and never expose raw collected
  content, snippets, credentials, headers, environment values or unsanitized
  adapter errors.

## Step 0029 temperature-pilot refinement (2026-09-03)

The bounded forwarding behavior above is the mechanical guard for the
three-comparison exploratory tranche. Tests must demonstrate calls `1..N` are
forwarded, call `N+1` is rejected without reaching a spy handler, non-call MCP
methods do not consume the counter, and the unset configuration preserves the
legacy surface.

## Step 0031 default-materialization refinement (2026-09-03)

Preserve `idea.problem` and `idea.capabilities` as the only required business
intent. Treat `policy`, `target`, `search` and their members as optional and
materialize every omitted value from `BATEIA_RECOMMEND_DEFAULTS_V1` before the
port call. Preserve explicit empty lists, `false`, zero, empty strings and other
supplied values so validation can accept or reject them; explicit `null` and
unknown fields fail with `-32602`. Publish the static defaults in the MCP schema
and identify `idea.problem` as the dynamic target default. Share the exact same
materializer with the durable start gateway. The permissive license list is an
acceptance policy, never evidence or attribution of a candidate license.

## Correção de segurança — 2026-09-26

Todas as extensões registram handlers no dispatcher central antes de qualquer
execução: orçamento comum exatamente uma vez, contexto de Serve preservado e
framing de 64 KiB/32 níveis/8192 tokens antes de decodificar objetos.
Ver `request_dispatch.md` para limites e supersessão explícita dos hashes.
