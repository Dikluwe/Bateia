# Coleta coordenada com exploração

L1, passo 0153. CollectInsightEvidence recebe contexto, planos, porta ContentProvider e pedido de exploração sem bundles. Validar contexto, pergunta/familiaridade/limites e 1–8 planos antes de efeitos. Forçar ReviewAlternatives no modo novo. Coletar via CollectEvidence; falha de resolução vira bundle incompleto com falha explícita, conforme modo legado. Propagar cancelamento sem completar todos planos depois do cancelamento. Explorar bundles coletados pela função existente; erro de integridade/normalização impede retorno de relatório útil. Não classificar relevância ou novidade e não buscar implicitamente outras fontes. Sem I/O direto.

Identidades cuja resolução não fixou commit imutável permanecem nos bundles/falhas e em unexplored_repositories; não entram no seletor de trechos. Acrescentar UNRESOLVED_SOURCE_BUNDLES. Outros erros de integridade continuam fatais.

## 0158 — Integridade e serialização estável

Aplicar o [contrato0158](../../../passos/0158-validar-origem-e-estabilizar-bundles.md) e o [ADR0022](../../../ADR/ADR-0022-integridade-de-coleta-e-serializacao-estavel.md). Domínio estabiliza requirement_coverage nil/vazio como array e preserva bytes populados, sem mutação. Coleta genérica compartilha admissão, conserva pin do plano e confronta identidade de metadados e Ref de extração antes de admitir conteúdo. Coleta coordenada valida todos os schemas1 antes de efeitos e propaga ErrIntegrity sem relatório útil. Somente a API genérica conserva schema ausente legado. Não recalcular registros/hashes históricos; representações instáveis passam à forma corrigida apenas em novas serializações.

Critérios compartilhados vigentes: [EVD-01–05](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#aceitação).
