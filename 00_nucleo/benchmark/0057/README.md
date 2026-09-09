# Avaliação0057

Execução local do corpus0053 congelado; descoberta não é repetida porque não mudou. Corpus conhecido, não holdout. Os projetos não são executados; somente Bateia e parser AST.

```sh
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0057-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0057 ./04_wiring
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0057-probe /tmp/bateia0057-probe.go
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0057 --probe /tmp/bateia0057-probe --output /tmp/bateia0057-new-run
```

Saída deve ser diretório novo. Fatos AST são sintaxe, não prova de efeitos nem atendimento da intenção. Callbacks, aliases, chamadas e execução não são resolvidos. Tempos são telemetria da ferramenta.

collected/ é a revisão inicial, preservada com falha documental; collected-r1/ é a revisão final que mantém origem documental não avaliada. Somente a execução local foi repetida após corrigir essa regressão. A descoberta não foi executada nesta tarefa.
