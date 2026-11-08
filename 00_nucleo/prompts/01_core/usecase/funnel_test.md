# Prompt proprietário — `01_core/usecase/funnel_test.go`

Origem L0 exclusiva de `01_core/usecase/funnel_test.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/usecase/funnel.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `SearchCandidates`, `FetchManifest`, `FetchTree`, `ExtractRepoFiles`, `FetchRepoMetadata`, `NewMockCacheStore`, `GetDiscovery`, `SaveDiscovery`, `GetManifest`, `SaveManifest`, `GetExtracted`, `SaveExtracted`, `ReportProgress`, `TestSearchFunnelExecution`, `TestSearchFunnelPolyglotPreservesPartialFailures`, `TestSearchFunnelDiscardsRepositoryWithoutSupportedManifest`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/usecase/funnel.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

## Passo 0035

Verificar injeção da porta semântica local e preservação de cobertura insuficiente/saúde não avaliada mesmo com resolução estática parcial.

## Instante de coleta — 0128

Verificar que manifesto reanalisado conserva o instante da extração fornecido por uma porta controlada, incluindo zero desconhecido. Não comparar a telemetria temporal do funil, que permanece fora deste recorte.
