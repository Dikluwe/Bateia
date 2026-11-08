# Market discovery use-case companion tests

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/market_discovery_test.go`

Verifique que o registry de descoberta é canônico e puramente estático, e que
um conector global preserva `MARKET_NOT_FILTERED` sem copiar locale para mercado.

## Descoberta global — 0126

Tentativas MARKET_NOT_FILTERED com PlannedCost positivo devem consultar a porta, preservando o aviso de ausência de filtro geográfico no plano. Tentativas sem orçamento não executam. Projeção separa status operacional da lacuna MARKET_NOT_FILTERED e nunca adiciona o mercado solicitado a ObservedMarkets se MarketFilterApplied for falso. Manter enums, hashes de requisição, limites e critérios de reúso. Testar execução durável com porta falsa, orçamento e publicação dos artefatos, além da projeção pura.
