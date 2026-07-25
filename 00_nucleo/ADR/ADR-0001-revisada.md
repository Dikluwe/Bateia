# ADR-0001 (rev. 1): Arquitetura Multimotor de Busca em Camadas, Decouplamento por Adaptadores e Cache como Cidadão de Primeira Classe

**Status**: PROPOSTO — aguardando validação do frescor das fontes de descoberta
**Data**: 2026-07-25 (revisão de 2026-07-24)
**Substitui**: ADR-0001 original (2026-07-24)
**Documentos relacionados**: Caso de Uso v2

## Contexto

A versão anterior da aplicação (GitDeepSearch) dependia exclusivamente da API REST oficial do GitHub e morreu por esgotamento de cota. A revisão desta ADR corrige três erros da versão original:

1. **Confusão conceitual entre busca e download.** São cotas separadas: Search API tem 30 req/min (autenticado) e apenas 10 req/min sem autenticação; a API de conteúdo (leitura de arquivos) tem 5.000 req/h. [^1^][^2^][^3^] O funil deve economizar *busca*, não *download* — o orçamento de download é 100x maior.
2. **raw.githubusercontent.com não é rota de fuga.** Desde 2025-05-08 o GitHub aplica rate limiting a requisições não autenticadas **explicitamente incluindo downloads do raw**. [^4^] A versão original desta ADR colocava o meio do funil inteiro sobre esse endpoint.
3. **Frescor das fontes de descoberta não foi tratado.** Fontes de topo de funil têm características de atualização muito diferentes, e isso determina qual pode ser primária:

| Fonte | Frescor | Custo | Limite | Adequação |
| :---- | :---- | :---- | :---- | :---- |
| Sourcegraph (sourcegraph.com) | Contínuo (indexa o GitHub corrente) | Gratuito com cadastro | Limites por conta, consulta complexa pode ser cortada | **Primária** |
| GH Archive (BigQuery) | ~horária | 1 TB/mês grátis de scan | Custo por byte escaneado | **Secundária** (metadados/eventos, não conteúdo) |
| bigquery-public-data.github_repos | Snapshot semanal, cobertura parcial (~54% dos repos com >15 estrelas, só ASCII <10 MB com licença detectada) | Idem | Idem | **Suplementar apenas** — nunca como fonte primária de descoberta [^5^] |

## Decisão

> 1. **Funil de busca com invariante de redução**: cada estágio reduz o volume em pelo menos uma ordem de grandeza.
>    * **Topo (Descoberta)**: Sourcegraph como fonte primária; GH Archive para enriquecimento por metadados (atividade, popularidade). `github_repos` no BigQuery rebaixado a papel suplementar (análise histórica/larga escala), nunca primário.
>    * **Meio (Triagem)**: leitura de manifestos (`go.mod`) e árvore de diretórios via **Contents API autenticada** (5.000/h), não via raw anônimo. GraphQL com batching de múltiplos arquivos por query quando o volume justificar (1 query = N arquivos).
>    * **Fundo (Extração)**: Contents API ou clone shallow (`git clone --depth 1 --filter=blob:none`) para os poucos repositórios validados. Para extração de múltiplos arquivos de um mesmo repo, clone shallow local é mais barato que N requisições de API. [^6^]
> 2. **Cache persistente como componente obrigatório (não otimização futura)**:
>    * Manifestos e arquivos extraídos são cacheados por chave `repo@commit` — conteúdo endereçado por commit é imutável, portanto o cache nunca expira.
>    * Resultados de descoberta são cacheados com TTL curto (horas), pois índices mudam.
>    * Extração é idempotente: interrupção por cota retoma do ponto de parada.
> 3. **Arquitetura de Adaptadores (Ports & Adapters)** — mantida da versão original, com uma salvaguarda:
>    * Contratos (interfaces Go) no 01_core, implementações no 03_infra, composição no 04_wiring. Parser no core (lógica pura, sem I/O).
>    * **Risco registrado — interface prematura**: a interface `SearchEngine` só será promovida a contrato estável do core quando o **segundo** adaptador de descoberta existir. Até lá, vive como detalhe do primeiro adaptador, para não vazar abstrações do Sourcegraph.
> 4. **Política de cota explícita**: a aplicação respeita os headers `x-ratelimit-remaining` e `x-ratelimit-reset`, aplica backoff exponencial em erros secundários, e **nunca** contorna limites de forma adversarial (sem rotação de tokens, sem scraping de HTML). [^3^]

## Prompts Afetados

