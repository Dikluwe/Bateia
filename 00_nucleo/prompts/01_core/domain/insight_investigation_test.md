# Contrato JSON da investigação de insights

**Camada:** L1  
**Arquivo proprietário:** `01_core/domain/insight_investigation_test.go`  
**Contrato:** passo 0105

Verificar a interoperabilidade do relatório JSON: estados de relatório, mecanismo e hipótese permanecem distintos após roundtrip. A origem declarada das avaliações e identidade da fonte não se perdem. O contrato independente não acrescenta run_id, decisão ou política de reúso.
