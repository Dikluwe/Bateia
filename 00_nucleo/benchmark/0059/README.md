# Experimento0059

Correção da exclusão sistemática do canal documental quando o canal de código ocupa o limite. Diagnóstico dirigido com alvos conhecidos, não holdout.

`fsnotify-p1-full-before.json` mostra Add em31º:28hits de código e quatro documentais. `colorful-p1-full-before.json` mostra BlendLab em79º:76hits de código e11documentais. Em ambos, o alvo é o terceiro candidato documental. O limite normal é24. Os campos de identidade, trecho documental e termos justificadores estão nesses arquivos. `fsnotify-p2-before.json` conserva o caso Remove, ausente na pergunta original mesmo sem corte; não atribuímos essa falha à competição por vagas.

O contrato0059 reserva no máximo um quarto do orçamento para candidatos documentais exclusivos, sem desalojar identificadores exatos nem aumentar o limite. Essa fração é política de seleção proposta, não valor ótimo demonstrado. Nenhuma mudança de consulta, sinônimos ou score.

A comparação deverá distinguir recuperação por código e por documentação: encontrar uma declaração via Doc não verifica seu corpo ou comportamento. A métrica congelada body_available mede disponibilidade de bytes, não análise semântica.

Reprodução local, com snapshots0053 disponíveis:

```sh
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0059 ./04_wiring
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0059-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0059-probe /tmp/bateia0059-probe.go
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0059 --probe /tmp/bateia0059-probe --output /tmp/bateia0059-reproduced
```

A saída deve ser nova. Recibos do runner conservam argumentos, hashes e observações. Tempos são telemetria do Bateia, não benchmarks dos projetos. Não há busca pública ou execução do código dos projetos neste experimento.
