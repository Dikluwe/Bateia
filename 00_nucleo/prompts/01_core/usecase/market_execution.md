# market_execution

Possuir investigação de mercado sem store/fila: capacidades selecionadas, planejamento, budget, consultas, normalização e projeção. Implementar ObservableMarketDiscoveryExecutor com progresso nos mesmos estágios e resultado detalhado para publicação pelo runtime. Cancelamento e erro de observação interrompem novas chamadas. Preservar lacunas geográficas, status e artefatos.
Contrato 0131/ADR-0014.

## Consolidação 0131
A entrada direta também é normalizada antes do dispatch, usando o mesmo validador já utilizado pela admissão durável. Não depender de efeitos do start para preencher QueryID. A porta configurada é executada; não é apenas um campo de configuração.
