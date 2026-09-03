# Oráculo independente L2 — tools MCP 0027

**Camada:** L2  
**Consumidor único:** `02_shell/mcp/market_discovery_step0027_oracle_test.go`

Congele por hash os arquivos certificados das tools 1–7; compare 1–5 e 7 por
deep equality e a tool 6 por deep equality após remover exatamente os cinco
sufixos aditivos contratados de `artifact.enum`. Exija apenas as tools 8–9 no
final com schemas e envelopes literais.
Listagem não chama Discover. Decoding fechado rejeita unknown fields, JSON
múltiplo, controles, duplicatas e valores fora do contrato antes da porta. Erros
públicos não ecoam canários. L2 não importa L3/L4.
