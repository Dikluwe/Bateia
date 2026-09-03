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

