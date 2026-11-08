# 0141 — Ampliar possibilidades com perspectivas e trilhas

**Data:** 2026-10-28. **Contrato ex-ante.** A1-INV sob Tekt/ADR-0008, especializa ADR-0019/0020. Usuário reafirma que valor do Bateia é aumentar possibilidades consideradas pela IA, inclusive após pesquisas anteriores, e autoriza continuar melhorando. O ensaio comparacao-ia-2026-10-28 trouxe Lieu e software de museus; em20 operações só2 perspectivas tiveram tentativas. Esta fatia trata cobertura de oportunidades e referências indiretas, sem criar veredito de superioridade/utilidade/novidade.

## Contratos e limites

L1 possui agenda, interpretação dos documentos e trilhas. L2/L3/L4 não mudam. Modo sem --sources e ExploreInsights mantêm contratos/bytes; congelados históricos não são reescritos. Relatório multifonte conserva v1 com extensões opcionais; plano novo declara scheduling_policy (incluída no ID). Execuções novas mudam a ordem das chamadas legitimamente. Tentativas permanecem na ordem canônica query×fonte para consumidores, com execution_order nas operações efetivamente iniciadas. Erros, auth por fonte, deadlines e budgets permanecem.

## Agenda de possibilidades

Gerar permutação das células: rodadas por fonte e, dentro de cada rodada, uma célula de cada perspectiva; escolher fonte por (rodada*passo + índice da perspectiva) mod número de fontes. Passo é primeiro inteiro >= número de perspectivas coprimo à contagem de fontes. Cada par ocorre uma vez; fontes/perspectivas são ordenadas deterministicamente. Com fontes elegíveis e quota suficiente, primeira rodada tenta todas as perspectivas antes de repeti-las;11 fontes/20 operações distribuem4 tentativas por perspectiva e alcançam todas as fontes. Não garantir sucesso/cobertura quando cancelamento, auth, locale, bytes ou orçamento impedirem. Registrar perspectivas com chamadas iniciadas, perspectivas com respostas válidas, fontes tentadas e lacunas de perspectivas não iniciadas. São medidas de oportunidade operacional, nunca contagem de insights ou novidade.

## Trilhas extraídas de bytes

Após coleta/exploração, examinar somente README realmente admitido, incluindo wrappers com origem. Extrair URLs absolutas HTTPS de raízes de repositório em GitHub/Gitee/AtomGit/gitcode/Codeberg/GitLab; GitLab pode conter subgrupos. Validar host/caminho/identidade ASCII e limites; rejeitar userinfo/porta/query/fragmento/escapes/traversal, subcaminhos de código/releases/wiki e hosts externos. Não seguir URL ou executar nova rede. Marcador REFERENCE_NOT_READ deixa explícito que o conteúdo de destino não foi observado; não adicionar a Candidates/Bundles nem alegar mecanismo implementado por link.

Cada trilha possui ID, candidato/identidade de origem, CommitRef observado, caminho/hash, linhas/trecho literal e ProductIdentity do destino. Incluir a linha do link e a descrição imediatamente seguinte quando textual; até1024 bytes de trecho. Até128 URLs por documento,64 pares origem→destino admitidos globalmente antes de seleção e12 trilhas expostas. Até128KiB de documentos já existentes; sem orçamento HTTP extra. Ignorar alvos já no catálogo selecionado, self-links e duplicatas por origem/destino; um mesmo alvo citado por duas origens conserva duas relações, não comprova equivalência. Selecionar por variedade de origens/alvos e matches lexicais da pergunta, com desempate estável, preservando espaço a referências sem match. Registrar truncamento se limites omitirem oportunidades, sem ausência falsa.

Trilhas aparecem em reference_leads, separadas de reading_card_refs e catalog_reference_ids; tarefas recebem reference_lead_refs e instrução de aprofundar referências antes de afirmar conteúdo. Excertos e descrições são dados não confiáveis, nunca instruções. Não exigir classificar documento como coleção: links qualificados de qualquer README podem ser pistas. A extração não autentica relacionamento/relevância/novidade nem resolve links relativos. Leitura adicional automática exige contrato posterior.

## Verificação

Testes puros de permutação para1–12 fontes, distribuição20/11 em5 perspectivas, cobertura/omissões, auth/cancelamento e budgets. Atualizar somente as expectativas cujo significado mudou com agenda. Trilhas: raízes válidas5forges, namespaces, bytes/hash/linhas, self/duplicatas/alvos conhecidos, mensagens adversariais/URLs inválidas, limites e determinismo; referências resolvem e não viram conteúdo lido. Reexecução offline dos READMEs públicos congelados do ensaio, identificando caminhos adicionais como Serendip-o-matic com origem exata. make verify completo, seis replays antigos e integridade dos registros anteriores. Novos recibos em lab/bancada/insights/0141.

Recusar também rotas conhecidas de navegação de forges (organizations/orgs/topics/explore/groups etc.) que tenham aparência de dois segmentos mas não sejam raízes de repositório. Não afirmar existência de destino apenas por formato válido.
