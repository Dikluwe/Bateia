# Análise de insights com observações específicas

**Camada:** L1. **Contrato:** passo 0134.

Implementar InvestigateSourceInsights para pedido v2, incluindo validação e canonicalização exatas do passo. Usar o classificador existente por InvestigateInsights sobre cópia v1 sem observações; não copiar política de evidência. Relatório v2 tem definição CATALOG_TEMPLATE, cobertura separada e declarações localizadas DECLARED_NOT_REASSESSED, incluindo trechos exatos de referências validadas. Conservar observações sem mecanismo com status/gap INCONCLUSIVE; não produzir hipóteses a partir de observações sozinhas. Nas hipóteses existentes projetar apenas pré-condições/trade-offs declarados da mesma fonte/classe e suas referências. Não interpretar texto por palavras-chave nem aumentar força de classificação/hipótese/alvo. IDs v2 independentes do legado e sensíveis ao pedido canônico/catálogo/contrato. Sem I/O, relógio, alteração do catálogo ou do pedido do chamador.
