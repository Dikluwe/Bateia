# Avaliação 0055

A avaliação usa o corpus e runner r2 congelados em ../0053. As famílias do catálogo foram definidas conhecendo o diagnóstico: este conjunto é de calibração, não holdout. Não interpretar melhora nele como generalização. Uma consulta dirigida conhecida de disponibilidade foi salva separadamente e excluída do recall.

```sh
cp 00_nucleo/benchmark/0053/probe.go.txt /tmp/bateia0055-probe.go
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia-0055 ./04_wiring
GOCACHE=/tmp/bateia-0048-cache go build -o /tmp/bateia0055-probe /tmp/bateia0055-probe.go
PYTHONDONTWRITEBYTECODE=1 python3 00_nucleo/benchmark/0053/runner.py --corpus 00_nucleo/benchmark/0053/corpus.json --freeze 00_nucleo/benchmark/0053/freeze.json --binary /tmp/bateia-0055 --probe /tmp/bateia0055-probe --output /tmp/bateia0055-new-run --discovery
```

Requer rede. Saída precisa ser diretório novo. Erros são unknown; tempos são telemetria da ferramenta/provedor. O código dos projetos-alvo não é executado. Preservar pergunta em Unverified não comprova suas restrições.

Sintaxe de consultas booleanas/content fundamentada na [documentação oficial Sourcegraph](https://sourcegraph.com/docs/code-search/queries), consultada em 2026-09-08. O planejador usa catálogo finito, não modelo semântico geral.

## Revisões preservadas

`collected/` é a primeira coleta com serialização defeituosa; `collected-r2/` é a rodada final após correção do contrato sintático (0055-v2). Não foi escolhido o melhor de resultados equivalentes: versões diferentes estão identificadas. O catálogo e as perguntas permaneceram iguais.

Os testes inicialmente tinham import não permitido no linteiro, fixture `null` usada como vazio e cleanup HTTP bloqueante. R0/r1 e logs foram preservados, expectativas de erro/admissão mantidas. A primeira suíte foi interrompida no cleanup; a suíte final passou. O primeiro build tinha acesso duplicado a `.Results`, corrigido sem mudança de contrato.

O diagnóstico Python recebeu HTTP403 e não foi avaliado; o diagnóstico pelo cliente HTTP Go foi salvo em provider-shape-diagnostic.json e não entra na métrica de descoberta.
