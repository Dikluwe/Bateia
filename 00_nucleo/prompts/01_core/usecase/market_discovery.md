# Market discovery use case

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/market_discovery.go`

Implemente a orquestração pura contratada: validação anterior a portas, registry
estático, hash/idempotência canônicos, plano produto cartesiano, budget justo,
dispatcher por `run_kind`, execução durável, classificação Unknown, consolidação,
merge por testemunha forte e projeção dos cinco artefatos. Reutilize as portas de
store e fila do Passo 0026; não crie pipeline, store, queue ou idempotência
paralelos. Zero I/O, MCP, filesystem, env, goroutine, mutex ou texto externo em
erros públicos.

