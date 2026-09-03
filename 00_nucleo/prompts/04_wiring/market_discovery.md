# Market discovery wiring

**Camada:** L4  
**Arquivo proprietário:** `04_wiring/market_discovery.go`

Componha os dois conectores repo-backed basais, o único conector Apple, replay,
registry, serviço L1, dispatcher e tools L2. Injete EXATAMENTE as mesmas instâncias
de store, fila, gerador de run ID, raiz e cursor secret de 0026; preserve o mesmo
`RecommendationService`. Apenas composição/configuração: nenhuma validação,
budget, hash, merge, transição, retry, parsing ou regra de negócio em L4.

