# ADR-0020 — Leitura de README com contexto por sítio

**Status:** ACEITO PARA O RECORTE 0139. **Data:** 2026-10-28. Relacionados: ADR-0018/0019, [passo0139](../passos/0139-ler-readmes-de-gitee-e-atomgit.md).

## Decisão

Ampliar leitura automática a Gitee e AtomGit por porta de resolução imutável própria, composta em L4 e implementada por HTTP em L3. A branch padrão declarada pela forge é resolvida para SHA; árvore e conteúdo são requisitados nesse SHA, com limites e verificação de blob. Metadata deve confirmar repositório público e mesma identidade antes da resolução. Sem seguimento de links/download_url/redirects, retries ou execução de candidatos.

A identidade legada CommitRef não contém forge. Em vez de modificar o contrato global de evidência e suas chaves, acrescentar site_contexts opcional ao relatório multifonte. Cada contexto possui ProductIdentity, candidate_id, bundle e cartões de leitura próprios. Os IDs de cartões incorporam identidade de sítio; tasks.reading_card_refs agrega esses cartões com os cartões GitHub. Bundles/exploration existentes continuam no escopo GitHub. Ausência de site_contexts conserva bytes dos resultados anteriores. Consumidores do novo conteúdo devem ler o wrapper, nunca interpretar seu bundle isolado como origem GitHub.

## Responsabilidades

L1 admite quais leitores pertencem a quais fontes, valida identidade/revisão resolvida, escolhe README pelo coletor existente, controla quotas e produz cartões com origem explícita. L3 autentica, resolve metadata/commit, requisita bytes e valida o protocolo e blob. L4 instancia/injeta leitores com as respectivas credenciais. Não há regra nova em L2, execução durável, política de reúso ou classificação de mecanismo.

## Limites

No máximo quatro GETs de leitura por referência, além das chamadas de catálogo contadas por max-calls; deadline por leitura e total existentes. Árvore raiz limitada e não recursiva: páginas incompletas/truncadas ou 100 entradas são recusadas conservadoramente. Conteúdo até32KiB, agregado até128KiB. Total de cartões respeita max-passages; cartões GitHub existentes têm precedência para preservar seu contrato. Um contexto de outro sítio fornece atéum cartão; se não houver espaço, bytes permanecem registrados e a lacuna READING_CARD_BUDGET é explícita.

Somente repositórios públicos com owner/repo ASCII são admitidos nesta leitura; nomes fora desse escopo continuam referência de catálogo com falha de leitura. AtomGit aceita hosts canônicos atomgit.com e gitcode.com nos respectivos endpoints fixos (api.atomgit.com e api.gitcode.com), sem fallback entre eles. Tokens são exclusivos dos adaptadores correspondentes. Auth/ausência/integridade/quota e timeout conservam resultado parcial. O SHA e hash de conteúdo fornecem localização e integridade; não autenticam autoria ou comprovam efeito/novidade.
