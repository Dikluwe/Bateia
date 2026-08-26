# Prompt proprietário — `02_shell/cli/idea_plan_test.go`

Origem L0 exclusiva de `02_shell/cli/idea_plan_test.go`. Este contrato foi separado de `00_nucleo/prompts/02_shell/cli/idea_plan.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `TestRunIdeaPlanMatchesGoldenFile`, `TestRunIdeaPlanWritesOutputFileWithoutStdout`, `TestRunIdeaPlanRejectsInvalidJSON`, `TestRunIdeaPlanRejectsEmptyIdea`, `TestRunIdeaPlanRequiresInput`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/02_shell/cli/idea_plan.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.
