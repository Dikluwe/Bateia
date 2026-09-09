# Experimento0058

Diagnóstico e correção da representação da consulta local. Corpus0053 e runner permanecem congelados. Comparação anterior:0057/collected-r1. Resultados são calibração de recuperação local; não medem descoberta pública ou comportamento executado.

- `diagnostic-before.json`: pergunta original, termos do plano, símbolo explícito; limite10000.
- `diagnostic-combined-before.json`: controle dirigido acrescentando termos à pergunta original.
- `diagnostic-after.json`: mesmas três consultas após correção. Intersection passa de ausente a7º na pergunta original.
- `candidate-freeze.json`: hashes do contrato, candidato e testes antes da execução.
- `comparison.json`: resultados por caso, antes/depois.
- `collected/`: resultados e recibos com argumentos, arquivos e hashes.
- `verification.json`: verificações e limites de independência A/B.

Reprodução a partir da raiz, com snapshots locais0053 disponíveis:

```sh
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0058 ./04_wiring
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0058-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0058-probe /tmp/bateia0058-probe.go
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0058 --probe /tmp/bateia0058-probe --output /tmp/bateia0058-reproduced
```

O diretório de saída deve ser novo. Tempos nos recibos são telemetria do Bateia, não benchmarks dos projetos. Não executamos código dos repositórios analisados.
