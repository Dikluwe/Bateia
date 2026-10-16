# Investigar mecanismos para insights

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/insight_investigation.go`  
**Contrato:** passo 0105

Validar e normalizar pedido limitado; commits imutáveis e bytes/hash de todas as fontes, sem assumir autenticidade externa. Usar extrator determinístico e catálogo ativo diretamente, sem RecommendationService, runtime ou decisão de reúso. Selecionar classes pedidas e gerar hipótese de referência por mecanismo com alvo UNKNOWN e nunca SUPPORTED. Explicitar baseline, compatibilidade e efeito não verificados; planejar baseline, aplicabilidade e experimento falsificável. Licença declarada não filtra valor investigativo nem autoriza adoção. Hashes causais incluem contexto/pergunta/fontes/catálogo. Resultado sem mecanismos é INCONCLUSIVE. Preservar pedido e comportamento da recomendação existente.
