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

## Correção de segurança — 2026-09-26

A extensão agora registra handlers no dispatcher comum; não usa readers que
interceptam chamadas. Todas as ferramentas/artefatos atravessam o orçamento
comum exatamente uma vez e recebem o contexto da sessão. Framing limitado é
aplicado antes de qualquer parsing. Ver `request_dispatch.md` para a supersessão.

## Projeção do registro no passo 0104

A lista pública de ferramentas deve usar o registro de descoberta injetado, inclusive quando síntese estiver habilitada. Projetar enum dos IDs configurados e limite de seleção de até dez. Preservar schema basal quando o registro contiver exatamente Apple/GitHub/Sourcegraph e comportamento legado quando não houver registro. O dispatcher continua único, com orçamento/contexto existentes; não adicionar parsing ou política de negócio.
