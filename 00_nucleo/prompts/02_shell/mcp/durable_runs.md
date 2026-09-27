# Durable recommendation MCP tools

Estenda o servidor MCP L2 com exatamente as quatro tools, schemas, envelopes e
mapeamento de erros de `00_nucleo/benchmark/0026/contract.md`. Preserve por
deep-equality as três tools do Passo 0025 e anuncie exatamente sete na ordem
contratada.

Dependa apenas da pequena porta L1 de start/get/artifact/cancel injetada. Faça
decode estrito com unknown fields rejeitados e validação runtime antes da porta.
Serialize JSON compacto no único item textual MCP; entries de artifact são valores
JSON, não strings. Nunca importe L3/L4, abra arquivo, crie worker, faça hash,
transição, paginação ou saneamento semântico. Erros físicos e canários viram apenas
as mensagens públicas exatas.

Este prompt é refinamento aditivo do servidor certificado; não edite o prompt
0025. Alterações mínimas em `server.go` e novo módulo de durable tools devem apontar
para esta causa e preservar a superfície síncrona anterior.

## Refinamento aditivo 0027

Preserve as sete tools e todos os erros 0026. Permita `run_kind` e status de
mercado na projeção de get-run e acrescente somente os cinco nomes contratados à
enum de artifacts. A autoria das tools 8–9 pertence a `market_discovery.md`.

## Refinamento 0031 — defaults compartilhados

Derive o schema de start do schema síncrono refinado: somente
`idempotency_key` e `idea` são obrigatórios. Use o mesmo materializador de
`policy`, `target` e `search` do gateway síncrono, preserve valores explícitos e
rejeite `null`, campos desconhecidos e valores inválidos antes de `Start`.

## Correção de segurança — 2026-09-26

A extensão agora registra handlers no dispatcher comum; não usa readers que
interceptam chamadas. Todas as ferramentas/artefatos atravessam o orçamento
comum exatamente uma vez e recebem o contexto da sessão. Framing limitado é
aplicado antes de qualquer parsing. Ver `request_dispatch.md` para a supersessão.
