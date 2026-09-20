# Execução 0043 — Evidências e reprodução

Regime A/B sem atestação física de isolamento. O agente de testes recebeu especificação e interfaces públicas, sem implementação candidata. Os registros declaram revisões de fixtures e mudanças de autoridade de leitura para migração de testes públicos antigos. Isso não é certificação independente do produto.

- `baseline.json`: hash do plano congelado, commit-base e hashes da worktree já alterada pelo 0041; não equivale a checkout limpo.
- `independent-tests.json`, arquivos `*-freeze.json` e `*-migration.json`: autoria, entradas, revisões e hashes dos testes.
- `synthetic-g-protocol.json`, `synthetic-g-freeze.json`, `synthetic-g-results.json`: casos e critérios definidos antes da execução. G possui somente negativos/opacos; não mede utilidade humana.
- `independent-test-results.json`: veredito do testador e limites.
- `tests.log`, `lint.log`, `validation.json`: verificações finais e identidade da implementação.
- `local-replay.tar.gz`: dossiê local e replay com entradas/saídas preservadas; nenhum software candidato foi executado.

Para reproduzir verificações, da raiz:

```bash
GOCACHE=/tmp/bateia-0043-go-cache go test ./...
crystalline-lint .
GOCACHE=/tmp/bateia-0043-go-cache go build -o /tmp/bateia-0043 ./04_wiring
```

Testes HTTP de adaptadores usam servidores locais e precisam de permissão de sockets. O comando legado que compila somente `main.go` omite arquivos do mesmo pacote; a compilação acima usa o pacote completo. Tempos nos logs pertencem ao Bateia/testes, não são benchmarks dos projetos candidatos.

Os testes novos podem ser localizados por `Step0043`; os testes de complexidade lexical têm autoria compartilhada. Fixtures de mecanismos preservam o catálogo público original byte a byte, pois a verificação de integridade não aceita reserialização JSON como identidade equivalente.
