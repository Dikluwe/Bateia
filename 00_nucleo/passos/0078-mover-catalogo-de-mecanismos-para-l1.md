# Passo 0078 — Catálogo de mecanismos em L1

## Objetivo

Colocar junto da extração de mecanismos a autoridade sobre o catálogo usado por ela.
O catálogo ativo passa a pertencer a `01_core/usecase`; L4 apenas o injeta nos
consumidores.

## Contrato

`ActiveMechanismRuleCatalog` retorna os bytes canônicos e sua identidade estável. Cada
chamada devolve uma cópia independente para impedir mutação compartilhada. Não há I/O em
runtime: o arquivo é incorporado na compilação.

## Aceite

- conteúdo e interpretação das regras permanecem iguais;
- mutar o retorno de uma chamada não afeta a próxima;
- L1 não acessa filesystem, ambiente ou rede;
- não existe segunda cópia ativa do catálogo em L4.
