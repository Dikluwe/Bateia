# Prompt proprietário — `03_infra/github/client_test.go`

Origem L0 exclusiva de `03_infra/github/client_test.go`. Este contrato foi separado de `00_nucleo/prompts/03_infra/github/client.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `TestClient_FetchRepoMetadata`, `TestClient_FetchRepoMetadata_ServerError`, `TestClient_FetchTree`, `TestClientFetchTreeRecursiveAtImmutableCommit`, `TestClientFetchTreeDoesNotInventFallbackOnFailure`, `TestClient_FetchManifest_GoMod`, `TestClient_FetchManifest_NotFound`, `TestClient_FetchManifest_UsesRootListingToAvoidManifestProbeStorm`, `TestClient_AuthorizationHeader`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/03_infra/github/client.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

## Instante de coleta — 0128

Verificar que FetchManifest continua retornando instante de coleta não zero e UTC registrado pelo adaptador, mesmo com parser puro.
