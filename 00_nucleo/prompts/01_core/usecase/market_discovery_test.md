# Market discovery use-case companion tests

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/market_discovery_test.go`

Verifique que o registry de descoberta é canônico e puramente estático, e que
um conector global preserva `MARKET_NOT_FILTERED` sem copiar locale para mercado.