| Prompt | Natureza da mudança |
| :---- | :---- |
| prompts/01-parser-e-structs.md | Mantido. Adicionar campo `commit` e `consultado_em` nas estruturas (rastreabilidade de frescor). |
| prompts/02-interfaces-busca.md | **Alterado**: interfaces `SearchEngine`/`FileFetcher` nascem no adaptador, promovidas ao core só com o 2º provedor. Adicionar contrato `CacheStore`. |
| prompts/03-adaptador-github-raw.md | **Substituído** por `03-adaptador-github-contents.md`: Contents API autenticada + GraphQL batching. Raw vira fallback interno deste adaptador, não um adaptador próprio. |
| prompts/04-adaptador-sourcegraph.md | Mantido (fonte primária de descoberta). |
| prompts/05-shell-cli.md | Mantido. Saída em relatório .md. |
| prompts/06-wiring-main.md | Mantido. Wiring agora inclui o `CacheStore`. |
| prompts/07-cache-store.md | **Novo**: cache em disco por `repo@commit`, TTL para descoberta, retomada idempotente. |

## Consequências

**Positivas**:

* Elimina o gargalo de cota na fase de descoberta sem criar um gargalo novo no meio do funil.
* Cache por `repo@commit` torna reexecuções gratuitas e a extração idempotente.
* Tabela de frescor torna explícita (e auditável) a escolha do provedor primário.
* Distinção busca/download impede que prompts derivados modelem orçamentos de cota errados.

**Negativas**:

* Conteúdo extraído exige autenticação (PAT) — o cenário "zero chaves" da ideia inicial não existe mais.
* Cache em disco adiciona um componente com estado (invalidação de descoberta, limpeza).
* GraphQL batching exige modelar queries dinâmicas (custo de complexidade no adaptador).
* Duas fontes distintas (indexador + GitHub) criam janela de TOCTOU: repositório pode sumir ou mudar entre descoberta e extração. Mitigação: marcar `stale`, seguir, e sempre registrar `repo@commit`.

**Neutras**:

* Chave do Sourcegraph pode ser necessária conforme volume; conta gratuita sujeita a cortes em consultas pesadas.

## Alternativas Consideradas

| Alternativa | Prós | Contras |
| :---- | :---- | :---- |
| Uso exclusivo da API oficial do GitHub | Simplicidade | 30 req/min na busca inviabiliza descoberta; foi a causa da morte do GitDeepSearch. |
| Raw anônimo como meio do funil (ADR original) | Sem chave, HTTP simples | Rate limited desde 2025-05 [^4^]; 1 GET por arquivo; sem batching. **Rejeitada.** |
| BigQuery `github_repos` como topo do funil | SQL arbitrário sobre conteúdo | Snapshot semanal, cobertura parcial, custo por byte; dados velhos geram "soluções exatas" desatualizadas. **Rebaixada a suplementar.** |
| Espelhamento local (clone em massa estilo pygit2) [^6^] | Throughput máximo, zero API depois do clone | Custo de disco/banda por repo inteiro; só compensa quando a extração quer muitos arquivos do mesmo repo — adotada **pontualmente** no fundo do funil. |
| **Adaptadores em funil + cache (escolhida)** | Orçamento de cota correto por estágio, frescor auditável, reexecução gratuita, core testável | Mais código inicial; estado do cache; TOCTOU entre fontes. |

## Questões em aberto (bloqueiam o status ACEITO)

1. Validar na prática o frescor e os limites da conta gratuita do Sourcegraph para as queries-alvo (executar 3 consultas reais e medir).
2. Definir o formato do relatório .md de saída (schema dos três cenários do Caso de Uso v2).
3. Medir o throughput real do GraphQL batching para leitura de manifestos (quantos arquivos por query antes do timeout de nó).

## Referências

[^1^]: GitHub Community Discussion #179480 — Search API: 30 req/min por PAT; GitHub App escala por instalação: https://github.com/orgs/community/discussions/179480
[^2^]: Latenode Community — busca autenticada 30/min, não autenticada 10/min, cota separada da core API: https://community.latenode.com/t/why-do-i-encounter-a-rate-limit-error-when-i-have-remaining-api-calls-on-github/25029
[^3^]: GitHub Docs — Rate limits for the REST API (primários: 5.000/h autenticado; secundários; headers `x-ratelimit-*`): https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
[^4^]: GitHub Blog Changelog (2025-05-08) — "Updated rate limits for unauthenticated requests", incluindo raw.githubusercontent.com: https://github.blog/changelog/2025-05-08-updated-rate-limits-for-unauthenticated-requests/
[^5^]: fhoffa/analyzing_github — características dos datasets: GH Archive (~horário), github_repos (snapshot semanal, cobertura parcial, ASCII <10 MB, licença detectada): https://github.com/fhoffa/analyzing_github
[^6^]: Stack Overflow — espelhamento local com pygit2 como alternativa de throughput para análise de muitos repositórios: https://stackoverflow.com/questions/13394077/is-there-a-way-to-increase-the-api-rate-limit-or-to-bypass-it-altogether-for-git
