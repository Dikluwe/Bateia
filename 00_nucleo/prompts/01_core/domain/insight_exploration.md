# Contratos de exploração para IA

**Camada:** L1. **Contrato:** passo 0136.

Declarar pedido de pergunta/familiaridade/bundles/limite, cartões de leitura com fonte/localizador/matches/razões, diagnósticos de cobertura da seleção, tarefas e pesquisas planejadas e relatório próprio. Distinguir trecho não interpretado, diversidade estrutural e novidade/efeito desconhecidos. Não modificar contratos anteriores, sem I/O/lógica.

## 0153 — Revisão opt-in
Adicionar ReviewAlternatives booleano omitido no JSON quando falso. O pedido padrão e seus hashes/saídas permanecem históricos. O campo é opção de contexto, não juízo sobre completude.

## Resultado de coleta para exploração

L1, passo 0153. InsightCollectionReport contém schema insight-collection-v1, bundles coletados e exploração desses mesmos bytes. Não certifica relevância, cobertura semântica, descoberta ou efeitos. Sem tipos de infraestrutura.

Identidades cuja resolução não fixou commit imutável permanecem nos bundles/falhas e em unexplored_repositories; não entram no seletor de trechos. Acrescentar UNRESOLVED_SOURCE_BUNDLES. Outros erros de integridade continuam fatais.

## 0155 — Ideias de algoritmos
AlgorithmIdeas booleano opcional omitido quando falso no pedido; AlgorithmAgenda []string opcional omitida no relatório legado. É guia de investigação, não lista de achados.
