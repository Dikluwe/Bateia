# Comparação consultiva de observações

**Camada:** L1. **Contrato:** passo 0135.

Implementar CompareInsightBaseline com pedido v1; delegar investigação a InvestigateSourceInsights v2 e baseline a normalizador próprio. Projetar análises/proveniência da fonte sem reescrever IDs. Agrupar observações por fonte/classe/tipo com status de disponibilidade, referências separadas e próximo passo que considera política já declarada antes de mudança. Emitir linhas mesmo sem mecanismo da fonte. Não emitir hipóteses ou veredictos de presença/ausência/equivalência/efeito. IDs canônicos incluem investigação, catálogo e baseline normalizada; não mutar caller. Sem I/O, relógio, avaliação por keywords ou execução.
