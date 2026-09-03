# Durable recommendation run domain

Defina em L1, sem I/O, os contratos versionados do ciclo operacional assíncrono
congelados em `00_nucleo/benchmark/0026/contract.md`: `RunID`, lifecycle separado
de `IdeaRunStatus`, estágios, progresso, registro persistível, projeção pública,
idempotência, transição CAS, índice de artifact, publicação e página.

Preserve exatamente enums, tags JSON, nullabilidade e arrays vazios do contrato.
`run_id` é opaco no formato `run_` + 32 hex minúsculos. Estados terminais são
imutáveis; `SUCCEEDED` não implica `COMPLETE`. Não duplique
`RecommendationRequest`, `RecommendationResult`, `ReuseDecisionReport`,
`SoftwareAlternativeReport`, `ExecutionFailure` ou `EvidenceSource`.

Este módulo não conhece diretório, arquivo, goroutine, mutex, MCP, segredo,
cursor físico, adapter ou variável de ambiente.

## Refinamento aditivo 0027

Refine este mesmo arquivo para `RunKind`, `MARKET_DISCOVERY`, estágios e status de
mercado exatamente como `00_nucleo/benchmark/0027/contract.md`. Preserve todos os
valores 0026; recommendation e market status são mutuamente exclusivos na
projeção pública. Store, CAS, lifecycle, run ID e artefatos continuam um único
domínio durável.
