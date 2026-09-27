# Contratos de comparação de declarações congeladas

Teste a comparação pura com snapshots sintéticos congelados, sem alegar veracidade. Verifique ordenação canônica sem deltas, determinismo, vetores não nulos e ausência de mutação. Mudanças apenas de rótulo ou tipo com fingerprint igual devem aparecer e exigir revisão; mudança apenas de relação `REFUTES` deve registrar aresta removida/adicionada e rever o consumidor histórico.

Teste nós adicionados/removidos com payload opcional correto e ordenação estável. Novos IDs não entram como nós antigos no plano; ausência não é `REVOKED` ou `EXPIRED`. Decisão candidata diferente aparece separada, enquanto a decisão histórica antes/depois e `decision_impact=NONE` são preservados. Mudança somente de limitações deve marcar metadados e pedir revisão conservadora da decisão.

Recuse schemas inválidos, hashes ausentes/adulterados, UTF-8 inválido, strings longas e orçamento agregado de ambos os grafos, retornando sempre resultado zero. Recuse mais de 100 raízes históricas sem descarte e aceite o limite exato de 100. Exercite entrada que caiba no orçamento mas cuja projeção de deltas mais tarefas expanda além do orçamento, sem resultado parcial nem serialização integral desprotegida.
