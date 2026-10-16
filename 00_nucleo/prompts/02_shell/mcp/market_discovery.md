# Market discovery MCP tools

**Camada:** L2  
**Arquivo proprietário:** `02_shell/mcp/market_discovery.go`

Adicione somente as tools 8 e 9, na ordem, com schemas, envelopes e validação
runtime exatos do contrato 0027. Preserve por deep equality as sete tools
certificadas e estenda somente a enum de artifacts contratada. Dependa de portas
L1 injetadas; `list` é estático e sem rede. Rejeite unknown fields antes da porta,
serialize um único item textual JSON compacto e mapeie erros somente às mensagens
públicas contratadas. Não importe L3/L4 nem implemente store, fila, rede, hash,
merge, budget ou saneamento semântico.


## Correção de segurança — 2026-09-26

A extensão agora registra handlers no dispatcher comum; não usa readers que
interceptam chamadas. Todas as ferramentas/artefatos atravessam o orçamento
comum exatamente uma vez e recebem o contexto da sessão. Framing limitado é
aplicado antes de qualquer parsing. Ver `request_dispatch.md` para a supersessão.

## Projeção do registro no passo 0104

A lista pública de ferramentas deve usar o registro de descoberta injetado, inclusive quando síntese estiver habilitada. Projetar enum dos IDs configurados e limite de seleção de até dez. Preservar schema basal quando o registro contiver exatamente Apple/GitHub/Sourcegraph e comportamento legado quando não houver registro. O dispatcher continua único, com orçamento/contexto existentes; não adicionar parsing ou política de negócio.
