# Validação 0041 — Autofagia e Katamari

A suíte completa e o linteiro passaram; compilação por `go build -o bateia ./04_wiring`. Logs em `tests.log` e `lint.log`; identidades em `validation.json`.

O experimento usa dois inventários gerados pelo comando real de autofagia e seis execuções do comando Katamari. `results.json` contém os resultados; `cli-cases.tar.gz` preserva entradas, saídas e comandos. `rules.json` contém o catálogo usado.

Para repetir, compile o binário e, na raiz do repositório, execute:

```sh
python3 00_nucleo/benchmark/0041/run_cases.py ./bateia /tmp/bateia-0041-new-cases
```

O script cria entradas próprias no diretório indicado. Use um diretório novo para preservar comparações anteriores. Avaliações externas PROVEN são fixtures fornecidas para testar o seletor; não são resultados de uma prova comportamental real. Nenhum candidato é instalado ou executado.

Controles: cobertura completa continua aguardando validação, conjunto mínimo dentro do universo fornecido, cobertura parcial explícita, candidato conflitante bloqueado, evidência parcial convertida em investigação e busca vazia sem sucesso inventado. Os testes Go acrescentam licença/operação/manutenção, fontes inválidas, avaliações faltantes, duplicação, ordem, revisão única por projeto, limite enumerável, no-change e hipótese de substituição.

O hash do binário identifica o executável testado; o pacote não é uma cópia completa do checkout. Tempos nos logs são telemetria dos testes, não desempenho dos candidatos. Não há alegação de avaliação independente.
