# Verificação 0054

Contrato: ../../passos/0054-preservar-intencao-na-descoberta.md.

O corpus, snapshots, runner r2 e expectativas vêm do passo 0053 e permanecem congelados. `benchmark0053-before.json` e o recibo final permitem conferir integridade. `baseline.json` registra a produção e contrato anteriores ao candidato. A/B é segregação de autoria sem atestação de isolamento físico.

Reprodução a partir da raiz:

```sh
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0054-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0054 ./04_wiring
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0054-probe /tmp/bateia0054-probe.go
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0054 --probe /tmp/bateia0054-probe --output /tmp/bateia0054-new-run --discovery
```

Requer rede pública; o diretório de saída precisa ser novo. O resultado de rede pode variar. Os tempos são telemetria do Bateia/provedor, não benchmarks dos projetos. Preservar literalmente uma pergunta não demonstra que o provedor compreendeu sua semântica.
