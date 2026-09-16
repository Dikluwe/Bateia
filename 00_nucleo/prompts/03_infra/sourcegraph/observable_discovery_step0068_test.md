# Prompt proprietário — `03_infra/sourcegraph/observable_discovery_step0068_test.go`

## Obrigação do oráculo 0068

Verificar por servidor controlado que somente provenance `technical-lexicon-v1` válida
ativa uma tentativa técnica com linguagem, escopo e `count` limitado, sem fallback.
Admitir apenas snippets retidos do mesmo arquivo que contenham os dois discriminadores;
rejeitar split-file, metadata, path, Repository match e apenas um discriminador. Exigir
Coverage, QueryPlan, diagnósticos, rejeições e AdmissionEvidence coerentes; distinguir
vazio de erro/timeout/cancelamento e preservar o comportamento não técnico.

Este arquivo é origem L0 somente do teste nomeado. O contrato produtivo permanece em
`00_nucleo/prompts/03_infra/sourcegraph/observable_discovery.md`.
