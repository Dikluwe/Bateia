# Market discovery domain

**Camada:** L1  
**Arquivo proprietário:** `01_core/domain/market_discovery.go`

Materialize somente os tipos, enums, portas e tags JSON públicos da revisão ativa
de `00_nucleo/benchmark/0027/contract.md`: capabilities/registry, query proposal,
produto, identidade, origem, repo link, tentativa, gap, resolução e cinco
artefatos. Preserve arrays vazios e as quatro dimensões independentes. Zero I/O,
MCP, filesystem, rede, env, goroutine ou mutex. Não invente política além do
contrato nem modele aplicação como `owner/repo`.


## 0126

Adicionar a constante ConnectorAtomGitCatalog com ID atomgit-repository-catalog, sem alterar formatos JSON, enums existentes ou política de investigação.

## 0154 — Descoberta contextual
QueryPolicy opcional omitida quando vazia no pedido do conector. Bridges podem aplicar política explícita de investigação sem mudar o pedido de mercado legado.
