# 01_core/domain/execution_runtime

**Camada:** L1
**Proprietário:** `01_core/domain/execution_runtime.go`
**Contrato:** passo 0129 e ADR-0012.

Porta de relógio e deadline; contexto transporta runtime explicitamente configurado. RuntimeNow sem porta retorna instante desconhecido. ElapsedMillis é transformação pura com resultado não negativo. Nenhum efeito concreto.
