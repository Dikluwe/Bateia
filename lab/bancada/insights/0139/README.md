# 0139 — Leitura automática de README entre forges

O comando multifonte agora lê o README de Gitee e AtomGit automaticamente. Sem flags novas:

```sh
go build -o /tmp/bateia ./04_wiring
/tmp/bateia idea insights discover \
  --question 'backpressure queue cancellation' \
  --sources github,gitee,atomgit \
  --max-calls 15 --attempt-timeout 15s --output /tmp/insights-readmes.json
```

AtomGit exige BATEIA_ATOMGIT_TOKEN no ambiente; Gitee usa opcionalmente BATEIA_GITEE_TOKEN. Não incluir tokens no comando ou compartilhar seus valores.

## Resultado consumido pela IA

- `readmes`: recibo por candidato/sítio, commit e caminho quando disponíveis; falhas explícitas.
- `site_contexts`: cada contexto Gitee/AtomGit contém candidate_id, identidade completa do sítio, bundle congelado e cartões. Mesmo owner/repo, SHA e texto em forges diferentes não colidem. Não separar bundle do wrapper ao interpretar origem.
- `bundles` e `exploration`: mantêm escopo GitHub para compatibilidade.
- `tasks.reading_card_refs`: inclui cartões GitHub e cartões dos novos contextos. `catalog_reference_ids` conserva metadados separados.

A leitura confirma que os bytes foram retornados pela API no SHA solicitado; não comprova o mecanismo descrito, autoria, relevância ou novidade. Catálogo continua CATALOG_METADATA_UNVERIFIED; novidade NOT_ASSESSED e efeito UNVERIFIED.

## Limites operacionais

Até quatro GETs por referência: metadata pública → commit da branch padrão → árvore raiz no SHA → conteúdo no mesmo SHA. Não segue download_url ou redirects. Max-calls limita catálogo, não esses GETs; deadlines total/por leitura e máximo de referências continuam valendo. Até32KiB por README/128KiB agregados e cartões limitados por max-passages (GitHub tem precedência). Se faltar espaço de cartão, registra READING_CARD_BUDGET e conserva os bytes. README regular raiz somente; sem subdiretórios/symlinks/submodules. Árvore com100 entradas, página incompleta/truncada é recusada por TREE_LIMIT. Identidades ASCII owner/repo são o escopo atual.

O host canônico AtomGit determina endpoint fixo: atomgit.com → api.atomgit.com; gitcode.com → api.gitcode.com. Não há fallback nem encaminhamento ao GitHub. Falhas externas mantêm PARTIAL e códigos seguros.

## Evidência e reprodução

Fixtures em `03_infra/forge/insight_reader_test.go` usam conectores reais e transporte HTTP controlado, sem serviço remoto. O caso positivo busca ambos catálogos e faz quatro GETs de leitura por sítio: nove chamadas por host contando as cinco consultas de catálogo. Dois contextos com bytes/revisão iguais geram cartões distintos. Há controles negativos de identidade, autenticação, branch, commit, truncamento/paginação, symlink, redirects, corpo/base64/tamanho/blob, conteúdo binário, cancelamento e limite exato32KiB/SHA256. L1 verifica budgets/refs/parcial; L4 confere composição.

```sh
GOCACHE=/tmp/bateia-review-go-cache go test ./01_core/domain ./01_core/usecase ./03_infra/forge ./04_wiring -run 'Test(SiteReadme|ForgeInsight|InsightSite)' -count=1
BATEIA_VERIFY_OUTPUT="$PWD/lab/bancada/insights/0139/final" make verify
```

`baseline.json` registra os bytes anteriores. `verification.json` e `lineage.json` consolidam preservação e linhagem; executar `python3 lab/bancada/insights/0139/verify.py` confere os recibos e os seis replays históricos. A execução inicial com sockets bloqueados está em `sandbox-limited`; `final` registra o gate completo aprovado com loopback. `focused-tests.log`, `lint-focused.log`, `final/verification.json` e os recibos históricos registram as verificações. Nenhuma execução autenticada online é alegada por essas fixtures; disponibilidade real de Gitee/AtomGit continua sem certificação neste recorte. Nenhum software candidato foi executado.

Protocolos consultados: [Gitee SDK](https://gitee.com/sdk/gitee5j/blob/main/docs/RepositoriesApi.md), [AtomGit conteúdo](https://docs.gitcode.com/en/docs/apis/get-api-v-5-repos-owner-repo-contents-path/), [árvore](https://docs.gitcode.com/en/docs/apis/get-api-v-5-repos-owner-repo-git-trees-sha/) e [commit](https://docs.gitcode.com/en/docs/apis/get-api-v-5-repos-owner-repo-commits-sha/).
