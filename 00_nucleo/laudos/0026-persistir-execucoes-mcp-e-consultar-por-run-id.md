# Laudo 0026 — Execuções MCP duráveis por `run_id`

**Data**: 2026-09-01  
**Estado**: `IMPLEMENTADO`  
**Veredito**: `MCP_DURABLE_RUNS_READY`  
**Regime Tekt**: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

## Resultado

O fragmento local `stdio` de execuções duráveis está aprovado para os hashes
registrados no certificado. A superfície MCP anuncia exatamente sete ferramentas,
preservando as três ferramentas do Passo 0025 e acrescentando `start`, `get_run`,
`get_artifact` e `cancel_run` nessa ordem.

O início é não bloqueante e idempotente, lifecycle e status de recomendação são
observáveis distintos, cancelamento tem um único vencedor, e recovery interrompe
órfãos sem inventar sucesso. Os seis artefatos seguros são publicados com
integridade, paginação determinística e cursor HMAC persistente. O store mantém
limites, permissões privadas e distinção entre ausência, indisponibilidade e
integridade.

## Reavaliação adversarial final

Três achados do verificador bloquearam versões intermediárias do candidato. As
mesmas testemunhas foram preservadas e repetidas após as correções:

- **M45**: o overlay original que observou aceitação de `cursor.key` hardlinked
  passou dez vezes. A regressão do candidato passou cinquenta vezes e cinco sob
  race; o inode externo não é lido nem alterado.
- **M55**: os arquivos L1 do Passo 0026 não importam `sync`, não contêm mutex,
  goroutine, timer ou I/O. A memorização concorrente foi delegada por porta
  estreita à fila L3; cancelamento passou 25 repetições e cinco sob race.
- **M44**: o binário real foi novamente iniciado com `BATEIA_RUN_ROOT` apontando
  para symlink. A configuração foi recusada antes de mutação: o destino externo
  permaneceu `0755` e vazio. A regressão passou cinquenta vezes e cinco sob race.

O gate discriminatório permanece em `60/60 = 1,0`, sem sobreviventes. P01–P05
permanecem `Preserved`; O01–O05 permanecem `Unknown` e não entram no denominador.

## Gates finais

- todos os 21 hashes do candidato e os 29 pares suite/prompt protegidos conferem;
- `go test -count=1 ./...` passou na execução independente, com o gate R6 em
  `57.991s`; outras duas execuções do coordenador também passaram;
- o gate R6 preservou nove processos de referência e rejeitou 270 processos
  mutantes em três ordens por três repetições;
- `go test -race` nos seis pacotes oficiais, `go vet ./...`, build e
  `git diff --check` passaram;
- `crystalline-lint .` terminou com exit 0. V12 está restrito a composição/testes
  L4 e V7 aos prompts dos quatro gates protegidos excluídos;
- `go list ./...` encontrou 29 pacotes normais e nenhum mutante sob `testdata`;
- o binário real listou as sete tools na ordem contratada e criou `cursor.key`
  com 32 bytes, modo `0600` e `nlink=1`;
- varredura dos canários adversariais nos arquivos de produção foi negativa;
- os hashes vivos de `recommendation.go` e `server.go` continuam iguais ao
  certificado do Passo 0025.

## Cadeia e limites

- manifesto final: `8d3c21973c0f3ed89d7eca5ecf7f5455454c4e54a0dc9de490575e80972439df`;
- selo R9: `f332530888df9f5c1b8cf02de49dbea3c269ac33d38b4bcae26e7326784924a9`;
- receipt discriminatório R9: `9b0e5fc0008041532d79de2ac68a4f772e22a0c3268f7ea7e3250bbf9bb5c3b3`;
- candidate receipt: `6802260a7b997b6d047083268d2f64bcbcdf0a0112380741c6f5b27e20fa8237`;
- verification receipt: `76c0c20aa66d00252a96d1cd78a78d0faa519b0113c63ff187d94dad3cff38d2`;
- certificado: `a04d6b07a939c1ce5fb7fb184e5cb56cc3cc34ed0b2b54c8087415a18a4a99ed`.

A execução ocorreu em worktree compartilhada. Refinamentos R7–R9 e correções
finais aconteceram depois de já existir candidato; a cadeia fornece hashes,
invalidações e testemunhas reproduzíveis, mas não alega independência temporal ou
isolamento material forte.

O laudo não generaliza testes de processo para durabilidade contra perda de
energia, todos os filesystems, ordem oculta de cancelamento, magnitude de entropia
ou shutdown forçado. Também não autoriza resume, replay, remoção, expiração,
multiusuário ou transporte remoto.
