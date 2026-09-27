# Tipos de rastreabilidade consultiva — 0098

Grafo decision-trace-v1, escopo RECORDED_DEPENDENCIES_ONLY, tipos INPUT/SOURCE/
CLAIM/DECISION. Arestas dirigidas da dependência para seu consumidor, inclusive
refutação. Fingerprints identificam conteúdo declarado, não autenticam a fonte.
Não promover evidência ou afirmar dependências completas. Report existente é
imutável; resultado incremental tem DecisionBefore==DecisionAfter e impactoNONE.

BuildInput contém report existente e opcionalmente ideia/versões de regras.
ImpactInput contém grafo congelado e mudanças explícitas por NodeID/reason.
Plano registra nós alterados, descendentes impactados, demais nós e tarefas
de revisão com ancestrais alterados. Expiração/revogação afetam aplicabilidade
atual, não apagam história nem tornam bytes anteriores falsos.
