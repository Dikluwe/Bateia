# Prompt proprietário — `01_core/usecase/semantic_expander_test.go`

Origem L0 exclusiva de `01_core/usecase/semantic_expander_test.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/usecase/semantic_expander.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `Expand`, `proposal`, `TestEnrichExpansionPlanPreservesLiteralAndInput`, `TestEnrichExpansionPlanDeduplicatesAndMergesTraceability`, `TestEnrichExpansionPlanRejectsMissingProvenance`, `TestEnrichExpansionPlanUsesNonGoalsAsBarrier`, `TestEvaluateExpansionStrategyIsDeterministicAndHandlesEmptyLists`, `TestEvaluateExpansionStrategyMeasuresCoverageContaminationAndNovelty`, `semanticExpanderFunc`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/usecase/semantic_expander.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.
