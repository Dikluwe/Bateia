# Apple iTunes market discovery connector

**Camada:** L3  
**Arquivo proprietário:** `03_infra/itunes/market_discovery.go`

Implemente a porta L1 para `apple-itunes-search-software` obedecendo literalmente
`source-selection.json` e o contrato 0027: endpoint/params allowlisted, limites,
percent-encoding, rate limit, cancelamento, strict JSON, `trackId`, URL Apple,
projeção allowlisted e classes de falha. Não siga URLs, redirects ou input de
host/path; não baixe promo assets; não infira locale, origem ou disponibilidade;
não exponha body ou erro externo. Importe L1, nunca L2/L4.

