# Exploração de fontes para auxiliar IA

**Camada:** L1. **Contrato:** passo 0136.

Implementar ExploreInsights e normalização própria. Conferir quotas, identidades/hash/path e cópias profundas, sem exigir avaliações/observações/classes. Delegar extração/seleção de trechos ao componente insight_passages. Preparar tarefas/direções com referências resolvíveis, regras de leitura de conteúdo não confiável, alternativas/inversão/analogia/recombinação/contraexemplo; declarar limites de corpus/domínios e separar hipótese de observação/experimento. IDs determinísticos próprios, sem alterar contratos existentes. Não executar fontes, buscar, chamar IA ou promover semântica/novidade/efeito.

## 0153 — Revisão opt-in
Quando ReviewAlternatives for verdadeiro, preparar tarefa WITHIN_SOURCE_ALTERNATIVES com todos cartões e COVERAGE_REVIEW que pede opções sustentadas, descartes justificados e pesquisas pendentes. Não inventar alternativas ou cobertura completa. Instruções distinguem mesma fonte de fonte diferente, pedem revisar cada fonte/condições e repetir a exploração após coleta útil; direções são apenas planos. Padrão e suas saídas permanecem iguais.

## 0155 — Ideias de algoritmos
AlgorithmIdeas ativa ReviewAlternatives na normalização e acrescenta tarefas ALGORITHM_PRINCIPLE, APPLICABILITY_BOUNDARY e DECISION_DELTA com refs válidos e lacuna sem cartões. Agenda explicita observação->princípio->condições->alvo->decisão->experimento. Instruir uso de formulate com evidência de fonte e alvo, sem promover novidade/semântica/efeito. Fora do modo, saída histórica igual.
