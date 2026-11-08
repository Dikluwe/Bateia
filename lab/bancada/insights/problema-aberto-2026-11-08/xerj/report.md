# Investigação independente XERJ

O problema exige alterar compromissos e preservação de observações, além de substituir planilhas. Três alternativas foram fundamentadas em fontes de software. Nenhuma biblioteca pode transformar uma cópia desatualizada de disponibilidade em promessa segura de uma peça física única. Esta última conclusão é síntese do agente aplicada ao problema, não benchmark dos projetos.

## Mudanças comuns propostas (síntese)

- Colar identificador estável/QR em cada peça; manter identidade separada de descrições e fotos. Duplicatas suspeitas são ligadas para revisão, não fundidas automaticamente.
- Registrar autor/dispositivo, oficina, referência causal e identificador único da operação. Correções apontam ao registro corrigido e conservam o original.
- Separar observado, reservado localmente, pedido remoto pendente, recebido e consumido. Exibir última sincronização e pendências. A peça separada vai para posição física marcada com a reserva.
- Confirmar disponibilidade com a oficina detentora antes de prometer deslocamento; durante isolamento pedidos remotos são pendentes. Se for indispensável consumir sem rede, limitar autoridade a um responsável/dispositivo por peça e transferi-la com confirmação, mantendo observações livres.

- Gravar localmente antes de mostrar sucesso; reenvio idempotente; preservar fila após reinício; distinguir erro de permissão de ausência de rede.
- Exportar JSON/CSV documentado mais operações originais, identidades, ligações de correção e manifesto de fotos com hashes; incluir registros ainda não sincronizados e conflitos. Importar em armazenamento vazio e verificar integridade.
- Miniaturas e fotos comprimidas com indicação explícita de ainda não enviadas; medir orçamento de armazenamento, não pressupor replicar todas as fotos em todos os celulares.
- Um serviço de sincronização hospedado pela cooperativa é uma hipótese operacional; custos, backup, atualização e segurança não foram medidos.

## A1 — Fichas locais replicadas com reconciliação explícita

Uma ficha mutável por peça no cliente PouchDB, sincronizada com servidor compatível; conflitos entram em fila humana, conservando observações e fotos em registros separados imutáveis. Não apresentar o vencedor padrão como disponibilidade confirmada.

O trabalho cotidiano edita fichas familiares; a equipe aceita resolver conflitos após reconexão, em vez de derivar tudo de eventos.

Condição: Boa candidata se disputas por ficha forem raras e houver responsável disponível para conciliá-las.

Restrições atendidas:
- Operações locais sem rede e sincronização posterior são documentadas.
- Anexos binários e JSON permitem construir interface Android e exportação.
- Pode começar com infraestrutura pequena, mas suficiência para 45 mil peças é hipótese.

Limites e falhas:
- Se a UI ler somente o vencedor padrão, reaparece apagamento aparente de observações.
- Compactação remove revisões não folha; histórico exigido precisa de arquivo de observações próprio.
- Não impede duas reservas offline da mesma peça.
- Duas pessoas generalistas precisam manter resolução, exportação e recuperação, além da UI.

Experimento que pode rejeitar: Dois celulares alteram localização, descrição e reserva da mesma peça offline; reconectar em ordens opostas. Rejeitar se alguma observação/foto sumir, conflito ficar invisível ou voluntários não conseguirem resolvê-lo em dois minutos.

Relação com baseline: materialmente modificada (B1, B4). A estratégia geral já existia; P-conflict e P-compact especificaram fila explícita e arquivo separado, evitando tratar revisão como histórico permanente.

