# 0124 — Investigar Laya usando o produto atual

O produto executou coleta nova e insights offline. O agente executou controles
adicionais no Lab; estes não são descobertas automáticas do extrator.

- plans.json e results/bundles.json: coleta no SHA de 0122; 24 itens incluindo
  metadata, 540 alvos de expansão pulados por orçamento. Seis alvos pedidos presentes.
- insights-input.json: bundle real com assessments vazio explícito.
- results/insights.json: INCONCLUSIVE, zero mecanismos e hipóteses.
- controls.json: seis controles autorais correlacionados; fontes transformadas
  identificadas, pergunta congelada. results/context-controls.json: SDK real.
- results/mechanism-controls.json: split upstream extraído por AST e orçamento
  real do tokenizador/build_sequence. Não executou treino completo.
- results/verification.json: recibo de integridade e verificação do projeto.

Reprodução na raiz, com ambiente/pesos isolados de 0123 disponíveis:

```sh
go build -o /tmp/bateia-product-0124 ./04_wiring
```

Em diretório temporário sem .env, sem tokens:

```sh
env -u GITHUB_TOKEN -u GH_TOKEN /tmp/bateia-product-0124 idea evidence collect --plan /repos/Antigravity/Bateia/lab/bancada/decisoes/0124/plans.json --output /repos/Antigravity/Bateia/lab/bancada/decisoes/0124/results/bundles.json --timeout 90s
```

Na raiz:

```sh
python3 lab/bancada/decisoes/0124/prepare.py
/tmp/bateia-product-0124 idea insights analyze --input lab/bancada/decisoes/0124/insights-input.json --output lab/bancada/decisoes/0124/results/insights.json
HF_HOME=/tmp/bateia-laya-0123/hf-cache HF_HUB_OFFLINE=1 TRANSFORMERS_OFFLINE=1 OMP_NUM_THREADS=8 /tmp/bateia-laya-0123/venv/bin/python lab/bancada/decisoes/0124/run_controls.py
HF_HOME=/tmp/bateia-laya-0123/hf-cache HF_HUB_OFFLINE=1 TRANSFORMERS_OFFLINE=1 /tmp/bateia-laya-0123/venv/bin/python lab/bancada/decisoes/0124/audit.py
python3 lab/bancada/decisoes/0124/verify.py
```

Reexecutar substitui somente arquivos desta bancada; coleta remota precisa de
rede e pode mudar metadata. Scripts não treinam nem medem acurácia independente.
Laudo: [0124](../../../../00_nucleo/laudos/0124-pesquisa-do-laya-com-bateia.md).
