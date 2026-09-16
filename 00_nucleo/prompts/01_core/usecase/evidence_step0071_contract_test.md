# Contrato independente — cobertura e avaliação, passo 0071

Derivar os oráculos somente do passo 0071 e do baseline público congelado; não ler nem
adaptar expectativas à implementação candidata.

Cobrir regra declarada customizada, capacidade fora do catálogo, fragmento lexical
inseguro, regra hipotética com e sem sinais, bundle completo e incompleto, catálogo
fornecido ao pipeline e classificação da investigação. `CATALOG_GAP` deve permanecer
`INCONCLUSIVE`, fora de gaps do candidato e distinto de `NOT_VERIFIED`; hipótese lexical
nunca produz `PROVEN`, `ABSENT` ou componente.

Comparar JSON de cobertura/assessment em vinte permutações determinísticas de
capacidades, regras, aliases e itens, verificando imutabilidade das entradas. Os testes
usam somente APIs públicas, nenhum I/O, rede, relógio ou aleatoriedade.
