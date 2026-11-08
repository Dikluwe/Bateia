# Prompt proprietário — `01_core/usecase/idea_discovery_test.go`

Origem L0 exclusiva de `01_core/usecase/idea_discovery_test.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/usecase/idea_discovery.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `Discover`, `discoveryPlan`, `candidate`, `TestDiscoverExpansionPlanDeduplicatesOriginsAndKeepsPartialFailure`, `TestDiscoverExpansionPlanFailsWhenAllExpansionsFail`, `TestDiscoverExpansionPlanDoesNotMutateInput`, `fakeObservableProvider`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/usecase/idea_discovery.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

## 0129

Cenários que pedem deadline total fornecem runtime de teste explicitamente. Manter resultado parcial, causa TOTAL_TIMEOUT e cancelamento; não remover as asserções por falta de host.
