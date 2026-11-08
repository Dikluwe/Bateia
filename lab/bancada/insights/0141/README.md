# 0141 — Perspectivas e trilhas para ampliar possibilidades

O foco é oferecer à IA caminhos adicionais de investigação. A agenda do modo `--sources` agora alterna perspectivas/fontes, e READMEs já coletados geram trilhas para projetos citados. Sem novas flags ou dependências, sem chamadas de modelo e sem HTTP adicional para extrair trilhas.

## Uso e saída

```sh
/tmp/bateia idea insights discover \
  --question 'serendipitous software discovery beyond semantic search and repository ranking' \
  --sources all --max-calls 20 --output /tmp/possibilidades.json
```

- `plan.scheduling_policy`: política explícita de rodízio perspectiva/fonte; integra ID do plano novo.
- `attempts`: permanecem em ordem query×fonte, com `execution_order` nas operações iniciadas.
- `opportunity_coverage`: perspectivas tentadas, perspectivas com resposta válida e fontes tentadas. Não conta ideias nem mede ganho conceitual.
- `reference_leads`: origem/candidato, commit, caminho/hash, linhas/trecho literal, destino e estado `REFERENCE_NOT_READ`.
- `tasks.reference_lead_refs`: liga tarefas da IA a essas trilhas, separado de cartões de leitura e referências de catálogo.

Sem --sources o contrato GitHub anterior é preservado. Resultados congelados anteriores não são reescritos. A nova agenda altera a ordem das operações de execuções multifonte novas e registra sua política no plano.

## Caso reproduzível, sem rede

O replay usa apenas os dois READMEs públicos congelados no ensaio comparacao-ia-2026-10-28, sem repetir busca online. O catálogo/leitor são fixtures em memória com versão marcada frozen-readme-replay-0141. Não é um ensaio de disponibilidade remota nem nova observação dos serviços.

```sh
GOCACHE=/tmp/bateia-review-go-cache go run ./lab/bancada/insights/0141/replay > /tmp/replay-possibilidades.json
```

O resultado em `frozen-corpus-result.json` mantém dois bundles lidos e extrai12 trilhas, incluindo `chnm/serendipomatic` das linhas106–107 do README de MuseCompNet/muse-tech-central. As tarefas passam a apontar esses caminhos. Nenhum destino foi lido automaticamente e12 trilhas não significam12 mecanismos/ideias validados. O índice continua testemunho indireto; software citado precisa de investigação própria.

O controle de agenda usa11 conectores em memória, todos elegíveis, sob20 operações: quatro tentativas em cada uma das cinco perspectivas e todas as11 fontes alcançadas. A agenda anterior, na mesma geometria20/11, cobria apenas duas perspectivas. Isso comprova distribuição de oportunidades; auth, locale, timeout e qualidade dos resultados podem reduzir cobertura real. As consultas continuam lexicais.

## Limites

Trilhas reconhecem URLs HTTPS absolutas de raízes GitHub/Gitee/AtomGit/gitcode/Codeberg/GitLab (incluindo subgrupos GitLab), com identidade ASCII e limites de caminho. Ignoram subcaminhos de código, credenciais/query/fragmento, escapes, hosts externos, rotas de navegação e destinos já selecionados. Formato aceito não comprova existência ou relação válida. Links relativos e texto de destino não são resolvidos.

Até128 URLs por documento,64 relações antes de seleção,12 trilhas expostas e1024 bytes por trecho. Cortes geram REFERENCE_LEADS_LIMITED. Seleção combina diversidade de destinos/origens e coincidências lexicais; referências sem coincidência continuam elegíveis. Origem→destino é preservada por relação, sem fundir sítios. Os bytes e limites de coleta anteriores permanecem32KiB/README e128KiB agregados.

## Verificação

Testes puros cobrem permutação1–12 fontes, cobertura/omissões, identidade/URLs adversariais, CRLF literal, hash/linhas, duplicação/self, destinos conhecidos, múltiplas origens, budgets e referências de tarefas sem leituras extras. `make verify` é o gate canônico, com loopback para fixtures gerais. O binário final reproduz os seis verificadores históricos; fontes/testes anteriores e arquivos de evidência são comparados com baseline.json. Logs, replay e recibos ficam nesta pasta. Nenhum candidato foi executado; novidade NOT_ASSESSED e efeito UNVERIFIED.

O recibo consolidado está em `verification.json` e a linhagem em `lineage.json`. Conferência offline: `python3 lab/bancada/insights/0141/verify.py`. O script verifica os seis recibos históricos e os testemunhos de links, sem consultar serviços externos.
