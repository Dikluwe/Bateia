# Recommendation domain

Defina os contratos JSON versionados e puros do gateway de recomendação síncrona:
`RecommendationSearch`, `RecommendationRequest`, `RecommendationSummary` e
`RecommendationResult`, exatamente conforme
`00_nucleo/benchmark/0025/contract.md`. Reuse `IdeaSpec`,
`ReuseDecisionPolicy`, `SoftwareAlternativeTarget`, `IdeaRunStatus`,
`ReuseDecisionReport`, `SoftwareAlternativeReport` e `ExecutionFailure`; não
duplique esses tipos. A entrada não contém provider, modo, token, URL, comando,
cabeçalho, caminho, replay, output, resume ou run id. A saída não contém bundles
nem conteúdo bruto coletado. Use tags JSON snake_case e arrays públicos vazios em
vez de `null`.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
