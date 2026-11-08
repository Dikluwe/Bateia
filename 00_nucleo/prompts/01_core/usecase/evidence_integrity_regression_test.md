# Regressões de integridade de coleta e estabilidade de bundles

**Camada:** L1. **Origem ex-ante:** [passo0158](../../../passos/0158-validar-origem-e-estabilizar-bundles.md) e [ADR0022](../../../ADR/ADR-0022-integridade-de-coleta-e-serializacao-estavel.md).

Domínio: serializar nil, vazio e populado por múltiplos roundtrips, incluindo entrada null legada. Exigir arrays não nulos, bytes estáveis, ausência de mutação e hash literal congelado para exemplo populado antes da mudança.

Casos de uso: portas contadas e controladas; recusar identidade/commit divergentes, nil e schema inválido, antes de produzir conteúdo. Plano inválido no final do lote deve impedir todas as portas; HEAD correto e metadado mais recente não alteram pin fixado. Fonte íntegra passa, falha operacional de resolução continua explícita, integridade é fatal na coordenação. Formulação com coverage nil/vazio/null mantém identidade em roundtrips, com controle literal de ID para exemplo populado capturado antes do patch. Chamador não é mutado. Sem I/O nos testes L1; nunca derive o resultado esperado de helpers produtivos de hash/selagem.

Critérios compartilhados vigentes: [EVD-01–05](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#aceitação).
