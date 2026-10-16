# 0107 — Pesquisa da qualidade do classificador usando Bateia

Esta bancada usa a ideia em `idea.json` como entrada do Bateia. A conclusão
é uma proposta de avaliação baseada em propriedades e afirmações, conservando
a agregação determinística atual. O relatório explicativo fica em
[pesquisa-classificador-0107](../../../../00_nucleo/relatorios/pesquisa-classificador-0107.md).

## Proveniência e escopo

- `automatic-plan.json`: 11 expansões produzidas por idea plan.
- `github-search-plan.json`: seis consultas manuais discriminadas; CheckList,
  Snorkel e FActScore são sementes oriundas da pesquisa anterior, sem alegação
  de descoberta cega. `github-discovery.json` registra seis candidatos reais.
- `academic-*.json` e recibos: quatro consultas, oito tentativas JOSS/zbMATH.
  Dois sucessos JOSS geraram seis candidatos, incluindo ruído de genética.
  Quatro respostas zbMATH não eram JSON e duas consultas JOSS receberam 429.
  Esses resultados não significam ausência de trabalhos.
- `catalog-screen.json`: triagem do Bateia. Score/ADVANCE não é prova de
  adequação ou qualidade científica; seleção manual em selection-review.json.
- `bundles.json`: coleta original Bateia incompleta em quatro repositórios,
  com 403 e resolução HEAD falha no CTF. Não foi substituída por coleta fictícia.
  `github-rate-receipt.json` registra quota core 60/60 usada, restante zero.
- `raw-fallback-receipts.json`: dez leituras complementares EXTERNAS pelo
  domínio público raw.githubusercontent.com, nos três SHAs fixados pelo
  Bateia. `review-bundles.json` acrescenta os bytes mantendo complete=false,
  failures e targets originais; não converte esse fallback em coleta Bateia.
- `source-observations.json`: inspeção manual com commit/path/linha/hash,
  sem execução de qualquer código dos candidatos. CTF tem artigo primário
  consultado separadamente; seu código não foi coletado em commit imutável.
- `locator-rules.json`: localizadores HIPOTÉTICOS manuais informados pela
  leitura dos próprios candidatos. Não são rótulos de referência independente.
  hash vazio é explícito: o CLI carrega catálogo padrão antes de ler --rules
  e, se hash for omitido, preserva o hash padrão incompatível (limitação observada).
- `assessment.json`: 18 avaliações de capacidades e três de alinhamento:
  quatro PARTIAL/HYPOTHESIS_ONLY e 17 INCONCLUSIVE; zero PROVEN.
- `evidence-metrics.json`: precisão null, sample_sufficient=false,
  inconclusive_rate=17/21, source_trace_rate=1 e budget_completion=0.
  Não há referência independente em independent-judgments.json. Rastreabilidade
  preenchida não equivale a acerto semântico; coleta parcial permanece parcial.
- `insight-input/result.json`: análise Bateia INCONCLUSIVE, zero mecanismos,
  24 lacunas RULE_NOT_IN_MECHANISM_CATALOG. Não forçamos metodologias de
  avaliação nas quatro classes fechadas de mecanismos existentes.
- `probe/`: oito intervenções controladas pelo CLI em fonte congelada 0106;
  caracterizam admissão, não medem qualidade semântica.

Nenhuma dependência dos candidatos foi instalada, modelo treinado, API de IA
chamada ou teste de candidato executado. Esta bancada usa somente ferramentas
existentes do Bateia; não altera seu código produtivo.

## Reprodução offline

Na raiz do Bateia:

```sh
go build -o /tmp/bateia-pesquisa ./04_wiring
python3 lab/bancada/pesquisa/0107/replay.py --binary /tmp/bateia-pesquisa
python3 lab/bancada/pesquisa/0107/probe_classifier.py \
  --binary /tmp/bateia-pesquisa --output-dir /tmp/bateia-probe-0107
```

Replay verifica hashes dos dez arquivos externos e compara cinco estágios
offline byte a byte. A sonda conserva entradas/saídas em diretório novo.
Os recibos de descoberta registram argv para repetição online em diretório
sem .env e com GITHUB_TOKEN/SOURCEGRAPH_TOKEN ausentes; use saídas novas.
Resultados de busca e quotas variam; não sobrescrever os registros originais.
