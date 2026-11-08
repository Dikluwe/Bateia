# 0138 — Ampliar descoberta de insights para outros sítios

**Data:** 2026-10-28. **Contrato ex-ante.** Usuário autoriza ampliar a descoberta diversificada para os outros sítios do Bateia, incluindo Gitee e AtomGit. L1 possui políticas; L2 transporte; L3 conectores existentes; L4 composição.

## Superfície e compatibilidade

`idea insights discover --sources all|github,gitee,atomgit,... --question ... --output ... [--max-calls N] [--locale en-US]` ativa relatório próprio `insight-site-discovery-v1`. Sem --sources, API/comando 0137 preserva comportamento/formato. Demais flags/budgets do pedido de descoberta são reutilizados; --max-calls e --locale só se aplicam ao modo multifonte. Factory de fontes é injetada e materializada após validação de flags/output, sem L3 na casca.

Reutilizar composição federada existente: GitHub, Sourcegraph, Gitee, AtomGit, Codeberg, GitLab, JOSS, zbMATH, iTunes, Google Play, F-Droid e Bitbucket quando workspaces estiverem configurados. Substituir no factory os anúncios capability-only de GitHub/Sourcegraph por bridges L1 sobre portas já existentes. Não mudar registro de mercado, conectores remotos, APIs antigas ou políticas de reúso. Até 12 bindings únicos; capacidades válidas pelo normalizador existente. Somente GitHub pode receber o leitor imutável 0137 neste recorte.

## Plano e execução

Resolver aliases ou IDs completos contra os bindings disponíveis; all é exclusivo. Recusar fonte ausente, vazia ou repetida e locale inválido antes de efeitos. Canonicalizar seleção independentemente da ordem; capacidades ordenadas por ID. Reutilizar as cinco consultas 0137. Percorrer query primeiro, fontes depois, para dar oportunidade a cada fonte antes de avançar a perspectiva. Limite global de 1–60 operações lógicas de catálogo, padrão 20, sob os deadlines total/por operação existentes. Até 64 candidatos únicos admitidos; limite por operação reduzido ao espaço restante. Reservar no máximo 128 KiB de conteúdo README admitido no conjunto, mantendo 32 KiB por documento e até oito referências.

Registrar todas as células fonte×query, inclusive BUDGET_EXHAUSTED, NOT_STARTED, UNSUPPORTED_LOCALE e AUTH_BLOCKED. Sem retries ocultos no core. AUTH_REQUIRED bloqueia as consultas seguintes daquele conector durante o run, sem repetir requisições inúteis; não transforma a fonte em vazia. Falhas de uma fonte não apagam resultados das demais. Cancelamento impede novas chamadas e conserva recibos/parcial. ErrorCode só usa allowlist segura; não copiar mensagens/URLs de erro arbitrárias. Validar resposta, cardinalidade, identidade/connector ID e metadados limitados (até 4 KiB por candidato), sem truncar silenciosamente.

Para conectores de aplicação que exigem filtro de mercado, usar o país da locale quando declarado entre os mercados suportados (en-US → US), registrando mercado no recibo e origem; caso contrário UNSUPPORTED_MARKET sem chamada. Fontes sem filtro recebem mercado vazio. Locale não é tradução nem prova de cobertura geográfica.

## Identidade e seleção

Usar ProductCandidate existente e conservar sua identidade `connector_id + external_id`. Nunca unir por owner/repo, nome ou link unilateral entre sítios. Mesmo owner/repo em Gitee e GitHub permanece distinto. Deduplicar dentro do mesmo conector, unir origens/query IDs e links de repositório com projeção canônica determinística; origem é selada pelo core com conector/versão/query/locale/hash da resposta. Metadata não confirma equivalência entre produtos. Preservar evidência acadêmica e links como declarações de catálogo.

Selecionar referências pelo ganho de fontes ainda não representadas, queries ainda não representadas, editor/autor ainda não representado e menos matches familiares literais; desempate por identidade canônica. Sem popularidade. Registrar razões/ordem da seleção e lacunas de fontes/queries sem referência. Variedade de sítios não comprova diversidade conceitual.

## Conteúdo e contexto para IA

Todos os selecionados geram referências de catálogo com origem CATALOG_METADATA_UNVERIFIED: nome, descrição, URL e proveniência continuam dados não confiáveis. Estes dados podem alimentar tarefas de pesquisa, sem serem apresentados como código observado. Somente GitHub vinculado a leitor pode converter owner/repo/revisões observadas em referência 0137 e coletar README, com detecção de conflito de SHA; origem deve ser o próprio registro GitHub, sem seguir links de outros catálogos. Outros sítios recebem CATALOG_ONLY. Não enviar credenciais de uma fonte a outra, instalar ou executar candidatos, nem inventar commit para catálogo.

Reutilizar ExploreInsights para bundles reais e expor suas referências de leitura separadas das referências de catálogo nas tarefas do novo relatório. Instruções à IA exigem investigar links/bytes antes de afirmar mecanismo e separar observação, declaração e hipótese. Recombinação/analogia não comprovam equivalência/domínio. Novidade NOT_ASSESSED, efeito UNVERIFIED sempre. Resultados operacionais distinguem catálogo disponível, contexto com README, falhas parciais, vazio válido e interrupção. ID inclui plano/seleção observados e recibos; plano determinístico inclui fontes/capacidades e budget global.

## Credenciais e verificação

Reutilizar variáveis já configuradas (GITHUB_TOKEN, SOURCEGRAPH_TOKEN, BATEIA_GITEE_TOKEN, BATEIA_ATOMGIT_TOKEN). Bloquear também os tokens Gitee/AtomGit em artefatos serializados, inclusive representação JSON escapada; nenhum valor de credencial entra em capabilities ou relatórios. Ausência de token AtomGit é lacuna de autenticação no conector existente.

Testes puros de bridges, aliases/capacidades/budgets, fairness, auth, parcial/cancelamento, deduplicação com namespaces/ordem, links sem promoção e separação catálogo/bytes/tarefas. CLI: validação antes de factory/rede, flags novas, modo legado, publicação protegida e segredos. HTTP em L3 com Gitee/AtomGit reais contra transporte de fixture; wiring confere lista federada e leitor só no GitHub sem rede. Ensaio público pequeno Gitee/AtomGit e outra fonte se disponível, sem tokens/execução de candidatos; preservar falhas e não confundir integração de fixture com disponibilidade online. Replays históricos intactos; make verify completo. Novos registros em lab/bancada/insights/0138.
