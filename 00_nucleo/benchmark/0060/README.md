# Experimento 0060

Tradução local explícita caminho/caminhos → path, registrada como hipótese. Diagnóstico dirigido: `original.json` não contém Remove, mesmo com limite 10000; `path-control.json` contém em 5º, combinando apenas path no corpo. Controle em inglês e controle de ação também preservados. Consulta pelo nome é controle de identidade, não descoberta.

A tradução também recupera Add, registros por caminho e Mkfifo. Ela não discrimina a ação de interromper notificações. A comparação congelada mede recuperação de alvos, não precisão semântica ou comportamento.

Reprodução a partir da raiz, com snapshots locais 0053 disponíveis:

```sh
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0060 ./04_wiring
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0060-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0060-probe /tmp/bateia0060-probe.go
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0060 --probe /tmp/bateia0060-probe --output /tmp/bateia0060-reproduced
```

A saída deve ser nova. Corpus e métricas permanecem congelados. Casos conhecidos são calibração, não holdout. Tempos são telemetria do Bateia, não benchmarks dos projetos. Não executamos busca pública ou código dos projetos.
