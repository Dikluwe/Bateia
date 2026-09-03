# Prompt proprietário — `02_shell/mcp/server_test.go`

Origem L0 exclusiva de `02_shell/mcp/server_test.go`. Este contrato foi separado de `00_nucleo/prompts/02_shell/mcp/server.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `TestMCPServer_Initialize`, `TestMCPServer_ToolsList`, `TestMCPServer_RejectsEmptySearchQuery`, `TestMCPServer_RejectsMalformedBenchmarkRepository`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- após o Passo 0025, `tools/list` deve anunciar exatamente, inclusive quando o recommender
  opcional for `nil`, `bateia_search`, `bateia_benchmark` e `bateia_recommend`, nessa ordem;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/02_shell/mcp/server.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.
