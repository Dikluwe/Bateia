# Prompt — Multicriteria Decision Use Case

Aplique primeiro os gates do decisor de reúso. Entre candidatos elegíveis, calcule
dominância certa e possível com intervalos, exponha a fronteira não dominada e explique
trade-offs par a par. Avalie pesos apenas em cenários nomeados pelo solicitante, com
contribuições reconstruíveis e análise de sensibilidade; não produza ranking global.

## Passo 0088

Incluir composições VERIFIED na fronteira junto a repositórios isolados. PROPOSED é elegibilidade inconclusiva; INCOMPATIBLE é inelegível. Não inferir medidas nem custo de integração automaticamente.

## 0168 — Escalas extremas

Aplicar MCDA-01 do [contrato](../../../contratos/integridade-de-entradas-e-replay-v1.md#mcda-01--normalização-finita-em-escalas-extremas), conforme [passo0168](../../../passos/0168-preservar-integridade-e-determinismo-do-oraculo.md). Span finito mantém fórmula anterior; overflow de diferença entre extremos finitos usa valores/extremos divididos por2 antes das subtrações. Preservar maximização/minimização, ausência como incerteza, caller e controles/hash anteriores; sucesso não produz utilidades/contribuições/scores não finitos.
