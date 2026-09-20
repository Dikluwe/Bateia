# Experimento 0061 — Indícios da ação por candidato

Contrato: passos/0061-discriminar-acoes-estaticas.md. A avaliação distingue start/stop pelo vocabulário finito do contrato. É uma hipótese sobre nomes e chamadas AST; não comprova resolução de chamadas, fluxo de execução, objeto ou efeitos. Documentação permanece fora da inspeção.

`diagnostic-before.json` preserva a projeção anterior do caso de interrupção: fontes por sobreposição lexical, sem discriminação da ação. `pair-contract.json` congela duas perguntas e a relação esperada de Add/Remove antes do ensaio; fonte única conhecida, calibração dirigida. Os testes B usam fixtures próprias e não leem esse material.

Reprodução do corpus congelado, com snapshots locais0053 disponíveis:

```sh
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0061 ./04_wiring
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0061-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0061-probe /tmp/bateia0061-probe.go
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0061 --probe /tmp/bateia0061-probe --output /tmp/bateia0061-reproduced
```

O diretório de saída deve ser novo. Para reproduzir o par dirigido, copie pair.go.txt para /tmp/bateia0061-pair.go, compile-o no módulo e execute passando o caminho de fonte registrado em pair-contract.json. O auxiliar somente lê e analisa a fonte; não executa o projeto.

Corpus e métricas não mudam. Recuperação por documentação não equivale a corpo analisado. Tempos são telemetria do Bateia, não benchmarks dos projetos. Sem rede, sem execução dos projetos, sem alegação de holdout ou isolamento físico atestado.
