# ADR-0015 — Observações localizadas na investigação de insights

**Status:** ACEITO PARA O RECORTE 0134. **Data:** 2026-10-28. **Relacionados:** ADR-0008, ADR-0010, passos 0133/0134.

## Contexto

O diagnóstico 0133 mostrou que detalhes de uma implementação real permaneciam numa revisão externa ao relatório. A projeção do catálogo usava limites de linguagem como pré-condições, e sua estratégia genérica não descrevia necessariamente as condições operacionais da fonte. Reescrever a saída v1 invalidaria recibos/IDs históricos e atingiria consumidores de extração e síntese compartilhados.

## Decisão

Adicionar pedido v2 ao mesmo comando público de insights, com relatório e fachada L1 próprios. Cada fonte pode conter observações declaradas com classe, tipo, texto e referências de arquivo/hash/linhas. O núcleo confere o vínculo aos bytes fornecidos e recupera trechos exatos; a semântica da declaração permanece não reavaliada. A casca escolhe a versão do contrato e mantém entrada/saída protegidas.

O relatório v2 separa definição de catálogo (`CATALOG_TEMPLATE`), escopo de extração, estado/testemunhos da classificação e observações (`DECLARED_NOT_REASSESSED`). Pré-condições de hipóteses vêm de declarações PRECONDITION da mesma fonte/classe; trade-offs declarados também mantêm referência. Sem declaração, o estado é desconhecido/lacuna, sem inferir condição de uma lista de linguagens.

A classificação existente continua única e é usada por delegação sobre cópia v1 sem observações. Mecanismos conservam suas identidades; o novo relatório e suas hipóteses têm identidade v2 própria, incluindo as observações canonicalizadas. Observações sem mecanismo são conservadas em análise inconclusiva, sem gerar hipóteses. Declarações não promovem classificação, hipótese, origem externa ou efeito no alvo.

## Consequências

O recorte permite levar detalhes específicos ao relatório sem reconstruir a avaliação semântica ou exigir IA. A curadoria das declarações continua necessária; declarações contraditórias podem coexistir com origem explícita. O vínculo a um trecho não confirma a veracidade da interpretação. O alvo continua UNKNOWN, e adoção de código mantém política própria.

Os contratos v1 de insights, mecanismos, recomendação e síntese durável permanecem intactos. Evoluir a extração compartilhada ou analisar automaticamente o alvo são cortes posteriores, não efeitos implícitos deste ADR. L1 conserva regras, L2 transporte/apresentação, L3 efeitos e L4 composição; não há novo runtime.
