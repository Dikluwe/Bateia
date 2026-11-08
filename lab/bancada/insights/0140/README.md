# 0140 — Leitura automática de README em Codeberg e GitLab

A descoberta multifonte agora coleta README das duas forges públicas e inclui seus cartões nas tarefas da IA. O modo sem --sources e as leituras GitHub/Gitee/AtomGit conservam seus contratos.

```sh
go build -o /tmp/bateia ./04_wiring
/tmp/bateia idea insights discover \
  --question 'backpressure queue cancellation' \
  --sources codeberg,gitlab \
  --max-calls 10 --attempt-timeout 15s --output /tmp/insights-readmes.json
```

`--sources all` também ativa os novos leitores. Sem novas flags ou credenciais: os adaptadores usam APIs anônimas de codeberg.org e gitlab.com. Instâncias privadas/self-hosted não integram este recorte.

## Origem, contexto e limites

`site_contexts` contém candidate_id, ProductIdentity, bundle e cartões de cada leitura; `tasks.reading_card_refs` inclui esses cartões. Mesmo owner/repo, commit e conteúdo em sítios distintos não colidem. `bundles`/`exploration` agregados continuam GitHub. No GitLab, `grupo/subgrupo/projeto` conserva Owner `grupo/subgrupo` e Repo `projeto`, namespace completo no cartão e identidade; não interpretar bundle separado do wrapper como origem GitHub. O seletor de passagens legado é reutilizado com validação e selagem da Source para subgrupos, sem flexibilizar seu normalizador público.

Até quatro GETs por referência: metadata pública, commit da branch padrão, árvore raiz no SHA, README no SHA. No Codeberg a revisão vem de branches/{branch}.commit.id; no GitLab repository/commits/{branch}.id. O caminho GitLab completo é escapado como ID do projeto. Não segue links/download_url/redirects nem usa tokens de outra forge. JSON até256KiB nas etapas iniciais/64KiB no arquivo, texto até32KiB por README/128KiB agregados, regular raiz somente. Árvore truncada/paginada ou100 entradas é recusada conservadoramente. Max-calls conta catálogo; leituras permanecem sob deadlines e máximo de referências. Max-passages limita cartões incluindo GitHub e os outros contextos; bytes sem espaço de cartão continuam com lacuna explícita.

Tamanho, base64, blob Git e, quando presente, content_sha256 são verificados; o GitLab exige também file_path e commit_id correspondentes ao pedido. Documentos convertidos pela API cujo blob não corresponda são rejeitados, sem fallback para baixar bytes por outra rota. Falhas externas preservam PARTIAL. Metadados continuam não verificados, novidade NOT_ASSESSED e efeito UNVERIFIED; ler uma descrição não comprova mecanismo ou utilidade.

## Verificação

Fixtures exercitam os conectores/leitores reais com transporte HTTP controlado: cinco buscas e quatro GETs de leitura por fonte, duas origens/contextos com cartões distintos, owner/repo iguais entre hosts e subgrupo GitLab. Negativos verificam identidade/publicidade/branch/SHA, paginação/symlink, corpo/base64/tamanho/blob, commit/caminho/hash divergentes, HTTP401/404/429, redirect e cancelamento. Limite exato32KiB e SHA256 também são verificados. Uma falha GitLab conserva o contexto Codeberg e suas referências de leitura.

```sh
GOCACHE=/tmp/bateia-review-go-cache go test ./01_core/usecase ./03_infra/forge ./04_wiring -run 'Test(Public|SiteReadme|ForgeInsight|InsightSite)' -count=1
BATEIA_VERIFY_OUTPUT="$PWD/lab/bancada/insights/0140/final" make verify
```

`baseline.json` congela os bytes anteriores. `focused-tests.log`, `lint-focused.log`, `final/verification.json`, replays históricos e recibo consolidado registram os controles. O verificador `python3 lab/bancada/insights/0140/verify.py` confere gate, preservação, linhagem e os seis recibos de replay. Apenas a expectativa do teste de composição foi atualizada de2 para4 leitores por sítio; demais testes anteriores são preservados. Estas fixtures não certificam disponibilidade online. Nenhum software candidato foi executado.

## Protocolos consultados

[GitLab projetos](https://docs.gitlab.com/api/projects/), [commits](https://docs.gitlab.com/api/commits/), [árvore](https://docs.gitlab.com/api/repositories/) e [arquivos](https://docs.gitlab.com/api/repository_files/). [Forgejo API](https://forgejo.org/docs/latest/user/api/usage/) e protocolos compatíveis [branch](https://docs.gitea.com/api/operations/repo-get-branch/), [árvore](https://docs.gitea.com/api/operations/get-tree/) e [conteúdo](https://docs.gitea.com/api/operations/repo-get-contents/). O OpenAPI servido pelo Codeberg não foi acessível pelo navegador de pesquisa; compatibilidade é testada por fixture e não apresentada como disponibilidade remota comprovada.
