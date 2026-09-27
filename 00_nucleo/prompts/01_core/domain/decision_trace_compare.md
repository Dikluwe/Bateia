# Contrato 0100 — comparar snapshots declarados

Representar entrada schema_version 1 com before/after congelados e saída
decision-trace-comparison-v1. Registrar deltas ordenados de nós completos e
arestas, hashes canônicos e mudanças de metadados. Deltas de nós têm status
ADDED, REMOVED ou DECLARATION_CHANGED e as declarações antes/depois presentes.

A decisão candidata é somente a declaração do segundo snapshot. O plano de
impacto usa o primeiro grafo e mantém a decisão histórica; decision_impact NONE.
Ausência de nó não afirma revogação, expiração ou perda factual de evidência.
Escopo RECORDED_DEPENDENCIES_ONLY, com limitações explícitas. Sem I/O ou execução.
