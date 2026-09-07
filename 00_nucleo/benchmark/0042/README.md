# Evidência diagnóstica 0042

Este diretório registra contraexemplos, não testes de correções aprovadas. Código de produção permaneceu inalterado durante este passo; o checkout já continha o 0041 pendente.

- `probe.go.txt`: programa sintético que chama APIs reais, sem busca remota ou execução de candidatos.
- `probe-results.json`: saída observada; contém comportamentos defeituosos.
- `catalog-compatibility.json`: IDs dos catálogos públicos comparados e sua interseção.
- `inventory.json`: hashes de arquivos Go de produção; inventário não significa revisão de cada arquivo.
- `review-ledger.json`: intervalos de leitura registrados e hashes.
- `validation.json`: identidade e limites da execução.
- `lint.log`: verificação estrutural deste passo; não certifica semântica.

Reprodução a partir da raiz do repositório, com dependências Go já disponíveis:

```bash
cp 00_nucleo/benchmark/0042/probe.go.txt /tmp/bateia-0042-probe.go
GOCACHE=/tmp/bateia-0042-go-cache go run /tmp/bateia-0042-probe.go
```

O teste de caminhos cria somente marcadores benignos temporários e os remove ao concluir. Instantes e outros metadados de relógio podem variar. A comparação dos catálogos é estática e não revalida todos os caminhos customizados. A revisão é do mesmo agente das alterações anteriores; não representa validação independente.
