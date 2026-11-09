# Prompt proprietário — `01_core/usecase/evidence_test.go`

Origem L0 exclusiva de `01_core/usecase/evidence_test.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/usecase/evidence.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `testBundle`, `TestDocumentationAloneDoesNotProveCapability`, `TestWeakSignalNeverBecomesAbsence`, `TestImplementedSymbolAndRelatedTestCanProve`, `TestFailureNeverBecomesAbsent`, `TestVendorNoiseCannotProve`, `TestAssessmentIsDeterministicAndDoesNotMutateBundle`, `TestPlanAcceptsResolutionPendingAndSelectsDeterministically`, `TestPlanIncludesExplicitCalibrationControls`, `FetchManifest`, `FetchTree`, `ExtractRepoFiles`, `FetchRepoMetadata`, `TestCollectDeduplicatesByHashAndPreservesPartialFailure`, `TestMissingRelevantTargetMakesAbsenceInconclusive`, `TestCommentDistantFromDeclarationDoesNotProduceSymbol`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/usecase/evidence.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

## Contrato 0040

Correspondências lexicais, declarações e testes vazios são evidência parcial, não comprovação comportamental. As expectativas anteriores de PROVEN e componentes derivados apenas desses sinais foram substituídas. Preservar verificações de fontes, localização, determinismo, licença e rejeição; identidade de execução muda com a versão de regras.

## Correção contratual 0043
O caso de README declarando execução local exige PARTIAL (declaração), não PROVEN (execução). Renomear o teste para tornar explícito esse limite; preservar demais casos.

## 0158 — Double de conteúdo com identidade

No recorte de integridade0158, os doubles que retornam extração bem-sucedida devem preencher Ref com owner/repo/commit do candidato efetivamente solicitado. Conservar conteúdo, falhas, quotas, testemunhos, mutantes e assertivas funcionais anteriores. Ref ausente ou divergente passa a ser controle negativo explícito nos testes próprios de integridade, sem dispensar a validação produtiva.

## 0162 — Preparação explícita do controle de deduplicação

O cenário TestCollectDeduplicatesByHashAndPreservesPartialFailure deve declarar os targets README.md/COPY.md no plano: sua árvore falha deliberadamente, e lista vazia não deve iniciar extração sob o contrato0162. Conservar a falha de árvore, os bytes iguais e todas as assertivas de deduplicação/incompletude. Ajustar somente essa preparação; não restaurar extração vazia na produção para atender ao double legado. Origem ex-ante: [passo0162](../../../passos/0162-aplicar-quotas-e-coerencia-de-fontes.md).
