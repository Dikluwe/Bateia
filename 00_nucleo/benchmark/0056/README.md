# Avaliação 0056

Mesmo corpus e runner r2 de 0053, conhecido pelo autor do contrato; não holdout. Nenhum catálogo ou alvo alterado. Comparar com 0055/collected-r2. Uma rodada pública final, sem escolher a melhor de repetições.

```sh
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0056-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0056 ./04_wiring
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0056-probe /tmp/bateia0056-probe.go
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0056 --probe /tmp/bateia0056-probe --output /tmp/bateia0056-new-run --discovery
```

Requer rede; saída nova. Erro/timeout é unknown. count solicita resultados do provedor, não número de repositórios; cobertura do índice permanece desconhecida. Não há download extra nem execução de código dos projetos. Tempos são telemetria.

Referência de count: [documentação Sourcegraph](https://sourcegraph.com/docs/code-search/queries), consultada2026-09-08.

A coleta usou o binário anterior à correção exclusiva do identificador de retomada live0056-v1. Consulta/admissão/local não mudaram. O recibo da coleta identifica esse binário; verification.json identifica o final. production-measured.json e idea-recommend-r0.go.txt preservam a revisão medida. Sem nova rede após essa correção.
