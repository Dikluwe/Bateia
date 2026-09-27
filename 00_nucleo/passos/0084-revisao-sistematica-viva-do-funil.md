# Passo 0084 — Revisão sistemática viva do funil

## Estado e propósito

- Estado: implementado e verificado em 2026-09-24.
- Depende dos passos 0020–0026, 0031, 0043, 0079 e 0083.
- Objetivo: projetar os artefatos existentes do Bateia como uma revisão auditável,
  versionada e incremental sem alterar os comandos simples nem promover modelos a
  autoridade de evidência.

## Entregas

1. protocolo versionado e congelado antes da seleção;
2. ledger que dê destino a todo registro e candidato;
3. fluxo quantitativo derivado exclusivamente do ledger;
4. razões terminais obrigatórias, estáveis e ligadas a regras;
5. certeza por claim separada de `EvidenceStatus` e da decisão;
6. índice vivo append-only com deltas entre versões;
7. aprendizado ativo somente para priorização, medido em corpus independente.

## Invariantes

- o protocolo inclui pergunta, busca, elegibilidade, screening, evidência e atualização;
- alteração do protocolo abre nova linhagem, não reescreve versão anterior;
- todo source record é único ou aponta para o registro canônico duplicado;
- todo candidato termina incluído, excluído ou inconclusivo com reason code estável;
- o fluxo é rejeitado se suas contagens não fecharem;
- certeza descreve robustez da conclusão e nunca substitui o status observado;
- active learning não elimina, inclui, cria evidência nem define stopping rule;
- rótulos de avaliação permanecem ocultos até a simulação revelar o item selecionado;
- todas as operações da primeira entrega são puras e offline.

## Superfície

`idea review protocol|synthesize|update|prioritize|evaluate` recebe e produz JSON estrito.
O relatório visual futuro será projeção do ledger, nunca estado independente.

## Gate

Testar protocolo e hash, conservação do ledger/fluxo, reason codes, separação entre status
e certeza, linhagem viva, prioridade sem exclusão e avaliação que não consulta gold antes
da seleção. Depois executar lint, suíte completa, vet, build e `git diff --check`.

## Resultado

- `idea review protocol` congela pergunta, busca, elegibilidade, screening, evidência,
  decisão, execução e política de atualização sob um hash verificável;
- `idea review synthesize` deriva ledger, fluxo, certeza por claim e registro de vieses;
- `idea review update` cria linhagens por protocolo e versões append-only com delta;
- `idea review prioritize` retorna somente uma fila de leitura, sempre com
  `can_exclude: false`;
- `idea review evaluate` revela o gold somente depois de cada seleção e não produz
  stopping rule;
- o corpus separado está em `lab/bancada/corpora/step0084-active-learning.json`;
- a busca simples e os comandos anteriores permanecem inalterados.
