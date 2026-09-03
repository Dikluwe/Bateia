# Oráculo independente da tool MCP 10 do Passo 0028

**Camada:** L2  
**Arquivo proprietário:** `02_shell/mcp/improvement_synthesis_step0028_oracle_test.go`

Receba do verificador capturas JSON-RPC de `tools/list`; não imponha construtor ou
nome privado ao servidor. Preserve as nove tools certificadas, aceitando somente
as extensões aditivas de get-run/get-artifact expressamente listadas no contrato
R3, e exija `bateia_start_improvement_synthesis` exatamente na posição dez.

Valide nome, descrição, schema estrito, ordem de required, enums, limites,
envelope textual compacto e erro `SOURCE_BASIS_UNAVAILABLE`. Tools/list não pode
iniciar run, reservar fila, ler source ou chamar rede. Erros não ecoam conteúdo,
hash divergente, path, proposta ou canário.
