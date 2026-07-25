# Caso de Uso: Motor de Busca de Código (v2)

**Status**: REVISADO — substitui a ideia inicial
**Data**: 2026-07-25
**Documentos relacionados**: ADR-0001 (revisada)

---

## 1. O que é o projeto

Um motor de busca de código para **uso pessoal de pesquisa em arquitetura de software**. Ele procura, lê e extrai exemplos de estrutura de projetos a partir de repositórios públicos, e entrega o resultado como um **relatório em Markdown** via linha de comando.

O objetivo não é buscar código por sintaxe (para isso existe grep). É responder perguntas de engenharia: *"como projetos reais organizam autenticação em Go?"*, *"quais camadas e manifestos um projeto desse tipo costuma ter que eu não pensei?"*

## 2. O problema atual

Quando um desenvolvedor precisa montar uma estrutura de código, procurar referências manualmente em dezenas de repositórios consome muito tempo.

Automatizar essa busca usando apenas a API oficial do GitHub falha por dois motivos distintos — e é importante não confundi-los:

1. **Busca** (Search API): limite de 30 requisições por minuto autenticado, 10 por minuto sem autenticação. Resultados limitados aos primeiros 1000 registros. [^1^][^2^]
2. **Download/leitura** (Contents API, raw): cota separada de 5.000 requisições por hora autenticado — e desde maio de 2025 o GitHub aplica limites também a acessos **não autenticados ao raw.githubusercontent.com**. [^3^][^7^]

O projeto de referência (GitDeepSearch) falhou porque tratou busca e download como o mesmo recurso e não previu nenhum dos dois limites.

## 3. Como o motor resolve o problema

A ferramenta usa um funil de três estágios, onde cada estágio tem **um modo de falha próprio e um orçamento de requisições próprio**:

1. **Descoberta (topo do funil)** — Localizar repositórios candidatos sem tocar na Search API do GitHub. Fonte primária: Sourcegraph (ou equivalente indexador). A cota do GitHub não é consumida nesta etapa.
2. **Triagem (meio do funil)** — Baixar apenas manifestos e arquivos de estrutura (`go.mod`, árvore de diretórios) dos candidatos, para decidir quais valem extração. Etapa cacheável: manifesto identificado por `repo@commit` não muda nunca.
3. **Extração (fundo do funil)** — Baixar os arquivos finais dos poucos repositórios validados, via Contents API ou clone shallow, dentro da cota de 5.000/h.

Invariante do funil: **cada estágio reduz o volume em pelo menos uma ordem de grandeza** (ex.: 1.000 candidatos → 50 triados → 5 extraídos). Se um estágio não filtra, ele não deveria existir.

## 4. O que o usuário faz com a ferramenta

O usuário é o próprio autor do projeto, em pesquisa de arquitetura. Três cenários:

1. **Copiar soluções exatas** — Pesquisa um padrão estrutural específico; o sistema entrega o código de projetos que já o implementaram, com referência `repo@commit` para rastreabilidade.
2. **Buscar ideias alternativas** — Pesquisa um conceito geral ("como fazer login de usuários em Go"); o sistema lista os métodos distintos encontrados, agrupados por abordagem.
3. **Descobrir necessidades ocultas** — Ao analisar a arquitetura de projetos similares, o usuário visualiza configurações, pastas e dependências que não havia previsto para o próprio projeto. **Este é o diferencial da ferramenta** — nenhum buscador de código existente se propõe a isso.

## 5. Comportamento em caso de falha

Um motor de busca é definido tanto pelo fracasso quanto pelo sucesso. Comportamentos obrigatórios:

| Falha | Comportamento |
| :---- | :---- |
| Busca retorna zero resultados | Informar e sugerir alargamento da query (remover filtros de linguagem/path). Nunca falhar em silêncio. |
| Busca retorna volume excessivo | Aplicar o corte do funil e **informar quantos candidatos foram descartados e por qual critério**. |
| Repositório sumiu entre descoberta e extração (404/410) | Marcar como `stale` no relatório, seguir com os demais. TOCTOU é esperado, não exceção. |
| Fonte do topo do funil indisponível | Abortar com erro claro. **Nunca** fazer fallback silencioso para a Search API do GitHub — isso recriaria o bug do GitDeepSearch. |
| Cota de download excedida na extração | Pausar até o reset informado no header `x-ratelimit-reset`, retomar do ponto de parada (extração é idempotente via cache). [^3^] |

## 6. Limites explícitos (o que a ferramenta NÃO faz)

- **Não entende semântica de código.** Não é análise estática, não resolve símbolos, não verifica se o código compila.
- **Não garante frescor absoluto.** O resultado depende da data de indexação da fonte de descoberta. Todo resultado carrega `repo@commit` + data da consulta; "atual" significa "o commit encontrado na data da consulta".
- **Não busca repositórios privados.**
- **Não contorna limites de plataforma de forma adversarial** (sem rotação de tokens, sem scraping de HTML). Quando a cota acaba, a ferramenta espera — não trapaceia.

## 7. Escopo da v1

- **Linguagem**: Go-first. O filtro de triagem entende `go.mod` e árvore de diretórios. Novas linguagens entram como **adaptadores de manifesto** (`Cargo.toml`, `package.json`, `pyproject.toml`), sem alterar o núcleo.
- **Interface**: CLI. Entrada: query + flags. Saída: relatório `.md` em disco.
- **Volume-alvo de sucesso**: uma consulta típica processa ~1.000 candidatos no topo e entrega 5–20 repositórios extraídos, sem exceder 10% da cota horária de download do GitHub.

## 8. Critérios de sucesso

1. Executar os três cenários da seção 4 sem esgotar nenhuma cota em uso normal.
2. Reexecutar a mesma consulta com cache quente sem fazer nenhuma requisição de rede nos estágios 2 e 3.
3. Trocar a fonte de descoberta (Sourcegraph → outro indexador) alterando apenas a camada de infraestrutura, com zero mudanças no núcleo.

---

[^1^]: GitHub Community Discussion #179480 — Search API limitada a 30 req/min por PAT: https://github.com/orgs/community/discussions/179480
[^2^]: Latenode Community — busca autenticada tem cota separada (30/min), não autenticada 10/min: https://community.latenode.com/t/why-do-i-encounter-a-rate-limit-error-when-i-have-remaining-api-calls-on-github/25029
[^3^]: GitHub Docs — Rate limits for the REST API (primários, secundários, headers): https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
[^7^]: GitHub Blog Changelog (2025-05-08) — limites estendidos a requisições não autenticadas, incluindo raw.githubusercontent.com: https://github.blog/changelog/2025-05-08-updated-rate-limits-for-unauthenticated-requests/
