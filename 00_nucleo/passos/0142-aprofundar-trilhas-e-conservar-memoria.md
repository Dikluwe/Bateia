# 0142 — Aprofundar trilhas e conservar memória de exploração

**Data:** 2026-10-28. **Contrato ex-ante.** Usuário autoriza aprofundamento automático de trilhas e memória. A1-INV sob Tekt. L1 governa estados/quotas, L2 transporta snapshots, L3 resolve/leia APIs e L4 compõe leitores. Preservar modo sem --sources, formatos/hashes históricos e bytes antigos.

## Jornada pública

No modo multifonte, aprofundar até2 referências por padrão; --follow-links N aceita0–4 (0 desativa), sem mudar max-calls de catálogo. --memory caminho lê memória standalone ou exploration_memory de relatório multifonte anterior. Resultado contém exploration_memory novo, e --output protege o arquivo usado como entrada: persistência é o snapshot no resultado publicado atomicamente, sem segundo arquivo de estado divergente. Não descobrir arquivo global oculto ou enviar dados a serviço externo. A CLI deve validar JSON/ID/limites e alias de saída antes de factory/rede.

Comando offline idea insights memory --input arquivo --target URL --state pending|discarded --reason texto --output memoria.json altera somente alvo já conhecido. EXAMINED só decorre de coleta íntegra; feedback não fabrica leitura. Descarte/reabertura exigem motivo. Memória registra estado atual e última justificativa/observação, não auditoria completa. Hash verifica integridade do snapshot, não autentica autoria de histórico importado.

## Estado e memória em L1

Schema insight-exploration-memory-v1, ID canônico, geração inteira limitada e até256 identidades por connector_id+external_id. Snapshot até1MiB, sem conteúdo completo de README. Cada entrada: ID derivado da identidade, estado PENDING/EXAMINED/DISCARDED, motivo, tentativas, último status de leitura seguro, primeira geração, última pergunta e testemunho opcional de trilha/última leitura (commit/caminho/hash). Validar identidades/URLs, duplicatas, estados, motivo do descarte, SHA/path/origem do testemunho, contadores e ID antes de efeitos; cópia profunda e ordem canônica. Histórico de leitura é declarado pelo snapshot, não evidência de bytes relidos nesta execução.

Toda seleção nova e trilha encontrada alimenta a memória; limite cheio conserva entradas existentes e registra MEMORY_CAPACITY, sem descarte silencioso. Coleta íntegra muda para EXAMINED; falha mantém PENDING com contador/status, nunca DISCARDED. Estado EXAMINED ou DISCARDED recebido impede releitura automática, com recibo PREVIOUSLY_EXAMINED/MEMORY_DISCARDED e refs de memória nas tarefas. Catálogo e links não viram EXAMINED. Permitir reabrir explicitamente.

## Aprofundamento limitado

Um salto: usar trilhas atuais e PENDING preservadas cujo testemunho exista, priorizar menos tentativas e primeiras gerações, e ordem das trilhas atuais para desempate. Deduplicar por destino; nenhum descendente das leituras adicionais entra no mesmo run. Só ler sítio presente na seleção/factory, com leitor correto; rejeitar URLs não canônicas, destinos já lidos/descartados e vínculos inválidos. Registrar omissões como SOURCE_NOT_SELECTED, READER_UNAVAILABLE, CONTENT_BUDGET, DOCUMENT_BUDGET ou FOLLOWUP_BUDGET.

Porta FollowReader opcional só GitHub (InsightSiteRepositoryReader); demais forges usam RepositoryReader já composto. GitHub ganha ResolveRepository estrito: metadata com private:false, full_name/html_url exatos e default_branch válida, commit SHA observado; erros seguros, corpos limitados, sem redirects/fallback de host. Delegar árvore/README aos métodos existentes com SHA, sem clone ou execução. Outros leitores continuam no protocolo atual. Readmes primários preservam escopo anterior; followups têm identidade própria, lead/entry refs, recibo, bundle/cartões em wrapper, sem misturar CommitRef de hosts distintos. Leituras adicionais compartilham deadline total/por leitura,32KiB/documento,128KiB agregados, até8 documentos e max-passages global. Não acrescentar destino a candidatos de catálogo nem fingir resultado de busca.

reference_leads permanece testemunho de citação; estado efetivo da leitura do destino vem de followups e exploration_memory. reading_card_refs inclui cartões adicionais com IDs por sítio e referência, memory_entry_refs aponta histórico/decisões. IA distingue citação, leitura atual, histórico declarado e hipótese; novidade NOT_ASSESSED e efeito UNVERIFIED.

## Verificação

Testes puros de validação/selagem/imobilidade, feedback e estados, duas execuções (não reler examinados, avançar pendentes), retry/falha/descartes, limites de documentos/bytes/cartões e cancelamento parcial. Transporte GitHub em fixture com metadata pública/identidade/SHA, auth/404/redirect/adulteração; composição confere FollowReader. CLI memoria standalone/relatório, validação antes de factory, proteção contra alias/segredos e publicação única. Replay offline de coleção→destinos documentais em fixture, sem afirmar disponibilidade pública. make verify completo, seis replays históricos e preservação de artefatos. Novas evidências em lab/bancada/insights/0142.

Na seleção primária, identidades ainda pendentes/desconhecidas têm precedência sobre examinadas/descartadas; estas só preenchem vagas restantes e não geram releitura. A memória gerada reserva4KiB abaixo do limite1MiB; se uma atualização exceder a capacidade, preservar o valor anterior e registrar MEMORY_CAPACITY.

O FollowReader GitHub é um adaptador dedicado sobre Client: árvore raiz não recursiva limitada e leitura direta Contents sem o fallback raw do extrator legado. Até4 GETs (metadata/commit/árvore/arquivo), corpo256KiB/64KiB, arquivo32KiB com tamanho/base64/blob verificados. O leitor primário legado permanece intacto.

Se catálogo não completar nenhuma consulta mas o backlog da memória produzir leitura válida, conservar PARTIAL com conteúdo atual; FAILED/ErrInsightDiscoveryUnavailable só quando também não houver followup coletado. Falha de busca não apaga progresso da memória.

Escopo da primeira memória: repositórios nas cinco forges com leitor. Outros catálogos conservam suas referências, mas não entram neste snapshot de exploração de trilhas.
