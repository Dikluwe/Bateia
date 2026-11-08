# 0138 — Descoberta de insights entre sítios

`--sources` amplia `idea insights discover` para os conectores federados já
disponíveis, incluindo Gitee e AtomGit. O modo sem a opção preserva o contrato
GitHub de 0137. Não exige repetir manualmente uma consulta por sítio.

## Uso

```sh
go build -o /tmp/bateia-insights-0138 ./04_wiring
/tmp/bateia-insights-0138 idea insights discover \
  --question 'backpressure queue cancellation' \
  --sources github,gitee,atomgit,codeberg,gitlab,joss,zbmath \
  --max-calls 20 --output /tmp/insights-sitios.json
```

`--sources all` seleciona também Sourcegraph, iTunes, Google Play e F-Droid;
Bitbucket aparece se BATEIA_BITBUCKET_WORKSPACES estiver configurado. São até
12 fontes. --max-calls é orçamento global, padrão20/máximo60; as cinco queries
alternam fontes em ordem canônica. --locale padrão en-US informa o locale,
sem tradução. Aplicações com filtro recebem o país suportado dessa locale;
sem suporte a célula registra UNSUPPORTED_MARKET.

## Dados para a IA

- `candidates`: ProductCandidate e origem CATALOG_METADATA_UNVERIFIED,
  identidade por conector/external ID, links e origens/hash de resposta.
- `tasks`: referências de catálogo em catalog_reference_ids e trechos de
  leitura em reading_card_refs, sem promover descrição a código observado.
- `readmes`: CATALOG_ONLY para sítios sem leitor imutável; o binding GitHub
  mantém coleta em SHA fixado, conforme 0137.
- `attempts`: células fonte×query, incluindo autenticação, quota, timeouts,
  orçamento, locale/mercado sem suporte e operações não iniciadas.

Mesmo owner/repo em Gitee, AtomGit e GitHub permanece distinto. Links externos
não autorizam leitura pelo cliente de outro host. O AtomGit sem
BATEIA_ATOMGIT_TOKEN recebe AUTH_REQUIRED e as próximas rotas AUTH_BLOCKED,
sem requisições repetidas. Gitee usa opcionalmente BATEIA_GITEE_TOKEN. Os
tokens são usados nos clientes correspondentes e bloqueados em artefatos.

## Ensaio público limitado em 2026-10-28

Executado sem tokens e sem `.env`, com all (11 fontes), máximo12 operações,
uma ocorrência por consulta, quatro referências/trechos, 45s totais e 2s por
operação. Das 55 células planejadas, 12 operações foram iniciadas e as restantes
receberam motivos explícitos de omissão. O resultado PARTIAL conserva:

- AtomGit: AUTH_REQUIRED e quatro AUTH_BLOCKED, sem autenticação disponível.
- Gitee, Codeberg e Sourcegraph: TIMED_OUT nesse ensaio de deadline curto.
- F-Droid, GitLab, JOSS e zbMATH: EMPTY no escopo consultado.
- Google Play: INTEGRITY_ERROR, sem candidatos admitidos.
- iTunes e GitHub: três candidatos no conjunto, de dois sítios; um README
  GitHub foi admitido e gerou um cartão de leitura.

O ensaio também mostra ruído: a loja retornou um aplicativo de exercícios
para dor nas costas e outro de fila. Isso não comprova relevância para
backpressure. Variedade de sítios/queries é estrutural; a IA deve interpretar
as pistas e buscar evidências antes de afirmar uma ideia útil.

A integração controlada usa os conectores reais Gitee/AtomGit com transporte
HTTP de fixture: autenticação, endpoints, limites e identidades homônimas foram
verificados. Isso não certifica disponibilidade online de nenhum serviço.
Nenhum candidato foi instalado, clonado ou executado.

## Reprodução e registros

```sh
python3 lab/bancada/insights/0138/verify.py --binary /tmp/bateia-insights-0138
```

O verificador confere o snapshot congelado, orçamento, origens/identidades,
separação catálogo/leitura e replay offline do contexto derivado. Não refaz
buscas online. Para novo ensaio público, copiar a bancada e executar
capture_online.py com um binário, preservando os registros originais.

- `online-result.json`, `online-check.json`, stdout/stderr: ensaio real.
- `public-snapshot-check.json`: verificação offline e replay do contexto.
- `focused-tests.log`, `lint-focused.log`, `final/`: testes e gate canônico.
- `baseline.json`, `lineage.json`, `verification.json`: preservação/consolidação.
- `legacy-v1-check.json`, `source-v2-check.json`, `comparison-check.json`,
  `exploration-check.json`, `github-discovery-check.json`: contratos históricos.

[Contrato](../../../../00_nucleo/passos/0138-ampliar-descoberta-de-insights-para-outros-sitios.md)
· [ADR](../../../../00_nucleo/ADR/ADR-0019-descoberta-de-insights-com-identidade-por-sitio.md)
· [Laudo](../../../../00_nucleo/laudos/0138-descoberta-multifonte-de-insights.md)
