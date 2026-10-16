# Registro configurado de descoberta

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/market_registry.go`  
**Contrato:** passo 0104

Definir NewMarketDiscoveryRegistry como admissão pura para ativar descoberta: exigir pelo menos um conector, delegar validação/canonicalização/IDs únicos a NewDiscoveryConnectorRegistry e não fixar número de fontes ou identidades de fornecedores. Preservar o construtor genérico, que ainda pode representar registro vazio. Retornar ErrInvalidRunInput para configuração inválida; nenhum I/O ou dependência de shell/infra/wiring.
