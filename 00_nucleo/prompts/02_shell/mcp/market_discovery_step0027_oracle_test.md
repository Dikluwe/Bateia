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

## Sucessor histórico O3

O digest basal antigo de `server.go` permanece no artefato de invalidação; o
teste vincula o digest sucessor exato e os artefatos O3. A garantia não é
reduzida a hash: tools 1–5 e 7 continuam comparadas por igualdade profunda, tool
6 aceita somente os cinco sufixos aditivos já contratados e tools 8–9 permanecem
exatas. A extensão de budget não permite mudar ordem, schema ou efeitos.
## Sucessor explícito 0031

Aceite os novos hashes de `server.go` e `durable_runs.go` somente quando
vinculados a `historical-contract-supersession.json`. Preserve por deep-equality
a relação entre as sete tools basais e as extensões 8–9; a mudança autorizada é
apenas o perfil de defaults/required das recomendações.
