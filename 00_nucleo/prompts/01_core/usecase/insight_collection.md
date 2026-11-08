# Coleta coordenada com exploração

L1, passo 0153. CollectInsightEvidence recebe contexto, planos, porta ContentProvider e pedido de exploração sem bundles. Validar contexto, pergunta/familiaridade/limites e 1–8 planos antes de efeitos. Forçar ReviewAlternatives no modo novo. Coletar via CollectEvidence; falha de resolução vira bundle incompleto com falha explícita, conforme modo legado. Propagar cancelamento sem completar todos planos depois do cancelamento. Explorar bundles coletados pela função existente; erro de integridade/normalização impede retorno de relatório útil. Não classificar relevância ou novidade e não buscar implicitamente outras fontes. Sem I/O direto.

Identidades cuja resolução não fixou commit imutável permanecem nos bundles/falhas e em unexplored_repositories; não entram no seletor de trechos. Acrescentar UNRESOLVED_SOURCE_BUNDLES. Outros erros de integridade continuam fatais.
