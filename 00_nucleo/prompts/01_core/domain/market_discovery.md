# Market discovery domain

**Camada:** L1  
**Arquivo proprietário:** `01_core/domain/market_discovery.go`

Materialize somente os tipos, enums, portas e tags JSON públicos da revisão ativa
de `00_nucleo/benchmark/0027/contract.md`: capabilities/registry, query proposal,
produto, identidade, origem, repo link, tentativa, gap, resolução e cinco
artefatos. Preserve arrays vazios e as quatro dimensões independentes. Zero I/O,
MCP, filesystem, rede, env, goroutine ou mutex. Não invente política além do
contrato nem modele aplicação como `owner/repo`.

