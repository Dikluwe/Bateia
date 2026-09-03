# Oráculo independente L1 — market discovery 0027

**Camada:** L1  
**Consumidor único:** `01_core/usecase/market_discovery_step0027_oracle_test.go`

Teste somente APIs públicas do contrato 0027. Use um modelo de referência cego
para validar registry, matriz locale×market×connector, controles cruzados,
ordenação causal, orçamento round-robin, identidade conservadora, quatro
dimensões, projeções seguras, idempotência e lifecycle compartilhados. Nenhum I/O,
rede, filesystem, mutex ou goroutine em produção L1. Unknown nunca é sucesso.