Evidência efetivamente lida:
- P-local: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/learn.md` linhas 7–7, localizador s0, passage ordinal 0. Trecho literal: "save data locally". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/019.json`.
- P-conflict: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/guides/conflicts.md` linhas 80–80, localizador s2, passage ordinal 0. Trecho literal: "To detect if a document is in conflict, you use the `{conflicts: true}` option when you `get()` it.". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/022.json`.
- P-compact: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/guides/compact-and-destroy.md` linhas 19–19, localizador b0, passage ordinal 0. Trecho literal: "non-leaf revisions will no longer be available.". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/017.json`.
- P-sync: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/guides/replication.md` linhas 111–111, localizador s3, passage ordinal 0. Trecho literal: "retry until the connection is re-established". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/021.json`.
- P-photo: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/_includes/api/save_attachment.html` linhas 5–5, localizador b0, passage ordinal 0. Trecho literal: "Attaches a binary object to a document.". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/018.json`.
- P-fetch: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/_includes/api/batch_fetch.html` linhas 143–144, localizador b4966, passage ordinal 0. Trecho literal: "  include_docs: true,\n  attachments: true,". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/032.json`.
- P-quota: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/faq.md` linhas 50–50, localizador s3, passage ordinal 0. Trecho literal: "your users will lose unsynced data.". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/018.json`.

## A2 — Livro de operações e custódia física com projeções

Cada chegada, movimentação, separação, uso e correção vira documento novo e imutável, com ID por dispositivo/operação, peça, autor e referências. Ficha atual é projeção reconstruível. Movimento é saída e recebimento ligados; peça em trânsito não é disponível. Reserva confirmada só pela custódia vigente; outros dispositivos registram pedidos pendentes.

Correção acrescenta explicação e compensação; não reescreve o que outra pessoa observou. Aceita estados provisórios e reconciliação de contradições físicas, privilegiando rastreabilidade.

Condição: Preferível quando histórico e correções são centrais e a cooperativa aceita uma regra explícita de custódia e confirmação física.

Restrições atendidas:
- Criações locais replicáveis evitam disputar a mesma revisão.
- Histórico fica em documentos vivos, não depende de revisões antigas compactáveis.
- Exportação de eventos e fotos pode reconstruir projeções sem fornecedor.

Limites e falhas:
- Eventos contraditórios continuam possíveis; log não resolve o mundo físico.
- IDs por relógio somente, como no exemplo ilustrativo da fonte, não bastam: relógios podem coincidir ou atrasar. Usar identidade aleatória/contador persistente e causalidade.
- Múltiplos documentos não são uma transação; modelar comando completo em um evento e projeções idempotentes, tratar recebimento faltante.
- Custódia rigorosa pode criar gargalo offline; se muitos voluntários precisarem confirmar a mesma peça simultaneamente sem coordenação, essa restrição de trabalho falha.
- Reprocessamento, migrações de evento e conflito de identidade aumentam responsabilidade dos desenvolvedores.

Experimento que pode rejeitar: Repetir duas vezes o mesmo envio, interromper entre saída e recebimento, corrigir consumo e usar duas reservas offline. Reconstruir tudo apenas do export. Rejeitar se saldo/estado depender da ordem de entrega, surgir disponibilidade duplicada confirmada ou a rotina de custódia atrasar reparos de forma impraticável.

Relação com baseline: já presente, materialmente detalhada (B2, B3). P-delta fundamenta documento por operação; P-atomic impõe unidade de escrita e tratamento de etapas faltantes. Custódia é síntese do agente, não capacidade alegada do produto.

