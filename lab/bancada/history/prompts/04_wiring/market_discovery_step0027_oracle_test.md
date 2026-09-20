# Oráculo independente L3/L4 — Apple, replay e composição 0027

**Camada:** L4  
**Consumidor único:** `04_wiring/market_discovery_step0027_oracle_test.go`

Por APIs públicas, componha o conector Apple L3 com transporte capturável, replay
sem transporte e o runtime L4. Exija endpoint e parâmetros allowlisted, escape,
parsing estrito, identidade/dimensões, SSRF/injection/error/empty/cancelamento,
replay byte-idêntico, artefatos seguros e reutilização por identidade das portas
duráveis 0026. L4 contém apenas composição.

