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