Evidência efetivamente lida:
- P-delta: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/guides/conflicts.md` linhas 142–142, localizador s3, passage ordinal 0. Trecho literal: "you never update or remove existing documents". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/017.json`.
- P-compact: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/guides/compact-and-destroy.md` linhas 19–19, localizador b0, passage ordinal 0. Trecho literal: "non-leaf revisions will no longer be available.". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/017.json`.
- P-sync: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/guides/replication.md` linhas 111–111, localizador s3, passage ordinal 0. Trecho literal: "retry until the connection is re-established". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/021.json`.
- P-photo: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/_includes/api/save_attachment.html` linhas 5–5, localizador b0, passage ordinal 0. Trecho literal: "Attaches a binary object to a document.". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/018.json`.
- P-atomic: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/_includes/api/batch_create.html` linhas 122–123, localizador b2716, passage ordinal 0. Trecho literal: "the smallest\natomic unit is the document.". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/033.json`.
- P-fetch: pouchdb@0ec5a68736e306814f491081eaf31692211af007, `docs/_includes/api/batch_fetch.html` linhas 143–144, localizador b4966, passage ordinal 0. Trecho literal: "  include_docs: true,\n  attachments: true,". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/032.json`.

## A3 — Ficha colaborativa mesclável para descrições; decisões físicas separadas

Manter por peça mapas/listas locais Automerge para descrições e observações, mesclados após rede; expor conflitos de campos com getConflicts e histórico. Reservas/movimentos continuam registros explícitos submetidos à custódia, nunca um campo disponível cujo vencedor arbitrário define a realidade. Fotos ficam referenciadas por manifesto e arquivo separado.

Investir em colaboração frequente sobre descrições sem exigir resolução de documento inteiro; separar o domínio mesclável das decisões exclusivas de estoque.

Condição: Só priorizar se a experiência mostrar edição descritiva concorrente frequente e a implementação básica não atender à colaboração.

Restrições atendidas:
- Fonte documenta edição local e merge posterior.
- Valores concorrentes e histórico podem ser inspecionados.
- Formato binário incremental dá base para persistência/export técnico; exportação legível e mídia ainda precisam ser implementadas.

Limites e falhas:
- Convergência de estruturas não valida reservas nem localização física.
- Vencedor determinístico também pode esconder campo concorrente se getConflicts não for usado.
- Biblioteca não entrega sozinha app Android, transporte, persistência durável, catálogo, anexos ou operação de servidor.
- Fonte exige stream confiável em ordem; reconexões e sincronização de muitos documentos aumentam trabalho para só dois mantenedores.
- Memória/CPU/storage para 45 mil documentos e fotos em celulares baratos não foram medidos.

Experimento que pode rejeitar: Protótipo futuro com 500 fichas e descrições concorrentes em três aparelhos baratos, reinício durante escrita e 72h offline. Rejeitar se conflitar silenciosamente, perder histórico, exceder orçamento do aparelho, ou exigir mais tempo de manutenção que A1/A2 sem melhorar tarefas.

Relação com baseline: materialmente modificada (B1, B3). O baseline não nomeava CRDT. Automerge veio da memória model; A-local/A-conflict mostram uma variante de colaboração que preserva conflitos mas exige limite de domínio. Não é tratada como novidade descoberta pelo catálogo.

Evidência efetivamente lida:
- A-local: automerge@09aeeabaa499c3aca50c48f33051a5accad0cc58, `rust/automerge/src/lib.rs` linhas 7–7, localizador code, passage ordinal 0. Trecho literal: "entirely locally but which can at any time be merged with other instances of". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/029.json`.
- A-conflict: automerge@09aeeabaa499c3aca50c48f33051a5accad0cc58, `javascript/src/implementation.ts` linhas 1006–1007, localizador code, passage ordinal 987. Trecho literal: "Automerge chooses one value\n * arbitrarily (but deterministically". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/025.json`.
- A-history: automerge@09aeeabaa499c3aca50c48f33051a5accad0cc58, `javascript/src/implementation.ts` linhas 1158–1158, localizador code, passage ordinal 1156. Trecho literal: "const history = getAllChanges(doc)". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/030.json`.
- A-transport: automerge@09aeeabaa499c3aca50c48f33051a5accad0cc58, `rust/automerge/src/sync.rs` linhas 4–4, localizador code, passage ordinal 0. Trecho literal: "it assumes a reliable in-order stream". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/029.json`.
- A-save: automerge@09aeeabaa499c3aca50c48f33051a5accad0cc58, `javascript/src/implementation.ts` linhas 829–829, localizador code, passage ordinal 781. Trecho literal: "Create binary save data to be appended to a save file". Testemunho: `/tmp/bateia-xerj-open-0152/xerj/audit/031.json`.

## Avaliação pequena antes de compromisso

Propor um ensaio de uma semana em duas oficinas, seis voluntários e três aparelhos baratos existentes, com 200 peças representativas e 500 fichas sintéticas adicionais. Não executar os projetos nesta investigação. Prototipar no máximo duas opções inicialmente (A1 e A2); A3 só entra se a necessidade descritiva justificar. O mesmo roteiro cobre chegada, QR duplicado, transferência com recebimento atrasado, duas reservas da mesma peça, consumo seguido de correção, foto pendente, reinício, relógio errado e reenvio duplicado. Manter uma lista de verdade física independente. Isolar aparelhos por 72 horas e reconectar em ordens diferentes.

Critérios de rejeição propostos: qualquer observação ou foto perdida; qualquer promessa confirmada em duplicidade; reconstrução diferente após importação limpa; ação cotidiana com mediana acima de 20 segundos após treinamento; resolução de disputa acima de dois minutos sem compreender autoria. Registrar também tempo de sincronização, bytes enviados, espaço, bateria, tarefas concluídas e horas de trabalho dos dois desenvolvedores. Os limites são hipóteses negociáveis com a cooperativa, não resultados medidos. Só depois testar carga de 45 mil peças e orçamento de fotos real.

A1 troca simplicidade inicial por conciliação de fichas; A2 troca familiaridade por histórico explícito e trabalho de custódia; A3 troca investimento técnico por colaboração descritiva mais granular. Não confundir custo de licença com custo recorrente de operação. Não há evidência suficiente para declarar vencedor.

## Proveniência, execução e limites

Baseline congelado antes da primeira ferramenta de pesquisa: `/tmp/bateia-xerj-open-0152/xerj/baseline.json`, SHA256 d990ae6c0753f56518e1ececc3d403fff70d4f34282d514464d8e0a2311c1527. B1–B4 eram hipóteses do agente sem pesquisa. Nenhuma alternativa é anunciada como novidade forçada. B4 não virou opção independente de plataforma pronta por falta de fonte adequada; B2 foi incorporada como regra de custódia sintetizada.

Descoberta nativa: audit/005 falhou por catálogo ausente; audit/006 sincronizou 514 corpora. audit/007–010 retornaram ambiguidade ou no_coverage e resultados irrelevantes. Consultas específicas posteriores usaram fontes adquiridas de URLs model: PouchDB (audit/012, pin 0ec5a68736e306814f491081eaf31692211af007) e Automerge (audit/016, pin 09aeeabaa499c3aca50c48f33051a5accad0cc58). O catálogo não descobriu esses candidatos. Não houve busca externa.

Indexações audit/015 e audit/020 concluíram pesquisáveis com exit 3 e arquivos binários/arquivos compactados ignorados (16 e 36 respectivamente). Busca ampla Automerge recuperou ruído de fixtures de benchmark; esses trechos não fundamentam alternativas. Evidência útil veio de buscas de símbolos e contratos. Documentação PouchDB contém exemplos e referências antigas; presença no pin atual não prova comportamento em qualquer Android atual. Nenhuma fonte foi executada. Não medi performance, preços, confiabilidade real, permissões, segurança ou aceitação dos voluntários.

Não houve intervenção humana, publicação, feedback ou leitura de outro braço. Todos os produtos/acessos de aquisição passaram pelo runner auditado com require_escalated autorizado pelo teste. Nó temporário exclusivo: porta 19691, data-dir `/tmp/bateia-xerj-open-0152/xerj/node-data`, startup audit/013. Encerramento será registrado no complemento abaixo.

Entregas: `/tmp/bateia-xerj-open-0152/xerj/report.md`, `/tmp/bateia-xerj-open-0152/xerj/alternatives.json`, `/tmp/bateia-xerj-open-0152/xerj/sources.json`. Evidências completas e stdout/stderr preservados em `/tmp/bateia-xerj-open-0152/xerj/audit`.

Encerramento: enviei SIGINT (Ctrl-C) à sessão exclusiva 94912. O runner terminou com exit 130/KeyboardInterrupt, por isso o registro inicial audit/013 pode permanecer pending; não reescrevi auditoria. A chamada audit/034 confirmou que nenhum nó XERJ responde em 127.0.0.1:19691. Esta falha era o resultado esperado da verificação de encerramento. Um comando local de consolidação de contagem falhou por registro sem campo kind; foi corrigido com leitura tolerante, sem afetar fontes ou entregas.

Contagem final dos registros do runner: {"other": 1, "setup": 8, "research": 24, "transport": 2}. Aquisições: duas; nenhuma outra fonte foi adquirida.
