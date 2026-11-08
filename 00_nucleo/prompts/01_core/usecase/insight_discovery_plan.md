# Plano de rotas de pesquisa

**Camada:** L1. **Contrato:** passo 0137.

Normalizar/validar pedido e produzir PlanInsightDiscovery com cinco rotas e IDs determinísticos conforme passo. Remover tokens familiares somente dos anchors exploratórios, preservar direta; não inventar tradução/domínio/negação. Validar quotas/durações antes de efeitos, copiar slices, não mutar entradas.

## 0154 — Descoberta contextual
Adicionar context-v2, validar política antes de efeitos e preservar default L1 vazio/lexical-v1. Nova política usa cinco conjunções: primeiros três termos; dois primeiros anchors; primeiro+último anchor; segundo+terceiro; anchor mediano+último. Quando poucos termos, usar disponíveis sem duplicação dentro da consulta. Remover familiaridade só dos anchors, preservar direta. Default novo cinco resultados/rota. Rotas continuam com perspectivas históricas mas derivação CONTEXTUAL_TERM_CONJUNCTION_V2, sem sufixos interpretativos ou nomes embutidos. IDs incluem política.

## 0155 — Ideias de algoritmos
Propagar AlgorithmIdeas na normalização do pedido de exploração, sem reinterpretar política de consultas ou quotas.
