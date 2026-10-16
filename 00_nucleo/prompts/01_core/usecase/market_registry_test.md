# Contrato do registro configurado

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/market_registry_test.go`  
**Contrato:** passo 0104

Verificar registro unitário e ampliado, canonicalização, resolução e isolamento das capacidades retornadas. Rejeitar registro vazio, conector nil, IDs duplicados e capacidades inválidas sem chamar Discover. Não substituir validação por contagem. Preservar comportamento do registro genérico vazio.
