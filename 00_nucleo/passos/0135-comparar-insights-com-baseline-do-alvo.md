# 0135 — Comparar insights com baseline do alvo

**Data:** 2026-10-28. **Contrato ex-ante.** Continuação autorizada de 0134; lacuna de comparação registrada em 0133. Ownership: investigação L1; transporte/apresentação L2; sem novo efeito ou runtime.

## Superfície e compatibilidade

Adicionar `idea insights compare --input pedido.json [--output resultado.json]`. Pedido próprio de versão `1` contém `investigation` (pedido completo v2 de insights) e `target_baseline` opcional. Chamar `CompareInsightBaseline` em L1; `analyze` e seus formatos/IDs v1/v2 permanecem iguais. JSON estrito, alias, segredos, escrita protegida e short writes usam os mecanismos da CLI. Ajuda lista o comando novo.

## Baseline como snapshot de conteúdo

`target_baseline` contém `repository` (`owner/repo`, associação ao alvo declarada pelo autor), `snapshot_sha256`, `files` e `observations`. Cada arquivo contém `path`, `content`, `content_hash`. Os hashes conferem os bytes; o hash do snapshot cobre o JSON canônico da lista de `{path,content_hash}` ordenada por path, com hashes `sha256:`. Chaves nessa ordem, UTF-8, sem espaços/newline e sem escaping HTML. Não é commit Git, prova de completude do repositório, autenticação externa ou medição em execução. Escopo de saída: `SUPPLIED_FILES_ONLY`.

Até 64 arquivos, 256 KiB por arquivo, 1 MiB de conteúdo total; caminhos relativos únicos, hashes conferidos, UTF-8 válido. Observações usam o contrato 0134 (32 por baseline, referências de até 20 linhas/4096 bytes, tipos/classes/textos válidos e classe solicitada), sem mutar entradas. Validar usando o normalizador/projetor compartilhado de observações, adaptando o digest do snapshot como chave interna de identidade e usando namespace `target_obs_` na saída. Nenhum digest do snapshot é apresentado como commit. Até 2 MiB no pedido total. Baseline ausente ou sem observações significa desconhecido nesse escopo, nunca ausência da capacidade.

## Resultado consultivo

Relatório `insight-baseline-comparison-v1` conserva pergunta/alvo normalizados, catálogo, identidade/hash da investigação v2, análises/proveniência/lacunas da fonte e baseline localizada por classe. Baseline tem repository/snapshot/escopo/origem `DECLARED_NOT_REASSESSED`; não contém caminhos físicos ou bytes completos. Não há classificação de mecanismo nem veredito de presença/ausência do alvo a partir do texto das observações.

Gerar uma linha para cada fonte/classe solicitada e cada tipo de observação: STRATEGY, PRECONDITION, FAILURE_MODE, TRADEOFF, LIMITATION. Agrupar por classe/tipo, nunca por semelhança textual. Cada linha tem fonte/classe/tipo, referências separadas a observações da fonte/alvo, identidade própria, lacunas e próximo passo de validação. `BOTH_DECLARED`, `SOURCE_ONLY`, `TARGET_ONLY`, `NO_DECLARATIONS` descrevem disponibilidade de declarações; não equivalência, diferença ou evidência de efeito. Sem classe observada na fonte ainda emitir linhas com lacunas. Status geral PARTIAL quando há ao menos uma linha BOTH_DECLARED; caso contrário INCONCLUSIVE. Nenhum resultado COMPLETE/SUPPORTED.

O próximo passo deve comparar as condições/políticas já declaradas antes de propor alteração, conferir diferenças nos bytes e planejar baseline/experimento. Não emitir nova hipótese genérica de adicionar uma capacidade, nem reinterpretar as hipóteses v2: este relatório consultivo não as projeta. Efeito/aplicabilidade permanecem desconhecidos e adoção é NOT_AN_ADOPTION_DECISION. Não executar a fonte/alvo, inferir ganhos, fundir contraevidência ou resolver declarações contraditórias por palavras-chave.

IDs do relatório/linhas incluem contrato, input canônico da investigação, catálogo e baseline canonicalizada (incluindo declarações). Reordenação de arquivos/observações/referências e prefixo opcional do hash não muda saída; mudança de bytes/declaração/identidade do alvo muda IDs. Hash de snapshot não cobre declarações, mas identidade causal do relatório cobre. Não truncar declarações.

## Aceite

Testes puros de baseline/hash/linhas/quotas/imutabilidade/determinismo/identidade causal, indisponibilidade parcial ou total de declarações, isolamento de fontes/classes e ausência de promoção de efeito/estado. Testes CLI pela superfície pública, JSON estrito, falhas preservando saída, aliases, segredos e short writes. Reproduzir os contratos v1/v2 e corpus histórico sem alteração. Caso curado conc × snapshot atual da fila do Bateia deve expor submissão bloqueante versus reserva/rejeição e cancelamento cooperativo, como declarações lado a lado. Nenhuma diferença semântica é deduzida automaticamente pelo programa. Documentar, executar `make verify` e conservar novos registros em `lab/bancada/insights/0135`.
