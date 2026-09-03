# Improvement synthesis MCP tool

**Camada:** L2  
**Arquivo proprietário:** `02_shell/mcp/improvement_synthesis.go`

Adicione somente a tool 10 `bateia_start_improvement_synthesis`, no final, com
descrição, schema, envelope e erros exatos do contrato 0028. Preserve as posições
1–9; mudanças aditivas de `run_kind`, status e enum de artifacts pertencem apenas
às extensões expressamente contratadas.

Faça decode JSON estrito, rejeite campos desconhecidos e controles e invoque uma
porta L1 pequena. Validação semântica e verificação de source refs permanecem em
L1, antes da reserva da fila. Serialize um único item textual JSON compacto; não
importe L3/L4 nem implemente store, fila, hash, resolução, síntese ou saneamento
semântico. Não ecoe proposta, conteúdo externo, hash divergente ou canário em
erros.
