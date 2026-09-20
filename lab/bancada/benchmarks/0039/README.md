# Diagnóstico 0039 — jornada principal de reúso

Este pacote preserva uma execução diagnóstica, não um benchmark independente nem uma aprovação de qualidade. Veja o [relatório](../../relatorios/analise-jornada-principal-0039.md) e o [passo](../../passos/0039-analisar-jornada-principal-de-reuso.md).

- `inputs-and-dossiers.tar.gz`: entradas, candidatos sintéticos, cinco dossiês locais, replay, logs e adaptador controlado de serviço/MCP.
- `results.json`: achados derivados dos dossiês.
- `analyze_results.py`: verifica que os defeitos relatados e as proteções foram reproduzidos. Passar significa reproduzir o diagnóstico, não corrigir o produto.
- `probe.go.txt`: fonte do adaptador controlado, mantida como texto fora das camadas Go.
- `validation.json`: identidade do binário, commit-base, hashes de fontes, versões de regras e hash do pacote.
- `journey-tests.log`: testes existentes das interfaces e jornada, aprovados com cache.

## Conferir o diagnóstico congelado

A partir da raiz do Bateia, use um diretório vazio:

```sh
mkdir -p /tmp/bateia-0039-review
tar -xzf 00_nucleo/benchmark/0039/inputs-and-dossiers.tar.gz -C /tmp/bateia-0039-review
python3 00_nucleo/benchmark/0039/analyze_results.py /tmp/bateia-0039-review
```

O arquivo compactado contém os conteúdos examinados nos bundles. O checkout sujo completo e o executável não estão arquivados: os hashes permitem identificar divergências, mas o commit-base sozinho não reproduz o analisador.

## Reexecutar com o analisador disponível

Exemplo do caso de falso positivo, com saída nova para preservar a execução original:

```sh
./bateia idea recommend --idea /tmp/bateia-0039-review/events-go.json --policy /tmp/bateia-0039-review/policy.json --mode local --local /tmp/bateia-0039-review/events-go --strategy literal --output /tmp/bateia-0039-review/new-events
```

Os outros pares são `notes-python.json`/`notes-python` e `events-negated.json`/`events-negated`. A comparação usa a ideia de eventos e dois argumentos `--local`, um para cada candidato. A autoanálise original usa a ideia e política públicas de `00_nucleo/benchmark/0012/experimento`, com `--local` apontando à raiz do Bateia. O replay usa `--mode replay --replay <dossiê-original>` em lugar dos argumentos locais, mantendo ideia, política e estratégia.

O adaptador `probe.go` arquivado usa caminhos `/tmp/bateia-0039`; para reexecutá-lo, extraia em um diretório de trabalho separado e ajuste esse caminho antes de `go run`. Ele exercita o serviço e o MCP reais com provedores controlados; não mede disponibilidade nem qualidade da busca pública.

Tempos e durações são telemetria da ferramenta. A igualdade de replay se restringe aos quatro JSONs indicados em `results.json`.
