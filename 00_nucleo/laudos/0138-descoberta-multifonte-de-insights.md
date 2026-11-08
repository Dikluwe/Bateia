# 0138 — Descoberta multifonte de insights

**Data:** 2026-10-28. **Contrato:** [passo 0138](../passos/0138-ampliar-descoberta-de-insights-para-outros-sitios.md). **Decisão:** [ADR-0019](../ADR/ADR-0019-descoberta-de-insights-com-identidade-por-sitio.md). **Evidências:** [bancada](../../lab/bancada/insights/0138/README.md).

## Entrega

`idea insights discover --sources all` ou lista de aliases/IDs ativa descoberta entre os conectores federados existentes. Inclui GitHub, Sourcegraph, Gitee, AtomGit, Codeberg, GitLab, JOSS, zbMATH, iTunes, Google Play e F-Droid; Bitbucket integra a lista quando workspaces estão configurados. Sem a flag, o comando e formato GitHub 0137 permanecem compatíveis.

O novo relatório `insight-site-discovery-v1` conserva plano, tentativas fonte×query, identidades/proveniência, seleção, recebimentos, bundles e tarefas para IA. O orçamento global padrão é vinte operações, até sessenta; até 64 candidatos únicos e oito referências. As cinco queries percorrem todas as fontes antes de avançar a perspectiva. Deadlines, locale/mercado sem suporte, autenticação, falhas, cancelamento e operações omitidas ficam explícitos. AUTH_REQUIRED bloqueia somente as rotas seguintes daquela fonte, sem repetir HTTP.

## Identidade e evidência

Candidatos são unidos somente dentro de connector_id + external_id, com todas as origens/hash de resposta seladas pelo núcleo. Mesmo owner/repo em GitHub/Gitee/AtomGit não é fundido. Seleção prioriza fontes ainda não representadas, ganho de queries, editor/autor e menos familiaridade lexical, sem popularidade. Essa variedade estrutural não comprova relevância, outro domínio ou novidade.

As tarefas separam catalog_reference_ids de reading_card_refs. Os primeiros são CATALOG_METADATA_UNVERIFIED; podem orientar pesquisa, sem serem apresentados como código observado. Links de outros catálogos não autorizam leitura pelo GitHub. Somente o binding GitHub recebe leitor imutável e mantém o coletor 0137. Os demais selecionados são CATALOG_ONLY. Estado de execução não promove hipóteses, novidade NOT_ASSESSED ou efeito UNVERIFIED.

## Arquitetura e credenciais

L1 possui plano/admissão, bridges GitHub/Sourcegraph, identidade/seleção e execução. Reutiliza a canonicalização de capacidades e os conectores existentes, sem modificar registro de mercado. L2 interpreta flags e publica resultados protegidos, construindo a factory somente no modo multifonte. L4 compõe as fontes e injeta leitor só no GitHub. L3 mantém protocolos HTTP. Não há cliente concreto em L1/L2 ou novo runtime.

As credenciais permanecem nos clientes correspondentes. O bloqueio compartilhado de artefatos agora também contempla BATEIA_GITEE_TOKEN/BATEIA_ATOMGIT_TOKEN, incluindo forma JSON escapada. A fixture de protocolos confere token Gitee na query e token AtomGit no header, sem cruzamento, e nenhuma requisição AtomGit quando sua credencial está ausente. Nenhum valor de token foi incluído nos registros.

## Ensaio público e limites observados

O binário real executou `--sources all` sem credenciais, sem `.env` e sem Bitbucket configurado: onze fontes e 55 células planejadas. O ensaio limitou-se a doze operações lógicas, um resultado por consulta, quatro referências/trechos, 45 s totais e 2 s por operação. O resultado foi PARTIAL com três candidatos de iTunes/GitHub e um README GitHub admitido.

AtomGit registrou AUTH_REQUIRED e quatro AUTH_BLOCKED. Gitee, Codeberg e Sourcegraph atingiram o deadline curto. F-Droid, GitLab, JOSS e zbMATH registraram vazio observado na consulta; Google Play registrou INTEGRITY_ERROR sem admissão. As células posteriores ficaram marcadas pelo orçamento. Esses estados não certificam indisponibilidade geral ou ausência de software.

A loja trouxe inclusive um aplicativo de exercícios para dor nas costas, ilustrando ruído lexical. O pacote conserva isso como pista de catálogo, sem afirmar pertinência à pergunta de backpressure. O README GitHub foi lido em SHA fixado; os cartões e suas identidades são reproduzíveis offline pelos bundles publicados. Não foram instalados, clonados ou executados candidatos, nem medida utilidade das ideias da IA.

## Verificação

`make verify` passou: lint sem violações e onze controles, auditorias sobre 375 arquivos produtivos, suíte raiz, vet, build e diff-check. O [recibo final](../../lab/bancada/insights/0138/final/verification.json) registra os comandos/códigos. A suíte usou loopback para fixtures HTTP e assinatura Git desativada somente no processo do gate.

Testes puros cobriram aliases/all/capacidades, budgets, round-robin, auth/bloqueio, locale/mercado, limites globais, parcial/cancelamento, identidades homônimas, projeção canônica e separação catálogo/leitura. Testes CLI verificaram factory lazy, formato legado e segredos escapados. A integração L3 usa conectores Gitee/AtomGit reais com transporte de fixture; o teste L4 verifica a composição real sem rede, inclusive Bitbucket condicional. Integração controlada não é evidência de disponibilidade online.

O verificador do snapshot confere identidade do plano, células/budget, origens, IDs por sítio, referências resolvíveis e igualdade do contexto reconstruído. O binário final também passou nos cinco verificadores históricos de 0106/0134/0135/0136/0137, preservando suas saídas nos respectivos escopos.

Comparada à baseline, a mudança preservou todos os testes Go preexistentes e 46 artefatos históricos inspecionados. Dos 838 arquivos Go anteriores, quatro arquivos de transporte/composição foram alterados. Quatorze novos e quatro alterados possuem prompts espelhados e [linhagem](../../lab/bancada/insights/0138/lineage.json). Gates selados com módulos próprios não foram promovidos pela suíte raiz.

## Continuidade

A descoberta entre sítios está integrada, com limitações verificáveis. Permanecem fora deste corte leitura imutável nos outros hosts, resolução semântica de equivalência, tradução de consultas e medição de utilidade/novidade. Fontes sem credencial, sem resposta no deadline ou somente com metadados continuam disponíveis para investigação com seus estados explícitos.
