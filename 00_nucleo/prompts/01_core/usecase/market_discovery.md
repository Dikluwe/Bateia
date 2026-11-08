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


## Descoberta global — 0126

Tentativas MARKET_NOT_FILTERED com PlannedCost positivo devem consultar a porta, preservando o aviso de ausência de filtro geográfico no plano. Tentativas sem orçamento não executam. Projeção separa status operacional da lacuna MARKET_NOT_FILTERED e nunca adiciona o mercado solicitado a ObservedMarkets se MarketFilterApplied for falso. Manter enums, hashes de requisição, limites e critérios de reúso. Testar execução durável com porta falsa, orçamento e publicação dos artefatos, além da projeção pura.

## Passo 0129 — ownership de capacidade
O serviço durável conserva store, fila, IDs, execução e lifecycle único. Containers internos separados possuem registro de descoberta, configuração de extração e admissão de síntese. A extração e ordenação dos mecanismos pertence à capacidade; o runtime publica suas saídas sem reinterpretá-las. APIs Configure, hashes, artefatos e transições permanecem compatíveis.

## Passo 0131
Aplicar o ownership do ADR-0014 e passo 0131. Jornada/decisões em L1, transporte/renderização em L2, efeitos em L3, composição em L4. Preservar os contratos observáveis, formatos, hashes, versões, ordem, retomada e cancelamento; testes seguem o proprietário ou Lab para auditoria geral. Não conservar regra duplicada na casca/runtime.
