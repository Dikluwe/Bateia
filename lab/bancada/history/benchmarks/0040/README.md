# Validação do passo 0040

[Passo](../../passos/0040-corrigir-autoridade-da-evidencia.md) · [Relatório](../../relatorios/correcoes-jornada-principal-0040.md)

## Artefatos

- `results.json`: resultados verificados dos contraexemplos e controles.
- `check_results.py`: verificações públicas sobre os dossiês congelados; não é avaliação independente.
- `dossiers.tar.gz`: cinco recomendações locais, replay, adaptador controlado serviço/MCP, entradas sintéticas, comandos e logs. Os caminhos absolutos dos comandos registram o ambiente original.
- `changed-production-sources.tar.gz`: os quatro arquivos de produção modificados neste passo, para inspeção da versão analisada.
- `validation.json`: commit-base, hashes do binário/fontes/pacote e resultados dos comandos de validação.
- `tests.log` e `lint.log`: suíte completa aprovada e linteiro sem violações.

## Conferir resultados

Na raiz do repositório, extraia em diretório vazio:

```sh
mkdir -p /tmp/bateia-0040-review
tar -xzf 00_nucleo/benchmark/0040/dossiers.tar.gz -C /tmp/bateia-0040-review
python3 00_nucleo/benchmark/0040/check_results.py /tmp/bateia-0040-review
```

Para uma nova execução, use as ideias e candidatos em `inputs/`, uma saída nova e o binário compilado com `go build -o bateia ./04_wiring`. Exemplo:

```sh
./bateia idea recommend --idea /tmp/bateia-0040-review/inputs/events-go.json --policy /tmp/bateia-0040-review/inputs/policy.json --mode local --local /tmp/bateia-0040-review/inputs/events-go --strategy literal --output /tmp/bateia-0040-new-events
```

Mudar o caminho local pode mudar a identidade do candidato. O adaptador `probe.go` arquivado requer ajustar seu diretório de saída antes de reexecutar. O pacote não contém o checkout completo nem o executável: os hashes identificam o estado, mas o commit-base sozinho não recompõe modificações anteriores.

O verificador do 0039 continua preservado: ele confirma o diagnóstico antigo. O verificador do 0040 exige os resultados corrigidos. Tempo de execução é telemetria da ferramenta, não benchmark do candidato.
